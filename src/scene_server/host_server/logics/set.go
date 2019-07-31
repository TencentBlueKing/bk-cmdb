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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
)

func (lgc *Logics) GetSetIDByCond(ctx context.Context, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	if err := parse.ParseCommonParams(cond, condc); err != nil {
		blog.Warnf("ParseCommonParams failed, err: %+v, rid: %s", err, lgc.rid)
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.NewFromMap(condc),
		Fields:    []string{common.BKSetIDField},
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKSetIDField).ToSearchSortArr(),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("GetSetIDByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSetIDByCond http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	setIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("GetSetIDByCond convert %s %s to integer error, set info:%+v, input:%+v,rid:%s", common.BKInnerObjIDSet, common.BKSetIDField, i, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
		}
		setIDArr = append(setIDArr, setID)
	}
	return setIDArr, nil
}

func (lgc *Logics) GetSetMapByCond(ctx context.Context, fields []string, cond mapstr.MapStr) (map[int64]mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryCondition{
		Condition: cond,
		Fields:    fields,
		SortArr:   metadata.NewSearchSortParse().String(common.BKSetIDField).ToSearchSortArr(),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDSet, query)
	if err != nil {
		blog.Errorf("GetSetMapByCond http do error, err:%s,objID:%s,input:%+v,rid:%s", err.Error(), common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSetMapByCond http reponse error, err code:%d, err msg:%s,objID:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, common.BKInnerObjIDSet, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	setMap := make(map[int64]mapstr.MapStr)
	for _, i := range result.Data.Info {
		setID, err := i.Int64(common.BKSetIDField)
		if err != nil {
			blog.Errorf("GetSetMapByCond convert %s %s to integer error, set info:%+v, input:%+v,rid:%s", common.BKInnerObjIDSet, common.BKSetIDField, i, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDSet, common.BKSetIDField, "int", err.Error())
		}

		setMap[setID] = i
	}
	return setMap, nil
}
