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
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/selector"
	"configcenter/src/common/util"
)

// createServiceInstances 创建服务实例
// 支持直接创建和通过模板创建，用 module 是否绑定模版信息区分两种情况
// 通过模板创建时，进程信息则表现为更新
func (ps *ProcServer) CreateServiceInstances(ctx *rest.Contexts) {
	input := metadata.CreateServiceInstanceForServiceTemplateInput{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var serviceInstanceIDs []int64
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		serviceInstanceIDs, err = ps.createServiceInstances(ctx, input)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(serviceInstanceIDs)
}

func (ps *ProcServer) createServiceInstances(ctx *rest.Contexts, input metadata.CreateServiceInstanceForServiceTemplateInput) ([]int64, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			blog.Errorf("createServiceInstances failed, parse bizID failed, metadata: %+v, err: %s, rid: %s", input.Metadata, err.Error(), rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommHTTPInputInvalid, common.MetadataField)
		}
	}

	// check hosts in business
	hostIDs := make([]int64, 0)
	hostIDHit := make(map[int64]bool)
	for _, instance := range input.Instances {
		if !util.InArray(instance.HostID, hostIDs) {
			hostIDs = append(hostIDs, instance.HostID)
			hostIDHit[instance.HostID] = false
		}
	}
	if err := ps.CheckHostInBusiness(ctx, bizID, hostIDs); err != nil {
		blog.ErrorJSON("createServiceInstances failed, CheckHostInBusiness failed, bizID: %s, hostIDs: %s, err: %s, rid: %s", bizID, hostIDs, err.Error(), rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotBelongBusiness, hostIDs, bizID)
	}

	module, err := ps.getModule(ctx, input.ModuleID)
	if err != nil {
		blog.Errorf("createServiceInstances failed, get module failed, moduleID: %d, err: %v, rid: %s", input.ModuleID, err, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, err.Error())
	}

	if bizID != module.BizID {
		blog.Errorf("createServiceInstances failed, module %d not belongs to biz %d, rid: %s", input.ModuleID, bizID, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, module.ModuleID, bizID)
	}

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
			blog.ErrorJSON("createServiceInstances failed, core service CreateServiceInstance failed, option: %s, err: %s, rid: %s", instance, err.Error(), rid)
			return nil, err
		}

		if len(inst.Processes) > 0 {
			if module.ServiceTemplateID == 0 {
				// if this service have process instance to create, then create it now.
				createProcessInput := &metadata.CreateRawProcessInstanceInput{
					BizID:             bizID,
					ServiceInstanceID: serviceInstance.ID,
					Processes:         inst.Processes,
				}
				if _, err = ps.createProcessInstances(ctx, createProcessInput); err != nil {
					blog.ErrorJSON("createServiceInstances failed, createProcessInstances failed, input: %s, err: %s, rid: %s", createProcessInput, err.Error(), rid)
					return nil, err
				}
			} else {
				// update process instance by templateID
				relationOption := &metadata.ListProcessInstanceRelationOption{
					BusinessID:         bizID,
					ServiceInstanceIDs: []int64{serviceInstance.ID},
					Page: metadata.BasePage{
						Limit: common.BKNoLimit,
					},
				}
				relationResult, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relationOption)
				if err != nil {
					blog.ErrorJSON("createServiceInstances failed, ListProcessInstanceRelation failed, option: %s, err: %s, rid: %s", relationOption, err.Error(), rid)
					return nil, err
				}
				templateID2ProcessID := make(map[int64]int64)
				for _, relation := range relationResult.Info {
					templateID2ProcessID[relation.ProcessTemplateID] = relation.ProcessID
				}

				processes := make([]map[string]interface{}, 0)
				for _, item := range inst.Processes {
					templateID := item.ProcessTemplateID
					processID, exist := templateID2ProcessID[templateID]
					if !exist {
						continue
					}
					processData := item.ProcessData
					processData[common.BKProcessIDField] = processID
					processes = append(processes, processData)
				}
				input := metadata.UpdateRawProcessInstanceInput{
					BizID: bizID,
					Raw:   processes,
				}
				_, err = ps.updateProcessInstances(ctx, input)
				if err != nil {
					blog.ErrorJSON("CreateServiceInstances failed, updateProcessInstances failed, input: %s, err: %s, rid: %s", input, err.Error(), rid)
					return nil, err
				}
			}
		}

		if err = ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstance.ID); err != nil {
			blog.ErrorJSON("createServiceInstances failed, reconstruct service instance name failed, instanceID: %d, err: %s, rid:  %s", serviceInstance.ID, err.Error(), rid)
			return nil, err
		}

		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	// update host by host apply rule conflict resolvers
	attributeIDs := make([]int64, 0)
	for _, rule := range input.HostApplyConflictResolvers {
		attributeIDs = append(attributeIDs, rule.AttributeID)
	}
	attrCond := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID, common.BKPropertyIDField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: attributeIDs,
			},
		},
	}
	attrRes, e := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDHost, attrCond)
	if e != nil {
		blog.ErrorJSON("ReadModelAttr failed, err: %s, attrCond: %s, rid: %s", e.Error(), attrCond, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := attrRes.CCError(); err != nil {
		blog.ErrorJSON("ReadModelAttr failed, err: %s, attrCond: %s, rid: %s", err.Error(), attrCond, rid)
		return nil, err
	}
	attrMap := make(map[int64]string)
	for _, attr := range attrRes.Data.Info {
		attrMap[attr.ID] = attr.PropertyID
	}

	hostAttrMap := make(map[int64]map[string]interface{})
	for _, rule := range input.HostApplyConflictResolvers {
		if hostAttrMap[rule.HostID] == nil {
			hostAttrMap[rule.HostID] = make(map[string]interface{})
		}
		hostAttrMap[rule.HostID][attrMap[rule.AttributeID]] = rule.PropertyValue
	}
	for hostID, hostData := range hostAttrMap {
		updateOption := &metadata.UpdateOption{
			Data: hostData,
			Condition: map[string]interface{}{
				common.BKHostIDField: hostID,
			},
		}
		updateResult, err := ps.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDHost, updateOption)
		if err != nil {
			blog.ErrorJSON("RunHostApplyRule, update host failed, option: %s, err: %s, rid: %s", updateOption, err.Error(), rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if err := updateResult.CCError(); err != nil {
			blog.ErrorJSON("RunHostApplyRule, update host response failed, option: %s, response: %s, rid: %s", updateOption, updateResult, rid)
			return nil, err
		}
	}
	return serviceInstanceIDs, nil
}

// CreateServiceInstancesPreview generate a preview for service instance creation related host transfer action
func (ps *ProcServer) CreateServiceInstancesPreview(ctx *rest.Contexts) {
	input := metadata.CreateServiceInstancePreviewInput{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	previews, err := ps.createServiceInstancesPreview(ctx, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(previews)
}

func (ps *ProcServer) createServiceInstancesPreview(ctx *rest.Contexts, input metadata.CreateServiceInstancePreviewInput) ([]metadata.HostTransferPreview, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			blog.Errorf("createServiceInstances failed, parse bizID failed, metadata: %+v, err: %s, rid: %s", input.Metadata, err.Error(), rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommHTTPInputInvalid, common.MetadataField)
		}
	}
	// check hosts in business
	hostIDs := input.HostIDs
	if err := ps.CheckHostInBusiness(ctx, bizID, hostIDs); err != nil {
		blog.ErrorJSON("createServiceInstances failed, CheckHostInBusiness failed, bizID: %s, hostIDs: %s, err: %s, rid: %s", bizID, hostIDs, err.Error(), rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceHostNotBelongBusiness, hostIDs, bizID)
	}
	// check module in business
	module, ccErr := ps.getModule(ctx, input.ModuleID)
	if ccErr != nil {
		blog.Errorf("createServiceInstances failed, get module failed, moduleID: %d, err: %v, rid: %s", input.ModuleID, ccErr, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, ccErr.Error())
	}
	if bizID != module.BizID {
		blog.Errorf("createServiceInstances failed, module %d not belongs to biz %d, rid: %s", input.ModuleID, bizID, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, module.ModuleID, bizID)
	}

	// get module service template
	var serviceTemplate *metadata.ServiceTemplateDetail
	if module.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		serviceTemplateDetails, err := ps.CoreAPI.CoreService().Process().ListServiceTemplateDetail(ctx.Kit.Ctx, ctx.Kit.Header, bizID, module.ServiceTemplateID)
		if err != nil {
			blog.Errorf("ListServiceTemplateDetail failed, err: %s, bizID: %d, serviceTemplateID: %d, rid: %s", err.Error(), bizID, module.ServiceTemplateID, ctx.Kit.Rid)
			return nil, err
		}
		for _, serviceTemplateDetail := range serviceTemplateDetails.Info {
			serviceTemplate = &serviceTemplateDetail
		}
	}

	// get default modules, if host pre module is default module, remove host from it
	moduleCondition := metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField},
		Condition: map[string]interface{}{
			common.BKDefaultField: map[string]interface{}{
				common.BKDBNE: common.DefaultFlagDefaultValue,
			},
		},
	}
	defaultModules, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, &moduleCondition)
	if err != nil {
		blog.ErrorJSON("ReadInstance of %s failed, filter: %s, err: %s, rid: %s", common.BKTableNameBaseModule, moduleCondition, err.Error(), rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	isDefaultModuleMap := make(map[int64]bool)
	for _, module := range defaultModules.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("parse module from db failed, module: %s, err: %s, rid: %s", module, err.Error(), rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		isDefaultModuleMap[moduleID] = true
	}

	// get host module info
	hostModuleRelationRes, err := ps.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header, &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		HostIDArr:     hostIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField, common.BKHostIDField},
	})
	if err != nil {
		blog.Errorf("GetHostModuleRelation failed, err: %v, bizID: %d, hostIDs: %+v, rid: %s", err, input.BizID, hostIDs, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrHostModuleConfigFailed, err.Error())
	}
	preModuleIDs := make([]int64, 0)
	hostModuleMap := make(map[int64][]int64)
	hostRemoveModuleMap := make(map[int64]int64)
	for _, relation := range hostModuleRelationRes.Data.Info {
		moduleID := relation.ModuleID
		hostID := relation.HostID
		if hostModuleMap[relation.HostID] == nil {
			hostModuleMap[relation.HostID] = []int64{}
		}
		if isDefaultModuleMap[moduleID] {
			hostRemoveModuleMap[hostID] = moduleID
			continue
		}
		preModuleIDs = append(preModuleIDs, moduleID)
		hostModuleMap[relation.HostID] = append(hostModuleMap[relation.HostID], moduleID)
	}
	preModuleIDs = util.IntArrayUnique(preModuleIDs)

	// get modules that enabled host apply rule
	moduleCondition.Condition = map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: preModuleIDs,
		},
		common.HostApplyEnabledField: true,
	}
	enabledModules, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, &moduleCondition)
	if err != nil {
		blog.ErrorJSON("ReadInstance of %s failed, filter: %s, err: %s, rid: %s", common.BKTableNameBaseModule, moduleCondition, err.Error(), rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	enableModuleMap := make(map[int64]bool)
	for _, module := range enabledModules.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.ErrorJSON("parse module from db failed, module: %s, err: %s, rid: %s", module, err.Error(), rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		enableModuleMap[moduleID] = true
	}
	hostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range hostModuleMap {
		host2Module := metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: make([]int64, 0),
		}
		if module.HostApplyEnabled {
			host2Module.ModuleIDs = append(host2Module.ModuleIDs, module.ModuleID)
		}
		for _, moduleID := range moduleIDs {
			if enableModuleMap[moduleID] {
				host2Module.ModuleIDs = append(host2Module.ModuleIDs, moduleID)
			}
		}
		hostModules = append(hostModules, host2Module)
	}

	// generate host apply plans
	ruleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: append(preModuleIDs, input.ModuleID),
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	rules, ccErr := ps.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, ruleOption)
	if ccErr != nil {
		blog.ErrorJSON("ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, ruleOption, ccErr.Error(), rid)
		return nil, ccErr
	}

	planOption := metadata.HostApplyPlanOption{
		Rules:       rules.Info,
		HostModules: hostModules,
	}
	hostApplyPlanResult, ccErr := ps.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(ctx.Kit.Ctx, ctx.Kit.Header, bizID, planOption)
	if ccErr != nil {
		blog.ErrorJSON("TransferHostWithAutoClearServiceInstance failed, generateApplyPlan failed, core service GenerateApplyPlan failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, planOption, ccErr.Error(), rid)
		return nil, ccErr
	}
	hostApplyPlanMap := make(map[int64]metadata.OneHostApplyPlan)
	for _, item := range hostApplyPlanResult.Plans {
		hostApplyPlanMap[item.HostID] = item
	}

	// generate preview for each host
	previews := make([]metadata.HostTransferPreview, 0)
	for _, hostID := range hostIDs {
		preview := metadata.HostTransferPreview{
			HostID:              hostID,
			FinalModules:        append(hostModuleMap[hostID], input.ModuleID),
			ToRemoveFromModules: make([]metadata.RemoveFromModuleInfo, 0),
			ToAddToModules: []metadata.AddToModuleInfo{
				{
					ModuleID:        input.ModuleID,
					ServiceTemplate: serviceTemplate,
				},
			},
			HostApplyPlan: hostApplyPlanMap[hostID],
		}
		if hostRemoveModuleMap[hostID] != 0 {
			preview.ToRemoveFromModules = []metadata.RemoveFromModuleInfo{{ModuleID: hostRemoveModuleMap[hostID]}}
		}
		previews = append(previews, preview)
	}
	return previews, nil
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
		ModuleIDs:  []int64{input.ModuleID},
		Page:       input.Page,
		SearchKey:  input.SearchKey,
		Selectors:  input.Selectors,
		HostIDs:    input.HostIDs,
	}
	serviceInstanceResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in module: %d failed, err: %v", input.ModuleID, err)
		return
	}

	serviceInstanceIDs := make([]int64, 0)
	for _, instance := range serviceInstanceResult.Info {
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
		if _, ok := processCountMap[relation.ServiceInstanceID]; !ok {
			processCountMap[relation.ServiceInstanceID] = 0
		}
		processCountMap[relation.ServiceInstanceID] += 1
	}

	// insert `process_count` field
	serviceInstanceDetails := make([]map[string]interface{}, 0)
	for _, instance := range serviceInstanceResult.Info {
		item, err := mapstr.Struct2Map(instance)
		if err != nil {
			blog.ErrorJSON("SearchServiceInstancesInModuleWeb failed, Struct2Map failed, serviceInstance: %s, err: %s, rid: %s", instance, err.Error(), ctx.Kit.Rid)
			ccErr := ctx.Kit.CCError.CCError(common.CCErrCommParseDBFailed)
			ctx.RespAutoError(ccErr)
			return
		}
		item["process_count"] = 0
		if count, ok := processCountMap[instance.ID]; ok {
			item["process_count"] = count
		}
		serviceInstanceDetails = append(serviceInstanceDetails, item)
	}
	result := metadata.MultipleMap{
		Count: serviceInstanceResult.Count,
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
		ModuleIDs:  []int64{input.ModuleID},
		Page:       input.Page,
		SearchKey:  input.SearchKey,
		Selectors:  input.Selectors,
		HostIDs:    input.HostIDs,
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

	if len(input.ServiceInstanceIDs) == 0 {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "delete service instances, service_instance_ids empty", "service_instance_ids")
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		// when a service instance is deleted, the related data should be deleted at the same time
		for _, serviceInstanceID := range input.ServiceInstanceIDs {
			serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceID)
			if err != nil {
				blog.Errorf("delete service instance failed, service instance not found, serviceInstanceIDs: %d", serviceInstanceID)
				return ctx.Kit.CCError.CCError(common.CCErrProcGetProcessInstanceFailed)
			}
			if serviceInstance.BizID != bizID {
				blog.Errorf("delete service instance failed, biz id from input and service instance not equal, serviceInstanceIDs: %d", serviceInstanceID)
				return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.MetadataField)
			}

			// step1: delete the service instance relation.
			option := &metadata.ListProcessInstanceRelationOption{
				BusinessID:         bizID,
				ServiceInstanceIDs: []int64{serviceInstanceID},
			}
			relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, option)
			if err != nil {
				blog.Errorf("delete service instance: %d, but list service instance relation failed.", serviceInstanceID)
				return ctx.Kit.CCError.CCError(common.CCErrProcGetProcessInstanceRelationFailed)
			}

			if len(relations.Info) > 0 {
				deleteOption := metadata.DeleteProcessInstanceRelationOption{
					ServiceInstanceIDs: []int64{serviceInstanceID},
				}
				err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
				if err != nil {
					blog.Errorf("delete service instance: %d, but delete service instance relations failed.", serviceInstanceID)
					return ctx.Kit.CCError.CCError(common.CCErrProcDeleteServiceInstancesFailed)
				}

				// step2: delete process instance belongs to this service instance.
				processIDs := make([]int64, 0)
				for _, r := range relations.Info {
					processIDs = append(processIDs, r.ProcessID)
				}
				if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, processIDs); err != nil {
					blog.Errorf("delete service instance: %d, but delete process instance failed.", serviceInstanceID)
					return ctx.Kit.CCError.CCError(common.CCErrProcDeleteServiceInstancesFailed)
				}
			}

			// step3: delete service instance.
			deleteOption := &metadata.CoreDeleteServiceInstanceOption{
				BizID:              bizID,
				ServiceInstanceIDs: []int64{serviceInstanceID},
			}
			err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
			if err != nil {
				blog.Errorf("delete service instance: %d failed, err: %v", serviceInstanceID, err)
				return ctx.Kit.CCError.CCError(common.CCErrProcDeleteServiceInstancesFailed)
			}

			// step4: check and move host from module if no serviceInstance on it
			filter := &metadata.ListServiceInstanceOption{
				BusinessID: bizID,
				HostIDs:    []int64{serviceInstance.HostID},
				ModuleIDs:  []int64{serviceInstance.ModuleID},
			}
			result, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, filter)
			if err != nil {
				blog.Errorf("get host related service instances failed, bizID: %d, serviceInstanceID: %d, err: %v", bizID, serviceInstance.HostID, err)
				return ctx.Kit.CCError.CCError(common.CCErrProcGetServiceInstancesFailed)
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
				blog.Errorf("remove host from module failed, option: %+v, err: %v", removeHostFromModuleOption, err)
				return ctx.Kit.CCError.CCError(common.CCErrHostMoveResourcePoolFail)
			}
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
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
	diffOption := metadata.DiffModuleWithTemplateOption{}
	if err := ctx.DecodeInto(&diffOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(diffOption.ModuleIDs) == 0 {
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_module_ids")
		ctx.RespAutoError(err)
		return
	}

	result := make([]*metadata.ModuleDiffWithTemplateDetail, 0)
	for _, moduleID := range diffOption.ModuleIDs {
		option := metadata.DiffOneModuleWithTemplateOption{
			Metadata: diffOption.Metadata,
			BizID:    diffOption.BizID,
			ModuleID: moduleID,
		}
		oneModuleResult, err := ps.diffServiceInstanceWithTemplate(ctx, option)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		result = append(result, oneModuleResult)
	}

	ctx.RespEntity(result)
}

