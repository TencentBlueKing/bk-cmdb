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
// 1. 将主机 bk_host_ids 从 remove_from 指定的拓扑节点移除
// 2. 移入到 add_to 指定的模块中
// 3. 自动删除主机在移除模块下的服务实例
// 4. 自动添加主机在新模块上的服务实例
// note:
// - 不允许 remove_from 和 add_to 同时为空
// - bk_host_ids 不允许为空
// - 如果 remove_from 指定为业务ID，则接口行为是：覆盖更新
// - 如果 remove_from 没有指定，仅仅是增量更新，无移除操作
// - 如果 add_to 没有指定，主机将仅仅从 remove_from 指定的模块中移除
// - 如果 add_to 是空先机/故障机/待回收模块中的一个，必须显式指定 remove_from(可指定成业务节点), 否则报主机不能属于互斥模块错误
// - 如果 add_to 是普通模块，主机当前数据空先机/故障机/待回收模块中的一个，必须显式指定 remove_from(可指定成业务节点), 否则报主机不能属于互斥模块错误
// - 模块同时出现在add_to和 remove_from 时，不会导致对应的服务实例被删除然后重新添加
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

	if option.RemoveFrom == nil && option.AddTo == nil {
		err := srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "add_to")
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
	for _, item := range transferPlans {
		err := s.runTransferPlans(srvData, bizID, item)
		hostTransferResult := HostTransferResult{
			HostID: item.HostID,
		}
		if err != nil {
			hostTransferResult.Code = err.GetCode()
			hostTransferResult.Message = err.Error()
			if firstErr == nil {
				firstErr = err
			}
		}
		transferResult = append(transferResult, hostTransferResult)
	}
	if firstErr != nil {
		response := metadata.RespError{
			Msg:     firstErr,
			ErrCode: firstErr.GetCode(),
			Data:    transferResult,
		}
		_ = resp.WriteEntity(response)
		return
	}

	_ = resp.WriteEntity(metadata.Response{Data: transferResult})
	return
}

func (s *Service) runTransferPlans(srvData *srvComm, bizID int64, transferPlan metadata.HostTransferPlan) errors.CCErrorCoder {
	// step1 compute to be delete service instances
	listServiceInstanceOption := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostIDs:    []int64{transferPlan.HostID},
		ModuleIDs:  transferPlan.RemoveFrom,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstances, ccErr := s.CoreAPI.CoreService().Process().ListServiceInstance(srvData.ctx, srvData.header, listServiceInstanceOption)
	if ccErr != nil {
		blog.ErrorJSON("runTransferPlans failed, ListServiceInstance failed, option: %s, err: %s, rid: %s", listServiceInstanceOption, ccErr.Error(), srvData.rid)
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
			blog.ErrorJSON("runTransferPlans failed, ListProcessInstanceRelation failed, option: %s, err: %s, rid: %s", listRelationOption, ccErr.Error(), srvData.rid)
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
				blog.ErrorJSON("runTransferPlans failed, DeleteProcessInstanceRelation failed, option: %s, err: %s, rid: %s", deleteRelationOption, ccErr.Error(), srvData.rid)
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
				blog.ErrorJSON("runTransferPlans failed, DeleteInstance of process failed, option: %s, err: %s, rid: %s", processDeleteOption, err.Error(), srvData.rid)
				return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
			}
			if deleteProcessResult.Result == false {
				blog.ErrorJSON("runTransferPlans failed, DeleteInstance of process failed, option: %s, result: %s, rid: %s", processDeleteOption, deleteProcessResult, srvData.rid)
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
			blog.ErrorJSON("runTransferPlans failed, DeleteServiceInstance failed, option: %s, err: %s, rid: %s", deleteServiceInstanceOption, ccErr.Error(), srvData.rid)
			return ccErr
		}
	}

	// step3 transfer host
	var transferHostResult *metadata.OperaterException
	var err error
	var option interface{}
	if transferPlan.ToInnerModule == true {
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
		blog.ErrorJSON("runTransferPlans failed, transfer hosts failed, option: %s, err: %s, rid: %s", option, err.Error(), srvData.rid)
		return srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if transferHostResult.Result == false {
		blog.ErrorJSON("runTransferPlans failed, transfer hosts failed, option: %s, result: %s, rid: %s", option, transferHostResult, srvData.rid)
		return errors.New(transferHostResult.Code, transferHostResult.ErrMsg)
	}
	return nil
}

