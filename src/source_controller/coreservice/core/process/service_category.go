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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

const (
	// firstLevelParentId serverCategory first level's parent id
	firstLevelParentId = 0
)

func (p *processOperation) CreateServiceCategory(kit *rest.Kit, category metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	// base attribute validate
	if field, err := category.Validate(); err != nil {
		blog.Errorf("CreateServiceCategory failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	var bizID int64
	var err error
	if bizID, err = p.validateBizID(kit, category.BizID); err != nil {
		blog.Errorf("CreateServiceCategory failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}

	category.BizID = bizID

	category.RootID = 0
	if category.ParentID > 0 {
		parentCategory, err := p.GetServiceCategory(kit, category.ParentID)
		if err != nil {
			blog.Errorf("CreateServiceCategory failed, parent id invalid, code: %d, category: %+v, err: %+v, rid: %s", common.CCErrCommParamsInvalid, category, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
		}
		category.RootID = parentCategory.RootID
	}

	// check name unique in business scope
	var count uint64
	filter := map[string]interface{}{
		"name": category.Name,
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: []int64{0, category.BizID},
		},
	}
	// judge levels of classification, ( bk_parent_id == 0 is 1 level, other is 2 level)
	if category.ParentID == 0 {
		filter[common.BKParentIDField] = firstLevelParentId
	} else {
		filter[common.BKParentIDField] = map[string]interface{}{
			common.BKDBNE: firstLevelParentId,
		}
	}
	if count, err = mongodb.Client().Table(common.BKTableNameServiceCategory).Find(filter).Count(kit.Ctx); nil != err {
		blog.Errorf("CreateServiceCategory failed, mongodb query failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("CreateServiceCategory failed, category name duplicated, already exist %d, rid: %s", count, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceServiceCategoryNameDuplicated, category.Name)
	}

	// generate id field
	id, err := mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameServiceCategory)
	if nil != err {
		blog.Errorf("CreateServiceCategory failed, generate id failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGenerateRecordIDFailed)
	}
	category.ID = int64(id)

	// set root id
	if category.RootID == 0 {
		category.RootID = category.ID
	}

	category.IsBuiltIn = false
	category.SupplierAccount = kit.SupplierAccount

	if err := mongodb.Client().Table(common.BKTableNameServiceCategory).Insert(kit.Ctx, &category); nil != err {
		blog.Errorf("mongodb failed, table: %s, category: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, category, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBInsertFailed)
	}

	return &category, nil
}

func (p *processOperation) IsServiceCategoryLeafNode(kit *rest.Kit, categoryID int64) (bool, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKParentIDField: categoryID,
	}
	count, err := mongodb.Client().Table(common.BKTableNameServiceCategory).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("IsServiceCategoryLeafNode failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, err, kit.Rid)
		return false, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	return count == 0, nil
}

func (p *processOperation) GetServiceCategory(kit *rest.Kit, categoryID int64) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	category := metadata.ServiceCategory{}

	filter := map[string]int64{
		common.BKFieldID: categoryID,
	}
	if err := mongodb.Client().Table(common.BKTableNameServiceCategory).Find(filter).One(kit.Ctx, &category); nil != err {
		blog.Errorf("GetServiceCategory failed, mongodb failed, table: %s, filter: %+v, category: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, category, err, kit.Rid)
		if mongodb.Client().IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &category, nil
}

func (p *processOperation) GetDefaultServiceCategory(kit *rest.Kit) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	category := metadata.ServiceCategory{}

	filter := map[string]interface{}{
		common.BKFieldName: common.DefaultServiceCategoryName,
		common.BKParentIDField: map[string]interface{}{
			common.BKDBNE: 0,
		},
		common.BKAppIDField: 0,
	}
	if err := mongodb.Client().Table(common.BKTableNameServiceCategory).Find(filter).One(kit.Ctx, &category); nil != err {
		blog.Errorf("GetDefaultServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, err, kit.Rid)
		if mongodb.Client().IsNotFoundError(err) {
			return nil, kit.CCError.CCError(common.CCErrCommNotFound)
		}
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	return &category, nil
}

func (p *processOperation) UpdateServiceCategory(kit *rest.Kit, categoryID int64, input metadata.ServiceCategory) (*metadata.ServiceCategory, errors.CCErrorCoder) {
	category, err := p.GetServiceCategory(kit, categoryID)
	if err != nil {
		return nil, err
	}

	// update fields to local object
	category.Name = input.Name
	if field, err := category.Validate(); err != nil {
		blog.Errorf("UpdateServiceCategory failed, validation failed, code: %d, err: %+v, rid: %s", common.CCErrCommParamsInvalid, err, kit.Rid)
		err := kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, field)
		return nil, err
	}

	if category.IsBuiltIn {
		blog.Errorf("UpdateServiceCategory failed, forbidden update built-in category, code: %d, rid: %s", common.CCErrCommOperateBuiltInItemForbidden, kit.Rid)
		err := kit.CCError.CCError(common.CCErrCommOperateBuiltInItemForbidden)
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
	count, e := mongodb.Client().Table(common.BKTableNameServiceCategory).Find(uniqueFilter).Count(kit.Ctx)
	if e != nil {
		blog.Errorf("UpdateServiceCategory failed, mongodb query failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, uniqueFilter, e, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if count > 0 {
		blog.Errorf("UpdateServiceCategory failed, category name duplicated, already exist %d, rid: %s", count, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCoreServiceServiceCategoryNameDuplicated, category.Name)
	}

	// do update
	filter := map[string]int64{common.BKFieldID: categoryID}
	if err := mongodb.Client().Table(common.BKTableNameServiceCategory).Update(kit.Ctx, filter, category); nil != err {
		blog.Errorf("UpdateServiceCategory failed, mongodb failed, table: %s, filter: %+v, category: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, category, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBUpdateFailed)
	}
	return category, nil
}

func (p *processOperation) ListServiceCategories(kit *rest.Kit, bizID int64, withStatistics bool) (*metadata.MultipleServiceCategoryWithStatistics, errors.CCErrorCoder) {
	filter := map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			common.BKDBIN: []int64{bizID, 0},
		},
	}
	categories := make([]metadata.ServiceCategory, 0)
	sort := "name"
	if err := mongodb.Client().Table(common.BKTableNameServiceCategory).Find(filter).Sort(sort).All(kit.Ctx, &categories); nil != err {
		blog.Errorf("ListServiceCategories failed, mongodb failed, filter: %+v, category: %+v, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, filter, categories, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}

	usageMap := map[int64]int64{}
	if withStatistics {
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
		if err := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(templateFilter).All(kit.Ctx, &serviceTemplates); nil != err {
			blog.Errorf("ListServiceCategories failed, find reference templates failed, mongodb failed, filter: %+v, table: %s, err: %+v, rid: %s", templateFilter, common.BKTableNameServiceTemplate, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		for _, tpl := range serviceTemplates {
			count, exist := usageMap[tpl.ServiceCategoryID]
			if !exist {
				usageMap[tpl.ServiceCategoryID] = 1
				continue
			}
			usageMap[tpl.ServiceCategoryID] = count + 1
		}

		modules := make([]metadata.ModuleInst, 0)
		if err := mongodb.Client().Table(common.BKTableNameBaseModule).Find(templateFilter).All(kit.Ctx, &modules); nil != err {
			blog.Errorf("ListServiceCategories failed, find reference modules failed, mongodb failed, filter: %+v, table: %s, err: %+v, rid: %s", templateFilter, common.BKTableNameBaseModule, err, kit.Rid)
			return nil, kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
		}
		for _, module := range modules {
			count, exist := usageMap[module.ServiceCategoryID]
			if !exist {
				usageMap[module.ServiceCategoryID] = 1
				continue
			}
			usageMap[module.ServiceCategoryID] = count + 1
		}
	}

	categoriesWithStatistics := make([]metadata.ServiceCategoryWithStatistics, 0)
	for _, category := range categories {
		count := usageMap[category.ID]
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

func (p *processOperation) DeleteServiceCategory(kit *rest.Kit, categoryID int64) errors.CCErrorCoder {
	category, err := p.GetServiceCategory(kit, categoryID)
	if err != nil {
		blog.Errorf("DeleteServiceCategory failed, GetServiceCategory failed, categoryID: %d, err: %+v, rid: %s", categoryID, err, kit.Rid)
		return err
	}

	// 允许全局模式下删除
	/*
		if category.IsBuiltIn == true {
			blog.Errorf("DeleteServiceCategory failed, forbidden delete built-in category, code: %d, rid: %s", common.CCErrCommOperateBuiltInItemForbidden, kit.Rid)
			err := kit.CCError.CCError(common.CCErrCommOperateBuiltInItemForbidden)
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
	childrenCount, e := mongodb.Client().Table(common.BKTableNameServiceCategory).Find(childrenFilter).Count(kit.Ctx)
	if nil != e {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, childrenFilter, e, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBSelectFailed)
	}
	if childrenCount > 0 {
		blog.Errorf("DeleteServiceCategory failed, forbidden delete category has children node, childrenCount: %d, code: %d, rid: %s", childrenCount, common.CCErrCommOperateBuiltInItemForbidden, kit.Rid)
		err := kit.CCError.CCError(common.CCErrCommRemoveRecordHasChildrenForbidden)
		return err
	}

	// category that referenced by service template shouldn't be removed
	usageFilter := map[string]int64{common.BKServiceCategoryIDField: category.ID}
	usageCount, e := mongodb.Client().Table(common.BKTableNameServiceTemplate).Find(usageFilter).Count(kit.Ctx)
	if e != nil {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, usageFilter, e, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceCategory failed, forbidden delete category be referenced by service template, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, kit.Rid)
		err := kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	// category that referenced by service template shouldn't be removed
	usageFilter = map[string]int64{common.BKServiceCategoryIDField: category.ID}
	usageCount, e = mongodb.Client().Table(common.BKTableNameBaseModule).Find(usageFilter).Count(kit.Ctx)
	if e != nil {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceTemplate, usageFilter, e, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommDBDeleteFailed)
	}
	if usageCount > 0 {
		blog.Errorf("DeleteServiceCategory failed, forbidden delete category be referenced by module, code: %d, rid: %s", common.CCErrCommRemoveRecordHasChildrenForbidden, kit.Rid)
		err := kit.CCError.CCError(common.CCErrCommRemoveReferencedRecordForbidden)
		return err
	}

	deleteFilter := map[string]int64{common.BKFieldID: category.ID}
	if e := mongodb.Client().Table(common.BKTableNameServiceCategory).Delete(kit.Ctx, deleteFilter); e != nil {
		blog.Errorf("DeleteServiceCategory failed, mongodb failed, table: %s, filter: %+v, err: %+v, rid: %s", common.BKTableNameServiceCategory, deleteFilter, e, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBDeleteFailed)
	}
	return nil
}
