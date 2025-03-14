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
	"sync"
	"time"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/common/watch"
	"configcenter/src/source_controller/cacheservice/event"
	watchcli "configcenter/src/source_controller/cacheservice/event/watch"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/types"
)

type loopWatchTask struct {
	name         string
	key          event.Key
	isMaster     discovery.ServiceManageInterface
	watchCli     *watchcli.Client
	tokenHandler types.TaskTokenHandler
	eventHandler EventHandler
	tenantChan   <-chan TenantEvent

	mu               sync.Mutex
	tenantCancelFunc map[string]context.CancelFunc
}

// run loop watch task
func (t *loopWatchTask) run() {
	for e := range t.tenantChan {
		switch e.EventType {
		case watch.Create:
			t.startTenantTask(e.TenantID, e.WatchOpts)
		case watch.Delete:
			t.stopTenantTask(e.TenantID)
		case watch.Update:
			t.stopTenantTask(e.TenantID)
			t.startTenantTask(e.TenantID, e.WatchOpts)
		}
	}
}

// startTenantTask start loop watch task for new tenant
func (t *loopWatchTask) startTenantTask(tenantID string, opts *watch.WatchEventOptions) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if _, exists := t.tenantCancelFunc[tenantID]; exists {
		return
	}

	ctx, cancel := context.WithCancel(util.SetDBReadPreference(context.Background(), common.SecondaryPreferredMode))
	kit := rest.NewKit().WithCtx(ctx).WithTenant(tenantID)

	t.tenantCancelFunc[tenantID] = cancel

	go t.loopWatch(kit, opts)
	blog.Infof("start tenant %s loop watch task %s, rid: %s", tenantID, t.name, kit.Rid)
}

// stopTenantTask stop loop watch task for removed or disabled tenant
func (t *loopWatchTask) stopTenantTask(tenantID string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if cancel, exists := t.tenantCancelFunc[tenantID]; exists {
		cancel()
		delete(t.tenantCancelFunc, tenantID)
		blog.Infof("stop tenant %s loop watch task %s", tenantID, t.name)
	}
}

// LoopWatch loop watch event flow
func (t *loopWatchTask) loopWatch(kit *rest.Kit, opts *watch.WatchEventOptions) {
	go func() {
		prevStatus := t.isMaster.IsMaster()
		for {
			select {
			case <-kit.Ctx.Done():
				return
			default:
			}

			isMaster := t.isMaster.IsMaster()
			if !isMaster {
				prevStatus = false
				blog.V(4).Infof("watch %s event, but not master, skip.", t.name)
				time.Sleep(time.Minute)
				continue
			}

			// need watch status changed, re-watch from the last cursor with renewed watch resource type
			if !prevStatus {
				prevStatus = isMaster
				var err error
				token, err := t.tokenHandler.GetStartWatchToken(kit.Ctx, kit.TenantID,
					mongodb.Dal("watch").Shard(kit.ShardOpts()))
				if err != nil {
					blog.Errorf("get %s start watch token failed, err: %v, rid: %s", t.name, err, kit.Rid)
					time.Sleep(500 * time.Millisecond)
					continue
				}
				opts.Cursor = token.Token
			}

			retryWrapper(5, func() error {
				return t.doWatch(kit, opts)
			})
		}
	}()
}

// EventHandler is the event handler
type EventHandler func(kit *rest.Kit, events []*watch.WatchEventDetail) error

// DoWatch do watch event for one step
func (t *loopWatchTask) doWatch(kit *rest.Kit, opts *watch.WatchEventOptions) error {
	var events []*watch.WatchEventDetail
	var err error
	if opts.Cursor == "" {
		lastEvent, err := t.watchCli.WatchFromNow(kit, t.key, opts)
		if err != nil {
			blog.Errorf("watch %s event from now failed, re-watch again, err: %v, rid: %s", t.name, err, kit.Rid)
			return err
		}
		events = []*watch.WatchEventDetail{lastEvent}
	} else {
		events, err = t.watchCli.WatchWithCursor(kit, t.key, opts)
		if err != nil {
			if ccErr, ok := err.(errors.CCErrorCoder); ok && ccErr.GetCode() == common.CCErrEventChainNodeNotExist {
				// the cursor does not exist, re-watch from now
				opts.Cursor = ""
				if err = t.tokenHandler.SetLastWatchToken(kit.Ctx, kit.TenantID,
					mongodb.Dal("watch").Shard(kit.ShardOpts()), &types.TokenInfo{Token: ""}); err != nil {
					blog.Errorf("reset %s watch token failed, err: %v, rid: %s", t.name, err, kit.Rid)
					return err
				}

				blog.Errorf("watch event failed, re-watch from now, err: %v, opt: %+v, rid: %s", err,
					*opts, kit.Rid)
				return ccErr
			}
			blog.Errorf("watch event failed, err: %v, opt: %+v, rid: %s", err, *opts, kit.Rid)
			return err
		}
	}

	if len(events) == 0 {
		return nil
	}

	if err = t.eventHandler(kit, events); err != nil {
		blog.Errorf("handle %s events failed, err: %v, events: %+v, rid: %s", t.name, err, events, kit.Rid)
		return err
	}

	opts.Cursor = events[len(events)-1].Cursor
	if err = t.tokenHandler.SetLastWatchToken(kit.Ctx, kit.TenantID, mongodb.Dal("watch").Shard(kit.ShardOpts()),
		&types.TokenInfo{Token: opts.Cursor}); err != nil {
		blog.Errorf("set %s watch token to %s failed, err: %v, rid: %s", t.name, opts.Cursor, err, kit.Rid)
		return err
	}
	return nil
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
