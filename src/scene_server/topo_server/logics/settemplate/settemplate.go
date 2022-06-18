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

package settemplate

import (
	"sync"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

type SetTemplate interface {
	DiffSetTplWithInst(kit *rest.Kit, bizID int64, setTemplateID int64,
		option metadata.DiffSetTplWithInstOption) (metadata.SetDiff, errors.CCErrorCoder)
	SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64,
		option metadata.SyncSetTplToInstOption) errors.CCErrorCoder
	GetLatestSyncTaskDetail(kit *rest.Kit, taskCond metadata.ListAPITaskDetail) (
		map[int64]*metadata.APITaskDetail, errors.CCErrorCoder)
	CheckSetInstUpdateToDateStatus(kit *rest.Kit, bizID int64, setTemplateID int64) (
		*metadata.SetTemplateUpdateToDateStatus, errors.CCErrorCoder)
	ListSetTemplateSyncHistory(kit *rest.Kit, option *metadata.ListSetTemplateSyncStatusOption) (
		*metadata.ListAPITaskSyncStatusResult, errors.CCErrorCoder)
	ListSetTemplateSyncStatus(kit *rest.Kit, option *metadata.ListSetTemplateSyncStatusOption) (
		*metadata.ListAPITaskSyncStatusResult, errors.CCErrorCoder)
}

func NewSetTemplate(client apimachinery.ClientSetInterface) SetTemplate {
	return &setTemplate{
		client: client,
	}
}

type setTemplate struct {
	client apimachinery.ClientSetInterface
}

func (st *setTemplate) getSetResult(kit *rest.Kit, bizID, setTemplateID, setID int64) (*metadata.ResponseSetInstance,
	errors.CCErrorCoder) {
	filter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKSetIDField:         setID,
		}),
		DisableCounter: true,
	}

	set := new(metadata.ResponseSetInstance)
	if err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDSet,
		filter, set); err != nil {
		blog.Errorf("get set failed, bizID: %d, setTemplateID: %d, setID: %d, err: %d, rid: %s", bizID, setTemplateID,
			setID, err, kit.Rid)
		return nil, err
	}
	if err := set.CCError(); err != nil {
		blog.Errorf("get set http reply failed, bizID: %d, setTemplateID: %d, setID: %d, filter: %+v, reply: %v, "+
			"rid: %s", bizID, setTemplateID, setID, filter, set, kit.Rid)
		return nil, err
	}

	if len(set.Data.Info) != 1 {
		blog.Errorf("get set failed, setID: %d, rid: %s", setID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField)
	}
	return set, nil
}

func (st *setTemplate) getModuleResult(kit *rest.Kit, bizID, setTemplateID, setID int64) ([]metadata.ModuleInst,
	errors.CCErrorCoder) {

	filter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKParentIDField:      setID,
		}),
		DisableCounter: true,
	}

	modules := new(metadata.ResponseModuleInstance)
	if err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		filter, modules); err != nil {
		blog.Errorf("list modules failed, bizID: %d, setTemplateID: %d, setID: %d, err: %v, rid: %s",
			bizID, setTemplateID, setID, err, kit.Rid)
		return nil, err
	}
	if err := modules.CCError(); err != nil {
		blog.Errorf("list module http reply failed, bizID: %d, setTemplateID: %d, setID: %d, filter: %+v, reply: %+v,"+
			" rid: %s", bizID, setTemplateID, setID, filter, modules, kit.Rid)
		return nil, err
	}

	if len(modules.Data.Info) == 0 {
		blog.Errorf("list module http reply failed, bizID: %d, setTemplateID: %d, setID: %d, filter: %+v, reply: %+v,"+
			" rid: %s", bizID, setTemplateID, setID, filter, modules, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKModuleIDField)
	}
	return modules.Data.Info, nil
}

