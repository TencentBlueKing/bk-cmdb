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
	"sort"
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
		blog.V(7).Infof("parse bizID from url failed, bizID: %s, err: %+v, rid: %s", bizIDStr, err, ctx.Kit.Rid)
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

	transferPlans, hostIDs, err := s.preTransferPlans(ctx.Kit, option, bizID)
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

	transToInnerOpt, transToNormalPlans := s.parseTransferPlans(bizID, option.IsRemoveFromAll,
		len(option.RemoveFromModules) == 0, transferPlans, svcInstMap, option.Options.HostApplyTransPropertyRule.Changed)

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		return s.transferHostWithAutoClearServiceInstance(ctx.Kit, bizID, option, transToInnerOpt, transToNormalPlans,
			svcInstMap, hostIDs)
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
	return
}

// parseTransferPlans aggregate transfer plans into transfer to inner/normal module options by module ids and increment
func (s *Service) parseTransferPlans(bizID int64, isRemoveFromAll, isRemoveFromNone bool,
	transferPlans []metadata.HostTransferPlan, svcInstMap map[int64]map[int64][]metadata.ProcessInstanceDetail,
	changed bool) (*metadata.TransferHostToInnerModule, map[string]*metadata.HostsModuleRelation) {

	// when hosts are removed from all or no current modules, the plans are the same, we use a special key for this case
	const sameKey = "same"
	transferToInnerHostIDs := make([]int64, 0)
	var innerModuleID int64
	transferToNormalPlans := make(map[string]*metadata.HostsModuleRelation)
	for _, plan := range transferPlans {
		// do not need to transfer, skip
		if len(plan.ToAddToModules) == 0 && len(plan.ToRemoveFromModules) == 0 {
			delete(svcInstMap, plan.HostID)
			continue
		}

		// transfer to inner modules
		if plan.IsTransferToInnerModule {
			innerModuleID = plan.FinalModules[0]
			transferToInnerHostIDs = append(transferToInnerHostIDs, plan.HostID)
			if len(svcInstMap[plan.HostID]) != 0 {
				delete(svcInstMap, plan.HostID)
			}
			continue
		}

		var transKey string
		var isIncrement bool
		var moduleIDs []int64

		if len(plan.ToRemoveFromModules) == 0 {
			isIncrement = true
			moduleIDs = plan.ToAddToModules
		} else {
			isIncrement = false
			moduleIDs = plan.FinalModules
		}

		if isRemoveFromAll || isRemoveFromNone {
			transKey = sameKey
		} else {
			// we use is increment and sorted module ids to aggregate hosts with the same transfer option
			sort.Slice(moduleIDs, func(i, j int) bool { return moduleIDs[i] < moduleIDs[j] })
			transKey = fmt.Sprintf("%v%v", isIncrement, moduleIDs)
		}

		if _, exists := transferToNormalPlans[transKey]; !exists {
			transferToNormalPlans[transKey] = &metadata.HostsModuleRelation{
				ApplicationID:            bizID,
				HostID:                   []int64{plan.HostID},
				DisableAutoCreateSvcInst: true,
				IsIncrement:              isIncrement,
				ModuleID:                 moduleIDs,
			}
		} else {
			transferToNormalPlans[transKey].HostID = append(transferToNormalPlans[transKey].HostID, plan.HostID)
		}

		// the mark here is to adapt to the host transfer scenario, the user does not need the attributes of the host
		// dimension to be automatically applied.
		transferToNormalPlans[transKey].DisableTransferHostAutoApply = !changed
		for moduleID := range svcInstMap[plan.HostID] {
			if !util.InArray(moduleID, plan.FinalModules) {
				delete(svcInstMap[plan.HostID], moduleID)
			}
		}
	}

	var transToInnerOpt *metadata.TransferHostToInnerModule
	if len(transferToInnerHostIDs) > 0 {
		transToInnerOpt = &metadata.TransferHostToInnerModule{ApplicationID: bizID, HostID: transferToInnerHostIDs,
			ModuleID: innerModuleID}
	}

	return transToInnerOpt, transferToNormalPlans
}

