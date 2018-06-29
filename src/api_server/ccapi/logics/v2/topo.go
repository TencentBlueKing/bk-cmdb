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
	"encoding/json"
	"fmt"
	"strconv"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/topo_service/manager"
)

func GetAppTopo(req *restful.Request, cc *api.APIResource, ownerID string, appID int64, conds common.KvMap) (map[string]interface{}, int) {
	apps, errCode := getAppInfo(req, cc, appID)
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
		"ApplicationID":   fmt.Sprintf("%d", appID),
		"Children":        make([]interface{}, 0),
	}

	moduleFields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}
	modules, errCode := getModulesByConds(req, cc, appID, conds, moduleFields, "")
	if 0 != errCode {
		return nil, 0
	}
	modulesMap, errCode := getModuleMap(modules, appID)
	if 0 != errCode {
		return nil, 0
	}

	sets, errCode := getSets(req, cc, ownerID, appID)
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
			"SetID":    fmt.Sprintf("%d", setID),
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
	strAppID := fmt.Sprintf("%d", appID)
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
			"SetID":         fmt.Sprintf("%d", setID),
			"ModuleID":      fmt.Sprintf("%d", moduleID),
			"ModuleName":    moduleName,
			"HostNum":       "0",
			"ApplicationID": strAppID,
			"ObjID":         common.BKInnerObjIDModule,
		})

	}

	return modulesMap, 0
}

// GetDefaultTopo 获取空闲机池topo
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
func SetModuleHostCount(data []map[string]interface{}, req *restful.Request, cc *api.APIResource) error {
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
			hostNum, getErr := GetModuleHostCount(appId, mouduleId, req, cc)
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
			SetModuleHostCount(children, req, cc)
		} else {
			children := make([]interface{}, 0)
			itemMap["Children"] = children
		}
	}
	return nil
}

// GetModuleHostCount get module host count
func GetModuleHostCount(appID, mouduleID interface{}, req *restful.Request, cc *api.APIResource) (int, error) {

	param := map[string]interface{}{
		common.BKAppIDField:    appID,
		common.BKModuleIDField: []interface{}{mouduleID},
	}
	paramJson, err := json.Marshal(param)
	if err != nil {
		blog.Error("getModuleHostCount Marshal json error:%v", err)
		return 0, nil
	}

	url := fmt.Sprintf("%s/host/v1/getmodulehostlist", cc.HostAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(paramJson))
	blog.Debug(rspV3, err)
	if nil != err {
		blog.Errorf("getModuleHostCount url:%s, params:%s, error:%s", url, string(paramJson), err.Error())
		return 0, err
	}
	rspV3Map := make(map[string]interface{})
	err = json.Unmarshal([]byte(rspV3), &rspV3Map)
	if nil != err {
		blog.Error("getmodulehostlist Unmarshal json error:%v, rspV3:%s", err, rspV3)
		return 0, nil
	}
	rspV3MapData, ok := rspV3Map["data"].([]interface{})
	if false == ok {
		blog.Error("assign error rspV3Map.data is not []interface{}, rspV3Map:%v", rspV3Map)
		return 0, nil
	}
	return len(rspV3MapData), nil
}

func getSets(req *restful.Request, cc *api.APIResource, ownerID string, appID int64) ([]interface{}, int) {
	url := cc.TopoAPI() + fmt.Sprintf("/topo/v1/set/search/%s/%d", ownerID, appID)
	var js params.SearchParams
	js.Page = common.KvMap{"start": 0, "limit": common.BKNoLimit}

	inputBody, err := json.Marshal(js)
	if nil != err {
		return nil, common.CCErrCommJSONMarshalFailed
	}

	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, inputBody)
	if err != nil {
		blog.Error("getSets  url:%s, params:%s error:%v", url, string(inputBody), err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}
	blog.Infof("getSets url:%s, reply:%s ", url, rspV3)
	return getRspV3DataInfo("getSets", rspV3)
}

func getModules(req *restful.Request, cc *api.APIResource, appID int64) ([]interface{}, int) {

	fields := []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}

	return getModulesByConds(req, cc, appID, nil, fields, "")
}

