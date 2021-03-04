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
通用模型实例表中的索引。新建模型的时候使用
*/

var (
	instanceDefaultIndex = []types.Index{
		{
			Keys: map[string]int32{
				common.BKObjIDField: 1,
			},
			Name:       "bkcc_idx_ObjID",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKOwnerIDField: 1,
			},
			Name:       "bkcc_idx_supplierAccount",
			Background: true,
		},
		{
			Keys: map[string]int32{
				common.BKInstIDField: 1,
			},
			Name:       "bkcc_idx_InstId",
			Background: true,
			Unique:     true,
		},
		{
			Keys: map[string]int32{
				common.BKInstNameField: 1,
			},
			Name:       "bkcc_idx_InstName",
			Background: true,
		},
	}
)
