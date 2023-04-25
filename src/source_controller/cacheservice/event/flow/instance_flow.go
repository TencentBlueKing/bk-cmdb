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
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	types2 "configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/types"
	"configcenter/src/thirdparty/monitor"
	"configcenter/src/thirdparty/monitor/meta"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func newInstanceFlow(ctx context.Context, opts flowOptions, getDeleteEventDetails getDeleteEventDetailsFunc,
	parseEvent parseEventFunc) error {

	flow, err := NewFlow(opts, getDeleteEventDetails, parseEvent)
	if err != nil {
		return err
	}
	instFlow := InstanceFlow{
		Flow:              flow,
		mainlineObjectMap: new(mainlineObjectMap),
	}

	mainlineObjectMap, err := instFlow.getMainlineObjectMap(ctx)
	if err != nil {
		blog.Errorf("run object instance watch, but get mainline objects failed, err: %v", err)
		return err
	}
	instFlow.mainlineObjectMap.Set(mainlineObjectMap)

	go instFlow.syncMainlineObjectMap()

	return instFlow.RunFlow(ctx)
}

// syncMainlineObjectMap refresh mainline object ID map every 5 minutes
func (f *InstanceFlow) syncMainlineObjectMap() {
	for {
		time.Sleep(time.Minute * 5)

		mainlineObjectMap, err := f.getMainlineObjectMap(context.Background())
		if err != nil {
			blog.Errorf("run object instance watch, but get mainline objects failed, err: %v", err)
			continue
		}
		f.mainlineObjectMap.Set(mainlineObjectMap)
		blog.V(5).Infof("run object instance watch, sync mainline object map done, map: %+v", f.mainlineObjectMap.Get())
	}
}

type mapStrWithOid struct {
	Oid    primitive.ObjectID     `bson:"_id"`
	MapStr map[string]interface{} `bson:",inline"`
}

type mainlineObjectMap struct {
	data map[string]struct{}
	lock sync.RWMutex
}

// Get TODO
func (m *mainlineObjectMap) Get() map[string]struct{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	data := make(map[string]struct{})
	for key, value := range m.data {
		data[key] = value
	}
	return data
}

