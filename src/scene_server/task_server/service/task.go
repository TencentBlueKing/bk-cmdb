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

package service

import (
	"context"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateTask create a task
func (s *Service) CreateTask(ctx *rest.Contexts) {
	input := new(metadata.CreateTaskRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	taskInfo, err := s.Logics.Create(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(taskInfo)
}

// CreateTaskBatch create task batch
func (s *Service) CreateTaskBatch(ctx *rest.Contexts) {
	input := make([]metadata.CreateTaskRequest, 0)
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	taskInfo, err := s.Logics.CreateBatch(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(taskInfo)
}

// ListTask list the task by input condition
func (s *Service) ListTask(ctx *rest.Contexts) {

	input := new(metadata.ListAPITaskRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	infos, cnt, err := s.Logics.List(ctx.Kit, ctx.Request.PathParameter("name"), input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(metadata.ListAPITaskData{
		Info:  infos,
		Count: int64(cnt),
		Page:  input.Page,
	})
}

// ListLatestTask list the latest task
func (s *Service) ListLatestTask(ctx *rest.Contexts) {
	input := new(metadata.ListAPITaskLatestRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	infos, err := s.Logics.ListLatestTask(ctx.Kit, ctx.Request.PathParameter("name"), input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(infos)
}

// DetailTask show a task detail
func (s *Service) DetailTask(ctx *rest.Contexts) {
	taskInfo, err := s.Logics.Detail(ctx.Kit, ctx.Request.PathParameter("task_id"))
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{"info": taskInfo})
}

// DeleteTask delete task by condition
func (s *Service) DeleteTask(ctx *rest.Contexts) {

	input := new(metadata.DeleteOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := s.Logics.DeleteTask(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(common.CCSuccessStr)
}

// ListLatestSyncStatus list latest api task sync status
func (s *Service) ListLatestSyncStatus(ctx *rest.Contexts) {
	input := new(metadata.ListLatestSyncStatusRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	infos, err := s.Logics.ListLatestSyncStatus(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(infos)
}

// ListSyncStatusHistory list api task sync status history
func (s *Service) ListSyncStatusHistory(ctx *rest.Contexts) {
	input := new(metadata.QueryCondition)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	infos, err := s.Logics.ListSyncStatusHistory(ctx.Kit, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(infos)
}

// ListFieldTemplateTasksStatus query task status by field template ID and object id
// get the status of all related tasks in reverse order of the time they were created,
// if there are "in progress" tasks. Prioritize returning a status of "in progress".
// if there is no "in progress" task status, return the latest task status, which may
// be "successful", "failed" or "queuing", etc.
func (s *Service) ListFieldTemplateTasksStatus(ctx *rest.Contexts) {
	input := new(metadata.ListFieldTmpltTaskStatusOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	objIDs := util.IntArrayUnique(input.ObjectIDs)
	query := mapstr.MapStr{
		common.BKInstIDField: input.ID,
		metadata.APITaskExtraField: mapstr.MapStr{
			common.BKDBIN: objIDs,
		},
		common.BKTaskTypeField: common.SyncFieldTemplateTaskFlag,
	}

	infos, err := s.Logics.ListFieldTemplateSyncStatus(ctx.Kit, query)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	objStatusMap := make(map[int64]string)
	for _, info := range infos {
		obj, err := util.GetInt64ByInterface(info.Extra)
		if err != nil {
			blog.Errorf("get instance id failed, obj: %+v, err: %v, rid: %s", obj, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		_, ok := objStatusMap[obj]
		if ok && string(info.Status) != string(metadata.APITaskStatusExecute) {
			continue
		}
		objStatusMap[obj] = string(info.Status)
	}

	result := make([]metadata.ListFieldTmpltTaskStatusResult, 0)
	for obj, status := range objStatusMap {
		result = append(result, metadata.ListFieldTmpltTaskStatusResult{
			ObjectID: obj,
			Status:   status,
		})
	}

	ctx.RespEntity(result)
}

// TimerDeleteHistoryTask delete api task history message
func (s *Service) TimerDeleteHistoryTask(ctx context.Context) {
	for {
		time.Sleep(time.Hour * 24)

		isMaster := s.Engine.ServiceManageInterface.IsMaster()
		if !isMaster {
			continue
		}

		rid := util.GenerateRID()

		blog.Infof("begin delete redundancy task, time: %v, rid: %s", time.Now(), rid)
		err := s.deleteRedundancyTask(ctx, rid)
		if err != nil {
			blog.Errorf("delete redundancy task failed, err: %v, rid: %s", err, rid)
			continue
		}
		blog.Infof("delete redundancy task completed, time: %v, rid: %s", time.Now(), rid)

		blog.Infof("begin delete redundancy task history, time: %v, rid: %s", time.Now(), rid)
		err = s.deleteRedundancyTaskHistory(ctx, rid)
		if err != nil {
			blog.Errorf("delete redundancy task history failed, err: %v, rid: %s", err, rid)
			continue
		}
		blog.Infof("delete redundancy task history completed, time: %v, rid: %s", time.Now(), rid)
	}
}

type taskIDInfo struct {
	TaskID string `bson:"task_id"`
}

// deleteRedundancyTask delete redundancy tasks from two month ago
func (s *Service) deleteRedundancyTask(ctx context.Context, rid string) error {
	cond := map[string]interface{}{
		common.LastTimeField: map[string]interface{}{
			common.BKDBLT: time.Now().AddDate(0, -2, 0),
		},
	}

	for {
		tasks := make([]taskIDInfo, 0)
		err := s.DB.Table(common.BKTableNameAPITask).Find(cond).Fields(common.BKTaskIDField).Limit(100).All(ctx, &tasks)
		if err != nil {
			blog.Errorf("get one month ago tasks failed, err: %v, cond: %#v, rid: %s", err, cond, rid)
			return err
		}

		if len(tasks) == 0 {
			blog.Infof("found no redundancy tasks, rid: %s", rid)
			return nil
		}

		var taskIDs []string
		for _, task := range tasks {
			taskIDs = append(taskIDs, task.TaskID)
		}

		deleteCond := map[string]interface{}{
			common.BKTaskIDField: map[string]interface{}{common.BKDBIN: taskIDs},
		}

		if err := s.DB.Table(common.BKTableNameAPITask).Delete(ctx, deleteCond); err != nil {
			blog.Errorf("delete redundancy task failed, err: %v, rid: %s", err, rid)
			return err
		}

		blog.Infof("delete %d redundancy tasks successful, rid: %s", len(tasks), rid)
		time.Sleep(time.Second * 20)
	}
}

// deleteRedundancyTaskHistory delete redundancy task history from two month ago
func (s *Service) deleteRedundancyTaskHistory(ctx context.Context, rid string) error {
	cond := map[string]interface{}{
		common.LastTimeField: map[string]interface{}{
			common.BKDBLT: time.Now().AddDate(-1, 0, 0),
		},
	}

	for {
		tasks := make([]taskIDInfo, 0)
		err := s.DB.Table(common.BKTableNameAPITaskSyncHistory).Find(cond).Fields(common.BKTaskIDField).Limit(100).
			All(ctx, &tasks)
		if err != nil {
			blog.Errorf("get two month ago task history failed, err: %v, cond: %#v, rid: %s", err, cond, rid)
			return err
		}

		if len(tasks) == 0 {
			blog.Infof("found no redundancy task history, rid: %s", rid)
			return nil
		}

		var taskIDs []string
		for _, task := range tasks {
			taskIDs = append(taskIDs, task.TaskID)
		}

		deleteCond := map[string]interface{}{
			common.BKTaskIDField: map[string]interface{}{common.BKDBIN: taskIDs},
		}

		if err := s.DB.Table(common.BKTableNameAPITaskSyncHistory).Delete(ctx, deleteCond); err != nil {
			blog.Errorf("delete redundancy task failed, err: %v, rid: %s", err, rid)
			return err
		}

		blog.Infof("delete %d redundancy task history successful, rid: %s", len(tasks), rid)
		time.Sleep(time.Second * 20)
	}
}
