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
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
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

// MixEventFlowOptions mix event flow options
type MixEventFlowOptions struct {
	MixKey       event.Key
	Key          event.Key
	WatchFields  []string
	Watch        stream.LoopInterface
	WatchDB      *local.Mongo
	CcDB         dal.DB
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
type parseEventFunc func(e *types.Event, id uint64, rid string) (*watch.ChainNode, []byte, bool, error)

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
	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct:     &es,
			Collection:      f.Key.Collection(),
			StartAfterToken: nil,
		},
	}
	if f.Key.Collection() == common.BKTableNameBaseHost {
		watchOpts.EventStruct = new(metadata.HostMapStr)
	}

	f.tokenHandler = newMixEventTokenHandler(f.MixKey, f.Key, f.WatchDB, f.metrics)

	startAtTime, err := f.tokenHandler.getStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get start watch time for %s failed, err: %v", f.Key.Collection(), err)
		return err
	}
	watchOpts.StartAtTime = startAtTime
	watchOpts.WatchFatalErrorCallback = f.tokenHandler.resetWatchToken
	watchOpts.Fields = f.WatchFields

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name:         fmt.Sprintf("%s_%s", f.MixKey.Namespace(), f.Key.Namespace()),
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

	if err := f.Watch.WithBatch(opts); err != nil {
		blog.Errorf("watch %s events, but watch batch failed, err: %v", f.MixKey.Namespace(), err)
		return err
	}

	return nil
}

// doBatch batch handle events
func (f *MixEventFlow) doBatch(es []*types.Event) (retry bool) {
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

	// get the lock to get sequences ids.
	// otherwise, we can not guarantee the multiple event's id is in the right order/sequences
	// it should be a natural increase order.
	if err = f.getLock(rid); err != nil {
		blog.Errorf("get %s lock failed, err: %v, rid: %s", f.MixKey.Namespace(), err, rid)
		return true
	}

	// release the lock when the job is done or failed.
	defer f.releaseLock(rid)

	// last event in original events is used to generate
	lastEvent := es[len(es)-1]
	lastTokenData := mapstr.MapStr{
		common.BKTokenField:       lastEvent.Token.Data,
		common.BKStartAtTimeField: lastEvent.ClusterTime,
	}

	// handle the rearranged events
	retry, err = f.handleEvents(events, lastTokenData, rid)
	if err != nil {
		return retry
	}

	hasError = false
	return false
}

