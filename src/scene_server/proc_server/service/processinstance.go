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
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	processhook "configcenter/src/thirdparty/hooks/process"
)

// CreateProcessInstances TODO
func (ps *ProcServer) CreateProcessInstances(ctx *rest.Contexts) {
	input := new(metadata.CreateRawProcessInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.Processes) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsIsInvalid, "not set processes"))
		blog.Infof("no process to create, return")
		return
	}
	if len(input.Processes) > common.BKMaxUpdateOrCreatePageSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "create process instances",
			common.BKMaxUpdateOrCreatePageSize))
		return
	}

	processIDs := make([]int64, 0)
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		processIDs, err = ps.createProcessInstances(ctx, input)
		if err != nil {
			blog.Errorf("create process instance failed, serviceInstanceID: %d, input: %+v, err: %+v", input.ServiceInstanceID, input, err)
			return err
		}

		// generate and save audit log after processes are created
		audit := auditlog.NewSvcInstAudit(ps.CoreAPI.CoreService())
		if err = audit.WithServiceInstanceByIDs(ctx.Kit, input.BizID, []int64{input.ServiceInstanceID},
			[]string{common.BKFieldID}); err != nil {
			return err
		}

		relations := make([]metadata.ProcessInstanceRelation, len(processIDs))
		for index, procID := range processIDs {
			relations[index] = metadata.ProcessInstanceRelation{
				BizID:             input.BizID,
				ProcessID:         procID,
				ServiceInstanceID: input.ServiceInstanceID,
			}
		}

		genProcAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditCreate)
		err = audit.WithProcByRelations(genProcAuditParam, relations, nil)
		if err != nil {
			return err
		}

		genAuditParam := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate)
		auditLogs := audit.GenerateAuditLog(genAuditParam)
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(processIDs)
}

func (ps *ProcServer) createProcessInstances(ctx *rest.Contexts, input *metadata.CreateRawProcessInstanceInput) ([]int64, errors.CCErrorCoder) {
	serviceInstance, err := ps.CoreAPI.CoreService().Process().GetServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceInstanceID)
	if err != nil {
		blog.Errorf("create process instance failed, get service instance by id failed, serviceInstanceID: %d, err: %v, rid: %s", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, err
	}
	if serviceInstance.BizID != input.BizID {
		blog.Errorf("create process instance with raw, biz id from input not equal with service instance, rid: %s", ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	if serviceInstance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		blog.Errorf("create process instance failed, create process instance on service instance initialized by template forbidden, serviceInstanceID: %d, err: %v, rid: %s", input.ServiceInstanceID, err, ctx.Kit.Rid)
		return nil, ctx.Kit.CCError.CCError(common.CCErrProcEditProcessInstanceCreateByTemplateForbidden)
	}

	processIDs := make([]int64, 0)
	processDatas := make([]map[string]interface{}, len(input.Processes))
	for idx, item := range input.Processes {
		now := time.Now()
		item.ProcessData[common.BKProcessIDField] = int64(0)
		item.ProcessData[common.BKAppIDField] = input.BizID
		item.ProcessData[common.BkSupplierAccount] = ctx.Kit.SupplierAccount
		item.ProcessData[common.CreateTimeField] = now
		item.ProcessData[common.LastTimeField] = now
		item.ProcessData[common.BKServiceInstanceIDField] = input.ServiceInstanceID

		processDatas[idx] = item.ProcessData
	}

	processIDs, err = ps.Logic.CreateProcessInstances(ctx.Kit, processDatas)
	if err != nil {
		blog.Errorf("create process instance failed, create process failed, serviceInstanceID: %d, processDatas: %+v, err: %v, rid: %s", input.ServiceInstanceID, processDatas, err, ctx.Kit.Rid)
		return nil, err
	}

	relations := make([]*metadata.ProcessInstanceRelation, len(processIDs))
	for idx, processID := range processIDs {
		relation := &metadata.ProcessInstanceRelation{
			BizID:             input.BizID,
			ProcessID:         processID,
			ProcessTemplateID: common.ServiceTemplateIDNotSet,
			ServiceInstanceID: serviceInstance.ID,
			HostID:            serviceInstance.HostID,
		}
		relations[idx] = relation
	}
	_, err = ps.CoreAPI.CoreService().Process().CreateProcessInstanceRelations(ctx.Kit.Ctx, ctx.Kit.Header, relations)
	if err != nil {
		blog.ErrorJSON("create service instance relations, CreateProcessInstanceRelations err: %s, relations:%s, rid: %s", err, relations, ctx.Kit.Rid)
		return nil, err
	}

	return processIDs, nil
}

// UpdateProcessInstancesByIDs TODO
func (ps *ProcServer) UpdateProcessInstancesByIDs(ctx *rest.Contexts) {
	input := metadata.UpdateProcessByIDsInput{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: input.ProcessIDs,
		},
	}
	reqParam := &metadata.QueryCondition{
		Condition: filter,
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	processResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, reqParam)
	if nil != err {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "UpdateProcessInstancesByIDs failed, "+
			"reqParam: %#v, err: %+v", reqParam, err)
		return
	}

	// generate audit log before processes are updated
	relOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID: input.BizID,
		ProcessIDs: input.ProcessIDs,
		Page:       metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header, relOpt)
	if err != nil {
		blog.Errorf("get process relations failed, option: %+v, err: %v, rid: %s", relOpt, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	svcInstIDs := make([]int64, 0)
	for _, relation := range relRes.Info {
		svcInstIDs = append(svcInstIDs, relation.ServiceInstanceID)
	}

	audit := auditlog.NewSvcInstAudit(ps.CoreAPI.CoreService())
	if err = audit.WithServiceInstanceByIDs(ctx.Kit, input.BizID, svcInstIDs, []string{common.BKFieldID}); err != nil {
		ctx.RespAutoError(err)
		return
	}

	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(ctx.Kit, metadata.AuditUpdate).
		WithUpdateFields(input.UpdateData)
	err = audit.WithProc(generateAuditParameter, processResult.Info, relRes.Info)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	auditLogs := audit.GenerateAuditLog(generateAuditParameter)

	// parse update data and update the processes
	raws := make([]map[string]interface{}, 0)
	for _, process := range processResult.Info {
		for k, v := range input.UpdateData {
			process[k] = v
		}
		raws = append(raws, process)
	}

	if len(raws) == 0 {
		ctx.RespEntity([]int64{})
		return
	}

	updateInput := metadata.UpdateRawProcessInstanceInput{
		BizID: input.BizID,
		Raw:   raws,
	}

	var result []int64
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		result, err = ps.updateProcessInstances(ctx, updateInput)
		if err != nil {
			return err
		}

		// save audit logs
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(result)
}

// UpdateProcessInstances TODO
func (ps *ProcServer) UpdateProcessInstances(ctx *rest.Contexts) {
	input := metadata.UpdateRawProcessInstanceInput{}
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.Raw) == 0 {
		ctx.RespEntity([]int64{})
		return
	}

	if len(input.Raw) > common.BKMaxUpdateOrCreatePageSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "update process instances",
			common.BKMaxUpdateOrCreatePageSize))
		return
	}
	// generate audit log before processes are updated
	auditLogs, err := ps.generateUpdateProcessAudit(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	var result []int64
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		// update process instances
		var err error
		result, err = ps.updateProcessInstances(ctx, input)
		if err != nil {
			return err
		}

		// save audit log
		audit := auditlog.NewSvcInstAudit(ps.CoreAPI.CoreService())
		if err := audit.SaveAuditLog(ctx.Kit, auditLogs...); err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(result)
}

