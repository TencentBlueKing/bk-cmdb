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
	"configcenter/src/common"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {

	// 先注册未规范化的索引，如果索引出现冲突旧，删除未规范化的索引
	registerIndexes(common.BKTableNameBaseHost, deprecatedHostBaseIndexes)
	registerIndexes(common.BKTableNameBaseHost, commHostBaseIndexes)

}

//  新加和修改后的索引,索引名字一定要用对应的前缀，CCLogicUniqueIdxNamePrefix|common.CCLogicIndexNamePrefix
var commHostBaseIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIP_bkCloudID",
		Keys: bson.D{
			{common.BKHostInnerIPField, 1},
			{common.BKCloudIDField, 1},
		},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKHostInnerIPField: map[string]string{common.BKDBType: "string"},
			common.BKCloudIDField:     map[string]string{common.BKDBType: "number"},
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bkHostInnerIPv6_bkCloudID",
		Keys: bson.D{
			{common.BKHostInnerIPv6Field, 1},
			{common.BKCloudIDField, 1},
		},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKHostInnerIPv6Field: map[string]string{common.BKDBType: "string"},
			common.BKCloudIDField:       map[string]string{common.BKDBType: "number"},
		},
	},
}

// deprecated 未规范化前的索引，只允许删除不允许新加和修改，
var deprecatedHostBaseIndexes = []types.Index{
	{
		Name: "bk_host_name_1",
		Keys: bson.D{{
			"bk_host_name", 1},
		},
		Background: true,
	},
	{
		Name: "bk_host_innerip_1",
		Keys: bson.D{{
			"bk_host_innerip", 1},
		},
		Background: true,
	},
	{
		Name: "bk_host_id_1_bk_supplier_account_1",
		Keys: bson.D{
			{"bk_host_id", 1},
			{"bk_supplier_account", 1},
		},
		Background: true,
	},
	/* 	{
		Name: "innerIP_platID",
		Keys: bson.D{{
			"bk_host_innerip", 1}
			"bk_cloud_id":     1,
		},
		Background: false,
	}, */
	{
		Name: "bk_supplier_account_1",
		Keys: bson.D{{
			"bk_supplier_account", 1},
		},
		Background: true,
	},
	{
		Name: "bk_cloud_id_1",
		Keys: bson.D{{
			"bk_cloud_id", 1},
		},
		Background: true,
	},
	{
		Name: "idx_unique_hostID",
		Keys: bson.D{{
			"bk_host_id", 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: "cloudInstID",
		Keys: bson.D{{
			"bk_cloud_inst_id", 1},
		},
		Background: true,
	},
	{
		Name: "bk_idx_bk_asset_id",
		Keys: bson.D{{
			"bk_asset_id", 1},
		},
		Background: true,
	},
	{
		Name: "bk_os_type_1",
		Keys: bson.D{{
			"bk_os_type", 1},
		},
		Background: true,
	},
}
