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
	errs "errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) CreateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.CreateHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var rule metadata.HostApplyRule
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rule, err = s.CoreAPI.CoreService().HostApplyRule().CreateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
			option)
		if err != nil {
			blog.ErrorJSON("CreateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, "+
				"option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(rule)
}

func (s *Service) UpdateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	ruleIDStr := ctx.Request.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField))
		return
	}

	option := metadata.UpdateHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var rule metadata.HostApplyRule
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		rule, err = s.CoreAPI.CoreService().HostApplyRule().UpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, ruleID, option)
		if err != nil {
			blog.ErrorJSON("UpdateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(rule)
}

func (s *Service) DeleteHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.DeleteHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(option.RuleIDs) == 0 {
		blog.Errorf("DeleteHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "host_apply_rule_ids"))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option); err != nil {
			blog.ErrorJSON("DeleteHostApplyRule failed, core service DeleteHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(make(map[string]interface{}))
}

func (s *Service) GetHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	ruleIDStr := ctx.Request.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField))
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().GetHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, ruleID)
	if err != nil {
		blog.ErrorJSON("GetHostApplyRule failed, core service GetHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, err.Error(), rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(rule)
}

func (s *Service) ListHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := checkIDs(option.ModuleIDs); err != nil {
		blog.Errorf("get module host apply rule failed, parameter bk_module_ids invalid, err: %v, rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids"))
		return
	}

	ruleResult, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
	if err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, core service ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(ruleResult)
}

func (s *Service) BatchCreateOrUpdateHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("BatchCreateOrUpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.BatchCreateOrUpdateApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var batchResult metadata.BatchCreateOrUpdateHostApplyRuleResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		batchResult, err = s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option)
		if err != nil {
			blog.ErrorJSON("BatchCreateOrUpdateHostApplyRule failed, coreservice BatchUpdateHostApplyRule failed, option: %s, result: %s, err: %s, rid:%s", option, batchResult, err, rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	var firstErr errors.CCErrorCoder
	for _, item := range batchResult.Items {
		if err := item.GetError(); err != nil {
			firstErr = err
			break
		}
	}
	if firstErr != nil {
		ctx.RespEntityWithError(batchResult, firstErr)
		return
	}
	ctx.RespEntity(batchResult)
}

// GenerateModuleApplyPlan generate module host apply rule plan
func (s *Service) GenerateModuleApplyPlan(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	planRequest := metadata.HostApplyModulesOption{}
	if err := ctx.DecodeInto(&planRequest); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if planRequest.BizID == 0 {
		blog.Errorf("generate module host apply rule plan failed, bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs(planRequest.ModuleIDs); err != nil {
		blog.Errorf("generate module host apply rule plan failed, bk_module_ids invalid, err: %v, rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids"))
		return
	}

	result, err := s.generateModuleApplyPlan(ctx, &planRequest)
	if err != nil {
		blog.Errorf("generate module apply plan failed, request: %s, err: %v, rid:%s", planRequest, err, rid)
		ctx.RespAutoError(err)
		return
	}

	for _, item := range result.Plans {
		if err := item.GetError(); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntity(result)
	return
}

func (s *Service) generateModuleApplyPlan(ctx *rest.Contexts, planRequest *metadata.HostApplyModulesOption) (
	metadata.HostApplyPlanResult, errors.CCErrorCoder) {

	rid := ctx.Kit.Rid

	relationReq := &metadata.HostModuleRelationRequest{
		ApplicationID: planRequest.BizID,
		ModuleIDArr:   planRequest.ModuleIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}
	if planRequest.HostIDs != nil {
		relationReq.HostIDArr = planRequest.HostIDs
	}
	hostRelations, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationReq)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, rid)
		return metadata.HostApplyPlanResult{}, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	hostModuleMap := make(map[int64][]int64)
	moduleIDs := make([]int64, 0)
	for _, item := range hostRelations.Info {
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
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
		planRequest.BizID, ruleOption)
	if ccErr != nil {
		blog.Errorf("list host apply rule failed, bizID: %d, opt: %v, err: %v, rid: %s",
			planRequest.BizID, ruleOption, ccErr, rid)
		return metadata.HostApplyPlanResult{}, ccErr
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
			rules.Info = append(rules.Info, metadata.HostApplyRule{BizID: planRequest.BizID, ModuleID: item.ModuleID,
				AttributeID: item.AttributeID, PropertyValue: item.PropertyValue, Creator: ctx.Kit.User,
				Modifier: ctx.Kit.User, CreateTime: now, LastTime: now, SupplierAccount: ctx.Kit.SupplierAccount,
			})
		}
	}

	// filter out removed rules
	finalRules := make([]metadata.HostApplyRule, 0)
	for _, item := range rules.Info {
		if util.InArray(item.ID, planRequest.RemoveRuleIDs) || util.InArray(item.ID, planRequest.IgnoreRuleIDs) {
			continue
		}
		finalRules = append(finalRules, item)
	}

	planOption := metadata.HostApplyPlanOption{
		Rules:       finalRules,
		HostModules: hostModules,
	}

	planResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(ctx.Kit.Ctx, ctx.Kit.Header,
		planRequest.BizID, planOption)
	if ccErr != nil {
		blog.Errorf("generate apply plan failed, bizID: %d, opt: %v, err: %v, rid: %s",
			planRequest.BizID, planOption, ccErr, rid)
		return planResult, ccErr
	}
	planResult.Rules = rules.Info
	return planResult, nil
}

