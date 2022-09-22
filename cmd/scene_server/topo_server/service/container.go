/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package service

import (
	"strconv"

	"configcenter/pkg/blog"
	"configcenter/pkg/common"
	"configcenter/pkg/http/rest"
	types2 "configcenter/pkg/kube/types"
	"configcenter/pkg/mapstr"
	"configcenter/pkg/metadata"
)

// ListContainer list container
func (s *Service) ListContainer(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	req := types2.ContainerQueryReq{}
	if err := ctx.DecodeInto(&req); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := mapstr.MapStr{
		common.BKAppIDField: bizID,
		common.BKFieldID:    req.PodID,
	}
	counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
		types2.BKTableNameBasePod, []map[string]interface{}{cond})
	if err != nil {
		blog.Errorf("get pod failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	if counts[0] != 1 {
		blog.Errorf("get pod failed, count: %d, cond: %v, err: %v, rid: %s", counts[0], cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, types2.BKPodIDField))
		return
	}

	cond, err = req.BuildCond()
	if err != nil {
		blog.Errorf("build query container condition failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if req.Page.EnableCount {
		counts, err := s.Engine.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header,
			types2.BKTableNameBaseContainer, []map[string]interface{}{cond})
		if err != nil {
			blog.Errorf("count container failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntityWithCount(counts[0], make([]mapstr.MapStr, 0))
		return
	}

	if req.Page.Sort == "" {
		req.Page.Sort = common.BKFieldID
	}

	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      req.Page,
		Fields:    req.Fields,
	}

	resp, err := s.Engine.CoreAPI.CoreService().Kube().ListContainer(ctx.Kit.Ctx, ctx.Kit.Header, query)
	if err != nil {
		blog.Errorf("find container failed, cond: %v, err: %v, rid: %s", query, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntityWithCount(0, resp.Info)
}
