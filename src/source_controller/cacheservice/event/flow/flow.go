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
	"strings"
	"sync"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"go.mongodb.org/mongo-driver/mongo"
)

type flowOptions struct {
	key         event.Key
	watch       stream.LoopInterface
	isMaster    discovery.ServiceManageInterface
	watchDB     *local.Mongo
	ccDB        dal.DB
	EventStruct interface{}
}

// oidCollKey key for oid to detail map. Since oid can duplicate in different collections, we need oid & coll for unique
type oidCollKey struct {
	oid  string
	coll string
}

func newFlow(ctx context.Context, opts flowOptions, getDeleteEventDetails getDeleteEventDetailsFunc,
	parseEvent parseEventFunc) error {

	flow, err := NewFlow(opts, getDeleteEventDetails, parseEvent)
	if err != nil {
		return err
	}

	return flow.RunFlow(ctx)
}

// NewFlow create a new event watch flow
func NewFlow(opts flowOptions, getDelEventDetails getDeleteEventDetailsFunc, parseEvent parseEventFunc) (Flow, error) {
	if getDelEventDetails == nil {
		return Flow{}, fmt.Errorf("getDeleteEventDetailsFunc is not set, key: %s", opts.key.Namespace())
	}

	if parseEvent == nil {
		return Flow{}, fmt.Errorf("parseEventFunc is not set, key: %s", opts.key.Namespace())
	}

	return Flow{
		flowOptions:           opts,
		metrics:               event.InitialMetrics(opts.key.Collection(), "watch"),
		getDeleteEventDetails: getDelEventDetails,
		parseEvent:            parseEvent,
		cursorQueue: &cursorQueue{
			cursorQueue: make(map[string]string),
		},
	}, nil
}

// Flow TODO
type Flow struct {
	flowOptions
	metrics               *event.EventMetrics
	tokenHandler          *flowTokenHandler
	getDeleteEventDetails getDeleteEventDetailsFunc
	parseEvent            parseEventFunc
	cursorQueue           *cursorQueue
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
func (c *cursorQueue) checkIfConflict(cursor string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, exists := c.cursorQueue[cursor]; exists {
		return true
	}

	if c.length <= 0 {
		c.head = cursor
		c.tail = cursor
		c.cursorQueue[cursor] = ""
		c.length++
		return false
	}

	// append cursor to the tail of the cursor queue
	c.cursorQueue[c.tail] = cursor
	c.cursorQueue[cursor] = ""
	c.tail = cursor

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

// RunFlow TODO
func (f *Flow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.watchDB, f.metrics)

	startAtTime, err := f.tokenHandler.getStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get start watch time for %s failed, err: %v", f.key.Collection(), err)
		return err
	}

	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct:             f.EventStruct,
			Collection:              f.key.Collection(),
			StartAfterToken:         nil,
			StartAtTime:             startAtTime,
			WatchFatalErrorCallback: f.tokenHandler.resetWatchToken,
		},
	}

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         f.key.Namespace(),
			WatchOpt:     watchOpts,
			TokenHandler: f.tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 10,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: batchSize,
	}

	if err := f.watch.WithBatch(opts); err != nil {
		blog.Errorf("run flow, but watch batch failed, err: %v", err)
		return err
	}

	return nil
}

