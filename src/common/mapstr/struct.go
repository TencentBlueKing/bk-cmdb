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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"configcenter/src/common/blog"
)

func setMapStrByStruct(targetType reflect.Type, targetValue reflect.Value, values MapStr, tagName string) error {

	switch targetType.Kind() {
	case reflect.Ptr:
		targetType = targetType.Elem()
		targetValue = targetValue.Elem()

		if targetType.Kind() == reflect.Ptr {
			return setMapStrByStruct(targetType, targetValue, values, tagName)
		}

	}

	numField := targetType.NumField()
	for i := 0; i < numField; i++ {
		structField := targetType.Field(i)
		tag, ok := structField.Tag.Lookup(tagName)
		if !ok && !structField.Anonymous {
			continue
		}

		if (0 == len(tag) || strings.Contains(tag, "ignoretomap")) && !structField.Anonymous {
			continue
		}
		tags := strings.Split(tag, ",")
		if 0 == len(tag) {
			tags = []string{structField.Name}
		}

		fieldValue := targetValue.FieldByName(structField.Name)
		if fieldValue.IsValid() && !fieldValue.CanInterface() {
			continue
		}

		if isEmptyValue(fieldValue) && strings.Contains(tag, "omitempty") {
			continue
		}

		switch structField.Type.Kind() {
		case reflect.String,
			reflect.Float32, reflect.Float64,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Complex64, reflect.Complex128,
			reflect.Array,
			reflect.Interface,
			reflect.Map,
			reflect.Slice,
			reflect.Bool:
			values.Set(tags[0], fieldValue.Interface())
		case reflect.Struct:
			innerMapStr := SetValueToMapStrByTagsWithTagName(fieldValue.Interface(), tagName)
			values.Set(tags[0], innerMapStr)

		case reflect.Ptr:

			innerValue := dealPointer(fieldValue, tags[0], tagName)
			values.Set(tags[0], innerValue)
		default:
			blog.Infof("[mapstr] invalide kind: %v for field %v", structField.Type.Kind(), tags[0])
		}

	}
	return nil
}

func setStructByMapStr(targetType reflect.Type, targetValue reflect.Value, values MapStr, tagName string) error {

	targetType = getTypeElem(targetType)
	targetValue = getValueElem(targetValue)

	numField := targetType.NumField()
	for i := 0; i < numField; i++ {
		structField := targetType.Field(i)
		tag, ok := structField.Tag.Lookup(tagName)
		if !ok {
			continue
		}

		if 0 == len(tag) || strings.Contains(tag, "ignoretostruct") {
			continue
		}

		tags := strings.Split(tag, ",")

		tagVal, ok := values[tags[0]]
		if !ok {
			continue
		}

		if nil == tagVal {
			continue
		}

		fieldValue := targetValue.FieldByName(structField.Name)
		if !fieldValue.CanSet() {
			continue
		}

		switch structField.Type.Kind() {
		default:
			return fmt.Errorf("unsupport the type %s %v", structField.Name, structField.Type.Kind())

		case reflect.Map:
			if _, err := setMapToReflectValue(fieldValue, reflect.ValueOf(tagVal)); nil != err {
				return err
			}

		case reflect.Interface:
			tmpVal := reflect.ValueOf(tagVal)
			switch tmpVal.Kind() {
			case reflect.Ptr:
				fieldValue.Set(tmpVal.Elem())
			default:
				fieldValue.Set(tmpVal)
			}

		case reflect.Struct:
			valMapStr, err := NewFromInterface(tagVal)
			if nil != err {
				return err
			}
			targetResult := reflect.New(structField.Type)
			if err := setStructByMapStr(structField.Type, targetResult, valMapStr, tagName); nil != err {
				return err
			}
			fieldValue.Set(targetResult.Elem())

		case reflect.Ptr:

			targetResult := reflect.New(structField.Type.Elem())
			switch t := tagVal.(type) {
			default:
				valMapStr, err := NewFromInterface(tagVal)
				if nil != err {
					return err
				}
				if err := setStructByMapStr(structField.Type, targetResult, valMapStr, tagName); nil != err {
					return err
				}
				fieldValue.Set(targetResult)
			case string:
				targetResult = getValueElem(targetResult)
				targetResult.SetString(t)
				fieldValue.Set(targetResult.Addr())
			}

		case reflect.Bool:
			fieldValue.SetBool(toBool(tagVal))
		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
			fieldValue.SetInt(int64(toInt(tagVal)))
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			fieldValue.SetUint(uint64(toUint(tagVal)))
		case reflect.Float32, reflect.Float64:
			fieldValue.SetFloat(toFloat(tagVal))
		case reflect.String:
			switch t := tagVal.(type) {
			case string:
				fieldValue.SetString(t)
			}

		}

	}

	return nil
}

func setMapToReflectValue(returnVal, inputVal reflect.Value) (reflect.Value, error) {

	mapKeys := inputVal.MapKeys()
	for _, key := range mapKeys {

		value := inputVal.MapIndex(key)
		if !returnVal.CanSet() || !value.IsValid() || !value.CanInterface() {
			continue
		}

		switch returnVal.Type().Elem().Kind() {
		default:
			return returnVal, errors.New("not support:" + returnVal.Type().Elem().Kind().String() + fmt.Sprintf(" value: %v", value.Interface()))
		case reflect.Interface:
			returnVal.Set(reflect.ValueOf(map[string]interface{}{key.String(): value.Interface()}))
		case reflect.String:
			returnVal.Set(reflect.ValueOf(map[string]string{key.String(): fmt.Sprintf("%v", value.Interface())}))
		case reflect.Int:
			returnVal.Set(reflect.ValueOf(map[string]int{key.String(): toInt(value.Interface())}))
		case reflect.Int8:
			returnVal.Set(reflect.ValueOf(map[string]int8{key.String(): int8(toInt(value.Interface()))}))
		case reflect.Int16:
			returnVal.Set(reflect.ValueOf(map[string]int16{key.String(): int16(toInt(value.Interface()))}))
		case reflect.Int32:
			returnVal.Set(reflect.ValueOf(map[string]int32{key.String(): int32(toInt(value.Interface()))}))
		case reflect.Int64:
			returnVal.Set(reflect.ValueOf(map[string]int64{key.String(): int64(toInt(value.Interface()))}))
		case reflect.Uint:
			returnVal.Set(reflect.ValueOf(map[string]uint{key.String(): uint(toUint(value.Interface()))}))
		case reflect.Uint16:
			returnVal.Set(reflect.ValueOf(map[string]uint16{key.String(): uint16(toUint(value.Interface()))}))
		case reflect.Uint32:
			returnVal.Set(reflect.ValueOf(map[string]uint32{key.String(): uint32(toUint(value.Interface()))}))
		case reflect.Uint64:
			returnVal.Set(reflect.ValueOf(map[string]uint64{key.String(): uint64(toUint(value.Interface()))}))
		case reflect.Uint8:
			returnVal.Set(reflect.ValueOf(map[string]uint8{key.String(): uint8(toUint(value.Interface()))}))
		case reflect.Float32:
			returnVal.Set(reflect.ValueOf(map[string]float32{key.String(): float32(toFloat(value.Interface()))}))
		case reflect.Float64:
			returnVal.Set(reflect.ValueOf(map[string]float64{key.String(): toFloat(value.Interface())}))
		}
	}

	return returnVal, nil
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}