func (st *setTemplate) DiffSetTplWithInst(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.DiffSetTplWithInstOption) (metadata.SetDiff, errors.CCErrorCoder) {

	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header, bizID,
		setTemplateID)
	if err != nil {
		blog.Errorf("list service templates failed, bizID: %d, setTemplateID: %d, err: %v, rid: %s", bizID,
			setTemplateID, err, kit.Rid)
		return metadata.SetDiff{}, err
	}
	serviceTemplateMap := make(map[int64]metadata.ServiceTemplate)
	for _, svcTpl := range serviceTemplates {
		serviceTemplateMap[svcTpl.ID] = svcTpl
	}

	set, err := st.getSetResult(kit, bizID, setTemplateID, option.SetID)
	if err != nil {
		return metadata.SetDiff{}, err
	}

	modules, err := st.getModuleResult(kit, bizID, setTemplateID, option.SetID)
	if err != nil {
		return metadata.SetDiff{}, err
	}

	topoTree, ccErr := st.client.CoreService().Mainline().SearchMainlineInstanceTopo(kit.Ctx, kit.Header, bizID, false)
	if ccErr != nil {
		blog.Errorf("ListSetTplRelatedSetsWeb failed, bizID: %d, err: %v, rid: %s", bizID, ccErr, kit.Rid)
		return metadata.SetDiff{}, ccErr
	}
	// diff
	moduleDiff := DiffServiceTemplateWithModules(serviceTemplates, modules)
	setDiff := metadata.SetDiff{
		ModuleDiffs: moduleDiff,
		SetID:       option.SetID,
	}
	setDiff.SetDetail = set.Data.Info[0]

	// add topo path info
	setPath := topoTree.TraversalFindNode(common.BKInnerObjIDSet, option.SetID)
	topoPath := make([]metadata.TopoInstanceNodeSimplify, 0)
	for _, pathNode := range setPath {
		nodeSimplify := metadata.TopoInstanceNodeSimplify{
			ObjectID:     pathNode.ObjectID,
			InstanceID:   pathNode.InstanceID,
			InstanceName: pathNode.InstanceName,
		}
		topoPath = append(topoPath, nodeSimplify)
	}
	setDiff.TopoPath = topoPath
	setDiff.UpdateNeedSyncField()
	return setDiff, nil
}

func (st *setTemplate) SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.SyncSetTplToInstOption) errors.CCErrorCoder {

	var (
		wg       sync.WaitGroup
		firstErr errors.CCErrorCoder
	)

	pipeline := make(chan bool, 10)
	setDiffs := make([]metadata.SetDiff, 0)

	for _, setID := range option.SetIDs {
		pipeline <- true
		wg.Add(1)

		go func(bizID, setTemplateID, setID int64) {
			defer func() {
				wg.Done()
				<-pipeline
			}()
			option := metadata.DiffSetTplWithInstOption{
				SetID: setID,
			}
			setDiff, err := st.DiffSetTplWithInst(kit, bizID, setTemplateID, option)
			if err != nil {
				blog.Errorf("diff set template with instance failed, bizID: %d, set template ID: %d, setID: %d, "+
					"err: %v, rid: %s", bizID, setTemplateID, setID, err, kit.Rid)
				if firstErr == nil {
					firstErr = err
				}
				return
			}
			setDiffs = append(setDiffs, setDiff)

		}(bizID, setTemplateID, setID)
	}
	wg.Wait()
	if firstErr != nil {
		return firstErr
	}

	for _, setDiff := range setDiffs {
		blog.V(3).Infof("dispatch synchronize task on set [%s](%d), rid: %s",
			setDiff.SetDetail.SetName, setDiff.SetID, kit.Rid)
		tasks := make([]metadata.SyncModuleTask, 0)
		for _, moduleDiff := range setDiff.ModuleDiffs {
			task := metadata.SyncModuleTask{
				Set:         setDiff.SetDetail,
				ModuleDiff:  moduleDiff,
				SetTopoPath: setDiff.TopoPath,
			}
			tasks = append(tasks, task)
		}
		taskDetail, err := st.DispatchTask4ModuleSync(kit, common.SyncSetTaskFlag, setDiff.SetID, tasks...)
		if err != nil {
			return err
		}
		if blog.V(3) {
			blog.InfoJSON("dispatch synchronize task on set [%s](%s) success, result: %s, rid: %s",
				setDiff.SetDetail.SetName, setDiff.SetID, taskDetail, kit.Rid)
		}
	}
	return nil
}

// DispatchTask4ModuleSync dispatch synchronize task
func (st *setTemplate) DispatchTask4ModuleSync(kit *rest.Kit, taskType string, setID int64,
	tasks ...metadata.SyncModuleTask) (metadata.APITaskDetail, errors.CCErrorCoder) {

	taskDetail := metadata.APITaskDetail{}
	tasksData := make([]interface{}, 0)
	for _, task := range tasks {
		tasksData = append(tasksData, task)
	}
	createTaskResult, err := st.client.TaskServer().Task().Create(kit.Ctx, kit.Header, taskType, setID, tasksData)
	if err != nil {
		blog.ErrorJSON("dispatch synchronize task failed, task: %s, err: %s, rid: %s", tasks, err.Error(), kit.Rid)
		return taskDetail, err
	}
	blog.InfoJSON("dispatch synchronize task success, task: %s, create result: %s, rid: %s",
		tasks, createTaskResult, kit.Rid)
	return createTaskResult, nil
}

