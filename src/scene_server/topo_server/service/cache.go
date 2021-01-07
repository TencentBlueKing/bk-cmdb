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
	"configcenter/src/source_controller/cacheservice/cache/topo_tree"
)

func (s *Service) SearchTopologyTree(ctx *rest.Contexts) {
	opt := new(topo_tree.SearchOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	topo, err := s.Engine.CoreAPI.CacheService().Cache().Topology().SearchTopologyTree(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(topo)
}

func (s *Service) SearchTopologyNodePath(ctx *rest.Contexts) {
	opt := new(topo_tree.SearchNodePathOption)
	appID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid biz id")
		return
	}

	if appID <= 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommParamsIsInvalid, "invalid biz id")
		return
	}

	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if len(opt.Nodes) > common.BKMaxPageSize {
		ctx.RespErrorCodeF(common.CCErrCommValExceedMaxFailed, "bk_nodes counts exceeded", "bk_nodes", common.BKMaxPageSize)
		return
	}

	opt.Business = appID

	paths, err := s.Engine.CoreAPI.CacheService().Cache().Topology().SearchTopologyNodePath(ctx.Kit.Ctx, ctx.Kit.Header, opt)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(paths)
}
