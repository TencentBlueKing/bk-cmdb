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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *Service) CreateTask(ctx *rest.Contexts) {
	input := new(metadata.CreateTaskRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	srvData := s.newSrvComm(ctx.Request.Request.Header)
	taskInfo, err := srvData.lgc.Create(srvData.ctx, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(taskInfo)
}

func (s *Service) ListTask(ctx *rest.Contexts) {

	input := new(metadata.ListAPITaskRequest)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}
	srvData := s.newSrvComm(ctx.Request.Request.Header)
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

func (s *Service) DetailTask(ctx *rest.Contexts) {
	srvData := s.newSrvComm(ctx.Request.Request.Header)
	taskInfo, err := srvData.lgc.Detail(srvData.ctx, ctx.Request.PathParameter("task_id"))
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(map[string]interface{}{"info": taskInfo})
}

func (s *Service) StatusToSuccess(ctx *rest.Contexts) {
	taskID := ctx.Request.PathParameter("task_id")
	subTaskID := ctx.Request.PathParameter("sub_task_id")

	srvData := s.newSrvComm(ctx.Request.Request.Header)
	err := srvData.lgc.ChangeStatusToSuccess(srvData.ctx, taskID, subTaskID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) StatusToFailure(ctx *rest.Contexts) {
	taskID := ctx.Request.PathParameter("task_id")
	subTaskID := ctx.Request.PathParameter("sub_task_id")
	input := &metadata.Response{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	srvData := s.newSrvComm(ctx.Request.Request.Header)
	err := srvData.lgc.ChangeStatusToFailure(srvData.ctx, taskID, subTaskID, input)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
