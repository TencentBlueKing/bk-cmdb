/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package flow TODO
package flow

import (
	"context"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/task"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/mongo"
)

type flowOptions struct {
	key         event.Key
	task        *task.Task
	EventStruct interface{}
}

func newFlow(ctx context.Context, opts flowOptions, parseEvent parseEventFunc) error {
	flow, err := NewFlow(opts, parseEvent)
	if err != nil {
		return err
	}

	return flow.RunFlow(ctx)
}

// NewFlow create a new event watch flow
func NewFlow(opts flowOptions, parseEvent parseEventFunc) (Flow, error) {
	if parseEvent == nil {
		return Flow{}, fmt.Errorf("parseEventFunc is not set, key: %s", opts.key.Namespace())
	}

	return Flow{
		flowOptions: opts,
		metrics:     event.InitialMetrics(opts.key.Collection(), "watch"),
		parseEvent:  parseEvent,
		cursorQueue: &cursorQueue{
			cursorQueue: make(map[string]string),
		},
	}, nil
}

// Flow TODO
type Flow struct {
	flowOptions
	metrics      *event.EventMetrics
	tokenHandler *flowTokenHandler
	parseEvent   parseEventFunc
	cursorQueue  *cursorQueue
}

// cursorQueue saves the specific amount of previous cursors to check if event is duplicated with previous batch's event
type cursorQueue struct {
	cursorQueue map[string]string
	head        string
	tail        string
	length      int64
	lock        sync.Mutex
}

// checkIfConflict check if the cursor is conflict with previous cursors, maintain the length of the queue
func (c *cursorQueue) checkIfConflict(uuid, cursor string) bool {
	dbCursor := uuid + "-" + cursor
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.cursorQueue[dbCursor]; exists {
		return true
	}

	if c.length <= 0 {
		c.head = dbCursor
		c.tail = dbCursor
		c.cursorQueue[dbCursor] = ""
		c.length++
		return false
	}

	// append cursor to the tail of the cursor queue
	c.cursorQueue[c.tail] = dbCursor
	c.cursorQueue[dbCursor] = ""
	c.tail = dbCursor

	if c.length < cursorQueueSize {
		c.length++
		return false
	}

	// when cursor queue length reaches the limit, remove the earliest cursor
	newHead := c.cursorQueue[c.head]
	delete(c.cursorQueue, c.head)
	c.head = newHead

	return false
}

const (
	batchSize       = 200
	cursorQueueSize = 50000
)

// RunFlow run event flow
func (f *Flow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.metrics)

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name: f.key.Namespace(),
			CollOpts: &types.WatchCollOptions{
				CollectionOptions: types.CollectionOptions{
					CollectionFilter: &types.CollectionFilter{
						Regex: fmt.Sprintf("_%s$", f.key.Collection()),
					},
					EventStruct: f.EventStruct,
				},
			},
			TokenHandler: f.tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.TaskBatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: batchSize,
	}

	err := f.task.AddLoopBatchTask(opts)
	if err != nil {
		blog.Errorf("run %s flow, but add loop batch task failed, err: %v", f.key.Namespace(), err)
		return err
	}
	return nil
}

