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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (p *ProcServer) GetServiceCategory(ctx *rest.Contexts) {
	meta := new(metadata.Metadata)
	if err := ctx.DecodeInto(meta); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(*meta)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get service category list, but get business id failed, err: %v", err)
		return
	}

	list, err := p.CoreAPI.CoreService().Process().ListServiceCategories(ctx.Ctx, ctx.Header, bizID, true)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPReadBodyFailed, "get service category list failed, err: %v", err)
		return
	}

	ctx.RespEntity(metadata.NewSuccessResp(list))
}

func (p *ProcServer) CreateServiceCategory(ctx *rest.Contexts) {
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

	category, err := p.CoreAPI.CoreService().Process().CreateServiceCategory(ctx.Ctx, ctx.Header, input)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service category failed, err: %v", err)
		return
	}

	ctx.RespEntity(metadata.NewSuccessResp(category))
}
