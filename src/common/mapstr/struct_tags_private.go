/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
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
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

func getZeroFieldValue(valueType reflect.Type) interface{} {

	switch valueType.Kind() {
	case reflect.Ptr:
		return getZeroFieldValue(valueType.Elem())
	case reflect.String:
		return ""
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return 0
	}

	return nil
}

func dealPointer(value reflect.Value, tag, tagName string) interface{} {

	if value.IsNil() {
		return getZeroFieldValue(value.Type())
	}

	value = value.Elem()

	switch value.Kind() {
	case reflect.Struct:
		if value.CanInterface() {
			innerMapStr := SetValueToMapStrByTagsWithTagName(value.Interface(), tagName)
			return MapStr{tag: innerMapStr}
		}
	case reflect.Ptr:
		return dealPointer(value.Elem(), tag, tagName)
	}

	if value.CanInterface() {
		return value.Interface()
	}

	return nil
}
func getValue(inputVal interface{}) interface{} {
	switch t := inputVal.(type) {
	case string, int, int8, int16, int32, int64, uint, uint16, uint32, uint64, uint8:
		return t
	case json.Number:
		if val, err := t.Int64(); nil == err {
			return val
		}
		if val, err := t.Float64(); nil == err {
			return val
		}
		return t.String()
	default:
		return t
	}
}
func setMapToReflectValue(returnVal, inputVal reflect.Value) reflect.Value {

	mapKeys := inputVal.MapKeys()
	for _, key := range mapKeys {

		value := inputVal.MapIndex(key)
		if !value.CanInterface() {
			continue
		}

		switch returnVal.Type().Elem().Kind() {
		default:
			panic("not support:" + returnVal.Type().Elem().Kind().String() + fmt.Sprintf(" value: %v", value.Interface()))
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
			returnVal.Set(reflect.ValueOf(map[string]uint{key.String(): uint(toInt(value.Interface()))}))
		case reflect.Uint16:
			returnVal.Set(reflect.ValueOf(map[string]uint16{key.String(): uint16(toInt(value.Interface()))}))
		case reflect.Uint32:
			returnVal.Set(reflect.ValueOf(map[string]uint32{key.String(): uint32(toInt(value.Interface()))}))
		case reflect.Uint64:
			returnVal.Set(reflect.ValueOf(map[string]uint64{key.String(): uint64(toInt(value.Interface()))}))
		case reflect.Uint8:
			returnVal.Set(reflect.ValueOf(map[string]uint8{key.String(): uint8(toInt(value.Interface()))}))
		case reflect.Float32:
			returnVal.Set(reflect.ValueOf(map[string]float32{key.String(): float32(toFloat(value.Interface()))}))
		case reflect.Float64:
			returnVal.Set(reflect.ValueOf(map[string]float64{key.String(): toFloat(value.Interface())}))
		}
	}

	return returnVal
}

func convertToInt(fieldName string, tagVal interface{}, fieldValue *reflect.Value) error {
	switch t := tagVal.(type) {
	default:
		return fmt.Errorf("unsuport the type %s tagVal %v", fieldName, reflect.TypeOf(tagVal).Kind())
	case float32:
		fieldValue.SetInt(int64(t))
	case float64:
		fieldValue.SetInt(int64(t))
	case int:
		fieldValue.SetInt(int64(t))
	case int16:
		fieldValue.SetInt(int64(t))
	case int32:
		fieldValue.SetInt(int64(t))
	case int64:
		fieldValue.SetInt(int64(t))
	case int8:
		fieldValue.SetInt(int64(t))
	case uint:
		fieldValue.SetInt(int64(t))
	case uint16:
		fieldValue.SetInt(int64(t))
	case uint32:
		fieldValue.SetInt(int64(t))
	case uint64:
		fieldValue.SetInt(int64(t))
	case uint8:
		fieldValue.SetInt(int64(t))
	}
	return nil
}
func getTypeElem(targetType reflect.Type) reflect.Type {
	switch targetType.Kind() {
	case reflect.Ptr:
		return getTypeElem(targetType.Elem())
	}
	return targetType
}
func getValueElem(targetValue reflect.Value) reflect.Value {
	switch targetValue.Kind() {
	case reflect.Ptr:
		return getValueElem(targetValue.Elem())
	}
	return targetValue
}

func parseStruct(targetType reflect.Type, targetValue reflect.Value, values MapStr, tagName string) error {

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
			setMapToReflectValue(fieldValue, reflect.ValueOf(tagVal))

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
			if err := parseStruct(structField.Type, targetResult, valMapStr, tagName); nil != err {
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
				if err := parseStruct(structField.Type, targetResult, valMapStr, tagName); nil != err {
					return err
				}
				fieldValue.Set(targetResult)
			case string:
				targetResult = getValueElem(targetResult)
				targetResult.SetString(t)
				fieldValue.Set(targetResult.Addr())
			}

		case reflect.Bool:
			fieldValue.SetBool(tagVal.(bool))

		case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
			if err := convertToInt(structField.Name, tagVal, &fieldValue); nil != err {
				return err
			}

		case reflect.Float32, reflect.Float64:
			if err := convertToInt(structField.Name, tagVal, &fieldValue); nil != err {
				return err
			}

		case reflect.String:
			switch t := tagVal.(type) {
			case string:
				fieldValue.SetString(t)
			}

		}

	}

	return nil
}