func (ps *ProcServer) diffServiceInstanceWithTemplate(ctx *rest.Contexts, diffOption metadata.DiffOneModuleWithTemplateOption) (*metadata.ModuleDiffWithTemplateDetail, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid

	// why we need validate metadata here?
	bizID := diffOption.BizID
	if bizID == 0 && diffOption.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*diffOption.Metadata)
		if err != nil {
			blog.ErrorJSON("diffServiceInstanceWithTemplate failed, parse biz id failed, metadata: %s, err: %v, rid: %s", diffOption.Metadata, err, rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommHTTPInputInvalid, common.MetadataField)
		}
	}
	diffOption.BizID = bizID

	if diffOption.ModuleID == 0 {
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, module id empty, option: %s, rid: %s", diffOption, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	module, err := ps.getModule(ctx, diffOption.ModuleID)
	if err != nil {
		blog.Errorf("diffServiceInstanceWithTemplate failed, getModule failed, moduleID: %d, err: %+v, rid: %s", diffOption.ModuleID, err, rid)
		return nil, err
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
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, ReadModelAttr failed, option: %s, err: %s, rid: %s", cond, e, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}

	// step2. get process templates
	listProcessTemplateOption := &metadata.ListProcessTemplatesOption{
		BusinessID:         module.BizID,
		ServiceTemplateIDs: []int64{module.ServiceTemplateID},
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listProcessTemplateOption)
	if err != nil {
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, ListProcessTemplates failed, option: %s, err: %s, rid: %s", listProcessTemplateOption, err, rid)
		return nil, err
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
		ModuleIDs:         []int64{diffOption.ModuleID},
	}
	serviceInstances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceOption)
	if err != nil {
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, ListServiceInstance failed, option: %s, err: %s, rid: %s", serviceOption, err, rid)
		return nil, err
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
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, ListProcessInstanceRelation failed, option: %s, err: %s, rid: %s", option, err.Error(), rid)
		return nil, err
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
					blog.Errorf("diffServiceInstanceWithTemplate failed, GetProcessInstanceWithID failed, processID: %d, err: %s, rid: %s", relation.ProcessID, err.Error(), rid)
					return nil, err
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
			if _, exist := processTemplateReferenced[templateID]; exist {
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
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, CalculateModuleAttributeDifference failed, module: %s, err: %s, rid: %s", module, err.Error(), rid)
		return nil, err
	}
	moduleDifference.ChangedAttributes = moduleChangedAttributes

	if len(moduleDifference.Added) > 0 ||
		len(moduleDifference.Changed) > 0 ||
		len(moduleDifference.Removed) > 0 ||
		len(moduleDifference.ChangedAttributes) > 0 {
		moduleDifference.HasDifference = true
	}

	moduleDifference.ModuleID = diffOption.ModuleID
	return &moduleDifference, nil
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
	if module.ServiceCategoryID == serviceTpl.ServiceCategoryID {
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
	syncOption := metadata.SyncServiceInstanceByTemplateOption{}
	if err := ctx.DecodeInto(&syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(syncOption.ModuleIDs) == 0 {
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_module_ids")
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		err := ps.syncServiceInstanceByTemplate(ctx, syncOption)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(make(map[string]interface{}))
}

func (ps *ProcServer) syncServiceInstanceByTemplate(ctx *rest.Contexts, syncOption metadata.SyncServiceInstanceByTemplateOption) errors.CCErrorCoder {
	rid := ctx.Kit.Rid

	bizID := syncOption.BizID
	if bizID == 0 && syncOption.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*syncOption.Metadata)
		if err != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, parse bizID from metadata failed, metadata: %s, err: %s, rid: %s", syncOption.Metadata, err.Error(), rid)
			return ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.MetadataField)
		}
	}

	modules, err := ps.getModules(ctx, syncOption.ModuleIDs)
	if err != nil {
		blog.Errorf("syncServiceInstanceByTemplate failed, getModule failed, moduleIDs: %+v, err: %s, rid: %s", syncOption.ModuleIDs, err.Error(), rid)
		return err
	}

	// step 0:
	// find service instances
	serviceInstanceOption := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		ModuleIDs:  syncOption.ModuleIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstanceResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, serviceInstanceOption)
	if err != nil {
		blog.ErrorJSON("syncServiceInstanceByTemplate failed, ListServiceInstance failed, option: %s, err: %s, rid: %s", serviceInstanceOption, err.Error(), rid)
		return err
	}
	if serviceInstanceResult.Count == 0 {
		blog.V(3).Infof("syncServiceInstanceByTemplate success, no service instance found, option: %+v, rid: %s", serviceInstanceOption, rid)
		return nil
	}
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	// step 1:
	// find all the process template according to the service template id
	serviceTemplateIDs := make([]int64, 0)
	serviceTemplateModuleMap := make(map[int64][]*metadata.ModuleInst)
	for _, module := range modules {
		serviceTemplateIDs = append(serviceTemplateIDs, module.ServiceTemplateID)
		serviceTemplateModuleMap[module.ServiceTemplateID] = append(serviceTemplateModuleMap[module.ServiceTemplateID], module)
	}
	processTemplateFilter := &metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: serviceTemplateIDs,
	}
	processTemplate, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, processTemplateFilter)
	if err != nil {
		blog.ErrorJSON("syncServiceInstanceByTemplate failed, ListProcessTemplates failed, option: %s, err: %s, rid: %s", processTemplateFilter, err.Error(), rid)
		return err
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
		blog.ErrorJSON("syncServiceInstanceByTemplate failed, ListProcessInstanceRelation failed, option: %s, err: %s, rid: %s", relationOption, err.Error(), rid)
		return err
	}
	procIDs := make([]int64, 0)
	for _, r := range relations.Info {
		procIDs = append(procIDs, r.ProcessID)
	}

	// step 3:
	// find all the process instance in process instance relation.
	processInstances, err := ps.Logic.ListProcessInstanceWithIDs(ctx.Kit, procIDs)
	if err != nil {
		blog.ErrorJSON("syncServiceInstanceByTemplate failed, ListProcessInstanceWithIDs failed, procIDs: %s, err: %s, rid: %s", procIDs, err.Error(), rid)
		return err
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
	// {ServiceInstanceID: isNameChanged}
	serviceInstanceNameChangedMap := make(map[int64]bool)
	// {ServiceInstanceID: ServiceTemplateID}
	serviceInstanceTemplateMap := make(map[int64]int64)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstance2ProcessMap[serviceInstance.ID] = make([]*metadata.Process, 0)
		serviceInstanceWithTemplateMap[serviceInstance.ID] = make(map[int64]bool)
		serviceInstance2HostMap[serviceInstance.ID] = serviceInstance.HostID
		serviceInstanceTemplateMap[serviceInstance.ID] = serviceInstance.ServiceTemplateID
	}
	processInstanceWithTemplateMap := make(map[int64]int64)
	for _, r := range relations.Info {
		p, exist := processInstanceMap[r.ProcessID]
		if !exist {
			// something is wrong, but can this process instance,
			// but we can find it in the process instance relation.
			blog.Warnf("force sync process instance according to process template: %d, but can not find the process instance: %d, rid: %s", r.ProcessTemplateID, r.ProcessID, rid)
			continue
		}
		serviceInstance2ProcessMap[r.ServiceInstanceID] = append(serviceInstance2ProcessMap[r.ServiceInstanceID], p)
		processInstanceWithTemplateMap[r.ProcessID] = r.ProcessTemplateID
		serviceInstanceWithTemplateMap[r.ServiceInstanceID][r.ProcessTemplateID] = true
	}

	// step 5:
	// compare the difference between process instance and process template from one service instance to another.
	removedProcessIDs := make([]int64, 0)
	for serviceInstanceID, processes := range serviceInstance2ProcessMap {
		for _, process := range processes {
			processTemplateID := processInstanceWithTemplateMap[process.ProcessID]
			template, exist := processTemplateMap[processTemplateID]
			if !exist || template.ServiceTemplateID != serviceInstanceTemplateMap[serviceInstanceID] {
				// this process template has already removed form the service template,
				// which means this process instance need to be removed from this service instance
				removedProcessIDs = append(removedProcessIDs, process.ProcessID)
				serviceInstanceNameChangedMap[serviceInstanceID] = true
				continue
			}

			// this process's bounded is still exist, need to check whether this process instance
			// need to be updated or not.
			proc, changed, isNamePortChanged := template.ExtractChangeInfo(process)
			if !changed {
				continue
			}
			if !serviceInstanceNameChangedMap[serviceInstanceID] {
				serviceInstanceNameChangedMap[serviceInstanceID] = isNamePortChanged
			}
			if err := ps.Logic.UpdateProcessInstance(ctx.Kit, process.ProcessID, proc); err != nil {
				blog.Errorf("syncServiceInstanceByTemplate failed, UpdateProcessInstance failed, processID:%d, err: %s, rid:%s", process.ProcessID, err.Error(), rid)
				return err
			}
		}
	}
	// remove processes whose template has been removed
	if len(removedProcessIDs) != 0 {
		if err := ps.Logic.DeleteProcessInstanceBatch(ctx.Kit, removedProcessIDs); err != nil {
			blog.Errorf("syncServiceInstanceByTemplate failed, DeleteProcessInstance failed, processID: %d, err: %s, rid: %s", removedProcessIDs, err.Error(), rid)
			return err
		}
		// remove process instance relation now.
		deleteOption := metadata.DeleteProcessInstanceRelationOption{}
		deleteOption.ProcessIDs = removedProcessIDs
		if err := ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption); err != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, DeleteProcessInstanceRelation failed, option: %s, err: %s, rid: %s", deleteOption, err.Error(), rid)
			return err
		}
	}

	// step 6:
	// check if a new process is added to the service template.
	// if true, then create a new process instance for every service instance with process template's default value.
	for processTemplateID, processTemplate := range processTemplateMap {
		for svcID, templates := range serviceInstanceWithTemplateMap {
			if processTemplate.ServiceTemplateID != serviceInstanceTemplateMap[svcID] {
				continue
			}
			if _, exist := templates[processTemplateID]; exist {
				continue
			}

			// we can not find this process template in all this service instance,
			// which means that a new process template need to be added to this service instance
			newProcess := processTemplate.NewProcess(bizID, ctx.Kit.SupplierAccount)
			processData, e := mapstruct.Struct2Map(newProcess)
			if e != nil {
				blog.ErrorJSON("SyncServiceInstanceByTemplate failed, Struct2Map failed, process: %s, err: %s, rid: %s", newProcess, e.Error(), rid)
				return ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
			}
			newProcessID, err := ps.Logic.CreateProcessInstance(ctx.Kit, processData)
			if err != nil {
				blog.ErrorJSON("syncServiceInstanceByTemplate failed, CreateProcessInstance failed, option: %s, err: %s, rid: %s", processData, err.Error(), rid)
				return err
			}

			relation := &metadata.ProcessInstanceRelation{
				BizID:             bizID,
				ProcessID:         newProcessID,
				ServiceInstanceID: svcID,
				ProcessTemplateID: processTemplateID,
				HostID:            serviceInstance2HostMap[svcID],
			}

			// create service instance relation, so that the process instance created upper can be related to this service instance.
			_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relation)
			if err != nil {
				blog.ErrorJSON("syncServiceInstanceByTemplate failed, CreateProcessInstanceRelation failed, relation: %s, err: %s, rid: %s", relation, err.Error(), rid)
				return err
			}
		}
	}

	// reconstruct service instance's name as it's dependence(first process's name + first process's port) changed
	for _, svcInstanceID := range serviceInstanceIDs {
		if !serviceInstanceNameChangedMap[svcInstanceID] {
			continue
		}
		if err := ps.CoreAPI.CoreService().Process().ReconstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, svcInstanceID); err != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, ReconstructServiceInstanceName failed, instanceID:%d, err:%s, rid:%s", svcInstanceID, err.Error(), rid)
			return err
		}
	}

	// get service templates
	serviceTemplates, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, &metadata.ListServiceTemplateOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: serviceTemplateIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	})
	if err != nil {
		blog.Errorf("syncServiceInstanceByTemplate failed, ListServiceTemplates failed, serviceTemplateIDs: %+v, err: %s, rid: %s", serviceTemplateIDs, err.Error(), rid)
		return err
	}

	// step 7:
	// update module service category and name field
	for _, serviceTemplate := range serviceTemplates.Info {
		updateModules := make([]int64, 0)
		for _, module := range serviceTemplateModuleMap[serviceTemplate.ID] {
			if module == nil {
				continue
			}
			if serviceTemplate.ServiceCategoryID != module.ServiceCategoryID || serviceTemplate.Name != module.ModuleName {
				updateModules = append(updateModules, module.ModuleID)
			}
		}
		if len(updateModules) == 0 {
			continue
		}
		moduleUpdateOption := &metadata.UpdateOption{
			Data: map[string]interface{}{
				common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
				common.BKModuleNameField:        serviceTemplate.Name,
			},
			Condition: map[string]interface{}{
				common.BKModuleIDField: map[string]interface{}{
					common.BKDBIN: updateModules,
				},
			},
		}
		resp, e := ps.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, moduleUpdateOption)
		if e != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, UpdateInstance failed, option: %s, err: %s, rid:%s", moduleUpdateOption, e.Error(), rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if ccErr := resp.CCError(); ccErr != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, UpdateInstance failed, option: %s, result: %s, rid: %s", moduleUpdateOption, resp, rid)
			return ccErr
		}
	}
	return nil
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
		HostIDs:    []int64{input.HostID},
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
	rid := ctx.Kit.Rid
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
		HostIDs:    []int64{input.HostID},
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
		blog.Errorf("ListServiceInstancesWithHostWeb failed, search mainline instance topo failed, bizID: %d, err: %+v, riz: %s", bizID, e, rid)
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

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		if err := ps.CoreAPI.CoreService().Label().AddLabel(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameServiceInstance, option); err != nil {
			blog.Errorf("ServiceInstanceAddLabels failed, option: %+v, err: %v", option, err)
			return ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
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

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		if err := ps.CoreAPI.CoreService().Label().RemoveLabel(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameServiceInstance, option); err != nil {
			blog.Errorf("ServiceInstanceRemoveLabels failed, option: %+v, err: %v", option, err)
			return ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
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
		listOption.ModuleIDs = []int64{*option.ModuleID}
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
			if _, exist := aggregationData[key]; !exist {
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

func (ps *ProcServer) DeleteServiceInstancePreview(ctx *rest.Contexts) {
	// step1. parse request parameters
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
	if len(input.ServiceInstanceIDs) == 0 {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "delete service instances, service_instance_ids empty", "service_instance_ids")
		return
	}

	// step2. get to be delete service instances and related hosts
	listOption := &metadata.ListServiceInstanceOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: input.ServiceInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	listToDeleteResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed, "generate delete preview failed, ListServiceInstance failed, listOption: %+v", listOption)
		return
	}
	hostIDs := make([]int64, 0)
	hostModules := make(map[int64][]int64)
	for _, instance := range listToDeleteResult.Info {
		hostIDs = append(hostIDs, instance.HostID)
		if _, exist := hostModules[instance.HostID]; !exist {
			hostModules[instance.HostID] = make([]int64, 0)
		}
		if !util.InArray(instance.ModuleID, hostModules[instance.HostID]) {
			hostModules[instance.HostID] = append(hostModules[instance.HostID], instance.ModuleID)
		}
	}

	// step3. get related hosts expired service instances(current minus deleted items)
	listOption = &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		HostIDs:    hostIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	listCurrentResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceFailed, "generate delete preview failed, get hosts related service instance failed, hostIDs: %+v", hostIDs)
		return
	}
	expiredHostModule := make(map[int64][]int64)
	for _, instance := range listCurrentResult.Info {
		// skip to be delete serviceInstance
		if util.InArray(instance.ID, input.ServiceInstanceIDs) {
			continue
		}

		if _, exist := expiredHostModule[instance.HostID]; !exist {
			expiredHostModule[instance.HostID] = make([]int64, 0)
		}
		expiredHostModule[instance.HostID] = append(expiredHostModule[instance.HostID], instance.ModuleID)
	}

	// get idle module
	idleModule, err := ps.getDefaultModule(ctx, bizID, common.DefaultResModuleFlag)
	if err != nil {
		blog.Errorf("generate delete preview failed, getDefaultModule failed, bizID: %d, err: %+v, rid: %s", bizID, err, ctx.Kit.Rid)
		ctx.RespWithError(err, common.CCErrGetModule, "generate delete preview failed, get idle module failed, bizID: %d", bizID)
		return
	}
	idleModuleID := idleModule.ModuleID

	preview := metadata.ServiceInstanceDeletePreview{
		ToMoveModuleHosts: make([]metadata.RemoveFromModuleHost, 0),
	}

	// check host remove from modules
	finalHostModules := make([]metadata.Host2Modules, 0)
	for hostID, moduleIDs := range hostModules {
		hostPreview := metadata.RemoveFromModuleHost{
			HostID:            hostID,
			RemoveFromModules: make([]int64, 0),
		}
		expiredModules, exist := expiredHostModule[hostID]
		if !exist {
			// host will be move to idle module
			hostPreview.MoveToIdle = true
			hostPreview.RemoveFromModules = moduleIDs
			hostPreview.FinalModules = []int64{idleModuleID}
			preview.ToMoveModuleHosts = append(preview.ToMoveModuleHosts, hostPreview)
			finalHostModules = append(finalHostModules, metadata.Host2Modules{
				HostID:    hostID,
				ModuleIDs: hostPreview.FinalModules,
			})
			continue
		}
		for _, moduleID := range moduleIDs {
			if !util.InArray(moduleID, expiredModules) {
				// host will be remove from module:expiredModules[hostID]
				hostPreview.RemoveFromModules = append(hostPreview.RemoveFromModules, moduleID)
			} else {
				hostPreview.FinalModules = append(hostPreview.FinalModules, moduleID)
			}
		}
		if len(hostPreview.RemoveFromModules) > 0 {
			preview.ToMoveModuleHosts = append(preview.ToMoveModuleHosts, hostPreview)
		}
		finalHostModules = append(finalHostModules, metadata.Host2Modules{
			HostID:    hostID,
			ModuleIDs: hostPreview.FinalModules,
		})
	}

	finalModules := make([]int64, 0)
	for _, item := range finalHostModules {
		finalModules = append(finalModules, item.ModuleIDs...)
	}
	listRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: finalModules,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	ruleResult, err := ps.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, bizID, listRuleOption)
	if err != nil {
		blog.Errorf("generate delete preview failed, ListHostApplyRule failed, option: %+v, err: %+v, rid: %s", listRuleOption, err, ctx.Kit.Rid)
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "generate delete preview failed, ListHostApplyRule failed, option: %+v", listRuleOption)
		return
	}
	hostApplyPlanOption := metadata.HostApplyPlanOption{
		HostModules: finalHostModules,
		Rules:       ruleResult.Info,
	}
	applyPlan, err := ps.CoreAPI.CoreService().HostApplyRule().GenerateApplyPlan(ctx.Kit.Ctx, ctx.Kit.Header, bizID, hostApplyPlanOption)
	if err != nil {
		blog.Errorf("generate delete preview failed, GenerateApplyPlan failed, option: %+v, err: %+v, rid: %s", hostApplyPlanOption, err, ctx.Kit.Rid)
		ctx.RespWithError(err, common.CCErrGetModule, "generate delete preview failed, GenerateApplyPlan failed, option: %+v", hostApplyPlanOption)
		return
	}
	preview.HostApplyPlan = applyPlan
	ctx.RespEntity(preview)
}