// Set TODO
func (m *mainlineObjectMap) Set(data map[string]struct{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data = data
}

// InstanceFlow TODO
type InstanceFlow struct {
	Flow
	*mainlineObjectMap
}

// RunFlow TODO
func (f *InstanceFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.watchDB, f.metrics)

	startAtTime, err := f.tokenHandler.getStartWatchTime(ctx)
	if err != nil {
		blog.Errorf("get start watch time for %s failed, err: %v", f.key.Collection(), err)
		return err
	}

	watchOpts := &types.WatchOptions{
		Options: types.Options{
			EventStruct: f.EventStruct,
			// watch all tables with the prefix of instance table
			CollectionFilter: map[string]interface{}{
				common.BKDBLIKE: event.ObjInstTablePrefixRegex,
			},
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

func (f *InstanceFlow) doBatch(es []*types.Event) (retry bool) {
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

	oidDetailMap, retry, err := f.getDeleteEventDetails(es, f.ccDB, f.metrics)
	if err != nil {
		blog.Errorf("get deleted event details failed, err: %v, rid: %s", err, rid)
		return retry
	}

	eventMap, oidIndexMap, aggregationEvents, err := f.classifyEvents(es, oidDetailMap, rid)
	if err != nil {
		blog.Errorf("get aggregation inst events failed, err: %v, rid: %s", err, rid)
		return false
	}

	eventLen := len(aggregationEvents)
	if eventLen == 0 {
		return false
	}

	ids, err := f.watchDB.NextSequences(context.Background(), f.key.Collection(), eventLen)
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.key.Collection(), err, rid)
		return true
	}

	pipe := redis.Client().Pipeline()
	oids := make([]string, 0)
	chainNodesMap := make(map[string][]*watch.ChainNode)
	lastChainNode := new(watch.ChainNode)
	for coll, events := range eventMap {
		key := f.getKeyByCollection(coll)
		cursorMap := make(map[string]struct{})
		hitConflict := false
		for _, e := range events {
			// collect event's basic metrics
			f.metrics.CollectBasic(e)

			idIndex := oidIndexMap[e.Oid+e.Collection]
			chainNode, detailBytes, retry, err := f.parseEvent(f.ccDB, key, e, oidDetailMap, ids[idIndex], rid)
			if err != nil {
				return retry
			}
			if chainNode == nil {
				continue
			}
			chainNode.SubResource = []string{gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()}

			if idIndex == eventLen-1 {
				lastChainNode = chainNode
			}

			// validate if the cursor is already exists, this is happens when the concurrent operation is very high.
			// which will generate the same operation event with same cluster time, and generate with the same cursor
			// in the end. if this happens, the last event will be used finally, and the former events with the same
			// cursor will be dropped, and it's acceptable.
			if _, exists := cursorMap[chainNode.Cursor]; exists {
				hitConflict = true
			}
			cursorMap[chainNode.Cursor] = struct{}{}

			oids = append(oids, e.ID())
			chainNodesMap[coll] = append(chainNodesMap[coll], chainNode)

			// if hit cursor conflict, the former cursor node's detail will be overwrite by the later one, so it
			// is not needed to remove the overlapped cursor node's detail again.
			pipe.Set(key.DetailKey(chainNode.Cursor), string(detailBytes), time.Duration(key.TTLSeconds())*time.Second)
		}

		if hitConflict {
			chainNodesMap[coll] = f.rearrangeEvents(chainNodesMap[coll], rid)
		}
	}

	lastTokenData := map[string]interface{}{
		common.BKTokenField:       aggregationEvents[eventLen-1].Token.Data,
		common.BKStartAtTimeField: aggregationEvents[eventLen-1].ClusterTime,
	}

	// if all events are invalid, set last token to the last events' token, do not need to retry for the invalid ones
	if len(chainNodesMap) == 0 {
		if err := f.tokenHandler.setLastWatchToken(context.Background(), lastTokenData); err != nil {
			f.metrics.CollectMongoError()
			return false
		}
		return false
	}

	lastTokenData[common.BKFieldID] = lastChainNode.ID
	lastTokenData[common.BKCursorField] = lastChainNode.Cursor
	lastTokenData[common.BKStartAtTimeField] = lastChainNode.ClusterTime

	// store details at first, in case those watching cmdb events read chain when details are not inserted yet
	if _, err := pipe.Exec(); err != nil {
		f.metrics.CollectRedisError()
		blog.Errorf("run flow, but insert details for %s failed, oids: %+v, err: %v, rid: %s,", f.key.Collection(),
			oids, err, rid)
		return true
	}

	retry, err = f.doInsertEvents(chainNodesMap, lastTokenData, rid)
	if err != nil {
		return retry
	}

	hasError = false
	return false
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
		if common.IsInnerModel(relation.ObjectID) {
			continue
		}
		objIDMap[relation.ObjectID] = struct{}{}
	}
	return objIDMap, nil
}

// classifyEvents classify events by their related key's collection
func (f *InstanceFlow) classifyEvents(es []*types.Event, oidDetailMap map[oidCollKey][]byte, rid string) (
	map[string][]*types.Event, map[string]int, []*types.Event, error) {

	mainlineColl := event.MainlineInstanceKey.Collection()
	commonColl := f.key.Collection()

	aggregationInstEvents, err := f.convertTableInstEvent(es, rid)
	if err != nil {
		blog.Errorf("get aggregation inst events failed, err: %v, rid: %s", err, rid)
		return nil, nil, nil, err
	}
	if len(aggregationInstEvents) == 0 {
		return nil, nil, nil, err
	}

	eventMap := make(map[string][]*types.Event)
	oidIndexMap := make(map[string]int)
	for index, e := range aggregationInstEvents {
		oidIndexMap[e.Oid+e.Collection] = index

		if e.OperationType == types.Delete {
			doc, exist := oidDetailMap[oidCollKey{oid: e.Oid, coll: e.Collection}]
			if !exist {
				blog.Errorf("run flow, received %s %s event, but delete doc[oid: %s] detail not exists, rid: %s",
					f.key.Collection(), e.OperationType, e.Oid, rid)
				continue
			}
			e.DocBytes = doc
		}

		objID := gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()
		if len(objID) == 0 {
			blog.Errorf("run flow, %s event[oid: %s] object id not exists, doc: %s, rid: %s",
				f.key.Collection(), e.Oid, string(e.DocBytes), rid)
			continue
		}

		if _, exists := f.mainlineObjectMap.Get()[objID]; exists {
			eventMap[mainlineColl] = append(eventMap[mainlineColl], e)
			continue
		}
		eventMap[commonColl] = append(eventMap[commonColl], e)
	}

	return eventMap, oidIndexMap, aggregationInstEvents, nil
}

