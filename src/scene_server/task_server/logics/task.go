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
	"net/http"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"

	"github.com/rs/xid"
)

// Create add task
func (lgc *Logics) Create(kit *rest.Kit, input *metadata.CreateTaskRequest) (metadata.APITaskDetail, error) {
	dbTask := metadata.APITaskDetail{}
	input.TaskType = strings.TrimSpace(input.TaskType)

	if input.TaskType == "" {
		return dbTask, kit.CCError.Errorf(common.CCErrCommParamsNeedString, "name")
	}

	if len(input.Data) == 0 {
		return dbTask, kit.CCError.Errorf(common.CCErrCommParamsNeedString, "data")
	}

	// check if there is another unfinished task already created, forbidden create duplicate tasks
	// TODO: replace with index when partial filter supports $in operator
	duplicateCond := mapstr.MapStr{
		common.BKTaskTypeField: input.TaskType,
		common.BKInstIDField:   input.InstID,
		common.BKStatusField: map[string]interface{}{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute,
				metadata.APITaskStatusExecute},
		},
	}
	cnt, err := lgc.db.Table(common.BKTableNameAPITask).Find(duplicateCond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("get duplicate tasks failed, err: %v, cond: %#v, rid: %s", err, duplicateCond, kit.Rid)
		return dbTask, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		return dbTask, kit.CCError.Errorf(common.CCErrTaskCreateConflict, input.InstID)
	}

	dbTask.TaskID = getStrTaskID("id")
	dbTask.User = kit.User
	dbTask.TaskType = input.TaskType
	dbTask.InstID = input.InstID
	dbTask.Header = GetDBHTTPHeader(kit.Header)
	dbTask.Status = metadata.APITaskStatusNew
	dbTask.CreateTime = time.Now()
	dbTask.LastTime = time.Now()
	for _, taskItem := range input.Data {
		dbTask.Detail = append(dbTask.Detail, metadata.APISubTaskDetail{
			SubTaskID: getStrTaskID("sid"),
			Data:      taskItem,
			Status:    metadata.APITaskStatusNew,
		})
	}
	err = lgc.db.Table(common.BKTableNameAPITask).Insert(kit.Ctx, dbTask)
	if err != nil {
		blog.Errorf("create task failed, data: %#v, err: %v, rid: %s", dbTask, err, kit.Rid)
		return dbTask, kit.CCError.Error(common.CCErrCommDBInsertFailed)
	}

	taskHistory := metadata.APITaskSyncStatus{
		TaskID:          dbTask.TaskID,
		TaskType:        input.TaskType,
		InstID:          input.InstID,
		Status:          metadata.APITaskStatusNew,
		Creator:         kit.User,
		CreateTime:      dbTask.CreateTime,
		LastTime:        dbTask.LastTime,
		SupplierAccount: kit.SupplierAccount,
	}

	if err := lgc.db.Table(common.BKTableNameAPITaskSyncHistory).Insert(kit.Ctx, taskHistory); err != nil {
		blog.Errorf("create task sync history failed, data: %#v, err: %v, rid: %s", taskHistory, err, kit.Rid)
		return dbTask, kit.CCError.Error(common.CCErrCommDBInsertFailed)
	}
	return dbTask, nil
}

