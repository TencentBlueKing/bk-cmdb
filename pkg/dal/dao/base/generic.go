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

// Package base defines the basic dao
package base

import (
	"database/sql"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	idgenerator "github.com/TencentBlueKing/bk-cmdb/pkg/dal/dao/id-generator"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/gen"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/orm/conv"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/table"
	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// IDSettable 可以被设置ID的类型，应该被指针类型实现
type IDSettable interface {
	SetID(id string)
}

// Generic 泛型通用DAO接口, T应为非指针类型, 且*T应该实现 IDSettable 接口,否则BatchCreate方法无法设置ID
type Generic[T types.Tabler] interface {
	// AutoTxn auto start transaction to execute fn then commit, rollback if error
	AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error, opts ...*sql.TxOptions) error
	// WithTx creates a new GenericDao with the given transaction
	WithTx(tx orm.Interface) Generic[T]

	// BatchCreate create given items, will auto generate id
	BatchCreate(kt *kit.Kit, items []T) (ids []string, err error)
	// List items by filter and page option
	List(kt *kit.Kit, opt *types.ListOption, opts ...Option) (*types.ListDetails[T], error)
	// Delete delete items by filter
	Delete(kt *kit.Kit, filterExpr *filter.Expression, opts ...Option) (updated int64, err error)
	// Update update model items by filter
	Update(kt *kit.Kit, filterExpr *filter.Expression, model *T, opts ...Option) (updated int64, err error)
	// UpdateByID update model item by id
	UpdateByID(kt *kit.Kit, id string, model *T) (updated int64, err error)
}

var _ Generic[table.TestModel] = new(GenericDao[table.TestModel])

// GenericDao dao for generic type
type GenericDao[T types.Tabler] struct {
	Orm   orm.Interface
	IDGen idgenerator.Interface
	Config
	log log.Logger
}

// NewGenericDao returns a new GenericDao
func NewGenericDao[T types.Tabler](orm orm.Interface, idGen idgenerator.Interface, opts ...Option) *GenericDao[T] {
	cfg := &Config{}
	for _, opt := range opts {
		opt(cfg)
	}

	d := &GenericDao[T]{
		Orm:    orm,
		IDGen:  idGen,
		Config: *cfg,
	}
	logger := log.With("table", d.GetTableName())
	d.log = logger
	return d
}

// GetTableName get table name
func (d *GenericDao[T]) GetTableName() types.Name {
	var t T
	return types.Name(t.TableName())
}

// AutoTxn auto start transaction to execute fn then commit, rollback if error
func (d *GenericDao[T]) AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error, opts ...*sql.TxOptions) error {
	return d.Orm.DB(orm.WithContext(kt)).Transaction(func(tx *gorm.DB) error {
		return fn(d.Orm.WithTx(tx))
	}, opts...)
}

// WithTx creates a new GenericDao with the given transaction
func (d *GenericDao[T]) WithTx(tx orm.Interface) Generic[T] {
	return &GenericDao[T]{
		Orm: tx,
		// IDGen should not use transaction, to avoid long time lock
		IDGen: d.IDGen,
	}
}

// BatchCreate ...
func (d *GenericDao[T]) BatchCreate(kt *kit.Kit, items []T) (ids []string, err error) {
	// type check, T should not be ptr, and its pointer type should implement SetID(string) method
	if err := checkIDSettable[T](); err != nil {
		return nil, err
	}

	// generate ids
	tableName := d.GetTableName()
	ids, err = d.IDGen.Batch(kt, tableName, uint64(len(items)))
	if err != nil {
		d.log.Error(kt, "fail to generate ids for table", "table", tableName, log.E(err))
		return nil, err
	}
	for i := range items {
		// type already checked
		any(&items[i]).(IDSettable).SetID(ids[i])
	}
	err = d.Orm.DB().Create(items).Error
	if err != nil {
		d.log.Error(kt, "fail to create table", log.E(err))
		return nil, err
	}
	return ids, nil
}

func checkIDSettable[T any]() error {
	var t T
	if _, ok := any(&t).(IDSettable); !ok {
		return fmt.Errorf("item's pointer type should implement SetID method, got %T", t)
	}
	return nil
}

