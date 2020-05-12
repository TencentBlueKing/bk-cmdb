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

func (p *PodManager) checkModuleExist(kit *rest.Kit, moduleID string) (bool, error) {

	filter := map[string]interface{}{
		common.BKModuleIDField: common.KvMap{common.BKDBIN: moduleID},
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
		blog.Errorf("checkModuleExist ReadInstance http do err. error: %s, input: %#v, rid: %s", err.Error(), query, kit.Rid)
		return false, fmt.Errorf("checkModuleExist ReadInstance http do err. error: %s, input: %#v, rid: %s", err.Error(), query, kit.Rid)
	}
	if !queryResult.Result {
		blog.Errorf("checkModuleExist ReadInstance result false, reply: %#v, input: %#v, rid: %s", queryResult, query, kit.Rid)
		return false, fmt.Errorf("checkModuleExist ReadInstance result false, reply: %#v, input: %#v, rid: %s", queryResult, query, kit.Rid)
	}

	if len(queryResult.Data.Info) == 0 {
		return false, nil
	}
	return true, nil
}
