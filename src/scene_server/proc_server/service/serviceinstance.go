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

func uniqueGeneralResult(origin []*metadata.ServiceTemplateGeneralDiff) *metadata.ServiceTemplateGeneralDiff {

	if len(origin) == 0 {
		return &metadata.ServiceTemplateGeneralDiff{}
	}

	addMap := make(map[int][]metadata.ProcessGeneralInfo)
	changedMap := make(map[int][]metadata.ProcessGeneralInfo)
	removedMap := make(map[string][]metadata.ProcessGeneralInfo)
	attrFlag := false

	for _, v := range origin {
		for _, add := range v.Added {
			if _, exist := addMap[add.Id]; !exist {
				addMap[add.Id] = append(addMap[add.Id], add)
			}
		}

		for _, changed := range v.Changed {
			if _, exist := changedMap[changed.Id]; !exist {
				changedMap[changed.Id] = append(changedMap[changed.Id], changed)
			}
		}

		for _, removed := range v.Removed {
			if _, exist := removedMap[removed.Name]; !exist {
				removedMap[removed.Name] = append(removedMap[removed.Name], removed)
			}
		}
		if !attrFlag && v.ChangedAttribute {
			attrFlag = true
		}
	}

	result := new(metadata.ServiceTemplateGeneralDiff)

	for _, add := range addMap {
		result.Added = append(result.Added, add...)
	}
	for _, changed := range changedMap {
		result.Changed = append(result.Changed, changed...)
	}
	for _, removed := range removedMap {
		result.Removed = append(result.Removed, removed...)
	}

	result.ChangedAttribute = attrFlag

	return result
}

func (ps *ProcServer) getServiceCategoryDiff(ctx *rest.Contexts, diffOption metadata.ProcessTemplateDiffOption) (
	*metadata.ServiceCategoryName, errors.CCErrorCoder) {

	rid := ctx.Kit.Rid
	if diffOption.ModuleIDs[0] == 0 {
		blog.ErrorJSON("module id empty, option: %s, rid: %s", diffOption, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	module, err := ps.getModule(ctx.Kit, diffOption.ModuleIDs[0])
	if err != nil {
		blog.Errorf("getModule failed, moduleID: %d, err: %+v, rid: %s", diffOption.ModuleIDs, err, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, "get none or multiple modules")
	}

	serviceCategory, err := ps.CoreAPI.CoreService().Process().GetServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header,
		module.ServiceCategoryID)
	if err != nil {
		blog.Errorf("serviceCategory failed, moduleID: %d,ServiceCategoryID: %d, err: %+v, rid: %s",
			diffOption.ModuleIDs, module.ServiceCategoryID, err, rid)
		return nil, err
	}

	res := new(metadata.ServiceCategoryName)
	if serviceCategory.ParentID != 0 {
		serviceParentCategory, err := ps.CoreAPI.CoreService().Process().GetServiceCategory(ctx.Kit.Ctx, ctx.Kit.Header,
			serviceCategory.ParentID)
		if err != nil {
			return nil, err
		}
		res.ParentName = serviceParentCategory.Name
	}

	res.Name = serviceCategory.Name
	return res, nil
}

// DiffServiceInstanceDetail 获取单个实例的详细信息
func (ps *ProcServer) DiffServiceInstanceDetail(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid
	option := new(metadata.ServiceInstanceDetailReq)
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	op := &metadata.DiffOption{
		BizID:             option.BizID,
		ModuleID:          option.ModuleID,
		ServiceTemplateId: option.ServiceTemplateId,
	}

	msg, bFlag := op.ServiceInstancesOptionValidate()
	if !bFlag {
		blog.Errorf("option req is invalid,option: %+v,err: %s,rid: %s", option, msg, rid)
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, msg)
		ctx.RespAutoError(err)
		return
	}

	result, err := ps.serviceInstanceDetailDiff(ctx, option)
	if err != nil {
		blog.ErrorJSON("get service instance detail failed, err: %s, option: %s, rid: %s", err, option, rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)

}

// ListDiffServiceInstances 计算指定进程模板涉及到的服务实例列表
func (ps *ProcServer) ListDiffServiceInstances(ctx *rest.Contexts) {

	rid := ctx.Kit.Rid
	option := new(metadata.ListServiceInstancesOption)
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	op := &metadata.DiffOption{
		BizID:             option.BizID,
		ModuleID:          option.ModuleID,
		ServiceTemplateId: option.ServiceTemplateId,
	}

	msg, bFlag := op.ServiceInstancesOptionValidate()
	if !bFlag {
		blog.Errorf("request option is invalid,option: %+v,err %s,rid: %s", option, msg, rid)
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, msg)
		ctx.RespAutoError(err)
		return
	}

	result, err := ps.ListServiceInstances(ctx, option)
	if err != nil {
		blog.Errorf("list service instances failed,option: %+v, err: %s,  rid: %s", option, err, rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// DiffServiceTemplateGeneralDiff List which process templates have changed.
func (ps *ProcServer) DiffServiceTemplateGeneral(ctx *rest.Contexts) {
	rid := ctx.Kit.Rid
	option := new(metadata.ServiceTemplateDiffOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	msg, bFlag := option.ServiceTemplateOptionValidate()
	if !bFlag {
		blog.Errorf("parameters is invalid,option: %+v,err is %s,rid: %s.", option, msg, rid)
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, msg)
		ctx.RespAutoError(err)
		return
	}
	var (
		wg       sync.WaitGroup
		firstErr errors.CCErrorCoder
	)
	pipeline := make(chan bool, 10)
	result := make([]*metadata.ServiceTemplateGeneralDiff, 0)

	processTemplates, err := ps.getProcessTemplate(ctx, option.BizID, option.ServiceTemplateId)
	if err != nil {
		blog.Errorf("get  processTemplates failed,option: %+v,err is %v,rid: %s.", option, err, rid)
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, msg)
		ctx.RespAutoError(err)
		return
	}

	for _, moduleID := range option.ModuleIDs {

		pipeline <- true
		wg.Add(1)

		go func(bizID, moduleID int64, processTemplates *metadata.MultipleProcessTemplate) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			op := metadata.DiffWithOneModuleOption{
				BizID:    bizID,
				ModuleID: moduleID,
			}
			oneModuleResult, err := ps.serviceTemplateGeneralDiff(ctx, op, processTemplates)
			if err != nil {
				blog.Errorf("calculate service template diff failed, err: %v, option: %+v, rid: %s", err, op, rid)
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			result = append(result, oneModuleResult)

		}(option.BizID, moduleID, processTemplates)
	}

	wg.Wait()
	if firstErr != nil {
		ctx.RespAutoError(firstErr)
		return
	}
	uniqueResult := uniqueGeneralResult(result)
	ctx.RespEntity(uniqueResult)
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
				Type: metadata.ServiceAdded,
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
				item.Type = metadata.ServiceRemoved
			} else {
				item.Type = metadata.ServiceChanged
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
				Type:              metadata.ServiceChanged,
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
				sInstance.Type = metadata.ServiceAdded
			} else {
				sInstance.Type = metadata.ServiceChanged
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

// getHostInfo 根据bizId和moduleid获取hostMap
func (ps *ProcServer) getHostInfo(ctx *rest.Contexts, bizId int64, moduleId int64) (
	map[int64]map[string]interface{}, errors.CCErrorCoder) {

	hostIDOpt := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizId},
		ModuleIDArr:      []int64{moduleId},
	}

	hostIDs, err := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, ctx.Kit.Header, hostIDOpt)
	if err != nil {
		return nil, err
	}

	hostMap, err := ps.Logic.GetHostIPMapByID(ctx.Kit, hostIDs)
	if err != nil {
		return nil, err
	}

	return hostMap, nil
}

