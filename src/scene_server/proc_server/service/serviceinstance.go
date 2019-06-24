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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) CreateProcessInstances(ctx *rest.Contexts) {
	input := new(metadata.CreateRawProcessInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	processIDs, err := ps.createProcessInstances(ctx, input)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcCreateProcessFailed, "create service instance failed, serviceInstanceID: %d, err: %+v", input.ServiceInstanceID, err)
		return
	}
	ctx.RespEntity(processIDs)
}

func (ps *ProcServer) createProcessInstances(ctx *rest.Contexts, input *metadata.CreateRawProcessInstanceInput) ([]int64, errors.CCErrorCoder) {
	bizID, e := metadata.BizIDFromMetadata(input.Metadata)
	if e != nil {
		blog.Errorf("create process instance with raw, parse biz id from metadata failed, err: %+v, rid: %s", e, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommHTTPInputInvalid, common.MetadataField)
	}

	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstanceID)
	if err != nil {
		blog.Errorf("create process instance failed, get service instance by id failed, serviceInstanceID: %d, err: %v", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, err
	}
	businessID, e := metadata.BizIDFromMetadata(serviceInstance.Metadata)
	if e != nil {
		blog.Errorf("create process instance with raw, parse biz id from service instance metadata failed, err: %+v, rid: %s", e, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParseBizIDFromMetadataInDBFailed, common.MetadataField)
	}
	if businessID != bizID {
		blog.Errorf("create process instance with raw, biz id from input not equal with service instance, err: %+v, rid: %s", e, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.MetadataField)
	}
	if serviceInstance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		blog.Errorf("create process instance failed, create process instance on service instance initialized by template forbidden, serviceInstanceID: %d, err: %v", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrProcEditProcessInstanceCreateByTemplateForbidden)
	}

	processIDs := make([]int64, 0)
	for _, process := range input.Processes {
		process.ProcessInfo.ProcessID = 0
		process.ProcessInfo.BusinessID = bizID
		process.ProcessInfo.SupplierAccount = ctx.Kit.SupplierAccount
		now := time.Now()
		process.ProcessInfo.CreateTime = now
		process.ProcessInfo.LastTime = now

		if err := ps.validateRawInstanceUnique(ctx, serviceInstance.ID, &process.ProcessInfo); err != nil {
			ctx.RespWithError(err, common.CCErrProcCreateProcessFailed, "create process instance failed, serviceInstanceID: %d, process: %+v, err: %v", input.ServiceInstanceID, process, err)
			return nil, err
		}

		processID, err := ps.Logic.CreateProcessInstance(ctx.Kit, &process.ProcessInfo)
		if err != nil {
			blog.Errorf("create process instance failed, create process failed, serviceInstanceID: %d, process: %+v, err: %v, rid: %s", input.ServiceInstanceID, process, err, ctx.Kit.Rid)
			return nil, err
		}

		relation := &metadata.ProcessInstanceRelation{
			Metadata:          input.Metadata,
			ProcessID:         processID,
			ProcessTemplateID: common.ServiceTemplateIDNotSet,
			ServiceInstanceID: serviceInstance.ID,
			HostID:            serviceInstance.HostID,
		}

		_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relation)
		if err != nil {
			blog.Errorf("create service instance relations, create process instance relation failed, serviceInstanceID: %d, relation: %+v, err: %v", input.ServiceInstanceID, relation, err)
			return nil, err
		}
		processIDs = append(processIDs, processID)
	}

	return processIDs, nil
}

