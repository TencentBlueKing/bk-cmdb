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
	"fmt"
	"strconv"
	"strings"
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
		ctx.RespWithError(err, common.CCErrProcCreateProcessFailed, "create process instance failed, serviceInstanceID: %d, input: %+v, err: %+v", input.ServiceInstanceID, input, err)
		return
	}

	if err := ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstanceID); err != nil {
		ctx.RespWithError(err, common.CCErrProcReconstructServiceInstanceNameFailed, "create process instance failed, reconstruct service instance name failed, instanceID: %d, err: %s", input.ServiceInstanceID, err.Error())
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
		blog.Errorf("create process instance failed, get service instance by id failed, serviceInstanceID: %d, err: %v, rid: %s", input.ServiceInstanceID, err, ctx.Kit.Rid)
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
		blog.Errorf("create process instance failed, create process instance on service instance initialized by template forbidden, serviceInstanceID: %d, err: %v, rid: %s", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrProcEditProcessInstanceCreateByTemplateForbidden)
	}

	processIDs := make([]int64, 0)
	for _, process := range input.Processes {
		process.ProcessInfo.ProcessID = int64(0)
		process.ProcessInfo.BusinessID = bizID
		process.ProcessInfo.SupplierAccount = ctx.Kit.SupplierAccount
		now := time.Now()
		process.ProcessInfo.CreateTime = now
		process.ProcessInfo.LastTime = now

		if err := ps.validateRawInstanceUnique(ctx, serviceInstance.ID, &process.ProcessInfo); err != nil {
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
			blog.Errorf("create service instance relations, create process instance relation failed, serviceInstanceID: %d, relation: %+v, err: %v, rid: %s", input.ServiceInstanceID, relation, err, ctx.Kit.Rid)
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
	processIDs = util.IntArrayUnique(processIDs)
	option := &metadata.ListProcessInstanceRelationOption{
		BusinessID: bizID,
		ProcessIDs: processIDs,
		Page:       metadata.BasePage{Limit: common.BKNoLimit},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPDoRequestFailed, "update process instance failed, search process instance relation failed, err: %+v", err)
		return
	}

	// make sure all process valid
	foundProcessIDs := make([]int64, 0)
	serviceInstanceIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		foundProcessIDs = append(foundProcessIDs, relation.ProcessID)
		serviceInstanceIDs = append(serviceInstanceIDs, relation.ServiceInstanceID)
	}
	invalidProcessIDs := make([]string, 0)
	for _, processID := range processIDs {
		if util.InArray(processID, foundProcessIDs) == false {
			invalidProcessIDs = append(invalidProcessIDs, strconv.FormatInt(processID, 10))
		}
	}
	if len(invalidProcessIDs) > 0 {
		msg := fmt.Sprintf("[%s: %s]", common.BKProcessIDField, strings.Join(invalidProcessIDs, ","))
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, msg)
		ctx.RespWithError(err, common.CCErrCommParamsIsInvalid, "update process instance failed, process %+v not found", invalidProcessIDs)
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

		processData := mapstr.MapStr{}
		if relation.ProcessTemplateID == common.ServiceTemplateIDNotSet {
			serviceInstanceID := relation.ServiceInstanceID
			if err := ps.validateRawInstanceUnique(ctx, serviceInstanceID, &process); err != nil {
				ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed, "update process instance failed, serviceInstanceID: %d, process: %+v, err: %v", serviceInstanceID, process, err)
				return
			}
			process.BusinessID = bizID
			process.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
			processBytes, err := json.Marshal(process)
			if err != nil {
				blog.Errorf("UpdateProcessInstances failed, json Marshal process failed, process: %+v, err: %+v, rid: %s", process, err, ctx.Kit.Rid)
				err := ctx.Kit.CCError.CCError(common.CCErrCommJsonEncode)
				ctx.RespWithError(err, common.CCErrCommJsonDecode, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
			}
			if err := json.Unmarshal(processBytes, &processData); nil != err && 0 != len(processBytes) {
				blog.Errorf("UpdateProcessInstances failed, json Unmarshal process failed, processData: %s, err: %+v, rid: %s", processData, err, ctx.Kit.Rid)
				err := ctx.Kit.CCError.CCError(common.CCErrCommJsonDecode)
				ctx.RespWithError(err, common.CCErrCommJsonDecode, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
			}
			processData.Remove(common.BKProcessIDField)
			processData.Remove(common.MetadataField)
			processData.Remove(common.LastTimeField)
			processData.Remove(common.CreateTimeField)
		} else {
			processTemplate, exist = processTemplateMap[relation.ProcessTemplateID]
			if exist == false {
				err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
				ctx.RespWithError(err, common.CCErrCommNotFound, "update process instance failed, process related template not found, relation: %+v, err: %v", relation, err)
				return
			}
			updateData := processTemplate.ExtractInstanceUpdateData(&process)
			processData = mapstr.MapStr(updateData)
		}

		if err := ps.Logic.UpdateProcessInstance(ctx.Kit, process.ProcessID, processData); err != nil {
			ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
			return
		}
	}

	serviceInstanceIDs = util.IntArrayUnique(serviceInstanceIDs)
	for _, svcInstanceID := range serviceInstanceIDs {
		if err := ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, svcInstanceID); err != nil {
			ctx.RespWithError(err, common.CCErrProcReconstructServiceInstanceNameFailed, "update process instance failed, reconstruct service instance name failed, instanceID: %d, err: %s", svcInstanceID, err.Error())
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
	for _, item := range result.Data.Info {
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

func (ps *ProcServer) getModule(ctx *rest.Contexts, moduleID int64) (*metadata.ModuleInst, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	moduleFilter := &metadata.QueryCondition{
		Condition: mapstr.MapStr(filter),
	}
	modules, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, moduleFilter)
	if err != nil {
		blog.Errorf("getModule failed, moduleID: %d, err: %s, rid: %s", moduleID, err.Error(), ctx.Kit.Rid)
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
	if processInfo.ProcessName != nil && (len(*processInfo.ProcessName) == 0 || len(*processInfo.ProcessName) > common.NameFieldMaxLength) {
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessNameField)
	}
	if processInfo.FuncName != nil && (len(*processInfo.FuncName) == 0 || len(*processInfo.ProcessName) > common.NameFieldMaxLength) {
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
		ServiceInstanceIDs: []int64{serviceInstance.ID},
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
		processName := ""
		if processInfo.ProcessName != nil {
			processName = *processInfo.ProcessName
		}
		return ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceProcessNameDuplicated, processName)
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
		startParamRegex := ""
		if processInfo.StartParamRegex != nil {
			startParamRegex = *processInfo.StartParamRegex
		}
		processName := ""
		if processInfo.FuncName != nil {
			processName = *processInfo.FuncName
		}
		return ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceFuncNameDuplicated, processName, startParamRegex)
	}
	return nil
}

