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
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/instapi"
	"encoding/json"
	"strings"

	simplejson "github.com/bitly/go-simplejson"
	restful "github.com/emicklei/go-restful"
)

// NewHostSearch new search host by mutiple condition
func NewHostSearch(req *restful.Request, data hostParse.NewHostCommonSearch, hostCtrl, objCtrl string) (interface{}, error) {
	conditionMap := data.Condition
	fieldMap := data.Field

	var hostCond, appCond, setCond, moduleCond, objectCond map[string]interface{}
	var hostField, appField, setField, moduleField []string

	setCondArr := make([]interface{}, 0)
	appCondArr := make([]interface{}, 0)
	hostCondArr := make([]interface{}, 0)
	moduleCondArr := make([]interface{}, 0)
	objectCondArr := make([]interface{}, 0)

	appIDArr := make([]int, 0)
	setIDArr := make([]int, 0)
	moduleIDArr := make([]int, 0)
	hostIDArr := make([]int, 0)
	objSetIDArr := make([]int, 0)
	moduleHostConfig := make(map[string][]int, 0)

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

	url := hostCtrl + "/host/v1/hosts/search"
	start := data.Page.Start
	limit := data.Page.Limit
	sort := data.Page.Sort
	body := make(map[string]interface{})
	body["start"] = start
	body["limit"] = limit
	body["sort"] = sort

	for key, object := range conditionMap {
		if key == common.BKInnerObjIDSet {
			setCond = object
		} else if key == common.BKInnerObjIDApp {
			appCond = object
		} else if key == common.BKInnerObjIDHost {
			hostCond = object
		} else if key == common.BKInnerObjIDModule {
			moduleCond = object
		} else {
			objectCond = object
		}
	}

	for key, object := range fieldMap {
		if key == common.BKInnerObjIDSet {
			setField = object
		} else if key == common.BKInnerObjIDApp {
			appField = object
		} else if key == common.BKInnerObjIDHost {
			hostField = object
		} else if key == common.BKInnerObjIDModule {
			moduleField = object
		}
	}

	if len(setCond) > 0 {
		setCondArr, _ = hostParse.ParseCondParams(setCond)
	}

	if len(appCond) > 0 {
		appCondArr, _ = hostParse.ParseCondParams(appCond)
	}

	if len(hostCond) > 0 {
		hostCondArr, _ = hostParse.ParseCondParams(hostCond)
	}

	if len(moduleCond) > 0 {
		moduleCondArr, _ = hostParse.ParseCondParams(moduleCond)
	}

	if len(objectCond) > 0 {
		objectCondArr, _ = hostParse.ParseCondParams(objectCond)
	}

	if len(objectCondArr) > 0 {
		objSetIDArr = GetSetIDByObjectCond(req, objCtrl, objectCondArr)
	}

	if len(appCondArr) > 0 {
		appIDArr, _ = GetAppIDByCond(req, objCtrl, appCondArr)
	}

	if len(setCondArr) > 0 || len(objectCondArr) > 0 {
		if len(appCondArr) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKAppIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = appIDArr
			setCondArr = append(setCondArr, cond)
		}
		if len(objectCondArr) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKSetIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = objSetIDArr
			setCondArr = append(setCondArr, cond)
		}
		setIDArr, _ = GetSetIDByCond(req, objCtrl, setCondArr)
	}

	if len(moduleCondArr) > 0 {
		if len(setCondArr) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKSetIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = setIDArr
			moduleCondArr = append(moduleCondArr, cond)
		}
		if len(appCondArr) > 0 {
			cond := make(map[string]interface{})
			cond["field"] = common.BKAppIDField
			cond["operator"] = common.BKDBIN
			cond["value"] = appIDArr
			moduleCondArr = append(moduleCondArr, cond)
		}
		//search module by cond
		moduleIDArr, _ = GetModuleIDByCond(req, objCtrl, moduleCondArr)
	}

	if len(appCondArr) > 0 {
		moduleHostConfig[common.BKAppIDField] = appIDArr
	}
	if len(setCondArr) > 0 {
		moduleHostConfig[common.BKSetIDField] = setIDArr
	}
	if len(moduleCondArr) > 0 {
		moduleHostConfig[common.BKModuleIDField] = moduleIDArr
	}

	hostIDArr, _ = GetHostIDByCond(req, hostCtrl, moduleHostConfig)

	if len(appCondArr) > 0 || len(setCondArr) > 0 || len(moduleCondArr) > 0 {
		cond := make(map[string]interface{})
		cond["field"] = common.BKHostIDField
		cond["operator"] = common.BKDBIN
		cond["value"] = hostIDArr
		hostCondArr = append(hostCondArr, cond)
	}

	hostField = append(hostField, common.BKHostIDField)
	body["fields"] = strings.Join(hostField, ",")
	condition := make(map[string]interface{})
	hostParse.ParseHostParams(hostCondArr, condition)
	hostParse.ParseHostIPParams(data.IP, condition)
	body["condition"] = condition
	bodyContent, _ := json.Marshal(body)
	blog.Info("Get Host By Cond url :%s", url)
	blog.Info("Get Host By Cond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("Get Host By Cond return :%s", string(reply))
	if err != nil {
		return nil, err
	}

	js, err := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()

	hostData := output["data"]
	hostResult, ok := hostData.(map[string]interface{})
	if false == ok {
		//cli.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, common.CC_Err_Comm_Host_Get_FAIL_STR, resp)
		return nil, err
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
	if len(appCondArr) > 0 {
		//get app fields

		exist := util.InArray(common.BKAppIDField, appCondArr)
		if 0 != len(appCondArr) && !exist {
			appField = append(appField, common.BKAppIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disAppIDArr
		cond[common.BKAppIDField] = celld
		fields := strings.Join(appField, ",")
		hostAppMap, _ = GetAppMapByCond(req, fields, objCtrl, cond)
	}
	if len(setField) > 0 {
		//get set fields

		exist := util.InArray(common.BKSetIDField, setField)
		if !exist && 0 != len(setField) {
			setField = append(setField, common.BKSetIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disSetIDArr
		cond[common.BKSetIDField] = celld
		fields := strings.Join(setField, ",")
		hostSetMap, _ = GetSetMapByCond(req, fields, objCtrl, cond)
	}
	if nil != moduleField {
		//get module fields

		exist := util.InArray(common.BKModuleIDField, moduleField)
		if !exist && 0 != len(moduleField) {
			moduleField = append(moduleField, common.BKModuleIDField)
		}
		cond := make(map[string]interface{})
		celld := make(map[string]interface{})
		celld[common.BKDBIN] = disModuleIDArr
		cond[common.BKModuleIDField] = celld
		fields := strings.Join(moduleField, ",")
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
		if ok && nil != setField {
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
		if ok && nil != moduleField {
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
