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

package logics

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetResoulePoolModuleID(pheader http.Header, condition interface{}) (int64, error) {
	query := metadata.QueryInput{
		Start:     0,
		Limit:     1,
		Sort:      common.BKModuleIDField,
		Fields:    common.BKModuleIDField,
		Condition: condition,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, pheader, &query)
	if err != nil || (err == nil && !result.Result) {
		return -1, fmt.Errorf("search host obj log failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return -1, errors.New("search resource pool, but get nil data")
	}

	return result.Data.Info[0].Int64(common.BKModuleIDField)
}

func (lgc *Logics) GetModuleByModuleID(pheader http.Header, appID, moduleID int64) ([]mapstr.MapStr, error) {
	query := &metadata.QueryInput{
		Start:     0,
		Limit:     1,
		Sort:      common.BKModuleIDField,
		Fields:    common.BKModuleIDField,
		Condition: hutil.NewOperation().WithAppID(appID).WithModuleID(moduleID).Data(),
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get modules with id failed, err: %v, result err: %s", err, result.ErrMsg)
	}

	return result.Data.Info, nil
}

func (lgc *Logics) GetModuleIDByCond(phader http.Header, cond []metadata.ConditionItem) ([]int64, error) {
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)

	query := &metadata.QueryInput{
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKModuleIDField,
		Fields:    common.BKModuleIDField,
		Condition: condc,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, phader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get module id failed, err: %v, %v", err, result.ErrMsg)
	}

	moduleIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		moduleID, err := i.Int64(common.BKModuleIDField)
		if err != nil {
			return nil, fmt.Errorf("invalid module id: %v", err)
		}
		moduleIDArr = append(moduleIDArr, moduleID)
	}
	return moduleIDArr, nil
}

func (lgc *Logics) GetModuleMapByCond(pheader http.Header, fields string, cond interface{}) (map[int64]types.MapStr, error) {

	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKModuleIDField,
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	moduleMap := make(map[int64]types.MapStr)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKModuleIDField)
		if err != nil {
			return nil, err
		}
		moduleMap[id] = info
	}

	return moduleMap, nil
}
