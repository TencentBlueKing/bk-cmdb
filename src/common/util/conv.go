/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetStrByInterface(a interface{}) string {
	if nil == a {
		return ""
	}
	return fmt.Sprintf("%v", a)
}

func GetIntByInterface(a interface{}) (int, error) {
	id := 0
	var err error
	switch val := a.(type) {
	case int:
		id = val
	case int32:
		id = int(val)
	case int64:
		id = int(val)
	case json.Number:
		var tmpID int64
		tmpID, err = val.Int64()
		id = int(tmpID)
	case float64:
		id = int(val)
	case float32:
		id = int(val)
	case string:
		var tmpID int64
		tmpID, err = strconv.ParseInt(a.(string), 10, 64)
		id = int(tmpID)
	default:
		err = errors.New("not numeric")

	}
	return id, err
}

func GetInt64ByInterface(a interface{}) (int64, error) {
	var id int64 = 0
	var err error
	switch a.(type) {
	case int:
		id = int64(a.(int))
	case uint:
		id = int64(a.(uint))
	case int32:
		id = int64(a.(int32))
	case uint32:
		id = int64(a.(uint32))
	case int64:
		id = int64(a.(int64))
	case uint64:
		id = int64(a.(uint64))
	case json.Number:
		var tmpID int64
		tmpID, err = a.(json.Number).Int64()
		id = int64(tmpID)
	case float64:
		tmpID := a.(float64)
		id = int64(tmpID)
	case float32:
		tmpID := a.(float32)
		id = int64(tmpID)
	case string:
		id, err = strconv.ParseInt(a.(string), 10, 64)
	default:
		err = errors.New("not numeric")

	}
	return id, err
}

func GetFloat64ByInterface(a interface{}) (float64, error) {
	switch i := a.(type) {
	case int:
		return float64(i), nil
	case int8:
		return float64(i), nil
	case int16:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case uint:
		return float64(i), nil
	case uint8:
		return float64(i), nil
	case uint16:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case string:
		return strconv.ParseFloat(i, 64)
	case json.Number:
		return i.Float64()
	default:
		return 0, errors.New("not numeric")
	}
}

func GetTypeSensitiveUInt64(v interface{}) (uint64, bool) {
	switch tv := v.(type) {
	case int8:
		return uint64(tv), true
	case int16:
		return uint64(tv), true
	case int32:
		return uint64(tv), true
	case int64:
		return uint64(tv), true
	case int:
		return uint64(tv), true
	case uint8:
		return uint64(tv), true
	case uint16:
		return uint64(tv), true
	case uint32:
		return uint64(tv), true
	case uint64:
		return uint64(tv), true
	case uint:
		return uint64(tv), true
	case float32:
		return uint64(tv), true
	case float64:
		return uint64(tv), true
	default:
		return 0, false
	}
}

func GetMapInterfaceByInerface(data interface{}) ([]interface{}, error) {
	values := make([]interface{}, 0)
	switch data.(type) {
	case []int:
		vs, _ := data.([]int)
		for _, v := range vs {
			values = append(values, v)
		}
	case []int32:
		vs, _ := data.([]int32)
		for _, v := range vs {
			values = append(values, v)
		}
	case []int64:
		vs, _ := data.([]int64)
		for _, v := range vs {
			values = append(values, v)
		}
	case []string:
		vs, _ := data.([]string)
		for _, v := range vs {
			values = append(values, v)
		}
	case []interface{}:
		values = data.([]interface{})
	default:
		return nil, errors.New("params value can not be empty")
	}

	return values, nil
}

// SliceStrToInt: 将字符串切片转换为整型切片
func SliceStrToInt(sliceStr []string) ([]int, error) {
	sliceInt := make([]int, 0)
	for _, idStr := range sliceStr {

		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return []int{}, err
		}
		sliceInt = append(sliceInt, id)
	}
	return sliceInt, nil
}

// SliceStrToInt64 将字符串切片转换为整型切片
func SliceStrToInt64(sliceStr []string) ([]int64, error) {
	sliceInt := make([]int64, 0)
	for _, idStr := range sliceStr {

		if idStr == "" {
			continue
		}

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return []int64{}, err
		}
		sliceInt = append(sliceInt, id)
	}
	return sliceInt, nil
}

