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
)

func init() {

	// 先注册未规范化的索引，如果索引出现冲突旧，删除未规范化的索引
	registerIndexes(common.BKTableNameProcessInstanceRelation, deprecatedProcessInstanceRelationIndexes)
	registerIndexes(common.BKTableNameProcessInstanceRelation, commProcessInstanceRelationIndexes)

}

//  新加和修改后的索引,索引名字一定要用对应的前缀，CCLogicUniqueIdxNamePrefix|common.CCLogicIndexNamePrefix

var commProcessInstanceRelationIndexes = []types.Index{}

// deprecated 未规范化前的索引，只允许删除不允许新加和修改，
var deprecatedProcessInstanceRelationIndexes = []types.Index{
	{
		Name: "idx_bkServiceInstanceID",
		Keys: map[string]int32{
			"service_instance_id": 1,
		},
		Background: true,
	},
	{
		Name: "idx_bkProcessTemplateID",
		Keys: map[string]int32{
			"process_template_id": 1,
		},
		Background: true,
	},
	{
		Name: "idx_bkBizID",
		Keys: map[string]int32{
			"bk_biz_id": 1,
		},
		Background: true,
	},
	{
		Name: "idx_bkProcessID",
		Keys: map[string]int32{
			"bk_process_id": 1,
		},
		Background: true,
	},
	{
		Name: "idx_bkHostID",
		Keys: map[string]int32{
			"bk_host_id": 1,
		},
		Background: true,
	},
	{
		Name: "idx_unique_serviceInstID_ProcID",
		Keys: map[string]int32{
			"service_instance_id": 1,
			"bk_process_id":       1,
		},
		Unique:     true,
		Background: true,
	},
	{
		Name: "idx_unique_procID_hostID",
		Keys: map[string]int32{
			"bk_process_id": 1,
			"bk_host_id":    1,
		},
		Unique:     true,
		Background: true,
	},
}
