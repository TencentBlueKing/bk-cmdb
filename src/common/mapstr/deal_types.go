/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mapstr

import (
	"reflect"
	"strings"
)

func dealPointer(value reflect.Value, tag, tagName string) interface{} {

	if value.IsNil() {
		return getZeroValue(value.Type())
	}

	value = value.Elem()

	switch value.Kind() {
	case reflect.Struct:
		if value.IsValid() && value.CanInterface() {
			innerMapStr := SetValueToMapStrByTagsWithTagName(value.Interface(), tagName)
			return MapStr{tag: innerMapStr}
		}
	case reflect.Ptr:
		return dealPointer(value.Elem(), tag, tagName)
	}

	if value.IsValid() && value.CanInterface() {
		return value.Interface()
	}

	return nil
}

func dealMap(value reflect.Value, tagName string) (MapStr, error) {
	mapKeys := value.MapKeys()
	mapResult := MapStr{}
	for _, key := range mapKeys {
		keyValue := value.MapIndex(key)
		switch keyValue.Kind() {
		default:
			if keyValue.IsValid() && keyValue.CanInterface() {
				mapResult.Set(key.String(), keyValue.Interface())
			}
		case reflect.Interface:
			subMapResult, err := convertInterfaceIntoMapStrByReflection(keyValue.Interface(), tagName)
			if nil != err {
				return nil, err
			}
			mapResult.Set(key.String(), subMapResult)
		case reflect.Struct:
			subMapResult, err := dealStruct(keyValue.Type(), keyValue, tagName)
			if nil != err {
				return nil, err
			}
			mapResult.Set(key.String(), subMapResult)
		case reflect.Map:
			subMapResult, err := dealMap(keyValue, tagName)
			if nil != err {
				return nil, err
			}
			mapResult.Set(key.String(), subMapResult)
		}
	}

	return mapResult, nil
}

func dealStruct(kind reflect.Type, value reflect.Value, tagName string) (MapStr, error) {

	mapResult := MapStr{}

	fieldNum := value.NumField()
	for idx := 0; idx < fieldNum; idx++ {

		fieldType := kind.Field(idx)
		fieldValue := value.Field(idx)

		switch fieldValue.Kind() {
		default:
			if fieldValue.IsValid() && fieldValue.CanInterface() {
				mapResult.Set(fieldType.Name, fieldValue.Interface())
			}
		case reflect.Interface:
			subMapResult, err := convertInterfaceIntoMapStrByReflection(fieldValue.Interface(), tagName)
			if nil != err {
				return nil, err
			}
			mapResult.Set(fieldType.Name, subMapResult)
		case reflect.Struct:
			subMapResult, err := dealStruct(fieldValue.Type(), fieldValue, tagName)
			if nil != err {
				return nil, err
			}

			tag, ok := fieldType.Tag.Lookup(tagName)
			if !ok {
				mapResult.Set(fieldType.Name, subMapResult)
				continue
			}
			if 0 == len(tag) || strings.Contains(tag, "ignoretomap") {
				continue
			}

			tags := strings.Split(tag, ",")
			mapResult.Set(tags[0], subMapResult)

		case reflect.Map:
			subMapResult, err := dealMap(fieldValue, tagName)
			if nil != err {
				return nil, err
			}
			mapResult.Set(fieldType.Name, subMapResult)
		}
	}

	return mapResult, nil
}
