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
	registerIndexes(common.BKTableNameBaseBizSet, commBizSetIndexes)
}

var commBizSetIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "biz_set_id",
		Keys: bson.D{{
			common.BKBizSetIDField, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "biz_set_name",
		Keys: bson.D{{
			common.BKBizSetNameField, 1},
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "biz_set_id_biz_set_name_owner_id",
		Keys: bson.D{
			{common.BKBizSetIDField, 1},
			{common.BKBizSetNameField, 1},
			{common.BKOwnerIDField, 1},
		},
		Background: true,
	},
}
