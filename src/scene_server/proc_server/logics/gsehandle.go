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
	"fmt"
	"net/http"

	"github.com/rs/xid"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (lgc *Logics) MatchProcessInstance(ctx context.Context, procOp *metadata.ProcessOperate, forward http.Header) (map[string]*metadata.ProcInstanceModel, error) {
	setConds := make(map[string]interface{})
	setConds[common.BKAppIDField] = procOp.ApplicationID
	setIDs, _, err := lgc.matchName(ctx, forward, procOp.SetName, common.BKInnerObjIDSet, common.BKSetIDField, common.BKSetNameField, setConds)
	if nil != err {
		blog.Errorf("MatchProcessInstance error:%s", err.Error())
		return nil, err
	}
	if 0 == len(setIDs) {
		return nil, nil
	}
	moduleConds := make(map[string]interface{}, 0)
	moduleConds[common.BKAppIDField] = procOp.ApplicationID
	moduleConds[common.BKSetIDField] = common.KvMap{common.BKDBIN: setIDs}
	moduleIDs, _, err := lgc.matchName(ctx, forward, procOp.ModuleName, common.BKInnerObjIDModule, common.BKModuleIDField, common.BKModuleNameField, moduleConds)
	if nil != err {
		blog.Errorf("MatchProcessInstance error:%s", err.Error())
		return nil, err
	}
	if 0 == len(moduleIDs) {
		return nil, nil
	}
	conds := make(map[string]interface{}, 0)
	conds[common.BKAppIDField] = procOp.ApplicationID
	conds[common.BKSetIDField] = common.KvMap{common.BKDBIN: setIDs}
	conds[common.BKModuleIDField] = common.KvMap{common.BKDBIN: moduleIDs}
	return lgc.matchID(ctx, forward, procOp.FuncID, procOp.InstanceID, conds)

}

func (lgc *Logics) RegisterProcInstanceToGse(moduleID int64, gseHost []metadata.GseHost, procInfo map[string]interface{}, forward http.Header) error {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward))
	if 0 == len(procInfo) {
		return defErr.Errorf(common.CCErrCommInstDataNil, "process ")
	}

	// process
	procName, ok := procInfo[common.BKProcessNameField].(string)
	if ok {
		blog.Errorf("process name not found, raw process:%+v", procInfo)
		return defErr.Errorf(common.CCErrCommInstFieldNotFound, "process name", "process")
	}
	// process
	procID, err := util.GetInt64ByInterface(procInfo[common.BKProcessIDField])
	if nil != err {
		blog.Errorf("process id not found, raw process:%+v", procInfo)
		return defErr.Errorf(common.CCErrCommInstFieldNotFound, "process id", "process")
	}

	// process
	appID, err := util.GetInt64ByInterface(procInfo[common.BKAppIDField])
	if nil != err {
		blog.Errorf("application id not found, raw process:%+v", procInfo)
		return defErr.Errorf(common.CCErrCommInstFieldNotFound, "application id", "process")
	}
	// under value allow nil , can ignore error
	pidFilePath, _ := procInfo[common.BKProcPidFile].(string)
	workPath, _ := procInfo[common.BKProcWorkPath].(string)
	startCmd, _ := procInfo[common.BKProcStartCmd].(string)
	stopCmd, _ := procInfo[common.BKProcStopCmd].(string)
	reloadCmd, _ := procInfo[common.BKProcReloadCmd].(string)
	restartCmd, _ := procInfo[common.BKProcRestartCmd].(string)

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
	return lgc.registerProcInstanceToGse(gseproc, forward)

}

func (lgc *Logics) registerProcInstanceToGse(gseproc *metadata.GseProcRequest, forward http.Header) error {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward))
	namespace := getGseProcNameSpace(gseproc.AppID, gseproc.ModuleID)
	gseproc.Meta.Namespace = namespace
	ret, err := lgc.CoreAPI.GseProcServer().RegisterProcInfo(context.Background(), forward, namespace, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) into gse failed. err: %v,", gseproc.Meta.Name, err)
		return lgc.CCErr.Error(util.GetLanguage(forward), common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && 0 != ret.Code {
		blog.Errorf("register process(%s) into gse failed. errcode: %d, errmsg: %s", gseproc.Meta.Name, ret.Code, ret.ErrMsg)
		return defErr.New(ret.Code, ret.ErrMsg)
	}

	gseproc.OpType = common.GSEProcOPRegister
	procInstDetail, err := lgc.CoreAPI.ProcController().RegisterProcInstanceDetail(context.Background(), forward, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) detail failed. err: %v", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !procInstDetail.Result {
		blog.Errorf("register process(%s) detail failed.  errcode: %d, errmsg: %s", gseproc.Meta.Name, procInstDetail.Code, procInstDetail.ErrMsg)
		return defErr.New(procInstDetail.Code, procInstDetail.ErrMsg)

	}

	return nil
}