// DiffServiceTemplateWithModules diff modules with template in one set
func DiffServiceTemplateWithModules(serviceTemplates []metadata.ServiceTemplate,
	modules []metadata.ModuleInst) []metadata.SetModuleDiff {

	svcTplMap := make(map[int64]metadata.ServiceTemplate, len(serviceTemplates))
	svcTplHitMap := make(map[int64]bool, len(serviceTemplates))
	for _, svcTpl := range serviceTemplates {
		svcTplMap[svcTpl.ID] = svcTpl
		svcTplHitMap[svcTpl.ID] = false
	}

	moduleMap := make(map[int64]metadata.ModuleInst)
	for _, module := range modules {
		moduleMap[module.ModuleID] = module
	}

	moduleDiffs := make([]metadata.SetModuleDiff, 0)
	for _, module := range modules {
		template, ok := svcTplMap[module.ServiceTemplateID]
		if !ok {
			moduleDiffs = append(moduleDiffs, metadata.SetModuleDiff{
				ModuleID:            module.ModuleID,
				ModuleName:          module.ModuleName,
				ServiceTemplateID:   module.ServiceTemplateID,
				ServiceTemplateName: "",
				DiffType:            metadata.ModuleDiffRemove,
			})
			continue
		}
		if _, ok := svcTplHitMap[module.ServiceTemplateID]; ok {
			svcTplHitMap[module.ServiceTemplateID] = true
		}
		diffType := metadata.ModuleDiffUnchanged
		if module.ModuleName != template.Name {
			diffType = metadata.ModuleDiffChanged
		}
		moduleDiffs = append(moduleDiffs, metadata.SetModuleDiff{
			ModuleID:            module.ModuleID,
			ModuleName:          module.ModuleName,
			ServiceTemplateID:   module.ServiceTemplateID,
			ServiceTemplateName: template.Name,
			DiffType:            diffType,
		})
	}

	for templateID, hit := range svcTplHitMap {
		if hit {
			continue
		}
		template := svcTplMap[templateID]
		moduleDiffs = append(moduleDiffs, metadata.SetModuleDiff{
			ModuleID:            0,
			ModuleName:          template.Name,
			ServiceTemplateID:   templateID,
			ServiceTemplateName: template.Name,
			DiffType:            metadata.ModuleDiffAdd,
		})
	}
	return moduleDiffs
}

// CheckSetTplInstLatest 检查通过集群模板 setTemplateID 实例化的集群是否都已经达到最新状态
func (st *setTemplate) CheckSetInstUpdateToDateStatus(kit *rest.Kit, bizID int64,
	setTemplateID int64) (*metadata.SetTemplateUpdateToDateStatus, errors.CCErrorCoder) {

	result := new(metadata.SetTemplateUpdateToDateStatus)
	result.SetTemplateID = setTemplateID
	result.NeedSync = false

	filter := &metadata.QueryCondition{
		Fields: []string{common.BKSetIDField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField:         bizID,
			common.BKSetTemplateIDField: setTemplateID,
		},
	}
	setResult := metadata.ResponseSetInstance{}
	err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDSet, filter,
		&setResult)
	if err != nil {
		blog.Errorf(" list set failed, option: %s, err: %s, rid: %s", filter, err, kit.Rid)
		return result, errors.CCHttpError
	}
	if ccErr := setResult.CCError(); ccErr != nil {
		blog.Errorf("list set failed, option: %s, response: %s, rid: %s", filter, setResult, kit.Rid)
		return result, ccErr
	}

	if len(setResult.Data.Info) == 0 {
		return result, nil
	}

	var setIDs []int64
	for _, item := range setResult.Data.Info {
		setIDs = append(setIDs, item.SetID)
	}

	needSync, err := st.isSyncRequired(kit, bizID, setTemplateID, setIDs, true)
	if err != nil {
		blog.Errorf("check set whether need sync failed, set: %+v, err: %s, rid: %s",
			setIDs, err.Error(), kit.Rid)
		return result, err
	}

	for _, setID := range setIDs {
		if !result.NeedSync {
			if needSync[setID] {
				result.NeedSync = true
			}
		}

		setStatus := metadata.SetUpdateToDateStatus{
			SetID:    setID,
			NeedSync: needSync[setID],
		}
		result.Sets = append(result.Sets, setStatus)
	}

	return result, nil
}