// GetStrValsFromArrMapInterfaceByKey get []string from []map[string]interface{}, Do not consider errors
func GetStrValsFromArrMapInterfaceByKey(arrI []interface{}, key string) []string {
	ret := make([]string, 0)
	for _, row := range arrI {
		mapRow, ok := row.(map[string]interface{})
		if ok {
			val, ok := mapRow[key].(string)
			if ok {
				ret = append(ret, val)
			}
		}
	}

	return ret
}

func ConverToInterfaceSlice(value interface{}) []interface{} {
	rflVal := reflect.ValueOf(value)
	for rflVal.CanAddr() {
		rflVal = rflVal.Elem()
	}
	if rflVal.Kind() != reflect.Slice {
		return []interface{}{value}
	}

	result := make([]interface{}, 0)
	for i := 0; i < rflVal.Len(); i++ {
		if rflVal.Index(i).CanInterface() {
			result = append(result, rflVal.Index(i).Interface())
		}
	}

	return result
}

// SplitStrField    split string field, remove empty string
func SplitStrField(str, sep string) []string {
	if "" == str {
		return nil
	}
	return strings.Split(str, sep)
}

func DecodeIntoTypeValue(fieldType reflect.Type, value interface{}) (reflect.Value, error) {
	var retVal reflect.Value
	switch fieldType.Kind() {
	case reflect.Bool:
		bl, ok := value.(bool)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(bl)
	case reflect.Int:

		val, err := GetInt64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(int(val))

	case reflect.Int8:
		val, err := GetInt64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(int8(val))
	case reflect.Int16:
		val, err := GetInt64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(int16(val))
	case reflect.Int32:
		val, err := GetInt64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(int32(val))
	case reflect.Int64:
		val, err := GetInt64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(int64(val))

	case reflect.Uint:
		val, ok := GetTypeSensitiveUInt64(value)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(uint(val))
	case reflect.Uint8:
		val, ok := GetTypeSensitiveUInt64(value)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(uint8(val))
	case reflect.Uint16:
		val, ok := GetTypeSensitiveUInt64(value)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(uint16(val))
	case reflect.Uint32:
		val, ok := GetTypeSensitiveUInt64(value)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(uint32(val))
	case reflect.Uint64:
		val, ok := GetTypeSensitiveUInt64(value)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(uint64(val))
	case reflect.Float32:
		val, err := GetFloat64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(float32(val))
	case reflect.Float64:
		val, err := GetFloat64ByInterface(value)
		if err != nil {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		retVal = reflect.ValueOf(float64(val))
	// case reflect.Complex64:
	// case reflect.Complex128:
	case reflect.Array:
		valArr, ok := value.([]interface{})
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		setValArr := reflect.MakeSlice(fieldType.Elem(), 0, 10)
		for _, valItem := range valArr {
			tmpVal, err := DecodeIntoTypeValue(fieldType.Elem(), valItem)
			if err != nil {
				return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
			}
			retVal.Set(reflect.Append(setValArr, tmpVal))
		}
	//case reflect.Chan:
	//case reflect.Func:
	case reflect.Interface:
		retVal = reflect.ValueOf(value)
	case reflect.Map:
		valMap, ok := value.(map[string]interface{})
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", reflect.TypeOf(value).Kind().String(), fieldType.Kind().String())
		}
		fmt.Println("\n \n test code", fieldType, fieldType.Elem().Kind().String(), fieldType.Elem().Elem().Kind().String())
		///setValMap := reflect.MakeMap(fieldType)
		for key, val := range valMap {
			fmt.Println(key, val)
		}
	case reflect.Slice:
		valArr, ok := value.([]interface{})
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", fieldType.Kind().String(), reflect.TypeOf(value).Kind().String())
		}
		setValArr := reflect.MakeSlice(fieldType, 0, len(valArr))
		for _, valItem := range valArr {
			tmpVal, err := DecodeIntoTypeValue(fieldType.Elem(), valItem)
			if err != nil {
				return retVal, fmt.Errorf("%s into Go value of type %s", fieldType.Kind().String(), reflect.TypeOf(value).Kind().String())
			}
			setValArr = reflect.Append(setValArr, tmpVal)
		}

		retVal = setValArr

	case reflect.String:
		val, ok := value.(string)
		if !ok {
			return retVal, fmt.Errorf("%s into Go value of type %s", fieldType.Kind().String(), reflect.TypeOf(value).Kind().String())
		}
		retVal = reflect.ValueOf(val)
	default:
		return retVal, fmt.Errorf("%s into Go value of type %s", fieldType.Kind().String(), reflect.TypeOf(value).Kind().String())
	}

	return retVal, nil
}
