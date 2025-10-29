/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package base

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	idgenerator "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/id-generator"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm/conv"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

// DynamicConstructor build for Dynamic dao
type DynamicConstructor interface {
	// Dynamic get new Dynamic dao with model, model should already be registered in structs
	Dynamic(model string) (Dynamic, error)
}

// Dynamic dao for dynamic model
type Dynamic interface {

	// AutoTxn auto start transaction to execute fn then commit, rollback if error
	AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error, opts ...*sql.TxOptions) error
	// WithTx creates a new Dynamic dao with the given transaction
	WithTx(tx orm.Interface) Dynamic

	/* operations below must be used after Dynamic */

	// BatchCreate batch create
	BatchCreate(kt *kit.Kit, items *structs.Slice) (ids []string, err error)
	List(kt *kit.Kit, opt *types.ListOption) (*types.DynamicListDetails, error)
	Delete(kt *kit.Kit, filterExpr filter.RuleFactory) (updated int, err error)
	Update(kt *kit.Kit, filterExpr filter.RuleFactory, model *structs.Struct) (updated int, err error)
	UpdateByID(kt *kit.Kit, id string, model *structs.Struct) (updated int, err error)
}

var _ Dynamic = &dynamicDao{}

// NewDynamicConstructor returns a new DynamicConstructor
func NewDynamicConstructor(orm orm.Interface, idg idgenerator.Interface) DynamicConstructor {
	return &dynamicDao{
		Orm:   orm,
		IDGen: idg,
	}
}

type dynamicDao struct {
	modelName    string
	modelBuilder *structs.Builder
	Orm          orm.Interface
	IDGen        idgenerator.Interface
}

// Dynamic get new Dynamic dao with model
func (d *dynamicDao) Dynamic(model string) (Dynamic, error) {
	// load model from structs
	builder, ok := structs.GetBuilder(model)
	if !ok {
		return nil, fmt.Errorf("dynamic model not found: %s", model)
	}
	return &dynamicDao{
		modelName:    model,
		modelBuilder: builder,
		Orm:          d.Orm,
		IDGen:        d.IDGen,
	}, nil
}

// AutoTxn auto start transaction to execute fn then commit, rollback if error
func (d *dynamicDao) AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error, opts ...*sql.TxOptions) error {
	return d.Orm.DB().WithContext(kt).Transaction(func(tx *gorm.DB) error {
		return fn(d.Orm.WithTx(tx))
	}, opts...)

}

// WithTx creates a new Dynamic dao with the given transaction
func (d *dynamicDao) WithTx(tx orm.Interface) Dynamic {
	return &dynamicDao{
		modelName:    d.modelName,
		modelBuilder: d.modelBuilder,
		Orm:          tx,
		IDGen:        d.IDGen,
	}
}

// BatchCreate batch create
func (d *dynamicDao) BatchCreate(kt *kit.Kit, items *structs.Slice) (ids []string, err error) {
	// 1. check builder
	_, err = d.getBuilderBySlice(items)
	if err != nil {
		return nil, err
	}

	// 2. check id field
	if !items.HaveField("ID") {
		return nil, fmt.Errorf("model %s not have id field", d.modelName)
	}
	// 3. generate id
	ids, err = d.IDGen.Batch(kt, table.Name(d.modelName), uint64(items.Len()))
	if err != nil {
		return nil, err
	}
	// 4. set id for each item
	for i := range items.Len() {
		s, err := items.GetStruct(i)
		if err != nil {
			return nil, err
		}

		if err = s.Set("ID", ids[i]); err != nil {
			log.Error(kt, "dynamic dao fail to set id", "model", d.modelName, log.E(err))
			return nil, fmt.Errorf("fail to set id: %w", err)
		}
	}
	// 5. create
	ret := d.Orm.DBContext(kt).Table(d.modelName).Create(items.Value())
	if ret.Error != nil {
		log.Error(kt, "dynamic dao fail to create items", "model", d.modelName, log.E(ret.Error))
		return nil, ret.Error
	}

	return ids, nil

}

