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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// GetDistinctField TODO
func (s *coreService) GetDistinctField(ctx *rest.Contexts) {
	option := new(metadata.DistinctFieldOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	ret, err := s.core.CommonOperation().GetDistinctField(ctx.Kit, option)

	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(ret)
}

// GetDistinctCount 根据条件获取指定表中满足条件数据的数量
func (s *coreService) GetDistinctCount(ctx *rest.Contexts) {
	option := new(metadata.DistinctFieldOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	rawErr := option.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	count, err := s.core.CommonOperation().GetDistinctCount(ctx.Kit, option)

	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(count)
}

// GroupRelResByIDs group related resource by ids
func (s *coreService) GroupRelResByIDs(cts *rest.Contexts) {
	opt := new(metadata.GroupRelResByIDsOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	kind := metadata.GroupByResKind(cts.Request.PathParameter("kind"))
	tableName, exists := metadata.CountResKindTableMap[kind]
	if !exists {
		blog.Errorf("%s kind is invalid, rid: %s", kind, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "kind"))
		return
	}

	filter := mapstr.MapStr{
		opt.IDField: mapstr.MapStr{common.BKDBIN: opt.IDs},
	}

	if len(opt.ExtraCond) > 0 {
		filter = mapstr.MapStr{
			common.BKDBAND: []mapstr.MapStr{filter, opt.ExtraCond},
		}
	}

	resources := make([]mapstr.MapStr, 0)
	err := mongodb.Client().Table(tableName).Find(filter).Fields(opt.IDField, opt.RelField).All(cts.Kit.Ctx, &resources)
	if err != nil {
		blog.Errorf("get all %s resource by filter(%+v) failed, err: %v, rid: %s", tableName, filter, err, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	result := make(map[int64][]interface{})
	for _, res := range resources {
		id, err := util.GetInt64ByInterface(res[opt.IDField])
		if err != nil {
			blog.Errorf("parse res(%+v) %s field failed, err: %v, rid: %s", res, opt.IDField, err, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, opt.IDField))
			return
		}
		result[id] = append(result[id], res[opt.RelField])
	}

	cts.RespEntity(result)
}
