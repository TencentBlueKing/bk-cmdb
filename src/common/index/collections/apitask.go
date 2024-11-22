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
	"configcenter/src/common/metadata"
	"configcenter/src/storage/dal/types"

	"go.mongodb.org/mongo-driver/bson"
)

func init() {

	// 先注册未规范化的索引，如果索引出现冲突旧，删除未规范化的索引
	registerIndexes(common.BKTableNameAPITask, deprecatedAPITaskIndexes)
	registerIndexes(common.BKTableNameAPITask, commAPITaskIndexes)
	registerIndexes(common.BKTableNameAPITaskSyncHistory, apiTaskSyncHistoryIndexes)

}

var commAPITaskIndexes = []types.Index{
	{
		Name:       common.CCLogicIndexNamePrefix + "lastTime",
		Keys:       bson.D{{common.LastTimeField, -1}},
		Background: true,
		// delete redundant tasks from 6 months ago
		ExpireAfterSeconds: 6 * 30 * 24 * 60 * 60,
	},
	{
		Name: common.CCLogicIndexNamePrefix + "taskType_status_createTime",
		Keys: bson.D{
			{common.BKTaskTypeField, 1},
			{common.BKStatusField, 1},
			{common.CreateTimeField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "tenantID_taskType_instID_extra",
		Keys: bson.D{
			{common.TenantID, 1},
			{common.BKTaskTypeField, 1},
			{common.BKInstIDField, 1},
			{metadata.APITaskExtraField, 1},
		},
		Background: true,
		Unique:     true,
		PartialFilterExpression: map[string]interface{}{
			common.BKStatusField: map[string]interface{}{
				common.BKDBIN: []metadata.APITaskStatus{metadata.APITaskStatusNew, metadata.APITaskStatusWaitExecute,
					metadata.APITaskStatusExecute},
			},
		},
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "tenantID_instID_taskType_createTime",
		Keys: bson.D{
			{common.TenantID, 1},
			{common.BKInstIDField, 1},
			{common.BKTaskTypeField, 1},
			{common.CreateTimeField, -1},
		},
		Background: true,
	},
}

var apiTaskSyncHistoryIndexes = []types.Index{
	{
		Name:       common.CCLogicIndexNamePrefix + "lastTime",
		Keys:       bson.D{{common.LastTimeField, -1}},
		Background: true,
		// delete redundant tasks from 6 months ago
		ExpireAfterSeconds: 6 * 30 * 24 * 60 * 60,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "tenantID_taskID_taskType",
		Keys: bson.D{
			{common.TenantID, 1},
			{common.BKTaskIDField, 1},
			{common.BKTaskTypeField, 1},
		},
		Background: true,
	},
	{
		Name: common.CCLogicUniqueIdxNamePrefix + "tenantID_instID_taskType_createTime",
		Keys: bson.D{
			{common.TenantID, 1},
			{common.BKInstIDField, 1},
			{common.BKTaskTypeField, 1},
			{common.CreateTimeField, -1},
		},
		Background: true,
	},
}

// deprecated 未规范化前的索引，只允许删除不允许新加和修改，
var deprecatedAPITaskIndexes = []types.Index{
	{
		Name: "idx_taskID",
		Keys: bson.D{{
			"task_id", 1},
		},
		Unique:     true,
		Background: true,
	},
}