func (ps *ProcServer) getProcDetail(ctx *rest.Contexts, relations *metadata.MultipleProcessInstanceRelation) (
	map[int64]*metadata.Process, errors.CCErrorCoder) {

	procIDs := make([]int64, 0)
	for _, r := range relations.Info {
		procIDs = append(procIDs, r.ProcessID)
	}

	processDetails, err := ps.Logic.ListProcessInstanceWithIDs(ctx.Kit, procIDs)
	if err != nil {
		blog.ErrorJSON("list process details failed err: %s, procIDs: %s, rid: %s", err, procIDs, ctx.Kit.Rid)
		return nil, err
	}

	procID2Detail := make(map[int64]*metadata.Process)
	for idx, p := range processDetails {
		procID2Detail[p.ProcessID] = &processDetails[idx]
	}
	return procID2Detail, nil
}

func (ps *ProcServer) getRelations(ctx *rest.Contexts, module *metadata.ModuleInst,
	sInstances []serviceListBrief) (*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder) {

	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range sInstances {
		//serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)

		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.id)
	}

	option := metadata.ListProcessInstanceRelationOption{
		BusinessID:         module.BizID,
		ServiceInstanceIDs: serviceInstanceIDs,
	}

	rs, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		blog.ErrorJSON("list process instance failed, option: %s, err: %s, rid: %s", option, err.Error(),
			ctx.Kit.Rid)
		return nil, err
	}

	return rs, nil
}

// initModuleDiffRes 初始化结果结构。 added和changed 用来做过滤
func initModuleDiffRes() (*metadata.ServiceTemplateGeneralDiff, map[int64]struct{}, map[int64]struct{}) {
	moduleDifference := &metadata.ServiceTemplateGeneralDiff{
		Changed: make([]metadata.ProcessGeneralInfo, 0),
		Added:   make([]metadata.ProcessGeneralInfo, 0),
		Removed: make([]metadata.ProcessGeneralInfo, 0),
	}
	added := make(map[int64]struct{})
	changed := make(map[int64]struct{})
	return moduleDifference, added, changed
}

func (ps *ProcServer) getProcIDDetailAndRelations(ctx *rest.Contexts, module *metadata.ModuleInst,
	serviceInstances []serviceListBrief) (map[int64][]metadata.ProcessInstanceRelation,
	map[int64]*metadata.Process, errors.CCErrorCoder) {
	relations, err := ps.getRelations(ctx, module, serviceInstances)
	if err != nil {
		blog.Errorf("get relations fail err: %v,rid: %s", err, ctx.Kit.Rid)
		return nil, nil, err
	}

	serviceRelationMap := make(map[int64][]metadata.ProcessInstanceRelation)
	for _, r := range relations.Info {
		serviceRelationMap[r.ServiceInstanceID] = append(serviceRelationMap[r.ServiceInstanceID], r)
	}

	procID2Detail, err := ps.getProcDetail(ctx, relations)
	if err != nil {
		blog.Errorf("get proc detail fail err: %v,rid: %s", err, ctx.Kit.Rid)
		return nil, nil, err
	}
	return serviceRelationMap, procID2Detail, nil
}

