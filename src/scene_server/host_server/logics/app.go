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
	"configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/service"
)

func (lgc *Logics) GetDefaultAppIDWithSupplier(supplierID int64, pheader http.Header) (int64, error) {
	cond := service.NewOperation().WithDefaultField(int64(common.DefaultAppFlag)).WithSupplierID(supplierID).Data()
	appDetails, err := lgc.GetAppDetails(common.BKAppIDField, cond, pheader)
	if err != nil {
		return -1, err
	}

	id, exist := appDetails[common.BKAppIDField].(int64)
	if !exist {
		return -1, errors.New("can not find bk biz field")
	}
	return id, nil
}

func (lgc *Logics) GetDefaultAppID(ownerID string, pheader http.Header) (int64, error) {
	cond := service.NewOperation().WithOwnerID(ownerID).WithDefaultField(int64(common.DefaultAppFlag)).Data()
	appDetails, err := lgc.GetAppDetails(common.BKAppIDField, cond, pheader)
	if err != nil {
		return -1, err
	}

	id, exist := appDetails[common.BKAppIDField].(int64)
	if !exist {
		return -1, errors.New("can not find bk biz field")
	}
	return id, nil
}

func (lgc *Logics) GetAppDetails(fields string, condition map[string]interface{}, pheader http.Header) (map[string]interface{}, error) {
	query := metadata.QueryInput{
		Condition: condition,
		Start:     0,
		Limit:     1,
		Fields:    fields,
		Sort:      common.BKAppIDField,
	}

	result, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKAppIDField, pheader, &query)
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

func (lgc *Logics) GetAppIDByCond(pheader http.Header, cond []interface{}) ([]int64, error) {
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
