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
	registerIndexes(common.BKTableNameBaseModule, commModuleBaseIndexes)
}

var commModuleBaseIndexes = []types.Index{
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkModuleName",
		Keys:                    bson.D{{common.BKModuleNameField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "default",
		Keys:                    bson.D{{common.BKDefaultField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkBizID",
		Keys:                    bson.D{{common.BKAppIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkSetID",
		Keys:                    bson.D{{common.BKSetIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkParentID",
		Keys:                    bson.D{{common.BKInstParentStr, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkBizID_bkSetID_bkModuleName",
		Keys:       bson.D{{common.BKAppIDField, 1}, {common.BKSetIDField, 1}, {common.BKModuleNameField, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKAppIDField:      map[string]string{common.BKDBType: "number"},
			common.BKSetIDField:      map[string]string{common.BKDBType: "number"},
			common.BKModuleNameField: map[string]string{common.BKDBType: "string"}},
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkModuleID_bkBizID",
		Keys:                    bson.D{{common.BKModuleIDField, 1}, {common.BKAppIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkModuleID",
		Keys:                    bson.D{{common.BKModuleIDField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "setTemplateID_serviceTemplateID",
		Keys:                    bson.D{{common.BKSetTemplateIDField, 1}, {common.BKServiceTemplateIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "setTemplateID",
		Keys:                    bson.D{{common.BKSetTemplateIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "serviceTemplateID",
		Keys:                    bson.D{{common.BKServiceTemplateIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
