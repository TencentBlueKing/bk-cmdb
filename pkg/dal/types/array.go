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

// Package types for data
package types

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ArrayElem defines types that supported by Array.
// NOTE: `NULL`(nil) element not supported currently
type ArrayElem interface {
	int32 | int64 | float32 | float64 | string | bool | []byte
}

// Array generic data type for database array.
// NOTE: if length of slice is 0, it is converted to empty array `{}`, including nil slice,
// if you want to convert to `NULL`, please use pointer of Array `*Array[T]` as data type
type Array[T ArrayElem] []T

// NewArray create a new array from given slice
func NewArray[T ArrayElem](s []T) Array[T] {
	a := make(Array[T], len(s))
	copy(a, s)
	return a
}

// Value return array value, implement driver.Valuer interface
func (arr Array[T]) Value() (driver.Value, error) {
	// use empty array as default to avoid value to string type '<nil>',
	if len(arr) == 0 {
		return driver.Value("{}"), nil
	}
	// use array encoder from lib/pq
	pqArr := pq.Array(([]T)(arr))
	return pqArr.Value()
}

// Scan value into Array[T], implements sql.Scanner interface
func (arr *Array[T]) Scan(value any) error {
	if arr == nil {
		return errors.New("can not scan to nil array")
	}

	// 用反射获取数组arr的数据类型
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("failed to unmarshal array value:", value))
	}

	// TODO use another decoder to support NULL(nil) array value
	pqArr := pq.Array([]T{})
	err := pqArr.Scan(bytes)
	if err != nil {
		return err
	}
	// copy data to current array
	sliceVal := reflect.ValueOf(pqArr).Elem()
	*arr = make([]T, sliceVal.Len())
	reflect.Copy(reflect.ValueOf(*arr), sliceVal)
	return nil
}

// GormDataType gorm common data type
func (Array[T]) GormDataType() string {
	return getArrayDataType[T]()
}

// GormDBDataType gorm db data type
func (Array[T]) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	return getArrayDataType[T]()
}

func getArrayDataType[T any]() string {
	t := reflect.TypeFor[T]()
	switch t.Kind() {
	case reflect.Bool:
		return "boolean[]"
	case reflect.Int32:
		return "integer[]"
	case reflect.Int64:
		return "bigint[]"
	case reflect.Float32:
		return "real[]"
	case reflect.Float64:
		return "double precision[]"
	case reflect.String:
		return "text[]"
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "bytea[]"
		}
		return "UNKNOWN Slice"
	default:
		return "UNKNOWN"
	}
}

// GormValue gorm value
func (arr Array[T]) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, err := arr.Value()
	if err != nil {
		data = fmt.Errorf("<Array Value() failed: %s>", err.Error())
	}
	switch v := data.(type) {
	case string:
		return gorm.Expr("?", v)
	case []byte:
		return gorm.Expr("?", string(v))
	}
	return gorm.Expr("?", fmt.Sprint(data))
}
