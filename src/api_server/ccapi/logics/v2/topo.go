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
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	//	"configcenter/src/scene_server/topo_server/topo_service/manager"
)

func (lgc *Logics) GetAppTopo(user string, pheader http.Header, appID int64, conds mapstr.MapStr) (map[string]interface{}, int) {

	apps, errCode := lgc.getAppInfo(user, pheader, appID)
	if 0 != errCode {
		return nil, 0
	}
	if 0 == len(apps) {
		blog.Errorf("GetAppTopo not find app by id:%d", appID)
		return nil, common.CCErrCommNotFound
	}
	appInfo, ok := apps[0].(map[string]interface{})
	if false == ok {
		blog.Errorf("GetAppTopo not find app by id %s, response:%v", apps)
		return nil, common.CCErrCommNotFound
	}

	appName, ok := appInfo[common.BKAppNameField]
	if nil == appName || false == ok {
		appName = ""
	}

	ret := map[string]interface{}{
		"Level":           3,
		"ApplicationName": appName,
		"ApplicationID":   strconv.FormatInt(appID, 10),
		"Children":        make([]interface{}, 0),
	}

	moduleFields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}
	modules, errCode := lgc.getModulesByConds(user, pheader, strconv.FormatInt(appID, 10), conds, moduleFields, "")
	if 0 != errCode {
		return nil, 0
	}
	modulesMap, errCode := getModuleMap(modules, appID)
	if 0 != errCode {
		return nil, 0
	}

	sets, errCode := lgc.getSets(user, pheader, appID)
	if 0 != errCode {
		return nil, 0
	}
	if 0 == len(sets) {
		blog.Errorf("GetAppTopo not find set by app id:%d", appID)
		return nil, common.CCErrCommNotFound
	}
	retSets := make([]map[string]interface{}, 0)
	for _, set := range sets {
		setInfo, ok := set.(map[string]interface{})
		if false == ok {
			blog.Error("GetAppTopo getSets  return info error, set info:%v  error", set)
			return nil, common.CCErrCommHTTPDoRequestFailed
		}

		setID, err := util.GetInt64ByInterface(setInfo[common.BKSetIDField])
		if nil != err {
			blog.Error("GetAppTopo get set id by getsets  return info  error, module info:%v  error", set)
			return nil, common.CCErrCommHTTPDoRequestFailed
		}
		setName, _ := setInfo[common.BKSetNameField]
		if nil == setName {
			setName = ""
		}
		moduleArr, ok := modulesMap[setID]
		if false == ok {
			moduleArr = make([]map[string]interface{}, 0)
		}
		retSets = append(retSets, map[string]interface{}{
			"SetName":  setName,
			"SetID":    strconv.FormatInt(setID, 10),
			"Children": moduleArr,
		})

	}
	ret["Children"] = retSets

	return ret, 0

}

func getModuleMap(modules []interface{}, appID int64) (map[int64][]map[string]interface{}, int) {
	// key is setID
	modulesMap := make(map[int64][]map[string]interface{})
	if 0 == len(modules) {
		blog.Errorf("GetAppTopo not find module by app id:%d", appID)
		return nil, common.CCErrCommNotFound
	}
	strAppID := strconv.FormatInt(appID, 10)
	for _, moduleI := range modules {
		module, ok := moduleI.(map[string]interface{})
		if false == ok {
			blog.Error("GetAppTopo getmodule  return info error, module info:%v  error", module)
			return nil, common.CCErrCommHTTPDoRequestFailed
		}

		setID, err := util.GetInt64ByInterface(module[common.BKSetIDField])
		if nil != err {
			blog.Error("GetAppTopo get set id by getmodues  return info  error, module info:%v  error", module)
			return nil, common.CCErrCommHTTPDoRequestFailed
		}
		moduleName, _ := module[common.BKModuleNameField]
		if nil == moduleName {
			moduleName = ""
		}
		moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
		if nil != err {
			blog.Error("GetAppTopo get module id by getmodues  return info  error, module info:%v  error", module)
			return nil, common.CCErrCommHTTPDoRequestFailed
		}
		_, ok = modulesMap[setID]
		if false == ok {
			modulesMap[setID] = make([]map[string]interface{}, 0)
		}
		modulesMap[setID] = append(modulesMap[setID], map[string]interface{}{
			"SetID":         strconv.FormatInt(setID, 10),
			"ModuleID":      strconv.FormatInt(moduleID, 10),
			"ModuleName":    moduleName,
			"HostNum":       "0",
			"ApplicationID": strAppID,
			"ObjID":         common.BKInnerObjIDModule,
		})

	}

	return modulesMap, 0
}

