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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
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
func (s *Service) TransferHostWithAutoClearServiceInstance(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	option := metadata.TransferHostWithAutoClearServiceInstanceOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("TransferHostWithAutoClearServiceInstance failed, parse request body failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(option.HostIDs) == 0 {
		err := srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "bk_host_ids")
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	if option.RemoveFromNode == nil && option.AddToModules == nil {
		err := srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "add_to_modules")
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.V(7).Infof("parse bizID from url failed, bizID: %s, err: %+v, rid: %s", bizIDStr, srvData.rid)
		err := srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	if len(option.AddToModules) != 0 {
		if ccErr := s.validateModules(srvData, bizID, option.AddToModules, "add_to_modules"); ccErr != nil {
			_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: ccErr})
			return
		}
	}
	if option.DefaultInternalModule != 0 {
		if ccErr := s.validateModules(srvData, bizID, []int64{option.DefaultInternalModule}, "default_internal_module"); ccErr != nil {
			_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: ccErr})
			return
		}
	}

	transferPlans, err := s.generateTransferPlans(srvData, bizID, option)
	if err != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, generateTransferPlans failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	type HostTransferResult struct {
		HostID  int64  `json:"bk_host_id"`
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	transferResult := make([]HostTransferResult, 0)
	var firstErr errors.CCErrorCoder
	for _, plan := range transferPlans {
		ccErr := s.runTransferPlans(srvData, bizID, plan)
		hostTransferResult := HostTransferResult{
			HostID: plan.HostID,
		}
		if ccErr == nil {
			// create or update related service instance
			for _, item := range option.Options.ServiceInstanceOptions {
				if item.HostID != plan.HostID {
					continue
				}
				if util.InArray(item.ModuleID, plan.FinalModules) == false {
					continue
				}
				if ccErr = s.createOrUpdateServiceInstance(srvData, bizID, plan.HostID, item); ccErr != nil {
					break
				}
			}
		}
		if ccErr != nil {
			hostTransferResult.Code = ccErr.GetCode()
			hostTransferResult.Message = ccErr.Error()
			if firstErr == nil {
				firstErr = ccErr
			}
		}
		transferResult = append(transferResult, hostTransferResult)
	}
	if firstErr != nil {
		response := metadata.Response{
			BaseResp: metadata.BaseResp{
				Result:      false,
				Code:        firstErr.GetCode(),
				ErrMsg:      firstErr.Error(),
				Permissions: nil,
			},
			Data: transferResult,
		}
		_ = resp.WriteEntity(response)
		return
	}

	_ = resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     transferResult,
	})
	return
}

