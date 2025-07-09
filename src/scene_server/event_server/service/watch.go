/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

package service

import (
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/watch"
)

// WatchEvent TODO
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

	resp, err := s.engine.CoreAPI.CacheService().Cache().Event().WatchEvent(ctx.Kit.Ctx, ctx.Kit.Header, options)
	if err != nil {
		blog.Errorf("watch event, but call cache service failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespString(resp)
}