// GetDefaultTopo get resource topo
func GetDefaultTopo(req *restful.Request, appID string, topoApi string) (map[string]interface{}, error) {
	defaultTopo := make(map[string]interface{})
	url := fmt.Sprintf("%s/topo/v1/topo/internal/%s/%s", topoApi, common.BKDefaultOwnerID, appID)
	res, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
	if err != nil {
		blog.Error("getDefaultTopo error:%v", err)
		return nil, err
	}

	//appIDInt ,_:= strconv.Atoi(appID)
	resMap := make(map[string]interface{})

	err = json.Unmarshal([]byte(res), &resMap)
	if resMap["result"].(bool) {

		resMapData, ok := resMap["data"].(map[string]interface{})
		if false == ok {
			blog.Error("assign error resMap:%v", resMap)
			return defaultTopo, nil
		}
		defaultTopo["Children"] = make([]map[string]interface{}, 0)
		resModule, ok := resMapData["module"].([]interface{})
		if false == ok {
			blog.Error("assign error resMapData:%v", resMapData)
			return defaultTopo, nil
		}
		for _, module := range resModule {
			Module, ok := module.(map[string]interface{})
			if false == ok {
				blog.Error("assign error module:%v", module)
				continue
			}
			moduleMap := map[string]interface{}{
				"ModuleID":      Module[common.BKModuleIDField],
				"ModuleName":    Module[common.BKModuleNameField],
				"ApplicationID": appID,
				"ObjID":         "module",
			}
			defaultTopo["Children"] = append(defaultTopo["Children"].([]map[string]interface{}), moduleMap)
		}
		defaultTopo["SetName"] = common.DefaultResSetName
		setIdInt, _ := util.GetIntByInterface(resMap["data"].(map[string]interface{})[common.BKSetIDField])
		setIdStr := strconv.Itoa(setIdInt)
		defaultTopo["SetID"] = setIdStr
		defaultTopo["ObjID"] = "set"
	}

	blog.Debug("defaultTopo:%v", defaultTopo)
	return defaultTopo, nil
}

// AppendDefaultTopo combin  idle pool set
func AppendDefaultTopo(topo map[string]interface{}, defaultTopo map[string]interface{}) map[string]interface{} {
	topoChildren, ok := topo["Children"].([]map[string]interface{})
	if false == ok {
		err := fmt.Sprintf("assign error topo.Children is not []map[string]interface{},topo:%v", topo)
		blog.Error(err)
		return nil
	}

	children := make([]map[string]interface{}, 0)
	children = append(children, defaultTopo)
	for _, child := range topoChildren {
		children = append(children, child)
	}
	topo["Children"] = children
	return topo
}

// SetModuleHostCount get set host count
func (lgc *Logics) SetModuleHostCount(data []map[string]interface{}, user string, pheader http.Header) error {
	blog.Debug("setModuleHostCount data: %v", data)
	for _, itemMap := range data {
		blog.Debug("ObjID: %s", itemMap)

		switch itemMap["ObjID"] {
		case common.BKInnerObjIDModule:

			mouduleId, getErr := util.GetIntByInterface(itemMap["ModuleID"])
			if nil != getErr {
				blog.Errorf("%v, %v", getErr, itemMap)
				return getErr
			}
			appId, getErr := util.GetIntByInterface(itemMap["ApplicationID"])
			if nil != getErr {
				blog.Errorf("%v, %v", getErr, itemMap)
				return getErr
			}
			blog.Debug("mouduleId: %v", mouduleId)
			hostNum, getErr := lgc.GetModuleHostCount(appId, mouduleId, user, pheader)
			if nil != getErr {
				return getErr
			}
			blog.Debug("hostNum: %v", hostNum)
			itemMap["HostNum"] = hostNum
		}

		if nil != itemMap["Children"] {
			children, ok := itemMap["Children"].([]map[string]interface{})
			if false == ok {
				children = make([]map[string]interface{}, 0)
			}
			lgc.SetModuleHostCount(children, user, pheader)
		} else {
			children := make([]interface{}, 0)
			itemMap["Children"] = children
		}
	}
	return nil
}

