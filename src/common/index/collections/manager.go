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

package collections

import (
	"fmt"

	"configcenter/src/storage/dal/types"
)

var tableNameIndexes = make(map[string][]types.Index, 0)

// registerIndexes 注册索引会有合法性检查，名字重复或索引使用的key重复会出现panic
func registerIndexes(tableName string, indexes []types.Index) {

	for _, newIdx := range indexes {
		for _, idx := range tableNameIndexes[tableName] {
			if idx.Name == newIdx.Name {
				panic(fmt.Sprintf("table(%s). index(%s) with the same name already exists", tableName, idx.Name))
			}
			if indexKeyEqual(newIdx, idx) {
				panic(fmt.Sprintf("table(%s).  index with the same keys. index: (%s, %s) ",
					tableName, idx.Name, newIdx.Name))
			}
		}
		tableNameIndexes[tableName] = append(tableNameIndexes[tableName], newIdx)

	}

}

// indexKeyEqual 索引对比， 索引名字不参与对比
func indexKeyEqual(toDBIndex, dbIndex types.Index) bool {

	if len(toDBIndex.Keys) != len(dbIndex.Keys) {
		return false
	}

	toDBIndexMap := toDBIndex.Keys.Map()
	dbIndexMap := dbIndex.Keys.Map()
	for key, val := range toDBIndexMap {
		dbVal, exists := dbIndexMap[key]
		if !exists {
			return false
		}
		if val != dbVal {
			return false
		}
	}

	return true

}

// Indexes 返回表数据库表中定义索引
func Indexes() map[string][]types.Index {
	return tableNameIndexes
}

// DeprecatedIndexName TODO
//  DeprecatedIndexName 未规范索引的名字。用来对需要删除的未规范索引
func DeprecatedIndexName() map[string][]string {
	return deprecatedIndexName
}