func (lgc *Logics) unregisterProcInstanceToGse(gseproc *metadata.GseProcRequest, forward http.Header) error {
	namespace := getGseProcNameSpace(gseproc.AppID, gseproc.ModuleID)
	gseproc.Meta.Namespace = namespace
	gseproc.OpType = common.GSEProcOPUnregister
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward))
	ret, err := lgc.CoreAPI.GseProcServer().UnRegisterProcInfo(context.Background(), forward, namespace, gseproc)
	if err != nil {
		blog.Errorf("register process(%s) into gse failed.  errcode: %d, errmsg: %s", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if 0 != ret.Code {
		blog.Errorf("register process(%s) into gse failed. errcode: %d, errmsg: %s", gseproc.Meta.Name, ret.Code, ret.ErrMsg)
		return defErr.New(ret.Code, ret.ErrMsg)
	}

	unregisterProcDetail := make([]interface{}, 0)
	for _, host := range gseproc.Hosts {
		unregisterProcDetail = append(unregisterProcDetail, common.KvMap{
			common.BKAppIDField:     gseproc.AppID,
			common.BKModuleIDField:  gseproc.ModuleID,
			common.BKHostIDField:    host.HostID,
			common.BKProcessIDField: gseproc.ProcID})
	}

	procInstDetail, err := lgc.CoreAPI.ProcController().DeleteProcInstanceDetail(context.Background(), forward, common.KvMap{common.BKDBOR: unregisterProcDetail})
	if err != nil {
		blog.Errorf("register process(%s) detail failed. err: %v", gseproc.Meta.Name, err)
		return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !procInstDetail.Result {
		blog.Errorf("register process(%s) detail failed. errcode: %d, errmsg: %s", gseproc.Meta.Name, procInstDetail.Code, procInstDetail.ErrMsg)
		return defErr.New(procInstDetail.Code, procInstDetail.ErrMsg)
	}

	return nil
}

func (lgc *Logics) OperateProcInstanceByGse(procOp *metadata.ProcessOperate, instModels map[string]*metadata.ProcInstanceModel, forward http.Header) (string, error) {
	var err error

	ccTaskID := fmt.Sprintf("%s-%s", common.BKSTRIDPrefix, xid.New().String())
	opProcInsts := make([]*metadata.ProcessOperateTask, 0)
	for _, model := range instModels {
		gseprocReq := new(metadata.GseProcRequest)
		// hostInfo
		gseprocReq.Hosts, err = lgc.GetHostForGse(model.ApplicationID, model.HostID, forward)
		if err != nil {
			blog.Warnf("OperateProcInstanceByGse getHostForGse failed. err: %v", err)
			continue
		}

		// get processinfo
		procInfo, err := lgc.GetProcessbyProcID(string(model.ProcID), forward)
		if err != nil {
			blog.Warnf("OperateProcInstanceByGse getProcessByProcID failed. err: %v", err)
			continue
		}

		// register process into gse
		if err := lgc.RegisterProcInstanceToGse(model.ModuleID, gseprocReq.Hosts, procInfo, forward); err != nil {
			blog.Warnf("OperateProcInstanceByGse register process into gse failed. err: %v", err)
			continue
		}

		procName, ok := procInfo[common.BKProcessNameField].(string)
		if !ok {
			blog.Warnf("OperateProcInstanceByGse convert process name to string failed")
			continue
		}

		gseprocReq.Meta.Name = procName
		gseprocReq.Meta.Namespace = getGseProcNameSpace(model.ApplicationID, model.ModuleID)
		gseprocReq.OpType = procOp.OpType

		gseRsp, err := lgc.CoreAPI.GseProcServer().OperateProcess(context.Background(), forward, gseprocReq.Meta.Namespace, gseprocReq)
		if err != nil || (err == nil || 0 != gseRsp.Code) {
			blog.Warnf("OperateProcInstanceByGse fail to operate process by gse process server. err: %v, errcode: %d, errmsg: %s", err, gseRsp.Code, gseRsp.ErrMsg)
			continue
		}

		taskID, ok := gseRsp.Result[common.BKGseTaskIdField].(string)
		if !ok {
			blog.Warnf("OperateProcInstanceByGse convert gse process operate taskid to string failed. value: %v", gseRsp.Result)
			continue
		}

		opProcInsts = append(opProcInsts, &metadata.ProcessOperateTask{
			OperateInfo: procOp,
			TaskID:      ccTaskID,
			GseTaskID:   taskID,
			Namespace:   gseprocReq.Meta.Namespace,
			Status:      metadata.ProcessOperateTaskStatusWaitOP,
		})
	}

	if 0 < len(opProcInsts) {
		ret, err := lgc.CoreAPI.ProcController().AddOperateTaskInfo(context.Background(), forward, opProcInsts)
		if nil != err {
			blog.Errorf("OperateProcInstanceByGse AddOperateTaskInfo http do  error:%s, input:%v", err.Error(), procOp)
			return "", lgc.CCErr.Error(util.GetLanguage(forward), common.CCErrCommHTTPDoRequestFailed)
		}
		if !ret.Result {
			blog.Errorf("OperateProcInstanceByGse AddOperateTaskInfo  error:%s, input:%v", ret.Result, procOp)
			return "", lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward)).New(ret.Code, ret.ErrMsg)
		}
	}

	return ccTaskID, nil
}

