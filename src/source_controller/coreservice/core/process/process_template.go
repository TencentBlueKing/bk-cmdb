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

package process

import (
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateProcessTemplate(ctx core.ContextParams, template metadata.ProcessTemplate) (*metadata.ProcessTemplate, error) {
	// base attribute validate
	if field, err := template.Validate(); err != nil {
		blog.Errorf("CreateProcessTemplate failed, validation failed, code: %d, field: %s, err: %+v, rid: %s", common.CCErrCommParamsInvalid, field, err, ctx.ReqID)
		err := ctx.Error.Errorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, template.Metadata); err != nil {
		blog.Errorf("CreateProcessTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// keep metadata clean
	template.Metadata = metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))

	// validate service template id field
	serviceTemplate, err := p.GetServiceTemplate(ctx, template.ServiceTemplateID)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, template id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "service_template_id")
	}

	// make sure biz id identical with service template
	serviceTemplateBizID, err := metadata.BizIDFromMetadata(serviceTemplate.Metadata)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, parse biz id from service template failed, code: %d, err: %+v, rid: %s", common.CCErrCommInternalServerError, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParseBizIDFromMetadataInDBFailed)
	}
	if bizID != serviceTemplateBizID {
		blog.Errorf("CreateProcessTemplate failed, validation failed, input bizID:%d not equal service template bizID:%d, rid: %s", bizID, serviceTemplateBizID, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommParamsInvalid, "metadata.label.bk_biz_id")
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameProcessTemplate)
	if nil != err {
		blog.Errorf("CreateProcessTemplate failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommGenerateRecordIDFailed)
	}
	template.ID = int64(id)

	template.Creator = ctx.User
	template.Modifier = ctx.User
	template.CreateTime = time.Now()
	template.LastTime = time.Now()
	template.SupplierAccount = ctx.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Insert(ctx.Context, &template); nil != err {
		blog.Errorf("CreateProcessTemplate failed, mongodb failed, table: %s, template: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, template, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBInsertFailed)
	}
	return &template, nil
}

func (p *processOperation) GetProcessTemplate(ctx core.ContextParams, templateID int64) (*metadata.ProcessTemplate, error) {
	template := metadata.ProcessTemplate{}

	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(filter).One(ctx.Context, &template); nil != err {
		blog.Errorf("GetProcessTemplate failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, filter, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	return &template, nil
}

func (p *processOperation) UpdateProcessTemplate(ctx core.ContextParams, templateID int64, input metadata.ProcessTemplate) (*metadata.ProcessTemplate, error) {
	template, err := p.GetProcessTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	if field, err := input.Validate(); err != nil {
		blog.Errorf("UpdateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.Errorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// update fields to local object
	if input.Property != nil {
		template.Property.Update(*input.Property)
	}

	// do update
	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Update(ctx, filter, template); nil != err {
		blog.Errorf("UpdateServiceTemplate failed, mongodb failed, table: %s, filter: %+v, template: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, template, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBUpdateFailed)
	}
	return template, nil
}

func (p *processOperation) ListProcessTemplates(ctx core.ContextParams, bizID int64, serviceTemplateID int64, processTemplateIDs *[]int64, limit metadata.BasePage) (*metadata.MultipleProcessTemplate, error) {
	md := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
	filter := map[string]interface{}{}
	filter["metadata"] = md.ToMapStr()

	if serviceTemplateID != 0 {
		filter["service_template_id"] = serviceTemplateID
	}

	if processTemplateIDs != nil {
		filter["process_template_id"] = map[string][]int64{"$in": *processTemplateIDs}
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListProcessTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, filter, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	templates := make([]metadata.ProcessTemplate, 0)
	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(filter).Start(
		uint64(limit.Start)).Limit(uint64(limit.Limit)).All(ctx.Context, &templates); nil != err {
		blog.Errorf("ListProcessTemplates failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessTemplate, err, ctx.ReqID)
		return nil, ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleProcessTemplate{
		Count: total,
		Info:  templates,
	}
	return result, nil
}

func (p *processOperation) DeleteProcessTemplate(ctx core.ContextParams, processTemplateID int64) error {
	template, err := p.GetProcessTemplate(ctx, processTemplateID)
	if err != nil {
		blog.Errorf("DeleteProcessTemplate failed, GetServiceTemplate failed, templateID: %d, err: %+v, rid: %s", processTemplateID, err, ctx.ReqID)
		return err
	}

	// service template that referenced by process template shouldn't be removed
	usageFilter := map[string]int64{"process_template_id": template.ID}
	usageCount, err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(usageFilter).Count(ctx.Context)
	if nil != err {
		blog.Errorf("DeleteProcessTemplate failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, usageFilter, err, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteProcessTemplate failed, forbidden delete process template be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	deleteFilter := map[string]int64{common.BKFieldID: template.ID}
	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteProcessTemplate failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, deleteFilter, err, ctx.ReqID)
		return ctx.Error.Errorf(common.CCErrCommDBSelectFailed)
	}
	return nil
}
