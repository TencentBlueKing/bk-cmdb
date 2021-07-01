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

func extractSetTemplateVersionFromTaskData(detail *metadata.APITaskDetail) (int64, error) {
	// TODO: better to implement with JSONPath
	if detail == nil {
		return 0, fmt.Errorf("detail field empty")
	}
	if len(detail.Detail) == 0 {
		return 0, fmt.Errorf("detail field empty")
	}
	detailData, ok := detail.Detail[0].Data.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("detail[0].data field")
	}
	version, ok := detailData["set_template_version"]
	if !ok {
		return 0, fmt.Errorf("detail[0].data.set_template_version field doesn't exist")
	}
	versionInt, err := util.GetInt64ByInterface(version)
	if err != nil {
		return 0, fmt.Errorf("parse set_template_version field failed, err: %+v", err)
	}
	return versionInt, nil
}

func (st *setTemplate) UpdateSetSyncStatus(kit *rest.Kit, setID int64) (metadata.SetTemplateSyncStatus, errors.CCErrorCoder) {
	setSyncStatus := metadata.SetTemplateSyncStatus{}
	set, err := st.GetOneSet(kit, setID)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, GetOneSet failed, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
		return setSyncStatus, err
	}
	if set.SetTemplateID == common.SetTemplateIDNotSet {
		blog.V(3).Infof("UpdateSetSyncStatus success, set not bound with template, setID: %d, rid: %s", setID, kit.Rid)
		return setSyncStatus, nil
	}
	option := metadata.DiffSetTplWithInstOption{
		SetIDs: []int64{set.SetID},
	}
	diff, err := st.DiffSetTplWithInst(kit, set.BizID, set.SetTemplateID, option)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, DiffSetTplWithInst failed, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
		return setSyncStatus, err
	}
	if len(diff) == 0 {
		blog.Errorf("UpdateSetSyncStatus failed, DiffSetTplWithInst result empty, setID: %d, rid: %s", setID, kit.Rid)
		return setSyncStatus, kit.CCError.CCError(common.CCErrCommInternalServerError)
	}
	setDiff := diff[0]

	detail, err := st.GetLatestSyncTaskDetail(kit, setID)
	if err != nil {
		return setSyncStatus, err
	}
	var syncStatus metadata.SyncStatus
	if detail == nil {
		if setDiff.NeedSync {
			syncStatus = metadata.SyncStatusWaiting
		} else {
			syncStatus = metadata.SyncStatusFinished
		}
	} else if !detail.Status.IsFinished() {
		syncStatus = metadata.SyncStatusSyncing
	} else if detail.Status.IsSuccessful() {
		if setDiff.NeedSync {
			syncStatus = metadata.SyncStatusWaiting
		} else {
			syncStatus = metadata.SyncStatusFinished
		}
	} else if detail.Status.IsFailure() {
		syncStatus = metadata.SyncStatusFailure
	} else {
		blog.ErrorJSON("unexpected task status: %s, rid: %s", detail, kit.Rid)
		return setSyncStatus, kit.CCError.CCError(common.CCErrCommInternalServerError)
	}

	// update sync status
	setSyncStatus = metadata.SetTemplateSyncStatus{
		SetID:           set.SetID,
		Name:            set.SetName,
		BizID:           set.BizID,
		SetTemplateID:   set.SetTemplateID,
		SupplierAccount: set.SupplierAccount,
		Status:          syncStatus,
	}
	setTemplateVersion := int64(0)
	if detail == nil {
		// no sync task has been run, just use
		setSyncStatus.Creator = set.Creator
		setSyncStatus.CreateTime = set.CreateTime
		setSyncStatus.LastTime = set.LastTime
		setSyncStatus.TaskID = ""
		setSyncStatus.Creator = kit.User
	} else {
		version, err := extractSetTemplateVersionFromTaskData(detail)
		if err != nil && blog.V(5) {
			blog.InfoJSON("extractSetTemplateVersionFromTaskData failed, detail: %s, err: %s", detail, err.Error())
		}
		setTemplateVersion = version
		setSyncStatus.Creator = detail.User
		setSyncStatus.CreateTime = metadata.Time{Time: detail.CreateTime}
		setSyncStatus.LastTime = metadata.Time{Time: detail.LastTime}
		setSyncStatus.TaskID = detail.TaskID
	}
	if setSyncStatus.Status == metadata.SyncStatusWaiting {
		setSyncStatus.TaskID = ""
	}

	if setTemplateVersion != 0 {
		if ccErr := st.UpdateSetVersion(kit, setID, setTemplateVersion); ccErr != nil {
			blog.Errorf("UpdateSetSyncStatus failed, UpdateSetVersion failed, setID: %d, setTemplateVersion: %d, err: %s, rid: %s", setID, setTemplateVersion, ccErr.Error(), kit.Rid)
			return setSyncStatus, ccErr
		}
	}

	err = st.client.CoreService().SetTemplate().UpdateSetTemplateSyncStatus(kit.Ctx, kit.Header, setID, setSyncStatus)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, UpdateSetTemplateSyncStatus failed, setID: %d, err: %s, rid: %s", setID, err.Error(), kit.Rid)
		return setSyncStatus, err
	}

	return setSyncStatus, nil
}

