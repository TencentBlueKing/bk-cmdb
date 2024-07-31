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

// Package general defines the general resource caching logics
package general

import (
	"fmt"

	"configcenter/pkg/cache/general"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/source_controller/cacheservice/cache/general/cache"
	fullsynccond "configcenter/src/source_controller/cacheservice/cache/general/full-sync-cond"
	"configcenter/src/source_controller/cacheservice/cache/general/types"
	"configcenter/src/source_controller/cacheservice/cache/general/watch"
	watchcli "configcenter/src/source_controller/cacheservice/event/watch"
	"configcenter/src/storage/stream"
)

// Cache defines the general resource caching logics
type Cache struct {
	fullSyncCond *fullsynccond.FullSyncCond
	cacheSet     map[general.ResType]*cache.Cache
}

// New Cache
func New(isMaster discovery.ServiceManageInterface, loopW stream.LoopInterface, watchCli *watchcli.Client) (*Cache,
	error) {

	cacheSet := cache.GetAllCache()

	fullSyncCondChMap := make(map[general.ResType]chan<- types.FullSyncCondEvent)
	for resType, cacheInst := range cacheSet {
		if err := watch.Init(cacheInst, isMaster, watchCli); err != nil {
			return nil, fmt.Errorf("init %s general resource watcher failed, err: %v", cacheInst.Key().Resource(), err)
		}
		fullSyncCondChMap[resType] = cacheInst.FullSyncCondCh()
	}

	fullSyncCondCli, err := fullsynccond.New(loopW, fullSyncCondChMap)
	if err != nil {
		return nil, fmt.Errorf("init full sync cond failed, err: %v", err)
	}

	return &Cache{
		fullSyncCond: fullSyncCondCli,
		cacheSet:     cacheSet,
	}, nil
}
