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
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type SetTemplate interface {
	DiffSetTplWithInst(kit *rest.Kit, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) (
		[]metadata.SetDiff, errors.CCErrorCoder)
	SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64,
		option metadata.SyncSetTplToInstOption) errors.CCErrorCoder
	UpdateSetSyncStatus(kit *rest.Kit, setTemplateID int64, setID []int64) ([]metadata.SetTemplateSyncStatus,
		errors.CCErrorCoder)
	GetLatestSyncTaskDetail(kit *rest.Kit, taskCond metadata.ListAPITaskDetail) (map[int64]*metadata.APITaskDetail,
		errors.CCErrorCoder)
	CheckSetInstUpdateToDateStatus(kit *rest.Kit, bizID int64, setTemplateID int64) (
		*metadata.SetTemplateUpdateToDateStatus, errors.CCErrorCoder)
	TriggerCheckSetTemplateSyncingStatus(kit *rest.Kit, bizID, setTemplateID int64, setID []int64) errors.CCErrorCoder
	ListSetTemplateSyncStatus(kit *rest.Kit, bizID int64, option metadata.ListSetTemplateSyncStatusOption) (
		metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder)
}

func NewSetTemplate(client apimachinery.ClientSetInterface) SetTemplate {
	return &setTemplate{
		client: client,
	}
}

type setTemplate struct {
	client apimachinery.ClientSetInterface
}

// DiffSetTplWithInst diff set template with instances
func (st *setTemplate) DiffSetTplWithInst(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.DiffSetTplWithInstOption) ([]metadata.SetDiff, errors.CCErrorCoder) {

	// get set template related service template
	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header,
		bizID, setTemplateID)
	if err != nil {
		blog.Errorf("list set tpl related svc tpl failed, bizID: %d, setTemplateID: %d, err: %v, rid: %s", bizID,
			setTemplateID, err, kit.Rid)
		return nil, err
	}
	serviceTemplateMap := make(map[int64]metadata.ServiceTemplate)
	for _, svcTpl := range serviceTemplates {
		serviceTemplateMap[svcTpl.ID] = svcTpl
	}

	// get set template related set and modules
	setMap, setModules, err := st.getSetTplRelatedInst(kit, bizID, setTemplateID, option.SetIDs)
	if err != nil {
		return nil, err
	}

	// get set topo path
	topoTree, ccErr := st.client.CoreService().Mainline().SearchMainlineInstanceTopo(kit.Ctx, kit.Header, bizID, false)
	if ccErr != nil {
		blog.Errorf("list set tpl related sets for web failed, bizID: %d, err: %s, rid: %s", bizID, ccErr, kit.Rid)
		return nil, ccErr
	}

	// diff
	setDiffs := make([]metadata.SetDiff, 0)
	for setID, modules := range setModules {
		moduleDiff := DiffServiceTemplateWithModules(serviceTemplates, modules)
		setDiff := metadata.SetDiff{
			ModuleDiffs: moduleDiff,
			SetID:       setID,
		}
		if set, ok := setMap[setID]; ok {
			setDiff.SetDetail = set
		}

		// add topo path info
		setPath := topoTree.TraversalFindNode(common.BKInnerObjIDSet, setID)
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
		setDiffs = append(setDiffs, setDiff)
	}
	return setDiffs, nil
}

func (st *setTemplate) getSetTplRelatedInst(kit *rest.Kit, bizID, setTemplateID int64, setIDs []int64) (
	map[int64]metadata.SetInst, map[int64][]metadata.ModuleInst, errors.CCErrorCoder) {

	// get set template related sets
	setIDs = util.IntArrayUnique(setIDs)
	sets, err := st.getSets(kit, setTemplateID, setIDs)
	if err != nil {
		blog.Errorf("get sets failed, setIDs: %+v, err: %v, rid: %s", setIDs, err, kit.Rid)
		return nil, nil, err
	}

	setMap := make(map[int64]metadata.SetInst)
	for _, set := range sets {
		if set.SetID == 0 {
			blog.Errorf("decode set instance result setID=0, data: %+v, rid: %s", set, kit.Rid)
			return nil, nil, kit.CCError.CCError(common.CCErrCommJSONMarshalFailed)
		}
		setMap[set.SetID] = set
	}

	// get set template related modules
	moduleFilter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKParentIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
		}),
	}
	modulesInstResult := metadata.ResponseModuleInstance{}
	if err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		moduleFilter, &modulesInstResult); err != nil {
		blog.Errorf("list modules failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %v, rid: %s", bizID,
			setTemplateID, setIDs, err, kit.Rid)
		return nil, nil, err
	}
	if err := modulesInstResult.CCError(); err != nil {
		blog.Errorf("list modules failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %v, rid: %s", bizID,
			setTemplateID, setIDs, err, kit.Rid)
		return nil, nil, err
	}

	// get set to modules map
	setModules := make(map[int64][]metadata.ModuleInst)
	// init before modules loop so that set with no modules could be initial correctly
	for _, setID := range setIDs {
		setModules[setID] = make([]metadata.ModuleInst, 0)
	}
	for _, module := range modulesInstResult.Data.Info {
		if _, exist := setModules[module.ParentID]; !exist {
			setModules[module.ParentID] = make([]metadata.ModuleInst, 0)
		}
		setModules[module.ParentID] = append(setModules[module.ParentID], module)
	}

	return setMap, setModules, nil
}

