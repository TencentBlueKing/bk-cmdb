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

// Package watch defines the general resource cache data watch logics
package watch

import (
	"fmt"
	"strings"

	"configcenter/pkg/cache/general/mapping"
	"configcenter/pkg/filter"
	"configcenter/pkg/tenant"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/cache/general/cache"
	cachetypes "configcenter/src/source_controller/cacheservice/cache/general/types"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	"configcenter/src/source_controller/cacheservice/event/loop"
	watchcli "configcenter/src/source_controller/cacheservice/event/watch"
)

// Watcher defines mongodb event watcher for general resource
type Watcher struct {
	cache       *cache.Cache
	loopWatcher *loop.LoopWatcher
}

// Init general resource mongodb event watcher
func Init(cache *cache.Cache, isMaster discovery.ServiceManageInterface, watchCli *watchcli.Client) error {
	watcher := &Watcher{
		cache:       cache,
		loopWatcher: loop.NewLoopWatcher(isMaster, watchCli),
	}

	cache.CacheChangeCh() <- struct{}{}

	if err := watcher.watch(); err != nil {
		return err
	}

	return nil
}

// watch general resource event for cache
func (w *Watcher) watch() error {
	resType := w.cache.Key().Resource()
	cursorType, err := mapping.GetCursorTypeByResType(resType)
	if err != nil {
		blog.Errorf("get cursor type by resource type %s failed, err: %v", resType, err)
		return err
	}

	name := fmt.Sprintf("%s%s:%s", common.BKCacheKeyV3Prefix, "common_res", resType)

	loopEventChan := make(chan loop.TenantEvent)

	go w.watchCacheChange(cursorType, name, loopEventChan)

	opts := &loop.LoopWatchTaskOptions{
		Name:         name,
		CursorType:   cursorType,
		TokenHandler: tokenhandler.NewSingleTokenHandler(name),
		EventHandler: w.handleEvents,
		TenantChan:   loopEventChan,
	}

	err = w.loopWatcher.AddLoopWatchTask(opts)
	if err != nil {
		blog.Errorf("watch %s brief biz topo cache failed, err: %v", cursorType, err)
		return err
	}

	return nil
}

func (w *Watcher) watchCacheChange(cursorType watch.CursorType, name string, loopEventChan chan<- loop.TenantEvent) {
	tenantChan := make(<-chan tenant.TenantEvent)
	for {
		select {
		case <-w.cache.CacheChangeCh():
			// get need watched resource when initialization and cache changes
			watchAll, tenantSubResMap := w.cache.NeedWatchRes()
			if watchAll {
				// watch all tenants' event
				tenantChan = tenant.NewTenantEventChan(name)
				continue
			}

			// watch specific tenants' event
			if len(tenantSubResMap) > 0 {
				for tenantID, subRes := range tenantSubResMap {
					loopEventChan <- loop.TenantEvent{
						EventType: watch.Update,
						TenantID:  tenantID,
						WatchOpts: &watch.WatchEventOptions{
							Resource: cursorType,
							Filter: watch.WatchEventFilter{
								SubResources: util.StrArrayUnique(subRes),
							},
						},
					}
				}
				continue
			}
			blog.Infof("watch %s event, but do not need cache, wait until it needs cache again", cursorType)
		case e, ok := <-tenantChan:
			if !ok {
				// tenant chan is closed, cancel all tenants' loop watch task
				for _, tenant := range tenant.GetAllTenants() {
					loopEventChan <- loop.TenantEvent{
						EventType: watch.Delete,
						TenantID:  tenant.TenantID,
					}
				}
				tenant.RemoveTenantEventChan(name)
				tenantChan = make(<-chan tenant.TenantEvent)
				continue
			}

			switch e.EventType {
			case tenant.Create:
				loopEventChan <- loop.TenantEvent{
					EventType: watch.Create,
					TenantID:  e.TenantID,
					WatchOpts: &watch.WatchEventOptions{Resource: cursorType},
				}
			case tenant.Delete:
				loopEventChan <- loop.TenantEvent{
					EventType: watch.Delete,
					TenantID:  e.TenantID,
				}
			}

		}
	}
}

func (w *Watcher) handleEvents(kit *rest.Kit, events []*watch.WatchEventDetail) error {
	if len(events) == 0 {
		return nil
	}

	upsertDataArr, delDataArr := w.aggregateEvent(events, kit.Rid)

	if err := w.cache.AddData(kit, upsertDataArr); err != nil {
		blog.Errorf("add %s cache data failed, err: %v, data: %+v, rid: %s", w.cache.Key().Resource(), err,
			upsertDataArr, kit.Rid)
		return err
	}

	if err := w.cache.RemoveData(kit, delDataArr); err != nil {
		blog.Errorf("delete %s cache data failed, err: %v, data: %+v, rid: %s", w.cache.Key().Resource(), err,
			delDataArr, kit.Rid)
		return err
	}
	return nil
}

func (w *Watcher) aggregateEvent(events []*watch.WatchEventDetail, rid string) ([]cachetypes.WatchEventData,
	[]cachetypes.WatchEventData) {

	upsertDataMap, delDataMap := make(map[string]cachetypes.WatchEventData), make(map[string]cachetypes.WatchEventData)

	for _, e := range events {
		if e.Detail == nil {
			continue
		}

		if e.ChainNode == nil {
			blog.Errorf("event %+v has no chain node, rid: %s", e, rid)
			continue
		}

		data := cachetypes.WatchEventData{
			ChainNode: e.ChainNode,
		}
		switch t := e.Detail.(type) {
		case watch.JsonString:
			data.Data = filter.JsonString(t)
		default:
			continue
		}

		uniqueKey := e.ChainNode.Oid + ":" + strings.Join(e.ChainNode.SubResource, ":")

		switch e.EventType {
		case watch.Create, watch.Update:
			upsertDataMap[uniqueKey] = data
		case watch.Delete:
			_, exists := upsertDataMap[uniqueKey]
			if exists {
				delete(upsertDataMap, uniqueKey)
				continue
			}
			delDataMap[uniqueKey] = data
		}
	}

	upsertDataArr, delDataArr := make([]cachetypes.WatchEventData, 0), make([]cachetypes.WatchEventData, 0)
	for _, data := range upsertDataMap {
		upsertDataArr = append(upsertDataArr, data)
	}
	for _, data := range delDataMap {
		delDataArr = append(delDataArr, data)
	}
	return upsertDataArr, delDataArr
}
