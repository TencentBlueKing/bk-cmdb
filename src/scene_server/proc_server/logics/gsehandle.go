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
	"errors"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// opGseProcInfo op gse proc data
type opGseProcInfo struct {
	ProcID    int64
	ModuleID  int64
	HostIDArr []int64
}

func (lgc *Logics) RegisterProcInstanceToGse(appID, moduleID, procID int64, gseHost []metadata.GseHost, procInfo map[string]interface{}, header http.Header) error {

	procName, ok := procInfo[common.BKProcessNameField].(string)
	if !ok {
		blog.Errorf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
		return errors.New("process name not found")
	}

	// under value allow nil , can ignore error, value is empty
	pidFilePath, ok := procInfo[common.BKProcPidFile].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
	}
	workPath, ok := procInfo[common.BKProcWorkPath].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
	}
	startCmd, ok := procInfo[common.BKProcStartCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
	}
	stopCmd, ok := procInfo[common.BKProcStopCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
	}
	reloadCmd, ok := procInfo[common.BKProcReloadCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
	}
	restartCmd, ok := procInfo[common.BKProcRestartCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name", appID, moduleID, procID)
	}

	namespace := getGseProcNameSpace(appID, moduleID)
	gseproc := new(metadata.GseProcRequest)
	gseproc.AppID = appID
	gseproc.ProcID = procID
	gseproc.Hosts = gseHost
	gseproc.ModuleID = moduleID
	gseproc.Meta.Namespace = namespace
	gseproc.Meta.Name = procName
	gseproc.Spec.Identity.PidPath = pidFilePath
	gseproc.Spec.Identity.SetupPath = workPath
	gseproc.Spec.Control.StartCmd = startCmd
	gseproc.Spec.Control.StopCmd = stopCmd
	gseproc.Spec.Control.ReloadCmd = reloadCmd
	gseproc.Spec.Control.RestartCmd = restartCmd
	return lgc.registerProcInstanceToGse(gseproc, header)

}

func (lgc *Logics) registerProcInstanceToGse(gseproc *metadata.GseProcRequest, header http.Header) error {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	namespace := getGseProcNameSpace(gseproc.AppID, gseproc.ModuleID)
	gseproc.Meta.Namespace = namespace
	ret, err := lgc.EsbServ.GseSrv().RegisterProcInfo(context.Background(), header, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) into gse failed. err: %v,", gseproc.Meta.Name, err)
		return lgc.CCErr.Error(util.GetLanguage(header), common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && 0 != ret.Code {
		blog.Errorf("register process(%s) into gse failed. errcode: %d, errmsg: %s", gseproc.Meta.Name, ret.Code, ret.Message)
		return defErr.New(ret.Code, ret.Message)
	}

	gseproc.OpType = common.GSEProcOPRegister
	procInstDetail, err := lgc.CoreAPI.ProcController().RegisterProcInstanceDetail(context.Background(), header, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) detail failed. err: %v", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !procInstDetail.Result {
		blog.Errorf("register process(%s) detail failed.  errcode: %d, errmsg: %s", gseproc.Meta.Name, procInstDetail.Code, procInstDetail.ErrMsg)
		return defErr.New(procInstDetail.Code, procInstDetail.ErrMsg)

	}

	return nil
}

