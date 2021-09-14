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
// 1. 将主机 bk_host_ids 从 remove_from_modules 指定的拓扑节点移除
// 2. 移入到 add_to_modules 指定的模块中
// 3. 自动删除主机在移除模块下的服务实例
// 4. 自动添加主机在新模块上的服务实例
// note:
// - 不允许 remove_from_modules 和 add_to_modules 同时为空
// - bk_host_ids 不允许为空
// - 如果 is_remove_from_all 为true，则接口行为是：覆盖更新
// - 如果 remove_from_modules 没有指定，仅仅是增量更新，无移除操作
// - 如果 add_to_modules 没有指定，主机将仅仅从 remove_from_modules 指定的模块中移除
// - 如果 add_to_modules 是空闲机/故障机/待回收模块中的一个，必须显式指定 remove_from_modules(可指定成业务节点),
// 否则报主机不能属于互斥模块错误
// - 如果 add_to_modules 是普通模块，主机当前数据空闲机/故障机/待回收模块中的一个，必须显式指定 remove_from_modules
// (可指定成业务节点), 否则报主机不能属于互斥模块错误
// - 模块同时出现在 add_to_modules 和 remove_from_modules 时，不会导致对应的服务实例被删除然后重新添加
// - 主机从 remove_from_modules 移除后如果不再属于其它模块， 默认转移到空闲机模块，default_internal_module 可以指定为空闲机/故障机/
// 待回收模块中的一个，表示主机移除全部模块后默认转移到的模块
// - 不允许 add_to_modules 和 default_internal_module 同时指定
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

	if ccErr := s.validateTransferHostWithAutoClearServiceInstanceOption(ctx.Kit, bizID, &option); ccErr != nil {
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

	// parse service instances to map[hostID->map[moduleID->processes]], skip those that do not belong to host modules
	svcInstMap := make(map[int64]map[int64][]metadata.ProcessInstanceDetail)
	svrInstOp := append(option.Options.ServiceInstanceOptions.Created, option.Options.ServiceInstanceOptions.Updated...)
	for _, svcInst := range svrInstOp {
		if _, exists := svcInstMap[svcInst.HostID]; !exists {
			svcInstMap[svcInst.HostID] = make(map[int64][]metadata.ProcessInstanceDetail)
		}
		svcInstMap[svcInst.HostID][svcInst.ModuleID] = svcInst.Processes
	}

	transferToInnerHostIDs, transferToNormalHostIDs := make([]int64, 0), make([]int64, 0)
	var innerModuleID int64
	normalModuleIDs := make([]int64, 0)
	separatePlans := make([]metadata.HostTransferPlan, 0)
	for _, plan := range transferPlans {
		if len(plan.ToAddToModules) == 0 && len(plan.ToRemoveFromModules) == 0 {
			delete(svcInstMap, plan.HostID)
			continue
		}
		if plan.IsTransferToInnerModule {
			innerModuleID = plan.FinalModules[0]
			transferToInnerHostIDs = append(transferToInnerHostIDs, plan.HostID)
			if len(svcInstMap[plan.HostID]) != 0 {
				delete(svcInstMap, plan.HostID)
			}
		} else {
			if option.IsRemoveFromAll || len(option.RemoveFromModules) == 0 {
				normalModuleIDs = plan.FinalModules
				transferToNormalHostIDs = append(transferToNormalHostIDs, plan.HostID)
			}
			separatePlans = append(separatePlans, plan)
			for moduleID := range svcInstMap[plan.HostID] {
				if !util.InArray(moduleID, plan.FinalModules) {
					delete(svcInstMap[plan.HostID], moduleID)
				}
			}
		}
	}

	transferResult := make([]metadata.HostTransferResult, 0)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), option.HostIDs)
		if err := audit.WithPrevious(ctx.Kit); err != nil {
			blog.Errorf("get prev module host config for audit failed, err: %v, HostIDs: %+v, rid: %s", err,
				option.HostIDs, ctx.Kit.Rid)
			return err
		}

		// hosts that are transferred to inner module or removed from all previous modules have the same destination
		// if there's no remove module, hosts are appended to same modules, can transfer together
		if len(transferToInnerHostIDs) > 0 {
			transferOpt := &metadata.TransferHostToInnerModule{
				ApplicationID: bizID,
				HostID:        transferToInnerHostIDs,
				ModuleID:      innerModuleID,
			}
			res, ccErr := s.CoreAPI.CoreService().Host().TransferToInnerModule(ctx.Kit.Ctx, ctx.Kit.Header, transferOpt)
			if ccErr != nil {
				blog.Errorf("transfer host failed, err: %v, res: %v, option: %#v, rid: %s", ccErr, res, transferOpt,
					ctx.Kit.Rid)
				return ccErr
			}
		}

		if len(transferToNormalHostIDs) > 0 && (option.IsRemoveFromAll || len(option.RemoveFromModules) == 0) {
			transferOpt := &metadata.HostsModuleRelation{
				ApplicationID:         bizID,
				HostID:                transferToNormalHostIDs,
				NeedAutoCreateSvcInst: false,
			}
			if option.IsRemoveFromAll {
				transferOpt.IsIncrement = false
				transferOpt.ModuleID = normalModuleIDs
			} else {
				transferOpt.IsIncrement = true
				transferOpt.ModuleID = option.AddToModules
			}
			res, ccErr := s.CoreAPI.CoreService().Host().TransferToNormalModule(ctx.Kit.Ctx, ctx.Kit.Header,
				transferOpt)
			if ccErr != nil {
				blog.Errorf("transfer host failed, err: %v, res: %v, option: %#v, rid: %s", ccErr, res, transferOpt,
					ctx.Kit.Rid)
				return ccErr
			}
		}

		var firstErr errors.CCErrorCoder
		pipeline := make(chan bool, 300)
		wg := sync.WaitGroup{}
		for _, plan := range separatePlans {
			if firstErr != nil {
				break
			}
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

				// transfer hosts in 2 scenario, add to modules and transfer to other modules
				transferOpt := &metadata.HostsModuleRelation{
					ApplicationID:         bizID,
					HostID:                []int64{plan.HostID},
					NeedAutoCreateSvcInst: false,
				}
				if len(plan.ToRemoveFromModules) == 0 {
					transferOpt.IsIncrement = true
					transferOpt.ModuleID = plan.ToAddToModules
				} else {
					transferOpt.IsIncrement = false
					transferOpt.ModuleID = plan.FinalModules
				}

				transRes, ccErr := s.CoreAPI.CoreService().Host().TransferToNormalModule(ctx.Kit.Ctx, ctx.Kit.Header,
					transferOpt)
				if ccErr != nil {
					blog.Errorf("transfer host failed, err: %v, res: %v, option: %#v, rid: %s", ccErr, transRes,
						transferOpt, ctx.Kit.Rid)
					return
				}

			}(plan)
		}
		wg.Wait()
		if firstErr != nil {
			return firstErr
		}

		// create or update related service instance
		moduleSvcInstMap := make(map[int64][]metadata.CreateServiceInstanceDetail)
		for hostID, moduleProcMap := range svcInstMap {
			for moduleID, processes := range moduleProcMap {
				moduleSvcInstMap[moduleID] = append(moduleSvcInstMap[moduleID], metadata.CreateServiceInstanceDetail{
					HostID:    hostID,
					Processes: processes,
				})
			}
		}

		wg = sync.WaitGroup{}
		for moduleID, svcInst := range moduleSvcInstMap {
			if firstErr != nil {
				break
			}
			pipeline <- true
			wg.Add(1)

			svrInstOpt := &metadata.CreateServiceInstanceInput{
				BizID:     bizID,
				ModuleID:  moduleID,
				Instances: svcInst,
			}
			go func(svrInstOpt *metadata.CreateServiceInstanceInput) {
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

				_, ccErr = s.CoreAPI.ProcServer().Service().CreateServiceInstance(ctx.Kit.Ctx,
					ctx.Kit.Header, svrInstOpt)
				if ccErr != nil {
					blog.ErrorJSON("create service instance failed, err: %v, option: %#v, rid: %s", ccErr, svrInstOpt,
						ctx.Kit.Rid)
					return
				}
			}(svrInstOpt)
		}
		wg.Wait()
		if firstErr != nil {
			return firstErr
		}

		// update host by host apply rule conflict resolvers
		err = s.updateHostByHostApplyConflictResolvers(ctx.Kit, option.Options.HostApplyConflictResolvers)
		if err != nil {
			blog.Errorf("update host by host apply rule conflict resolvers(%#v) failed, err: %v, rid: %s",
				option.Options.HostApplyConflictResolvers, err, ctx.Kit.Rid)
			return err
		}

		if err := audit.SaveAudit(ctx.Kit); err != nil {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, save audit log failed, err: %s, " +
				"HostIDs: %+v, rid: %s", err.Error(), option.HostIDs, ctx.Kit.Rid)
			return err
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

