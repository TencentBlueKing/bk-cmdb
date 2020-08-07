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
	"context"
	"net/http"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type SetTemplate interface {
	DiffSetTplWithInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) ([]metadata.SetDiff, errors.CCErrorCoder)
	SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64, option metadata.SyncSetTplToInstOption) errors.CCErrorCoder
	UpdateSetSyncStatus(kit *rest.Kit, setID int64) (metadata.SetTemplateSyncStatus, errors.CCErrorCoder)
	GetLatestSyncTaskDetail(kit *rest.Kit, setID int64) (*metadata.APITaskDetail, errors.CCErrorCoder)
	CheckSetInstUpdateToDateStatus(kit *rest.Kit, bizID int64, setTemplateID int64) (metadata.SetTemplateUpdateToDateStatus, errors.CCErrorCoder)
	TriggerCheckSetTemplateSyncingStatus(kit *rest.Kit, bizID, setTemplateID, setID int64) errors.CCErrorCoder
}

func NewSetTemplate(client apimachinery.ClientSetInterface) SetTemplate {
	return &setTemplate{
		client: client,
	}
}

type setTemplate struct {
	client apimachinery.ClientSetInterface
}

func (st *setTemplate) DiffSetTplWithInst(ctx context.Context, header http.Header, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption) ([]metadata.SetDiff, errors.CCErrorCoder) {
	rid := util.GetHTTPCCRequestID(header)

	ccError := util.GetDefaultCCError(header)
	if ccError == nil {
		return nil, errors.GlobalCCErrorNotInitialized
	}

	setTemplate, err := st.client.CoreService().SetTemplate().GetSetTemplate(ctx, header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, GetSetTemplate failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}

	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(ctx, header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, ListSetTplRelatedSvcTpl failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}
	serviceTemplateMap := make(map[int64]metadata.ServiceTemplate)
	for _, svcTpl := range serviceTemplates {
		serviceTemplateMap[svcTpl.ID] = svcTpl
	}

	setIDs := util.IntArrayUnique(option.SetIDs)
	setFilter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKSetIDField: map[string]interface{}{
				common.BKDBIN: setIDs,
			},
		}),
	}
	setInstResult, e := st.client.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDSet, setFilter)
	if e != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, list sets failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %s, rid: %s", bizID, setTemplateID, option.SetIDs, e.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}
	if len(setInstResult.Data.Info) != len(setIDs) {
		blog.Errorf("DiffSetTemplateWithInstances failed, some setID invalid, input IDs: %+v, valid IDs: %+v, rid: %s", setIDs, setInstResult.Data.Info, rid)
		return nil, ccError.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids")
	}
	setMap := make(map[int64]metadata.SetInst)
	for _, setInstance := range setInstResult.Data.Info {
		set := metadata.SetInst{}
		if err := mapstruct.Decode2Struct(setInstance, &set); err != nil {
			blog.Errorf("DiffSetTemplateWithInstances failed, decode set instance failed, set: %+v, err: %s, rid: %s", setInstance, err.Error(), rid)
			return nil, ccError.CCError(common.CCErrCommJSONMarshalFailed)
		}
		if set.SetID == 0 {
			blog.Errorf("DiffSetTemplateWithInstances failed, decode set instance result setID=0, data: %+v, rid: %s", setInstance, rid)
			return nil, ccError.CCError(common.CCErrCommJSONMarshalFailed)
		}
		setMap[set.SetID] = set
	}

	moduleFilter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKParentIDField: map[string]interface{}{
				common.BKDBIN: option.SetIDs,
			},
		}),
	}
	modulesInstResult, e := st.client.CoreService().Instance().ReadInstance(ctx, header, common.BKInnerObjIDModule, moduleFilter)
	if e != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, list modules failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %s, rid: %s", bizID, setTemplateID, option.SetIDs, e.Error(), rid)
		return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
	}

	setModules := make(map[int64][]metadata.ModuleInst)
	// init before modules loop so that set with no modules could be initial correctly
	for _, setID := range option.SetIDs {
		setModules[setID] = make([]metadata.ModuleInst, 0)
	}
	for _, moduleInstance := range modulesInstResult.Data.Info {
		module := metadata.ModuleInst{}
		if err := mapstruct.Decode2Struct(moduleInstance, &module); err != nil {
			blog.Errorf("DiffSetTemplateWithInstances failed, decode module instance failed, module: %+v, err: %s, rid: %s", moduleInstance, err.Error(), rid)
			return nil, ccError.CCError(common.CCErrCommDBSelectFailed)
		}
		if _, exist := setModules[module.ParentID]; !exist {
			setModules[module.ParentID] = make([]metadata.ModuleInst, 0)
		}
		setModules[module.ParentID] = append(setModules[module.ParentID], module)
	}

	topoTree, ccErr := st.client.CoreService().Mainline().SearchMainlineInstanceTopo(ctx, header, bizID, false)
	if ccErr != nil {
		blog.Errorf("ListSetTplRelatedSetsWeb failed, bizID: %d, err: %s, rid: %s", bizID, ccErr.Error(), rid)
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
		setDiff.SetTemplateVersion = setTemplate.Version
		setDiff.UpdateNeedSyncField()
		setDiffs = append(setDiffs, setDiff)
	}
	return setDiffs, nil
}

