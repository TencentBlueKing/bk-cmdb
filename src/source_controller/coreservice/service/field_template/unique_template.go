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
	"configcenter/src/common/util"
	"configcenter/src/storage/driver/mongodb"
)

// CreateFieldTemplateUniques create field template uniques.
func (s *service) CreateFieldTemplateUniques(ctx *rest.Contexts) {
	uniques := make([]metadata.FieldTemplateUnique, 0)
	if err := ctx.DecodeInto(&uniques); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(uniques) == 0 {
		ctx.RespEntity(metadata.RspIDs{IDs: make([]int64, 0)})
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

	if err := s.checkUniques(ctx.Kit, uniques, templateID); err != nil {
		blog.Errorf("check field template uniques failed, uniques: %v, err: %v, rid: %s", uniques, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ids, err := mongodb.Client().NextSequences(ctx.Kit.Ctx, common.BKTableNameObjectUniqueTemplate, len(uniques))
	if err != nil {
		blog.Errorf("get sequence id on the table (%s) failed, err: %v, rid: %s",
			common.BKTableNameObjectUniqueTemplate, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error()))
		return
	}

	result := make([]int64, len(ids))
	now := time.Now()
	for idx := range uniques {
		uniques[idx].ID = int64(ids[idx])
		uniques[idx].OwnerID = ctx.Kit.SupplierAccount
		uniques[idx].Creator = ctx.Kit.User
		uniques[idx].Modifier = ctx.Kit.User
		uniques[idx].CreateTime = &metadata.Time{Time: now}
		uniques[idx].LastTime = &metadata.Time{Time: now}

		result[idx] = int64(ids[idx])
	}

	if err = mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Insert(ctx.Kit.Ctx, uniques); err != nil {
		blog.Errorf("save field template unique failed, data: %v, err: %v, rid: %s", uniques, err, ctx.Kit.Rid)
		if mongodb.Client().IsDuplicatedError(err) {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, mongodb.GetDuplicateKey(err)))
			return
		}
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBInsertFailed))
		return
	}

	ctx.RespEntity(metadata.RspIDs{IDs: result})
}

// checkUniques check if same unique rule has existed. issue #5240
// the method parameter uniques need to be checked with each other, and be checked with the uniques in db.
func (s *service) checkUniques(kit *rest.Kit, uniques []metadata.FieldTemplateUnique, templateID int64) error {
	if len(uniques) == 0 {
		return nil
	}

	attrIDs := make([]int64, 0)
	for idx, unique := range uniques {
		if unique.TemplateID != templateID {
			blog.Errorf("unique template id is invalid, data: %v, template id: %d, rid: %s", unique, templateID,
				kit.Rid)
			return kit.CCError.New(common.CCErrCommParamsInvalid, "uniques")
		}

		if err := unique.Validate(); err.ErrCode != 0 {
			return err.ToCCError(kit.CCError)
		}

		if err := s.checkKeys(kit, unique, uniques, idx); err != nil {
			return err
		}

		attrIDs = append(attrIDs, unique.Keys...)
	}

	attrIDs = util.IntArrayUnique(attrIDs)
	if err := s.isTmplUniquesLegal(kit, attrIDs, uniques); err != nil {
		return err
	}

	existedUniques, err := s.findUniqueByTemplateID(kit, templateID)
	if err != nil {
		blog.Errorf("find uniques by template id failed, id: %v, err: %v, rid: %s", templateID, err, kit.Rid)
		return err
	}

	for _, unique := range uniques {
		if err := s.checkKeys(kit, unique, existedUniques, noNeedSkip); err != nil {
			return err
		}
	}

	return nil
}

const noNeedSkip = -1

