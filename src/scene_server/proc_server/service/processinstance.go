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
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
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
	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var e error
		bizID, e = metadata.BizIDFromMetadata(*input.Metadata)
		if e != nil {
			blog.Errorf("create process instance with raw, parse biz id from metadata failed, err: %+v, rid: %s", e, ctx.Kit.Rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommHTTPInputInvalid, common.MetadataField)
		}
	}

	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstanceID)
	if err != nil {
		blog.Errorf("create process instance failed, get service instance by id failed, serviceInstanceID: %d, err: %v, rid: %s", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, err
	}
	if serviceInstance.BizID != bizID {
		blog.Errorf("create process instance with raw, biz id from input not equal with service instance, rid: %s", ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.MetadataField)
	}
	if serviceInstance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		blog.Errorf("create process instance failed, create process instance on service instance initialized by template forbidden, serviceInstanceID: %d, err: %v, rid: %s", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrProcEditProcessInstanceCreateByTemplateForbidden)
	}

	processIDs := make([]int64, 0)
	for _, item := range input.Processes {
		now := time.Now()
		item.ProcessData[common.BKProcessIDField] = int64(0)
		item.ProcessData[common.BKAppIDField] = bizID
		item.ProcessData[common.BkSupplierAccount] = ctx.Kit.SupplierAccount
		item.ProcessData[common.CreateTimeField] = now
		item.ProcessData[common.LastTimeField] = now

		if err := ps.validateRawInstanceUnique(ctx, serviceInstance.ID, item.ProcessData); err != nil {
			return nil, err
		}

		processID, err := ps.Logic.CreateProcessInstance(ctx.Kit, item.ProcessData)
		if err != nil {
			blog.Errorf("create process instance failed, create process failed, serviceInstanceID: %d, process: %+v, err: %v, rid: %s", input.ServiceInstanceID, item, err, ctx.Kit.Rid)
			return nil, err
		}

		relation := &metadata.ProcessInstanceRelation{
			BizID:             bizID,
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
	input := metadata.UpdateRawProcessInstanceInput{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	result, err := ps.updateProcessInstances(ctx, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (ps *ProcServer) updateProcessInstances(ctx *rest.Contexts, input metadata.UpdateRawProcessInstanceInput) ([]int64, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			blog.Errorf("update process instance failed, parse business id failed, err: %+v, rid: %s", err, rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPInputInvalid)
		}
	}
	processIDs := make([]int64, 0)
	input.Processes = make([]metadata.Process, 0)
	for _, pData := range input.Raw {
		process := metadata.Process{}
		if err := mapstr.DecodeFromMapStr(&process, pData); err != nil {
			blog.ErrorJSON("update process instance failed, unmarshal request body failed, data: %s, err: %s, rid: %s", pData, err.Error(), rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
		}
		input.Processes = append(input.Processes, process)

		if process.ProcessID == 0 {
			blog.Errorf("update process instance failed, process_id invalid, rid: %s", rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
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
		blog.ErrorJSON("update process instance failed, search process instance relation failed, option: %s, err: %+v, rid: %s", option, err, rid)
		return nil, err
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
		blog.Errorf("update process instance failed, process %+v not found", invalidProcessIDs)
		msg := fmt.Sprintf("[%s: %s]", common.BKProcessIDField, strings.Join(invalidProcessIDs, ","))
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, msg)
		return nil, err
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
			blog.ErrorJSON("update process instance failed, get process template failed, processTemplateID: %d, err: %s, rid: %s", relation.ProcessTemplateID, err, rid)
			return nil, err
		}
		processTemplateMap[relation.ProcessTemplateID] = processTemplate
	}

	process2ServiceInstanceMap := make(map[int64]*metadata.ProcessInstanceRelation)
	for _, relation := range relations.Info {
		process2ServiceInstanceMap[relation.ProcessID] = &relation
	}

	var processTemplate *metadata.ProcessTemplate
	for idx, process := range input.Processes {
		// 单独提取需要被重置成 nil 的字段
		raw := input.Raw[idx]
		clearFields := make([]string, 0)
		for key, value := range raw {
			if value == nil {
				clearFields = append(clearFields, key)
			}
		}
		clearFields = metadata.FilterValidFields(clearFields)

		relation, exist := process2ServiceInstanceMap[process.ProcessID]
		if exist == false {
			err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
			blog.ErrorJSON("update process instance failed, process related service instance not found, process: %s, err: %s, rid: %s", process, err, rid)
			return nil, err
		}

		processData := make(map[string]interface{})
		if relation.ProcessTemplateID == common.ServiceTemplateIDNotSet {
			serviceInstanceID := relation.ServiceInstanceID
			process.BusinessID = bizID
			var err error
			processData, err = mapstruct.Struct2Map(process)
			if nil != err {
				blog.Errorf("UpdateProcessInstances failed, json Unmarshal process failed, processData: %s, err: %+v, rid: %s", processData, err, ctx.Kit.Rid)
				return nil, ctx.Kit.CCError.CCError(common.CCErrCommJsonDecode)
			}
			if err := ps.validateRawInstanceUnique(ctx, serviceInstanceID, processData); err != nil {
				blog.Errorf("update process instance failed, serviceInstanceID: %d, process: %+v, err: %v, rid: %s", serviceInstanceID, process, err, rid)
				return nil, err
			}
			delete(processData, common.BKProcessIDField)
			delete(processData, common.MetadataField)
			delete(processData, common.LastTimeField)
			delete(processData, common.CreateTimeField)
		} else {
			processTemplate, exist = processTemplateMap[relation.ProcessTemplateID]
			if exist == false {
				err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
				blog.Errorf("update process instance failed, process related template not found, relation: %+v, err: %v, rid: %s", relation, err, rid)
				return nil, err
			}
			processData = processTemplate.ExtractInstanceUpdateData(&process)
			clearFields = processTemplate.GetEditableFields(clearFields)
		}
		// set field value as nil
		for _, field := range clearFields {
			processData[field] = nil
		}

		if err := ps.Logic.UpdateProcessInstance(ctx.Kit, process.ProcessID, processData); err != nil {
			blog.Errorf("update process failed, processID: %d, process: %+v, err: %v, rid: %s", process.ProcessID, process, err, rid)
			return nil, err
		}
	}

	serviceInstanceIDs = util.IntArrayUnique(serviceInstanceIDs)
	for _, svcInstanceID := range serviceInstanceIDs {
		if err := ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, svcInstanceID); err != nil {
			blog.Errorf("update process instance failed, reconstruct service instance name failed, instanceID: %d, err: %s, rid: %s", svcInstanceID, err.Error(), rid)
			return nil, err
		}
	}

	return processIDs, nil
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
		blog.ErrorJSON("CheckHostInBusiness failed, GetHostModuleRelation failed, filter: %s, err: %s, rid: %s", hostConfigFilter, err.Error(), ctx.Kit.Rid)
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

func (ps *ProcServer) validateRawInstanceUnique(ctx *rest.Contexts, serviceInstanceID int64, processData map[string]interface{}) errors.CCErrorCoder {
	rid := ctx.Kit.Rid

	processInfo := metadata.Process{}
	if err := mapstruct.Decode2Struct(processData, &processInfo); err != nil {
		blog.ErrorJSON("validateRawInstanceUnique failed, Decode2Struct failed, process: %s, err: %s, rid: %s", processData, err.Error(), rid)
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}
	if processInfo.ProcessName != nil && (len(*processInfo.ProcessName) == 0 || len(*processInfo.ProcessName) > common.NameFieldMaxLength) {
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessNameField)
	}
	if processInfo.FuncName != nil && (len(*processInfo.FuncName) == 0 || len(*processInfo.ProcessName) > common.NameFieldMaxLength) {
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFuncName)
	}
	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, get service instance failed, bk_biz_id: %d, err: %v, rid: %s", serviceInstance.BizID, err, rid)
		return err
	}

	// find process under service instance
	bizID := serviceInstance.BizID
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
		blog.Errorf("validateRawInstanceUnique failed, get relation under service instance failed, err: %v, rid: %s", serviceInstance.BizID, err, ctx.Kit.Rid)
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
		blog.Errorf("validateRawInstanceUnique failed, bk_process_name duplicated under service instance, bk_biz_id: %d, err: %v, rid: %s", serviceInstance.BizID, err, ctx.Kit.Rid)
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

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete process instance in service instance failed, err: %v", err)
			return
		}
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
	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var e error
		bizID, e = metadata.BizIDFromMetadata(*input.Metadata)
		if e != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list process instances with host, but parse biz id failed, err: %+v", e)
			return
		}
	}

	if input.ServiceInstanceID == 0 {
		err := ctx.Kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKServiceInstanceIDField)
		ctx.RespAutoError(err)
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
	filter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: processIDs,
		},
	}
	reqParam := &metadata.QueryCondition{
		Condition: filter,
	}
	processResult, ccErr := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDProc, reqParam)
	if nil != ccErr {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list process instance property failed, bizID: %d, processIDs: %+v, err: %+v", bizID, processIDs, ccErr)
		return
	}

	processIDPropertyMap := map[int64]mapstr.MapStr{}
	for _, process := range processResult.Data.Info {
		processIDVal, exist := process.Get(common.BKProcessIDField)
		if exist == false {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "list process instance failed, parse bk_process_id from process property failed, field not exist, bizID: %d, processIDs: %+v", bizID, processIDs)
			return
		}
		processID, err := util.GetInt64ByInterface(processIDVal)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "list process instance failed, parse bk_process_id from process property failed, parse field to int64 failed, bizID: %d, processIDs: %+v, process: %+v, err: %+v", bizID, processIDs, process, err)
			return
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

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "remove template binding on module failed, parse business id failed, err: %+v", err)
			return
		}
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
