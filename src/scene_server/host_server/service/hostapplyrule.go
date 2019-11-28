/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) CreateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.CreateHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("CreateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().CreateHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("CreateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(rule))
}

func (s *Service) UpdateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleIDStr := req.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.UpdateHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("UpdateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().UpdateHostApplyRule(srvData.ctx, srvData.header, bizID, ruleID, option)
	if err != nil {
		blog.ErrorJSON("UpdateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(rule))
}

func (s *Service) DeleteHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.DeleteHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	if err := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(srvData.ctx, srvData.header, bizID, option); err != nil {
		blog.ErrorJSON("DeleteHostApplyRule failed, core service DeleteHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(make(map[string]interface{})))
}

func (s *Service) GetHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleIDStr := req.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().GetHostApplyRule(srvData.ctx, srvData.header, bizID, ruleID)
	if err != nil {
		blog.ErrorJSON("GetHostApplyRule failed, core service GetHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(rule))
}

func (s *Service) ListHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.ListHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("ListHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleResult, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, core service ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(ruleResult))
}

func (s *Service) BatchCreateOrUpdateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("BatchCreateOrUpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.BatchCreateOrUpdateApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("BatchCreateOrUpdateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	batchResult, err := s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("BatchCreateOrUpdateHostApplyRule failed, coreservice BatchUpdateHostApplyRule failed, option: %s, result: %s, err: %s, rid:%s", option, batchResult, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	response := meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     batchResult,
	}
	var firstErr errors.CCErrorCoder
	for _, item := range batchResult.Items {
		if err := item.GetError(); err != nil {
			firstErr = err
			break
		}
	}
	if firstErr != nil {
		response.BaseResp = meta.BaseResp{
			Result:      false,
			Code:        firstErr.GetCode(),
			ErrMsg:      firstErr.Error(),
			Permissions: nil,
		}
	}

	_ = resp.WriteEntity(response)
}

func (s *Service) GenerateApplyPlan(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	planRequest := meta.HostApplyPlanRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&planRequest); err != nil {
		blog.Errorf("GenerateApplyPlan failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	result, err := s.generateApplyPlan(srvData, bizID, planRequest)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, generateApplyPlan failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, planRequest, err, rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	response := meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	}
	_ = resp.WriteEntity(response)
	return
}

func (s *Service) generateApplyPlan(srvData *srvComm, bizID int64, planRequest meta.HostApplyPlanRequest) (meta.HostApplyPlanResult, errors.CCErrorCoder) {
	rid := srvData.rid
	var planResult meta.HostApplyPlanResult

	relationRequest := &meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   planRequest.ModuleIDs,
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostRelations, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, relationRequest)
	if err != nil {
		blog.Errorf("generateApplyPlan failed, err: %+v, rid: %s", err, rid)
		return planResult, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if hostRelations.Code != 0 {
		blog.ErrorJSON("generateApplyPlan failed, response failed, filter: %s, response: %s, err: %s, rid: %s", relationRequest, hostRelations, err, rid)
		return planResult, errors.New(hostRelations.Code, hostRelations.ErrMsg)
	}
	hostIDs := make([]int64, 0)
	for _, item := range hostRelations.Data.Info {
		hostIDs = append(hostIDs, item.HostID)
	}
	relationRequest = &meta.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostRelations, err = s.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, relationRequest)
	if err != nil {
		blog.Errorf("generateApplyPlan failed, err: %+v, rid: %s", err, rid)
		return planResult, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if hostRelations.Code != 0 {
		blog.ErrorJSON("generateApplyPlan failed, response failed, filter: %s, response: %s, err: %s, rid: %s", relationRequest, hostRelations, err, rid)
		return planResult, errors.New(hostRelations.Code, hostRelations.ErrMsg)
	}
	hostModuleMap := make(map[int64][]int64)
	moduleIDs := make([]int64, 0)
	for _, item := range hostRelations.Data.Info {
		if _, exist := hostModuleMap[item.HostID]; exist == false {
			hostModuleMap[item.HostID] = make([]int64, 0)
		}
		hostModuleMap[item.HostID] = append(hostModuleMap[item.HostID], item.ModuleID)
		moduleIDs = append(moduleIDs, item.ModuleID)
	}
	hostModules := make([]meta.Host2Modules, 0)
	for hostID, moduleIDs := range hostModuleMap {
		hostModules = append(hostModules, meta.Host2Modules{
			HostID:    hostID,
			ModuleIDs: moduleIDs,
		})
	}

	ruleOption := meta.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page: meta.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(srvData.ctx, srvData.header, bizID, ruleOption)
	if ccErr != nil {
		blog.ErrorJSON("generateApplyPlan failed, ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, ruleOption, ccErr.Error(), rid)
		return planResult, ccErr
	}

	now := time.Now()
	if len(planRequest.AdditionalRules) > 0 {
	OuterLoop:
		for _, item := range planRequest.AdditionalRules {
			for index, rule := range rules.Info {
				if item.ModuleID == rule.ModuleID && item.AttributeID == rule.AttributeID {
					rules.Info[index].PropertyValue = item.PropertyValue
					continue OuterLoop
				}
			}
			rules.Info = append(rules.Info, meta.HostApplyRule{
				ID:              0,
				BizID:           bizID,
				ModuleID:        item.ModuleID,
				AttributeID:     item.AttributeID,
				PropertyValue:   item.PropertyValue,
				Creator:         srvData.user,
				Modifier:        srvData.user,
				CreateTime:      now,
				LastTime:        now,
				SupplierAccount: srvData.ownerID,
			})
		}
	}

	// filter out removed rules
	if planRequest.RemoveRuleIDs == nil {
		planRequest.RemoveRuleIDs = make([]int64, 0)
	}
	if planRequest.IgnoreRuleIDs == nil {
		planRequest.IgnoreRuleIDs = make([]int64, 0)
	}
	finalRules := make([]meta.HostApplyRule, 0)
	for _, item := range rules.Info {
		if util.InArray(item.ID, planRequest.RemoveRuleIDs) == true {
			continue
		}
		if util.InArray(item.ID, planRequest.IgnoreRuleIDs) == true {
			continue
		}
		finalRules = append(finalRules, item)
	}

	planOption := meta.HostApplyPlanOption{
		Rules:             finalRules,
		HostModules:       hostModules,
		ConflictResolvers: planRequest.ConflictResolvers,
	}

	planResult, ccErr = s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(srvData.ctx, srvData.header, bizID, planOption)
	if err != nil {
		blog.ErrorJSON("generateApplyPlan failed, core service GenerateApplyPlan failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, planOption, ccErr.Error(), rid)
		return planResult, ccErr
	}
	planResult.Rules = rules.Info
	return planResult, nil
}

func (s *Service) RunHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	planRequest := meta.HostApplyPlanRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&planRequest); err != nil {
		blog.Errorf("GenerateApplyPlan failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	planResult, err := s.generateApplyPlan(srvData, bizID, planRequest)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, generateApplyPlan failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, planRequest, err, rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// enable host apply on module
	moduleUpdateOption := &meta.UpdateOption{
		Data: map[string]interface{}{
			common.HostApplyEnabledField: true,
		},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: planRequest.ModuleIDs,
			},
		},
	}
	updateModuleResult, err := s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, moduleUpdateOption)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, UpdateInstance of module http failed, option: %s, err: %v, rid:%s", moduleUpdateOption, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	if ccErr := updateModuleResult.CCError(); ccErr != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, UpdateInstance of module failed, option: %s, result: %s, rid:%s", moduleUpdateOption, updateModuleResult, rid)
		result := &meta.RespError{Msg: ccErr}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// save rules to database
	rulesOption := make([]meta.CreateOrUpdateApplyRuleOption, 0)
	for _, rule := range planResult.Rules {
		rulesOption = append(rulesOption, meta.CreateOrUpdateApplyRuleOption{
			AttributeID:   rule.AttributeID,
			ModuleID:      rule.ModuleID,
			PropertyValue: rule.PropertyValue,
		})
	}
	saveRuleOption := meta.BatchCreateOrUpdateApplyRuleOption{
		Rules: rulesOption,
	}
	if _, ccErr := s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(srvData.ctx, srvData.header, bizID, saveRuleOption); ccErr != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, BatchUpdateHostApplyRule failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, saveRuleOption, ccErr, rid)
		result := &meta.RespError{Msg: ccErr}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// delete rules
	if len(planRequest.RemoveRuleIDs) > 0 {
		deleteRuleOption := meta.DeleteHostApplyRuleOption{
			RuleIDs: planRequest.RemoveRuleIDs,
		}
		if ccErr := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(srvData.ctx, srvData.header, bizID, deleteRuleOption); ccErr != nil {
			blog.ErrorJSON("GenerateApplyPlan failed, DeleteHostApplyRule failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, deleteRuleOption, ccErr, rid)
			result := &meta.RespError{Msg: ccErr}
			_ = resp.WriteError(http.StatusBadRequest, result)
			return
		}
	}

	hostApplyResults := make([]meta.HostApplyResult, 0)
	for _, plan := range planResult.Plans {
		hostApplyResult := meta.HostApplyResult{
			HostID: plan.HostID,
		}
		if len(plan.UpdateFields) == 0 {
			continue
		}
		updateData := make(map[string]interface{})
		for _, field := range plan.UpdateFields {
			updateData[field.PropertyID] = field.PropertyValue
		}
		updateOption := &meta.UpdateOption{
			Data: updateData,
			Condition: map[string]interface{}{
				common.BKHostIDField: plan.HostID,
			},
		}
		updateResult, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, updateOption)
		if err != nil {
			blog.ErrorJSON("RunHostApplyRule, update host failed, option: %s, err: %s, rid: %s", updateOption, err.Error(), rid)
			ccErr := srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
			hostApplyResult.SetError(ccErr)
			hostApplyResults = append(hostApplyResults, hostApplyResult)
			continue
		}
		if ccErr := updateResult.CCError(); ccErr != nil {
			blog.ErrorJSON("RunHostApplyRule, update host response failed, option: %s, response: %s, rid: %s", updateOption, updateResult, rid)
			hostApplyResult.SetError(ccErr)
			hostApplyResults = append(hostApplyResults, hostApplyResult)
			continue
		}
		hostApplyResults = append(hostApplyResults, hostApplyResult)
	}

	var ccErr errors.CCErrorCoder
	for _, item := range hostApplyResults {
		if err := item.GetError(); err != nil {
			ccErr = err
			break
		}
	}
	if ccErr != nil {
		result := &meta.RespError{Msg: ccErr}
		result.Data = hostApplyResults
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	result := meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     hostApplyResults,
	}
	_ = resp.WriteEntity(result)
	return
}