func (s *service) checkKeys(kit *rest.Kit, target metadata.FieldTemplateUnique,
	comparators []metadata.FieldTemplateUnique, skipIdx int) error {

	keysMap := make(map[int64]struct{})
	for _, key := range target.Keys {
		if _, exists := keysMap[key]; exists {
			blog.Errorf("unique key is invalid, unique: %v, rid: %s", target, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
		}

		keysMap[key] = struct{}{}
	}

	for idx, comparator := range comparators {
		if idx == skipIdx || target.TemplateID != comparator.TemplateID ||
			(target.ID != 0 && target.ID == comparator.ID) {
			continue
		}

		cnt := 0
		for _, existedKey := range comparator.Keys {
			if _, exist := keysMap[existedKey]; exist {
				cnt++
			}
		}

		if len(keysMap) == cnt || (len(keysMap) > len(comparator.Keys) && len(comparator.Keys) == cnt) {
			blog.Errorf("unique key is invalid, unique: %v, rid: %s", target, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
		}
	}

	return nil
}

func (s *service) isTmplUniquesLegal(kit *rest.Kit, attrIDs []int64, uniques []metadata.FieldTemplateUnique) error {
	cond := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{common.BKDBIN: attrIDs},
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	fields := []string{common.BKFieldID, common.BKPropertyTypeField}

	attrs := make([]metadata.FieldTemplateAttr, 0)
	err := mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(cond).Fields(fields...).All(kit.Ctx, &attrs)
	if err != nil {
		blog.Errorf("find template attributes failed, cond: %v, err: %v, rid: %v", cond, err, kit.Rid)
		return err
	}

	if len(attrs) != len(attrIDs) {
		blog.Errorf("can not find all attributes, cond: %v, attrs count: %d, ids count: %d, rid: %v", cond, len(attrs),
			len(attrIDs), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
	}

	attrIDToType := make(map[int64]string)
	for _, attr := range attrs {
		attrIDToType[attr.ID] = attr.PropertyType
	}

	for _, unique := range uniques {
		keys := unique.Keys
		if len(keys) == 1 {
			keyType := attrIDToType[keys[0]]
			if keyType != common.FieldTypeSingleChar && keyType != common.FieldTypeInt && keyType !=
				common.FieldTypeFloat {

				blog.Errorf("unique attribute type is invalid, attr: %v, rid: %v", unique, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
			}

			continue
		}

		for _, key := range keys {
			keyType := attrIDToType[key]
			if keyType != common.FieldTypeSingleChar && keyType != common.FieldTypeInt && keyType !=
				common.FieldTypeFloat && keyType != common.FieldTypeDate && keyType != common.FieldTypeList {

				blog.Errorf("unique attribute type is invalid, attr: %v, rid: %v", unique, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
			}
		}
	}

	return nil
}

func (s *service) findUniqueByTemplateID(kit *rest.Kit, id int64) ([]metadata.FieldTemplateUnique, error) {
	cond := map[string]interface{}{
		common.BKTemplateID: id,
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	uniques := make([]metadata.FieldTemplateUnique, 0)

	err := mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Find(cond).All(kit.Ctx, &uniques)
	if err != nil {
		blog.Errorf("find field template uniques failed, cond: %v, err: %v, rid: %v", cond, err, kit.Rid)
		return nil, err
	}

	return uniques, nil
}

// ListFieldTemplateUnique list field template unique.
func (s *service) ListFieldTemplateUnique(ctx *rest.Contexts) {
	opt := new(metadata.CommonQueryOption)
	if err := ctx.DecodeInto(opt); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if rawErr := opt.Validate(); rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	filter, err := opt.ToMgo()
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, err.Error()))
		return
	}

	filter = util.SetQueryOwner(filter, ctx.Kit.SupplierAccount)

	if opt.Page.EnableCount {
		count, err := mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Find(filter).Count(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("count field template uniques failed, err: %v, filter: %+v, rid: %v", err, filter, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
			return
		}

		ctx.RespEntity(metadata.FieldTemplateInfo{Count: count})
		return
	}

	uniques := make([]metadata.FieldTemplateUnique, 0)
	err = mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Find(filter).Start(uint64(opt.Page.Start)).
		Limit(uint64(opt.Page.Limit)).Sort(opt.Page.Sort).Fields(opt.Fields...).All(ctx.Kit.Ctx, &uniques)
	if err != nil {
		blog.Errorf("list field template uniques failed, err: %v, filter: %+v, rid: %v", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(metadata.FieldTemplateUniqueInfo{Info: uniques})
}

// DeleteFieldTemplateUniques delete field template uniques
func (s *service) DeleteFieldTemplateUniques(ctx *rest.Contexts) {
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
	if templateID == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKTemplateID))
		return
	}

	cond := util.SetModOwner(opt.Condition, ctx.Kit.SupplierAccount)
	cond[common.BKTemplateID] = templateID

	if err := mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Delete(ctx.Kit.Ctx, cond); err != nil {
		blog.Errorf("delete field template uniques failed, cond: %v, err: %v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// UpdateFieldTemplateUniques update field template uniques
func (s *service) UpdateFieldTemplateUniques(ctx *rest.Contexts) {
	uniques := make([]metadata.FieldTemplateUnique, 0)
	if err := ctx.DecodeInto(&uniques); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if len(uniques) == 0 {
		ctx.RespEntity(nil)
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
	for _, unique := range uniques {
		if unique.ID == 0 {
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKFieldID))
			return
		}

		ids = append(ids, unique.ID)
	}

	cond := mapstr.MapStr{
		common.BKFieldID:    mapstr.MapStr{common.BKDBIN: ids},
		common.BKTemplateID: templateID,
	}
	cond = util.SetModOwner(cond, ctx.Kit.SupplierAccount)
	dbTmplUniques := make([]metadata.FieldTemplateUnique, 0)
	err = mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Find(cond).All(ctx.Kit.Ctx, &dbTmplUniques)
	if err != nil {
		blog.Errorf("list field template uniques failed, filter: %+v, err: %v, rid: %v", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	if len(ids) != len(dbTmplUniques) {
		blog.Errorf("field template uniques are invalid, data: %v, err: %v, rid: %v", uniques, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "uniques"))
		return
	}

	dbUniqueMap := make(map[int64]metadata.FieldTemplateUnique)
	for _, dbUnique := range dbTmplUniques {
		dbUniqueMap[dbUnique.ID] = dbUnique
	}

	for idx := range uniques {
		dbUnique := dbUniqueMap[uniques[idx].ID]
		uniques[idx].OwnerID = dbUnique.OwnerID
		uniques[idx].Creator = dbUnique.Creator
		uniques[idx].CreateTime = dbUnique.CreateTime
	}

	if err := s.updateFieldTemplateUniques(ctx.Kit, templateID, uniques); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *service) updateFieldTemplateUniques(kit *rest.Kit, templateID int64,
	uniques []metadata.FieldTemplateUnique) error {

	if err := s.checkUniques(kit, uniques, templateID); err != nil {
		blog.Errorf("check field template unique failed, unique: %v, err: %v, rid: %s", uniques, err, kit.Rid)
		return err
	}

	now := time.Now()
	for _, unique := range uniques {
		unique.Modifier = kit.User
		unique.LastTime = &metadata.Time{Time: now}

		cond := mapstr.MapStr{
			common.BKFieldID:    unique.ID,
			common.BKTemplateID: templateID,
		}

		err := mongodb.Client().Table(common.BKTableNameObjectUniqueTemplate).Update(kit.Ctx, cond, unique)
		if err != nil {
			blog.Errorf("update field template unique failed, data: %v, err: %v, rid: %s", unique, err, kit.Rid)
			return err
		}
	}

	return nil
}
