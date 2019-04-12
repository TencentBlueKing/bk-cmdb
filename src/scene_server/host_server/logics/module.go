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
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	parse "configcenter/src/common/paraparse"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetResoulePoolModuleID(ctx context.Context, condition mapstr.MapStr) (int64, errors.CCError) {
	query := &metadata.QueryCondition{
		Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    []string{common.BKModuleIDField},
		Condition: condition,
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetResoulePoolModuleID http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return -1, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetResoulePoolModuleID http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return -1, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		blog.Errorf("GetResoulePoolModuleID http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return -1, lgc.ccErr.Errorf(common.CCErrTopoGetAppFaild, lgc.ccLang.Languagef("host_resource_pool_not_exist"))
	}

	return result.Data.Info[0].Int64(common.BKModuleIDField)
}

func (lgc *Logics) GetNormalModuleByModuleID(ctx context.Context, appID, moduleID int64) ([]mapstr.MapStr, errors.CCError) {
    query := &metadata.QueryCondition{
        Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
        SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
        Fields:    []string{common.BKModuleIDField},
        Condition: hutil.NewOperation().WithAppID(appID).WithModuleID(moduleID).Data(),
    }

    result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
    if err != nil {
        blog.Errorf("GetNormalModuleByModuleID http do error, err:%s,input:%#v,rid:%s", err.Error(), query, lgc.rid)
        return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
    }
    if !result.Result {
        blog.Errorf("GetNormalModuleByModuleID http response error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
        return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
    }

    return result.Data.Info, nil
}

func (lgc *Logics) GetModuleIDByCond(ctx context.Context, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	parse.ParseCommonParams(cond, condc)

	query := &metadata.QueryCondition{
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    []string{common.BKModuleIDField},
		Condition: mapstr.NewFromMap(condc),
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetModuleIDByCond http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetModuleIDByCond http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	moduleIDArr := make([]int64, 0)
	for _, i := range result.Data.Info {
		moduleID, err := i.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("GetModuleIDByCond convert  module id to int error, err:%s, module:%+v,input:%+v,rid:%s", err.Error(), i, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		moduleIDArr = append(moduleIDArr, moduleID)
	}
	return moduleIDArr, nil
}

func (lgc *Logics) GetModuleMapByCond(ctx context.Context, fields []string, cond mapstr.MapStr) (map[int64]types.MapStr, errors.CCError) {

	query := &metadata.QueryCondition{
		Condition: cond,
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKModuleIDField).ToSearchSortArr(),
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDModule, query)
	if err != nil {
		blog.Errorf("GetModuleMapByCond http do error, err:%s,input:%+v,rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetModuleMapByCond http reponse error, err code:%d, err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	moduleMap := make(map[int64]types.MapStr)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKModuleIDField)
		if err != nil {
			blog.Errorf("GetModuleMapByCond convert  module id to int error, err:%s, module:%+v,input:%+v,rid:%s", err.Error(), info, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDModule, common.BKModuleIDField, "int", err.Error())
		}
		moduleMap[id] = info
	}

	return moduleMap, nil
}
