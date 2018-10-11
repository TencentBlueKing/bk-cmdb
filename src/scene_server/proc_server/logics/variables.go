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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type Variables struct {
	logic   *Logics
	header  http.Header
	ownerID string
	appID   int64
}

func (lgc *Logics) NewVariables(pheader http.Header, ownerID string, appID int64) *Variables {
	return &Variables{
		logic:   lgc,
		header:  pheader,
		ownerID: ownerID,
		appID:   appID,
	}
}

func (v *Variables) GetAppVariables() types.MapStr {
	data := types.MapStr{}
	cond := types.MapStr{common.BKAppIDField: v.appID}

	input := metadata.QueryInput{Condition: cond}
	result, err := v.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, v.header, &input)

	if nil != err || !result.Result || 0 == len(result.Data.Info) {
		blog.Errorf("get app info use template variables, error %v. error code:%d, error message:%s", err, result.Code, result.ErrMsg)
		return data
	}

	return result.Data.Info[0]
}

func (v *Variables) GetSetVariables(setName string) (int64, types.MapStr) {
	data := types.MapStr{}
	cond := types.MapStr{common.BKAppIDField: v.appID,
		common.BKSetNameField: setName}

	input := metadata.QueryInput{Condition: cond}
	result, err := v.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, v.header, &input)

	if nil != err || !result.Result || 0 == len(result.Data.Info) {
		return 0, data
	}

	setInfo := result.Data.Info[0]
	setID, err := util.GetInt64ByInterface(setInfo[common.BKSetIDField])
	if nil != err {
		return 0, data
	}
	return setID, setInfo
}

func (v *Variables) GetModuleVariables(setID int64, moduleName string) (int64, types.MapStr) {
	data := types.MapStr{}
	cond := types.MapStr{common.BKAppIDField: v.appID,
		common.BKSetIDField:   setID,
		common.BKSetNameField: moduleName}

	input := metadata.QueryInput{Condition: cond}
	result, err := v.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, v.header, &input)

	if nil != err || !result.Result || 0 == len(result.Data.Info) {
		blog.Errorf("get data error %v", err)
		return 0, data
	}

	moduleInfo := result.Data.Info[0]
	moduleID, err := util.GetInt64ByInterface(moduleInfo[common.BKSetIDField])
	if nil != err {
		return 0, data
	}
	return moduleID, moduleInfo
}

func (v *Variables) GetHostVariables(HostID int64) types.MapStr {
	data := types.MapStr{}
	cond := types.MapStr{common.BKHostIDField: HostID}

	input := metadata.QueryInput{Condition: cond}
	result, err := v.logic.CoreAPI.HostController().Host().GetHosts(context.Background(), v.header, &input)

	if nil != err || !result.Result || 0 == len(result.Data.Info) {
		blog.Errorf("get data error %v", err)
		return data
	}

	return result.Data.Info[0]
}

func (v *Variables) GetProcessVariables(funcID int64) types.MapStr {
	data := types.MapStr{}
	cond := types.MapStr{common.BKAppIDField: v.appID, common.BKFuncIDField: funcID}

	input := metadata.QueryInput{Condition: cond}
	result, err := v.logic.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDProc, v.header, &input)

	if nil != err || !result.Result || 0 == len(result.Data.Info) {
		return data
	}

	return result.Data.Info[0]
}

func (v *Variables) GetHostIDByInst(setID, moduleID, funcID, instID int64) int64 {
	return 0
}

func (v *Variables) GetStandVariables(setName, moduleName string, funcID, instID int64) types.MapStr {
	allVariables := types.MapStr{}
	appVariables := v.GetAppVariables()
	if 0 == len(appVariables) {
		return allVariables
	}
	setID, setVariables := v.GetSetVariables(setName)
	if 0 == len(setVariables) || 0 == setID {
		return allVariables
	}

	moduleID, moduleVariables := v.GetModuleVariables(setID, moduleName)
	if 0 == moduleID || 0 == len(moduleVariables) {
		return allVariables
	}
	hostID := v.GetHostIDByInst(setID, moduleID, funcID, instID)
	if 0 == hostID {
		return allVariables
	}

	hostVariables := v.GetHostVariables(hostID)
	if 0 == len(hostVariables) {
		return allVariables
	}

	procVariables := v.GetProcessVariables(hostID)
	if 0 == len(procVariables) {
		return allVariables
	}

	for key, val := range procVariables {
		allVariables[key] = val
	}

	for key, val := range hostVariables {
		allVariables[key] = val
	}

	for key, val := range moduleVariables {
		allVariables[key] = val
	}

	for key, val := range setVariables {
		allVariables[key] = val
	}

	for key, val := range appVariables {
		allVariables[key] = val
	}

	return allVariables
}