func (lgc *Logics) unregisterProcInstanceToGse(gseproc *metadata.GseProcRequest, header http.Header) error {
	namespace := getGseProcNameSpace(gseproc.AppID, gseproc.ModuleID)
	gseproc.Meta.Namespace = namespace
	gseproc.OpType = common.GSEProcOPUnregister
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	ret, err := lgc.EsbServ.GseSrv().UnRegisterProcInfo(context.Background(), header, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) into gse failed. err: %s", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if 0 != ret.Code {
		blog.Errorf("register process(%s) into gse failed. errcode: %d, errmsg: %s", gseproc.Meta.Name, ret.Code, ret.Message)
		return defErr.New(ret.Code, ret.Message)
	}

	unregisterProcDetail := make([]interface{}, 0)
	for _, host := range gseproc.Hosts {
		unregisterProcDetail = append(unregisterProcDetail, mapstr.MapStr{
			common.BKAppIDField:     gseproc.AppID,
			common.BKModuleIDField:  gseproc.ModuleID,
			common.BKHostIDField:    host.HostID,
			common.BKProcessIDField: gseproc.ProcID})
	}

	procInstDetail, err := lgc.CoreAPI.ProcController().DeleteProcInstanceDetail(context.Background(), header, mapstr.MapStr{common.BKDBOR: unregisterProcDetail})
	if err != nil {
		blog.Errorf("register process(%s) detail failed. err: %v", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !procInstDetail.Result {
		blog.Errorf("register process(%s) detail failed. errcode: %d, errmsg: %s", gseproc.Meta.Name, procInstDetail.Code, procInstDetail.ErrMsg)
		return defErr.New(procInstDetail.Code, procInstDetail.ErrMsg)
	}

	return nil
}

func (lgc *Logics) getOperateProcInstanceData(ctx context.Context, procOp *metadata.ProcessOperate, instModels map[string]*metadata.ProcInstanceModel, header http.Header) ([]*metadata.GseProcRequest, error) {

	allHostIDArr := make([]int64, 0)
	allProcIDArr := make([]int64, 0)
	opGseProcInstMap := make(map[string]*opGseProcInfo, 0)
	for _, model := range instModels {
		key := getGseOpInstKey(model.ModuleID, model.ProcID)
		gseProcInst, ok := opGseProcInstMap[key]
		if !ok {
			gseProcInst = &opGseProcInfo{
				ProcID:   model.ProcID,
				ModuleID: model.ModuleID,
			}
			opGseProcInstMap[key] = gseProcInst
		}
		allHostIDArr = append(allHostIDArr, model.HostID)
		allProcIDArr = append(allProcIDArr, model.ProcID)
		gseProcInst.HostIDArr = append(gseProcInst.HostIDArr, model.HostID)
	}

	// get processinfo
	procInfoArr, err := lgc.GetProcbyProcIDArr(ctx, allProcIDArr, header)
	if err != nil {
		blog.Warnf("OperateProcInstanceByGse getProcessByProcID failed. err: %s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
		return nil, err
	}
	procInfoMap := make(map[int64]mapstr.MapStr, 0)
	for _, procInfo := range procInfoArr {
		procID, procIDErr := procInfo.Int64(common.BKProcessIDField)
		if procIDErr != nil {
			blog.Warnf("OperateProcInstanceByGse getProcessByProcID failed. err: %s, logID:%s", procIDErr.Error(), util.GetHTTPCCRequestID(header))
			continue
		}
		procInfoMap[procID] = procInfo
	}

	gseHostArr, err := lgc.GetHostForGse(ctx, procOp.ApplicationID, allHostIDArr, header)
	// register process into gse
	if err != nil {
		blog.Errorf("OperateProcInstanceByGse register process into gse failed. err: %v, rid:%s", err, util.GetHTTPCCRequestID(header))
		return nil, err
	}
	hostInfoMap := make(map[int64]*metadata.GseHost, 0)
	for _, hostInfo := range gseHostArr {
		hostInfoMap[hostInfo.HostID] = hostInfo
	}

	gseReqArr := make([]*metadata.GseProcRequest, 0)
	for _, opGseProcInfo := range opGseProcInstMap {
		procInfo, ok := procInfoMap[opGseProcInfo.ProcID]
		if !ok {
			continue
		}
		hostInfoArr := make([]*metadata.GseHost, 0)
		for _, hostID := range opGseProcInfo.HostIDArr {
			hostInfo, ok := hostInfoMap[hostID]
			if !ok {
				continue
			}
			hostInfoArr = append(hostInfoArr, hostInfo)
		}
		procName, err := procInfo.String(common.BKProcessNameField)
		if nil != err {
			blog.Warnf("OperateProcInstanceByGse convert process name to string failed, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
			continue
		}
		gseprocReq := new(metadata.GseProcRequest)
		gseprocReq.Meta.Name = procName
		gseprocReq.Meta.Namespace = getGseProcNameSpace(procOp.ApplicationID, opGseProcInfo.ModuleID)
		gseprocReq.OpType = procOp.OpType
		gseReqArr = append(gseReqArr, gseprocReq)
	}

	return gseReqArr, nil
}

func (lgc *Logics) OperateProcInstanceByGse(ctx context.Context, procOp *metadata.ProcessOperate, instModels map[string]*metadata.ProcInstanceModel, header http.Header) (string, error) {

	ccTaskID := getTaskID()
	opProcInsts := make([]*metadata.ProcessOperateTask, 0)

	gseReqArr, err := lgc.getOperateProcInstanceData(ctx, procOp, instModels, header)
	if nil != err {
		return "", err
	}

	mustNeedHeader := getMustNeedHeader(header)
	cacheTaskInfo := opProcTask{}
	cacheTaskInfo.TaskID = ccTaskID
	cacheTaskInfo.Header = mustNeedHeader
	for _, gseReq := range gseReqArr {
		gseRsp, err := lgc.EsbServ.GseSrv().OperateProcess(context.Background(), header, gseReq)
		status := metadata.ProcOpTaskStatusWaitOP
		detail := make(map[string]metadata.ProcessOperateTaskDetail, 0)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse fail to operate process by gse process server. err: %v, logID:%s", err, util.GetHTTPCCRequestID(header))
			status = metadata.ProcOpTaskStatusHTTPErr
			detail["http_request_error"] = metadata.ProcessOperateTaskDetail{
				Errcode: common.CCErrCommHTTPDoRequestFailed,
				ErrMsg:  err.Error(),
			}
		} else if !gseRsp.Result {
			blog.Errorf("OperateProcInstanceByGse fail to operate process by gse process server. errcode: %d, errmsg: %s, logID:%s", gseRsp.Code, gseRsp.Message, util.GetHTTPCCRequestID(header))
			status = metadata.ProcOpTaskStatusErr
			detail["gse_error_message"] = metadata.ProcessOperateTaskDetail{
				Errcode: gseRsp.Code,
				ErrMsg:  gseRsp.Message,
			}
		}

		taskID, ok := gseRsp.Data[common.BKGseTaskIDField].(string)
		if !ok {
			blog.Warnf("OperateProcInstanceByGse convert gse process operate taskid to string failed. value: %v", gseRsp.Data)
			status = metadata.ProcOpTaskStatusNotTaskIDErr
			detail["not_foud_gse_task_id"] = metadata.ProcessOperateTaskDetail{
				Errcode: common.CCErrCommNotFound,
				ErrMsg:  gseRsp.Message,
			}
		}

		opProcInsts = append(opProcInsts, &metadata.ProcessOperateTask{
			OperateInfo: procOp,
			TaskID:      ccTaskID,
			GseTaskID:   taskID,
			Namespace:   gseReq.Meta.Namespace,
			Status:      status,
			Host:        gseReq.Hosts,
			ProcName:    gseReq.Meta.Name,
			Detail:      detail,
		})
		cacheTaskInfo.GseTaskIDArr = append(cacheTaskInfo.GseTaskIDArr, taskID)
	}

	if 0 < len(opProcInsts) {
		ret, err := lgc.CoreAPI.ProcController().AddOperateTaskInfo(context.Background(), header, opProcInsts)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse AddOperateTaskInfo http do  error:%s, input:%v", err.Error(), procOp)
			return "", lgc.CCErr.Error(util.GetLanguage(header), common.CCErrCommHTTPDoRequestFailed)
		}
		if !ret.Result {
			blog.Errorf("OperateProcInstanceByGse AddOperateTaskInfo  error:%s, input:%v", ret.Result, procOp)
			return "", lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(ret.Code, ret.ErrMsg)
		}
		cacheInfoStr, err := json.Marshal(cacheTaskInfo)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse cache OperateTaskInfo json marshal error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
			return "", err
		}
		_, err = lgc.cache.SAdd(common.RedisProcSrvQueryProcOPResultKey, string(cacheInfoStr)).Result()
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse cache TaskIDInfo  error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(header))
			return "", err
		}
	}

	return ccTaskID, nil
}

