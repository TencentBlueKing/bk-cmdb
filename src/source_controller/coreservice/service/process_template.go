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
)

func (s *coreService) CreateProcessTemplate(ctx *rest.Contexts) {
	template := metadata.ProcessTemplate{}
	if err := ctx.DecodeInto(&template); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateProcessTemplate(ctx.Kit, template)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) GetProcessTemplate(ctx *rest.Contexts) {
	processTemplateIDStr := ctx.Request.PathParameter(common.BKProcessTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("GetProcessTemplate failed, path parameter `%s` empty, rid: %s", processTemplateIDStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField))
		return
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcessTemplateIDField, processTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField))
		return
	}

	result, err := s.core.ProcessOperation().GetProcessTemplate(ctx.Kit, processTemplateID)
	if err != nil {
		blog.Errorf("GetProcessTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListProcessTemplates(ctx *rest.Contexts) {
	// filter parameter
	fp := metadata.ListProcessTemplatesOption{}
	if err := ctx.DecodeInto(&fp); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if fp.BusinessID == 0 {
		blog.Errorf("ListServiceTemplates failed, business id can't be empty, bk_biz_id: %d, rid: %s", fp.BusinessID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	result, err := s.core.ProcessOperation().ListProcessTemplates(ctx.Kit, fp)
	if err != nil {
		blog.Errorf("ListProcessTemplates failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateProcessTemplate(ctx *rest.Contexts) {
	processTemplateIDStr := ctx.Request.PathParameter(common.BKProcessTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("UpdateProcessTemplate failed, path parameter `%s` empty, rid: %s", common.BKProcessTemplateIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField))
		return
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcessTemplateIDField, processTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField))
		return
	}
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}
	result, err := s.core.ProcessOperation().UpdateProcessTemplate(ctx.Kit, processTemplateID, data)
	if err != nil {
		blog.Errorf("UpdateProcessTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) DeleteProcessTemplate(ctx *rest.Contexts) {
	processTemplateIDStr := ctx.Request.PathParameter(common.BKProcessTemplateIDField)
	if len(processTemplateIDStr) == 0 {
		blog.Errorf("DeleteProcessTemplate failed, path parameter `%s` empty, rid: %s", common.BKProcessTemplateIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField))
		return
	}

	processTemplateID, err := strconv.ParseInt(processTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteProcessTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKProcessTemplateIDField, processTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKProcessTemplateIDField))
		return
	}

	if err := s.core.ProcessOperation().DeleteProcessTemplate(ctx.Kit, processTemplateID); err != nil {
		blog.Errorf("DeleteProcessTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *coreService) BatchDeleteProcessTemplate(ctx *rest.Contexts) {
	input := struct {
		ProcessTemplateIDs []int64 `json:"process_template_ids" field:"process_template_ids"`
	}{}

	if err := ctx.DecodeInto(&input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// TODO: replace with batch delete interface
	for _, id := range input.ProcessTemplateIDs {
		if err := s.core.ProcessOperation().DeleteProcessTemplate(ctx.Kit, id); err != nil {
			blog.Errorf("BatchDeleteProcessTemplate failed, templateID: %d, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}
	ctx.RespEntity(nil)
}