func (s *Service) getModuleRes(kit *rest.Kit, bizID int64, moduleIDs []int64,
	srvTemplateIDs []int64) ([]metadata.ModuleInst, error) {
	moduleFilter := &metadata.QueryCondition{
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Fields: []string{common.BKModuleIDField, common.HostApplyEnabledField, common.BKServiceTemplateIDField},
		Condition: mapstr.MapStr{
			common.BKAppIDField: bizID,
		},
		DisableCounter: true,
	}

	if len(moduleIDs) != 0 {
		moduleFilter.Condition[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: moduleIDs,
		}
	}

	if len(srvTemplateIDs) != 0 {
		moduleFilter.Condition[common.BKServiceTemplateIDField] = map[string]interface{}{
			common.BKDBIN: srvTemplateIDs,
		}
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	if err := s.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		moduleFilter, &moduleRes); err != nil {
		blog.Errorf("get module failed, filter: %v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, err
	}

	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get module failed, filter: %v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, err
	}

	return moduleRes.Data.Info, nil
}

func (s *Service) getEnabledModuleRules(kit *rest.Kit, bizID int64, ids []int64) ([]metadata.HostApplyRule, error) {

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: ids,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}

	mouleRules, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header, bizID, ruleOption)
	if err != nil {
		blog.Errorf("list host apply rule failed, bizID: %d, opt: %v, err: %v, rid: %s",
			bizID, ruleOption, err, kit.Rid)
		return nil, err
	}

	return mouleRules.Info, nil
}

func generateCondition(dataStr string, hostIDs []int64) (map[string]interface{}, map[string]interface{}) {
	data := make(map[string]interface{})
	_ = json.Unmarshal([]byte(dataStr), &data)

	andCond := make([]map[string]interface{}, 0)

	for key, value := range data {
		andCond = append(andCond, map[string]interface{}{
			key: map[string]interface{}{common.BKDBNE: value},
		})
	}
	mergeCond := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs},
		common.BKDBAND:       andCond,
	}
	return mergeCond, data
}

func (s *Service) updateHostPlan(planResult metadata.HostApplyPlanResult, kit *rest.Kit) (
	[]metadata.HostApplyResult, errors.CCErrorCoder) {
	var (
		wg       sync.WaitGroup
		firstErr errors.CCErrorCoder
	)

	// update host instances, allow partial success
	updateMap := make(map[string][]int64, 0)
	for _, plan := range planResult.Plans {
		if len(plan.UpdateFields) == 0 {
			continue
		}
		dataStr := plan.GetUpdateDataStr()
		updateMap[dataStr] = append(updateMap[dataStr], plan.HostID)
	}

	hostApplyResults := make([]metadata.HostApplyResult, 0)
	pipeline := make(chan bool, 5)
	for dataStr, hostIDs := range updateMap {

		pipeline <- true
		wg.Add(1)

		go func(dataStr string, hostIDs []int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			mergeCond, data := generateCondition(dataStr, hostIDs)
			counts, cErr := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
				common.BKTableNameBaseHost, []map[string]interface{}{mergeCond})
			if cErr != nil {
				if firstErr == nil {
					firstErr = cErr
				}
				blog.Errorf("get hosts count failed, filter: %+v, err: %v, rid: %s", mergeCond, cErr, kit.Rid)
				return
			}
			if counts[0] == 0 {
				blog.V(5).Infof("no hosts founded, filter: %+v, rid: %s", mergeCond, kit.Rid)
				return
			}

			// If there is no eligible host, then return directly.
			updateOp := &metadata.UpdateOption{Data: data, Condition: mergeCond}

			_, err := s.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header,
				common.BKInnerObjIDHost, updateOp)
			if err != nil {
				blog.Errorf("update host failed, option: %s, err: %v, rid: %s", updateOp, err, kit.Rid)
				for _, hostID := range hostIDs {
					hostApplyResult := metadata.HostApplyResult{HostID: hostID}
					hostApplyResult.SetError(kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
					hostApplyResults = append(hostApplyResults, hostApplyResult)
				}
				if firstErr == nil {
					firstErr = errors.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
				}
				return
			}

			for _, hostID := range hostIDs {
				hostApplyResult := metadata.HostApplyResult{HostID: hostID}
				hostApplyResults = append(hostApplyResults, hostApplyResult)
			}

		}(dataStr, hostIDs)
	}

	wg.Wait()
	if firstErr != nil {
		return nil, firstErr
	}

	return hostApplyResults, nil
}

