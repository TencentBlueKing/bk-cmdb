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

func (s *coreService) UpdateSetTemplateSyncStatus(ctx *rest.Contexts) {
	setIDStr := ctx.Request.PathParameter(common.BKSetIDField)
	setID, err := strconv.ParseInt(setIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKSetIDField))
		return
	}

	option := metadata.SetTemplateSyncStatus{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}

	if err := s.core.SetTemplateOperation().UpdateSetTemplateSyncStatus(ctx.Kit, setID, option); err != nil {
		blog.Errorf("UpdateSetTemplateSyncStatus failed, setID: %d, option: %+v, err: %+v, rid: %s", setID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) ListSetTemplateSyncStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListSetTemplateSyncStatusOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}
	option.BizID = bizID

	result, err := s.core.SetTemplateOperation().ListSetTemplateSyncStatus(ctx.Kit, option)
	if err != nil {
		blog.Errorf("ListSetTemplateSyncStatus failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListSetTemplateSyncHistory(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.ListSetTemplateSyncStatusOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}
	option.BizID = bizID

	result, err := s.core.SetTemplateOperation().ListSetTemplateSyncHistory(ctx.Kit, option)
	if err != nil {
		blog.Errorf("ListSetTemplateSyncHistory failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) DeleteSetTemplateSyncStatus(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	option := metadata.DeleteSetTemplateSyncStatusOption{}
	if err := ctx.DecodeInto(&option); nil != err {
		ctx.RespAutoError(err)
		return
	}
	option.BizID = bizID

	ccErr := s.core.SetTemplateOperation().DeleteSetTemplateSyncStatus(ctx.Kit, option)
	if ccErr != nil {
		blog.Errorf("DeleteSetTemplateSyncStatus failed, bizID: %d, option: %+v, err: %+v, rid: %s", bizID, option, ccErr, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

func (s *coreService) ModifySetTemplateSyncStatus(ctx *rest.Contexts) {
	setIDStr := ctx.Request.PathParameter(common.BKSetIDField)
	setID, err := strconv.ParseInt(setIDStr, 10, 64)
	if err != nil {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsNeedInt, common.BKSetIDField))
		return
	}
	strStatus := ctx.Request.PathParameter(common.BKStatusField)
	if strStatus == "" {
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsLostField, common.BKStatusField))
		return
	}
	var status metadata.SyncStatus
	switch strStatus {
	case "waiting":
		status = metadata.SyncStatusWaiting
	case "syncing":
		status = metadata.SyncStatusSyncing
	case "finished":
		status = metadata.SyncStatusFinished
	case "failure":
		status = metadata.SyncStatusFailure
	default:
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKStatusField))
		return

	}

	ccErr := s.core.SetTemplateOperation().ModifySetTemplateSyncStatus(ctx.Kit, setID, status)
	if ccErr != nil {
		blog.Errorf("ModifySetTemplateSyncStatus failed, setID: %+v, status: %s, err: %+v, rid: %s", setID, status, ccErr, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}
