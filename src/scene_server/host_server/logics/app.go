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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	params "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetDefaultAppIDWithSupplier(kit *rest.Kit) (int64, errors.CCError) {
	cond := hutil.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).Data()
	appDetails, err := lgc.GetAppDetails(kit, common.BKAppIDField, cond)
	if err != nil {
		return -1, err
	}

	id, err := util.GetInt64ByInterface(appDetails[common.BKAppIDField])
	if nil != err {
		blog.ErrorJSON("GetDefaultAppIDWithSupplier failed, parse bk_biz_id field from app details failed, inst:%s, err:%s, rid:%s", appDetails, err.Error(), kit.Rid)
		return -1, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}
	return id, nil
}

func (lgc *Logics) GetDefaultAppID(kit *rest.Kit) (int64, errors.CCError) {
	cond := hutil.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).Data()
	appDetails, err := lgc.GetAppDetails(kit, common.BKAppIDField, cond)
	if err != nil {
		return -1, err
	}

	id, err := appDetails.Int64(common.BKAppIDField)
	if nil != err {
		blog.ErrorJSON("GetDefaultAppID failed, parse bk_biz_id from app detail failed, inst:%s, err:%s, rid:%s", appDetails, err.Error(), kit.Rid)
		return -1, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
	}
	return id, nil
}

func (lgc *Logics) GetAppDetails(kit *rest.Kit, fields string, condition map[string]interface{}) (types.MapStr, errors.CCError) {
	fields = fields + "," + common.BkSupplierAccount
	input := &metadata.QueryCondition{
		Condition: condition,
		Fields:    strings.Split(fields, ","),
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, input)
	if err != nil {
		blog.Errorf("GetAppDetail http do error, err:%s, input:%+v, rid:%s", err.Error(), condition, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetAppDetail http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, condition, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	for idx, biz := range result.Data.Info {
		if kit.SupplierAccount == biz[common.BkSupplierAccount].(string) {
			return result.Data.Info[idx], nil
		}
	}

	return nil, errors.New(common.CCErrCommBizNotFoundError, "find resource pool biz failed")
}

func (lgc *Logics) IsHostExistInApp(kit *rest.Kit, appID, hostID int64) (bool, errors.CCErrorCoder) {
	conf := metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
	}

	result, err := lgc.CoreAPI.CoreService().Host().GetHostModulesIDs(kit.Ctx, kit.Header, &conf)
	if err != nil {
		blog.Errorf("IsHostExistInApp http do error, err:%s, input:%+v, rid:%s", err.Error(), hostID, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if err := result.CCError(); err != nil {
		blog.Errorf("IsHostExistInApp http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, hostID, kit.Rid)
		return false, err
	}

	if result.Data == nil {
		return false, nil
	}

	if len(result.Data) == 0 {
		return false, nil
	}

	return true, nil
}

func (lgc *Logics) GetSingleApp(kit *rest.Kit, cond mapstr.MapStr) (mapstr.MapStr, errors.CCError) {
	cond.Set(common.BKDataStatusField, mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled})
	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Start: 0, Limit: 1, Sort: common.BKAppIDField},
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)

	if err != nil {
		blog.Errorf("GetSingleApp http do error, err:%s, input:%+v, rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetSingleApp http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return nil, nil
	}
	return result.Data.Info[0], nil
}

func (lgc *Logics) GetAppIDByCond(kit *rest.Kit, cond []metadata.ConditionItem) ([]int64, errors.CCError) {
	condc := make(map[string]interface{})
	if err := params.ParseCommonParams(cond, condc); err != nil {
		blog.Errorf("ParseCommonParams failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, err

	}
	condMap := mapstr.NewFromMap(condc)
	condMap.Set(common.BKDataStatusField, mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled})

	query := &metadata.QueryCondition{
		Condition: condMap,
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKAppIDField},
		Fields:    []string{common.BKAppIDField},
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("GetAppIDByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetAppIDByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}

	appIDs := make([]int64, 0)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKAppIDField)
		if err != nil {
			blog.ErrorJSON("GetAppIDByCond failed, convert bk_biz_id to int error, inst:%s  input:%s, err:%s, rid:%s", info, query, err.Error(), kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		appIDs = append(appIDs, id)
	}

	return appIDs, nil
}

func (lgc *Logics) GetAppMapByCond(kit *rest.Kit, fields []string, cond mapstr.MapStr) (map[int64]types.MapStr, errors.CCError) {

	if cond == nil {
		cond = mapstr.New()
	}
	cond.Set(common.BKDataStatusField, mapstr.MapStr{common.BKDBNE: common.DataStatusDisabled})
	query := &metadata.QueryCondition{
		Condition: cond,
		Page:      metadata.BasePage{Start: 0, Limit: common.BKNoLimit, Sort: common.BKAppIDField},
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDApp, query)
	if err != nil {
		blog.Errorf("GetAppMapByCond http do error, err:%s, input:%+v, rid:%s", err.Error(), query, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("GetAppMapByCond http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, query, kit.Rid)
		return nil, kit.CCError.New(result.Code, result.ErrMsg)
	}
	appMap := make(map[int64]types.MapStr)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKAppIDField)
		if err != nil {
			blog.Errorf("GetAppMapByCond http response format error,convert bk_biz_id to int error, err:%s, inst:%+v  input:%+v, rid:%s", err.Error(), info, query, kit.Rid)
			return nil, kit.CCError.Errorf(common.CCErrCommInstFieldConvertFail, common.BKInnerObjIDApp, common.BKAppIDField, "int", err.Error())
		}
		appMap[id] = info
	}

	return appMap, nil
}

func (lgc *Logics) ExistInnerModule(kit *rest.Kit, moduleIDArr []int64) (bool, errors.CCErrorCoder) {
	input := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKDefaultField: map[string]interface{}{
				"$ne": common.DefaultFlagDefaultValue,
			},
			common.BKModuleIDField: map[string]interface{}{
				"$in": moduleIDArr,
			},
		},
	}
	result, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, common.BKInnerObjIDModule, input)
	if err != nil {
		blog.Errorf("ExistInnerModule http do error, err:%s, input:%+v, rid:%s", err.Error(), input, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("ExistInnerModule http response error, err code:%d, err msg:%s, input:%+v, rid:%s", result.Code, result.ErrMsg, input, kit.Rid)
		return false, result.CCError()
	}

	exist := result.Data.Count > 0
	return exist, nil
}
