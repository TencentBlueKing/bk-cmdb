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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

func (st *setTemplate) GetOneSet(params types.ContextParams, setID int64) (metadata.SetInst, errors.CCErrorCoder) {
	set := metadata.SetInst{}

	filter := map[string]interface{}{
		common.BKSetIDField: setID,
	}
	qc := &metadata.QueryCondition{
		Condition: filter,
	}
	instResult, err := st.client.CoreService().Instance().ReadInstance(params.Context, params.Header, common.BKInnerObjIDSet, qc)
	if err != nil {
		blog.ErrorJSON("GetOneSet failed, db select failed, filter: %s, err: %s, rid: %s", filter, err.Error(), params.ReqID)
		return set, params.Err.CCError(common.CCErrCommDBSelectFailed)
	}
	if instResult.Result == false || instResult.Code != 0 {
		blog.ErrorJSON("GetOneSet failed, read instance failed, filter: %s, instResult: %s, rid: %s", filter, instResult, params.ReqID)
		return set, errors.NewCCError(instResult.Code, instResult.ErrMsg)
	}
	if len(instResult.Data.Info) == 0 {
		blog.ErrorJSON("GetOneSet failed, not found, filter: %s, instResult: %s, rid: %s", filter, instResult, params.ReqID)
		return set, params.Err.CCError(common.CCErrCommNotFound)
	}
	if len(instResult.Data.Info) > 1 {
		blog.ErrorJSON("GetOneSet failed, got multiple, filter: %s, instResult: %s, rid: %s", filter, instResult, params.ReqID)
		return set, params.Err.CCError(common.CCErrCommGetMultipleObject)
	}
	if err := mapstruct.Decode2Struct(instResult.Data.Info[0], &set); err != nil {
		blog.ErrorJSON("GetOneSet failed, unmarshal set failed, instResult: %s, err: %s, rid: %s", instResult, err.Error(), params.ReqID)
		return set, params.Err.CCError(common.CCErrCommJSONUnmarshalFailed)
	}
	return set, nil
}

func (st *setTemplate) UpdateSetSyncStatus(params types.ContextParams, setID int64) (metadata.SetTemplateSyncStatus, errors.CCErrorCoder) {
	setSyncStatus := metadata.SetTemplateSyncStatus{}
	set, err := st.GetOneSet(params, setID)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, GetOneSet failed, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return setSyncStatus, err
	}
	if set.SetTemplateID == common.SetTemplateIDNotSet {
		blog.V(3).Infof("UpdateSetSyncStatus success, set not bound with template, setID: %d, rid: %s", setID, params.ReqID)
		return setSyncStatus, nil
	}
	option := metadata.DiffSetTplWithInstOption{
		SetIDs: []int64{set.SetID},
	}
	diff, err := st.DiffSetTplWithInst(params.Context, params.Header, set.BizID, set.SetTemplateID, option)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, DiffSetTplWithInst failed, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return setSyncStatus, err
	}
	if len(diff) == 0 {
		blog.Errorf("UpdateSetSyncStatus failed, DiffSetTplWithInst result empty, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return setSyncStatus, params.Err.CCError(common.CCErrCommInternalServerError)
	}
	setDiff := diff[0]

	detail, err := st.GetLatestSyncTaskDetail(params, setID)
	if err != nil {
		return setSyncStatus, err
	}
	syncStatus := metadata.SyncStatusWaiting
	if detail == nil {
		if setDiff.NeedSync == true {
			syncStatus = metadata.SyncStatusWaiting
		} else {
			syncStatus = metadata.SyncStatusFinished
		}
	} else if detail.Status.IsFinished() == false {
		syncStatus = metadata.SyncStatusSyncing
	} else if detail.Status.IsSuccessful() == true {
		if setDiff.NeedSync == true {
			syncStatus = metadata.SyncStatusWaiting
		} else {
			syncStatus = metadata.SyncStatusFinished
		}
	} else if detail.Status.IsFailure() == true {
		syncStatus = metadata.SyncStatusFailure
	} else {
		blog.ErrorJSON("unexpected task status: %s, rid: %s", detail, params.ReqID)
		return setSyncStatus, params.Err.CCError(common.CCErrCommInternalServerError)
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
	if detail == nil {
		// no sync task has been run, just use
		setSyncStatus.Creator = set.Creator
		setSyncStatus.CreateTime = set.CreateTime
		setSyncStatus.LastTime = set.LastTime
		setSyncStatus.TaskID = ""
		setSyncStatus.Creator = params.User
	} else {
		setSyncStatus.Creator = detail.User
		setSyncStatus.CreateTime = metadata.Time{Time: detail.CreateTime}
		setSyncStatus.LastTime = metadata.Time{Time: detail.LastTime}
		setSyncStatus.TaskID = detail.TaskID
	}
	if setSyncStatus.Status == metadata.SyncStatusWaiting {
		setSyncStatus.TaskID = ""
	}
	err = st.client.CoreService().SetTemplate().UpdateSetTemplateSyncStatus(params.Context, params.Header, setID, setSyncStatus)
	if err != nil {
		blog.Errorf("UpdateSetSyncStatus failed, UpdateSetTemplateSyncStatus failed, setID: %d, err: %s, rid: %s", setID, err.Error(), params.ReqID)
		return setSyncStatus, err
	}

	return setSyncStatus, nil
}

func (st *setTemplate) GetLatestSyncTaskDetail(params types.ContextParams, setID int64) (*metadata.APITaskDetail, errors.CCErrorCoder) {
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

	listResult, err := st.client.TaskServer().Task().ListTask(params.Context, params.Header, common.SyncSetTaskName, &listTaskOption)
	if err != nil {
		blog.ErrorJSON("list set sync tasks failed, option: %s, rid: %s", listTaskOption, params.ReqID)
		return nil, params.Err.CCError(common.CCErrTaskListTaskFail)
	}
	if listResult == nil || len(listResult.Data.Info) == 0 {
		blog.InfoJSON("list set sync tasks result empty, option: %s, result: %s, rid: %s", listTaskOption, listTaskOption, params.ReqID)
		return nil, nil
	}
	taskDetail := &listResult.Data.Info[0]
	clearSetSyncTaskDetail(taskDetail)
	return taskDetail, nil
}

func clearSetSyncTaskDetail(detail *metadata.APITaskDetail) {
	detail.Header = http.Header{}
	for taskIdx := range detail.Detail {
		subTaskDetail, ok := detail.Detail[taskIdx].Data.(map[string]interface{})
		if ok == false {
			blog.Warnf("clearSetSyncTaskDetail expect map[string]interface{}, got unexpected type, data: %+v", detail.Detail[taskIdx].Data)
			detail.Detail[taskIdx].Data = nil
		}
		delete(subTaskDetail, "header")
	}
}
