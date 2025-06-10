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
	"strings"
	"time"

	synctypes "configcenter/pkg/synchronize/types"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/watch"
	"configcenter/src/storage/stream/task"
	"configcenter/src/storage/stream/types"
)

// watchDB watch db events for resource that are not watched by flow
func (w *Watcher) watchDB(resType synctypes.ResType) (*task.Task, error) {
	collOpts, err := w.genWatchCollOptions(resType)
	if err != nil {
		return nil, err
	}
	handler := w.tokenHandlers[resType]

	opts := &types.LoopBatchTaskOptions{
		WatchTaskOptions: &types.WatchTaskOptions{
			Name:         string(resType),
			CollOpts:     collOpts,
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

	return task.NewLoopBatchTask(opts)
}

func (w *Watcher) genWatchCollOptions(resType synctypes.ResType) (*types.WatchCollOptions, error) {
	switch resType {
	case synctypes.ServiceInstance:
		collections := make([]string, 0)
		for tenantID := range w.tenantMap {
			collections = append(collections, common.GenTenantTableName(tenantID, common.BKTableNameServiceInstance))
		}
		return &types.WatchCollOptions{
			CollectionOptions: types.CollectionOptions{
				CollectionFilter: &types.CollectionFilter{
					Regex: strings.Join(collections, "|"),
				},
				EventStruct: new(metadata.ServiceInstance),
			},
		}, nil
	}
	return nil, fmt.Errorf("not supported watch db resType: %s", resType)
}

// handleDBEvents handle db events
func (w *Watcher) handleDBEvents(resType synctypes.ResType, es []*types.Event) bool {
	kit := rest.NewKit()

	tenantEventInfos := make(map[string][]*synctypes.EventInfo, 0)
	for _, e := range es {
		kit = kit.WithTenant(e.TenantID)
		eventType := watch.ConvertOperateType(e.OperationType)
		eventInfo, needSync := w.metadata.ParseEventDetail(kit, eventType, resType, e.Oid, e.DocBytes)
		if !needSync {
			continue
		}
		tenantEventInfos[e.TenantID] = append(tenantEventInfos[e.TenantID], eventInfo)
	}

	// push incremental sync data to transfer medium
	for tenantID, eventInfos := range tenantEventInfos {
		kit = kit.WithTenant(tenantID)
		err := w.pushSyncData(kit, eventInfos)
		if err != nil {
			return true
		}
	}

	return false
}
