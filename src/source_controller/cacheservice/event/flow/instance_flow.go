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

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/storage/dal/mongo/sharding"
	dbtypes "configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream/types"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func newInstanceFlow(ctx context.Context, opts flowOptions, parseEvent parseEventFunc) error {
	flow, err := NewFlow(opts, parseEvent)
	if err != nil {
		return err
	}
	instFlow := InstanceFlow{
		Flow: flow,
		mainlineObjectMap: &mainlineObjectMap{
			data: make(map[string]map[string]struct{}),
		},
	}

	err = tenant.ExecForAllTenants(func(tenantID string) error {
		mainlineObjMap, err := instFlow.getMainlineObjectMap(ctx, tenantID)
		if err != nil {
			blog.Errorf("run object instance watch, but get tenant %s mainline objects failed, err: %v", tenantID, err)
			return err
		}
		instFlow.mainlineObjectMap.Set(tenantID, mainlineObjMap)

		go instFlow.syncMainlineObjectMap(tenantID)
		return nil
	})
	if err != nil {
		return err
	}

	return instFlow.RunFlow(ctx)
}

// syncMainlineObjectMap refresh mainline object ID map every 5 minutes
func (f *InstanceFlow) syncMainlineObjectMap(tenantID string) {
	for {
		time.Sleep(time.Minute * 5)

		mainlineObjMap, err := f.getMainlineObjectMap(context.Background(), tenantID)
		if err != nil {
			blog.Errorf("run object instance watch, but get tenant %s mainline objects failed, err: %v", tenantID, err)
			continue
		}
		f.mainlineObjectMap.Set(tenantID, mainlineObjMap)

		blog.V(5).Infof("sync tenant %s mainline obj map done, map: %+v", tenantID, mainlineObjMap)
	}
}

type mapStrWithOid struct {
	Oid    primitive.ObjectID     `bson:"_id"`
	MapStr map[string]interface{} `bson:",inline"`
}

type mainlineObjectMap struct {
	data map[string]map[string]struct{}
	lock sync.RWMutex
}

// Get mainline object ID map for db
func (m *mainlineObjectMap) Get(tenantID string) map[string]struct{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	mainlineMap, exists := m.data[tenantID]
	if !exists {
		return make(map[string]struct{})
	}
	return mainlineMap
}

