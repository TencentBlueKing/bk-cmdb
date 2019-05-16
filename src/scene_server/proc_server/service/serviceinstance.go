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
// if can also create process instance for a service instance at the same time.
func (p *ProcServer) CreateServiceInstancesForTemplate(ctx *rest.Contexts) {
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
		temp, err := p.CoreAPI.CoreService().Process().CreateServiceInstance(ctx.Ctx, ctx.Header, instance)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed,
				"create service instance for template: %d, moduleID: %d, failed, err: %v",
				input.TemplateID, input.ModuleID, err)
			return
		}

		// if this service have process instance to create, then create it now.
		// TODO:
		if len(inst.Processes) != 0 {
			// p.CoreAPI.CoreService().Instance().CreateInstance(ctx.Ctx, ctx.Header, common.BKProcessObjectName)
		}

		serviceInstanceIDs = append(serviceInstanceIDs, temp.ID)
	}

	ctx.RespEntity(metadata.NewSuccessResp(serviceInstanceIDs))
}
