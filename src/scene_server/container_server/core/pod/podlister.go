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

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// Lister pod lister
type Lister struct {
	clientSet apimachinery.ClientSetInterface
}

// NewLister create lister
func NewLister() *Lister {
	return &Lister{}
}

func (li *Lister) getSetModuleIDs(kit *rest.Kit, bizID int64, setIDs []int64) ([]int64, error) {
	// query module ids
	filter := map[string]interface{}{
		common.BKAppIDField: bizID,
		common.BKSetIDField: common.KvMap{
			common.BKDBIN: setIDs,
		},
	}
	query := &metadata.QueryCondition{
		Condition: filter,
		Page: metadata.BasePage{
			Start: 0,
			Limit: common.BKNoLimit,
		},
	}

	queryResult, err := li.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, query)
	if err != nil || !queryResult.Result {
		blog.Errorf("getSetModuleIDs query set instances with %#v failed, resp %#v, err %s", query, queryResult, err.Error())
		return nil, fmt.Errorf("getSetModuleIDs query set instances with %#v failed, resp %#v, err %s", query, queryResult, err.Error())
	}

	var moduleids []int64
	for _, module := range queryResult.Data.Info {
		moduleID, err := module.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("get %s from module %v, err %s", common.BKModuleIDField, module, err.Error())
			return nil, fmt.Errorf("get %s from module %v, err %s", common.BKModuleIDField, module, err.Error())
		}
		moduleids = append(moduleids, moduleID)
	}
	return moduleids, nil
}

// ListPod list pod
func (li *Lister) ListPod(kit *rest.Kit, option metadata.ListPods) (*metadata.ListPodsResult, error) {

	var moduleIDs []int64
	if len(option.SetIDs) != 0 {
		setModuleIDs, err := li.getSetModuleIDs(kit, option.BizID, option.SetIDs)
		if err != nil {
			return nil, err
		}
		moduleIDs = append(moduleIDs, setModuleIDs...)
	} else if len(option.ModuleIDs) != 0 {
		moduleIDs = append(moduleIDs, option.ModuleIDs...)
	}

	moduleFilter := map[string]interface{}{
		common.BKAppIDField: option.BizID,
	}
	if len(moduleIDs) != 0 {
		moduleFilter[common.BKModuleIDField] = common.KvMap{
			common.BKDBIN: moduleIDs,
		}
	}

	if option.PodPropertyFilter != nil {
		errKey, err := option.PodPropertyFilter.Validate()
		if err != nil {
			blog.Errorf("pod_property_filter validate failed, errkey %s, err %s", errKey, err.Error())
			return nil, err
		}
		mgoFilter, key, err := option.PodPropertyFilter.ToMgo()
		if err != nil {
			return nil, fmt.Errorf("invalid key:pod_property_filter.%s, err: %s", key, err)
		}
		for key, val := range mgoFilter {
			moduleFilter[key] = val
		}
	}

	query := &metadata.QueryCondition{
		Condition: moduleFilter,
		Fields:    option.Fields,
		Page:      option.Page,
	}
	queryResult, err := li.clientSet.CoreService().Instance().ReadInstance(
		kit.Ctx, kit.Header, common.BKInnerObjIDPod, query)
	if err != nil {
		blog.Errorf("read pod instance failed, err %s", err.Error())
		return nil, fmt.Errorf("read pod instance failed, err %s", err.Error())
	}
	if !queryResult.Result {
		blog.Errorf("read pod instance return false, result %#v", queryResult)
		return nil, fmt.Errorf("read pod instance return false, result %#v", queryResult)
	}
	return &metadata.ListPodsResult{
		Count: queryResult.Data.Count,
		Info:  queryResult.Data.Info,
	}, nil
}
