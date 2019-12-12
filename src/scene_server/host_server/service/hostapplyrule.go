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
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
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
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := metadata.CreateHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("CreateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().CreateHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("CreateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(rule))
}

func (s *Service) UpdateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleIDStr := req.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := metadata.UpdateHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("UpdateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().UpdateHostApplyRule(srvData.ctx, srvData.header, bizID, ruleID, option)
	if err != nil {
		blog.ErrorJSON("UpdateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(rule))
}

func (s *Service) DeleteHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := metadata.DeleteHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	if err := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(srvData.ctx, srvData.header, bizID, option); err != nil {
		blog.ErrorJSON("DeleteHostApplyRule failed, core service DeleteHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(make(map[string]interface{})))
}

func (s *Service) GetHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleIDStr := req.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().GetHostApplyRule(srvData.ctx, srvData.header, bizID, ruleID)
	if err != nil {
		blog.ErrorJSON("GetHostApplyRule failed, core service GetHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, err.Error(), rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(rule))
}

func (s *Service) ListHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := metadata.ListHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("ListHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	if len(option.ModuleIDs) == 0 {
		blog.Errorf("ListHostApplyRule failed, parameter bk_module_ids empty, rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids")}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleResult, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, core service ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(ruleResult))
}

func (s *Service) BatchCreateOrUpdateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("BatchCreateOrUpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := metadata.BatchCreateOrUpdateApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("BatchCreateOrUpdateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	batchResult, err := s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("BatchCreateOrUpdateHostApplyRule failed, coreservice BatchUpdateHostApplyRule failed, option: %s, result: %s, err: %s, rid:%s", option, batchResult, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	response := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
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
		response.BaseResp = metadata.BaseResp{
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
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	planRequest := metadata.HostApplyPlanRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&planRequest); err != nil {
		blog.Errorf("GenerateApplyPlan failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	if len(planRequest.ModuleIDs) == 0 {
		blog.Errorf("GenerateApplyPlan failed, bk_module_ids shouldn't empty, err: %v, rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids")}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	result, err := s.generateApplyPlan(srvData, bizID, planRequest)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, generateApplyPlan failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, planRequest, err, rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	var ccErr errors.CCErrorCoder
	for _, item := range result.Plans {
		if err := item.GetError(); err != nil {
			ccErr = err
			break
		}
	}
	if ccErr != nil {
		response := &metadata.RespError{Msg: ccErr}
		response.Data = result
		_ = resp.WriteError(http.StatusBadRequest, response)
		return
	}

	response := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result,
	}
	_ = resp.WriteEntity(response)
	return
}

func (s *Service) generateApplyPlan(srvData *srvComm, bizID int64, planRequest metadata.HostApplyPlanRequest) (metadata.HostApplyPlanResult, errors.CCErrorCoder) {
	rid := srvData.rid
	var planResult metadata.HostApplyPlanResult

	relationRequest := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   planRequest.ModuleIDs,
		Page: metadata.BasePage{
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
		if planRequest.HostIDs != nil {
			if util.InArray(item.HostID, planRequest.HostIDs) == false {
				continue
			}
		}
		hostIDs = append(hostIDs, item.HostID)
	}
	relationRequest = &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
		Page: metadata.BasePage{
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
		if util.InArray(item.HostID, hostIDs) == false {
			continue
		}
		if _, exist := hostModuleMap[item.HostID]; exist == false {
			hostModuleMap[item.HostID] = make([]int64, 0)
		}
		hostModuleMap[item.HostID] = append(hostModuleMap[item.HostID], item.ModuleID)
		moduleIDs = append(moduleIDs, item.ModuleID)
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range hostModuleMap {
		hostModules = append(hostModules, metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: moduleIDs,
		})
	}

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: moduleIDs,
		Page: metadata.BasePage{
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
			rules.Info = append(rules.Info, metadata.HostApplyRule{
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
	finalRules := make([]metadata.HostApplyRule, 0)
	for _, item := range rules.Info {
		if util.InArray(item.ID, planRequest.RemoveRuleIDs) == true {
			continue
		}
		if util.InArray(item.ID, planRequest.IgnoreRuleIDs) == true {
			continue
		}
		finalRules = append(finalRules, item)
	}

	planOption := metadata.HostApplyPlanOption{
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
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	planRequest := metadata.HostApplyPlanRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(&planRequest); err != nil {
		blog.Errorf("GenerateApplyPlan failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	planResult, err := s.generateApplyPlan(srvData, bizID, planRequest)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, generateApplyPlan failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, planRequest, err, rid)
		result := &metadata.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// enable host apply on module
	moduleUpdateOption := &metadata.UpdateOption{
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: planRequest.ModuleIDs,
			},
		},
		Data: map[string]interface{}{
			common.HostApplyEnabledField: true,
		},
	}
	updateModuleResult, err := s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, moduleUpdateOption)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, UpdateInstance of module http failed, option: %s, err: %v, rid:%s", moduleUpdateOption, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	if ccErr := updateModuleResult.CCError(); ccErr != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, UpdateInstance of module failed, option: %s, result: %s, rid:%s", moduleUpdateOption, updateModuleResult, rid)
		result := &metadata.RespError{Msg: ccErr}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// save rules to database
	rulesOption := make([]metadata.CreateOrUpdateApplyRuleOption, 0)
	for _, rule := range planResult.Rules {
		rulesOption = append(rulesOption, metadata.CreateOrUpdateApplyRuleOption{
			AttributeID:   rule.AttributeID,
			ModuleID:      rule.ModuleID,
			PropertyValue: rule.PropertyValue,
		})
	}
	saveRuleOption := metadata.BatchCreateOrUpdateApplyRuleOption{
		Rules: rulesOption,
	}
	if _, ccErr := s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(srvData.ctx, srvData.header, bizID, saveRuleOption); ccErr != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, BatchUpdateHostApplyRule failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, saveRuleOption, ccErr, rid)
		result := &metadata.RespError{Msg: ccErr}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	// delete rules
	if len(planRequest.RemoveRuleIDs) > 0 {
		deleteRuleOption := metadata.DeleteHostApplyRuleOption{
			RuleIDs: planRequest.RemoveRuleIDs,
		}
		if ccErr := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(srvData.ctx, srvData.header, bizID, deleteRuleOption); ccErr != nil {
			blog.ErrorJSON("GenerateApplyPlan failed, DeleteHostApplyRule failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, deleteRuleOption, ccErr, rid)
			result := &metadata.RespError{Msg: ccErr}
			_ = resp.WriteError(http.StatusBadRequest, result)
			return
		}
	}

	hostApplyResults := make([]metadata.HostApplyResult, 0)
	for _, plan := range planResult.Plans {
		hostApplyResult := metadata.HostApplyResult{
			HostID: plan.HostID,
		}
		if len(plan.UpdateFields) == 0 {
			continue
		}
		updateOption := &metadata.UpdateOption{
			Data: plan.GetUpdateData(),
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
		result := &metadata.RespError{Msg: ccErr}
		result.Data = hostApplyResults
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	result := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     hostApplyResults,
	}
	_ = resp.WriteEntity(result)
	return
}

// ListHostRelatedApplyRule 返回主机关联的规则信息（仅返回启用模块的规则）
func (s *Service) ListHostRelatedApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := metadata.ListHostRelatedApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	result, ccErr := s.listHostRelatedApplyRule(srvData, bizID, option)
	if ccErr != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}
	response := metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result,
	}
	_ = resp.WriteEntity(response)
}

func (s *Service) listHostRelatedApplyRule(srvData *srvComm, bizID int64, option metadata.ListHostRelatedApplyRuleOption) (map[int64][]metadata.HostApplyRule, errors.CCErrorCoder) {
	rid := srvData.rid

	relationOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     option.HostIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relationResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, relationOption)
	if err != nil {
		blog.Errorf("listHostRelatedApplyRule failed, GetHostModuleRelation failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, relationOption, err, rid)
		return nil, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := relationResult.CCError(); ccErr != nil {
		blog.Errorf("listHostRelatedApplyRule failed, GetHostModuleRelation failed, option: %s, result: %s, rid: %s", relationOption, relationResult, rid)
		return nil, ccErr
	}
	hostModuleIDMap := make(map[int64][]int64)
	moduleIDs := make([]int64, 0)
	for _, item := range relationResult.Data.Info {
		moduleIDs = append(moduleIDs, item.ModuleID)
		if _, exist := hostModuleIDMap[item.HostID]; exist == false {
			hostModuleIDMap[item.HostID] = make([]int64, 0)
		}
		hostModuleIDMap[item.HostID] = append(hostModuleIDMap[item.HostID], item.ModuleID)
	}

	// filter enabled modules
	moduleFilter := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
			common.HostApplyEnabledField: true,
		},
	}
	moduleResult, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, moduleFilter)
	if err != nil {
		blog.ErrorJSON("listHostRelatedApplyRule failed, ReadInstance of module failed, option: %s, err: %s, rid: %s", moduleFilter, err.Error(), rid)
		return nil, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := moduleResult.CCError(); ccErr != nil {
		blog.ErrorJSON("listHostRelatedApplyRule failed, ReadInstance of module failed, filter: %s, result: %s, rid: %s", moduleFilter, moduleResult, rid)
		return nil, ccErr
	}
	validModuleIDs := make([]int64, 0)
	for _, item := range moduleResult.Data.Info {
		module := struct {
			ModuleID int64 `mapstructure:"bk_module_id" json:"bk_module_id"`
		}{}
		if err := mapstruct.Decode2Struct(item, &module); err != nil {
			blog.ErrorJSON("listHostRelatedApplyRule failed, ReadInstance of module failed, parse module data failed, filter: %s, item: %s, rid: %s", moduleFilter, item, rid)
			return nil, srvData.ccErr.CCError(common.CCErrCommParseDBFailed)
		}
		validModuleIDs = append(validModuleIDs, module.ModuleID)
	}

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: validModuleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	ruleResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(srvData.ctx, srvData.header, bizID, ruleOption)
	if ccErr != nil {
		blog.ErrorJSON("listHostRelatedApplyRule failed, ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, ccErr.Error(), rid)
		return nil, ccErr
	}
	// moduleID -> []hostApplyRule
	moduleRules := make(map[int64][]metadata.HostApplyRule)
	for _, item := range ruleResult.Info {
		if _, exist := moduleRules[item.ModuleID]; exist == false {
			moduleRules[item.ModuleID] = make([]metadata.HostApplyRule, 0)
		}
		moduleRules[item.ModuleID] = append(moduleRules[item.ModuleID], item)
	}

	// hostID -> []moduleIDs
	result := make(map[int64][]metadata.HostApplyRule)
	for _, hostID := range option.HostIDs {
		if _, exist := result[hostID]; exist == false {
			result[hostID] = make([]metadata.HostApplyRule, 0)
		}
		moduleIDs, exist := hostModuleIDMap[hostID]
		if exist == false {
			continue
		}
		for _, moduleID := range moduleIDs {
			rules, exist := moduleRules[moduleID]
			if exist == true {
				result[hostID] = append(result[hostID], rules...)
			}
		}
	}
	return result, nil
}
