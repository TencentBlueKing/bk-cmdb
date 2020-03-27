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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (s *Service) SearchVpc(ctx *rest.Contexts) {
	vpcOpt := new(metadata.SearchVpcOption)
	if err := ctx.DecodeInto(vpcOpt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountID)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountID))
		return
	}

	result, err := s.Logics.SearchVpc(ctx.Kit, accountID, vpcOpt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) CreateSyncTask(ctx *rest.Contexts) {
	task := new(metadata.CloudSyncTask)
	if err := ctx.DecodeInto(task); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.Logics.CreateSyncTask(ctx.Kit, task)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) SearchSyncTask(ctx *rest.Contexts) {
	option := metadata.SearchCloudOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.Logics.SearchSyncTask(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) UpdateSyncTask(ctx *rest.Contexts) {
	taskIDStr := ctx.Request.PathParameter(common.BKCloudSyncTaskID)
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudSyncTaskID))
		return
	}

	option := map[string]interface{}{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err = s.Logics.UpdateSyncTask(ctx.Kit, taskID, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) DeleteSyncTask(ctx *rest.Contexts) {
	taskIDStr := ctx.Request.PathParameter(common.BKCloudSyncTaskID)
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudSyncTaskID))
		return
	}

	err = s.Logics.DeleteSyncTask(ctx.Kit, taskID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) SearchSyncHistory(ctx *rest.Contexts) {
	option := metadata.SearchSyncHistoryOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.Logics.SearchSyncHistory(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *Service) SearchSyncRegion(ctx *rest.Contexts) {
	option := metadata.SearchSyncRegionOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.Logics.SearchSyncRegion(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}