// GetModuleHostCount get module host count
func (lgc *Logics) GetModuleHostCount(appID, mouduleID interface{}, user string, pheader http.Header) (int, error) {

	param := mapstr.MapStr{
		common.BKAppIDField:    appID,
		common.BKModuleIDField: []interface{}{mouduleID},
	}

	result, err := lgc.CoreAPI.HostServer().HostSearchByModuleID(context.Background(), pheader, param)
	if nil != err {
		blog.Errorf("getModuleHostCount , error:%v", err)
		return 0, err
	}
	rspV3MapData, ok := result.Data.([]interface{})
	if false == ok {
		blog.Error("assign error rspV3Map.data is not []interface{}, rspV3Map:%v", result.Data)
		return 0, nil
	}
	return len(rspV3MapData), nil
}

func (lgc *Logics) getSets(user string, pheader http.Header, appID int64) ([]interface{}, int) {

	page := mapstr.MapStr{"start": 0, "limit": common.BKNoLimit}
	param := &params.SearchParams{Page: page}
	appIDStr := strconv.FormatInt(appID, 10)
	result, err := lgc.CoreAPI.TopoServer().Instance().SearchSet(context.Background(), user, appIDStr, pheader, param)
	if err != nil {
		blog.Error("get sets   error:%v", err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}
	return getRspV3DataInfo("getSets", result.Result, result.Code, result.Data)
}

func (lgc *Logics) getModules(user string, pheader http.Header, appID int64) ([]interface{}, int) {

	fields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}

	return lgc.getModulesByConds(user, pheader, strconv.FormatInt(appID, 10), nil, fields, "")
}

func (lgc *Logics) getModulesByConds(user string, pheader http.Header, appIDStr string, conds map[string]interface{}, fields []string, sort string) ([]interface{}, int) {

	searchParams := mapstr.MapStr{}
	searchParams["fields"] = fields
	searchParams["condition"] = conds
	searchParams["page"] = mapstr.MapStr{"start": 0, "limit": common.BKNoLimit, "sort": sort}
	result, err := lgc.CoreAPI.TopoServer().OpenAPI().SearchModuleByApp(context.Background(), appIDStr, pheader, searchParams)
	if err != nil {
		blog.Error("getModules   error:%v", err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}
	return getRspV3DataInfo("getModules", result.Result, result.Code, result.Data)
}

func (lgc *Logics) getAppInfo(user string, pheader http.Header, appID int64) ([]interface{}, int) {

	page := mapstr.MapStr{"start": 0, "limit": 2}
	condition := mapstr.MapStr{common.BKAppIDField: appID}
	param := &params.SearchParams{Condition: condition, Page: page}

	result, err := lgc.CoreAPI.TopoServer().Instance().SearchApp(context.Background(), user, pheader, param)

	if nil != err {
		blog.Errorf("get app info error:%v", err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}

	return getRspV3DataInfo("getAppInfo", result.Result, result.Code, result.Data)
}

func getRspV3DataInfo(logPrex string, result bool, code int, data interface{}) ([]interface{}, int) {

	if false == result {
		return nil, code
	}
	dataMap, ok := data.(map[string]interface{})
	if false == ok {
		blog.Errorf("%s rspV3 json get data.info error, body:%s  error:%s", logPrex, dataMap)
		return nil, common.CCErrCommJSONUnmarshalFailed
	}
	dataInfo, ok := dataMap["info"].([]interface{})
	if false == ok {
		blog.Errorf("%s rspV3 json get data.info error, body:%s  error:%s", logPrex, dataMap)
		return nil, common.CCErrCommJSONUnmarshalFailed
	}
	return dataInfo, 0
}

func (lgc *Logics) CheckAppTopoIsThreeLevel(user string, pheader http.Header) (bool, error) {

	result, err := lgc.CoreAPI.TopoServer().Object().SelectModel(context.Background(), util.GetOwnerID(pheader), pheader)
	if nil != err {
		blog.Errorf("CheckAppTopoIsThreeLevel  error:%s ", err.Error())
		return false, err
	}
	if false == result.Result {
		blog.Errorf("CheckAppTopoIsThreeLevel reply:%s, error:%s ", result.ErrMsg)
		return false, fmt.Errorf(result.ErrMsg)
	}

	for _, item := range result.Data {
		if common.BKInnerObjIDApp == item.ObjID {
			if common.BKInnerObjIDSet != item.NextObj {
				return false, nil
			}
		}

		if common.BKInnerObjIDSet == item.ObjID {
			if common.BKInnerObjIDModule != item.NextObj {
				return false, nil
			}
		}
	}

	return true, nil
}
