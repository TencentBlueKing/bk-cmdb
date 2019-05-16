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
	"errors"
	"fmt"
	"strconv"

	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	template := metadata.ProcessTemplate{}
	if err := mapstr.SetValueToStructByTags(&template, data); err != nil {
		blog.Errorf("CreateProcessTemplate failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, fmt.Errorf("decode request body failed, err: %v", err)
	}

	result, err := s.core.ProcessOperation().CreateProcessTemplate(params, template)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processTemplateIDField := "process_template_id"
	processTemplateIDStr := pathParams(processTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("GetProcessTemplate failed, path parameter `%s` empty", processTemplateIDStr)
		return nil, fmt.Errorf("path parameter `%s` empty", processTemplateIDField)
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v", processTemplateIDField, processTemplateIDStr, err)
		return nil, fmt.Errorf("convert path parameter %s to int failed, value: %d, err: %v", processTemplateIDField, processTemplateID, err)
	}

	result, err := s.core.ProcessOperation().GetProcessTemplate(params, processTemplateID)
	if err != nil {
		blog.Errorf("GetProcessTemplate failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListProcessTemplates(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := struct {
		Metadata           metadata.Metadata    `json:"metadata" field:"metadata"`
		ServiceTemplateID  int64                `json:"service_template_id" field:"service_template_id"`
		ProcessTemplateIDs *[]int64             `json:"process_template_ids" field:"process_template_ids"`
		Limit              metadata.SearchLimit `json:"limit" field:"limit"`
	}{}

	if err := mapstr.SetValueToStructByTags(&fp, data); err != nil {
		blog.Errorf("ListProcessTemplates failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, fmt.Errorf("decode request body failed, err: %v", err)
	}

	bizID, err := metadata.BizIDFromMetadata(fp.Metadata)
	if err != nil {
		blog.Errorf("ListServiceTemplates failed, parse business id from metadata failed, metadata: %+v, err: %v", fp.Metadata, err)
		return nil, fmt.Errorf("parse business id from metadata failed, err: %v", err)
	}
	if bizID == 0 {
		blog.Errorf("ListServiceTemplates failed, business id can't be empty, metadata: %+v, err: %v", fp.Metadata, err)
		return nil, errors.New("business id can't be empty")
	}

	result, err := s.core.ProcessOperation().ListProcessTemplates(params, bizID, fp.ServiceTemplateID, fp.ProcessTemplateIDs, fp.Limit)
	if err != nil {
		blog.Errorf("ListProcessTemplates failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processTemplateIDField := "process_template_id"
	processTemplateIDStr := pathParams(processTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("UpdateProcessTemplate failed, path parameter `%s` empty", processTemplateIDField)
		return nil, fmt.Errorf("path parameter `%s` empty", processTemplateIDField)
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v", processTemplateIDField, processTemplateIDStr, err)
		return nil, fmt.Errorf("convert path parameter %s to int failed, value: %s, err: %v", processTemplateIDField, processTemplateIDStr, err)
	}

	template := metadata.ProcessTemplate{}
	if err := mapstr.SetValueToStructByTags(&template, data); err != nil {
		blog.Errorf("UpdateProcessTemplate failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, fmt.Errorf("decode request body failed, err: %v", err)
	}

	result, err := s.core.ProcessOperation().UpdateProcessTemplate(params, processTemplateID, template)
	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, err: %+v", err)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteProcessTemplate(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processTemplateIDField := "process_template_id"
	processTemplateIDStr := pathParams(processTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("DeleteProcessTemplate failed, path parameter `%s` empty", processTemplateIDField)
		return nil, fmt.Errorf("path parameter `%s` empty", processTemplateIDField)
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v", processTemplateIDField, processTemplateIDStr, err)
		return nil, fmt.Errorf("convert path parameter %s to int failed, value: %s, err: %v", processTemplateIDField, processTemplateIDStr, err)
	}

	if err := s.core.ProcessOperation().DeleteProcessTemplate(params, processTemplateID); err != nil {
		blog.Errorf("DeleteProcessTemplate failed, err: %+v", err)
		return nil, err
	}

	return nil, nil
}
