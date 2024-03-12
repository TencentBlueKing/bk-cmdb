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

// Package custom defines the custom resource caching logics
package custom

import (
	"fmt"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	"configcenter/src/source_controller/cacheservice/cache/custom/watch"
	"configcenter/src/storage/stream"
)

// Cache defines the custom resource caching logics
type Cache struct {
	cacheSet *cache.CacheSet
}

// New Cache
func New(isMaster discovery.ServiceManageInterface, loopW stream.LoopInterface) (*Cache, error) {
	t := &Cache{
		cacheSet: cache.New(isMaster),
	}

	if err := watch.Init(loopW, t.cacheSet); err != nil {
		return nil, fmt.Errorf("initialize custom resource watcher failed, err: %v", err)
	}

	t.cacheSet.LoopRefreshCache()
	return t, nil
}
