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

package settemplate

import (
	"fmt"
	"reflect"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// GetSets TODO
func (st *setTemplate) GetSets(kit *rest.Kit, setTemplateID int64, setIDs []int64) ([]metadata.SetInst,
	errors.CCErrorCoder) {

	filter := &metadata.QueryCondition{}
	filter.Condition = mapstr.MapStr{
		common.BKSetIDField:         map[string]interface{}{common.BKDBIN: setIDs},
		common.BKSetTemplateIDField: setTemplateID,
	}

	instResult := new(metadata.ResponseSetInstance)
	err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDSet, filter,
		instResult)
	if err != nil {
		blog.Errorf("GetSets failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if ccErr := instResult.CCError(); ccErr != nil {
		blog.Errorf("GetSets failed, read instance failed, filter: %s, instResult: %s, err: %s, rid: %s",
			filter, instResult, ccErr.Error(), kit.Rid)
		return nil, ccErr
	}

	if len(instResult.Data.Info) == 0 {
		blog.Errorf("GetSets failed, set not found, filter: %s, instResult: %s, rid: %s", filter, instResult, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if instResult.Data.Count != len(setIDs) {
		blog.Errorf("GetSets failed, some setID invalid, input IDs: %+v, valid ,IDs: %+v, rid: %s",
			setIDs, instResult.Data.Info, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids")
	}

	return instResult.Data.Info, nil
}

// isSyncRequired Note: If the parameter isInterrupt is true, it will return if a state to be synchronized is found.
// At this time, the rest of the cluster state will be set to synchronized by default. If you need to return all pending
// synchronization status state setId, you need to set this parameter to false.
func (st *setTemplate) isSyncRequired(kit *rest.Kit, bizID, setTemplateID int64, setIDs []int64,
	setMap map[int64]mapstr.MapStr, isInterrupt bool, attrIdPropertyIdMap map[int64]string,
	setTemplateAttrValueMap map[int64]interface{}) (map[int64]bool, errors.CCErrorCoder) {

	if len(setIDs) == 0 {
		blog.Errorf("array of set_id is empty, rid: %s", kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKSetIDField)
	}

	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header, bizID,
		setTemplateID)
	if err != nil {
		blog.Errorf(" list set template and service template related failed, bizID: %d, "+
			"setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), kit.Rid)
		return nil, err
	}

	svcTplCnt := int64(len(serviceTemplates))
	svcTplMap := make(map[int64]metadata.ServiceTemplate, svcTplCnt)
	for _, serviceTemplate := range serviceTemplates {
		svcTplMap[serviceTemplate.ID] = serviceTemplate
	}

	moduleFilter := &metadata.QueryCondition{
		Page: metadata.BasePage{Limit: common.BKNoLimit},
		Fields: []string{common.BKSetIDField, common.BKModuleIDField, common.BKSetTemplateIDField,
			common.BKModuleNameField, common.BKServiceTemplateIDField},
		Condition: mapstr.MapStr(map[string]interface{}{
			common.BKSetTemplateIDField: setTemplateID,
			common.BKSetIDField:         map[string]interface{}{common.BKDBIN: setIDs}}),
	}

	modulesInstResult := new(metadata.ResponseModuleInstance)
	if err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDModule,
		moduleFilter, modulesInstResult); err != nil {
		blog.Errorf("list modules failed, bizID: %s, setTemplateID: %s, setIDs: %s, err: %s, rid: %s",
			bizID, setTemplateID, setIDs, err, kit.Rid)
		return nil, err
	}

	if err := modulesInstResult.CCError(); err != nil {
		blog.Errorf("list module http reply failed, bizID: %s, setTemplateID: %s, setIDs: %s, filter: %s, "+
			"reply: %s, rid: %s", bizID, setTemplateID, setIDs, moduleFilter, modulesInstResult, kit.Rid)
		return nil, err
	}

	setModules := make(map[int64][]metadata.ModuleInst, len(modulesInstResult.Data.Info))
	for _, module := range modulesInstResult.Data.Info {
		if _, exist := setModules[module.SetID]; !exist {
			setModules[module.SetID] = make([]metadata.ModuleInst, 0)
		}
		setModules[module.SetID] = append(setModules[module.SetID], module)
	}

	checkResult := make(map[int64]bool, len(setIDs))
	for _, setID := range setIDs {
		module := setModules[setID]
		checkResult[setID] = diffModuleServiceTplAndAttrs(svcTplCnt, int64(len(module)), svcTplMap, module,
			attrIdPropertyIdMap, setMap[setID], setTemplateAttrValueMap)
		if isInterrupt && checkResult[setID] {
			return checkResult, nil
		}
	}
	return checkResult, nil
}

// diffModuleServiceTplAndAttrs check different of modules with template in one set
func diffModuleServiceTplAndAttrs(serviceTplCnt, moduleCnt int64, serviceTemplates map[int64]metadata.ServiceTemplate,
	modules []metadata.ModuleInst, attrIdPropertyIdMap map[int64]string, setMap mapstr.MapStr,
	setTemplateAttrValueMap map[int64]interface{}) bool {

	if serviceTplCnt != moduleCnt {
		return true
	}

	// 对比集群模板与集群属性值是否有差异
	for setTemplateAttrID, value := range setTemplateAttrValueMap {
		if !reflect.DeepEqual(value, setMap[attrIdPropertyIdMap[setTemplateAttrID]]) {
			return true
		}
	}
	/*
		depend on logic in func DiffServiceTemplateWithModules
		if the number of the module and the template is not the same, it changed
		if the name of the module and the template is not the same, it changed
		this function only use to check if module and template are the same
	*/

	for _, module := range modules {
		template, ok := serviceTemplates[module.ServiceTemplateID]
		if !ok {
			return true
		}
		if template.Name != module.ModuleName {
			return true
		}
	}

	return false
}

// GetLatestSyncTaskDetail TODO
func (st *setTemplate) GetLatestSyncTaskDetail(kit *rest.Kit,
	taskCond metadata.ListAPITaskDetail) (map[int64]*metadata.APITaskDetail, errors.CCErrorCoder) {

	if len(taskCond.InstID) == 0 {
		blog.Errorf("set id is empty, rid: %s", kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTaskListTaskFail)
	}

	latestTaskResult := make(map[int64]*metadata.APITaskDetail)

	setRelatedTaskFilter := map[string]interface{}{
		common.BKInstIDField:   map[string]interface{}{common.BKDBIN: taskCond.InstID},
		common.BKTaskTypeField: common.SyncSetTaskFlag,
	}
	listTaskOption := new(metadata.ListAPITaskLatestRequest)
	listTaskOption.Condition = setRelatedTaskFilter
	listTaskOption.Fields = taskCond.Fields

	listResult, err := st.client.TaskServer().Task().ListLatestTask(kit.Ctx, kit.Header, common.SyncSetTaskFlag,
		listTaskOption)
	if err != nil {
		blog.Errorf("list set sync tasks failed, option: %s, err: %v, rid: %s", listTaskOption, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTaskListTaskFail)
	}

	if len(listResult) == 0 {
		return latestTaskResult, nil
	}

	for _, APITask := range listResult {
		if len(taskCond.Fields) == 0 {
			clearSetSyncTaskDetail(&APITask)
		}

		if APITask.InstID != 0 {
			latestTaskResult[APITask.InstID] = &APITask
		}
	}

	return latestTaskResult, nil
}

func clearSetSyncTaskDetail(detail *metadata.APITaskDetail) {
	detail.Header = util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)
	for taskIdx := range detail.Detail {
		subTaskDetail, ok := detail.Detail[taskIdx].Data.(map[string]interface{})
		if !ok {
			blog.Warnf("clearSetSyncTaskDetail expect map[string]interface{}, got unexpected type, data: %+v",
				detail.Detail[taskIdx].Data)

			detail.Detail[taskIdx].Data = nil
		}
		delete(subTaskDetail, "header")
	}
}

func (st *setTemplate) getSetMapStrByOption(kit *rest.Kit, option *metadata.ListSetTemplateSyncStatusOption,
	fields []string) ([]mapstr.MapStr, errors.CCErrorCoder) {

	filter := &metadata.QueryCondition{
		Fields: fields,
		Condition: map[string]interface{}{
			common.BKSetTemplateIDField: option.SetTemplateID,
			common.BKAppIDField:         option.BizID,
		},
		Page:           option.Page,
		DisableCounter: true,
	}

	if len(option.SetIDs) > 0 {
		filter.Condition[common.BKSetIDField] = map[string]interface{}{
			common.BKDBIN: option.SetIDs,
		}
	}

	if len(option.SearchKey) > 0 {
		filter.Condition[common.BKSetNameField] = map[string]interface{}{
			common.BKDBLIKE:    fmt.Sprintf(".*%s.*", option.SearchKey),
			common.BKDBOPTIONS: "i",
		}
	}
	set, err := st.client.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, filter)
	if err != nil {
		blog.Errorf("get set failed, option: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed, err.Error())
	}

	return set.Info, nil
}

// ListSetTemplateSyncStatus batch search set template sync status
func (st *setTemplate) ListSetTemplateSyncStatus(kit *rest.Kit, option *metadata.ListSetTemplateSyncStatusOption) (
	*metadata.ListAPITaskSyncStatusResult, errors.CCErrorCoder) {

	// validate option
	err := option.Validate(kit.CCError)
	if err != nil {
		blog.Errorf("parse set condition failed, err: %v, cond: %#v, rid: %s", err, option, kit.Rid)
		return nil, err
	}

	// 获取指定集群模板的属性ID及属性值
	attrIDs, setTemplateAttrValueMap, cErr := st.getSetTemplateAttrIdAndPropertyValue(kit, option.BizID,
		option.SetTemplateID)
	if cErr != nil {
		return nil, cErr
	}

	// 获取集群 attrID 与 propertyID的映射关系
	propertyIDs, attrIdPropertyIdMap, cErr := st.getSetAttrIDAndPropertyID(kit, attrIDs)
	if cErr != nil {
		return nil, cErr
	}

	fields := []string{common.BKSetIDField}
	if len(propertyIDs) > 0 {
		fields = append(fields, propertyIDs...)
	}

	sets, err := st.getSetMapStrByOption(kit, option, fields)
	if err != nil {
		blog.Errorf("list set failed, option: %+v, err: %v, rid: %s", option, err, kit.Rid)
		return nil, err
	}

	if len(sets) == 0 {
		return &metadata.ListAPITaskSyncStatusResult{Count: 0, Info: make([]metadata.APITaskSyncStatus, 0)}, nil
	}

	setIDs := make([]int64, 0)
	setMap := make(map[int64]mapstr.MapStr)

	for _, set := range sets {
		setID, err := util.GetInt64ByInterface(set[common.BKSetIDField])
		if err != nil {
			return nil, kit.CCError.CCErrorf(common.CCErrTopoSetSelectFailed)
		}
		setIDs = append(setIDs, setID)
		setMap[setID] = set
	}

	// get latest sync set template api task sync status by sets
	option.SetIDs = setIDs
	statusCond, err := option.ToStatusCond(kit.CCError)
	if err != nil {
		blog.Errorf("parse status condition failed, err: %v, cond: %#v, rid: %s", err, option, kit.Rid)
		return nil, err
	}

	statusOpt := &metadata.ListLatestSyncStatusRequest{
		Condition:     statusCond.Condition,
		Fields:        statusCond.Fields,
		TimeCondition: statusCond.TimeCondition,
	}

	taskStatusRes, err := st.client.TaskServer().Task().ListLatestSyncStatus(kit.Ctx, kit.Header, statusOpt)
	if err != nil {
		blog.Errorf("list latest sync status failed, option: %#v, err: %v, rid: %s", statusOpt, err, kit.Rid)
		return nil, err
	}

	// compare sets with set templates to get their sync status
	statusMap, err := st.isSyncRequired(kit, option.BizID, option.SetTemplateID, setIDs, setMap, false,
		attrIdPropertyIdMap, setTemplateAttrValueMap)
	if err != nil {
		blog.Errorf("check if set need sync failed, err: %v, set ids: %+v, rid: %s", err, setIDs, kit.Rid)
		return &metadata.ListAPITaskSyncStatusResult{}, err
	}

	reformatStatuses, err := st.rearrangeSetTempSyncStatus(kit, option, taskStatusRes, statusMap)
	if err != nil {
		return &metadata.ListAPITaskSyncStatusResult{}, err
	}

	return &metadata.ListAPITaskSyncStatusResult{Count: int64(len(sets)), Info: reformatStatuses}, nil
}

// rearrangeSetTempSyncStatus set status by actual status and do another round of filter by status
func (st *setTemplate) rearrangeSetTempSyncStatus(kit *rest.Kit, option *metadata.ListSetTemplateSyncStatusOption,
	taskStatuses []metadata.APITaskSyncStatus, statusMap map[int64]bool) ([]metadata.APITaskSyncStatus,
	errors.CCErrorCoder) {

	statusFilterMap := make(map[metadata.APITaskStatus]struct{})
	if len(option.Status) > 0 {
		for _, status := range option.Status {
			statusFilterMap[status] = struct{}{}
		}
	}

	statuses := make([]metadata.APITaskSyncStatus, 0)
	statusExistsMap := make(map[int64]struct{})
	for _, status := range taskStatuses {
		statusExistsMap[status.InstID] = struct{}{}
		// if current status and api task status does not match, use current status
		if statusMap[status.InstID] && status.Status.IsSuccessful() {
			status.Status = metadata.APITAskStatusNeedSync
		} else if !statusMap[status.InstID] && !status.Status.IsSuccessful() {
			status.Status = metadata.APITaskStatusSuccess
		}

		// only returns the statuses that matches the status filter after comparing with current status
		if len(option.Status) > 0 {
			if _, exists := statusFilterMap[status.Status]; !exists {
				continue
			}
		}
		statuses = append(statuses, status)
	}

	// compensate for the sets that hasn't been synced before, or its latest sync task is already outdated
	compensateSetIDs := make([]int64, 0)
	for _, setID := range option.SetIDs {
		if _, exists := statusExistsMap[setID]; !exists {
			compensateSetIDs = append(compensateSetIDs, setID)
		}
	}

	setOpt := &metadata.QueryCondition{
		Fields: []string{common.BKSetIDField, common.CreatorField, common.CreateTimeField,
			common.LastTimeField},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		Condition:      mapstr.MapStr{common.BKSetIDField: mapstr.MapStr{common.BKDBIN: compensateSetIDs}},
		DisableCounter: true,
	}
	setRes := new(metadata.ResponseSetInstance)
	if err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDSet, setOpt, &setRes); err != nil {
		blog.Errorf("get sets failed, err: %v, opt: %#v, rid: %s", err, setOpt, kit.Rid)
		return nil, err
	}
	if err := setRes.CCError(); err != nil {
		blog.Errorf("get sets failed, err: %v, opt: %#v, rid: %s", err, setOpt, kit.Rid)
		return nil, err
	}

	for _, set := range setRes.Data.Info {
		status := metadata.APITaskSyncStatus{
			InstID:     set.SetID,
			Creator:    set.Creator,
			CreateTime: set.CreateTime.Time,
			LastTime:   set.LastTime.Time,
		}

		if statusMap[status.InstID] {
			status.Status = metadata.APITAskStatusNeedSync
		} else {
			status.Status = metadata.APITaskStatusSuccess
		}

		// only returns the statuses that matches the status filter after comparing with current status
		if len(option.Status) > 0 {
			if _, exists := statusFilterMap[status.Status]; !exists {
				continue
			}
		}
		statuses = append(statuses, status)
	}

	return statuses, nil
}

// ListSetTemplateSyncHistory list set template sync history
func (st *setTemplate) ListSetTemplateSyncHistory(kit *rest.Kit, option *metadata.ListSetTemplateSyncStatusOption) (
	*metadata.ListAPITaskSyncStatusResult, errors.CCErrorCoder) {

	setCond, err := option.ToSetCond(kit.CCError)
	if err != nil {
		blog.Errorf("parse set condition failed, err: %v, cond: %#v, rid: %s", err, option, kit.Rid)
		return nil, err
	}

	distinctOpt := &metadata.DistinctFieldOption{
		TableName: common.BKTableNameBaseSet,
		Field:     common.BKSetIDField,
		Filter:    setCond,
	}

	rawSetIDs, err := st.client.CoreService().Common().GetDistinctField(kit.Ctx, kit.Header, distinctOpt)
	if err != nil {
		blog.Errorf("get biz ids failed, err: %v, opt: %#v, rid: %s", err, distinctOpt, kit.Rid)
		return nil, err
	}

	if len(rawSetIDs) == 0 {
		return &metadata.ListAPITaskSyncStatusResult{Count: 0, Info: make([]metadata.APITaskSyncStatus, 0)}, nil
	}

	setIDs, ccErr := util.SliceInterfaceToInt64(rawSetIDs)
	if ccErr != nil {
		blog.Errorf("parse set ids to int failed, err: %v, raw ids: %+v, rid: %s", err, rawSetIDs, kit.Rid)
		return nil, err
	}

	option.SetIDs = setIDs
	statusCond, err := option.ToStatusCond(kit.CCError)
	if err != nil {
		blog.Errorf("parse status condition failed, err: %v, cond: %#v, rid: %s", err, option, kit.Rid)
		return nil, err
	}

	taskStatusRes, err := st.client.TaskServer().Task().ListSyncStatusHistory(kit.Ctx, kit.Header, statusCond)
	if err != nil {
		blog.Errorf("list sync status history failed, option: %#v, err: %v, rid: %s", statusCond, err, kit.Rid)
		return nil, err
	}

	return taskStatusRes, nil
}
