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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateProcessTemplate(ctx core.ContextParams, template metadata.ProcessTemplate) (*metadata.ProcessTemplate, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := template.Validate(); err != nil {
		blog.Errorf("CreateProcessTemplate failed, validation failed, code: %d, field: %s, err: %+v, rid: %s", common.CCErrCommParamsInvalid, field, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}
	*template.Property.ProcessName.AsDefaultValue = true
	*template.Property.FuncName.AsDefaultValue = true
	if template.Property != nil && template.Property.ProcessName.Value != nil {
		template.ProcessName = *template.Property.ProcessName.Value
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, template.BizID); err != nil {
		blog.Errorf("CreateProcessTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	template.BizID = bizID

	// validate service template id field
	serviceTemplate, err := p.GetServiceTemplate(ctx, template.ServiceTemplateID)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, template id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_template_id")
	}

	// make sure biz id identical with service template
	if bizID != serviceTemplate.BizID {
		blog.Errorf("CreateProcessTemplate failed, validation failed, input bizID:%d not equal service template bizID:%d, rid: %s", bizID, serviceTemplate.BizID, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	if err := p.processNameUniqueValidate(ctx, &template); err != nil {
		return nil, err
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameProcessTemplate)
	if nil != err {
		blog.Errorf("CreateProcessTemplate failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	template.ID = int64(id)

	template.Creator = ctx.User
	template.Modifier = ctx.User
	template.CreateTime = time.Now()
	template.LastTime = time.Now()
	template.SupplierAccount = ctx.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Insert(ctx.Context, &template); nil != err {
		blog.Errorf("CreateProcessTemplate failed, mongodb failed, table: %s, template: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, template, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return &template, nil
}

func (p *processOperation) processNameUniqueValidate(ctx core.ContextParams, template *metadata.ProcessTemplate) errors.CCErrorCoder {
	// process name unique
	processName := ""
	if template.Property.ProcessName.Value != nil {
		processName = *template.Property.ProcessName.Value
	}
	processNameFilter := map[string]interface{}{
		common.BKServiceTemplateIDField:  template.ServiceTemplateID,
		"property.bk_process_name.value": processName,
	}
	if template.ID != 0 {
		processNameFilter[common.BKFieldID] = map[string]interface{}{
			common.BKDBNE: template.ID,
		}
	}
	count, err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(processNameFilter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, check process_name unique failed, err: %+v, rid: %s", err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		return ctx.Error.CCErrorf(common.CCErrCoreServiceProcessNameDuplicated, processName)
	}

	// func name unique
	funcName := ""
	if template.Property.FuncName.Value != nil {
		funcName = *template.Property.FuncName.Value
	}
	startRegex := ""
	if template.Property.StartParamRegex.Value != nil {
		startRegex = *template.Property.StartParamRegex.Value
	}
	funcNameFilter := map[string]interface{}{
		common.BKServiceTemplateIDField:       template.ServiceTemplateID,
		"property.bk_func_name.value":         funcName,
		"property.bk_start_param_regex.value": startRegex,
	}
	if template.ID != 0 {
		funcNameFilter[common.BKFieldID] = map[string]interface{}{
			common.BKDBNE: template.ID,
		}
	}
	count, err = p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(funcNameFilter).Count(ctx.Context)
	if err != nil {
		blog.Errorf("CreateProcessTemplate failed, check func_name unique failed, err: %+v, rid: %s", err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		return ctx.Error.CCErrorf(common.CCErrCoreServiceFuncNameDuplicated, funcName, startRegex)
	}
	return nil
}

func (p *processOperation) GetProcessTemplate(ctx core.ContextParams, templateID int64) (*metadata.ProcessTemplate, errors.CCErrorCoder) {
	template := metadata.ProcessTemplate{}

	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(filter).One(ctx.Context, &template); nil != err {
		blog.Errorf("GetProcessTemplate failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, filter, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &template, nil
}

func (p *processOperation) UpdateProcessTemplate(ctx core.ContextParams, templateID int64, rawProperty map[string]interface{}) (*metadata.ProcessTemplate, errors.CCErrorCoder) {
	template, err := p.GetProcessTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	property := metadata.ProcessProperty{}
	if err := mapstr.DecodeFromMapStr(&property, rawProperty); err != nil {
		blog.ErrorJSON("UpdateProcessTemplate failed, unmarshal failed, property: %s, err: %s, rid: %s", property, err, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommJSONUnmarshalFailed)
		return nil, err
	}

	// update fields to local object
	template.Property.Update(property, rawProperty)

	if field, err := template.Validate(); err != nil {
		blog.Errorf("UpdateProcessTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	*template.Property.ProcessName.AsDefaultValue = true
	*template.Property.FuncName.AsDefaultValue = true
	template.Modifier = ctx.User
	template.LastTime = time.Now()

	if err := p.processNameUniqueValidate(ctx, template); err != nil {
		return nil, err
	}
	if template.Property != nil {
		if template.Property.ProcessName.Value != nil {
			template.ProcessName = *template.Property.ProcessName.Value
		}
	}

	// do update
	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Update(ctx, filter, &template); nil != err {
		blog.Errorf("UpdateProcessTemplate failed, mongodb failed, table: %s, filter: %+v, template: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, filter, template, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return template, nil
}

func (p *processOperation) ListProcessTemplates(ctx core.ContextParams, option metadata.ListProcessTemplatesOption) (*metadata.MultipleProcessTemplate, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
	}

	if len(option.ServiceTemplateIDs) != 0 {
		filter[common.BKServiceTemplateIDField] = map[string]interface{}{
			common.BKDBIN: option.ServiceTemplateIDs,
		}
	}

	if option.ProcessTemplateIDs != nil {
		filter[common.BKProcessTemplateIDField] = map[string][]int64{
			common.BKDBIN: option.ProcessTemplateIDs,
		}
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListProcessTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, filter, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	templates := make([]metadata.ProcessTemplate, 0)

	// ex: "-id,name"
	sort := "-id"
	if len(option.Page.Sort) > 0 {
		sort = option.Page.Sort
	}

	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Find(filter).Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(sort).All(ctx.Context, &templates); nil != err {
		blog.Errorf("ListProcessTemplates failed, mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameProcessTemplate, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleProcessTemplate{
		Count: total,
		Info:  templates,
	}
	return result, nil
}

func (p *processOperation) DeleteProcessTemplate(ctx core.ContextParams, processTemplateID int64) errors.CCErrorCoder {
	template, err := p.GetProcessTemplate(ctx, processTemplateID)
	if err != nil {
		blog.Errorf("DeleteProcessTemplate failed, GetProcessTemplate failed, templateID: %d, err: %+v, rid: %s", processTemplateID, err, ctx.ReqID)
		return err
	}

	updateFilter := map[string]int64{
		common.BKProcessTemplateIDField: template.ID,
	}
	updateDoc := map[string]interface{}{
		common.BKProcessTemplateIDField: common.ServiceTemplateIDNotSet,
	}
	e := p.dbProxy.Table(common.BKTableNameProcessInstanceRelation).Update(ctx.Context, updateFilter, updateDoc)
	if nil != e {
		blog.Errorf("DeleteProcessTemplate failed, clear process instance templateID field failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessInstanceRelation, updateFilter, e, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	deleteFilter := map[string]int64{common.BKFieldID: template.ID}
	if err := p.dbProxy.Table(common.BKTableNameProcessTemplate).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteProcessTemplate failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameProcessTemplate, deleteFilter, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	return nil
}
