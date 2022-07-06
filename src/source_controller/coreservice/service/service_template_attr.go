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

package service

import (
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// CreateServiceTemplateAttrs create service template attributes
func (s *coreService) CreateServiceTemplateAttrs(ctx *rest.Contexts) {
	opt := new(metadata.CreateSvcTempAttrsOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ids, err := s.core.ProcessOperation().CreateServiceTemplateAttrs(ctx.Kit, opt)
	if err != nil {
		blog.Errorf("create service template attributes(%+v) failed, err: %v, rid: %s", opt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := metadata.RspIDs{
		IDs: make([]int64, len(ids)),
	}

	for idx, id := range ids {
		result.IDs[idx] = int64(id)
	}
	ctx.RespEntity(result)
}

// UpdateServiceTemplateAttribute update service template attribute
func (s *coreService) UpdateServiceTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.UpdateServTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.core.ProcessOperation().UpdateServTempAttr(ctx.Kit, option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// DeleteServiceTemplateAttribute delete service template attribute
func (s *coreService) DeleteServiceTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.DeleteServTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.core.ProcessOperation().DeleteServiceTemplateAttribute(ctx.Kit, option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// ListServiceTemplateAttribute list service template attribute
func (s *coreService) ListServiceTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.ListServTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	data, err := s.core.ProcessOperation().ListServiceTemplateAttribute(ctx.Kit, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(data)
}
