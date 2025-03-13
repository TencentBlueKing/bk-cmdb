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
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

// Watch full sync cond, send the event to the corresponding channel
func (f *FullSyncCond) Watch() error {
	tokenHandler := tokenhandler.NewMemoryTokenHandler()

	if err := f.initFullSyncCond(); err != nil {
		return err
	}

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name: "full-sync-cond",
			CollOpts: &types.WatchCollOptions{
				CollectionOptions: types.CollectionOptions{
					CollectionFilter: &types.CollectionFilter{
						Regex: fullsynccond.BKTableNameFullSyncCond,
					},
					EventStruct: new(fullsynccond.FullSyncCond),
				},
			},
			TokenHandler: tokenHandler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 5,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.TaskBatchHandler{
			DoBatch: f.doBatch,
		},
		BatchSize: 200,
	}

	err := f.task.AddLoopBatchTask(opts)
	if err != nil {
		blog.Errorf("add watch full sync cond task failed, err: %v", err)
		return err
	}

	return nil
}

// doBatch batch handle full sync cond event
func (f *FullSyncCond) doBatch(dbInfo *types.DBInfo, es []*types.Event) bool {
	// aggregate full sync cond event
	condMap := f.aggregateEvent(es)

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
func (f *FullSyncCond) aggregateEvent(es []*types.Event) map[types.OperType]map[string]*fullsynccond.FullSyncCond {
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
			cond, ok := e.Document.(*fullsynccond.FullSyncCond)
			if !ok {
				blog.Errorf("event document %+v type is invalid", e.Document)
				continue
			}

			condKey := genFullSyncCondUniqueKey(cond)

			// if full sync cond with same unique key is inserted before, treat these event as not exists
			_, exists := condMap[types.Insert][condKey]
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

	return fmt.Sprintf("%s:%s:%s", cond.Resource, cond.TenantID, cond.SubResource)
}

// initFullSyncCond get all full sync cond from db and send initialize event to channel
func (f *FullSyncCond) initFullSyncCond() error {
	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	rid := util.GenerateRID()

	fullSyncCondMap := make(map[general.ResType][]*fullsynccond.FullSyncCond)

	err := mongodb.Dal().ExecForAllDB(func(db local.DB) error {
		cond := make(mapstr.MapStr)

		for {
			paged := make([]*fullsynccond.FullSyncCond, 0)
			err := db.Table(fullsynccond.BKTableNameFullSyncCond).Find(cond).Sort(fullsynccond.IDField).
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

		return nil
	})
	if err != nil {
		return err
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
