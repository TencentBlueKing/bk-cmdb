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
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
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

// CreateFieldTemplateTask create field template task batch
func (s *Service) CreateFieldTemplateTask(ctx *rest.Contexts) {
	input := make([]metadata.CreateTaskRequest, 0)
	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	taskInfo, err := s.Logics.CreateFieldTemplateBatch(ctx.Kit, input)
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

// ListFieldTmplTaskSyncResult list field template task sync result, it will first query the task table,
// and if it cannot find it, it will query the task history table.
func (s *Service) ListFieldTmplTaskSyncResult(ctx *rest.Contexts) {
	input := new(metadata.ListFieldTmplSyncTaskStatusOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := input.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	input.ObjectIDs = util.IntArrayUnique(input.ObjectIDs)
	taskDetailMap, taskHistoryMap, err := s.Logics.ListFieldTemplateSyncResult(ctx.Kit, input.ID, input.ObjectIDs)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	result := make([]metadata.ListFieldTmplTaskSyncResult, 0)
	for _, id := range input.ObjectIDs {
		taskDetail, exist := taskDetailMap[id]
		if exist {
			syncResult := metadata.ListFieldTmplTaskSyncResult{
				ObjectID: id,
				Status:   taskDetail.Status,
				SyncTime: taskDetail.LastTime,
			}
			if len(taskDetail.Detail) != 0 && taskDetail.Detail[0].Response != nil {
				syncResult.FailMsg = taskDetail.Detail[0].Response.ErrMsg
			}

			result = append(result, syncResult)
			continue
		}

		taskHistory, exist := taskHistoryMap[id]
		if exist {
			result = append(result, metadata.ListFieldTmplTaskSyncResult{
				ObjectID: id,
				Status:   taskHistory.Status,
				SyncTime: taskHistory.LastTime,
			})
		}
	}

	ctx.RespEntity(result)
}