// generateUpdateProcessAudit generate audit logs for process update operation
func (ps *ProcServer) generateUpdateProcessAudit(kit *rest.Kit, input metadata.UpdateRawProcessInstanceInput) (
	[]metadata.AuditLog, error) {

	// generate audit log before processes are updated
	// get process ids
	procIDs := make([]int64, 0)
	procMap := make(map[int64]mapstr.MapStr)
	for _, proc := range input.Raw {
		procID, err := util.GetInt64ByInterface(proc[common.BKProcessIDField])
		if err != nil {
			blog.Errorf("parse process(%+v) id failed, err: %v, rid: %s", proc, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
		}
		procIDs = append(procIDs, procID)
		procMap[procID] = proc
	}

	// get process relations, then get service instance ids by relations, set service instance data in audit logs
	relOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID: input.BizID,
		ProcessIDs: procIDs,
		Page:       metadata.BasePage{Limit: common.BKNoLimit},
	}
	relRes, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, relOpt)
	if err != nil {
		blog.Errorf("get process relations failed, option: %+v, err: %v, rid: %s", relOpt, err, kit.Rid)
		return nil, err
	}

	svcInstIDs := make([]int64, 0)
	procRelationMap := make(map[int64]metadata.ProcessInstanceRelation)
	for _, relation := range relRes.Info {
		svcInstIDs = append(svcInstIDs, relation.ServiceInstanceID)
		procRelationMap[relation.ProcessID] = relation
	}

	audit := auditlog.NewSvcInstAudit(ps.CoreAPI.CoreService())
	if err = audit.WithServiceInstanceByIDs(kit, input.BizID, svcInstIDs, []string{common.BKFieldID}); err != nil {
		return nil, err
	}

	// get process data before updating, generate audit logs by these
	procOpt := &metadata.QueryCondition{
		Condition: map[string]interface{}{common.BKProcessIDField: map[string]interface{}{common.BKDBIN: procIDs}},
		Page:      metadata.BasePage{Limit: common.BKNoLimit},
	}
	procRes, rawErr := ps.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header,
		common.BKInnerObjIDProc, procOpt)
	if rawErr != nil {
		blog.Errorf("get process data failed, option: %+v, err: %v, rid: %s", procOpt, rawErr, kit.Rid)
		return nil, rawErr
	}

	for _, data := range procRes.Info {
		procID, err := util.GetInt64ByInterface(data[common.BKProcessIDField])
		if err != nil {
			blog.Errorf("parse previous process(%+v) id failed, err: %v, rid: %s", data, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
		}

		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).
			WithUpdateFields(procMap[procID])
		err = audit.WithProc(generateAuditParameter, []mapstr.MapStr{data},
			[]metadata.ProcessInstanceRelation{procRelationMap[procID]})
		if err != nil {
			return nil, err
		}
	}

	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate)
	auditLogs := audit.GenerateAuditLog(generateAuditParameter)

	return auditLogs, nil
}