func (lgc *Logics) QueryProcessOperateResult(ctx context.Context, taskID string, header http.Header) (succ, waitExec []string, exceErrMap map[string]string, err error) {

	waitExecArr, execErrMap, err := lgc.handleOPProcTask(ctx, header, taskID)
	if nil != err || 0 < len(execErrMap) || 0 < len(waitExecArr) {
		return nil, waitExecArr, execErrMap, err
	}

	dat := new(metadata.QueryInput)
	dat.Condition = mapstr.MapStr{common.BKTaskIDField: taskID}
	dat.Limit = common.BKNoLimit
	succ = make([]string, 0)

	ret, err := lgc.CoreAPI.ProcController().SearchOperateTaskInfo(ctx, header, dat)
	dat.Start += dat.Limit
	if nil != err {
		blog.Errorf("QueryProcessOperateResult http search task info taskID:%s  http do error:%s logID:%s", taskID, err.Error(), util.GetHTTPCCRequestID(header))
		return nil, nil, nil, lgc.CCErr.Error(util.GetLanguage(header), common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("QueryProcessOperateResult http search task info taskID:%s error:%s logID:%s", taskID, ret.ErrMsg, util.GetHTTPCCRequestID(header))
		return nil, nil, nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header)).New(ret.Code, ret.ErrMsg)
	}

	for _, item := range ret.Data.Info {
		for key, info := range item.Detail {
			switch info.Errcode {
			case int(metadata.ProcOpTaskStatusExecuteing):
				waitExecArr = append(waitExecArr, key)
			case int(metadata.ProcOpTaskStatusSucc):
				succ = append(succ, key)
			default:
				exceErrMap[key] = info.ErrMsg
			}
		}

		if 0 != len(exceErrMap) {
			return nil, nil, exceErrMap, nil
		}
		if 0 != len(exceErrMap) {
			return nil, waitExecArr, nil, err
		}
	}

	return succ, nil, nil, nil

}

