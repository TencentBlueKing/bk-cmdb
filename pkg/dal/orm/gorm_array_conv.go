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

package orm

import (
	"database/sql/driver"
	"fmt"
	"reflect"

	"github.com/lib/pq"
	"gorm.io/gorm/clause"

	"github.com/TencentBlueKing/bk-cmdb/pkg/filter"
)

// arrayRuleToClauseExpr convert to postgresql array expression
// ref: https://www.postgresql.org/docs/current/functions-array.html
func arrayRuleToClauseExpr(rule *filter.AtomRule) (clause.Expression, error) {
	arrVal, err := buildArraySQL(rule.Value)
	if err != nil {
		return nil, fmt.Errorf("build array value sql failed: %w", err)
	}
	switch filter.OpType(rule.Op) {
	case filter.ArrayEqual:
		return buildArrayEqual(rule, arrVal)
	case filter.ArrayNotEqual:
		return buildArrayNotEqual(rule, arrVal)
	case filter.ArrayContains:
		return buildArrayContains(rule, arrVal)
	case filter.ArraySubset:
		return buildArraySubset(rule, arrVal)
	case filter.ArrayOverlap:
		return buildArrayOverlap(rule, arrVal)
	default:
		return nil, fmt.Errorf("not support array op %s", rule.Op)
	}
}

func buildArrayEqual(rule *filter.AtomRule, val string) (clause.Expression, error) {
	return NewArrayQuery(rule.Field).Equal(val), nil
}

func buildArrayNotEqual(rule *filter.AtomRule, val string) (clause.Expression, error) {
	return NewArrayQuery(rule.Field).NotEqual(val), nil
}

func buildArrayOverlap(rule *filter.AtomRule, val string) (clause.Expression, error) {
	return NewArrayQuery(rule.Field).Overlap(val), nil
}

func buildArraySubset(rule *filter.AtomRule, val string) (clause.Expression, error) {
	return NewArrayQuery(rule.Field).Subset(val), nil
}

func buildArrayContains(rule *filter.AtomRule, val string) (clause.Expression, error) {
	return NewArrayQuery(rule.Field).Contains(val), nil
}

var rTypeStringSlice = reflect.TypeFor[[]string]()
var rTypeInt32Slice = reflect.TypeFor[[]int32]()
var rTypeInt64Slice = reflect.TypeFor[[]int64]()
var rTypeFloat32Slice = reflect.TypeFor[[]float32]()
var rTypeFloat64Slice = reflect.TypeFor[[]float64]()
var rTypeByteaSlice = reflect.TypeFor[[][]byte]()
var rTypeBoolSlice = reflect.TypeFor[[]bool]()

// ArrayQuery array query
type ArrayQuery struct {
	column any
	value  any
	op     string
}

// NewArrayQuery ...
func NewArrayQuery(column any) *ArrayQuery {
	return &ArrayQuery{
		column: column,
	}
}

// Equal equal operator =
func (a *ArrayQuery) Equal(arrayText any) *ArrayQuery {
	a.op = "="
	a.value = arrayText
	return a
}

// NotEqual not equal operator <>
func (a *ArrayQuery) NotEqual(arrayText any) *ArrayQuery {
	a.op = "<>"
	a.value = arrayText
	return a
}

// Contains contains operator @>
func (a *ArrayQuery) Contains(arrayText any) *ArrayQuery {
	a.op = "@>"
	a.value = arrayText
	return a
}

// Subset subset operator <@
func (a *ArrayQuery) Subset(arrayText any) *ArrayQuery {
	a.op = "<@"
	a.value = arrayText
	return a
}

// Overlap overlap operator &&
func (a *ArrayQuery) Overlap(arrayText any) *ArrayQuery {
	a.op = "&&"
	a.value = arrayText
	return a
}

// Build sql, implements clause.Expression interface
func (a ArrayQuery) Build(builder clause.Builder) {

	builder.WriteQuoted(a.column)
	_ = builder.WriteByte(' ')
	_, _ = builder.WriteString(a.op)
	_ = builder.WriteByte(' ')
	builder.AddVar(builder, a.value)
}

