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

func (lgc *Logics) RegisterProcInstanceToGse(ctx context.Context, appID, moduleID, procID int64, gseHost []metadata.GseHost, procInfo map[string]interface{}) error {

	procName, ok := procInfo[common.BKProcessNameField].(string)
	if !ok {
		blog.Errorf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,proc info:%+v,rid:%s", appID, moduleID, procID, procInfo, lgc.rid)
		//	// CCErrCommInstFieldConvFail  convert %s  field %s to %s error %s
		return lgc.ccErr.Errorf(common.CCErrCommInstFieldConvFail, common.BKInnerObjIDProc, common.BKProcessNameField, "string", "not string")
	}

	// under value allow nil , can ignore error, value is empty
	pidFilePath, ok := procInfo[common.BKProcPidFile].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,rid:%s", appID, moduleID, procID, lgc.rid)
	}
	workPath, ok := procInfo[common.BKProcWorkPath].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,rid:%s", appID, moduleID, procID, lgc.rid)
	}
	startCmd, ok := procInfo[common.BKProcStartCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,rid:%s", appID, moduleID, procID, lgc.rid)
	}
	stopCmd, ok := procInfo[common.BKProcStopCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,rid:%s", appID, moduleID, procID, lgc.rid)
	}
	reloadCmd, ok := procInfo[common.BKProcReloadCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,rid:%s", appID, moduleID, procID, lgc.rid)
	}
	restartCmd, ok := procInfo[common.BKProcRestartCmd].(string)
	if !ok {
		blog.Warnf("register process to gse, appID(%d) moduleID(%d) procID(%d)  not found process name,rid:%s", appID, moduleID, procID, lgc.rid)
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
	return lgc.registerProcInstanceToGse(ctx, gseproc)

}

func (lgc *Logics) registerProcInstanceToGse(ctx context.Context, gseproc *metadata.GseProcRequest) error {
	defErr := lgc.ccErr
	namespace := getGseProcNameSpace(gseproc.AppID, gseproc.ModuleID)
	gseproc.Meta.Namespace = namespace
	ret, err := lgc.esbServ.GseSrv().RegisterProcInfo(ctx, lgc.header, gseproc)
	if err != nil {
		blog.Errorf("registerProcInstanceToGse RegisterProcInfo http do error.register process(%s) into gse failed. err: %v,input:%+v,rid:%s", gseproc.Meta.Name, err, gseproc, lgc.rid)
		return lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && 0 != ret.Code {
		blog.Errorf("registerProcInstanceToGse RegisterProcInfo http reply error.register process(%s) into gse failed. errcode: %d, errmsg: %s,input:%+v,rid:%s", gseproc.Meta.Name, ret.Code, ret.Message, gseproc, lgc.rid)
		return defErr.New(ret.Code, ret.Message)
	}

	gseproc.OpType = common.GSEProcOPRegister
	procInstDetail, err := lgc.CoreAPI.ProcController().RegisterProcInstanceDetail(ctx, lgc.header, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) detail failed. err: %v", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !procInstDetail.Result {
		blog.Errorf("register process(%s) detail failed.  errcode: %d, errmsg: %s", gseproc.Meta.Name, procInstDetail.Code, procInstDetail.ErrMsg)
		return defErr.New(procInstDetail.Code, procInstDetail.ErrMsg)

	}

	return nil
}