func (f *Flow) doBatch(dbInfo *types.DBInfo, es []*types.Event) (retry bool) {
	eventLen := len(es)
	if eventLen == 0 {
		return false
	}

	rid := es[0].ID()
	hasError := true

	// collect event related metrics
	start := time.Now()
	defer func() {
		if retry {
			f.metrics.CollectRetryError()
		}
		if hasError {
			return
		}
		f.metrics.CollectCycleDuration(time.Since(start))
	}()

	ids, err := dbInfo.WatchDB.NextSequences(context.Background(), f.key.ChainCollection(), eventLen)
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.key.ChainCollection(), err, rid)
		return true
	}

	chainNodes := make(map[string][]*watch.ChainNode, 0)
	oids := make([]string, eventLen)
	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	cursorMap := make(map[string]struct{})
	hitConflict := false
	for index, e := range es {
		// collect event's basic metrics
		f.metrics.CollectBasic(e)

		tenant, chainNode, detail, retry, err := f.parseEvent(dbInfo.CcDB, f.key, e, ids[index], rid)
		if err != nil {
			return retry
		}
		if chainNode == nil {
			continue
		}

		// if hit cursor conflict, the former cursor node's detail will be overwrite by the later one, so it
		// is not needed to remove the overlapped cursor node's detail again.
		ttl := time.Duration(f.key.TTLSeconds()) * time.Second
		pipe.Set(f.key.DetailKey(tenant, chainNode.Cursor), string(detail.eventInfo), ttl)
		pipe.Set(f.key.GeneralResDetailKey(tenant, chainNode), string(detail.resDetail), ttl)

		// validate if the cursor already exists in the batch, this happens when the concurrency is very high.
		// which will generate the same operation event with same cluster time, and generate with the same cursor
		// in the end. if this happens, the last event will be used finally, and the former events with the same
		// cursor will be dropped, and it's acceptable.
		exists := false
		if _, exists = cursorMap[chainNode.Cursor]; exists {
			hitConflict = true
		}

		// if the cursor is conflict with another cursor in the former batches, skip it
		if f.cursorQueue.checkIfConflict(dbInfo.UUID, chainNode.Cursor) && !exists {
			f.metrics.CollectConflict()
			continue
		}

		cursorMap[chainNode.Cursor] = struct{}{}
		oids[index] = e.ID()
		chainNodes[tenant] = append(chainNodes[tenant], chainNode)
	}
	lastTokenData := map[string]interface{}{
		common.BKTokenField:       es[eventLen-1].Token.Data,
		common.BKStartAtTimeField: es[eventLen-1].ClusterTime,
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		err = f.tokenHandler.setLastWatchToken(context.Background(), dbInfo.UUID, dbInfo.WatchDB, lastTokenData)
		if err != nil {
			f.metrics.CollectMongoError()
			return false
		}
		return false
	}

	// store details at first, in case those watching cmdb events read chain when details are not inserted yet
	if _, err := pipe.Exec(); err != nil {
		f.metrics.CollectRedisError()
		blog.Errorf("run flow, but insert details for %s failed, oids: %+v, err: %v, rid: %s,", f.key.Collection(),
			oids, err, rid)
		return true
	}

	if hitConflict {
		// update the chain nodes with picked chain nodes, so that we can handle them later.
		chainNodes = f.rearrangeEvents(chainNodes, rid)
	}

	retry, err = f.doInsertEvents(dbInfo, chainNodes, lastTokenData, rid)
	if err != nil {
		return retry
	}

	blog.Infof("insert watch event for %s success, oids: %v, rid: %s", f.key.Collection(), oids, rid)
	hasError = false
	return false
}

// rearrangeEvents remove the earlier chain nodes with the same cursor with a later one
func (f *Flow) rearrangeEvents(chainNodeMap map[string][]*watch.ChainNode, rid string) map[string][]*watch.ChainNode {
	pickedChainNodeMap := make(map[string][]*watch.ChainNode)

	for tenantID, chainNodes := range chainNodeMap {
		pickedChainNodes := make([]*watch.ChainNode, 0)
		conflictNodes := make([]*watch.ChainNode, 0)
		reminder := make(map[string]struct{})
		for i := len(chainNodes) - 1; i >= 0; i-- {
			chainNode := chainNodes[i]
			if _, exists := reminder[chainNode.Cursor]; exists {
				conflictNodes = append(conflictNodes, chainNode)
				// skip this event, because it has been replaced the one later.
				continue
			}

			reminder[chainNode.Cursor] = struct{}{}
			pickedChainNodes = append(pickedChainNodes, chainNode)
		}

		// reverse the picked chain nodes to their origin order
		for i, j := 0, len(pickedChainNodes)-1; i < j; i, j = i+1, j-1 {
			pickedChainNodes[i], pickedChainNodes[j] = pickedChainNodes[j], pickedChainNodes[i]
		}

		blog.WarnJSON("got tenant %s conflict cursor with chain nodes: %s, replaced with nodes: %s, rid: %s",
			tenantID, conflictNodes, pickedChainNodes, rid)

		pickedChainNodeMap[tenantID] = pickedChainNodes
	}

	return pickedChainNodeMap
}

