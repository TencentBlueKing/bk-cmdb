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
	case int32:
		id = int64(a.(int32))
	case int64:
		id = int64(a.(int64))
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

func GetMapInterfaceByInerface(data interface{}) ([]interface{}, error) {
	var values []interface{}
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
	var ret []string
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
	rflval := reflect.ValueOf(value)
	for rflval.CanAddr() {
		rflval = rflval.Elem()
	}
	if rflval.Kind() != reflect.Slice {
		return []interface{}{value}
	}

	result := []interface{}{}
	for i := 0; i < rflval.Len(); i++ {
		if rflval.Index(i).CanInterface() {
			result = append(result, rflval.Index(i).Interface())
		}
	}

	return result
}
