/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *coreService) CreateServiceCategory(ctx *rest.Contexts) {
	category := metadata.ServiceCategory{}
	if err := ctx.DecodeInto(&category); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateServiceCategory(ctx.Kit, category)
	if err != nil {
		blog.Errorf("CreateServiceCategory failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) GetServiceCategory(ctx *rest.Contexts) {
	serviceCategoryIDStr := ctx.Request.PathParameter(common.BKServiceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("GetServiceCategory failed, path parameter `%s` empty, rid: %s", common.BKServiceCategoryIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField))
		return
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceCategoryIDField, serviceCategoryIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField))
		return
	}

	result, err := s.core.ProcessOperation().GetServiceCategory(ctx.Kit, serviceCategoryID)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) GetDefaultServiceCategory(ctx *rest.Contexts) {
	result, err := s.core.ProcessOperation().GetDefaultServiceCategory(ctx.Kit)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListServiceCategories(ctx *rest.Contexts) {
	// filter parameter
	fp := struct {
		BusinessID     int64 `json:"bk_biz_id" field:"bk_biz_id"`
		WithStatistics bool  `json:"with_statistics" field:"with_statistics"`
	}{}

	if err := ctx.DecodeInto(&fp); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceCategories failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	result, err := s.core.ProcessOperation().ListServiceCategories(ctx.Kit, fp.BusinessID, fp.WithStatistics)
	if err != nil {
		blog.Errorf("ListServiceCategories failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	lang := s.Language(ctx.Kit.Header)
	// translate
	for index := range result.Info {
		if result.Info[index].ServiceCategory.IsBuiltIn {
			result.Info[index].ServiceCategory.Name = s.TranslateServiceCategory(lang, &result.Info[index].ServiceCategory)
		}
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateServiceCategory(ctx *rest.Contexts) {
	serviceCategoryIDStr := ctx.Request.PathParameter(common.BKServiceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("UpdateServiceCategory failed, path parameter `%s` empty, rid: %s", common.BKServiceCategoryIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField))
		return
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceCategoryIDField, serviceCategoryIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField))
		return
	}

	category := metadata.ServiceCategory{}
	if err := ctx.DecodeInto(&category); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().UpdateServiceCategory(ctx.Kit, serviceCategoryID, category)
	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) DeleteServiceCategory(ctx *rest.Contexts) {
	serviceCategoryIDStr := ctx.Request.PathParameter(common.BKServiceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("DeleteServiceCategory failed, path parameter `%s` empty, rid: %s", common.BKServiceCategoryIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField))
		return
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceCategoryIDField, serviceCategoryIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField))
		return
	}

	if err := s.core.ProcessOperation().DeleteServiceCategory(ctx.Kit, serviceCategoryID); err != nil {
		blog.Errorf("DeleteServiceCategory failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}
