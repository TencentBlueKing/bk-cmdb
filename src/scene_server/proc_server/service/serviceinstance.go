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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"strconv"
)

func (ps *ProcServer) CreateServiceInstancesWithRaw(ctx *rest.Contexts) {
	input := new(metadata.CreateServiceInstanceForServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"create service instance with raw , moduleID: %d, but get business id failed, err: %v", input.ModuleID, err)
		return
	}

	ps.createServiceInstances(ctx, input)
}

func (ps *ProcServer) CreateProcessInstancesWithRaw(ctx *rest.Contexts) {
	input := new(metadata.CreateRawProcessInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"create process instance with raw , but get business id failed, err: %v", err)
		return
	}

	ps.createProcessInstancesRaw(ctx, input)
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
		data := mapstr.NewFromStruct(process, "field")
		data.Remove(common.BKProcessIDField)
		data.Remove(common.MetadataField)
		data.Remove(common.LastTimeField)
		data.Remove(common.CreateTimeField)
		err := ps.Logic.UpdateProcessInstance(ctx.Kit, processID, data)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed, "update process failed, processID: %d, process: %+v, err: %v", process.ProcessID, process, err)
			return
		}
	}

	ctx.RespEntity(processIDs)
}

func (ps *ProcServer) CreateServiceInstancesWithTemplate(ctx *rest.Contexts) {
	input := new(metadata.CreateServiceInstanceForServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create service instance with template : %d, moduleID: %d, but get business id failed, err: %v", input.TemplateID, input.ModuleID, err)
		return
	}

	ps.createServiceInstances(ctx, input)
}

// create service instance batch, which must belongs to a same module and service template.
// if needed, it also create process instance for a service instance at the same time.
func (ps *ProcServer) createServiceInstances(ctx *rest.Contexts, input *metadata.CreateServiceInstanceForServiceTemplateInput) {

	serviceInstanceIDs := make([]int64, 0)
	for _, inst := range input.Instances {
		instance := &metadata.ServiceInstance{
			Metadata:          input.Metadata,
			Name:              input.Name,
			ServiceTemplateID: input.TemplateID,
			ModuleID:          input.ModuleID,
			HostID:            inst.HostID,
		}

		// create service instance at first
		temp, err := ps.CoreAPI.CoreService().Process().CreateServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, instance)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed,
				"create service instance for template: %d, moduleID: %d, failed, err: %v",
				input.TemplateID, input.ModuleID, err)
			return
		}

		// if this service have process instance to create, then create it now.
		for _, detail := range inst.Processes {
			id, err := ps.Logic.CreateProcessInstance(ctx.Kit, &detail.ProcessInfo)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
					"create service instance, for template: %d, moduleID: %d, but create process failed, err: %v",
					input.TemplateID, input.ModuleID, err)
				return
			}

			relation := &metadata.ProcessInstanceRelation{
				Metadata:          input.Metadata,
				ProcessID:         int64(id),
				ProcessTemplateID: detail.ProcessTemplateID,
				ServiceInstanceID: temp.ID,
				HostID:            inst.HostID,
			}

			_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relation)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
					"create service instance relations, for template: %d, moduleID: %d, err: %v",
					input.TemplateID, input.ModuleID, err)
				return
			}
		}

		serviceInstanceIDs = append(serviceInstanceIDs, temp.ID)
	}

	ctx.RespEntity(serviceInstanceIDs)
}

