/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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
	"errors"
	"fmt"
	"reflect"

	"configcenter/src/common/blog"
)

// GetUpdateFieldsWithOption get update fields with option, it will return fields with assigned values,
// masking fields that are not allowed to be updated
func GetUpdateFieldsWithOption(data interface{}, opts *FieldOption) (map[string]interface{}, error) {
	if data == nil {
		return nil, errors.New("can not get update fields with option, data is nil")
	}

	if opts == nil {
		return nil, errors.New("can not get update fields with option, opts is nil")
	}

	updateFields, err := getUpdateFields(data)
	if err != nil {
		return nil, err
	}

	toUpdate := make(map[string]interface{})
	for tag, value := range updateFields {
		if opts.NeedIgnored(tag) {
			// this is a field which is need to be ignored,
			// which means do not need to be updated.
			continue
		}

		toUpdate[tag] = value
	}

	return toUpdate, nil
}

// getUpdateFields get update fields from object.
// 1. the input v can only be struct or *struct.
// 2. fields with a bson value of ',inline' must be pointers or struct.
// 3. if the value is not pointers type, it will ignore this field, unless the value of bson tag is ',inline'.
// 4. when the value of bson tag is ',inline', it will pick up the value inside and level with the outer layer.
// 5. except for the field whose bson value of tag is', inline', it only determines whether the outermost
// field needs to be updated.
func getUpdateFields(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return map[string]interface{}{}, nil
	}

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		v = value.Elem().Interface()
		value = reflect.ValueOf(v)
	}

	if value.Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported value type: %v", value.Kind())
	}

	kv := make(map[string]interface{})
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		name := field.Name
		tag := field.Tag.Get("bson")
		if tag == "" {
			return nil, fmt.Errorf("field: %s do not have a 'bson' tag", name)
		}

		fieldVal := value.FieldByName(name)
		if tag == ",inline" {
			val := fieldVal.Interface()
			mapper, err := getUpdateFields(val)
			if err != nil {
				return nil, err
			}

			for k, v := range mapper {
				kv[k] = v
			}
			continue
		}

		if fieldVal.Kind() != reflect.Ptr {
			blog.V(4).Infof("field is not pointer type, type: %v, field: %s, skip", fieldVal.Kind(), name)
			continue
		}

		if fieldVal.IsNil() {
			blog.V(4).Infof("field %s value is nil, skip", name)
			continue
		}

		val := fieldVal.Elem().Interface()
		kv[tag] = val
	}

	return kv, nil
}

// FieldOption is to define which field need to be:
// 1. be ignored, which means not be updated even its value is not nil.
// NOTE:
// 1. The map's key is the structs' 'bson' tag of that field.
type FieldOption struct {
	ignored map[string]struct{}
}

// NewFieldOptions create a blank option instances for add keys
// to be updated when update data.
func NewFieldOptions() *FieldOption {
	return &FieldOption{
		ignored: make(map[string]struct{}),
	}
}

// NeedIgnored check if this field does not need to be updated.
func (f *FieldOption) NeedIgnored(field string) bool {
	_, ok := f.ignored[field]
	return ok
}

// AddIgnoredFields add fields which do not need to be updated even it
// do has a value.
func (f *FieldOption) AddIgnoredFields(fields ...string) *FieldOption {
	for _, one := range fields {
		f.ignored[one] = struct{}{}
	}

	return f
}
