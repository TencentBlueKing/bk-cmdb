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
	"context"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/selector"
	"configcenter/src/common/util"
)

// createServiceInstances 创建服务实例
// 支持直接创建和通过模板创建，用 module 是否绑定模版信息区分两种情况
func (ps *ProcServer) CreateServiceInstances(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid
	input := new(metadata.CreateServiceInstanceForServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create service instance with template : %d, moduleID: %d, but get business id failed, err: %v", input.ModuleID, err)
			return
		}
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

	if bizID != module.BizID {
		err := ctx.Kit.CCError.Errorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, module.ModuleID, bizID)
		ctx.RespWithError(err, common.CCErrCoreServiceHasModuleNotBelongBusiness, "create service instance failed, module %d not belongs to biz %d, err: %v", input.ModuleID, bizID, err)
		return
	}

	header := ctx.Kit.Header
	tx, e := ps.TransactionClient.Start(context.Background())
	if e != nil {
		blog.Errorf("start transaction failed, err: %+v", e)
		return
	}
	header = tx.TxnInfo().IntoHeader(header)
	ctx.Kit.Header = header

	defer func() {
		if err != nil {
			if txErr := tx.Abort(ctx.Kit.Ctx); txErr != nil {
				blog.Errorf("create service instance failed, abort translation failed, err: %v, rid: %s", txErr, rid)
			}
		} else {
			if txErr := tx.Commit(ctx.Kit.Ctx); txErr != nil {
				blog.Errorf("create service instance failed, transaction commit failed, err: %v, rid: %s", txErr, rid)
			}
		}
	}()

	serviceInstanceIDs := make([]int64, 0)
	for _, inst := range input.Instances {
		instance := &metadata.ServiceInstance{
			BizID:             bizID,
			Name:              input.Name,
			ServiceTemplateID: module.ServiceTemplateID,
			ModuleID:          input.ModuleID,
			HostID:            inst.HostID,
		}

		var serviceInstance *metadata.ServiceInstance
		// create service instance at first
		serviceInstance, err = ps.CoreAPI.CoreService().Process().CreateServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, instance)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service instance failed, moduleID: %d, err: %s", input.ModuleID, err.Error())
			return
		}

		if module.ServiceTemplateID == 0 && len(inst.Processes) > 0 {
			// if this service have process instance to create, then create it now.
			createProcessInput := &metadata.CreateRawProcessInstanceInput{
				BizID:             bizID,
				ServiceInstanceID: serviceInstance.ID,
				Processes:         inst.Processes,
			}
			if _, err = ps.createProcessInstances(ctx, createProcessInput); err != nil {
				ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "create service instance failed, create process instances failed, moduleID: %d, err: %s", input.ModuleID, err.Error())
				return
			}
		}
		if module.ServiceTemplateID == 0 {
			if err = ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstance.ID); err != nil {
				ctx.RespWithError(err, common.CCErrProcReconstructServiceInstanceNameFailed, "create service instance failed, reconstruct service instance name failed, instanceID: %d, err: %s", serviceInstance.ID, err.Error())
				return
			}
		}

		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	ctx.RespEntity(serviceInstanceIDs)
}

func (ps *ProcServer) SearchServiceInstancesInModuleWeb(ctx *rest.Contexts) {
	input := new(metadata.GetServiceInstanceInModuleInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get service instances in module, but parse biz id failed, err: %v", err)
			return
		}
	}

	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		ModuleID:   input.ModuleID,
		Page:       input.Page,
		SearchKey:  input.SearchKey,
		Selectors:  input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in module: %d failed, err: %v", input.ModuleID, err)
		return
	}

	serviceInstanceIDs := make([]int64, 0)
	for _, instance := range instances.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, instance.ID)
	}
	listRelationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: serviceInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, listRelationOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance relations failed, list option: %+v, err: %v", listRelationOption, err)
		return
	}

	// service_instance_id -> process count
	processCountMap := make(map[int64]int)
	for _, relation := range relations.Info {
		if _, ok := processCountMap[relation.ServiceInstanceID]; ok == false {
			processCountMap[relation.ServiceInstanceID] = 0
		}
		processCountMap[relation.ServiceInstanceID] += 1
	}

	// insert `process_count` field
	serviceInstanceDetails := make([]map[string]interface{}, 0)
	for _, instance := range instances.Info {
		item, err := mapstr.Struct2Map(instance)
		if err != nil {
			blog.ErrorJSON("SearchServiceInstancesInModuleWeb failed, Struct2Map failed, serviceInstance: %s, err: %s, rid: %s", instance, err.Error(), ctx.Kit.Rid)
			ccErr := ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed)
			ctx.RespAutoError(ccErr)
			return
		}
		item["process_count"] = 0
		if count, ok := processCountMap[instance.ID]; ok == true {
			item["process_count"] = count
		}
		serviceInstanceDetails = append(serviceInstanceDetails, item)
	}
	result := metadata.MultipleMap{
		Count: instances.Count,
		Info:  serviceInstanceDetails,
	}
	ctx.RespEntity(result)
}

