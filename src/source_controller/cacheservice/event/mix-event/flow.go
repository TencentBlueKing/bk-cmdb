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

package mixevent

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/task"
	"configcenter/src/storage/stream/types"

	"go.mongodb.org/mongo-driver/mongo"
)

// MixEventFlowOptions mix event flow options
type MixEventFlowOptions struct {
	MixKey       event.Key
	Key          event.Key
	WatchFields  []string
	Task         *task.Task
	EventLockTTL time.Duration
	EventLockKey string
}

// MixEventFlow mix event flow
type MixEventFlow struct {
	MixEventFlowOptions
	metrics         *event.EventMetrics
	tokenHandler    *mixEventHandler
	rearrangeEvents rearrangeEventsFunc
	parseEvent      parseEventFunc
}

// rearrangeEventsFunc function type for rearranging mix events
type rearrangeEventsFunc func(rid string, es []*types.Event) ([]*types.Event, error)

// parseEventFunc function type for parsing mix event into chain node and detail
type parseEventFunc func(e *types.Event, id uint64, rid string) (string, *watch.ChainNode, []byte, bool, error)

// NewMixEventFlow create a new mix event watch flow
func NewMixEventFlow(opts MixEventFlowOptions, rearrangeEvents rearrangeEventsFunc, parseEvent parseEventFunc) (
	MixEventFlow, error) {

	if rearrangeEvents == nil {
		return MixEventFlow{}, fmt.Errorf("rearrangeEvents is not set, mix key: %s, key: %s", opts.MixKey.Namespace(),
			opts.Key.Namespace())
	}

	if parseEvent == nil {
		return MixEventFlow{}, fmt.Errorf("parseEvent is not set, mix key: %s, key: %s", opts.MixKey.Namespace(),
			opts.Key.Namespace())
	}

	return MixEventFlow{
		MixEventFlowOptions: opts,
		metrics:             event.InitialMetrics(opts.Key.Collection(), opts.MixKey.Namespace()),
		rearrangeEvents:     rearrangeEvents,
		parseEvent:          parseEvent,
	}, nil
}

const batchSize = 500

