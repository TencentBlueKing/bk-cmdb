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

	"go.mongodb.org/mongo-driver/bson"
)

/*
通用模型实例表中的索引。新建模型的时候使用
*/

var (
	instanceDefaultIndexes = []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "bkObjId",
			Keys: bson.D{
				{"bk_obj_id", 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkSupplierAccount_bkObjID",
			Keys: bson.D{
				{"bk_supplier_account", 1},
				{"bk_obj_id", 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkInstId_bkSupplierAccount",
			Keys: bson.D{
				{"bk_inst_id", 1},
				{"bk_supplier_account", 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkInstName_bkSupplierAccount",
			Keys: bson.D{
				{"bk_inst_name", 1},
				{"bk_supplier_account", 1},
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkInstId",
			Keys: bson.D{
				{"bk_inst_id", 1},
			},
			Background: true,
			// 新加 2021年03月11日
			Unique: true,
		},
	}
)

// MainLineInstanceUniqueIndex 建表前需要先建立预定义主线模型的唯一索引
func MainLineInstanceUniqueIndex() []types.Index {

	return []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkParentID_bkInstName",
			Keys: bson.D{
				{"bk_parent_id", 1},
				{"bk_inst_name", 1},
			},
			Background: false,
			Unique:     true,
			PartialFilterExpression: map[string]interface{}{
				"bk_inst_name": map[string]interface{}{"$type": "string"},
				"bk_parent_id": map[string]interface{}{"$type": "number"},
			},
		},
	}
}

// InstanceUniqueIndex 建表前需要先建立预非主线模型的唯一索引
func InstanceUniqueIndex() []types.Index {

	return []types.Index{
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkInstName",
			Keys: bson.D{
				{"bk_inst_name", 1},
			},
			Background: false,
			Unique:     true,
			PartialFilterExpression: map[string]interface{}{
				common.BKInstNameField: map[string]string{common.BKDBType: "string"},
			},
		},
	}
}
