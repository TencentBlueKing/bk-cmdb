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
	"fmt"

	"configcenter/pkg/inst/logics"
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
	if err := mongodb.Shard(ctx.Kit.ShardOpts()).Table(common.BKTableNameObjAttDes).Find(cond).All(ctx.Kit.Ctx,
		&allAttr); err != nil {
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

	if err := s.updateInsts(ctx.Kit, opt, attr, attrTypeMap); err != nil {
		blog.Errorf("update instance %s id rule failed, err: %v, rid: %s", opt.PropertyID, err, ctx.Kit.Rid)
		ctx.RespAutoError(fmt.Errorf("update instance %s id rule failed", opt.PropertyID))
		return
	}

	ctx.RespEntity(nil)
}

func (s *service) updateInsts(kit *rest.Kit, opt *metadata.UpdateInstIDRuleOption, attr metadata.Attribute,
	attrTypeMap map[string]string) error {

	idField := common.GetInstIDField(opt.ObjID)
	cond := mapstr.MapStr{common.BKObjIDField: opt.ObjID, idField: mapstr.MapStr{common.BKDBIN: opt.IDs}}

	table, err := logics.GetObjInstTableFromCache(kit, s.clientSet, opt.ObjID)
	if err != nil {
		blog.Errorf("get object(%s) instance table name failed, err: %v, rid: %s", opt.ObjID, err, kit.Rid)
		return err
	}

	insts := make([]mapstr.MapStr, 0)
	if err = mongodb.Shard(kit.ShardOpts()).Table(table).Find(cond).All(kit.Ctx, &insts); err != nil {
		blog.Errorf("find instances failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	for _, inst := range insts {
		val, exist := inst.Get(opt.PropertyID)
		if exist && val != "" {
			continue
		}

		val, err := instances.GetIDRuleVal(kit, inst, attr, attrTypeMap)
		if err != nil {
			blog.Errorf("get id rule val failed, inst: %+v, attr: %+v, err: %v, rid: %s", inst, attr, err, kit.Rid)
			return err
		}

		id, exist := inst.Get(idField)
		if !exist {
			blog.Errorf("get instance %s value failed, inst: %+v, rid: %s", idField, inst, kit.Rid)
			return err
		}
		idInt64, err := util.GetInt64ByInterface(id)
		if err != nil {
			blog.Errorf("get instance %s value failed, inst: %+v, err: %v, rid: %s", idField, inst, err, kit.Rid)
			return err
		}

		cond = mapstr.MapStr{common.BKObjIDField: opt.ObjID, idField: idInt64}
		data := mapstr.MapStr{opt.PropertyID: val}
		if err = mongodb.Shard(kit.ShardOpts()).Table(table).Update(kit.Ctx, cond, data); err != nil {
			blog.Errorf("update instance failed, cond: %+v, data: %+v, err: %v, rid: %s", cond, data, err, kit.Rid)
			return err
		}
	}
	return nil
}
