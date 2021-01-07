/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/watch"
	ewatcher "configcenter/src/scene_server/event_server/watcher"
	"configcenter/src/source_controller/cacheservice/event"
)

func (s *Service) WatchEvent(ctx *rest.Contexts) {
	resource := ctx.Request.PathParameter("resource")
	options := new(watch.WatchEventOptions)
	if err := ctx.DecodeInto(&options); err != nil {
		blog.Errorf("watch event, but decode request body failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		time.Sleep(500 * time.Millisecond)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	options.Resource = watch.CursorType(resource)

	if err := options.Validate(); err != nil {
		blog.Errorf("watch event, but got invalid request options, err: %v, rid: %s", err, ctx.Kit.Rid)
		time.Sleep(500 * time.Millisecond)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPInputInvalid))
		return
	}

	key, err := event.GetResourceKeyWithCursorType(options.Resource)
	if err != nil {
		blog.Errorf("watch event, but get resource key with cursor type failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		time.Sleep(500 * time.Millisecond)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPInputInvalid))
		return
	}

	// build a resource watcher.
	watcher := ewatcher.NewWatcher(s.ctx, s.cache)

	// watch with cursor
	if len(options.Cursor) != 0 {
		events, err := watcher.WatchWithCursor(key, options, ctx.Kit.Rid)
		if err != nil {
			blog.Errorf("watch event with cursor failed, cursor: %s, err: %v, rid: %s", options.Cursor, err, ctx.Kit.Rid)
			time.Sleep(500 * time.Millisecond)
			ctx.RespAutoError(err)
			return
		}

		// if not events is hit, then we return user's cursor, so that they can watch with this cursor again.
		ctx.RespEntity(s.generateResp(options.Cursor, options.Resource, events))
		return
	}

	// watch with start from
	if options.StartFrom != 0 {
		events, err := watcher.WatchWithStartFrom(key, options, ctx.Kit.Rid)
		if err != nil {
			blog.Errorf("watch event with start from: %s, err: %v, rid: %s", time.Unix(options.StartFrom, 0).Format(time.RFC3339), err, ctx.Kit.Rid)
			time.Sleep(500 * time.Millisecond)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPInputInvalid))
			return
		}

		ctx.RespEntity(s.generateResp("", options.Resource, events))
		return
	}

	// watch from now
	events, err := watcher.WatchFromNow(key, options, ctx.Kit.Rid)
	if err != nil {
		blog.Errorf("watch event from now, err: %v, rid: %s", err, ctx.Kit.Rid)
		time.Sleep(500 * time.Millisecond)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPInputInvalid))
		return
	}

	ctx.RespEntity(s.generateResp("", options.Resource, []*watch.WatchEventDetail{events}))
}

func (s *Service) generateResp(startCursor string, rsc watch.CursorType, events []*watch.WatchEventDetail) *watch.WatchResp {
	result := new(watch.WatchResp)
	if len(events) == 0 {
		result.Watched = false
		if len(startCursor) == 0 {
			result.Events = []*watch.WatchEventDetail{
				{
					Cursor:   watch.NoEventCursor,
					Resource: rsc,
				},
			}
		} else {
			// if user's watch with a start cursor, but we do not find event after this cursor,
			// then we return this start cursor directly, so that they can watch with this cursor for next round.
			result.Events = []*watch.WatchEventDetail{
				{
					Cursor:   startCursor,
					Resource: rsc,
				},
			}
		}

	} else {
		if events[0].Cursor == watch.NoEventCursor {
			result.Watched = false

			if len(startCursor) == 0 {
				// user watch with start form time, or watch from now, then return with NoEventCursor cursor.
				result.Events = []*watch.WatchEventDetail{
					{
						Cursor:   watch.NoEventCursor,
						Resource: rsc,
					},
				}
			} else {
				// if user's watch with a start cursor, but hit a NoEventCursor cursor,
				// then we return this start cursor directly, so that they can watch with this cursor for next round.
				result.Events = []*watch.WatchEventDetail{
					{
						Cursor:   startCursor,
						Resource: rsc,
					},
				}
			}

		} else if events[0].Detail == nil {
			// compatible for event happens but not hit(with different event type), last cursor is returned with no detail
			result.Watched = false
			result.Events = []*watch.WatchEventDetail{events[0]}
		} else {
			result.Watched = true
			result.Events = events
		}
	}

	return result
}
