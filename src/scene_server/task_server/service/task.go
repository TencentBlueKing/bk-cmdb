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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

const step uint64 = 1000

// CreateTask create a task
func (s *Service) CreateTask(ctx *rest.Contexts) {
	input := new(metadata.CreateTaskRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	srvData := s.newSrvComm(ctx.Kit.Header)
	taskInfo, err := srvData.lgc.Create(srvData.ctx, input)
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
	srvData := s.newSrvComm(ctx.Kit.Header)
	infos, cnt, err := srvData.lgc.List(srvData.ctx, ctx.Request.PathParameter("name"), input)
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

	srvData := s.newSrvComm(ctx.Kit.Header)
	infos, err := srvData.lgc.ListLatestTask(srvData.ctx, ctx.Request.PathParameter("name"), input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(infos)
}

// DetailTask show a task detail
func (s *Service) DetailTask(ctx *rest.Contexts) {
	srvData := s.newSrvComm(ctx.Kit.Header)
	taskInfo, err := srvData.lgc.Detail(srvData.ctx, ctx.Request.PathParameter("task_id"))
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{"info": taskInfo})
}

// DeleteTask delete task by condition
func (s *Service) DeleteTask(ctx *rest.Contexts) {
	srvData := s.newSrvComm(ctx.Kit.Header)

	input := new(metadata.DeleteOption)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err := srvData.lgc.DeleteTask(srvData.ctx, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(common.CCSuccessStr)
}

// StatusToSuccess change the task which task_id in path status to success
func (s *Service) StatusToSuccess(ctx *rest.Contexts) {
	taskID := ctx.Request.PathParameter("task_id")
	subTaskID := ctx.Request.PathParameter("sub_task_id")

	srvData := s.newSrvComm(ctx.Kit.Header)
	err := srvData.lgc.ChangeStatusToSuccess(srvData.ctx, taskID, subTaskID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// StatusToSuccess change the task which task_id in path status to failure
func (s *Service) StatusToFailure(ctx *rest.Contexts) {
	taskID := ctx.Request.PathParameter("task_id")
	subTaskID := ctx.Request.PathParameter("sub_task_id")
	input := &metadata.Response{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	srvData := s.newSrvComm(ctx.Kit.Header)
	err := srvData.lgc.ChangeStatusToFailure(srvData.ctx, taskID, subTaskID, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// TimerDeleteHistoryTask delete apitask history message
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
	}
}

// deleteRedundancyTask delete redundancy tasks from a month ago
func (s *Service) deleteRedundancyTask(ctx context.Context, rid string) error {

	for {
		cond := map[string]interface{}{
			common.LastTimeField: map[string]interface{}{
				common.BKDBLT: time.Now().AddDate(0, -1, 0),
			},
		}

		tasks := make([]metadata.APITaskDetail, 0)
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
