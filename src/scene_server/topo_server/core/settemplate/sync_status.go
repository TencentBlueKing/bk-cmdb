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
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/lock"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/redis"
)

func (st *setTemplate) GetOneSet(kit *rest.Kit, setID int64) (metadata.SetInst, errors.CCErrorCoder) {
	set := metadata.SetInst{}

	filter := map[string]interface{}{
		common.BKSetIDField: setID,
	}
	qc := &metadata.QueryCondition{
		Condition: filter,
	}
	instResult, err := st.client.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, qc)
	if err != nil {
		blog.ErrorJSON("GetOneSet failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return set, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := instResult.CCError(); ccErr != nil {
		blog.ErrorJSON("GetOneSet failed, read instance failed, filter: %s, instResult: %s, rid: %s", filter, instResult, kit.Rid)
		return set, ccErr
	}
	if len(instResult.Data.Info) == 0 {
		blog.ErrorJSON("GetOneSet failed, not found, filter: %s, instResult: %s, rid: %s", filter, instResult, kit.Rid)
		return set, kit.CCError.CCError(common.CCErrCommNotFound)
	}
	if len(instResult.Data.Info) > 1 {
		blog.ErrorJSON("GetOneSet failed, got multiple, filter: %s, instResult: %s, rid: %s", filter, instResult, kit.Rid)
		return set, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	if err := mapstruct.Decode2StructWithHook(instResult.Data.Info[0], &set); err != nil {
		blog.ErrorJSON("GetOneSet failed, unmarshal set failed, instResult: %s, err: %s, rid: %s", instResult, err.Error(), kit.Rid)
		return set, kit.CCError.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	return set, nil
}

func (st *setTemplate) GetSets(kit *rest.Kit, setTemplateID int64, setIDs []int64) ([]metadata.SetInst, errors.CCErrorCoder) {
	filter := &metadata.QueryCondition{}
	filter.Condition = mapstr.MapStr{
		common.BKSetIDField:         map[string]interface{}{common.BKDBIN: setIDs},
		common.BKSetTemplateIDField: setTemplateID,
	}
	instResult := metadata.ResponseSetInstance{}
	err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header, common.BKInnerObjIDSet, filter, &instResult)
	if err != nil {
		blog.ErrorJSON("GetSets failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if ccErr := instResult.CCError(); ccErr != nil {
		blog.ErrorJSON("GetSets failed, read instance failed, filter: %s, instResult: %s, rid: %s", filter, instResult, kit.Rid)
		return nil, ccErr
	}

	if len(instResult.Data.Info) == 0 {
		blog.ErrorJSON("GetSets failed, not found, filter: %s, instResult: %s, rid: %s", filter, instResult, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if instResult.Data.Count != len(setIDs) {
		blog.Errorf("GetSets failed, some setID invalid, input IDs: %+v, valid ,IDs: %+v, rid: %s",
			setIDs, instResult.Data.Info, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "bk_set_ids")
	}

	return instResult.Data.Info, nil
}

func getSetIDFromTaskDetail(kit *rest.Kit, detail metadata.APITaskDetail) (int64, error) {
	setID, err := strconv.ParseInt(detail.Flag[len("set_template_sync:"):], 10, 64)
	if err != nil {
		blog.Errorf("getSetIDFromTaskDetail failed, err: %+v, rid: %s", err, kit.Rid)
		return 0, err
	}
	return setID, nil
}

func (st *setTemplate) isSyncRequired(kit *rest.Kit, bizID int64, setTemplateID int64, setIDs []int64) (map[int64]bool, errors.CCErrorCoder) {
	serviceTemplates, err := st.client.CoreService().SetTemplate().ListSetTplRelatedSvcTpl(kit.Ctx, kit.Header, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("DiffSetTemplateWithInstances failed, ListSetTplRelatedSvcTpl failed, bizID: %d, "+
			"setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), kit.Rid)
		return nil, err
	}

	serviceTemplateCnt := int64(len(serviceTemplates))
	serviceTemplateMap := make(map[int64]metadata.ServiceTemplate, serviceTemplateCnt)
	for _, serviceTemplate := range serviceTemplates {
		serviceTemplateMap[serviceTemplate.ID] = serviceTemplate
	}

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
		blog.ErrorJSON("DiffSetTemplateWithInstances failed, list modules failed, bizID: %s, setTemplateID: %s,"+
			" setIDs: %s, err: %s, rid: %s", bizID, setTemplateID, setIDs, err, kit.Rid)
		return nil, err
	}
	if err := modulesInstResult.CCError(); err != nil {
		blog.ErrorJSON("DiffSetTemplateWithInstances failed, list module http reply failed, bizID: %s, "+
			"setTemplateID: %s, setIDs: %s, filter: %s, reply: %s, rid: %s", bizID,
			setTemplateID, setIDs, moduleFilter, modulesInstResult, kit.Rid)
		return nil, err
	}

	setModulesCnt := int64(modulesInstResult.Data.Count)
	setModules := make(map[int64][]metadata.ModuleInst, setModulesCnt)
	for _, module := range modulesInstResult.Data.Info {
		if _, exist := setModules[module.ParentID]; !exist {
			setModules[module.ParentID] = make([]metadata.ModuleInst, 0)
		}
		setModules[module.ParentID] = append(setModules[module.ParentID], module)
	}

	checkResult := make(map[int64]bool, len(setIDs))
	for idx, module := range setModules {
		NeedSync := DiffModuleServiceTpl(serviceTemplateCnt, serviceTemplateMap, setModulesCnt, module)
		checkResult[idx] = NeedSync
	}

	return checkResult, nil
}

func DiffModuleServiceTpl(serviceTplCnt int64, serviceTemplates map[int64]metadata.ServiceTemplate, moduleCnt int64,
	modules []metadata.ModuleInst) bool {

	if serviceTplCnt != moduleCnt {
		return true
	}

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

func (st *setTemplate) UpdateSetSyncStatus(kit *rest.Kit, setTemplateID int64, setID []int64) ([]metadata.SetTemplateSyncStatus, errors.CCErrorCoder) {
	var setSyncStatus []metadata.SetTemplateSyncStatus

	if setTemplateID == common.SetTemplateIDNotSet {
		blog.V(3).Infof("UpdateSetSyncStatus success, set not bound with template, setID: %d, rid: %s", setID, kit.Rid)
		return setSyncStatus, nil
	}

	sets, err := st.GetSets(kit, setTemplateID, setID)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, GetSets failed, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
		return setSyncStatus, err
	}

	bizID := sets[0].BizID
	checkNeedSyncResult, err := st.isSyncRequired(kit, bizID, setTemplateID, setID)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, check sync required failed, templateID: %d, setID: %d, err: %s, rid: %s",
			setTemplateID, setID, err.Error(), kit.Rid)
		return setSyncStatus, err
	}

	if len(checkNeedSyncResult) == 0 {
		blog.Errorf("UpdateSetSyncStatus failed, checkNeedSyncResult empty, templateID: %d, setID: %d, rid: %s",
			setTemplateID, setID, kit.Rid)
		return setSyncStatus, kit.CCError.CCError(common.CCErrCommInternalServerError)
	}

	details, err := st.GetLatestSyncTaskDetail(kit, setID)
	if err != nil {
		return setSyncStatus, err
	}
	for _, set := range sets {
		// update sync status
		syncStatus := metadata.SetTemplateSyncStatus{
			SetID:           set.SetID,
			Name:            set.SetName,
			BizID:           set.BizID,
			SetTemplateID:   set.SetTemplateID,
			SupplierAccount: set.SupplierAccount,
			Creator:         kit.User,
			CreateTime:      set.CreateTime,
			LastTime:        set.LastTime,
			TaskID:          "",
			Status:          metadata.SyncStatusFinished,
		}

		if checkNeedSyncResult[set.SetID] {
			syncStatus.Status = metadata.SyncStatusWaiting
		}

		if _, ok := details[set.SetID]; !ok {
			setSyncStatus = append(setSyncStatus, syncStatus)
			continue
		}

		syncStatus.Creator = details[set.SetID].User
		syncStatus.CreateTime = metadata.Time{Time: details[set.SetID].CreateTime}
		syncStatus.LastTime = metadata.Time{Time: details[set.SetID].LastTime}
		syncStatus.TaskID = details[set.SetID].TaskID

		if !details[set.SetID].Status.IsFinished() {
			syncStatus.Status = metadata.SyncStatusSyncing
		}

		if details[set.SetID].Status.IsFailure() {
			syncStatus.Status = metadata.SyncStatusFailure
		}

		setSyncStatus = append(setSyncStatus, syncStatus)
	}

	err = st.client.CoreService().SetTemplate().UpdateManySetTemplateSyncStatus(kit.Ctx, kit.Header, setSyncStatus)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, UpdateSetTemplateSyncStatus failed, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
		return setSyncStatus, err
	}

	return setSyncStatus, nil
}

func (st *setTemplate) UpdateSetVersion(kit *rest.Kit, setID, setTemplateVersion int64) errors.CCErrorCoder {
	updateSetOption := &metadata.UpdateOption{
		Data: map[string]interface{}{
			common.BKSetTemplateVersionField: setTemplateVersion,
		},
		Condition: map[string]interface{}{
			common.BKSetIDField: setID,
		},
	}
	updateSetResult, err := st.client.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, common.BKInnerObjIDSet, updateSetOption)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, UpdateInstance of set failed, option: %+v, err: %s, rid: %s", updateSetOption, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if ccErr := updateSetResult.CCError(); ccErr != nil {
		blog.Errorf("UpdateSetSyncStatus failed, UpdateInstance failed, option: %+v, result: %+v, rid: %s", updateSetOption, updateSetResult, kit.Rid)
		return ccErr
	}
	return nil
}