// CreateBatch create task batch
func (lgc *Logics) CreateBatch(kit *rest.Kit, tasks []metadata.CreateTaskRequest) ([]metadata.APITaskDetail,
	error) {

	if len(tasks) == 0 {
		return make([]metadata.APITaskDetail, 0), nil
	}

	now := time.Now()
	dbTask := metadata.APITaskDetail{
		User:       kit.User,
		Header:     GetDBHTTPHeader(kit.Header),
		Status:     metadata.APITaskStatusNew,
		CreateTime: now,
		LastTime:   now,
	}

	taskHistory := metadata.APITaskSyncStatus{
		Status:          metadata.APITaskStatusNew,
		Creator:         kit.User,
		CreateTime:      now,
		LastTime:        now,
		SupplierAccount: kit.SupplierAccount,
	}

	dbTasks := make([]metadata.APITaskDetail, len(tasks))
	taskHistories := make([]metadata.APITaskSyncStatus, len(tasks))
	taskTypes := make([]string, len(tasks))
	instIDs := make([]int64, len(tasks))
	for index, task := range tasks {
		task.TaskType = strings.TrimSpace(task.TaskType)
		if task.TaskType == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.BKTaskTypeField)
		}

		if len(task.Data) == 0 {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "data")
		}

		taskTypes[index] = task.TaskType
		instIDs[index] = task.InstID

		dbTask.TaskID = getStrTaskID("id")
		dbTask.TaskType = task.TaskType
		dbTask.InstID = task.InstID
		dbTask.Detail = make([]metadata.APISubTaskDetail, 0)
		for _, taskItem := range task.Data {
			dbTask.Detail = append(dbTask.Detail, metadata.APISubTaskDetail{
				SubTaskID: getStrTaskID("sid"),
				Data:      taskItem,
				Status:    metadata.APITaskStatusNew,
			})
		}
		dbTasks[index] = dbTask

		taskHistory.TaskID = dbTask.TaskID
		taskHistory.TaskType = task.TaskType
		taskHistory.InstID = task.InstID
		taskHistories[index] = taskHistory
	}

	// check if there is another unfinished task already created, forbidden create duplicate tasks
	// TODO: replace with index when partial filter supports $in operator
	duplicateCond := mapstr.MapStr{
		common.BKTaskTypeField: map[string]interface{}{common.BKDBIN: taskTypes},
		common.BKInstIDField:   map[string]interface{}{common.BKDBIN: instIDs},
		common.BKStatusField: map[string]interface{}{
			common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute,
				metadata.APITaskStatusExecute},
		},
	}

	duplicateIDs, err := lgc.db.Table(common.BKTableNameAPITask).Distinct(kit.Ctx, common.BKInstIDField, duplicateCond)
	if err != nil {
		blog.Errorf("get duplicate tasks failed, err: %v, cond: %#v, rid: %s", err, duplicateCond, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommDBInsertFailed)
	}

	if len(duplicateIDs) > 0 {
		return nil, kit.CCError.Errorf(common.CCErrTaskCreateConflict, duplicateIDs)
	}

	if err = lgc.db.Table(common.BKTableNameAPITask).Insert(kit.Ctx, dbTasks); err != nil {
		blog.Errorf("create tasks(%#v) failed, err: %v, rid: %s", dbTasks, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommDBInsertFailed)
	}

	if err := lgc.db.Table(common.BKTableNameAPITaskSyncHistory).Insert(kit.Ctx, taskHistories); err != nil {
		blog.Errorf("create task sync history failed, data: %#v, err: %v, rid: %s", taskHistories, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommDBInsertFailed)
	}
	return dbTasks, nil
}

