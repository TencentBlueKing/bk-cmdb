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
	"fmt"
	"strconv"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

/*
transfer模块 实现带实例自动清除的主机转移操作
*/

// TransferHostWithAutoClearServiceInstance 主机转移接口(带服务实例自动清除功能)
// 1. 将主机 bk_host_ids 从 remove_from_node 指定的拓扑节点移除
// 2. 移入到 add_to_modules 指定的模块中
// 3. 自动删除主机在移除模块下的服务实例
// 4. 自动添加主机在新模块上的服务实例
// note:
// - 不允许 remove_from_node 和 add_to_modules 同时为空
// - bk_host_ids 不允许为空
// - 如果 remove_from_node 指定为业务ID，则接口行为是：覆盖更新
// - 如果 remove_from_node 没有指定，仅仅是增量更新，无移除操作
// - 如果 add_to_modules 没有指定，主机将仅仅从 remove_from_node 指定的模块中移除
// - 如果 add_to_modules 是空先机/故障机/待回收模块中的一个，必须显式指定 remove_from_node(可指定成业务节点), 否则报主机不能属于互斥模块错误
// - 如果 add_to_modules 是普通模块，主机当前数据空先机/故障机/待回收模块中的一个，必须显式指定 remove_from_node(可指定成业务节点), 否则报主机不能属于互斥模块错误
// - 模块同时出现在 add_to_modules 和 remove_from_node 时，不会导致对应的服务实例被删除然后重新添加
func (s *Service) TransferHostWithAutoClearServiceInstance(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.V(7).Infof("parse bizID from url failed, bizID: %s, err: %+v, rid: %s", bizIDStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	option := metadata.TransferHostWithAutoClearServiceInstanceOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if ccErr := s.validateTransferHostWithAutoClearServiceInstanceOption(ctx.Kit, bizID, option); ccErr != nil {
		ctx.RespAutoError(ccErr)
		return
	}

	transferPlans, err := s.generateTransferPlans(ctx.Kit, bizID, false, option)
	if err != nil {
		blog.ErrorJSON("generate transfer plans failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	type HostTransferResult struct {
		HostID  int64  `json:"bk_host_id"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	transferResult := make([]HostTransferResult, 0)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {

		// get service instance modules
		moduleIDs := make([]int64, 0)
		for _, item := range option.Options.ServiceInstanceOptions {
			moduleIDs = append(moduleIDs, item.ModuleID)
		}
		modules, err := s.getModules(ctx.Kit, bizID, moduleIDs)
		if err != nil {
			blog.ErrorJSON("TransferHostWithAutoClearServiceInstance, get modules failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), ctx.Kit.Rid)
			return err
		}
		moduleMap := make(map[int64]int64)
		for _, mod := range modules {
			moduleMap[mod.ModuleID] = mod.ServiceTemplateID
		}

		audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), option.HostIDs)
		if err := audit.WithPrevious(ctx.Kit); err != nil {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, get prev module host config for audit failed, err: %s, HostIDs: %+v, rid: %s", err.Error(), option.HostIDs, ctx.Kit.Rid)
			return err
		}

		var transferHostResult *metadata.OperaterException
		var transferOpt interface{}
		var transferErr error
		if transferPlans[0].IsTransferToInnerModule == true {
			transferOption := &metadata.TransferHostToInnerModule{
				ApplicationID: bizID,
				HostID:        option.HostIDs,
				ModuleID:      transferPlans[0].FinalModules[0],
			}
			transferOpt = transferOption
			transferHostResult, transferErr = s.CoreAPI.CoreService().Host().TransferToInnerModule(ctx.Kit.Ctx, ctx.Kit.Header, transferOption)
		} else {
			transferOption := &metadata.HostsModuleRelation{
				ApplicationID: bizID,
				HostID:        option.HostIDs,
				ModuleID:      transferPlans[0].FinalModules,
				IsIncrement:   false,
			}
			transferOpt = transferOption
			transferHostResult, transferErr = s.CoreAPI.CoreService().Host().TransferToNormalModule(ctx.Kit.Ctx, ctx.Kit.Header, transferOption)
		}

		if transferErr != nil {
			blog.ErrorJSON("runTransferPlans failed, transfer hosts failed, option: %s, err: %s, rid: %s", transferOpt, transferErr.Error(), ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if err := transferHostResult.Error(); err != nil {
			blog.ErrorJSON("runTransferPlans failed, transfer hosts failed, option: %s, result: %s, rid: %s", transferOpt, transferHostResult, ctx.Kit.Rid)
			return err
		}

		var firstErr errors.CCErrorCoder
		pipeline := make(chan bool, 300)
		wg := sync.WaitGroup{}
		for _, plan := range transferPlans {
			pipeline <- true
			wg.Add(1)
			go func(plan metadata.HostTransferPlan) {
				var ccErr errors.CCErrorCoder
				defer func() {
					if ccErr != nil {
						if firstErr == nil {
							firstErr = ccErr
						}
					}
					<-pipeline
					wg.Done()
				}()

				// create or update related service instance
				for _, item := range option.Options.ServiceInstanceOptions {
					if item.HostID != plan.HostID {
						continue
					}
					if util.InArray(item.ModuleID, plan.FinalModules) == false {
						continue
					}
					serviceTemplateID, exist := moduleMap[item.ModuleID]
					if !exist {
						blog.ErrorJSON("TransferHostWithAutoClearServiceInstance, but can not find module: %d, bizID: %s, option: %s, err: %s, rid: %s", item.ModuleID, bizID, option, err.Error(), ctx.Kit.Rid)
						ccErr = errors.New(common.CCErrCommParamsInvalid, fmt.Sprintf("module %d not exist", item.ModuleID))
						return
					}
					if ccErr = s.createOrUpdateServiceInstance(ctx, bizID, plan.HostID, serviceTemplateID, item); ccErr != nil {
						return
					}
				}

			}(plan)
		}
		wg.Wait()

		if firstErr != nil {
			return firstErr
		}

		// update host by host apply rule conflict resolvers
		attributeIDs := make([]int64, 0)
		for _, rule := range option.Options.HostApplyConflictResolvers {
			attributeIDs = append(attributeIDs, rule.AttributeID)
		}
		attrCond := &metadata.QueryCondition{
			Fields: []string{common.BKFieldID, common.BKPropertyIDField},
			Page:   metadata.BasePage{Limit: common.BKNoLimit},
			Condition: map[string]interface{}{
				common.BKFieldID: map[string]interface{}{
					common.BKDBIN: attributeIDs,
				},
			},
		}
		attrRes, ccErr := s.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDHost, attrCond)
		if ccErr != nil {
			blog.ErrorJSON("ReadModelAttr failed, err: %v, attrCond: %s, rid: %s", ccErr, attrCond, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if ccErr = attrRes.CCError(); ccErr != nil {
			blog.ErrorJSON("ReadModelAttr failed, err: %s, attrCond: %s, rid: %s", ccErr.Error(), attrCond, ctx.Kit.Rid)
			return ccErr
		}
		attrMap := make(map[int64]string)
		for _, attr := range attrRes.Data.Info {
			attrMap[attr.ID] = attr.PropertyID
		}

		if err := audit.SaveAudit(ctx.Kit); err != nil {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, save audit log failed, err: %s, HostIDs: %+v, rid: %s", err.Error(), option.HostIDs, ctx.Kit.Rid)
			return err
		}

		hostAttrMap := make(map[int64]map[string]interface{})
		for _, rule := range option.Options.HostApplyConflictResolvers {
			if hostAttrMap[rule.HostID] == nil {
				hostAttrMap[rule.HostID] = make(map[string]interface{})
			}
			hostAttrMap[rule.HostID][attrMap[rule.AttributeID]] = rule.PropertyValue
		}

		for hostID, hostData := range hostAttrMap {
			updateOption := &metadata.UpdateOption{
				Data: hostData,
				Condition: map[string]interface{}{
					common.BKHostIDField: hostID,
				},
			}
			updateResult, err := s.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDHost, updateOption)
			if err != nil {
				blog.ErrorJSON("RunHostApplyRule, update host failed, option: %s, err: %v, rid: %s", updateOption, err, ctx.Kit.Rid)
				return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
			}
			if ccErr = updateResult.CCError(); ccErr != nil {
				blog.ErrorJSON("RunHostApplyRule, update host response failed, option: %s, response: %s, rid: %s", updateOption, updateResult, ctx.Kit.Rid)
				return ccErr
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(transferResult, txnErr)
		return
	}
	ctx.RespEntity(transferResult)
	return
}

func (s *Service) createOrUpdateServiceInstance(ctx *rest.Contexts, bizID int64, hostID int64, svcTemplateID int64, serviceInstanceOption metadata.CreateServiceInstanceOption) errors.CCErrorCoder {
	rid := ctx.Kit.Rid

	if svcTemplateID == common.ServiceTemplateIDNotSet {
		input := map[string]interface{}{
			common.BKAppIDField: bizID,
			"bk_module_id":      serviceInstanceOption.ModuleID,
			"instances": []map[string]interface{}{
				{
					"bk_host_id": hostID,
					"processes":  serviceInstanceOption.Processes,
				},
			},
		}
		result, err := s.CoreAPI.ProcServer().Service().CreateServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input)
		if err != nil {
			blog.ErrorJSON("createServiceInstance failed, http failed, option: %s, err: %s, rid: %s", input, err.Error(), ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if result.Result == false || result.Code != 0 {
			blog.ErrorJSON("createServiceInstance failed, option: %s, response: %s, rid: %s", input, result, ctx.Kit.Rid)
			return errors.New(result.Code, result.ErrMsg)
		}
	} else {
		// update process instances
		// update process instance by templateID
		relationOption := &metadata.ListProcessInstanceRelationOption{
			BusinessID: bizID,
			HostID:     hostID,
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		relationResult, err := s.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationOption)
		if err != nil {
			blog.ErrorJSON("update process instance failed, list process relation failed, option: %s, err: %s, rid: %s", relationOption, err.Error(), ctx.Kit.Rid)
			return err
		}
		templateID2ProcessID := make(map[int64]int64)
		for _, relation := range relationResult.Info {
			templateID2ProcessID[relation.ProcessTemplateID] = relation.ProcessID
		}

		processes := make([]map[string]interface{}, 0)
		for _, item := range serviceInstanceOption.Processes {
			templateID := item.ProcessTemplateID
			processID, exist := templateID2ProcessID[templateID]
			if exist == false {
				continue
			}
			process := item.ProcessInfo
			process[common.BKProcessIDField] = processID
			processes = append(processes, process)
		}

		if len(processes) == 0 {
			return nil
		}

		updateProcessOption := map[string]interface{}{
			"bk_biz_id": bizID,
			"processes": processes,
		}
		result, e := s.CoreAPI.ProcServer().Process().UpdateProcessInstance(ctx.Kit.Ctx, ctx.Kit.Header, updateProcessOption)
		if e != nil {
			blog.ErrorJSON("updateProcessInstances failed, input: %s, err: %s, rid: %s", updateProcessOption, e.Error(), rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if result.Result == false || result.Code != 0 {
			blog.ErrorJSON("UpdateProcessInstance failed, option: %s, result: %s, rid: %s", updateProcessOption, result, ctx.Kit.Rid)
			return errors.New(result.Code, result.ErrMsg)
		}
	}
	return nil
}

func (s *Service) generateTransferPlans(kit *rest.Kit, bizID int64, withHostApply bool,
	option metadata.TransferHostWithAutoClearServiceInstanceOption) ([]metadata.HostTransferPlan, errors.CCErrorCoder) {
	rid := kit.Rid

	// get host module config
	hostModuleOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     option.HostIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField, common.BKHostIDField},
	}
	hostModuleResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header,
		hostModuleOption)
	if err != nil {
		blog.ErrorJSON("get host module relation failed, option: %s, err: %s, rid: %s", hostModuleOption, err, rid)
		err := kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		return nil, err
	}
	if err := hostModuleResult.CCError(); err != nil {
		blog.ErrorJSON("get host module relation failed, option: %s, err: %s, rid: %s", hostModuleOption, err, rid)
		return nil, err
	}

	hostModulesIDMap := make(map[int64][]int64)
	for _, item := range hostModuleResult.Data.Info {
		if _, exist := hostModulesIDMap[item.HostID]; exist == false {
			hostModulesIDMap[item.HostID] = make([]int64, 0)
		}
		hostModulesIDMap[item.HostID] = append(hostModulesIDMap[item.HostID], item.ModuleID)
	}

	// get inner modules and default inner module to transfer when hosts is removed from all modules
	innerModules, ccErr := s.getInnerModules(kit, bizID)
	if ccErr != nil {
		return nil, ccErr
	}

	innerModuleIDMap := make(map[int64]struct{}, 0)
	defaultInternalModuleID := int64(0)
	for _, module := range innerModules {
		innerModuleIDMap[module.ModuleID] = struct{}{}
		if module.Default == int64(common.DefaultResModuleFlag) {
			defaultInternalModuleID = module.ModuleID
		}
	}
	if defaultInternalModuleID == 0 {
		blog.InfoJSON("default internal module ID not found, bizID: %s, modules: %s, rid: %s", bizID, innerModules, rid)
	}

	if option.DefaultInternalModule != 0 {
		if _, exists := innerModuleIDMap[option.DefaultInternalModule]; !exists {
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "default_internal_module")
		}
	}
	if option.DefaultInternalModule != 0 {
		defaultInternalModuleID = option.DefaultInternalModule
	}

	transferPlans := make([]metadata.HostTransferPlan, 0)
	for hostID, currentInModules := range hostModulesIDMap {
		transferPlan := generateTransferPlan(currentInModules, option.RemoveFromModules, option.AddToModules,
			defaultInternalModuleID)
		transferPlan.HostID = hostID
		// check module compatibility
		finalModuleCount := len(transferPlan.FinalModules)
		for _, moduleID := range transferPlan.FinalModules {
			if _, exists := innerModuleIDMap[moduleID]; !exists {
				continue
			}
			if finalModuleCount != 1 {
				return nil, kit.CCError.CCError(common.CCErrHostTransferFinalModuleConflict)
			}
			transferPlan.IsTransferToInnerModule = true
		}
		transferPlans = append(transferPlans, transferPlan)
	}

	// if do not need host apply, then return directly.
	if !withHostApply {
		return transferPlans, nil
	}

	// generate host apply plans
	finalModuleIDs := make([]int64, 0)
	for _, item := range transferPlans {
		finalModuleIDs = append(finalModuleIDs, item.FinalModules...)
	}

	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: finalModuleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header, bizID, ruleOption)
	if ccErr != nil {
		blog.ErrorJSON("list apply rule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, ruleOption, ccErr, rid)
		return transferPlans, ccErr
	}

	moduleCondition := metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField},
		Condition: map[string]interface{}{
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: finalModuleIDs,
			},
			common.HostApplyEnabledField: true,
		},
	}
	enabledModules, err := s.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, &moduleCondition)
	if err != nil {
		blog.ErrorJSON("get host apply enabled modules failed, filter: %s, err: %s, rid: %s", moduleCondition, err, rid)
		return transferPlans, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	enableModuleMap := make(map[int64]bool)
	for _, item := range enabledModules.Data.Info {
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			blog.ErrorJSON("parse module from db failed, module: %s, err: %s, rid: %s", item, err, rid)
			return transferPlans, kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		enableModuleMap[moduleID] = true
	}

	hostModules := make([]metadata.Host2Modules, 0)
	for _, item := range transferPlans {
		host2Module := metadata.Host2Modules{
			HostID:    item.HostID,
			ModuleIDs: make([]int64, 0),
		}
		for _, moduleID := range item.FinalModules {
			if _, exist := enableModuleMap[moduleID]; exist {
				host2Module.ModuleIDs = append(host2Module.ModuleIDs, moduleID)
			}
		}
		hostModules = append(hostModules, host2Module)
	}

	planOption := metadata.HostApplyPlanOption{
		Rules:             rules.Info,
		HostModules:       hostModules,
		ConflictResolvers: option.Options.HostApplyConflictResolvers,
	}

	hostApplyPlanResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(kit.Ctx, kit.Header, bizID,
		planOption)
	if ccErr != nil {
		blog.ErrorJSON("generate apply plan failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, planOption, ccErr,
			rid)
		return transferPlans, ccErr
	}
	hostApplyPlanMap := make(map[int64]metadata.OneHostApplyPlan)
	for _, item := range hostApplyPlanResult.Plans {
		hostApplyPlanMap[item.HostID] = item
	}
	for index, transferPlan := range transferPlans {
		applyPlan, ok := hostApplyPlanMap[transferPlan.HostID]
		if !ok {
			continue
		}
		transferPlans[index].HostApplyPlan = applyPlan
	}

	return transferPlans, nil
}

// generateTransferPlan 实现计算主机将从哪个模块移除，添加到哪个模块，最终在哪些模块
// param hostID: 主机ID
// param currentIn: 主机当前所属模块
// param removeFrom: 从哪些模块中移除
// param addTo: 添加到哪些模块
// param defaultInternalModuleID: 默认内置模块ID
func generateTransferPlan(currentIn []int64, removeFrom []int64, addTo []int64,
	defaultInternalModuleID int64) metadata.HostTransferPlan {

	removeFromModuleMap := make(map[int64]struct{})
	for _, moduleID := range removeFrom {
		removeFromModuleMap[moduleID] = struct{}{}
	}

	// 主机最终所在模块列表，包括当前所在模块和新增模块，不包括移出模块
	finalModules := make([]int64, 0)
	finalModuleMap := make(map[int64]struct{})
	// 主机将会被移出的模块列表，包括当前所在模块里在移出模块中的模块
	realRemoveModules := make([]int64, 0)
	realRemoveModuleMap := make(map[int64]struct{})
	for _, moduleID := range currentIn {
		if _, exists := removeFromModuleMap[moduleID]; exists {
			if _, exists := realRemoveModuleMap[moduleID]; !exists {
				realRemoveModuleMap[moduleID] = struct{}{}
				realRemoveModules = append(realRemoveModules, moduleID)
			}
			continue
		}
		if _, exists := finalModuleMap[moduleID]; exists {
			continue
		}
		finalModuleMap[moduleID] = struct{}{}
		finalModules = append(finalModules, moduleID)
	}

	// 主机将会被新加到的模块列表，包括新增模块里不在当前模块的模块
	realAddModules := make([]int64, 0)
	for _, moduleID := range addTo {
		if _, exists := removeFromModuleMap[moduleID]; exists {
			continue
		}
		if _, exists := finalModuleMap[moduleID]; exists {
			continue
		}
		finalModuleMap[moduleID] = struct{}{}
		finalModules = append(finalModules, moduleID)
		realAddModules = append(realAddModules, moduleID)
	}

	if len(finalModules) == 0 {
		finalModules = []int64{defaultInternalModuleID}
	}

	return metadata.HostTransferPlan{
		FinalModules:        finalModules,
		ToRemoveFromModules: realRemoveModules,
		ToAddToModules:      realAddModules,
	}
}

func (s *Service) validateModules(kit *rest.Kit, bizID int64, moduleIDs []int64, field string) errors.CCErrorCoder {
	if len(moduleIDs) == 0 {
		return nil
	}

	moduleIDs = util.IntArrayUnique(moduleIDs)
	filter := []map[string]interface{}{{
		common.BKAppIDField: bizID,
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: moduleIDs,
		},
	}}

	moduleCounts, err := s.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameBaseModule, filter)
	if err != nil {
		return err
	}

	if len(moduleCounts) == 0 || moduleCounts[0] != int64(len(moduleIDs)) {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
	}
	return nil
}

func (s *Service) getModules(kit *rest.Kit, bizID int64, ids []int64) ([]metadata.ModuleInst, errors.CCErrorCoder) {
	query := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{
			common.BKModuleIDField,
			common.BKServiceTemplateIDField,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: ids,
			},
		},
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	err := s.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule, query,
		&moduleRes)
	if err != nil {
		blog.Errorf("get modules failed, input: %#v, err: %v, rid:%s", query, err, kit.Rid)
		return nil, err
	}
	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get modules failed, input: %#v, err: %v, rid:%s", query, err, kit.Rid)
		return nil, err
	}

	return moduleRes.Data.Info, nil
}

func (s *Service) getInnerModules(kit *rest.Kit, bizID int64) ([]metadata.ModuleInst, errors.CCErrorCoder) {
	query := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{
			common.BKModuleIDField,
			common.BKDefaultField,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: common.DefaultFlagDefaultValue,
			},
		},
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	err := s.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule, query,
		&moduleRes)
	if err != nil {
		blog.Errorf("get modules failed, input: %#v, err: %v, rid:%s", query, err, kit.Rid)
		return nil, err
	}
	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get modules failed, input: %#v, err: %v, rid:%s", query, err, kit.Rid)
		return nil, err
	}

	return moduleRes.Data.Info, nil
}

// TransferHostWithAutoClearServiceInstancePreview generate a preview of changes for
// TransferHostWithAutoClearServiceInstance operation
// 接口请求参数跟转移是一致的
// 主机从模块删除时提供了将要删除的服务实例信息
// 主机添加到新模块时，提供了模块对应的服务模板（如果有）
func (s *Service) TransferHostWithAutoClearServiceInstancePreview(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.V(7).Infof("parse bizID from url failed, bizID: %s, err: %+v, rid: %s", bizIDStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField))
		return
	}

	option := metadata.TransferHostWithAutoClearServiceInstanceOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if ccErr := s.validateTransferHostWithAutoClearServiceInstanceOption(ctx.Kit, bizID, option); ccErr != nil {
		ctx.RespAutoError(ccErr)
		return
	}

	transferPlans, ccErr := s.generateTransferPlans(ctx.Kit, bizID, true, option)
	if ccErr != nil {
		blog.ErrorJSON("generate transfer plans failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, ccErr,
			ctx.Kit.Rid)
		ctx.RespAutoError(ccErr)
		return
	}

	addModuleIDs := make([]int64, 0)
	removeModuleIDs := make([]int64, 0)
	for _, plan := range transferPlans {
		addModuleIDs = append(addModuleIDs, plan.ToAddToModules...)
		removeModuleIDs = append(removeModuleIDs, plan.ToRemoveFromModules...)
	}

	// get to remove service instances
	moduleHostSrvInstMap := make(map[int64]map[int64][]metadata.ServiceInstance)
	if len(removeModuleIDs) > 0 {
		listSrvInstOption := &metadata.ListServiceInstanceOption{
			BusinessID: bizID,
			HostIDs:    option.HostIDs,
			ModuleIDs:  removeModuleIDs,
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}

		srvInstResult, ccErr := s.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			listSrvInstOption)
		if ccErr != nil {
			blog.ErrorJSON("list service instance failed, bizID: %s, option: %s, err: %s, rid: %s", bizID,
				listSrvInstOption, ccErr, ctx.Kit.Rid)
			ctx.RespAutoError(ccErr)
			return
		}
		for _, item := range srvInstResult.Info {
			if _, exist := moduleHostSrvInstMap[item.ModuleID]; !exist {
				moduleHostSrvInstMap[item.ModuleID] = make(map[int64][]metadata.ServiceInstance, 0)
			}
			if _, exist := moduleHostSrvInstMap[item.ModuleID][item.HostID]; !exist {
				moduleHostSrvInstMap[item.ModuleID][item.HostID] = make([]metadata.ServiceInstance, 0)
			}
			moduleHostSrvInstMap[item.ModuleID][item.HostID] = append(moduleHostSrvInstMap[item.ModuleID][item.HostID],
				item)
		}
	}

	moduleServiceTemplateMap := make(map[int64]metadata.ServiceTemplateDetail)
	if len(addModuleIDs) > 0 {
		// get add to modules
		modules, ccErr := s.getModules(ctx.Kit, bizID, addModuleIDs)
		if ccErr != nil {
			blog.Errorf("get modules failed, err: %v, bizID: %d, module ids: %+v, rid: %s", ccErr, bizID, addModuleIDs,
				ctx.Kit.Rid)
			ctx.RespAutoError(ccErr)
			return
		}

		// get service template related to add modules
		serviceTemplateIDs := make([]int64, 0)
		for _, module := range modules {
			if module.ServiceTemplateID == common.ServiceTemplateIDNotSet {
				continue
			}
			serviceTemplateIDs = append(serviceTemplateIDs, module.ServiceTemplateID)
		}

		serviceTemplateDetails, ccErr := s.CoreAPI.CoreService().Process().ListServiceTemplateDetail(ctx.Kit.Ctx,
			ctx.Kit.Header, bizID, serviceTemplateIDs...)
		if ccErr != nil {
			blog.Errorf("list service template detail failed, bizID: %s, option: %s, err: %s, rid: %s", bizID,
				serviceTemplateIDs, ccErr, ctx.Kit.Rid)
			ctx.RespAutoError(ccErr)
			return
		}

		serviceTemplateMap := make(map[int64]metadata.ServiceTemplateDetail)
		for _, templateDetail := range serviceTemplateDetails.Info {
			serviceTemplateMap[templateDetail.ServiceTemplate.ID] = templateDetail
		}
		for _, module := range modules {
			templateDetail, exist := serviceTemplateMap[module.ServiceTemplateID]
			if exist {
				moduleServiceTemplateMap[module.ModuleID] = templateDetail
			}
		}
	}

	previews := make([]metadata.HostTransferPreview, 0)
	for _, plan := range transferPlans {
		preview := metadata.HostTransferPreview{
			HostID:              plan.HostID,
			FinalModules:        plan.FinalModules,
			ToRemoveFromModules: make([]metadata.RemoveFromModuleInfo, 0),
			ToAddToModules:      make([]metadata.AddToModuleInfo, 0),
			HostApplyPlan:       plan.HostApplyPlan,
		}

		for _, moduleID := range plan.ToRemoveFromModules {
			removeInfo := metadata.RemoveFromModuleInfo{
				ModuleID:         moduleID,
				ServiceInstances: make([]metadata.ServiceInstance, 0),
			}
			hostSrvInstMap, exist := moduleHostSrvInstMap[moduleID]
			if !exist {
				continue
			}

			serviceInstances, exist := hostSrvInstMap[moduleID]
			if exist {
				removeInfo.ServiceInstances = append(removeInfo.ServiceInstances, serviceInstances...)
			}
			preview.ToRemoveFromModules = append(preview.ToRemoveFromModules, removeInfo)
		}

		for _, moduleID := range plan.ToAddToModules {
			addInfo := metadata.AddToModuleInfo{
				ModuleID:        moduleID,
				ServiceTemplate: nil,
			}
			serviceTemplateDetail, exist := moduleServiceTemplateMap[moduleID]
			if exist {
				addInfo.ServiceTemplate = &serviceTemplateDetail
			}
			preview.ToAddToModules = append(preview.ToAddToModules, addInfo)
		}
		previews = append(previews, preview)
	}
	ctx.RespEntity(previews)
	return
}

func (s *Service) validateTransferHostWithAutoClearServiceInstanceOption(kit *rest.Kit, bizID int64,
	option metadata.TransferHostWithAutoClearServiceInstanceOption) errors.CCErrorCoder {

	if len(option.HostIDs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_host_ids")
	}

	if len(option.RemoveFromModules) == 0 && len(option.AddToModules) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "remove_from_modules or add_to_modules")
	}

	if option.DefaultInternalModule != 0 && len(option.AddToModules) != 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "add_to_modules & default_internal_module")
	}

	if len(option.AddToModules) != 0 {
		return s.validateModules(kit, bizID, option.AddToModules, "add_to_modules")
	}

	if option.DefaultInternalModule != 0 {
		return s.validateModules(kit, bizID, []int64{option.DefaultInternalModule}, "default_internal_module")
	}

	return nil
}
