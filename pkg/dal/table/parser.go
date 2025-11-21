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

// Package table ...
package table

import (
	"context"
	"encoding/json"
	"reflect"
	"sync"
	"time"

	"gorm.io/gorm/schema"

	"github.com/TencentBlueKing/bk-cmdb/pkg/dal/types"
	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/structs"
)

type tableFieldStore struct {
	tableMap sync.Map
}

// Store table field definition
func (m *tableFieldStore) Store(table types.Name, val *ModelTypeInfo) {
	m.tableMap.Store(table, val)
}

// CAS compare and swap table field definition
func (m *tableFieldStore) CAS(table types.Name, old, new *ModelTypeInfo) bool {
	return m.tableMap.CompareAndSwap(table, old, new)
}

// Load table field definition
func (m *tableFieldStore) Load(table types.Name) (*ModelTypeInfo, bool) {
	value, ok := m.tableMap.Load(table)
	if !ok {
		return nil, false
	}
	modelInfo, ok := value.(*ModelTypeInfo)
	return modelInfo, ok
}

// Delete table field definition
func (m *tableFieldStore) Delete(table types.Name) {
	m.tableMap.Delete(table)
}

// tableFields store table field definition
var tableFields = tableFieldStore{}

// FieldInfo store field type information
type FieldInfo struct {
	Name        string
	ColumnType  filter.FieldType
	ReflectType reflect.Type
}

// ModelTypeInfo store field type information of table fields
type ModelTypeInfo struct {
	fields       []FieldInfo
	fieldTypeMap map[string]FieldInfo
	builder      *structs.Builder
}

// GetBuilder return builder
func (mt *ModelTypeInfo) GetBuilder() *structs.Builder {
	return mt.builder
}

// GetFieldMetas return fields
func (mt *ModelTypeInfo) GetFieldMetas() []FieldInfo {
	return mt.fields
}

// GetColumnTypeMap return column type map
func (mt *ModelTypeInfo) GetColumnTypeMap() map[string]filter.FieldType {
	var colMap = make(map[string]filter.FieldType, len(mt.fields))
	for i := range mt.fields {
		colMap[mt.fields[i].Name] = mt.fields[i].ColumnType
	}
	return colMap
}

// GetFieldByName return field type by field name
func (mt *ModelTypeInfo) GetFieldByName(name string) (*FieldInfo, bool) {
	info, ok := mt.fieldTypeMap[name]
	return &info, ok
}

// NewModelTypeInfo creates a new ModelTypeInfo with the given fields
func NewModelTypeInfo(fields []FieldInfo) *ModelTypeInfo {
	return &ModelTypeInfo{
		fields:       fields,
		fieldTypeMap: mergeColumnTypes(fields),
	}
}

// GetModelTypeInfo try to get table field definition from cache
// Note: table must be registered first
func GetModelTypeInfo(table types.Name) (*ModelTypeInfo, bool) {
	return tableFields.Load(table)
}

// load all static table struct, note it should only be used to register static table during startup
// if table parse failed, it will panic
func (m *tableFieldStore) loadAllStaticTableStruct() {
	for tableName, tableStruct := range GetAllStaticTables() {
		m.loadTableStruct(tableName, tableStruct)
	}
}

// loadTableStruct load single table struct note it should only be used to register static table during startup,
// if table parse failed, it will panic
func (m *tableFieldStore) loadTableStruct(tableName types.Name, tableStruct any) {
	ft, err := parseTableStruct(tableStruct)
	if err != nil {
		log.Error(context.TODO(), "fail to parse table for struct", "name", tableName, log.E(err))
		panic(err)
	}
	m.Store(tableName, ft)
}

func mergeColumnTypes(columns []FieldInfo) map[string]FieldInfo {
	columnMap := make(map[string]FieldInfo)
	for _, column := range columns {
		columnMap[column.Name] = column
	}
	return columnMap
}

// used by schema.ParseWithSpecialTableName
var cacheStore = sync.Map{}

// used by schema.ParseWithSpecialTableName
var namingStrategy = schema.NamingStrategy{IdentifierMaxLength: 64}

func parseTableStruct(tableStruct any) (*ModelTypeInfo, error) {
	// ref: gorm.io/gen ConvertStructs
	tableSchema, err := schema.ParseWithSpecialTableName(tableStruct, &cacheStore, namingStrategy, "")
	if err != nil {
		return nil, err
	}

	infos := make([]FieldInfo, 0)
	for _, field := range tableSchema.Fields {
		column := field.DBName
		columnType := getColumnType(field.FieldType)
		infos = append(infos, FieldInfo{Name: column, ColumnType: columnType, ReflectType: field.FieldType})
	}

	ft := NewModelTypeInfo(infos)

	return ft, nil
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

// onModelBuilderChange update table field definition when model builder changes
func (m *tableFieldStore) onModelBuilderChange(name string, newBuilder *structs.Builder) {
	tableStruct := newBuilder.New().Value()
	ft, err := parseTableStruct(tableStruct)
	if err != nil {
		log.Error(context.TODO(), "fail to parse table from struct",
			"table", name, "struct", tableStruct, log.E(err))
		return
	}
	// same table should have same field definition, it should be ok to store directly
	m.Store(types.Name(name), ft)
}

func init() {
	tableFields.loadAllStaticTableStruct()
	// register change handler to update table field definition when model builder changes
	structs.RegisterChangeHandler(tableFields.onModelBuilderChange)
}