func (ps *ProcServer) UpdateProcessInstances(ctx *rest.Contexts) {
	input := new(metadata.UpdateRawProcessInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update process instance failed, parse business id failed, err: %+v", err)
		return
	}
	processIDs := make([]int64, 0)
	for _, process := range input.Processes {
		if process.ProcessID == 0 {
			ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "update process instance failed, process_id invalid", common.BKProcessIDField)
			return
		}
		processIDs = append(processIDs, process.ProcessID)
	}
	option := &metadata.ListProcessInstanceRelationOption{
		BusinessID: bizID,
		ProcessIDs: &processIDs,
		Page:       metadata.BasePage{Limit: common.BKNoLimit},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPDoRequestFailed, "update process instance failed, search process instance relation failed, err: %+v", err)
		return
	}

	processTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	for _, relation := range relations.Info {
		if relation.ProcessTemplateID == common.ServiceTemplateIDNotSet {
			continue
		}
		if _, exist := processTemplateMap[relation.ProcessTemplateID]; exist == true {
			continue
		}
		processTemplate, err := ps.CoreAPI.CoreService().Process().GetProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, relation.ProcessTemplateID)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPDoRequestFailed, "update process instance failed, search process instance relation failed, err: %+v", err)
			return
		}
		processTemplateMap[relation.ProcessTemplateID] = processTemplate
	}

	process2ServiceInstanceMap := make(map[int64]*metadata.ProcessInstanceRelation)
	for _, relation := range relations.Info {
		process2ServiceInstanceMap[relation.ProcessID] = &relation
	}

	var processTemplate *metadata.ProcessTemplate
	for _, process := range input.Processes {
		relation, exist := process2ServiceInstanceMap[process.ProcessID]
		if exist == false {
			err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
			ctx.RespWithError(err, common.CCErrCommParamsInvalid, "update process instance failed, process related service instance not found, process: %+v, err: %v", process, err)
			return
		}
		if relation.ProcessTemplateID == 0 {
			serviceInstanceID := relation.ServiceInstanceID
			if err := ps.validateRawInstanceUnique(ctx, serviceInstanceID, &process); err != nil {
				ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed, "update process instance failed, serviceInstanceID: %d, process: %+v, err: %v", serviceInstanceID, process, err)
				return
			}
		} else {
			processTemplate, exist = processTemplateMap[relation.ProcessTemplateID]
			if exist == false {
				err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
				ctx.RespWithError(err, common.CCErrCommNotFound, "update process instance failed, process related template not found, relation: %+v, err: %v", relation, err)
				return
			}
			processTemplate.InstanceUpdate(&process)
		}

		processID := process.ProcessID
		process.BusinessID = bizID
		process.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
		processBytes, err := json.Marshal(process)
		if err != nil {
			blog.Errorf("UpdateProcessInstances failed, json Marshal process failed, process: %+v, err: %+v", process, err)
			err := ctx.Kit.CCError.CCError(common.CC_ERR_Comm_JSON_ENCODE)
			ctx.RespWithError(err, common.CC_ERR_Comm_JSON_DECODE, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
		}
		processData := mapstr.MapStr{}
		if err := json.Unmarshal(processBytes, &processData); nil != err && 0 != len(processBytes) {
			blog.Errorf("UpdateProcessInstances failed, json Unmarshal process failed, processData: %s, err: %+v", processData, err)
			err := ctx.Kit.CCError.CCError(common.CC_ERR_Comm_JSON_DECODE)
			ctx.RespWithError(err, common.CC_ERR_Comm_JSON_DECODE, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
		}
		processData.Remove(common.BKProcessIDField)
		processData.Remove(common.MetadataField)
		processData.Remove(common.LastTimeField)
		processData.Remove(common.CreateTimeField)
		if err := ps.Logic.UpdateProcessInstance(ctx.Kit, processID, processData); err != nil {
			ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
			return
		}
	}

	ctx.RespEntity(processIDs)
}

func (ps *ProcServer) CheckHostInBusiness(ctx *rest.Contexts, bizID int64, hostIDs []int64) errors.CCErrorCoder {
	hostIDHit := make(map[int64]bool)
	for _, hostID := range hostIDs {
		hostIDHit[hostID] = false
	}
	hostConfigFilter := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
	}
	result, err := ps.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, hostConfigFilter)
	if err != nil {
		e, ok := err.(errors.CCErrorCoder)
		if ok == true {
			return e
		} else {
			return ctx.Kit.CCError.CCError(common.CCErrWebGetHostFail)
		}
	}
	for _, item := range result.Data {
		hostIDHit[item.HostID] = true
	}
	invalidHost := make([]int64, 0)
	for hostID, hit := range hostIDHit {
		if hit == false {
			invalidHost = append(invalidHost, hostID)
		}
	}
	if len(invalidHost) > 0 {
		return ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotBelongBusiness, invalidHost, bizID)
	}
	return nil
}

// createServiceInstances 创建服务实例
// 支持直接创建和通过模板创建，用 module 是否绑定模版信息区分两种情况
func (ps *ProcServer) CreateServiceInstances(ctx *rest.Contexts) {
	input := new(metadata.CreateServiceInstanceForServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create service instance with template : %d, moduleID: %d, but get business id failed, err: %v", input.ModuleID, err)
		return
	}

	// check hosts in business
	hostIDs := make([]int64, 0)
	hostIDHit := make(map[int64]bool)
	for _, instance := range input.Instances {
		if util.InArray(instance.HostID, hostIDs) == false {
			hostIDs = append(hostIDs, instance.HostID)
			hostIDHit[instance.HostID] = false
		}
	}
	if err := ps.CheckHostInBusiness(ctx, bizID, hostIDs); err != nil {
		ctx.RespWithError(err, common.CCErrCoreServiceHostNotBelongBusiness, "create service instance failed, host %+v not belong to business %d, hostIDs: %+v, err: %v", hostIDs, bizID, err)
		return
	}

	module, err := ps.getModule(ctx, input.ModuleID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrTopoGetModuleFailed, "create service instance failed, get module failed, moduleID: %d, err: %v", input.ModuleID, err)
		return
	}

	if module.BizID != bizID {
		err := ctx.Kit.CCError.Errorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, module.ModuleID, bizID)
		ctx.RespWithError(err, common.CCErrCoreServiceHasModuleNotBelongBusiness, "create service instance failed, module %d not belongs to biz %d, err: %v", input.ModuleID, bizID, err)
		return
	}

	serviceInstanceIDs := make([]int64, 0)
	for _, inst := range input.Instances {
		instance := &metadata.ServiceInstance{
			Metadata:          input.Metadata,
			Name:              input.Name,
			ServiceTemplateID: module.ServiceTemplateID,
			ModuleID:          input.ModuleID,
			HostID:            inst.HostID,
		}

		// create service instance at first
		serviceInstance, err := ps.CoreAPI.CoreService().Process().CreateServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, instance)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service instance failed, moduleID: %d, err: %v", input.ModuleID, err)
			return
		}

		if module.ServiceTemplateID == 0 && len(inst.Processes) > 0 {
			// if this service have process instance to create, then create it now.
			createProcessInput := &metadata.CreateRawProcessInstanceInput{
				Metadata:          input.Metadata,
				ServiceInstanceID: serviceInstance.ID,
				Processes:         inst.Processes,
			}
			if _, err := ps.createProcessInstances(ctx, createProcessInput); err != nil {
				ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service instance failed, create process instances failed, moduleID: %d, err: %v", input.ModuleID, err)
				return
			}
		}

		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	ctx.RespEntity(serviceInstanceIDs)
}