// List list task
func (lgc *Logics) List(kit *rest.Kit, name string, input *metadata.ListAPITaskRequest) ([]metadata.APITaskDetail,
	uint64, error) {

	if input == nil {
		input = new(metadata.ListAPITaskRequest)
	}
	if input.Condition == nil {
		input.Condition = mapstr.New()
	}
	input.Condition.Set("name", name)
	if input.Page.IsIllegal() {
		return nil, 0, kit.CCError.Errorf(common.CCErrCommPageLimitIsExceeded)
	}
	cnt, err := lgc.db.Table(common.BKTableNameAPITask).Find(input.Condition).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("list task failed, data: %#v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	rows := make([]metadata.APITaskDetail, 0)
	err = lgc.db.Table(common.BKTableNameAPITask).Find(input.Condition).Start(uint64(input.Page.Start)).
		Limit(uint64(input.Page.Limit)).Sort(input.Page.Sort).All(kit.Ctx, &rows)

	return rows, cnt, nil
}

// ListLatestTask list latest task
func (lgc *Logics) ListLatestTask(kit *rest.Kit, name string,
	input *metadata.ListAPITaskLatestRequest) ([]metadata.APITaskDetail, error) {
	/*
		aggregateCond parameter of aggregate to search the latest created task in input.Condition need.
		because multiple results of the same task may be at the front end of sorting by
		create_time field, use group to get the first result of each task
	*/
	aggregateCond := []map[string]interface{}{
		{common.BKDBSort: map[string]interface{}{common.CreateTimeField: -1}},
		{common.BKDBGroup: map[string]interface{}{
			"_id": "$bk_inst_id",
			"doc": map[string]interface{}{"$first": "$$ROOT"},
		}},
		{common.BKDBReplaceRoot: map[string]interface{}{"newRoot": "$doc"}},
	}

	if input == nil {
		input = &metadata.ListAPITaskLatestRequest{}
	}

	if input.Condition == nil {
		input.Condition = mapstr.New()
	}

	if len(name) != 0 {
		input.Condition.Set("name", name)
	}

	if len(input.Condition) != 0 {
		aggregateCond = append([]map[string]interface{}{{common.BKDBMatch: input.Condition}}, aggregateCond...)
	}

	if len(input.Fields) != 0 {
		cond := map[string]int64{}
		for _, field := range input.Fields {
			cond[field] = 1
		}
		aggregateCond = append(aggregateCond, map[string]interface{}{
			common.BKDBProject: cond,
		})
	}

	result := make([]metadata.APITaskDetail, 0)
	if err := lgc.db.Table(common.BKTableNameAPITask).AggregateAll(kit.Ctx, aggregateCond, &result); err != nil {
		blog.Errorf("list latest task failed, aggregateCond: %v, err: %v, rid: %v", aggregateCond, err, kit.Rid)
		return nil, err
	}

	return result, nil
}

// Detail  task detail
func (lgc *Logics) Detail(kit *rest.Kit, taskID string) (*metadata.APITaskDetail, error) {

	condition := mapstr.New()
	condition.Set("task_id", taskID)

	rows := make([]metadata.APITaskDetail, 0)
	err := lgc.db.Table(common.BKTableNameAPITask).Find(condition).All(kit.Ctx, &rows)
	if err != nil {
		blog.Errorf("get task detail failed, data: %#v, err: %v, rid: %s", condition, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	if len(rows) == 0 {
		return nil, nil
	}
	return &rows[0], nil
}

// DeleteTask delete task
func (lgc *Logics) DeleteTask(kit *rest.Kit, taskCond *metadata.DeleteOption) error {
	if len(taskCond.Condition) == 0 {
		blog.Errorf("task condition is empty, rid: %s", kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommInstDataNil, "task condition")
	}

	err := lgc.db.Table(common.BKTableNameAPITask).Delete(kit.Ctx, taskCond.Condition)
	if err != nil {
		blog.Errorf("delete task failed, err: %v, cond: %#v, rid: %s", err, taskCond, kit.Rid)
		return err
	}

	return nil
}

// ChangeStatusToSuccess task status change to success
func (lgc *Logics) ChangeStatusToSuccess(kit *rest.Kit, taskID, subTaskID string) error {

	if taskID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "task_id")
	}
	if subTaskID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "sub_task_id")
	}
	return lgc.changeStatus(kit, taskID, subTaskID, metadata.APITaskStatusSuccess, nil)
}

// ChangeStatusToFailure change status to failure
func (lgc *Logics) ChangeStatusToFailure(kit *rest.Kit, taskID, subTaskID string,
	errResponse *metadata.Response) error {

	if taskID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "task_id")
	}
	if subTaskID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "sub_task_id")
	}
	if errResponse == nil || errResponse.Code == 0 {
		return kit.CCError.CCError(common.CCErrTaskErrResponseemptyFail)
	}
	return lgc.changeStatus(kit, taskID, subTaskID, metadata.APITAskStatusFail, errResponse)
}

