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
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
)

func (lgc *Logics) GetSetIDByCond(pheader http.Header, cond []metadata.ConditionItem) ([]int64, error) {
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)

	query := &metadata.QueryInput{
		Condition: condc,
		Fields:    common.BKSetIDField,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKSetIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	setIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			return nil, fmt.Errorf("invalid set id: %v", err)
		}
		setIDArr = append(setIDArr, setID)
	}
	return setIDArr, nil
}

func (lgc *Logics) GetSetMapByCond(pheader http.Header, fields string, cond interface{}) (map[int64]mapstr.MapStr, error) {
	query := &metadata.QueryInput{
		Condition: cond,
		Fields:    fields,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKModuleIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	setMap := make(map[int64]mapstr.MapStr)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			return nil, fmt.Errorf("invalid set id %v", err)
		}

		setMap[setID] = i
	}
	return setMap, nil
}