func (ps *ProcServer) getModule(ctx *rest.Contexts, moduleID int64) (*metadata.ModuleInst, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	moduleFilter := &metadata.QueryCondition{
		Condition: mapstr.MapStr(filter),
	}
	modules, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, moduleFilter)
	if err != nil {
		blog.Errorf("getModule failed, moduleID: %d, err: %+v, rid: %s", moduleID, err, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, err)
	}
	if len(modules.Data.Info) == 0 {
		blog.Errorf("getModule failed, moduleID: %d, err: %+v, rid: %s", moduleID, "not found", ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, "not found")
	}
	if len(modules.Data.Info) > 1 {
		blog.Errorf("getModule failed, moduleID: %d, err: %+v, rid: %s", moduleID, "get multiple", ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, "get multiple modules")
	}
	module := modules.Data.Info[0]
	moduleInst := &metadata.ModuleInst{}
	if err := module.ToStructByTag(moduleInst, "field"); err != nil {
		blog.Errorf("getModule failed, marshal json failed, moduleID: %d, err: %+v, rid: %s", moduleID, err, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}
	return moduleInst, nil
}

func (ps *ProcServer) validateRawInstanceUnique(ctx *rest.Contexts, serviceInstanceID int64, processInfo *metadata.Process) errors.CCErrorCoder {
	if len(processInfo.ProcessName) == 0 || len(processInfo.ProcessName) > common.NameFieldMaxLength {
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessNameField)
	}
	if len(processInfo.FuncName) == 0 || len(processInfo.ProcessName) > common.NameFieldMaxLength {
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFuncName)
	}
	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, get service instance failed, metadata: %+v, err: %v, rid: %s", serviceInstance.Metadata, err, ctx.Kit.Rid)
		return err
	}

	// find process under service instance
	bizID, e := metadata.BizIDFromMetadata(serviceInstance.Metadata)
	if e != nil {
		blog.Errorf("validateRawInstanceUnique failed, parse business id from metadata failed, metadata: %+v, err: %v, rid: %s", serviceInstance.Metadata, e, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommParseBizIDFromMetadataInDBFailed)
	}
	relationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: &[]int64{serviceInstance.ID},
		ProcessTemplateID:  common.ServiceTemplateIDNotSet,
		HostID:             serviceInstance.HostID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationOption)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, get relation under service instance failed, err: %v, rid: %s", serviceInstance.Metadata, err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	existProcessIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		existProcessIDs = append(existProcessIDs, relation.ProcessID)
	}
	otherProcessIDs := existProcessIDs
	if processInfo.ProcessID != 0 {
		otherProcessIDs = make([]int64, 0)
		for _, processID := range existProcessIDs {
			if processID != processInfo.ProcessID {
				otherProcessIDs = append(otherProcessIDs, processID)
			}
		}
	}
	// process name unique
	processNameFilter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: otherProcessIDs,
		},
		common.BKProcessNameField: processInfo.ProcessName,
	}
	processNameFilterCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(processNameFilter),
	}
	listResult, e := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKProcessObjectName, processNameFilterCond)
	if e != nil {
		blog.Errorf("validateRawInstanceUnique failed, search process with bk_process_name failed, filter: %+v, err: %v, rid: %s", processNameFilter, e, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if listResult.Data.Count > 0 {
		blog.Errorf("validateRawInstanceUnique failed, bk_process_name duplicated under service instance, err: %v, rid: %s", serviceInstance.Metadata, err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCoreServiceProcessNameDuplicated)
	}

	// func name unique
	funcNameFilter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: otherProcessIDs,
		},
		common.BKStartParamRegex: processInfo.StartParamRegex,
		common.BKFuncName:        processInfo.FuncName,
	}
	funcNameFilterCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(funcNameFilter),
	}
	listFuncNameResult, e := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKProcessObjectName, funcNameFilterCond)
	if e != nil {
		blog.Errorf("validateRawInstanceUnique failed, search process with func name failed, filter: %+v, err: %v, rid: %s", funcNameFilterCond, e, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if listFuncNameResult.Data.Count > 0 {
		blog.Errorf("validateRawInstanceUnique failed, bk_func_name and bk_start_param_regex duplicated under service instance, err: %v, rid: %s", err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCoreServiceFuncNameDuplicated)
	}
	return nil
}

func (ps *ProcServer) DeleteProcessInstance(ctx *rest.Contexts) {
	input := new(metadata.DeleteProcessInstanceInServiceInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete process instance in service instance failed, err: %v", err)
		return
	}

	// delete process relation at the same time.
	deleteOption := metadata.DeleteProcessInstanceRelationOption{}
	deleteOption.ProcessIDs = &input.ProcessInstanceIDs
	err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "delete process instance: %v, but delete instance relation failed.", input.ProcessInstanceIDs)
		return
	}

	if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, input.ProcessInstanceIDs); err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "delete process instance:%v failed, err: %v", input.ProcessInstanceIDs, err)
		return
	}

	ctx.RespEntity(nil)
}