func (ps *ProcServer) validateRawInstanceUnique(ctx *rest.Contexts, serviceInstanceID int64, processInfo *metadata.Process) errors.CCError {
	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, get service instance failed, metadata: %+v, err: %v, rid: %s", serviceInstance.Metadata, err, ctx.Kit.Rid)
		return err
	}

	// find process under service instance
	bizID, err := metadata.BizIDFromMetadata(serviceInstance.Metadata)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, parse business id from metadata failed, metadata: %+v, err: %v, rid: %s", serviceInstance.Metadata, err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommParseBizIDFromMetadataInDBFailed)
	}
	relationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: &[]int64{serviceInstance.ID},
		ProcessTemplateID:  common.ServiceTemplateIDNotSet,
		HostID:             serviceInstance.ID,
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
		otherProcessIDs := make([]int64, 0)
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
	listResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKProcessObjectName, processNameFilterCond)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, search process with bk_process_name failed, filter: %+v, err: %v, rid: %s", processNameFilter, err, ctx.Kit.Rid)
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
		common.BKStartParamRegex: processInfo.ProcessName,
		common.BKFuncName:        processInfo.FuncName,
	}
	funcNameFilterCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(funcNameFilter),
	}
	listFuncNameResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKProcessObjectName, funcNameFilterCond)
	if err != nil {
		blog.Errorf("validateRawInstanceUnique failed, search process with func name failed, filter: %+v, err: %v, rid: %s", funcNameFilterCond, err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	if listFuncNameResult.Data.Count > 0 {
		blog.Errorf("validateRawInstanceUnique failed, bk_func_name and bk_start_param_regex duplicated under service instance, err: %v, rid: %s", err, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCoreServiceFuncNameDuplicated)
	}
	return nil
}

func (ps *ProcServer) createProcessInstancesRaw(ctx *rest.Contexts, input *metadata.CreateRawProcessInstanceInput) {
	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstanceID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
			"create process instance failed, get service instance by id failed, serviceInstanceID: %d, err: %v",
			input.ServiceInstanceID, err)
		return
	}
	if serviceInstance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		ctx.RespWithError(err, common.CCErrProcEditProcessInstanceCreateByTemplateForbidden,
			"create process instance failed, create process instance on service instance initialized by template forbidden, serviceInstanceID: %d, err: %v",
			input.ServiceInstanceID, err)
		return
	}

	processIDs := make([]int64, 0)
	for _, process := range input.Processes {
		process.ProcessInfo.ProcessID = 0
		if err := ps.validateRawInstanceUnique(ctx, serviceInstance.ID, &process.ProcessInfo); err != nil {
			ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
				"create process instance failed, serviceInstanceID: %d, process: %+v, err: %v",
				input.ServiceInstanceID, process, err)
			return
		}

		processID, err := ps.Logic.CreateProcessInstance(ctx.Kit, &process.ProcessInfo)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
				"create process instance failed, create process failed, serviceInstanceID: %d, process: %+v, err: %v",
				input.ServiceInstanceID, process, err)
			return
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
			ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
				"create service instance relations, create process instance relation failed, serviceInstanceID: %d, relation: %+v, err: %v",
				input.ServiceInstanceID, relation, err)
			return
		}
		processIDs = append(processIDs, processID)
	}

	ctx.RespEntity(processIDs)
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