func (ps *ProcServer) updateProcessInstances(ctx *rest.Contexts, input metadata.UpdateRawProcessInstanceInput) ([]int64, errors.CCErrorCoder) {
	rid := ctx.Kit.Rid
	bizID := input.BizID

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
	hostIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		foundProcessIDs = append(foundProcessIDs, relation.ProcessID)
		if relation.ProcessTemplateID != common.ServiceTemplateIDNotSet {
			hostIDs = append(hostIDs, relation.HostID)
		}
	}
	invalidProcessIDs := make([]string, 0)
	for _, processID := range processIDs {
		if !util.InArray(processID, foundProcessIDs) {
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
		if _, exist := processTemplateMap[relation.ProcessTemplateID]; exist {
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
	for i := range relations.Info {
		process2ServiceInstanceMap[relations.Info[i].ProcessID] = &relations.Info[i]
	}

	hostMap, err := ps.Logic.GetHostIPMapByID(ctx.Kit, hostIDs)
	if err != nil {
		return nil, err
	}

	processDataMap := make(map[int64]map[string]interface{})
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
		if !exist {
			err := ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessIDField)
			blog.ErrorJSON("update process instance failed, process related service instance not found, process: %s, err: %s, rid: %s", process, err, rid)
			return nil, err
		}

		var processData map[string]interface{}
		if relation.ProcessTemplateID == common.ServiceTemplateIDNotSet {
			process.BusinessID = bizID
			var err error
			processData, err = mapstruct.Struct2Map(process)
			if nil != err {
				blog.Errorf("UpdateProcessInstances failed, json Unmarshal process failed, processData: %s, err: %+v, rid: %s", processData, err, ctx.Kit.Rid)
				return nil, ctx.Kit.CCError.CCError(common.CCErrCommJsonDecode)
			}
			delete(processData, common.BKProcessIDField)
			delete(processData, common.MetadataField)
			delete(processData, common.LastTimeField)
			delete(processData, common.CreateTimeField)
		} else {
			processTemplate, exist := processTemplateMap[relation.ProcessTemplateID]
			if !exist {
				err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
				blog.Errorf("update process instance failed, process related template not found, relation: %+v, err: %v, rid: %s", relation, err, rid)
				return nil, err
			}
			var compareErr error
			processData, compareErr = processTemplate.ExtractInstanceUpdateData(&process, hostMap[relation.HostID])
			if compareErr != nil {
				blog.ErrorJSON("extract process(%s) update data failed, err: %s, rid: %s", process, err, rid)
				return nil, errors.New(common.CCErrCommParamsInvalid, compareErr.Error())
			}
			clearFields = processTemplate.GetEditableFields(clearFields)
		}
		// set field value as nil
		for _, field := range clearFields {
			processData[field] = nil
		}

		processInfo := new(metadata.Process)
		if err := mapstr.DecodeFromMapStr(&processInfo, processData); err != nil {
			blog.ErrorJSON("parse update process data failed, data: %s, err: %v, rid: %s", processData, err, rid)
			return nil, ctx.Kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
		}

		if err := ps.validateProcessInstance(ctx.Kit, processInfo); err != nil {
			blog.ErrorJSON("validate update process failed, err: %s, data: %s, rid: %s", err, processInfo, rid)
			return nil, err
		}
		processDataMap[process.ProcessID] = processData
	}

	var wg sync.WaitGroup
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)

	for processID := range processDataMap {
		pipeline <- true
		wg.Add(1)

		go func(processID int64, processData map[string]interface{}) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			err := ps.Logic.UpdateProcessInstance(ctx.Kit, processID, processData)
			if err != nil {
				blog.ErrorJSON("UpdateProcessInstance failed, processID: %s, process: %s, err: %s, rid: %s", processID, processData, err, rid)
				if firstErr == nil {
					firstErr = err
				}
				return
			}

		}(processID, processDataMap[processID])
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return processIDs, nil
}

// checkHostsInModule check if hosts are in the business module, can only create service instance for hosts in it
func (ps *ProcServer) checkHostsInModule(kit *rest.Kit, bizID, moduleID int64, hostIDs []int64) errors.CCErrorCoder {
	hostFilter := &metadata.DistinctHostIDByTopoRelationRequest{
		ApplicationIDArr: []int64{bizID},
		ModuleIDArr:      []int64{moduleID},
		HostIDArr:        hostIDs,
	}

	hitHostIDs, err := ps.CoreAPI.CoreService().Host().GetDistinctHostIDByTopology(kit.Ctx, kit.Header, hostFilter)
	if err != nil {
		blog.ErrorJSON("check host in module failed, filter: %s, err: %s, rid: %s", hostFilter, err, kit.Rid)
		return err
	}

	hostIDHit := make(map[int64]struct{})
	for _, id := range hitHostIDs {
		hostIDHit[id] = struct{}{}
	}

	invalidHost := make([]int64, 0)
	for _, hostID := range hostIDs {
		if _, exists := hostIDHit[hostID]; !exists {
			invalidHost = append(invalidHost, hostID)
		}
	}
	if len(invalidHost) > 0 {
		return kit.CCError.CCErrorf(common.CCErrHostModuleConfigNotMatch, invalidHost)
	}
	return nil
}

func (ps *ProcServer) getModule(kit *rest.Kit, moduleID int64) (*metadata.ModuleInst, errors.CCErrorCoder) {
	filter := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKModuleIDField: moduleID,
		},
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	err := ps.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, filter, moduleRes)
	if err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, err
	}

	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", filter, err, kit.Rid)
		return nil, err
	}

	if len(moduleRes.Data.Info) != 1 {
		blog.Errorf("get not one module by id $d, data: %#v, rid: %s", moduleID, moduleRes.Data.Info, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoGetModuleFailed, "get none or multiple modules")
	}

	return &moduleRes.Data.Info[0], nil
}

var (
	ipRegex = `^((1?\d{1,2}|2[0-4]\d|25[0-5])[.]){3}(1?\d{1,2}|2[0-4]\d|25[0-5])$`
)

func (ps *ProcServer) validateProcessInstance(kit *rest.Kit, process *metadata.Process) errors.CCErrorCoder {
	if process.ProcessName != nil && (len(*process.ProcessName) == 0 ||
		len(*process.ProcessName) > common.NameFieldMaxLength) {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessNameField)
	}
	if process.FuncName != nil && (len(*process.FuncName) == 0 ||
		len(*process.ProcessName) > common.NameFieldMaxLength) {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFuncName)
	}

	// validate that process bind info must have ip and port and protocol
	for _, bindInfo := range process.BindInfo {
		if bindInfo.Std.IP == nil || len(*bindInfo.Std.IP) == 0 {
			if err := processhook.ValidateProcessBindIPEmptyHook(); err != nil {
				return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKProcBindInfo+"."+common.BKIP)
			}
		} else {
			matched, err := regexp.MatchString(ipRegex, *bindInfo.Std.IP)
			if err != nil || !matched {
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcBindInfo+"."+common.BKIP)
			}
		}

		port := (*metadata.PropertyPortValue)(bindInfo.Std.Port)
		if err := port.Validate(); err != nil {
			return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKProcBindInfo+"."+common.BKPort)
		}

		protocol := (*metadata.ProtocolType)(bindInfo.Std.Protocol)
		if err := protocol.Validate(); err != nil {
			return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKProcBindInfo+"."+common.BKProtocol)
		}
	}

	return nil
}