func (st *setTemplate) GetLatestSyncTaskDetail(kit *rest.Kit, setID []int64) (map[int64]*metadata.APITaskDetail, errors.CCErrorCoder) {
	var setIndex []string
	for _, item := range setID {
		setIndex = append(setIndex, metadata.GetSetTemplateSyncIndex(item))
	}
	setRelatedTaskFilter := map[string]interface{}{
		// "detail.data.set.bk_set_id": setID,
		"flag": map[string]interface{}{common.BKDBIN: setIndex},
	}
	listTaskOption := metadata.ListAPITaskLatestRequest{
		Condition: mapstr.MapStr(setRelatedTaskFilter),
	}

	listResult, err := st.client.TaskServer().Task().ListLatestTask(kit.Ctx, kit.Header, common.SyncSetTaskName, &listTaskOption)
	if err != nil {
		blog.ErrorJSON("list set sync tasks failed, option: %s, err: %v, rid: %s", listTaskOption, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTaskListTaskFail)
	}
	if listResult == nil || len(listResult.Data) == 0 {
		blog.InfoJSON("list set sync tasks result empty, option: %s, result: %s, rid: %s", listTaskOption, listTaskOption, kit.Rid)
		return nil, nil
	}

	latestTaskResult := make(map[int64]*metadata.APITaskDetail)
	for _, APITask := range listResult.Data {
		clearSetSyncTaskDetail(&APITask)
		setID, err := getSetIDFromTaskDetail(kit, APITask)
		if err != nil {
			blog.Errorf("get setID from task failed, err: %+v, rid: %s", err, kit.Rid)
		}
		latestTaskResult[setID] = &APITask
	}
	return latestTaskResult, nil
}