func (st *setTemplate) UpdateSetVersion(kit *rest.Kit, setID, setTemplateVersion int64) errors.CCErrorCoder {
	updateSetOption := &metadata.UpdateOption{
		Data: map[string]interface{}{
			"set_template_version": setTemplateVersion,
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

func (st *setTemplate) GetLatestSyncTaskDetail(kit *rest.Kit, setID int64) (*metadata.APITaskDetail, errors.CCErrorCoder) {
	setRelatedTaskFilter := map[string]interface{}{
		// "detail.data.set.bk_set_id": setID,
		"flag": metadata.GetSetTemplateSyncIndex(setID),
	}
	listTaskOption := metadata.ListAPITaskRequest{
		Condition: mapstr.MapStr(setRelatedTaskFilter),
		Page: metadata.BasePage{
			Sort:  "-create_time",
			Limit: 1,
		},
	}

	listResult, err := st.client.TaskServer().Task().ListTask(kit.Ctx, kit.Header, common.SyncSetTaskName, &listTaskOption)
	if err != nil {
		blog.ErrorJSON("list set sync tasks failed, option: %s, err: %v, rid: %s", listTaskOption, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTaskListTaskFail)
	}
	if listResult == nil || len(listResult.Data.Info) == 0 {
		blog.InfoJSON("list set sync tasks result empty, option: %s, result: %s, rid: %s", listTaskOption, listTaskOption, kit.Rid)
		return nil, nil
	}
	taskDetail := &listResult.Data.Info[0]
	clearSetSyncTaskDetail(taskDetail)
	return taskDetail, nil
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
func (st *setTemplate) TriggerCheckSetTemplateSyncingStatus(kit *rest.Kit, bizID, setTemplateID, setID int64) errors.CCErrorCoder {
	setTempLock := lock.NewLocker(redis.Client())
	key := lock.GetLockKey(lock.CheckSetTemplateSyncFormat, setID)
	locked, err := setTempLock.Lock(key, time.Minute)
	if err != nil {
		blog.Errorf("get sync set template  lock error. set template id: %d, setID: %d, err: %s, rid: %s", setTemplateID, setID, err.Error(), kit.Rid)
		return kit.CCError.CCError(common.CCErrCommRedisOPErr)
	}
	if locked {
		defer setTempLock.Unlock()
		_, err := st.UpdateSetSyncStatus(kit, setID)
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
				st.TriggerCheckSetTemplateSyncingStatus(kit.NewKit(), info.BizID, info.SetTemplateID, info.SetID)
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
