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

package watch

import (
	"context"
	"encoding/json"
	"time"

	synctypes "configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/transfer-service/sync/util"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

// resTypeWatchOptMap is cmdb data sync resource type to db watch options map
var resTypeWatchOptMap = map[synctypes.ResType]types.Options{
	synctypes.ServiceInstance: {
		EventStruct: new(metadata.ServiceInstance),
		Collection:  common.BKTableNameServiceInstance,
	},
}

// watchDB watch db events for resource that are not watched by flow
func (w *Watcher) watchDB(resType synctypes.ResType) error {
	handler := w.tokenHandlers[resType]

	token, err := handler.getWatchTokenInfo(context.Background(), common.BKStartAtTimeField)
	if err != nil {
		blog.Errorf("get %s watch db token info failed, err: %v", resType, err)
		return err
	}

	startAtTime := &types.TimeStamp{Sec: uint32(time.Now().Unix())}
	if token.StartAtTime != nil {
		startAtTime = &types.TimeStamp{
			Sec:  uint32(token.StartAtTime.Unix()),
			Nano: uint32(token.StartAtTime.Nanosecond()),
		}
	}

	watchOpts := resTypeWatchOptMap[resType]
	watchOpts.StartAtTime = startAtTime
	watchOpts.WatchFatalErrorCallback = handler.resetWatchToken

	opts := &types.LoopBatchOptions{
		LoopOptions: types.LoopOptions{
			Name: string(resType),
			WatchOpt: &types.WatchOptions{
				Options: watchOpts,
			},
			TokenHandler: handler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 3,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.BatchHandler{
			DoBatch: func(es []*types.Event) (retry bool) {
				return w.handleDBEvents(resType, watchOpts.Collection, es)
			},
		},
		BatchSize: common.BKMaxLimitSize,
	}

	if err = w.loopW.WithBatch(opts); err != nil {
		blog.Errorf("watch %s events from db failed, err: %v", resType, err)
		return err
	}

	return nil
}

// handleDBEvents handle db events
func (w *Watcher) handleDBEvents(resType synctypes.ResType, coll string, es []*types.Event) (retry bool) {
	kit := util.NewKit()

	// get deleted event oid to detail map
	delOids := make([]string, 0)
	for _, e := range es {
		if e.OperationType == types.Delete {
			delOids = append(delOids, e.Oid)
		}
	}

	delOidMap := make(map[string]json.RawMessage)
	if len(delOids) > 0 {
		cond := mapstr.MapStr{
			"oid":  mapstr.MapStr{common.BKDBIN: delOids},
			"coll": coll,
		}
		archives := make([]delArchiveInfo, 0)
		err := mongodb.Client().Table(common.BKTableNameDelArchive).Find(cond).All(kit.Ctx, &archives)
		if err != nil {
			blog.Errorf("get del archive failed, err: %v, cond: %+v, rid: %s", err, cond, kit.Rid)
			return true
		}

		for _, archive := range archives {
			if archive.Detail == nil {
				continue
			}

			detail, err := json.Marshal(archive.Detail)
			if err != nil {
				blog.Errorf("marshal del archive detail failed, err: %v, archive: %+v, rid: %s", err, archive, kit.Rid)
				return true
			}
			delOidMap[archive.Oid] = detail
		}
	}

	eventInfos := make([]*synctypes.EventInfo, 0)
	for _, e := range es {
		eventType := watch.ConvertOperateType(e.OperationType)
		if eventType == watch.Delete {
			delDetail, exists := delOidMap[e.Oid]
			if !exists {
				continue
			}
			e.DocBytes = delDetail
		}

		eventInfo, needSync := w.metadata.ParseEventDetail(eventType, resType, e.Oid, e.DocBytes)
		if !needSync {
			continue
		}
		eventInfos = append(eventInfos, eventInfo)
	}

	// push incremental sync data to transfer medium
	err := w.pushSyncData(kit, eventInfos)
	if err != nil {
		return true
	}

	return false
}

type delArchiveInfo struct {
	Oid    string        `bson:"oid"`
	Detail mapstr.MapStr `bson:"detail"`
}
