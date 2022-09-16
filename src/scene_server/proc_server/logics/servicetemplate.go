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

package logics

import (
	"reflect"
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logic) getModuleMapStr(kit *rest.Kit, cond map[string]interface{}) ([]mapstr.MapStr, errors.CCErrorCoder) {

	option := &metadata.QueryCondition{
		Condition:      cond,
		DisableCounter: true,
	}

	modules, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, option)
	if err != nil {
		blog.Errorf("get modules failed, option: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrHostGetModuleFail, err.Error())
	}

	return modules.Info, nil
}

// GetSvcTempSyncStatus 获取服务模板的同步状态
func (lgc *Logic) GetSvcTempSyncStatus(kit *rest.Kit, bizID int64, moduleFilter map[string]interface{},
	isPartial bool) ([]metadata.SvcTempSyncStatus, []metadata.ModuleSyncStatus, errors.CCErrorCoder) {

	// module detail
	modules, err := lgc.getModuleMapStr(kit, moduleFilter)
	if err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, nil, err
	}

	moduleSyncStatuses := make([]metadata.ModuleSyncStatus, 0)
	svcTempSyncStatuses := make([]metadata.SvcTempSyncStatus, 0)

	if len(modules) == 0 {
		return svcTempSyncStatuses, moduleSyncStatuses, nil
	}

	svcTempModuleMap := make(map[int64][]mapstr.MapStr)
	for index, module := range modules {
		serviceTemplateID, err := util.GetInt64ByInterface(module[common.BKServiceTemplateIDField])
		if err != nil {
			blog.Errorf("get serviceTemplateID failed, module: %+v, err: %v, rid: %s", module, err, kit.Rid)
			return nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKServiceTemplateIDField)
		}
		if serviceTemplateID != common.ServiceTemplateIDNotSet {
			svcTempModuleMap[serviceTemplateID] = append(svcTempModuleMap[serviceTemplateID], modules[index])
		}
	}

	if len(svcTempModuleMap) == 0 {
		return svcTempSyncStatuses, moduleSyncStatuses, nil
	}

	svcTempIDs := make([]int64, 0)
	for svcTempID := range svcTempModuleMap {
		svcTempIDs = append(svcTempIDs, svcTempID)
	}

	svcTempOpt := metadata.ListServiceTemplateOption{
		BusinessID: bizID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		ServiceTemplateIDs: svcTempIDs,
	}

	templates, err := lgc.CoreAPI.CoreService().Process().ListServiceTemplates(kit.Ctx, kit.Header, &svcTempOpt)
	if err != nil {
		blog.Errorf("list service templates failed, input: %+v, err: %v, rid: %s", svcTempOpt, err, kit.Rid)
		return nil, nil, err
	}

	for _, svcTemp := range templates.Info {
		modules := svcTempModuleMap[svcTemp.ID]

		needSync, statuses, err := lgc.getSvcTempSyncStatus(kit, &svcTemp, modules, isPartial)
		if err != nil {
			blog.Errorf("get service template sync status failed, template: %+v, modules: %+v, err: %v, rid: %s",
				svcTemp, modules, err, kit.Rid)
			return nil, nil, err
		}

		if isPartial {
			svcTempSyncStatuses = append(svcTempSyncStatuses, metadata.SvcTempSyncStatus{
				ServiceTemplateID: svcTemp.ID,
				NeedSync:          needSync,
			})
		} else {
			moduleSyncStatuses = append(moduleSyncStatuses, statuses...)
		}
	}

	return svcTempSyncStatuses, moduleSyncStatuses, nil
}

func getModuleNameAndID(kit *rest.Kit, module mapstr.MapStr) (string, int64, errors.CCErrorCoder) {

	moduleName := util.GetStrByInterface(module[common.BKModuleNameField])

	moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
	if err != nil {
		blog.Errorf("get moduleID failed, module: %#v, err: %v, rid: %s", module, err, kit.Rid)
		return "", 0, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed, common.BKModuleIDField)
	}
	return moduleName, moduleID, nil
}

func (lgc *Logic) getSvcTemplateAttrsInfo(kit *rest.Kit, svcTemp *metadata.ServiceTemplate) ([]string,
	map[int64]interface{}, map[int64]string, errors.CCErrorCoder) {
	attrIDs, srvTempAttrValueMap, cErr := lgc.getSrvTemplateAttrIdAndPropertyValue(kit, svcTemp.BizID, svcTemp.ID)
	if cErr != nil {
		return []string{}, nil, nil, cErr
	}

	propertyIDs, attrIdPropertyIdMap, cErr := lgc.getModuleAttrIDAndPropertyID(kit, attrIDs)
	if cErr != nil {
		return []string{}, nil, nil, cErr
	}
	return propertyIDs, srvTempAttrValueMap, attrIdPropertyIdMap, nil
}

