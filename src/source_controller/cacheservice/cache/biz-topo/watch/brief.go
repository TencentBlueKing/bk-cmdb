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

	"configcenter/pkg/tenant"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/key"
	topolgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/topo"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	"configcenter/src/source_controller/cacheservice/event/loop"

	"github.com/tidwall/gjson"
)

type briefWatcher struct {
	cacheKey    key.Key
	loopWatcher *loop.LoopWatcher
}

func (w *Watcher) watchBrief() error {
	cacheKey, exists := key.TopoKeyMap[types.BriefType]
	if !exists {
		return fmt.Errorf("watch brief but topo type %s is invalid", types.BriefType)
	}

	watcher := &briefWatcher{
		loopWatcher: loop.NewLoopWatcher(w.isMaster, w.watchCli),
		cacheKey:    cacheKey,
	}

	cursorTypes := []watch.CursorType{watch.Set, watch.Module, watch.MainlineInstance}
	for _, cursorType := range cursorTypes {
		err := watcher.watchEvents(cursorType)
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *briefWatcher) watchEvents(cursorType watch.CursorType) error {
	loopEventChan := make(chan loop.TenantEvent)

	name := fmt.Sprintf("biz topo %s %s", w.cacheKey.Type(), cursorType)
	tenantChan := tenant.NewTenantEventChan(name)
	go func() {
		for e := range tenantChan {
			switch e.EventType {
			case tenant.Create:
				loopEventChan <- loop.TenantEvent{
					EventType: watch.Create,
					TenantID:  e.TenantID,
					WatchOpts: &watch.WatchEventOptions{
						EventTypes: []watch.EventType{watch.Create, watch.Delete},
						Fields:     []string{common.BKAppIDField},
						Resource:   cursorType,
					},
				}
			case tenant.Delete:
				loopEventChan <- loop.TenantEvent{
					EventType: watch.Delete,
					TenantID:  e.TenantID,
				}
			}
		}
	}()

	opts := &loop.LoopWatchTaskOptions{
		Name:         name,
		CursorType:   cursorType,
		TokenHandler: tokenhandler.NewMixTokenHandler(w.cacheKey.Namespace(), string(cursorType)),
		EventHandler: w.handleEvents,
		TenantChan:   loopEventChan,
	}

	err := w.loopWatcher.AddLoopWatchTask(opts)
	if err != nil {
		blog.Errorf("watch %s brief biz topo cache failed, err: %v, rid: %s", cursorType, err)
		return err
	}

	return nil
}

func (w *briefWatcher) handleEvents(kit *rest.Kit, events []*watch.WatchEventDetail) error {
	bizList := make([]int64, 0)

	for _, event := range events {
		if event.Detail == nil {
			continue
		}

		var bizID int64
		switch t := event.Detail.(type) {
		case watch.JsonString:
			bizID = gjson.Get(string(t), common.BKAppIDField).Int()
		}

		if bizID == 0 {
			continue
		}

		blog.V(4).Infof("watch brief biz topo cache, received resource:%s, tenant: %s, biz: %d, cursor: %s, rid: %s",
			event.Resource, kit.TenantID, bizID, event.Cursor, kit.Rid)

		bizList = append(bizList, bizID)
	}

	bizList = util.IntArrayUnique(bizList)
	topolgc.AddRefreshBizTopoTask(types.BriefType, map[string][]int64{kit.TenantID: bizList}, kit.Rid)
	return nil
}
