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

// Package settemplate TODO
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
	"configcenter/src/common/util"
)

// SetTemplate TODO
type SetTemplate interface {
	DiffSetTplWithInst(kit *rest.Kit, bizID int64, setTemplateID int64, option metadata.DiffSetTplWithInstOption,
		serviceTemplates []metadata.ServiceTemplate) (metadata.SetDiff, errors.CCErrorCoder)
	SetWithDeleteModulesRelation(kit *rest.Kit, bizID int64, setTemplateID int64,
		option metadata.SetWithHostFlagOption) (map[int64][]int64, errors.CCErrorCoder)
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

// NewSetTemplate TODO
func NewSetTemplate(client apimachinery.ClientSetInterface) SetTemplate {
	return &setTemplate{
		client: client,
	}
}

type setTemplate struct {
	client apimachinery.ClientSetInterface
}

func getSetInstanceResult(kit *rest.Kit, setMap mapstr.MapStr) (metadata.SetInst, errors.CCErrorCoder) {

	if setMap == nil {
		return metadata.SetInst{}, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, "set")
	}
	bizID, err := util.GetInt64ByInterface(setMap[common.BKAppIDField])
	if err != nil {
		return metadata.SetInst{}, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, common.BKAppIDField)
	}

	setTemplateID, err := util.GetInt64ByInterface(setMap[common.BKSetTemplateIDField])
	if err != nil {
		return metadata.SetInst{}, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, common.BKSetTemplateIDField)
	}

	setID, err := util.GetInt64ByInterface(setMap[common.BKSetIDField])
	if err != nil {
		return metadata.SetInst{}, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, common.BKSetIDField)
	}

	set := metadata.SetInst{
		BizID:         bizID,
		SetID:         setID,
		SetName:       util.GetStrByInterface(setMap[common.BKSetNameField]),
		SetTemplateID: setTemplateID,
	}
	return set, nil

}

func (st *setTemplate) getModuleResult(kit *rest.Kit, bizID, setTemplateID, setID int64) ([]metadata.ModuleInst,
	errors.CCErrorCoder) {

	filter := &metadata.QueryCondition{
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		Fields: []string{common.BKModuleIDField, common.BKModuleNameField, common.BKServiceTemplateIDField},
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

// getSetTemplateAttrIdAndPropertyValue 获取集群模板的属性id以及对应的属性值
func (st *setTemplate) getSetTemplateAttrIdAndPropertyValue(kit *rest.Kit, bizID, setTemplateID int64) ([]int64,
	map[int64]interface{}, errors.CCErrorCoder) {

	option := &metadata.ListSetTempAttrOption{
		BizID:  bizID,
		ID:     setTemplateID,
		Fields: []string{common.BKAttributeIDField, common.BKPropertyValueField},
	}

	data, err := st.client.CoreService().SetTemplate().ListSetTemplateAttribute(kit.Ctx, kit.Header, option)
	if err != nil {
		blog.Errorf("list set template attributes failed, option: %+v, bizID: %d, set templateID: %d, err: %v, rid: %s",
			option, bizID, setTemplateID, err, kit.Rid)
		return nil, nil, err
	}

	attrIDs := make([]int64, 0)
	setTemplateAttrValueMap := make(map[int64]interface{})
	for _, attr := range data.Attributes {
		attrIDs = append(attrIDs, attr.AttributeID)
		setTemplateAttrValueMap[attr.AttributeID] = attr.PropertyValue
	}

	return attrIDs, setTemplateAttrValueMap, nil
}

// getSetAttrIDAndPropertyID 根据模块属性ID获取对应的propertyID列表以及属性ID与propertyID的对应关系
func (st *setTemplate) getSetAttrIDAndPropertyID(kit *rest.Kit, attrIDs []int64) ([]string, map[int64]string,
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
			metadata.AttributeFieldIsEditable: true,
		},
		DisableCounter: true,
	}

	res, err := st.client.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDSet, option)
	if err != nil {
		blog.Errorf("read set attribute failed, option: %#v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeSelectFailed)
	}

	propertyIDs := make([]string, 0)
	for _, attrs := range res.Info {
		propertyIDs = append(propertyIDs, attrs.PropertyID)
		attrIdPropertyMap[attrs.ID] = attrs.PropertyID
	}

	return propertyIDs, attrIdPropertyMap, nil
}