func (lgc *Logic) getSvcTempSyncStatus(kit *rest.Kit, svcTemp *metadata.ServiceTemplate, modules []mapstr.MapStr,
	isPartial bool) (bool, []metadata.ModuleSyncStatus, errors.CCErrorCoder) {

	var wg sync.WaitGroup
	var firstErr errors.CCErrorCoder

	isFinish, needSync, isProcTempMapSet, pipeline := false, false, false, make(chan bool, 10)
	procTempMap, statuses := make(map[int64]*metadata.ProcessTemplate), make([]metadata.ModuleSyncStatus, 0)

	propertyIDs, srvTempAttrValueMap, attrIdPropertyIdMap, cErr := lgc.getSvcTemplateAttrsInfo(kit, svcTemp)
	if cErr != nil {
		return false, nil, cErr
	}
	for _, module := range modules {
		if isFinish {
			break
		}

		moduleName := util.GetStrByInterface(module[common.BKModuleNameField])
		moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
		if err != nil {
			blog.Errorf("get moduleID failed, module: %#v, err: %v, rid: %s", module, err, kit.Rid)
			return false, nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed, common.BKModuleIDField)
		}

		if moduleName != svcTemp.Name {
			if isPartial {
				return true, statuses, nil
			}

			statuses = append(statuses, metadata.ModuleSyncStatus{ModuleID: moduleID, NeedSync: true})
			needSync = true
			continue
		}

		// get process templates to compare with the processes of the module
		if !isProcTempMapSet {
			procTempOpt := &metadata.ListProcessTemplatesOption{
				BusinessID: svcTemp.BizID, ServiceTemplateIDs: []int64{svcTemp.ID}}

			procTemps, err := lgc.CoreAPI.CoreService().Process().ListProcessTemplates(kit.Ctx, kit.Header, procTempOpt)
			if err != nil {
				blog.Errorf("list process templates failed, input: %#v, err: %v, rid: %s", procTempOpt, err, kit.Rid)
				return false, nil, err
			}

			for idx, procTemp := range procTemps.Info {
				procTempMap[procTemp.ID] = &procTemps.Info[idx]
			}
			isProcTempMapSet = true
		}

		pipeline <- true
		wg.Add(1)

		go func(module mapstr.MapStr, bizID, serviceTemplateID, moduleID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			attrNeedSync := lgc.getModuleAttrNeedSync(module, propertyIDs, attrIdPropertyIdMap, srvTempAttrValueMap)

			moduleNeedSync, err := lgc.getModuleProcessSyncStatus(kit, bizID, serviceTemplateID, moduleID, procTempMap)
			if err != nil {
				blog.Errorf("get module(%+v) process sync status failed, err: %v, rid: %s", module, err, kit.Rid)
				if firstErr == nil {
					firstErr = err
				}
				isFinish = true
				return
			}

			if moduleNeedSync || attrNeedSync {
				if isPartial {
					isFinish, needSync = true, true
					return
				}

				statuses = append(statuses, metadata.ModuleSyncStatus{ModuleID: moduleID, NeedSync: true})
				needSync = true
				return
			}

			statuses = append(statuses, metadata.ModuleSyncStatus{ModuleID: moduleID, NeedSync: false})
		}(module, svcTemp.BizID, svcTemp.ID, moduleID)
	}

	wg.Wait()
	if firstErr != nil {
		return false, nil, firstErr
	}
	return needSync, statuses, nil
}

// getSrvTemplateAttrIdAndPropertyValue 获取服务模板的属性id以及对应的属性值
func (lgc *Logic) getSrvTemplateAttrIdAndPropertyValue(kit *rest.Kit, bizID, serviceTemplateID int64) ([]int64,
	map[int64]interface{}, errors.CCErrorCoder) {

	option := &metadata.ListServTempAttrOption{
		BizID:  bizID,
		ID:     serviceTemplateID,
		Fields: []string{common.BKAttributeIDField, common.BKPropertyValueField},
	}

	data, cErr := lgc.Engine.CoreAPI.CoreService().Process().ListServiceTemplateAttribute(kit.Ctx, kit.Header, option)
	if cErr != nil {
		blog.Errorf("list service template attributes failed, bizID: %d, service template id: %d, err: %v, rid: %s",
			bizID, serviceTemplateID, cErr, kit.Rid)
		return nil, nil, cErr
	}

	attrIDs := make([]int64, 0)
	srvTemplateAttrValueMap := make(map[int64]interface{})

	for _, attr := range data.Attributes {
		attrIDs = append(attrIDs, attr.AttributeID)
		srvTemplateAttrValueMap[attr.AttributeID] = attr.PropertyValue
	}

	return attrIDs, srvTemplateAttrValueMap, nil
}

