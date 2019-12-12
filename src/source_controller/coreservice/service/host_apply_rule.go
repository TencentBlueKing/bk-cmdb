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

func (s *coreService) CreateHostApplyRule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.CreateHostApplyRuleOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("CreateHostApplyRule failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.HostApplyRuleOperation().CreateHostApplyRule(params, bizID, option)
	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) UpdateHostApplyRule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	ruleIDStr := pathParams(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)
	}

	option := metadata.UpdateHostApplyRuleOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("UpdateHostApplyRule failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.HostApplyRuleOperation().UpdateHostApplyRule(params, bizID, ruleID, option)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, ruleID: %d, option: %+v, err: %+v, rid: %s", ruleID, option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}

func (s *coreService) DeleteHostApplyRule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.DeleteHostApplyRuleOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	if err := s.core.HostApplyRuleOperation().DeleteHostApplyRule(params, bizID, option.RuleIDs...); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, bizID: %d, ruleID: %d, err: %+v, rid: %s", bizID, option.RuleIDs, err, params.ReqID)
		return nil, err
	}
	return nil, nil
}

func (s *coreService) GetHostApplyRule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	hostApplyRuleIDStr := pathParams(common.HostApplyRuleIDField)
	hostApplyRuleID, err := strconv.ParseInt(hostApplyRuleIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)
	}

	rule, err := s.core.HostApplyRuleOperation().GetHostApplyRule(params, bizID, hostApplyRuleID)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, bizID: %d, ruleID: %d, err: %+v, rid: %s", bizID, hostApplyRuleID, err, params.ReqID)
		return nil, err
	}
	return rule, nil
}

func (s *coreService) ListHostApplyRule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.ListHostApplyRuleOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("ListHostApplyRule failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	hostApplyRuleResult, err := s.core.HostApplyRuleOperation().ListHostApplyRule(params, bizID, option)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return hostApplyRuleResult, nil
}

func (s *coreService) GenerateApplyPlan(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.HostApplyPlanOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("GenerateApplyPlan failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	applyPlans, err := s.core.HostApplyRuleOperation().GenerateApplyPlan(params, bizID, option)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return applyPlans, nil
}

func (s *coreService) SearchRuleRelatedModules(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.SearchRuleRelatedModulesOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	modules, err := s.core.HostApplyRuleOperation().SearchRuleRelatedModules(params, bizID, option)
	if err != nil {
		blog.Errorf("SearchRuleRelatedModules failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, params.ReqID)
		return nil, err
	}
	return modules, nil
}

func (s *coreService) BatchUpdateHostApplyRule(params core.ContextParams, pathParams, queryParams ParamsGetter, data mapstr.MapStr) (interface{}, error) {
	bizIDStr := pathParams(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		return nil, params.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	option := metadata.BatchCreateOrUpdateApplyRuleOption{}
	if err := mapstruct.Decode2Struct(data, &option); err != nil {
		blog.Errorf("BatchUpdateHostApplyRule failed, decode request body failed, body: %+v, err: %v, rid: %s", data, err, params.ReqID)
		return nil, params.Error.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	result, err := s.core.HostApplyRuleOperation().BatchUpdateHostApplyRule(params, bizID, option)
	if err != nil {
		blog.Errorf("BatchUpdateHostApplyRule failed, option: %+v, err: %+v, rid: %s", option, err, params.ReqID)
		return nil, err
	}
	return result, nil
}
