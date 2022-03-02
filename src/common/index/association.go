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
)

/*
关联关系表的索引， 新加关联关系表的时候使用
*/

var (
	associationDefaultIndexes = []types.Index{
		{
			Name: common.CCLogicIndexNamePrefix + "bkObjId_bkInstID",
			Keys: map[string]int32{
				"bk_obj_id":  1,
				"bk_inst_id": 1,
			},
			Background: true,
		},
		{
			Name: common.CCLogicUniqueIdxNamePrefix + "id",
			Keys: map[string]int32{
				"id": 1,
			},
			Unique:     true,
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkAsstObjId_bkAsstInstId",
			Keys: map[string]int32{
				"bk_asst_obj_id":  1,
				"bk_asst_inst_id": 1,
			},
			Background: true,
		},
		{
			Name: common.CCLogicIndexNamePrefix + "bkAsstID",
			Keys: map[string]int32{
				"bk_asst_id": 1,
			},
			Background: true,
		},
	}
)
