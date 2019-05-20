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
	"configcenter/src/common"
	"strconv"

	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	category := metadata.ServiceCategory{}
	if err := mapstr.DecodeFromMapStr(&category, data); err != nil {
		blog.Errorf("CreateServiceCategory failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().CreateServiceCategory(params, category)
	if err != nil {
		blog.Errorf("CreateServiceCategory failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceCategoryIDField := "service_category_id"
	serviceCategoryIDStr := pathParams(serviceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("GetServiceCategory failed, path parameter `%s` empty", serviceCategoryIDField)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, serviceCategoryIDField)
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v", serviceCategoryIDField, serviceCategoryIDStr, err)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, serviceCategoryIDField)
	}

	result, err := s.core.ProcessOperation().GetServiceCategory(params, serviceCategoryID)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListServiceCategories(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := struct {
		Metadata       metadata.Metadata `json:"metadata" field:"metadata"`
		WithStatistics bool              `json:"with_statistics" field:"with_statistics"`
	}{}

	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListServiceCategories failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	bizID, err := metadata.BizIDFromMetadata(fp.Metadata)
	if err != nil {
		blog.Errorf("ListServiceCategories failed, parse business id from metadata failed, metadata: %+v, err: %v", fp.Metadata, err)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}
	if bizID == 0 {
		blog.Errorf("ListServiceCategories failed, business id can't be empty, metadata: %+v, err: %v", fp.Metadata, err)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	result, err := s.core.ProcessOperation().ListServiceCategories(params, bizID, fp.WithStatistics)
	if err != nil {
		blog.Errorf("ListServiceCategories failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceCategoryIDField := "service_category_id"
	serviceCategoryIDStr := pathParams(serviceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("UpdateServiceCategory failed, path parameter `%s` empty", serviceCategoryIDField)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, serviceCategoryIDField)
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v", serviceCategoryIDField, serviceCategoryIDStr, err)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, serviceCategoryIDField)
	}

	category := metadata.ServiceCategory{}
	if err := mapstr.DecodeFromMapStr(&category, data); err != nil {
		blog.Errorf("UpdateServiceCategory failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().UpdateServiceCategory(params, serviceCategoryID, category)
	if err != nil {
		blog.Errorf("UpdateServiceCategory failed, err: %+v", err)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteServiceCategory(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceCategoryIDField := "service_category_id"
	serviceCategoryIDStr := pathParams(serviceCategoryIDField)
	if len(serviceCategoryIDStr) == 0 {
		blog.Errorf("DeleteServiceCategory failed, path parameter `%s` empty", serviceCategoryIDField)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, serviceCategoryIDField)
	}

	serviceCategoryID, err := strconv.ParseInt(serviceCategoryIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteServiceCategory failed, convert path parameter %s to int failed, value: %s, err: %v", serviceCategoryIDField, serviceCategoryIDStr, err)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, serviceCategoryIDField)
	}

	if err := s.core.ProcessOperation().DeleteServiceCategory(params, serviceCategoryID); err != nil {
		blog.Errorf("DeleteServiceCategory failed, err: %+v", err)
		return nil, err
	}

	return nil, nil
}
