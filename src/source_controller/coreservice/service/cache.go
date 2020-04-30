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
	"configcenter/src/source_controller/coreservice/cache/topo_tree"
)

func (s *coreService) SearchTopologyTree(ctx *rest.Contexts) {
	opt := new(topo_tree.SearchOption)
	if err := ctx.DecodeInto(&opt); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := opt.Validate(); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "search topology tree, but request parameter is invalid: %v", err)
		return
	}

	topo, err := s.cache.SearchTopologyTree(opt)
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
