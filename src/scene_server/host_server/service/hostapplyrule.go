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
	"strconv"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/json"
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

	if len(option.ModuleIDs) == 0 {
		blog.Errorf("ListHostApplyRule failed, parameter bk_module_ids empty, rid:%s", err, rid)
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

func (s *Service) GenerateApplyPlan(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	planRequest := metadata.HostApplyPlanRequest{}
	if err := ctx.DecodeInto(&planRequest); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(planRequest.ModuleIDs) == 0 {
		blog.Errorf("GenerateApplyPlan failed, bk_module_ids shouldn't empty, err: %v, rid:%s", err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, "bk_module_ids"))
		return
	}
	result, err := s.generateApplyPlan(ctx, bizID, planRequest)
	if err != nil {
		blog.ErrorJSON("GenerateApplyPlan failed, generateApplyPlan failed, bizID: %s, request: %s, err: %v, rid:%s", bizID, planRequest, err, rid)
		ctx.RespAutoError(err)
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
		ctx.RespAutoError(response)
		return
	}

	ctx.RespEntity(result)
	return
}

func (s *Service) generateApplyPlan(ctx *rest.Contexts, bizID int64, planRequest metadata.HostApplyPlanRequest) (
	metadata.HostApplyPlanResult, errors.CCErrorCoder) {

	rid := ctx.Kit.Rid

	relationReq := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
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
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
		ruleOption)
	if ccErr != nil {
		blog.Errorf("list host apply rule failed, bizID: %d, opt: %#v, err: %v, rid: %s", bizID, ruleOption, ccErr, rid)
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
			rules.Info = append(rules.Info, metadata.HostApplyRule{BizID: bizID, ModuleID: item.ModuleID,
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
		Rules:             finalRules,
		HostModules:       hostModules,
		ConflictResolvers: planRequest.ConflictResolvers,
	}

	planResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
		planOption)
	if ccErr != nil {
		blog.Errorf("generate apply plan failed, bizID: %d, opt: %#v, err: %v, rid: %s", bizID, planOption, ccErr, rid)
		return planResult, ccErr
	}
	planResult.Rules = rules.Info
	return planResult, nil
}

func (s *Service) RunHostApplyRule(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid

	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GenerateApplyPlan failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	planRequest := metadata.HostApplyPlanRequest{}
	if err := ctx.DecodeInto(&planRequest); nil != err {
		ctx.RespAutoError(err)
		return
	}

	planResult, err := s.generateApplyPlan(ctx, bizID, planRequest)
	if err != nil {
		blog.Errorf("generate apply plan failed, bizID: %s, req: %s, err: %v, rid: %s", bizID, planRequest, err, rid)
		ctx.RespAutoError(err)
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// enable host apply on module
		op := &metadata.UpdateOption{
			Condition: map[string]interface{}{
				common.BKModuleIDField: map[string]interface{}{common.BKDBIN: planRequest.ModuleIDs}},
			Data: map[string]interface{}{common.HostApplyEnabledField: true},
		}
		
		_, err := s.Engine.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDModule, op)
		if err != nil {
			blog.Errorf("update instance of module failed, option: %s, err: %v, rid: %s", op, err, rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
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

		saveRuleOp := metadata.BatchCreateOrUpdateApplyRuleOption{Rules: rulesOption}
		if _, ccErr := s.CoreAPI.CoreService().HostApplyRule().BatchUpdateHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
			bizID, saveRuleOp); ccErr != nil {
			blog.Errorf("update host rule failed, bizID: %s, req: %s, err: %v, rid: %s", bizID, saveRuleOp, ccErr, rid)
			return ccErr
		}

		// delete rules
		if len(planRequest.RemoveRuleIDs) > 0 {
			delOp := metadata.DeleteHostApplyRuleOption{
				RuleIDs: planRequest.RemoveRuleIDs,
			}

			if ccErr := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID,
				delOp); ccErr != nil {
				blog.Errorf("delete apply rule failed, bizID: %s, req: %s, err: %v, rid: %s", bizID, delOp, ccErr, rid)
				return ccErr
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(&metadata.RespError{Msg: txnErr})
		return
	}

	// update host operation is not done in a transaction, since the successfully updated hosts need not roll back
	ctx.Kit.Header.Del(common.TransactionIdHeader)

	hostApplyResults, err := s.updateHostPlan(planResult, ctx.Kit)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(hostApplyResults)
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
