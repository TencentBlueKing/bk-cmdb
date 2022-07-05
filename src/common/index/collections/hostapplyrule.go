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
	registerIndexes(common.BKTableNameHostApplyRule, deprecatedHostApplyRuleIndexes)
	registerIndexes(common.BKTableNameHostApplyRule, commHostApplyRuleIndexes)

}

//  新加和修改后的索引,索引名字一定要用对应的前缀，CCLogicUniqueIdxNamePrefix|common.CCLogicIndexNamePrefix
var commHostApplyRuleIndexes = []types.Index{

	{
		Name: common.CCLogicUniqueIdxNamePrefix + "bizID_ModuleID_serviceTemplateID_attrID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKModuleIDField, 1},
			{common.BKServiceTemplateIDField, 1},
			{common.BKAttributeIDField, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "host_property_under_service_template",
		Keys: bson.D{
			{common.BKServiceTemplateIDField, 1},
			{common.BKAttributeIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bizID_serviceTemplateID_attrID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKServiceTemplateIDField, 1},
			{common.BKAttributeIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "bizID_moduleID_attrID",
		Keys: bson.D{
			{common.BKAppIDField, 1},
			{common.BKModuleIDField, 1},
			{common.BKAttributeIDField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "moduleID_attrID",
		Keys: bson.D{
			{common.BKModuleIDField, 1},
			{common.BKAttributeIDField, 1},
		},
		Background: true,
	},
}

// deprecated 未规范化前的索引，只允许删除不允许新加和修改，
var deprecatedHostApplyRuleIndexes = []types.Index{
	{
		Name: "bk_biz_id",
		Keys: bson.D{{
			"bk_biz_id", 1},
		},
		Background: false,
	},
	{
		Name: "id",
		Keys: bson.D{{
			"id", 1},
		},
		Unique:     true,
		Background: false,
	},
	{
		Name: "bk_module_id",
		Keys: bson.D{{
			"bk_module_id", 1},
		},
		Background: false,
	},
}