// convertTableInstEvent convert table inst event to inst event
// 当收到所有的事件后，通过集合名称来获取objID，然后通过objID查询模型引用表cc_ModelQuoteRelation，判断当前事件是否为表格实例事件
// 取出表格实例事件中父模型的实例id，构造父模型objID->instID的map，然后去重实例id，查询源模型实例数据，构造源模型实例的更新事件
// 其中通过记录index->event的map，来保证表格实例事件聚合后，当前批次所有事件的时序性，表格实例事件聚合成一个源模型实例，使用最后的一个event
func (f *InstanceFlow) convertTableInstEvent(es []*types.Event, rid string) ([]*types.Event, error) {
	if len(es) == 0 {
		return es, nil
	}

	notContainTableInstEventsMap := make(map[int]*types.Event, 0)
	srcObjIDInstIDsMap := make(map[string][]int64, 0)
	instIDEventMap := make(map[int64]*types.Event)
	instIDIndexMap := make(map[int64]int, 0)
	for index, e := range es {
		objID, err := common.GetInstObjIDByTableName(e.Collection, gjson.Get(string(e.DocBytes),
			common.BkSupplierAccount).Str)
		if err != nil {
			blog.Errorf("collection name is illegal, err: %v, rid: %s", err, rid)
			return nil, err
		}

		modelQuoteRel := make([]metadata.ModelQuoteRelation, 0)
		queryCond := mapstr.MapStr{
			common.BKDestModelField: objID,
		}
		err = f.ccDB.Table(common.BKTableNameModelQuoteRelation).Find(queryCond).All(context.TODO(), &modelQuoteRel)
		if err != nil {
			blog.Errorf("get model quote relation failed, err: %v, rid: %s", err, rid)
			return nil, err
		}
		if len(modelQuoteRel) == 0 {
			notContainTableInstEventsMap[index] = e
			continue
		}
		if len(modelQuoteRel) != 1 {
			return nil, fmt.Errorf("model quote relation not unique, rel: %v", modelQuoteRel)
		}

		if modelQuoteRel[0].SrcModel == "" {
			return nil, fmt.Errorf("src model objID is illegal, rel: %v", modelQuoteRel)
		}

		if modelQuoteRel[0].PropertyID == "" {
			return nil, fmt.Errorf("table field property id is illegal, rel: %v", modelQuoteRel)
		}

		srcObjID := modelQuoteRel[0].SrcModel
		instID := gjson.Get(string(e.DocBytes), common.BKInstIDField).Int()
		srcObjIDInstIDsMap[srcObjID] = append(srcObjIDInstIDsMap[srcObjID], instID)

		instIDIndexMap[instID] = index
		instIDEventMap[instID] = e
	}

	return f.convertToInstEvents(notContainTableInstEventsMap, srcObjIDInstIDsMap, instIDEventMap, instIDIndexMap, rid)
}

