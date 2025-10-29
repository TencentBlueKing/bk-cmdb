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
	"log/slog"

	"gorm.io/gorm"

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

// IDSetterModel 对类型指针的约束
type IDSetterModel[T any] interface {
	table.Tabler
	SetID(id string)
	// SetID requires a pointer to T
	*T
}

// Table returns a slog attribute for the given table name
func Table[T ~string](name T) slog.Attr {
	return slog.String("table", string(name))
}

// Generic generic dao interface
type Generic[T any, PT IDSetterModel[T]] interface {
	// AutoTxn auto start transaction to execute fn then commit, rollback if error
	AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error, opts ...*sql.TxOptions) error
	// WithTx creates a new GenericDao with the given transaction
	WithTx(tx orm.Interface) Generic[T, PT]
	BatchCreate(kt *kit.Kit, items []T) (ids []string, err error)
	List(kt *kit.Kit, opt *types.ListOption) (*types.ListDetails[T], error)
	Delete(kt *kit.Kit, filterExpr *filter.Expression) (updated int, err error)
	Update(kt *kit.Kit, filterExpr *filter.Expression, model PT) (updated int, err error)
	UpdateByID(kt *kit.Kit, id string, model PT) (updated int, err error)
}

var _ Generic[table.TestModel, *table.TestModel] = new(GenericDao[table.TestModel, *table.TestModel])

// GenericDao dao for generic type
type GenericDao[T any, PT IDSetterModel[T]] struct {
	Orm   orm.Interface
	IDGen idgenerator.Interface
}

// NewGenericDao returns a new GenericDao
func NewGenericDao[T any, PT IDSetterModel[T]](orm orm.Interface, idGen idgenerator.Interface) *GenericDao[T, PT] {
	return &GenericDao[T, PT]{
		Orm:   orm,
		IDGen: idGen,
	}
}

// GetTableName get table name
func (d *GenericDao[T, PT]) GetTableName() table.Name {
	return table.Name(PT(new(T)).TableName())
}

// AutoTxn auto start transaction to execute fn then commit, rollback if error
func (d *GenericDao[T, PT]) AutoTxn(kt *kit.Kit, fn func(tx orm.Interface) error,
	opts ...*sql.TxOptions) error {

	return d.Orm.DB().WithContext(kt).Transaction(func(tx *gorm.DB) error {
		return fn(d.Orm.WithTx(tx))
	}, opts...)
}

// WithTx creates a new GenericDao with the given transaction
func (d *GenericDao[T, PT]) WithTx(tx orm.Interface) Generic[T, PT] {
	return &GenericDao[T, PT]{
		Orm: tx,
		// IDGen should not use transaction, to avoid long time lock
		IDGen: d.IDGen,
	}
}

// BatchCreate ...
func (d *GenericDao[T, PT]) BatchCreate(kt *kit.Kit, items []T) (ids []string, err error) {
	// generate ids
	tableName := d.GetTableName()
	ids, err = d.IDGen.Batch(kt, tableName, uint64(len(items)))
	if err != nil {
		log.Error(kt, "fail to generate ids for table", "table", tableName.String(), log.E(err))
		return nil, err
	}
	for i := range items {
		PT(&items[i]).SetID(ids[i])
	}
	err = d.Orm.DB().Create(items).Error
	if err != nil {
		log.Error(kt, "fail to create table", Table(tableName), log.E(err))
		return nil, err
	}
	return ids, nil
}