// handleEvents handle the rearranged events, parse them into chain nodes and details, then insert into db and redis
func (f *MixEventFlow) handleEvents(events []*types.Event, lastTokenData mapstr.MapStr, rid string) (bool, error) {
	eventIDs, err := f.WatchDB.NextSequences(context.Background(), f.MixKey.Collection(), len(events))
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.Key.ChainCollection(), err, rid)
		return true, err
	}

	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	needSaveDetails := false

	chainNodes := make([]*watch.ChainNode, 0)
	oids := make([]string, 0)
	cursorMap := make(map[string]struct{})
	for index, e := range events {
		// collect event's basic metrics
		f.metrics.CollectBasic(e)
		oids = append(oids, e.ID())

		chainNode, detailBytes, retry, err := f.parseEvent(e, eventIDs[index], rid)
		if err != nil {
			return retry, err
		}

		if chainNode == nil {
			continue
		}

		if len(detailBytes) > 0 {
			// if hit cursor conflict, the former cursor node's detail will be overwritten by the later one, so it
			// is not needed to remove the overlapped cursor node's detail again.
			pipe.Set(f.MixKey.DetailKey(chainNode.Cursor), string(detailBytes),
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

		chainNodes = append(chainNodes, chainNode)
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		if err := f.tokenHandler.setLastWatchToken(context.Background(), lastTokenData); err != nil {
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

	retry, err := f.doInsertEvents(chainNodes, lastTokenData, rid)
	if err != nil {
		return retry, err
	}

	blog.Infof("insert %s event for %s success, oids: %v, rid: %s", f.MixKey.Namespace(), f.Key.Collection(), oids, rid)
	return false, nil

}

func (f *MixEventFlow) doInsertEvents(chainNodes []*watch.ChainNode, lastTokenData map[string]interface{}, rid string) (
	bool, error) {

	count := len(chainNodes)

	if count == 0 {
		return false, nil
	}

	watchDBClient := f.WatchDB.GetDBClient()

	session, err := watchDBClient.StartSession()
	if err != nil {
		blog.Errorf("watch %s events, but start session failed, coll: %s, err: %v, rid: %s", f.MixKey.Namespace(),
			f.Key.Collection(), err, rid)
		return true, err
	}
	defer session.EndSession(context.Background())

	// insert events into db in an transaction
	txnErr, conflictError := f.insertEvents(session, chainNodes, lastTokenData, rid)

	if txnErr != nil {
		blog.Errorf("do insert %s events failed, err: %v, rid: %s", f.MixKey.Namespace(), txnErr, rid)

		rid = rid + ":" + chainNodes[0].Oid
		if conflictError != nil && len(chainNodes) >= 1 {
			monitor.Collect(&meta.Alarm{
				RequestID: rid,
				Type:      meta.EventFatalError,
				Detail: fmt.Sprintf("host identifier, but got conflict %s cursor with chain nodes",
					f.Key.Collection()),
				Module:    types2.CC_MODULE_CACHESERVICE,
				Dimension: map[string]string{"retry_conflict_nodes": "yes"},
			})

			var conflictNode *watch.ChainNode
			// get the conflict cursor
			for idx := range chainNodes {
				if strings.Contains(conflictError.Error(), chainNodes[idx].Cursor) {
					// record conflict node
					conflictNode = chainNodes[idx]
					// remove the conflict cursor
					chainNodes = append(chainNodes[0:idx], chainNodes[idx+1:]...)
					break
				}
			}

			if conflictNode == nil {
				// this should not happen
				// reduce event's one by one, then retry again.
				blog.ErrorJSON("watch %s events, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
					f.MixKey.Namespace(), f.Key.Collection(), chainNodes[0], chainNodes[1:], rid)

				// retry insert events
				return f.doInsertEvents(chainNodes[1:], lastTokenData, rid)
			}

			blog.ErrorJSON("watch %s events, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
				f.MixKey.Namespace(), f.Key.Collection(), conflictNode, chainNodes, rid)

			// retry insert events
			return f.doInsertEvents(chainNodes, lastTokenData, rid)
		}

		// if an error occurred, roll back and re-watch again
		blog.Warnf("do insert %s events, do retry insert with rid: %s", f.MixKey.Namespace(), rid)
		return true, err
	}

	return false, nil
}

// insertEvents insert events and last watch token
func (f *MixEventFlow) insertEvents(session mongo.Session, chainNodes []*watch.ChainNode,
	lastTokenData map[string]interface{}, rid string) (error, error) {

	// conflictError record the conflict cursor error
	var conflictError error

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			blog.Errorf("watch %s events, but start transaction failed, coll: %s, err: %v, rid: %s",
				f.MixKey.Namespace(), f.Key.Collection(), err, rid)
			return err
		}

		if err := f.WatchDB.Table(f.MixKey.ChainCollection()).Insert(sc, chainNodes); err != nil {
			blog.ErrorJSON("watch %s events, but insert chain nodes for %s failed, nodes: %s, err: %v, rid: %s",
				f.Key.Collection(), f.MixKey.Namespace(), chainNodes, err, rid)
			f.metrics.CollectMongoError()
			_ = session.AbortTransaction(context.Background())

			if event.IsConflictError(err) {
				conflictError = err
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
		if err := session.CommitTransaction(context.Background()); err != nil {
			blog.Errorf("watch %s events, but commit mongo transaction failed, err: %v", f.MixKey.Namespace(), err)
			f.metrics.CollectMongoError()
			return err
		}
		return nil
	})

	return txnErr, conflictError
}

func (f *MixEventFlow) getLock(rid string) error {
	timeout := time.After(f.EventLockTTL)
	for {
		select {
		case <-timeout:
			return fmt.Errorf("get %s: %s lock timeout", f.MixKey.Namespace(), f.EventLockKey)
		default:
		}

		success, err := redis.Client().SetNX(context.Background(), f.EventLockKey, 1, f.EventLockTTL).Result()
		if err != nil {
			blog.Errorf("get %s: %s lock, err: %v, rid: %s", f.MixKey.Namespace(), f.EventLockKey, err, rid)
			return err
		}

		if !success {
			blog.V(3).Infof("get %s: %s lock failed, retry later, rid: %s", f.MixKey.Namespace(), f.EventLockKey, rid)
			time.Sleep(300 * time.Millisecond)
			continue
		}
		// get lock success.
		return nil
	}

}

func (f *MixEventFlow) releaseLock(rid string) {
	_, err := redis.Client().Del(context.Background(), f.EventLockKey).Result()
	if err != nil {
		blog.Errorf("delete %s lock key: %s failed, err: %v, rid: %s", f.MixKey.Namespace(), f.EventLockKey, err, rid)
		return
	}
	return
}
