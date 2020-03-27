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
	"configcenter/src/common/mapstruct"
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
		blog.Errorf("GetServiceTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetServiceTemplateWithStatistics(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
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

	template, err := s.core.ProcessOperation().GetServiceTemplate(params, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetServiceTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	// related service instance count
	serviceInstanceFilter := map[string]interface{}{
		common.BKServiceTemplateIDField: template.ID,
	}
	serviceInstanceCount, err := s.db.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(params.Context)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, filter: %+v, err: %+v, rid: %s", serviceInstanceFilter, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}

	// related service template count
	processRelationFilter := map[string]interface{}{
		common.BKServiceTemplateIDField: template.ID,
	}
	processRelationCount, err := s.db.Table(common.BKTableNameProcessInstanceRelation).Find(processRelationFilter).Count(params.Context)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, filter: %+v, err: %+v, rid: %s", serviceInstanceFilter, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	result := metadata.ServiceTemplateWithStatistics{
		Template:             *template,
		ServiceInstanceCount: int64(serviceInstanceCount),
		ProcessInstanceCount: int64(processRelationCount),
	}
	return result, nil
}

func (s *coreService) ListServiceTemplateDetail(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	if len(bizIDStr) == 0 {
		blog.Errorf("ListServiceTemplateDetail failed, path parameter `%s` empty, rid: %s", common.BKAppIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListServiceTemplateDetail failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKAppIDField, bizIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	input := struct {
		ServiceTemplateIDs []int64 `json:"service_template_ids" mapstructure:"service_template_ids"`
	}{}
	if err := mapstruct.Decode2Struct(data, &input); err != nil {
		blog.ErrorJSON("ListServiceTemplateDetail failed, unmarshal request body failed, value: %s, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: input.ServiceTemplateIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceTemplateResult, ccErr := s.core.ProcessOperation().ListServiceTemplates(params, option)
	if ccErr != nil {
		blog.Errorf("ListServiceTemplateDetail failed, ListServiceTemplate failed, err: %+v, rid: %s", ccErr, params.ReqID)
		return nil, ccErr
	}
	srvTplIDs := make([]int64, 0)
	for _, item := range serviceTemplateResult.Info {
		srvTplIDs = append(srvTplIDs, item.ID)
	}

	listProcessTemplateOption := metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: srvTplIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	listProcResult, ccErr := s.core.ProcessOperation().ListProcessTemplates(params, listProcessTemplateOption)
	if ccErr != nil {
		blog.Errorf("ListServiceTemplateDetail failed, ListProcessTemplates failed, err: %+v, rid: %s", ccErr, params.ReqID)
		return nil, ccErr
	}
	serviceProcessTemplateMap := make(map[int64][]metadata.ProcessTemplate)
	for _, item := range listProcResult.Info {
		if _, exist := serviceProcessTemplateMap[item.ServiceTemplateID]; exist == false {
			serviceProcessTemplateMap[item.ServiceTemplateID] = make([]metadata.ProcessTemplate, 0)
		}
		serviceProcessTemplateMap[item.ServiceTemplateID] = append(serviceProcessTemplateMap[item.ServiceTemplateID], item)
	}

	templateDetails := make([]metadata.ServiceTemplateDetail, 0)
	for _, item := range serviceTemplateResult.Info {
		templateDetail := metadata.ServiceTemplateDetail{
			ServiceTemplate:  item,
			ProcessTemplates: make([]metadata.ProcessTemplate, 0),
		}
		processTemplates, exist := serviceProcessTemplateMap[item.ID]
		if exist == true {
			templateDetail.ProcessTemplates = processTemplates
		}
		templateDetails = append(templateDetails, templateDetail)
	}
	result := metadata.MultipleServiceTemplateDetail{
		Count: serviceTemplateResult.Count,
		Info:  templateDetails,
	}
	return result, nil
}

func (s *coreService) ListServiceTemplates(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := metadata.ListServiceTemplateOption{}

	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListServiceTemplates failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().ListServiceTemplates(params, fp)
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
