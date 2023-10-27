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
	"fmt"

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

// AddModelWithMultipleBizIDCondition 此函数与上面函数的区别是此函数适用于多个bizId的场景。当传入多个biz id的场景，需要对每一个bizID进行校验,
// 如果加入条件的biz id是单个，请用上面 AddModelBizIDCondition 进行操作。
func AddModelWithMultipleBizIDCondition(cond mapstr.MapStr, modelBizIDs []int64) error {
	var modelBizIDOrCondArr []mapstr.MapStr

	if len(modelBizIDs) <= 1 {
		return fmt.Errorf("biz id must be set")
	}

	bizIds := make([]int64, 0)
	for _, bizId := range modelBizIDs {
		// 此场景下不允许biz id为0
		if bizId <= 0 {
			return fmt.Errorf("biz id illegal biz id: %d", bizId)

		} else {
			bizIds = append(bizIds, bizId)
		}
	}

	if len(bizIds) > 0 {
		// special business model and global shared model
		modelBizIDOrCondArr = []mapstr.MapStr{
			{common.BKAppIDField: mapstr.MapStr{common.BKDBIN: bizIds}},
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
	return nil
}
