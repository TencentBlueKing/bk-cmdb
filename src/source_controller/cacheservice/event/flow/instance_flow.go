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

package flow

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
)

func newInstanceFlow(ctx context.Context, opts flowOptions) error {
	flow := InstanceFlow{
		Flow: Flow{
			flowOptions: opts,
			metrics:     event.InitialMetrics(opts.key.Collection(), "watch"),
		},
		collKeyMap: map[string]event.Key{
			event.MainlineInstanceKey.Collection(): event.MainlineInstanceKey,
		},
	}

	// refresh mainline object ID map every 5 minutes
	var err error
	flow.mainlineObjectMap, err = flow.getMainlineObjectMap(ctx)
	if err != nil {
		blog.Errorf("run object instance watch, but get mainline objects failed, err: %v", err)
		return err
	}

	go func() {
		t := time.NewTicker(time.Minute * 5)
		for range t.C {
			flow.mainlineObjectMap, err = flow.getMainlineObjectMap(ctx)
			if err != nil {
				blog.Errorf("run object instance watch, but get mainline objects failed, err: %v", err)
			}
		}
	}()

	return flow.RunFlow(ctx)
}

type InstanceFlow struct {
	Flow
	mainlineObjectMap map[string]struct{}
	// collKeyMap the collection to key map of the sub keys, not including the flow's default key
	collKeyMap map[string]event.Key
}

func (f *InstanceFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	opts, err := f.generateLoopBatchOptions(ctx)
	if err != nil {
		return err
	}

	// watch all tables with the prefix of instance table
	opts.WatchOpt.Collection = ""
	opts.WatchOpt.CollectionFilter = map[string]interface{}{
		common.BKDBLIKE: event.ObjInstTablePrefixRegex,
	}
	opts.EventHandler.DoBatch = f.doBatchWrapper(f.doBatch)

	if err := f.watch.WithBatch(opts); err != nil {
		blog.Errorf("run flow, but watch batch failed, err: %v", err)
		return err
	}

	return nil
}

func (f *InstanceFlow) doBatch(es []*types.Event, rid string) (bool, error) {
	oidDetailMap, retry, err := f.getDeleteEventDetails(es)
	if err != nil {
		blog.Errorf("get deleted event details failed, err: %v, rid: %s", err, rid)
		return retry, err
	}

	eventMap := f.classifyEvents(es, oidDetailMap, rid)
	for coll, events := range eventMap {
		key, exists := f.collKeyMap[coll]
		if !exists {
			// if no sub keys are hit, use the default key to handle events
			retry, err = f.batchHandleEvents(events, oidDetailMap, rid)
			if err != nil {
				return retry, err
			}
			continue
		}

		flow := InstanceFlow{
			Flow: Flow{
				flowOptions:  f.flowOptions,
				metrics:      f.metrics,
				tokenHandler: f.tokenHandler,
			},
		}
		flow.key = key
		retry, err = flow.batchHandleEvents(events, oidDetailMap, rid)
		if err != nil {
			return retry, err
		}
	}

	return false, nil
}

func (f *InstanceFlow) getMainlineObjectMap(ctx context.Context) (map[string]struct{}, error) {
	relations := make([]metadata.Association, 0)
	filter := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := f.ccDB.Table(common.BKTableNameObjAsst).Find(filter).Fields(common.BKObjIDField).All(ctx, &relations)
	if err != nil {
		blog.Errorf("get mainline topology association failed, err: %v", err)
		return nil, err
	}

	objIDMap := make(map[string]struct{}, 0)
	for _, relation := range relations {
		if common.IsInnerMainlineModel(relation.ObjectID) {
			continue
		}
		objIDMap[relation.ObjectID] = struct{}{}
	}
	return objIDMap, nil
}

// classifyEvents classify events by their related key's collection
func (f *InstanceFlow) classifyEvents(es []*types.Event, oidDetailMap map[string][]byte,
	rid string) map[string][]*types.Event {

	mainlineColl := event.MainlineInstanceKey.Collection()
	commonColl := f.key.Collection()

	eventMap := make(map[string][]*types.Event)
	for _, e := range es {
		if e.OperationType == types.Delete {
			doc, exist := oidDetailMap[e.Oid]
			if !exist {
				blog.Errorf("run flow, received %s event, but delete doc[oid: %s] detail not exists, rid: %s",
					f.key.Collection(), e.Oid, rid)
				continue
			}
			e.DocBytes = doc
		}

		objID := gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()
		if _, exists := f.mainlineObjectMap[objID]; exists {
			eventMap[mainlineColl] = append(eventMap[mainlineColl], e)
			continue
		}
		eventMap[commonColl] = append(eventMap[commonColl], e)
	}

	return eventMap
}

func (f *InstanceFlow) batchHandleEvents(es []*types.Event, oidDetailMap map[string][]byte, rid string) (bool, error) {
	ids, err := f.watchDB.NextSequences(context.Background(), f.key.ChainCollection(), len(es))
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.key.ChainCollection(), err, rid)
		return true, err
	}

	chainNodes := make([]*watch.ChainNode, 0)
	oids := make([]string, len(es))
	// process events into db chain nodes to store in db and details to store in redis
	pipe := redis.Client().Pipeline()
	lastTokenData := make(map[string]interface{})
	cursorMap := make(map[string]struct{})
	hitConflict := false
	for index, e := range es {
		// collect event's basic metrics
		f.metrics.CollectBasic(e)
		lastTokenData[common.BKTokenField] = e.Token.Data
		lastTokenData[common.BKStartAtTimeField] = e.ClusterTime

		chainNode, detailBytes, retry, err := f.parseEvent(e, oidDetailMap, ids[index], rid)
		if err != nil {
			return retry, err
		}

		// set sub resource for chain node, the sub resource of instance watch is object id
		chainNode.SubResource = gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()

		// validate if the cursor is already exists, this is happens when the concurrent operation is very high.
		// which will generate the same operation event with same cluster time, and generate with the same cursor
		// in the end. if this happens, the last event will be used finally, and the former events with the same
		// cursor will be dropped, and it's acceptable.
		if _, exists := cursorMap[chainNode.Cursor]; exists {
			hitConflict = true
		}
		cursorMap[chainNode.Cursor] = struct{}{}

		oids[index] = e.ID()
		chainNodes = append(chainNodes, chainNode)

		// if hit cursor conflict, the former cursor node's detail will be overwrite by the later one, so it
		// is not needed to remove the overlapped cursor node's detail again.
		pipe.Set(f.key.DetailKey(chainNode.Cursor), string(detailBytes), time.Duration(f.key.TTLSeconds())*time.Second)
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodes) == 0 {
		if err := f.tokenHandler.setLastWatchToken(context.Background(), lastTokenData); err != nil {
			f.metrics.CollectMongoError()
			return false, err
		}
		return false, nil
	}

	// store details at first, in case those watching cmdb events read chain when details are not inserted yet
	if _, err := pipe.Exec(); err != nil {
		f.metrics.CollectRedisError()
		blog.Errorf("run flow, but insert details for %s failed, oids: %+v, err: %v, rid: %s,", f.key.Collection(),
			oids, err, rid)
		return true, err
	}

	if hitConflict {
		// update the chain nodes with picked chain nodes, so that we can handle them later.
		chainNodes = f.rearrangeEvents(chainNodes, rid)
	}

	retry, err := f.doInsertEvents(chainNodes, lastTokenData, rid)
	if err != nil {
		return retry, err
	}

	blog.Infof("insert watch event for %s success, oids: %v, rid: %s", f.key.Collection(), oids, rid)
	return false, nil
}
