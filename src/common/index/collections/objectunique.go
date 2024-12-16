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
	registerIndexes(common.BKTableNameObjUnique, commObjectUniqueIndexes)
}

var commObjectUniqueIndexes = []types.Index{
	{
		Name: common.CCLogicIndexNamePrefix + "bkTemplateID",
		Keys: bson.D{
			{
				common.BKTemplateID, 1,
			},
		},
		Background: true,
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkObjID",
		Keys:                    bson.D{{common.BKObjIDField, 1}},
		Background:              false,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "ID",
		Keys:                    bson.D{{common.BKFieldID, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
