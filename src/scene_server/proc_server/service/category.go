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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (ps *ProcServer) GetServiceCategory(ctx *rest.Contexts) {
	blog.Debug("here-----------")
	meta := new(metadata.MetadataWrapper)
	if err := ctx.DecodeInto(meta); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(meta.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get service category list, but get business id failed, err: %v", err)
		return
	}

	list, err := ps.CoreAPI.CoreService().Process().ListServiceCategories(ctx.Kit.Ctx, ctx.Kit.Header, bizID, true)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPReadBodyFailed, "get service category list failed, err: %v", err)
		return
	}

	ctx.RespEntity(list)
}

func (ps *ProcServer) CreateServiceCategory(ctx *rest.Contexts) {
	input := new(metadata.ServiceCategory)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create service category, but get business id failed, err: %v", err)
		return
	}

	category, err := ps.CoreAPI.CoreService().Process().CreateServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service category failed, err: %v", err)
		return
	}

	ctx.RespEntity(category)
}

func (ps *ProcServer) UpdateServiceCategory(ctx *rest.Contexts) {
	input := new(metadata.ServiceCategory)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update service category, but get business id failed, err: %v", err)
		return
	}

	category, err := ps.CoreAPI.CoreService().Process().UpdateServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, input.ID, input)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "update service category failed, err: %v", err)
		return
	}

	ctx.RespEntity(category)
}

func (ps *ProcServer) DeleteServiceCategory(ctx *rest.Contexts) {
	input := new(metadata.DeleteCategoryInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete service category, but get business id failed, err: %v", err)
		return
	}

	err = ps.CoreAPI.CoreService().Process().DeleteServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, input.ID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "delete service category failed, err: %v", err)
		return
	}

	ctx.RespEntity(nil)
}