func (lgc *Logics) unregisterProcInstanceToGse(ctx context.Context, gseproc *metadata.GseProcRequest) error {
	namespace := getGseProcNameSpace(gseproc.AppID, gseproc.ModuleID)
	gseproc.Meta.Namespace = namespace
	gseproc.OpType = common.GSEProcOPUnregister
	defErr := lgc.ccErr
	ret, err := lgc.esbServ.GseSrv().UnRegisterProcInfo(ctx, lgc.header, gseproc)
	if err != nil {
		blog.Errorf("unregisterProcInstanceToGse UnRegisterProcInfo http do error.unregister process(%s) into gse failed.  errcode: %d, errmsg: %s,input:%+v,rid:%s", gseproc.Meta.Name, err, gseproc, lgc.rid)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if 0 != ret.Code {
		blog.Errorf("unregisterProcInstanceToGse UnRegisterProcInfo http reply error.register process(%s) into gse failed. errcode: %d, errmsg: %s,input:%+v,rid:%s", gseproc.Meta.Name, ret.Code, ret.Message, gseproc, lgc.rid)
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
	delConds := mapstr.MapStr{common.BKDBOR: unregisterProcDetail}
	procInstDetail, err := lgc.CoreAPI.ProcController().DeleteProcInstanceDetail(ctx, lgc.header, delConds)
	if err != nil {
		blog.Errorf("DeleteProcInstanceDetail DeleteProcInstanceDetail http do error.unregister process(%s) detail failed. err: %v,input:%+v,rid:%s", gseproc.Meta.Name, err, delConds, lgc.rid)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !procInstDetail.Result {
		blog.Errorf("DeleteProcInstanceDetail DeleteProcInstanceDetail http do error.unregister process(%s) detail failed. errcode: %d, errmsg: %s,input:%+v,rid:%s", gseproc.Meta.Name, procInstDetail.Code, procInstDetail.ErrMsg, delConds, lgc.rid)
		return defErr.New(procInstDetail.Code, procInstDetail.ErrMsg)
	}

	return nil
}

func (lgc *Logics) getOperateProcInstanceData(ctx context.Context, procOp *metadata.ProcessOperate, instModels map[string]*metadata.ProcInstanceModel) ([]*metadata.GseProcRequest, error) {

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
	procInfoArr, err := lgc.GetProcbyProcIDArr(ctx, allProcIDArr)
	if err != nil {
		blog.Warnf("OperateProcInstanceByGse getProcessByProcID failed. err: %s, rid:%s", err.Error(), lgc.rid)
		return nil, err
	}
	procInfoMap := make(map[int64]mapstr.MapStr, 0)
	for _, procInfo := range procInfoArr {
		procID, procIDErr := procInfo.Int64(common.BKProcessIDField)
		if procIDErr != nil {
			blog.Warnf("OperateProcInstanceByGse getProcessByProcID failed. err: %s, rid:%s", procIDErr.Error(), lgc.rid)
			continue
		}
		procInfoMap[procID] = procInfo
	}

	gseHostArr, err := lgc.GetHostForGse(ctx, procOp.ApplicationID, allHostIDArr)
	// register process into gse
	if err != nil {
		blog.Errorf("OperateProcInstanceByGse register process into gse failed. err: %v", err, lgc.rid)
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
			blog.Warnf("OperateProcInstanceByGse convert process name to string failed, error:%s, proc info:%+v,rid:%s", err.Error(), opGseProcInfo, lgc.rid)
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

func (lgc *Logics) OperateProcInstanceByGse(ctx context.Context, procOp *metadata.ProcessOperate, instModels map[string]*metadata.ProcInstanceModel) (string, error) {

	ccTaskID := getTaskID()
	opProcInsts := make([]*metadata.ProcessOperateTask, 0)

	gseReqArr, err := lgc.getOperateProcInstanceData(ctx, procOp, instModels)
	if nil != err {
		return "", err
	}

	mustNeedHeader := getMustNeedHeader(lgc.header)
	cacheTaskInfo := opProcTask{}
	cacheTaskInfo.TaskID = ccTaskID
	cacheTaskInfo.Header = mustNeedHeader
	for _, gseReq := range gseReqArr {
		gseRsp, err := lgc.esbServ.GseSrv().OperateProcess(ctx, lgc.header, gseReq)
		status := metadata.ProcOpTaskStatusWaitOP
		detail := make(map[string]metadata.ProcessOperateTaskDetail, 0)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse fail to operate process by gse process server. err: %v,input:%+v, rid:%s", err, gseReq, lgc.rid)
			status = metadata.ProcOpTaskStatusHTTPErr
			detail["http_request_error"] = metadata.ProcessOperateTaskDetail{
				Errcode: common.CCErrCommHTTPDoRequestFailed,
				ErrMsg:  err.Error(),
			}
		} else if !gseRsp.Result {
			blog.Errorf("OperateProcInstanceByGse fail to operate process by gse process server. errcode: %d, errmsg: %s,input:%+v,rid:%s", gseRsp.Code, gseRsp.Message, gseReq, lgc.rid)
			status = metadata.ProcOpTaskStatusErr
			detail["gse_error_message"] = metadata.ProcessOperateTaskDetail{
				Errcode: gseRsp.Code,
				ErrMsg:  gseRsp.Message,
			}
		}

		taskID, ok := gseRsp.Data[common.BKGseTaskIDField].(string)
		if !ok {
			blog.Warnf("OperateProcInstanceByGse convert gse process operate taskid to string failed. value: %v,rid:%s", gseRsp.Data, lgc.rid)
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
		ret, err := lgc.CoreAPI.ProcController().AddOperateTaskInfo(ctx, lgc.header, opProcInsts)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse AddOperateTaskInfo http do  error:%s, input:%+v,rid:%s", err.Error(), procOp, lgc.rid)
			return "", lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !ret.Result {
			blog.Errorf("OperateProcInstanceByGse AddOperateTaskInfo  error:%s, input:%+v,rid:%s", ret.Result, procOp, lgc.rid)
			return "", lgc.ccErr.New(ret.Code, ret.ErrMsg)
		}
		cacheInfoStr, err := json.Marshal(cacheTaskInfo)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse cache OperateTaskInfo json marshal error:%s, logID:%s", err.Error(), lgc.rid)
			return "", err
		}
		_, err = lgc.cache.SAdd(common.RedisProcSrvQueryProcOPResultKey, string(cacheInfoStr)).Result()
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse cache TaskIDInfo  error:%s, logID:%s", err.Error(), lgc.rid)
			return "", lgc.ccErr.Errorf(common.CCErrCommUtilHandleFail, "redis sadd", err.Error())
		}
	}

	return ccTaskID, nil
}

func (lgc *Logics) QueryProcessOperateResult(ctx context.Context, taskID string) (succ, waitExec []string, exceErrMap map[string]string, err error) {

	waitExecArr, execErrMap, err := lgc.handleOPProcTask(ctx, taskID)
	if nil != err || 0 < len(execErrMap) || 0 < len(waitExecArr) {
		return nil, waitExecArr, execErrMap, err
	}

	dat := new(metadata.QueryInput)
	dat.Condition = mapstr.MapStr{common.BKTaskIDField: taskID}
	dat.Limit = common.BKNoLimit
	succ = make([]string, 0)
	waitExec = make([]string, 0)
	exceErrMap = make(map[string]string, 0)

	ret, err := lgc.CoreAPI.ProcController().SearchOperateTaskInfo(ctx, lgc.header, dat)
	dat.Start += dat.Limit
	if nil != err {
		blog.Errorf("QueryProcessOperateResult http do error. search task info taskID:%s  http do error:%s,input:%+v,rid:%s", taskID, err.Error(), dat, lgc.rid)
		return nil, nil, nil, lgc.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("QueryProcessOperateResult http reply error. search task info taskID:%s error:%s,input:%+v,rid:%s", taskID, ret.ErrMsg, dat, lgc.rid)
		return nil, nil, nil, lgc.ccErr.New(ret.Code, ret.ErrMsg)
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

func (lgc *Logics) GetHostForGse(ctx context.Context, appID int64, hostIDArr []int64) ([]*metadata.GseHost, error) {
	gseHosts := make([]*metadata.GseHost, 0)
	// get bk_supplier_id from applicationbase
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appID
	reqParam := new(metadata.QueryCondition)
	reqParam.Condition = condition
	defErr := lgc.ccErr
	appRet, err := lgc.CoreAPI.CoreService().Instance().ReadInstance(ctx, lgc.header, common.BKInnerObjIDApp, reqParam)
	if err != nil {
		blog.Errorf("GetHostForGse SearchObjects http do error.get application failed. condition: %+v, err: %v,,rid:%s", reqParam, err, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !appRet.Result {
		blog.Errorf("GetHostForGse SearchObjects http reply error.get application failed. condition: %+v,   errcode: %d, errmsg: %s,rid:%s", reqParam, appRet.Code, appRet.ErrMsg, lgc.rid)
		return nil, defErr.New(appRet.Code, appRet.ErrMsg)
	}

	supplierID := int64(0)
	if len(appRet.Data.Info) >= 1 {
		tmp, ok := appRet.Data.Info[0].Get(common.BKSupplierIDField)
		if !ok {
			blog.Errorf("there is no supplierID in appID(%d),input:%+v,app info:%+v,rid:%s", appID, appRet.Data.Info[0], reqParam, lgc.rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldNotFound, "supplierID", "application", err.Error())
		}
		supplierID, err = util.GetInt64ByInterface(tmp)
		if !ok {
			blog.Errorf("there is no supplierID in appID(%d) not integer,input:%+v,app info:%+v,rid:%s", appID, appRet.Data.Info[0], reqParam, lgc.rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "supplierID", "int", err.Error())
		}
	}
	hostQuery := new(metadata.QueryInput)
	hostQuery.Condition = mapstr.MapStr{common.BKHostIDField: mapstr.MapStr{common.BKDBIN: hostIDArr}}
	hostQuery.Limit = len(hostIDArr)
	// get host info
	hostRet, err := lgc.CoreAPI.HostController().Host().GetHosts(ctx, lgc.header, hostQuery)
	if err != nil {
		blog.Errorf("get host by hostid(%+v) failed. err: %v,input:%+v,rid:%s", hostIDArr, err, hostQuery, lgc.rid)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !hostRet.Result {
		blog.Errorf("get host by hostid(%+v) failed.   errcode: %d, errmsg: %s,input:%+v,rid:%s", hostIDArr, hostRet.Code, hostRet.ErrMsg, hostQuery, lgc.rid)
		return nil, defErr.New(hostRet.Code, hostRet.ErrMsg)
	}

	for _, hostInfo := range hostRet.Data.Info {
		hostIp, err := hostInfo.String(common.BKHostInnerIPField)
		if err != nil {
			blog.Errorf("get inner ip from host info error, host info:%+v,input:%+v,rid:%s", hostInfo, hostQuery, lgc.rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "innerIP", "string", err.Error())
		}

		cloudID, err := hostInfo.Int64(common.BKCloudIDField)
		if nil != err {
			blog.Errorf("get cloud id from host info error, host info:%+v,input:%+v,rid:%s", hostInfo, hostQuery, lgc.rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "cloud id", "int", err.Error())
		}

		hostID, err := hostInfo.Int64(common.BKHostIDField)
		if nil != err {
			blog.Errorf("get inner id from host info error, host info:%+v,input:%+v,rid:%s", hostInfo, hostQuery, lgc.rid)
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "host id", "int", err.Error())
		}

		gseHost := &metadata.GseHost{}
		gseHost.HostID = hostID
		gseHost.Ip = hostIp
		gseHost.BkCloudId = cloudID
		gseHost.BkSupplierId = supplierID

		gseHosts = append(gseHosts, gseHost)
	}

	return gseHosts, nil
}