// calculateGeneralDiff 计算每个进程模板的分类，分为三类:1、新增。2、变更。3、删除
func (ps *ProcServer) calculateGeneralDiff(ctx *rest.Contexts, module *metadata.ModuleInst,
	hostMap map[int64]map[string]interface{}, pTemplateMap map[int64]*metadata.ProcessTemplate,
	serviceInstances []serviceListBrief) (*metadata.ServiceTemplateGeneralDiff, errors.CCErrorCoder) {

	serviceRelationMap, procID2Detail, err := ps.getProcIDDetailAndRelations(ctx, module, serviceInstances)
	if err != nil {
		return nil, err
	}

	processTemplateReferenced := make(map[int64]struct{})

	moduleDifference, added, changed := initModuleDiffRes()

	for _, serviceInst := range serviceInstances {

		relations := serviceRelationMap[serviceInst.id]

		for _, relation := range relations {
			processTemplateReferenced[relation.ProcessTemplateID] = struct{}{}
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
				moduleDifference.Removed = append(moduleDifference.Removed, metadata.ProcessGeneralInfo{
					Id:   int(relation.ProcessTemplateID),
					Name: processName,
				})
				continue
			}

			_, isChanged, diffErr := ps.Logic.DiffWithProcessTemplate(property.Property, process,
				hostMap[serviceInst.hostId], map[string]metadata.Attribute{}, false)
			if diffErr != nil {
				blog.Errorf("compare template failed, process id: %d  err: %v, rid: %s", relation.ProcessID, err,
					ctx.Kit.Rid)
				return nil, errors.New(common.CCErrCommParamsInvalid, diffErr.Error())
			}

			if !isChanged {
				continue
			}

			if _, ok := changed[relation.ProcessTemplateID]; !ok {
				moduleDifference.Changed = append(moduleDifference.Changed, metadata.ProcessGeneralInfo{
					Id:   int(relation.ProcessTemplateID),
					Name: property.ProcessName,
				})
				changed[relation.ProcessTemplateID] = struct{}{}
			}
		}

		// check whether a new process template has been added.
		for templateID, processTemplate := range pTemplateMap {
			if _, ok := processTemplateReferenced[templateID]; ok {
				continue
			}
			if _, ok := added[templateID]; !ok {
				moduleDifference.Added = append(moduleDifference.Added, metadata.ProcessGeneralInfo{
					Id:   int(templateID),
					Name: processTemplate.ProcessName,
				})
				added[templateID] = struct{}{}
			}
		}
	}

	//这里得单独处理一下第一次added processTemplate的场景，第一次added的时候模块下面的主机还没有实例化
	if len(hostMap) > 0 && len(serviceInstances) == 0 {
		for templateID, processTemplate := range pTemplateMap {
			moduleDifference.Added = append(moduleDifference.Added, metadata.ProcessGeneralInfo{
				Id:   int(templateID),
				Name: processTemplate.ProcessName,
			})
			added[templateID] = struct{}{}
		}
	}

	return moduleDifference, nil
}

func (ps *ProcServer) getProcessTemplate(ctx *rest.Contexts, bizId int64, serviceTemplateID int64) (
	*metadata.MultipleProcessTemplate, errors.CCErrorCoder) {

	listProcessTemplateOption := &metadata.ListProcessTemplatesOption{
		Page: metadata.BasePage{
			Sort: common.BKFieldID,
		},
	}
	if bizId != 0 {
		listProcessTemplateOption.BusinessID = bizId
	}

	if serviceTemplateID != 0 {
		listProcessTemplateOption.ServiceTemplateIDs = []int64{serviceTemplateID}
	}

	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header,
		listProcessTemplateOption)
	if err != nil {
		blog.ErrorJSON("istProcessTemplates failed, option: %s, err: %s, rid: %s", listProcessTemplateOption,
			err, ctx.Kit.Rid)
		return nil, err
	}

	return processTemplates, nil
}

// hostAndServiceInstsOption 查询实例的指定字段
type hostAndServiceInstsOption struct {
	ProcessTemplateId int64
	BizID             int64
	ModuleID          int64
}

// hostIDs: 模块下的 hostid 列表
// relations: 进程id、服务实例id与进程模板之间的关系表
// serviceRelationMap:服务实例Id与relations 的Map
// hostMap:以 hostid 与host信息的Map 其中host信息只关心ip相关信息
// processTemplates: 进程模板信息
// pTemplateMap: 进程模板id为key的processTemplates Map
// hostWithSrvInstMap: 由于模块下的host并不一定全部实例化 serviceinstance中的hostid map.

// getHostAndServiceInstances 获取后续计算实例列表的的基本信息
func (ps *ProcServer) getHostAndServiceInstances(ctx *rest.Contexts, option *hostAndServiceInstsOption,
	module *metadata.ModuleInst, serviceInstances []serviceListBrief) (
	map[int64][]metadata.ProcessInstanceRelation, *metadata.MultipleProcessInstanceRelation,
	map[int64]map[string]interface{}, []int64, *metadata.MultipleProcessTemplate, map[int64]*metadata.ProcessTemplate,
	map[int64]struct{}, errors.CCErrorCoder) {

	cond := &metadata.ListProcessTemplatesOption{
		BusinessID:         module.BizID,
		ServiceTemplateIDs: []int64{module.ServiceTemplateID},
		Page: metadata.BasePage{
			Sort: common.BKFieldID,
		},
		ProcessTemplateIDs: []int64{option.ProcessTemplateId},
	}

	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, cond)
	if err != nil {
		blog.Errorf("list process templates failed, option: %s, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		return nil, nil, nil, []int64{}, nil, nil, nil, err
	}

	pTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	for idx, pTemplate := range processTemplates.Info {
		pTemplateMap[pTemplate.ID] = &processTemplates.Info[idx]
	}

	// construct map {hostID ==> host}
	var hostIdArray []int64
	for _, v := range serviceInstances {
		hostIdArray = append(hostIdArray, v.hostId)

	}

	hostIDOpt := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{option.BizID},
		ModuleIDArr:      []int64{option.ModuleID},
		HostIDArr:        hostIdArray,
	}

	hostIDs, err := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(ctx.Kit.Ctx, ctx.Kit.Header, hostIDOpt)
	if err != nil {
		blog.Errorf("get host ids failed, err: %v, option: %#v, rid: %s", err, hostIDOpt, ctx.Kit.Rid)
		return nil, nil, nil, []int64{}, nil, nil, nil, err
	}

	hostMap, err := ps.Logic.GetHostIPMapByID(ctx.Kit, hostIDs)
	if err != nil {
		blog.Errorf("get host info by id failed, option: %#v, err: %v, rid: %s", hostIDOpt, err, ctx.Kit.Rid)
		return nil, nil, nil, []int64{}, nil, nil, nil, err
	}

	// construct map {ServiceInstanceID ==> []ProcessInstanceRelation}
	serviceInstanceIDs := make([]int64, 0)
	hostWithSrvInstMap := make(map[int64]struct{})

	for _, serviceInstance := range serviceInstances {

		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.id)
		hostWithSrvInstMap[serviceInstance.hostId] = struct{}{}
	}
	serviceRelationMap := make(map[int64][]metadata.ProcessInstanceRelation)
	relations := new(metadata.MultipleProcessInstanceRelation)
	if len(serviceInstanceIDs) > 0 {
		op := metadata.ListProcessInstanceRelationOption{
			BusinessID:         module.BizID,
			ServiceInstanceIDs: serviceInstanceIDs,
			ProcessTemplateID:  option.ProcessTemplateId,
		}

		relations, err = ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, &op)
		if err != nil {
			blog.Errorf("list process instance relation failed, option: %s, err: %v, rid: %s", option,
				err, ctx.Kit.Rid)
			return nil, nil, nil, []int64{}, nil, nil, nil, err
		}

		for _, r := range relations.Info {
			serviceRelationMap[r.ServiceInstanceID] = append(serviceRelationMap[r.ServiceInstanceID], r)
		}
	}

	return serviceRelationMap, relations, hostMap, hostIDs, processTemplates, pTemplateMap, hostWithSrvInstMap, nil
}