func (ps *ProcServer) SearchServiceInstancesInModule(ctx *rest.Contexts) {
	input := new(metadata.GetServiceInstanceInModuleInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get service instances in module, but parse biz id failed, err: %v", err)
		return
	}

	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		ModuleID:   input.ModuleID,
		Page:       input.Page,
		WithName:   input.WithName,
		SearchKey:  input.SearchKey,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in module: %d failed, err: %v", input.ModuleID, err)
		return
	}

	ctx.RespEntity(instances)
}

func (ps *ProcServer) DeleteServiceInstance(ctx *rest.Contexts) {
	input := new(metadata.DeleteServiceInstanceOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete service instances, but parse biz id failed, err: %v", err)
		return
	}

	// when a service instance is deleted, the related data should be deleted at the same time
	for _, serviceInstanceID := range input.ServiceInstanceIDs {
		serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed, "delete service instance failed, service instance not found, serviceInstanceIDs: %d", serviceInstanceID)
			return
		}
		businessID, e := metadata.BizIDFromMetadata(serviceInstance.Metadata)
		if e != nil {
			ctx.RespWithError(err, common.CCErrCommParseBizIDFromMetadataInDBFailed, "delete service instance failed, parse biz id from service instance metadata failed, serviceInstanceIDs: %d, err: %+v", serviceInstanceID, e)
			return
		}
		if businessID != bizID {
			err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.MetadataField)
			ctx.RespWithError(err, common.CCErrCommParamsInvalid, "delete service instance failed, biz id from input and service instance not equal, serviceInstanceIDs: %d", serviceInstanceID)
			return
		}

		// step1: delete the service instance relation.
		option := &metadata.ListProcessInstanceRelationOption{
			BusinessID:         bizID,
			ServiceInstanceIDs: &[]int64{serviceInstanceID},
		}
		relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed, "delete service instance: %d, but list service instance relation failed.", serviceInstanceID)
			return
		}

		deleteOption := metadata.DeleteProcessInstanceRelationOption{
			ServiceInstanceIDs: &[]int64{serviceInstanceID},
		}
		err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d, but delete service instance relations failed.", serviceInstanceID)
			return
		}

		// step2: delete process instance belongs to this service instance.
		var processIDs []int64
		for _, r := range relations.Info {
			processIDs = append(processIDs, r.ProcessID)
		}
		if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, processIDs); err != nil {
			ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d, but delete process instance failed.", serviceInstanceID)
			return
		}

		// step3: delete service instance.
		deleteSvcInstOption := &metadata.DeleteServiceInstanceOption{
			ServiceInstanceIDs: []int64{serviceInstanceID},
		}
		err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, deleteSvcInstOption)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d failed, err: %v", serviceInstanceID, err)
			return
		}

		// step4: check and move host from module if no serviceInstance on it
		filter := &metadata.ListServiceInstanceOption{
			BusinessID: bizID,
			HostID:     serviceInstance.HostID,
			ModuleID:   serviceInstance.ModuleID,
		}
		result, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, filter)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get host related service instances failed, bizID: %d, serviceInstanceID: %d, err: %v", bizID, serviceInstance.HostID, err)
			return
		}
		if len(result.Info) == 0 {
			continue
		}
		// just remove host from this module
		removeHostFromModuleOption := metadata.RemoveHostsFromModuleOption{
			ApplicationID: bizID,
			HostID:        serviceInstance.HostID,
			ModuleID:      serviceInstance.ModuleID,
		}
		if _, err := ps.CoreAPI.CoreService().Host().RemoveFromModule(ctx.Kit.Ctx, ctx.Kit.Header, &removeHostFromModuleOption); err != nil {
			ctx.RespWithError(err, common.CCErrHostMoveResourcePoolFail, "remove host from module failed, option: %+v, err: %v", removeHostFromModuleOption, err)
			return
		}
	}
	ctx.RespEntity(nil)
}

