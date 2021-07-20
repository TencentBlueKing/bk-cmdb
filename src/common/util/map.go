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
	"sort"
)

func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// SortedMapIntKeys get sorted int keys slice from map[int]map[string]interface{}
func SortedMapIntKeys(data map[int]map[string]interface{}) []int {
	keys := make([]int, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Sort(IntSlice(keys))
	return keys
}

// SortedMapInt64Keys get sorted int64 keys slice from map[int64]map[string]interface{}
func SortedMapInt64Keys(data map[int64]map[string]interface{}) []int64 {
	keys := make([]int64, 0)
	for k := range data {
		keys = append(keys, k)
	}
	sort.Sort(Int64Slice(keys))
	return keys
}
