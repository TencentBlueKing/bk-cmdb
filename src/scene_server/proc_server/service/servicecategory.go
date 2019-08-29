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

func (ps *ProcServer) ListServiceCategoryWithStatistics(ctx *rest.Contexts) {
	ps.listServiceCategory(ctx, true)
}

func (ps *ProcServer) ListServiceCategory(ctx *rest.Contexts) {
	ps.listServiceCategory(ctx, false)
}

func (ps *ProcServer) listServiceCategory(ctx *rest.Contexts, withStatistics bool) {
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

	listOption := metadata.ListServiceCategoriesOption{
		BusinessID:     bizID,
		WithStatistics: withStatistics,
	}
	/*
		if ps.AuthManager.Enabled() == true {
			authorizedCategoryIDs, err := ps.AuthManager.ListAuthorizedServiceCategoryIDs(ctx.Kit.Ctx, ctx.Kit.Header, bizID)
			if err != nil {
				blog.Errorf("ListAuthorizedServiceCategoryIDs failed, bizID: %d, err: %+v, rid: %s", bizID, err, ctx.Kit.Rid)
				err := ctx.Kit.CCError.Error(common.CCErrCommListAuthorizedResourcedFromIAMFailed)
				ctx.RespAutoError(err)
				return
			}
			if listOption.ServiceCategoryIDs != nil {
				ids := make([]int64, 0)
				for _, id := range listOption.ServiceCategoryIDs {
					if util.InArray(id, authorizedCategoryIDs) == true {
						ids = append(ids, id)
					}
				}
				listOption.ServiceCategoryIDs = ids
			} else {
				listOption.ServiceCategoryIDs = authorizedCategoryIDs
			}
		}
	*/

	list, err := ps.CoreAPI.CoreService().Process().ListServiceCategories(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
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

	/*
		if err := ps.AuthManager.RegisterServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, *category); err != nil {
			blog.Errorf("create service category success, but register to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			err := ctx.Kit.CCError.CCError(common.CCErrCommRegistResourceToIAMFailed)
			ctx.RespAutoError(err)
			return
		}
	*/

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

	/*
		if err := ps.AuthManager.UpdateRegisteredServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, *category); err != nil {
			blog.Errorf("update service category success, but update register to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			err := ctx.Kit.CCError.CCError(common.CCErrCommRegistResourceToIAMFailed)
			ctx.RespAutoError(err)
			return
		}
	*/

	ctx.RespEntity(category)
}

func (ps *ProcServer) DeleteServiceCategory(ctx *rest.Contexts) {
	input := new(metadata.DeleteCategoryInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete service category, but get business id failed, err: %v", err)
		return
	}
	_ = bizID

	/*
		// generate iam resource
		iamResources, err := ps.AuthManager.MakeResourcesByServiceCategoryIDs(ctx.Kit.Ctx, ctx.Kit.Header, meta.Delete, bizID, input.ID)
		if err != nil {
			blog.Errorf("make iam resource by service category failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			err := ctx.Kit.CCError.CCError(common.CCErrCommRegistResourceToIAMFailed)
			ctx.RespAutoError(err)
			return
		}
	*/

	err = ps.CoreAPI.CoreService().Process().DeleteServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, input.ID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "delete service category failed, err: %v", err)
		return
	}

	/*
		// deregister iam resource
		if err := ps.AuthManager.Authorize.DeregisterResource(ctx.Kit.Ctx, iamResources...); err != nil {
			blog.Errorf("delete service category success, but deregister from iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
			err := ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed)
			ctx.RespAutoError(err)
			return
		}
	*/

	ctx.RespEntity(nil)
}
