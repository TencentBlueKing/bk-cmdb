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

func (s *coreService) CreateServiceInstance(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	instance := metadata.ServiceInstance{}
	if err := mapstr.DecodeFromMapStr(&instance, data); err != nil {
		blog.Errorf("CreateServiceInstance failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().CreateServiceInstance(params, instance)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) GetServiceInstance(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceInstanceIDStr := pathParams(common.BKServiceInstanceIDField)
	if len(serviceInstanceIDStr) == 0 {
		blog.Errorf("GetServiceInstance failed, path parameter `%s` empty, rid: %s", common.BKServiceInstanceIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}

	serviceInstanceID, err := strconv.ParseInt(serviceInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceInstance failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceInstanceIDField, serviceInstanceIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}

	result, err := s.core.ProcessOperation().GetServiceInstance(params, serviceInstanceID)
	if err != nil {
		blog.Errorf("GetServiceInstance failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) ListServiceInstances(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	// filter parameter
	fp := struct {
		BusinessID        int64             `json:"bk_biz_id" field:"bk_biz_id"`
		ServiceTemplateID int64             `json:"service_template_id"`
		HostID            int64             `json:"host_id"`
		Page              metadata.BasePage `json:"page" field:"page"`
	}{}

	if err := mapstr.DecodeFromMapStr(&fp, data); err != nil {
		blog.Errorf("ListServiceInstances failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceTemplates failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	result, err := s.core.ProcessOperation().ListServiceInstance(params, fp.BusinessID, fp.ServiceTemplateID, fp.HostID, fp.Page)
	if err != nil {
		blog.Errorf("ListServiceInstance failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateServiceInstance(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceInstanceIDStr := pathParams(common.BKServiceInstanceIDField)
	if len(serviceInstanceIDStr) == 0 {
		blog.Errorf("UpdateServiceInstance failed, path parameter `%s` empty, rid: %s", common.BKServiceInstanceIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}

	serviceInstanceID, err := strconv.ParseInt(serviceInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceInstance failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceInstanceIDField, serviceInstanceIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}

	instance := metadata.ServiceInstance{}
	if err := mapstr.DecodeFromMapStr(&instance, data); err != nil {
		blog.Errorf("UpdateServiceInstance failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.ProcessOperation().UpdateServiceInstance(params, serviceInstanceID, instance)
	if err != nil {
		blog.Errorf("UpdateServiceInstance failed, err: %+v, rid: %s", err, params.ReqID)
		return nil, err
	}

	return result, nil
}

func (s *coreService) DeleteServiceInstance(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	serviceInstanceIDStr := pathParams(common.BKServiceInstanceIDField)
	if len(serviceInstanceIDStr) == 0 {
		blog.Errorf("DeleteServiceInstance failed, path parameter `%s` empty, rid: %s", common.BKServiceInstanceIDField, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}

	serviceInstanceID, err := strconv.ParseInt(serviceInstanceIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteServiceInstance failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceInstanceIDField, serviceInstanceIDStr, err, params.ReqID)
		return nil, params.Error.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
	}

	if err := s.core.ProcessOperation().DeleteServiceInstance(params, serviceInstanceID); err != nil {
		blog.Errorf("DeleteServiceInstance failed, err: %+v, rid: %s", err, common.BKServiceInstanceIDField)
		return nil, err
	}

	return nil, nil
}
