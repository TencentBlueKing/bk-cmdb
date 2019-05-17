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

	err = p.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstancID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcDeleteServiceInstancesFailed, "delete service instance: %d failed, err: %v", input.ServiceInstancID, err)
		return
	}

	ctx.RespEntity(metadata.NewSuccessResp(nil))
}