func (lgc *Logics) QueryProcessOperateResult(ctx context.Context, taskID string, forward http.Header) (succ, waitExec []string, mapExceErr map[string]string, err error) {

	dat := new(metadata.QueryInput)
	dat.Limit = 200
	succ = make([]string, 0)
	waitExec = make([]string, 0)
	mapExceErr = make(map[string]string, 0)

	for {
		ret, err := lgc.CoreAPI.ProcController().SearchOperateTaskInfo(ctx, forward, dat)
		if nil != err {
			blog.Errorf("QueryProcessOperateResult http search task info taskID:%s  http do error:%s logID:%s", taskID, err.Error(), util.GetHTTPCCRequestID(forward))
			return nil, nil, nil, lgc.CCErr.Error(util.GetLanguage(forward), common.CCErrCommHTTPDoRequestFailed)
		}
		if !ret.Result {
			blog.Errorf("QueryProcessOperateResult http search task info taskID:%s error:%s logID:%s", taskID, ret.ErrMsg, util.GetHTTPCCRequestID(forward))
			return nil, nil, nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward)).New(ret.Code, ret.ErrMsg)
		}
		if 0 == ret.Data.Count {
			return nil, nil, nil, nil
		}
		for _, item := range ret.Data.Info {
			itemSucc, itemWaitExce, itemMapExecErr, err := lgc.handleGseTaskResult(ctx, &item, forward)
			if nil != err {
				return nil, nil, nil, err
			}
			if 0 != len(itemMapExecErr) {
				return nil, nil, itemMapExecErr, nil
			}
			if 0 != len(itemWaitExce) {
				return nil, itemWaitExce, nil, err
			}
			succ = append(succ, itemSucc...)
		}
	}
	return succ, nil, nil, nil

}

