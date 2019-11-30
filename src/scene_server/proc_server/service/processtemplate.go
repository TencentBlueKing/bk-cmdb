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
	"strconv"

	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// create a process template for a service template.
func (ps *ProcServer) CreateProcessTemplateBatch(ctx *rest.Contexts) {
	input := new(metadata.CreateProcessTemplateBatchInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "create process template, but get business id failed, err: %v", err)
			return
		}
	}

	// authorize
	if err := ps.AuthManager.AuthorizeByServiceTemplateID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, input.ServiceTemplateID); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCheckAuthorizeFailed, "authorize by service template id failed, id: %d, err: %+v", input.ServiceTemplateID, err)
		return
	}

	ids := make([]int64, 0)
	for _, process := range input.Processes {
		t := &metadata.ProcessTemplate{
			BizID:             bizID,
			ServiceTemplateID: input.ServiceTemplateID,
			Property:          process.Spec,
		}

		temp, err := ps.CoreAPI.CoreService().Process().CreateProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, t)
		if err != nil {
			ctx.RespWithError(err, common.CCErrProcCreateProcessTemplateFailed, "create process template failed, template: +%v", *t)
			return
		}

		ids = append(ids, temp.ID)
	}

	ctx.RespEntity(ids)
}

func (ps *ProcServer) DeleteProcessTemplateBatch(ctx *rest.Contexts) {
	input := new(metadata.DeleteProcessTemplateBatchInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "delete process template: %v, but get business id failed, err: %v", input.ProcessTemplates, err)
			return
		}
	}
	input.BizID = bizID

	// authorize by service template
	listOption := &metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ProcessTemplateIDs: input.ProcessTemplates,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	serviceTemplateIDs := make([]int64, 0)
	for _, processTemplate := range processTemplates.Info {
		serviceTemplateIDs = append(serviceTemplateIDs, processTemplate.ServiceTemplateID)
	}
	if err := ps.AuthManager.AuthorizeByServiceTemplateID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, serviceTemplateIDs...); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCheckAuthorizeFailed, "authorize by service template id failed, id: %+v, err: %+v", serviceTemplateIDs, err)
		return
	}

	err = ps.CoreAPI.CoreService().Process().DeleteProcessTemplateBatch(ctx.Kit.Ctx, ctx.Kit.Header, input.ProcessTemplates)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed, "delete process template: %v failed",
			input.ProcessTemplates)
		return
	}
	ctx.RespEntity(nil)
}

func (ps *ProcServer) UpdateProcessTemplate(ctx *rest.Contexts) {
	input := new(metadata.UpdateProcessTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update process template, but get business id failed, err: %v, input: %+v",
				err, input)
			return
		}
	}
	input.BizID = bizID

	if input.Property == nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update process template, but property empty, input: %+v", input)
		return
	}
	if input.ProcessTemplateID <= 0 {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "update process template, but get nil process template, input: %+v", input)
		return
	}

	listOption := &metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ProcessTemplateIDs: []int64{input.ProcessTemplateID},
	}
	processTemplates, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, listOption)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	serviceTemplateIDs := make([]int64, 0)
	for _, processTemplate := range processTemplates.Info {
		serviceTemplateIDs = append(serviceTemplateIDs, processTemplate.ServiceTemplateID)
	}
	if err := ps.AuthManager.AuthorizeByServiceTemplateID(ctx.Kit.Ctx, ctx.Kit.Header, meta.Update, serviceTemplateIDs...); err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommCheckAuthorizeFailed, "authorize by service template id failed, id: %+v, err: %+v", serviceTemplateIDs, err)
		return
	}

	template, err := ps.CoreAPI.CoreService().Process().UpdateProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, input.ProcessTemplateID, input.Property)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcUpdateProcessTemplateFailed, "update process template: %v failed.", input)
		return
	}
	ctx.RespEntity(template)
}

func (ps *ProcServer) GetProcessTemplate(ctx *rest.Contexts) {
	input := &struct {
		Metadata metadata.Metadata `json:"metadata"`
		BizID    int64             `json:"bk_biz_id"`
	}{}
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 {
		var err error
		bizID, err = metadata.BizIDFromMetadata(input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get process template, but get business id failed, err: %v, input: %+v",
				err, input)
			return
		}
	}

	templateID, err := strconv.ParseInt(ctx.Request.PathParameter("processTemplateID"), 10, 64)
	if err != nil {
		ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get process template, but get process template id failed, err: %v", err)
		return
	}

	template, err := ps.CoreAPI.CoreService().Process().GetProcessTemplate(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get process template: %v failed, err: %v.", input, err)
		return
	}
	ctx.RespEntity(template)
}

func (ps *ProcServer) ListProcessTemplate(ctx *rest.Contexts) {
	input := new(metadata.ListProcessTemplateWithServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	bizID := input.BizID
	if bizID == 0 && input.Metadata != nil {
		var err error
		bizID, err = metadata.BizIDFromMetadata(*input.Metadata)
		if err != nil {
			ctx.RespErrorCodeOnly(common.CCErrCommHTTPInputInvalid, "get process template, but get business id failed, err: %v, input: %+v", err, input)
			return
		}
	}

	option := &metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: []int64{input.ServiceTemplateID},
		Page:               input.Page,
	}
	if input.ProcessTemplatesIDs != nil {
		option.ProcessTemplateIDs = input.ProcessTemplatesIDs
	}
	template, err := ps.CoreAPI.CoreService().Process().ListProcessTemplates(ctx.Kit.Ctx, ctx.Kit.Header, option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplateFailed, "get process template: %v failed", input)
		return
	}
	ctx.RespEntity(template)
}
