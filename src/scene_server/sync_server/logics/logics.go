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
	"configcenter/src/common/blog"
	"configcenter/src/scene_server/sync_server/logics/full-text-search"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/scheduler"
)

// Logics defines the struct that contains all sync logics
type Logics struct {
	FullTextSearch fulltextsearch.SyncI
	scheduler      *scheduler.Scheduler
}

// New Logics instance
func New(engine *backbone.Engine, conf *Config) (*Logics, error) {
	lgc := new(Logics)

	var err error
	lgc.FullTextSearch, err = fulltextsearch.New(conf.FullTextSearch, engine.CoreAPI.CacheService().Cache())
	if err != nil {
		blog.Errorf("new full text search logics failed, err: %v", err)
		return nil, err
	}

	watchTasks := lgc.FullTextSearch.GetWatchTasks()
	if len(watchTasks) > 0 {
		lgc.scheduler, err = scheduler.New(mongodb.Dal(), mongodb.Dal("watch"), engine.ServiceManageInterface)
		if err != nil {
			blog.Errorf("new watch task scheduler failed, err: %v", err)
			return nil, err
		}

		if err = lgc.scheduler.AddTasks(watchTasks...); err != nil {
			blog.Errorf("add event watch tasks failed, err: %v", err)
			return nil, err
		}

		if err = lgc.scheduler.Start(); err != nil {
			blog.Errorf("start event watch task scheduler failed, err: %v", err)
			return nil, err
		}
	}

	return lgc, nil
}

// Config defines synchronization logics configuration
type Config struct {
	FullTextSearch *fulltextsearch.Config `mapstructure:"fullTextSearch"`
}

// Stop logics
func (lgc *Logics) Stop() {
	lgc.scheduler.Stop()
}