// List ...
func (d *GenericDao[T, PT]) List(kt *kit.Kit, opt *types.ListOption) (*types.ListDetails[T], error) {
	tableName := d.GetTableName()
	fieldsInfo, ok := table.GetModelTypeInfo(kt, tableName)
	if !ok {
		return nil, fmt.Errorf("fail to get table fields for %s", tableName)
	}
	option := filter.NewExprOption(filter.RuleFields(fieldsInfo.GetColumnTypeMap()))
	err := opt.Validate(option, types.NewDefaultPageOption())
	if err != nil {
		return nil, fmt.Errorf("invalid list option: %w", err)
	}

	expr, err := conv.Filter(opt.Filter)
	if err != nil {
		log.Error(kt, "fail to convert filter to clause expression", "opt", opt, log.E(err))
		return nil, fmt.Errorf("fail to convert filter to clause expression: %v", err)
	}

	db := d.Orm.DB().Model(new(T)).Where(expr)
	if len(opt.Fields) > 0 {
		db = db.Select(opt.Fields)
	}

	// for count request
	if opt.Page.Count {
		var count int64
		err = db.Count(&count).Error
		if err != nil {
			log.Error(kt, "fail to count table", Table(tableName), log.E(err))
			return nil, err
		}
		return &types.ListDetails[T]{Count: uint64(count)}, nil
	}

	data := make([]T, 0)
	// for list request
	err = db.Find(&data).Error
	if err != nil {
		log.Error(kt, "fail to list table", Table(tableName),
			log.E(err), slog.Any("expr", expr))
		return nil, err
	}

	results := &types.ListDetails[T]{
		Details: data,
	}
	return results, nil
}

// Delete ...
func (d *GenericDao[T, PT]) Delete(kt *kit.Kit, filterExpr *filter.Expression) (deleted int, err error) {
	fieldsInfo, ok := table.GetModelTypeInfo(kt, d.GetTableName())
	if !ok {
		return 0, fmt.Errorf("fail to get table fields for table")
	}
	option := filter.NewExprOption(filter.RuleFields(fieldsInfo.GetColumnTypeMap()))
	err = filterExpr.Validate(option)
	if err != nil {
		return 0, fmt.Errorf("invalid delete option: %w", err)
	}

	expr, err := conv.Filter(filterExpr)
	if err != nil {
		log.Error(kt, "fail to convert filter to clause expression",
			log.E(err), slog.Any("filter", filterExpr))
		return 0, fmt.Errorf("fail to convert filter to clause expression: %v", err)
	}

	do := gen.Use(d.Orm.DB()).WithContext(kt).TestModel
	ret, err := do.Clauses(expr).Delete()
	if err != nil {
		log.Error(kt, "fail to delete table", log.E(err))
		return 0, err
	}
	return int(ret.RowsAffected), nil
}

// Update ...
func (d *GenericDao[T, PT]) Update(kt *kit.Kit, filterExpr *filter.Expression, model PT) (
	updated int, err error) {

	fieldInfo, ok := table.GetModelTypeInfo(kt, d.GetTableName())
	if !ok {
		return 0, fmt.Errorf("fail to get table fields for table")
	}
	option := filter.NewExprOption(filter.RuleFields(fieldInfo.GetColumnTypeMap()))
	err = filterExpr.Validate(option)
	if err != nil {
		return 0, fmt.Errorf("invalid update option: %w", err)
	}

	expr, err := conv.Filter(filterExpr)
	if err != nil {
		log.Error(kt, "fail to convert filter to clause expression",
			log.E(err), slog.Any("filter", filterExpr))
		return 0, fmt.Errorf("fail to convert filter to clause expression: %v", err)
	}

	do := gen.Use(d.Orm.DB()).WithContext(kt).TestModel
	ret, err := do.Clauses(expr).Updates(model)
	if err != nil {
		log.Error(kt, "fail to update table", log.E(err))
		return 0, err
	}
	return int(ret.RowsAffected), nil
}

// UpdateByID ...
func (d *GenericDao[T, PT]) UpdateByID(kt *kit.Kit, id string, model PT) (
	updated int, err error) {

	info, err := gen.Use(d.Orm.DB()).WithContext(kt).TestModel.
		Where(gen.TestModel.BaseID.Eq(id)).
		Updates(model)
	if err != nil {
		log.Error(kt, "fail to update table", log.E(err))
		return 0, err
	}
	return int(info.RowsAffected), nil
}
