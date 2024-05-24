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
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"configcenter/pkg/cache/general/mapping"
	"configcenter/pkg/filter"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/cache/general/cache"
	cachetypes "configcenter/src/source_controller/cacheservice/cache/general/types"
	tokenhandler "configcenter/src/source_controller/cacheservice/cache/token-handler"
	"configcenter/src/source_controller/cacheservice/event"
	watchcli "configcenter/src/source_controller/cacheservice/event/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

// Watcher defines mongodb event watcher for general resource
type Watcher struct {
	cache        *cache.Cache
	isMaster     discovery.ServiceManageInterface
	watchCli     *watchcli.Client
	tokenHandler *tokenhandler.SingleHandler
	eventKey     event.Key
}

// Init general resource mongodb event watcher
func Init(cache *cache.Cache, isMaster discovery.ServiceManageInterface, watchCli *watchcli.Client) error {
	watcher := &Watcher{
		cache:    cache,
		isMaster: isMaster,
		watchCli: watchCli,
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

	w.eventKey, err = event.GetResourceKeyWithCursorType(cursorType)
	if err != nil {
		blog.Errorf("get event key with cursor type %s failed, err: %v", cursorType, err)
		return err
	}

	name := fmt.Sprintf("%s%s:%s", common.BKCacheKeyV3Prefix, "common_res", resType)
	w.tokenHandler = tokenhandler.NewSingleTokenHandler(name, mongodb.Client())

	ctx := util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode)
	go w.loopWatch(ctx, cursorType)
	return nil
}

func (w *Watcher) loopWatch(ctx context.Context, cursorType watch.CursorType) {
	prevStatus := w.isMaster.IsMaster()

	opts := &watch.WatchEventOptions{
		Resource: cursorType,
	}

	for {
		// get need watched resource when initialization and cache changes, wait until the cache needs watch
		select {
		case <-w.cache.CacheChangeCh():
			for {
				watchAll, subRes := w.cache.NeedWatchRes()
				if watchAll || len(subRes) > 0 {
					opts.Filter.SubResources = subRes
					break
				}
				blog.Infof("watch %s event, but do not need cache, wait until it needs cache again", cursorType)
				prevStatus = false
				select {
				case <-w.cache.CacheChangeCh():
				}
			}
		default:
		}

		isMaster := w.isMaster.IsMaster()
		if !isMaster {
			prevStatus = false
			blog.V(4).Infof("watch %s event, but not master, skip.", cursorType)
			time.Sleep(time.Minute)
			continue
		}

		// need watch status changed, re-watch from the last cursor with renewed watch resource type
		if !prevStatus {
			prevStatus = isMaster
			var err error
			opts.Cursor, err = w.tokenHandler.GetStartWatchToken(ctx)
			if err != nil {
				blog.Errorf("get %s start watch token failed, err: %v", cursorType, err)
				time.Sleep(500 * time.Millisecond)
				continue
			}
			select {
			case w.cache.CacheChangeCh() <- struct{}{}:
			default:
			}
			continue
		}

		retryWrapper(5, func() error {
			return w.doWatch(ctx, opts)
		})
	}
}

func (w *Watcher) doWatch(ctx context.Context, opts *watch.WatchEventOptions) error {
	kit := &rest.Kit{
		Rid:             util.GenerateRID(),
		Header:          make(http.Header),
		Ctx:             ctx,
		CCError:         errors.NewFromCtx(errors.EmptyErrorsSetting).CreateDefaultCCErrorIf("zh-cn"),
		User:            common.CCSystemOperatorUserName,
		SupplierAccount: common.BKSuperOwnerID,
	}

	var events []*watch.WatchEventDetail
	var err error
	if opts.Cursor == "" {
		lastEvent, err := w.watchCli.WatchFromNow(kit, w.eventKey, opts)
		if err != nil {
			blog.Errorf("watch %s event from now failed, re-watch again, err: %v, rid: %s", opts.Resource, err, kit.Rid)
			return err
		}
		events = []*watch.WatchEventDetail{lastEvent}
	} else {
		events, err = w.watchCli.WatchWithCursor(kit, w.eventKey, opts)
		if err != nil {
			if ccErr, ok := err.(errors.CCErrorCoder); ok && ccErr.GetCode() == common.CCErrEventChainNodeNotExist {
				// the cursor does not exist, re-watch from now
				opts.Cursor = ""
				if err = w.tokenHandler.ResetWatchToken(types.TimeStamp{Sec: uint32(time.Now().Unix())}); err != nil {
					blog.Errorf("reset %s watch token failed, err: %v, rid: %s", opts.Resource, err, kit.Rid)
					return err
				}

				blog.Errorf("watch event failed, re-watch from now, err: %v, opt: %+v, rid: %s", err, opts, kit.Rid)
				return ccErr
			}
			blog.Errorf("watch event failed, err: %v, opt: %+v, rid: %s", err, opts, kit.Rid)
			return err
		}
	}

	if len(events) == 0 {
		return nil
	}

	upsertDataArr, delDataArr := w.aggregateEvent(events)

	if err = w.cache.AddData(ctx, upsertDataArr, kit.Rid); err != nil {
		blog.Errorf("add %s cache data failed, err: %v, data: %+v, rid: %s", opts.Resource, err, upsertDataArr, kit.Rid)
		return err
	}

	if err = w.cache.RemoveData(ctx, delDataArr, kit.Rid); err != nil {
		blog.Errorf("delete %s cache data failed, err: %v, data: %+v, rid: %s", opts.Resource, err, delDataArr, kit.Rid)
		return err
	}

	opts.Cursor = events[len(events)-1].Cursor
	if err = w.tokenHandler.SetLastWatchToken(ctx, opts.Cursor); err != nil {
		blog.Errorf("set %s watch token to %s failed, err: %v, rid: %s", opts.Resource, opts.Cursor, err, kit.Rid)
		return err
	}
	return nil
}

func (w *Watcher) aggregateEvent(events []*watch.WatchEventDetail) ([]cachetypes.WatchEventData,
	[]cachetypes.WatchEventData) {

	upsertDataMap, delDataMap := make(map[string]cachetypes.WatchEventData), make(map[string]cachetypes.WatchEventData)

	for _, e := range events {
		if e.Detail == nil {
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

func retryWrapper(maxRetry int, handler func() error) {
	for retry := 0; retry < maxRetry; retry++ {
		err := handler()
		if err == nil {
			return
		}
		time.Sleep(500 * time.Millisecond * time.Duration(retry))
	}
}