// SyncSetTplToInst sync set template to instances
func (st *setTemplate) SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.SyncSetTplToInstOption) errors.CCErrorCoder {

	diffOption := metadata.DiffSetTplWithInstOption{
		SetIDs: option.SetIDs,
	}
	setDiffs, err := st.DiffSetTplWithInst(kit, bizID, setTemplateID, diffOption)
	if err != nil {
		return err
	}

	for _, setDiff := range setDiffs {
		blog.V(3).Infof("dispatch synchronize task on set[%#v], rid: %s", setDiff, kit.Rid)
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
		blog.V(3).Infof("dispatch sync task on set[%#v]) success, result: %#v, rid: %s", setDiff, taskDetail, kit.Rid)

		// 修改任务到同步中的状态
		if err = st.client.CoreService().SetTemplate().ModifySetTemplateSyncStatus(kit.Ctx, kit.Header, setDiff.SetID,
			metadata.SyncStatusSyncing); err != nil {
			blog.Errorf("update set sync status failed, setID: %d, err: %v, rid: %s", setDiff.SetID, err, kit.Rid)
			return err
		}

		// 定时更新 SetTemplateSyncStatus 状态，优化加载
		go func(setID int64) {
			// 指数增长轮询间隔
			duration := 200 * time.Millisecond
			maxDuration := 10 * time.Second
			ticker := time.NewTimer(duration)
			timeoutTimer := time.NewTimer(5 * time.Minute)

			for {
				select {
				case <-timeoutTimer.C:
					blog.Errorf("poll set sync status timeout, setID: %d, rid: %s", setID, kit.Rid)
					return
				case <-ticker.C:
					setSyncStatus, err := st.UpdateSetSyncStatus(kit.NewKit(), setTemplateID, []int64{setID})
					if err != nil {
						blog.Errorf("update set sync status failed, setID: %d, err: %s, rid: %s", setID, err, kit.Rid)
						return
					}
					if len(setSyncStatus) == 0 {
						blog.Errorf("update set sync status returns empty, setID: %d, err: %v, rid: %s", setID, err,
							kit.Rid)
						return
					}
					if setSyncStatus[0].Status.IsFinished() {
						return
					}

					// set next timer
					duration = duration * 2
					if duration > maxDuration {
						duration = maxDuration
					}
					ticker = time.NewTimer(duration)
				}
			}
		}(setDiff.SetID)
	}
	return nil
}

// DispatchTask4ModuleSync dispatch synchronize task for module sync
func (st *setTemplate) DispatchTask4ModuleSync(kit *rest.Kit, indexKey string, setID int64,
	tasks ...metadata.SyncModuleTask) (metadata.APITaskDetail, errors.CCErrorCoder) {

	taskDetail := metadata.APITaskDetail{}
	tasksData := make([]interface{}, 0)
	for _, task := range tasks {
		tasksData = append(tasksData, task)
	}

	createTaskResult, err := st.client.TaskServer().Task().Create(kit.Ctx, kit.Header, common.SyncSetTaskName,
		indexKey, setID, tasksData)
	if err != nil {
		blog.ErrorJSON("dispatch synchronize task failed, task: %s, err: %s, rid: %s", tasks, err, kit.Rid)
		return taskDetail, errors.CCHttpError
	}
	if ccErr := createTaskResult.CCError(); ccErr != nil {
		blog.ErrorJSON("dispatch synchronize task failed, task: %s, err: %s, rid: %s", tasks, ccErr, kit.Rid)
		return taskDetail, ccErr
	}

	blog.InfoJSON("dispatch synchronize task success, task: %s, result: %s, rid: %s", tasks, createTaskResult, kit.Rid)
	return createTaskResult.Data, nil
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

// CheckSetInstUpdateToDateStatus 检查通过集群模板 setTemplateID 实例化的集群是否都已经达到最新状态
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
		blog.Errorf("list set failed, option: %#v, err: %v, rid: %s", filter, err, kit.Rid)
		return result, errors.CCHttpError
	}
	if ccErr := setResult.CCError(); ccErr != nil {
		blog.Errorf("list set failed, option: %#v, err: %v, rid: %s", filter, ccErr, kit.Rid)
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
		blog.Errorf("check if set need sync failed, setIDs: %+v, err: %v, rid: %s", setIDs, err, kit.Rid)
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
