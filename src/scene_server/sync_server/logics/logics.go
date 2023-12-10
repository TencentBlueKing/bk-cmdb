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

// Package logics defines the synchronization logics
package logics

import (
	"configcenter/src/common/backbone"
	"configcenter/src/scene_server/sync_server/logics/full-text-search"
	"configcenter/src/storage/stream"
)

// Logics defines the struct that contains all sync logics
type Logics struct {
	FullTextSearch fulltextsearch.SyncI
}

// New Logics instance
func New(engine *backbone.Engine, conf *Config, watcher stream.LoopInterface) (*Logics, error) {
	lgc := new(Logics)
	var err error

	lgc.FullTextSearch, err = fulltextsearch.New(conf.FullTextSearch, engine.CoreAPI.CacheService().Cache(), watcher)
	if err != nil {
		return nil, err
	}

	return lgc, nil
}

// Config defines synchronization logics configuration
type Config struct {
	FullTextSearch *fulltextsearch.Config `mapstructure:"fullTextSearch"`
}
