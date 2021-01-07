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

	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) CreateServiceTemplate(ctx *rest.Contexts) {
	option := new(metadata.CreateServiceTemplateOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	newTemplate := &metadata.ServiceTemplate{
		BizID:             option.BizID,
		Name:              option.Name,
		ServiceCategoryID: option.ServiceCategoryID,
		SupplierAccount:   ctx.Kit.SupplierAccount,
	}

	var tpl *metadata.ServiceTemplate
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		tpl, err = ps.CoreAPI.CoreService().Process().CreateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, newTemplate)
		if err != nil {
			blog.Errorf("create service template failed, err: %v", err)
			return err
		}

		// register service template resource creator action to iam
		if auth.EnableAuthorize() {
			iamInstance := metadata.IamInstanceWithCreator{
				Type:    string(iam.BizProcessServiceTemplate),
				ID:      strconv.FormatInt(tpl.ID, 10),
				Name:    tpl.Name,
				Creator: ctx.Kit.User,
			}
			_, err = ps.AuthManager.Authorizer.RegisterResourceCreatorAction(ctx.Kit.Ctx, ctx.Kit.Header, iamInstance)
			if err != nil {
				blog.Errorf("register created service template to iam failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				return err
			}
		}

		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(tpl)
}

func (ps *ProcServer) GetServiceTemplate(ctx *rest.Contexts) {
	templateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	templateID, err := util.GetInt64ByInterface(templateIDStr)
	if err != nil {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "create service template failed, err: %v", common.BKServiceTemplateIDField, err)
		return
	}
	template, err := ps.CoreAPI.CoreService().Process().GetServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get service template failed, err: %v", err)
		return
	}

	ctx.RespEntity(template)
}

// GetServiceTemplateDetail return more info than GetServiceTemplate
func (ps *ProcServer) GetServiceTemplateDetail(ctx *rest.Contexts) {
	templateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	templateID, err := util.GetInt64ByInterface(templateIDStr)
	if err != nil {
		ctx.RespErrorCodeF(common.CCErrCommParamsInvalid, "create service template failed, err: %v", common.BKServiceTemplateIDField, err)
		return
	}
	templateDetail, err := ps.CoreAPI.CoreService().Process().GetServiceTemplateWithStatistics(ctx.Kit.Ctx, ctx.Kit.Header, templateID)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "get service template failed, err: %v", err)
		return
	}

	ctx.RespEntity(templateDetail)
}

func (ps *ProcServer) UpdateServiceTemplate(ctx *rest.Contexts) {
	option := new(metadata.UpdateServiceTemplateOption)
	if err := ctx.DecodeInto(option); err != nil {
		ctx.RespAutoError(err)
		return
	}

	updateParam := &metadata.ServiceTemplate{
		ID:                option.ID,
		Name:              option.Name,
		ServiceCategoryID: option.ServiceCategoryID,
	}

	var tpl *metadata.ServiceTemplate
	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		tpl, err = ps.CoreAPI.CoreService().Process().UpdateServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, option.ID, updateParam)
		if err != nil {
			blog.Errorf("update service template failed, err: %v", err)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(tpl)
}

func (ps *ProcServer) ListServiceTemplates(ctx *rest.Contexts) {
	input := new(metadata.ListServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if input.Page.Limit > common.BKMaxPageSize {
		ctx.RespErrorCodeOnly(common.CCErrCommPageLimitIsExceeded, "list service template, but page limit:%d is over limited.", input.Page.Limit)
		return
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:        input.BizID,
		Page:              input.Page,
		ServiceCategoryID: &input.ServiceCategoryID,
	}
	temp, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "list service template failed, input: %+v", input)
		return
	}

	ctx.RespEntity(temp)
}

func (ps *ProcServer) ListServiceTemplatesWithDetails(ctx *rest.Contexts) {
	input := new(metadata.ListServiceTemplateInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	if input.Page.Limit > common.BKMaxPageSize {
		ctx.RespErrorCodeOnly(common.CCErrCommPageLimitIsExceeded, "list service template, but page limit:%d is over limited.", input.Page.Limit)
		return
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:        input.BizID,
		Page:              input.Page,
		ServiceCategoryID: &input.ServiceCategoryID,
		Search:            input.Search,
	}
	listResult, err := ps.CoreAPI.CoreService().Process().ListServiceTemplates(ctx.Kit.Ctx, ctx.Kit.Header, &option)
	if err != nil {
		ctx.RespWithError(err, common.CCErrCommHTTPDoRequestFailed, "list service template failed, input: %+v", input)
		return
	}

	// generate count conditions
	filters := make([]map[string]interface{}, len(listResult.Info))
	for idx, serviceTemplate := range listResult.Info {
		filters[idx] = map[string]interface{}{
			common.BKServiceTemplateIDField: serviceTemplate.ID,
		}
	}

	// process templates reference count
	processTemplateCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameProcessTemplate, filters)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetProcessTemplatesFailed, "count process template by filters: %+v failed.", filters)
		return
	}

	// module reference count
	moduleCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameBaseModule, filters)
	if err != nil {
		ctx.RespWithError(err, common.CCErrTopoModuleSelectFailed, "count process template by filters: %+v failed.", filters)
		return
	}

	// service instance reference count
	serviceInstanceCounts, err := ps.CoreAPI.CoreService().Count().GetCountByFilter(ctx.Kit.Ctx, ctx.Kit.Header, common.BKTableNameServiceInstance, filters)
	if err != nil {
		ctx.RespWithError(err, common.CCErrProcGetServiceInstancesFailed, "count process template by filters: %+v failed.", filters)
		return
	}

	details := make([]metadata.ListServiceTemplateWithDetailResult, 0)
	for idx, serviceTemplate := range listResult.Info {
		details = append(details, metadata.ListServiceTemplateWithDetailResult{
			ServiceTemplate:      serviceTemplate,
			ProcessTemplateCount: processTemplateCounts[idx],
			ServiceInstanceCount: serviceInstanceCounts[idx],
			ModuleCount:          moduleCounts[idx],
		})
	}

	ctx.RespEntityWithCount(int64(listResult.Count), details)
}

// a service template can be delete only when it is not be used any more,
// which means that no process instance belongs to it.
func (ps *ProcServer) DeleteServiceTemplate(ctx *rest.Contexts) {
	input := new(metadata.DeleteServiceTemplatesInput)
	if err := ctx.DecodeInto(input); err != nil {
		ctx.RespAutoError(err)
		return
	}

	txnErr := ps.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ps.EnableTxn, ctx.Kit.Header, func() error {
		err := ps.CoreAPI.CoreService().Process().DeleteServiceTemplate(ctx.Kit.Ctx, ctx.Kit.Header, input.ServiceTemplateID)
		if err != nil {
			blog.Errorf("delete service template: %d failed", input.ServiceTemplateID)
			return ctx.Kit.CCError.CCError(common.CCErrProcDeleteServiceTemplateFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}
