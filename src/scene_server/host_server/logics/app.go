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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
)

func (lgc *Logics) GetDefaultAppIDWithSupplier(pheader http.Header) (int64, error) {
	cond := hutil.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).WithOwnerID(util.GetOwnerID(pheader)).Data()
	cond[common.BKDBAND] = []mapstr.MapStr{
		mapstr.MapStr{common.BKOwnerIDField: util.GetOwnerID(pheader)},
	}
	appDetails, err := lgc.GetAppDetails(common.BKAppIDField, cond, pheader)
	if err != nil {
		return -1, err
	}

	id, err := util.GetInt64ByInterface(appDetails[common.BKAppIDField])
	if nil != err {
		return -1, errors.New("can not find bk biz field")
	}
	return id, nil
}

func (lgc *Logics) GetDefaultAppID(ownerID string, pheader http.Header) (int64, error) {
	cond := hutil.NewOperation().WithOwnerID(ownerID).WithDefaultField(int64(common.DefaultAppFlag)).Data()
	blog.Infof("get default app id cond: %v", cond)
	cond[common.BKDBAND] = []mapstr.MapStr{
		mapstr.MapStr{common.BKOwnerIDField: util.GetOwnerID(pheader)},
	}
	appDetails, err := lgc.GetAppDetails(common.BKAppIDField, cond, pheader)
	if err != nil {
		return -1, err
	}

	id, err := appDetails.Int64(common.BKAppIDField) //[common.BKAppIDField]
	if nil != err {
		return -1, errors.New("can not find bk biz field")
	}
	return id, nil
}

func (lgc *Logics) GetAppDetails(fields string, condition map[string]interface{}, pheader http.Header) (types.MapStr, error) {
	query := metadata.QueryInput{
		Condition: condition,
		Start:     0,
		Limit:     1,
		Fields:    fields,
		Sort:      common.BKAppIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, pheader, &query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get default appid failed, err: %v, %v", err, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return make(map[string]interface{}), nil
	}

	return result.Data.Info[0], nil
}

func (lgc *Logics) IsHostExistInApp(appID, hostID int64, pheader http.Header) (bool, error) {
	conf := metadata.ModuleHostConfigParams{
		ApplicationID: appID,
		HostID:        hostID,
	}

	result, err := lgc.CoreAPI.HostController().Module().GetHostModulesIDs(context.Background(), pheader, &conf)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("get host module ids failed, err: %v, %v", err, result.ErrMsg)
		return false, errors.New(lgc.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader)).Languagef("host_search_module_fail_with_errmsg", err.Error()))
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
func (lgc *Logics) ExistHostIDSInApp(ctx context.Context, appID int64, hostIDArray []int64, header http.Header) ([]int64, error) {
	// v3.3.x TODO lgc.rid
	rid := util.GetHTTPCCRequestID(header)
	// v3.3.x TODO lgc.ccErr
	defErr := lgc.Engine.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	conf := map[string][]int64{
		common.BKAppIDField:  []int64{appID},
		common.BKHostIDField: hostIDArray,
	}

	result, err := lgc.CoreAPI.HostController().Module().GetModulesHostConfig(ctx, header, conf)
	if err != nil {
		blog.Errorf("ExistHostIDSInApp http do error. err:%s, input:%#v,rid:%s", err.Error(), conf, rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !result.Result {
		blog.Errorf("ExistHostIDSInApp http reply error. err code:%d,err msg:%s, input:%#v,rid:%s", result.Code, result.ErrMsg, conf, rid)
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

func (lgc *Logics) GetSingleApp(pheader http.Header, cond interface{}) (mapstr.MapStr, error) {
	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     1,
		Sort:      common.BKAppIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("get app failed, err: %v, %v", err, result.ErrMsg)
	}

	if len(result.Data.Info) == 0 {
		return nil, nil
	}
	return result.Data.Info[0], nil
}

func (lgc *Logics) GetAppIDByCond(pheader http.Header, cond []metadata.ConditionItem) ([]int64, error) {
	condc := make(map[string]interface{})
	params.ParseCommonParams(cond, condc)
	query := &metadata.QueryInput{
		Condition: condc,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKHostIDField,
		Fields:    common.BKAppIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}

	appIDs := make([]int64, 0)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKAppIDField)
		if err != nil {
			return nil, err
		}
		appIDs = append(appIDs, id)
	}

	return appIDs, nil
}

func (lgc *Logics) GetAppMapByCond(pheader http.Header, fields string, cond interface{}) (map[int64]types.MapStr, error) {

	query := &metadata.QueryInput{
		Condition: cond,
		Start:     0,
		Limit:     common.BKNoLimit,
		Sort:      common.BKAppIDField,
		Fields:    fields,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		return nil, fmt.Errorf("%v, %v", err, result.ErrMsg)
	}
	appMap := make(map[int64]types.MapStr)
	for _, info := range result.Data.Info {
		id, err := info.Int64(common.BKAppIDField)
		if err != nil {
			return nil, err
		}
		appMap[id] = info
	}

	return appMap, nil
}
