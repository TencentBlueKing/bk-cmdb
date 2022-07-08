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

package index

import (
	"reflect"

	"configcenter/src/storage/dal/types"
)

// FindIndexByIndexFields 根据索引中用到的字段找db中对应的索引，
// 注意： 索引字段的排序方式，是否background，unique，partialFilterExpression 有存在不同的情况。
//       由于mongodb 不允许同一组字段加多个索引
func FindIndexByIndexFields(keys map[string]int32, indexList []types.Index) (dbIndex types.Index, exists bool) {
	for _, idx := range indexList {
		if len(idx.Keys) != len(keys) {
			continue
		}
		exists = true
		for key := range idx.Keys {
			if _, keyExists := keys[key]; !keyExists {
				exists = false
				break
			}
		}
		if exists {
			return idx, exists
		}

	}

	return types.Index{}, false
}

// IndexEqual 索引对比， 索引名字不参与对比
func IndexEqual(toDBIndex, dbIndex types.Index) bool {
	if toDBIndex.Background != dbIndex.Background {
		return false
	}
	if toDBIndex.Unique != dbIndex.Unique {
		return false
	}
	if toDBIndex.ExpireAfterSeconds != dbIndex.ExpireAfterSeconds {
		return false
	}
	if len(toDBIndex.Keys) != len(dbIndex.Keys) {
		return false
	}

	if len(toDBIndex.PartialFilterExpression) != len(dbIndex.PartialFilterExpression) {
		return false
	}

	for key, val := range toDBIndex.Keys {
		dbVal, exists := dbIndex.Keys[key]
		if !exists {
			return false
		}
		if val != dbVal {
			return false
		}
	}

	// NOTICE: 对比逻辑不严谨, 如果是cc 代码产生的唯一索引，不存在问题
	for key, val := range toDBIndex.PartialFilterExpression {
		dbVal, exists := dbIndex.PartialFilterExpression[key]
		if !exists {
			return false
		}
		if !reflect.DeepEqual(val, dbVal) {
			return false
		}
	}

	return true

}
