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
	"configcenter/src/common/errors"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (s *coreService) CreateProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	relation := &metadata.ProcessInstanceRelation{}
	if err := mapstr.DecodeFromMapStr(relation, data); err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().CreateProcessInstanceRelation(params, relation)
	if err != nil {
		blog.Errorf("CreateProcessInstanceRelation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processInstanceIDStr := pathParams(common.BKProcIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("GetProcessInstanceRelation failed, path parameter `%s` empty, rid: %s", common.BKProcIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcIDField)
	}

	serviceTemplateID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcIDField, processInstanceIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcIDField)
	}

	result, err := s.core.ProcessOperation().GetProcessInstanceRelation(params, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetProcessInstanceRelation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := metadata.ListProcessInstanceRelationOption{}

	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommHTTPReadBodyFailed)
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListProcessInstanceRelation failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	blog.Debug("fp: %v", fp.ServiceInstanceIDs)
	result, err := s.core.ProcessOperation().ListProcessInstanceRelation(params, fp)
	if err != nil {
		blog.Errorf("ListProcessInstanceRelation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	processInstanceIDStr := pathParams(common.BKProcIDField)
	if len(processInstanceIDStr) == 0 {
		blog.Errorf("UpdateProcessInstanceRelation failed, path parameter `%s` empty, rid: %s", common.BKProcIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcIDField)
	}

	processInstanceID, err := strconv.ParseInt(processInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcIDField, processInstanceIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKProcIDField)
	}

	relation := metadata.ProcessInstanceRelation{}
	if err := mapstr.DecodeFromMapStr(&relation, data); err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().UpdateProcessInstanceRelation(params, processInstanceID, relation)
	if err != nil {
		blog.Errorf("UpdateProcessInstanceRelation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteProcessInstanceRelation(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	option := metadata.DeleteProcessInstanceRelationOption{}
	if err := mapstr.DecodeFromMapStr(&option, data); err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if err := s.core.ProcessOperation().DeleteProcessInstanceRelation(params, option); err != nil {
		blog.Errorf("DeleteProcessInstanceRelation failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return nil, nil
}

func (s *coreService) CreateProcessInstance(params core.ContextParams, process *metadata.Process) (*metadata.Process, errors.CCErrorCoder) {
	processBytes, err := json.Marshal(process)
	if err != nil {
		return nil, params.Error.CCError(common.CCErrCommJsonEncode)
	}
	mData := mapstr.MapStr{}
	if err := json.Unmarshal(processBytes, &mData); nil != err && 0 != len(processBytes) {
		return nil, params.Error.CCError(common.CCErrCommJsonDecode)
	}
	inputParam := metadata.CreateModelInstance{
		Data: mData,
	}
	result, err := s.core.InstanceOperation().CreateModelInstance(params, common.BKProcessObjectName, inputParam)
	if err != nil {
		blog.Errorf("CreateProcessInstance failed, CreateModelInstance failed, inputParam: %+v, err: %+v, rid: %s", inputParam, err, params.ReqID)
		return nil, params.Error.CCError(common.CCErrProcCreateProcessFailed)
	}
	process.ProcessID = int64(result.Created.ID)
	return process, nil
}
