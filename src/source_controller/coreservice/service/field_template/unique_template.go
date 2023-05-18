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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
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

	if err := s.checkUniques(ctx.Kit, uniques); err != nil {
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

func (s *service) checkUniques(kit *rest.Kit, uniques []metadata.FieldTemplateUnique) error {
	templateIDs := make([]int64, 0)
	templateIDMap := make(map[int64]struct{})
	for idx, unique := range uniques {
		if err := unique.Validate(); err.ErrCode != 0 {
			return err.ToCCError(kit.CCError)
		}

		if err := s.checkKeys(kit, unique, uniques, idx); err != nil {
			return err
		}

		if _, exist := templateIDMap[unique.TemplateID]; !exist {
			templateIDs = append(templateIDs, unique.TemplateID)
			templateIDMap[unique.TemplateID] = struct{}{}
		}
	}

	existedUniques, err := s.findUniqueByTemplateIDs(kit, templateIDs)
	if err != nil {
		blog.Errorf("find uniques by template ids failed, ids: %v, err: %v, rid: %s", templateIDs, err, kit.Rid)
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
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
		}

		keysMap[key] = struct{}{}
	}

	exist, err := s.isTemplateAttrsExist(kit, target.Keys)
	if err != nil {
		return err
	}
	if !exist {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
	}

	for idx, comparator := range comparators {
		if idx == skipIdx || target.TemplateID != comparator.TemplateID ||
			(target.TemplateID != 0 && target.ID == comparator.ID) {
			continue
		}

		cnt := 0
		for _, existedKey := range comparator.Keys {
			if _, exist := keysMap[existedKey]; exist {
				cnt++
			}
		}

		if len(keysMap) == cnt || (len(keysMap) > len(comparator.Keys) && len(comparator.Keys) == cnt) {
			return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKObjectUniqueKeys)
		}
	}

	return nil
}

func (s *service) isTemplateAttrsExist(kit *rest.Kit, ids []int64) (bool, error) {
	cond := map[string]interface{}{
		common.BKFieldID: map[string]interface{}{common.BKDBIN: ids},
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)

	cnt, err := mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count template attributes failed, cond: %v, err: %v, rid: %v", cond, err, kit.Rid)
		return false, err
	}

	return int(cnt) == len(ids), nil
}

func (s *service) findUniqueByTemplateIDs(kit *rest.Kit, ids []int64) ([]metadata.FieldTemplateUnique, error) {
	cond := map[string]interface{}{
		common.BKTemplateID: map[string]interface{}{common.BKDBIN: ids},
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
		Limit(uint64(opt.Page.Limit)).Fields(opt.Fields...).All(ctx.Kit.Ctx, &uniques)
	if err != nil {
		blog.Errorf("list field template uniques failed, err: %v, filter: %+v, rid: %v", err, filter, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	ctx.RespEntity(metadata.FieldTemplateUniqueInfo{Info: uniques})
}