func clearSetSyncTaskDetail(detail *metadata.APITaskDetail) {
	detail.Header = util.BuildHeader(common.CCSystemOperatorUserName, common.BKDefaultOwnerID)
	for taskIdx := range detail.Detail {
		subTaskDetail, ok := detail.Detail[taskIdx].Data.(map[string]interface{})
		if !ok {
			blog.Warnf("clearSetSyncTaskDetail expect map[string]interface{}, got unexpected type, data: %+v", detail.Detail[taskIdx].Data)
			detail.Detail[taskIdx].Data = nil
		}
		delete(subTaskDetail, "header")
	}
}

// TriggerCheckSetTemplateSyncingStatus  触发对正在同步中任务的状态改变处理
func (st *setTemplate) TriggerCheckSetTemplateSyncingStatus(kit *rest.Kit, bizID, setTemplateID int64, setID []int64) errors.CCErrorCoder {
	setTempLock := lock.NewLocker(redis.Client())
	key := lock.GetLockKey(lock.CheckSetTemplateSyncFormat, setID)
	locked, err := setTempLock.Lock(key, time.Minute)
	if err != nil {
		blog.Errorf("get sync set template  lock error. set template id: %d, setID: %d, err: %s, rid: %s", setTemplateID, setID, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommRedisOPErr)
	}
	if locked {
		defer setTempLock.Unlock()
		_, err := st.UpdateSetSyncStatus(kit, setTemplateID, setID)
		if err != nil {
			return err
		}

	} else {
		blog.Warnf("skip task, reason not get lock . set template id: %d, setID: %d, rid: %s", setTemplateID, setID, kit.Rid)
	}
	return nil
}