// Set mainline object ID map for db
func (m *mainlineObjectMap) Set(tenantID string, data map[string]struct{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data[tenantID] = data
}

// InstanceFlow TODO
type InstanceFlow struct {
	Flow
	*mainlineObjectMap
}

// RunFlow TODO
func (f *InstanceFlow) RunFlow(ctx context.Context) error {
	blog.Infof("start run flow for key: %s.", f.key.Namespace())

	f.tokenHandler = NewFlowTokenHandler(f.key, f.metrics)

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name: f.key.Namespace(),
			CollOpts: &types.WatchCollOptions{
				CollectionOptions: types.CollectionOptions{
					CollectionFilter: &types.CollectionFilter{
						Regex: fmt.Sprintf("_%s", common.BKObjectInstShardingTablePrefix),
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

func (f *InstanceFlow) doBatch(dbInfo *types.DBInfo, es []*types.Event) (retry bool) {
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

	eventMap, oidIndexMap, aggregationEvents, err := f.classifyEvents(es, rid)
	if err != nil {
		blog.Errorf("get aggregation inst events failed, err: %v, rid: %s", err, rid)
		return false
	}

	eventLen := len(aggregationEvents)
	if eventLen == 0 {
		return false
	}

	ids, err := dbInfo.WatchDB.NextSequences(context.Background(), f.key.Collection(), eventLen)
	if err != nil {
		blog.Errorf("get %s event ids failed, err: %v, rid: %s", f.key.Collection(), err, rid)
		return true
	}

	pipe := redis.Client().Pipeline()
	oids := make([]string, 0)
	chainNodesMap := make(map[string]map[string][]*watch.ChainNode)
	for coll, events := range eventMap {
		key := f.getKeyByCollection(coll)
		cursorMap := make(map[string]struct{})
		hitConflict := false
		for _, e := range events {
			// collect event's basic metrics
			f.metrics.CollectBasic(e)

			idIdx := oidIndexMap[e.Oid+e.Collection]
			tenantID, chainNode, detail, retry, err := f.parseEvent(dbInfo.CcDB, key, e, ids[idIdx], rid)
			if err != nil {
				return retry
			}
			if chainNode == nil {
				continue
			}
			chainNode.SubResource = []string{gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()}

			// validate if the cursor is already exists, this is happens when the concurrent operation is very high.
			// which will generate the same operation event with same cluster time, and generate with the same cursor
			// in the end. if this happens, the last event will be used finally, and the former events with the same
			// cursor will be dropped, and it's acceptable.
			if _, exists := cursorMap[chainNode.Cursor]; exists {
				hitConflict = true
			}
			cursorMap[chainNode.Cursor] = struct{}{}

			oids = append(oids, e.ID())
			_, exists := chainNodesMap[coll]
			if !exists {
				chainNodesMap[coll] = make(map[string][]*watch.ChainNode)
			}
			chainNodesMap[coll][tenantID] = append(chainNodesMap[coll][tenantID], chainNode)

			// if hit cursor conflict, the former cursor node's detail will be overwrite by the later one, so it
			// is not needed to remove the overlapped cursor node's detail again.
			ttl := time.Duration(key.TTLSeconds()) * time.Second
			pipe.Set(key.DetailKey(tenantID, chainNode.Cursor), string(detail.eventInfo), ttl)
			pipe.Set(key.GeneralResDetailKey(tenantID, chainNode), string(detail.resDetail), ttl)
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

	retry, err = f.doInsertEvents(dbInfo, chainNodesMap, lastTokenData, rid)
	if err != nil {
		return retry
	}

	hasError = false
	return false
}

func (f *InstanceFlow) getMainlineObjectMap(ctx context.Context, tenantID string) (map[string]struct{}, error) {
	relations := make([]metadata.Association, 0)
	filter := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameObjAsst).Find(filter).
		Fields(common.BKObjIDField).All(ctx, &relations)
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
func (f *InstanceFlow) classifyEvents(es []*types.Event, rid string) (map[string][]*types.Event, map[string]int,
	[]*types.Event, error) {

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

		objID := gjson.GetBytes(e.DocBytes, common.BKObjIDField).String()
		if len(objID) == 0 {
			blog.Errorf("run flow, %s event[oid: %s] object id not exists, doc: %s, rid: %s",
				f.key.Collection(), e.Oid, string(e.DocBytes), rid)
			continue
		}

		if _, exists := f.mainlineObjectMap.Get(e.TenantID)[objID]; exists {
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

	tenantObjIDInstIDsMap := make(map[string]map[string][]int64)
	tenantObjIDsMap := make(map[string][]string)
	instIDEventMap := make(map[int64]*types.Event)
	instIDIndexMap := make(map[int64]int)
	for index, e := range es {
		objID, err := common.GetInstObjIDByTableName(e.ParsedColl, e.TenantID)
		if err != nil {
			blog.Errorf("collection name is illegal, err: %v, rid: %s", err, rid)
			return nil, err
		}
		tenantObjIDsMap[e.TenantID] = append(tenantObjIDsMap[e.TenantID], objID)

		instID := gjson.Get(string(e.DocBytes), common.BKInstIDField).Int()

		_, exists := tenantObjIDInstIDsMap[e.TenantID]
		if !exists {
			tenantObjIDInstIDsMap[e.TenantID] = make(map[string][]int64)
		}
		tenantObjIDInstIDsMap[e.TenantID][objID] = append(tenantObjIDInstIDsMap[e.TenantID][objID], instID)

		instIDIndexMap[instID] = index
		instIDEventMap[instID] = e
	}

	notContainTableInstEventsMap := make(map[int]*types.Event)
	srcObjIDInstIDsMap := make(map[string]map[string][]int64)
	for tenantID, objIDs := range tenantObjIDsMap {
		modelQuoteRel := make([]metadata.ModelQuoteRelation, 0)
		queryCond := mapstr.MapStr{
			common.BKDestModelField: mapstr.MapStr{common.BKDBIN: objIDs},
		}
		err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(common.BKTableNameModelQuoteRelation).
			Find(queryCond).All(context.TODO(), &modelQuoteRel)
		if err != nil {
			blog.Errorf("get model quote relation failed, err: %v, rid: %s", err, rid)
			return nil, err
		}

		objSrcObjIDMap := make(map[string]string)
		for _, rel := range modelQuoteRel {
			if rel.SrcModel == "" {
				return nil, fmt.Errorf("src model objID is illegal, rel: %v", modelQuoteRel)
			}
			if rel.PropertyID == "" {
				return nil, fmt.Errorf("table field property id is illegal, rel: %v", modelQuoteRel)
			}
			objSrcObjIDMap[rel.DestModel] = rel.SrcModel
		}

		for _, objID := range objIDs {
			srcObjID, exists := objSrcObjIDMap[objID]
			if !exists {
				for _, instID := range tenantObjIDInstIDsMap[tenantID][objID] {
					notContainTableInstEventsMap[instIDIndexMap[instID]] = instIDEventMap[instID]
				}
				continue
			}
			srcObjIDInstIDsMap[tenantID][srcObjID] = append(srcObjIDInstIDsMap[tenantID][srcObjID],
				tenantObjIDInstIDsMap[tenantID][objID]...)
		}
	}

	return f.convertToInstEvents(notContainTableInstEventsMap, srcObjIDInstIDsMap, instIDEventMap, instIDIndexMap, rid)
}

func (f *InstanceFlow) convertToInstEvents(es map[int]*types.Event, srcObjIDInstIDsMap map[string]map[string][]int64,
	instIDEventMap map[int64]*types.Event, instIDIndexMap map[int64]int, rid string) ([]*types.Event, error) {

	for tenantID, objInstIDMap := range srcObjIDInstIDsMap {
		for objID, instIDs := range objInstIDMap {
			if len(instIDs) == 0 {
				continue
			}

			tableName := common.GetInstTableName(objID, tenantID)
			filter := mapstr.MapStr{
				common.GetInstIDField(objID): mapstr.MapStr{
					common.BKDBIN: util.IntArrayUnique(instIDs),
				},
			}
			findOpts := dbtypes.NewFindOpts().SetWithObjectID(true)
			insts := make([]mapStrWithOid, 0)
			err := mongodb.Shard(sharding.NewShardOpts().WithTenant(tenantID)).Table(tableName).Find(filter, findOpts).
				All(context.TODO(), &insts)
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
					CollectionInfo: types.CollectionInfo{
						Collection: common.GenTenantTableName(tenantID, tableName),
						ParsedColl: tableName,
						TenantID:   tenantID,
					},
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

func (f *InstanceFlow) doInsertEvents(dbInfo *types.DBInfo, chainNodesMap map[string]map[string][]*watch.ChainNode,
	lastTokenData map[string]interface{}, rid string) (bool, error) {

	if len(chainNodesMap) == 0 {
		return false, nil
	}

	session, err := dbInfo.WatchDB.GetDBClient().StartSession()
	if err != nil {
		blog.Errorf("run flow, but start session failed, coll: %s, err: %v, rid: %s", f.key.Collection(), err, rid)
		return true, err
	}
	defer session.EndSession(context.Background())

	// retry insert the event node with remove the first event node,
	// which means the first one's cursor is conflicted with the former's batch operation inserted nodes.
	retryWithReduce := false
	var conflictColl, conflictTenantID string

	txnErr := mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		if err = session.StartTransaction(); err != nil {
			blog.Errorf("run flow, but start transaction failed, coll: %s, err: %v, rid: %s", f.key.Collection(),
				err, rid)
			return err
		}

		for coll, chainNodeInfo := range chainNodesMap {
			key := f.getKeyByCollection(coll)
			err, retryWithReduce, conflictTenantID = f.insertChainNodes(sc, session, key, chainNodeInfo, rid)
			if err != nil {
				if retryWithReduce {
					conflictColl = coll
				}
				return err
			}
		}

		if err := f.tokenHandler.setLastWatchToken(sc, dbInfo.UUID, dbInfo.WatchDB, lastTokenData); err != nil {
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
			chainNodesMap[conflictColl] = event.ReduceChainNode(chainNodesMap[conflictColl], conflictTenantID,
				f.getKeyByCollection(conflictColl).Collection(), txnErr, f.metrics, rid)
			if len(chainNodesMap[conflictColl]) == 0 {
				delete(chainNodesMap, conflictColl)
			}
			if len(chainNodesMap) == 0 {
				return false, nil
			}
			return f.doInsertEvents(dbInfo, chainNodesMap, lastTokenData, rid)
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
