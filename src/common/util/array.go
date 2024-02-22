// Package util TODO
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
	"fmt"
	"reflect"
	"strings"
)

// InArray TODO
func InArray(obj interface{}, target interface{}) bool {
	if target == nil {
		return false
	}

	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// ArrayUnique TODO
func ArrayUnique(a interface{}) (ret []interface{}) {
	ret = make([]interface{}, 0)
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		v := va.Index(i).Interface()
		if !InArray(v, ret) {
			ret = append(ret, v)
		}
	}
	return ret
}

// StrArrayUnique get unique string array
func StrArrayUnique(a []string) (ret []string) {
	ret = make([]string, 0)
	length := len(a)
	for i := 0; i < length; i++ {
		if !Contains(ret, a[i]) {
			ret = append(ret, a[i])
		}
	}
	return ret
}

// IntArrayUnique get unique int array
func IntArrayUnique(a []int64) (ret []int64) {
	unique := make(map[int64]struct{})
	ret = make([]int64, 0)
	for _, val := range a {
		if _, exists := unique[val]; exists {
			continue
		}
		unique[val] = struct{}{}
		ret = append(ret, val)
	}

	return ret
}

// BoolArrayUnique TODO
func BoolArrayUnique(a []bool) (ret []bool) {
	ret = make([]bool, 0)
	trueExist := false
	falseExist := false
	for _, item := range a {
		if item == true {
			trueExist = true
		}
		if item == false {
			falseExist = true
		}
	}
	if trueExist {
		ret = append(ret, true)
	}
	if falseExist {
		ret = append(ret, false)
	}
	return ret
}

// RemoveDuplicatesAndEmpty TODO
func RemoveDuplicatesAndEmpty(slice []string) (ret []string) {
	ret = make([]string, 0)
	for _, a := range slice {
		if strings.TrimSpace(a) != "" && !Contains(ret, a) {
			ret = append(ret, a)
		}
	}
	return
}

// StrArrDiff TODO
func StrArrDiff(slice1 []string, slice2 []string) []string {
	diffStr := make([]string, 0)
	for _, i := range slice1 {
		isIn := false
		for _, j := range slice2 {
			if i == j {
				isIn = true
				break
			}
		}
		if !isIn {
			diffStr = append(diffStr, i)
		}
	}
	return diffStr
}

// IntArrDiff 返回slice1与slice2的差集,存在于slice1，不存在于slice2，使用该方法要注意参数的传入顺序
func IntArrDiff(slice1 []int64, slice2 []int64) []int64 {
	diffInt := make([]int64, 0)

	intMap2 := make(map[int64]struct{})
	for _, num2 := range slice2 {
		intMap2[num2] = struct{}{}
	}

	for _, num1 := range slice1 {
		if _, found := intMap2[num1]; !found {
			diffInt = append(diffInt, num1)
		}
	}
	return diffInt
}

// IntArrIntersection TODO
func IntArrIntersection(slice1 []int64, slice2 []int64) []int64 {
	intersectInt := make([]int64, 0)
	intMap := make(map[int64]bool)
	for _, i := range slice1 {
		intMap[i] = true
	}
	for _, j := range slice2 {
		if _, ok := intMap[j]; ok == true {
			intersectInt = append(intersectInt, j)
		}
	}
	return intersectInt
}

// PrettyIPStr TODO
func PrettyIPStr(ips []string) string {
	if len(ips) > 2 {
		return fmt.Sprintf("%s ...", strings.Join(ips[:2], ","))
	}
	return strings.Join(ips, ",")
}

// ReverseArrayString reverse the slice's element from tail to head.
func ReverseArrayString(t []string) []string {
	if len(t) == 0 {
		return t
	}
	for i, j := 0, len(t)-1; i < j; i, j = i+1, j-1 {
		t[i], t[j] = t[j], t[i]
	}
	return t
}

// RemoveDuplicatesAndEmptyByMap remove duplicate element and empty element by map
func RemoveDuplicatesAndEmptyByMap(target []string) []string {
	result := make([]string, 0)
	tempMap := map[string]struct{}{}
	for _, item := range target {
		if item == "" {
			continue
		}

		if _, exist := tempMap[item]; exist {
			continue
		}
		tempMap[item] = struct{}{}
		result = append(result, item)
	}

	return result
}

// IntArrComplementary calculates the complement of subset relative to target
func IntArrComplementary(target []int64, subset []int64) []int64 {
	complementaryInt := make([]int64, 0)
	intMap := make(map[int64]struct{})
	for _, i := range subset {
		intMap[i] = struct{}{}
	}
	for _, j := range target {
		if _, exist := intMap[j]; !exist {
			complementaryInt = append(complementaryInt, j)
		}
	}
	return complementaryInt
}

// IntArrDeleteElements  the same elements in target and sub are deleted from target.
func IntArrDeleteElements(target, sub []int64) []int64 {
	if len(sub) == 0 {
		return target
	}

	templateMap := make(map[int64]struct{})
	for _, id := range target {
		templateMap[id] = struct{}{}
	}
	for _, id := range sub {
		if _, ok := templateMap[id]; ok {
			delete(templateMap, id)
		}
	}
	result := make([]int64, 0)
	for id := range templateMap {
		result = append(result, id)
	}
	return result
}