// this function works to find differences between the service template and service instances in a module.
// compared to the service template's process template, a process instance in the service instance may
// contains several differences, like as follows:
// unchanged: the process instance's property values are same with the process template it belongs.
// changed: the process instance's property values are not same with the process template it belongs.
// add: a new process template is added, compared to the service instance belongs to this service template.
// deleted: a process is already deleted, compared to the service instance belongs to this service template.
func (ps *ProcServer) DiffServiceInstanceWithTemplate(ctx *rest.Contexts) {
	diffOption := new(metadata.DiffServiceInstanceWithTemplateOption)
	if err := ctx.DecodeInto(diffOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// why we need validate metadata here?
	if _, err := metadata.BizIDFromMetadata(diffOption.Metadata); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "find difference between service template and process instances, but parse biz id failed, err: %v", err)
		return
	}

	if diffOption.ModuleID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "find difference between service template and process instances, but got empty service template id or module id")
		return
	}
	module, err := ps.getModule(ctx, diffOption.ModuleID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrTopoGetModuleFailed, "find difference between service template and process instances failed, get module by id:%d failed, err: %+v", diffOption.ModuleID, err)
		return
	}

	// step 1:
	// find process object's attribute
	cond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKObjIDField: common.BKInnerObjIDProc,
		}),
	}
	attrResult, e := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDProc, cond)
	if e != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
			"find difference between service template: %d and process instances, bizID: %d, but get process attributes failed, err: %v",
			module.ServiceTemplateID, module.BizID, e)
		return
	}
	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}

	// step2. get process templates
	listProcessTemplateOption := &metadata.ListProcessTemplatesOption{
		BusinessID:        module.BizID,
		ServiceTemplateID: module.ServiceTemplateID,
	}
	processTemplates, e := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listProcessTemplateOption)
	if e != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed, "find difference between service template: %d and process instances, bizID: %d, but get process templates failed, err: %v", module.ServiceTemplateID, module.BizID, e)
		return
	}

	// step 3:
	// find process instance's relations, which allows us know the relationship between
	// process instance and it's template, service instance, etc.
	pTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	serviceRelationMap := make(map[int64][]metadata.ProcessInstanceRelation)
	for idx, pTemplate := range processTemplates.Info {
		pTemplateMap[pTemplate.ID] = &processTemplates.Info[idx]

		option := metadata.ListProcessInstanceRelationOption{
			BusinessID:        module.BizID,
			ProcessTemplateID: pTemplate.ID,
		}

		relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
				"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance relations failed, err: %v",
				module.ServiceTemplateID, module.BizID, diffOption.ModuleID, err)
			return
		}

		for _, r := range relations.Info {
			serviceRelationMap[r.ServiceInstanceID] = append(serviceRelationMap[r.ServiceInstanceID], r)
		}
	}

	// step 4:
	// find all the service instances belongs to this service template and this module.
	// which contains the process instances details at the same time.
	serviceOption := &metadata.ListServiceInstanceOption{
		BusinessID:        module.BizID,
		ServiceTemplateID: module.ServiceTemplateID,
		ModuleID:          diffOption.ModuleID,
	}
	serviceInstances, e := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceOption)
	if e != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed,
			"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance failed, err: %v",
			module.ServiceTemplateID, module.BizID, diffOption.ModuleID, e)
		return
	}

	// step 5: compare the process instance with it's process template one by one in a service instance.
	type recorder struct {
		ProcessID        int64
		ProcessName      string
		ServiceInstance  *metadata.ServiceInstance
		ChangedAttribute []metadata.ProcessChangedAttribute
	}
	removed := make(map[int64][]recorder)
	changed := make(map[int64][]recorder)
	unchanged := make(map[int64][]recorder)
	added := make(map[int64]bool, 0)
	processTemplateReferenced := make(map[int64]int64)
	for _, serviceInstance := range serviceInstances.Info {
		relations := serviceRelationMap[serviceInstance.ID]

		for _, relation := range relations {
			// record the used process template for checking whether a new process template has been added to service template.
			processTemplateReferenced[relation.ProcessTemplateID] += 1

			process, err := ps.Logic.GetProcessInstanceWithID(ctx.Kit, relation.ProcessID)
			if err != nil {
				if err.GetCode() == common.CCErrCommNotFound {
					process = new(metadata.Process)
				} else {
					ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed,
						"get difference between with process template and process instance in a service instance, but get process instance: %d failed, %v", err)
					return
				}
			}

			property, exist := pTemplateMap[relation.ProcessTemplateID]
			if !exist {
				// process's template doesn't exist means the template has already been removed.
				removed[relation.ProcessTemplateID] = append(removed[relation.ProcessTemplateID], recorder{
					ProcessID:       relation.ProcessID,
					ProcessName:     process.ProcessName,
					ServiceInstance: &serviceInstance,
				})
				continue
			}

			diff := ps.Logic.DiffWithProcessTemplate(property.Property, process, attributeMap)
			if len(diff) == 0 {
				// nothing changed
				unchanged[relation.ProcessTemplateID] = append(unchanged[relation.ProcessTemplateID], recorder{
					ProcessID:       relation.ProcessID,
					ProcessName:     process.ProcessName,
					ServiceInstance: &serviceInstance,
				})
				continue
			}

			// something has already changed.
			changed[relation.ProcessTemplateID] = append(changed[relation.ProcessTemplateID], recorder{
				ProcessID:        relation.ProcessID,
				ProcessName:      process.ProcessName,
				ServiceInstance:  &serviceInstance,
				ChangedAttribute: diff,
			})
		}

		// check whether a new process template has been added.
		for templateID := range pTemplateMap {
			if _, exist := processTemplateReferenced[templateID]; exist == true {
				continue
			}
			// the process template does not exist in all the service instances,
			// which means a new process template is added.
			added[templateID] = true
		}
	}

	// it's time to rearrange the data
	differences := metadata.ProcessTemplateWithInstancesDifference{
		Unchanged: make([]metadata.ServiceInstanceDifferenceDetail, 0),
		Changed:   make([]metadata.ServiceInstanceDifferenceDetail, 0),
		Added:     make([]metadata.ServiceInstanceDifferenceDetail, 0),
		Removed:   make([]metadata.ServiceInstanceDifferenceDetail, 0),
	}

	for removedID, records := range removed {
		if len(records) == 0 {
			continue
		}
		processTemplateName := records[0].ProcessName

		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, record := range records {
			serviceInstances = append(serviceInstances, metadata.ServiceDifferenceDetails{ServiceInstance: *record.ServiceInstance})
		}
		differences.Removed = append(differences.Removed, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    removedID,
			ProcessTemplateName:  processTemplateName,
			ServiceInstanceCount: len(serviceInstances),
			ServiceInstances:     serviceInstances,
		})
	}

	for unchangedID, records := range unchanged {
		if len(records) == 0 {
			continue
		}
		processTemplateName := records[0].ProcessName
		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, record := range records {
			serviceInstances = append(serviceInstances, metadata.ServiceDifferenceDetails{ServiceInstance: *record.ServiceInstance})
		}
		differences.Unchanged = append(differences.Unchanged, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    unchangedID,
			ProcessTemplateName:  processTemplateName,
			ServiceInstanceCount: len(serviceInstances),
			ServiceInstances:     serviceInstances,
		})
	}

	for changedID, records := range changed {
		if len(records) == 0 {
			continue
		}
		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, record := range records {
			serviceInstances = append(serviceInstances, metadata.ServiceDifferenceDetails{
				ServiceInstance:   *record.ServiceInstance,
				ChangedAttributes: record.ChangedAttribute,
			})
		}
		differences.Changed = append(differences.Changed, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    changedID,
			ProcessTemplateName:  records[0].ProcessName,
			ServiceInstanceCount: len(serviceInstances),
			ServiceInstances:     serviceInstances,
		})
	}

	for addedID := range added {
		sInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, s := range serviceInstances.Info {
			sInstances = append(sInstances, metadata.ServiceDifferenceDetails{ServiceInstance: s})
		}

		differences.Added = append(differences.Added, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    addedID,
			ProcessTemplateName:  pTemplateMap[addedID].ProcessName,
			ServiceInstanceCount: len(sInstances),
			ServiceInstances:     sInstances,
		})
	}

	ctx.RespEntity(differences)
}