func (st *setTemplate) getSetAttributesResult(kit *rest.Kit, bizID, setTemplateID, setID int64) (
	[]metadata.AttributeFields, mapstr.MapStr, errors.CCErrorCoder) {

	attrValues := make([]metadata.AttributeFields, 0)

	option := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKAppIDField:         bizID,
			common.BKSetIDField:         setID,
		},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}
	set, err := st.client.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, option)
	if err != nil {
		blog.Errorf("get set failed, option: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return attrValues, nil, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, err.Error())
	}
	if len(set.Info) == 0 {
		return attrValues, nil, nil
	}

	// 1、获取指定集群模板的属性ID及属性值 get set template attributes
	attrIDs, setTemplateAttrValueMap, cErr := st.getSetTemplateAttrIdAndPropertyValue(kit, bizID, setTemplateID)
	if cErr != nil {
		return attrValues, nil, cErr
	}
	if len(attrIDs) == 0 {
		return attrValues, set.Info[0], nil
	}

	// 2、获取集群 attrID 与 propertyID的映射关系
	propertyIDs, attrIdPropertyIdMap, cErr := st.getSetAttrIDAndPropertyID(kit, attrIDs)
	if cErr != nil {
		return attrValues, nil, cErr
	}

	if len(propertyIDs) == 0 {
		return attrValues, set.Info[0], nil
	}

	// 3、根据propertyID 获取对应集群实例的值
	setPropertyValue := make(map[string]interface{})
	for _, propertyID := range propertyIDs {
		if _, ok := set.Info[0][propertyID]; ok {
			setPropertyValue[propertyID] = set.Info[0][propertyID]
		}
	}

	// 4、整理数据
	for id, attr := range setTemplateAttrValueMap {
		attrValues = append(attrValues, metadata.AttributeFields{
			ID:                    id,
			TemplatePropertyValue: attr,
			InstancePropertyValue: setPropertyValue[attrIdPropertyIdMap[id]],
		})
	}

	return attrValues, set.Info[0], nil
}

// SetWithDeleteModulesRelation 获取涉及到的每个集群下面删除的模块列表
func (st *setTemplate) SetWithDeleteModulesRelation(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.SetWithHostFlagOption) (map[int64][]int64, errors.CCErrorCoder) {

	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header, bizID,
		setTemplateID)
	if err != nil {
		blog.Errorf("list service templates failed, bizID: %d, setTemplateID: %d, err: %v, rid: %s",
			bizID, setTemplateID, err, kit.Rid)
		return nil, err
	}

	setIDs := util.IntArrayUnique(option.SetIDs)
	queryCond := []map[string]interface{}{
		{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKSetIDField:         map[string]interface{}{common.BKDBIN: setIDs}},
	}

	rst, err := st.client.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameBaseSet,
		queryCond)
	if err != nil {
		blog.Errorf("get duplicate set count failed, filter: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, err
	}

	if int(rst[0]) != len(setIDs) {
		blog.Errorf("some setIDs invalid, input IDs: %+v, set num: %+v, rid: %s", setIDs, rst[0], kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids")
	}

	op := &metadata.QueryCondition{
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKParentIDField:      map[string]interface{}{common.BKDBIN: option.SetIDs},
		}),
	}

	modules := metadata.ResponseModuleInstance{}
	if err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule, op,
		&modules); err != nil {
		blog.Errorf("list modules failed, bizID: %d, setTemplateID: %d, setIDs: %+v, err: %v, rid: %s", bizID,
			setTemplateID, option.SetIDs, err, kit.Rid)
		return nil, err
	}
	if err := modules.CCError(); err != nil {
		blog.Errorf("list module http reply failed, bizID: %d, setTemplateID: %d, setIDs: %+v, filter: %+v, reply: %s,"+
			" rid: %s", bizID, setTemplateID, option.SetIDs, op, modules, kit.Rid)
		return nil, err
	}

	setModules := make(map[int64][]metadata.ModuleInst)
	// init before modules loop so that set with no modules could be initial correctly
	for _, setID := range option.SetIDs {
		setModules[setID] = make([]metadata.ModuleInst, 0)
	}

	for _, module := range modules.Data.Info {
		if _, exist := setModules[module.ParentID]; !exist {
			setModules[module.ParentID] = make([]metadata.ModuleInst, 0)
		}
		setModules[module.ParentID] = append(setModules[module.ParentID], module)
	}

	setModuleResults := make(map[int64][]int64, 0)
	for setID, mods := range setModules {
		moduleIDs := getDeleteModulesWithServiceTemplate(serviceTemplates, mods)
		setModuleResults[setID] = moduleIDs
	}
	return setModuleResults, nil
}