// DeleteProcessInstance TODO
func (ps *ProcServer) DeleteProcessInstance(ctx *rest.Contexts) {
	input := new(metadata.DeleteProcessInstanceInServiceInstanceInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(input.ProcessInstanceIDs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKProcessIDField))
		return
	}

	if len(input.ProcessInstanceIDs) > common.BKMaxDeletePageSize {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommXXExceedLimit, "delete process instance",
			common.BKMaxDeletePageSize))
		return
	}

	listOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID: input.BizID,
		ProcessIDs: input.ProcessInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		listOption)
	if err != nil {
		blog.Errorf("list process relation failed, option: %#v, err: %v, rid: %s", listOption, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	templateProcessIDs := make([]string, 0)
	for _, relation := range relations.Info {
		// get processes that are created by template, can not delete them
		if relation.ProcessTemplateID != common.ServiceTemplateIDNotSet {
			templateProcessIDs = append(templateProcessIDs, strconv.FormatInt(relation.ProcessID, 10))
		}
	}

	if len(templateProcessIDs) > 0 {
		invalidProcesses := strings.Join(templateProcessIDs, ",")
		blog.Errorf("some process: %s are initialized by template, rid: %s", invalidProcesses, ctx.Kit.Rid)
		err := ctx.Kit.CCError.CCErrorf(common.CCErrCoreServiceShouldNotRemoveProcessCreateByTemplate, invalidProcesses)
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		return ps.deleteProcessInstance(ctx.Kit, input.BizID, input.ProcessInstanceIDs, relations.Info)
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (ps *ProcServer) deleteProcessInstance(kit *rest.Kit, bizID int64, procIDs []int64,
	relations []metadata.ProcessInstanceRelation) error {

	if len(procIDs) == 0 {
		return nil
	}

	// set process data for audit log before they are deleted
	audit := auditlog.NewSvcInstAudit(ps.CoreAPI.CoreService())
	genAuditParam := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	if err := audit.WithProcByRelations(genAuditParam, relations, nil); err != nil {
		return err
	}

	// delete process relations at the same time.
	deleteOpt := metadata.DeleteProcessInstanceRelationOption{
		BusinessID: &bizID,
		ProcessIDs: procIDs,
	}
	err := ps.CoreAPI.CoreService().Process().DeleteProcessInstanceRelation(kit.Ctx, kit.Header, deleteOpt)
	if err != nil {
		blog.Errorf("delete process relation failed, err: %v, option: %#v, rid: %s", err, deleteOpt, kit.Rid)
		return err
	}

	if err := ps.Logic.DeleteProcessInstanceBatch(kit, procIDs); err != nil {
		blog.Errorf("delete process instance by ids: %+v failed, err: %v", procIDs, err)
		return err
	}

	// skip checking for service instance to delete or saving audit logs if no relation exists
	if len(relations) == 0 {
		return nil
	}

	updatedSvcInstIDs, delSvcInstIDs, err := ps.checkIfSvcInstNeedCascadeDelete(kit, bizID, relations)
	if err != nil {
		return err
	}

	// get service instances by ids for audit logs
	serviceInstances, err := audit.GetSvcInstByIDs(kit, bizID, append(updatedSvcInstIDs, delSvcInstIDs...), nil)
	if err != nil {
		return err
	}
	svcInstMap := make(map[int64]metadata.ServiceInstance)
	for _, svcInst := range serviceInstances {
		svcInstMap[svcInst.ID] = svcInst
	}

	// generate audit logs for updated service instances
	auditLogs := make([]metadata.AuditLog, 0)
	if len(updatedSvcInstIDs) > 0 {
		updatedServiceInstances := make([]metadata.ServiceInstance, len(updatedSvcInstIDs))
		for index, svcInstID := range updatedSvcInstIDs {
			updatedServiceInstances[index] = svcInstMap[svcInstID]
		}
		audit.WithServiceInstance(updatedServiceInstances)
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate)
		auditLogs = append(auditLogs, audit.GenerateAuditLog(generateAuditParameter)...)
	}

	if len(delSvcInstIDs) > 0 {
		// generate service instance audit log before they are deleted
		deletedServiceInstances := make([]metadata.ServiceInstance, len(delSvcInstIDs))
		for index, svcInstID := range delSvcInstIDs {
			deletedServiceInstances[index] = svcInstMap[svcInstID]
		}
		audit.WithServiceInstance(deletedServiceInstances)
		auditLogs = append(auditLogs, audit.GenerateAuditLog(genAuditParam)...)

		// remove the service instances whose last process is deleted
		deleteOption := &metadata.CoreDeleteServiceInstanceOption{
			BizID:              bizID,
			ServiceInstanceIDs: delSvcInstIDs,
		}
		err = ps.CoreAPI.CoreService().Process().DeleteServiceInstance(kit.Ctx, kit.Header, deleteOption)
		if err != nil {
			blog.Errorf("delete service instances: %+v failed, err: %v, rid: %s", delSvcInstIDs, err, kit.Rid)
			return err
		}
	}

	// save audit log
	if err := audit.SaveAuditLog(kit, auditLogs...); err != nil {
		return err
	}

	return nil
}

// checkIfSvcInstNeedCascadeDelete returns to be updated and cascade deleted svc inst ids after processes are deleted
func (ps *ProcServer) checkIfSvcInstNeedCascadeDelete(kit *rest.Kit, bizID int64,
	relations []metadata.ProcessInstanceRelation) ([]int64, []int64, errors.CCErrorCoder) {

	if len(relations) == 0 {
		return make([]int64, 0), make([]int64, 0), nil
	}

	// get service instance to processes relations after processes are deleted to check if they have other processes
	svcInstExistsMap := make(map[int64]struct{}, 0)
	serviceInstanceIDs := make([]int64, 0)
	for _, relation := range relations {
		// get service instances to check if all of their processes are deleted
		svcInstExistsMap[relation.ServiceInstanceID] = struct{}{}
		serviceInstanceIDs = append(serviceInstanceIDs, relation.ServiceInstanceID)
	}

	svcOpt := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: serviceInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	svcRelations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, svcOpt)
	if err != nil {
		blog.Errorf("list service relation failed, option: %#v, err: %v, rid: %s", svcOpt, err, kit.Rid)
		return nil, nil, err
	}

	// exclude those service instances that has other process instances
	updatedSvcInstIDs := make([]int64, 0)
	for _, relation := range svcRelations.Info {
		if _, exists := svcInstExistsMap[relation.ServiceInstanceID]; exists {
			delete(svcInstExistsMap, relation.ServiceInstanceID)
			updatedSvcInstIDs = append(updatedSvcInstIDs, relation.ServiceInstanceID)
		}
	}

	if len(svcInstExistsMap) == 0 {
		return updatedSvcInstIDs, make([]int64, 0), nil
	}

	delSvcInstIDs := make([]int64, 0)
	for serviceInstanceID := range svcInstExistsMap {
		delSvcInstIDs = append(delSvcInstIDs, serviceInstanceID)
	}
	return updatedSvcInstIDs, delSvcInstIDs, nil
}

// ListProcessInstances TODO
func (ps *ProcServer) ListProcessInstances(ctx *rest.Contexts) {
	input := new(metadata.ListProcessInstancesOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	processInstanceList, err := ps.Logic.ListProcessInstances(ctx.Kit, input.BizID, input.ServiceInstanceID, nil)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(processInstanceList)
}

// ListProcessInstancesNameIDsInModule get the process id list with its name in a module
func (ps *ProcServer) ListProcessInstancesNameIDsInModule(ctx *rest.Contexts) {
	input := new(metadata.ListProcessInstancesNameIDsOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	option := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameServiceInstance,
		Field:     common.BKFieldID,
		Filter: map[string]interface{}{
			common.BKAppIDField:    input.BizID,
			common.BKModuleIDField: input.ModuleID,
		},
	}
	sIDs, err := ps.CoreAPI.CoreService().Common().GetDistinctField(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		blog.Errorf("GetDistinctField failed, err:%s, option:%#v, rid:%s", err, *option, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if len(sIDs) == 0 {
		ctx.RespEntityWithCount(0, []map[string][]int64{})
		return
	}

	serviceInstanceIDs := make([]int64, len(sIDs))
	for idx, sID := range sIDs {
		if ID, err := strconv.ParseInt(fmt.Sprintf("%v", sID), 10, 64); err == nil {
			serviceInstanceIDs[idx] = ID
		}
	}
	listRelationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID:         input.BizID,
		ServiceInstanceIDs: serviceInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		listRelationOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
			"ListProcessInstancesNameIDsInModule failed, list option: %+v, err: %+v", listRelationOption, err)
		return
	}

	processIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		processIDs = append(processIDs, relation.ProcessID)
	}

	filter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: processIDs,
		},
	}
	if input.ProcessName != "" {
		filter[common.BKProcessNameField] = map[string]interface{}{
			common.BKDBLIKE:    input.ProcessName,
			common.BKDBOPTIONS: "i",
		}
	}
	sort := common.BKProcessNameField
	if input.Page.Sort == "-"+common.BKProcessNameField {
		sort = input.Page.Sort
	}
	reqParam := &metadata.QueryCondition{
		Condition: filter,
		Fields:    []string{common.BKProcessIDField, common.BKProcessNameField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
			Sort:  sort,
		},
	}
	processResult, ccErr := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, reqParam)
	if nil != ccErr {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed,
			"ListProcessInstancesNameIDsInModule failed, reqParam: %#v, err: %+v", reqParam, ccErr)
		return
	}

	processNameIDs := make(map[string][]int64)
	sortedProcessNames := make([]string, 0)

	for _, process := range processResult.Info {
		processID, err := process.Int64(common.BKProcessIDField)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "ListProcessInstancesNameIDsInModule failed, "+
				"process: %#v, err: %+v", process, err)
			return
		}
		processName, err := process.String(common.BKProcessNameField)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "ListProcessInstancesNameIDsInModule failed, "+
				"process: %#v, err: %+v", process, err)
			return
		}
		if _, ok := processNameIDs[processName]; !ok {
			processNameIDs[processName] = make([]int64, 0)
			sortedProcessNames = append(sortedProcessNames, processName)
		}
		processNameIDs[processName] = append(processNameIDs[processName], processID)
	}

	startIndex := input.Page.Start
	if startIndex >= len(sortedProcessNames) {
		ctx.RespEntityWithCount(int64(len(sortedProcessNames)), []map[string][]int64{})
		return
	}

	endindex := startIndex + input.Page.Limit
	if endindex > len(sortedProcessNames) {
		endindex = len(sortedProcessNames)
	}

	ret := make([]metadata.ProcessInstanceNameIDs, endindex-startIndex)
	for idx, name := range sortedProcessNames[startIndex:endindex] {
		ret[idx] = metadata.ProcessInstanceNameIDs{
			ProcessName: name,
			ProcessIDs:  processNameIDs[name],
		}
	}

	ctx.RespEntityWithCount(int64(len(sortedProcessNames)), ret)
}

