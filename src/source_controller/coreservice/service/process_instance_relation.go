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
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	relation := metadata.ProcessInstanceRelation{}
	if err := mapstr.SetValueToStructByTags(&relation, data); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, fmt.Errorf("decode request body failed, err: %v", err)
	}

	result, err := s.core.ProcessOperation().CreateProcessInstanceRelation(params, relation)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processInstanceIDField := "process_instance_id"
	processInstanceIDStr := pathParams(processInstanceIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("GetProcessInstanceRelation failed, path parameter `%s` empty", processInstanceIDField)
		return nil, fmt.Errorf("path parameter `%s` empty", processInstanceIDField)
	}

	serviceTemplateID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v", processInstanceIDField, processInstanceIDStr, err)
		return nil, fmt.Errorf("convert path parameter %s to int failed, value: %s, err: %v", processInstanceIDField, processInstanceIDStr, err)
	}

	result, err := s.core.ProcessOperation().GetProcessInstanceRelation(params, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := struct {
		Metadata          metadata.Metadata `json:"metadata" field:"metadata"`
		ServiceInstanceID int64             `json:"service_instance_id" field:"service_instance_id"`
		HostID            int64             `json:"host_id" field:"host_id"`
		Page              metadata.BasePage `json:"page" field:"page"`
	}{}

	if err := mapstr.SetValueToStructByTags(&fp, data); err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, params.Error.Errorf(common.CCErrCommHTTPReadBodyFailed)
	}

	bizID, err := metadata.BizIDFromMetadata(fp.Metadata)
	if err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, parse business id from metadata failed, metadata: %+v, err: %v", fp.Metadata, err)
		return nil, fmt.Errorf("parse business id from metadata failed, err: %v", err)
	}
	if bizID == 0 {
		blog.Errorf("ListProcessInstanceRelation failed, business id can't be empty, metadata: %+v, err: %v", fp.Metadata, err)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	result, err := s.core.ProcessOperation().ListProcessInstanceRelation(params, bizID, fp.ServiceInstanceID, fp.HostID, fp.Page)
	if err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, err: %+v", err)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processInstanceIDField := "process_instance_id"
	processInstanceIDStr := pathParams(processInstanceIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("UpdateProcessInstanceRelation failed, path parameter `%s` empty", processInstanceIDField)
		return nil, fmt.Errorf("path parameter `%s` empty", processInstanceIDField)
	}

	processInstanceID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v", processInstanceIDField, processInstanceIDStr, err)
		return nil, fmt.Errorf("convert path parameter %s to int failed, value: %s, err: %v", processInstanceIDField, processInstanceIDStr, err)
	}

	relation := metadata.ProcessInstanceRelation{}
	if err := mapstr.SetValueToStructByTags(&relation, data); err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v", data, err)
		return nil, fmt.Errorf("decode request body failed, err: %v", err)
	}

	result, err := s.core.ProcessOperation().UpdateProcessInstanceRelation(params, processInstanceID, relation)
	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, err: %+v", err)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processInstanceIDField := "process_instance_id"
	processInstanceIDStr := pathParams(processInstanceIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("DeleteProcessInstanceRelation failed, path parameter `%s` empty", processInstanceIDField)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, "process_instance_id")
	}

	processInstanceID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v", processInstanceIDField, processInstanceIDStr, err)
		return nil, fmt.Errorf("convert path parameter %s to int failed, value: %s, err: %v", processInstanceIDField, processInstanceIDStr, err)
	}

	if err := s.core.ProcessOperation().DeleteProcessInstanceRelation(params, processInstanceID); err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, err: %+v", err)
		return nil, err
	}

	return nil, nil
}
