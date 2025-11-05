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

package structs

import (
	"reflect"
	"time"
)

// FieldType defines the struct field type.
type FieldType string

const (
	// Int64Type is the int64 struct field type.
	Int64Type FieldType = "int64"
	// StringType is the string struct field type.
	StringType FieldType = "string"
	// Float64Type is the float64 struct field type.
	Float64Type FieldType = "float64"
	// BoolType is the bool struct field type.
	BoolType FieldType = "bool"
	// TimeType is the time.Time struct field type.
	TimeType FieldType = "time"
	// MapType is the map[string]any struct field type.
	MapType FieldType = "map"
)

// fieldTypeRegistry is the mapping of field type to reflection type.
var fieldTypeRegistry = map[FieldType]reflect.Type{
	Int64Type:   reflect.TypeFor[int64](),
	StringType:  reflect.TypeFor[string](),
	Float64Type: reflect.TypeFor[float64](),
	BoolType:    reflect.TypeFor[bool](),
	TimeType:    reflect.TypeFor[time.Time](),
	MapType:     reflect.TypeFor[map[string]any](),
}

// RegisterFieldType registers field type and its corresponding reflection type to the registry.
func RegisterFieldType(fieldType FieldType, typ reflect.Type) {
	fieldTypeRegistry[fieldType] = typ
}

// GetFieldType get reflection type by field type.
func GetFieldType(fieldType FieldType) (reflect.Type, bool) {
	typ, exists := fieldTypeRegistry[fieldType]
	return typ, exists
}

// Field is the metadata of a struct field.
type Field struct {
	// Name is the field name, the field name must be exported.
	Name string
	// Type is the field type.
	Type FieldType
	// IsSlice defines whether the field is a slice.
	IsSlice bool
	// Tags is the struct tag key-value pairs, for example json->test will be transformed to json:"test" in struct tag.
	Tags map[string]string
	// Anonymous defines whether the field is an anonymous field.
	Anonymous bool

	// Validator is an optional validation function for the field value.
	Validator func(any) error
}
