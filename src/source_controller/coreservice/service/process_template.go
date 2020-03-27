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

func (s *coreService) CreateProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	template := metadata.ProcessTemplate{}
	if err := mapstr.DecodeFromMapStr(&template, data); err != nil {
		blog.Errorf("CreateProcessTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().CreateProcessTemplate(params, template)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processTemplateIDStr := pathParams(common.BKProcessTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("GetProcessTemplate failed, path parameter `%s` empty, rid: %s", processTemplateIDStr, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField)
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcessTemplateIDField, processTemplateIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField)
	}

	result, err := s.core.ProcessOperation().GetProcessTemplate(params, processTemplateID)
	if err != nil {
		blog.Errorf("GetProcessTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListProcessTemplates(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := metadata.ListProcessTemplatesOption{}
	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListProcessTemplates failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceTemplates failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	result, err := s.core.ProcessOperation().ListProcessTemplates(params, fp)
	if err != nil {
		blog.Errorf("ListProcessTemplates failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processTemplateIDStr := pathParams(common.BKProcessTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("UpdateProcessTemplate failed, path parameter `%s` empty, rid: %s", common.BKProcessTemplateIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField)
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcessTemplateIDField, processTemplateIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField)
	}

	result, err := s.core.ProcessOperation().UpdateProcessTemplate(params, processTemplateID, data)
	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processTemplateIDStr := pathParams(common.BKProcessTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("DeleteProcessTemplate failed, path parameter `%s` empty, rid: %s", common.BKProcessTemplateIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField)
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcessTemplateIDField, processTemplateIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField)
	}

	if err := s.core.ProcessOperation().DeleteProcessTemplate(params, processTemplateID); err != nil {
		blog.Errorf("DeleteProcessTemplate failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) BatchDeleteProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	input := struct {
		ProcessTemplateIDs []int64 `json:"process_template_ids" field:"process_template_ids"`
	}{}

	if err := mapstr.DecodeFromMapStr(&input, data); err != nil {
		blog.Errorf("BatchDeleteProcessTemplate failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	// TODO: replace with batch delete interface
	for _, id := range input.ProcessTemplateIDs {
		if err := s.core.ProcessOperation().DeleteProcessTemplate(params, id); err != nil {
			blog.Errorf("BatchDeleteProcessTemplate failed, templateID: %d, err: %s, rid: %s", id, err.Error(), params.ReqID)
			return nil, err
		}
	}
	return nil, nil
}
