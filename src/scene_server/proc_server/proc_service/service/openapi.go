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
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"
)

func (ps *ProcServer) GetProcessPortByApplicationID(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	//get appID
	pathParams := req.PathParameters()
	appID, err := util.GetInt64ByInterface(pathParams[common.BKAppIDField])
	if err != nil {
		blog.Errorf("fail to get appid from pathparameter. err: %s,queryString:%+v,rid:%s", err.Error(), pathParams, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	modules := make([]map[string]interface{}, 0)
	if err := json.NewDecoder(req.Request.Body).Decode(&modules); err != nil {
		blog.Errorf("fail to decode request body. err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
	moduleToProcessesMap := make(map[int64][]mapstr.MapStr)
	for _, module := range modules {
		moduleName, ok := module[common.BKModuleNameField].(string)
		if !ok {
			blog.Warnf("assign error module['ModuleName'] is not string, module:%+v,rid:%s", module, srvData.rid)
			continue
		}

		processes, getErr := ps.getProcessesByModuleName(srvData.header, moduleName, appID)
		if getErr != nil {
			blog.Errorf("GetProcessesByModuleName failed int GetProcessPortByApplicationID, err: %s,input:%+v,rid:%s", getErr.Error(), module, srvData.rid)
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

	blog.V(5).Infof("moduleToProcessesMap: %+v,rid:%s", moduleToProcessesMap, srvData.rid)
	moduleHostConfigs, err := ps.getModuleHostConfigsByAppID(appID, srvData.header)
	if err != nil {
		blog.Errorf("getModuleHostConfigsByAppID failed in GetProcessPortByApplicationID, err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}

	blog.V(5).Infof("moduleHostConfigs:%v,rid:%s", moduleHostConfigs, srvData.rid)
	// 根据AppID获取AppInfo
	appInfoMap, err := ps.getAppInfoByID(appID, srvData.header)
	if err != nil {
		blog.Errorf("getAppInfoByID failed . err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}
	appInfo, ok := appInfoMap[appID]
	if !ok {
		blog.Errorf("GetProcessPortByApplicationID error : can not find app by id: %d,rid:%s", appID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}

	hostMap, err := ps.getHostMapByAppID(srvData.header, moduleHostConfigs)
	if err != nil {
		blog.Errorf("getHostMapByAppID failed in GetProcessPortByApplicationID. err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
		return
	}
	blog.V(5).Infof("GetProcessPortByApplicationID  hostMap:%+v,rid:%s", hostMap, srvData.rid)

	hostProcs := make(map[int64][]mapstr.MapStr, 0)
	for _, moduleHostConf := range moduleHostConfigs {
		hostID, ok := moduleHostConf[common.BKHostIDField]
		if !ok {
			blog.Errorf("fail to get hostID in GetProcessPortByApplicationID. error, field %s not found, data:%#v,rid:%s", common.BKHostIDField, moduleHostConf, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
			return
		}

		moduleID, ok := moduleHostConf[common.BKModuleIDField]
		if !ok {
			blog.Errorf("fail to get moduleID in GetProcessPortByApplicationID. err: field %s not found ,data:%#v,rid:%s", common.BKModuleIDField, moduleHostConf, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByApplicationIDFail)})
			return
		}

		procs, ok := hostProcs[hostID]
		if !ok {
			procs = make([]mapstr.MapStr, 0)
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
			processes = make([]mapstr.MapStr, 0)
		}
		host[common.BKProcField] = processes
		host[common.BKAppNameField] = appInfo[common.BKAppNameField]
		host[common.BKAppIDField] = appID
		retData = append(retData, host)
	}

	blog.V(5).Infof("GetProcessPortByApplicationID: %+v,rid:%s", retData, srvData.rid)
	resp.WriteEntity(meta.NewSuccessResp(retData))
}

//根据IP获取进程端口
func (ps *ProcServer) GetProcessPortByIP(req *restful.Request, resp *restful.Response) {
	srvData := ps.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	reqParam := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&reqParam); err != nil {
		blog.Errorf("fail to decode request body in GetProcessPortByIP. err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	ipArr := reqParam[common.BKIPArr]
	hostCondition := map[string]interface{}{common.BKHostInnerIPField: map[string]interface{}{"$in": ipArr}}
	hostData, hostIdArr, err := ps.getHostMapByCond(srvData.header, hostCondition)
	if err != nil {
		blog.Errorf("fail to getHostMapByCond in GetProcessPortByIP. err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
		return
	}
	// 获取appId
	configCondition := map[string][]int64{
		common.BKHostIDField: hostIdArr,
	}
	confArr, err := ps.getConfigByCond(srvData.header, configCondition)
	if err != nil {
		blog.Errorf("fail to getConfigByCond in GetProcessPortByIP. err: %s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
		return
	}
	blog.V(5).Infof("configArr: %+v,rid:%s", confArr, srvData.rid)
	//根据业务id获取进程
	resultData := make([]interface{}, 0)
	for _, item := range confArr {
		appId := item[common.BKAppIDField]
		moduleId := item[common.BKModuleIDField]
		hostId := item[common.BKHostIDField]
		//业务
		appData, err := ps.getAppInfoByID(appId, srvData.header)
		if err != nil {
			blog.Errorf("fail to getAppInfoByID in GetProcessPortByIP. err: %s,rid:%s", err.Error(), srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		//模块
		moduleData, err := ps.getModuleMapByCond(srvData.header, nil, map[string]interface{}{
			common.BKModuleIDField: moduleId,
			common.BKAppIDField:    appId,
		})
		if err != nil {
			blog.Errorf("fail to getModuleMap in GetProcessPortByIP. err: %s,appID:%+v,moduleID:%+v,rid:%s", err.Error(), appId, moduleId, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		moduleName, err := moduleData[moduleId].String(common.BKModuleNameField)
		if nil != err {
			blog.Errorf("fail to getModuleMap in GetProcessPortByIP not found moduleName. err: %s,rid:%s", err.Error(), srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		blog.V(5).Infof("moduleData:%+v,rid:%s", moduleData, srvData.rid)

		//进程
		procData, err := ps.getProcessMapByAppID(appId, srvData.header)
		if err != nil {
			blog.Errorf("fail to getProcessMapByAppID. err: %s,rid:%s", err.Error(), srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
			return
		}
		blog.V(5).Infof("procData: %+v,rid:%s", procData, srvData.rid)
		//获取绑定关系
		for _, itemProcData := range procData {
			result := make(map[string]interface{})
			procID, err := itemProcData.Int64(common.BKProcessIDField)
			if err != nil {
				blog.Warnf("fail to get procid in procdata(%+v),rid:%s", itemProcData, srvData.rid)
				continue
			}
			procModuleData, err := ps.getProcessBindModule(appId, procID, srvData.header)
			if err != nil {
				blog.Errorf("fail to getProcessBindModule in GetProcessPortByIP. err: %s,appID:%v,procID:%v,rid:%d", err.Error(), appId, procID, srvData.rid)
				resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcGetByIP)})
				return
			}

			for _, procMod := range procModuleData {
				itemMap, _ := procMod.(map[string]interface{})[common.BKModuleNameField].(string)
				blog.V(5).Infof("process module, %+v,rid:%s", itemMap, srvData.rid)
				if itemMap == moduleName {
					result[common.BKAppNameField], err = appData[appId].String(common.BKAppNameField)
					if nil != err {
						blog.Warnf("not foud app name,err:%s,rid:%s", err, srvData.rid)
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
					delete(itemProcData, common.BKProcessIDField)
					result["process"] = itemProcData
					resultData = append(resultData, result)
				}
			}
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(resultData))
}

// 根据模块获取所有关联的进程，建立Map ModuleToProcesses
func (ps *ProcServer) getProcessesByModuleName(forward http.Header, moduleName string, appID int64) ([]mapstr.MapStr, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr
	procData := make([]mapstr.MapStr, 0)
	params := mapstr.MapStr{
		common.BKAppIDField:      appID,
		common.BKModuleNameField: moduleName,
	}

	ret, err := ps.CoreAPI.ObjectController().OpenAPI().GetProcessesByModuleName(srvData.ctx, forward, params)
	if nil != err {
		blog.Errorf("getProcessesByModuleName http do error.  err:%s, input:%+v,rid:%s", err.Error(), params, srvData.rid)
		return procData, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getProcessesByModuleName http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, params, srvData.rid)
		return procData, defErr.New(ret.Code, ret.ErrMsg)
	}

	return ret.Data, nil
}

func (ps *ProcServer) getModuleHostConfigsByAppID(appID int64, forward http.Header) (moduleHostConfigs []map[string]int64, err error) {
	return ps.getConfigByCond(forward, map[string][]int64{
		common.BKAppIDField: []int64{int64(appID)},
	})
}

func (ps *ProcServer) getConfigByCond(forward http.Header, cond map[string][]int64) ([]map[string]int64, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr

	configArr := make([]map[string]int64, 0)
	ret, err := ps.CoreAPI.HostController().Module().GetModulesHostConfig(srvData.ctx, forward, cond)
	if nil != err {
		blog.Errorf("getConfigByCond http do error.  err:%s, input:%+v,rid:%s", err.Error(), cond, srvData.rid)
		return configArr, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getConfigByCond http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, cond, srvData.rid)
		return configArr, defErr.New(ret.Code, ret.ErrMsg)
	}

	for _, mdhost := range ret.Data {
		data := make(map[string]int64)
		data[common.BKAppIDField] = mdhost.AppID
		data[common.BKSetIDField] = mdhost.SetID
		data[common.BKModuleIDField] = mdhost.ModuleID
		data[common.BKHostIDField] = mdhost.HostID
		configArr = append(configArr, data)
	}

	return configArr, nil
}

func (ps *ProcServer) getAppInfoByID(appID int64, forward http.Header) (map[int64]mapstr.MapStr, error) {
	return ps.getAppMapByCond(forward, nil, map[string]interface{}{
		common.BKAppIDField: map[string]interface{}{
			"$in": []int64{appID},
		},
	})
}

func (ps *ProcServer) getAppMapByCond(forward http.Header, fields []string, cond mapstr.MapStr) (map[int64]mapstr.MapStr, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr

	appMap := make(map[int64]mapstr.MapStr, 0)
	input := new(meta.QueryCondition)
	input.Condition = cond
	input.Fields = fields
	input.SortArr = []meta.SearchSort{meta.SearchSort{Field: common.BKAppIDField}}
	input.Limit.Offset = 0
	input.Limit.Limit = common.BKNoLimit

	ret, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDApp, input)
	if nil != err {
		blog.Errorf("getAppMapByCond http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		return appMap, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("getAppMapByCond http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		return appMap, defErr.New(ret.Code, ret.ErrMsg)
	}
	for _, info := range ret.Data.Info {
		appID, err := info.Int64(common.BKAppIDField)
		if nil != err {
			continue
		}
		appMap[appID] = info
	}

	return appMap, nil
}

func (ps *ProcServer) getHostMapByAppID(forward http.Header, configData []map[string]int64) (map[int64]map[string]interface{}, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr
	hostIDArr := make([]int64, 0)
	for _, config := range configData {
		hostIDArr = append(hostIDArr, config[common.BKHostIDField])
	}

	hostMapCondition := map[string]interface{}{
		common.BKHostIDField: map[string]interface{}{
			"$in": hostIDArr,
		},
	}

	hostMap := make(map[int64]map[string]interface{})

	// build host controller
	input := new(meta.QueryInput)
	input.Fields = fmt.Sprintf("%s,%s,%s,%s", common.BKHostIDField, common.BKHostInnerIPField, common.BKCloudIDField, common.BKHostOuterIPField)
	input.Condition = hostMapCondition
	ret, err := ps.CoreAPI.HostController().Host().GetHosts(srvData.ctx, forward, input)
	if err != nil {
		blog.Errorf("getAppMapByCond http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		return hostMap, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if err == nil && !ret.Result {
		blog.Errorf("getAppMapByCond http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		return hostMap, defErr.New(ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		hostID, err := info.Int64(common.BKHostIDField)
		if nil != err {
			continue
		}
		hostMap[hostID] = info
	}

	return hostMap, nil
}

func (ps *ProcServer) getHostMapByCond(forward http.Header, condition map[string]interface{}) (map[int64]map[string]interface{}, []int64, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr

	hostMap := make(map[int64]map[string]interface{})
	hostIdArr := make([]int64, 0)

	input := new(meta.QueryInput)
	input.Fields = ""
	input.Condition = condition
	ret, err := ps.CoreAPI.HostController().Host().GetHosts(srvData.ctx, forward, input)
	if err != nil {
		blog.Errorf("getHostMapByCond http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		return hostMap, hostIdArr, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if err == nil && !ret.Result {
		blog.Errorf("getHostMapByCond http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		return hostMap, hostIdArr, defErr.New(ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		host_id, err := info.Int64(common.BKHostIDField)
		if nil != err {
			blog.Errorf("getHostMapByCond hostID not integer, err:%s,Search input:%+v, host info:%+v,rid:%s", err.Error(), input, info, srvData.rid)
			// CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
			return nil, nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDHost, common.BKHostIDField, "string", err.Error())
		}

		hostMap[host_id] = info
		hostIdArr = append(hostIdArr, host_id)
	}

	return hostMap, hostIdArr, nil
}

func (ps *ProcServer) getModuleMapByCond(forward http.Header, field []string, cond mapstr.MapStr) (map[int64]mapstr.MapStr, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr
	moduleMap := make(map[int64]mapstr.MapStr)
	input := new(meta.QueryCondition)
	input.Fields = field
	input.SortArr = []meta.SearchSort{meta.SearchSort{Field: common.BKModuleIDField}}
	input.Limit.Limit = common.BKNoLimit
	input.Condition = cond
	ret, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, input)
	if err != nil {
		blog.Errorf("getModuleMapByCond http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		return moduleMap, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if err == nil && !ret.Result {
		blog.Errorf("getHostMapByCond http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		return moduleMap, defErr.New(ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		moduleID, err := info.Int64(common.BKModuleIDField)
		if nil != err {
			blog.Warnf("fail to get moduleid in getModuleMapByCond. info: %+v,rid:%s", info, srvData.rid)
		} else {
			moduleMap[moduleID] = info
		}
	}

	return moduleMap, nil
}

func (ps *ProcServer) getProcessMapByAppID(appID int64, forward http.Header) (map[int64]mapstr.MapStr, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr

	procMap := make(map[int64]mapstr.MapStr)
	condition := mapstr.MapStr{
		common.BKAppIDField: appID,
	}

	input := new(meta.QueryCondition)
	input.Condition = condition
	ret, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDProc, input)
	if err != nil {
		blog.Errorf("getProcessMapByAppID http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		return procMap, defErr.Error(common.CCErrCommHTTPDoRequestFailed)

	}
	if !ret.Result {
		blog.Errorf("getProcessMapByAppID http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", ret.Code, ret.Result, input, srvData.rid)
		return procMap, defErr.New(ret.Code, ret.ErrMsg)
	}

	for _, info := range ret.Data.Info {
		processID, err := info.Int64(common.BKProcessIDField)
		if nil != err {
			blog.Warnf("fail to get appid in getProcessMapByAppID. info: %+v", info)
		} else {
			procMap[processID] = info
		}
	}

	return procMap, nil
}

func (ps *ProcServer) getProcessBindModule(appId, procId int64, forward http.Header) ([]interface{}, error) {
	srvData := ps.newSrvComm(forward)
	defErr := srvData.ccErr

	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	input := new(meta.QueryCondition)
	input.Condition = condition
	objModRet, err := ps.CoreAPI.CoreService().Instance().ReadInstance(srvData.ctx, srvData.header, common.BKInnerObjIDModule, input)
	if err != nil {
		blog.Errorf("getProcessMapByAppID  ReadInstance http do error.  err:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)

	}
	if !objModRet.Result {
		blog.Errorf("getProcessMapByAppID ReadInstance http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", objModRet.Code, objModRet.Result, input, srvData.rid)
		return nil, defErr.New(objModRet.Code, objModRet.ErrMsg)
	}

	moduleArr := objModRet.Data.Info
	condition[common.BKProcessIDField] = procId
	procRet, err := ps.CoreAPI.ProcController().GetProc2Module(srvData.ctx, forward, condition)
	if err != nil {
		blog.Errorf("getProcessMapByAppID http do error.  err:%s, input:%+v,rid:%s", err.Error(), condition, srvData.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)

	}
	if !procRet.Result {
		blog.Errorf("getProcessMapByAppID http reply  error. err code:%d err msg:%s, input:%+v,rid:%s", procRet.Code, procRet.Result, condition, srvData.rid)
		return nil, defErr.New(procRet.Code, procRet.ErrMsg)
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
				blog.Errorf("GetProcessBindModule get module default error:%s,rid:%d", err.Error(), srvData.rid)
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

	blog.V(5).Infof("getProcessBindModule result: %+v,rid:%s", result, srvData.rid)
	return result, nil
}
