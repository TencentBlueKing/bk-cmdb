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
	registerIndexes(common.BKTableNameBaseSet, commSetBaseIndexes)
}

var commSetBaseIndexes = []types.Index{
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkParentID",
		Keys:                    bson.D{{common.BKInstParentStr, 1}},
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
		Name:                    common.CCLogicIndexNamePrefix + "bkSetName",
		Keys:                    bson.D{{common.BKSetNameField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkParentID_bkSetName",
		Keys:                    bson.D{{common.BKInstParentStr, 1}, {common.BKSetNameField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:       common.CCLogicUniqueIdxNamePrefix + "bkBizID_bkSetName_bkParentID",
		Keys:       bson.D{{common.BKAppIDField, 1}, {common.BKSetNameField, 1}, {common.BKInstParentStr, 1}},
		Unique:     true,
		Background: true,
		PartialFilterExpression: map[string]interface{}{
			common.BKSetNameField:  map[string]string{common.BKDBType: "string"},
			common.BKInstParentStr: map[string]string{common.BKDBType: "number"},
			common.BKAppIDField:    map[string]string{common.BKDBType: "number"}},
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkSetID_bkBizID",
		Keys:                    bson.D{{common.BKSetIDField, 1}, {common.BKAppIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkSetID",
		Keys:                    bson.D{{common.BKSetIDField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