// ListProcessRelatedInfo list process related info according to condition
func (ps *ProcServer) ListProcessRelatedInfo(ctx *rest.Contexts) {

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("ListProcessRelatedInfo failed, parse bk_biz_id error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(metadata.ListProcessRelatedInfoOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	// get moduleIDs
	moduleIDs := input.Module.ModuleIDs
	if len(input.Set.SetIDs) > 0 {
		filter := map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: input.Set.SetIDs,
			},
		}
		if len(input.Module.ModuleIDs) > 0 {
			filter[common.BKModuleIDField] = map[string]interface{}{
				common.BKDBIN: input.Module.ModuleIDs,
			}
		}

		param := &metadata.QueryCondition{
			Condition: filter,
			Fields:    []string{common.BKModuleIDField},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
			},
		}

		moduleResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
			common.BKInnerObjIDModule, param)
		if nil != err {
			blog.Errorf("ListProcessRelatedInfo failed, coreservice http ReadInstance fail, param: %v, err: %v, "+
				"rid:%s", param, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
			return
		}

		if len(moduleResult.Info) == 0 {
			ctx.RespEntityWithCount(0, []interface{}{})
			return
		}

		mIDs := make([]int64, len(moduleResult.Info))
		for idx, info := range moduleResult.Info {
			mID, _ := info.Int64(common.BKModuleIDField)
			mIDs[idx] = mID
		}

		moduleIDs = mIDs
	}

	// get serviceIntanceIDs
	serviceIntanceIDs := input.ServiceInstance.IDs
	if len(input.ServiceInstance.IDs) > 0 || len(moduleIDs) > 0 {
		filter := map[string]interface{}{
			common.BKAppIDField: bizID,
		}

		if len(input.ServiceInstance.IDs) > 0 {
			filter[common.BKFieldID] = map[string]interface{}{
				common.BKDBIN: input.ServiceInstance.IDs,
			}
		}

		if len(moduleIDs) > 0 {
			filter[common.BKModuleIDField] = map[string]interface{}{
				common.BKDBIN: moduleIDs,
			}
		}

		option := &metadata.DistinctFieldOption{
			TableName: common.BKTableNameServiceInstance,
			Field:     common.BKFieldID,
			Filter:    filter,
		}

		sIDs, err := ps.CoreAPI.CoreService().Common().GetDistinctField(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			blog.Errorf("GetDistinctField failed, err:%s, option:%#v, rid:%s", err, *option, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if len(sIDs) == 0 {
			ctx.RespEntityWithCount(0, []interface{}{})
			return
		}

		srvInstIDs := make([]int64, len(sIDs))
		for idx, sID := range sIDs {
			if ID, err := strconv.ParseInt(fmt.Sprintf("%v", sID), 10, 64); err == nil {
				srvInstIDs[idx] = ID
			}
		}

		serviceIntanceIDs = srvInstIDs
	}

	// get processIDs
	var processIDs []int64
	if len(serviceIntanceIDs) > 0 {
		filter := map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKServiceInstanceIDField: map[string]interface{}{
				common.BKDBIN: serviceIntanceIDs,
			},
		}

		option := &metadata.DistinctFieldOption{
			TableName: common.BKTableNameProcessInstanceRelation,
			Field:     common.BKProcessIDField,
			Filter:    filter,
		}

		pIDs, err := ps.CoreAPI.CoreService().Common().GetDistinctField(ctx.Kit.Ctx, ctx.Kit.Header, option)
		if err != nil {
			blog.Errorf("GetDistinctField failed, err:%s, option:%#v, rid:%s", err, *option, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if len(pIDs) == 0 {
			ctx.RespEntityWithCount(0, []interface{}{})
			return
		}

		procIDs := make([]int64, len(pIDs))
		for idx, pID := range pIDs {
			if ID, err := strconv.ParseInt(fmt.Sprintf("%v", pID), 10, 64); err == nil {
				procIDs[idx] = ID
			}
		}

		processIDs = procIDs
	}

	// process detail
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
	}

	if len(processIDs) > 0 {
		filter[common.BKProcessIDField] = map[string]interface{}{
			common.BKDBIN: processIDs,
		}
	}

	propertyFilter := make(map[string]interface{})
	if input.ProcessPropertyFilter != nil {
		mgoFilter, key, err := input.ProcessPropertyFilter.ToMgo()
		if err != nil {
			blog.ErrorJSON("ListProcessRelatedInfo failed, ToMgo err:%s, ProcessPropertyFilter:%s, rid:%s", err,
				input.ProcessPropertyFilter, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()+fmt.Sprintf(", "+
				"host_property_filter.%s", key)))
			return
		}
		if len(mgoFilter) > 0 {
			propertyFilter = mgoFilter
		}
	}

	finalFilter := make(map[string]interface{})
	if len(propertyFilter) > 0 {
		finalFilter[common.BKDBAND] = []map[string]interface{}{filter, propertyFilter}
	} else {
		finalFilter = filter
	}

	fields := []string{}
	if len(input.Fields) > 0 {
		fields = input.Fields
		fields = append(fields, common.BKProcessIDField)
		fields = append(fields, common.BKProcessNameField)
		fields = append(fields, common.BKFuncIDField)
	}

	sort := input.Page.Sort
	if sort == "" {
		sort = common.BKProcessIDField
	}
	reqParam := &metadata.QueryCondition{
		Fields: fields,
		Page: metadata.BasePage{
			Sort:  sort,
			Limit: input.Page.Limit,
			Start: input.Page.Start,
		},
		Condition: finalFilter,
	}

	processResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, reqParam)
	if nil != err {
		blog.Errorf("ListProcessRelatedInfo failed, coreservice http ReadInstance fail, reqParam: %v, err: %v, "+
			"rid:%s", *reqParam, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	if len(processResult.Info) == 0 {
		ctx.RespEntityWithCount(0, []interface{}{})
		return
	}

	processIDsNeed := make([]int64, len(processResult.Info))
	processDetailMap := map[int64]interface{}{}
	for idx, process := range processResult.Info {
		processID, _ := process.Int64(common.BKProcessIDField)
		processIDsNeed[idx] = processID
		processDetailMap[processID] = process
	}

	ps.listProcessRelatedInfo(ctx, bizID, processIDsNeed, processDetailMap, int64(processResult.Count))
}

// listProcessRelatedInfo list process related info according to process info
func (ps *ProcServer) listProcessRelatedInfo(ctx *rest.Contexts, bizID int64, processIDs []int64,
	processDetailMap map[int64]interface{}, totalCnt int64) {

	// objID array
	srvinstArr := make([]int64, 0)
	hostArr := make([]int64, 0)
	moduleArr := make([]int64, 0)
	setArr := make([]int64, 0)

	// procID => objID map
	procSrvinstMap := make(map[int64]int64)
	procTemplateMap := make(map[int64]int64)
	procHostMap := make(map[int64]int64)
	srvinstModuleMap := make(map[int64]int64)
	moduleSetMap := make(map[int64]int64)

	// objID => objDetail map
	srvinstDetailMap := make(map[int64]metadata.ServiceInstanceDetailOfP)
	hostDetailMap := make(map[int64]metadata.HostDetailOfP)
	moduleDetailMap := make(map[int64]metadata.ModuleDetailOfP)
	setDetailMap := make(map[int64]metadata.SetDetailOfP)

	// get ID of serviceInstance, host, processTemplate and their process relation map
	listRelationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID: bizID,
		ProcessIDs: processIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, ccErr := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		listRelationOption)
	if ccErr != nil {
		ctx.RespWithError(ccErr, ccErr.GetCode(), "ListProcessInstanceRelation failed, option: %+v, err: %+v",
			listRelationOption, ccErr)
		return
	}

	for _, relation := range relations.Info {
		srvinstArr = append(srvinstArr, relation.ServiceInstanceID)
		hostArr = append(hostArr, relation.HostID)
		procSrvinstMap[relation.ProcessID] = relation.ServiceInstanceID
		procTemplateMap[relation.ProcessID] = relation.ProcessTemplateID
		procHostMap[relation.ProcessID] = relation.HostID
		procTemplateMap[relation.ProcessID] = relation.ProcessTemplateID
	}
	srvinstArr = util.IntArrayUnique(srvinstArr)

	// service instance detail
	instOpt := &metadata.ListServiceInstanceOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: srvinstArr,
	}
	instances, ccErr := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header, instOpt)
	if ccErr != nil {
		ctx.RespWithError(ccErr, ccErr.GetCode(), "ListServiceInstance failed, instOpt:%#v, err: %v", instOpt, ccErr)
		return
	}

	for _, inst := range instances.Info {
		srvinstDetailMap[inst.ID] = metadata.ServiceInstanceDetailOfP{
			ID:   inst.ID,
			Name: inst.Name,
		}
		srvinstModuleMap[inst.ID] = inst.ModuleID
		moduleArr = append(moduleArr, inst.ModuleID)
	}
	moduleArr = util.IntArrayUnique(moduleArr)

	// host detail
	hostParam := &metadata.QueryCondition{
		Fields: []string{common.BKHostIDField, common.BKCloudIDField, common.BKHostInnerIPField},
		Condition: map[string]interface{}{common.BKHostIDField: map[string]interface{}{
			common.BKDBIN: hostArr,
		},
		},
	}

	hostResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDHost, hostParam)
	if nil != err {
		blog.Errorf("ListProcessRelatedInfo failed, coreservice http ReadInstance fail, param: %v, err: %v, rid:%s",
			*hostParam, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	for _, host := range hostResult.Info {
		hostID, _ := host.Int64(common.BKHostIDField)
		cloudID, _ := host.Int64(common.BKCloudIDField)
		innerIP, _ := host.String(common.BKHostInnerIPField)
		hostDetailMap[hostID] = metadata.HostDetailOfP{
			HostID:  hostID,
			CloudID: cloudID,
			InnerIP: innerIP,
		}
	}

	// module detail
	moduleParam := &metadata.QueryCondition{
		Fields: []string{common.BKModuleIDField, common.BKModuleNameField, common.BKSetIDField},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKModuleIDField: map[string]interface{}{
				common.BKDBIN: moduleArr,
			},
		},
	}

	moduleResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDModule, moduleParam)
	if nil != err {
		blog.Errorf("ListProcessRelatedInfo failed, coreservice http ReadInstance fail, param: %v, err: %v, rid:%s",
			*moduleParam, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	for _, module := range moduleResult.Info {
		moduleID, _ := module.Int64(common.BKModuleIDField)
		moduleName, _ := module.String(common.BKModuleNameField)
		moduleDetailMap[moduleID] = metadata.ModuleDetailOfP{
			ModuleID:   moduleID,
			ModuleName: moduleName,
		}

		setID, _ := module.Int64(common.BKSetIDField)
		moduleSetMap[moduleID] = setID
		setArr = append(setArr, setID)

	}
	setArr = util.IntArrayUnique(setArr)

	// set detail
	setParam := &metadata.QueryCondition{
		Fields: []string{common.BKSetIDField, common.BKSetNameField, common.BKSetEnvField},
		Condition: map[string]interface{}{
			common.BKAppIDField: bizID,
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setArr,
			},
		},
	}

	setResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDSet, setParam)
	if nil != err {
		blog.Errorf("ListProcessRelatedInfo failed, coreservice http ReadInstance fail, param: %v, err: %v, rid:%s",
			*setParam, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	for _, set := range setResult.Info {
		setID, _ := set.Int64(common.BKSetIDField)
		setName, _ := set.String(common.BKSetNameField)
		setEnv, _ := set.String(common.BKSetEnvField)
		setDetailMap[setID] = metadata.SetDetailOfP{
			SetID:   setID,
			SetName: setName,
			SetEnv:  setEnv,
		}
	}

	// construct the final result
	ret := make([]metadata.ListProcessRelatedInfoResult, len(processIDs))

	for idx, processID := range processIDs {

		srvinstID := procSrvinstMap[processID]
		moduleID := srvinstModuleMap[srvinstID]
		setID := moduleSetMap[moduleID]

		hostDetail := hostDetailMap[procHostMap[processID]]
		srvinstDetail := srvinstDetailMap[srvinstID]
		moduleDetail := moduleDetailMap[moduleID]
		setDetail := setDetailMap[setID]

		info := metadata.ListProcessRelatedInfoResult{
			Set:             setDetail,
			Module:          moduleDetail,
			Host:            hostDetail,
			ServiceInstance: srvinstDetail,
			ProcessTemplate: metadata.ProcessTemplateDetailOfP{
				ID: procTemplateMap[processID],
			},
			Process: processDetailMap[processID],
		}
		ret[idx] = info
	}

	ctx.RespEntityWithCount(totalCnt, ret)
}

// ListProcessInstancesDetailsByIDs get process instances details and relation by their ids
func (ps *ProcServer) ListProcessInstancesDetailsByIDs(ctx *rest.Contexts) {
	input := new(metadata.ListProcessInstancesDetailsByIDsOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: input.ProcessIDs,
		},
	}
	reqParam := &metadata.QueryCondition{
		Condition: filter,
		Page:      input.Page,
	}
	processResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, reqParam)
	if nil != err {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "ListProcessInstancesDetailsByIDs failed, "+
			"reqParam: %#v, err: %+v", reqParam, err)
		return
	}

	processIDPropertyMap := map[int64]mapstr.MapStr{}
	sortedprocessIDs := make([]int64, 0)
	for _, process := range processResult.Info {
		processID, err := process.Int64(common.BKProcessIDField)
		if err != nil {
			ctx.RespWithError(err, common.CCErrCommParseDataFailed, "ListProcessInstancesDetailsByIDs failed, "+
				"process: %#v, err: %+v", process, err)
			return
		}
		processIDPropertyMap[processID] = process
		sortedprocessIDs = append(sortedprocessIDs, processID)
	}

	listRelationOption := &metadata.ListProcessInstanceRelationOption{
		BusinessID: input.BizID,
		ProcessIDs: sortedprocessIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	relations, err := ps.CoreAPI.CoreService().Process().ListProcessInstanceRelation(ctx.Kit.Ctx, ctx.Kit.Header,
		listRelationOption)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessInstanceRelationFailed,
			"ListProcessInstancesDetailsByIDs failed, list option: %+v, err: %+v", listRelationOption, err)
		return
	}

	processIDRelationMap := make(map[int64]metadata.ProcessInstanceRelation)
	serviceInstanceIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		processIDRelationMap[relation.ProcessID] = relation
		serviceInstanceIDs = append(serviceInstanceIDs, relation.ServiceInstanceID)
	}

	option := &metadata.ListServiceInstanceOption{
		BusinessID:         input.BizID,
		ServiceInstanceIDs: serviceInstanceIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstanceResult, err := ps.CoreAPI.CoreService().Process().ListServiceInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "ListProcessInstancesDetailsByIDs failed, "+
			"option: %#v, err: %v", option, err)
		return
	}
	serviceInstanceIDNames := make(map[int64]string)
	for _, instance := range serviceInstanceResult.Info {
		serviceInstanceIDNames[instance.ID] = instance.Name
	}

	processInstanceList := make([]metadata.ProcessInstanceDetailByID, 0)
	for _, id := range sortedprocessIDs {
		processDetail := metadata.ProcessInstanceDetailByID{
			ProcessID: id,
			Property:  processIDPropertyMap[id],
		}
		relation, exist := processIDRelationMap[id]
		if exist {
			processDetail.Relation = relation
			processDetail.ServiceInstanceName = serviceInstanceIDNames[relation.ServiceInstanceID]

		}
		processInstanceList = append(processInstanceList, processDetail)
	}

	ctx.RespEntityWithCount(int64(processResult.Count), processInstanceList)
}