func (lgc *Logics) changeStatus(kit *rest.Kit, taskID, subTaskID string, status metadata.APITaskStatus,
	errResponse *metadata.Response) error {

	condition := mapstr.MapStr{"task_id": taskID}

	rows := make([]metadata.APITaskDetail, 0)
	err := lgc.db.Table(common.BKTableNameAPITask).Find(condition).All(kit.Ctx, &rows)
	if err != nil {
		blog.Errorf("get task status failed, input: %#v, err: %v, rid:%s", condition, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	if len(rows) == 0 {
		blog.Errorf("get task status, input: %#v, task not found, rid: %s", condition, kit.Rid)
		return kit.CCError.CCError(common.CCErrTaskNotFound)
	}

	existSubTask := false
	canChangeStatus := false
	for _, subTask := range rows[0].Detail {
		if subTask.SubTaskID == subTaskID {
			existSubTask = true
			if (subTask.Status == metadata.APITaskStatusNew || subTask.Status == metadata.APITaskStatusExecute) &&
				subTask.Status != status {
				canChangeStatus = true
			}
			break
		}
	}
	if !existSubTask {
		return kit.CCError.CCError(common.CCErrTaskSubTaskNotFound)
	}
	if !canChangeStatus {
		return kit.CCError.CCError(common.CCErrTaskStatusNotAllowChangeTo)
	}

	now := time.Now()
	updateCondition := mapstr.MapStr{"task_id": taskID, "detail.sub_task_id": subTaskID}
	updateData := mapstr.MapStr{"detail.$.status": status, common.LastTimeField: now}
	needUpdateHistory := false
	if status == metadata.APITAskStatusFail {
		// 任务的一个子任务失败，则任务失败
		updateData.Set("status", status)
		updateData.Set("detail.$.response", errResponse)

		needUpdateHistory = true
	}
	err = lgc.db.Table(common.BKTableNameAPITask).Update(kit.Ctx, updateCondition, updateData)
	if err != nil {
		blog.Errorf("change task status failed, data: %#v, err: %v, rid:%s", updateData, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBUpdateFailed)
	}

	// 如果修改为执行成功。 判断是否所有的子项都成功。如果都成功，则任务完成
	if status == metadata.APITaskStatusSuccess {
		allSuccess := true
		for _, subTask := range rows[0].Detail {
			if subTask.SubTaskID != subTaskID && subTask.Status != metadata.APITaskStatusSuccess {
				allSuccess = false
				break
			}
		}
		if allSuccess {
			updateCondition := mapstr.New()
			updateCondition.Set("task_id", taskID)
			updateData := mapstr.New()
			updateData.Set("status", metadata.APITaskStatusSuccess)
			updateData.Set(common.LastTimeField, now)

			needUpdateHistory = true

			err = lgc.db.Table(common.BKTableNameAPITask).Update(kit.Ctx, updateCondition, updateData)
			if err != nil {
				blog.Errorf("change task status failed, data: %#v, err: %v, rid:%s", updateData, err, kit.Rid)
				return kit.CCError.Error(common.CCErrCommDBUpdateFailed)
			}
		}
	}

	updateHistoryData := mapstr.MapStr{common.BKStatusField: status, common.LastTimeField: now}
	if needUpdateHistory {
		err = lgc.db.Table(common.BKTableNameAPITaskSyncHistory).Update(kit.Ctx, updateCondition, updateHistoryData)
		if err != nil {
			blog.Errorf("change history task status(%#v) failed, err: %v, rid: %s", updateHistoryData, err, kit.Rid)
			return kit.CCError.Error(common.CCErrCommDBUpdateFailed)
		}
	}
	return nil
}

// GetDBHTTPHeader TODO
func GetDBHTTPHeader(header http.Header) http.Header {

	newHeader := make(http.Header, 0)
	newHeader.Add(common.BKHTTPCCRequestID, header.Get(common.BKHTTPCCRequestID))
	newHeader.Add(common.BKHTTPCookieLanugageKey, header.Get(common.BKHTTPCookieLanugageKey))
	newHeader.Add(common.BKHTTPHeaderUser, header.Get(common.BKHTTPHeaderUser))
	newHeader.Add(common.BKHTTPLanguage, header.Get(common.BKHTTPLanguage))
	newHeader.Add(common.BKHTTPOwner, header.Get(common.BKHTTPOwner))
	newHeader.Add(common.BKHTTPOwnerID, header.Get(common.BKHTTPOwnerID))
	newHeader.Add(common.BKHTTPRequestAppCode, header.Get(common.BKHTTPRequestAppCode))
	newHeader.Add(common.BKHTTPRequestRealIP, header.Get(common.BKHTTPRequestRealIP))

	return newHeader
}

func getStrTaskID(prefix string) string {
	prefix = strings.TrimSpace(prefix)
	if prefix != "" {
		prefix = prefix + ":"
	}
	return prefix + xid.New().String()
}

// ListLatestSyncStatus list latest api task sync status
func (lgc *Logics) ListLatestSyncStatus(kit *rest.Kit, input *metadata.ListLatestSyncStatusRequest) (
	[]metadata.APITaskSyncStatus, error) {

	aggrCond := []map[string]interface{}{
		{common.BKDBSort: map[string]interface{}{common.CreateTimeField: -1}},
		{common.BKDBGroup: map[string]interface{}{
			"_id": "$bk_inst_id",
			"doc": map[string]interface{}{"$first": "$$ROOT"},
		}},
		{common.BKDBReplaceRoot: map[string]interface{}{"newRoot": "$doc"}},
	}

	var err error
	if input.TimeCondition != nil {
		input.Condition, err = input.TimeCondition.MergeTimeCondition(input.Condition)
		if err != nil {
			blog.Errorf("merge time condition failed, err: %#v, input: %s, rid: %s", err, input, kit.Rid)
			return nil, err
		}
	}
	if len(input.Condition) != 0 {
		aggrCond = append([]map[string]interface{}{{common.BKDBMatch: input.Condition}}, aggrCond...)
	}

	if len(input.Fields) != 0 {
		cond := map[string]int64{}
		for _, field := range input.Fields {
			cond[field] = 1
		}
		aggrCond = append(aggrCond, map[string]interface{}{
			common.BKDBProject: cond,
		})
	}

	result := make([]metadata.APITaskSyncStatus, 0)
	if err := lgc.db.Table(common.BKTableNameAPITaskSyncHistory).AggregateAll(kit.Ctx, aggrCond, &result); err != nil {
		blog.Errorf("list latest sync status failed, cond: %#v, err: %v, rid: %v", aggrCond, err, kit.Rid)
		return nil, err
	}

	return result, nil
}

// ListSyncStatusHistory list api task sync status history
func (lgc *Logics) ListSyncStatusHistory(kit *rest.Kit, input *metadata.QueryCondition) (
	*metadata.ListAPITaskSyncStatusResult, error) {

	var err error
	if input.TimeCondition != nil {
		input.Condition, err = input.TimeCondition.MergeTimeCondition(input.Condition)
		if err != nil {
			blog.Errorf("merge time condition failed, err: %#v, input: %s, rid: %s", err, input, kit.Rid)
			return nil, err
		}
	}

	dbQuery := lgc.db.Table(common.BKTableNameAPITaskSyncHistory).Find(input.Condition)

	if input.Page.Start != 0 {
		dbQuery.Start(uint64(input.Page.Start))
	}

	if input.Page.Limit != 0 {
		dbQuery.Limit(uint64(input.Page.Limit))
	}

	if len(input.Page.Sort) != 0 {
		dbQuery.Sort(input.Page.Sort)
	}

	result := &metadata.ListAPITaskSyncStatusResult{
		Count: 0,
		Info:  make([]metadata.APITaskSyncStatus, 0),
	}
	if !input.DisableCounter {
		count, err := dbQuery.Count(kit.Ctx)
		if err != nil {
			blog.Errorf("get sync status history count failed, input: %#v, err: %v, rid: %v", input, err, kit.Rid)
			return nil, err
		}
		if count == 0 {
			return result, nil
		}
		result.Count = int64(count)
	}

	if len(input.Fields) != 0 {
		dbQuery.Fields(input.Fields...)
	}

	if err := dbQuery.All(kit.Ctx, &result.Info); err != nil {
		blog.Errorf("list sync status history failed, input: %#v, err: %v, rid: %v", input, err, kit.Rid)
		return nil, err
	}

	return result, nil
}
