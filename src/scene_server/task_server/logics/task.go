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
	input.Name = strings.TrimSpace(input.Name)

	if input.Name == "" {
		return dbTask, kit.CCError.Errorf(common.CCErrCommParamsNeedString, "name")
	}

	if len(input.Data) == 0 {
		return dbTask, kit.CCError.Errorf(common.CCErrCommParamsNeedString, "data")
	}

	dbTask.TaskID = getStrTaskID("id")
	dbTask.Name = input.Name
	dbTask.User = kit.User
	dbTask.Flag = input.Flag
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
	err := lgc.db.Table(common.BKTableNameAPITask).Insert(kit.Ctx, dbTask)
	if err != nil {
		blog.Errorf("create task failed, data: %#v, err: %v, rid: %s", dbTask, err, kit.Rid)
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

	dbTasks := make([]metadata.APITaskDetail, len(tasks))
	for index, task := range tasks {
		task.Name = strings.TrimSpace(task.Name)
		if task.Name == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "name")
		}

		if len(task.Data) == 0 {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "data")
		}

		dbTask.TaskID = getStrTaskID("id")
		dbTask.Name = task.Name
		dbTask.Flag = task.Flag
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
	}

	err := lgc.db.Table(common.BKTableNameAPITask).Insert(kit.Ctx, dbTasks)
	if err != nil {
		blog.Errorf("create tasks(%#v) failed, err: %v, rid: %s", dbTasks, err, kit.Rid)
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
func (lgc *Logics) ChangeStatusToFailure(kit *rest.Kit, taskID, subTaskID string, errResponse *metadata.Response) error {

	if taskID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "task_id")
	}
	if subTaskID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "sub_task_id")
	}
	if errResponse == nil || errResponse.Code == 0 {
		return kit.CCError.CCError(common.CCErrTaskErrResponseEmtpyFail)
	}
	return lgc.changeStatus(kit, taskID, subTaskID, metadata.APITAskStatusFail, errResponse)
}

func (lgc *Logics) changeStatus(kit *rest.Kit, taskID, subTaskID string, status metadata.APITaskStatus,
	errResponse *metadata.Response) error {

	condition := mapstr.New()
	condition.Set("task_id", taskID)

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
			if (subTask.Status == metadata.APITaskStatusNew || subTask.Status == metadata.APITaskStatuExecute) &&
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

	updateCondition := mapstr.New()
	updateCondition.Set("task_id", taskID)
	updateCondition.Set("detail.sub_task_id", subTaskID)
	updateData := mapstr.New()
	updateData.Set("detail.$.status", status)
	updateData.Set(common.LastTimeField, time.Now())
	if status == metadata.APITAskStatusFail {
		// 任务的一个子任务失败，则任务失败
		updateData.Set("status", status)
		updateData.Set("detail.$.response", errResponse)
	}
	err = lgc.db.Table(common.BKTableNameAPITask).Update(kit.Ctx, updateCondition, updateData)
	if err != nil {
		blog.Errorf("change task status failed, data: %#v, err: %v, rid:%s", updateData, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
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
			updateData.Set(common.LastTimeField, time.Now())

			err = lgc.db.Table(common.BKTableNameAPITask).Update(kit.Ctx, updateCondition, updateData)
			if err != nil {
				blog.Errorf("change task status failed, data: %#v, err: %v, rid:%s", updateData, err, kit.Rid)
				return kit.CCError.Error(common.CCErrCommDBSelectFailed)
			}
		}
	}
	return nil
}

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
	return "task:" + prefix + xid.New().String()
}