// SyncServiceInstanceByTemplate sync the service instance with it's bounded service template.
// It keeps the processes exactly same with the process template in the service template,
// which means the number of process is same, and the process instance's info is also exactly same.
// It contains several scenarios in a service instance:
// 1. add a new process
// 2. update a process
// 3. removed a process
func (ps *ProcServer) SyncServiceInstanceByTemplate(ctx *rest.Contexts) {
	syncOption := new(metadata.SyncServiceInstanceByTemplateOption)
	if err := ctx.DecodeInto(syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(syncOption.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "force sync service instance according to service template, but parse biz id failed, err: %v", err)
		return
	}

	module, err := ps.getModule(ctx, syncOption.ModuleID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrTopoGetModuleFailed, "force sync service instance according to service template, get module by id:%d failed, err: %+v", syncOption.ModuleID, err)
		return
	}

	// step 0:
	// find service instances
	serviceInstanceOption := &metadata.ListServiceInstanceOption{
		BusinessID:        bizID,
		ModuleID:          syncOption.ModuleID,
		ServiceTemplateID: module.ServiceTemplateID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		WithName: false,
	}
	serviceInstanceResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "sync service instance with template: %d failed, get service instances failed, err: %v", module.ServiceTemplateID, err)
		return
	}
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	// step 1:
	// find all the process template according to the service template id
	processTemplateFilter := &metadata.ListProcessTemplatesOption{
		BusinessID:        bizID,
		ServiceTemplateID: module.ServiceTemplateID,
	}
	processTemplate, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, processTemplateFilter)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed, "force sync service instance according to service template: %d, but list process template failed, err: %v", module.ServiceTemplateID, err)
		return
	}
	processTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	for idx, t := range processTemplate.Info {
		processTemplateMap[t.ID] = &processTemplate.Info[idx]
	}

	// step2:
	// find all the process instances relations for the usage of getting process instances.
	relationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: &serviceInstanceIDs,
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed, "force sync service instance according to service template: %d, but list process template failed, err: %v", module.ServiceTemplateID, err)
		return
	}
	procIDs := make([]int64, 0)
	for _, r := range relations.Info {
		procIDs = append(procIDs, r.ProcessID)
	}

	// step 3:
	// find all the process instance in process instance relation.
	processInstances, err := ps.Logic.ListProcessInstanceWithIDs(ctx.Kit, procIDs)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed, "force sync service instance according to service template: %d, but list process instance: %v failed, err: %v", module.ServiceTemplateID, procIDs, err)
		return
	}
	processInstanceMap := make(map[int64]*metadata.Process)
	for idx, p := range processInstances {
		processInstanceMap[p.ProcessID] = &processInstances[idx]
	}

	// step 4:
	// rearrange the service instance with process instance.
	// {ServiceInstanceID: []Process}
	serviceInstance2ProcessMap := make(map[int64][]*metadata.Process)
	// {ServiceInstanceID: {ProcessTemplateID: true}}
	serviceInstanceWithTemplateMap := make(map[int64]map[int64]bool)
	// {ServiceInstanceID: HostID}
	serviceInstance2HostMap := make(map[int64]int64)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstance2ProcessMap[serviceInstance.ID] = make([]*metadata.Process, 0)
		serviceInstanceWithTemplateMap[serviceInstance.ID] = make(map[int64]bool)
		serviceInstance2HostMap[serviceInstance.ID] = serviceInstance.HostID
	}
	processInstanceWithTemplateMap := make(map[int64]int64)
	for _, r := range relations.Info {
		p, exist := processInstanceMap[r.ProcessID]
		if !exist {
			// something is wrong, but can this process instance,
			// but we can find it in the process instance relation.
			blog.Warnf("force sync service instance according to service template: %d, but can not find the process instance: %d", module.ServiceTemplateID, r.ProcessID)
			continue
		}
		serviceInstance2ProcessMap[r.ServiceInstanceID] = append(serviceInstance2ProcessMap[r.ServiceInstanceID], p)
		processInstanceWithTemplateMap[r.ProcessID] = r.ProcessTemplateID
		serviceInstanceWithTemplateMap[r.ServiceInstanceID][r.ProcessTemplateID] = true
	}

	// step 5:
	// compare the difference between process instance and process template from one service instance to another.
	for svcInstanceID, processes := range serviceInstance2ProcessMap {
		for _, process := range processes {
			processTemplateID := processInstanceWithTemplateMap[process.ProcessID]
			template, exist := processTemplateMap[processTemplateID]
			if exist == false {
				// this process template has already removed form the service template,
				// which means this process instance need to be removed from this service instance
				if err := ps.Logic.DeleteProcessInstance(ctx.Kit, process.ProcessID); err != nil {
					ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "force sync service instance according to service template: %d, but delete process instance: %d with template: %d failed, err: %v", module.ServiceTemplateID, process.ProcessID, template.ID, err)
					return
				}

				// remove process instance relation now.
				deleteOption := metadata.DeleteProcessInstanceRelationOption{}
				deleteOption.ProcessIDs = &[]int64{process.ProcessID}
				if err := ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption); err != nil {
					ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "force sync service instance according to service template: %d, but delete process instance relation: %d with template: %d failed, err: %v", module.ServiceTemplateID, process.ProcessID, template.ID, err)
					return
				}
				continue
			}

			// this process's bounded is still exist, need to check whether this process instance
			// need to be updated or not.
			proc, changed := template.ExtractChangeInfo(process)
			if !changed {
				continue
			}
			if err := ps.Logic.UpdateProcessInstance(ctx.Kit, process.ProcessID, proc); err != nil {
				ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed, "force sync service instance according to service template: %d, service instance: %d, but update process instance with template: %d failed, err: %v, process: %v", module.ServiceTemplateID, svcInstanceID, template.ID, err, proc)
				return
			}
		}
	}

	// step 6:
	// check if a new process is added to the service template.
	// if true, then create a new process instance for every service instance with process template's default value.
	for processTemplateID, processTemplate := range processTemplateMap {
		for svcID, templates := range serviceInstanceWithTemplateMap {
			if _, exist := templates[processTemplateID]; exist == true {
				continue
			}

			// we can not find this process template in all this service instance,
			// which means that a new process template need to be added to this service instance
			newProcessData := processTemplate.NewProcess(bizID, ctx.Kit.SupplierAccount)
			newProcessID, err := ps.Logic.CreateProcessInstance(ctx.Kit, newProcessData)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed, "force sync service instance according to service template: %d, but create process instance with template: %d failed, err: %v", module.ServiceTemplateID, processTemplateID, err)
				return
			}

			relation := &metadata.ProcessInstanceRelation{
				Metadata:          syncOption.Metadata,
				ProcessID:         int64(newProcessID),
				ServiceInstanceID: svcID,
				ProcessTemplateID: processTemplateID,
				HostID:            serviceInstance2HostMap[svcID],
			}

			// create service instance relation, so that the process instance created upper can be related to this service instance.
			_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relation)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed, "force sync service instance according to service template: %d, but create process instance relation with template: %d failed, err: %v", module.ServiceTemplateID, processTemplateID, err)
				return
			}
		}
	}

	// Finally, we do the force sync successfully.
	ctx.RespEntity(nil)
}

