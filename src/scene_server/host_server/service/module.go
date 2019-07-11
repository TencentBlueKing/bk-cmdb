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
	"configcenter/src/auth"
	"encoding/json"
	"fmt"
	"net/http"

	authmeta "configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// HostModuleRelation transfer host to module specify by bk_module_id (in the same business)
// move a business host to a module.
func (s *Service) HostModuleRelation(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	config := new(metadata.HostsModuleRelation)
	if err := json.NewDecoder(req.Request.Body).Decode(config); err != nil {
		blog.Errorf("add host and module relation failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	for _, moduleID := range config.ModuleID {
		module, err := srvData.lgc.GetNormalModuleByModuleID(srvData.ctx, config.ApplicationID, moduleID)
		if err != nil {
			blog.Errorf("add host and module relation, but get module with id[%d] failed, err: %v,param:%+v,rid:%s", moduleID, err, config, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
			return
		}

		if len(module) == 0 {
			blog.Errorf("add host and module relation, but get empty module with id[%d],input:%+v,rid:%s", moduleID, config, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
			return
		}
	}

	audit := srvData.lgc.NewHostModuleLog(config.HostID)
	if err := audit.WithPrevious(srvData.ctx); err != nil {
		blog.Errorf("host module relation, get prev module host config failed, err: %v,param:%+v,rid:%s", err, config, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.MoveBizHostToModule, config.HostID...); err != nil {
		blog.Errorf("check move host to module authorization failed, hosts: %+v, err: %v", config.HostID, err)
		if err != auth.NoAuthorizeError {
			resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		}
		resp.WriteEntity(s.AuthManager.GenEditBizHostNoPermissionResp(config.HostID))
		return
	}
	// auth: deregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, config.HostID...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v", config.HostID, err)
		resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	result, err := s.CoreAPI.CoreService().Host().TransferHostModule(srvData.ctx, srvData.header, config)
	if err != nil {
		blog.Errorf("add host module relation, but add config failed, err: %v, %v,input:%+v,rid:%s", err, result.ErrMsg, config, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("add host module relation, but add config failed, err: %v, %v.input:%+v,rid:%s", err, result.ErrMsg, config, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg), Data: result.Data})
		return
	}

	if err := audit.SaveAudit(srvData.ctx, config.ApplicationID, srvData.user, ""); err != nil {
		blog.Errorf("host module relation, save audit log failed, err: %v,input:%+v,rid:%s", err, config, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, config.HostID...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v", config.HostID, err)
		resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) MoveHost2EmptyModule(req *restful.Request, resp *restful.Response) {
	s.moveHostToModuleByName(req, resp, common.DefaultResModuleName)
}

func (s *Service) MoveHost2FaultModule(req *restful.Request, resp *restful.Response) {
	s.moveHostToModuleByName(req, resp, common.DefaultFaultModuleName)
}

func (s *Service) MoveHostToResourcePool(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("move host to resource pool failed with decode body err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if 0 == len(conf.HostID) {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.MoveHostFromModuleToResPool, conf.HostID...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", conf.HostID, err)
		if err != auth.NoAuthorizeError {
			resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(s.AuthManager.GenMoveBizHostToResourcePoolNoPermissionResp(conf.HostID))
		return
	}
	// auth: deregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, conf.HostID...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v", conf.HostID, err)
		resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}
	exceptionArr, err := srvData.lgc.MoveHostToResourcePool(srvData.ctx, conf)
	if err != nil {
		blog.Errorf("move host to resource pool failed, err:%s, input:%#v, rid:%s", err.Error(), conf, srvData.rid)
		resp.WriteEntity(&metadata.RespError{Msg: err, Data: exceptionArr})
		return
	}

	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, conf.HostID...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v", conf.HostID, err)
		resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

// AssignHostToApp transfer host from resource pool to idle module
func (s *Service) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("assign host to app failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	// auth: check target business update priority
	// if err := s.AuthManager.AuthorizeByBusinessID(srvData.ctx, srvData.header, authmeta.Update, conf.ApplicationID); err != nil {
	// 	blog.Errorf("AssignHostToApp failed, authorize on business update failed, business: %d, err: %v, rid:%s", conf.ApplicationID, err, srvData.rid)
	// 	resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }
	//
	// // auth: check host transfer priority
	// if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, authmeta.MoveResPoolHostToBizIdleModule, conf.HostID...); err != nil {
	// 	blog.Errorf("AssignHostToApp failed, authorize on host transfer failed, hosts: %+v, err: %v,rid:%s", conf.HostID, err, srvData.rid)
	// 	resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
	// 	return
	// }

	// auth: deregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, conf.HostID...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v, rid:%s", conf.HostID, err, srvData.rid)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	exceptionArr, err := srvData.lgc.AssignHostToApp(srvData.ctx, conf)
	if err != nil {
		blog.Errorf("assign host to app, but assign to app http do error. err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err, Data: exceptionArr})
	}

	// register host to new business
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, conf.HostID...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v, rid:%s", conf.HostID, err, srvData.rid)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) AssignHostToAppModule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	data := new(metadata.HostToAppModule)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("assign host to app module failed with decode body err: %v, rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	appID, _, moduleID, err := srvData.lgc.GetTopoIDByName(srvData.ctx, data)
	if nil != err {
		blog.Errorf("get app  topology id by name error:%s, msg: applicationName:%s, setName:%s, moduleName:%s,input;%+v,rid:%s", err.Error(), data.AppName, data.SetName, data.ModuleName, data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrAddHostToModuleFailStr, "search application module not found ")})
		return
	}

	if 0 == appID || 0 == moduleID {
		// get default app
		ownerAppID, err := srvData.lgc.GetDefaultAppID(srvData.ctx)
		if err != nil {
			blog.Errorf("assign host to app module, but get resource pool failed, err: %v,input:%+v,rid:%s", err, data, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
			return
		}
		if 0 == ownerAppID {
			blog.Errorf("assign host to app module, but get resource pool failed, err: %v,input:%+v,rid:%s", err, data, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrAddHostToModuleFailStr)})
			return
		}

		// get idle module
		mConds := mapstr.New()
		mConds.Set(common.BKDefaultField, common.DefaultResModuleFlag)
		mConds.Set(common.BKModuleNameField, common.DefaultResModuleName)
		mConds.Set(common.BKAppIDField, ownerAppID)
		ownerModuleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, mConds)
		if nil != err {
			blog.Errorf("assign host to app module, but get unused host pool failed, ownerid[%v], err: %v,input:%+v,param:%+v,rid:%s", ownerModuleID, err, data, mConds, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrAddHostToModuleFailStr, err.Error())})
			return
		}
		appID = ownerAppID
		moduleID = ownerModuleID
		data.AppName = common.DefaultAppName
		data.SetName = ""
		data.ModuleName = common.DefaultResModuleName

	}

	// TODO hoffer host can not exist, not exist create
	// check authorization
	hostIDArr := make([]int64, 0)
	existNewAddHost := false
	for _, ip := range data.Ips {
		hostID, err := s.ip2hostID(srvData, ip, data.PlatID)
		if err != nil {
			blog.Errorf("invalid ip:%v, err:%v, rid:%s", ip, err, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrAddHostToModuleFailStr, err.Error())})
			return
		}

		if hostID == 0 {
			existNewAddHost = true
			continue
		}

		hostIDArr = append(hostIDArr, hostID)
	}

	// auth: check authorization
	if existNewAddHost == true {
		/*
			// 检查注册到资源池的权限
			if err := s.AuthManager.AuthorizeAddToResourcePool(srvData.ctx, srvData.header); err != nil {
				blog.Errorf("check host authorization for add to resource pool failed, err: %v", err)
				resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
				return
			}
		*/
		// 检查转移主机到目标业务的权限
		// auth: check target business update priority
		// if err := s.AuthManager.AuthorizeByBusinessID(srvData.ctx, srvData.header, authmeta.Update, appID); err != nil {
		// 	blog.Errorf("AssignHostToApp failed, authorize on business update failed, business: %d, err: %v, rid:%s", appID, err, srvData.rid)
		// 	resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
		// 	return
		// }
	}

	// auth: deregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v", hostIDArr, err)
		resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	var errmsg []string
	for index, ip := range data.Ips {
		host := make(map[string]interface{})
		if index < len(data.HostName) {
			host[common.BKHostNameField] = data.HostName[index]
		}
		if "" != data.OsType {
			host[common.BKOSTypeField] = data.OsType
		}
		host[common.BKCloudIDField] = data.PlatID

		// dispatch to app
		err := srvData.lgc.EnterIP(srvData.ctx, util.GetOwnerID(req.Request.Header), appID, moduleID, ip, data.PlatID, host, data.IsIncrement)
		if nil != err {
			blog.Errorf("%s add host error: %s,input:%+v,rid:%s", ip, err.Error(), data, srvData.rid)
			errmsg = append(errmsg, fmt.Sprintf("%s add host error: %s", ip, err.Error()))
		}
	}
	if 0 == len(errmsg) {
		// auth: register hosts
		if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, hostIDArr...); err != nil {
			blog.Errorf("register host to iam failed, hosts: %+v, err: %v, errmsg:%#v, rid:%s", hostIDArr, err, errmsg, srvData.rid)
			resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
			return
		}
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	blog.Errorf("assign host to app module failed, err: %v,rid:%s", errmsg, srvData.rid)
	resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrAddHostToModuleFailStr)})
	return

}