func (st *setTemplate) SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64, option metadata.SyncSetTplToInstOption) errors.CCErrorCoder {
	rid := util.GetHTTPCCRequestID(kit.Header)

	diffOption := metadata.DiffSetTplWithInstOption{
		SetIDs: option.SetIDs,
	}
	setDiffs, err := st.DiffSetTplWithInst(kit.Ctx, kit.Header, bizID, setTemplateID, diffOption)
	if err != nil {
		return err
	}

	for _, setDiff := range setDiffs {
		indexKey := metadata.GetSetTemplateSyncIndex(setDiff.SetID)
		blog.V(3).Infof("dispatch synchronize task on set [%s](%d), rid: %s", setDiff.SetDetail.SetName, setDiff.SetID, rid)
		tasks := make([]metadata.SyncModuleTask, 0)
		for _, moduleDiff := range setDiff.ModuleDiffs {
			task := metadata.SyncModuleTask{
				Set:                setDiff.SetDetail,
				ModuleDiff:         moduleDiff,
				SetTopoPath:        setDiff.TopoPath,
				SetTemplateVersion: setDiff.SetTemplateVersion,
			}
			tasks = append(tasks, task)
		}
		taskDetail, err := st.DispatchTask4ModuleSync(kit.Ctx, kit.Header, indexKey, tasks...)
		if err != nil {
			return err
		}
		if blog.V(3) {
			blog.InfoJSON("dispatch synchronize task on set [%s](%s) success, result: %s, rid: %s", setDiff.SetDetail.SetName, setDiff.SetID, taskDetail, rid)
		}

		// 修改任务到同步中的状态
		err = st.client.CoreService().SetTemplate().ModifySetTemplateSyncStatus(kit.Ctx, kit.Header, setDiff.SetID, metadata.SyncStatusSyncing)
		//_, err = st.UpdateSetSyncStatus(kit, setDiff.SetID)
		if err != nil {
			blog.Errorf("UpdateSetSyncStatus failed, setID: %d, err: %s", setDiff.SetID, err.Error())
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
					blog.Errorf("poll UpdateSetSyncStatus timeout, setID: %d", setID)
					return
				case <-ticker.C:
					setSyncStatus, err := st.UpdateSetSyncStatus(kit.NewKit(), setID)
					if err != nil {
						blog.Errorf("UpdateSetSyncStatus failed, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
						return
					}
					if setSyncStatus.Status.IsFinished() {
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

func (st *setTemplate) DispatchTask4ModuleSync(ctx context.Context, header http.Header, indexKey string, tasks ...metadata.SyncModuleTask) (metadata.APITaskDetail, errors.CCErrorCoder) {
	taskDetail := metadata.APITaskDetail{}
	rid := util.GetHTTPCCRequestID(header)
	tasksData := make([]interface{}, 0)
	for _, task := range tasks {
		tasksData = append(tasksData, task)
	}
	createTaskResult, err := st.client.TaskServer().Task().Create(ctx, header, common.SyncSetTaskName, indexKey, tasksData)
	if err != nil {
		blog.ErrorJSON("dispatch synchronize task failed, task: %s, err: %s, rid: %s", tasks, err.Error(), rid)
		return taskDetail, errors.CCHttpError
	}
	if ccErr := createTaskResult.CCError(); ccErr != nil {
		blog.ErrorJSON("dispatch synchronize task failed, task: %s, result: %s, rid: %s", tasks, createTaskResult, rid)
		return taskDetail, ccErr
	}
	blog.InfoJSON("dispatch synchronize task success, task: %s, create result: %s, rid: %s", tasks, createTaskResult, rid)
	taskDetail = createTaskResult.Data
	return taskDetail, nil
}

// DiffServiceTemplateWithModules diff modules with template in one set
func DiffServiceTemplateWithModules(serviceTemplates []metadata.ServiceTemplate, modules []metadata.ModuleInst) []metadata.SetModuleDiff {
	svcTplMap := make(map[int64]metadata.ServiceTemplate)
	svcTplHitMap := make(map[int64]bool)
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
func (st *setTemplate) CheckSetInstUpdateToDateStatus(kit *rest.Kit, bizID int64, setTemplateID int64) (metadata.SetTemplateUpdateToDateStatus, errors.CCErrorCoder) {
	rid := kit.Rid

	result := metadata.SetTemplateUpdateToDateStatus{}
	setTemplate, ccErr := st.client.CoreService().SetTemplate().GetSetTemplate(kit.Ctx, kit.Header, bizID, setTemplateID)
	if ccErr != nil {
		blog.Errorf("CheckSetInstUpdateToDateStatus failed, GetSetTemplate failed, bizID: %d, setTemplateID: %d, err: %+v, rid: %s", bizID, setTemplateID, ccErr, rid)
		return result, ccErr
	}
	result.SetTemplateVersion = setTemplate.Version
	result.SetTemplateID = setTemplateID
	result.NeedSync = false

	filter := &metadata.QueryCondition{
		Fields: []string{common.BKSetIDField, common.BKSetTemplateVersionField},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Condition: map[string]interface{}{
			common.BKAppIDField:         bizID,
			common.BKSetTemplateIDField: setTemplateID,
		},
	}
	setResult, err := st.client.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, filter)
	if err != nil {
		blog.Errorf("CheckSetInstUpdateToDateStatus failed, ReadInstance of set failed, option: %+v, err: %+v, rid: %s", filter, err, rid)
		return result, errors.CCHttpError
	}
	if ccErr := setResult.CCError(); ccErr != nil {
		blog.Errorf("CheckSetInstUpdateToDateStatus failed, ReadInstance of set failed, option: %+v, response: %+v, rid: %s", filter, setResult, rid)
		return result, ccErr
	}

	for _, item := range setResult.Data.Info {
		set := metadata.SetInst{}
		if err := mapstruct.Decode2Struct(item, &set); err != nil {
			blog.ErrorJSON("CheckSetInstUpdateToDateStatus failed, unmarshal set data failed, set: %s, err: %s, rid: %s", item, err.Error(), rid)
			return result, kit.CCError.CCError(common.CCErrCommParseDBFailed)
		}
		needSync := set.SetTemplateVersion != setTemplate.Version
		setStatus := metadata.SetUpdateToDateStatus{
			SetID:              set.SetID,
			SetTemplateVersion: set.SetTemplateVersion,
			NeedSync:           needSync,
		}
		if needSync {
			result.NeedSync = true
		}
		result.Sets = append(result.Sets, setStatus)
	}

	return result, nil
}
