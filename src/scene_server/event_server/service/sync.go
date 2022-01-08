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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/event_server/sync/hostidentifier"
)

// SyncHostIdentifier sync host identifier, add hostInfo message to redis fail host list
func (s *Service) SyncHostIdentifier(ctx *rest.Contexts) {
	hostInfo := new(hostidentifier.HostInfo)
	if err := ctx.DecodeInto(&hostInfo); err != nil {
		blog.Errorf("decode request body err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result, err := s.engine.CoreAPI.CoreService().Host().GetHostByID(ctx.Kit.Ctx, ctx.Kit.Header, hostInfo.HostID)
	if err != nil {
		blog.Errorf("sync hostIdentifier http do error, hostID: %d, err: %v, rid: %s",
			hostInfo.HostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed))
		return

	}
	if !result.Result {
		blog.Errorf("get host by id error, hostID: %d, err: %v, rid: %s", hostInfo.HostID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(result.Code, result.ErrMsg))
		return
	}

	innerIP := util.GetStrByInterface(result.Data[common.BKHostInnerIPField])
	if innerIP == "" {
		blog.Errorf("the host value have not innerIP, host: %v, rid: %s", result.Data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommNotFound))
		return
	}

	cloudID, err := util.GetInt64ByInterface(result.Data[common.BKCloudIDField])
	if err != nil {
		blog.Errorf("the host value have not cloudID, host: %v, rid: %s", result.Data, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	hostInfo.CloudID = cloudID
	hostInfo.HostInnerIP = innerIP

	if err := s.cache.LPush(ctx.Kit.Ctx, hostidentifier.RedisFailHostListName, hostInfo).Err(); err != nil {
		blog.Errorf("add hostInfo to redis list error, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