func (ps *ProcServer) SearchServiceInstancesInModule(ctx *rest.Contexts) {
	input := new(metadata.GetServiceInstanceInModuleInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get service instances in module, but parse biz id failed, err: %v", err)
			return
		}
	}

	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		ModuleID:   input.ModuleID,
		Page:       input.Page,
		SearchKey:  input.SearchKey,
		Selectors:  input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in module: %d failed, err: %v", input.ModuleID, err)
		return
	}

	ctx.RespEntity(instances)
}

func (ps *ProcServer) ListServiceInstancesDetails(ctx *rest.Contexts) {
	input := new(metadata.ListServiceInstanceDetailRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil || bizID == 0 {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get service instances in module, but parse biz id failed, err: %v", err)
			return
		}
	}

	option := &metadata.ListServiceInstanceDetailOption{
		BusinessID:         bizID,
		ModuleID:           input.ModuleID,
		SetID:              input.SetID,
		HostID:             input.HostID,
		ServiceInstanceIDs: input.ServiceInstanceIDs,
		Page:               input.Page,
		Selectors:          input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstanceDetail(ctx.Kit.Ctx, ctx.Kit.Header, option)
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

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete service instances, but parse biz id failed, err: %v", err)
			return
		}
	}
	input.BizID = bizID

	// when a service instance is deleted, the related data should be deleted at the same time
	for _, serviceInstanceID := range input.ServiceInstanceIDs {
		serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed, "delete service instance failed, service instance not found, serviceInstanceIDs: %d", serviceInstanceID)
			return
		}
		if serviceInstance.BizID != bizID {
			err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.MetadataField)
			ctx.RespWithError(err, common.CCErrCommParamsInvalid, "delete service instance failed, biz id from input and service instance not equal, serviceInstanceIDs: %d", serviceInstanceID)
			return
		}

		// step1: delete the service instance relation.
		option := &metadata.ListProcessInstanceRelationOption{
			BusinessID:         bizID,
			ServiceInstanceIDs: []int64{serviceInstanceID},
		}
		relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed, "delete service instance: %d, but list service instance relation failed.", serviceInstanceID)
			return
		}

		if len(relations.Info) > 0 {
			deleteOption := metadata.DeleteProcessInstanceRelationOption{
				ServiceInstanceIDs: []int64{serviceInstanceID},
			}
			err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d, but delete service instance relations failed.", serviceInstanceID)
				return
			}

			// step2: delete process instance belongs to this service instance.
			processIDs := make([]int64, 0)
			for _, r := range relations.Info {
				processIDs = append(processIDs, r.ProcessID)
			}
			if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, processIDs); err != nil {
				ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d, but delete process instance failed.", serviceInstanceID)
				return
			}
		}

		// step3: delete service instance.
		deleteOption := &metadata.CoreDeleteServiceInstanceOption{
			BizID:              bizID,
			ServiceInstanceIDs: []int64{serviceInstanceID},
		}
		err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
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
		if len(result.Info) != 0 {
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
	diffOption := new(metadata.DiffModuleWithTemplateOption)
	if err := ctx.DecodeInto(diffOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// why we need validate metadata here?
	bizID := diffOption.BizID
	if bizID == 0 && diffOption.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*diffOption.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "find difference between service template and process instances, but parse biz id failed, err: %v", err)
			return
		}
	}
	diffOption.BizID = bizID

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
	for idx, pTemplate := range processTemplates.Info {
		pTemplateMap[pTemplate.ID] = &processTemplates.Info[idx]
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

	// step 5:
	// construct map {ServiceInstanceID ==> []ProcessInstanceRelation}
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstances.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}
	option := metadata.ListProcessInstanceRelationOption{
		BusinessID:         module.BizID,
		ServiceInstanceIDs: serviceInstanceIDs,
	}

	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
			"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get service instance relations failed, err: %v",
			module.ServiceTemplateID, module.BizID, diffOption.ModuleID, err)
		return
	}
	serviceRelationMap := make(map[int64][]metadata.ProcessInstanceRelation)
	for _, r := range relations.Info {
		serviceRelationMap[r.ServiceInstanceID] = append(serviceRelationMap[r.ServiceInstanceID], r)
	}

	// step 5: compare the process instance with it's process template one by one in a service instance.
	type recorder struct {
		ProcessID        int64
		ProcessName      string
		Process          *metadata.Process
		ServiceInstance  *metadata.ServiceInstance
		ChangedAttribute []metadata.ProcessChangedAttribute
	}
	removed := make(map[string][]recorder)
	changed := make(map[int64][]recorder)
	unchanged := make(map[int64][]recorder)
	added := make(map[int64][]recorder)
	processTemplateReferenced := make(map[int64]int64)
	for idx, serviceInstance := range serviceInstances.Info {
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
			processName := ""
			if process.ProcessName != nil {
				processName = *process.ProcessName
			}
			property, exist := pTemplateMap[relation.ProcessTemplateID]
			if !exist {
				// process's template doesn't exist means the template has already been removed.
				removed[processName] = append(removed[processName], recorder{
					ProcessID:       relation.ProcessID,
					Process:         process,
					ProcessName:     processName,
					ServiceInstance: &serviceInstances.Info[idx],
				})
				continue
			}

			changedAttributes := ps.Logic.DiffWithProcessTemplate(property.Property, process, attributeMap)
			if len(changedAttributes) == 0 {
				// nothing changed
				unchanged[relation.ProcessTemplateID] = append(unchanged[relation.ProcessTemplateID], recorder{
					ProcessID:       relation.ProcessID,
					ProcessName:     processName,
					ServiceInstance: &serviceInstances.Info[idx],
				})
				continue
			}

			// something has already changed.
			changed[relation.ProcessTemplateID] = append(changed[relation.ProcessTemplateID], recorder{
				ProcessID:        relation.ProcessID,
				ProcessName:      processName,
				ServiceInstance:  &serviceInstances.Info[idx],
				ChangedAttribute: changedAttributes,
			})
		}

		// check whether a new process template has been added.
		for templateID, processTemplate := range pTemplateMap {
			if _, exist := processTemplateReferenced[templateID]; exist == true {
				continue
			}
			// the process template does not exist in all the service instances,
			// which means a new process template is added.
			record := recorder{
				ProcessName:     processTemplate.ProcessName,
				ServiceInstance: &serviceInstances.Info[idx],
			}
			added[templateID] = append(added[templateID], record)
		}
	}

	// it's time to rearrange the data
	moduleDifference := metadata.ModuleDiffWithTemplateDetail{
		Unchanged:     make([]metadata.ServiceInstanceDifference, 0),
		Changed:       make([]metadata.ServiceInstanceDifference, 0),
		Added:         make([]metadata.ServiceInstanceDifference, 0),
		Removed:       make([]metadata.ServiceInstanceDifference, 0),
		HasDifference: false,
	}

	for _, records := range removed {
		if len(records) == 0 {
			continue
		}
		processTemplateName := records[0].ProcessName

		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for idx := range records {
			item := metadata.ServiceDifferenceDetails{
				ServiceInstance: *records[idx].ServiceInstance,
				Process:         records[idx].Process,
			}
			serviceInstances = append(serviceInstances, item)
		}
		moduleDifference.Removed = append(moduleDifference.Removed, metadata.ServiceInstanceDifference{
			ProcessTemplateID:    0,
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
		moduleDifference.Unchanged = append(moduleDifference.Unchanged, metadata.ServiceInstanceDifference{
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
		moduleDifference.Changed = append(moduleDifference.Changed, metadata.ServiceInstanceDifference{
			ProcessTemplateID:    changedID,
			ProcessTemplateName:  records[0].ProcessName,
			ServiceInstanceCount: len(serviceInstances),
			ServiceInstances:     serviceInstances,
		})
	}

	for addedID, records := range added {
		sInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, s := range records {
			sInstances = append(sInstances, metadata.ServiceDifferenceDetails{ServiceInstance: *s.ServiceInstance})
		}

		moduleDifference.Added = append(moduleDifference.Added, metadata.ServiceInstanceDifference{
			ProcessTemplateID:    addedID,
			ProcessTemplateName:  pTemplateMap[addedID].ProcessName,
			ServiceInstanceCount: len(sInstances),
			ServiceInstances:     sInstances,
		})
	}

	moduleChangedAttributes, err := ps.CalculateModuleAttributeDifference(ctx.Kit.Ctx, ctx.Kit.Header, *module)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed,
			"get difference between with module and service template failed, diff module attributes failed, moduleID: %d, %v", module.ModuleID, err)
		return
	}
	moduleDifference.ChangedAttributes = moduleChangedAttributes

	if len(moduleDifference.Added) > 0 ||
		len(moduleDifference.Changed) > 0 ||
		len(moduleDifference.Removed) > 0 ||
		len(moduleDifference.ChangedAttributes) > 0 {
		moduleDifference.HasDifference = true
	}

	ctx.RespEntity(moduleDifference)
}

func (ps *ProcServer) CalculateModuleAttributeDifference(ctx context.Context, header http.Header, module metadata.ModuleInst) ([]metadata.ModuleChangedAttribute, errors.CCErrorCoder) {
	rid := util.ExtractRequestIDFromContext(ctx)

	changedAttributes := make([]metadata.ModuleChangedAttribute, 0)
	if module.ServiceTemplateID == common.ServiceTemplateIDNotSet {
		return changedAttributes, nil
	}
	serviceTpl, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(ctx, header, module.ServiceTemplateID)
	if err != nil {
		return nil, err
	}

	// just for better performance
	if module.ServiceCategoryID == serviceTpl.ServiceCategoryID &&
		module.ModuleName == serviceTpl.Name {
		return changedAttributes, nil
	}

	// find process object's attribute
	filter := &metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKObjIDField: common.BKInnerObjIDModule,
		}),
	}
	attrResult, e := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx, header, common.BKInnerObjIDProc, filter)
	if e != nil {
		blog.Errorf("read module attributes failed, filter: %+v, err: %+v, rid: %s", rid)
		return nil, errors.New(common.CCErrCommDBSelectFailed, "db select failed")
	}
	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}
	if module.ServiceCategoryID != serviceTpl.ServiceCategoryID {
		field := "service_category_id"
		changedAttribute := metadata.ModuleChangedAttribute{
			ID:                    attributeMap[field].ID,
			PropertyID:            field,
			PropertyName:          attributeMap[field].PropertyName,
			PropertyValue:         module.ServiceCategoryID,
			TemplatePropertyValue: serviceTpl.ServiceCategoryID,
		}
		changedAttributes = append(changedAttributes, changedAttribute)
	}
	if module.ModuleName != serviceTpl.Name {
		field := "bk_module_name"
		changedAttribute := metadata.ModuleChangedAttribute{
			ID:                    attributeMap[field].ID,
			PropertyID:            field,
			PropertyName:          attributeMap[field].PropertyName,
			PropertyValue:         module.ModuleName,
			TemplatePropertyValue: serviceTpl.Name,
		}
		changedAttributes = append(changedAttributes, changedAttribute)
	}
	return changedAttributes, nil
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

	bizID := syncOption.BizID
	if bizID == 0 && syncOption.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*syncOption.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "force sync service instance according to service template, but parse biz id failed, err: %v", err)
			return
		}
	}
	syncOption.BizID = bizID

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
		ServiceInstanceIDs: serviceInstanceIDs,
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
			blog.Warnf("force sync service instance according to service template: %d, but can not find the process instance: %d, rid: %s", module.ServiceTemplateID, r.ProcessID, ctx.Kit.Rid)
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
				deleteOption.ProcessIDs = []int64{process.ProcessID}
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
				BizID:             bizID,
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

	// reconstruct service instance's name as it's dependence(first process's + first process's port) changed
	for _, svcInstanceID := range serviceInstanceIDs {
		if err := ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, svcInstanceID); err != nil {
			ctx.RespWithError(err, common.CCErrProcReconstructServiceInstanceNameFailed, "sync service instance failed, reconstruct service instance name failed, instanceID: %d, err: %s", svcInstanceID, err.Error())
			return
		}
	}

	// get service template
	serviceTemplate, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, module.ServiceTemplateID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// step 7:
	// update module service category and name field
	moduleUpdateOption := &metadata.UpdateOption{
		Data: map[string]interface{}{
			common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
			common.BKModuleNameField:        serviceTemplate.Name,
		},
		Condition: map[string]interface{}{
			common.BKModuleIDField: module.ModuleID,
		},
	}
	resp, e := ps.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, moduleUpdateOption)
	if e != nil {
		ctx.RespWithError(e, common.CCErrTopoModuleUpdateFailed, "sync module field failed, http failed, option: %+v", moduleUpdateOption)
		return
	}
	if resp.Result == false || resp.Code != 0 {
		err := ctx.Kit.CCError.New(resp.Code, resp.ErrMsg)
		ctx.RespWithError(err, common.CCErrTopoModuleUpdateFailed, "sync module service category and name failed, option: %+v", moduleUpdateOption)
		return
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

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service instances with host, but parse biz id failed, err: %v", err)
			return
		}
	}
	input.BizID = bizID

	if input.HostID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service instances with host, but got empty host id. input: %+v", input)
		return
	}

	option := metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostID:     input.HostID,
		SearchKey:  input.SearchKey,
		Page:       input.Page,
		Selectors:  input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list service instance failed, bizID: %d, hostID: %d", bizID, input.HostID, err)
		return
	}

	ctx.RespEntity(instances)
}

