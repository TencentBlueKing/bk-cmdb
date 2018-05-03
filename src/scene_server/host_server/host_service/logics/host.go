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
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/language"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/instapi"
	"encoding/json"
	"errors"
	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
	"strings"
)

// HostSearch search host by mutiple condition
func HostSearch(req *restful.Request, data hostParse.HostCommonSearch, hostCtrl, objCtrl string) (interface{}, error) {
	var hostCond, appCond, setCond, moduleCond, mainlineCond hostParse.SearchCondition
	objectCondMap := make(map[string][]interface{}, 0)
	appIDArr := make([]int, 0)
	setIDArr := make([]int, 0)
	moduleIDArr := make([]int, 0)
	hostIDArr := make([]int, 0)
	instAsstHostIDArr := make([]int, 0)
	objSetIDArr := make([]int, 0)
	disAppIDArr := make([]int, 0)
	disSetIDArr := make([]int, 0)
	disModuleIDArr := make([]int, 0)
	hostAppConfig := make(map[int]int)
	hostSetConfig := make(map[int][]int)
	hostModuleConfig := make(map[int][]int)
	hostModuleMap := make(map[int]interface{})
	hostSetMap := make(map[int]interface{})
	hostAppMap := make(map[int]interface{})

	result := make(map[string]interface{})
	totalInfo := make([]interface{}, 0)
	moduleHostConfig := make(map[string][]int, 0)

	url := hostCtrl + "/host/v1/hosts/search"
	start := data.Page.Start
	limit := data.Page.Limit
	sort := data.Page.Sort
	body := make(map[string]interface{})
	body["start"] = start
	body["limit"] = limit
	body["sort"] = sort
	for _, object := range data.Condition {
		if object.ObjectID == common.BKInnerObjIDHost {
			hostCond = object
		} else if object.ObjectID == common.BKInnerObjIDSet {
			setCond = object
		} else if object.ObjectID == common.BKInnerObjIDModule {
			moduleCond = object
		} else if object.ObjectID == common.BKInnerObjIDApp {
			appCond = object
		} else if object.ObjectID == common.BKINnerObjIDObject {
			mainlineCond = object
		} else {
			objectCondMap[object.ObjectID] = object.Condition
		}
	}

	//search appID by cond
	if -1 != data.AppID && 0 != data.AppID {
		cond := make(map[string]interface{})
		cond["field"] = common.BKAppIDField
		cond["operator"] = common.BKDBEQ
		cond["value"] = data.AppID
		appCond.Condition = append(appCond.Condition, cond)
	}
	if len(appCond.Condition) > 0 {
		appIDArr, _ = GetAppIDByCond(req, objCtrl, appCond.Condition)
	}
	//search mainline object by cond
	if len(mainlineCond.Condition) > 0 {
		objSetIDArr = GetSetIDByObjectCond(req, objCtrl, mainlineCond.Condition)
	}
	//search set by appcond
	if len(setCond.Condition) > 0 || len(mainlineCond.Condition) > 0 {
		if len(appCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKAppIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = appIDArr
			setCond.Condition = append(setCond.Condition, cond)
		}
		if len(mainlineCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKSetIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = objSetIDArr
			setCond.Condition = append(setCond.Condition, cond)
		}
		setIDArr, _ = GetSetIDByCond(req, objCtrl, setCond.Condition)
	}

	//search host id by object
	if len(objectCondMap) > 0 {
		for objID, objCond := range objectCondMap {
			instIDArr := GetObjectInstByCond(req, objID, objCtrl, objCond)
			instAsstHostIDArr = GetHostIDByInstID(req, objID, objCtrl, instIDArr)
		}

	}
	if len(moduleCond.Condition) > 0 {
		if len(setCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKSetIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = setIDArr
			moduleCond.Condition = append(moduleCond.Condition, cond)
		}
		if len(appCond.Condition) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKAppIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = appIDArr
			moduleCond.Condition = append(moduleCond.Condition, cond)
		}
		//search module by cond
		moduleIDArr, _ = GetModuleIDByCond(req, objCtrl, moduleCond.Condition)
	}

	if len(appCond.Condition) > 0 {
		moduleHostConfig[common.BKAppIDField] = appIDArr
	}
	if len(setCond.Condition) > 0 {
		moduleHostConfig[common.BKSetIDField] = setIDArr
	}
	if len(moduleCond.Condition) > 0 {
		moduleHostConfig[common.BKModuleIDField] = moduleIDArr
	}
	if len(objectCondMap) > 0 {
		moduleHostConfig[common.BKHostIDField] = instAsstHostIDArr
	}
	hostIDArr, _ = GetHostIDByCond(req, hostCtrl, moduleHostConfig)

	if len(appCond.Condition) > 0 || len(setCond.Condition) > 0 || len(moduleCond.Condition) > 0 || -1 != data.AppID {
		cond := make(map[string]interface{})
		cond["field"] = common.BKHostIDField
		cond["operator"] = common.BKDBIN
		cond["value"] = hostIDArr
		hostCond.Condition = append(hostCond.Condition, cond)
	}
	if 0 != len(hostCond.Fields) {
		hostCond.Fields = append(hostCond.Fields, common.BKHostIDField)
	}
	body["fields"] = strings.Join(hostCond.Fields, ",")
	condition := make(map[string]interface{})
	hostParse.ParseHostParams(hostCond.Condition, condition)
	hostParse.ParseHostIPParams(data.Ip, condition)
	body["condition"] = condition
	bodyContent, _ := json.Marshal(body)
	blog.Info("Get Host By Cond url :%s", url)
	blog.Info("Get Host By Cond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("Get Host By Cond return :%s", string(reply))
	if err != nil {
		//cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return nil, errors.New(common.CC_Err_Comm_Host_Get_FAIL_STR)
	}

	js, err := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	hostData := output["data"]
	hostResult, ok := hostData.(map[string]interface{})
	if false == ok {
		//cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return nil, errors.New(common.CC_Err_Comm_Host_Get_FAIL_STR)
	}

	// deal the host
	instapi.Inst.InitInstHelper(hostCtrl, objCtrl)
	hostResult, retStrErr := instapi.Inst.GetInstDetailsSub(req, common.BKInnerObjIDHost, common.BKDefaultOwnerID, hostResult, map[string]interface{}{
		"start": 0,
		"limit": common.BKNoLimit,
		"sort":  "",
	})

	if common.CCSuccess != retStrErr {
		blog.Error("failed to replace association object, error code is %d", retStrErr)
	}

	cnt := hostResult["count"]
	hostInfo := hostResult["info"].([]interface{})
	result["count"] = cnt
	resHostIDArr := make([]int, 0)
	queryCond := make(map[string]interface{})
	for _, j := range hostInfo {
		host := j.(map[string]interface{})
		hostID, _ := host[common.BKHostIDField].(json.Number).Int64()
		resHostIDArr = append(resHostIDArr, int(hostID))

		queryCond[common.BKHostIDField] = resHostIDArr
	}
	mhconfig, _ := GetConfigByCond(req, hostCtrl, queryCond)
	blog.Info("get modulehostconfig map:%v", mhconfig)
	for _, mh := range mhconfig {
		hostID := mh[common.BKHostIDField]
		hostAppConfig[hostID] = mh[common.BKAppIDField]
		hostSetConfig[hostID] = append(hostSetConfig[hostID], mh[common.BKSetIDField])
		hostModuleConfig[hostID] = append(hostModuleConfig[hostID], mh[common.BKModuleIDField])
		disAppIDArr = append(disAppIDArr, mh[common.BKAppIDField])
		disSetIDArr = append(disSetIDArr, mh[common.BKSetIDField])
		disModuleIDArr = append(disModuleIDArr, mh[common.BKModuleIDField])
	}
	if nil != appCond.Fields {
		//get app fields

		exist := util.InArray(common.BKAppIDField, appCond.Fields)
		if 0 != len(appCond.Fields) && !exist {
			appCond.Fields = append(appCond.Fields, common.BKAppIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disAppIDArr
		cond[common.BKAppIDField] = celld
		fields := strings.Join(appCond.Fields, ",")
		hostAppMap, _ = GetAppMapByCond(req, fields, objCtrl, cond)
	}
	if nil != setCond.Fields {
		//get set fields

		exist := util.InArray(common.BKSetIDField, setCond.Fields)
		if !exist && 0 != len(setCond.Fields) {
			setCond.Fields = append(setCond.Fields, common.BKSetIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disSetIDArr
		cond[common.BKSetIDField] = celld
		fields := strings.Join(setCond.Fields, ",")
		hostSetMap, _ = GetSetMapByCond(req, fields, objCtrl, cond)
	}
	if nil != moduleCond.Fields {
		//get module fields

		exist := util.InArray(common.BKModuleIDField, moduleCond.Fields)
		if !exist && 0 != len(moduleCond.Fields) {
			moduleCond.Fields = append(moduleCond.Fields, common.BKModuleIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disModuleIDArr
		cond[common.BKModuleIDField] = celld
		fields := strings.Join(moduleCond.Fields, ",")
		hostModuleMap, _ = GetModuleMapByCond(req, fields, objCtrl, cond)
	}

	//com host info
	for _, j := range hostInfo {
		host := j.(map[string]interface{})
		hostID, _ := host[common.BKHostIDField].(json.Number).Int64()
		hostID32 := int(hostID)
		hostData := make(map[string]interface{})
		//appdata
		appInfo, ok := hostAppMap[hostAppConfig[hostID32]]
		if ok {
			hostData[common.BKInnerObjIDApp] = appInfo
		} else {
			hostData[common.BKInnerObjIDApp] = make(map[string]interface{})
		}
		//setdata
		hostSetIDArr, ok := hostSetConfig[hostID32]
		if ok && nil != setCond.Fields {
			setNameArr := make([]string, 0)
			for _, setID := range hostSetIDArr {
				setInfo, ok := hostSetMap[setID]
				if false == ok {
					continue
				}
				data, ok := setInfo.(map[string]interface{})
				if false == ok {
					continue
				}
				setName, ok := data[common.BKSetNameField].(string)
				if false == ok {
					continue
				}
				setNameArr = append(setNameArr, setName)
			}
			setNameStr := strings.Join(util.StrArrayUnique(setNameArr), ",")
			hostData[common.BKInnerObjIDSet] = map[string]string{common.BKSetNameField: setNameStr}
		} else {
			hostData[common.BKInnerObjIDSet] = make(map[string]interface{})
		}
		//moduledata
		hostModuleIDArr, ok := hostModuleConfig[hostID32]
		if ok && nil != moduleCond.Fields {
			moduleNameArr := make([]string, 0)
			for _, setID := range hostModuleIDArr {
				moduleInfo, ok := hostModuleMap[setID]
				if false == ok {
					continue
				}
				data, ok := moduleInfo.(map[string]interface{})
				if false == ok {
					continue
				}
				moduleName, ok := data[common.BKModuleNameField].(string)
				if false == ok {
					continue
				}
				moduleNameArr = append(moduleNameArr, moduleName)
			}
			moduleNameStr := strings.Join(util.StrArrayUnique(moduleNameArr), ",")
			hostData[common.BKInnerObjIDModule] = map[string]string{common.BKModuleNameField: moduleNameStr}
		} else {
			hostData[common.BKInnerObjIDModule] = make(map[string]interface{})
		}

		hostData[common.BKInnerObjIDHost] = j
		totalInfo = append(totalInfo, hostData)
	}

	result["info"] = totalInfo
	result["count"] = cnt

	return result, err
}

func GetHostInfoByConds(req *restful.Request, hostURL string, conds map[string]interface{}, defLang language.DefaultCCLanguageIf) ([]interface{}, error) {
	hostURL = hostURL + "/host/v1/hosts/search"
	getParams := make(map[string]interface{})
	getParams["fields"] = nil
	getParams["condition"] = conds
	getParams["start"] = 0
	getParams["limit"] = common.BKNoLimit
	getParams["sort"] = common.BKHostIDField
	blog.Info("get host info by conds url:%s", hostURL)
	blog.Info("get host info by conds params:%v", getParams)
	isSucess, message, iRetData := GetHttpResult(req, hostURL, common.HTTPSelectPost, getParams)
	blog.Info("get host info by conds return:%v", iRetData)
	if !isSucess {
		msg := defLang.Languagef("host_search_fail_with_errmsg", message)
		blog.Error(msg)
		return nil, errors.New(msg)
	}
	if nil == iRetData {
		return nil, nil
	}
	retData := iRetData.(map[string]interface{})
	data, _ := retData["info"]
	if nil == data {
		return nil, nil
	}
	return data.([]interface{}), nil
}
