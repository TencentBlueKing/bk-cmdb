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
	registerIndexes(common.BKTableNameProcessInstanceRelation, commProcessInstanceRelationIndexes)
}

var commProcessInstanceRelationIndexes = []types.Index{
	{
		Name:                    common.CCLogicIndexNamePrefix + "serviceInstanceID",
		Keys:                    bson.D{{common.BKServiceInstanceIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "processTemplateID",
		Keys:                    bson.D{{common.BKProcessTemplateIDField, 1}},
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
		Name:                    common.CCLogicIndexNamePrefix + "bkProcessID",
		Keys:                    bson.D{{common.BKProcessIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "serviceInstanceID_bkProcessID",
		Keys:                    bson.D{{common.BKServiceInstanceIDField, 1}, {common.BKProcessIDField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicUniqueIdxNamePrefix + "bkProcessID_bkHostID",
		Keys:                    bson.D{{common.BKProcessIDField, 1}, {common.BKHostIDField, 1}},
		Unique:                  true,
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
	{
		Name:                    common.CCLogicIndexNamePrefix + "bkHostID",
		Keys:                    bson.D{{common.BKHostIDField, 1}},
		Background:              true,
		PartialFilterExpression: make(map[string]interface{}),
	},
}
