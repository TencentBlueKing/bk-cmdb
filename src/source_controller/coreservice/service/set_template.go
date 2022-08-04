/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// CreateSetTemplate TODO
func (s *coreService) CreateSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.CreateSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.SetTemplateOperation().CreateSetTemplate(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("CreateSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// UpdateSetTemplate TODO
func (s *coreService) UpdateSetTemplate(ctx *rest.Contexts) {
	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	option := metadata.UpdateSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.SetTemplateOperation().UpdateSetTemplate(ctx.Kit, setTemplateID, option)
	if err != nil {
		blog.Errorf("UpdateSetTemplate failed, setTemplateID: %d, option: %+v, err: %+v, rid: %s", setTemplateID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

// DeleteSetTemplate TODO
func (s *coreService) DeleteSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.DeleteSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.SetTemplateOperation().DeleteSetTemplate(ctx.Kit, bizID, option); err != nil {
		blog.Errorf("UpdateSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// GetSetTemplate TODO
func (s *coreService) GetSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	setTemplate, err := s.core.SetTemplateOperation().GetSetTemplate(ctx.Kit, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("GetSetTemplate failed, bizID: %d, setTemplateID: %d, err: %+v, rid: %s", bizID, setTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(setTemplate)
}

// ListSetTemplate TODO
func (s *coreService) ListSetTemplate(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListSetTemplateOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	setTemplateResult, err := s.core.SetTemplateOperation().ListSetTemplate(ctx.Kit, bizID, option)
	if err != nil {
		blog.Errorf("ListSetTemplate failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(setTemplateResult)
}

// CountSetTplInstances TODO
func (s *coreService) CountSetTplInstances(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.CountSetTplInstOption{}
	if err := ctx.DecodeInto(&option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	filter := map[string]interface{}{
		common.BKSetTemplateIDField: map[string]interface{}{
			common.BKDBIN: option.SetTemplateIDs,
		},
		common.BKAppIDField: bizID,
	}
	pipeline := []map[string]interface{}{
		{common.BKDBMatch: filter},
		{common.BKDBGroup: map[string]interface{}{
			"_id":                 "$" + common.BKSetTemplateIDField,
			"set_instances_count": map[string]interface{}{common.BKDBSum: 1}},
		},
	}
	result := make([]metadata.CountSetTplInstItem, 0)
	if err := mongodb.Client().Table(common.BKTableNameBaseSet).AggregateAll(ctx.Kit.Ctx, pipeline, &result); err != nil {
		if mongodb.Client().IsNotFoundError(err) == true {
			result = make([]metadata.CountSetTplInstItem, 0)
		} else {
			blog.Errorf("CountSetTplInstances failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommDBSelectFailed))
			return
		}
	}

	ctx.RespEntity(result)
}

// ListSetServiceTemplateRelations TODO
func (s *coreService) ListSetServiceTemplateRelations(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	relations, err := s.core.SetTemplateOperation().ListSetServiceTemplateRelations(ctx.Kit, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetServiceTemplateRelations failed, bizID: %d, setTemplateID: %+v, err: %+v, rid: %s", bizID, setTemplateID, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(relations)
}

// ListSetTplRelatedSvcTpl TODO
func (s *coreService) ListSetTplRelatedSvcTpl(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	setTemplateIDStr := ctx.Request.PathParameter(common.BKSetTemplateIDField)
	setTemplateID, err := strconv.ParseInt(setTemplateIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetTemplateIDField))
		return
	}

	serviceTemplates, err := s.core.SetTemplateOperation().ListSetTplRelatedSvcTpl(ctx.Kit, bizID, setTemplateID)
	if err != nil {
		blog.Errorf("ListSetTplRelatedSvcTpl failed, bizID: %d, setTemplateID: %d, err: %s, rid: %s", bizID, setTemplateID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(serviceTemplates)
}

// CreateSetTemplateAttribute create set template attributes
func (s *coreService) CreateSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.CreateSetTempAttrsOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	ids, err := s.core.SetTemplateOperation().CreateSetTempAttr(ctx.Kit, option)
	if err != nil {
		blog.Errorf("create set template attributes(%+v) failed, err: %v, rid: %s", option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	result := metadata.RspIDs{
		IDs: make([]int64, len(ids)),
	}

	for idx, id := range ids {
		result.IDs[idx] = int64(id)
	}
	ctx.RespEntity(result)
}

// UpdateSetTemplateAttribute update set template attribute
func (s *coreService) UpdateSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.UpdateSetTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.core.SetTemplateOperation().UpdateSetTempAttr(ctx.Kit, option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// DeleteSetTemplateAttribute delete set template attribute
func (s *coreService) DeleteSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.DeleteSetTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	if err := s.core.SetTemplateOperation().DeleteSetTemplateAttribute(ctx.Kit, option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// ListSetTemplateAttribute list set template attribute
func (s *coreService) ListSetTemplateAttribute(ctx *rest.Contexts) {
	option := new(metadata.ListSetTempAttrOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if err := option.Validate(); err.ErrCode != 0 {
		ctx.RespAutoError(err.ToCCError(ctx.Kit.CCError))
		return
	}

	data, err := s.core.SetTemplateOperation().ListSetTemplateAttribute(ctx.Kit, option)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(data)
}
