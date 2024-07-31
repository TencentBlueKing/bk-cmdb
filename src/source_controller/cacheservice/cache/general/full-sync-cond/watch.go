/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package fullsynccond

import (
	"context"
	"fmt"
	"strconv"
	"time"

	fullsynccond "configcenter/pkg/cache/full-sync-cond"
	"configcenter/pkg/cache/general"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	cachetypes "configcenter/src/source_controller/cacheservice/cache/general/types"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

// Watch full sync cond, send the event to the corresponding channel
func (f *FullSyncCond) Watch() error {
	tokenHandler := tokenhandler.NewMemoryTokenHandler()

	startAtTime := &types.TimeStamp{Sec: uint32(time.Now().Unix())}

	if err := f.initFullSyncCond(); err != nil {
		return err
	}

	loopOptions := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name: "full-sync-cond",
			WatchOpt: &types.WatchOptions{
				Options: types.Options{
					Filter:      make(mapstr.MapStr),
					EventStruct: new(fullsynccond.FullSyncCond),
					Collection:  fullsynccond.BKTableNameFullSyncCond,
					StartAtTime: startAtTime,
				},
			},
			TokenHandler: tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 5,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: 200,
	}

	if err := f.loopW.WithBatch(loopOptions); err != nil {
		blog.Errorf("watch full sync cond failed, err: %v", err)
		return err
	}

	return nil
}

// doBatch batch handle full sync cond event
func (f *FullSyncCond) doBatch(es []*types.Event) (retry bool) {
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)

	// get deleted full sync cond oid to info map
	delOids := make([]string, 0)
	for _, e := range es {
		if e.OperationType == types.Delete {
			delOids = append(delOids, e.Oid)
		}
	}

	delOidCondMap := make(map[string]*fullsynccond.FullSyncCond)
	if len(delOids) > 0 {
		filter := mapstr.MapStr{
			"oid":  mapstr.MapStr{common.BKDBIN: delOids},
			"coll": fullsynccond.BKTableNameFullSyncCond,
		}
		archives := make([]delArchive, 0)
		err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(filter).All(ctx, &archives)
		if err != nil {
			blog.Errorf("get deleted full sync cond failed, err: %v, oids: %+v", err, delOids)
			return true
		}

		for _, archive := range archives {
			if archive.Detail == nil {
				continue
			}
			delOidCondMap[archive.Oid] = archive.Detail
		}
	}

	// aggregate full sync cond event
	condMap := f.aggregateEvent(es, delOidCondMap)

	// generate full sync cond event
	resEventMap := make(map[general.ResType]map[cachetypes.EventType][]*fullsynccond.FullSyncCond)

	for op, condInfo := range condMap {
		var eventType cachetypes.EventType
		switch op {
		case types.Insert, types.Update:
			eventType = cachetypes.Upsert
		case types.Delete:
			eventType = cachetypes.Delete
		}
		for _, cond := range condInfo {
			eventMap, exists := resEventMap[cond.Resource]
			if !exists {
				eventMap = make(map[cachetypes.EventType][]*fullsynccond.FullSyncCond)
			}
			eventMap[eventType] = append(eventMap[eventType], cond)
			resEventMap[cond.Resource] = eventMap
		}
	}

	// add full sync cond event to channel
	for resType, eventMap := range resEventMap {
		ch, exists := f.chMap[resType]
		if !exists {
			blog.Infof("%s resource type is invalid, events: %+v", resType, eventMap)
			continue
		}

		ch <- cachetypes.FullSyncCondEvent{EventMap: eventMap}
	}

	return false
}

