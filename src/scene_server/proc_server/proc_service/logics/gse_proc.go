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
	if 0 == len(procInfo) {
		return fmt.Errorf("process info is nil")
	}

	// process
	procName, ok := procInfo[common.BKProcessNameField].(string)
	if ok {
		return fmt.Errorf("process name not found")
	}
	// process
	procID, err := util.GetInt64ByInterface(procInfo[common.BKProcIDField])
	if nil != err {
		return fmt.Errorf("process id  not found")
	}

	// process
	appID, err := util.GetInt64ByInterface(procInfo[common.BKAppIDField])
	if nil != err {
		return fmt.Errorf("application id  not found")
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

	procInstDetail, err := lgc.CoreAPI.ProcController().RegisterProcInstanceDetail(context.Background(), forward, gseproc)
	if err != nil || (err == nil && !procInstDetail.Result) {
		return fmt.Errorf("register process(%s) detail failed. err: %v, errcode: %d, errmsg: %s", procName, err, procInstDetail.Code, procInstDetail.ErrMsg)
	}

	ret, err := lgc.CoreAPI.GseProcServer().RegisterProcInfo(context.Background(), forward, namespace, gseproc)
	if err != nil || (err == nil && !ret.Result) {
		return fmt.Errorf("register process(%s) into gse failed. err: %v, errcode: %d, errmsg: %s", procName, err, ret.Code, ret.ErrMsg)
	}

	return nil
}

func (lgc *Logics) OperateProcInstanceByGse(procOp *metadata.ProcessOperate, instModels map[string]*metadata.ProcInstanceModel, forward http.Header) (map[string]string, error) {
	var err error
	model_TaskId := make(map[string]string)
	for key, model := range instModels {
		gseprocReq := new(metadata.GseProcRequest)
		// hostInfo
		gseprocReq.Hosts, err = lgc.GetHostForGse(model.ApplicationID, model.HostID, forward)
		if err != nil {
			blog.Warnf("getHostForGse failed. err: %v", err)
			continue
		}

		// get processinfo
		procInfo, err := lgc.GetProcessbyProcID(string(model.ProcID), forward)
		if err != nil {
			blog.Warnf("getProcessByProcID failed. err: %v", err)
			continue
		}

		// register process into gse
		if err := lgc.RegisterProcInstanceToGse(model.ModuleID, gseprocReq.Hosts, procInfo, forward); err != nil {
			blog.Warnf("register process into gse failed. err: %v", err)
			continue
		}

		procName, ok := procInfo[common.BKProcessNameField].(string)
		if !ok {
			blog.Warnf("convert process name to string failed")
			continue
		}

		gseprocReq.Meta.Name = procName
		gseprocReq.Meta.Namespace = getGseProcNameSpace(model.ApplicationID, model.ModuleID)
		gseprocReq.OpType = procOp.OpType

		gseRsp, err := lgc.CoreAPI.GseProcServer().OperateProcess(context.Background(), forward, gseprocReq.Meta.Namespace, gseprocReq)
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

func (lgc *Logics) GetHostForGse(appId, hostId int64, forward http.Header) ([]metadata.GseHost, error) {
	gseHosts := make([]metadata.GseHost, 0)
	// get bk_supplier_id from applicationbase
	condition := make(map[string]interface{})
	condition[common.BKAppIDField] = appId
	reqParam := new(metadata.QueryInput)
	reqParam.Condition = condition
	appRet, err := lgc.CoreAPI.ObjectController().Instance().SearchObjects(context.Background(), common.BKInnerObjIDApp, forward, reqParam)
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
	hostRet, err := lgc.CoreAPI.HostController().Host().GetHostByID(context.Background(), string(hostId), forward)
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

	hostID, err := util.GetInt64ByInterface(hostRet.Data[common.BKHostIDField])
	if nil != err {
		return nil, fmt.Errorf("convert hostID to int failed")
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
	parseConds, notParse, err := ParseProcInstMatchCondition(match, true)
	if nil != err {
		blog.Errorf("matchName  parse set regex %s error %s", match, err.Error())
		return nil, nil, err
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
		return nil, nil, fmt.Errorf("get % error:%s", objID, err.Error())
	}
	if !ret.Result {
		blog.Errorf("matchName get %s instance error:%s", objID, ret.ErrMsg)
		return nil, nil, fmt.Errorf("get % error:%s", objID, ret.ErrMsg)
	}
	var rangeRegexrole *RegexRole
	if notParse {
		rangeRegexrole, err = NewRegexRole(match, true)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", match, err.Error())
			return nil, nil, fmt.Errorf("regex role %s parse error %s", match, err.Error())
		}
	}
	for _, inst := range ret.Data.Info {
		ID, err := inst.Int64(instIDKey)
		if nil != err {
			blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instIDKey)
			return nil, nil, fmt.Errorf("get %s id by int64 error", objID)
		}
		if notParse {
			name, err := inst.String(instNameKey)
			if nil != err {
				blog.Errorf("matchName %s info %v get key %s by int error", objID, inst, instNameKey)
				return nil, nil, fmt.Errorf("get %s id by int64 error", objID)
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

	funcIDConds, funcIDNotParse, err := ParseProcInstMatchCondition(funcIDMath, false)
	if nil != err {
		blog.Errorf("matchID  parse funcID regex %s error %s", funcIDMath, err.Error())
		return nil, err
	}
	var funcRegexRole *RegexRole
	if funcIDNotParse {
		funcRegexRole, err = NewRegexRole(funcIDMath, false)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", funcIDMath, err.Error())
			return nil, fmt.Errorf("regex role %s parse error %s", funcIDMath, err.Error())
		}
	}
	hostConds, hostIDNotParse, err := ParseProcInstMatchCondition(HostIDMatch, false)
	if nil != err {
		blog.Errorf("matchID  parse host instance id regex %s error %s", HostIDMatch, err.Error())
		return nil, err
	}
	var hostRegexRole *RegexRole
	if funcIDNotParse {
		hostRegexRole, err = NewRegexRole(HostIDMatch, false)
		if nil != err {
			blog.Errorf("regex role %s parse error %s", HostIDMatch, err.Error())
			return nil, fmt.Errorf("regex role %s parse error %s", HostIDMatch, err.Error())
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
		return nil, fmt.Errorf("get set error:%s", err.Error())
	}
	if !ret.Result {
		blog.Errorf("matchID get set instance error:%s", ret.ErrMsg)
		return nil, fmt.Errorf("get set error:%s", ret.ErrMsg)
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