// ListServiceInstancesWithHostWeb will return topo level info for each service instance
// api only for web frontend
func (ps *ProcServer) ListServiceInstancesWithHostWeb(ctx *rest.Contexts) {
	input := new(metadata.ListServiceInstancesWithHostInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service instances with host, but parse biz id failed, err: %v", err)
			return
		}
	}

	if input.HostID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service instances with host, but got empty host id. input: %+v", input)
		return
	}

	option := metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostID:     input.HostID,
		SearchKey:  input.SearchKey,
		Page:       input.Page,
		Selectors:  input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list service instance failed, bizID: %d, hostID: %d", bizID, input.HostID, err)
		return
	}

	topoRoot, e := ps.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx, ctx.Kit.Header, bizID, false)
	if e != nil {
		blog.Errorf("ListServiceInstancesWithHostWeb failed, search mainline instance topo failed, bizID: %d, err: %+v, riz: %s", bizID, e, ctx.Kit.Rid)
		err := ctx.Kit.CCError.Errorf(common.CCErrTopoMainlineSelectFailed)
		ctx.RespAutoError(err)
		return
	}

	serviceInstances := make([]metadata.ServiceInstanceWithTopoPath, 0)
	for _, instance := range instances.Info {
		topoPath := topoRoot.TraversalFindModule(instance.ModuleID)
		nodes := make([]metadata.TopoInstanceNodeSimplify, 0)
		for _, topoNode := range topoPath {
			node := metadata.TopoInstanceNodeSimplify{
				ObjectID:     topoNode.ObjectID,
				InstanceID:   topoNode.InstanceID,
				InstanceName: topoNode.InstanceName,
			}
			nodes = append(nodes, node)
		}
		serviceInstance := metadata.ServiceInstanceWithTopoPath{
			ServiceInstance: instance,
			TopoPath:        nodes,
		}
		serviceInstances = append(serviceInstances, serviceInstance)
	}

	result := map[string]interface{}{
		"count": instances.Count,
		"info":  serviceInstances,
	}
	ctx.RespEntity(result)
}