func getModulesByConds(req *restful.Request, cc *api.APIResource, appID int64, conds map[string]interface{}, fields []string, sort string) ([]interface{}, int) {
	url := cc.TopoAPI() + fmt.Sprintf("/topo/v1/openapi/module/searchByApp/%d", appID)
	var js params.SearchParams
	js.Fields = fields
	js.Condition = conds
	js.Page = common.KvMap{"start": 0, "limit": common.BKNoLimit, "sort": sort}

	inputBody, err := json.Marshal(js)
	if nil != err {
		return nil, common.CCErrCommJSONMarshalFailed
	}

	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, inputBody)
	if err != nil {
		blog.Error("getModules  url:%s, params:%s error:%v", url, string(inputBody), err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}
	blog.Infof("getModules url:%s, reply:%s ", url, rspV3)
	return getRspV3DataInfo("getModules", rspV3)
}

func getAppInfo(req *restful.Request, cc *api.APIResource, appID int64) ([]interface{}, int) {
	// build v3 parameters

	var js params.SearchParams
	js.Page = common.KvMap{"start": 0, "limit": common.BKNoLimit}

	js.Condition = map[string]interface{}{
		common.BKAppIDField: appID,
	}
	js.Page = common.KvMap{
		"limit": 2,
		"start": 0,
	}
	inputBody, err := json.Marshal(js)
	if nil != err {
		return nil, common.CCErrCommJSONMarshalFailed
	}
	url := fmt.Sprintf("%s/topo/v1/app/search/"+common.BKDefaultOwnerID, cc.TopoAPI())
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, inputBody)
	blog.Infof("getAppInfo url:%s, input:%s ", url, rspV3)
	if err != nil {
		blog.Error("getAppInfo  url:%s, params:%s error:%v", url, string(inputBody), err)
		return nil, common.CCErrCommHTTPDoRequestFailed
	}
	return getRspV3DataInfo("getAppInfo", rspV3)
}

func getRspV3DataInfo(logPrex, rspV3 string) ([]interface{}, int) {

	js, err := simplejson.NewJson([]byte(rspV3))
	if nil != err {
		blog.Infof("%s rspV3 new json body:%s  error:%s", logPrex, rspV3, err.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	result, err := js.Get("result").Bool()

	if nil != err {
		blog.Infof("%s rspV3 json get result error, body:%s  error:%s", logPrex, rspV3, err.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	if false == result {
		errCode, err := js.Get(common.HTTPBKAPIErrorCode).Int()
		if nil != err {
			blog.Infof("%s rspV3 json get result errorcode error, body:%s  error:%s", logPrex, rspV3, err.Error())
			return nil, common.CCErrCommJSONUnmarshalFailed
		}

		return nil, errCode
	}
	data, err := js.Get("data").Get("info").Array()
	if nil != err {
		blog.Infof("%s rspV3 json get data.info error, body:%s  error:%s", logPrex, rspV3, err.Error())
		return nil, common.CCErrCommJSONUnmarshalFailed
	}

	return data, 0
}

func CheckAppTopoIsThreeLevel(req *restful.Request, cc *api.APIResource) (bool, error) {
	type mainLineItem struct {
		Result bool                   `json:"result"`
		Code   int                    `json:"bk_error_code"`
		ErrMsg string                 `json:"bk_error_msg"`
		Data   []manager.TopoModelRsp `json:"data"`
	}
	url := fmt.Sprintf("%s/topo/v1/model/%s", cc.TopoAPI(), util.GetActionOnwerID(req))
	rspV3, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
	blog.V(0).Infof("getAppInfo url:%s, input:%s ", url, rspV3)
	if nil != err {
		blog.Errorf("CheckAppTopoIsThreeLevel url:%s, error:%s ", url, err.Error())
		return false, err
	}

	resMainLine := &mainLineItem{}
	err = json.Unmarshal([]byte(rspV3), resMainLine)
	if nil != err {
		blog.Errorf("CheckAppTopoIsThreeLevel reply:%s, error:%s ", rspV3, err.Error())
		return false, fmt.Errorf("check mainline topology reply:%s", rspV3)
	}
	if false == resMainLine.Result {
		blog.Errorf("CheckAppTopoIsThreeLevel reply:%s, error:%s ", rspV3, err.Error())
		return false, fmt.Errorf(resMainLine.ErrMsg)
	}

	for _, item := range resMainLine.Data {
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
