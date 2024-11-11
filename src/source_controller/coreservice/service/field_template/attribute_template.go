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

package fieldtmpl

import (
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	attrvalid "configcenter/src/common/valid/attribute"
	"configcenter/src/source_controller/coreservice/core/model"
	"configcenter/src/storage/driver/mongodb"
)

// ListFieldTemplateAttr list field template attributes.
func (s *service) ListFieldTemplateAttr(cts *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := cts.DecodeInto(opt); err != nil {
		cts.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		cts.RespAutoError(rawErr.ToCCError(cts.Kit.CCError))
		return
	}

	filter, err := opt.ToMgo()
	if err != nil {
		cts.RespAutoError(cts.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(filter).Count(cts.Kit.Ctx)
		if err != nil {
			blog.Errorf("count field template attr failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
			cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		cts.RespEntity(metadata.FieldTemplateAttrInfo{Count: count})
		return
	}

	attrTemplates := make([]metadata.FieldTemplateAttr, 0)
	err = mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(filter).Start(uint64(opt.Page.Start)).
		Limit(uint64(opt.Page.Limit)).Sort(opt.Page.Sort).Fields(opt.Fields...).All(cts.Kit.Ctx, &attrTemplates)
	if err != nil {
		blog.Errorf("list field template attributes failed, err: %v, filter: %+v, rid: %v", err, filter, cts.Kit.Rid)
		cts.RespAutoError(cts.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	cts.RespEntity(metadata.FieldTemplateAttrInfo{Info: attrTemplates})
}

// CreateFieldTemplateAttrs create field template attributes.
func (s *service) CreateFieldTemplateAttrs(ctx *rest.Contexts) {
	attrs := make([]metadata.FieldTemplateAttr, 0)
	if err := ctx.DecodeInto(&attrs); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(attrs) == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, "attributes"))
		return
	}

	templateID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKTemplateID), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse %s, err: %v, rid: %s", common.BKTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKTemplateID))
		return
	}
	if templateID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKTemplateID))
		return
	}

	ids, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, common.BKTableNameObjAttDesTemplate, len(attrs))
	if err != nil {
		blog.Errorf("get sequence id on the table (%s) failed, err: %v, rid: %s", common.BKTableNameObjAttDesTemplate,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error()))
		return
	}

	result := make([]int64, len(ids))
	now := time.Now()
	for idx := range attrs {
		attrs[idx].ID = int64(ids[idx])
		attrs[idx].OwnerID = ctx.Kit.SupplierAccount
		attrs[idx].Creator = ctx.Kit.User
		attrs[idx].Modifier = ctx.Kit.User
		attrs[idx].CreateTime = &metadata.Time{Time: now}
		attrs[idx].LastTime = &metadata.Time{Time: now}

		// 目前字段组合模版属性只支持枚举多选，枚举多选的默认值在option中，default值需要为nil
		if attrs[idx].PropertyType == common.FieldTypeEnumMulti && attrs[idx].Default != nil {
			attrs[idx].Default = nil
		}

		if attrs[idx].TemplateID != templateID {
			blog.Errorf("attribute template id is invalid, data: %v, template id: %d, rid: %s", attrs[idx], templateID,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, "attributes"))
			return
		}

		if err := validateFieldTemplateAttr(ctx.Kit, &attrs[idx]); err != nil {
			ctx.RespAutoError(err)
			return
		}

		if !model.SatisfyMongoFieldLimit(attrs[idx].PropertyID) {
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, common.BKPropertyIDField))
			return
		}

		result[idx] = int64(ids[idx])
	}

	if err = mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Insert(ctx.Kit.Ctx, attrs); err != nil {
		blog.Errorf("save field template attribute failed, data: %v, err: %v, rid: %s", attrs, err, ctx.Kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err)))
			return
		}
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}

	ctx.RespEntity(metadata.RspIDs{IDs: result})
}

func validateFieldTemplateAttr(kit *rest.Kit, attr *metadata.FieldTemplateAttr) error {
	if err := attr.Validate(); err.ErrCode != 0 {
		blog.Errorf("field template attribute(%+v) is invalid, err: %v, rid: %s", attr, err, kit.Rid)
		return err.ToCCError(kit.CCError)
	}

	var extraOpt interface{}
	switch attr.PropertyType {
	case common.FieldTypeEnum, common.FieldTypeEnumMulti:
		extraOpt = &attr.IsMultiple
	default:
		extraOpt = attr.Default
	}

	err := attrvalid.ValidPropertyOption(kit, attr.PropertyType, attr.Option, extraOpt)
	if err != nil {
		blog.Errorf("validate field template attribute(%+v) option failed, err: %v, rid: %s", attr, err, kit.Rid)
		return err
	}

	return nil
}