func (ps *ProcServer) GetServiceInstancesInModule(ctx *rest.Contexts) {
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
	// when a service instance is deleted, the related data should be deleted at the same time:
	// 1. service instance relation need to be deleted.
	// 2. process instance belongs to this service instance should be deleted.

	for _, serviceInstanceID := range input.ServiceInstanceIDs {
		serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed,
				"delete service instance failed, service instance not found, serviceInstanceIDs: %d", serviceInstanceID)
			return
		}

		// Firstly, delete the service instance relation.
		option := &metadata.ListProcessInstanceRelationOption{
			BusinessID:         bizID,
			ServiceInstanceIDs: &[]int64{serviceInstanceID},
		}
		relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
				"delete service instance: %d, but list service instance relation failed.", serviceInstanceID)
			return
		}

		deleteOption := metadata.DeleteProcessInstanceRelationOption{
			ServiceInstanceIDs: &[]int64{serviceInstanceID},
		}
		err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed,
				"delete service instance: %d, but delete service instance relations failed.", serviceInstanceID)
			return
		}

		// Secondly, delete process instance belongs to this service instance.
		var processIDs []int64
		for _, r := range relations.Info {
			processIDs = append(processIDs, r.ProcessID)
		}
		if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, processIDs); err != nil {
			ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed,
				"delete service instance: %d, but delete process instance failed.", serviceInstanceID)
			return
		}

		// Finally, delete service instance.
		deleteSvcInstOption := &metadata.DeleteServiceInstanceOption{
			ServiceInstanceIDs: []int64{serviceInstanceID},
		}
		err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, deleteSvcInstOption)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d failed, err: %v", serviceInstanceID, err)
			return
		}

		// check and move host from module if no serviceInstance on it
		filter := &metadata.ListServiceInstanceOption{
			BusinessID: bizID,
			HostID:     serviceInstance.HostID,
		}
		result, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, filter)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get host related service instances failed, bizID: %d, serviceIntanceID: %d, err: %v", bizID, serviceInstance.HostID, err)
			return
		}

		var moduleHasServiceInstance bool
		for _, instance := range result.Info {
			if instance.ModuleID == serviceInstance.ModuleID {
				moduleHasServiceInstance = true
			}
		}
		if moduleHasServiceInstance == false {
			// just remove host from this module
			removeHostFromModuleOption := metadata.RemoveHostsFromModuleOption{
				ApplicationID: bizID,
				HostID:        serviceInstance.HostID,
				ModuleID:      serviceInstance.ModuleID,
			}
			if _, err := ps.CoreAPI.CoreService().Host().RemoveHostFromModule(ctx.Kit.Ctx, ctx.Kit.Header, &removeHostFromModuleOption); err != nil {
				ctx.RespWithError(err, common.CCErrHostMoveResourcePoolFail, "remove host from module failed, option: %+v, err: %v", removeHostFromModuleOption, err)
				return
			}
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
func (ps *ProcServer) FindDifferencesBetweenServiceAndProcessInstance(ctx *rest.Contexts) {
	input := new(metadata.FindServiceTemplateAndInstanceDifferenceOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "find difference between service template and process instances, but parse biz id failed, err: %v", err)
		return
	}

	// step 1:
	// find process object's attribute
	attrResult, err := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDProc, new(metadata.QueryCondition))
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
			"find difference between service template: %d and process instances, bizID: %d, but get process attributes failed, err: %v",
			input.ServiceTemplateID, bizID, err)
		return
	}

	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}

	// step 2:
	// find all the process template in this service template, for compare usage.
	listProcOption := &metadata.ListProcessTemplatesOption{
		BusinessID:        bizID,
		ServiceTemplateID: input.ServiceTemplateID,
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listProcOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
			"find difference between service template: %d and process instances, bizID: %d, but get process templates failed, err: %v",
			input.ServiceTemplateID, bizID, err)
		return
	}

	// step 3:
	// find process instance's relations, which allows us know the relationship between
	// process instance and it's template, service instance, etc.
	pTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	serviceRelationMap := make(map[int64][]metadata.ProcessInstanceRelation)
	for _, pTemplate := range processTemplates.Info {
		pTemplateMap[pTemplate.ID] = &pTemplate

		option := metadata.ListProcessInstanceRelationOption{
			BusinessID:        bizID,
			ProcessTemplateID: pTemplate.ID,
		}

		relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
				"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance relations failed, err: %v",
				input.ServiceTemplateID, bizID, input.ModuleID, err)
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
		BusinessID:        bizID,
		ServiceTemplateID: input.ServiceTemplateID,
		ModuleID:          input.ModuleID,
	}
	serviceInstances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed,
			"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance failed, err: %v",
			input.ServiceTemplateID, bizID, input.ModuleID, err)
		return
	}

	// step 5: compare the process instance with it's process template one by one in a service instance.
	differences := make([]*metadata.ServiceProcessInstanceDifference, 0)
	processTemplatesUsing := make(map[int64]bool)
	for _, serviceInstance := range serviceInstances.Info {
		// get the process instance relation
		relations := serviceRelationMap[serviceInstance.ID]

		if len(relations) == 0 {
			// There is no relations in this service instance, which means no process instances.
			// Normally, this can not be happy.
			// TODO: what???
			differences = append(differences, &metadata.ServiceProcessInstanceDifference{
				ServiceInstanceID:   serviceInstance.ID,
				ServiceInstanceName: serviceInstance.Name,
				HostID:              serviceInstance.HostID,
				Differences:         metadata.NewDifferenceDetail(),
			})
			continue
		}

		// now, we can compare the differences between process template and process instance.
		diff := &metadata.ServiceProcessInstanceDifference{
			ServiceInstanceID:   serviceInstance.ID,
			ServiceInstanceName: serviceInstance.Name,
			HostID:              serviceInstance.HostID,
			Differences:         metadata.NewDifferenceDetail(),
		}
		for _, r := range relations {
			// remember what process template is using, so that we can check whether a new process template has
			// been added or not.
			processTemplatesUsing[r.ProcessTemplateID] = true

			// find the process instance now.
			processInstance, err := ps.Logic.GetProcessInstanceWithID(ctx.Kit, r.ProcessID)
			if err != nil {
				if err.GetCode() == common.CCErrCommNotFound {
					processInstance = new(metadata.Process)
				} else {
					ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed,
						"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get process instance: %d failed, err: %v",
						input.ServiceTemplateID, bizID, input.ModuleID, r.ProcessID, err)
					return
				}
			}

			// let's check if the process instance bounded process template is still exist in it's service template
			// if not exist, that means that this process has already been removed from service template.
			pTemplate, exist := pTemplateMap[r.ProcessTemplateID]
			if !exist {
				// the process instance's bounded process template has already been removed from this service template.
				diff.Differences.Removed = append(diff.Differences.Removed, metadata.ProcessDifferenceDetail{
					ProcessTemplateID: r.ProcessTemplateID,
					ProcessInstance:   *processInstance,
				})
				differences = append(differences, diff)
				continue
			}

			diff := &metadata.ServiceProcessInstanceDifference{
				ServiceInstanceID:   serviceInstance.ID,
				ServiceInstanceName: serviceInstance.Name,
				HostID:              serviceInstance.HostID,
				Differences:         metadata.NewDifferenceDetail(),
			}

			if pTemplate.Property == nil {
				continue
			}

			diffAttributes := ps.Logic.DiffWithProcessTemplate(pTemplate.Property, processInstance, attributeMap)
			if len(diffAttributes) == 0 {
				// the process instance's value is exactly same with the process template's value
				diff.Differences.Unchanged = append(diff.Differences.Unchanged, metadata.ProcessDifferenceDetail{
					ProcessTemplateID: pTemplate.ID,
					ProcessInstance:   *processInstance,
				})
			} else {
				// the process instance's value is not same with the process template's value
				diff.Differences.Changed = append(diff.Differences.Changed, metadata.ProcessDifferenceDetail{
					ProcessTemplateID: pTemplate.ID,
					ProcessInstance:   *processInstance,
					ChangedAttributes: diffAttributes,
				})
			}

		}

		// it's time to see whether a new process template has been added.
		for _, t := range processTemplates.Info {
			if _, exist := processTemplatesUsing[t.ID]; exist {
				continue
			}

			// this process template does not exist in this template's all service instances.
			// so it's a new one to be added.
			if t.Property == nil {
				continue
			}
			diff.Differences.Added = append(diff.Differences.Added, metadata.ProcessDifferenceDetail{
				ProcessTemplateID: t.ID,
				ProcessInstance:   *ps.Logic.NewProcessInstanceFromProcessTemplate(t.Property),
			})

		}

		differences = append(differences, diff)
	}

	ctx.RespEntity(differences)
}

