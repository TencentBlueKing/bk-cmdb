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

func (s *coreService) CreateServiceTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	template := metadata.ServiceTemplate{}
	if err := mapstr.DecodeFromMapStr(&template, data); err != nil {
		blog.Errorf("CreateServiceTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().CreateServiceTemplate(params, template)
	if err != nil {
		blog.Errorf("CreateServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetServiceTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceTemplateIDStr := pathParams(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("GetServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	result, err := s.core.ProcessOperation().GetServiceTemplate(params, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListServiceTemplates(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := struct {
		BusinessID        int64             `json:"bk_biz_id" field:"bk_biz_id"`
		ServiceCategoryID int64             `json:"service_category_id" field:"service_category_id"`
		Page              metadata.BasePage `json:"page" field:"page"`
	}{}

	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListServiceTemplates failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceTemplates failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	result, err := s.core.ProcessOperation().ListServiceTemplates(params, fp.BusinessID, fp.ServiceCategoryID, fp.Page)
	if err != nil {
		blog.Errorf("ListServiceTemplates failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateServiceTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceTemplateIDStr := pathParams(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("UpdateServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	template := metadata.ServiceTemplate{}
	if err := mapstr.DecodeFromMapStr(&template, data); err != nil {
		blog.Errorf("UpdateServiceTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().UpdateServiceTemplate(params, serviceTemplateID, template)
	if err != nil {
		blog.Errorf("UpdateServiceTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteServiceTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceTemplateIDStr := pathParams(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("DeleteServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	if err := s.core.ProcessOperation().DeleteServiceTemplate(params, serviceTemplateID); err != nil {
		blog.Errorf("DeleteServiceTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return nil, nil
}
