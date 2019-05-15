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
	"configcenter/src/common/mapstr"
	"fmt"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/source_controller/coreservice/core"
)

func (p *processOperation) CreateServiceCategory(ctx core.ContextParams, category metadata.ServiceCategory) (*metadata.ServiceCategory, error) {
	// base attribute validate
	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("validation failed, err: %+v", err)
	}
	
	if _, err := p.validateBizID(ctx, category.Metadata); err != nil {
		return nil, err
	}

	category.RootID = 0
	if category.ParentID > 0 {
		parentCategory, err := p.GetServiceCategory(ctx, category.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent id invalid, err: %+v", err)
		}
		category.RootID = parentCategory.RootID
	}

	// generate id field
	id, err := p.dbProxy.NextSequence(ctx, common.BKTableNameServiceCategory)
	if nil != err {
		return nil, err
	}
	category.ID = int64(id)
	
	// set root id
	if category.RootID == 0 {
		category.RootID = category.ID
	}
	
	category.SupplierAccount = ctx.SupplierAccount
	
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Insert(ctx.Context, &category); nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, err.Error(), ctx.ReqID)
		return  nil, err
	}
	
	return &category, nil
}

func (p *processOperation) GetServiceCategory(ctx core.ContextParams, categoryID int64) (*metadata.ServiceCategory, error) {
	category := metadata.ServiceCategory{}
	
	filter := condition.CreateCondition()
	filter.Field(common.BKFieldID).Eq(categoryID)
	
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter.ToMapStr()).One(ctx.Context, &category); nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, err.Error(), ctx.ReqID)
		return  nil, err
	}
	
	return &category, nil
}

func (p *processOperation) UpdateServiceCategory(ctx core.ContextParams, categoryID int64, input metadata.ServiceCategory) (*metadata.ServiceCategory, error) {
	category, err := p.GetServiceCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	// update fields to local object
	category.Name = input.Name
	
	// do update
	filter := condition.CreateCondition()
	filter.Field(common.BKFieldID).Eq(categoryID)
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Update(ctx, filter.ToMapStr(), category); nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, err.Error(), ctx.ReqID)
		return  nil, err
	}
	return category, nil
}

func (p *processOperation) ListServiceCategories(ctx core.ContextParams, bizID int64, withStatistics bool) (*metadata.MultipleServiceCategory, error) {
	md := metadata.NewMetaDataFromBusinessID(strconv.FormatInt(bizID, 10))
	filter := map[string]mapstr.MapStr{
		"metadata": md.ToMapStr(),
	}
	
	
	categories := make([]metadata.ServiceCategory, 0)
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Find(filter).All(ctx.Context, &categories); nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, err.Error(), ctx.ReqID)
		return  nil, err
	}
	
	result := &metadata.MultipleServiceCategory{
		Count: int64(len(categories)),
		Info:  categories,
	}
	return result, nil
}

func (p *processOperation) DeleteServiceCategory(ctx core.ContextParams, categoryID int64) error {
	_, err := p.GetServiceCategory(ctx, categoryID)
	if err != nil {
		return err
	}
	
	// TODO: support delete sub tree at same time 

	filter := condition.CreateCondition()
	filter.Field(common.BKFieldID).Eq(categoryID)
	if err := p.dbProxy.Table(common.BKTableNameServiceCategory).Delete(ctx, filter.ToMapStr()); nil != err {
		blog.Errorf("mongodb failed, table: %s, err: %+v, rid: %s", common.BKTableNameServiceCategory, err.Error(), ctx.ReqID)
		return  err
	}
	return nil
}
