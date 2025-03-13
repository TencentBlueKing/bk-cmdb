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
	"fmt"
	"time"

	synctypes "configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/stream/types"
)

// resTypeWatchOptMap is cmdb data sync resource type to db watch options map
var resTypeWatchOptMap = map[synctypes.ResType]*types.WatchCollOptions{
	synctypes.ServiceInstance: {
		CollectionOptions: types.CollectionOptions{
			CollectionFilter: &types.CollectionFilter{
				Regex: fmt.Sprintf("_%s$", common.BKTableNameServiceInstance),
			},
			EventStruct: new(metadata.ServiceInstance),
		},
	},
}

// watchDB watch db events for resource that are not watched by flow
func (w *Watcher) watchDB(resType synctypes.ResType) error {
	handler := w.tokenHandlers[resType]

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name:         string(resType),
			CollOpts:     resTypeWatchOptMap[resType],
			TokenHandler: handler,
			RetryOptions: &types.RetryOptions{
				MaxRetryCount: 3,
				RetryDuration: 1 * time.Second,
			},
		},
		EventHandler: &types.TaskBatchHandler{
			DoBatch: func(dbInfo *types.DBInfo, es []*types.Event) bool {
				return w.handleDBEvents(resType, es)
			},
		},
		BatchSize: common.BKMaxLimitSize,
	}

	err := w.task.AddLoopBatchTask(opts)
	if err != nil {
		blog.Errorf("watch %s events from db failed, err: %v", resType, err)
		return err
	}

	return nil
}

// handleDBEvents handle db events
func (w *Watcher) handleDBEvents(resType synctypes.ResType, es []*types.Event) bool {
	kit := rest.NewKit()

	eventInfos := make([]*synctypes.EventInfo, 0)
	for _, e := range es {
		eventType := watch.ConvertOperateType(e.OperationType)
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