// RunFlow run mix event flow
func (f *MixEventFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run %s event flow for key: %s.", f.MixKey.Namespace(), f.Key.Namespace())

	es := make(map[string]interface{})

	f.tokenHandler = newMixEventTokenHandler(f.MixKey, f.Key, f.metrics)

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name: fmt.Sprintf("%s_%s", f.MixKey.Namespace(), f.Key.Namespace()),
			CollOpts: &types.WatchCollOptions{
				CollectionOptions: types.CollectionOptions{
					CollectionFilter: &types.CollectionFilter{
						Regex: fmt.Sprintf("_%s$", f.Key.Collection()),
					},
					EventStruct: &es,
					Fields:      f.WatchFields,
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

	if f.Key.Collection() == common.BKTableNameBaseHost {
		opts.CollOpts.EventStruct = new(metadata.HostMapStr)
	}

	err := f.Task.AddLoopBatchTask(opts)
	if err != nil {
		blog.Errorf("watch %s events, but add watch batch task failed, err: %v", f.MixKey.Namespace(), err)
		return err
	}
	return nil
}

// doBatch batch handle events
func (f *MixEventFlow) doBatch(dbInfo *types.DBInfo, es []*types.Event) (retry bool) {
	if len(es) == 0 {
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

	// rearranging mix events
	events, err := f.rearrangeEvents(rid, es)
	if err != nil {
		blog.Errorf("rearrange %s events failed, will retry, err: %v, rid: %s", f.MixKey.Namespace(), err, rid)
		return true
	}

	if len(events) == 0 {
		return false
	}

	// get the lock to get sequences ids.
	// otherwise, we can not guarantee the multiple event's id is in the right order/sequences
	// it should be a natural increase order.
	if err = f.getLock(dbInfo.UUID, rid); err != nil {
		blog.Errorf("get %s lock failed, err: %v, rid: %s", f.MixKey.Namespace(), err, rid)
		return true
	}

	// release the lock when the job is done or failed.
	defer f.releaseLock(dbInfo.UUID, rid)

	// last event in original events is used to generate
	lastEvent := es[len(es)-1]
	lastTokenData := mapstr.MapStr{
		common.BKTokenField:       lastEvent.Token.Data,
		common.BKStartAtTimeField: lastEvent.ClusterTime,
	}

	// handle the rearranged events
	retry, err = f.handleEvents(dbInfo, events, lastTokenData, rid)
	if err != nil {
		return retry
	}

	hasError = false
	return false
}

// handleEvents handle the rearranged events, parse them into chain nodes and details, then insert into db and redis
func (f *MixEventFlow) handleEvents(dbInfo *types.DBInfo, events []*types.Event, lastTokenData mapstr.MapStr,
	rid string) (bool, error) {

	eventIDs, err := dbInfo.WatchDB.NextSequences(context.Background(), f.MixKey.Collection(), len(events))
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.Key.ChainCollection(), err, rid)
		return true, err
	}

	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	needSaveDetails := false

	chainNodes := make(map[string][]*watch.ChainNode, 0)
	oids := make([]string, 0)
	cursorMap := make(map[string]struct{})
	for index, e := range events {
		// collect event's basic metrics
		f.metrics.CollectBasic(e)
		oids = append(oids, e.ID())

		tenantID, chainNode, detailBytes, retry, err := f.parseEvent(e, eventIDs[index], rid)
		if err != nil {
			return retry, err
		}

		if chainNode == nil {
			continue
		}

		if len(detailBytes) > 0 {
			// if hit cursor conflict, the former cursor node's detail will be overwritten by the later one, so it
			// is not needed to remove the overlapped cursor node's detail again.
			pipe.Set(f.MixKey.DetailKey(tenantID, chainNode.Cursor), string(detailBytes),
				time.Duration(f.MixKey.TTLSeconds())*time.Second)
			needSaveDetails = true
		}

		// validate if the cursor is already exists, this is happens when the concurrent operation is very high.
		// which will generate the same operation event with same cluster time, and generate with the same cursor
		// in the end. if this happens, drop this event directly, because we only care this host's identifier is
		// changed or not.
		if _, exists := cursorMap[chainNode.Cursor]; exists {
			// skip this event.
			continue
		}
		cursorMap[chainNode.Cursor] = struct{}{}

		chainNodes[tenantID] = append(chainNodes[tenantID], chainNode)
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		err = f.tokenHandler.setLastWatchToken(context.Background(), dbInfo.UUID, dbInfo.WatchDB, lastTokenData)
		if err != nil {
			f.metrics.CollectMongoError()
			return false, err
		}
		return false, err
	}

	// store details at first, in case those watching cmdb events read chain when details are not inserted yet
	if needSaveDetails {
		if _, err := pipe.Exec(); err != nil {
			f.metrics.CollectRedisError()
			blog.Errorf("run flow, but insert details for %s failed, oids: %+v, err: %v, rid: %s,", f.Key.Collection(),
				oids, err, rid)
			return true, err
		}
	}

	retry, err := f.doInsertEvents(dbInfo, chainNodes, lastTokenData, rid)
	if err != nil {
		return retry, err
	}

	blog.Infof("insert %s event for %s success, oids: %v, rid: %s", f.MixKey.Namespace(), f.Key.Collection(), oids, rid)
	return false, nil
}

func (f *MixEventFlow) doInsertEvents(dbInfo *types.DBInfo, chainNodeMap map[string][]*watch.ChainNode,
	lastTokenData map[string]interface{}, rid string) (bool, error) {

	if len(chainNodeMap) == 0 {
		return false, nil
	}

	session, err := dbInfo.WatchDB.GetDBClient().StartSession()
	if err != nil {
		blog.Errorf("watch %s events, but start session failed, coll: %s, err: %v, rid: %s", f.MixKey.Namespace(),
			f.Key.Collection(), err, rid)
		return true, err
	}
	defer session.EndSession(context.Background())

	// insert events into db in an transaction
	txnErr, conflictError, conflictTenantID := f.insertEvents(dbInfo, session, chainNodeMap, lastTokenData, rid)
	if txnErr != nil {
		blog.Errorf("do insert %s events failed, err: %v, rid: %s", f.MixKey.Namespace(), txnErr, rid)

		if conflictError != nil && len(chainNodeMap[conflictTenantID]) >= 1 {
			chainNodeMap = event.ReduceChainNode(chainNodeMap, conflictTenantID,
				f.MixKey.Namespace()+":"+f.Key.Collection(), txnErr, f.metrics, rid)
			return f.doInsertEvents(dbInfo, chainNodeMap, lastTokenData, rid)
		}

		// if an error occurred, roll back and re-watch again
		blog.Warnf("do insert %s events, do retry insert with rid: %s", f.MixKey.Namespace(), rid)
		return true, err
	}

	return false, nil
}

// insertEvents insert events and last watch token
func (f *MixEventFlow) insertEvents(dbInfo *types.DBInfo, session mongo.Session,
	chainNodeMap map[string][]*watch.ChainNode, lastToken map[string]interface{}, rid string) (error, error, string) {

	// conflictError record the conflict cursor error
	var conflictError error
	var conflictTenantID string

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			blog.Errorf("watch %s events, but start transaction failed, coll: %s, err: %v, rid: %s",
				f.MixKey.Namespace(), f.Key.Collection(), err, rid)
			return err
		}

		for tenantID, chainNodes := range chainNodeMap {
			if len(chainNodes) == 0 {
				continue
			}

			shardingDB := mongodb.Dal("watch").Shard(sharding.NewShardOpts().WithTenant(tenantID))

			// insert chain nodes into db
			if err := shardingDB.Table(f.MixKey.ChainCollection()).Insert(sc, chainNodes); err != nil {
				blog.ErrorJSON("watch %s events, but insert chain nodes for %s failed, nodes: %s, err: %v, rid: %s",
					f.Key.Collection(), f.MixKey.Namespace(), chainNodes, err, rid)
				f.metrics.CollectMongoError()
				_ = session.AbortTransaction(context.Background())

				if event.IsConflictError(err) {
					conflictError = err
				}
				return err
			}

			// set last watch event info
			lastNode := chainNodes[len(chainNodes)-1]
			lastNodeInfo := map[string]interface{}{
				common.BKFieldID:     lastNode.ID,
				common.BKCursorField: lastNode.Cursor,
			}

			filter := map[string]interface{}{
				"_id":            f.MixKey.Collection(),
				common.BKFieldID: mapstr.MapStr{common.BKDBLT: lastNode.ID},
			}

			if err := shardingDB.Table(common.BKTableNameLastWatchEvent).Update(sc, filter, lastNodeInfo); err != nil {
				blog.Errorf("insert tenant %s mix event %s coll %s last event info(%+v) failed, err: %v, rid: %s",
					tenantID, f.MixKey.Namespace(), f.Key.Collection(), lastNodeInfo, err, rid)
				f.metrics.CollectMongoError()
				_ = session.AbortTransaction(context.Background())
				return err
			}
		}

		if err := f.tokenHandler.setLastWatchToken(sc, dbInfo.UUID, dbInfo.WatchDB, lastToken); err != nil {
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())
			return err
		}

		// Use context.Background() to ensure that the commit can complete successfully even if the context passed to
		// mongo.WithSession is changed to have a timeout.
		if err := session.CommitTransaction(context.Background()); err != nil {
			blog.Errorf("watch %s events, but commit mongo transaction failed, err: %v", f.MixKey.Namespace(), err)
			f.metrics.CollectMongoError()
			return err
		}
		return nil
	})

	return txnErr, conflictError, conflictTenantID
}

func (f *MixEventFlow) getLock(uuid, rid string) error {
	lockKey := f.EventLockKey + ":" + uuid
	timeout := time.After(f.EventLockTTL)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("get %s: %s lock timeout", f.MixKey.Namespace(), lockKey)
		default:
		}

		success, err := redis.Client().SetNX(context.Background(), lockKey, 1, f.EventLockTTL).Result()
		if err != nil {
			blog.Errorf("get %s: %s lock, err: %v, rid: %s", f.MixKey.Namespace(), lockKey, err, rid)
			return err
		}

		if !success {
			blog.V(3).Infof("get %s: %s lock failed, retry later, rid: %s", f.MixKey.Namespace(), lockKey, rid)
			time.Sleep(300 * time.Millisecond)
			continue
		}
		// get lock success.
		return nil
	}

}

func (f *MixEventFlow) releaseLock(uuid, rid string) {
	lockKey := f.EventLockKey + ":" + uuid
	_, err := redis.Client().Del(context.Background(), lockKey).Result()
	if err != nil {
		blog.Errorf("delete %s lock key: %s failed, err: %v, rid: %s", f.MixKey.Namespace(), lockKey, err, rid)
		return
	}
	return
}
