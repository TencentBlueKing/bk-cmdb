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
关联关系表的索引， 新加关联关系表的时候使用
*/

var (
	associationDefaultIndexes = []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "bkObjId_bkInstID",
			Keys: bson.D{
				{common.BKObjIDField, 1},
				{common.BKInstIDField, 1},
			},
			Background: true,
		},
		{
			Name:       common.CCLogicUniqueIdxNamePrefix + "id",
			Keys:       bson.D{{common.BKFieldID, 1}},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkAsstObjId_bkAsstInstId",
			Keys: bson.D{
				{common.BKAsstObjIDField, 1},
				{common.BKAsstInstIDField, 1},
			},
			Background: true,
		},
		{
			Name:       common.CCLogicIndexNamePrefix + "bkAsstID",
			Keys:       bson.D{{common.AssociationKindIDField, 1}},
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "bkInstID_bkAsstInstID_bkObjAsstID",
			Keys: bson.D{
				{common.BKInstIDField, 1},
				{common.BKAsstInstIDField, 1},
				{common.AssociationObjAsstIDField, 1},
			},
			Unique:     true,
			Background: true,
		},
	}
)