// List ...
func (d *dynamicDao) List(kt *kit.Kit, opt *types.ListOption) (*types.DynamicListDetails, error) {
	// 1. check modelBuilder
	if err := d.refreshBuilder(); err != nil {
		return nil, fmt.Errorf("dynamic list failed: %w", err)
	}

	modelInfo, ok := table.GetModelTypeInfo(kt, table.Name(d.modelName))
	if !ok {
		return nil, fmt.Errorf("model type info not found for model: %s", d.modelName)
	}
	typeMap := modelInfo.GetColumnTypeMap()
	eo := filter.NewExprOption(filter.RuleFields(typeMap))
	// validate opt
	err := opt.Validate(eo, types.NewDefaultPageOption())
	if err != nil {
		return nil, err
	}

	expr, err := conv.Filter(opt.Filter)
	if err != nil {
		log.Error(kt, "fail to convert filter to clause expression", "opt", opt, log.E(err))
		return nil, fmt.Errorf("fail to convert filter to clause expression: %v", err)
	}

	db := d.Orm.DB().Table(d.modelName).Where(expr)
	if len(opt.Fields) > 0 {
		db = db.Select(opt.Fields)
	}

	// for count request
	if opt.Page.Count {
		var count int64
		err = db.Count(&count).Error
		if err != nil {
			log.Error(kt, "fail to count table", Table(d.modelName), log.E(err))
			return nil, err
		}
		return &types.DynamicListDetails{Count: uint64(count)}, nil
	}

	data := d.modelBuilder.NewSlice(0, int(opt.Page.Limit))
	// for list request
	err = db.Find(data.Pointer()).Error
	if err != nil {
		log.Error(kt, "fail to list table", Table(d.modelName), "expr", expr, log.E(err))
		return nil, err
	}

	results := &types.DynamicListDetails{
		Details: data,
	}
	return results, nil
}

// Delete ...
func (d *dynamicDao) Delete(kt *kit.Kit, flt filter.RuleFactory) (deleted int, err error) {
	expr, err := d.convFilter(kt, flt)
	if err != nil {
		return 0, fmt.Errorf("fail to convert filter to clause expression for delete: %w", err)
	}
	ret := d.Orm.DBContext(kt).Table(d.modelName).Clauses(expr).Delete(nil)
	return int(ret.RowsAffected), ret.Error
}

// Update ...
func (d *dynamicDao) Update(kt *kit.Kit, flt filter.RuleFactory, model *structs.Struct) (
	updated int, err error) {

	expr, err := d.convFilter(kt, flt)
	if err != nil {
		return 0, fmt.Errorf("fail to convert filter to clause expression for update: %w", err)
	}
	ret := d.Orm.DBContext(kt).Table(d.modelName).Clauses(expr).Updates(model.Value())
	return int(ret.RowsAffected), ret.Error
}

// UpdateByID update by id
func (d *dynamicDao) UpdateByID(kt *kit.Kit, id string, model *structs.Struct) (updated int, err error) {
	// 1. check modelBuilder
	_, err = d.getBuilderByStruct(model)
	if err != nil {
		return 0, fmt.Errorf("dynamic update by id fail: %w", err)
	}

	ret := d.Orm.DBContext(kt).Table(d.modelName).Where("id = ?", id).Updates(model.Value())
	return int(ret.RowsAffected), ret.Error
}

func (d *dynamicDao) getBuilderBySlice(items *structs.Slice) (*structs.Builder, error) {
	err := d.refreshBuilder()
	if err != nil {
		return nil, err
	}
	if !d.modelBuilder.OfSlice(items) {
		return nil, fmt.Errorf("items is not slice of model: %s", d.modelName)
	}
	return d.modelBuilder, nil
}

func (d *dynamicDao) getBuilderByStruct(items *structs.Struct) (*structs.Builder, error) {
	err := d.refreshBuilder()
	if err != nil {
		return nil, err
	}
	if !d.modelBuilder.Of(items) {
		return nil, fmt.Errorf("items is not struct of model: %s", d.modelName)
	}
	return d.modelBuilder, nil
}

// refreshBuilder refresh model builder if invalid
func (d *dynamicDao) refreshBuilder() error {
	// check is modelBuilder set
	if d.modelName == "" {
		return fmt.Errorf("model builder is unset")
	}

	// refresh model builder if invalid
	if d.modelBuilder == nil || d.modelBuilder.Invalid() {
		builder, ok := structs.GetBuilder(d.modelName)
		if !ok {
			return fmt.Errorf("dynamic model not found after old builder invalid %s", d.modelName)
		}
		d.modelBuilder = builder
	}
	return nil
}
func (d *dynamicDao) convFilter(kt *kit.Kit, flt filter.RuleFactory) (clause.Expression, error) {
	// check is modelBuilder up to date
	if err := d.refreshBuilder(); err != nil {
		return nil, err
	}
	// get model struct info
	modelInfo, ok := table.GetModelTypeInfo(kt, table.Name(d.modelName))
	if !ok {
		return nil, fmt.Errorf("model type info not found for deleting model: %s", d.modelName)
	}

	// convert filter
	typeMap := modelInfo.GetColumnTypeMap()
	eo := filter.NewExprOption(filter.RuleFields(typeMap))
	if err := flt.Validate(eo); err != nil {
		return nil, err
	}

	return conv.Filter(flt)
}
