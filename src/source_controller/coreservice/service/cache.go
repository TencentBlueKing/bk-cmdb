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
	"configcenter/src/source_controller/coreservice/cache/topo_tree"
)

func (s *coreService) SearchTopologyTreeInCache(ctx *rest.Contexts) {
	opt := new(topo_tree.SearchOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "search topology tree, but request parameter is invalid: %v", err)
		return
	}

	topo, err := s.cacheSet.Topology.SearchTopologyTree(opt)
	if err != nil {
		if err == topo_tree.OverHeadError {
			ctx.RespWithError(err, common.SearchTopoTreeScanTooManyData, "search topology tree failed, err: %v", err)
			return
		}
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search topology tree failed, err: %v", err)
		return
	}
	ctx.RespEntity(topo)
}

func (s *coreService) SearchHostWithInnerIPInCache(ctx *rest.Contexts) {
	opt := new(metadata.SearchHostWithInnerIPOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}
	host, err := s.cacheSet.Host.GetHostWithInnerIP(ctx.Kit.Ctx, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search host with inner ip in cache, but get host failed, err: %v", err)
		return
	}
	ctx.RespString(host)
}

func (s *coreService) SearchHostWithHostIDInCache(ctx *rest.Contexts) {
	opt := new(metadata.SearchHostWithIDOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	host, err := s.cacheSet.Host.GetHostWithID(ctx.Kit.Ctx, opt)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search host with id in cache, but get host failed, err: %v", err)
		return
	}
	ctx.RespString(host)
}

func (s *coreService) SearchBusinessInCache(ctx *rest.Contexts) {
	bizID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKAppIDField), 10, 64)
	biz, err := s.cacheSet.Business.GetBusiness(bizID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search biz with id in cache, but get biz failed, err: %v", err)
		return
	}
	ctx.RespString(biz)
}

func (s *coreService) SearchSetInCache(ctx *rest.Contexts) {
	setID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKSetIDField), 10, 64)
	set, err := s.cacheSet.Business.GetSet(setID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search set with id in cache failed, err: %v", err)
		return
	}
	ctx.RespString(set)
}

func (s *coreService) SearchModuleInCache(ctx *rest.Contexts) {
	moduleID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKModuleIDField), 10, 64)
	module, err := s.cacheSet.Business.GetModule(moduleID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search module with id in cache failed, err: %v", err)
		return
	}
	ctx.RespString(module)
}

func (s *coreService) SearchCustomLayerInCache(ctx *rest.Contexts) {
	objID := ctx.Request.PathParameter(common.BKObjIDField)
	instID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKInstIDField), 10, 64)
	inst, err := s.cacheSet.Business.GetCustomLevelDetail(objID, instID)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommDBSelectFailed, "search custom layer with id in cache failed, err: %v", err)
		return
	}
	ctx.RespString(inst)
}