// aggregateEvent aggregate full sync cond event
func (f *FullSyncCond) aggregateEvent(es []*types.Event,
	delOidCondMap map[string]*fullsynccond.FullSyncCond) map[types.OperType]map[string]*fullsynccond.FullSyncCond {

	condMap := make(map[types.OperType]map[string]*fullsynccond.FullSyncCond)
	supportedOps := []types.OperType{types.Insert, types.Update, types.Delete}
	for _, op := range supportedOps {
		condMap[op] = make(map[string]*fullsynccond.FullSyncCond)
	}

	for _, e := range es {
		switch e.OperationType {
		case types.Insert, types.Update:
			cond, ok := e.Document.(*fullsynccond.FullSyncCond)
			if !ok {
				blog.Errorf("event document %+v type is invalid", e.Document)
				continue
			}

			condKey := genFullSyncCondUniqueKey(cond)

			// if full sync cond with same unique key is deleted before, treat these event as an upsert event
			_, exists := condMap[types.Delete][condKey]
			if exists {
				delete(condMap[types.Delete], condKey)
				condMap[types.Update][condKey] = cond
				continue
			}

			// if full sync cond with same unique key is inserted before, treat these event as an insert event
			_, exists = condMap[types.Insert][condKey]
			if exists {
				condMap[types.Insert][condKey] = cond
				continue
			}

			condMap[e.OperationType][condKey] = cond
		case types.Delete:
			cond, exists := delOidCondMap[e.Oid]
			if !exists {
				blog.Errorf("delete event %s has no matching del archive", e.Oid)
				continue
			}

			condKey := genFullSyncCondUniqueKey(cond)

			// if full sync cond with same unique key is inserted before, treat these event as not exists
			_, exists = condMap[types.Insert][condKey]
			if exists {
				delete(condMap[types.Insert], condKey)
				continue
			}

			// if full sync cond with same unique key is updated before, treat these event as delete event
			_, exists = condMap[types.Update][condKey]
			if exists {
				delete(condMap[types.Update], condKey)
				condMap[types.Delete][condKey] = cond
				continue
			}

			condMap[types.Delete][condKey] = cond
		}
	}

	return condMap
}

func genFullSyncCondUniqueKey(cond *fullsynccond.FullSyncCond) string {
	if !cond.IsAll {
		return strconv.FormatInt(cond.ID, 10)
	}

	return fmt.Sprintf("%s:%s:%s", cond.Resource, cond.SupplierAccount, cond.SubResource)
}

type delArchive struct {
	Oid    string                     `bson:"oid"`
	Detail *fullsynccond.FullSyncCond `bson:"detail"`
}

// initFullSyncCond get all full sync cond from db and send initialize event to channel
func (f *FullSyncCond) initFullSyncCond() error {
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	rid := util.GenerateRID()

	fullSyncCondMap := make(map[general.ResType][]*fullsynccond.FullSyncCond)

	cond := make(mapstr.MapStr)

	for {
		paged := make([]*fullsynccond.FullSyncCond, 0)
		err := mongodb.Client().Table(fullsynccond.BKTableNameFullSyncCond).Find(cond).Sort(fullsynccond.IDField).
			Limit(cachetypes.PageSize).All(ctx, &paged)
		if err != nil {
			blog.Errorf("paged get full sync cond data failed, cond: %+v, err: %v, rid: %s", cond, err, rid)
			return err
		}

		for _, data := range paged {
			fullSyncCondMap[data.Resource] = append(fullSyncCondMap[data.Resource], data)
		}

		if len(paged) < cachetypes.PageSize {
			break
		}

		cond[fullsynccond.IDField] = mapstr.MapStr{common.BKDBGT: paged[len(paged)-1].ID}
	}

	// send init full sync cond event to channel
	for resType, fullSyncConds := range fullSyncCondMap {
		ch, exists := f.chMap[resType]
		if !exists {
			blog.Infof("%s resource type is invalid, conds: %+v", resType, fullSyncConds)
			continue
		}

		ch <- cachetypes.FullSyncCondEvent{
			EventMap: map[cachetypes.EventType][]*fullsynccond.FullSyncCond{cachetypes.Init: fullSyncConds},
		}
	}

	return nil
}
