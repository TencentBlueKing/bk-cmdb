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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateServiceCategory(ctx core.ContextParams, category metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := category.Validate(); err != nil {
		blog.Errorf("CreateServiceCategory failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(ctx, category.BizID); err != nil {
		blog.Errorf("CreateServiceCategory failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	category.BizID = bizID

	category.RootID = 0
	if category.ParentID > 0 {
		parentCategory, err := p.GetServiceCategory(ctx, category.ParentID)
		if err != nil {
			blog.Errorf("CreateServiceCategory failed, parent id invalid, code: %d, category: %+v, err: %+v, rid: %s", common.CCErrCommParamsInvalid, category, err, ctx.ReqID)
			return nil, ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		}
		category.RootID = parentCategory.RootID
	}

	// check name unique in business scope
	var count uint64
	filter := map[string]interface{}{
		common.BKParentIDField: category.ParentID,
		"name":                 category.Name,
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: []int64{0, category.BizID},
		},
	}
	if count, err = p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter).Count(ctx); nil != err {
		blog.Errorf("CreateServiceCategory failed, mongodb query failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("CreateServiceCategory failed, category name duplicated, already exist %d, rid: %s", count, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCoreServiceServiceCategoryNameDuplicated, category.Name)
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameServiceCategory)
	if nil != err {
		blog.Errorf("CreateServiceCategory failed, generate id failed, err: %+v, rid: %s", err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	category.ID = int64(id)

	// set root id
	if category.RootID == 0 {
		category.RootID = category.ID
	}

	category.IsBuiltIn = false
	category.SupplierAccount = ctx.SupplierAccount

	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Insert(ctx.Context, &category); nil != err {
		blog.Errorf("mongodb failed, table: %s, category: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, category, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return &category, nil
}

func (p *processOperation) IsServiceCategoryLeafNode(ctx core.ContextParams, categoryID int64) (bool, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKParentIDField: categoryID,
	}
	count, err := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter).Count(ctx)
	if err != nil {
		blog.Errorf("IsServiceCategoryLeafNode failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, err, ctx.ReqID)
		return false, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	return count == 0, nil
}

func (p *processOperation) GetServiceCategory(ctx core.ContextParams, categoryID int64) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	category := metadata.ServiceCategory{}

	filter := map[string]int64{
		common.BKFieldID: categoryID,
	}
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter).One(ctx.Context, &category); nil != err {
		blog.Errorf("GetServiceCategory failed, mongodb failed, table: %s, filter: %+v, category: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, category, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &category, nil
}

func (p *processOperation) GetDefaultServiceCategory(ctx core.ContextParams) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	category := metadata.ServiceCategory{}

	filter := map[string]interface{}{
		common.BKFieldName: common.DefaultServiceCategoryName,
		common.BKParentIDField: map[string]interface{}{
			common.BKDBNE: 0,
		},
		common.BKAppIDField: 0,
	}
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter).One(ctx.Context, &category); nil != err {
		blog.Errorf("GetDefaultServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, err, ctx.ReqID)
		if p.dbProxy.IsNotFoundError(err) {
			return nil, ctx.Error.CCError(common.CCErrCommNotFound)
		}
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &category, nil
}

func (p *processOperation) UpdateServiceCategory(ctx core.ContextParams, categoryID int64, input metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	category, err := p.GetServiceCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// update fields to local object
	category.Name = input.Name
	if field, err := category.Validate(); err != nil {
		blog.Errorf("UpdateServiceCategory failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, ctx.ReqID)
		err := ctx.Error.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	if category.IsBuiltIn == true {
		blog.Errorf("UpdateServiceCategory failed, forbidden update built-in category, code: %d, rid: %s", common.CCErrCommOperateBuiltInItemForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommOperateBuiltInItemForbidden)
		return nil, err
	}

	// check name unique in business scope
	uniqueFilter := map[string]interface{}{
		common.BKFieldName:     category.Name,
		common.BKParentIDField: category.ParentID,
		common.BKFieldID: map[string]interface{}{
			common.BKDBNE: categoryID,
		},
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: []int64{0, category.BizID},
		},
	}
	count, e := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(uniqueFilter).Count(ctx)
	if e != nil {
		blog.Errorf("UpdateServiceCategory failed, mongodb query failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, uniqueFilter, e, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("UpdateServiceCategory failed, category name duplicated, already exist %d, rid: %s", count, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCoreServiceServiceCategoryNameDuplicated, category.Name)
	}

	// do update
	filter := map[string]int64{common.BKFieldID: categoryID}
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Update(ctx, filter, category); nil != err {
		blog.Errorf("UpdateServiceCategory failed, mongodb failed, table: %s, filter: %+v, category: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, category, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return category, nil
}

func (p *processOperation) ListServiceCategories(ctx core.ContextParams, bizID int64, withStatistics bool) (*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: []int64{bizID, 0},
		},
	}
	categories := make([]metadata.ServiceCategory, 0)
	sort := "name"
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter).Sort(sort).All(ctx.Context, &categories); nil != err {
		blog.Errorf("ListServiceCategories failed, mongodb failed, filter: %+v, category: %+v, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, categories, err, ctx.ReqID)
		return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	usageMap := map[int64]int64{}
	if withStatistics == true {
		categoryIDs := make([]int64, 0)
		for _, category := range categories {
			categoryIDs = append(categoryIDs, category.ID)
		}
		templateFilter := map[string]interface{}{
			common.BKServiceCategoryIDField: map[string]interface{}{
				common.BKDBIN: categoryIDs,
			},
			common.BKAppIDField: bizID,
		}
		serviceTemplates := make([]metadata.ServiceTemplate, 0)
		if err := p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(templateFilter).All(ctx.Context, &serviceTemplates); nil != err {
			blog.Errorf("ListServiceCategories failed, find reference templates failed, mongodb failed, filter: %+v, table: %s, err: %+v, rid: %s", common.BKTableNameServiceTemplate, serviceTemplates, err, ctx.ReqID)
			return nil, ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		for _, tpl := range serviceTemplates {
			count, exist := usageMap[tpl.ServiceCategoryID]
			if exist == false {
				usageMap[tpl.ServiceCategoryID] = 1
				continue
			}
			usageMap[tpl.ServiceCategoryID] = count + 1
		}
	}

	categoriesWithStatistics := make([]metadata.ServiceCategoryWithStatistics, 0)
	for _, category := range categories {
		count, _ := usageMap[category.ID]
		categoriesWithStatistics = append(categoriesWithStatistics, metadata.ServiceCategoryWithStatistics{
			ServiceCategory: category,
			UsageAmount:     count,
		})
	}
	result := &metadata.MultipleServiceCategoryWithStatistics{
		Count: int64(len(categoriesWithStatistics)),
		Info:  categoriesWithStatistics,
	}
	return result, nil
}

func (p *processOperation) DeleteServiceCategory(ctx core.ContextParams, categoryID int64) errors.CCErrorCoder {
	category, err := p.GetServiceCategory(ctx, categoryID)
	if err != nil {
		blog.Errorf("DeleteServiceCategory failed, GetServiceCategory failed, categoryID: %d, err: %+v, rid: %s", categoryID, err, ctx.ReqID)
		return err
	}

	// 允许全局模式下删除
	/*
		if category.IsBuiltIn == true {
			blog.Errorf("DeleteServiceCategory failed, forbidden delete built-in category, code: %d, rid: %s", common.CCErrCommOperateBuiltInItemForbidden, ctx.ReqID)
			err := ctx.Error.CCError(common.CCErrCommOperateBuiltInItemForbidden)
			return err
		}
	*/

	// category that has sub category shouldn't be removed
	childrenFilter := map[string]interface{}{
		common.BKParentIDField: category.ID,
		common.BKFieldID: map[string]interface{}{
			common.BKDBNE: category.ID,
		},
	}
	childrenCount, e := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(childrenFilter).Count(ctx.Context)
	if nil != e {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, childrenFilter, e, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if childrenCount > 0 {
		blog.Errorf("DeleteServiceCategory failed, forbidden delete category has children node, childrenCount: %d, code: %d, rid: %s", childrenCount, common.CCErrCommOperateBuiltInItemForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveRecordHasChildrenForbidden)
		return err
	}

	// category that referenced by service template shouldn't be removed
	usageFilter := map[string]int64{common.BKServiceCategoryIDField: category.ID}
	usageCount, e := p.dbProxy.Table(common.BKTableNameServiceTemplate).Find(usageFilter).Count(ctx.Context)
	if e != nil {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, usageFilter, e, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceCategory failed, forbidden delete category be referenced by service template, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	// category that referenced by service template shouldn't be removed
	usageFilter = map[string]int64{common.BKServiceCategoryIDField: category.ID}
	usageCount, e = p.dbProxy.Table(common.BKTableNameBaseModule).Find(usageFilter).Count(ctx.Context)
	if e != nil {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, usageFilter, e, ctx.ReqID)
		return ctx.Error.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceCategory failed, forbidden delete category be referenced by module, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, ctx.ReqID)
		err := ctx.Error.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	deleteFilter := map[string]int64{common.BKFieldID: category.ID}
	if e := p.dbProxy.Table(common.BKTableNameServiceCategory).Delete(ctx, deleteFilter); e != nil {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, deleteFilter, e, ctx.ReqID)
		return ctx.Error.CCError(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