func (ps *ProcServer) DeleteProcessInstance(ctx *rest.Contexts) {
	input := new(metadata.DeleteProcessInstanceInServiceInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete process instance in service instance failed, err: %v", err)
		return
	}

	listOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID: bizID,
		ProcessIDs: input.ProcessInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	templateProcessIDs := make([]string, 0)
	serviceInstanceIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		if relation.ProcessTemplateID != common.ServiceTemplateIDNotSet {
			templateProcessIDs = append(templateProcessIDs, strconv.FormatInt(relation.ProcessID, 10))
		}
		serviceInstanceIDs = append(serviceInstanceIDs, relation.ServiceInstanceID)
	}
	if len(templateProcessIDs) > 0 {
		invalidProcesses := strings.Join(templateProcessIDs, ",")
		blog.Errorf("DeleteProcessInstance failed, some process:%s initialized by template, rid: %s", invalidProcesses, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceShouldNotRemoveProcessCreateByTemplate, invalidProcesses)
		ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "delete process instance: %v, but delete instance relation failed.", input.ProcessInstanceIDs)
		return
	}

	// delete process relation at the same time.
	deleteOption := metadata.DeleteProcessInstanceRelationOption{}
	deleteOption.ProcessIDs = input.ProcessInstanceIDs
	err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "delete process instance: %v, but delete instance relation failed.", input.ProcessInstanceIDs)
		return
	}

	if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, input.ProcessInstanceIDs); err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "delete process instance:%v failed, err: %v", input.ProcessInstanceIDs, err)
		return
	}

	serviceInstanceIDs = util.IntArrayUnique(serviceInstanceIDs)
	for _, svcInstanceID := range serviceInstanceIDs {
		if err := ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, svcInstanceID); err != nil {
			ctx.RespWithError(err, common.CCErrProcReconstructServiceInstanceNameFailed, "delete instance failed, reconstruct service instance name failed, serviceInstanceID: %d, err: %s", svcInstanceID, err.Error())
			return
		}
	}

	ctx.RespEntity(nil)
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
		ServiceInstanceIDs: []int64{input.ServiceInstanceID},
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

var UnbindServiceTemplateOnModuleEnable = true

func (ps *ProcServer) RemoveTemplateBindingOnModule(ctx *rest.Contexts) {
	if UnbindServiceTemplateOnModuleEnable == true {
		ctx.RespErrorCodeOnly(common.CCErrProcUnbindModuleServiceTemplateDisabled, "unbind service template from module disabled")
		return
	}

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
