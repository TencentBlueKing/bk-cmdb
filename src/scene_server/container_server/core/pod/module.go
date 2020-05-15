/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package pod

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

func (p *PodManager) checkModuleIDs(kit *rest.Kit, bizID int64, moduleIDs []int64) (bool, error) {

	// get unique module ids
	var uniqueModuleIDs []int64
	moduleIDMap := make(map[int64]bool)
	for _, moduleID := range moduleIDs {
		if _, ok := moduleIDMap[moduleID]; !ok {
			uniqueModuleIDs = append(uniqueModuleIDs, moduleID)
		}
	}

	// query module ids
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKModuleIDField: common.KvMap{
			common.BKDBIN: uniqueModuleIDs,
		},
	}
	query := &metadata.QueryCondition{
		Condition: filter,
		Page: metadata.BasePage{
			Start: 0,
			Limit: common.BKNoLimit,
		},
	}
	queryResult, err := p.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("checkModuleIDs ReadInstance http do err. error: %s, input: %#v, rid: %s", err.Error(), query, kit.Rid)
		return false, fmt.Errorf("checkModuleIDs ReadInstance http do err. error: %s, input: %#v, rid: %s", err.Error(), query, kit.Rid)
	}
	if !queryResult.Result {
		blog.Errorf("checkModuleIDs ReadInstance result false, reply: %#v, input: %#v, rid: %s", queryResult, query, kit.Rid)
		return false, fmt.Errorf("checkModuleIDs ReadInstance result false, reply: %#v, input: %#v, rid: %s", queryResult, query, kit.Rid)
	}

	// return false if not all modules exists or not all modules belong to this business
	if queryResult.Data.Count != len(uniqueModuleIDs) {
		return false, nil
	}
	return true, nil
}
