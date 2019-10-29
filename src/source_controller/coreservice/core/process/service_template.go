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
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateServiceTemplate(ctx core.ContextParams, template metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := template.Validate(); err != nil {
		blog.Errorf("CreateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, template.BizID); err != nil {
		blog.Errorf("CreateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// keep metadata clean
	template.BizID = bizID

	// validate service category id field
	category, err := p.GetServiceCategory(ctx, template.ServiceCategoryID)
	if err != nil {
		blog.Errorf("CreateServiceTemplate failed, category id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, "service_category_id")
	}
	isLeafNode, ccErr := p.IsServiceCategoryLeafNode(ctx, category.ID)
	if ccErr != nil {
		blog.Errorf("UpdateServiceTemplate failed, check leaf node failed, err: %+v, rid: %s", ccErr, ctx.ReqID)
		return nil, ccErr
	}
	if isLeafNode == false {
		return nil, ctx.Error.CCError(common.CCErrCoreServiceOnlyNodeServiceCategoryAvailable)
	}

	// make sure biz id identical with category
	// categoryBizID 0 and 1 is default category
	if bizID != category.BizID && category.BizID != 0 {
		blog.Errorf("CreateServiceTemplate failed, validation failed, input bizID:%d not equal category bizID:%d, rid: %s", bizID, category.BizID, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	// check name field unique under business
	nameUniqueFilter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKFieldName:  template.Name,
	}
	count, err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(nameUniqueFilter).Count(ctx)
	if err != nil {
		blog.Errorf("CreateServiceTemplate failed, count same name instance failed, filter: %+v, err: %+v, rid: %s", nameUniqueFilter, err, ctx.ReqID)
		return nil, ctx.Error.CCError(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("CreateServiceTemplate failed, service instance name duplicated, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDuplicateItem, common.BKFieldName)
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameServiceTemplate)
	if nil != err {
		blog.Errorf("CreateServiceTemplate failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	template.ID = int64(id)

	template.Creator = ctx.User
	template.Modifier = ctx.User
	template.CreateTime = time.Now()
	template.LastTime = time.Now()
	template.SupplierAccount = ctx.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Insert(ctx.Context, &template); nil != err {
		blog.Errorf("CreateServiceTemplate failed, mongodb failed, table: %s, template: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, template, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBInsertFailed)
	}
	return &template, nil
}

func (p *processOperation) GetServiceTemplate(ctx core.ContextParams, templateID int64) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	template := metadata.ServiceTemplate{}

	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).One(ctx.Context, &template); nil != err {
		blog.Errorf("GetServiceTemplate failed, mongodb failed, table: %s, filter: %+v, template: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, template, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &template, nil
}

// UpdateServiceTemplate
// not support update name field yet, so don't need validate name unique before update
func (p *processOperation) UpdateServiceTemplate(ctx core.ContextParams, templateID int64, input metadata.ServiceTemplate) (*metadata.ServiceTemplate, errors.CCErrorCoder) {
	template, err := p.GetServiceTemplate(ctx, templateID)
	if err != nil {
		return nil, err
	}

	// update fields to local object
	// template.Name = input.Name
	if input.ServiceCategoryID != 0 {
		// 允许模块的服务分类信息与模板的服务分类信息不一致，模块同步按钮会调整模块的分类信息, 详情见 issue #2927
		template.ServiceCategoryID = input.ServiceCategoryID

		// validate service category id field
		category, err := p.GetServiceCategory(ctx, template.ServiceCategoryID)
		if err != nil {
			blog.Errorf("UpdateServiceTemplate failed, category id invalid, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
			return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
		if category.BizID != 0 && category.BizID != template.BizID {
			blog.Errorf("UpdateServiceTemplate failed, category biz id and template not equal, err: %+v, rid: %s", err, ctx.ReqID)
			return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKServiceCategoryIDField)
		}
		isLeafNode, err := p.IsServiceCategoryLeafNode(ctx, template.ServiceCategoryID)
		if err != nil {
			blog.Errorf("UpdateServiceTemplate failed, check leaf node failed, err: %+v, rid: %s", err, ctx.ReqID)
			return nil, err
		}
		if isLeafNode == false {
			return nil, ctx.Error.CCError(common.CCErrCoreServiceOnlyNodeServiceCategoryAvailable)
		}
	}

	if field, err := template.Validate(); err != nil {
		blog.Errorf("UpdateServiceTemplate failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	// do update
	filter := map[string]int64{common.BKFieldID: templateID}
	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Update(ctx, filter, template); nil != err {
		blog.Errorf("UpdateServiceTemplate failed, mongodb failed, table: %s, filter: %+v, template: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, template, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return template, nil
}

func (p *processOperation) ListServiceTemplates(ctx core.ContextParams, option metadata.ListServiceTemplateOption) (*metadata.MultipleServiceTemplate, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField: option.BusinessID,
	}

	// filter with matching any sub category
	if option.ServiceCategoryID != nil && *option.ServiceCategoryID > 0 {
		categoriesWithSts, err := p.ListServiceCategories(ctx, option.BusinessID, false)
		if err != nil {
			blog.Errorf("ListServiceTemplates failed, ListServiceCategories failed, err: %+v, rid: %s", err, ctx.ReqID)
			return nil, err
		}
		childrenIDs := make([]int64, 0)
		childrenIDs = append(childrenIDs, *option.ServiceCategoryID)
		for {
			pre := len(childrenIDs)
			for _, categoryWithSts := range categoriesWithSts.Info {
				category := categoryWithSts.ServiceCategory
				if category.ParentID == 0 {
					continue
				}
				if util.InArray(category.ParentID, childrenIDs) == true && util.InArray(category.ID, childrenIDs) == false {
					childrenIDs = append(childrenIDs, category.ID)
				}
			}
			if pre == len(childrenIDs) {
				break
			}
		}
		filter[common.BKServiceCategoryIDField] = map[string][]int64{
			common.BKDBIN: childrenIDs,
		}
	}

	if option.ServiceTemplateIDs != nil {
		filter[common.BKFieldID] = map[string][]int64{
			common.BKDBIN: option.ServiceTemplateIDs,
		}
	}

	if len(option.Search) > 0 {
		filter[common.BKFieldName] = map[string]interface{}{
			common.BKDBLIKE: fmt.Sprintf(".*%s.*", option.Search),
		}
	}

	var total uint64
	var err error
	if total, err = p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).Count(ctx.Context); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	sort := "-id"
	if len(option.Page.Sort) > 0 {
		sort = option.Page.Sort
	}
	templates := make([]metadata.ServiceTemplate, 0)
	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(filter).Start(uint64(option.Page.Start)).Limit(uint64(option.Page.Limit)).Sort(sort).All(ctx.Context, &templates); nil != err {
		blog.Errorf("ListServiceTemplates failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, filter, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	result := &metadata.MultipleServiceTemplate{
		Count: total,
		Info:  templates,
	}
	return result, nil
}

func (p *processOperation) DeleteServiceTemplate(ctx core.ContextParams, serviceTemplateID int64) errors.CCErrorCoder {
	template, err := p.GetServiceTemplate(ctx, serviceTemplateID)
	if err != nil {
		blog.Errorf("DeleteServiceTemplate failed, GetServiceTemplate failed, templateID: %d, err: %+v, rid: %s", serviceTemplateID, err, ctx.ReqID)
		return err
	}

	// service template that referenced by process template shouldn't be removed
	usageFilter := map[string]int64{
		common.BKServiceTemplateIDField: template.ID,
	}
	usageCount, e := p.dbProxy.Table(common.BKTableNameServiceInstance).Find(usageFilter).Count(ctx.Context)
	if nil != e {
		blog.Errorf("DeleteServiceTemplate failed, mongodb failed, table: %s, process template usageFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, usageFilter, e, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceTemplate failed, forbidden delete service template be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	// service template that referenced by module shouldn't be removed
	usageCount, e = p.dbProxy.Table(common.BKTableNameBaseModule).Find(usageFilter).Count(ctx.Context)
	if nil != e {
		blog.Errorf("DeleteServiceTemplate failed, mongodb failed, table: %s, module usageFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceInstance, usageFilter, e, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceTemplate failed, forbidden delete service template be referenced, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	deleteFilter := map[string]int64{common.BKFieldID: template.ID}
	if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Delete(ctx, deleteFilter); nil != err {
		blog.Errorf("DeleteServiceTemplate failed, mongodb failed, table: %s, deleteFilter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, deleteFilter, err, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
