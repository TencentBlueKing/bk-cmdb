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

package table

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

type tableFieldStore struct {
	tableMap sync.Map
}

var tableFields = tableFieldStore{}

// FieldInfo store field type information
type FieldInfo struct {
	Name        string
	ColumnType  filter.FieldType
	ReflectType reflect.Type
}

// ModelTypeInfo store field type information of table fields
type ModelTypeInfo struct {
	fields       []*FieldInfo
	fieldTypeMap map[string]*FieldInfo
	builder      *structs.Builder
}

// GetBuilder return builder
func (mt *ModelTypeInfo) GetBuilder() *structs.Builder {
	return mt.builder
}

// GetFieldMetas return fields
func (mt *ModelTypeInfo) GetFieldMetas() []*FieldInfo {
	return mt.fields
}

// GetColumnTypeMap return column type map
func (mt *ModelTypeInfo) GetColumnTypeMap() map[string]filter.FieldType {
	var colMap = make(map[string]filter.FieldType, len(mt.fields))
	for i := 0; i < len(mt.fields); i++ {
		colMap[mt.fields[i].Name] = mt.fields[i].ColumnType
	}
	return colMap
}

// GetFieldByName return field type by field name
func (mt *ModelTypeInfo) GetFieldByName(name string) (*FieldInfo, bool) {
	info, ok := mt.fieldTypeMap[name]
	return info, ok
}

// NewModelTypeInfo creates a new ModelTypeInfo with the given fields
func NewModelTypeInfo(fields []*FieldInfo) *ModelTypeInfo {
	return &ModelTypeInfo{
		fields:       fields,
		fieldTypeMap: mergeColumnTypes(fields),
	}
}

// GetModelTypeInfo try to get table field definition from cache
// Note: table must be registered first
func GetModelTypeInfo(ctx context.Context, table Name) (*ModelTypeInfo, bool) {
	return tableFields.getModelInfo(ctx, table)
}

// getModelInfo try to get table fields type map from cache, if not exist, try to parse table struct and store to cache
func (m *tableFieldStore) getModelInfo(ctx context.Context, tableName Name) (*ModelTypeInfo, bool) {
	value, ok := m.tableMap.Load(tableName)
	if !ok {
		return m.reloadTableStruct(ctx, tableName)
	}

	modelInfo, ok := value.(*ModelTypeInfo)
	if !ok {
		log.Error(ctx, "try get cached table value, but type mismatch",
			"type", reflect.TypeOf(value).String(), "table", tableName)
		return nil, false
	}
	if modelInfo == nil {
		log.Error(ctx, "try get cached table model info, but got nil", "table", tableName)
		return nil, false
	}
	if modelInfo.builder != nil && modelInfo.builder.Invalid() {
		return m.reloadTableStruct(ctx, tableName)
	}
	return modelInfo, true
}

func (m *tableFieldStore) reloadTableStruct(ctx context.Context, tableName Name) (*ModelTypeInfo, bool) {
	tableStruct, registered := tableRegistry[tableName]
	if !registered {
		log.Warn(ctx, "try to get unregistered table", "table", tableName)
		return nil, false
	}
	if builder, isBuilder := tableStruct.(*structs.Builder); isBuilder {
		if builder.Invalid() {
			// try to reload builder from structs if invalid
			var exists bool
			builder, exists = structs.GetBuilder(builder.Name())
			if !exists {
				log.Error(ctx, "fail to get builder from struct after invalid", "table", builder.Name())
				return nil, false
			}
		}
		tableStruct = builder.New().Value()
	}
	columns, err := parseTableStruct(tableStruct)
	if err != nil {
		log.Error(ctx, "fail to parse table from struct", "table", tableName, "struct", tableStruct, log.E(err))
		return nil, false
	}

	ft := NewModelTypeInfo(columns)
	// same table should have same field definition, it should be ok to store directly
	m.tableMap.Store(tableName, ft)
	return ft, true
}

func mergeColumnTypes(columns []*FieldInfo) map[string]*FieldInfo {
	columnMap := make(map[string]*FieldInfo)
	for _, column := range columns {
		columnMap[column.Name] = column
	}
	return columnMap
}

func parseTableStruct(tableStruct any) ([]*FieldInfo, error) {
	// ref: gorm.io/gen ConvertStructs
	dummyDB, _ := gorm.Open(tests.DummyDialector{})
	stmt := gorm.Statement{DB: dummyDB}
	err := stmt.Parse(tableStruct)
	if err != nil {
		return nil, err
	}

	infos := make([]*FieldInfo, 0)
	for _, field := range stmt.Schema.Fields {
		column := field.DBName
		columnType := getColumnType(field.FieldType)
		infos = append(infos, &FieldInfo{Name: column, ColumnType: columnType, ReflectType: field.FieldType})
	}

	return infos, nil
}

var jsonNumberReflectType = reflect.TypeFor[json.Number]()
var timeReflectType = reflect.TypeFor[time.Time]()

func getColumnType(field reflect.Type) filter.FieldType {
	var fieldType filter.FieldType

	switch field {
	case jsonNumberReflectType:
		return filter.Numeric
	case timeReflectType:
		return filter.Time
	default:
		// continue check kind below
	}

	fieldKind := field.Kind()
	switch fieldKind {
	case reflect.String:
		fieldType = filter.String
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		fieldType = filter.Numeric
	case reflect.Bool:
		fieldType = filter.Boolean
	default:
		fieldType = filter.Any
	}
	return fieldType
}
