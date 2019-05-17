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
    "strings"

    "configcenter/src/common"
    "configcenter/src/common/blog"
    "configcenter/src/common/errors"
    "configcenter/src/common/mapstr"
    types "configcenter/src/common/mapstr"
    "configcenter/src/common/metadata"
    "configcenter/src/common/paraparse"
    "configcenter/src/common/util"
    hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetDefaultAppIDWithSupplier(ctx context.Context) (int64, errors.CCError) {
	cond := hutil.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).WithOwnerID(util.GetOwnerID(lgc.header)).Data()
	cond[common.BKDBAND] = []mapstr.MapStr{
		mapstr.MapStr{common.BKOwnerIDField: util.GetOwnerID(lgc.header)},
	}
	appDetails, err := lgc.GetAppDetails(ctx, common.BKAppIDField, cond)
	if err != nil {
		return -1, err
	}

	id, err := util.GetInt64ByInterface(appDetails[common.BKAppIDField])
	if nil != err {
		blog.Errorf("GetDefaultAppID http response format error,convert bk_biz_id to int error, err:%s, inst:%+v, rid:%s", err.Error(), appDetails, lgc.rid)
		return -1, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}
	return id, nil
}

func (lgc *Logics) GetDefaultAppID(ctx context.Context, ownerID string) (int64, errors.CCError) {
	cond := hutil.NewOperation().WithOwnerID(ownerID).WithDefaultField(int64(common.DefaultAppFlag)).Data()
	cond[common.BKDBAND] = []mapstr.MapStr{
		mapstr.MapStr{common.BKOwnerIDField: util.GetOwnerID(lgc.header)},
	}
	appDetails, err := lgc.GetAppDetails(ctx, common.BKAppIDField, cond)
	if err != nil {
		return -1, err
	}

	id, err := appDetails.Int64(common.BKAppIDField)
	if nil != err {
		blog.Errorf("GetDefaultAppID http response format error,convert bk_biz_id to int error, err:%s, inst:%+v, rid:%s", err.Error(), appDetails, lgc.rid)
		return -1, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}
	return id, nil
}

func (lgc *Logics) GetAppDetails(ctx context.Context, fields string, condition map[string]interface{}) (types.MapStr, errors.CCError) {

	input := &metadata.QueryCondition{
		Condition: condition,
		Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
		SortArr:   metadata.NewSearchSortParse().String(common.BKAppIDField).ToSearchSortArr(),
		Fields:    strings.Split(fields, ","),
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDApp, input)
	if err != nil {
		blog.Errorf("GetAppDetail http do error, err:%s, input:%+v, rid:%s", err.Error(), condition, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetAppDetail http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, condition, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return make(map[string]interface{}), nil
	}

	return result.Data.Info[0], nil
}

func (lgc *Logics) IsHostExistInApp(ctx context.Context, appID, hostID int64) (bool, errors.CCError) {
	conf := metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
	}

	result, err := lgc.CoreAPI.HostController().Module().GetHostModulesIDs(ctx, lgc.header, &conf)
	if err != nil {
		blog.Errorf("IsHostExistInApp http do error, err:%s, input:%+v, rid:%s", err.Error(), hostID, lgc.rid)
		return false, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("IsHostExistInApp http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, hostID, lgc.rid)
		return false, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if result.Data == nil {
		return false, nil
	}

	if len(result.Data) == 0 {
		return false, nil
	}

	return true, nil
}

// ExistHostIDSInApp exist host id in app return []int64 don't exist in app hostID, error handle logic error
func (lgc *Logics) ExistHostIDSInApp(ctx context.Context, appID int64, hostIDArray []int64) ([]int64, error) {
    defErr := lgc.ccErr
	conf := map[string][]int64{
		common.BKAppIDField:  []int64{appID},
		common.BKHostIDField: hostIDArray,
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx,lgc.header,  conf)
	if err != nil {
		blog.Errorf("ExistHostIDSInApp http do error. err:%s, input:%#v,rid:%s", err.Error(), conf, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("ExistHostIDSInApp http reply error. err code:%d,err msg:%s, input:%#v,rid:%s", result.Code, result.ErrMsg, conf, lgc.rid)
		return nil, defErr.New(result.Code, result.ErrMsg)
	}
	hostIDMap := make(map[int64]bool, 0)
	for _, row := range result.Data {
		hostIDMap[row.HostID] = true
	}
	var notExistHOstID []int64
	for _, hostID := range hostIDArray {
		_, ok := hostIDMap[hostID]
		if !ok {
			notExistHOstID = append(notExistHOstID, hostID)
		}
	}

	return notExistHOstID, nil
}

func (lgc *Logics) GetSingleApp(ctx context.Context, cond mapstr.MapStr) (mapstr.MapStr, errors.CCError) {
	query := &metadata.QueryCondition{
		Condition: cond,
		Limit:     metadata.SearchLimit{Offset: 0, Limit: 1},
		SortArr:   metadata.NewSearchSortParse().String(common.BKAppIDField).ToSearchSortArr(),
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDApp, query)

	if err != nil {
		blog.Errorf("GetSingleApp http do error, err:%s, input:%+v, rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSingleApp http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return nil, nil
	}
	return result.Data.Info[0], nil
}

func (lgc *Logics) GetAppIDByCond(ctx context.Context, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	params.ParseCommonParams(cond, condc)
	query := &metadata.QueryCondition{
		Condition: mapstr.NewFromMap(condc),
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKAppIDField).ToSearchSortArr(),
		Fields:    []string{common.BKAppIDField},
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("GetAppIDByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetAppIDByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}

	appIDs := make([]int64, 0)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKAppIDField)
		if err != nil {
			blog.Errorf("GetAppIDByCond http response format error,convert bk_biz_id to int error, err:%s, inst:%+v  input:%+v, rid:%s", err.Error(), info, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		appIDs = append(appIDs, id)
	}

	return appIDs, nil
}

func (lgc *Logics) GetAppMapByCond(ctx context.Context, fields []string, cond mapstr.MapStr) (map[int64]types.MapStr, errors.CCError) {

	query := &metadata.QueryCondition{
		Condition: cond,
		Limit:     metadata.SearchLimit{Offset: 0, Limit: common.BKNoLimit},
		SortArr:   metadata.NewSearchSortParse().String(common.BKAppIDField).ToSearchSortArr(),
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("GetAppMapByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), query, lgc.rid)
		return nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetAppMapByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, lgc.rid)
		return nil, lgc.ccErr.New(result.Code, result.ErrMsg)
	}
	appMap := make(map[int64]types.MapStr)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKAppIDField)
		if err != nil {
			blog.Errorf("GetAppMapByCond http response format error,convert bk_biz_id to int error, err:%s, inst:%+v  input:%+v, rid:%s", err.Error(), info, query, lgc.rid)
			return nil, lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		appMap[id] = info
	}

	return appMap, nil
}