func (f *InstanceFlow) convertToInstEvents(es map[int]*types.Event, srcObjIDInstIDsMap map[string][]int64,
	instIDEventMap map[int64]*types.Event, instIDIndexMap map[int64]int, rid string) ([]*types.Event, error) {
	for objID, instIDs := range srcObjIDInstIDsMap {
		if len(instIDs) == 0 {
			continue
		}

		tableName := common.GetInstTableName(objID, gjson.Get(string(instIDEventMap[instIDs[0]].DocBytes),
			common.BkSupplierAccount).Str)
		filter := mapstr.MapStr{
			common.GetInstIDField(objID): mapstr.MapStr{
				common.BKDBIN: util.IntArrayUnique(instIDs),
			},
		}
		findOpts := dbtypes.NewFindOpts().SetWithObjectID(true)
		insts := make([]mapStrWithOid, 0)
		err := f.ccDB.Table(tableName).Find(filter, findOpts).All(context.TODO(), &insts)
		if err != nil {
			blog.Errorf("get src model inst failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		for _, inst := range insts {
			doc, err := json.Marshal(inst.MapStr)
			if err != nil {
				blog.Errorf("marshal inst to byte failed, err: %v, rid: %s", err, rid)
				continue
			}

			instID, err := util.GetInt64ByInterface(inst.MapStr[common.GetInstIDField(objID)])
			if err != nil {
				blog.Errorf("get inst id failed, err: %v, rid: %s", err, rid)
				return nil, err
			}

			instEvent := &types.Event{
				Oid:           inst.Oid.Hex(),
				Document:      inst.MapStr,
				DocBytes:      doc,
				OperationType: "update",
				Collection:    tableName,
				ClusterTime: types.TimeStamp{
					Sec:  instIDEventMap[instID].ClusterTime.Sec,
					Nano: instIDEventMap[instID].ClusterTime.Nano,
				},
				Token: instIDEventMap[instID].Token,
				ChangeDesc: &types.ChangeDescription{
					UpdatedFields: make(map[string]interface{}, 0),
					RemovedFields: make([]string, 0),
				},
			}

			es[instIDIndexMap[instID]] = instEvent
		}
	}

	keys := make([]int, 0)
	for key := range es {
		keys = append(keys, key)
	}
	sort.Sort(sort.IntSlice(keys))

	aggregationInstEvents := make([]*types.Event, 0)
	for _, v := range keys {
		aggregationInstEvents = append(aggregationInstEvents, es[v])
	}
	return aggregationInstEvents, nil
}

func (f *InstanceFlow) doInsertEvents(chainNodesMap map[string][]*watch.ChainNode, lastTokenData map[string]interface{},
	rid string) (bool, error) {

	if len(chainNodesMap) == 0 {
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
	var conflictColl string

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			blog.Errorf("run flow, but start transaction failed, coll: %s, err: %v, rid: %s", f.key.Collection(),
				err, rid)
			return err
		}

		for coll, chainNodes := range chainNodesMap {
			if len(chainNodes) == 0 {
				continue
			}
			key := f.getKeyByCollection(coll)

			if err := f.watchDB.Table(key.ChainCollection()).Insert(sc, chainNodes); err != nil {
				blog.ErrorJSON("run flow, but insert chain nodes for %s failed, nodes: %s, err: %v, rid: %s",
					key.Collection(), chainNodes, err, rid)
				f.metrics.CollectMongoError()
				_ = session.AbortTransaction(context.Background())

				if event.IsConflictError(err) {
					// set retry with reduce flag and retry later
					retryWithReduce = true
					conflictColl = coll
				}
				return err
			}
		}

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

		if retryWithReduce {
			chainNodes := chainNodesMap[conflictColl]
			if len(chainNodes) == 0 {
				return false, nil
			}
			key := f.getKeyByCollection(conflictColl)

			rid = rid + ":" + chainNodes[0].Oid
			monitor.Collect(&meta.Alarm{
				RequestID: rid,
				Type:      meta.EventFatalError,
				Detail:    fmt.Sprintf("run event flow, but got conflict %s cursor with chain nodes", key.Collection()),
				Module:    types2.CC_MODULE_CACHESERVICE,
				Dimension: map[string]string{"retry_conflict_nodes": "yes"},
			})

			// no need to retry because the only one chain node conflicts with the nodes in db
			if len(chainNodes) <= 1 {
				return false, nil
			}

			for index, reducedChainNode := range chainNodes {
				if isConflictChainNode(reducedChainNode, txnErr) {
					chainNodes = append(chainNodes[:index], chainNodes[index+1:]...)

					chainNodesMap[conflictColl] = chainNodes

					// need do with retry with reduce
					blog.ErrorJSON("run flow, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
						key.Collection(), reducedChainNode, chainNodes, rid)

					return f.doInsertEvents(chainNodesMap, lastTokenData, rid)
				}
			}

			// when no cursor conflict node is found, discard the first node and try to insert the others
			blog.ErrorJSON("run flow, insert %s event with reduce node %s, remain nodes: %s, rid: %s",
				key.Collection(), chainNodes[0], chainNodes[1:], rid)

			chainNodesMap[conflictColl] = chainNodes[1:]
			return f.doInsertEvents(chainNodesMap, lastTokenData, rid)
		}

		// if an error occurred, roll back and re-watch again
		return true, err
	}

	return false, nil
}

func (f *InstanceFlow) getKeyByCollection(collection string) event.Key {
	switch collection {
	case event.MainlineInstanceKey.Collection():
		return event.MainlineInstanceKey
	default:
		return f.key
	}
}