// ListHostRelatedApplyRule 返回主机关联的规则信息（仅返回启用模块的规则）
func (s *Service) ListHostRelatedApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListHostRelatedApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	result, ccErr := s.listHostRelatedApplyRule(ctx, bizID, option)
	if ccErr != nil {
		blog.Errorf("ListHostRelatedApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommJSONUnmarshalFailed))
		return
	}
	ctx.RespEntity(result)
}

func (s *Service) listHostRelatedApplyRule(ctx *rest.Contexts, bizID int64,
	option metadata.ListHostRelatedApplyRuleOption) (map[int64][]metadata.HostApplyRule, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid

	relationOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     option.HostIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}
	relationResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		relationOption)
	if err != nil {
		blog.Errorf("get host module relation failed, option: %+v, err: %v, rid: %s", relationOption, err, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	hostModuleIDMap := make(map[int64][]int64)
	moduleIDs := make([]int64, 0)
	for _, item := range relationResult.Info {
		moduleIDs = append(moduleIDs, item.ModuleID)
		hostModuleIDMap[item.HostID] = append(hostModuleIDMap[item.HostID], item.ModuleID)
	}

	// filter enabled modules
	moduleFilter := &metadata.QueryCondition{
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Fields: []string{common.BKModuleIDField},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
			common.HostApplyEnabledField: true,
		},
	}
	moduleResult, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, moduleFilter)
	if err != nil {
		blog.Errorf("get module failed, option: %#v, err: %v, rid: %s", moduleFilter, err, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	validModuleIDs := make([]int64, 0)
	for _, item := range moduleResult.Info {
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("parse module id failed, err: %v, module: %#v, rid: %s", err, item, rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		validModuleIDs = append(validModuleIDs, moduleID)
	}

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: validModuleIDs,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	ruleResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
		bizID, ruleOption)
	if ccErr != nil {
		blog.Errorf("list host apply rule failed, bizID: %d, option: %#v, err: %v, rid: %s", bizID, option, ccErr, rid)
		return nil, ccErr
	}
	// moduleID -> []hostApplyRule
	moduleRules := make(map[int64][]metadata.HostApplyRule)
	for _, item := range ruleResult.Info {
		moduleRules[item.ModuleID] = append(moduleRules[item.ModuleID], item)
	}

	// hostID -> []moduleIDs
	result := make(map[int64][]metadata.HostApplyRule)
	for _, hostID := range option.HostIDs {
		moduleIDs, exist := hostModuleIDMap[hostID]
		if !exist {
			continue
		}
		for _, moduleID := range moduleIDs {
			rules, exist := moduleRules[moduleID]
			if exist {
				result[hostID] = append(result[hostID], rules...)
			}
		}
	}
	return result, nil
}

// GetTemplateHostApplyStatus get service template host apply status
func (s *Service) GetTemplateHostApplyStatus(ctx *rest.Contexts) {
	param := metadata.GetHostApplyStatusParam{}
	if err := ctx.DecodeInto(&param); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if param.ApplicationID == 0 {
		blog.Errorf("bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	moduleFilter := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField, common.BKServiceTemplateIDField},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: param.ModuleIDs,
			},
		},
	}
	moduleResult, err := s.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, moduleFilter)
	if err != nil {
		blog.Errorf("get module failed, option: %v, err: %v, rid: %s", moduleFilter, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	moduleToTemplate := make(map[int64]int64)
	templateIDs := make([]int64, 0)
	for _, item := range moduleResult.Info {
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("parse bk_module_id failed, err: %v, module: %v, rid: %s", err, item, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		templateID, err := util.GetInt64ByInterface(item[common.BKServiceTemplateIDField])
		if err != nil {
			blog.Errorf("parse service_template_id failed, err: %v, module: %v, rid: %s", err, item, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if templateID == 0 {
			blog.Errorf("get service template from module fail, err: %v, module: %v, rid: %s", err, item, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrFindServiceTemplateByModuleFailed, moduleID))
			return
		}

		moduleToTemplate[moduleID] = templateID
		templateIDs = append(templateIDs, templateID)
	}

	templateToStatus, err := s.getSrvTemplateApplyStatus(ctx.Kit, param.ApplicationID, templateIDs)
	if err != nil {
		blog.Errorf("get service template host apply status failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	result := make([]*metadata.HostApplyStatusResult, 0)
	for _, moduleID := range param.ModuleIDs {
		status := &metadata.HostApplyStatusResult{
			ModuleID:         moduleID,
			HostApplyEnabled: templateToStatus[moduleToTemplate[moduleID]],
		}
		result = append(result, status)
	}

	ctx.RespEntity(result)
}

func (s *Service) getSrvTemplateApplyStatus(kit *rest.Kit, bizID int64, ids []int64) (map[int64]bool, error) {

	option := metadata.ListServiceTemplateOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: ids,
	}
	templteResult, err := s.CoreAPI.CoreService().Process().ListServiceTemplates(kit.Ctx, kit.Header, &option)
	if err != nil {
		blog.Errorf("get service template failed, option: %v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, err
	}

	templateToStatus := make(map[int64]bool)
	for _, template := range templteResult.Info {
		templateToStatus[template.ID] = template.HostApplyEnabled
	}

	return templateToStatus, nil
}

// GenerateTemplateApplyPlan generate service template host apply plan
func (s *Service) GenerateTemplateApplyPlan(ctx *rest.Contexts) {
	planRequest := metadata.HostApplyServiceTemplateOption{}
	if err := ctx.DecodeInto(&planRequest); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if planRequest.BizID == 0 {
		blog.Errorf("generate service template host apply plan, bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs(planRequest.ServiceTemplateIDs); err != nil {
		blog.Errorf("generate service template host apply plan failed, service_template_ids invalid, err: %v, rid: %s",
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "service_template_ids"))
		return
	}

	result, err := s.generateServiceTemplateApplyPlan(ctx.Kit, &planRequest)
	if err != nil {
		blog.Errorf("generate service template apply plan failed, request: %v, err: %v, rid: %s",
			planRequest, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, item := range result.Plans {
		if err := item.GetError(); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntity(result)
	return
}

func (s *Service) generateServiceTemplateApplyPlan(kit *rest.Kit, option *metadata.HostApplyServiceTemplateOption) (
	metadata.HostApplyPlanResult, errors.CCErrorCoder) {

	// 1.找出模版对应的最终rule
	rules, err := s.findSrvTemplateRule(kit, option.BizID, option.ServiceTemplateIDs)
	if err != nil {
		blog.Errorf("list service template host apply rule failed, err: %v, rid: %s", err, kit.Rid)
		return metadata.HostApplyPlanResult{}, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	templateRules := getFinalRule(rules, option)

	// 2.将模版的rule赋值给对应的模块
	moduleRes, err := s.getModuleRes(kit, option.BizID, nil, option.ServiceTemplateIDs)
	if err != nil {
		blog.Errorf("get module resource failed, err: %v, rid: %s", err, kit.Rid)
		return metadata.HostApplyPlanResult{}, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(moduleRes) == 0 {
		return metadata.HostApplyPlanResult{Rules: templateRules}, nil
	}

	moduleIDs := make([]int64, 0)
	tempToModules := make(map[int64][]int64)
	for _, module := range moduleRes {
		moduleIDs = append(moduleIDs, module.ModuleID)
		tempToModules[module.ServiceTemplateID] = append(tempToModules[module.ServiceTemplateID], module.ModuleID)
	}

	finalRules := make([]metadata.HostApplyRule, 0)
	for _, rule := range templateRules {
		moduleIDs, exist := tempToModules[rule.ServiceTemplateID]
		if !exist {
			continue
		}
		for _, moduleID := range moduleIDs {
			rule.ModuleID = moduleID
			finalRules = append(finalRules, rule)
		}
	}

	// 3.查询模块与主机关系
	relationReq := &metadata.HostModuleRelationRequest{
		ApplicationID: option.BizID,
		ModuleIDArr:   moduleIDs,
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}

	hostRelations, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, relationReq)
	if err != nil {
		blog.Errorf("get host module relation failed, err: %v, rid: %s", err, kit.Rid)
		return metadata.HostApplyPlanResult{}, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	hostModuleMap := make(map[int64][]int64)
	for _, item := range hostRelations.Info {
		hostModuleMap[item.HostID] = append(hostModuleMap[item.HostID], item.ModuleID)
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range hostModuleMap {
		hostModules = append(hostModules, metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: moduleIDs,
		})
	}

	// 4.生成预览结果
	planOption := metadata.HostApplyPlanOption{
		Rules:       finalRules,
		HostModules: hostModules,
	}

	planResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(kit.Ctx, kit.Header, option.BizID,
		planOption)
	if ccErr != nil {
		blog.Errorf("generate apply plan failed, bizID: %d, opt: %v, err: %v, rid: %s", option.BizID, planOption,
			ccErr, kit.Rid)
		return planResult, ccErr
	}
	planResult.Rules = templateRules
	return planResult, nil
}

func (s *Service) findSrvTemplateRule(kit *rest.Kit, bizID int64, ids []int64) ([]metadata.HostApplyRule, error) {
	ruleOption := metadata.ListHostApplyRuleOption{
		ServiceTemplateIDs: ids,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header, bizID, ruleOption)
	if err != nil {
		blog.Errorf("list service template apply rule failed, opt: %v, err: %v, rid: %s", ruleOption, err, kit.Rid)
		return nil, err
	}

	return rule.Info, nil
}

func getFinalRule(rules []metadata.HostApplyRule, option *metadata.HostApplyServiceTemplateOption) []metadata.
	HostApplyRule {

	keyToRule := make(map[string]metadata.HostApplyRule)
	for _, rule := range rules {
		key := ruleKey(rule.ServiceTemplateID, rule.AttributeID)
		keyToRule[key] = rule
	}

	if len(option.AdditionalRules) > 0 {
		for _, item := range option.AdditionalRules {
			key := ruleKey(item.ServiceTemplateID, item.AttributeID)
			if rule, exsit := keyToRule[key]; exsit {
				rule.PropertyValue = item.PropertyValue
				keyToRule[key] = rule
				continue
			}

			keyToRule[key] = metadata.HostApplyRule{BizID: option.BizID, ServiceTemplateID: item.ServiceTemplateID,
				AttributeID: item.AttributeID, PropertyValue: item.PropertyValue}
		}
	}

	finalRules := make([]metadata.HostApplyRule, 0)
	for _, rule := range keyToRule {
		if util.InArray(rule.ID, option.RemoveRuleIDs) || util.InArray(rule.ID, option.IgnoreRuleIDs) {
			continue
		}
		finalRules = append(finalRules, rule)
	}
	return finalRules
}

func ruleKey(id, attrID int64) string {
	return fmt.Sprintf("%d:%d", id, attrID)
}

// GetServiceTemplateHostApplyRule get service template host apply rule
func (s *Service) GetServiceTemplateHostApplyRule(ctx *rest.Contexts) {

	option := metadata.ListHostApplyRuleOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if option.ApplicationID == 0 {
		blog.Errorf("get service template rule failed, bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs(option.ServiceTemplateIDs); err != nil {
		blog.Errorf("get service template rule failed,service_template_ids invalid, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "service_template_ids"))
		return
	}

	ruleResult, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
		option.ApplicationID, option)
	if err != nil {
		blog.Errorf("list host apply rule failed, option: %s, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(ruleResult)
}

// GetModuleInvalidHostCount get module invalid host count
func (s *Service) GetModuleInvalidHostCount(ctx *rest.Contexts) {
	planRequest := metadata.InvalidHostCountOption{}
	if err := ctx.DecodeInto(&planRequest); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if planRequest.ApplicationID == 0 {
		blog.Errorf("get module invalid host count failed, bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs([]int64{planRequest.ID}); err != nil {
		blog.Errorf("get module invalid host count failed, id invalid, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "id"))
		return
	}

	option := &metadata.HostApplyModulesOption{
		HostApplyPlanBase: metadata.HostApplyPlanBase{
			BizID: planRequest.ApplicationID,
		},
		ModuleIDs: []int64{planRequest.ID},
	}
	result, err := s.generateModuleApplyPlan(ctx, option)
	if err != nil {
		blog.Errorf("generate module apply plan failed, request: %s, err: %v, rid:%s", planRequest, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, item := range result.Plans {
		if err := item.GetError(); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntity(&metadata.InvalidHostCountResult{
		Count: result.UnresolvedConflictCount,
	})
}

// GetServiceTemplateInvalidHostCount get service template invalid host count
func (s *Service) GetServiceTemplateInvalidHostCount(ctx *rest.Contexts) {
	planRequest := metadata.InvalidHostCountOption{}
	if err := ctx.DecodeInto(&planRequest); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if planRequest.ApplicationID == 0 {
		blog.Errorf("get service template invalid host count failed, bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs([]int64{planRequest.ID}); err != nil {
		blog.Errorf("get service template invalid host count failed, id invalid, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "id"))
		return
	}

	option := &metadata.HostApplyServiceTemplateOption{
		HostApplyPlanBase: metadata.HostApplyPlanBase{
			BizID: planRequest.ApplicationID,
		},
		ServiceTemplateIDs: []int64{planRequest.ID},
	}
	result, err := s.generateServiceTemplateApplyPlan(ctx.Kit, option)
	if err != nil {
		blog.Errorf("generate service template apply plan failed, request: %v, err: %v, rid:%s",
			planRequest, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	for _, item := range result.Plans {
		if err := item.GetError(); err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntity(&metadata.InvalidHostCountResult{
		Count: result.UnresolvedConflictCount,
	})
	return
}

// GetServiceTemplateHostApplyRuleCount get service template host apply rule count
func (s *Service) GetServiceTemplateHostApplyRuleCount(ctx *rest.Contexts) {
	option := metadata.HostApplyRuleCountOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if option.ApplicationID == 0 {
		blog.Errorf("get service template host apply rule count failed, bk_biz_id shouldn't empty, rid:%s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs(option.ServiceTemplateIDs); err != nil {
		blog.Errorf("get service template host apply rule count failed, service_template_ids invalid, err: %v, rid: %s",
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "service_template_ids"))
		return
	}

	filters := make([]map[string]interface{}, 0)
	for _, serviceTemplateID := range option.ServiceTemplateIDs {
		filter := map[string]interface{}{
			common.BKAppIDField:             option.ApplicationID,
			common.BKServiceTemplateIDField: serviceTemplateID,
		}
		filters = append(filters, filter)

	}

	counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKTableNameHostApplyRule, filters)
	if err != nil || len(counts) != len(option.ServiceTemplateIDs) {
		blog.Errorf("get count failed, filter: %s, err: %v, rid: %s", filters, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := make([]metadata.HostApplyRuleCountResult, 0)
	for idx, serviceTemplateID := range option.ServiceTemplateIDs {
		templateToCount := metadata.HostApplyRuleCountResult{
			ServiceTemplateID: serviceTemplateID,
			Count:             counts[idx],
		}
		result = append(result, templateToCount)
	}

	ctx.RespEntity(result)
}

func checkIDs(ids []int64) error {
	if len(ids) == 0 {
		return errs.New("the parameters length is 0")
	}

	if util.InArray(0, ids) {
		return errs.New("the parameters can not have 0 value")
	}

	return nil
}

// GetModuleFinalRules get module final rules priority from template
func (s *Service) GetModuleFinalRules(ctx *rest.Contexts) {
	option := metadata.ModuleFinalRulesParam{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if option.ApplicationID == 0 {
		blog.Errorf("get module final rules failed, bk_biz_id shouldn't empty, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_biz_id"))
		return
	}

	if err := checkIDs(option.ModuleIDs); err != nil {
		blog.Errorf("get module final rules failed, bk_module_ids invalid, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids"))
		return
	}

	rules, err := s.getRulesPriorityFromTemplate(ctx.Kit, option.ModuleIDs, option.ApplicationID)
	if err != nil {
		blog.Errorf("get module rule failed, err: %v, rid: %s", err, ctx.Kit)
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(rules)
}
