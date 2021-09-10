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
	"sync"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (lgc *Logic) GetSvcTempSyncStatus(kit *rest.Kit, bizID int64, moduleCond map[string]interface{}, isPartial bool) (
	[]metadata.SvcTempSyncStatus, []metadata.ModuleSyncStatus, errors.CCErrorCoder) {

	moduleFilter := &metadata.QueryCondition{
		Condition: moduleCond,
		Fields: []string{common.BKModuleIDField, common.BKServiceTemplateIDField, common.BKModuleNameField,
			common.BKServiceCategoryIDField},
	}

	moduleRes := new(metadata.ResponseModuleInstance)
	err := lgc.CoreAPI.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDModule, moduleFilter, moduleRes)
	if err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, nil, err
	}

	if err := moduleRes.CCError(); err != nil {
		blog.Errorf("get module failed, filter: %#v, err: %v, rid: %s", moduleFilter, err, kit.Rid)
		return nil, nil, err
	}

	moduleSyncStatuses := make([]metadata.ModuleSyncStatus, 0)
	svcTempSyncStatuses := make([]metadata.SvcTempSyncStatus, 0)
	if len(moduleRes.Data.Info) == 0 {
		return svcTempSyncStatuses, moduleSyncStatuses, nil
	}

	svcTempModuleMap := make(map[int64][]*metadata.ModuleInst)
	for index, module := range moduleRes.Data.Info {
		if module.ServiceTemplateID != common.ServiceTemplateIDNotSet {
			svcTempModuleMap[module.ServiceTemplateID] = append(svcTempModuleMap[module.ServiceTemplateID],
				&moduleRes.Data.Info[index])
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
		blog.ErrorJSON("list service templates failed, err: %s, input: %s, rid: %s", err, svcTempOpt, kit.Rid)
		return nil, nil, err
	}

	for _, svcTemp := range templates.Info {
		modules := svcTempModuleMap[svcTemp.ID]

		needSync, statuses, err := lgc.getSvcTempSyncStatus(kit, &svcTemp, modules, isPartial)
		if err != nil {
			blog.ErrorJSON("get service template sync status failed, err: %s, template: %s, modules: %s, rid: %s",
				err, svcTemp, modules, kit.Rid)
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

func (lgc *Logic) getSvcTempSyncStatus(kit *rest.Kit, svcTemp *metadata.ServiceTemplate, modules []*metadata.ModuleInst,
	isPartial bool) (bool, []metadata.ModuleSyncStatus, errors.CCErrorCoder) {

	statuses := make([]metadata.ModuleSyncStatus, 0)
	var wg sync.WaitGroup
	var firstErr errors.CCErrorCoder
	pipeline := make(chan bool, 10)
	isFinish, needSync, isProcTempMapSet := false, false, false
	procTempMap := make(map[int64]*metadata.ProcessTemplate)

	for _, module := range modules {
		if isFinish {
			break
		}

		// check if module info has difference with service template
		if module.ModuleName != svcTemp.Name || module.ServiceCategoryID != svcTemp.ServiceCategoryID {
			if isPartial {
				return true, statuses, nil
			}

			statuses = append(statuses, metadata.ModuleSyncStatus{ModuleID: module.ModuleID, NeedSync: true})
			needSync = true
			continue
		}

		// get process templates to compare with the processes of the module
		if !isProcTempMapSet {
			procTempOpt := &metadata.ListProcessTemplatesOption{
				BusinessID:         svcTemp.BizID,
				ServiceTemplateIDs: []int64{svcTemp.ID},
			}

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

		module.BizID = svcTemp.BizID
		module.ServiceTemplateID = svcTemp.ID

		go func(module *metadata.ModuleInst) {
			defer func() {
				wg.Done()
				<-pipeline
			}()

			moduleNeedSync, err := lgc.getModuleProcessSyncStatus(kit, module, procTempMap)
			if err != nil {
				blog.ErrorJSON("get module(%s) process sync status failed, err: %s, rid: %s", module, err, kit.Rid)
				if firstErr == nil {
					firstErr = err
				}
				isFinish = true
				return
			}

			if moduleNeedSync {
				if isPartial {
					isFinish = true
					needSync = true
					return
				}

				statuses = append(statuses, metadata.ModuleSyncStatus{ModuleID: module.ModuleID, NeedSync: true})
				needSync = true
				return
			}

			statuses = append(statuses, metadata.ModuleSyncStatus{ModuleID: module.ModuleID, NeedSync: false})
		}(module)
	}

	wg.Wait()
	if firstErr != nil {
		return false, nil, firstErr
	}
	return needSync, statuses, nil
}

func (lgc *Logic) getModuleProcessSyncStatus(kit *rest.Kit, module *metadata.ModuleInst,
	procTempMap map[int64]*metadata.ProcessTemplate) (bool, errors.CCErrorCoder) {

	// get service instances by module
	svcInstOpt := &metadata.ListServiceInstanceOption{
		BusinessID:        module.BizID,
		ServiceTemplateID: module.ServiceTemplateID,
		ModuleIDs:         []int64{module.ModuleID},
		Page:              metadata.BasePage{Limit: common.BKNoLimit},
	}
	serviceInstances, err := lgc.CoreAPI.CoreService().Process().ListServiceInstance(kit.Ctx, kit.Header, svcInstOpt)
	if err != nil {
		blog.ErrorJSON("list service instance failed, option: %s, err: %s, rid: %s", svcInstOpt, err, kit.Rid)
		return false, err
	}

	// get host ids by module
	hostIDFilter := []map[string]interface{}{{
		common.BKModuleIDField: module.ModuleID,
		common.BKAppIDField:    module.BizID,
	}}
	hostIDCount, err := lgc.CoreAPI.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameModuleHostConfig, hostIDFilter)
	if err != nil {
		blog.Errorf("get host id count failed, err: %v, option: %#v, rid: %s", err, hostIDFilter, kit.Rid)
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

	// get process instance relations by service instances
	procRelOpt := metadata.ListProcessInstanceRelationOption{
		BusinessID:         module.BizID,
		ServiceInstanceIDs: serviceInstanceIDs,
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
	}

	relations, err := lgc.CoreAPI.CoreService().Process().ListProcessInstanceRelation(kit.Ctx, kit.Header, &procRelOpt)
	if err != nil {
		blog.ErrorJSON("list process instance relation failed, option: %s, err: %s, rid: %s", procRelOpt, err, kit.Rid)
		return false, err
	}

	if len(relations.Info) == 0 {
		return len(procTempMap) != 0, err
	}

	// find all the process instance detail by ids
	procIDs := make([]int64, 0)
	processTemplateReferenced := make(map[int64]struct{})
	for _, relation := range relations.Info {
		procIDs = append(procIDs, relation.ProcessID)

		// record the used process template for checking whether a new process template has been added to service template.
		processTemplateReferenced[relation.ProcessTemplateID] = struct{}{}
	}

	// check whether a new process template has been added
	for templateID := range procTempMap {
		if _, exist := processTemplateReferenced[templateID]; exist {
			continue
		}
		return true, nil
	}

	processDetails, err := lgc.ListProcessInstanceWithIDs(kit, procIDs)
	if err != nil {
		blog.Errorf("list process instance with IDs failed, err: %v, procIDs: %+v, rid: %s", err, procIDs, kit.Rid)
		return false, err
	}

	procID2Detail := make(map[int64]*metadata.Process)
	for idx, p := range processDetails {
		procID2Detail[p.ProcessID] = &processDetails[idx]
	}

	// get host map for process bind info compare use
	hostMap, err := lgc.GetHostIPMapByID(kit, hostIDs)
	if err != nil {
		blog.Errorf("get host map failed, err: %v, ids: %+v, rid: %s", err, hostIDs, kit.Rid)
		return false, err
	}

	// compare the process instance with it's process template one by one
	for _, relation := range relations.Info {
		process, ok := procID2Detail[relation.ProcessID]
		if !ok {
			blog.ErrorJSON("process doesn't exist, id: %d, rid: %s", err, relation.ProcessID, kit.Rid)
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

	return false, nil
}