func buildArraySQL(value any) (sql string, err error) {

	if value == nil {
		return "", nil
	}

	rVal := reflect.ValueOf(value)
	rKind := rVal.Kind()
	if rKind == reflect.Pointer {
		rVal = rVal.Elem()
		rKind = rVal.Kind()
	}

	var elemKind reflect.Kind
	var elemType reflect.Type
	switch rKind {
	case reflect.Slice, reflect.Array:
		elemType = rVal.Type().Elem()
		elemKind = elemType.Kind()
	default:
		return "", fmt.Errorf("value is not slice or array: %T", value)
	}

	var valuer driver.Valuer
	switch elemKind {
	case reflect.String:
		valuer = buildStringArray(elemType, rVal)
	case reflect.Int, reflect.Int32, reflect.Int64:
		valuer = buildIntArray(elemType, rVal)
	case reflect.Float32, reflect.Float64:
		valuer = buildFloatArray(elemType, rVal)
	case reflect.Bool:
		valuer = buildBoolArray(elemType, rVal)
	case reflect.Array, reflect.Slice:
		if elemType.Elem().Kind() == reflect.Uint8 {
			valuer = buildByteaArray(elemType.Elem(), rVal)
			break
		}
		fallthrough
	default:
		return "", fmt.Errorf("not support array elem kind %s", elemKind)
	}

	val, err := valuer.Value()
	if err != nil {
		return "", err
	}

	s, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("valuer.Value() is not string: %T", val)
	}
	return s, nil
}

func buildIntArray(elemType reflect.Type, rVal reflect.Value) driver.Valuer {
	if elemType == rTypeInt32Slice {
		return pq.Int32Array(rVal.Interface().([]int32))
	}
	if elemType == rTypeInt64Slice {
		return pq.Int64Array(rVal.Interface().([]int64))
	}
	// convert to int64 slice
	newSlice := reflect.MakeSlice(rTypeInt64Slice, rVal.Len(), rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		newSlice.Index(i).SetInt(rVal.Index(i).Int())
	}
	return pq.Int64Array(newSlice.Interface().([]int64))

}

func buildByteaArray(elemType reflect.Type, rVal reflect.Value) (valuer driver.Valuer) {
	if elemType == rTypeByteaSlice {
		return pq.ByteaArray(rVal.Interface().([][]byte))
	}
	newSlice := reflect.MakeSlice(rTypeByteaSlice, rVal.Len(), rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		newSlice.Index(i).SetBytes(rVal.Index(i).Bytes())
	}
	return pq.ByteaArray(newSlice.Interface().([][]byte))
}

func buildBoolArray(elemType reflect.Type, rVal reflect.Value) (valuer driver.Valuer) {
	if elemType == rTypeBoolSlice {
		return pq.BoolArray(rVal.Interface().([]bool))
	}
	newSlice := reflect.MakeSlice(rTypeBoolSlice, rVal.Len(), rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		newSlice.Index(i).SetBool(rVal.Index(i).Bool())
	}
	return pq.BoolArray(newSlice.Interface().([]bool))
}

func buildStringArray(elemType reflect.Type, rVal reflect.Value) (valuer driver.Valuer) {
	if elemType == rTypeStringSlice {
		return pq.StringArray(rVal.Interface().([]string))
	}
	newSlice := reflect.MakeSlice(rTypeStringSlice, rVal.Len(), rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		newSlice.Index(i).SetString(rVal.Index(i).String())
	}
	return pq.StringArray(newSlice.Interface().([]string))
}

func buildFloatArray(elemType reflect.Type, rVal reflect.Value) driver.Valuer {
	if elemType == rTypeFloat32Slice {
		return pq.Float32Array(rVal.Interface().([]float32))
	}
	if elemType == rTypeFloat64Slice {
		return pq.Float64Array(rVal.Interface().([]float64))
	}
	newSlice := reflect.MakeSlice(rTypeFloat64Slice, rVal.Len(), rVal.Len())
	for i := 0; i < rVal.Len(); i++ {
		newSlice.Index(i).SetFloat(rVal.Index(i).Float())
	}
	return pq.Float64Array(newSlice.Interface().([]float64))

}