func (s *Service) transferHostWithAutoClearServiceInstance(
	kit *rest.Kit,
	bizID int64,
	option metadata.TransferHostWithAutoClearServiceInstanceOption,
	transToInnerOpt *metadata.TransferHostToInnerModule,
	transToNormalPlans map[string]*metadata.HostsModuleRelation,
	svcInstMap map[int64]map[int64][]metadata.ProcessInstanceDetail,
	hostIDs []int64) error {

	audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), option.HostIDs)
	if err := audit.WithPrevious(kit); err != nil {
		blog.Errorf("generate host transfer audit failed, err: %v, HostIDs: %+v, rid: %s", err, option.HostIDs, kit.Rid)
		return err
	}

	// hosts that are transferred to inner module or removed from all previous modules have the same destination
	// if there's no remove module, hosts are appended to same modules, can transfer together
	if transToInnerOpt != nil {
		res, err := s.CoreAPI.CoreService().Host().TransferToInnerModule(kit.Ctx, kit.Header, transToInnerOpt)
		if err != nil {
			blog.Errorf("transfer host failed, err: %v, res: %v, opt: %#v, rid: %s", err, res, transToInnerOpt, kit.Rid)
			return err
		}
	}

	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 20)
	wg := sync.WaitGroup{}
	for _, plan := range transToNormalPlans {
		if firstErr != nil {
			break
		}
		pipeline <- true
		wg.Add(1)
		go func(kit *rest.Kit, plan *metadata.HostsModuleRelation) {
			defer func() {
				<-pipeline
				wg.Done()
			}()

			// transfer hosts in 2 scenario, add to modules and transfer to other modules
			res, ccErr := s.CoreAPI.CoreService().Host().TransferToNormalModule(kit.Ctx, kit.Header, plan)
			if ccErr != nil {
				if firstErr == nil {
					firstErr = ccErr
				}
				blog.Errorf("transfer host failed, err: %v, res: %v, opt: %#v, rid: %s", ccErr, res, plan, kit.Rid)
				return
			}

		}(kit, plan)
	}
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}

	if err := s.upsertServiceInstance(kit, bizID, svcInstMap); err != nil {
		blog.Errorf("upsert service instance(%#v) failed, err: %v, svcInstMap: %#v, rid: %s", err, svcInstMap, kit.Rid)
		return err
	}
	// update host properties according to specified rules.
	err := s.updateHostApplyByRule(kit, option.Options.HostApplyTransPropertyRule, hostIDs)
	if err != nil {
		blog.Errorf("update host properties according to specified rules(%#v) failed, err: %v, rid: %s",
			option.Options.HostApplyTransPropertyRule, err, kit.Rid)
		return err
	}

	if err := audit.SaveAudit(kit); err != nil {
		blog.Errorf("save audit log failed, err: %v, HostIDs: %+v, rid: %s", err, option.HostIDs, kit.Rid)
		return err
	}

	return nil
}

// upsertServiceInstance create or update related service instance in host transfer option
func (s *Service) upsertServiceInstance(kit *rest.Kit, bizID int64,
	svcInstMap map[int64]map[int64][]metadata.ProcessInstanceDetail) error {

	moduleSvcInstMap := make(map[int64][]metadata.CreateServiceInstanceDetail)
	for hostID, moduleProcMap := range svcInstMap {
		for moduleID, processes := range moduleProcMap {
			moduleSvcInstMap[moduleID] = append(moduleSvcInstMap[moduleID], metadata.CreateServiceInstanceDetail{
				HostID:    hostID,
				Processes: processes,
			})
		}
	}

	wg := sync.WaitGroup{}
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 20)
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
			defer func() {
				<-pipeline
				wg.Done()
			}()

			instances := svrInstOpt.Instances

			total := len(instances)
			for start := 0; start < total; start += common.BKMaxUpdateOrCreatePageSize {
				// 这里需要进行分批处理，一次处理100个
				var tmpInstances []metadata.CreateServiceInstanceDetail
				if total-start >= common.BKMaxUpdateOrCreatePageSize {
					tmpInstances = instances[start : start+common.BKMaxUpdateOrCreatePageSize]
				} else {
					tmpInstances = instances[start:total]
				}

				svrInstOpt.Instances = tmpInstances

				_, ccErr := s.CoreAPI.ProcServer().Service().CreateServiceInstance(kit.Ctx, kit.Header, svrInstOpt)
				if ccErr != nil {
					if firstErr == nil {
						firstErr = ccErr
					}
					blog.Errorf("create service instances failed, option: %#v, err: %v, rid: %s", ccErr, svrInstOpt,
						kit.Rid)
					return
				}
			}
		}(svrInstOpt)
	}
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}
	return nil
}

