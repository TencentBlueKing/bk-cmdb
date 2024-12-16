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
	registerIndexes(common.BKTableNameDelArchive, commDelArchiveIndexes)
	registerIndexes(common.BKTableNameKubeDelArchive, commKubeDelArchiveIndexes)
}

var commDelArchiveIndexes = []types.Index{
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "oid_coll",
		Keys:                    bson.D{{"oid", 1}, {"coll", 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "coll",
		Keys:                    bson.D{{"coll", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "oid",
		Keys:                    bson.D{{"oid", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "time",
		Keys:                    bson.D{{common.FieldTypeTime, -1}},
		Background:              true,
		ExpireAfterSeconds:      7 * 24 * 60 * 60,
		PartialFilterExpression: make(map[string]interface{}),
	},
}

var commKubeDelArchiveIndexes = []types.Index{
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "coll_oid",
		Keys: bson.D{
			{"coll", 1},
			{"oid", 1},
		},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "coll",
		Keys:                    bson.D{{"coll", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "oid",
		Keys:                    bson.D{{"oid", 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "time",
		Keys:                    bson.D{{common.FieldTypeTime, -1}},
		Background:              true,
		ExpireAfterSeconds:      2 * 24 * 60 * 60,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
