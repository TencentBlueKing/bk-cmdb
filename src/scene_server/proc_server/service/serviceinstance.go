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
	"strconv"
	"sync"

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
// 通过模板创建时，进程信息则表现为更新
func (ps *ProcServer) CreateServiceInstances(ctx *rest.Contexts) {
	input := metadata.CreateServiceInstanceInput{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	var serviceInstanceIDs []int64
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
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

func (ps *ProcServer) createServiceInstances(ctx *rest.Contexts, input metadata.CreateServiceInstanceInput) ([]int64,
	errors.CCErrorCoder) {

	if len(input.Instances) == 0 {
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "instances")
	}

	rid := ctx.Kit.Rid
	bizID := input.BizID
	moduleID := input.ModuleID

	// check if hosts are in the business module, and check if module is in the business
	module, err := ps.getModule(ctx.Kit, moduleID)
	if err != nil {
		blog.Errorf("get module failed, moduleID: %d, err: %v, rid: %s", moduleID, err, rid)
		return nil, err
	}

	if bizID != module.BizID {
		blog.Errorf("module %d has biz id %d, not belongs to biz %d, rid: %s", moduleID, module.BizID, bizID, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceHasModuleNotBelongBusiness, moduleID, bizID)
	}

	if module.Default != 0 {
		blog.Errorf("can not create service instance for inner module %d, rid: %s", moduleID, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	// check if process exists, can not create service instance with no process
	if module.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		procTempFilter := []map[string]interface{}{{common.BKServiceTemplateIDField: module.ServiceTemplateID}}
		count, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKTableNameProcessTemplate, procTempFilter)
		if err != nil {
			blog.Errorf("count service template(%d) proc failed, err: %v, rid: %s", module.ServiceTemplateID, err, rid)
			return nil, err
		}

		if count[0] == 0 {
			blog.Errorf("service template(%d) has no process template, rid: %s", module.ServiceTemplateID, rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
	}

	hostIDs := make([]int64, len(input.Instances))
	for idx, instance := range input.Instances {
		hostIDs[idx] = instance.HostID

		if module.ServiceTemplateID == common.ServiceTemplateIDNotSet && len(instance.Processes) == 0 {
			blog.Errorf("create srv inst(%#v) in module(%d) with no process, rid: %s", instance, module.ModuleID, rid)
			return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "instances.processes")
		}
	}

	// check if hosts are in the business module
	hostIDs = util.IntArrayUnique(hostIDs)
	if err := ps.checkHostsInModule(ctx.Kit, bizID, moduleID, hostIDs); err != nil {
		blog.Errorf("check hosts(%+v) in biz %d module %d failed, err: %v, rid: %s", hostIDs, bizID, moduleID, err, rid)
		return nil, err
	}

	// create service instances
	serviceInstances := make([]*metadata.ServiceInstance, len(input.Instances))
	for idx, inst := range input.Instances {
		instance := &metadata.ServiceInstance{
			BizID:             bizID,
			Name:              inst.ServiceInstanceName,
			ServiceTemplateID: module.ServiceTemplateID,
			ModuleID:          moduleID,
			HostID:            inst.HostID,
		}
		serviceInstances[idx] = instance
	}

	serviceInstances, err = ps.CoreAPI.CoreService().Process().CreateServiceInstances(ctx.Kit.Ctx, ctx.Kit.Header,
		serviceInstances)
	if err != nil {
		blog.ErrorJSON("create service instances(%s) failed, err: %s, rid: %s", serviceInstances, err, rid)
		return nil, err
	}

	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstances {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	if err := ps.upsertProcesses(ctx, serviceInstanceIDs, bizID, module.ServiceTemplateID, input.Instances); err != nil {
		return nil, err
	}

	return serviceInstanceIDs, nil
}

func (ps *ProcServer) upsertProcesses(ctx *rest.Contexts, serviceInstanceIDs []int64, bizID int64,
	serviceTemplateID int64, instances []metadata.CreateServiceInstanceDetail) errors.CCErrorCoder {

	instanceIDsUpdate := make([]int64, 0)
	instanceProcessesUpdateMap := make(map[int64][]metadata.ProcessInstanceDetail)
	for idx, inst := range instances {
		if len(inst.Processes) == 0 {
			continue
		}

		svcInstID := serviceInstanceIDs[idx]
		if serviceTemplateID == 0 {
			// if this service have process instance to create, then create it now.
			createProcInput := &metadata.CreateRawProcessInstanceInput{
				BizID:             bizID,
				ServiceInstanceID: svcInstID,
				Processes:         inst.Processes,
			}
			if _, err := ps.createProcessInstances(ctx, createProcInput); err != nil {
				blog.ErrorJSON("create process failed, input: %s, err: %s, rid: %s", createProcInput, err, ctx.Kit.Rid)
				return err
			}

			// if no service instance name is set and have processes under it, update it
			if inst.ServiceInstanceName == "" {
				err := ps.updateServiceInstanceName(ctx, svcInstID, inst.HostID, inst.Processes[0].ProcessData)
				if err != nil {
					blog.ErrorJSON("update service instance name failed, id: %s, hostID: %s, process: %s, err: %s, "+
						"rid: %s", svcInstID, inst.HostID, inst.Processes[0].ProcessData, err, ctx.Kit.Rid)
					return err
				}
			}
		} else {
			instanceIDsUpdate = append(instanceIDsUpdate, svcInstID)
			instanceProcessesUpdateMap[svcInstID] = inst.Processes
		}
	}

	if len(instanceIDsUpdate) == 0 {
		return nil
	}

	// update processes which have process template
	relOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: instanceIDsUpdate,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relOpt)
	if err != nil {
		blog.ErrorJSON("list process relation failed, option: %s, err: %s, rid: %s", relOpt, err, ctx.Kit.Rid)
		return err
	}

	templateID2ProcessID := make(map[int64]map[int64]int64)
	for _, relation := range relRes.Info {
		if templateID2ProcessID[relation.ProcessTemplateID] == nil {
			templateID2ProcessID[relation.ProcessTemplateID] = make(map[int64]int64)
		}
		templateID2ProcessID[relation.ProcessTemplateID][relation.ServiceInstanceID] = relation.ProcessID
	}

	processesUpdate := make([]map[string]interface{}, 0)
	for instanceID, processes := range instanceProcessesUpdateMap {
		for _, proc := range processes {
			if instProcMap, exist := templateID2ProcessID[proc.ProcessTemplateID]; exist {
				if processID, exist := instProcMap[instanceID]; exist {
					processData := proc.ProcessData
					processData[common.BKProcessIDField] = processID
					processesUpdate = append(processesUpdate, processData)
				}
			}
		}
	}

	if len(processesUpdate) > 0 {
		input := metadata.UpdateRawProcessInstanceInput{
			BizID: bizID,
			Raw:   processesUpdate,
		}
		if _, err = ps.updateProcessInstances(ctx, input); err != nil {
			blog.ErrorJSON("update process instances failed, input: %s, err: %s, rid: %s", input, err, ctx.Kit.Rid)
			return err
		}
	}

	return nil
}

func (ps *ProcServer) updateServiceInstanceName(ctx *rest.Contexts, serviceInstanceID, hostID int64, processData map[string]interface{}) errors.CCErrorCoder {
	firstProcess := new(metadata.Process)
	if err := mapstr.DecodeFromMapStr(firstProcess, processData); err != nil {
		blog.ErrorJSON("updateServiceInstanceName failed, Decode2Struct failed, process: %s, err: %s, rid: %s", processData, err.Error(), ctx.Kit.Rid)
		return ctx.Kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}

	hostMap, err := ps.Logic.GetHostIPMapByID(ctx.Kit, []int64{hostID})
	if err != nil {
		blog.Errorf("updateServiceInstanceName failed, getHostIPMapByID failed, hostID: %d, err: %v, rid: %s", hostID, err, ctx.Kit.Rid)
		return err
	}
	host := hostMap[hostID]

	srvInstNameParams := &metadata.SrvInstNameParams{
		ServiceInstanceID: serviceInstanceID,
		Host:              host,
		Process:           firstProcess,
	}

	return ps.CoreAPI.CoreService().Process().ConstructServiceInstanceName(ctx.Kit.Ctx, ctx.Kit.Header, srvInstNameParams)
}

// SearchHostWithNoServiceInstance used for ui to get hosts that has no service instance and can create one
func (ps *ProcServer) SearchHostWithNoServiceInstance(ctx *rest.Contexts) {
	input := new(metadata.SearchHostWithNoSvcInstInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if input.BizID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAppIDField))
		return
	}

	if input.ModuleID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField))
		return
	}

	// get hosts that has service instances, exclude them when creating service instances
	svcInstOpt := &metadata.ListServiceInstanceOption{
		BusinessID: input.BizID,
		ModuleIDs:  []int64{input.ModuleID},
		Page:       metadata.BasePage{Limit: common.BKNoLimit},
		HostIDs:    input.HostIDs,
	}
	svcInstRes, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, svcInstOpt)
	if err != nil {
		blog.Errorf("list service instance failed, err: %v, input: %#v, rid: %s", err, input, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	hasSvcHostIDMap := make(map[int64]struct{})
	for _, instance := range svcInstRes.Info {
		hasSvcHostIDMap[instance.HostID] = struct{}{}
	}

	// if input hosts are specified, exclude the previous hosts. If all of them have service instances, return nothing
	inputHostIDs := make([]int64, 0)
	if len(input.HostIDs) > 0 {
		for _, hostID := range input.HostIDs {
			if _, exists := hasSvcHostIDMap[hostID]; !exists {
				inputHostIDs = append(inputHostIDs, hostID)
			}
		}

		if len(inputHostIDs) == 0 {
			ctx.RespEntity(metadata.SearchHostWithNoSvcInstOutput{HostIDs: inputHostIDs})
			return
		}
	}

	// get all hosts that are in the module and satisfies the input rules
	hostOpt := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{input.BizID},
		ModuleIDArr:      []int64{input.ModuleID},
		HostIDArr:        inputHostIDs,
	}
	hostRes, err := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, ctx.Kit.Header, hostOpt)
	if err != nil {
		blog.Errorf("get host ids failed, err: %v, option: %#v, rid: %s", err, hostOpt, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// exclude the hosts with service instances, since input hosts are filtered before, we do not need to filter here
	if len(input.HostIDs) > 0 {
		ctx.RespEntity(metadata.SearchHostWithNoSvcInstOutput{HostIDs: hostRes})
		return
	}

	hostIDs := make([]int64, 0)
	for _, hostID := range hostRes {
		if _, exists := hasSvcHostIDMap[hostID]; !exists {
			hostIDs = append(hostIDs, hostID)
		}
	}
	ctx.RespEntity(metadata.SearchHostWithNoSvcInstOutput{HostIDs: hostIDs})
}

func (ps *ProcServer) SearchServiceInstancesInModuleWeb(ctx *rest.Contexts) {
	input := new(metadata.GetServiceInstanceInModuleInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
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

func (ps *ProcServer) SearchServiceInstancesBySetTemplate(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("SearchServiceInstancesBySetTemplate failed, parse bk_biz_id error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}
	input := new(metadata.GetServiceInstanceBySetTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if input.SetTemplateID == 0 {
		blog.Errorf("SearchServiceInstancesBySetTemplate failed, lost input params SetTemplateID, rid: %s", ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsLostField, "set_template_id"))
		return
	}

	// query modules by set_template_id
	cond := mapstr.MapStr{
		common.BKAppIDField:         bizID,
		common.BKSetTemplateIDField: input.SetTemplateID,
	}
	qc := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: cond,
	}
	moduleInsts, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header, common.BKInnerObjIDModule, qc)
	if err != nil {
		blog.Errorf("SearchServiceInstancesBySetTemplate failed, http request failed, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed))
		return
	}
	if !moduleInsts.Result {
		blog.ErrorJSON("SearchServiceInstancesBySetTemplate failed, ReadInstance failed, filter: %s, response: %s, rid: %s", qc, moduleInsts, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(moduleInsts.Code, moduleInsts.ErrMsg))
		return
	}

	// get the list of module by moduleInsts
	modules := make([]int64, moduleInsts.Data.Count)
	for _, moduleInst := range moduleInsts.Data.Info {
		moduleID, err := util.GetInt64ByInterface(moduleInst[common.BKModuleIDField])
		if err != nil {
			blog.ErrorJSON("SearchServiceInstancesBySetTemplate failed, GetInt64ByInterface failed, moduleInst: %s, err: %#v, rid: %s", moduleInsts, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.New(moduleInsts.Code, moduleInsts.ErrMsg))
			return
		}
		modules = append(modules, moduleID)
	}

	// set return the list of service instances sorted by id
	input.Page.Sort = "id"
	// query serviceInstances
	option := &metadata.ListServiceInstanceOption{
		BusinessID: bizID,
		ModuleIDs:  modules,
		Page:       input.Page,
	}
	serviceInstanceResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		blog.ErrorJSON("SearchServiceInstancesBySetTemplate failed, ListServiceInstance failed, filter: %s, err: %s, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in set_template_id: %d failed, err: %v", input.SetTemplateID, err)
		return
	}

	ctx.RespEntity(serviceInstanceResult)
}

func (ps *ProcServer) SearchServiceInstancesInModule(ctx *rest.Contexts) {
	input := new(metadata.GetServiceInstanceInModuleInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	option := &metadata.ListServiceInstanceOption{
		BusinessID: input.BizID,
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
	input := new(metadata.ListServiceInstanceDetailOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstanceDetail(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "get service instance in module: %d failed, err: %v", input.ModuleID, err)
		return
	}

	ctx.RespEntity(instances)
}

// UpdateServiceInstances update instances in one biz
func (ps *ProcServer) UpdateServiceInstances(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceInstances failed, parse bk_biz_id error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	option := new(metadata.UpdateServiceInstanceOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := ps.CoreAPI.CoreService().Process().UpdateServiceInstances(ctx.Kit.Ctx, ctx.Kit.Header, bizID, option); err != nil {
			blog.Errorf("UpdateServiceInstances failed, err:%s, bizID:%d, option:%#v, rid:%s",
				err, bizID, *option, ctx.Kit.Rid)
			return err
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}

	ctx.RespEntity(nil)
}

func (ps *ProcServer) DeleteServiceInstance(ctx *rest.Contexts) {
	input := new(metadata.DeleteServiceInstanceOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// TODO confirm if we need to validate the limit of the ids
	if len(input.ServiceInstanceIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "service_instance_ids"))
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		return ps.deleteServiceInstance(ctx.Kit, input.BizID, input.ServiceInstanceIDs)
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (ps *ProcServer) deleteServiceInstance(kit *rest.Kit, bizID int64, svcInstIDs []int64) error {
	if len(svcInstIDs) == 0 {
		return nil
	}

	// check if all service instances are exists in the business
	cntOpt := []map[string]interface{}{{
		common.BKAppIDField: bizID,
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: svcInstIDs}},
	}
	svcInstCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameServiceInstance, cntOpt)
	if err != nil {
		blog.Errorf("get service instances(%+v) count failed, err: %v, rid: %s", svcInstIDs, err, kit.Rid)
		return err
	}

	if svcInstCounts[0] != int64(len(svcInstIDs)) {
		blog.ErrorJSON("service instance ids(%+v) not all exists in business, rid: %s", svcInstIDs, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "service_instance_ids")
	}

	// when a service instance is deleted, the related data should be deleted at the same time
	// step1: delete the service instance relation.
	relationOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: svcInstIDs,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}
	relationRes, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header,
		relationOpt)
	if err != nil {
		blog.Errorf("list service instance(%+v) relations failed, err: %v, rid: %s", svcInstIDs, err, kit.Rid)
		return err
	}

	if len(relationRes.Info) > 0 {
		delOpt := metadata.DeleteProcessInstanceRelationOption{
			ServiceInstanceIDs: svcInstIDs,
		}
		err = ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(kit.Ctx, kit.Header, delOpt)
		if err != nil {
			blog.Errorf("delete service instance(%+v) relation failed, err: %v, rid: %s", svcInstIDs, err, kit.Rid)
			return err
		}

		// step2: delete process instance belongs to this service instance.
		processIDs := make([]int64, 0)
		for _, r := range relationRes.Info {
			processIDs = append(processIDs, r.ProcessID)
		}
		if err := ps.Logic.DeleteProcessInstanceBatch(kit, processIDs); err != nil {
			blog.Errorf("delete process instances(%+v) failed, err: %v, rid: %s", processIDs, err, kit.Rid)
			return err
		}
	}

	// step3: delete service instance.
	deleteOption := &metadata.CoreDeleteServiceInstanceOption{
		BizID:              bizID,
		ServiceInstanceIDs: svcInstIDs,
	}
	err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(kit.Ctx, kit.Header, deleteOption)
	if err != nil {
		blog.Errorf("delete service instances: %+v failed, err: %v, rid: %s", svcInstIDs, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrProcDeleteServiceInstancesFailed)
	}
	return nil
}

// this function works to find differences between the service template and service instances in a module.
// compared to the service template's process template, a process instance in the service instance may
// contains several differences, like as follows:
// unchanged: the process instance's property values are same with the process template it belongs.
// changed: the process instance's property values are not same with the process template it belongs.
// add: a new process template is added, compared to the service instance belongs to this service template.
// deleted: a process is already deleted, compared to the service instance belongs to this service template.
func (ps *ProcServer) DiffServiceInstanceWithTemplate(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid
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

	var wg sync.WaitGroup
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	result := make([]*metadata.ModuleDiffWithTemplateDetail, 0)
	isFinish := false

	for _, moduleID := range diffOption.ModuleIDs {
		if isFinish {
			break
		}
		pipeline <- true
		wg.Add(1)

		go func(bizID, moduleID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			option := metadata.DiffOneModuleWithTemplateOption{
				BizID:    bizID,
				ModuleID: moduleID,
			}
			oneModuleResult, err := ps.diffServiceInstanceWithTemplate(ctx, option)
			if err != nil {
				blog.ErrorJSON("diffServiceInstanceWithTemplate failed, err: %s, option: %s, rid: %s", err, option, rid)
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			result = append(result, oneModuleResult)

			// judge whether need compare partial and finish in advance
			if diffOption.PartialCompare {
				if oneModuleResult.HasDifference {
					isFinish = true
				}
			}

		}(diffOption.BizID, moduleID)
	}

	wg.Wait()
	if firstErr != nil {
		ctx.RespAutoError(firstErr)
		return
	}

	ctx.RespEntity(result)
}

func (ps *ProcServer) diffServiceInstanceWithTemplate(ctx *rest.Contexts, diffOption metadata.DiffOneModuleWithTemplateOption) (*metadata.ModuleDiffWithTemplateDetail, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid

	if diffOption.ModuleID == 0 {
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, module id empty, option: %s, rid: %s", diffOption, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	module, err := ps.getModule(ctx.Kit, diffOption.ModuleID)
	if err != nil {
		blog.Errorf("diffServiceInstanceWithTemplate failed, getModule failed, moduleID: %d, err: %+v, rid: %s", diffOption.ModuleID, err, rid)
		return nil, err
	}

	if module.ServiceTemplateID == 0 {
		blog.Errorf("module %d has no service template, option: %s, rid: %s", diffOption.ModuleID, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	// step1. get process templates
	listProcessTemplateOption := &metadata.ListProcessTemplatesOption{
		BusinessID:         module.BizID,
		ServiceTemplateIDs: []int64{module.ServiceTemplateID},
		Page: metadata.BasePage{
			Sort: common.BKFieldID,
		},
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listProcessTemplateOption)
	if err != nil {
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, ListProcessTemplates failed, option: %s, err: %s, rid: %s", listProcessTemplateOption, err, rid)
		return nil, err
	}

	// step 2:
	// find process instance's relations, which allows us know the relationship between
	// process instance and it's template, service instance, etc.
	pTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	for idx, pTemplate := range processTemplates.Info {
		pTemplateMap[pTemplate.ID] = &processTemplates.Info[idx]
	}

	// step 3:
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

	// step4. compare module and service template TODO: remove this when updating template includes syncing module
	moduleDifference := &metadata.ModuleDiffWithTemplateDetail{
		Unchanged:     make([]metadata.ServiceInstanceDifference, 0),
		Changed:       make([]metadata.ServiceInstanceDifference, 0),
		Added:         make([]metadata.ServiceInstanceDifference, 0),
		Removed:       make([]metadata.ServiceInstanceDifference, 0),
		HasDifference: false,
	}

	moduleChangedAttributes, err := ps.CalculateModuleAttributeDifference(ctx.Kit.Ctx, ctx.Kit.Header, *module)
	if err != nil {
		blog.ErrorJSON("calculate module attribute difference failed, module: %s, err: %s, rid: %s", module, err, rid)
		return nil, err
	}
	moduleDifference.ChangedAttributes = moduleChangedAttributes

	// if there is no service instance and no process template, then there's no need to compare the process changes
	if len(serviceInstances.Info) == 0 && len(processTemplates.Info) == 0 {
		return moduleDifference, nil
	}

	// step 6:
	// construct map {hostID ==> host}
	hostIDOpt := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{diffOption.BizID},
		ModuleIDArr:      []int64{diffOption.ModuleID},
	}
	hostIDs, err := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, ctx.Kit.Header, hostIDOpt)
	if err != nil {
		blog.Errorf("get host ids failed, err: %v, option: %#v, rid: %s", err, hostIDOpt, ctx.Kit.Rid)
		return nil, err
	}

	hostMap, err := ps.Logic.GetHostIPMapByID(ctx.Kit, hostIDs)
	if err != nil {
		return nil, err
	}

	// if no service instance is found, need to create all the service instances
	if len(serviceInstances.Info) == 0 {
		srvInstNameSuffix := ""
		proc := processTemplates.Info[0].Property
		if proc != nil {
			if proc.ProcessName.Value != nil && len(*proc.ProcessName.Value) > 0 {
				srvInstNameSuffix += "_" + processTemplates.Info[0].ProcessName
			}
			for _, bindInfo := range proc.BindInfo.Value {
				if bindInfo.Std != nil && bindInfo.Std.Port.Value != nil {
					srvInstNameSuffix += "_" + *bindInfo.Std.Port.Value
					break
				}
			}
		}

		svrInstDiffs := make([]metadata.ServiceDifferenceDetails, 0)
		for _, host := range hostMap {
			svrInstDiffs = append(svrInstDiffs, metadata.ServiceDifferenceDetails{
				ServiceInstance: metadata.SrvInstBriefInfo{
					ID:        0,
					Name:      util.GetStrByInterface(host[common.BKHostInnerIPField]) + srvInstNameSuffix,
					SvcTempID: module.ServiceTemplateID,
				},
				Flag: metadata.ServiceAdded,
			})
		}

		for templateID, processTemplate := range pTemplateMap {
			moduleDifference.Added = append(moduleDifference.Added, metadata.ServiceInstanceDifference{
				ProcessTemplateID:    templateID,
				ProcessTemplateName:  processTemplate.ProcessName,
				ServiceInstanceCount: len(svrInstDiffs),
				ServiceInstances:     svrInstDiffs,
			})
		}
		return moduleDifference, nil
	}

	// step 5:
	// construct map {ServiceInstanceID ==> []ProcessInstanceRelation}
	serviceInstanceIDs := make([]int64, 0)
	hostWithSrvInstMap := make(map[int64]struct{})
	for _, serviceInstance := range serviceInstances.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
		hostWithSrvInstMap[serviceInstance.HostID] = struct{}{}
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

	// step 7: find all the process instance detail by ids
	procIDs := make([]int64, 0)
	for _, r := range relations.Info {
		procIDs = append(procIDs, r.ProcessID)
	}

	processDetails, err := ps.Logic.ListProcessInstanceWithIDs(ctx.Kit, procIDs)
	if err != nil {
		blog.ErrorJSON("diffServiceInstanceWithTemplate failed, ListProcessInstanceWithIDs err:%s, procIDs: %s, rid: %s", err, procIDs, rid)
		return nil, err
	}
	procID2Detail := make(map[int64]*metadata.Process)
	for idx, p := range processDetails {
		procID2Detail[p.ProcessID] = &processDetails[idx]
	}

	// step 8: find process object's attribute
	cond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKObjIDField: common.BKInnerObjIDProc,
		}),
	}
	attrResult, e := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, cond)
	if e != nil {
		blog.ErrorJSON("read model attr failed, option: %s, err: %s, rid: %s", cond, e, rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}

	// step 9: compare the process instance with it's process template one by one in a service instance.
	for idx, serviceInstance := range serviceInstances.Info {
		relations := serviceRelationMap[serviceInstance.ID]

		for _, relation := range relations {
			// record the used process template for checking whether a new process template has been added to service template.
			processTemplateReferenced[relation.ProcessTemplateID] += 1

			process, ok := procID2Detail[relation.ProcessID]
			if !ok {
				process = new(metadata.Process)
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

			changedAttributes, isChanged, diffErr := ps.Logic.DiffWithProcessTemplate(property.Property, process,
				hostMap[serviceInstance.HostID], attributeMap, true)
			if diffErr != nil {
				blog.Errorf("diff with process template failed, process ID: %d  err: %v, rid: %s",
					relation.ProcessID, err, rid)
				return nil, errors.New(common.CCErrCommParamsInvalid, diffErr.Error())
			}

			if !isChanged {
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

	// step 10: handle all the service instances that need to be added
	if len(processTemplates.Info) > 0 {
		srvInstNameSuffix := ""
		proc := processTemplates.Info[0].Property
		if proc != nil {
			if proc.ProcessName.Value != nil && len(*proc.ProcessName.Value) > 0 {
				srvInstNameSuffix += "_" + processTemplates.Info[0].ProcessName
			}
			for _, bindInfo := range proc.BindInfo.Value {
				if bindInfo.Std != nil && bindInfo.Std.Port.Value != nil {
					srvInstNameSuffix += "_" + *bindInfo.Std.Port.Value
					break
				}
			}
		}

		for _, hostID := range hostIDs {
			if _, exists := hostWithSrvInstMap[hostID]; exists {
				continue
			}

			srvInstName := util.GetStrByInterface(hostMap[hostID][common.BKHostInnerIPField]) + srvInstNameSuffix
			for templateID, processTemplate := range pTemplateMap {
				record := recorder{
					ProcessName: processTemplate.ProcessName,
					ServiceInstance: &metadata.ServiceInstance{
						ID:                0,
						Name:              srvInstName,
						ServiceTemplateID: module.ServiceTemplateID,
					},
				}
				added[templateID] = append(added[templateID], record)
			}
		}
	}

	// it's time to rearrange the data
	for _, records := range removed {
		if len(records) == 0 {
			continue
		}
		processTemplateName := records[0].ProcessName

		serviceInstances := make([]metadata.ServiceDifferenceDetails, 0)
		for _, record := range records {
			item := metadata.ServiceDifferenceDetails{
				ServiceInstance: metadata.SrvInstBriefInfo{
					ID:        record.ServiceInstance.ID,
					Name:      record.ServiceInstance.Name,
					SvcTempID: record.ServiceInstance.ServiceTemplateID,
				},
				Process: record.Process,
			}
			if len(processTemplates.Info) == 0 {
				item.Flag = metadata.ServiceRemoved
			} else {
				item.Flag = metadata.ServiceChanged
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
			serviceInstances = append(serviceInstances, metadata.ServiceDifferenceDetails{ServiceInstance: metadata.SrvInstBriefInfo{
				ID:        record.ServiceInstance.ID,
				Name:      record.ServiceInstance.Name,
				SvcTempID: record.ServiceInstance.ServiceTemplateID,
			}})
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
				ServiceInstance: metadata.SrvInstBriefInfo{
					ID:        record.ServiceInstance.ID,
					Name:      record.ServiceInstance.Name,
					SvcTempID: record.ServiceInstance.ServiceTemplateID,
				},
				ChangedAttributes: record.ChangedAttribute,
				Flag:              metadata.ServiceChanged,
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
			sInstance := metadata.ServiceDifferenceDetails{
				ServiceInstance: metadata.SrvInstBriefInfo{
					ID:        s.ServiceInstance.ID,
					Name:      s.ServiceInstance.Name,
					SvcTempID: s.ServiceInstance.ServiceTemplateID,
				},
			}
			if s.ServiceInstance.ID == 0 {
				sInstance.Flag = metadata.ServiceAdded
			} else {
				sInstance.Flag = metadata.ServiceChanged
			}
			sInstances = append(sInstances, sInstance)
		}

		moduleDifference.Added = append(moduleDifference.Added, metadata.ServiceInstanceDifference{
			ProcessTemplateID:    addedID,
			ProcessTemplateName:  pTemplateMap[addedID].ProcessName,
			ServiceInstanceCount: len(sInstances),
			ServiceInstances:     sInstances,
		})
	}

	if len(moduleDifference.Added) > 0 ||
		len(moduleDifference.Changed) > 0 ||
		len(moduleDifference.Removed) > 0 ||
		len(moduleDifference.ChangedAttributes) > 0 {
		moduleDifference.HasDifference = true
	}

	moduleDifference.ModuleID = diffOption.ModuleID
	return moduleDifference, nil
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
	if module.ServiceCategoryID == serviceTpl.ServiceCategoryID && module.ModuleName == serviceTpl.Name {
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
		field := common.BKServiceCategoryIDField
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
		field := common.BKModuleNameField
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
// 1. add/update/remove a new process
// 2. sync module name and service category with service template
func (ps *ProcServer) SyncServiceInstanceByTemplate(ctx *rest.Contexts) {
	syncOpt := metadata.SyncServiceInstanceByTemplateOption{}
	if err := ctx.DecodeInto(&syncOpt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := syncOpt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	syncOneModuleOpt := metadata.SyncOneModuleBySvcTempOption{
		BizID:             syncOpt.BizID,
		ServiceTemplateID: syncOpt.ServiceTemplateID,
	}
	tasks := make([]metadata.CreateTaskRequest, 0)
	for _, moduleID := range syncOpt.ModuleIDs {
		syncOneModuleOpt.ModuleID = moduleID
		tasks = append(tasks, metadata.CreateTaskRequest{
			TaskType: common.SyncModuleTaskFlag,
			InstID:   moduleID,
			Data:     []interface{}{syncOneModuleOpt},
		})
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		taskRes, err := ps.CoreAPI.TaskServer().Task().CreateBatch(ctx.Kit.Ctx, ctx.Kit.Header, tasks)
		if err != nil {
			blog.Errorf("create service template sync task(%#v) failed, err: %v, rid: %s", tasks, err, ctx.Kit.Rid)
			return err
		}
		blog.V(4).Infof("successfully created service template sync task: %#v, rid: %s", taskRes, ctx.Kit.Rid)
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DoSyncServiceInstanceTask do sync one module's service instance by service template task
func (ps *ProcServer) DoSyncServiceInstanceTask(ctx *rest.Contexts) {
	syncOption := metadata.SyncOneModuleBySvcTempOption{}
	if err := ctx.DecodeInto(&syncOption); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := ps.doSyncServiceInstanceTask(ctx.Kit, syncOption); err != nil {
			blog.Errorf("do sync service instance task(%#v) failed, err: %v, rid: %s", syncOption, err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (ps *ProcServer) doSyncServiceInstanceTask(kit *rest.Kit,
	syncOption metadata.SyncOneModuleBySvcTempOption) errors.CCErrorCoder {

	module, err := ps.getModule(kit, syncOption.ModuleID)
	if err != nil {
		blog.Errorf("get module failed, moduleID: %d, err: %v, rid: %s", syncOption.ModuleID, err, kit.Rid)
		return err
	}

	// step 1:
	// find service instances
	svcInstOpt := &metadata.ListServiceInstanceOption{
		BusinessID:        syncOption.BizID,
		ModuleIDs:         []int64{syncOption.ModuleID},
		ServiceTemplateID: syncOption.ServiceTemplateID,
		Page:              metadata.BasePage{Limit: common.BKNoLimit},
	}
	svcInstRes, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(kit.Ctx, kit.Header, svcInstOpt)
	if err != nil {
		blog.ErrorJSON("list service instance failed, option: %s, err: %s, rid: %s", svcInstOpt, err, kit.Rid)
		return err
	}

	// {ServiceInstanceID: []Process}
	serviceInstance2ProcessMap := make(map[int64][]*metadata.Process)
	// {ServiceInstanceID: {ProcessTemplateID: true}}
	serviceInstanceWithTemplateMap := make(map[int64]map[int64]struct{})
	// {ServiceInstanceID: HostID}
	serviceInstance2HostMap := make(map[int64]int64)
	hostWithSrvInstMap := make(map[int64]struct{})
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range svcInstRes.Info {
		serviceInstance2ProcessMap[serviceInstance.ID] = make([]*metadata.Process, 0)
		serviceInstanceWithTemplateMap[serviceInstance.ID] = make(map[int64]struct{})
		serviceInstance2HostMap[serviceInstance.ID] = serviceInstance.HostID
		hostWithSrvInstMap[serviceInstance.HostID] = struct{}{}
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	// step 2:
	// get all host ids in the module
	hostOpt := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{syncOption.BizID},
		ModuleIDArr:      []int64{syncOption.ModuleID},
	}
	hostIDs, err := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, hostOpt)
	if err != nil {
		blog.Errorf("get host ids failed, err: %v, option: %#v, rid: %s", err, hostOpt, kit.Rid)
		return err
	}

	// find all the process template under the service template
	procTempOpt := &metadata.ListProcessTemplatesOption{
		BusinessID:         syncOption.BizID,
		ServiceTemplateIDs: []int64{syncOption.ServiceTemplateID},
	}
	procTemps, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(kit.Ctx, kit.Header, procTempOpt)
	if err != nil {
		blog.ErrorJSON("list process templates failed, option: %s, err: %s, rid: %s", procTempOpt, err, kit.Rid)
		return err
	}

	processTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	for idx, t := range procTemps.Info {
		processTemplateMap[t.ID] = &procTemps.Info[idx]
	}

	// step 3: handle all the service instances that need to be added
	srvInstToAdd := make([]*metadata.ServiceInstance, 0)
	if len(procTemps.Info) > 0 {
		for _, hostID := range hostIDs {
			if _, exists := hostWithSrvInstMap[hostID]; exists {
				continue
			}
			instance := &metadata.ServiceInstance{
				BizID:             syncOption.BizID,
				ServiceTemplateID: module.ServiceTemplateID,
				ModuleID:          module.ModuleID,
				HostID:            hostID,
			}
			srvInstToAdd = append(srvInstToAdd, instance)
		}
	}

	_, err = ps.CoreAPI.CoreService().Process().CreateServiceInstances(kit.Ctx, kit.Header, srvInstToAdd)
	if err != nil {
		blog.Errorf("create service instances(%#v) failed, err: %v, rid: %s", srvInstToAdd, err, kit.Rid)
		return err
	}

	// step 4:
	// update module service category and name field
	serviceTemplate, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(kit.Ctx, kit.Header,
		syncOption.ServiceTemplateID)
	if err != nil {
		blog.Errorf("get service template(%d) failed, err: %v, rid: %s", syncOption.ServiceTemplateID, err, kit.Rid)
		return err
	}
	if serviceTemplate.ServiceCategoryID != module.ServiceCategoryID ||
		serviceTemplate.Name != module.ModuleName {

		moduleUpdateOption := &metadata.UpdateOption{
			Data: map[string]interface{}{
				common.BKServiceCategoryIDField: serviceTemplate.ServiceCategoryID,
				common.BKModuleNameField:        serviceTemplate.Name,
			},
			Condition: map[string]interface{}{
				common.BKModuleIDField: syncOption.ModuleID,
			},
		}
		resp, e := ps.CoreAPI.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
			moduleUpdateOption)
		if e != nil {
			blog.Errorf("update module failed, option: %#v, err: %v, rid: %s", moduleUpdateOption, e, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if ccErr := resp.CCError(); ccErr != nil {
			blog.Errorf("update module failed, option: %#v, err: %v, rid: %s", moduleUpdateOption, ccErr, kit.Rid)
			return ccErr
		}
	}

	if len(serviceInstanceIDs) == 0 {
		return nil
	}

	// step5:
	// find all the process instances relations for the usage of getting process instances.
	relationOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         syncOption.BizID,
		ServiceInstanceIDs: serviceInstanceIDs,
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, relationOpt)
	if err != nil {
		blog.ErrorJSON("list process relation failed, option: %s, err: %s, rid: %s", relationOpt, err, kit.Rid)
		return err
	}
	procIDs := make([]int64, 0)
	for _, r := range relations.Info {
		procIDs = append(procIDs, r.ProcessID)
	}

	// step 6:
	// find all the process instance in process instance relation.
	processInstances, err := ps.Logic.ListProcessInstanceWithIDs(kit, procIDs)
	if err != nil {
		blog.ErrorJSON("list process instance with IDs failed, procIDs: %s, err: %s, rid: %s", procIDs, err, kit.Rid)
		return err
	}
	processInstanceMap := make(map[int64]*metadata.Process)
	for idx, p := range processInstances {
		processInstanceMap[p.ProcessID] = &processInstances[idx]
	}

	// step 7:
	// rearrange the service instance with process instance.
	processInstanceWithTemplateMap := make(map[int64]int64)
	for _, r := range relations.Info {
		p, exist := processInstanceMap[r.ProcessID]
		if !exist {
			// something is wrong, but can this process instance,
			// but we can find it in the process instance relation.
			blog.Warnf("but can not find the process instance: %d, rid: %s", r.ProcessTemplateID, r.ProcessID, kit.Rid)
			continue
		}
		if _, exist := serviceInstanceWithTemplateMap[r.ServiceInstanceID]; !exist {
			// something is wrong, service instance is not exist, but we can find it in the process instance relation
			blog.Warnf("relation: %#v has a service instance that is not exist, rid: %s", r, kit.Rid)
			continue
		}
		serviceInstance2ProcessMap[r.ServiceInstanceID] = append(serviceInstance2ProcessMap[r.ServiceInstanceID], p)
		processInstanceWithTemplateMap[r.ProcessID] = r.ProcessTemplateID
		serviceInstanceWithTemplateMap[r.ServiceInstanceID][r.ProcessTemplateID] = struct{}{}
	}

	// step 8:
	// construct map {hostID ==> host}
	hostMap, err := ps.Logic.GetHostIPMapByID(kit, hostIDs)
	if err != nil {
		return err
	}

	// step 9:
	// compare the difference between process instance and process template from one service instance to another.
	removedProcessIDs := make([]int64, 0)
	removedSvrInstIDs := make([]int64, 0)
	var wg sync.WaitGroup
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	for serviceInstanceID, processes := range serviceInstance2ProcessMap {
		if len(procTemps.Info) == 0 {
			removedSvrInstIDs = append(removedSvrInstIDs, serviceInstanceID)
			for _, process := range processes {
				removedProcessIDs = append(removedProcessIDs, process.ProcessID)
			}
			continue
		}

		for _, process := range processes {
			processTemplateID := processInstanceWithTemplateMap[process.ProcessID]
			template, exist := processTemplateMap[processTemplateID]
			if !exist || template.ServiceTemplateID != syncOption.ServiceTemplateID {
				// this process template has already removed form the service template,
				// which means this process instance need to be removed from this service instance
				removedProcessIDs = append(removedProcessIDs, process.ProcessID)
				continue
			}
			pipeline <- true
			wg.Add(1)

			go func(process *metadata.Process, host map[string]interface{}) {
				defer func() {
					wg.Done()
					<-pipeline
				}()

				// this process's bounded is still exist, need to check whether this process instance
				// need to be updated or not.
				proc, changed, err := template.ExtractChangeInfo(process, host)
				if err != nil {
					blog.ErrorJSON("extract process(%s) change info failed, err: %s, rid: %s", process, err, kit.Rid)
					if firstErr == nil {
						firstErr = errors.New(common.CCErrCommParamsInvalid, err.Error())
					}
					return
				}

				if !changed {
					return
				}

				if err := ps.Logic.UpdateProcessInstance(kit, process.ProcessID, proc); err != nil {
					blog.ErrorJSON("UpdateProcessInstance failed, processID: %s, process: %s, err: %s, rid: %s",
						process.ProcessID, proc, err, kit.Rid)
					if firstErr == nil {
						firstErr = err
					}
					return
				}

			}(process, hostMap[serviceInstance2HostMap[serviceInstanceID]])
		}
	}

	wg.Wait()
	if firstErr != nil {
		return firstErr
	}

	// remove processes whose template has been removed
	if len(removedProcessIDs) != 0 {
		if err := ps.Logic.DeleteProcessInstanceBatch(kit, removedProcessIDs); err != nil {
			blog.Errorf("delete process failed, processIDs: %+v, err: %s, rid: %s", removedProcessIDs, err, kit.Rid)
			return err
		}
		// remove process instance relation now.
		deleteOption := metadata.DeleteProcessInstanceRelationOption{}
		deleteOption.ProcessIDs = removedProcessIDs
		err := ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(kit.Ctx, kit.Header, deleteOption)
		if err != nil {
			blog.ErrorJSON("delete process relation failed, option: %s, err: %s, rid: %s", deleteOption, err, kit.Rid)
			return err
		}
	}

	// delete service instances whose processes are all removed
	if len(removedSvrInstIDs) > 0 {
		deleteOption := &metadata.CoreDeleteServiceInstanceOption{
			BizID:              syncOption.BizID,
			ServiceInstanceIDs: removedSvrInstIDs,
		}
		err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(kit.Ctx, kit.Header, deleteOption)
		if err != nil {
			blog.Errorf("delete service instances: %+v failed, err: %v, rid: %s", removedSvrInstIDs, err, kit.Rid)
			return kit.CCError.CCError(common.CCErrProcDeleteServiceInstancesFailed)
		}
	}

	// step 10:
	// check if a new process is added to the service template.
	// if true, then create a new process instance for every service instance with process template's default value.
	processDatas := make([]map[string]interface{}, 0)
	procRelations := make([]*metadata.ProcessInstanceRelation, 0)
	for processTemplateID, processTemplate := range processTemplateMap {
		for svcID, templates := range serviceInstanceWithTemplateMap {
			if processTemplate.ServiceTemplateID != syncOption.ServiceTemplateID {
				continue
			}
			if _, exist := templates[processTemplateID]; exist {
				continue
			}

			// we can not find this process template in all this service instance,
			// which means that a new process template need to be added to this service instance
			newProcess, generateErr := processTemplate.NewProcess(syncOption.BizID, kit.SupplierAccount,
				hostMap[serviceInstance2HostMap[svcID]])
			if generateErr != nil {
				blog.ErrorJSON("generate process instance by template %s failed, err: %s, rid: %s", processTemplate,
					generateErr, kit.Rid)
				return errors.New(common.CCErrCommParamsInvalid, generateErr.Error())
			}
			processDatas = append(processDatas, newProcess.Map())
			procRelations = append(procRelations, &metadata.ProcessInstanceRelation{
				BizID:             syncOption.BizID,
				ServiceInstanceID: svcID,
				ProcessTemplateID: processTemplateID,
				HostID:            serviceInstance2HostMap[svcID],
			})
		}
	}

	if len(processDatas) > 0 {
		// create process instances in batch
		processIDs, err := ps.Logic.CreateProcessInstances(kit, processDatas)
		if err != nil {
			blog.ErrorJSON("create process failed, err: %s, processDatas: %s, rid: %s", err, processDatas, kit.Rid)
			return kit.CCError.CCError(common.CCErrSyncServiceInstanceByTemplateFailed)
		}

		if len(processIDs) != len(procRelations) {
			blog.Error("the count of processIDs is not equal to the count of procInstRelations, rid: %s", kit.Rid)
			return nil
		}

		// create process instance relations in batch
		for idx, processID := range processIDs {
			procRelations[idx].ProcessID = processID
		}
		_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelations(kit.Ctx, kit.Header, procRelations)
		if err != nil {
			blog.ErrorJSON("create process relations(%s) failed, err: %s, rid: %s", err, procRelations, kit.Rid)
			return err
		}
	}
	return nil
}

// FindServiceTemplateSyncStatus find service template sync status
func (ps *ProcServer) FindServiceTemplateSyncStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("parse biz id %s failed, err: %v, rid: %s", bizIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.FindServiceTemplateSyncStatusOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(option.ModuleIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "bk_module_ids"))
		return
	}

	if option.ServiceTemplateID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKServiceTemplateIDField))
		return
	}

	// get latest sync service template api task sync status by modules
	statusOpt := &metadata.ListLatestSyncStatusRequest{
		Condition: map[string]interface{}{
			common.BKInstIDField:   map[string]interface{}{common.BKDBIN: option.ModuleIDs},
			common.BKTaskTypeField: common.SyncModuleTaskFlag,
		},
		Fields: []string{common.BKInstIDField, common.CreateTimeField, common.LastTimeField, common.CreatorField,
			common.BKStatusField},
	}

	taskStatusRes, err := ps.CoreAPI.TaskServer().Task().ListLatestSyncStatus(ctx.Kit.Ctx, ctx.Kit.Header, statusOpt)
	if err != nil {
		blog.Errorf("list latest sync status failed, option: %#v, err: %v, rid: %s", statusOpt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// compare modules with their service templates to get their sync status
	moduleCond := map[string]interface{}{
		common.BKAppIDField:             bizID,
		common.BKServiceTemplateIDField: option.ServiceTemplateID,
		common.BKModuleIDField:          map[string]interface{}{common.BKDBIN: option.ModuleIDs},
	}
	_, statuses, err := ps.Logic.GetSvcTempSyncStatus(ctx.Kit, bizID, moduleCond, false)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	statusMap := make(map[int64]bool)
	for _, status := range statuses {
		statusMap[status.ModuleID] = status.NeedSync
	}

	statusExistsMap := make(map[int64]struct{})
	for index, status := range taskStatusRes {
		statusExistsMap[status.InstID] = struct{}{}
		// if current status and api task status does not match, use current status
		if statusMap[status.InstID] && status.Status.IsSuccessful() {
			taskStatusRes[index].Status = metadata.APITAskStatusNeedSync
		} else if !statusMap[status.InstID] && !status.Status.IsSuccessful() {
			taskStatusRes[index].Status = metadata.APITaskStatusSuccess
		}
	}

	// compensate for the modules that hasn't been synced before, or its latest sync task is already outdated
	compensateModuleIDs := make([]int64, 0)
	for _, moduleID := range option.ModuleIDs {
		if _, exists := statusExistsMap[moduleID]; !exists {
			compensateModuleIDs = append(compensateModuleIDs, moduleID)
		}
	}

	compensateStatuses, err := ps.compensateSvcTempSyncStatus(ctx.Kit, compensateModuleIDs, statusMap)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(append(taskStatusRes, compensateStatuses...))
}

// compensateSvcTempSyncStatus compensate sync status for the modules with no sync task
func (ps *ProcServer) compensateSvcTempSyncStatus(kit *rest.Kit, moduleIDs []int64, statusMap map[int64]bool) (
	[]metadata.APITaskSyncStatus, error) {

	moduleOpt := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField, common.CreatorField, common.CreateTimeField,
			common.LastTimeField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		Condition:      mapstr.MapStr{common.BKModuleIDField: mapstr.MapStr{common.BKDBIN: moduleIDs}},
		DisableCounter: true,
	}
	moduleRes := new(metadata.ResponseModuleInstance)
	if err := ps.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, moduleOpt, &moduleRes); err != nil {
		blog.Errorf("get modules failed, err: %v, opt: %#v, rid: %s", err, moduleOpt, kit.Rid)
		return nil, err
	}
	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get modules failed, err: %v, opt: %#v, rid: %s", err, moduleOpt, kit.Rid)
		return nil, err
	}

	statuses := make([]metadata.APITaskSyncStatus, 0)
	for _, module := range moduleRes.Data.Info {
		status := metadata.APITaskSyncStatus{
			InstID:     module.ModuleID,
			Creator:    module.Creator,
			CreateTime: module.CreateTime.Time,
			LastTime:   module.LastTime.Time,
		}

		if statusMap[status.InstID] {
			status.Status = metadata.APITAskStatusNeedSync
		} else {
			status.Status = metadata.APITaskStatusSuccess
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

func (ps *ProcServer) ListServiceInstancesWithHost(ctx *rest.Contexts) {
	input := new(metadata.ListServiceInstancesWithHostInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if input.HostID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service instances with host, but got empty host id. input: %+v", input)
		return
	}

	option := metadata.ListServiceInstanceOption{
		BusinessID: input.BizID,
		HostIDs:    []int64{input.HostID},
		SearchKey:  input.SearchKey,
		Page:       input.Page,
		Selectors:  input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list service instance failed, bizID: %d, hostID: %d", input.BizID, input.HostID, err)
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

	if input.HostID == 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "list service instances with host, but got empty host id. input: %+v", input)
		return
	}

	option := metadata.ListServiceInstanceOption{
		BusinessID: input.BizID,
		HostIDs:    []int64{input.HostID},
		SearchKey:  input.SearchKey,
		Page:       input.Page,
		Selectors:  input.Selectors,
	}
	instances, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "list service instance failed, bizID: %d, hostID: %d", input.BizID, input.HostID, err)
		return
	}

	topoRoot, e := ps.CoreAPI.CoreService().Mainline().SearchMainlineInstanceTopo(ctx.Kit.Ctx, ctx.Kit.Header, input.BizID, false)
	if e != nil {
		blog.Errorf("search mainline instance topo failed, bizID: %d, err: %v, rid: %s", input.BizID, e, rid)
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

// ServiceInstanceUpdateLabels Update service instance label operation.
func (ps *ProcServer) ServiceInstanceUpdateLabels(ctx *rest.Contexts) {
	option := new(selector.LabelUpdateOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		if err := ps.CoreAPI.CoreService().Label().UpdateLabel(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKTableNameServiceInstance, option); err != nil {
			blog.Errorf("serviceInstance update labels failed, option: %+v, err: %v,rid: %s", option, err,
				ctx.Kit.Rid)
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

func (ps *ProcServer) ServiceInstanceAddLabels(ctx *rest.Contexts) {
	option := selector.LabelAddOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
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

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
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

	if option.BizID == 0 {
		ctx.RespErrorCodeF(common.CCErrCommParamsIsInvalid, "list service instance label, but got invalid biz id: 0", "bk_biz_id")
		return
	}

	listOption := &metadata.ListServiceInstanceOption{
		BusinessID: option.BizID,
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