func (s *Service) updateHostApplyByRule(kit *rest.Kit, rule metadata.HostApplyTransRules,
	hostIDs []int64) errors.CCErrorCoder {

	if !rule.Changed {
		return nil
	}

	attributeIDs := make([]int64, 0)
	for _, rule := range rule.FinalRules {
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

	planResult := make([]metadata.CreateHostApplyRuleOption, 0)
	for _, r := range rule.FinalRules {
		planResult = append(planResult, metadata.CreateHostApplyRuleOption{
			AttributeID:   r.AttributeID,
			PropertyValue: r.PropertyValue,
		})
	}
	attributes := make([]metadata.HostAttribute, 0)

	for _, rule := range planResult {
		attributes = append(attributes, metadata.HostAttribute{
			AttributeID:   rule.AttributeID,
			PropertyValue: rule.PropertyValue,
		})
	}
	if err := s.updateHostAttributes(kit, attributes, hostIDs); err != nil {
		blog.Errorf("update attributes failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}
func (s *Service) preTransferPlans(kit *rest.Kit, option metadata.TransferHostWithAutoClearServiceInstanceOption,
	bizID int64) ([]metadata.HostTransferPlan, []int64, errors.CCErrorCoder) {

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
		return nil, nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
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
		return nil, nil, ccErr
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
			return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "default_internal_module")
		}
		defaultInternalModuleID = option.DefaultInternalModule
	}

	hostIDs := make([]int64, 0)
	transferPlans := make([]metadata.HostTransferPlan, 0)
	for hostID, currentInModules := range hostModulesIDMap {
		// if host is currently in inner module and is going to append to another module, transfer to that module
		if len(option.RemoveFromModules) == 0 {
			option.RemoveFromModules = innerModuleIDs
		}

		transferPlan := generateTransferPlan(currentInModules, option.RemoveFromModules, option.AddToModules,
			defaultInternalModuleID)
		transferPlan.HostID = hostID
		hostIDs = append(hostIDs, hostID)
		// check module compatibility
		finalModuleCount := len(transferPlan.FinalModules)
		for _, moduleID := range transferPlan.FinalModules {
			if _, exists := innerModuleIDMap[moduleID]; !exists {
				continue
			}
			if finalModuleCount != 1 {
				return nil, nil, kit.CCError.CCError(common.CCErrHostTransferFinalModuleConflict)
			}
			transferPlan.IsTransferToInnerModule = true
		}
		transferPlans = append(transferPlans, transferPlan)
	}

	return transferPlans, hostIDs, nil
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

func (s *Service) generateTransferPlans(kit *rest.Kit, bizID int64,
	option metadata.TransferHostWithAutoClearServiceInstanceOption) (
	[]metadata.HostTransferPlan, errors.CCErrorCoder) {
	transferPlans, _, err := s.preTransferPlans(kit, option, bizID)
	if err != nil {
		return nil, err
	}

	// generate host apply plans
	return s.generateHostApplyPlans(kit, bizID, transferPlans)
}

func (s *Service) generateHostApplyPlans(kit *rest.Kit, bizID int64, plans []metadata.HostTransferPlan) (
	[]metadata.HostTransferPlan, errors.CCErrorCoder) {

	if len(plans) == 0 {
		return plans, nil
	}

	// get final modules' host apply rules
	finalModuleIDs := make([]int64, 0)
	for _, item := range plans {
		finalModuleIDs = append(finalModuleIDs, item.FinalModules...)
	}

	rules, err := s.getRulesPriorityFromTemplate(kit, finalModuleIDs, bizID)
	if err != nil {
		blog.Errorf("get module rule failed, err: %v, rid: %s", err, kit.Rid)
		return plans, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	// generate host apply plans only generate new module
	hostModules := make([]metadata.Host2Modules, 0)
	for _, item := range plans {
		host2Module := metadata.Host2Modules{
			HostID:    item.HostID,
			ModuleIDs: make([]int64, 0),
		}
		for _, moduleID := range item.ToAddToModules {
			host2Module.ModuleIDs = append(host2Module.ModuleIDs, moduleID)
		}
		hostModules = append(hostModules, host2Module)
	}

	planOpt := metadata.HostApplyPlanOption{
		Rules:       rules,
		HostModules: hostModules,
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

func (s *Service) getRemovedServiceInstance(ctx *rest.Contexts, bizID int64, removeModuleIDs []int64,
	option metadata.TransferHostWithAutoClearServiceInstanceOption) (map[int64]map[int64][]metadata.ServiceInstance,
	error) {
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
		blog.Errorf("list service instance failed, bizID: %d, option: %s, err: %v, rid: %s", bizID,
			listSrvInstOption, ccErr, ctx.Kit.Rid)
		return nil, ccErr
	}
	moduleHostSrvInstMap := make(map[int64]map[int64][]metadata.ServiceInstance)

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
	return moduleHostSrvInstMap, nil
}

func (s *Service) getModuleServiceTemplate(ctx *rest.Contexts, bizID int64, addModuleIDs []int64) (
	map[int64]metadata.ServiceTemplateDetail, error) {
	// get add to modules
	modules, ccErr := s.getModules(ctx.Kit, bizID, addModuleIDs)
	if ccErr != nil {
		blog.Errorf("get modules failed, bizID: %d, module ids: %+v, err: %v, rid: %s", bizID, addModuleIDs, ccErr,
			ctx.Kit.Rid)
		return nil, ccErr
	}

	// get service template related to add modules
	serviceTemplateIDs := make([]int64, 0)
	for _, module := range modules {
		if module.ServiceTemplateID == common.ServiceTemplateIDNotSet {
			continue
		}
		serviceTemplateIDs = append(serviceTemplateIDs, module.ServiceTemplateID)
	}
	moduleServiceTemplateMap := make(map[int64]metadata.ServiceTemplateDetail)

	if len(serviceTemplateIDs) > 0 {
		serviceTemplateDetails, ccErr := s.CoreAPI.CoreService().Process().ListServiceTemplateDetail(ctx.Kit.Ctx,
			ctx.Kit.Header, bizID, serviceTemplateIDs...)
		if ccErr != nil {
			blog.Errorf("list service template detail failed, bizID: %d, option: %s, err: %s, rid: %s", bizID,
				serviceTemplateIDs, ccErr, ctx.Kit.Rid)
			return nil, ccErr
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
	return moduleServiceTemplateMap, nil
}

// 接口请求参数跟转移是一致的
// 主机从模块删除时提供了将要删除的服务实例信息
// 主机添加到新模块时，提供了模块对应的服务模板（如果有）

// TransferHostWithAutoClearServiceInstancePreview generate a preview of changes for operation
func (s *Service) TransferHostWithAutoClearServiceInstancePreview(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.V(7).Infof("parse bizID from url failed, bizID: %s, err: %+v, rid: %s", bizIDStr, err, ctx.Kit.Rid)
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

	transferPlans, ccErr := s.generateTransferPlans(ctx.Kit, bizID, option)
	if ccErr != nil {
		blog.Errorf("generate plans fail, bizID: %s, option: %+v, err: %v, rid: %s", bizID, option, ccErr, ctx.Kit.Rid)
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
		moduleHostSrvInstMap, err = s.getRemovedServiceInstance(ctx, bizID, removeModuleIDs, option)
		if err != nil {
			ctx.RespAutoError(ccErr)
			return
		}
	}

	moduleServiceTemplateMap := make(map[int64]metadata.ServiceTemplateDetail)
	if len(addModuleIDs) > 0 {

		moduleServiceTemplateMap, err = s.getModuleServiceTemplate(ctx, bizID, addModuleIDs)
		if err != nil {
			ctx.RespAutoError(ccErr)
			return
		}
	}

	previews := getPreviewsResult(transferPlans, moduleServiceTemplateMap, moduleHostSrvInstMap)
	ctx.RespEntity(previews)
	return
}

func getPreviewsResult(transferPlans []metadata.HostTransferPlan,
	moduleServiceTemplateMap map[int64]metadata.ServiceTemplateDetail,
	moduleHostSrvInstMap map[int64]map[int64][]metadata.ServiceInstance) []metadata.HostTransferPreview {
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
				preview.ToRemoveFromModules = append(preview.ToRemoveFromModules, removeInfo)
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
	return previews
}

func (s *Service) getRulesPriorityFromTemplate(kit *rest.Kit, moduleIDs []int64, bizID int64) (
	[]metadata.HostApplyRule, error) {

	moduleRes, err := s.getModuleRelateHostApply(kit, bizID, moduleIDs, nil)
	if err != nil {
		return nil, err
	}

	// 1.过滤出需要查询主机应用规则的模版id和模块id
	enabledModuleIDs := make([]int64, 0)
	enabledModuleIDMap := make(map[int64]bool)
	tempToModMap := make(map[int64][]int64)
	srvTmpIDs := make([]int64, 0)
	for _, module := range moduleRes {
		if module.ServiceTemplateID != 0 {
			tempToModMap[module.ServiceTemplateID] = append(tempToModMap[module.ServiceTemplateID], module.ModuleID)
			srvTmpIDs = append(srvTmpIDs, module.ServiceTemplateID)

			if module.HostApplyEnabled {
				enabledModuleIDMap[module.ModuleID] = true
			}
			continue
		}

		if module.HostApplyEnabled {
			enabledModuleIDs = append(enabledModuleIDs, module.ModuleID)
		}
	}

	enableSrvTemplateIDs := make([]int64, 0)
	if len(srvTmpIDs) != 0 {
		srvTempStatus, err := s.getSrvTemplateApplyStatus(kit, bizID, srvTmpIDs)
		if err != nil {
			blog.Errorf("get service template host apply status failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		for templateID, status := range srvTempStatus {
			if status {
				enableSrvTemplateIDs = append(enableSrvTemplateIDs, templateID)
				continue
			}
			for _, moduleID := range tempToModMap[templateID] {
				if enabledModuleIDMap[moduleID] {
					enabledModuleIDs = append(enabledModuleIDs, moduleID)
				}
			}
		}
	}

	// 2.查询有模版并且模版开启主机自动应用的规则
	rules := make([]metadata.HostApplyRule, 0)
	if len(enableSrvTemplateIDs) != 0 {
		srvTemplateRules, err := s.findSrvTemplateRule(kit, bizID, enableSrvTemplateIDs)
		if err != nil {
			blog.Errorf("list service template host apply rule failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		for _, rule := range srvTemplateRules {
			moduleIDs, exist := tempToModMap[rule.ServiceTemplateID]
			if !exist {
				continue
			}

			for _, moduleID := range moduleIDs {
				rule.ModuleID = moduleID
				rules = append(rules, rule)
			}
		}
	}

	// 3.查询没有模版，以及有模版但是模版没有开启主机自动应用的模块的规则
	if len(enabledModuleIDs) != 0 {
		moduleRules, err := s.getEnabledModuleRules(kit, bizID, enabledModuleIDs)
		if err != nil {
			blog.Errorf("get module host apply rule failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		if len(moduleRules) != 0 {
			rules = append(rules, moduleRules...)
		}
	}

	return rules, nil
}