func (ps *ProcServer) ServiceInstanceAddLabels(ctx *rest.Contexts) {
	option := selector.LabelAddOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := ps.CoreAPI.CoreService().Label().AddLabel(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameServiceInstance, option); err != nil {
		ctx.RespWithError(err, common.CCErrCommDBUpdateFailed, "ServiceInstanceAddLabels failed, option: %+v, err: %v", option, err)
		return
	}
	ctx.RespEntity(nil)
}

func (ps *ProcServer) ServiceInstanceRemoveLabels(ctx *rest.Contexts) {
	option := selector.LabelRemoveOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if err := ps.CoreAPI.CoreService().Label().RemoveLabel(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameServiceInstance, option); err != nil {
		ctx.RespWithError(err, common.CCErrCommDBUpdateFailed, "ServiceInstanceRemoveLabels failed, option: %+v, err: %v", option, err)
		return
	}
	ctx.RespEntity(nil)
}

// ServiceInstanceLabelsAggregation aggregation instance's labels
func (ps *ProcServer) ServiceInstanceLabelsAggregation(ctx *rest.Contexts) {
	option := metadata.LabelAggregationOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := option.BizID
	if bizID == 0 {
		var err error
		bizID, err = option.Metadata.ParseBizID()
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
	}

	if bizID == 0 {
		ctx.RespErrorCodeF(common.CCErrCommParamsIsInvalid, "list service instance label, but got invalid biz id: 0", "bk_biz_id")
		return
	}

	listOption := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
	}
	if option.ModuleID != nil {
		listOption.ModuleID = *option.ModuleID
	}
	instanceRst, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	// TODO: how to move aggregation into label service
	aggregationData := make(map[string][]string)
	for _, inst := range instanceRst.Info {
		for key, value := range inst.Labels {
			if _, exist := aggregationData[key]; exist == false {
				aggregationData[key] = make([]string, 0)
			}
			aggregationData[key] = append(aggregationData[key], value)
		}
	}
	for key := range aggregationData {
		aggregationData[key] = util.StrArrayUnique(aggregationData[key])
	}
	ctx.RespEntity(aggregationData)
}
