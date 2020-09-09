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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// 新建云账户
func (s *coreService) CreateAccount(ctx *rest.Contexts) {
	account := metadata.CloudAccount{}
	if err := ctx.DecodeInto(&account); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.CloudOperation().CreateAccount(ctx.Kit, &account)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// 查询云账户
func (s *coreService) SearchAccount(ctx *rest.Contexts) {
	option := metadata.SearchCloudOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.CloudOperation().SearchAccount(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// 更新云账户
func (s *coreService) UpdateAccount(ctx *rest.Contexts) {
	//get accountID
	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountID)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountID))
		return
	}

	option := mapstr.MapStr{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	err = s.core.CloudOperation().UpdateAccount(ctx.Kit, accountID, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// 删除云账户
func (s *coreService) DeleteAccount(ctx *rest.Contexts) {
	//get accountID
	accountIDStr := ctx.Request.PathParameter(common.BKCloudAccountID)
	accountID, err := strconv.ParseInt(accountIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKCloudAccountID))
		return
	}

	err = s.core.CloudOperation().DeleteAccount(ctx.Kit, accountID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// 查询云厂商账户配置
func (s *coreService) SearchAccountConf(ctx *rest.Contexts) {
	option := metadata.SearchCloudOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.CloudOperation().SearchAccountConf(ctx.Kit, &option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}
