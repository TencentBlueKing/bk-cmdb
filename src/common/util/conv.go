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
	case int8:
		return int64(a.(int8)), nil
	case int16:
		return int64(a.(int16)), nil
	case int32:
		id = int64(a.(int32))
	case int64:
		id = int64(a.(int64))
	case uint:
		id = int64(a.(uint))
	case uint8:
		return int64(a.(uint8)), nil
	case uint16:
		return int64(a.(uint16)), nil
	case uint32:
		id = int64(a.(uint32))
	case uint64:
		id = int64(a.(uint64))
	case json.Number:
		id, err = a.(json.Number).Int64()
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

// SliceInterfaceToInt64 将interface切片转化为int64切片,且interface的真实类型可以是任何整数类型.
// 失败则返回nil,error.
func SliceInterfaceToInt64(faceSlice []interface{}) ([]int64,error){
	// 预分配空间.
	var results = make([]int64,len(faceSlice))

	// 转化操作.
	for i,item := range faceSlice{
		switch val := item.(type) {
		case int64:
			results[i] = val
		case int:
			results[i] = int64(val)
		case int8:
			results[i] = int64(val)
		case int16:
			results[i] = int64(val)
		case int32:
			results[i] = int64(val)
		case uint:
			results[i] = int64(val)
		case uint8:
			results[i] = int64(val)
		case uint16:
			results[i] = int64(val)
		case uint32:
			results[i] = int64(val)
		case uint64:
			results[i] = int64(val)
		default:
			return nil,errors.New("can't convert to int64")
		}
	}
	return results,nil
}


// SliceInterfaceToBool将interface切片转化为string切片,且interface的真实类型必须是string.
// 失败则返回nil,error.
func SliceInterfaceToString(faceSlice []interface{}) ([]string,error){
	// 预分配空间.
	var results = make([]string,len(faceSlice))

	// 转化操作.
	for i,item := range faceSlice{
		var ok bool

		//如果转化失败则返回错误.
		if results[i],ok = item.(string) ; !ok {
			return nil,errors.New("can't convert to string")
		}

	}
	return results,nil
}

// SliceInterfaceToBool将interface切片转化为bool切片,且interface的真实类型必须是bool.
// 失败则返回nil,error.
func SliceInterfaceToBool(faceSlice []interface{}) ([]bool,error){
	// 预分配空间.
	var results = make([]bool,len(faceSlice))

	// 转化操作.
	for i,item := range faceSlice{
		 var ok bool

		 //如果转化失败则返回错误.
		 if results[i],ok = item.(bool) ; !ok {
			 return nil,errors.New("can't convert to bool")
		 }

	}
	return results,nil
}