func (f *Flow) doBatch(es []*types.Event) (retry bool) {
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

	oidDetailMap, retry, err := f.getDeleteEventDetails(es, f.ccDB, f.metrics)
	if err != nil {
		blog.Errorf("get deleted event details failed, err: %v, rid: %s", err, rid)
		return retry
	}

	ids, err := f.watchDB.NextSequences(context.Background(), f.key.ChainCollection(), eventLen)
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.key.ChainCollection(), err, rid)
		return true
	}

	chainNodes := make([]*watch.ChainNode, 0)
	oids := make([]string, eventLen)
	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	cursorMap := make(map[string]struct{})
	hitConflict := false
	for index, e := range es {
		// collect event's basic metrics
		f.metrics.CollectBasic(e)

		chainNode, detailBytes, retry, err := f.parseEvent(f.ccDB, f.key, e, oidDetailMap, ids[index], rid)
		if err != nil {
			return retry
		}
		if chainNode == nil {
			continue
		}

		// if hit cursor conflict, the former cursor node's detail will be overwrite by the later one, so it
		// is not needed to remove the overlapped cursor node's detail again.
		pipe.Set(f.key.DetailKey(chainNode.Cursor), string(detailBytes), time.Duration(f.key.TTLSeconds())*time.Second)

		// validate if the cursor already exists in the batch, this happens when the concurrency is very high.
		// which will generate the same operation event with same cluster time, and generate with the same cursor
		// in the end. if this happens, the last event will be used finally, and the former events with the same
		// cursor will be dropped, and it's acceptable.
		exists := false
		if _, exists = cursorMap[chainNode.Cursor]; exists {
			hitConflict = true
		}

		// if the cursor is conflict with another cursor in the former batches, skip it
		if f.cursorQueue.checkIfConflict(chainNode.Cursor) && !exists {
			f.metrics.CollectConflict()
			continue
		}

		cursorMap[chainNode.Cursor] = struct{}{}
		oids[index] = e.ID()
		chainNodes = append(chainNodes, chainNode)
	}
	lastTokenData := map[string]interface{}{
		common.BKTokenField:       es[eventLen-1].Token.Data,
		common.BKStartAtTimeField: es[eventLen-1].ClusterTime,
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		if err := f.tokenHandler.setLastWatchToken(context.Background(), lastTokenData); err != nil {
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

	retry, err = f.doInsertEvents(chainNodes, lastTokenData, rid)
	if err != nil {
		return retry
	}

	blog.Infof("insert watch event for %s success, oids: %v, rid: %s", f.key.Collection(), oids, rid)
	hasError = false
	return false
}

// rearrangeEvents remove the earlier chain nodes with the same cursor with a later one
func (f *Flow) rearrangeEvents(chainNodes []*watch.ChainNode, rid string) []*watch.ChainNode {
	pickedChainNodes := make([]*watch.ChainNode, 0)
	conflictNodes := make([]*watch.ChainNode, 0)
	reminder := make(map[string]struct{})
	for i := len(chainNodes) - 1; i >= 0; i-- {
		chainNode := chainNodes[i]
		if _, exists := reminder[chainNode.Cursor]; exists {
			conflictNodes = append(conflictNodes, chainNode)
			// skip this event, because it has been replaced the the one later.
			continue
		}

		reminder[chainNode.Cursor] = struct{}{}
		pickedChainNodes = append(pickedChainNodes, chainNode)
	}

	// reverse the picked chain nodes to their origin order
	for i, j := 0, len(pickedChainNodes)-1; i < j; i, j = i+1, j-1 {
		pickedChainNodes[i], pickedChainNodes[j] = pickedChainNodes[j], pickedChainNodes[i]
	}

	blog.WarnJSON("got conflict got conflict cursor with chain nodes: %s, replaced with nodes: %s, rid: %s",
		conflictNodes, pickedChainNodes, rid)

	return pickedChainNodes
}

func (f *Flow) doInsertEvents(chainNodes []*watch.ChainNode, lastTokenData map[string]interface{}, rid string) (
	bool, error) {

	count := len(chainNodes)

	if count == 0 {
		return false, nil
	}

	watchDBClient := f.watchDB.GetDBClient()

	session, err := watchDBClient.StartSession()
	if err != nil {
		blog.Errorf("run flow, but start session failed, coll: %s, err: %v, rid: %s", f.key.Collection(), err, rid)
		return true, err
	}
	defer session.EndSession(context.Background())

	// retry insert the event node with remove the first event node,
	// which means the first one's cursor is conflicted with the former's batch operation inserted nodes.
	retryWithReduce := false

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			blog.Errorf("run flow, but start transaction failed, coll: %s, err: %v, rid: %s", f.key.Collection(),
				err, rid)
			return err
		}

		if err := f.watchDB.Table(f.key.ChainCollection()).Insert(sc, chainNodes); err != nil {
			blog.ErrorJSON("run flow, but insert chain nodes for %s failed, nodes: %s, err: %v, rid: %s",
				f.key.Collection(), chainNodes, err, rid)
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())

			if event.IsConflictError(err) {
				// set retry with reduce flag and retry later
				retryWithReduce = true
			}
			return err
		}

		lastNode := chainNodes[len(chainNodes)-1]
		lastTokenData[common.BKFieldID] = lastNode.ID
		lastTokenData[common.BKCursorField] = lastNode.Cursor
		lastTokenData[common.BKStartAtTimeField] = lastNode.ClusterTime
		if err := f.tokenHandler.setLastWatchToken(sc, lastTokenData); err != nil {
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

		rid = rid + ":" + chainNodes[0].Oid
		if retryWithReduce {
			monitor.Collect(&meta.Alarm{
				RequestID: rid,
				Type:      meta.EventFatalError,
				Detail:    fmt.Sprintf("run event flow, but got conflict %s cursor with chain nodes", f.key.Collection()),
				Module:    types2.CC_MODULE_CACHESERVICE,
				Dimension: map[string]string{"retry_conflict_nodes": "yes"},
			})

			if len(chainNodes) <= 1 {
				return false, nil
			}

			for index, reducedChainNode := range chainNodes {
				if isConflictChainNode(reducedChainNode, txnErr) {
					f.metrics.CollectConflict()
					chainNodes = append(chainNodes[:index], chainNodes[index+1:]...)

					// need do with retry with reduce
					blog.ErrorJSON("run flow, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
						f.key.Collection(), reducedChainNode, chainNodes, rid)

					return f.doInsertEvents(chainNodes, lastTokenData, rid)
				}
			}

			// when no cursor conflict node is found, discard the first node and try to insert the others
			blog.ErrorJSON("run flow, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
				f.key.Collection(), chainNodes[0], chainNodes[1:], rid)

			return f.doInsertEvents(chainNodes[1:], lastTokenData, rid)
		}

		// if an error occurred, roll back and re-watch again
		return true, err
	}

	return false, nil
}

func isConflictChainNode(chainNode *watch.ChainNode, err error) bool {
	return strings.Contains(err.Error(), chainNode.Cursor) && strings.Contains(err.Error(), "index_cursor")
}
