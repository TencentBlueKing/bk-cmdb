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
	"sync"

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

	// BatchCreate create given items, will auto generate id
	BatchCreate(kt *kit.Kit, items *structs.Slice) (ids []string, err error)
	List(kt *kit.Kit, opt *types.ListOption, opts ...Option) (*types.DynamicListDetails, error)
	Delete(kt *kit.Kit, filterExpr filter.RuleFactory, opts ...Option) (deleted int64, err error)
	Update(kt *kit.Kit, filterExpr filter.RuleFactory, model *structs.Struct, opts ...Option) (updated int64, err error)
	UpdateByID(kt *kit.Kit, id string, model *structs.Struct) (updated int64, err error)
}

var _ Dynamic = &dynamicDao{}
var _ DynamicConstructor = &dynamicConstructor{}

type dynamicConstructor struct {
	orm   orm.Interface
	idgen idgenerator.Interface
	lock  sync.Mutex
}

// NewDynamicConstructor returns a new DynamicConstructor
func NewDynamicConstructor(orm orm.Interface, idg idgenerator.Interface) DynamicConstructor {
	constructor := &dynamicConstructor{
		orm:   orm,
		idgen: idg,
		lock:  sync.Mutex{},
	}
	return constructor
}

type dynamicDao struct {
	modelName string
	orm       orm.Interface
	idgen     idgenerator.Interface
	log       log.Logger
}

// Dynamic get new Dynamic dao with model
func (d *dynamicConstructor) Dynamic(model string) (Dynamic, error) {
	// check model builder
	_, ok := structs.GetBuilder(model)
	if !ok {
		return nil, fmt.Errorf("dynamic model not found: %s", model)
	}

	logger := log.With("model", model)
	dynamic := &dynamicDao{
		modelName: model,
		orm:       d.orm,
		idgen:     d.idgen,
		log:       logger,
	}
	return dynamic, nil
}

// AutoTxn auto start transaction to execute fn then commit, rollback if error
func (d *dynamicDao) AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error, opts ...*sql.TxOptions) error {
	return d.orm.DB(orm.WithContext(kt)).Transaction(func(tx *gorm.DB) error {
		return fn(d.orm.WithTx(tx))
	}, opts...)

}

// WithTx creates a new Dynamic dao with the given transaction
func (d *dynamicDao) WithTx(tx orm.Interface) Dynamic {
	return &dynamicDao{
		modelName: d.modelName,
		orm:       tx,
		idgen:     d.idgen,
	}
}

// BatchCreate ...
func (d *dynamicDao) BatchCreate(kt *kit.Kit, items *structs.Slice) (ids []string, err error) {
	builder, ok := d.GetBuilder()
	if !ok {
		return nil, fmt.Errorf("dynamic model not found: %s", d.modelName)
	}
	// 1. check type
	if !builder.OfSlice(items) {
		return nil, fmt.Errorf("items is not slice of model: %s", d.modelName)
	}

	// 2. check id field
	if !items.HaveField("ID") {
		return nil, fmt.Errorf("model %s should have ID field", d.modelName)
	}
	// 3. generate id
	ids, err = d.idgen.Batch(kt, types.Name(d.modelName), uint64(items.Len()))
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
			d.log.Error(kt, "dynamic dao fail to set id", log.E(err))
			return nil, fmt.Errorf("fail to set id: %w", err)
		}
	}
	// 5. create
	ret := d.orm.DB(orm.WithContext(kt)).Table(d.modelName).Create(items.Value())
	if ret.Error != nil {
		d.log.Error(kt, "dynamic dao fail to create items", log.E(ret.Error))
		return nil, ret.Error
	}

	return ids, nil
}