// DeleteFieldTemplateAttrs delete field template attributes
func (s *service) DeleteFieldTemplateAttrs(ctx *rest.Contexts) {
	opt := new(metadata.DeleteOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	templateID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKTemplateID), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse %s, err: %v, rid: %s", common.BKTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKTemplateID))
		return
	}
	if opt.Condition == nil {
		opt.Condition = make(map[string]interface{})
	}

	cond := opt.Condition
	cond[common.BKTemplateID] = templateID

	attrs := make([]metadata.FieldTemplateAttr, 0)
	err = mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(cond).Fields(common.BKFieldID).
		All(ctx.Kit.Ctx, &attrs)
	if err != nil {
		blog.Errorf("find field template attribute failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(attrs) == 0 {
		ctx.RespEntity(nil)
		return
	}

	attrIDs := make([]int64, 0)
	for _, attr := range attrs {
		attrIDs = append(attrIDs, attr.ID)
	}
	countCond := mapstr.MapStr{common.BKObjectUniqueKeys: mapstr.MapStr{common.BKDBIN: attrIDs}}

	count, err := mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Find(countCond).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("count field template unique failed, filter: %+v, err: %v, rid: %v", countCond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	if count != 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCoreServiceFieldTemplateHasUnique))
		return
	}

	if err := mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Delete(ctx.Kit.Ctx, cond); err != nil {
		blog.Errorf("delete field template attributes failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// UpdateFieldTemplateAttrs update field template attributes
func (s *service) UpdateFieldTemplateAttrs(ctx *rest.Contexts) {
	attrs := make([]metadata.FieldTemplateAttr, 0)
	if err := ctx.DecodeInto(&attrs); err != nil {
		ctx.RespAutoError(err)
		return
	}
	templateID, err := strconv.ParseInt(ctx.Request.PathParameter(common.BKTemplateID), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse %s, err: %v, rid: %s", common.BKTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKTemplateID))
		return
	}
	if templateID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKTemplateID))
		return
	}

	ids := make([]int64, 0)
	for _, attr := range attrs {
		if attr.ID == 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
			return
		}

		ids = append(ids, attr.ID)
	}
	if len(ids) == 0 {
		ctx.RespEntity(nil)
		return
	}

	cond := mapstr.MapStr{
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: ids},
		common.BKTemplateID: templateID,
	}
	dbTmplAttrs := make([]metadata.FieldTemplateAttr, 0)
	err = mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(cond).All(ctx.Kit.Ctx, &dbTmplAttrs)
	if err != nil {
		blog.Errorf("list field template attributes failed, filter: %+v, err: %v, rid: %v", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if len(ids) != len(dbTmplAttrs) {
		blog.Errorf("field template attributes are invalid, data: %v, err: %v, rid: %v", dbTmplAttrs, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "attributes"))
		return
	}

	dbAttrMap := make(map[int64]metadata.FieldTemplateAttr)
	for _, dbAttr := range dbTmplAttrs {
		dbAttrMap[dbAttr.ID] = dbAttr
	}

	for idx := range attrs {
		dbAttr := dbAttrMap[attrs[idx].ID]
		attrs[idx].PropertyID = dbAttr.PropertyID
		attrs[idx].PropertyType = dbAttr.PropertyType
		attrs[idx].OwnerID = dbAttr.OwnerID
		attrs[idx].Creator = dbAttr.Creator
		attrs[idx].CreateTime = dbAttr.CreateTime

		// 目前字段组合模版属性只支持枚举多选，枚举多选的默认值在option中，default值需要为nil
		if dbAttr.PropertyType == common.FieldTypeEnumMulti && attrs[idx].Default != nil {
			attrs[idx].Default = nil
		}
	}

	if err := s.updateFieldTemplateAttrs(ctx.Kit, templateID, attrs); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *service) updateFieldTemplateAttrs(kit *rest.Kit, templateID int64, attrs []metadata.FieldTemplateAttr) error {
	now := time.Now()
	for idx, attr := range attrs {
		if err := validateFieldTemplateAttr(kit, &attrs[idx]); err != nil {
			return err
		}

		attr.Modifier = kit.User
		attr.LastTime = &metadata.Time{Time: now}

		cond := mapstr.MapStr{
			common.BKFieldID:    attr.ID,
			common.BKTemplateID: templateID,
		}

		err := mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Update(kit.Ctx, cond, attr)
		if err != nil {
			blog.Errorf("update field template attribute failed, data: %v, err: %v, rid: %s", attr, err, kit.Rid)
			return err
		}
	}

	return nil
}