// GetHostModuleRelation  query host and module relation,
// hostID can empty
func (s *Service) GetHostModuleRelation(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	data := new(metadata.HostModuleRelationParameter)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("Transfer host across business failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	var cond metadata.HostModuleRelationRequest
	if data.AppID != 0 {
		cond.ApplicationID = data.AppID
	}
	if len(data.HostID) > 0 {
		cond.HostIDArr = data.HostID
	}

	configArr, err := srvData.lgc.GetHostModuleRelation(srvData.ctx, cond)
	if err != nil {
		blog.Errorf("GetHostModuleRelation logcis err:%s,cond:%#v,rid:%s", err.Error(), cond, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(configArr))
	return
}

// TransferHostAcrossBusiness  Transfer host across business,
// delete old business  host and module relation
func (s *Service) TransferHostAcrossBusiness(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	data := new(metadata.TransferHostAcrossBusinessParameter)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("Transfer host across business failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	err := srvData.lgc.TransferHostAcrossBusiness(srvData.ctx, data.SrcAppID, data.DstAppID, data.HostID, data.DstModuleIDArr)
	if err != nil {
		blog.Errorf("TransferHostAcrossBusiness logcis err:%s,input:%#v,rid:%s", err.Error(), data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
	return
}

// DeleteHostFromBusiness delete host from business
// dangerous operation
func (s *Service) DeleteHostFromBusiness(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)
	data := new(metadata.DeleteHostFromBizParameter)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("DeleteHostFromBizParameter failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	exceptionArr, err := srvData.lgc.DeleteHostFromBusiness(srvData.ctx, data.AppID, data.HostIDArr)
	if err != nil {
		blog.Errorf("DeleteHostFromBusiness logcis err:%s,input:%#v,rid:%s", err.Error(), data, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err, Data: exceptionArr})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
	return
}

// move host to idle or fault module under the same business.
func (s *Service) moveHostToModuleByName(req *restful.Request, resp *restful.Response, moduleName string) {
	pheader := req.Request.Header
	srvData := s.newSrvComm(pheader)
	defErr := srvData.ccErr
	ctx := srvData.ctx
	rid := srvData.rid
	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("move host to module %s failed with decode body err: %v,rid: %s", moduleName, err, rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	conds := make(map[string]interface{})
	var moduleNameLogKey string
	var action authmeta.Action
	if common.DefaultResModuleName == moduleName {
		// 空闲机
		moduleNameLogKey = "idle"
		conds[common.BKDefaultField] = common.DefaultResModuleFlag
		conds[common.BKModuleNameField] = common.DefaultResModuleName
		action = authmeta.MoveHostToBizIdleModule
	} else {
		// 故障机器
		moduleNameLogKey = "fault"
		conds[common.BKDefaultField] = common.DefaultFaultModuleFlag
		conds[common.BKModuleNameField] = common.DefaultFaultModuleName
		action = authmeta.MoveHostToBizFaultModule
	}
	conds[common.BKAppIDField] = conf.ApplicationID
	moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, conds)
	if err != nil {
		blog.Errorf("move host to module %s, get module id err: %v", moduleName, err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrAddHostToModuleFailStr, conds[common.BKModuleNameField].(string)+" not foud ")})
		return
	}

	audit := srvData.lgc.NewHostModuleLog(conf.HostID)
	if err := audit.WithPrevious(srvData.ctx); err != nil {
		blog.Errorf("move host to module %s, get prev module host config failed, err: %v", moduleName, err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, action, conf.HostID...); err != nil {
		blog.Errorf("auth host from iam failed, hosts: %+v, err: %v", conf.HostID, err)
		if err != auth.NoAuthorizeError {
			resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(s.AuthManager.GenEditBizHostNoPermissionResp(conf.HostID))
		return
	}
	// auth: deregister hosts
	if err := s.AuthManager.DeregisterHostsByID(srvData.ctx, srvData.header, conf.HostID...); err != nil {
		blog.Errorf("deregister host from iam failed, hosts: %+v, err: %v", conf.HostID, err)
		resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommUnRegistResourceToIAMFailed)})
		return
	}

	transferInput := &metadata.TransferHostToInnerModule{
		ApplicationID: conf.ApplicationID,
		HostID:        conf.HostID,
		ModuleID:      moduleID,
	}
	result, err := s.CoreAPI.CoreService().Host().TransferHostToInnerModule(ctx, pheader, transferInput)
	if err != nil {
		blog.Errorf("moveHostToModuleByName TransferHostToDefaultModule http do error. input:%#v,condition:%#v,err:%v,rid:%s", conf, transferInput, err.Error(), rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("moveHostToModuleByName TransferHostToDefaultModule http reply error. input:%#v,condition:%#v,err:%#v,rid:%s", conf, transferInput, result, rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.New(result.Code, result.ErrMsg), Data: result.Data})
		return
	}
	// auth: register hosts
	if err := s.AuthManager.RegisterHostsByID(srvData.ctx, srvData.header, conf.HostID...); err != nil {
		blog.Errorf("register host to iam failed, hosts: %+v, err: %v,rid:%s", conf.HostID, err, srvData.rid)
		resp.WriteEntity(&metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommRegistResourceToIAMFailed)})
		return
	}

	if err := audit.SaveAudit(srvData.ctx, conf.ApplicationID, srvData.user, "host to "+moduleNameLogKey+" module"); err != nil {
		blog.Errorf("move host to module %s, save audit log failed, err: %v,input:%+v,rid:%s", moduleName, err, conf, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}