// List ...
func (d *dynamicDao) List(kt *kit.Kit, listOpt *types.ListOption, opts ...Option) (*types.DynamicListDetails, error) {
	builder, ok := d.GetBuilder()
	if !ok {
		return nil, fmt.Errorf("dynamic model not found: %s", d.modelName)
	}

	// get model struct info
	modelInfo, ok := table.GetModelTypeInfo(types.Name(d.modelName))
	if !ok {
		d.log.Error(kt, "model type info not found for dynamic list ")
		return nil, fmt.Errorf("model type info not found for dynamic list: %s", d.modelName)
	}

	// get option and validate list option for dynamic
	typeMap := modelInfo.GetColumnTypeMap()
	eo, po := GetConfig(typeMap, opts)
	if err := listOpt.Validate(eo, po); err != nil {
		return nil, err
	}
	// convert filter
	expr, err := conv.Filter(listOpt.Filter)
	if err != nil {
		d.log.Error(kt, "fail to convert filter to clause expression", "listOpt", listOpt, log.E(err))
		return nil, fmt.Errorf("fail to convert filter to clause expression: %w", err)
	}

	db := d.orm.DB().Table(d.modelName).Where(expr)
	// for count request
	if listOpt.Page.Count {
		var count int64
		err = db.Count(&count).Error
		if err != nil {
			d.log.Error(kt, "fail to count table", log.E(err))
			return nil, err
		}
		return &types.DynamicListDetails{Count: uint64(count)}, nil
	}
	// for list request
	if len(listOpt.Fields) > 0 {
		db = db.Select(listOpt.Fields)
	}
	data := builder.NewSlice(0, int(listOpt.Page.Limit))
	err = db.Clauses(conv.Page(listOpt.Page, po)...).Find(data.Pointer()).Error
	if err != nil {
		d.log.Error(kt, "fail to list table", "expr", expr, log.E(err))
		return nil, err
	}

	results := &types.DynamicListDetails{
		Details: data,
	}
	return results, nil
}

// Delete ...
func (d *dynamicDao) Delete(kt *kit.Kit, flt filter.RuleFactory, opts ...Option) (deleted int64, err error) {
	expr, err := d.convFilter(kt, flt, opts)
	if err != nil {
		return 0, fmt.Errorf("fail to convert filter to clause expression for delete: %w", err)
	}
	ret := d.orm.DB(orm.WithContext(kt)).Table(d.modelName).Clauses(expr).Delete(nil)
	return ret.RowsAffected, ret.Error
}

// Update ...
func (d *dynamicDao) Update(kt *kit.Kit, flt filter.RuleFactory, model *structs.Struct, opts ...Option) (
	updated int64, err error) {

	builder, ok := d.GetBuilder()
	if !ok {
		return 0, fmt.Errorf("dynamic model not found: %s", d.modelName)
	}
	// 1.conv filter
	expr, err := d.convFilter(kt, flt, opts)
	if err != nil {
		return 0, fmt.Errorf("fail to convert filter to clause expression for update: %w", err)
	}
	// 2. check builder type
	if !builder.OfStruct(model) {
		return 0, fmt.Errorf("items is not a struct instance of model: %s", d.modelName)
	}
	ret := d.orm.DB(orm.WithContext(kt)).Table(d.modelName).Clauses(expr).Updates(model.Value())
	return ret.RowsAffected, ret.Error
}

// UpdateByID update by id
func (d *dynamicDao) UpdateByID(kt *kit.Kit, id string, model *structs.Struct) (updated int64, err error) {
	// 1. get modelBuilder
	builder, ok := d.GetBuilder()
	if !ok {
		return 0, fmt.Errorf("dynamic model not found: %s", d.modelName)
	}
	if !builder.OfStruct(model) {
		return 0, fmt.Errorf("update items is not struct of model: %s", d.modelName)
	}
	ret := d.orm.DB(orm.WithContext(kt)).Table(d.modelName).Where("id = ?", id).Updates(model.Value())
	return ret.RowsAffected, ret.Error
}

func (d *dynamicDao) convFilter(kt *kit.Kit, flt filter.RuleFactory, opts []Option) (clause.Expression, error) {
	// get model struct info
	modelInfo, ok := table.GetModelTypeInfo(types.Name(d.modelName))
	if !ok {
		d.log.Error(kt, "model type info not found")
		return nil, fmt.Errorf("model type info not found: %s", d.modelName)
	}
	typeMap := modelInfo.GetColumnTypeMap()
	// ignore page option here
	eo, _ := GetConfig(typeMap, opts)
	if err := flt.Validate(eo); err != nil {
		return nil, err
	}
	// convert filter
	exp, err := conv.Filter(flt)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

// GetBuilder ...
func (d *dynamicDao) GetBuilder() (*structs.Builder, bool) {
	return structs.GetBuilder(d.modelName)
}