func (lgc *Logics) GetHostForGse(ctx context.Context, appID int64, hostIDArr []int64, header http.Header) ([]*metadata.GseHost, error) {
	gseHosts := make([]*metadata.GseHost, 0)
	// get bk_supplier_id from applicationbase
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	reqParam := new(metadata.QueryInput)
	reqParam.Condition = condition
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(header))
	appRet, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, header, reqParam)
	if err != nil {
		blog.Errorf("get application failed. condition: %+v, err: %v ", reqParam, err)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !appRet.Result {
		blog.Errorf("get application failed. condition: %+v,   errcode: %d, errmsg: %s", reqParam, appRet.Code, appRet.ErrMsg)
		return nil, defErr.New(appRet.Code, appRet.ErrMsg)
	}

	supplierID := int64(0)
	if len(appRet.Data.Info) >= 1 {
		tmp, ok := appRet.Data.Info[0].Get(common.BKSupplierIDField)
		if !ok {
			blog.Errorf("there is no supplierID in appID(%d)", appID)
			return nil, defErr.Errorf(common.CCErrCommInstFieldNotFound, "supplierID", "application", err.Error())
		}
		supplierID, err = util.GetInt64ByInterface(tmp)
		if !ok {
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "supplierID", "int", err.Error())
		}
	}
	hostQuery := new(metadata.QueryInput)
	hostQuery.Condition = mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDArr}}
	hostQuery.Limit = len(hostIDArr)
	// get host info
	hostRet, err := lgc.CoreAPI.HostController().Host().GetHosts(context.Background(), header, hostQuery)
	if err != nil {
		blog.Errorf("get host by hostid(%d) failed. err: %v ", hostIDArr, err)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !hostRet.Result {
		blog.Errorf("get host by hostid(%d) failed.   errcode: %d, errmsg: %s", hostIDArr, hostRet.Code, hostRet.ErrMsg)
		return nil, defErr.New(hostRet.Code, hostRet.ErrMsg)
	}

	for _, hostInfo := range hostRet.Data.Info {
		hostIp, ok := hostInfo[common.BKHostInnerIPField].(string)
		if !ok {
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "innerIP", "string", err.Error())
		}

		cloudId, err := util.GetInt64ByInterface(hostInfo[common.BKCloudIDField])
		if nil != err {
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "cloud id", "int", err.Error())
		}

		hostID, err := util.GetInt64ByInterface(hostInfo[common.BKHostIDField])
		if nil != err {
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "host id", "int", err.Error())
		}

		gseHost := &metadata.GseHost{}
		gseHost.HostID = hostID
		gseHost.Ip = hostIp
		gseHost.BkCloudId = cloudId
		gseHost.BkSupplierId = supplierID

		gseHosts = append(gseHosts, gseHost)
	}

	return gseHosts, nil
}
