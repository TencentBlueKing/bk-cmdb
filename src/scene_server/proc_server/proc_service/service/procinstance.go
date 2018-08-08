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
	"regexp"
	"strconv"

	"github.com/emicklei/go-restful"
	"github.com/gin-gonic/gin/json"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (ps *ProcServer) OperateProcessInstance(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)
	namespace := req.PathParameter("namespace")

	procOpParam := new(meta.ProcessOperate)
	if err := json.NewDecoder(req.Request.Body).Decode(procOpParam); err != nil {
		blog.Errorf("fail to decode process operation parameter. err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	forward := req.Request.Header

	procInstModel, err := ps.matchProcessInstance(procOpParam, forward)
	if err != nil {
		blog.Errorf("match process instance failed in OperateProcessInstance. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcOperateFaile)})
		return
	}

	result, err := ps.operateProcInstanceByGse(procOpParam, procInstModel, namespace, forward)
	if err != nil {
		blog.Errorf("operate process failed. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcOperateFaile)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) QueryProcessOperateResult(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := ps.CCErr.CreateDefaultCCErrorIf(language)

	namespace := req.PathParameter("namespace")

	model_TaskId := make(map[string]string)
	if err := json.NewDecoder(req.Request.Body).Decode(&model_TaskId); err != nil {
		blog.Errorf("fail to decode QueryProcessOperateResult request body. err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	result := make([]map[string]interface{}, 0)
	for modkey, taskId := range model_TaskId {
		ret, err := ps.CoreAPI.GseProcServer().QueryProcOperateResult(context.Background(), req.Request.Header, namespace, taskId)
		if err != nil || (err == nil && !ret.Result) {
			rspMap := make(map[string]interface{})
			rspMap[taskId] = &meta.BaseResp{Code: 115, ErrMsg: "please retry query again"}
			rspMap[modkey] = taskId
			result = append(result, rspMap)
		}

		ret.Data[modkey] = taskId
		result = append(result, ret.Data)
	}

	resp.WriteEntity(meta.NewSuccessResp(result))
}

func (ps *ProcServer) operateProcInstanceByGse(procOp *meta.ProcessOperate, instModels map[string]*meta.ProcInstanceModel, namespace string, forward http.Header) (map[string]string, error) {
	var err error
	model_TaskId := make(map[string]string)
	for key, model := range instModels {
		gseprocReq := new(meta.GseProcRequest)
		// hostInfo
		gseprocReq.Hosts, err = ps.getHostForGse(string(model.ApplicationID), string(model.HostId), forward)
		if err != nil {
			blog.Warnf("getHostForGse failed. err: %v", err)
			continue
		}

		// get processinfo
		procInfo, err := ps.getProcessbyProcID(string(model.ProcID), forward)
		if err != nil {
			blog.Warnf("getProcessByProcID failed. err: %v", err)
			continue
		}

		// register process into gse
		if err := ps.registerProcInstanceToGse(namespace, procInfo, forward); err != nil {
			blog.Warnf("register process into gse failed. err: %v", err)
			continue
		}

		procName, ok := procInfo[common.BKProcessNameField].(string)
		if !ok {
			blog.Warnf("convert process name to string failed")
			continue
		}

		gseprocReq.Meta.Name = procName
		gseprocReq.Meta.Namespace = namespace
		gseprocReq.OpType = procOp.OpType

		gseRsp, err := ps.CoreAPI.GseProcServer().OperateProcess(context.Background(), forward, gseprocReq.Meta.Namespace, gseprocReq)
		if err != nil || (err == nil || !gseRsp.Result) {
			blog.Warnf("fail to operate process by gse process server. err: %v, errcode: %d, errmsg: %s", err, gseRsp.Code, gseRsp.ErrMsg)
			continue
		}

		taskId, ok := gseRsp.Data[common.BKGseTaskIdField].(string)
		if !ok {
			blog.Warnf("convert gse process operate taskid to string failed. value: %v", gseRsp.Data)
			continue
		}

		model_TaskId[key] = taskId
	}

	return model_TaskId, nil
}

func (ps *ProcServer) registerProcInstanceToGse(namespace string, procInfo map[string]interface{}, forward http.Header) error {
	// process
	procName, _ := procInfo[common.BKProcessNameField].(string)
	pidFilePath, _ := procInfo[common.BKProcPidFile].(string)
	workPath, _ := procInfo[common.BKProcWorkPath].(string)
	startCmd, _ := procInfo[common.BKProcStartCmd].(string)
	stopCmd, _ := procInfo[common.BKProcStopCmd].(string)
	reloadCmd, _ := procInfo[common.BKProcReloadCmd].(string)
	restartCmd, _ := procInfo[common.BKProcRestartCmd].(string)

	gseproc := new(meta.GseProcRequest)
	gseproc.Meta.Namespace = namespace
	gseproc.Meta.Name = procName
	gseproc.Spec.Identity.PidPath = pidFilePath
	gseproc.Spec.Identity.SetupPath = workPath
	gseproc.Spec.Control.StartCmd = startCmd
	gseproc.Spec.Control.StopCmd = stopCmd
	gseproc.Spec.Control.ReloadCmd = reloadCmd
	gseproc.Spec.Control.RestartCmd = restartCmd

	ret, err := ps.CoreAPI.GseProcServer().RegisterProcInfo(context.Background(), forward, namespace, gseproc)
	if err != nil || (err == nil && !ret.Result) {
		return fmt.Errorf("register process(%s) into gse failed. err: %v, errcode: %d, errmsg: %s", procName, err, ret.Code, ret.ErrMsg)
	}

	return nil
}

func (ps *ProcServer) getHostForGse(appId, hostId string, forward http.Header) ([]meta.GseHost, error) {
	gseHosts := make([]meta.GseHost, 0)
	// get bk_supplier_id from applicationbase
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	reqParam := new(meta.QueryInput)
	reqParam.Condition = condition
	appRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, forward, reqParam)
	if err != nil || (err == nil && !appRet.Result) {
		return nil, fmt.Errorf("get application failed. condition: %+v, err: %v, errcode: %d, errmsg: %s", reqParam, err, appRet.Code, appRet.ErrMsg)
	}

	supplierID := 0
	if len(appRet.Data.Info) >= 1 {
		tmp, ok := appRet.Data.Info[0].Get(common.BKSupplierIDField)
		if !ok {
			return nil, fmt.Errorf("there is no supplierID in appID(%s)", appId)
		}

		supplierID, ok = tmp.(int)
		if !ok {
			return nil, fmt.Errorf("convert supplierID to int failed.")
		}
	}

	// get host info
	hostRet, err := ps.CoreAPI.HostController().Host().GetHostByID(context.Background(), hostId, forward)
	if err != nil || (err == nil && !hostRet.Result) {
		return nil, fmt.Errorf("get host by hostid(%s) failed. err: %v, errcode: %d, errmsg: %s", hostId, err, hostRet.Code, hostRet.ErrMsg)
	}

	hostIp, ok := hostRet.Data[common.BKHostInnerIPField].(string)
	if !ok {
		return nil, fmt.Errorf("convert host ip to string failed.")
	}

	cloudId, ok := hostRet.Data[common.BKCloudIDField].(int)
	if !ok {
		return nil, fmt.Errorf("convert cloudid to int failed")
	}

	var gseHost meta.GseHost
	gseHost.Ip = hostIp
	gseHost.BkCloudId = cloudId
	gseHost.BkSupplierId = supplierID

	gseHosts = append(gseHosts, gseHost)

	return gseHosts, nil
}

func (ps *ProcServer) createProcInstanceModel(appId, procId, moduleName, ownerId string, forward http.Header) error {
	u64AppId, err := strconv.ParseUint(appId, 10, 64)
	if err != nil {
		return fmt.Errorf("fail to parse appId into uint64. err: %v", err)
	}

	u64ProcId, err := strconv.ParseUint(procId, 10, 64)
	if err != nil {
		return fmt.Errorf("fail to parse procId into uint64. err: %v", err)
	}
	// get moduleid from cc_ModuleBase
	modIdCond := make(map[string]interface{})
	modIdCond[common.BKModuleNameField] = moduleName
	modIdCond[common.BKAppIDField] = appId
	modIdCond[common.BKOwnerIDField] = ownerId
	modIdSearchParam := new(meta.QueryInput)
	modIdSearchParam.Condition = modIdCond
	modIdRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDModule, forward, modIdSearchParam)
	if err != nil || (err == nil && !modIdRet.Result) {
		return fmt.Errorf("fail to search module info when create process instance. err: %v, errcode: %d, errmsg: %s", err, modIdRet.Code, modIdRet.ErrMsg)
	}

	setIdArr := make([]uint64, 0)
	modIdArr := make([]int64, 0)
	for _, modItem := range modIdRet.Data.Info {
		modId, ok := modItem[common.BKModuleIDField].(int64)
		if !ok {
			blog.Warnf("fail to convert module id to uint64. value: %v", modItem[common.BKModuleIDField])
		} else {
			modIdArr = append(modIdArr, modId)
		}

		setId, ok := modItem[common.BKSetIDField].(uint64)
		if !ok {
			blog.Warnf("fail to convert module id to uint64. value: %v", modItem[common.BKSetIDField])
		} else {
			setIdArr = append(setIdArr, setId)
		}
	}

	// get set name
	setId_Name := make(map[uint64]string)
	for _, setId := range setIdArr {
		setIdCond := make(map[string]interface{})
		setIdCond[common.BKSetIDField] = setId
		setIdCond[common.BKAppIDField] = appId
		setIdCond[common.BKOwnerIDField] = ownerId
		setIdSearchParam := new(meta.QueryInput)
		setIdRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDSet, forward, setIdSearchParam)
		if err != nil || (err == nil && !setIdRet.Result) {
			blog.Warnf("fail to search set info by condition(%+v), err: %v, errcode: %d, errmsg: %s", setIdSearchParam, err, setIdRet.Code, setIdRet.ErrMsg)
			continue
		}

		for _, setInfo := range setIdRet.Data.Info {
			setName, ok := setInfo[common.BKSetNameField].(string)
			if !ok {
				blog.Warnf("fail to convert setname to string. value(+v)", setInfo)
			} else {
				setId_Name[setId] = setName
			}
		}
	}

	// get hostid from cc_ModuleHostConfig
	hostModConf, err := ps.getConfigByCond(forward, map[string][]int64{
		common.BKModuleIDField: modIdArr,
	})

	// get proc from cc_Process
	procCond := make(map[string]interface{})
	procCond[common.BKOwnerIDField] = ownerId
	procCond[common.BKAppIDField] = appId
	procCond[common.BKProcIDField] = procId
	procParam := new(meta.QueryInput)
	procParam.Condition = procCond
	procRet, err := ps.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDProc, forward, procParam)
	if err != nil || (err == nil && !procRet.Result) {
		return fmt.Errorf("fail to search process object by search param(%+v). err: %v, errcode: %d, errmsg: %s", procParam, err, procRet.Code, procRet.ErrMsg)
	}

	// create procInstance
	procInstModels := make([]*meta.ProcInstanceModel, 0)
	for _, procInfo := range procRet.Data.Info {
		searchProcId, ok := procInfo[common.BKProcIDField].(string)
		if !ok {
			blog.Warnf("fail to convert procid into string. value: %+v", procInfo)
			continue
		}

		if searchProcId != procId {
			blog.Warnf("the processid(%s) get from db is not equal the procId(%s) in parameter", searchProcId, procId)
			continue
		}

		u64FuncId, ok := procInfo[common.BKFuncIDField].(uint64)
		if !ok {
			blog.Warnf("fail to convert funcId into uint64")
			continue
		}

		procInstNum, ok := procInfo[common.BKProcInstNum].(int)
		if !ok {
			blog.Warnf("convert process instance number to int failed")
			continue
		}

		for _, hostMod := range hostModConf {
			for i := 0; i < procInstNum; i++ {
				instModel := new(meta.ProcInstanceModel)
				instModel.ApplicationID = u64AppId
				instModel.ProcID = u64ProcId
				instModel.FuncID = u64FuncId
				instModel.ModuleName = moduleName
				instModel.SetID = uint64(hostMod[common.BKSetIDField])
				instModel.SetName = setId_Name[instModel.SetID]
				instModel.ModuleID = uint64(hostMod[common.BKModuleIDField])
				instModel.HostId = uint64(hostMod[common.BKHostIDField])
				instModel.InstanceID = 0 //TODO use unified id generation method, uri will provide

				procInstModels = append(procInstModels, instModel)
			}
		}
	}
	if 0 == len(procInstModels) {
		return nil
	}
	// save into db
	instModelRet, err := ps.CoreAPI.ProcController().CreateProcInstanceModel(context.Background(), forward, procInstModels)
	if err != nil || (err == nil && !instModelRet.Result) {
		return fmt.Errorf("fail to save process instance model into db. err: %v, errcode: %d, errmsg: %s", err, instModelRet.Code, instModelRet.ErrMsg)
	}

	return nil
}

func (ps *ProcServer) deleteProcInstanceModel(appId, procId, moduleName string, forward http.Header) error {
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	condition[common.BKProcIDField] = procId
	condition[common.BKModuleNameField] = moduleName

	ret, err := ps.CoreAPI.ProcController().DeleteProcInstanceModel(context.Background(), forward, condition)
	if err != nil || (err == nil && !ret.Result) {
		return fmt.Errorf("fail to delete process instance model. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	return nil
}

func (ps *ProcServer) getProcInstanceModel(appId string, forward http.Header) (map[string]*meta.ProcInstanceModel, error) {
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId

	ret, err := ps.CoreAPI.ProcController().GetProcInstanceModel(context.Background(), forward, condition)
	if err != nil || (err == nil && !ret.Result) {
		return nil, fmt.Errorf("fail to get process instance model from db. err: %v, errcode: %d, errmsg: %s", err, ret.Code, ret.ErrMsg)
	}

	result := make(map[string]*meta.ProcInstanceModel)
	for _, instModel := range ret.Data {
		key := fmt.Sprintf("%s.%s.%s.%s", instModel.SetName, instModel.ModuleName, instModel.FuncID, instModel.ModuleID)
		result[key] = &instModel
	}

	return result, nil
}

func (ps *ProcServer) matchProcessInstance(procOp *meta.ProcessOperate, forward http.Header) (map[string]*meta.ProcInstanceModel, error) {
	allProcInst, err := ps.getProcInstanceModel(procOp.ApplicationID, forward)
	if err != nil {
		blog.Errorf("match process instance failed! err: %v", err)
		return nil, err
	}

	pattern := ""

	if "*" == procOp.SetName {
		pattern = pattern + ".+?"
	} else {
		pattern = pattern + procOp.SetName
	}
	pattern = pattern + "\\."

	if "*" == procOp.ModuleName {
		pattern = pattern + ".+?"
	} else {
		pattern = pattern + procOp.ModuleName
	}
	pattern = pattern + "\\."

	if "*" == procOp.FuncID {
		pattern = pattern + ".+?"
	} else {
		pattern = pattern + procOp.FuncID
	}
	pattern = pattern + "\\."

	if "*" == procOp.InstanceID {
		pattern = pattern + ".+?"
	} else {
		pattern = pattern + procOp.InstanceID
	}

	result := make(map[string]*meta.ProcInstanceModel)
	for key, instModel := range allProcInst {
		bMatch, err := regexp.MatchString(pattern, key)
		if !bMatch || err != nil {
			blog.Warnf("cat not match(%s) by pattern(%s)", key, pattern)
		} else {
			result[key] = instModel
		}
	}

	return result, nil
}
