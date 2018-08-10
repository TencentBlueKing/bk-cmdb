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

	procInstModel, err := ps.Logics.MatchProcessInstance(context.Background(), procOpParam, req.Request.Header)
	if err != nil {
		blog.Errorf("match process instance failed in OperateProcessInstance. err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrProcOperateFaile)})
		return
	}

	result, err := ps.Logics.OperateProcInstanceByGse(procOpParam, procInstModel, namespace, forward)
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

func (ps *ProcServer) operateProcInstanceByGse_deletefunc(procOp *meta.ProcessOperate, instModels map[string]*meta.ProcInstanceModel, namespace string, forward http.Header) (map[string]string, error) {
	var err error
	model_TaskId := make(map[string]string)
	for key, model := range instModels {
		gseprocReq := new(meta.GseProcRequest)
		// hostInfo
		gseprocReq.Hosts, err = ps.Logics.GetHostForGse(string(model.ApplicationID), string(model.HostID), forward)
		if err != nil {
			blog.Warnf("getHostForGse failed. err: %v", err)
			continue
		}

		// get processinfo
		procInfo, err := ps.Logics.GetProcessbyProcID(string(model.ProcID), forward)
		if err != nil {
			blog.Warnf("getProcessByProcID failed. err: %v", err)
			continue
		}

		// register process into gse
		if err := ps.Logics.RegisterProcInstanceToGse(namespace, procInfo, forward); err != nil {
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

func (ps *ProcServer) registerProcInstanceToGse_deletefunc(namespace string, procInfo map[string]interface{}, forward http.Header) error {
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

func (ps *ProcServer) getHostForGse_deletefunc(appId, hostId string, forward http.Header) ([]meta.GseHost, error) {
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

	supplierID := int64(0)
	if len(appRet.Data.Info) >= 1 {
		tmp, ok := appRet.Data.Info[0].Get(common.BKSupplierIDField)
		if !ok {
			return nil, fmt.Errorf("there is no supplierID in appID(%s)", appId)
		}

		supplierID, err = util.GetInt64ByInterface(tmp)
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

	cloudId, err := util.GetInt64ByInterface(hostRet.Data[common.BKCloudIDField])
	if nil != err {
		return nil, fmt.Errorf("convert cloudid to int failed")
	}

	var gseHost meta.GseHost
	gseHost.Ip = hostIp
	gseHost.BkCloudId = cloudId
	gseHost.BkSupplierId = supplierID

	gseHosts = append(gseHosts, gseHost)

	return gseHosts, nil
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
	// TODO use mongodb
	return nil, nil
}