func (s *Service) generateTransferPlans(srvData *srvComm, bizID int64, option metadata.TransferHostWithAutoClearServiceInstanceOption) ([]metadata.HostTransferPlan, errors.CCErrorCoder) {
	// step1. resolve host remove from modules
	removeFromModules := make([]int64, 0)
	if option.RemoveFrom != nil {
		topoTree, ccErr := s.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(srvData.ctx, srvData.header, bizID, false)
		if ccErr != nil {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, SearchMainlineInstanceTopo failed, bizID: %d, err: %s, rid: %s", bizID, ccErr.Error(), srvData.rid)
			return nil, ccErr
		}
		topoNodePath := topoTree.TraversalFindNode(option.RemoveFrom.ObjectID, option.RemoveFrom.InstanceID)
		if len(topoNodePath) == 0 {
			blog.Errorf("TransferHostWithAutoClearServiceInstance failed, remove_from invalid, bizID: %d, rid: %s", bizID, srvData.rid)
			err := srvData.ccErr.CCErrorf(common.CCErrCommParamsInvalid, "remove_from")
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
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, GetHostModuleRelation failed, option: %s, err: %s, rid: %s", hostModuleOption, err.Error(), srvData.rid)
		err := srvData.ccErr.CCError(common.CCErrCommHTTPDoRequestFailed)
		return nil, err
	}
	if hostModuleResult.Result == false {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, GetHostModuleRelation failed, option: %s, result: %s, rid: %s", hostModuleOption, hostModuleResult, srvData.rid)
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
	innerModules, ccErr := s.GetInnerModules(*srvData, bizID)
	if ccErr != nil {
		return nil, ccErr
	}
	innerModuleIDs := make([]int64, 0)
	for _, module := range innerModules {
		innerModuleIDs = append(innerModuleIDs, module.ModuleID)
	}

	transferPlans := make([]metadata.HostTransferPlan, 0)
	for hostID, currentInModules := range hostModulesIDMap {
		transferPlan := generateTransferPlan(currentInModules, removeFromModules, option.AddTo)
		transferPlan.HostID = hostID
		// check module compatibility
		finalModuleCount := len(transferPlan.FinalModules)
		for _, moduleID := range transferPlan.FinalModules {
			if util.InArray(moduleID, innerModuleIDs) && finalModuleCount != 1 {
				return nil, srvData.ccErr.CCError(common.CCErrHostTransferFinalModuleConflict)
			}
			if util.InArray(moduleID, innerModuleIDs) && finalModuleCount == 1 {
				transferPlan.ToInnerModule = true
			}
		}
		transferPlans = append(transferPlans, transferPlan)
	}
	return transferPlans, nil
}

// generateTransferPlan 实现计算主机将从哪个模块移除，添加到哪个模块，最终在哪些模块
// param hostID: 主机ID
// param currentIn: 主机当前所属模块
// param removeFrom: 从哪些模块中移除
// param addTo: 添加到哪些模块
func generateTransferPlan(currentIn []int64, removeFrom []int64, addTo []int64) metadata.HostTransferPlan {
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
	plan.FinalModules = finalModules

	// 主机将会被移出的模块列表
	realRemoveModules := make([]int64, 0)
	for _, moduleID := range removeFrom {
		if util.InArray(moduleID, finalModules) {
			continue
		}
		realRemoveModules = append(realRemoveModules, moduleID)
	}
	realRemoveModules = util.IntArrayUnique(realRemoveModules)
	plan.RemoveFrom = realRemoveModules

	// 主机将会被新加到的模块列表
	realAddModules := make([]int64, 0)
	for _, moduleID := range addTo {
		if util.InArray(moduleID, currentIn) {
			continue
		}
		realAddModules = append(realAddModules, moduleID)
	}
	realAddModules = util.IntArrayUnique(realAddModules)
	plan.AddTo = realAddModules

	return plan
}

func (s *Service) GetInnerModules(srvData srvComm, bizID int64) ([]metadata.ModuleInst, errors.CCErrorCoder) {
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
				common.BKDBNE: 0,
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
