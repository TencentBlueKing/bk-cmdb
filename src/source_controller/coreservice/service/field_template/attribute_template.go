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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
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

	filter = util.SetQueryOwner(filter, cts.Kit.SupplierAccount)

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
		Limit(uint64(opt.Page.Limit)).Fields(opt.Fields...).All(cts.Kit.Ctx, &attrTemplates)
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

		if attrs[idx].TemplateID != templateID {
			blog.Errorf("attribute template id is invalid, data: %v, template id: %d, rid: %s", attrs[idx], templateID,
				ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsInvalid, "attributes"))
			return
		}

		if err := attrs[idx].Validate(); err.ErrCode != 0 {
			blog.Errorf("field template attribute is invalid, data: %v, err: %v, rid: %s", attrs[idx], err, ctx.Kit.Rid)
			ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
			return
		}

		err := valid.ValidPropertyOption(attrs[idx].PropertyType, attrs[idx].Option, attrs[idx].IsMultiple,
			attrs[idx].Default, ctx.Kit.Rid, ctx.Kit.CCError)
		if err != nil {
			blog.Errorf("validate field template attribute option failed, data: %v, err: %v, rid: %s", attrs[idx], err,
				ctx.Kit.Rid)
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
