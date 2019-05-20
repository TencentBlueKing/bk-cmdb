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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// create service instance batch, which must belongs to a same module and service template.
// if needed, it also create process instance for a service instance at the same time.
func (p *ProcServer) CreateServiceInstances(ctx *rest.Contexts) {
	input := new(metadata.CreateServiceInstanceForServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid,
			"create service instance for template: %d, moduleID: %d, but get business id failed, err: %v",
			input.TemplateID, input.ModuleID, err)
		return
	}

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
		temp, err := p.CoreAPI.CoreService().Process().CreateServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, instance)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed,
				"create service instance for template: %d, moduleID: %d, failed, err: %v",
				input.TemplateID, input.ModuleID, err)
			return
		}

		// if this service have process instance to create, then create it now.
		for _, detail := range inst.Processes {
			id, err := p.Logic.CreateProcessInstance(ctx.Kit, &detail.ProcessInfo)
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

			_, err = p.CoreAPI.CoreService().Process().CreateProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relation)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcCreateProcessFailed,
					"create service instance relations, for template: %d, moduleID: %d, err: %v",
					input.TemplateID, input.ModuleID, err)
				return
			}
		}

		serviceInstanceIDs = append(serviceInstanceIDs, temp.ID)
	}

	ctx.RespEntity(metadata.NewSuccessResp(serviceInstanceIDs))
}

func (p *ProcServer) DeleteProcessInstanceInServiceInstance(ctx *rest.Contexts) {
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

	if err := p.Logic.DeleteProcessInstanceBatch(ctx.Kit, input.ProcessInstanceIDs); err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteProcessFailed, "delete process instance:%v failed, err: %v", input.ProcessInstanceIDs, err)
		return
	}

	ctx.RespEntity(metadata.NewSuccessResp(nil))
}

func (p *ProcServer) GetServiceInstancesInModule(ctx *rest.Contexts) {
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
	}
	instances, err := p.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in module: %d failed, err: %v", input.ModuleID, err)
		return
	}

	ctx.RespEntity(metadata.NewSuccessResp(instances))
}

func (p *ProcServer) DeleteServiceInstance(ctx *rest.Contexts) {
	input := new(metadata.DeleteServiceInstanceOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	_, err := metadata.BizIDFromMetadata(input.Metadata)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete service instances, but parse biz id failed, err: %v", err)
		return
	}

	err = p.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstanceID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d failed, err: %v", input.ServiceInstanceID, err)
		return
	}

	ctx.RespEntity(metadata.NewSuccessResp(nil))
}

// this function works to find differences between the service template and service instances in a module.
// compared to the service template's process template, a process instance in the service instance may
// contains several differences, like as follows:
// unchanged: the process instance's property values are same with the process template it belongs.
// changed: the process instance's property values are not same with the process template it belongs.
// add: a new process template is added, compared to the service instance belongs to this service template.
// deleted: a process is already deleted, compared to the service instance belongs to this service template.
func (p *ProcServer) FindDifferencesBetweenServiceAndProcessInstance(ctx *rest.Contexts) {
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
	attrResult, err := p.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDProc, new(metadata.QueryCondition))
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
	processTemplates, err := p.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listProcOption)
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
			Metadata:          input.Metadata,
			ProcessTemplateID: pTemplate.ID,
		}

		relations, err := p.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &option)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed,
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
	serviceInstances, err := p.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceOption)
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
				// Differences: ???,
			})
			continue
		}

		// now, we can compare the differences between process template and process instance.
		for _, r := range relations {
			// remember what process template is using, so that we can check whether a new process template has
			// been added or not.
			processTemplatesUsing[r.ProcessTemplateID] = true

			// find the process instance now.
			processInstance, err := p.Logic.GetProcessInstanceWithID(ctx.Kit, r.ProcessID)
			if err != nil {
				ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed,
					"find difference between service template: %d and process instances, bizID: %d, moduleID: %d, but get process instance: %d failed, err: %v",
					input.ServiceTemplateID, bizID, input.ModuleID, r.ProcessID, err)
				return
			}

			// let's check if the process instance bounded process template is still exit in it's service template
			// if not exist, that means that this process has already been removed from service template.
			pTemplate, exist := pTemplateMap[r.ProcessTemplateID]
			if !exist {
				// the process instance's bounded process template has already been removed from this service template.
				diff := &metadata.ServiceProcessInstanceDifference{
					ServiceInstanceID:   serviceInstance.ID,
					ServiceInstanceName: serviceInstance.Name,
					HostID:              serviceInstance.HostID,
					Differences: &metadata.DifferenceDetail{
						Removed: &metadata.ProcessDifferenceDetail{
							ProcessInstance: *processInstance,
						},
					},
				}
				differences = append(differences, diff)
				continue
			}

			diff := &metadata.ServiceProcessInstanceDifference{
				ServiceInstanceID:   serviceInstance.ID,
				ServiceInstanceName: serviceInstance.Name,
				HostID:              serviceInstance.HostID,
			}

			if pTemplate.Property == nil {
				continue
			}

			diffAttributes := p.Logic.GetDifferenceInProcessTemplateAndInstance(pTemplate.Property, processInstance, attributeMap)
			if len(diffAttributes) == 0 {
				// the process instance's value is exactly same with the process template's value
				diff.Differences.Unchanged = &metadata.ProcessDifferenceDetail{
					ProcessTemplateID: pTemplate.ID,
					ProcessInstance:   *processInstance,
				}
			} else {
				// the process instance's value is not same with the process template's value
				diff.Differences.Changed = &metadata.ProcessDifferenceDetail{
					ProcessTemplateID: pTemplate.ID,
					ProcessInstance:   *processInstance,
					ChangedAttributes: diffAttributes,
				}
			}

			differences = append(differences, diff)
		}

		// it's time to see if a new process template has been added.
		for _, t := range processTemplates.Info {
			if _, exist := processTemplatesUsing[t.ID]; exist {
				continue
			}

			// this process template does not exist in this template's all service instances.
			// so it's a new one to be added.
			if t.Property == nil {
				continue
			}

			differences = append(differences, &metadata.ServiceProcessInstanceDifference{
				ServiceInstanceID:   serviceInstance.ID,
				ServiceInstanceName: serviceInstance.Name,
				HostID:              serviceInstance.HostID,
				Differences: &metadata.DifferenceDetail{
					Added: &metadata.ProcessDifferenceDetail{
						ProcessTemplateID: t.ID,
						ProcessInstance:   *p.Logic.NewProcessInstanceFromProcessTemplate(t.Property),
					},
				},
			})

		}
	}

	ctx.RespEntity(differences)
}
