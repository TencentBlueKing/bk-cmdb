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

package util

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
)

// AddModelBizIDCondition add model bizID condition according to bizID value
func AddModelBizIDCondition(cond mapstr.MapStr, modelBizID int64) {
	var modelBizIDOrCondArr []mapstr.MapStr
	if modelBizID > 0 {
		// special business model and global shared model
		modelBizIDOrCondArr = []mapstr.MapStr{
			{common.BKAppIDField: modelBizID},
			{common.BKAppIDField: 0},
			{common.BKAppIDField: mapstr.MapStr{common.BKDBExists: false}},
		}
	} else {
		// global shared model
		modelBizIDOrCondArr = []mapstr.MapStr{
			{common.BKAppIDField: 0},
			{common.BKAppIDField: mapstr.MapStr{common.BKDBExists: false}},
		}
	}

	if _, exists := cond[common.BKDBOR]; !exists {
		cond[common.BKDBOR] = modelBizIDOrCondArr
	} else {
		andCondArr := []map[string]interface{}{
			{common.BKDBOR: modelBizIDOrCondArr},
		}

		andCond, exists := cond[common.BKDBAND]
		if !exists {
			cond[common.BKDBAND] = andCondArr
		} else {
			cond[common.BKDBAND] = append(andCondArr, map[string]interface{}{common.BKDBAND: andCond})
		}
	}
	delete(cond, common.BKAppIDField)
}
