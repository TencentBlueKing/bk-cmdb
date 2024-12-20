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
	"sync"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateServiceInstance create service instance
func (p *processOperation) CreateServiceInstance(kit *rest.Kit, instance *metadata.ServiceInstance) (
	*metadata.ServiceInstance, errors.CCErrorCoder) {

	if err := p.validateCreateSvcInstData(kit, instance); err != nil {
		return nil, err
	}

	// generate id field
	id, err := mongodb.Shard(kit.SysShardOpts()).NextSequence(kit.Ctx, common.BKTableNameServiceInstance)
	if err != nil {
		blog.Errorf("generate id failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	instance.ID = int64(id)
	instance.Creator = kit.User
	instance.Modifier = kit.User
	instance.CreateTime = time.Now()
	instance.LastTime = time.Now()
	instance.TenantID = kit.TenantID

	if err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Insert(kit.Ctx,
		&instance); err != nil {
		blog.Errorf("create service instance(%+v) failed, err: %v, rid: %s", instance, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	// get host data for instance name and bind IP
	host := metadata.HostMapStr{}
	filter := map[string]interface{}{common.BKHostIDField: instance.HostID}
	fields := []string{common.BKHostInnerIPField, common.BKHostOuterIPField, common.BKHostInnerIPv6Field,
		common.BKHostOuterIPv6Field}
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(filter).Fields(fields...).
		One(kit.Ctx, &host)
	if err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	firstTemplateProcess, ccErr := p.createSvcInstProcesses(kit, instance, host)
	if ccErr != nil {
		return nil, ccErr
	}

	if instance.Name == "" {
		name, err := p.ConstructServiceInstanceName(kit, instance.ID, host, firstTemplateProcess)
		if err != nil {
			blog.Errorf("construct service instance(%#v) name failed, err: %v, rid: %s", err, instance, kit.Rid)
			return nil, err
		}
		instance.Name = name
	}

	return instance, nil
}

// createSvcInstProcesses create processes in service instance by service template
func (p *processOperation) createSvcInstProcesses(kit *rest.Kit, instance *metadata.ServiceInstance,
	host metadata.HostMapStr) (*metadata.Process, errors.CCErrorCoder) {

	if instance.ServiceTemplateID == common.ServiceTemplateIDNotSet {
		return nil, nil
	}

	listProccTempOption := metadata.ListProcessTemplatesOption{
		BusinessID:         instance.BizID,
		ServiceTemplateIDs: []int64{instance.ServiceTemplateID},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
			Sort:  "id",
		},
	}
	listProcTplResult, ccErr := p.ListProcessTemplates(kit, listProccTempOption)
	if ccErr != nil {
		blog.Errorf("get process templates failed, opt: %+v, err: %v, rid: %s", listProccTempOption, ccErr, kit.Rid)
		return nil, ccErr
	}

	if len(listProcTplResult.Info) == 0 {
		blog.Errorf("create service instance for service template with no process is not allowed")
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
	}

	processes := make([]*metadata.Process, len(listProcTplResult.Info))
	relations := make([]*metadata.ProcessInstanceRelation, len(listProcTplResult.Info))
	templateIDs := make([]int64, len(listProcTplResult.Info))
	for idx, processTemplate := range listProcTplResult.Info {
		processData, err := processTemplate.NewProcess(kit.CCError, instance.BizID, instance.ID,
			kit.TenantID, host)
		if err != nil {
			blog.ErrorJSON("generate process instance by template %s failed, err: %s, rid: %s", processTemplate, err,
				kit.Rid)
			return nil, errors.New(common.CCErrCommParamsInvalid, err.Error())
		}
		processes[idx] = processData
		templateIDs[idx] = processTemplate.ID
	}

	processes, ccErr = p.dependence.CreateProcessInstances(kit, processes)
	if ccErr != nil {
		blog.Errorf("create process instances failed, processes: %#v, err: %v, rid: %s", processes, ccErr, kit.Rid)
		return nil, ccErr
	}

	for idx, process := range processes {
		relation := &metadata.ProcessInstanceRelation{
			BizID:             instance.BizID,
			ProcessID:         process.ProcessID,
			ServiceInstanceID: instance.ID,
			ProcessTemplateID: templateIDs[idx],
			HostID:            instance.HostID,
			TenantID:          kit.TenantID,
		}
		relations[idx] = relation
	}
	relations, ccErr = p.CreateProcessInstanceRelations(kit, relations)
	if ccErr != nil {
		blog.Errorf("create process relation failed, relations: %#v, err: %v, rid: %s", relations, ccErr, kit.Rid)
		return nil, ccErr
	}

	return processes[0], nil
}

func (p *processOperation) validateCreateSvcInstData(kit *rest.Kit,
	instance *metadata.ServiceInstance) errors.CCErrorCoder {

	// base attribute validate
	if field, err := instance.Validate(); err != nil {
		blog.Errorf("validate service instance(%+v) failed, err: %v, rid: %s", instance, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(kit, instance.BizID); err != nil {
		blog.Errorf("validate biz id(%d) failed, err: %v, rid: %s", instance.BizID, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	instance.BizID = bizID

	// validate module id field
	module, err := p.validateModuleID(kit, instance.ModuleID)
	if err != nil {
		blog.Errorf("module id %d is invalid, err: %v, rid: %s", instance.ModuleID, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}

	if module.ServiceTemplateID != instance.ServiceTemplateID {
		blog.Errorf("module template id %d and instance template %d not equal, err: %v, rid: %s",
			module.ServiceTemplateID, instance.ServiceTemplateID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCoreServiceModuleAndServiceInstanceTemplateNotCoincide)
	}

	// validate service template id field
	var serviceTemplate *metadata.ServiceTemplate
	if instance.ServiceTemplateID > 0 {
		st, err := p.GetServiceTemplate(kit, instance.ServiceTemplateID)
		if err != nil {
			blog.Errorf("service template id %d is invalid, err: %v, rid: %s", instance.ServiceTemplateID, err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField)
		}
		serviceTemplate = st
	}

	// validate host id field
	innerIP, err := p.validateHostID(kit, instance.HostID)
	if err != nil {
		blog.Errorf("host id %d is invalid, err: %v, rid: %s", instance.HostID, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
	}

	// make sure biz id identical with service template
	if serviceTemplate != nil && serviceTemplate.BizID != bizID {
		blog.Errorf("input biz id %d != service template biz id %d, rid: %s", bizID, serviceTemplate.BizID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_biz_id")
	}

	// check unique `template_id + module_id + host_id`
	if instance.ServiceTemplateID != 0 {
		filter := map[string]interface{}{
			common.BKModuleIDField:          instance.ModuleID,
			common.BKHostIDField:            instance.HostID,
			common.BKServiceTemplateIDField: instance.ServiceTemplateID,
		}
		count, err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(filter).
			Count(kit.Ctx)
		if err != nil {
			blog.Errorf("list service instance failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
			return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
		}
		if count > 0 {
			return kit.CCError.CCErrorf(common.CCErrCoreServiceInstanceAlreadyExist, innerIP)
		}
	}
	return nil
}

// CreateServiceInstances TODO
func (p *processOperation) CreateServiceInstances(kit *rest.Kit,
	instances []*metadata.ServiceInstance) ([]*metadata.ServiceInstance, errors.CCErrorCoder) {
	var wg sync.WaitGroup
	var lock sync.RWMutex
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	insts := make([]*metadata.ServiceInstance, len(instances))

	for idx := range instances {
		pipeline <- true
		wg.Add(1)

		go func(idx int, instance *metadata.ServiceInstance) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			inst, err := p.CreateServiceInstance(kit, instance)
			if err != nil {
				blog.ErrorJSON("CreateServiceInstance failed, idx: %s, instance: %s, err: %s, rid: %s", idx, instance,
					err, kit.Rid)
				if firstErr == nil {
					firstErr = err
				}
				return
			}

			lock.Lock()
			insts[idx] = inst
			lock.Unlock()

		}(idx, instances[idx])
	}

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	return insts, nil
}

// GetServiceInstance TODO
func (p *processOperation) GetServiceInstance(kit *rest.Kit, instanceID int64) (*metadata.ServiceInstance,
	errors.CCErrorCoder) {
	instance := metadata.ServiceInstance{}

	filter := map[string]int64{common.BKFieldID: instanceID}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(filter).One(kit.Ctx,
		&instance); err != nil {
		blog.Errorf("find service instance failed, instance: %v, err: %v, rid: %s", instance, err, kit.Rid)
		if mongodb.IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &instance, nil
}

// UpdateServiceInstances TODO
func (p *processOperation) UpdateServiceInstances(kit *rest.Kit, bizID int64,
	option *metadata.UpdateServiceInstanceOption) errors.CCErrorCoder {
	for _, data := range option.Data {
		needUpdate := data.Update
		needUpdate[common.LastTimeField] = time.Now()
		needUpdate[common.ModifierField] = kit.User

		filter := map[string]int64{
			common.BKAppIDField: bizID,
			common.BKFieldID:    data.ServiceInstanceID,
		}
		if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Update(kit.Ctx, filter,
			needUpdate); err != nil {
			blog.Errorf("update service instance failed, err: %v, filter: %v, needUpdate: %v, rid: %s", err, filter,
				needUpdate, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
		}
	}

	return nil
}

// ListServiceInstance TODO
func (p *processOperation) ListServiceInstance(kit *rest.Kit,
	option metadata.ListServiceInstanceOption) (*metadata.MultipleServiceInstance, errors.CCErrorCoder) {
	if option.BusinessID == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	filter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
	}

	if option.ServiceTemplateID != 0 {
		filter[common.BKServiceTemplateIDField] = option.ServiceTemplateID
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
		blog.Errorf("selector validate failed, selectors: %v, key: %s, err: %v, rid: %s", option.Selectors, key,
			err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	if len(option.Selectors) != 0 {
		labelFilter, err := option.Selectors.ToMgoFilter()
		if err != nil {
			blog.Errorf("selectors to filer failed, selectors: %+v, err: %v, rid: %s", option.Selectors, err,
				kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "labels")
		}
		filter = util.MergeMaps(filter, labelFilter)
	}

	var total uint64
	var err error
	if total, err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(filter).Count(kit.
		Ctx); err != nil {
		blog.Errorf("find service instance failed, mongodb failed, table: %s, filter: %+v, err: %v, rid: %s",
			common.BKTableNameServiceInstance, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	result := new(metadata.MultipleServiceInstance)

	instances := make([]metadata.ServiceInstance, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(filter).
		Fields(option.Fields...).Sort(option.Page.Sort).Start(uint64(option.Page.Start)).
		Limit(uint64(option.Page.Limit)).All(kit.Ctx, &instances); err != nil {
		blog.Errorf("find service instance failed, table: %s, filter: %+v, err: %v, rid: %s",
			common.BKTableNameServiceInstance, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result = &metadata.MultipleServiceInstance{
		Count: total,
		Info:  instances,
	}
	return result, nil
}

// ListServiceInstanceDetail list service instance detail
func (p *processOperation) ListServiceInstanceDetail(kit *rest.Kit, option metadata.ListServiceInstanceDetailOption) (
	*metadata.MultipleServiceInstanceDetail, errors.CCErrorCoder) {

	total, serviceInstanceDetails, err := p.listSvcInstForDetail(kit, option)
	if err != nil {
		return nil, err
	}

	if len(serviceInstanceDetails) == 0 {
		return &metadata.MultipleServiceInstanceDetail{Count: total, Info: serviceInstanceDetails}, nil
	}

	// filter process instances
	serviceInstanceIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstanceDetails {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
	}

	relations := make([]metadata.ProcessInstanceRelation, 0)
	relationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: map[string]interface{}{common.BKDBIN: serviceInstanceIDs},
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).
		All(kit.Ctx, &relations); err != nil {
		blog.Errorf("list process relations failed, filter: %v, err: %v, rid: %s", relationFilter, err, kit.Rid)
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
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseProcess).Find(processFilter).All(kit.Ctx,
		&processes); err != nil {
		blog.Errorf("list process failed, filter: %+v, err: %v, rid: %s", processFilter, err, kit.Rid)
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
			blog.Warnf("process's relation not found, process: %+v, rid: %s", process, kit.Rid)
			continue
		}
		if _, ok := serviceInstanceMap[relation.ServiceInstanceID]; !ok {
			serviceInstanceMap[relation.ServiceInstanceID] = make([]metadata.ProcessInstanceNG, 0)
		}
		processInstance := metadata.ProcessInstanceNG{
			Process:  process,
			Relation: relation,
		}
		serviceInstanceMap[relation.ServiceInstanceID] = append(serviceInstanceMap[relation.ServiceInstanceID],
			processInstance)
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

func (p *processOperation) listSvcInstForDetail(kit *rest.Kit, option metadata.ListServiceInstanceDetailOption) (uint64,
	[]metadata.ServiceInstanceDetail, errors.CCErrorCoder) {

	if option.BusinessID <= 0 {
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	if option.Page.IsIllegal() {
		return 0, nil, kit.CCError.CCError(common.CCErrCommPageLimitIsExceeded)
	}

	// set query params
	filter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
	}
	if option.ModuleID > 0 {
		filter[common.BKModuleIDField] = option.ModuleID
	}

	if option.HostID > 0 && len(option.HostList) > 0 {
		blog.Errorf("list service instance failed, parameters bk_host_id and bk_host_list cannot be set at the "+
			"same time, rid: %s", kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_host_id and bk_host_list cannot be "+
			"set at the same time")
	}

	if option.HostID > 0 {
		filter[common.BKHostIDField] = option.HostID
	}

	// Only one parameter between bk_host_list and bk_host_id can take effect,bk_host_id is not recommend to use.
	if len(option.HostList) > 0 {
		filter[common.BKHostIDField] = map[string]interface{}{
			common.BKDBIN: option.HostList,
		}
	}

	if option.ServiceInstanceIDs != nil {
		filter[common.BKFieldID] = map[string]interface{}{
			common.BKDBIN: option.ServiceInstanceIDs,
		}
	}
	if key, err := option.Selectors.Validate(); err != nil {
		blog.Errorf("selector validate failed, selectors: %v, key: %s, err: %v, rid: %s", option.Selectors, key,
			err, kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	if len(option.Selectors) != 0 {
		labelFilter, err := option.Selectors.ToMgoFilter()
		if err != nil {
			blog.Errorf("selectors to filer failed, selectors: %v, err: %v, rid: %s", option.Selectors, err, kit.Rid)
			return 0, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "labels")
		}
		filter = util.MergeMaps(filter, labelFilter)
	}

	var total uint64
	var err error
	if total, err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(filter).
		Count(kit.Ctx); err != nil {
		blog.Errorf("find service instance failed, filter: %v, err: %v, rid: %s", filter, err, kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	serviceInstances := make([]metadata.ServiceInstance, 0)
	start := uint64(option.Page.Start)
	limit := uint64(option.Page.Limit)
	query := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(filter).Start(start).
		Limit(limit)
	if len(option.Page.Sort) > 0 {
		query = query.Sort(option.Page.Sort)
	} else {
		query = query.Sort(common.BKFieldID)
	}
	if err := query.All(kit.Ctx, &serviceInstances); err != nil {
		blog.Errorf("list service instance failed, filter: %+v, err: %v, rid: %s", filter, err, kit.Rid)
		return 0, nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	serviceInstanceDetails := make([]metadata.ServiceInstanceDetail, 0, len(serviceInstances))
	for _, serviceInstance := range serviceInstances {
		serviceInstanceDetails = append(serviceInstanceDetails, metadata.ServiceInstanceDetail{
			ServiceInstance: serviceInstance,
		})
	}
	return total, serviceInstanceDetails, nil
}

// DeleteServiceInstance TODO
func (p *processOperation) DeleteServiceInstance(kit *rest.Kit, serviceInstanceIDs []int64) errors.CCErrorCoder {
	for _, serviceInstanceID := range serviceInstanceIDs {
		instance, err := p.GetServiceInstance(kit, serviceInstanceID)
		if err != nil {
			blog.Errorf("get service instance failed, instanceID: %d, err: %v, rid: %s", serviceInstanceID, err,
				kit.Rid)
			return err
		}

		// service template that referenced by process template shouldn't be removed
		usageFilter := map[string]int64{common.BKServiceInstanceIDField: instance.ID}
		usageCount, e := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).
			Find(usageFilter).Count(kit.Ctx)
		if e != nil {
			blog.Errorf("count process instance relation failed, filter: %v, err: %v, rid: %s", usageFilter, e,
				kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		if usageCount > 0 {
			blog.Errorf("forbidden delete service instance be referenced, usageCount: %d, rid: %s", usageCount,
				kit.Rid)
			err := kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
			return err
		}

		serviceInstanceFilter := map[string]int64{common.BKFieldID: instance.ID}
		if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Delete(kit.Ctx,
			serviceInstanceFilter); err != nil {
			blog.Errorf("delete service instance failed, deleteFilter: %+v, err: %v, rid: %s", serviceInstanceFilter,
				err, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
		}
	}
	return nil
}

// generateServiceInstanceName get service instance's name, format: `IP + first process name + first process port`
// 可能应用场景：1. 查询服务实例时组装名称；2. 更新进程信息时根据组装名称直接更新到 `name` 字段
// issue: https://github.com/TencentBlueKing/bk-cmdb/issues/2485
func (p *processOperation) generateServiceInstanceName(kit *rest.Kit, instanceID int64) (string, errors.CCErrorCoder) {

	// get instance
	instance := metadata.ServiceInstance{}
	instanceFilter := map[string]interface{}{
		common.BKFieldID: instanceID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(instanceFilter).One(kit.Ctx,
		&instance); err != nil {
		blog.Errorf("find service instance failed, filter: %v, err: %v, rid: %s", instanceFilter, err, kit.Rid)
		if mongodb.IsNotFoundError(err) {
			return "", kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}
		return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	// get host inner ip
	host := metadata.HostMapStr{}

	hostFilter := map[string]interface{}{
		common.BKHostIDField: instance.HostID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(hostFilter).One(kit.Ctx,
		&host); err != nil {
		blog.Errorf("find host failed, filter: %v, err: %v, rid: %s", hostFilter, err, kit.Rid)
		if mongodb.IsNotFoundError(err) {
			return "", kit.CCError.CCErrorf(common.CCErrCommNotFound)
		}
		return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	instanceName := util.GetStrByInterface(host[common.BKHostInnerIPField])

	// get first process instance relation
	relation := metadata.ProcessInstanceRelation{}
	relationFilter := map[string]interface{}{
		common.BKServiceInstanceIDField: instance.ID,
	}
	order := "id"
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Find(relationFilter).
		Sort(order).One(kit.Ctx,
		&relation); err != nil {
		// relation not found means no process in service instance, service instance's name will only contains ip in that case
		if !mongodb.IsNotFoundError(err) {
			blog.Errorf("find process instance relation failed, filter: %v, err: %v, rid: %s", relationFilter, err,
				kit.Rid)
			return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
	}

	if relation.ProcessID != 0 {
		// get process instance
		process := metadata.Process{}
		processFilter := map[string]interface{}{
			common.BKProcIDField: relation.ProcessID,
		}
		if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseProcess).Find(processFilter).One(kit.Ctx,
			&process); err != nil {
			blog.Errorf("find process failed, filter: %v, err: %v, rid: %s", processFilter, err, kit.Rid)
			if mongodb.IsNotFoundError(err) {
				return "", kit.CCError.CCErrorf(common.CCErrCommNotFound)
			}
			return "", kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}

		if process.ProcessName != nil && len(*process.ProcessName) > 0 {
			instanceName += fmt.Sprintf("_%s", *process.ProcessName)
		}
		for _, bindInfo := range process.BindInfo {
			if bindInfo.Std != nil && bindInfo.Std.Port != nil {
				instanceName += fmt.Sprintf("_%s", *bindInfo.Std.Port)
				break
			}
		}

	}
	return instanceName, nil
}

// ConstructServiceInstanceName construct service instance name
// if no service instance name is defined, use the following rule to construct service name:
// hostInnerIP(if exist) + firstProcessName(if exist) + firstProcessPort(if exist)
func (p *processOperation) ConstructServiceInstanceName(kit *rest.Kit, instanceID int64, host map[string]interface{},
	process *metadata.Process) (string, errors.CCErrorCoder) {

	serviceInstanceName := p.genServiceInstanceName(host, process)

	if err := p.updateServiceInstanceName(kit, instanceID, serviceInstanceName); err != nil {
		return "", err
	}

	return serviceInstanceName, nil
}

// genServiceInstanceName generate service instance name using the following rule:
// hostInnerIP(if exist) + firstProcessName(if exist) + firstProcessPort(if exist)
func (p *processOperation) genServiceInstanceName(host map[string]interface{}, process *metadata.Process) string {
	serviceInstanceName := util.GetStrByInterface(host[common.BKHostInnerIPField])
	if serviceInstanceName == "" {
		serviceInstanceName = util.GetStrByInterface(host[common.BKHostInnerIPv6Field])
	}

	if process != nil {
		if process.ProcessName != nil && len(*process.ProcessName) > 0 {
			serviceInstanceName += fmt.Sprintf("_%s", *process.ProcessName)
		}
		for _, bindInfo := range process.BindInfo {
			if bindInfo.Std != nil && bindInfo.Std.Port != nil {
				serviceInstanceName += fmt.Sprintf("_%s", *bindInfo.Std.Port)
				break
			}
		}
	}

	return serviceInstanceName
}

// ReconstructServiceInstanceName do reconstruct service instance name after process name or process port changed
func (p *processOperation) ReconstructServiceInstanceName(kit *rest.Kit, instanceID int64) errors.CCErrorCoder {
	name, err := p.generateServiceInstanceName(kit, instanceID)
	if err != nil {
		blog.Errorf("generate service instance %d name failed, err: %v, rid: %s", instanceID, err, kit.Rid)
		return err
	}

	return p.updateServiceInstanceName(kit, instanceID, name)
}

func (p *processOperation) updateServiceInstanceName(kit *rest.Kit, instanceID int64,
	serviceInstanceName string) errors.CCErrorCoder {
	if serviceInstanceName == "" {
		return nil
	}

	filter := map[string]interface{}{
		common.BKFieldID: instanceID,
	}
	doc := map[string]interface{}{
		common.BKFieldName: serviceInstanceName,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Update(kit.Ctx, filter, doc)
	if err != nil {
		blog.Errorf("update instance name failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}

// GetBusinessDefaultSetModuleInfo TODO
// GetDefaultModuleIDs get business's default module id, default module type specified by DefaultResModuleFlag
// be careful: it doesn't ensure business have all default module or set
func (p *processOperation) GetBusinessDefaultSetModuleInfo(kit *rest.Kit,
	bizID int64) (metadata.BusinessDefaultSetModuleInfo, errors.CCErrorCoder) {
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
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(defaultModuleCond).Fields(
		common.BKModuleIDField, common.BKDefaultField).All(kit.Ctx, &modules)
	if err != nil {
		blog.Errorf("get default module failed, err: %v, filter: %+v, rid: %s", err, defaultModuleCond, kit.Rid)
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
	err = mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseSet).Find(defaultSetCond).Fields(
		common.BKSetIDField).All(kit.Ctx, &sets)
	if err != nil {
		blog.Errorf("get default set failed, err: %v, filter: %+v, rid: %s", err, defaultSetCond, kit.Rid)
		return defaultSetModuleInfo, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	for _, set := range sets {
		defaultSetModuleInfo.IdleSetID = set.SetID
	}

	return defaultSetModuleInfo, nil
}

// AutoCreateServiceInstanceModuleHost create service instance on host under module base on service template
func (p *processOperation) AutoCreateServiceInstanceModuleHost(kit *rest.Kit, hostIDs []int64,
	moduleIDs []int64) errors.CCErrorCoder {

	// get all related resource info for service instance creation
	params, ccErr := p.prepareAutoCreateSvcInst(kit, hostIDs, moduleIDs)
	if ccErr != nil {
		return ccErr
	}

	if params == nil {
		return nil
	}

	// generate service instances & processes & relations data to create
	serviceInstances, processes, relations, ccErr := p.generateAutoCreateSvcInstData(kit, params)
	if ccErr != nil {
		return ccErr
	}

	if len(serviceInstances) == 0 {
		return nil
	}

	// create service instances
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Insert(kit.Ctx,
		serviceInstances); err != nil {
		blog.Errorf("create service instances failed, err: %v, instance: %+v, rid: %s", err, serviceInstances, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	// create processes
	processes, ccErr = p.dependence.CreateProcessInstances(kit, processes)
	if ccErr != nil {
		blog.Errorf("create process instances(%#v) failed, err: %v, rid: %s", processes, ccErr, kit.Rid)
		return ccErr
	}

	svcInstProcMap := make(map[int64][]mapstr.MapStr)
	svcInstRelMap := make(map[int64][]metadata.ProcessInstanceRelation)
	for i, process := range processes {
		relations[i].ProcessID = process.ProcessID
		svcInstProcMap[process.ServiceInstanceID] = append(svcInstProcMap[process.ServiceInstanceID],
			mapstr.SetValueToMapStrByTags(process))
		svcInstRelMap[process.ServiceInstanceID] = append(svcInstRelMap[process.ServiceInstanceID], *relations[i])
	}

	// create process relations
	relations, ccErr = p.CreateProcessInstanceRelations(kit, relations)
	if ccErr != nil {
		blog.Errorf("create process relations(%#v) failed, err: %v, rid: %s", relations, ccErr, kit.Rid)
		return ccErr
	}

	audit := auditlog.NewSvcInstAuditGenerator()
	auditLogs := make([]metadata.AuditLog, 0)
	genAuditParam := auditlog.NewGenerateAuditCommonParameter(kit.NewKit(), metadata.AuditCreate)

	for _, instance := range serviceInstances {
		// generate audit log for created service instance
		audit.WithServiceInstance([]metadata.ServiceInstance{instance})

		if err := audit.WithProc(genAuditParam, svcInstProcMap[instance.ID], svcInstRelMap[instance.ID]); err != nil {
			return err
		}
		logs := audit.GenerateAuditLog(genAuditParam)
		auditLogs = append(auditLogs, logs...)
	}

	// save service instance audit logs for host transfer, use a new kit to keep it out of the txn
	if err := p.dependence.CreateAuditLogDependence(kit.NewKit(), auditLogs...); err != nil {
		return kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return nil
}

type autoCreateSvcInstParams struct {
	modules         []metadata.ModuleInst
	hostMap         map[int64]metadata.HostMapStr
	existSvcInstMap map[int64]map[int64]struct{}
	procTempMap     map[int64][]metadata.ProcessTemplate
}

func (p *processOperation) prepareAutoCreateSvcInst(kit *rest.Kit, hostIDs, moduleIDs []int64) (
	*autoCreateSvcInstParams, errors.CCErrorCoder) {

	// list modules
	moduleFilter := map[string]interface{}{
		common.BKModuleIDField:          map[string]interface{}{common.BKDBIN: moduleIDs},
		common.BKDefaultField:           common.NormalModuleFlag,
		common.BKServiceTemplateIDField: map[string]interface{}{common.BKDBNE: common.ServiceTemplateIDNotSet},
	}

	modules := make([]metadata.ModuleInst, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(moduleFilter).Fields(
		common.BKModuleIDField, common.BKAppIDField, common.BKServiceTemplateIDField).All(kit.Ctx,
		&modules); err != nil {
		blog.Errorf("get module failed, err: %v, cond: %v, rid: %s", err, moduleFilter, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	serviceTemplateIDs := make([]int64, 0)
	for _, module := range modules {
		if module.ServiceTemplateID != common.ServiceTemplateIDNotSet {
			serviceTemplateIDs = append(serviceTemplateIDs, module.ServiceTemplateID)
		}
	}
	if len(serviceTemplateIDs) == 0 {
		return nil, nil
	}

	// list hosts
	hosts := make([]metadata.HostMapStr, 0)
	fields := []string{common.BKHostIDField, common.BKHostInnerIPField, common.BKHostInnerIPv6Field,
		common.BKHostOuterIPField}
	hostFilter := map[string]interface{}{common.BKHostIDField: map[string]interface{}{common.BKDBIN: hostIDs}}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(hostFilter).Fields(fields...).
		All(kit.Ctx, &hosts); err != nil {
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	hostMap := make(map[int64]metadata.HostMapStr)
	for _, host := range hosts {
		hostID, err := util.GetInt64ByInterface(host[common.BKHostIDField])
		if err != nil {
			blog.Errorf("parse host id failed, err: %v, host: %+v, rid: %s", err, host, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKHostIDField)
		}
		hostMap[hostID] = host
	}

	// list exist service instances
	serviceInstanceFilter := map[string]interface{}{
		common.BKModuleIDField: map[string]interface{}{common.BKDBIN: moduleIDs},
		common.BKHostIDField:   map[string]interface{}{common.BKDBIN: hostIDs},
	}

	serviceInstances := make([]metadata.ServiceInstance, 0)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).
		Fields(common.BKModuleIDField, common.BKHostIDField).All(kit.Ctx, &serviceInstances); err != nil {
		blog.Errorf("list service instance failed, filter: %+v, err: %v, rid: %s", serviceInstanceFilter, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	existServiceInstanceMap := make(map[int64]map[int64]struct{})
	for _, serviceInstance := range serviceInstances {
		if _, ok := existServiceInstanceMap[serviceInstance.HostID]; !ok {
			existServiceInstanceMap[serviceInstance.HostID] = make(map[int64]struct{}, 0)
		}
		existServiceInstanceMap[serviceInstance.HostID][serviceInstance.ModuleID] = struct{}{}
	}

	// list process templates in modules' service templates
	listProcTempOption := metadata.ListProcessTemplatesOption{
		BusinessID:         modules[0].BizID,
		ServiceTemplateIDs: serviceTemplateIDs,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}
	listProcTplResult, ccErr := p.ListProcessTemplates(kit, listProcTempOption)
	if ccErr != nil {
		blog.Errorf("get process templates failed, option: %+v, err: %v, rid: %s", listProcTempOption, ccErr, kit.Rid)
		return nil, ccErr
	}

	if len(listProcTplResult.Info) == 0 {
		return nil, nil
	}

	procTempMap := make(map[int64][]metadata.ProcessTemplate)
	for _, processTemplate := range listProcTplResult.Info {
		procTempMap[processTemplate.ServiceTemplateID] = append(procTempMap[processTemplate.ServiceTemplateID],
			processTemplate)
	}
	return &autoCreateSvcInstParams{modules, hostMap, existServiceInstanceMap, procTempMap}, nil
}

func (p *processOperation) generateAutoCreateSvcInstData(kit *rest.Kit, params *autoCreateSvcInstParams) (
	[]metadata.ServiceInstance, []*metadata.Process, []*metadata.ProcessInstanceRelation, errors.CCErrorCoder) {

	now := time.Now()

	// generate service instances
	serviceInstances := make([]metadata.ServiceInstance, 0)
	for _, module := range params.modules {
		processTemplates := params.procTempMap[module.ServiceTemplateID]
		if len(processTemplates) == 0 {
			blog.Warnf("service template(%d) has no process template, rid: %s", module.ServiceTemplateID, kit.Rid)
			continue
		}

		for hostID := range params.hostMap {
			if _, exist := params.existSvcInstMap[hostID][module.ModuleID]; exist {
				continue
			}

			serviceInstances = append(serviceInstances, metadata.ServiceInstance{
				BizID:             module.BizID,
				ServiceTemplateID: module.ServiceTemplateID,
				HostID:            hostID,
				ModuleID:          module.ModuleID,
				Creator:           kit.User,
				Modifier:          kit.User,
				CreateTime:        now,
				LastTime:          now,
				TenantID:          kit.TenantID,
			})
		}
	}

	if len(serviceInstances) == 0 {
		return serviceInstances, make([]*metadata.Process, 0), make([]*metadata.ProcessInstanceRelation, 0), nil
	}

	ids, err := mongodb.Shard(kit.SysShardOpts()).NextSequences(kit.Ctx, common.BKTableNameServiceInstance,
		len(serviceInstances))
	if err != nil {
		blog.Errorf("generate service instance ids failed, err: %v, rid: %s", err, kit.Rid)
		return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}

	// generate processes & relations
	processes := make([]*metadata.Process, 0)
	relations := make([]*metadata.ProcessInstanceRelation, 0)

	for i, instance := range serviceInstances {
		instance.ID = int64(ids[i])

		processTemplates := params.procTempMap[instance.ServiceTemplateID]
		host := params.hostMap[instance.HostID]

		var firstProc *metadata.Process
		for idx, procTemp := range processTemplates {
			processData, err := procTemp.NewProcess(kit.CCError, instance.BizID, int64(ids[i]),
				kit.TenantID, host)
			if err != nil {
				blog.ErrorJSON("generate process by template %s failed, err: %s, rid: %s", procTemp, err, kit.Rid)
				return nil, nil, nil, errors.New(common.CCErrCommParamsInvalid, err.Error())
			}

			processes = append(processes, processData)
			if idx == 0 {
				firstProc = processData
			}

			relations = append(relations, &metadata.ProcessInstanceRelation{
				BizID:             instance.BizID,
				ServiceInstanceID: instance.ID,
				ProcessTemplateID: procTemp.ID,
				HostID:            instance.HostID,
				TenantID:          kit.TenantID,
			})
		}

		instance.Name = p.genServiceInstanceName(host, firstProc)
		serviceInstances[i] = instance
	}

	return serviceInstances, processes, relations, nil
}

// RemoveTemplateBindingOnModule TODO
func (p *processOperation) RemoveTemplateBindingOnModule(kit *rest.Kit, moduleID int64) errors.CCErrorCoder {
	moduleFilter := map[string]interface{}{
		common.BKModuleIDField: moduleID,
	}
	moduleSimple := struct {
		ServiceTemplateID int64 `field:"service_template_id" bson:"service_template_id" json:"service_template_id"`
		ServiceCategoryID int64 `field:"service_category_id" bson:"service_category_id" json:"service_category_id"`
		BizID             int64 `field:"bk_biz_id" bson:"bk_biz_id" json:"bk_biz_id"`
	}{}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Find(moduleFilter).One(kit.Ctx,
		&moduleSimple); err != nil {
		blog.Errorf("get module by id failed, moduleID: %d, err: %v, rid: %s", moduleID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if moduleSimple.ServiceTemplateID == 0 {
		return kit.CCError.CCError(common.CCErrCoreServiceModuleNotBoundWithTemplate)
	}

	// clear template id field on module
	resetServiceTemplateIDOption := map[string]interface{}{
		common.BKServiceTemplateIDField: common.ServiceTemplateIDNotSet,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseModule).Update(kit.Ctx, moduleFilter,
		resetServiceTemplateIDOption); err != nil {
		blog.Errorf("update service_template_id on module failed, module: %d, err: %v, rid: %s", moduleID, err,
			kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}

	// clear service instance template
	serviceInstanceFilter := map[string]int64{
		common.BKModuleIDField: moduleID,
	}
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameServiceInstance).Update(kit.Ctx,
		serviceInstanceFilter, resetServiceTemplateIDOption); err != nil {
		blog.Errorf("update service_template_id on service instance failed, module: %d, err: %v, rid: %s", moduleID,
			err, kit.Rid)
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
		blog.Errorf("list service instance failed, option: %v, err: %v, rid: %s", listOption, err, kit.Rid)
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
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameProcessInstanceRelation).Update(kit.Ctx,
		processInstanceRelationFilter, resetProcessTemplateIDOption); err != nil {
		blog.Errorf("update service_template_id on process instance relation failed, moduleID: %d, err: %v, rid: %s",
			moduleID, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBUpdateFailed)
	}
	return nil
}
