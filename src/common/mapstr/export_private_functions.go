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
	"fmt"
	"reflect"
)

func dealMap(value reflect.Value) (MapStr, error) {
	mapKeys := value.MapKeys()
	mapResult := MapStr{}
	for _, key := range mapKeys {
		keyValue := value.MapIndex(key)
		switch keyValue.Kind() {
		default:
			mapResult.Set(key.String(), keyValue.Interface())
		case reflect.Interface:
			subMapResult, err := convertInterfaceIntoMapStrByReflection(keyValue.Interface())
			if nil != err {
				return nil, err
			}
			mapResult.Set(key.String(), subMapResult)
		case reflect.Struct:
			subMapResult, err := dealStruct(keyValue.Type(), keyValue)
			if nil != err {
				return nil, err
			}
			mapResult.Set(key.String(), subMapResult)
		case reflect.Map:
			subMapResult, err := dealMap(keyValue)
			if nil != err {
				return nil, err
			}
			mapResult.Set(key.String(), subMapResult)
		}
	}

	return mapResult, nil
}

func dealStruct(kind reflect.Type, value reflect.Value) (MapStr, error) {

	mapResult := MapStr{}

	fieldNum := value.NumField()
	for idx := 0; idx < fieldNum; idx++ {

		fieldType := kind.Field(idx)
		fieldValue := value.Field(idx)

		switch fieldValue.Kind() {
		default:
			if fieldValue.CanInterface() {
				mapResult.Set(fieldType.Name, fieldValue.Interface())
			}
		case reflect.Interface:
			subMapResult, err := convertInterfaceIntoMapStrByReflection(fieldValue.Interface())
			if nil != err {
				return nil, err
			}
			mapResult.Set(fieldType.Name, subMapResult)
		case reflect.Struct:
			subMapResult, err := dealStruct(fieldValue.Type(), fieldValue)
			if nil != err {
				return nil, err
			}
			mapResult.Set(fieldType.Name, subMapResult)
		case reflect.Map:
			subMapResult, err := dealMap(fieldValue)
			if nil != err {
				return nil, err
			}
			mapResult.Set(fieldType.Name, subMapResult)
		}
	}

	return mapResult, nil
}
func convertInterfaceIntoMapStrByReflection(target interface{}) (MapStr, error) {

	value := reflect.ValueOf(target)
	switch value.Kind() {
	case reflect.Map:
		return dealMap(value)
	case reflect.Struct:
		return dealStruct(value.Type(), value)
	}

	return nil, fmt.Errorf("no support the kind(%s)", value.Kind())
}
