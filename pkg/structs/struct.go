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

// Struct defines a dynamic struct.
type Struct struct {
	name string
	// data is the underlying dynamic struct instance pointer.
	data any
	// val is the reflection value of the struct instance which is used for runtime manipulation.
	val reflect.Value

	// validators maps struct field names to their corresponding validation functions.
	validators map[string]func(any) error
	// fieldIndexMap maps struct field names to their respective field indices in the struct type.
	fieldIndexMap map[string]int
}

// Pointer returns the underlying dynamic struct instance in the form of pointer.
func (s *Struct) Pointer() any {
	return s.data
}

// Value returns the underlying dynamic struct instance in the form of struct type.
func (s *Struct) Value() any {
	return s.val.Interface()
}

// Get returns the value of the specified struct field.
func (s *Struct) Get(fieldName string) (Valuer, error) {
	fieldIndex, exists := s.fieldIndexMap[fieldName]
	if !exists {
		return nil, fmt.Errorf("field %s not found", fieldName)
	}

	field := s.val.Field(fieldIndex)
	if !field.IsValid() {
		return nil, fmt.Errorf("field %s is invalid", fieldName)
	}

	if field.Kind() == reflect.Struct && fieldName != field.Type().Name() &&
		s.val.Type().Field(fieldIndex).Anonymous {

		v := field.FieldByName(fieldName)
		if v.IsValid() {
			return NewValue(v), nil
		}

		return nil, fmt.Errorf("field %s is invalid inside anonymous struct", fieldName)
	}
	return NewValue(field), nil
}

// Set sets the value of the specified struct field.
func (s *Struct) Set(fieldName string, value any) error {
	fieldIndex, exists := s.fieldIndexMap[fieldName]
	if !exists {
		return fmt.Errorf("field %s not found", fieldName)
	}

	field := s.val.Field(fieldIndex)
	if !field.IsValid() {
		return fmt.Errorf("field %s is invalid", fieldName)
	}
	if !field.CanSet() {
		return fmt.Errorf("field %s cannot be set", fieldName)
	}

	valueType := reflect.TypeOf(value)

	// directly set the field value if the types is the same.
	if valueType == field.Type() {
		field.Set(reflect.ValueOf(value))
		return nil
	}

	if field.Kind() == reflect.Struct && fieldName != field.Type().Name() &&
		s.val.Type().Field(fieldIndex).Anonymous {

		v := field.FieldByName(fieldName)
		if !v.IsValid() {
			return fmt.Errorf("field %s is invalid inside anonymous struct", fieldName)
		}
		// replace to the nested field.
		valueType = v.Type()
		field = v
	}

	// check if the value's type can be converted to the field's type.
	if !valueType.ConvertibleTo(field.Type()) {
		return fmt.Errorf("cannot set field %s, type mismatch", fieldName)
	}

	// convert the value to the field's type and set the struct field value.
	field.Set(reflect.ValueOf(value).Convert(field.Type()))

	return nil
}

// Validate executes all validators for their corresponding struct field values.
func (s *Struct) Validate() error {
	for fieldName, validator := range s.validators {
		field := s.val.Field(s.fieldIndexMap[fieldName])
		if !field.IsValid() {
			return fmt.Errorf("field %s is invalid", fieldName)
		}

		if err := validator(field.Interface()); err != nil {
			return fmt.Errorf("validate field %s failed, err: %w", fieldName, err)
		}
	}
	return nil
}

// Name get the name of the struct.
func (s *Struct) Name() string {
	return s.name
}

// HaveField checks if the struct has given field.
func (s *Struct) HaveField(field string) bool {
	_, ok := s.fieldIndexMap[field]
	return ok
}