// this function works to find differences between the service template and service instances in a module.
// compared to the service template's process template, a process instance in the service instance may
// contains several differences, like as follows:
// unchanged: the process instance's property values are same with the process template it belongs.
// changed: the process instance's property values are not same with the process template it belongs.
// add: a new process template is added, compared to the service instance belongs to this service template.
// deleted: a process is already deleted, compared to the service instance belongs to this service template.
func (ps *ProcServer) DiffServiceInstanceWithTemplate(ctx *rest.Contexts) {
	input := new(metadata.FindServiceTemplateAndInstanceDifferenceOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "find difference between service template and process instances, but parse biz id failed, err: %v", err)
		return
	}

	if input.ServiceTemplateID == 0 || input.ModuleID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "find difference between service template and process instances, but got empty service template id or module id")
		return
	}

	// step 1:
	// find process object's attribute
	cond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKObjIDField: common.BKInnerObjIDProc,
		}),
	}
	attrResult, err := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDProc, cond)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
			"find difference between service template: %d and process instances, bizID: %d, but get process attributes failed, err: %v",
			input.ServiceTemplateID, bizID, err)
		return
	}
	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}

	// step2. get process templates
	listProcessTemplateOption := &metadata.ListProcessTemplatesOption{
		BusinessID:        bizID,
		ServiceTemplateID: input.ServiceTemplateID,
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listProcessTemplateOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed, "find difference between service template: %d and process instances, bizID: %d, but get process templates failed, err: %v", input.ServiceTemplateID, bizID, err)
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
			BusinessID:        bizID,
			ProcessTemplateID: pTemplate.ID,
		}

		relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
				"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance relations failed, err: %v",
				input.ServiceTemplateID, bizID, input.ModuleID, err)
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
		BusinessID:        bizID,
		ServiceTemplateID: input.ServiceTemplateID,
		ModuleID:          input.ModuleID,
	}
	serviceInstances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed,
			"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance failed, err: %v",
			input.ServiceTemplateID, bizID, input.ModuleID, err)
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
	added := make([]int64, 0)
	usedProcessTemplate := make(map[int64]bool)
	for _, serviceInstance := range serviceInstances.Info {
		relations := serviceRelationMap[serviceInstance.ID]

		for _, relation := range relations {
			// record the used process template for checking whether a new process template has been added to service template.
			usedProcessTemplate[relation.ProcessTemplateID] = true

			property, exist := pTemplateMap[relation.ProcessTemplateID]
			if !exist {
				// this process's template is not exist in this service template's,
				// which means this process template has already been removed from the service template.
				removed[relation.ProcessTemplateID] = append(removed[relation.ProcessTemplateID], recorder{
					ProcessID:       relation.ProcessID,
					ServiceInstance: &serviceInstance,
				})
				continue
			}
			// this process instance's template is still exist in the service template.
			// now, we need to check if the process instance's has been changed compared with it's process template
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
		for t := range pTemplateMap {
			if _, exist := usedProcessTemplate[t]; exist == true {
				continue
			}
			// the process template does not exist in all the service instances,
			// which means a new process template is added.
			added = append(added, t)
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
		var processName string
		var gotName bool
		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, record := range records {
			if !gotName {
				process, err := ps.Logic.GetProcessInstanceWithID(ctx.Kit, record.ProcessID)
				if err != nil {
					if err.GetCode() == common.CCErrCommNotFound {
						process = new(metadata.Process)
					} else {
						ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed,
							"get difference between with process template and process instance in a service instance, but get process instance: %d failed, %v", err)
						return
					}
				}
				processName = process.ProcessName
				gotName = true
			}

			serviceInstances = append(serviceInstances, metadata.ServiceDifferenceDetails{ServiceInstance: *record.ServiceInstance})
		}
		differences.Removed = append(differences.Removed, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    removedID,
			ProcessTemplateName:  processName,
			ServiceInstanceCount: len(records),
			ServiceInstances:     serviceInstances,
		})
	}

	for unchangedID, records := range unchanged {
		if len(records) == 0 {
			continue
		}
		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, record := range records {
			serviceInstances = append(serviceInstances, metadata.ServiceDifferenceDetails{ServiceInstance: *record.ServiceInstance})
		}
		differences.Unchanged = append(differences.Unchanged, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    unchangedID,
			ProcessTemplateName:  records[0].ProcessName,
			ServiceInstanceCount: len(records),
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
			ServiceInstanceCount: len(records),
			ServiceInstances:     serviceInstances,
		})
	}

	for _, addedID := range added {
		sInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, s := range serviceInstances.Info {
			sInstances = append(sInstances, metadata.ServiceDifferenceDetails{ServiceInstance: s})
		}

		differences.Added = append(differences.Added, metadata.ServiceInstanceDifferenceDetail{
			ProcessTemplateID:    addedID,
			ProcessTemplateName:  pTemplateMap[addedID].ProcessName,
			ServiceInstanceCount: int(serviceInstances.Count),
			ServiceInstances:     sInstances,
		})
	}

	ctx.RespEntity(differences)
}

