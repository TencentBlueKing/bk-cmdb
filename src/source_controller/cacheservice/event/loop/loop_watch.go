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

// Package loop defines event loop watcher
package loop

import (
	"context"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common/blog"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	watchcli "configcenter/src/source_controller/cacheservice/event/watch"
	"configcenter/src/storage/stream/types"
)

// LoopWatcher is the loop watch event flow client
type LoopWatcher struct {
	isMaster discovery.ServiceManageInterface
	watchCli *watchcli.Client
}

// NewLoopWatcher new loop watch event flow client
func NewLoopWatcher(isMaster discovery.ServiceManageInterface, watchCli *watchcli.Client) *LoopWatcher {
	return &LoopWatcher{
		isMaster: isMaster,
		watchCli: watchCli,
	}
}

// LoopWatchTaskOptions is the loop watch event flow task options
type LoopWatchTaskOptions struct {
	Name         string
	CursorType   watch.CursorType
	TokenHandler types.TaskTokenHandler
	EventHandler EventHandler
	TenantChan   <-chan TenantEvent
}

// TenantEvent is the tenant change event for loop watch task
type TenantEvent struct {
	EventType watch.EventType
	TenantID  string
	WatchOpts *watch.WatchEventOptions
}

// AddLoopWatchTask add a loop watch task
func (w *LoopWatcher) AddLoopWatchTask(opts *LoopWatchTaskOptions) error {
	key, err := event.GetResourceKeyWithCursorType(opts.CursorType)
	if err != nil {
		blog.Errorf("get task %s resource key with cursor type %s failed, err: %v", opts.Name, opts.CursorType, err)
		return err
	}

	task := &loopWatchTask{
		name:             opts.Name,
		key:              key,
		isMaster:         w.isMaster,
		watchCli:         w.watchCli,
		tokenHandler:     opts.TokenHandler,
		eventHandler:     opts.EventHandler,
		tenantChan:       opts.TenantChan,
		tenantCancelFunc: make(map[string]context.CancelFunc),
	}
	go task.run()

	return nil
}
