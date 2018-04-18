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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

func GetAppTopo(req *restful.Request, cc *api.APIResource, ownerID string, appID int64) (map[string]interface{}, int) {
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

	modules, errCode := getModules(req, cc, appID)
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
		fmt.Println("moduleName =====", moduleName, module)
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
	url := cc.TopoAPI() + fmt.Sprintf("/topo/v1/openapi/module/searchByApp/%d", appID)
	var js params.SearchParams
	js.Fields = []string{common.BKModuleIDField, common.BKModuleNameField, common.BKAppIDField, common.BKSetIDField, common.BKDefaultField}
	js.Condition = common.KvMap{common.BKAppIDField: appID}
	js.Page = common.KvMap{"start": 0, "limit": common.BKNoLimit}

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
