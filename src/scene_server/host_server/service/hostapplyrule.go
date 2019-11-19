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

	batchResult := meta.BatchCreateOrUpdateHostApplyRuleResult{
		Items: make([]meta.CreateOrUpdateHostApplyRuleResult, 0),
	}
	var firstErr errors.CCErrorCoder
	for index, item := range option.Rules {
		var rule meta.HostApplyRule
		var err errors.CCErrorCoder
		if item.RuleID > 0 {
			updateOption := meta.UpdateHostApplyRuleOption{
				PropertyValue: item.PropertyValue,
			}
			rule, err = s.CoreAPI.CoreService().HostApplyRule().UpdateHostApplyRule(srvData.ctx, srvData.header, bizID, item.RuleID, updateOption)
			if err != nil {
				blog.ErrorJSON("BatchCreateOrUpdateHostApplyRule failed, core service UpdateHostApplyRule failed, bizID: %s, ruleID: %s, option: %s, err: %s, rid: %s", bizID, item.RuleID, updateOption, err.Error(), rid)
			}
		} else {
			createOption := meta.CreateHostApplyRuleOption{
				AttributeID:   item.AttributeID,
				ModuleID:      item.ModuleID,
				PropertyValue: item.PropertyValue,
			}
			rule, err = s.CoreAPI.CoreService().HostApplyRule().CreateHostApplyRule(srvData.ctx, srvData.header, bizID, createOption)
			if err != nil {
				blog.ErrorJSON("BatchCreateOrUpdateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, createOption, err.Error(), rid)
			}
		}
		itemResult := meta.CreateOrUpdateHostApplyRuleResult{
			Index:   index,
			Rule:    rule,
			ErrCode: 0,
			ErrMsg:  "",
		}
		if err != nil {
			itemResult.ErrCode = err.GetCode()
			itemResult.ErrMsg = err.Error()
			if firstErr == nil {
				firstErr = err
			}
		}
		batchResult.Items = append(batchResult.Items, itemResult)
	}
	response := meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     batchResult,
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

	_ = resp.WriteEntity(result)
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
		for _, item := range planRequest.AdditionalRules {
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

	planOption := meta.HostApplyPlanOption{
		Rules:             rules.Info,
		HostModules:       hostModules,
		ConflictResolvers: planRequest.ConflictResolvers,
	}

	planResult, ccErr = s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(srvData.ctx, srvData.header, bizID, planOption)
	if err != nil {
		blog.ErrorJSON("generateApplyPlan failed, core service GenerateApplyPlan failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, planOption, ccErr.Error(), rid)
		return planResult, ccErr
	}

	return planResult, nil
}