func (st *setTemplate) ListSetTemplateSyncStatus(kit *rest.Kit, bizID int64,
	option metadata.ListSetTemplateSyncStatusOption) (metadata.MultipleSetTemplateSyncStatus, errors.CCErrorCoder) {

	responseInfo := metadata.MultipleSetTemplateSyncStatus{}

	filterTemp := &metadata.QueryCondition{
		Page:      option.Page,
		Condition: mapstr.MapStr{common.BKSetTemplateIDField: option.SetTemplateID, common.BKAppIDField: bizID},
	}
	if len(option.SetIDs) != 0 {
		filterTemp.Condition[common.BKSetIDField] = mapstr.MapStr{common.BKDBIN: option.SetIDs}
	}
	filterTemp.Fields = []string{common.BKSetIDField, common.BKSetNameField, common.BkSupplierAccount}

	var setInfoResp metadata.ResponseSetInstance
	err := st.client.CoreService().Instance().ReadInstanceStruct(kit.Ctx, kit.Header,
		common.BKInnerObjIDSet, filterTemp, &setInfoResp)
	if err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, core service find set template failed,"+
			" option: %s, err: %s, rid: %s", filterTemp, err, kit.Rid)
		return metadata.MultipleSetTemplateSyncStatus{}, err
	}
	if err := setInfoResp.CCError(); err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, core service find set template http reply failed,"+
			" option: %s, err: %s, rid: %s", filterTemp, err, kit.Rid)
		return metadata.MultipleSetTemplateSyncStatus{}, err
	}

	responseInfo.Count = int64(setInfoResp.Data.Count)
	setIDs := make([]int64, len(setInfoResp.Data.Info))
	responseInfo.Info = make([]metadata.SetTemplateSyncStatus, len(setInfoResp.Data.Info))

	for idx, setInfo := range setInfoResp.Data.Info {
		setIDs[idx] = setInfo.SetID
		// setInfoResp 只返回了部分字段，新加字段注意修改
		responseInfo.Info[idx] = metadata.SetTemplateSyncStatus{
			SetID:           setInfo.SetID,
			Name:            setInfo.SetName,
			BizID:           bizID,
			SetTemplateID:   option.SetTemplateID,
			SupplierAccount: setInfo.SupplierAccount,
		}
	}

	// 使用存在模块
	option.SetIDs = setIDs
	result, err := st.client.CoreService().SetTemplate().ListSetTemplateSyncStatus(kit.Ctx, kit.Header, bizID, option)
	if err != nil {
		blog.ErrorJSON("ListSetTemplateSyncStatus failed, core service search failed, option: %s, err: %s, rid: %s",
			option, err.Error(), kit.Rid)
		return metadata.MultipleSetTemplateSyncStatus{}, err
	}

	setTempSyncMap := make(map[int64]metadata.SetTemplateSyncStatus, len(result.Info))
	// 处理当前需要同步任务的状态
	for _, info := range result.Info {
		setTempSyncMap[info.SetID] = info
		if !info.Status.IsFinished() {
			go func(info metadata.SetTemplateSyncStatus) {
				setID := []int64{info.SetID}
				st.TriggerCheckSetTemplateSyncingStatus(kit.NewKit(), info.BizID, info.SetTemplateID, setID)
			}(info)
		}

	}
	// 如果在同步表中有数据，使用同步表中的数据
	for idx, row := range responseInfo.Info {
		if newRow, ok := setTempSyncMap[row.SetID]; ok {
			responseInfo.Info[idx] = newRow
		}
	}

	return responseInfo, nil

}
