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
	registerIndexes(common.BKTableNameTopoGraphics, commTopoGraphicsIndexes)
}

var commTopoGraphicsIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "scopeType_scopeID_nodeType_bkObjID_bkInstID",
		Keys: bson.D{
			{"scope_type", 1},
			{"scope_id", 1},
			{"node_type", 1},
			{common.BKObjIDField, 1},
			{"bk_inst_id", 1},
		},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
