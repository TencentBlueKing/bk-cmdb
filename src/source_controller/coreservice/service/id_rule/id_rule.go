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

// Package idrule package
package idrule

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/storage/driver/mongodb"
)

// UpdateInstIDRule update instance id rule field
func (s *service) UpdateInstIDRule(ctx *rest.Contexts) {
	opt := new(metadata.UpdateInstIDRuleOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}
	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := mapstr.MapStr{common.BKObjIDField: opt.ObjID}
	allAttr := make([]metadata.Attribute, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(ctx.Kit.Ctx, &allAttr); err != nil {
		blog.Errorf("find attribute failed, cond: %+v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	var attr metadata.Attribute
	attrTypeMap := make(map[string]string)
	for _, attribute := range allAttr {
		if attribute.PropertyID == opt.PropertyID {
			attr = attribute
		}
		attrTypeMap[attribute.PropertyID] = attribute.PropertyType
	}

	if attr.PropertyID != opt.PropertyID {
		blog.Errorf("%s id rule attribute %s not exists, rid: %s", opt.ObjID, opt.PropertyID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKPropertyIDField))
		return
	}

	idField := common.GetInstIDField(opt.ObjID)
	cond = mapstr.MapStr{common.BKObjIDField: opt.ObjID, idField: mapstr.MapStr{common.BKDBIN: opt.IDs}}
	table := common.GetInstTableName(opt.ObjID, ctx.Kit.TenantID)
	insts := make([]mapstr.MapStr, 0)
	if err := mongodb.Client().Table(table).Find(cond).All(ctx.Kit.Ctx, &insts); err != nil {
		blog.Errorf("find instances failed, cond: %+v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	for _, inst := range insts {
		val, exist := inst.Get(opt.PropertyID)
		if exist && val != "" {
			continue
		}

		val, err := instances.GetIDRuleVal(ctx.Kit.Ctx, inst, attr, attrTypeMap)
		if err != nil {
			blog.Errorf("get id rule val failed, inst: %+v, attr: %+v, err: %v, rid: %s", inst, attr, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, err.Error()))
			return
		}

		id, exist := inst.Get(idField)
		if !exist {
			blog.Errorf("get instance %s value failed, inst: %+v, rid: %s", idField, inst, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}
		idInt64, err := util.GetInt64ByInterface(id)
		if err != nil {
			blog.Errorf("get instance %s value failed, inst: %+v, err: %v, rid: %s", idField, inst, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cond = mapstr.MapStr{common.BKObjIDField: opt.ObjID, idField: idInt64}
		data := mapstr.MapStr{opt.PropertyID: val}
		if err = mongodb.Client().Table(table).Update(ctx.Kit.Ctx, cond, data); err != nil {
			blog.Errorf("update instance failed, cond: %+v, data: %+v, err: %v, rid: %s", cond, data, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBUpdateFailed))
			return
		}
	}
	ctx.RespEntity(nil)
}
