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
	"fmt"
	"reflect"
)

// Slice defines a dynamic slice.
type Slice struct {
	// name is the identifier of the slice.
	name string
	// data is the underlying slice instance pointer.
	data any
	// val is the reflection value of the slice instance which is used for runtime manipulation.
	val reflect.Value

	// validators maps struct field names to their corresponding validation functions.
	validators map[string]func(any) error
	// fieldIndexMap maps struct field names to their respective field indices in the struct type.
	fieldIndexMap map[string]int
}

// Pointer returns the underlying dynamic slice instance in the form of pointer.
func (s *Slice) Pointer() any {
	return s.data
}

// Value returns the underlying dynamic slice instance in the form of slice type.
func (s *Slice) Value() any {
	return s.val.Interface()
}

// Len returns the length of the slice.
func (s *Slice) Len() int {
	return s.val.Len()
}

// Cap returns the capacity of the slice.
func (s *Slice) Cap() int {
	return s.val.Cap()
}

// Get returns slice element value by index.
func (s *Slice) Get(index int) (any, error) {
	if index < 0 || index >= s.val.Len() {
		return nil, fmt.Errorf("index out of range, max index: %d", s.val.Len()-1)
	}
	return s.val.Index(index).Interface(), nil
}

// GetStruct returns slice element value in the form of struct.
func (s *Slice) GetStruct(index int) (*Struct, error) {
	if index < 0 || index >= s.val.Len() {
		return nil, fmt.Errorf("index out of range, max index: %d", s.val.Len()-1)
	}

	val := s.val.Index(index).Addr()

	return &Struct{
		data:          val.Interface(),
		val:           val.Elem(),
		validators:    s.validators,
		fieldIndexMap: s.fieldIndexMap,
	}, nil
}

// Set sets the value of the slice element at the specified index.
func (s *Slice) Set(index int, value any) error {
	if index < 0 || index >= s.val.Len() {
		return fmt.Errorf("index out of range, max index: %d", s.val.Len()-1)
	}

	if reflect.TypeOf(value) != s.val.Type().Elem() {
		return fmt.Errorf("cannot set slice element, type %s != %s", reflect.TypeOf(value), s.val.Type().Elem())
	}

	s.val.Index(index).Set(reflect.ValueOf(value))

	return nil
}

// SetStruct sets the value of the slice element in the form of struct.
func (s *Slice) SetStruct(index int, value *Struct) error {
	if index < 0 || index >= s.val.Len() {
		return fmt.Errorf("index out of range, max index: %d", s.val.Len()-1)
	}

	if value.val.Type() != s.val.Type().Elem() {
		return fmt.Errorf("cannot set slice element, type %s != %s", value.val.Type().Elem(), s.val.Type().Elem())
	}

	s.val.Index(index).Set(value.val)

	return nil
}

// Append appends a value to the end of the slice.
func (s *Slice) Append(value ...any) error {
	for i, v := range value {
		if reflect.TypeOf(v) != s.val.Type().Elem() {
			return fmt.Errorf("cannot append to slice, type %s != %s at index %d", reflect.TypeOf(v),
				s.val.Type().Elem(), i)
		}

		s.val = reflect.Append(s.val, reflect.ValueOf(v))
	}
	s.data = s.val.Interface()
	return nil
}

// Validate executes all validators for all slice elements.
func (s *Slice) Validate() error {
	for i := 0; i < s.val.Len(); i++ {
		for fieldName, validator := range s.validators {
			field := s.val.Index(i).Field(s.fieldIndexMap[fieldName])
			if !field.IsValid() {
				return fmt.Errorf("field %s is invalid", fieldName)
			}

			if err := validator(field.Interface()); err != nil {
				return fmt.Errorf("validate field %s failed, err: %w", fieldName, err)
			}
		}
	}
	return nil
}

// HaveField checks if the slice element has given field.
func (s *Slice) HaveField(field string) bool {
	_, ok := s.fieldIndexMap[field]
	return ok
}