// getProcAndAttributeMap 通过进程id获取进程的详细信息，注意此时的进程信息是同步之前信息
func (ps *ProcServer) getProcAndAttributeMap(ctx *rest.Contexts, info []metadata.ProcessInstanceRelation) (
	map[int64]*metadata.Process, map[string]metadata.Attribute, errors.CCErrorCoder) {

	procIDs := make([]int64, 0)
	for _, r := range info {
		procIDs = append(procIDs, r.ProcessID)
	}

	// find all the process instance detail by ids
	processDetails, err := ps.Logic.ListProcessInstanceWithIDs(ctx.Kit, procIDs)
	if err != nil {
		blog.ErrorJSON("list process instance with ids fail, err:%s, procIDs: %s, rid: %s", err, procIDs,
			ctx.Kit.Rid)
		return nil, nil, err
	}

	procID2Detail := make(map[int64]*metadata.Process)
	for idx, p := range processDetails {
		procID2Detail[p.ProcessID] = &processDetails[idx]
	}

	// find process object's attribute
	cond := &metadata.QueryCondition{
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKObjIDField: common.BKInnerObjIDProc,
		}),
	}
	attrResult, e := ps.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, cond)
	if e != nil {
		blog.ErrorJSON("read model attr failed, option: %s, err: %s, rid: %s", cond, e, ctx.Kit.Rid)
		return nil, nil, ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	attributeMap := make(map[string]metadata.Attribute)
	for _, attr := range attrResult.Data.Info {
		attributeMap[attr.PropertyID] = attr
	}
	return procID2Detail, attributeMap, nil
}

func (ps *ProcServer) getServiceInstances(ctx *rest.Contexts, module *metadata.ModuleInst, serviceId []int64,
	field []string) (*metadata.MultipleServiceInstance, errors.CCErrorCoder) {

	var count int
	serviceInsts := new(metadata.MultipleServiceInstance)

	for {
		option := &metadata.ListServiceInstanceOption{
			BusinessID:        module.BizID,
			ServiceTemplateID: module.ServiceTemplateID,
			Fields:            field,

			ModuleIDs: []int64{module.ModuleID},
			Page: metadata.BasePage{
				Limit: common.BKMaxInstanceLimit,
				Start: count,
				Sort:  "id",
			},
			ServiceInstanceIDs: serviceId,
		}

		sInsts, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			blog.Errorf(" list service instances failed, option: %s, err: %v, rid: %s", option, err, ctx.Kit.Rid)
			return nil, err
		}
		tempLen := len(sInsts.Info)
		serviceInsts.Count += uint64(tempLen)
		count += len(sInsts.Info)
		serviceInsts.Info = append(serviceInsts.Info, sInsts.Info...)

		if len(sInsts.Info) < common.BKMaxInstanceLimit {
			break
		}
	}

	return serviceInsts, nil
}

// getServiceListBrief 格式化获取到的字段
func getServiceListBrief(briefInfo []metadata.ServiceInstance) ([]serviceListBrief, error) {
	insts := make([]serviceListBrief, 0)
	for _, k := range briefInfo {
		hostId, err := util.GetInt64ByInterface(k.HostID)
		if err != nil {
			return nil, err
		}
		id, err := util.GetInt64ByInterface(k.ID)
		if err != nil {
			return nil, err
		}
		name := util.GetStrByInterface(k.Name)

		insts = append(insts, serviceListBrief{
			hostId: hostId,
			id:     id,
			name:   name,
		})
	}
	return insts, nil
}

func (ps *ProcServer) getBriefinsts(ctx *rest.Contexts, module *metadata.ModuleInst, serviceInstanceId []int64,
	fields []string) ([]serviceListBrief, errors.CCErrorCoder) {

	serviceInts, err := ps.getServiceInstances(ctx, module, serviceInstanceId, fields)
	if err != nil {
		return nil, errors.New(common.CCErrCommDBSelectFailed, err.Error())
	}

	insts, e := getServiceListBrief(serviceInts.Info)
	if e != nil {
		return nil, errors.New(common.CCErrCommDBSelectFailed, e.Error())
	}
	return insts, nil
}

// diffServiceCategoryResult 获取服务分类的差异结果
func (ps *ProcServer) diffServiceCategoryResult(ctx *rest.Contexts, module *metadata.ModuleInst) (
	*metadata.ServiceInstanceDetailResult, errors.CCErrorCoder) {

	diffDetails := new(metadata.ServiceInstanceDetailResult)

	moduleAttr, err := ps.CalculateModuleAttributeDifference(ctx.Kit.Ctx, ctx.Kit.Header, *module)
	if err != nil {
		blog.Errorf("calc module attr diff failed, module: %s, err: %v, rid: %s", module, err, ctx.Kit.Rid)
		return nil, err
	}

	diffDetails.ModuleAttribute = moduleAttr
	diffDetails.Type = metadata.ServiceOthers
	return diffDetails, nil
}

