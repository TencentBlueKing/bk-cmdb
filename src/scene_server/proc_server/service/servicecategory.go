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
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (ps *ProcServer) ListServiceCategoryWithStatistics(ctx *rest.Contexts) {
	result, err := ps.listServiceCategory(ctx, true)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (ps *ProcServer) ListServiceCategory(ctx *rest.Contexts) {
	result, err := ps.listServiceCategory(ctx, false)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if result == nil {
		blog.Errorf("ListServiceCategory result unexpected nil, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrProcGetServiceCategoryFailed))
	}
	data := metadata.MultipleServiceCategory{
		Count: result.Count,
	}
	for _, item := range result.Info {
		data.Info = append(data.Info, item.ServiceCategory)
	}
	ctx.RespEntity(data)
}

func (ps *ProcServer) listServiceCategory(ctx *rest.Contexts, withStatistics bool) (*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	meta := &struct {
		Metadata metadata.Metadata `json:"metadata"`
		BizID    int64             `json:"bk_biz_id"`
	}{}
	if err := ctx.DecodeInto(meta); err != nil {
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	bizID := meta.BizID
	if bizID == 0 {
		var err error
		bizID, err = metadata.BizIDFromMetadata(meta.Metadata)
		if err != nil {
			blog.Errorf("get service category list failed, parse biz id from metadata failed, metadata: %+v, err: %+v, rid: %s", meta.Metadata, err, rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommHTTPInputInvalid, common.MetadataField)
		}
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

	list, ccErr := ps.CoreAPI.CoreService().Process().ListServiceCategories(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if ccErr != nil {
		blog.Errorf("CoreService ListServiceCategories failed, listOption: %+v, err: %s, rid: %s", listOption, ccErr.Error(), rid)
		return nil, ccErr
	}

	return list, nil
}

func (ps *ProcServer) CreateServiceCategory(ctx *rest.Contexts) {
	input := new(metadata.CreateServiceCategoryOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create service category, but get business id failed, err: %v", err)
			return
		}
	}

	newCategory := &metadata.ServiceCategory{
		BizID:    bizID,
		Name:     input.Name,
		ParentID: input.ParentID,
	}
	category, err := ps.CoreAPI.CoreService().Process().CreateServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, newCategory)
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

	err := ps.CoreAPI.CoreService().Process().DeleteServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header, input.ID)
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
