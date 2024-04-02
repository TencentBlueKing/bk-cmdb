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

// Package cache defines the custom resource redis cache logics
package cache

import (
	"configcenter/src/apimachinery/discovery"
)

// CacheSet is the set of custom resource caches
type CacheSet struct {
	Label       *PodLabelCache
	SharedNsRel *SharedNsRelCache
}

// New CacheSet
func New(isMaster discovery.ServiceManageInterface) *CacheSet {
	return &CacheSet{
		Label:       NewPodLabelCache(isMaster),
		SharedNsRel: NewSharedNsRelCache(isMaster),
	}
}

// LoopRefreshCache loop refresh all caches
func (c *CacheSet) LoopRefreshCache() {
	go c.Label.loopRefreshCache()
	go c.SharedNsRel.loopRefreshCache()
}
