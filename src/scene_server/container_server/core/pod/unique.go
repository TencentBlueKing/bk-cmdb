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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

func (p *PodManager) getPodAttrDes(kit *rest.Kit) ([]metadata.Attribute, error) {

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKObjIDField: common.BKInnerObjIDPod,
		},
	}

	queryResult, err := p.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, common.BKInnerObjIDPod, query)
	if err != nil {
		blog.Errorf("read pod attr failed, query %#v, err %s", query, err.Error())
		return nil, fmt.Errorf("read pod attr failed, query %#v, err %s", query, err.Error())
	}
	if !queryResult.Result {
		blog.Errorf("read pod attr return false, query %#v, queryResult %#v, err %s",
			query, queryResult, err.Error())
	}

	return queryResult.Data.Info, nil
}

func (p *PodManager) getPodUnique(kit *rest.Kit) ([]metadata.ObjectUnique, error) {

	query := metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKObjIDField: common.BKInnerObjIDPod,
		},
	}

	queryResult, err := p.clientSet.CoreService().Model().ReadModelAttrUnique(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("read pod uniques failed, query %#v, err %s", query, err.Error())
		return nil, fmt.Errorf("read pod uniques failed, query %#v, err %s", query, err.Error())
	}
	if !queryResult.Result {
		blog.Errorf("read pod uniques return false, query %#v, queryResult %#v, err %s",
			query, queryResult, err.Error())
	}
	return queryResult.Data.Info, nil
}

// validateCondition validate condition for update pod
func validateCondition(cond mapstr.MapStr, uniques []metadata.ObjectUnique, attrs []metadata.Attribute) bool {

	propertyIDMap := make(map[int64]string)
	for _, attr := range attrs {
		propertyIDMap[attr.ID] = attr.PropertyID
	}

	var uniKeysArr [][]string
	for _, uni := range uniques {
		var tmpArr []string
		for _, key := range uni.Keys {
			PropertyID, ok := propertyIDMap[int64(key.ID)]
			if !ok {
				blog.Errorf("pod unique key id is not in pod attr")
				continue
			}
			tmpArr = append(tmpArr, PropertyID)
		}
		uniKeysArr = append(uniKeysArr, tmpArr)
	}

	condMap := make(map[string]bool)
	for key := range cond {
		condMap[key] = true
	}

	for _, uniKeys := range uniKeysArr {
		hasAllKey := true
		for _, key := range uniKeys {
			if _, ok := condMap[key]; !ok {
				hasAllKey = false
			}
		}
		if hasAllKey {
			return true
		}
	}

	// inst id is also nice
	if _, ok := condMap[common.BKInstIDField]; ok {
		return true
	}

	return false
}