func (ps *ProcServer) serviceInstanceDetailDiff(ctx *rest.Contexts, option *metadata.ServiceInstanceDetailReq) (
	*metadata.ServiceInstanceDetailResult, errors.CCErrorCoder) {

	rid := ctx.Kit.Rid

	module, err := ps.getModuleInfo(ctx, option.ModuleID)
	if err != nil {
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	diffDetails := new(metadata.ServiceInstanceDetailResult)

	// 此处直接返回的是服务分类内容
	if option.ServiceCategory {
		result, err := ps.diffServiceCategoryResult(ctx, module)
		if err != nil {
			return nil, err
		}
		return result, nil
	}

	op := &hostAndServiceInstsOption{
		BizID:             option.BizID,
		ProcessTemplateId: option.ProcessTemplateId,
		ModuleID:          option.ModuleID,
	}

	fileds := []string{common.BKFieldID, common.BKHostIDField}
	insts, err := ps.getBriefinsts(ctx, module, []int64{option.ServiceInstanceId}, fileds)
	if err != nil {
		return nil, err
	}

	_, relations, hostMap, _, _, pTemplateMap, _, err := ps.getHostAndServiceInstances(ctx, op, module, insts)
	if err != nil {
		return nil, err
	}

	procID2Detail, attributeMap, err := ps.getProcAndAttributeMap(ctx, relations.Info)
	if err != nil {
		return nil, err
	}

	// record the used process template for checking whether a new process template has been added to service template.
	processTemplateReferenced := make(map[int64]struct{})

	for _, relation := range relations.Info {
		processTemplateReferenced[relation.ProcessTemplateID] = struct{}{}

		process, ok := procID2Detail[relation.ProcessID]
		if !ok {
			process = new(metadata.Process)
		}
		processName := ""
		if process.ProcessName != nil {
			processName = *process.ProcessName
		}
		property, exist := pTemplateMap[relation.ProcessTemplateID]

		//  当有多个删除模板的场景下，每个被删除的模板在关系表中的id都是0，所以只能通过name进行区分.
		if !exist && option.ProcessTemplateId == 0 && processName == option.ProcessTemplateName {
			diffDetails.Type = metadata.ServiceRemoved
			diffDetails.Process = process
			continue
		}
		//  无论是删除场景下 ProcessTemplateId 为0 还是没有找到请求的ProcessTemplateId 直接跳过就好
		if !exist {
			continue
		}
		id := insts[0].hostId

		changedAttributes, isChanged, e := ps.Logic.DiffWithProcessTemplate(property.Property, process,
			hostMap[id], attributeMap, true)
		if e != nil {
			blog.Errorf("diff process template failed, id: %d, err: %v, rid: %s", relation.ProcessID, e, rid)
			return nil, errors.New(common.CCErrCommParamsInvalid, e.Error())
		}

		if !isChanged {
			continue
		}

		diffDetails.ChangedAttributes = changedAttributes
		diffDetails.Type = metadata.ServiceChanged
	}

	if _, exist := processTemplateReferenced[option.ProcessTemplateId]; !exist {
		diffDetails.Type = metadata.ServiceAdded
	}

	return diffDetails, nil
}

// 服务分类改变场景下，模块下面的每个实例必然会随之变动,所以只要取前500个即可.
func listServiceCategoryInstances(count uint64, serviceInts []serviceListBrief) *metadata.ListServiceInstancesResult {

	result := new(metadata.ListServiceInstancesResult)

	if count > metadata.ServiceInstancesMaxNum {
		result.TotalCount = metadata.ServiceInstancestotalCount
	} else {
		result.TotalCount = strconv.FormatUint(count, 10)
	}

	for index, v := range serviceInts {
		if index >= metadata.ServiceInstancesMaxNum {
			break
		}
		result.ServiceInstances = append(result.ServiceInstances, metadata.ServiceInstancesInfo{
			Id:   v.id,
			Name: v.name,
		})
	}
	return result
}

func (ps *ProcServer) getModuleInfo(ctx *rest.Contexts, moduleId int64) (*metadata.ModuleInst, errors.CCErrorCoder) {

	module, err := ps.getModule(ctx.Kit, moduleId)
	if err != nil {
		blog.Errorf(" get module failed, option: %d, err: %v, rid: %s", moduleId, err, ctx.Kit.Rid)
		return nil, err
	}

	if module.ServiceTemplateID == 0 {
		blog.Errorf("module %d has no service template, option: %s, rid: %s", moduleId, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	return module, nil
}

// ListServiceInstances 获取最多500个有变化的实例，大于500时只显示500+

func (ps *ProcServer) ListServiceInstances(ctx *rest.Contexts,
	opt *metadata.ListServiceInstancesOption) (*metadata.ListServiceInstancesResult, errors.CCErrorCoder) {

	mod, err := ps.getModuleInfo(ctx, opt.ModuleID)
	if err != nil {
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	field := []string{common.BKFieldID, common.BKHostIDField, common.BKFieldName}
	serviceInstances, err := ps.getServiceInstances(ctx, mod, nil, field)

	insts, e := getServiceListBrief(serviceInstances.Info)
	if e != nil {
		return nil, errors.New(common.CCErrCommDBSelectFailed, e.Error())
	}

	result := new(metadata.ListServiceInstancesResult)
	if opt.ServiceCategory {
		result = listServiceCategoryInstances(serviceInstances.Count, insts)
		return result, nil
	}

	op := &hostAndServiceInstsOption{BizID: opt.BizID, ModuleID: opt.ModuleID, ProcessTemplateId: opt.ProcessTemplateId}

	serviceRelationMap, relations, hMap, hostIDs, pTemplates, pTemplateMap, hostWithSrvInstMap, err :=
		ps.getHostAndServiceInstances(ctx, op, mod, insts)
	if err != nil {
		return nil, err
	}

	procID2Detail, attributeMap, err := ps.getProcAndAttributeMap(ctx, relations.Info)
	if err != nil {
		return nil, err
	}

	// record the used process template for checking whether a new process template has been added to service template.
	processTemplateReferenced := make(map[int64]struct{})

	num, flag := 0, false

	// compare the process instance with it's process template one by one in a service instance.
	for _, serviceInstance := range insts {
		if num > metadata.ServiceInstancesMaxNum {
			break
		}
		//relations := serviceRelationMap[serviceInstance[common.BKFieldID].(int64)]
		relations := serviceRelationMap[serviceInstance.id]

		for _, relation := range relations {

			processTemplateReferenced[relation.ProcessTemplateID] = struct{}{}

			process, ok := procID2Detail[relation.ProcessID]
			if !ok {
				process = new(metadata.Process)
			}
			processName := ""
			if process.ProcessName != nil {
				processName = *process.ProcessName
			}

			p, exist := pTemplateMap[relation.ProcessTemplateID]

			// 此处如果不存在，那么说明模板已经被删除了，记录涉及到的实例信息，只要记录实例信息即可，此步骤不需要记录所有的内容 当有多个
			// 进程模板都被删除的场景下，关系表中的进程模板id即relation.ProcessTemplateID都为0，我们只需要计算前端指定的进程模板即可
			if !exist && opt.ProcessTemplateId == 0 && opt.ProcessTemplateName == processName {

				result.ServiceInstances = append(result.ServiceInstances, metadata.ServiceInstancesInfo{
					Id: serviceInstance.id, Name: serviceInstance.name})
				flag = true
				break
			}

			if !exist {
				continue
			}

			_, isChanged, dErr := ps.Logic.DiffWithProcessTemplate(p.Property, process, hMap[serviceInstance.hostId], attributeMap, false)
			if dErr != nil {
				blog.Errorf("diff fail, process ID: %d err: %v, rid: %s", relation.ProcessID, err, ctx.Kit.Rid)
				return nil, errors.New(common.CCErrCommParamsInvalid, dErr.Error())
			}

			if !isChanged {
				continue
			}

			result.ServiceInstances = append(result.ServiceInstances, metadata.ServiceInstancesInfo{
				Id: serviceInstance.id, Name: serviceInstance.name})
			flag = true
		}
		if flag {
			num++
			flag = false
			continue
		}
		_, exist := processTemplateReferenced[opt.ProcessTemplateId]
		if exist {
			continue
		}

		result.ServiceInstances = append(result.ServiceInstances, metadata.ServiceInstancesInfo{
			Id: serviceInstance.id, Name: serviceInstance.name})
		num++
	}

	// 加入没有实例化的主机.
	if len(serviceInstances.Info) == 0 && len(hMap) > 0 {
		result.ServiceInstances = ps.handleAddedServiceInsts(mod, hMap, opt, hostIDs, hostWithSrvInstMap, pTemplates)
		num = len(hMap)
	}
	if num > metadata.ServiceInstancesMaxNum {
		result.TotalCount = metadata.ServiceInstancestotalCount
	} else {
		result.TotalCount = strconv.FormatInt(int64(num), 10)
	}

	return result, nil
}

func (ps *ProcServer) handleAddedServiceInsts(module *metadata.ModuleInst, hostMap map[int64]map[string]interface{},
	option *metadata.ListServiceInstancesOption, hostIDs []int64, hostWithSrvInstMap map[int64]struct{},
	processTemplates *metadata.MultipleProcessTemplate) []metadata.ServiceInstancesInfo {
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
	added := make([]metadata.ServiceInstancesInfo, 0)
	for _, hostID := range hostIDs {
		if _, exists := hostWithSrvInstMap[hostID]; exists {
			continue
		}

		srvInstName := util.GetStrByInterface(hostMap[hostID][common.BKHostInnerIPField]) + srvInstNameSuffix
		added = append(added, metadata.ServiceInstancesInfo{
			Name: srvInstName,
		})
	}
	return added
}

type serviceListBrief struct {
	hostId int64
	id     int64
	name   string
}

func (ps *ProcServer) serviceTemplateGeneralDiff(ctx *rest.Contexts, op metadata.DiffWithOneModuleOption,
	pTemplates *metadata.MultipleProcessTemplate) (*metadata.ServiceTemplateGeneralDiff, errors.CCErrorCoder) {

	rid := ctx.Kit.Rid
	if op.ModuleID == 0 {
		blog.Errorf("module id empty, option: %+v, rid: %s", op, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	module, err := ps.getModule(ctx.Kit, op.ModuleID)
	if err != nil {
		blog.Errorf("get module detail failed, moduleID: %d, err: %+v, rid: %s", op.ModuleID, err, rid)
		return nil, err
	}

	if module.ServiceTemplateID == 0 {
		blog.Errorf("module %d has no service template, option: %s, rid: %s", op.ModuleID, rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	sInsts, err := ps.getServiceInstances(ctx, module, nil, []string{common.BKFieldID, common.BKHostIDField})
	if err != nil {
		return nil, err
	}

	flag, err := ps.calculateModuleAttributeGeneralDifference(ctx.Kit.Ctx, ctx.Kit.Header, *module)
	if err != nil {
		blog.Errorf("calculate attribute difference failed, module: %s, err: %v, rid: %s", module, err, rid)
		return nil, err
	}

	moduleDifference := &metadata.ServiceTemplateGeneralDiff{
		Changed: make([]metadata.ProcessGeneralInfo, 0),
		Added:   make([]metadata.ProcessGeneralInfo, 0),
		Removed: make([]metadata.ProcessGeneralInfo, 0),
	}
	moduleDifference.ChangedAttribute = flag

	// processTemplates->pTemplateMap
	pTemplateMap := make(map[int64]*metadata.ProcessTemplate)

	for idx, pTemplate := range pTemplates.Info {
		pTemplateMap[pTemplate.ID] = &pTemplates.Info[idx]
	}

	hostMap, err := ps.getHostInfo(ctx, op.BizID, op.ModuleID)
	if err != nil {
		blog.Errorf("get host fail option: %+v,err: %v,rid: %s", op, err, rid)
		return nil, err
	}

	insts, e := getServiceListBrief(sInsts.Info)
	if e != nil {
		return nil, errors.New(common.CCErrCommDBSelectFailed, e.Error())
	}

	diffTemp, err := ps.calculateGeneralDiff(ctx, module, hostMap, pTemplateMap, insts)
	if err != nil {
		return nil, err
	}

	moduleDifference.Removed = diffTemp.Removed
	moduleDifference.Added = diffTemp.Added
	moduleDifference.Changed = diffTemp.Changed

	return moduleDifference, nil
}

// calculateModuleAttributeGeneralDifference calculate whether the service classification has changed.
func (ps *ProcServer) calculateModuleAttributeGeneralDifference(ctx context.Context, header http.Header,
	module metadata.ModuleInst) (bool, errors.CCErrorCoder) {

	if module.ServiceTemplateID == common.ServiceTemplateIDNotSet {
		return false, nil
	}

	serviceTpl, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(ctx, header, module.ServiceTemplateID)
	if err != nil {
		return false, err
	}

	if module.ServiceCategoryID == serviceTpl.ServiceCategoryID && module.ModuleName == serviceTpl.Name {
		return false, nil
	}

	return true, nil
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

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
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

	// {ServiceInstanceID: []Process}
	serviceInstance2ProcessMap := make(map[int64][]*metadata.Process)
	// {ServiceInstanceID: {ProcessTemplateID: true}}
	serviceInstanceWithTemplateMap := make(map[int64]map[int64]bool)
	// {ServiceInstanceID: HostID}
	serviceInstance2HostMap := make(map[int64]int64)
	// {ServiceInstanceID: ServiceTemplateID}
	serviceInstanceTemplateMap := make(map[int64]int64)
	hostWithSrvInstMap := make(map[int64]struct{})
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstance2ProcessMap[serviceInstance.ID] = make([]*metadata.Process, 0)
		serviceInstanceWithTemplateMap[serviceInstance.ID] = make(map[int64]bool)
		serviceInstance2HostMap[serviceInstance.ID] = serviceInstance.HostID
		serviceInstanceTemplateMap[serviceInstance.ID] = serviceInstance.ServiceTemplateID
		hostWithSrvInstMap[serviceInstance.HostID] = struct{}{}
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

	srvTempWithProTempMap := make(map[int64]struct{})
	processTemplateMap := make(map[int64]*metadata.ProcessTemplate)
	for idx, t := range processTemplate.Info {
		processTemplateMap[t.ID] = &processTemplate.Info[idx]
		srvTempWithProTempMap[t.ServiceTemplateID] = struct{}{}
	}

	// step 2:
	// construct map {moduleID -> []hostID}
	hostRelationOpt := &metadata.HostModuleRelationRequest{
		ApplicationID: bizID,
		ModuleIDArr:   syncOption.ModuleIDs,
		Fields:        []string{common.BKModuleIDField, common.BKHostIDField},
	}
	hostRelationRes, rawErr := ps.Logic.CoreAPI.CoreService().Host().GetHostModuleRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		hostRelationOpt)
	if rawErr != nil {
		blog.Errorf("get host relation failed, err: %v, input: %#v, rid: %s", rawErr, hostRelationOpt, ctx.Kit.Rid)
		return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := hostRelationRes.CCError(); err != nil {
		blog.Errorf("get host relation failed, err: %v, input: %#v, rid: %s", rawErr, hostRelationOpt, ctx.Kit.Rid)
		return err
	}

	hostIDs := make([]int64, len(hostRelationRes.Data.Info))
	moduleHostMap := make(map[int64][]int64)
	for index, relation := range hostRelationRes.Data.Info {
		hostIDs[index] = relation.HostID
		moduleHostMap[relation.ModuleID] = append(moduleHostMap[relation.ModuleID], relation.HostID)
	}

	// step 3: handle all the service instances that need to be added
	srvInstToAdd := make([]*metadata.ServiceInstance, 0)
	for _, module := range modules {
		if _, exists := srvTempWithProTempMap[module.ServiceTemplateID]; !exists {
			continue
		}

		for _, hostID := range moduleHostMap[module.ModuleID] {
			if _, exists := hostWithSrvInstMap[hostID]; exists {
				continue
			}
			instance := &metadata.ServiceInstance{
				BizID:             bizID,
				ServiceTemplateID: module.ServiceTemplateID,
				ModuleID:          module.ModuleID,
				HostID:            hostID,
			}
			srvInstToAdd = append(srvInstToAdd, instance)
		}
	}

	_, err = ps.CoreAPI.CoreService().Process().CreateServiceInstances(ctx.Kit.Ctx, ctx.Kit.Header, srvInstToAdd)
	if err != nil {
		blog.Errorf("create service instances(%#v) failed, err: %v, rid: %s", srvInstToAdd, err, rid)
		return err
	}

	// get service templates
	serviceTemplates, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header,
		&metadata.ListServiceTemplateOption{
			BusinessID:         bizID,
			ServiceTemplateIDs: serviceTemplateIDs,
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		})
	if err != nil {
		blog.Errorf("list service templates failed, ids: %+v, err: %v, rid: %s", serviceTemplateIDs, err, rid)
		return err
	}

	// step 4:
	// update module service category and name field TODO: remove this
	for _, serviceTemplate := range serviceTemplates.Info {
		updateModules := make([]int64, 0)
		for _, module := range serviceTemplateModuleMap[serviceTemplate.ID] {
			if module == nil {
				continue
			}
			if serviceTemplate.ServiceCategoryID != module.ServiceCategoryID ||
				serviceTemplate.Name != module.ModuleName {
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
		resp, e := ps.CoreAPI.CoreService().Instance().UpdateInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDModule, moduleUpdateOption)
		if e != nil {
			blog.Errorf("update module failed, option: %#v, err: %v, rid: %s", moduleUpdateOption, e, rid)
			return ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}
		if ccErr := resp.CCError(); ccErr != nil {
			blog.Errorf("update module failed, option: %#v, err: %v, rid: %s", moduleUpdateOption, ccErr, rid)
			return ccErr
		}
	}

	if len(serviceInstanceIDs) == 0 {
		return nil
	}

	// step5:
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

	// step 6:
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

	// step 7:
	// rearrange the service instance with process instance.
	processInstanceWithTemplateMap := make(map[int64]int64)
	for _, r := range relations.Info {
		p, exist := processInstanceMap[r.ProcessID]
		if !exist {
			// something is wrong, but can this process instance,
			// but we can find it in the process instance relation.
			blog.Warnf("force sync process instance according to process template: %d, but can not find the process instance: %d, rid: %s", r.ProcessTemplateID, r.ProcessID, rid)
			continue
		}
		if _, exist := serviceInstanceWithTemplateMap[r.ServiceInstanceID]; !exist {
			// something is wrong, service instance is not exist, but we can find it in the process instance relation
			blog.Warnf("relation: %#v has a service instance that is not exist, rid: %s", r, rid)
			continue
		}
		serviceInstance2ProcessMap[r.ServiceInstanceID] = append(serviceInstance2ProcessMap[r.ServiceInstanceID], p)
		processInstanceWithTemplateMap[r.ProcessID] = r.ProcessTemplateID
		serviceInstanceWithTemplateMap[r.ServiceInstanceID][r.ProcessTemplateID] = true
	}

	// step 8:
	// construct map {hostID ==> host}
	hostMap, err := ps.Logic.GetHostIPMapByID(ctx.Kit, hostIDs)
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
		for _, process := range processes {
			processTemplateID := processInstanceWithTemplateMap[process.ProcessID]
			template, exist := processTemplateMap[processTemplateID]
			if !exist || template.ServiceTemplateID != serviceInstanceTemplateMap[serviceInstanceID] {
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
					blog.ErrorJSON("sync service instance, but extract process change info failed, err: %s, "+
						"process: %s, rid: %s", err, process, rid)
					if firstErr == nil {
						firstErr = errors.New(common.CCErrCommParamsInvalid, err.Error())
					}
					return
				}

				if !changed {
					return
				}

				if err := ps.Logic.UpdateProcessInstance(ctx.Kit, process.ProcessID, proc); err != nil {
					blog.ErrorJSON("UpdateProcessInstance failed, processID: %s, process: %s, err: %s, rid: %s",
						process.ProcessID, proc, err, rid)
					if firstErr == nil {
						firstErr = err
					}
					return
				}

			}(process, hostMap[serviceInstance2HostMap[serviceInstanceID]])
		}

		if _, exists := srvTempWithProTempMap[serviceInstanceTemplateMap[serviceInstanceID]]; !exists {
			removedSvrInstIDs = append(removedSvrInstIDs, serviceInstanceID)
		}
	}

	wg.Wait()
	if firstErr != nil {
		return firstErr
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

	// delete service instances whose processes are all removed
	if len(removedSvrInstIDs) > 0 {
		deleteOption := &metadata.CoreDeleteServiceInstanceOption{
			BizID:              bizID,
			ServiceInstanceIDs: removedSvrInstIDs,
		}
		err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, deleteOption)
		if err != nil {
			blog.Errorf("delete service instances: %+v failed, err: %v, rid: %s", removedSvrInstIDs, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrProcDeleteServiceInstancesFailed)
		}
	}

	// step 10:
	// check if a new process is added to the service template.
	// if true, then create a new process instance for every service instance with process template's default value.
	processDatas := make([]map[string]interface{}, 0)
	procInstRelations := make([]*metadata.ProcessInstanceRelation, 0)
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
			newProcess, generateErr := processTemplate.NewProcess(bizID, ctx.Kit.SupplierAccount,
				hostMap[serviceInstance2HostMap[svcID]])
			if generateErr != nil {
				blog.ErrorJSON("sync service instance by template, but generate process instance by template "+
					"%s failed, err: %s, rid: %s", processTemplate, generateErr, rid)
				return errors.New(common.CCErrCommParamsInvalid, generateErr.Error())
			}
			processDatas = append(processDatas, newProcess.Map())
			procInstRelations = append(procInstRelations, &metadata.ProcessInstanceRelation{
				BizID:             bizID,
				ServiceInstanceID: svcID,
				ProcessTemplateID: processTemplateID,
				HostID:            serviceInstance2HostMap[svcID],
			})
		}
	}

	if len(processDatas) > 0 {
		// create process instances in batch
		processIDs, err := ps.Logic.CreateProcessInstances(ctx.Kit, processDatas)
		if err != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, CreateProcessInstances err: %s, processDatas: %s, rid: %s", err, processDatas, ctx.Kit.Rid)
			return ctx.Kit.CCError.CCError(common.CCErrSyncServiceInstanceByTemplateFailed)
		}

		if len(processIDs) != len(procInstRelations) {
			blog.Error("syncServiceInstanceByTemplate failed, the count of processIDs must be equal with the count of procInstRelations")
			return nil
		}

		// create process instance relations in batch
		for idx, processID := range processIDs {
			procInstRelations[idx].ProcessID = processID
		}
		_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelations(ctx.Kit.Ctx, ctx.Kit.Header, procInstRelations)
		if err != nil {
			blog.ErrorJSON("syncServiceInstanceByTemplate failed, CreateProcessInstanceRelations err: %s, relations: %s, rid: %s", err, procInstRelations, ctx.Kit.Rid)
			return err
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