func (lgc *Logics) handleGseTaskResult(ctx context.Context, item *metadata.ProcessOperateTask, forward http.Header) (succ, waitExec []string, mapExceErr map[string]string, err error) {
	isWaiting := false
	isErr := false
	isChangeStatus := false
	taskStatusData := item.Detail
	succ = make([]string, 0)
	waitExec = make([]string, 0)
	mapExceErr = make(map[string]string, 0)
	if item.Status == metadata.ProcessOperateTaskStatusWaitOP || item.Status == metadata.ProcessOperateTaskStatusRuning {
		gseRet, err := lgc.CoreAPI.GseProcServer().QueryProcOperateResult(ctx, forward, item.Namespace, item.GseTaskID)
		if err != nil {
			blog.Errorf("QueryProcessOperateResult query task info from gse  error, taskID:%s, gseTaskID:%s, error:%s logID:%s", item.TaskID, item.GseTaskID, gseRet.ErrMsg, util.GetHTTPCCRequestID(forward))
			return nil, nil, nil, lgc.CCErr.Error(util.GetLanguage(forward), common.CCErrCommHTTPDoRequestFailed)
		} else if 0 != gseRet.Code {
			blog.Errorf("QueryProcessOperateResult query task info from gse failed,  taskID:%s, gseTaskID:%s, gse return error:%s, error code:%d logID:%s", item.TaskID, item.GseTaskID, gseRet.ErrMsg, gseRet.Code, util.GetHTTPCCRequestID(forward))
			return nil, nil, nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward)).New(gseRet.Code, gseRet.ErrMsg)
		}
		taskStatusData = gseRet.Result
		isChangeStatus = true
	}

	for key, item := range taskStatusData {
		if 0 == item.Errcode {
			succ = append(succ, key)
		} else {
			if int(metadata.ProcessOperateTaskStatusRuning) == item.Errcode {
				waitExec = append(waitExec, key)
				isWaiting = true
			} else {
				mapExceErr[key] = item.ErrMsg
				isErr = true
			}
		}
	}
	if isChangeStatus {
		updateConds := new(metadata.UpdateParams)
		data := common.KvMap{"detail": taskStatusData}
		updateConds.Condition = common.KvMap{"task_id": item.TaskID, "bk_task_id": item.GseTaskID}
		if isErr {
			data["status"] = metadata.ProcessOperateTaskStatusErr
		} else if isWaiting {
			data["status"] = metadata.ProcessOperateTaskStatusWaitOP
		} else {
			data["status"] = metadata.ProcessOperateTaskStatusSucc
		}
		updateConds.Data = data
		updateRet, err := lgc.CoreAPI.ProcController().UpdateOperateTaskInfo(ctx, forward, updateConds)
		if err != nil {
			blog.Errorf("QueryProcessOperateResult update task http do error, taskID:%s, gseTaskID:%s, error:%s logID:%s", item.TaskID, item.TaskID, item.GseTaskID, updateRet.ErrMsg, util.GetHTTPCCRequestID(forward))
			return nil, nil, nil, lgc.CCErr.Error(util.GetLanguage(forward), common.CCErrCommHTTPDoRequestFailed)
		} else if !updateRet.Result {
			blog.Errorf("QueryProcessOperateResult update task  reply error,  taskID:%s, gseTaskID:%s, gse return error:%s, error code:%d logID:%s", item.TaskID, item.GseTaskID, updateRet.ErrMsg, updateRet.Code, util.GetHTTPCCRequestID(forward))
			return nil, nil, nil, lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward)).New(updateRet.Code, updateRet.ErrMsg)
		}
	}
	return
}

func (lgc *Logics) GetHostForGse(appId, hostId int64, forward http.Header) ([]metadata.GseHost, error) {
	gseHosts := make([]metadata.GseHost, 0)
	// get bk_supplier_id from applicationbase
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	reqParam := new(metadata.QueryInput)
	reqParam.Condition = condition
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward))
	appRet, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, forward, reqParam)
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
			blog.Errorf("there is no supplierID in appID(%s)", appId)
			return nil, defErr.Errorf(common.CCErrCommInstFieldNotFound, "supplierID", "application")
		}
		supplierID, err = util.GetInt64ByInterface(tmp)
		if !ok {
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "application", "supplierID", "int", err.Error())
		}
	}

	// get host info
	hostRet, err := lgc.CoreAPI.HostController().Host().GetHostByID(context.Background(), string(hostId), forward)
	if err != nil {
		blog.Errorf("get host by hostid(%s) failed. err: %v ", hostId, err)
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	} else if err == nil && !hostRet.Result {
		blog.Errorf("get host by hostid(%s) failed.   errcode: %d, errmsg: %s", hostId, hostRet.Code, hostRet.ErrMsg)
		return nil, defErr.New(hostRet.Code, hostRet.ErrMsg)
	}

	hostIp, ok := hostRet.Data[common.BKHostInnerIPField].(string)
	if !ok {
		return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "innerIP", "string", err.Error())
	}

	cloudId, err := util.GetInt64ByInterface(hostRet.Data[common.BKCloudIDField])
	if nil != err {
		return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "cloud id", "int", err.Error())
	}

	hostID, err := util.GetInt64ByInterface(hostRet.Data[common.BKHostIDField])
	if nil != err {
		return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "host", "host id", "int", err.Error())
	}

	var gseHost metadata.GseHost
	gseHost.HostID = hostID
	gseHost.Ip = hostIp
	gseHost.BkCloudId = cloudId
	gseHost.BkSupplierId = supplierID

	gseHosts = append(gseHosts, gseHost)

	return gseHosts, nil
}