func (s *Service) createOrUpdateServiceInstance(srvData *srvComm, bizID int64, hostID int64, serviceInstanceOption metadata.CreateServiceInstanceOption) errors.CCErrorCoder {
	rid := srvData.rid
	module, err := s.getModule(srvData, bizID, serviceInstanceOption.ModuleID)
	if err != nil {
		blog.Errorf("createOrUpdateServiceInstance failed, get module failed, bizID: %d, moduleID: %d, err: %s, rid: %s",
			bizID, module, err.Error(), rid)
		return err
	}
	if module.ServiceTemplateID == common.ServiceTemplateIDNotSet {
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
		result, err := s.CoreAPI.ProcServer().Service().CreateServiceInstance(srvData.ctx, srvData.header, input)
		if err != nil {
			blog.ErrorJSON("createServiceInstance failed, http failed, option: %s, err: %s, rid: %s", input, err.Error(), srvData.rid)
			return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if result.Result == false || result.Code != 0 {
			blog.ErrorJSON("createServiceInstance failed, option: %s, response: %s, rid: %s", input, result, srvData.rid)
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
		relationResult, err := s.CoreAPI.CoreService().Process().ListProcessInstanceRelation(srvData.ctx, srvData.header, relationOption)
		if err != nil {
			blog.ErrorJSON("update process instance failed, list process relation failed, option: %s, err: %s, rid: %s", relationOption, err.Error(), srvData.rid)
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
		updateProcessOption := map[string]interface{}{
			"bk_biz_id": bizID,
			"processes": processes,
		}
		result, e := s.CoreAPI.ProcServer().Process().UpdateProcessInstance(srvData.ctx, srvData.header, updateProcessOption)
		if e != nil {
			blog.ErrorJSON("updateProcessInstances failed, input: %s, err: %s, rid: %s", updateProcessOption, e.Error(), rid)
			return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if result.Result == false || result.Code != 0 {
			blog.ErrorJSON("UpdateProcessInstance failed, option: %s, result: %s, rid: %s", updateProcessOption, result, srvData.rid)
			return errors.New(result.Code, result.ErrMsg)
		}
	}
	return nil
}

func (s *Service) runTransferPlans(srvData *srvComm, bizID int64, transferPlan metadata.HostTransferPlan) errors.CCErrorCoder {
	rid := srvData.rid

	// step1 compute to be delete service instances
	listServiceInstanceOption := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostIDs:    []int64{transferPlan.HostID},
		ModuleIDs:  transferPlan.ToRemoveFromModules,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstances, ccErr := s.CoreAPI.CoreService().Process().ListServiceInstance(srvData.ctx, srvData.header, listServiceInstanceOption)
	if ccErr != nil {
		blog.ErrorJSON("runTransferPlans failed, ListServiceInstance failed, option: %s, err: %s, rid: %s", listServiceInstanceOption, ccErr.Error(), rid)
		return ccErr
	}
	serviceInstanceIDs := make([]int64, 0)
	for _, instance := range serviceInstances.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, instance.ID)
	}

	// clear service instance if necessary
	if len(serviceInstanceIDs) > 0 {
		// step2.1 delete process instance relation
		listRelationOption := &metadata.ListProcessInstanceRelationOption{
			BusinessID:         bizID,
			ServiceInstanceIDs: serviceInstanceIDs,
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}
		relationResult, ccErr := s.CoreAPI.CoreService().Process().ListProcessInstanceRelation(srvData.ctx, srvData.header, listRelationOption)
		if ccErr != nil {
			blog.ErrorJSON("runTransferPlans failed, ListProcessInstanceRelation failed, option: %s, err: %s, rid: %s", listRelationOption, ccErr.Error(), rid)
			return ccErr
		}
		processIDs := make([]int64, 0)
		for _, relation := range relationResult.Info {
			processIDs = append(processIDs, relation.ProcessID)
		}

		if len(processIDs) > 0 {
			deleteRelationOption := metadata.DeleteProcessInstanceRelationOption{
				BusinessID:         &bizID,
				ServiceInstanceIDs: serviceInstanceIDs,
			}
			ccErr = s.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(srvData.ctx, srvData.header, deleteRelationOption)
			if ccErr != nil {
				blog.ErrorJSON("runTransferPlans failed, DeleteProcessInstanceRelation failed, option: %s, err: %s, rid: %s", deleteRelationOption, ccErr.Error(), rid)
				return ccErr
			}

			// step2.2 delete process instance
			processDeleteOption := &metadata.DeleteOption{
				Condition: map[string]interface{}{
					common.BKProcessIDField: map[string]interface{}{
						common.BKDBIN: processIDs,
					},
				},
			}
			deleteProcessResult, err := s.CoreAPI.CoreService().Instance().DeleteInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, processDeleteOption)
			if err != nil {
				blog.ErrorJSON("runTransferPlans failed, DeleteInstance of process failed, option: %s, err: %s, rid: %s", processDeleteOption, err.Error(), rid)
				return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
			}
			if deleteProcessResult.Result == false {
				blog.ErrorJSON("runTransferPlans failed, DeleteInstance of process failed, option: %s, result: %s, rid: %s", processDeleteOption, deleteProcessResult, rid)
				return errors.New(deleteProcessResult.Code, deleteProcessResult.ErrMsg)
			}
		}

		// step2.3 delete service instance
		deleteServiceInstanceOption := &metadata.CoreDeleteServiceInstanceOption{
			BizID:              bizID,
			ServiceInstanceIDs: serviceInstanceIDs,
		}
		ccErr = s.CoreAPI.CoreService().Process().DeleteServiceInstance(srvData.ctx, srvData.header, deleteServiceInstanceOption)
		if ccErr != nil {
			blog.ErrorJSON("runTransferPlans failed, DeleteServiceInstance failed, option: %s, err: %s, rid: %s", deleteServiceInstanceOption, ccErr.Error(), rid)
			return ccErr
		}
	}

	// step3 transfer host
	var transferHostResult *metadata.OperaterException
	var err error
	var option interface{}
	if transferPlan.IsTransferToInnerModule == true {
		transferOption := &metadata.TransferHostToInnerModule{
			ApplicationID: bizID,
			HostID:        []int64{transferPlan.HostID},
			ModuleID:      transferPlan.FinalModules[0],
		}
		option = transferOption
		transferHostResult, err = s.CoreAPI.CoreService().Host().TransferToInnerModule(srvData.ctx, srvData.header, transferOption)
	} else {
		transferOption := &metadata.HostsModuleRelation{
			ApplicationID: bizID,
			HostID:        []int64{transferPlan.HostID},
			ModuleID:      transferPlan.FinalModules,
			IsIncrement:   false,
		}
		option = transferOption
		transferHostResult, err = s.CoreAPI.CoreService().Host().TransferToNormalModule(srvData.ctx, srvData.header, transferOption)
	}
	if err != nil {
		blog.ErrorJSON("runTransferPlans failed, transfer hosts failed, option: %s, err: %s, rid: %s", option, err.Error(), rid)
		return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if transferHostResult.Result == false {
		blog.ErrorJSON("runTransferPlans failed, transfer hosts failed, option: %s, result: %s, rid: %s", option, transferHostResult, rid)
		return errors.New(transferHostResult.Code, transferHostResult.ErrMsg)
	}

	// update host by host apply rule
	applyPlan := transferPlan.HostApplyPlan
	if len(applyPlan.UpdateFields) > 0 {
		updateOption := &metadata.UpdateOption{
			Data: applyPlan.GetUpdateData(),
			Condition: map[string]interface{}{
				common.BKHostIDField: applyPlan.HostID,
			},
		}
		updateResult, err := s.CoreAPI.CoreService().Instance().UpdateInstance(srvData.ctx, srvData.header, common.BKInnerObjIDHost, updateOption)
		if err != nil {
			blog.ErrorJSON("RunHostApplyRule, update host failed, option: %s, err: %s, rid: %s", updateOption, err.Error(), rid)
			return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if ccErr := updateResult.CCError(); ccErr != nil {
			blog.ErrorJSON("RunHostApplyRule, update host response failed, option: %s, response: %s, rid: %s", updateOption, updateResult, rid)
			return ccErr
		}
	}

	return nil
}

func (s *Service) generateTransferPlans(srvData *srvComm, bizID int64, option metadata.TransferHostWithAutoClearServiceInstanceOption) ([]metadata.HostTransferPlan, errors.CCErrorCoder) {
	rid := srvData.rid

	// step1. resolve host remove from modules
	removeFromModules := make([]int64, 0)
	if option.RemoveFromNode != nil {
		topoTree, ccErr := s.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(srvData.ctx, srvData.header, bizID, false)
		if ccErr != nil {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, SearchMainlineInstanceTopo failed, bizID: %d, err: %s, rid: %s", bizID, ccErr.Error(), rid)
			return nil, ccErr
		}
		topoNodePath := topoTree.TraversalFindNode(option.RemoveFromNode.ObjectID, option.RemoveFromNode.InstanceID)
		if len(topoNodePath) == 0 {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, remove_from_node invalid, bizID: %d, rid: %s", bizID, rid)
			err := srvData.ccErr.CCErrorf(common.CCErrCommParamsInvalid, "remove_from_node")
			return nil, err
		}
		topoNodePath[0].DeepFirstTraversal(func(node *metadata.TopoInstanceNode) {
			if node.ObjectID == common.BKInnerObjIDModule {
				removeFromModules = append(removeFromModules, node.InstanceID)
			}
		})
	}

	// step2. get host module config
	hostModuleOption := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     option.HostIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	hostModuleResult, err := s.CoreAPI.CoreService().Host().GetHostModuleRelation(srvData.ctx, srvData.header, hostModuleOption)
	if err != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, GetHostModuleRelation failed, option: %s, err: %s, rid: %s", hostModuleOption, err.Error(), rid)
		err := srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		return nil, err
	}
	if hostModuleResult.Result == false {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, GetHostModuleRelation failed, option: %s, result: %s, rid: %s", hostModuleOption, hostModuleResult, rid)
		err := errors.New(hostModuleResult.Code, hostModuleResult.ErrMsg)
		return nil, err
	}
	hostModulesIDMap := make(map[int64][]int64)
	for _, item := range hostModuleResult.Data.Info {
		if _, exist := hostModulesIDMap[item.HostID]; exist == false {
			hostModulesIDMap[item.HostID] = make([]int64, 0)
		}
		hostModulesIDMap[item.HostID] = append(hostModulesIDMap[item.HostID], item.ModuleID)
	}

	// get inner modules
	innerModules, ccErr := s.getInnerModules(*srvData, bizID)
	if ccErr != nil {
		return nil, ccErr
	}
	innerModuleIDs := make([]int64, 0)
	defaultInternalModuleID := int64(0)
	for _, module := range innerModules {
		innerModuleIDs = append(innerModuleIDs, module.ModuleID)
		if module.Default == int64(common.DefaultResModuleFlag) {
			defaultInternalModuleID = module.ModuleID
		}
	}
	if defaultInternalModuleID == 0 {
		blog.InfoJSON("TransferHostWithAutoClearServiceInstance failed, defaultInternalModuleID not found, bizID: %s, innerModules: %s, rid: %s", bizID, innerModules, rid)
	}
	if option.DefaultInternalModule != 0 && util.InArray(option.DefaultInternalModule, innerModuleIDs) == false {
		return nil, srvData.ccErr.CCErrorf(common.CCErrCommParamsInvalid, "default_internal_module")
	}
	if option.DefaultInternalModule != 0 {
		defaultInternalModuleID = option.DefaultInternalModule
	}

	transferPlans := make([]metadata.HostTransferPlan, 0)
	for hostID, currentInModules := range hostModulesIDMap {
		transferPlan := generateTransferPlan(currentInModules, removeFromModules, option.AddToModules, defaultInternalModuleID)
		transferPlan.HostID = hostID
		// check module compatibility
		finalModuleCount := len(transferPlan.FinalModules)
		for _, moduleID := range transferPlan.FinalModules {
			if util.InArray(moduleID, innerModuleIDs) && finalModuleCount != 1 {
				return nil, srvData.ccErr.CCError(common.CCErrHostTransferFinalModuleConflict)
			}
			if util.InArray(moduleID, innerModuleIDs) && finalModuleCount == 1 {
				transferPlan.IsTransferToInnerModule = true
			}
		}
		transferPlans = append(transferPlans, transferPlan)
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
	rules, ccErr := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(srvData.ctx, srvData.header, bizID, ruleOption)
	if ccErr != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, generateApplyPlan failed, ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, ruleOption, ccErr.Error(), rid)
		return transferPlans, ccErr
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for _, item := range transferPlans {
		host2Module := metadata.Host2Modules{
			HostID:    item.HostID,
			ModuleIDs: item.FinalModules,
		}
		hostModules = append(hostModules, host2Module)
	}
	planOption := metadata.HostApplyPlanOption{
		Rules:             rules.Info,
		HostModules:       hostModules,
		ConflictResolvers: option.Options.HostApplyConflictResolvers,
	}

	hostApplyPlanResult, ccErr := s.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(srvData.ctx, srvData.header, bizID, planOption)
	if err != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, generateApplyPlan failed, core service GenerateApplyPlan failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, planOption, ccErr.Error(), rid)
		return transferPlans, ccErr
	}
	hostApplyPlanMap := make(map[int64]metadata.OneHostApplyPlan)
	for _, item := range hostApplyPlanResult.Plans {
		hostApplyPlanMap[item.HostID] = item
	}
	for index, transferPlan := range transferPlans {
		applyPlan, ok := hostApplyPlanMap[transferPlan.HostID]
		if ok == false {
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
func generateTransferPlan(currentIn []int64, removeFrom []int64, addTo []int64, defaultInternalModuleID int64) metadata.HostTransferPlan {
	plan := metadata.HostTransferPlan{}

	// 主机最终所在模块列表
	finalModules := make([]int64, 0)
	for _, moduleID := range currentIn {
		if util.InArray(moduleID, removeFrom) {
			continue
		}
		finalModules = append(finalModules, moduleID)
	}
	finalModules = append(finalModules, addTo...)
	finalModules = util.IntArrayUnique(finalModules)
	if len(finalModules) == 0 {
		finalModules = []int64{defaultInternalModuleID}
	}
	plan.FinalModules = finalModules

	// 主机将会被移出的模块列表
	realRemoveModules := make([]int64, 0)
	for _, moduleID := range currentIn {
		if util.InArray(moduleID, finalModules) {
			continue
		}
		realRemoveModules = append(realRemoveModules, moduleID)
	}
	realRemoveModules = util.IntArrayUnique(realRemoveModules)
	plan.ToRemoveFromModules = realRemoveModules

	// 主机将会被新加到的模块列表
	realAddModules := make([]int64, 0)
	for _, moduleID := range finalModules {
		if util.InArray(moduleID, currentIn) {
			continue
		}
		realAddModules = append(realAddModules, moduleID)
	}
	realAddModules = util.IntArrayUnique(realAddModules)
	plan.ToAddToModules = realAddModules

	return plan
}

func (s *Service) validateModules(srvData *srvComm, bizID int64, moduleIDs []int64, field string) errors.CCErrorCoder {
	if len(moduleIDs) == 0 {
		return nil
	}
	moduleIDs = util.IntArrayUnique(moduleIDs)
	modules, err := s.getModules(srvData, bizID, moduleIDs)
	if err != nil {
		return err
	}
	if len(modules) != len(moduleIDs) {
		return srvData.ccErr.CCErrorf(common.CCErrCommParamsInvalid, field)
	}
	return nil
}

func (s *Service) getModule(srvData *srvComm, bizID int64, moduleID int64) (metadata.ModuleInst, errors.CCErrorCoder) {
	modules, err := s.getModules(srvData, bizID, []int64{moduleID})
	if err != nil {
		return metadata.ModuleInst{}, err
	}
	if len(modules) == 0 {
		blog.Errorf("getModule failed, not found, module ID: %d, rid: %s", moduleID, srvData.rid)
		return metadata.ModuleInst{}, srvData.ccErr.CCError(common.CCErrCommNotFound)
	}
	return modules[0], nil
}

func (s *Service) getModules(srvData *srvComm, bizID int64, moduleIDs []int64) ([]metadata.ModuleInst, errors.CCErrorCoder) {
	query := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Fields: []string{
			common.BKModuleIDField,
			common.BKDefaultField,
			common.BKModuleNameField,
			common.BKAppIDField,
			common.BKSetIDField,
			common.BKServiceTemplateIDField,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleIDs,
			},
		},
	}
	result, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.ErrorJSON("GetModules failed, http do error, input:%+v, err:%s, rid:%s", query, err.Error(), srvData.rid)
		return nil, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.ErrorJSON("GetModules failed, result failed, input:%+v, response: %s, rid:%s", query, result, srvData.rid)
		return nil, errors.New(result.Code, result.ErrMsg)
	}

	modules := make([]metadata.ModuleInst, 0)
	for _, item := range result.Data.Info {
		module := metadata.ModuleInst{}
		if err := mapstruct.Decode2Struct(item, &module); err != nil {
			return nil, srvData.ccErr.CCError(common.CCErrCommJSONUnmarshalFailed)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

func (s *Service) getInnerModules(srvData srvComm, bizID int64) ([]metadata.ModuleInst, errors.CCErrorCoder) {
	query := &metadata.QueryCondition{
		Limit: metadata.SearchLimit{
			Limit: common.BKNoLimit,
		},
		Fields: []string{
			common.BKModuleIDField,
			common.BKDefaultField,
			common.BKModuleNameField,
			common.BKAppIDField,
			common.BKSetIDField,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: common.DefaultFlagDefaultValue,
			},
		},
	}
	result, err := s.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.ErrorJSON("GetInnerModules failed, http do error, input:%+v, err:%s, rid:%s", query, err.Error(), srvData.rid)
		return nil, srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.ErrorJSON("GetInnerModules failed, result failed, input:%+v, response: %s, rid:%s", query, result, srvData.rid)
		return nil, errors.New(result.Code, result.ErrMsg)
	}

	modules := make([]metadata.ModuleInst, 0)
	for _, item := range result.Data.Info {
		module := metadata.ModuleInst{}
		if err := mapstruct.Decode2Struct(item, &module); err != nil {
			return nil, srvData.ccErr.CCError(common.CCErrCommJSONUnmarshalFailed)
		}
		modules = append(modules, module)
	}

	return modules, nil
}

// TransferHostWithAutoClearServiceInstancePreview generate a preview of changes for TransferHostWithAutoClearServiceInstance operation
// 接口请求参数跟转移是一致的
// 主机从模块删除时提供了将要删除的服务实例信息
// 主机添加到新模块时，提供了模块对应的服务模板（如果有）
func (s *Service) TransferHostWithAutoClearServiceInstancePreview(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	option := metadata.TransferHostWithAutoClearServiceInstanceOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("TransferHostWithAutoClearServiceInstancePreview failed, parse request body failed, err: %v, rid: %s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(option.HostIDs) == 0 {
		err := srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "bk_host_ids")
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	if option.RemoveFromNode == nil && option.AddToModules == nil {
		err := srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "add_to_modules")
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.V(7).Infof("parse bizID from url failed, bizID: %s, err: %+v, rid: %s", bizIDStr, srvData.rid)
		err := srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	if len(option.AddToModules) != 0 {
		if ccErr := s.validateModules(srvData, bizID, option.AddToModules, "add_to_modules"); ccErr != nil {
			_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: ccErr})
			return
		}
	}
	if option.DefaultInternalModule != 0 {
		if ccErr := s.validateModules(srvData, bizID, []int64{option.DefaultInternalModule}, "default_internal_module"); ccErr != nil {
			_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: ccErr})
			return
		}
	}

	transferPlans, ccErr := s.generateTransferPlans(srvData, bizID, option)
	if ccErr != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstancePreview failed, generateTransferPlans failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, ccErr.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: ccErr})
		return
	}
	addModuleIDs := make([]int64, 0)
	removeModuleIDs := make([]int64, 0)
	for _, plan := range transferPlans {
		addModuleIDs = append(addModuleIDs, plan.ToAddToModules...)
		removeModuleIDs = append(removeModuleIDs, plan.ToRemoveFromModules...)
	}

	// get to remove service instances
	listSrvInstOption := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostIDs:    option.HostIDs,
		ModuleIDs:  removeModuleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	srvInstResult, ccErr := s.CoreAPI.CoreService().Process().ListServiceInstance(srvData.ctx, srvData.header, listSrvInstOption)
	if ccErr != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstancePreview failed, ListServiceInstance failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, listSrvInstOption, ccErr.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: ccErr})
		return
	}
	moduleServiceInstanceMap := make(map[int64][]metadata.ServiceInstance)
	for _, item := range srvInstResult.Info {
		if _, exist := moduleServiceInstanceMap[item.ModuleID]; exist == false {
			moduleServiceInstanceMap[item.ModuleID] = make([]metadata.ServiceInstance, 0)
		}
		moduleServiceInstanceMap[item.ModuleID] = append(moduleServiceInstanceMap[item.ModuleID], item)
	}

	// get add to modules
	modules, ccErr := s.getModules(srvData, bizID, addModuleIDs)
	if ccErr != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstancePreview failed, ListServiceInstance failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, listSrvInstOption, ccErr.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: ccErr})
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
	serviceTemplateDetails, ccErr := s.CoreAPI.CoreService().Process().ListServiceTemplateDetail(srvData.ctx, srvData.header, bizID, serviceTemplateIDs...)
	if ccErr != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstancePreview failed, ListServiceTemplateDetail failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, listSrvInstOption, ccErr.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: ccErr})
		return
	}
	serviceTemplateMap := make(map[int64]metadata.ServiceTemplateDetail)
	for _, templateDetail := range serviceTemplateDetails.Info {
		serviceTemplateMap[templateDetail.ServiceTemplate.ID] = templateDetail
	}
	moduleServiceTemplateMap := make(map[int64]metadata.ServiceTemplateDetail)
	for _, module := range modules {
		templateDetail, exist := serviceTemplateMap[module.ServiceTemplateID]
		if exist == true {
			moduleServiceTemplateMap[module.ModuleID] = templateDetail
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
			serviceInstances, exist := moduleServiceInstanceMap[moduleID]
			if exist {
				for _, serviceInstance := range serviceInstances {
					if serviceInstance.HostID == preview.HostID {
						removeInfo.ServiceInstances = append(removeInfo.ServiceInstances, serviceInstance)
					}
				}
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

	_ = resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     previews,
	})
	return
}
