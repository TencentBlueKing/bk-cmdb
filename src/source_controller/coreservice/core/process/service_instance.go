/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package process

import (
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (p *processOperation) CreateServiceInstance(kit *rest.Kit, instance metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := instance.Validate(); err != nil {
		blog.Errorf("CreateServiceInstance failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(kit, instance.BizID); err != nil {
		blog.Errorf("CreateServiceInstance failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	instance.BizID = bizID

	// validate module id field
	module, err := p.validateModuleID(kit, instance.ModuleID)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, module id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	if module.ServiceTemplateID != instance.ServiceTemplateID {
		blog.Errorf("CreateServiceInstance failed, module template id and instance template not equal, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCoreServiceModuleAndServiceInstanceTemplateNotCoincide)
	}

	// validate service template id field
	var serviceTemplate *metadata.ServiceTemplate
	if instance.ServiceTemplateID > 0 {
		st, err := p.GetServiceTemplate(kit, instance.ServiceTemplateID)
		if err != nil {
			blog.Errorf("CreateServiceInstance failed, service_template_id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		serviceTemplate = st
	}

	// validate host id field
	innerIP, err := p.validateHostID(kit, instance.HostID)
	if err != nil {
		blog.Errorf("CreateServiceInstance failed, host id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
	}

	// make sure biz id identical with service template
	if serviceTemplate != nil && serviceTemplate.BizID != bizID {
		blog.Errorf("CreateServiceInstance failed, validation failed, input bizID:%d not equal service template bizID:%d, rid: %s", bizID, serviceTemplate.BizID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// check unique `template_id + module_id + host_id`
	if instance.ServiceTemplateID != 0 {
		serviceInstanceFilter := map[string]interface{}{
			common.BKModuleIDField:          instance.ModuleID,
			common.BKHostIDField:            instance.HostID,
			common.BKServiceTemplateIDField: instance.ServiceTemplateID,
		}
		count, err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(kit.Ctx)
		if err != nil {
			blog.Errorf("CreateServiceInstance failed, list service instance failed, filter: %+v, err: %+v, rid: %s", serviceInstanceFilter, err, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if count > 0 {
			return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceInstanceAlreadyExist, innerIP)
		}
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(kit.Ctx, common.BKTableNameServiceInstance)
	if nil != err {
		blog.Errorf("CreateServiceInstance failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	instance.ID = int64(id)
	instance.Creator = kit.User
	instance.Modifier = kit.User
	instance.CreateTime = time.Now()
	instance.LastTime = time.Now()
	instance.SupplierAccount = kit.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Insert(kit.Ctx, &instance); nil != err {
		blog.Errorf("CreateServiceInstance failed, mongodb failed, table: %s, instance: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, instance, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	if instance.ServiceTemplateID != common.ServiceTemplateIDNotSet {
		listProcessTemplateOption := metadata.ListProcessTemplatesOption{
			BusinessID:         module.BizID,
			ServiceTemplateIDs: []int64{module.ServiceTemplateID},
			Page: metadata.BasePage{
				Limit: common.BKNoLimit,
				Sort:  "id",
			},
		}
		listProcTplResult, ccErr := p.ListProcessTemplates(kit, listProcessTemplateOption)
		if ccErr != nil {
			blog.Errorf("CreateServiceInstance failed, get process templates failed, listProcessTemplateOption: %+v, err: %+v, rid: %s", listProcessTemplateOption, ccErr, kit.Rid)
			return nil, ccErr
		}
		for _, processTemplate := range listProcTplResult.Info {
			processData := processTemplate.NewProcess(module.BizID, kit.SupplierAccount)
			process, ccErr := p.dependence.CreateProcessInstance(kit, processData)
			if ccErr != nil {
				blog.Errorf("CreateServiceInstance failed, create process instance failed, process: %+v, err: %+v, rid: %s", processData, ccErr, kit.Rid)
				return nil, ccErr
			}
			relation := &metadata.ProcessInstanceRelation{
				BizID:             bizID,
				ProcessID:         process.ProcessID,
				ServiceInstanceID: instance.ID,
				ProcessTemplateID: processTemplate.ID,
				HostID:            instance.HostID,
				SupplierAccount:   kit.SupplierAccount,
			}
			relation, ccErr = p.CreateProcessInstanceRelation(kit, relation)
			if ccErr != nil {
				blog.Errorf("CreateServiceInstance failed, create process relation failed, relation: %+v, err: %+v, rid: %s", relation, ccErr, kit.Rid)
				return nil, ccErr
			}
		}
	}

	if err := p.ReconstructServiceInstanceName(kit, instance.ID); err != nil {
		blog.Errorf("CreateServiceInstance failed, reconstruct instance name failed, instance: %+v, err: %s, rid: %s", instance, err.Error(), kit.Rid)
		return nil, err
	}

	// transfer host to target module
	transferConfig := &metadata.HostsModuleRelation{
		ApplicationID: bizID,
		HostID:        []int64{instance.HostID},
		ModuleID:      []int64{instance.ModuleID},
		IsIncrement:   true,
	}
	if _, err := p.dependence.TransferHostModuleDep(kit, transferConfig); err != nil {
		blog.Errorf("CreateServiceInstance failed, transfer host module failed, transfer: %+v, instance: %+v, err: %+v, rid: %s", transferConfig, instance, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostTransferModule)
	}

	return &instance, nil
}

func (p *processOperation) GetServiceInstance(kit *rest.Kit, instanceID int64) (*metadata.ServiceInstance, errors.CCErrorCoder) {
	instance := metadata.ServiceInstance{}

	filter := map[string]int64{common.BKFieldID: instanceID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).One(kit.Ctx, &instance); nil != err {
		blog.Errorf("GetServiceInstance failed, mongodb failed, table: %s, instance: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, instance, err, kit.Rid)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &instance, nil
}

func (p *processOperation) UpdateServiceInstance(kit *rest.Kit, instanceID int64, input metadata.ServiceInstance) (*metadata.ServiceInstance, errors.CCErrorCoder) {
	instance, err := p.GetServiceInstance(kit, instanceID)
	if err != nil {
		return nil, err
	}

	if field, err := input.Validate(); err != nil {
		blog.Errorf("UpdateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// update fields to original object
	instance.Name = input.Name

	// do update
	filter := map[string]int64{common.BKFieldID: instanceID}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Update(kit.Ctx, filter, instance); nil != err {
		blog.Errorf("UpdateServiceTemplate failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceInstance, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return instance, nil
}

func (p *processOperation) ListServiceInstance(kit *rest.Kit, option metadata.ListServiceInstanceOption) (*metadata.MultipleServiceInstance, errors.CCErrorCoder) {
	if option.BusinessID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	filter := map[string]interface{}{
		common.BKAppIDField:      option.BusinessID,
		common.BkSupplierAccount: kit.SupplierAccount,
	}
	filter = util.SetQueryOwner(filter, kit.SupplierAccount)

	if option.ServiceTemplateID != 0 {
		filter[common.BKServiceTemplateIDField] = option.ServiceTemplateID
	}

	// use param ServiceTemplateIDs when ServiceTemplateID is not set
	if len(option.ServiceTemplateIDs) > 0 && option.ServiceTemplateID == 0 {
		filter[common.BKServiceTemplateIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceTemplateIDs,
		}
	}

	if len(option.HostIDs) > 0 {
		filter[common.BKHostIDField] = map[string]interface{}{
			common.BKDBIN: option.HostIDs,
		}
	}

	if len(option.ModuleIDs) != 0 {
		filter[common.BKModuleIDField] = map[string]interface{}{
			common.BKDBIN: option.ModuleIDs,
		}
	}

	if option.ServiceInstanceIDs != nil {
		filter[common.BKFieldID] = map[string]interface{}{
			common.BKDBIN: option.ServiceInstanceIDs,
		}
	}

	if option.SearchKey != nil {
		filter[common.BKFieldName] = map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf(".*%s.*", *option.SearchKey),
		}
	}

	if key, err := option.Selectors.Validate(); err != nil {
		blog.Errorf("ListServiceInstance failed, selector validate failed, selectors: %+v, key: %s, err: %+v, rid: %s", option.Selectors, key, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	if len(option.Selectors) != 0 {
		labelFilter, err := option.Selectors.ToMgoFilter()
		if err != nil {
			blog.Errorf("ListServiceInstance failed, selectors to filer failed, selectors: %+v, err: %+v, rid: %s", option.Selectors, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "labels")
		}
		filter = util.MergeMaps(filter, labelFilter)
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Count(kit.Ctx); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	instances := make([]metadata.ServiceInstance, 0)
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Start(
		uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).All(kit.Ctx, &instances); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleServiceInstance{
		Count: total,
		Info:  instances,
	}
	return result, nil
}

func (p *processOperation) ListServiceInstanceDetail(kit *rest.Kit, option metadata.ListServiceInstanceDetailOption) (*metadata.MultipleServiceInstanceDetail, errors.CCErrorCoder) {
	if option.BusinessID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	if option.Page.Limit > common.BKMaxPageSize {
		return nil, kit.CCError.CCError(common.CCErrCommOverLimit)
	}

	moduleFilter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
	}
	if option.SetID != 0 {
		moduleFilter[common.BKSetIDField] = option.SetID
	}
	if option.ModuleID != 0 {
		moduleFilter[common.BKModuleIDField] = option.ModuleID
	}
	modules := make([]metadata.ModuleInst, 0)
	if err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(moduleFilter).All(kit.Ctx, &modules); err != nil {
		blog.Errorf("ListServiceInstanceDetail failed, list modules failed, filter: %+v, err: %+v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	targetModuleIDs := make([]int64, 0)
	moduleCategoryMap := make(map[int64]int64)
	for _, module := range modules {
		targetModuleIDs = append(targetModuleIDs, module.ModuleID)
		moduleCategoryMap[module.ModuleID] = module.ServiceCategoryID
	}

	if len(targetModuleIDs) == 0 {
		result := &metadata.MultipleServiceInstanceDetail{
			Count: 0,
			Info:  make([]metadata.ServiceInstanceDetail, 0),
		}
		return result, nil
	}

	filter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
		common.BKModuleIDField: map[string]interface{}{
			common.BKDBIN: targetModuleIDs,
		},
	}
	if option.HostID != 0 {
		filter[common.BKHostIDField] = option.HostID
	}

	if option.ServiceInstanceIDs != nil {
		filter[common.BKFieldID] = map[string]interface{}{
			common.BKDBIN: option.ServiceInstanceIDs,
		}
	}

	if key, err := option.Selectors.Validate(); err != nil {
		blog.Errorf("ListServiceInstance failed, selector validate failed, selectors: %+v, key: %s, err: %+v, rid: %s", option.Selectors, key, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	if len(option.Selectors) != 0 {
		labelFilter, err := option.Selectors.ToMgoFilter()
		if err != nil {
			blog.Errorf("ListServiceInstance failed, selectors to filer failed, selectors: %+v, err: %+v, rid: %s", option.Selectors, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "labels")
		}
		filter = util.MergeMaps(filter, labelFilter)
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Count(kit.Ctx); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	serviceInstances := make([]metadata.ServiceInstance, 0)
	serviceInstanceDetails := make([]metadata.ServiceInstanceDetail, 0)
	start := uint64(option.Page.Start)
	limit := uint64(option.Page.Limit)
	query := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(filter).Start(start).Limit(limit)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	}
	if err := query.All(kit.Ctx, &serviceInstances); nil != err {
		blog.Errorf("ListServiceInstance failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	for _, serviceInstance := range serviceInstances {
		serviceInstanceDetails = append(serviceInstanceDetails, metadata.ServiceInstanceDetail{
			ServiceInstance: serviceInstance,
		})
	}

	if len(serviceInstances) == 0 {
		result := &metadata.MultipleServiceInstanceDetail{
			Count: total,
			Info:  serviceInstanceDetails,
		}
		return result, nil
	}

	// filter process instances
	serviceInstanceIDs := make([]int64, 0)
	for idx, serviceInstance := range serviceInstanceDetails {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
		// set service_category_id field
		serviceInstanceDetails[idx].ServiceCategoryID = moduleCategoryMap[serviceInstance.ModuleID]
	}

	relations := make([]metadata.ProcessInstanceRelation, 0)
	relationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: map[string]interface{}{
			common.BKDBIN: serviceInstanceIDs,
		},
	}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).All(kit.Ctx, &relations); err != nil {
		blog.Errorf("ListServiceInstanceDetail failed, list processRelations failed, err: %+v, rid: %s", relationFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	processIDs := make([]int64, 0)
	for _, relation := range relations {
		processIDs = append(processIDs, relation.ProcessID)
	}
	processes := make([]metadata.Process, 0)
	processFilter := map[string]interface{}{
		common.BKProcessIDField: map[string]interface{}{
			common.BKDBIN: processIDs,
		},
	}
	if err := p.dbProxy.Table(common.BKTableNameBaseProcess).Find(processFilter).All(kit.Ctx, &processes); err != nil {
		blog.Errorf("ListServiceInstanceDetail failed, list process failed, filter: %+v, err: %s, rid: %s", processFilter, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}
	// processID -> relation
	processRelationMap := make(map[int64]metadata.ProcessInstanceRelation)
	for _, relation := range relations {
		processRelationMap[relation.ProcessID] = relation
	}
	// serviceInstanceID -> []ProcessInstance
	serviceInstanceMap := make(map[int64][]metadata.ProcessInstanceNG)
	for _, process := range processes {
		relation, ok := processRelationMap[process.ProcessID]
		if !ok {
			blog.Warnf("ListServiceInstanceDetail got unexpected state, process's relation not found, process: %+v, rid: %s", process, kit.Rid)
			continue
		}
		if _, ok := serviceInstanceMap[relation.ServiceInstanceID]; !ok {
			serviceInstanceMap[relation.ServiceInstanceID] = make([]metadata.ProcessInstanceNG, 0)
		}
		processInstance := metadata.ProcessInstanceNG{
			Process:  process,
			Relation: relation,
		}
		serviceInstanceMap[relation.ServiceInstanceID] = append(serviceInstanceMap[relation.ServiceInstanceID], processInstance)
	}

	for idx, serviceInstance := range serviceInstanceDetails {
		processInfo, ok := serviceInstanceMap[serviceInstance.ID]
		if !ok {
			continue
		}
		serviceInstanceDetails[idx].ProcessInstances = processInfo
	}

	result := &metadata.MultipleServiceInstanceDetail{
		Count: total,
		Info:  serviceInstanceDetails,
	}
	return result, nil
}

func (p *processOperation) DeleteServiceInstance(kit *rest.Kit, serviceInstanceIDs []int64) errors.CCErrorCoder {
	for _, serviceInstanceID := range serviceInstanceIDs {
		instance, err := p.GetServiceInstance(kit, serviceInstanceID)
		if err != nil {
			blog.Errorf("DeleteServiceInstance failed, GetServiceInstance failed, instanceID: %d, err: %+v, rid: %s", serviceInstanceID, err, kit.Rid)
			return err
		}

		// service template that referenced by process template shouldn't be removed
		usageFilter := map[string]int64{common.BKServiceInstanceIDField: instance.ID}
		usageCount, e := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(usageFilter).Count(kit.Ctx)
		if nil != e {
			blog.Errorf("DeleteServiceInstance failed, mongodb failed, table: %s, usageFilter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, usageFilter, e, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		if usageCount > 0 {
			blog.Errorf("DeleteServiceInstance failed, forbidden delete service instance be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, kit.Rid)
			err := kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
			return err
		}

		serviceInstanceFilter := map[string]int64{common.BKFieldID: instance.ID}
		if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Delete(kit.Ctx, serviceInstanceFilter); nil != err {
			blog.Errorf("DeleteServiceInstance failed, mongodb failed, table: %s, deleteFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, serviceInstanceFilter, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
		}
	}
	return nil
}

// GetServiceInstanceName get service instance's name, format: `IP + first process name + first process port`
// 可能应用场景：1. 查询服务实例时组装名称；2. 更新进程信息时根据组装名称直接更新到 `name` 字段
// issue: https://github.com/Tencent/bk-cmdb/issues/2485
func (p *processOperation) generateServiceInstanceName(kit *rest.Kit, instanceID int64) (string, errors.CCErrorCoder) {

	// get instance
	instance := metadata.ServiceInstance{}
	instanceFilter := map[string]interface{}{
		common.BKFieldID: instanceID,
	}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(instanceFilter).One(kit.Ctx, &instance); err != nil {
		blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, instanceFilter, err, kit.Rid)
		if p.dbProxy.IsNotFoundError(err) {
			return "", kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}
		return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	// get host inner ip
	host := struct {
		InnerIP string `json:"bk_host_innerip" bson:"bk_host_innerip"`
		HostID  int    `json:"bk_host_id" bson:"bk_host_id"`
	}{}

	hostFilter := map[string]interface{}{
		common.BKHostIDField: instance.HostID,
	}
	if err := p.dbProxy.Table(common.BKTableNameBaseHost).Find(hostFilter).One(kit.Ctx, &host); err != nil {
		blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameBaseHost, hostFilter, err, kit.Rid)
		if p.dbProxy.IsNotFoundError(err) {
			return "", kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}
		return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	instanceName := host.InnerIP

	// get first process instance relation
	relation := metadata.ProcessInstanceRelation{}
	relationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: instance.ID,
	}
	order := "id"
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).Sort(order).One(kit.Ctx, &relation); err != nil {
		// relation not found means no process in service instance, service instance's name will only contains ip in that case
		if !p.dbProxy.IsNotFoundError(err) {
			blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, relationFilter, err, kit.Rid)
			return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
	}

	if relation.ProcessID != 0 {
		// get process instance
		process := metadata.Process{}
		processFilter := map[string]interface{}{
			common.BKProcIDField: relation.ProcessID,
		}
		if err := p.dbProxy.Table(common.BKTableNameBaseProcess).Find(processFilter).One(kit.Ctx, &process); err != nil {
			blog.Errorf("GetServiceInstanceName failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameBaseProcess, processFilter, err, kit.Rid)
			if p.dbProxy.IsNotFoundError(err) {
				return "", kit.CCError.CCErrorf(common.CCErrCommNotFound)
			}
			return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		if process.ProcessName != nil && len(*process.ProcessName) > 0 {
			instanceName += fmt.Sprintf("_%s", *process.ProcessName)
		}
		if process.Port != nil && len(*process.Port) > 0 {
			instanceName += fmt.Sprintf("_%s", *process.Port)
		}
	}
	return instanceName, nil
}

// ReconstructServiceInstanceName do reconstruct service instance name after process name or process port changed
func (p *processOperation) ReconstructServiceInstanceName(kit *rest.Kit, instanceID int64) errors.CCErrorCoder {
	name, err := p.generateServiceInstanceName(kit, instanceID)
	if err != nil {
		blog.Errorf("ReconstructServiceInstanceName failed, generate instance name failed, err: %s, rid: %s", err.Error(), kit.Rid)
		return err
	}
	filter := map[string]interface{}{
		common.BKFieldID: instanceID,
	}
	doc := map[string]interface{}{
		common.BKFieldName: name,
	}
	e := p.dbProxy.Table(common.BKTableNameServiceInstance).Update(kit.Ctx, filter, doc)
	if e != nil {
		blog.Errorf("ReconstructServiceInstanceName failed, update instance name failed, err: %+v, rid: %s", e, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

// GetDefaultModuleIDs get business's default module id, default module type specified by DefaultResModuleFlag
// be careful: it doesn't ensure business have all default module or set
func (p *processOperation) GetBusinessDefaultSetModuleInfo(kit *rest.Kit, bizID int64) (metadata.BusinessDefaultSetModuleInfo, errors.CCErrorCoder) {
	defaultSetModuleInfo := metadata.BusinessDefaultSetModuleInfo{}

	// find and fill default module
	defaultModuleCond := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKDefaultField: map[string]interface{}{
			common.BKDBNE: common.DefaultFlagDefaultValue,
		},
	}
	modules := make([]struct {
		ModuleID   int64 `bson:"bk_module_id"`
		ModuleFlag int   `bson:"default"`
	}, 0)
	err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(defaultModuleCond).Fields(common.BKModuleIDField, common.BKDefaultField).All(kit.Ctx, &modules)
	if nil != err {
		blog.Errorf("get default module failed, err: %+v, filter: %+v, rid: %s", err, defaultModuleCond, kit.Rid)
		return defaultSetModuleInfo, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	for _, module := range modules {
		if module.ModuleFlag == common.DefaultResModuleFlag {
			defaultSetModuleInfo.IdleModuleID = module.ModuleID
		}
		if module.ModuleFlag == common.DefaultFaultModuleFlag {
			defaultSetModuleInfo.FaultModuleID = module.ModuleID
		}
		if module.ModuleFlag == common.DefaultRecycleModuleFlag {
			defaultSetModuleInfo.RecycleModuleID = module.ModuleID
		}
	}

	// find and fill default set
	defaultSetCond := map[string]interface{}{
		common.BKAppIDField:   bizID,
		common.BKDefaultField: common.DefaultResSetFlag,
	}
	sets := make([]struct {
		SetID int64 `bson:"bk_set_id"`
	}, 0)
	err = p.dbProxy.Table(common.BKTableNameBaseSet).Find(defaultSetCond).Fields(common.BKSetIDField).All(kit.Ctx, &sets)
	if nil != err {
		blog.Errorf("get default set failed, err: %+v, filter: %+v, rid: %s", err, defaultSetCond, kit.Rid)
		return defaultSetModuleInfo, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	for _, set := range sets {
		defaultSetModuleInfo.IdleSetID = set.SetID
	}

	return defaultSetModuleInfo, nil
}

// AutoCreateServiceInstanceModuleHost create service instance on host under module base on service template
func (p *processOperation) AutoCreateServiceInstanceModuleHost(kit *rest.Kit, hostID int64, moduleID int64) (*metadata.ServiceInstance, errors.CCErrorCoder) {
	moduleFilter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	module := struct {
		BizID             int64  `bson:"bk_biz_id"`
		ModuleID          int64  `bson:"bk_module_id"`
		ModuleName        string `bson:"bk_module_name"`
		SupplierAccount   string `bson:"bk_supplier_account"`
		ServiceTemplateID int64  `bson:"service_template_id"`
		ServiceCategoryID int64  `bson:"service_category_id"`
	}{}
	var err error
	if err = p.dbProxy.Table(common.BKTableNameBaseModule).Find(moduleFilter).One(kit.Ctx, &module); err != nil {
		blog.ErrorJSON("AutoCreateServiceInstanceModuleHost failed, get module failed, err: %+v, cond: %#v, rid: %s", err, moduleFilter, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if module.ServiceTemplateID == common.ServiceTemplateIDNotSet {
		blog.Infof("AutoCreateServiceInstanceModuleHost do nothing, ServiceTemplateID is %d, rid: %s", common.ServiceTemplateIDNotSet, kit.Rid)
		return nil, nil
	}

	now := time.Now()
	serviceInstanceData := &metadata.ServiceInstance{
		BizID:             module.BizID,
		ServiceTemplateID: module.ServiceTemplateID,
		HostID:            hostID,
		ModuleID:          moduleID,
		Creator:           kit.User,
		Modifier:          kit.User,
		CreateTime:        now,
		LastTime:          now,
		SupplierAccount:   kit.SupplierAccount,
	}
	var ccErr errors.CCErrorCoder
	serviceInstance, ccErr := p.CreateServiceInstance(kit, *serviceInstanceData)
	if ccErr != nil {
		if ccErr.GetCode() == common.CCErrCoreServiceInstanceAlreadyExist {
			serviceInstanceFilter := map[string]interface{}{
				common.BKModuleIDField:          serviceInstanceData.ModuleID,
				common.BKHostIDField:            serviceInstanceData.HostID,
				common.BKServiceTemplateIDField: serviceInstanceData.ServiceTemplateID,
			}
			serviceInstance = &metadata.ServiceInstance{}
			if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).One(kit.Ctx, serviceInstance); err != nil {
				blog.Errorf("AutoCreateServiceInstanceModuleHost failed, get exist service instance failed, serviceInstanceData: %+v, err: %+v, rid: %s", serviceInstanceData, ccErr, kit.Rid)
				return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
			}
			return serviceInstance, nil
		} else {
			blog.Errorf("AutoCreateServiceInstanceModuleHost failed, create service instance failed, serviceInstance: %+v, err: %+v, rid: %s", serviceInstance, ccErr, kit.Rid)
			return nil, ccErr
		}
	}

	return serviceInstance, nil
}

func (p *processOperation) RemoveTemplateBindingOnModule(kit *rest.Kit, moduleID int64) errors.CCErrorCoder {
	moduleFilter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	moduleSimple := struct {
		ServiceTemplateID int64 `field:"service_template_id" bson:"service_template_id" json:"service_template_id"`
		ServiceCategoryID int64 `field:"service_category_id" bson:"service_category_id" json:"service_category_id"`
		BizID             int64 `field:"bk_biz_id" bson:"bk_biz_id" json:"bk_biz_id"`
	}{}
	if err := p.dbProxy.Table(common.BKTableNameBaseModule).Find(moduleFilter).One(kit.Ctx, &moduleSimple); err != nil {
		blog.Errorf("RemoveTemplateBindingOnModule failed, get module by id failed, moduleID: %d, err: %+v, rid: %s", moduleID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if moduleSimple.ServiceTemplateID == 0 {
		return kit.CCError.CCError(common.CCErrCoreServiceModuleNotBoundWithTemplate)
	}

	// clear template id field on module
	resetServiceTemplateIDOption := map[string]interface{}{
		common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
	}
	if err := p.dbProxy.Table(common.BKTableNameBaseModule).Update(kit.Ctx, moduleFilter, resetServiceTemplateIDOption); err != nil {
		blog.Errorf("remove template binding on module failed, reset service_template_id on module failed, module: %d, err: %+v, rid: %s", moduleID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	// clear service instance template
	serviceInstanceFilter := map[string]int64{
		common.BKModuleIDField: moduleID,
	}
	if err := p.dbProxy.Table(common.BKTableNameServiceInstance).Update(kit.Ctx, serviceInstanceFilter, resetServiceTemplateIDOption); err != nil {
		blog.Errorf("remove template binding on module failed, reset service_template_id on service instance failed, module: %d, err: %+v, rid: %s", moduleID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	listOption := metadata.ListServiceInstanceOption{
		BusinessID:         moduleSimple.BizID,
		ModuleIDs:          []int64{moduleID},
		SearchKey:          nil,
		ServiceInstanceIDs: nil,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceInstanceResult, err := p.ListServiceInstance(kit, listOption)
	if err != nil {
		blog.Errorf("ListServiceInstance failed, option: %+v, err: %s, rid: %s", listOption, err.Error(), kit.Rid)
		return err
	}
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstanceResult.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	// clear process template id on relation
	processInstanceRelationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: map[string]interface{}{
			common.BKDBIN: serviceInstanceIDs,
		},
	}
	resetProcessTemplateIDOption := map[string]int64{
		common.BKProcessTemplateIDField: common.ServiceTemplateIDNotSet,
	}
	if err := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Update(kit.Ctx, processInstanceRelationFilter, resetProcessTemplateIDOption); err != nil {
		blog.Errorf("remove template binding on module failed, reset service_template_id on process instance relation failed, module: %d, err: %+v, rid: %s", moduleID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}
