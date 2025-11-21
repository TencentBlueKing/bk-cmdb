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

// Valuer represents a struct field value.
type Valuer interface {
	Raw() any
	Int64() int64
	String() string
	Float64() float64
	Bool() bool
	Time() time.Time
	Map() map[string]any
}

// NewValue creates a new Valuer.
// NOTE: the valuer only support the method with the same type, calling other method will result in panic.
func NewValue(v reflect.Value) Valuer {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &int64Value{value[int64]{raw: v.Int()}}
	case reflect.String:
		return &stringValue{value[string]{raw: v.String()}}
	case reflect.Float32, reflect.Float64:
		return &float64Value{value[float64]{raw: v.Float()}}
	case reflect.Bool:
		return &boolValue{value[bool]{raw: v.Bool()}}
	default:
		raw := v.Interface()

		switch t := raw.(type) {
		case time.Time:
			return &timeValue{value[time.Time]{raw: t}}
		case map[string]any:
			return &mapValue{value[map[string]any]{raw: t}}
		default:
			return &value[any]{raw: t}
		}
	}
}

type value[T any] struct {
	raw T
}

// Raw returns the raw value.
func (v *value[T]) Raw() any { return v.raw }

// Int64 not supported by value[T].
func (v *value[T]) Int64() int64 { panic("value type is not int64") }

// String not supported by value[T].
func (v *value[T]) String() string { panic("value type is not string") }

// Float64 not supported by value[T].
func (v *value[T]) Float64() float64 { panic("value type is not float64") }

// Bool not supported by value[T].
func (v *value[T]) Bool() bool { panic("value type is not bool") }

// Time not supported by value[T].
func (v *value[T]) Time() time.Time { panic("value type is not time.Time") }

// Map not supported by value[T].
func (v *value[T]) Map() map[string]any { panic("value type is not map[string]any") }

type int64Value struct {
	value[int64]
}

// Int64 returns the int64 value.
func (v *int64Value) Int64() int64 { return v.raw }

type stringValue struct {
	value[string]
}

// String returns the string value.
func (v *stringValue) String() string { return v.raw }

type float64Value struct {
	value[float64]
}

// Float64 returns the float64 value.
func (v *float64Value) Float64() float64 { return v.raw }

type boolValue struct {
	value[bool]
}

// Bool returns the bool value.
func (v *boolValue) Bool() bool { return v.raw }

type timeValue struct {
	value[time.Time]
}

// Time returns the time.Time value.
func (v *timeValue) Time() time.Time { return v.raw }

type mapValue struct {
	value[map[string]any]
}

// Map returns the map[string]any value.
func (v *mapValue) Map() map[string]any { return v.raw }
