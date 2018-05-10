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

package openapi

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"reflect"

	"github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var proc *procAction = &procAction{}

type procAction struct {
	base.BaseAction
}

//根据业务ID获取进程端口
func (cli *procAction) GetProcessPortByApplicationID(req *restful.Request, resp *restful.Response) {

	// 获取AppID
	blog.Error("GetProcessPortByApplicationID start")
	pathParams := req.PathParameters()
	appID, err := strconv.Atoi(pathParams[common.BKAppIDField])

	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	value, _ := ioutil.ReadAll(req.Request.Body)
	bodyData := make([]map[string]interface{}, 0)
	err = json.Unmarshal([]byte(value), &bodyData)
	blog.Debug("bodyData:%v", bodyData)
	if nil != err {
		blog.Error("GetProcessPortByApplicationID error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}

	modules := bodyData
	// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
	moduleToProcessesMap := make(map[int][]interface{})
	for _, module := range modules {
		moduleName, ok := module[common.BKModuleNameField].(string)
		if false == ok {
			blog.Error("assign error module['ModuleName'] is not string, module:%v", module)
			continue
		}
		processes, err := GetProcessesByModuleName(req, moduleName, cli.CC.ObjCtrl())
		if nil != err {
			msg := fmt.Sprintf("GetProcessesByModuleName error:%v", err)
			blog.Error(msg)
			cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, msg, resp)
		}
		if len(processes) > 0 {
			moduleToProcessesMap[int(module[common.BKModuleIDField].(float64))] = processes
		}
	}
	blog.Debug("moduleToProcessesMap%v", moduleToProcessesMap)
	moduleHostConfigs, err := getModuleHostConfigsByAppID(appID, req)
	if nil != err {
		blog.Error("GetProcessPortByApplicationID error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}
	blog.Debug("moduleHostConfigs:%v", moduleHostConfigs)
	// 根据AppID获取AppInfo
	appInfoMap, err := getAppInfoByID(appID, req)
	if nil != err {
		blog.Error("GetProcessPortByApplicationID error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}
	appInfoTemp, ok := appInfoMap[appID]
	if !ok {
		blog.Error("GetProcessPortByApplicationID error : can not find app by id: %d", appID)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}

	appInfo := appInfoTemp.(map[string]interface{})

	hostMap, err := getHostMapByAppID(moduleHostConfigs, req)
	if nil != err {
		blog.Error("GetProcessPortByApplicationID error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}
	blog.Debug("GetProcessPortByApplicationID  hostMap:%v", hostMap)

	hostProcs := make(map[int][]interface{}, 0)
	for _, moduleHostConfig := range moduleHostConfigs {
		hostID, errHostID := util.GetIntByInterface(moduleHostConfig[common.BKHostIDField])
		if nil != errHostID {
			err := defErr.Error(common.CCErrProcSelectBindToMoudleFaile)
			blog.Error("GetProcessPortByApplicationID error :%v", err)
			cli.ResponseFailed(common.CCErrProcSelectBindToMoudleFaile, err.Error(), resp)
			return
		}

		moduleID, ok := moduleHostConfig[common.BKModuleIDField]
		if false == ok {
			err := defErr.Error(common.CCErrProcSelectBindToMoudleFaile)
			blog.Error("GetProcessPortByApplicationID error :%v", err)
			defErr.Error(common.CCErrProcSelectBindToMoudleFaile)
			cli.ResponseFailed(common.CCErrProcSelectBindToMoudleFaile, err.Error(), resp)
			return
		}
		procs, ok := hostProcs[hostID]
		if false == ok {
			procs = make([]interface{}, 0)
		}
		processes, ok := moduleToProcessesMap[moduleID]
		if true == ok {
			hostProcs[hostID] = append(procs, processes...)
		}

	}

	resData := make([]interface{}, 0)
	for hostID, host := range hostMap {
		processes, ok := hostProcs[hostID]
		if false == ok {
			processes = make([]interface{}, 0)
		}
		host[common.BKProcField] = processes
		host[common.BKAppNameField] = appInfo[common.BKAppNameField]
		host[common.BKAppIDField] = appID
		resData = append(resData, host)
	}

	blog.Error("resDta:%v", resData)
	cli.ResponseSuccess(resData, resp)
}

//根据IP获取进程端口
func (cli *procAction) getProcessPortByIP(req *restful.Request, resp *restful.Response) {
	// 获取AppId
	blog.Debug("getProcessPortByIP start")
	value, _ := ioutil.ReadAll(req.Request.Body)
	input := make(map[string]interface{})
	err := json.Unmarshal([]byte(value), &input)
	if nil != err {
		blog.Error("getProcessPortByIP error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}
	ipArr := input["ipArr"]
	hostCondition := map[string]interface{}{common.BKHostInnerIPField: map[string]interface{}{"$in": ipArr}}
	hostData, hostIdArr, err := getHostMapByCond(req, hostCondition)
	if nil != err {
		blog.Error("getProcessPortByIP error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}
	// 获取appId
	configCondition := map[string]interface{}{
		common.BKHostIDField: hostIdArr,
	}
	confArr, err := getConfigByCond(req, cli.CC.HostCtrl(), configCondition)
	if nil != err {
		blog.Error("getProcessPortByIP error :%v", err)
		cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
		return
	}
	//根据业务id获取进程
	blog.Debug("configArr:%v", confArr)
	resultData := make([]interface{}, 0)
	for _, item := range confArr {
		appId := item[common.BKAppIDField]
		moduleId := item[common.BKModuleIDField]
		hostId := item[common.BKHostIDField]
		//业务
		appData, err := getAppInfoByID(appId, req)
		if nil != err {
			blog.Error("getProcessPortByIP getAppInfoById error :%v", err)
			cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
			return
		}
		//模块
		moduleData, err := GetModuleMapByCond(req, "", cli.CC.ObjCtrl(), map[string]interface{}{
			common.BKModuleIDField: moduleId,
			common.BKAppIDField:    appId,
		})
		moduleName := moduleData[moduleId].(map[string]interface{})[common.BKModuleNameField]
		blog.Debug("moduleData:%v", moduleData)

		//进程
		procData, err := getProcessMapByAppID(appId, req)
		if nil != err {
			blog.Error("getProcessPortByIP getProcessMapByAppId error :%v", err)
			cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
			return
		}
		blog.Debug("procDatass:%v", procData)
		//获取绑定关系
		result := make(map[string]interface{})
		for _, procDataV := range procData {
			procId, _ := util.GetIntByInterface(procDataV[common.BKProcIDField].(json.Number))
			procModuleData, err := GetProcessBindModule(req, appId, int(procId))
			if nil != err {
				blog.Error("getProcessPortByIP GetProcessBindModule error :%v", err)
				cli.ResponseFailed(common.CC_Err_Comm_GET_PROC_FAIL, common.CC_Err_Comm_GET_PROC_FAIL_STR, resp)
				return
			}
			blog.Debug("procModuleData:%v", procModuleData)
			blog.Debug("moduleName %s", moduleName)
			for _, procModuleV := range procModuleData {

				itemMap, _ := procModuleV.(map[string]interface{})[common.BKModuleNameField].(string)
				blog.Debug("process module, %v", itemMap)
				if itemMap == moduleName {
					result[common.BKAppNameField] = appData[appId].(map[string]interface{})[common.BKAppNameField]
					result[common.BKAppIDField] = appId
					result[common.BKHostIDField] = hostId
					result[common.BKHostInnerIPField] = hostData[hostId].(map[string]interface{})[common.BKHostInnerIPField]
					result[common.BKHostOuterIPField] = hostData[hostId].(map[string]interface{})[common.BKHostOuterIPField]
					if procDataV[common.BKBindIP].(string) == "第一内网IP" {
						procDataV[common.BKBindIP] = hostData[hostId].(map[string]interface{})[common.BKHostInnerIPField]
					}
					if procDataV[common.BKBindIP].(string) == "第一公网IP" {
						procDataV[common.BKBindIP] = hostData[hostId].(map[string]interface{})[common.BKHostOuterIPField]
					}
					delete(procDataV, common.BKAppIDField)
					delete(procDataV, common.BKProcIDField)
					result["process"] = procDataV
					resultData = append(resultData, result)
				}
			}
		}
	}
	cli.ResponseSuccess(resultData, resp)
}

//get modulemap by cond
func GetModuleMapByCond(req *restful.Request, fields string, objUrl string, cond interface{}) (map[int]interface{}, error) {
	moduleMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKModuleIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objUrl + "/object/v1/insts/module/search"
	blog.Info("GetModuleMapByCond url :%s", url)
	blog.Info("GetModuleMapByCond content :%s", string(bodyContent))
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Info("GetModuleMapByCond return :%s", string(reply))
	if err != nil {
		return moduleMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	moduleData := output["data"]
	moduleResult := moduleData.(map[string]interface{})
	moduleInfo := moduleResult["info"].([]interface{})
	for _, i := range moduleInfo {
		module := i.(map[string]interface{})
		moduleId, _ := module[common.BKModuleIDField].(json.Number).Int64()
		moduleMap[int(moduleId)] = i
	}
	return moduleMap, nil
}

//get process bind module
type ProcessCResult struct {
	Result  bool           `json:"result"`
	Code    int            `json:"code"`
	Message interface{}    `json:"message"`
	Data    map[string]int `json:"data"`
}

type ProcessResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type ModuleSResult struct {
	Result  bool        `json:"result"`
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    ModuleData  `json:"data"`
}

type ModuleData struct {
	Count int                      `json:"count"`
	Info  []map[string]interface{} `json:"info"`
}

type ProcModuleConfig struct {
	ApplicationID int
	ModuleName    string
	processID     int
}

type ProcModuleResult struct {
	Result  bool               `json:"result"`
	Code    int                `json:"code"`
	Message interface{}        `json:"message"`
	Data    []ProcModuleConfig `json:"data"`
}

func GetProcessBindModule(req *restful.Request, appId int, procId int) ([]interface{}, error) {
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	sCondJson, _ := json.Marshal(searchParams)

	gModuleUrl := proc.CC.ObjCtrl() + "/object/v1/insts/module/search"
	blog.Info("get module query url: %v", gModuleUrl)
	blog.Info("get module query params: %v", string(sCondJson))
	gModuleRe, err := httpcli.ReqHttp(req, gModuleUrl, common.HTTPSelectPost, []byte(sCondJson))
	blog.Info("get module query return: %v", gModuleRe)
	if nil != err {
		blog.Error("GetProcessBindModule Module  error :%v", err)
		return nil, nil
	}
	var modules ModuleSResult
	err = json.Unmarshal([]byte(gModuleRe), &modules)
	if nil != err {
		blog.Error("GetProcessBindModule Module  error :%v", err)
		return nil, nil
	}
	moduleArr := modules.Data.Info
	gProc2ModuleUrl := proc.CC.ProcCtrl() + "/process/v1/module/search"
	condition[common.BKProcIDField] = procId
	sCondJson, _ = json.Marshal(condition)
	blog.Info("get module config query url: %v", gModuleUrl)
	blog.Info("get module config params: %v", string(sCondJson))
	gPorc2ModuleRe, err := httpcli.ReqHttp(req, gProc2ModuleUrl, common.HTTPSelectPost, []byte(sCondJson))
	blog.Info("get module config return: %v", gPorc2ModuleRe)
	if nil != err {
		blog.Error("get module config params  error :%v", err)
		return nil, nil
	}
	var pro2Module ProcModuleResult
	err = json.Unmarshal([]byte(gPorc2ModuleRe), &pro2Module)
	if nil != err {
		blog.Error("get module config params  error :%v", err)
		return nil, nil
	}
	procModuleData := pro2Module.Data
	disModuleNameArr := make([]string, 0)
	for _, i := range moduleArr {
		if !util.InArray(i[common.BKModuleNameField], disModuleNameArr) {
			moduleName, ok := i[common.BKModuleNameField].(string)
			if false == ok {
				continue
			}
			isDefault64, err := util.GetInt64ByInterface(i[common.BKDefaultField])
			if nil != err {
				blog.Errorf("GetProcessBindModule get module default error:%s", err.Error())
				continue

			} else {
				if 0 != isDefault64 {
					continue
				}
			}
			disModuleNameArr = append(disModuleNameArr, moduleName)
		}
	}
	result := make([]interface{}, 0)
	for _, j := range disModuleNameArr {
		num := 0
		isBind := 0
		for _, k := range moduleArr {
			moduleName, ok := k[common.BKModuleNameField].(string)
			if false == ok {
				continue
			}
			if j == moduleName {
				num++
			}
		}
		for _, m := range procModuleData {
			if j == m.ModuleName {
				isBind = 1
				break
			}
		}
		data := make(map[string]interface{})
		data[common.BKModuleNameField] = j
		data["set_num"] = num
		data["is_bind"] = isBind
		result = append(result, data)
	}
	blog.Debug("result:%v", result)
	return result, nil
}

func getProcessMapByAppID(appId int, req *restful.Request) (map[int]map[string]interface{}, error) {
	procMap := make(map[int]map[string]interface{})
	condition := map[string]interface{}{
		common.BKAppIDField: appId,
	}
	searchParams := map[string]interface{}{
		"fields":    "",
		"condition": condition,
	}
	//build host controller url
	url := proc.CC.ObjCtrl() + "/object/v1/insts/process/search"

	inputJson, _ := json.Marshal(searchParams)
	processInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	if nil != err {
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(processInfo))
	res, _ := js.Map()
	resData := res["data"].(map[string]interface{})
	resDataInfo := resData["info"].([]interface{})
	for _, item := range resDataInfo {
		proc := item.(map[string]interface{})
		appId, _ := proc[common.BKAppIDField].(json.Number).Int64()
		procMap[int(appId)] = proc
	}
	return procMap, nil
}

func getHostMapByAppID(configData []map[string]int, req *restful.Request) (map[int]map[string]interface{}, error) {
	hostIDArr := make([]int, 0)

	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			"$in": hostIDArr,
		},
	}

	hostMap := make(map[int]map[string]interface{})

	// build host controller url
	url := proc.CC.HostCtrl() + "/host/v1/hosts/search"
	searchParams := map[string]interface{}{
		"fields":    fmt.Sprintf("%s,%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField, common.BKHostOuterIPField),
		"condition": hostMapCondition,
	}
	inputJson, _ := json.Marshal(searchParams)

	appInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	blog.Infof("url:%s, params:%v, reply:%s", url, string(inputJson), string(appInfo))
	if nil != err {
		return nil, err
	}
	js, err := simplejson.NewJson([]byte(appInfo))

	res, _ := js.Map()
	resData := res["data"].(map[string]interface{})
	resDataInfo := resData["info"].([]interface{})

	for _, item := range resDataInfo {
		host := item.(map[string]interface{})
		hostID, _ := host[common.BKHostIDField].(json.Number).Int64()
		hostMap[int(hostID)] = host
	}
	return hostMap, nil
}

func getAppInfoByID(appID int, req *restful.Request) (map[int]interface{}, error) {
	appMap, err := getAppMapByCond(req, "", proc.CC.ObjCtrl(), map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			"$in": []int{appID},
		},
	})
	return appMap, err
}

func getModulesByAppID(appID int, req *restful.Request) (modules []map[string]interface{}, err error) {
	blog.Debug("getModuleIDArrByAppID start, appID:%d", appID)

	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID

	searchParams := make(map[string]interface{})
	searchParams["condition"] = condition
	searchParams["fields"] = fmt.Sprintf("%s,%s", common.BKModuleIDField, common.BKModuleNameField)
	searchParams["start"] = 0
	searchParams["limit"] = 0
	searchParams["sort"] = ""

	sURL := proc.CC.ObjCtrl() + "/object/v1/insts/module/search"
	inputJson, _ := json.Marshal(searchParams)
	moduleRes, err := httpcli.ReqHttp(req, sURL, "POST", []byte(inputJson))

	blog.Debug("search modules end, moduleRes:%v", moduleRes)

	if nil != err {
		return nil, err
	}

	resData := make([]map[string]interface{}, 0)

	var rst map[string]interface{}
	if jsErr := json.Unmarshal([]byte(moduleRes), &rst); nil == jsErr {
		if rst["result"].(bool) {
			modules := (rst["data"].(map[string]interface{}))["info"]
			for _, module := range modules.([]interface{}) {
				resData = append(resData, module.(map[string]interface{}))
			}
			return resData, nil
		} else {
			return nil, errors.New(rst["message"].(string))
		}
	} else {
		return nil, jsErr
	}
}

func getModuleHostConfigsByAppID(appID int, req *restful.Request) (moduleHostConfigs []map[string]int, err error) {
	configData, err := getConfigByCond(req, proc.CC.HostCtrl(), map[string]interface{}{
		common.BKAppIDField: []interface{}{appID},
	})
	return configData, err
}

func getAppMapByCond(req *restful.Request, fields string, objURL string, cond interface{}) (map[int]interface{}, error) {
	appMap := make(map[int]interface{})
	condition := make(map[string]interface{})
	condition["fields"] = fields
	condition["sort"] = common.BKAppIDField
	condition["start"] = 0
	condition["limit"] = 0
	condition["condition"] = cond
	bodyContent, _ := json.Marshal(condition)
	url := objURL + "/object/v1/insts/" + common.BKInnerObjIDApp + "/search"
	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Debug("getAppMapByCond, url:%s, reply:%s, body:%s", url, reply, string(bodyContent))
	if err != nil {
		return appMap, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, _ := js.Map()
	appData := output["data"]
	appResult := appData.(map[string]interface{})
	appInfo := appResult["info"].([]interface{})
	for _, i := range appInfo {
		app := i.(map[string]interface{})
		appID, _ := app[common.BKAppIDField].(json.Number).Int64()
		appMap[int(appID)] = i
	}
	return appMap, nil
}

func getConfigByCond(req *restful.Request, hostURL string, cond interface{}) ([]map[string]int, error) {
	configArr := make([]map[string]int, 0)
	bodyContent, _ := json.Marshal(cond)
	url := hostURL + "/host/v1/meta/hosts/module/config/search"

	reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(bodyContent))
	blog.Debug("getConfigByCond, url:%s, reply:%s, body:%s", url, reply, string(bodyContent))
	if err != nil {
		return configArr, err
	}
	js, _ := simplejson.NewJson([]byte(reply))
	output, err := js.Map()
	configData := output["data"]
	configInfo, _ := configData.([]interface{})
	for _, mh := range configInfo {
		celh := mh.(map[string]interface{})

		hostID, _ := celh[common.BKHostIDField].(json.Number).Int64()
		setID, _ := celh[common.BKSetIDField].(json.Number).Int64()
		moduleID, _ := celh[common.BKModuleIDField].(json.Number).Int64()
		appID, _ := celh[common.BKAppIDField].(json.Number).Int64()
		data := make(map[string]int)
		data[common.BKAppIDField] = int(appID)
		data[common.BKSetIDField] = int(setID)
		data[common.BKModuleIDField] = int(moduleID)
		data[common.BKHostIDField] = int(hostID)
		configArr = append(configArr, data)
	}
	return configArr, nil
}

func getHostMapByCond(req *restful.Request, condition map[string]interface{}) (map[int]interface{}, []int, error) {
	hostMap := make(map[int]interface{})
	hostIdArr := make([]int, 0)

	// build host controller url
	url := proc.CC.HostCtrl() + "/host/v1/hosts/search"
	searchParams := map[string]interface{}{
		"fields":    "",
		"condition": condition,
	}
	inputJson, err := json.Marshal(searchParams)
	if nil != err {
		return nil, nil, err
	}

	appInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	blog.Debug("getHostMapByCond, url:%s, reply:%s, body:%s", url, appInfo, string(inputJson))
	if nil != err {
		return hostMap, hostIdArr, err
	}

	js, err := simplejson.NewJson([]byte(appInfo))
	if nil != err {
		return nil, nil, err
	}

	res, err := js.Map()
	if nil != err {
		return nil, nil, err
	}

	resData := res["data"].(map[string]interface{})
	resDataInfo := resData["info"].([]interface{})

	for _, item := range resDataInfo {
		host := item.(map[string]interface{})
		host_id, err := host[common.BKHostIDField].(json.Number).Int64()
		if nil != err {
			return nil, nil, err
		}

		hostMap[int(host_id)] = host
		hostIdArr = append(hostIdArr, int(host_id))
	}
	return hostMap, hostIdArr, nil
}

// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
func GetProcessesByModuleName(req *restful.Request, moduleName string, objUrl string) ([]interface{}, error) {
	procData := make([]interface{}, 0)
	url := objUrl + "/object/v1/openapi/proc/getProcModule"
	searchParams := map[string]interface{}{
		common.BKModuleNameField: moduleName,
	}
	inputJson, err := json.Marshal(searchParams)
	if nil != err {
		return nil, err
	}
	procInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
	blog.Info("url:%s, reply:%s", url, string(procInfo))
	if nil != err {
		blog.Error("get process by module err:%v", err)
		return procData, err
	}

	js, err := simplejson.NewJson([]byte(procInfo))
	if nil != err {
		blog.Error("simplejson error :%v", err)
		return nil, err
	}
	res, err := js.Map()
	if nil != err {
		blog.Error("json -> map error :%v", err)
		return nil, err
	}
	result, err := js.Get("result").Bool()
	if nil != err {
		blog.Error("assign error res['result'] is not bool , res:%v", res)
	}
	if result {

		procData, err := js.Get("data").Array() //.([]interface{})

		if nil != err {
			blog.Error("assign error res['data'] is not []interface{} , res:%v", reflect.TypeOf(res["data"]))
		}
		return procData, nil
	} else {
		message, ok := res[common.HTTPBKAPIErrorMessage].(string)
		if false == ok {
			blog.Error("assign error res['message'] is not string , res:%v", res)
		}
		return nil, errors.New(message)
	}
	blog.Infof("procData %v", procData)
	return procData, nil
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/GetProcessPortByApplicationID/{" + common.BKAppIDField + "}", Params: nil, Handler: proc.GetProcessPortByApplicationID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/openapi/GetProcessPortByIP", Params: nil, Handler: proc.getProcessPortByIP})
	proc.CreateAction()
}