// getModuleAttrIDAndPropertyID 根据模块属性ID获取对应的propertyID列表以及属性ID与propertyID的对应关系
func (lgc *Logic) getModuleAttrIDAndPropertyID(kit *rest.Kit, attrIDs []int64) ([]string, map[int64]string,
	errors.CCErrorCoder) {

	attrIdPropertyMap := make(map[int64]string)
	if len(attrIDs) == 0 {
		return []string{}, attrIdPropertyMap, nil
	}

	option := &metadata.QueryCondition{
		Fields: []string{common.BKFieldID, common.BKPropertyIDField},
		Page:   metadata.BasePage{Limit: common.BKNoLimit},
		Condition: map[string]interface{}{
			common.BKFieldID: map[string]interface{}{
				common.BKDBIN: attrIDs,
			},
		},
		DisableCounter: true,
	}

	res, err := lgc.CoreAPI.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDModule, option)
	if err != nil {
		blog.Errorf("read model attribute failed, err: %v, option: %#v, rid: %s", err, option, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeSelectFailed)
	}
	propertyIDs := make([]string, 0)
	for _, attrs := range res.Info {
		propertyIDs = append(propertyIDs, attrs.PropertyID)
		attrIdPropertyMap[attrs.ID] = attrs.PropertyID
	}

	return propertyIDs, attrIdPropertyMap, nil
}

// getModuleAttrNeedSync 获取模块与服务模板之间属性是否有差异
func (lgc *Logic) getModuleAttrNeedSync(module mapstr.MapStr, propertyIDs []string,
	attrIdPropertyIdMap map[int64]string, srvTemplateAttrValueMap map[int64]interface{}) bool {

	if len(propertyIDs) == 0 {
		return false
	}
	// 对比是否有不同
	for id, attr := range srvTemplateAttrValueMap {
		if !reflect.DeepEqual(attr, module[attrIdPropertyIdMap[id]]) {
			return true
		}
	}
	return false
}

func (lgc *Logic) getServiceInstancesAndHostIdCount(kit *rest.Kit, bizID, serviceTemplateID, moduleID int64) (
	*metadata.MultipleServiceInstance, []int64, errors.CCErrorCoder) {
	// get service instances by module
	svcInstOpt := &metadata.ListServiceInstanceOption{
		BusinessID:        bizID,
		ServiceTemplateID: serviceTemplateID,
		ModuleIDs:         []int64{moduleID},
		Page:              metadata.BasePage{Limit: common.BKNoLimit},
	}
	serviceInstances, err := lgc.CoreAPI.CoreService().Process().ListServiceInstance(kit.Ctx, kit.Header, svcInstOpt)
	if err != nil {
		blog.ErrorJSON("list service instance failed, option: %s, err: %s, rid: %s", svcInstOpt, err, kit.Rid)
		return nil, nil, err
	}

	// get host ids by module
	hostIDFilter := []map[string]interface{}{{
		common.BKModuleIDField: moduleID,
		common.BKAppIDField:    bizID,
	}}
	hostIDCount, err := lgc.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, hostIDFilter)
	if err != nil {
		blog.Errorf("get host id count failed, err: %v, option: %#v, rid: %s", err, hostIDFilter, kit.Rid)
		return nil, nil, err
	}
	return serviceInstances, hostIDCount, nil
}

func (lgc *Logic) getProcessInstanceRelation(kit *rest.Kit, bizID int64, serviceInstanceIDs []int64) (
	*metadata.MultipleProcessInstanceRelation, errors.CCErrorCoder) {

	// get process instance relations by service instances
	procRelOpt := metadata.ListProcessInstanceRelationOption{
		BusinessID:         bizID,
		ServiceInstanceIDs: serviceInstanceIDs,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}

	relations, err := lgc.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, &procRelOpt)
	if err != nil {
		blog.ErrorJSON("list process instance relation failed, option: %s, err: %s, rid: %s", procRelOpt, err, kit.Rid)
		return nil, err
	}
	return relations, nil
}