func (f *Flow) doInsertEvents(dbInfo *types.DBInfo, chainNodeMap map[string][]*watch.ChainNode,
	lastTokenData map[string]interface{}, rid string) (bool, error) {

	if len(chainNodeMap) == 0 {
		return false, nil
	}
	coll := f.key.Collection()

	session, err := dbInfo.WatchDB.GetDBClient().StartSession()
	if err != nil {
		blog.Errorf("run flow, but start session failed, coll: %s, err: %v, rid: %s", coll, err, rid)
		return true, err
	}
	defer session.EndSession(context.Background())

	// retry insert the event node with remove the first event node,
	// which means the first one's cursor is conflicted with the former's batch operation inserted nodes.
	retryWithReduce := false
	var conflictTenantID string

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			blog.Errorf("run flow, but start transaction failed, coll: %s, err: %v, rid: %s", coll, err, rid)
			return err
		}

		err, retryWithReduce, conflictTenantID = f.insertChainNodes(sc, session, f.key, chainNodeMap, rid)
		if err != nil {
			return err
		}

		if err = f.tokenHandler.setLastWatchToken(sc, dbInfo.UUID, dbInfo.WatchDB, lastTokenData); err != nil {
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())
			return err
		}

		// Use context.Background() to ensure that the commit can complete successfully even if the context passed to
		// mongo.WithSession is changed to have a timeout.
		if err = session.CommitTransaction(context.Background()); err != nil {
			blog.Errorf("run flow, but commit mongo transaction failed, err: %v", err)
			f.metrics.CollectMongoError()
			return err
		}
		return nil
	})

	if txnErr != nil {
		blog.Errorf("do insert flow events failed, err: %v, rid: %s", txnErr, rid)
		if retryWithReduce {
			chainNodeMap = event.ReduceChainNode(chainNodeMap, conflictTenantID, coll, txnErr, f.metrics, rid)
			if len(chainNodeMap) == 0 {
				return false, nil
			}
			return f.doInsertEvents(dbInfo, chainNodeMap, lastTokenData, rid)
		}
		// if an error occurred, roll back and re-watch again
		return true, err
	}
	return false, nil
}

// insertChainNodes insert chain nodes and last event info into db
func (f *Flow) insertChainNodes(ctx context.Context, session mongo.Session, key event.Key,
	chainNodeMap map[string][]*watch.ChainNode, rid string) (error, bool, string) {

	for tenantID, chainNodes := range chainNodeMap {
		if len(chainNodes) == 0 {
			continue
		}

		shardingDB := mongodb.Dal("watch").Shard(sharding.NewShardOpts().WithTenant(tenantID))

		// insert chain nodes into db
		if err := shardingDB.Table(key.ChainCollection()).Insert(ctx, chainNodes); err != nil {
			blog.ErrorJSON("run flow, but insert tenant %s chain nodes for %s failed, nodes: %s, err: %s, rid: %s",
				tenantID, key.Collection(), chainNodes, err, rid)
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())

			if event.IsConflictError(err) {
				return err, true, tenantID
			}
			return err, false, ""
		}

		// set last watch event info for tenant
		lastNode := chainNodes[len(chainNodes)-1]
		lastNodeInfo := map[string]interface{}{
			common.BKFieldID:     lastNode.ID,
			common.BKCursorField: lastNode.Cursor,
		}

		filter := map[string]interface{}{
			"_id": key.Collection(),
		}
		if err := shardingDB.Table(common.BKTableNameLastWatchEvent).Update(ctx, filter, lastNodeInfo); err != nil {
			blog.Errorf("insert %s last event info(%+v) for coll %s failed, err: %v, rid: %s", tenantID, lastNodeInfo,
				key.Collection(), err, rid)
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())
			return err, false, ""
		}
	}

	return nil, false, ""
}
