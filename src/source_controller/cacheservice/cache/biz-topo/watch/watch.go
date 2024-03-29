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

// Package watch defines the biz topology cache data watch logics
package watch

import (
	"configcenter/src/source_controller/cacheservice/cache/custom/cache"
	"configcenter/src/storage/stream"
)

// Watcher defines mongodb event watcher for biz topology
type Watcher struct {
	loopW    stream.LoopInterface
	cacheSet *cache.CacheSet
}

// New  biz topology mongodb event watcher
func New(loopW stream.LoopInterface, cacheSet *cache.CacheSet) (*Watcher, error) {
	watcher := &Watcher{
		loopW:    loopW,
		cacheSet: cacheSet,
	}

	if err := watcher.watchKube(); err != nil {
		return nil, err
	}

	return watcher, nil
}
