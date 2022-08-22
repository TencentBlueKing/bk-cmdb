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
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/kube/types"
	"configcenter/src/storage/driver/mongodb"
)

// FindInst find instance with table name and condition
func (s *coreService) FindInst(ctx *rest.Contexts) {
	req := types.QueryReq{}
	if err := ctx.DecodeInto(&req); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := req.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := req.Condition
	instItems := make([]mapstr.MapStr, 0)
	err := mongodb.Client().Table(req.Table).Find(cond.Condition).Start(uint64(cond.Page.Start)).
		Limit(uint64(cond.Page.Limit)).Sort(cond.Page.Sort).Fields(cond.Fields...).All(ctx.Kit.Ctx, &instItems)
	if err != nil {
		blog.Errorf("search instance failed, table: %s, cond: %v, rid: %s", req.Table, cond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := &metadata.QueryResult{
		Info: instItems,
	}
	ctx.RespEntity(result)
}
