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

func (s *coreService) CreateServiceTemplate(ctx *rest.Contexts) {
	template := metadata.ServiceTemplate{}
	if err := ctx.DecodeInto(&template); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().CreateServiceTemplate(ctx.Kit, template)
	if err != nil {
		blog.Errorf("CreateServiceCategory failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) GetServiceTemplate(ctx *rest.Contexts) {
	serviceTemplateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("GetServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	result, err := s.core.ProcessOperation().GetServiceTemplate(ctx.Kit, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetServiceTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) GetServiceTemplateWithStatistics(ctx *rest.Contexts) {
	serviceTemplateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("GetServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	template, err := s.core.ProcessOperation().GetServiceTemplate(ctx.Kit, serviceTemplateID)
	if err != nil {
		blog.Errorf("GetServiceTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// related service instance count
	serviceInstanceFilter := map[string]interface{}{
		common.BKServiceTemplateIDField: template.ID,
	}
	serviceInstanceCount, err := mongodb.Client().Table(common.BKTableNameServiceInstance).Find(serviceInstanceFilter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, filter: %+v, err: %+v, rid: %s", serviceInstanceFilter, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}

	// related service template count
	processRelationFilter := map[string]interface{}{
		common.BKServiceTemplateIDField: template.ID,
	}
	processRelationCount, err := mongodb.Client().Table(common.BKTableNameProcessInstanceRelation).Find(processRelationFilter).Count(ctx.Kit.Ctx)
	if err != nil {
		blog.Errorf("GetServiceCategory failed, filter: %+v, err: %+v, rid: %s", serviceInstanceFilter, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommDBSelectFailed))
		return
	}
	result := metadata.ServiceTemplateWithStatistics{
		Template:             *template,
		ServiceInstanceCount: int64(serviceInstanceCount),
		ProcessInstanceCount: int64(processRelationCount),
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListServiceTemplateDetail(ctx *rest.Contexts) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	if len(bizIDStr) == 0 {
		blog.Errorf("ListServiceTemplateDetail failed, path parameter `%s` empty, rid: %s", common.BKAppIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListServiceTemplateDetail failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKAppIDField, bizIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
		return
	}

	input := struct {
		ServiceTemplateIDs []int64 `json:"service_template_ids" mapstructure:"service_template_ids"`
	}{}
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	option := metadata.ListServiceTemplateOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: input.ServiceTemplateIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	serviceTemplateResult, ccErr := s.core.ProcessOperation().ListServiceTemplates(ctx.Kit, option)
	if ccErr != nil {
		blog.Errorf("ListServiceTemplateDetail failed, ListServiceTemplate failed, err: %+v, rid: %s", ccErr, ctx.Kit.Rid)
		ctx.RespAutoError(ccErr)
		return
	}
	srvTplIDs := make([]int64, 0)
	for _, item := range serviceTemplateResult.Info {
		srvTplIDs = append(srvTplIDs, item.ID)
	}

	listProcessTemplateOption := metadata.ListProcessTemplatesOption{
		BusinessID:         bizID,
		ServiceTemplateIDs: srvTplIDs,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	listProcResult, ccErr := s.core.ProcessOperation().ListProcessTemplates(ctx.Kit, listProcessTemplateOption)
	if ccErr != nil {
		blog.Errorf("ListServiceTemplateDetail failed, ListProcessTemplates failed, err: %+v, rid: %s", ccErr, ctx.Kit.Rid)
		ctx.RespAutoError(ccErr)
		return
	}
	serviceProcessTemplateMap := make(map[int64][]metadata.ProcessTemplate)
	for _, item := range listProcResult.Info {
		if _, exist := serviceProcessTemplateMap[item.ServiceTemplateID]; exist == false {
			serviceProcessTemplateMap[item.ServiceTemplateID] = make([]metadata.ProcessTemplate, 0)
		}
		serviceProcessTemplateMap[item.ServiceTemplateID] = append(serviceProcessTemplateMap[item.ServiceTemplateID], item)
	}

	templateDetails := make([]metadata.ServiceTemplateDetail, 0)
	for _, item := range serviceTemplateResult.Info {
		templateDetail := metadata.ServiceTemplateDetail{
			ServiceTemplate:  item,
			ProcessTemplates: make([]metadata.ProcessTemplate, 0),
		}
		processTemplates, exist := serviceProcessTemplateMap[item.ID]
		if exist == true {
			templateDetail.ProcessTemplates = processTemplates
		}
		templateDetails = append(templateDetails, templateDetail)
	}
	result := metadata.MultipleServiceTemplateDetail{
		Count: serviceTemplateResult.Count,
		Info:  templateDetails,
	}
	ctx.RespEntity(result)
}

func (s *coreService) ListServiceTemplates(ctx *rest.Contexts) {
	// filter parameter
	fp := metadata.ListServiceTemplateOption{}

	if err := ctx.DecodeInto(&fp); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().ListServiceTemplates(ctx.Kit, fp)
	if err != nil {
		blog.Errorf("ListServiceTemplates failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
}

func (s *coreService) UpdateServiceTemplate(ctx *rest.Contexts) {
	serviceTemplateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("UpdateServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	template := metadata.ServiceTemplate{}
	if err := ctx.DecodeInto(&template); err != nil {
		ctx.RespAutoError(err)
		return
	}

	result, err := s.core.ProcessOperation().UpdateServiceTemplate(ctx.Kit, serviceTemplateID, template)
	if err != nil {
		blog.Errorf("UpdateServiceTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

func (s *coreService) DeleteServiceTemplate(ctx *rest.Contexts) {
	serviceTemplateIDStr := ctx.Request.PathParameter(common.BKServiceTemplateIDField)
	if len(serviceTemplateIDStr) == 0 {
		blog.Errorf("DeleteServiceTemplate failed, path parameter `%s` empty, rid: %s", common.BKServiceTemplateIDField, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	serviceTemplateID, err := strconv.ParseInt(serviceTemplateIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteServiceTemplate failed, convert path parameter %s to int failed, value: %s, err: %v, rid: %s", common.BKServiceTemplateIDField, serviceTemplateIDStr, err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceTemplateIDField))
		return
	}

	if err := s.core.ProcessOperation().DeleteServiceTemplate(ctx.Kit, serviceTemplateID); err != nil {
		blog.Errorf("DeleteServiceTemplate failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}