// ListProcessInstancesDetails get process instances details by their ids
func (ps *ProcServer) ListProcessInstancesDetails(ctx *rest.Contexts) {

	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("ListProcessRelatedInfo failed, parse bk_biz_id error, err: %s, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_biz_id"))
		return
	}

	input := new(metadata.ListProcessInstancesDetailsOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := input.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: input.ProcessIDs,
		},
	}

	reqParam := &metadata.QueryCondition{
		Condition: filter,
		Fields:    input.Fields,
	}

	processResult, err := ps.CoreAPI.CoreService().Instance().ReadInstance(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDProc, reqParam)
	if nil != err {
		blog.Errorf("ListProcessInstancesDetails failed, coreservice http ReadInstance fail, reqParam: %v, err: %v, "+
			"rid:%s", *reqParam, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return
	}

	ctx.RespEntity(processResult.Info)
}

// UnbindServiceTemplateOnModuleEnable TODO
var UnbindServiceTemplateOnModuleEnable = true

// RemoveTemplateBindingOnModule TODO
func (ps *ProcServer) RemoveTemplateBindingOnModule(ctx *rest.Contexts) {
	if UnbindServiceTemplateOnModuleEnable {
		ctx.RespErrorCodeOnly(common.CCErrProcUnbindModuleServiceTemplateDisabled, "unbind service template from module disabled")
		return
	}

	input := new(metadata.RemoveTemplateBindingOnModuleOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	module, err := ps.getModule(ctx.Kit, input.ModuleID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrTopoGetModuleFailed, "create service instance failed, get module failed, moduleID: %d, err: %v", input.ModuleID, err)
		return
	}
	if module.BizID != input.BizID {
		err := ctx.Kit.CCError.CCError(common.CCErrCommNotFound)
		ctx.RespWithError(err, common.CCErrCommNotFound, "create service instance failed, get module failed, moduleID: %d, err: %v", input.ModuleID, err)
		return
	}

	var response *metadata.RemoveTemplateBoundOnModuleResult
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		response, err = ps.CoreAPI.CoreService().Process().RemoveTemplateBindingOnModule(ctx.Kit.Ctx, ctx.Kit.Header, input.ModuleID)
		if err != nil {
			blog.Errorf("remove template binding on module failed, parse business id failed, err: %+v", err)
			return ctx.Kit.CCError.CCError(common.CCErrProcRemoveTemplateBindingOnModule)
		}
		return nil
	})

	if txnErr != nil {
		blog.Errorf("RemoveTemplateBindingOnModule failed, err: %v, rid: %s", txnErr, ctx.Kit.Rid)
		return
	}
	ctx.RespEntity(response)
}