// List ...
func (d *GenericDao[T]) List(kt *kit.Kit, listOpt *types.ListOption, opts ...Option) (*types.ListDetails[T], error) {
	// get model struct info
	tableName := d.GetTableName()
	modelInfo, ok := table.GetModelTypeInfo(tableName)
	if !ok {
		d.log.Error(kt, "model type info not found for generic list")
		return nil, fmt.Errorf("model type info not found for generic list: %s", tableName)
	}

	// get option and validate list option for dynamic
	typeMap := modelInfo.GetColumnTypeMap()
	eo, po := GetConfig(typeMap, opts)
	if err := listOpt.Validate(eo, po); err != nil {
		return nil, err
	}

	expr, err := conv.Filter(listOpt.Filter)
	if err != nil {
		d.log.Error(kt, "fail to convert filter to clause expression", "listOpt", listOpt, log.E(err))
		return nil, fmt.Errorf("fail to convert filter to clause expression: %w", err)
	}

	db := d.Orm.DB().Model(new(T)).Where(expr)

	// for count request
	if listOpt.Page.Count {
		var count int64
		err = db.Count(&count).Error
		if err != nil {
			d.log.Error(kt, "fail to count table", log.E(err))
			return nil, err
		}
		return &types.ListDetails[T]{Count: uint64(count)}, nil
	}
	// for list request
	if len(listOpt.Fields) > 0 {
		db = db.Select(listOpt.Fields)
	}
	data := make([]T, 0)
	err = db.Clauses(conv.Page(listOpt.Page, po)...).Find(&data).Error
	if err != nil {
		d.log.Error(kt, "fail to list table", "expr", expr, log.E(err))
		return nil, err
	}

	results := &types.ListDetails[T]{
		Details: data,
	}
	return results, nil
}

// Delete ...
func (d *GenericDao[T]) Delete(kt *kit.Kit, flt *filter.Expression, opts ...Option) (deleted int64, err error) {
	expr, _, err := d.ConvFilter(kt, flt, opts)
	if err != nil {
		d.log.Error(kt, "generic delete fail to convert filter to clause expression", log.E(err), "filter", flt)
		return 0, err
	}

	do := gen.Use(d.Orm.DB()).WithContext(kt).TestModel
	ret, err := do.Clauses(expr).Delete()
	if err != nil {
		d.log.Error(kt, "fail to delete table", log.E(err))
		return 0, err
	}
	return ret.RowsAffected, nil
}

// Update ...
func (d *GenericDao[T]) Update(kt *kit.Kit, flt *filter.Expression, model *T, opts ...Option) (
	updated int64, err error) {

	expr, _, err := d.ConvFilter(kt, flt, opts)
	if err != nil {
		d.log.Error(kt, "generic update fail to convert filter to clause expression", log.E(err), "filter", flt)
		return 0, err
	}

	do := gen.Use(d.Orm.DB()).WithContext(kt).TestModel
	ret, err := do.Clauses(expr).Updates(model)
	if err != nil {
		d.log.Error(kt, "fail to update table", log.E(err))
		return 0, err
	}
	return ret.RowsAffected, nil
}

// UpdateByID ...
func (d *GenericDao[T]) UpdateByID(kt *kit.Kit, id string, model *T) (updated int64, err error) {
	info, err := gen.Use(d.Orm.DB()).WithContext(kt).TestModel.
		Where(gen.TestModel.BaseID.Eq(id)).
		Updates(model)
	if err != nil {
		d.log.Error(kt, "fail to update table", log.E(err))
		return 0, err
	}
	return info.RowsAffected, nil
}

// ConvFilter conv and check filter
func (d *GenericDao[T]) ConvFilter(kt *kit.Kit, flt filter.RuleFactory, opts []Option) (
	clause.Expression, *filter.ExprOption, error) {

	// get table fields definition
	tableName := d.GetTableName()
	fieldsInfo, ok := table.GetModelTypeInfo(tableName)
	if !ok {
		return nil, nil, fmt.Errorf("fail to get table fields for %s", tableName)
	}
	// ignore page option here
	eo, _ := GetConfig(fieldsInfo.GetColumnTypeMap(), opts)
	if err := flt.Validate(eo); err != nil {
		return nil, nil, err
	}
	// convert filter to clause expression
	expr, err := conv.Filter(flt)
	if err != nil {
		return nil, nil, fmt.Errorf("fail to convert filter to clause expression: %w", err)
	}

	return expr, eo, nil
}