func (lgc *Logics) matchName(ctx context.Context, forward http.Header, match, objID, instIDKey, instNameKey string, conds map[string]interface{}) (instIDs []int64, data map[int64]mapstr.MapStr, err error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward))
	parseConds, notParse, err := ParseProcInstMatchCondition(match, true)
	if nil != err {
		blog.Errorf("matchName  parse set regex %s error %s", match, err.Error())
		return nil, nil, defErr.Errorf(common.CCErrCommUtilHandleFail, fmt.Sprintf("parse math %s", match), err.Error())
	}
	query := new(metadata.QueryInput)
	query.Limit = common.BKNoLimit
	if nil != parseConds {
		if nil == conds {
			conds = make(map[string]interface{})
		}
		conds[instNameKey] = parseConds

	}
	query.Condition = conds
	ret, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(ctx, objID, forward, query)
	if nil != err {
		blog.Errorf("matchName get %s instance error:%s", objID, err.Error())
		return nil, nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("matchName get %s instance error:%s", objID, ret.ErrMsg)
		return nil, nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	var rangeRegexrole *RegexRole
	if notParse {
		rangeRegexrole, err = NewRegexRole(match, true)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", match, err.Error())
			return nil, nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "match string", match, "regex role", err.Error())
		}
	}
	for _, inst := range ret.Data.Info {
		ID, err := inst.Int64(instIDKey)
		if nil != err {
			blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instIDKey)
			return nil, nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, objID, "inst id", "int", err.Error())
		}
		if notParse {
			name, err := inst.String(instNameKey)
			if nil != err {
				blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instNameKey)
				return nil, nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, objID, "inst name", "string", err.Error())
			}
			if rangeRegexrole.MatchStr(name) {
				instIDs = append(instIDs, ID)
				data[ID] = inst
			}
		} else {
			instIDs = append(instIDs, ID)
			data[ID] = inst
		}
	}

	return instIDs, data, nil
}

func (lgc *Logics) matchID(ctx context.Context, forward http.Header, funcIDMath, HostIDMatch string, conds map[string]interface{}) (data map[string]*metadata.ProcInstanceModel, err error) {
	defErr := lgc.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(forward))
	funcIDConds, funcIDNotParse, err := ParseProcInstMatchCondition(funcIDMath, false)
	if nil != err {
		blog.Errorf("matchID  parse funcID regex %s error %s", funcIDMath, err.Error())
		return nil, defErr.Errorf(common.CCErrCommUtilHandleFail, fmt.Sprintf("parse math %s", funcIDMath), err.Error())
	}
	var funcRegexRole *RegexRole
	if funcIDNotParse {
		funcRegexRole, err = NewRegexRole(funcIDMath, false)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", funcIDMath, err.Error())
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "match string", funcIDMath, "regex role", err.Error())
		}
	}
	hostConds, hostIDNotParse, err := ParseProcInstMatchCondition(HostIDMatch, false)
	if nil != err {
		blog.Errorf("matchID  parse host instance id regex %s error %s", HostIDMatch, err.Error())
		return nil, defErr.Errorf(common.CCErrCommUtilHandleFail, fmt.Sprintf("parse math %s", HostIDMatch), err.Error())
	}
	var hostRegexRole *RegexRole
	if funcIDNotParse {
		hostRegexRole, err = NewRegexRole(HostIDMatch, false)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", HostIDMatch, err.Error())
			return nil, defErr.Errorf(common.CCErrCommInstFieldConvFail, "match string", HostIDMatch, "regex role", err.Error())
		}
	}
	if nil != funcIDConds {
		if nil == conds {
			conds = make(map[string]interface{})
		}
		conds[common.BKFuncIDField] = funcIDConds

	}
	if nil != hostConds {
		if nil == conds {
			conds = make(map[string]interface{})
		}
		conds[common.BKFuncIDField] = hostConds

	}
	query := new(metadata.QueryInput)
	query.Limit = common.BKNoLimit
	query.Condition = conds
	ret, err := lgc.CoreAPI.ProcController().GetProcInstanceModel(ctx, forward, query)
	if nil != err {
		blog.Errorf("matchID get set instance error:%s", err.Error())
		return nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !ret.Result {
		blog.Errorf("matchID get set instance error:%s", ret.ErrMsg)
		return nil, defErr.New(ret.Code, ret.ErrMsg)
	}
	data = make(map[string]*metadata.ProcInstanceModel, 0)
	for _, item := range ret.Data.Info {
		isAdd := true
		if funcIDNotParse {
			isAdd = funcRegexRole.MatchInt64(item.FuncID)
		}
		if hostIDNotParse {
			isAdd = hostRegexRole.MatchInt64(int64(item.HostInstanID))
		}
		if isAdd {
			data[fmt.Sprintf("%d.%d.%d.%d", item.SetID, item.ModuleID, item.FuncID, item.HostID)] = &item
		}
	}

	return data, nil
}