func (ps *ProcServer) ListServiceInstancesWithHost(ctx *rest.Contexts) {
	input := new(metadata.ListServiceInstancesWithHostInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"list service instances with host, but parse biz id failed, err: %v", err)
		return
	}

	if input.HostID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"list service instances with host, but got empty host id. input: %+v", err)
		return
	}

	option := metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostID:     input.HostID,
		WithName:   input.WithName,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list service instance failed, bizID: %d, hostID: %d",
			bizID, input.HostID, err)
		return
	}

	ctx.RespEntity(instances)
}

func (ps *ProcServer) AddProcessInstanceToServiceInstance(ctx *rest.Contexts) {
	input := new(metadata.ListServiceInstancesWithHostInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"list service instances with host, but parse biz id failed, err: %v", err)
		return
	}

	if input.HostID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"list service instances with host, but got empty host id. input: %+v", err)
		return
	}

	option := metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostID:     input.HostID,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list service instance failed, bizID: %d, hostID: %d",
			bizID, input.HostID, err)
		return
	}

	ctx.RespEntity(instances)
}

func (ps *ProcServer) ListProcessInstances(ctx *rest.Contexts) {
	input := new(metadata.ListProcessInstancesOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"list process instances with host, but parse biz id failed, err: %+v", err)
		return
	}

	// list process instance relation
	relationOption := metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: &[]int64{input.ServiceInstanceID},
	}
	relationsResult, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &relationOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list process instance relation failed, bizID: %d, serviceInstanceID: %d, err: %+v",
			bizID, input.ServiceInstanceID, err)
		return
	}

	processIDs := make([]int64, 0)
	for _, relation := range relationsResult.Info {
		processIDs = append(processIDs, relation.ProcessID)
	}

	cond := condition.CreateCondition()
	cond.Field(common.BKProcessIDField).In(processIDs)
	reqParam := new(metadata.QueryCondition)
	reqParam.Condition = cond.ToMapStr()
	processResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != err {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list process instance property failed, bizID: %d, processIDs: %+v, err: %+v", bizID, processIDs, err)
		return
	}

	processIDPropertyMap := map[int64]mapstr.MapStr{}
	for _, process := range processResult.Data.Info {
		processIDVal, exist := process.Get(common.BKProcessIDField)
		if exist == false {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "list process instance failed, parse bk_process_id from process property failed, field not exist, bizID: %d, processIDs: %+v", bizID, processIDs)
		}
		processID, err := util.GetInt64ByInterface(processIDVal)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "list process instance failed, parse bk_process_id from process property failed, parse field to int64 failed, bizID: %d, processIDs: %+v, process: %+v, err: %+v", bizID, processIDs, process, err)
		}
		processIDPropertyMap[processID] = process
	}

	processInstanceList := make([]metadata.ProcessInstance, 0)
	for _, relation := range relationsResult.Info {
		processInstance := metadata.ProcessInstance{
			Property: nil,
			Relation: relation,
		}
		process, exist := processIDPropertyMap[relation.ProcessID]
		if exist == true {
			processInstance.Property = process
		}
		processInstanceList = append(processInstanceList, processInstance)
	}

	ctx.RespEntity(processInstanceList)
}

func (ps *ProcServer) RemoveTemplateBindingOnModule(ctx *rest.Contexts) {
	input := new(metadata.RemoveTemplateBindingOnModuleOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "remove template binding on module failed, parse business id failed, err: %+v", err)
		return
	}

	module, err := ps.getModule(ctx, input.ModuleID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrTopoGetModuleFailed, "create service instance failed, get module failed, moduleID: %d, err: %v", input.ModuleID, err)
		return
	}
	if module.BizID != bizID {
		err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
		ctx.RespWithError(err, common.CCErrCommNotFound, "create service instance failed, get module failed, moduleID: %d, err: %v", input.ModuleID, err)
		return
	}

	response, err := ps.CoreAPI.CoreService().Process().RemoveTemplateBindingOnModule(ctx.Kit.Ctx, ctx.Kit.Header, input.ModuleID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcRemoveTemplateBindingOnModule, "remove template binding on module failed, parse business id failed, err: %+v", err)
		return
	}
	ctx.RespEntity(response)
}
