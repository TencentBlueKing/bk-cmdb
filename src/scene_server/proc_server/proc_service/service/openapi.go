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
package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) GetProcessPortByApplicationID(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	//get appID
	pathParams := req.PathParameters()
	appID, err := strconv.Atoi(pathParams[common.BKAppIDField])
	if err != nil {
		blog.Errorf("fail to get appid from pathparameter. err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	modules := make([]map[string]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&modules); err != nil {
		blog.Errorf("fail to decode request body. err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	forward := req.Request.Header
	// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
	moduleToProcessesMap := make(map[int64][]interface{})
	for _, module := range modules {
		moduleName, ok := module[common.BKModuleNameField].(string)
		if !ok {
			blog.Warnf("assign error module['ModuleName'] is not string, module:%v", module)
			continue
		}

		processes, getErr := ps.getProcessesByModuleName(forward, moduleName)
		if getErr != nil {
			blog.Errorf("GetProcessesByModuleName failed int GetProcessPortByApplicationID, err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
			return
		}
		if len(processes) > 0 {
			moduleID, err := util.GetInt64ByInterface(module[common.BKModuleIDField])
			if nil == err {
				moduleToProcessesMap[moduleID] = processes
			}

		}

	}

	blog.V(3).Infof("moduleToProcessesMap: %v", moduleToProcessesMap)
	moduleHostConfigs, err := ps.getModuleHostConfigsByAppID(appID, forward)
	if err != nil {
		blog.Errorf("getModuleHostConfigsByAppID failed in GetProcessPortByApplicationID, err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}

	blog.V(3).Infof("moduleHostConfigs:%v", moduleHostConfigs)
	// 根据AppID获取AppInfo
	appInfoMap, err := ps.getAppInfoByID(appID, forward)
	if err != nil {
		blog.Errorf("getAppInfoByID failed in GetProcessPortByApplicationID. err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}
	appInfTemp, ok := appInfoMap[appID]
	if !ok {
		blog.Errorf("GetProcessPortByApplicationID error : can not find app by id: %d", appID)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}

	appInfo := appInfTemp.(mapstr.MapStr)

	hostMap, err := ps.getHostMapByAppID(forward, moduleHostConfigs)
	if err != nil {
		blog.Errorf("getHostMapByAppID failed in GetProcessPortByApplicationID. err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}
	blog.V(3).Infof("GetProcessPortByApplicationID  hostMap:%v", hostMap)

	hostProcs := make(map[int][]interface{}, 0)
	for _, moduleHostConf := range moduleHostConfigs {
		hostID, ok := moduleHostConf[common.BKHostIDField]
		if !ok {
			blog.Errorf("fail to get hostID in GetProcessPortByApplicationID. err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
			return
		}

		moduleID, ok := moduleHostConf[common.BKModuleIDField]
		if !ok {
			blog.Errorf("fail to get moduleID in GetProcessPortByApplicationID. err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
			return
		}

		procs, ok := hostProcs[hostID]
		if !ok {
			procs = make([]interface{}, 0)
		}

		processes, ok := moduleToProcessesMap[int64(moduleID)]
		if ok {
			hostProcs[hostID] = append(procs, processes...)
		}

	}

	retData := make([]interface{}, 0)
	for hostID, host := range hostMap {
		processes, ok := hostProcs[hostID]
		if !ok {
			processes = make([]interface{}, 0)
		}
		host[common.BKProcField] = processes
		host[common.BKAppNameField] = appInfo[common.BKAppNameField]
		host[common.BKAppIDField] = appID
		retData = append(retData, host)
	}

	blog.V(3).Infof("GetProcessPortByApplicationID: %+v", retData)
	resp.WriteEntity(meta.NewSuccessResp(retData))
}

//根据IP获取进程端口
func (ps *ProcServer) GetProcessPortByIP(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("fail to decode request body in GetProcessPortByIP. err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	forward := req.Request.Header
	ipArr := reqParam[common.BKIPArr]
	hostCondition := map[string]interface{}{common.BKHostInnerIPField: map[string]interface{}{"$in": ipArr}}
	hostData, hostIdArr, err := ps.getHostMapByCond(forward, hostCondition)
	if err != nil {
		blog.Errorf("fail to getHostMapByCond in GetProcessPortByIP. err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
		return
	}
	// 获取appId
	configCondition := map[string][]int64{
		common.BKHostIDField: hostIdArr,
	}
	confArr, err := ps.getConfigByCond(forward, configCondition)
	if err != nil {
		blog.Errorf("fail to getConfigByCond in GetProcessPortByIP. err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
		return
	}
	blog.V(3).Infof("configArr: %+v", confArr)
	//根据业务id获取进程
	resultData := make([]interface{}, 0)
	for _, item := range confArr {
		appId := item[common.BKAppIDField]
		moduleId := item[common.BKModuleIDField]
		hostId := item[common.BKHostIDField]
		//业务
		appData, err := ps.getAppInfoByID(appId, forward)
		if err != nil {
			blog.Errorf("fail to getAppInfoByID in GetProcessPortByIP. err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		//模块
		moduleData, err := ps.getModuleMapByCond(forward, "", map[string]interface{}{
			common.BKModuleIDField: moduleId,
			common.BKAppIDField:    appId,
		})
		if err != nil {
			blog.Errorf("fail to getModuleMap in GetProcessPortByIP. err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		moduleName, err := moduleData[moduleId].(mapstr.MapStr).String(common.BKModuleNameField)
		if nil != err {
			blog.Errorf("fail to getModuleMap in GetProcessPortByIP not found moduleName. err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		blog.V(3).Infof("moduleData:%v", moduleData)

		//进程
		procData, err := ps.getProcessMapByAppID(appId, forward)
		if err != nil {
			blog.Errorf("fail to getProcessMapByAppID. err: %s", err.Error())
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		blog.V(3).Infof("procData: %v", procData)
		//获取绑定关系
		result := make(map[string]interface{})
		for _, itemProcData := range procData {
			procId, err := util.GetIntByInterface(itemProcData[common.BKProcIDField])
			if err != nil {
				blog.Warnf("fail to get procid in procdata(%+v)", itemProcData)
				continue
			}
			procModuleData, err := ps.getProcessBindModule(appId, procId, forward)
			if err != nil {
				blog.Errorf("fail to getProcessBindModule in GetProcessPortByIP. err: %s", err.Error())
				resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
				return
			}

			for _, procMod := range procModuleData {
				itemMap, _ := procMod.(map[string]interface{})[common.BKModuleNameField].(string)
				blog.V(3).Infof("process module, %v", itemMap)
				if itemMap == moduleName {
					result[common.BKAppNameField], err = appData[appId].(mapstr.MapStr).String(common.BKAppNameField)
					if nil != err {
						blog.Warnf("not foud app name")
						continue
					}
					result[common.BKAppIDField] = appId
					result[common.BKHostIDField] = hostId
					result[common.BKHostInnerIPField] = hostData[hostId][common.BKHostInnerIPField]
					result[common.BKHostOuterIPField] = hostData[hostId][common.BKHostOuterIPField]
					switch t := itemProcData[common.BKBindIP].(type) {
					case string:
						if t == "第一内网IP" {
							itemProcData[common.BKBindIP] = hostData[hostId][common.BKHostInnerIPField]
						}
						if t == "第一公网IP" {
							itemProcData[common.BKBindIP] = hostData[hostId][common.BKHostOuterIPField]
						}
					}

					delete(itemProcData, common.BKAppIDField)
					delete(itemProcData, common.BKProcIDField)
					result["process"] = itemProcData
					resultData = append(resultData, result)
				}
			}
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(resultData))
}

// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
func (ps *ProcServer) getProcessesByModuleName(forward http.Header, moduleName string) ([]interface{}, error) {
	procData := make([]interface{}, 0)
	params := map[string]interface{}{
		common.BKModuleNameField: moduleName,
	}

	ret, err := ps.CoreAPI.ObjectController().OpenAPI().GetProcessesByModuleName(context.Background(), forward, params)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("get process by module failed. err: %s, errcode: %d, errmsg: %s", err.Error(), ret.Code, ret.ErrMsg)
		return procData, err
	}

	procData = append(procData, ret.Data)
	return procData, nil
}

func (ps *ProcServer) getModuleHostConfigsByAppID(appID int, forward http.Header) (moduleHostConfigs []map[string]int, err error) {
	return ps.getConfigByCond(forward, map[string][]int64{
		common.BKAppIDField: []int64{int64(appID)},
	})
}

func (ps *ProcServer) getConfigByCond(forward http.Header, cond map[string][]int64) ([]map[string]int, error) {
	configArr := make([]map[string]int, 0)
	ret, err := ps.CoreAPI.HostController().Module().GetModulesHostConfig(context.Background(), forward, cond)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("fail to get module host config. err:%v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		return configArr, fmt.Errorf("fail to get module host config")
	}

	for _, mdhost := range ret.Data {
		data := make(map[string]int)
		data[common.BKAppIDField] = int(mdhost.AppID)
		data[common.BKSetIDField] = int(mdhost.SetID)
		data[common.BKModuleIDField] = int(mdhost.ModuleID)
		data[common.BKHostIDField] = int(mdhost.HostID)
		configArr = append(configArr, data)
	}

	return configArr, nil
}

func (ps *ProcServer) getAppInfoByID(appID int, forward http.Header) (map[int]interface{}, error) {
	return ps.getAppMapByCond(forward, "", map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			"$in": []int{appID},
		},
	})
}

func (ps *ProcServer) getAppMapByCond(forward http.Header, fields string, cond interface{}) (map[int]interface{}, error) {
	appMap := make(map[int]interface{})
	input := new(meta.QueryInput)
	input.Condition = cond
	input.Fields = fields
	input.Sort = common.BKAppIDField
	input.Start = 0
	input.Limit = 0
	ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, forward, input)
	if err != nil || (err == nil && !ret.Result) {
		blog.Errorf("fail to get appinfo by condition. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
		return appMap, fmt.Errorf("fail to get appid by condition")
	}

	for _, info := range ret.Data.Info {
		appID, err := info.Int64(common.BKAppIDField)
		if nil != err {
			continue
		}
		appMap[int(appID)] = info
	}

	return appMap, nil
}

func (ps *ProcServer) getHostMapByAppID(forward http.Header, configData []map[string]int) (map[int]map[string]interface{}, error) {
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

	// build host controller
	input := new(meta.QueryInput)
	input.Fields = fmt.Sprintf("%s,%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField, common.BKHostOuterIPField)
	input.Condition = hostMapCondition
	ret, err := ps.CoreAPI.HostController().Host().GetHosts(context.Background(), forward, input)
	if err != nil || (err == nil && !ret.Result) {
		return hostMap, fmt.Errorf("fail to gethosts. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		hostID, err := util.GetInt64ByInterface(info[common.BKHostIDField])
		if nil != err {
			continue
		}
		hostMap[int(hostID)] = info
	}

	return hostMap, nil
}

func (ps *ProcServer) getHostMapByCond(forward http.Header, condition map[string]interface{}) (map[int]map[string]interface{}, []int64, error) {
	hostMap := make(map[int]map[string]interface{})
	hostIdArr := make([]int64, 0)

	input := new(meta.QueryInput)
	input.Fields = ""
	input.Condition = condition
	ret, err := ps.CoreAPI.HostController().Host().GetHosts(context.Background(), forward, input)
	if err != nil || (err == nil && !ret.Result) {
		return hostMap, hostIdArr, fmt.Errorf("fail to getHostMapByCondition. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		host_id, err := util.GetInt64ByInterface(info[common.BKHostIDField])
		if nil != err {
			return nil, nil, fmt.Errorf("fail to get hostid")
		}

		hostMap[int(host_id)] = info
		hostIdArr = append(hostIdArr, host_id)
	}

	return hostMap, hostIdArr, nil
}

func (ps *ProcServer) getModuleMapByCond(forward http.Header, field string, cond interface{}) (map[int]interface{}, error) {
	moduleMap := make(map[int]interface{})
	input := new(meta.QueryInput)
	input.Fields = field
	input.Sort = common.BKModuleIDField
	input.Start = 0
	input.Limit = 0
	input.Condition = cond
	ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, forward, input)
	if err != nil || (err == nil && !ret.Result) {
		return moduleMap, fmt.Errorf("fail to getModuleMapByCond. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		moduleId, err := info.Int64(common.BKModuleIDField)
		if nil != err {
			blog.Warnf("fail to get moduleid in getModuleMapByCond. info: %v", info)
		} else {
			moduleMap[int(moduleId)] = info
		}
	}

	return moduleMap, nil
}

func (ps *ProcServer) getProcessMapByAppID(appId int, forward http.Header) (map[int]map[string]interface{}, error) {
	procMap := make(map[int]map[string]interface{})
	condition := map[string]interface{}{
		common.BKAppIDField: appId,
	}

	input := new(meta.QueryInput)
	input.Condition = condition
	input.Fields = ""
	ret, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDProc, forward, input)
	if err != nil || (err == nil && !ret.Result) {
		return procMap, fmt.Errorf("fail to getProcessMapByAppID. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		appId, err := info.Int64(common.BKAppIDField)
		if nil != err {
			blog.Warnf("fail to get appid in getProcessMapByAppID. info: %+v", info)
		} else {
			procMap[int(appId)] = info
		}
	}

	return procMap, nil
}

func (ps *ProcServer) getProcessBindModule(appId, procId int, forward http.Header) ([]interface{}, error) {
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	input := new(meta.QueryInput)
	input.Condition = condition
	objModRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, forward, input)
	if err != nil || (err == nil && !objModRet.Result) {
		return nil, fmt.Errorf("fail to get module by appid(%d). err: %v, errcode: %d, errmsg: %s", err, objModRet.Code, objModRet.ErrMsg)
	}

	moduleArr := objModRet.Data.Info
	condition[common.BKProcIDField] = procId
	procRet, err := ps.CoreAPI.ProcController().GetProc2Module(context.Background(), forward, condition)
	if err != nil || (err == nil && !procRet.Result) {
		return nil, fmt.Errorf("fail to GetProc2Module in getProcessBindModule. err: %v, errcode: %d, errmsg: %s", err, procRet.Code, procRet.ErrMsg)
	}

	procModuleData := procRet.Data
	disModuleNameArr := make([]string, 0)
	for _, modArr := range moduleArr {
		if !util.InArray(modArr[common.BKModuleNameField], disModuleNameArr) {
			moduleName, ok := modArr[common.BKModuleNameField].(string)
			if !ok {
				continue
			}
			isDefault64, err := util.GetInt64ByInterface(modArr[common.BKDefaultField])
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
	for _, disModName := range disModuleNameArr {
		num := 0
		isBind := 0
		for _, module := range moduleArr {
			moduleName, ok := module[common.BKModuleNameField].(string)
			if !ok {
				continue
			}
			if disModName == moduleName {
				num++
			}
		}
		for _, procMod := range procModuleData {
			if disModName == procMod.ModuleName {
				isBind = 1
				break
			}
		}

		data := make(map[string]interface{})
		data[common.BKModuleNameField] = disModName
		data["set_num"] = num
		data["is_bind"] = isBind
		result = append(result, data)
	}

	blog.V(3).Infof("getProcessBindModule result: %+v", result)
	return result, nil
}