// Force sync the service instance with it's bounded service template.
// It keeps the processes exactly same with the process template in the service template,
// which means the number of process is same, and the process instance's info is also exactly same.
// It contains several scenarios in a service instance:
// 1. add a new process
// 2. update a process
// 3. removed a process

func (ps *ProcServer) SyncServiceInstanceByTemplate(ctx *rest.Contexts) {
	input := new(metadata.ForceSyncServiceInstanceWithTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"force sync service instance according to service template, but parse biz id failed, err: %v", err)
		return
	}

	// step 0:
	// find service instances
	serviceInstanceOption := &metadata.ListServiceInstanceOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: &input.ServiceInstances,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		WithName: false,
	}
	serviceInstanceResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "sync service instance with template: %d failed, get service instances failed, err: %v", input.ServiceTemplateID, err)
		return
	}
	for _, serviceInstance := range serviceInstanceResult.Info {
		if serviceInstance.ServiceTemplateID != input.ServiceTemplateID {
			ctx.RespWithError(err, common.CCErrCommParamsInvalid, "sync service instance with template: %d failed, instance %d doesn't come from template %d, err: %v", serviceInstance.ID, input.ServiceTemplateID, err)
			return
		}
	}

	// step 1:
	// find all the process template according to the service template id
	option := &metadata.ListProcessTemplatesOption{
		BusinessID:        bizID,
		ServiceTemplateID: input.ServiceTemplateID,
	}
	processTemplate, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
			"force sync service instance according to service template: %d, but list process template failed, err: %v",
			input.ServiceTemplateID, err)
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
		ServiceInstanceIDs: &input.ServiceInstances,
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
			"force sync service instance according to service template: %d, but list process template failed, err: %v",
			input.ServiceTemplateID, err)
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
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed,
			"force sync service instance according to service template: %d, but list process instance: %v failed, err: %v",
			input.ServiceTemplateID, procIDs, err)
		return
	}
	processInstanceMap := make(map[int64]*metadata.Process)
	for idx, p := range processInstances {
		processInstanceMap[p.ProcessID] = &processInstances[idx]
	}

	// step 4:
	// rearrange the service instance with process instance.
	serviceInstanceWithProcessMap := make(map[int64][]*metadata.Process)
	serviceInstanceWithTemplateMap := make(map[int64]map[int64]bool)
	serviceInstanceWithHostMap := make(map[int64]int64)
	processInstanceWithTemplateMap := make(map[int64]int64)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstanceWithTemplateMap[serviceInstance.ID] = make(map[int64]bool)
		serviceInstanceWithHostMap[serviceInstance.ID] = serviceInstance.HostID
		serviceInstanceWithProcessMap[serviceInstance.ID] = make([]*metadata.Process, 0)
	}
	for _, r := range relations.Info {
		p, exist := processInstanceMap[r.ProcessID]
		if !exist {
			// something is wrong, but can this process instance,
			// but we can find it in the process instance relation.
			blog.Warnf("force sync service instance according to service template: %d, but can not find the process instance: %d",
				input.ServiceTemplateID, r.ProcessID)
			continue
		}
		serviceInstanceWithProcessMap[r.ServiceInstanceID] = append(serviceInstanceWithProcessMap[r.ServiceInstanceID], p)
		processInstanceWithTemplateMap[r.ProcessID] = r.ProcessTemplateID
		serviceInstanceWithTemplateMap[r.ServiceInstanceID][r.ProcessTemplateID] = true
	}

	// step 5:
	// compare the difference between process instance and process template from one service instance to another.
	for svcInstanceID, processes := range serviceInstanceWithProcessMap {
		for _, process := range processes {
			template, exist := processTemplateMap[processInstanceWithTemplateMap[process.ProcessID]]
			if !exist {
				// this process template has already removed form the service template,
				// which means this process instance need to be removed from this service instance
				if err := ps.Logic.DeleteProcessInstance(ctx.Kit, process.ProcessID); err != nil {
					ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed,
						"force sync service instance according to service template: %d, but delete process instance: %d with template: %d failed, err: %v",
						input.ServiceTemplateID, process.ProcessID, template.ID, err)
					return
				}

				// remove process instance relation now.
				deleteOption := metadata.DeleteProcessInstanceRelationOption{}
				deleteOption.ProcessIDs = &[]int64{process.ProcessID}
				if err := ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption); err != nil {
					ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed,
						"force sync service instance according to service template: %d, but delete process instance relation: %d with template: %d failed, err: %v",
						input.ServiceTemplateID, process.ProcessID, template.ID, err)
				}
				continue
			}

			// this process's bounded is still exist, need to check whether this process instance
			// need to be updated or not.
			proc, changed := ps.Logic.CheckProcessTemplateAndInstanceIsDifferent(template.Property, process)
			if !changed {
				// nothing is changed.
				continue
			}

			// process template has already changed, this process instance need to be updated.
			if err := ps.Logic.UpdateProcessInstance(ctx.Kit, process.ProcessID, proc); err != nil {
				ctx.RespWithError(err, common.CCErrProcUpdateProcessFailed,
					"force sync service instance according to service template: %d, service instance: %d, but update process instance with template: %d failed, err: %v, process: %v",
					input.ServiceTemplateID, svcInstanceID, template.ID, err, proc)
				return
			}
		}
	}

	// step 6:
	// check if a new process is added to the service template.
	// if true, then create a new process instance for every service instance with process template's default value.
	for id, pt := range processTemplateMap {
		for svcID, templates := range serviceInstanceWithTemplateMap {
			if _, exist := templates[id]; exist {
				// nothing is changed.
				continue
			}

			// we can not find this process template in all this service instance,
			// which means that a new process template need to be added to this service instance
			process, err := ps.Logic.CreateProcessInstance(ctx.Kit, ps.Logic.NewProcessInstanceFromProcessTemplate(pt.Property))
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
					"force sync service instance according to service template: %d, but create process instance with template: %d failed, err: %v",
					input.ServiceTemplateID, id, err)
				return
			}

			relation := &metadata.ProcessInstanceRelation{
				Metadata:          input.Metadata,
				ProcessID:         int64(process),
				ServiceInstanceID: svcID,
				ProcessTemplateID: id,
				HostID:            serviceInstanceWithHostMap[svcID],
			}

			// create service instance relation, so that the process instance created upper can be related to this service instance.
			_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relation)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
					"force sync service instance according to service template: %d, but create process instance relation with template: %d failed, err: %v",
					input.ServiceTemplateID, id, err)
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
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"remove template binding on module failed, parse business id failed, err: %+v", err)
		return
	}
	queryCondition := metadata.QueryCondition{
		Condition: mapstr.New(),
	}
	queryCondition.Condition.Set(common.BKModuleIDField, input.ModuleID)
	queryCondition.Condition.Set(common.BKAppIDField, bizID)
	result, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, &queryCondition)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"remove template binding on module failed, get module failed, err: %+v", err)
		return
	}
	if result.Data.Count == 0 || len(result.Data.Info) == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommNotFound, "remove template binding on module failed, get module result in not found, filter: %+v", queryCondition)
		return
	}
	moduleSimple := struct {
		ServiceTemplateID int64 `field:"service_template_id" bson:"service_template_id" json:"service_template_id"`
		ServiceCategoryID int64 `field:"service_category_id" bson:"service_category_id" json:"service_category_id"`
	}{}
	if err := result.Data.Info[0].ToStructByTag(&moduleSimple, "field"); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParseDBFailed, "remove template binding on module failed, parse module info from db failed, module: %+v, err: %+v", result.Data.Info, err)
		return
	}

	if moduleSimple.ServiceTemplateID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrProcModuleNotBindWithTemplate, "remove template binding on module failed, module doesn't bind with template yet, module: %+v, err: %+v", result.Data.Info, err)
		return
	}

	data := mapstr.New()
	data.Set(common.BKServiceTemplateIDField, common.ServiceTemplateIDNotSet)
	updateOption := metadata.UpdateOption{
		Data:      data,
		Condition: queryCondition.Condition,
	}
	updateResult, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, &updateOption)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPDoRequestFailed, "remove template binding on module failed, reset service_template_id attribute failed, module: %+v, err: %+v", result.Data.Info, err)
		return
	}
	ctx.RespEntity(updateResult)
}
