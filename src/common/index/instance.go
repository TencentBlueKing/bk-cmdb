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
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"
)

/*
通用模型实例表中的索引。新建模型的时候使用
*/

var (
	instanceDefaultIndexes = []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "bkObjId",
			Keys: map[string]int32{
				"bk_obj_id": 1,
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkSupplierAccount",
			Keys: map[string]int32{
				"bk_supplier_account": 1,
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkInstId",
			Keys: map[string]int32{
				"bk_inst_id": 1,
			},
			Background: true,
			// 新加 2021年03月11日
			Unique: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkInstName",
			Keys: map[string]int32{
				"bk_inst_name": 1,
			},
			Background: false,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkParentID",
			Keys: map[string]int32{
				"bk_parent_id": 1,
			},
			Background: false,
		},

		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkParentID_bkInstName",
			Keys: map[string]int32{
				"bk_parent_id": 1,
				"bk_inst_name": 1,
			},
			Background: false,
			Unique:     true,
		},
	}
)

// 创建唯一索引的时候可以被忽略的索引
func IngoreInstanceUniqueIndex(dbIndex types.Index) bool {
	if !dbIndex.Unique {
		return false
	}
	// 后需优化支持多个
	if len(dbIndex.Keys) != 2 {
		return false
	}

	if _, exist := dbIndex.Keys["bk_parent_id"]; !exist {
		return false
	}

	if _, exist := dbIndex.Keys["bk_inst_name"]; !exist {
		return false
	}

	return true

}
