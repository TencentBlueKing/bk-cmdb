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
	"strconv"
)

func (s *coreService) CreateSyncTask(ctx *rest.Contexts) {
	task := metadata.CloudSyncTask{}
	if err := ctx.DecodeInto(&task); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.CloudOperation().CreateSyncTask(ctx.Kit, &task)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) SearchSyncTask(ctx *rest.Contexts) {
	option := metadata.SearchCloudOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.CloudOperation().SearchSyncTask(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateSyncTask(ctx *rest.Contexts) {}

func (s *coreService) DeleteSyncTask(ctx *rest.Contexts) {
	taskIDStr := ctx.Request.PathParameter(common.BKCloudSyncTaskID)
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudSyncTaskID))
		return
	}

	err = s.core.CloudOperation().DeleteSyncTask(ctx.Kit, taskID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
