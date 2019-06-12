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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	category := metadata.ServiceCategory{}
	if err := mapstr.DecodeFromMapStr(&category, data); err != nil {
		blog.Errorf("CreateServiceCategory failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().CreateServiceCategory(params, category)
	if err != nil {
		blog.Errorf("CreateServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceCategoryIDStr := pathParams(common.BKServiceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("GetServiceCategory failed, path parameter `%s` empty, rid: %s", common.BKServiceCategoryIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceCategoryIDField, serviceCategoryIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}

	result, err := s.core.ProcessOperation().GetServiceCategory(params, serviceCategoryID)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetDefaultServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	result, err := s.core.ProcessOperation().GetDefaultServiceCategory(params)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListServiceCategories(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := struct {
		BusinessID     int64 `json:"bk_biz_id" field:"bk_biz_id"`
		WithStatistics bool  `json:"with_statistics" field:"with_statistics"`
	}{}

	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListServiceCategories failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceCategories failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id")
	}

	result, err := s.core.ProcessOperation().ListServiceCategories(params, fp.BusinessID, fp.WithStatistics)
	if err != nil {
		blog.Errorf("ListServiceCategories failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceCategoryIDStr := pathParams(common.BKServiceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("UpdateServiceCategory failed, path parameter `%s` empty, rid: %s", common.BKServiceCategoryIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceCategoryIDField, serviceCategoryIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}

	category := metadata.ServiceCategory{}
	if err := mapstr.DecodeFromMapStr(&category, data); err != nil {
		blog.Errorf("UpdateServiceCategory failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().UpdateServiceCategory(params, serviceCategoryID, category)
	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceCategoryIDStr := pathParams(common.BKServiceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("DeleteServiceCategory failed, path parameter `%s` empty, rid: %s", common.BKServiceCategoryIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceCategoryIDField, serviceCategoryIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
	}

	if err := s.core.ProcessOperation().DeleteServiceCategory(params, serviceCategoryID); err != nil {
		blog.Errorf("DeleteServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return nil, nil
}