func (s *Service) updateHostByHostApplyConflictResolvers(kit *rest.Kit,
	resolvers []metadata.HostApplyConflictResolver) errors.CCErrorCoder {

	if len(resolvers) == 0 {
		return nil
	}

	attributeIDs := make([]int64, 0)
	for _, rule := range resolvers {
		attributeIDs = append(attributeIDs, rule.AttributeID)
	}
	attCond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID, common.BKPropertyIDField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: attributeIDs,
			},
		},
	}

	attrRes, err := s.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDHost, attCond)
	if err != nil {
		blog.Errorf("read model attr failed, err: %v, attrCond: %#v, rid: %s", err, attCond, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	attrMap := make(map[int64]string)
	for _, attr := range attrRes.Info {
		attrMap[attr.ID] = attr.PropertyID
	}

	hostAttrMap := make(map[int64]map[string]interface{})
	for _, rule := range resolvers {
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
		_, err := s.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header,
			common.BKInnerObjIDHost, updateOption)
		if err != nil {
			blog.ErrorJSON("update host failed, option: %#v, err: %v, rid: %s", updateOption, err, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
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
		Page:          metadata.BasePage{Limit: common.BKNoLimit},
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}
	hostModuleResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(kit.Ctx, kit.Header, hostModuleOption)
	if err != nil {
		blog.ErrorJSON("get host module relation failed, option: %s, err: %s, rid: %s", hostModuleOption, err, rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	hostModulesIDMap := make(map[int64][]int64)
	for _, item := range hostModuleResult.Info {
		if _, exist := hostModulesIDMap[item.HostID]; !exist {
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
	innerModuleIDs := make([]int64, 0)
	defaultInternalModuleID := int64(0)
	for _, module := range innerModules {
		innerModuleIDMap[module.ModuleID] = struct{}{}
		innerModuleIDs = append(innerModuleIDs, module.ModuleID)
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
		defaultInternalModuleID = option.DefaultInternalModule
	}

	transferPlans := make([]metadata.HostTransferPlan, 0)
	for hostID, currentInModules := range hostModulesIDMap {
		// if host is currently in inner module and is going to append to another module, transfer to that module
		if len(option.RemoveFromModules) == 0 {
			option.RemoveFromModules = innerModuleIDs
		}

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
	return s.generateHostApplyPlans(kit, bizID, transferPlans, option.Options.HostApplyConflictResolvers)
}

func (s *Service) generateHostApplyPlans(kit *rest.Kit, bizID int64, plans []metadata.HostTransferPlan,
	resolvers []metadata.HostApplyConflictResolver) ([]metadata.HostTransferPlan, errors.CCErrorCoder) {

	if len(plans) == 0 {
		return plans, nil
	}

	// get final modules' host apply rules
	finalModuleIDs := make([]int64, 0)
	for _, item := range plans {
		finalModuleIDs = append(finalModuleIDs, item.FinalModules...)
	}

	ruleOpt := metadata.ListHostApplyRuleOption{
		ModuleIDs: finalModuleIDs,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header, bizID, ruleOpt)
	if ccErr != nil {
		blog.Errorf("list apply rule failed, bizID: %s, option: %#v, err: %s, rid: %s", bizID, ruleOpt, ccErr, kit.Rid)
		return plans, ccErr
	}

	// get modules that enabled host apply
	moduleCondition := metadata.QueryCondition{
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Fields: []string{common.BKModuleIDField},
		Condition: map[string]interface{}{
			common.BKModuleIDField:       map[string]interface{}{common.BKDBIN: finalModuleIDs},
			common.HostApplyEnabledField: true,
		},
	}
	enabledModules, err := s.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, &moduleCondition)
	if err != nil {
		blog.ErrorJSON("get apply enabled modules failed, filter: %s, err: %s, rid: %s", moduleCondition, err, kit.Rid)
		return plans, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	enableModuleMap := make(map[int64]bool)
	for _, item := range enabledModules.Info {
		moduleID, err := util.GetInt64ByInterface(item[common.BKModuleIDField])
		if err != nil {
			blog.ErrorJSON("parse module from db failed, module: %s, err: %s, rid: %s", item, err, kit.Rid)
			return plans, kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		enableModuleMap[moduleID] = true
	}

	// generate host apply plans
	hostModules := make([]metadata.Host2Modules, 0)
	for _, item := range plans {
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

	planOpt := metadata.HostApplyPlanOption{
		Rules:             rules.Info,
		HostModules:       hostModules,
		ConflictResolvers: resolvers,
	}

	hostApplyPlanResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(kit.Ctx, kit.Header, bizID,
		planOpt)
	if ccErr != nil {
		blog.Errorf("generate apply plan failed, biz: %d, opt: %#v, err: %v, rid: %s", bizID, planOpt, ccErr, kit.Rid)
		return plans, ccErr
	}

	hostApplyPlanMap := make(map[int64]metadata.OneHostApplyPlan)
	for _, item := range hostApplyPlanResult.Plans {
		hostApplyPlanMap[item.HostID] = item
	}
	for index, transferPlan := range plans {
		if applyPlan, ok := hostApplyPlanMap[transferPlan.HostID]; ok {
			plans[index].HostApplyPlan = applyPlan
		}
	}

	return plans, nil
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
	// 主机将会被移出的模块列表，包括当前所在模块里在移出模块且不在新增模块中的模块
	realRemoveModuleMap := make(map[int64]struct{})
	finalModules := make([]int64, 0)
	finalModuleMap := make(map[int64]struct{})
	currentModuleMap := make(map[int64]struct{})
	for _, moduleID := range currentIn {
		currentModuleMap[moduleID] = struct{}{}
		if _, exists := finalModuleMap[moduleID]; exists {
			continue
		}
		if _, exists := removeFromModuleMap[moduleID]; exists {
			if _, exists := realRemoveModuleMap[moduleID]; !exists {
				realRemoveModuleMap[moduleID] = struct{}{}
			}
			continue
		}
		finalModuleMap[moduleID] = struct{}{}
		finalModules = append(finalModules, moduleID)
	}

	// 主机将会被新加到的模块列表，包括新增模块里不在当前模块的模块
	realAddModules := make([]int64, 0)
	for _, moduleID := range addTo {
		if _, exists := finalModuleMap[moduleID]; exists {
			continue
		}
		finalModuleMap[moduleID] = struct{}{}
		finalModules = append(finalModules, moduleID)
		delete(realRemoveModuleMap, moduleID)
		if _, exists := currentModuleMap[moduleID]; exists {
			continue
		}
		realAddModules = append(realAddModules, moduleID)
	}

	realRemoveModules := make([]int64, 0)
	for moduleID := range realRemoveModuleMap {
		realRemoveModules = append(realRemoveModules, moduleID)
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

	if ccErr := s.validateTransferHostWithAutoClearServiceInstanceOption(ctx.Kit, bizID, &option); ccErr != nil {
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

		if len(serviceTemplateIDs) > 0 {
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

			serviceInstances, exist := hostSrvInstMap[plan.HostID]
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
	option *metadata.TransferHostWithAutoClearServiceInstanceOption) errors.CCErrorCoder {

	if option == nil {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_host_ids")
	}

	if len(option.HostIDs) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_host_ids")
	}

	if option.IsRemoveFromAll {
		moduleFilter := &metadata.DistinctFieldOption{
			TableName: common.BKTableNameModuleHostConfig,
			Field:     common.BKModuleIDField,
			Filter: map[string]interface{}{
				common.BKAppIDField:  bizID,
				common.BKHostIDField: map[string]interface{}{common.BKDBIN: option.HostIDs},
			},
		}
		rawModuleIDs, ccErr := s.CoreAPI.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, moduleFilter)
		if ccErr != nil {
			blog.Errorf("get host module ids failed, err: %v, filter: %#v, rid: %s", ccErr, moduleFilter, kit.Rid)
			return ccErr
		}

		moduleIDs, err := util.SliceInterfaceToInt64(rawModuleIDs)
		if err != nil {
			blog.Errorf("parse module ids(%#v) failed, err: %v, rid: %s", rawModuleIDs, err, kit.Rid)
			return ccErr
		}

		option.RemoveFromModules = moduleIDs
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