// DiffSetTplWithInst TODO
func (st *setTemplate) DiffSetTplWithInst(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.DiffSetTplWithInstOption, serviceTemplates []metadata.ServiceTemplate) (
	metadata.SetDiff, errors.CCErrorCoder) {

	attrs, setMap, ccErr := st.getSetAttributesResult(kit, bizID, setTemplateID, option.SetID)
	if ccErr != nil {
		blog.Errorf("get set attrs failed, bizID: %d, set template id: %d, setID: %d, err: %v, rid: %s", bizID,
			setTemplateID, option.SetID, ccErr, kit.Rid)
		return metadata.SetDiff{}, ccErr
	}

	set, err := getSetInstanceResult(kit, setMap)
	if err != nil {
		blog.Errorf("get set info failed, bizID: %d, setTemplateID: %d, setID: %d, setMap: %+v, err: %v, rid: %s",
			bizID, setTemplateID, option.SetID, setMap, err, kit.Rid)
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

	// diff service template and modules
	moduleDiff := DiffServiceTemplateWithModules(serviceTemplates, modules)
	setDiff := metadata.SetDiff{
		ModuleDiffs: moduleDiff,
		SetID:       option.SetID,
	}
	setDiff.SetDetail = set
	setDiff.Attributes = attrs

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

// SyncSetTplToInst TODO
func (st *setTemplate) SyncSetTplToInst(kit *rest.Kit, bizID int64, setTemplateID int64,
	option metadata.SyncSetTplToInstOption) errors.CCErrorCoder {

	var (
		wg       sync.WaitGroup
		firstErr errors.CCErrorCoder
	)

	pipeline := make(chan bool, 10)
	setDiffs := make([]metadata.SetDiff, 0)

	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header, bizID,
		setTemplateID)
	if err != nil {
		blog.Errorf("list service templates failed, bizID: %d, setTemplateID: %d, err: %v, rid: %s", bizID,
			setTemplateID, err, kit.Rid)
		return err
	}

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
			setDiff, err := st.DiffSetTplWithInst(kit, bizID, setTemplateID, option, serviceTemplates)
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

// getDeleteModulesWithServiceTemplate 返回的是集群和module列表的关系
func getDeleteModulesWithServiceTemplate(serviceTemplates []metadata.ServiceTemplate,
	modules []metadata.ModuleInst) []int64 {

	svcTplMap := make(map[int64]metadata.ServiceTemplate, len(serviceTemplates))
	for _, svcTpl := range serviceTemplates {
		svcTplMap[svcTpl.ID] = svcTpl
	}
	moduleIDs := make([]int64, 0)
	for _, module := range modules {
		_, ok := svcTplMap[module.ServiceTemplateID]
		if !ok {
			moduleIDs = append(moduleIDs, module.ModuleID)
		}
	}

	return moduleIDs
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

// CheckSetInstUpdateToDateStatus TODO
// CheckSetTplInstLatest 检查通过集群模板 setTemplateID 实例化的集群是否都已经达到最新状态
func (st *setTemplate) CheckSetInstUpdateToDateStatus(kit *rest.Kit, bizID int64,
	setTemplateID int64) (*metadata.SetTemplateUpdateToDateStatus, errors.CCErrorCoder) {

	result := new(metadata.SetTemplateUpdateToDateStatus)
	result.SetTemplateID = setTemplateID
	result.NeedSync = false

	// 获取指定集群模板的属性ID及属性值
	attrIDs, setTemplateAttrValueMap, cErr := st.getSetTemplateAttrIdAndPropertyValue(kit, bizID, setTemplateID)
	if cErr != nil {
		return &metadata.SetTemplateUpdateToDateStatus{}, cErr
	}

	// 获取集群 attrID 与 propertyID的映射关系
	propertyIDs, attrIdPropertyIdMap, cErr := st.getSetAttrIDAndPropertyID(kit, attrIDs)
	if cErr != nil {
		return &metadata.SetTemplateUpdateToDateStatus{}, cErr
	}
	if len(propertyIDs) == 0 {
		return &metadata.SetTemplateUpdateToDateStatus{}, nil
	}
	fields := []string{common.BKSetIDField}
	fields = append(fields, propertyIDs...)

	// get set detail
	option := &metadata.QueryCondition{
		Fields: fields,
		Condition: map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKAppIDField:         bizID,
		},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}
	set, err := st.client.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, option)
	if err != nil {
		blog.Errorf("get set failed, option: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, err.Error())
	}

	if len(set.Info) == 0 {
		return result, nil
	}

	setIDs := make([]int64, 0)
	setMap := make(map[int64]mapstr.MapStr)

	for _, set := range set.Info {
		setID, err := util.GetInt64ByInterface(set[common.BKSetIDField])
		if err != nil {
			return result, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed)
		}
		setIDs = append(setIDs, setID)
		setMap[setID] = set
	}

	needSync, ccErr := st.isSyncRequired(kit, bizID, setTemplateID, setIDs, setMap, true, attrIdPropertyIdMap,
		setTemplateAttrValueMap)
	if ccErr != nil {
		blog.Errorf("check set whether need sync failed, setIDs: %+v, err: %v, rid: %s", setIDs, ccErr, kit.Rid)
		return result, ccErr
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