func (lgc *Logic) getServiceRelationMapAndProcessDetails(kit *rest.Kit,
	relations *metadata.MultipleProcessInstanceRelation) (map[int64][]metadata.ProcessInstanceRelation,
	map[int64]*metadata.Process, errors.CCErrorCoder) {

	// find all the process instance detail by ids
	procIDs := make([]int64, 0)
	for _, relation := range relations.Info {
		procIDs = append(procIDs, relation.ProcessID)
	}

	serviceRelationMap := make(map[int64][]metadata.ProcessInstanceRelation)
	for _, r := range relations.Info {
		serviceRelationMap[r.ServiceInstanceID] = append(serviceRelationMap[r.ServiceInstanceID], r)
	}
	// check whether a new process template has been added

	processDetails, err := lgc.ListProcessInstanceWithIDs(kit, procIDs)
	if err != nil {
		blog.Errorf("list process instance with IDs failed, err: %v, procIDs: %+v, rid: %s", err, procIDs, kit.Rid)
		return nil, nil, err
	}

	procID2Detail := make(map[int64]*metadata.Process)
	for idx, p := range processDetails {
		procID2Detail[p.ProcessID] = &processDetails[idx]
	}
	return serviceRelationMap, procID2Detail, nil
}

func (lgc *Logic) getModuleProcessSyncStatus(kit *rest.Kit, bizID, serviceTemplateID, moduleID int64,
	procTempMap map[int64]*metadata.ProcessTemplate) (bool, errors.CCErrorCoder) {

	serviceInstances, hostIDCount, err := lgc.getServiceInstancesAndHostIdCount(kit, bizID, serviceTemplateID, moduleID)
	if err != nil {
		return false, err
	}

	// need to sync when the module has process templates and hosts but has no service instances
	if len(serviceInstances.Info) == 0 {
		if hostIDCount[0] == 0 || len(procTempMap) == 0 {
			return false, nil
		}
		return true, nil
	}

	if len(procTempMap) == 0 {
		return true, nil
	}

	if int64(len(serviceInstances.Info)) != hostIDCount[0] {
		return true, nil
	}

	serviceInstanceIDs := make([]int64, 0)
	hostIDs := make([]int64, 0)
	for _, serviceInstance := range serviceInstances.Info {
		serviceInstanceIDs = append(serviceInstanceIDs, serviceInstance.ID)
		hostIDs = append(hostIDs, serviceInstance.HostID)
	}

	relations, err := lgc.getProcessInstanceRelation(kit, bizID, serviceInstanceIDs)
	if err != nil {
		return false, err
	}

	if len(relations.Info) == 0 {
		return len(procTempMap) != 0, err
	}

	serviceRelationMap, procID2Detail, err := lgc.getServiceRelationMapAndProcessDetails(kit, relations)
	if err != nil {
		return false, err
	}

	// get host map for process bind info compare use
	hostMap, err := lgc.GetHostIPMapByID(kit, hostIDs)
	if err != nil {
		blog.Errorf("get host map failed, err: %v, ids: %+v, rid: %s", err, hostIDs, kit.Rid)
		return false, err
	}

	for _, serviceInst := range serviceInstances.Info {
		relations := serviceRelationMap[serviceInst.ID]
		processTemplateReferenced := make(map[int64]struct{})

		// compare the process instance with it's process template one by one
		for _, relation := range relations {
			processTemplateReferenced[relation.ProcessTemplateID] = struct{}{}
			process, ok := procID2Detail[relation.ProcessID]
			if !ok {
				blog.Errorf("process doesn't exist, id: %d, rid: %s", relation.ProcessID, kit.Rid)
				return false, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcIDField)
			}

			property, exist := procTempMap[relation.ProcessTemplateID]
			if !exist {
				return true, nil
			}

			_, isChanged, diffErr := lgc.DiffWithProcessTemplate(property.Property, process, hostMap[relation.HostID],
				map[string]metadata.Attribute{}, false)
			if diffErr != nil {
				blog.Errorf("diff process %d with template failed, err: %v, rid: %s", relation.ProcessID, diffErr, kit.Rid)
				return false, errors.New(common.CCErrCommParamsInvalid, diffErr.Error())
			}

			if isChanged {
				return true, nil
			}
		}
		for templateID := range procTempMap {
			if _, exist := processTemplateReferenced[templateID]; exist {
				continue
			}
			return true, nil
		}
	}
	return false, nil
}
