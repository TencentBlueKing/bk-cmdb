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

	"configcenter/src/ac"
	authmeta "configcenter/src/ac/meta"
	"configcenter/src/auth"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/logics"
)

// HostModuleRelation transfer host to module specify by bk_module_id (in the same business)
// move a business host to a module.
func (s *Service) TransferHostModule(ctx *rest.Contexts) {
	config := new(metadata.HostsModuleRelation)
	if err := ctx.DecodeInto(&config); nil != err {
		ctx.RespAutoError(err)
		return
	}

	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	for _, moduleID := range config.ModuleID {
		module, err := lgc.GetNormalModuleByModuleID(ctx.Kit.Ctx, config.ApplicationID, moduleID)
		if err != nil {
			blog.Errorf("add host and module relation, but get module with id[%d] failed, err: %v,param:%+v,rid:%s", moduleID, err, config, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}

		if len(module) == 0 {
			blog.Errorf("add host and module relation, but get empty module with id[%d],input:%+v,rid:%s", moduleID, config, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrTopoModuleIDNotfoundFailed))
			return
		}
	}

	audit := lgc.NewHostModuleLog(config.HostID)
	if err := audit.WithPrevious(ctx.Kit.Ctx); err != nil {
		blog.Errorf("host module relation, get prev module host config failed, err: %v,param:%+v,rid:%s", err, config, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.MoveBizHostToModule, config.HostID...); err != nil {
		blog.Errorf("check move host to module authorization failed, hosts: %+v, err: %v", config.HostID, err)
		if err != ac.NoAuthorizeError {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, config.HostID)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, auth.NoAuthorizeError)
		return
	}

	var result *metadata.OperaterException
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		result, err = s.CoreAPI.CoreService().Host().TransferToNormalModule(ctx.Kit.Ctx, ctx.Kit.Header, config)
		if err != nil {
			blog.Errorf("add host module relation, but add config failed, err: %v, %v,input:%+v,rid:%s", err, result.ErrMsg, config, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		
		if !result.Result {
			blog.Errorf("add host module relation, but add config failed, err: %v, %v.input:%+v,rid:%s", err, result.ErrMsg, config, ctx.Kit.Rid)
			return ctx.Kit.CCError.New(result.Code, result.ErrMsg)
		}

		if err := audit.SaveAudit(ctx.Kit.Ctx); err != nil {
			blog.Errorf("host module relation, save audit log failed, err: %v,input:%+v,rid:%s", err, config, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(result.Data,txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) MoveHost2IdleModule(ctx *rest.Contexts) {
	s.moveHostToDefaultModule(ctx, common.DefaultResModuleFlag)
}

func (s *Service) MoveHost2FaultModule(ctx *rest.Contexts) {
	s.moveHostToDefaultModule(ctx, common.DefaultFaultModuleFlag)
}

func (s *Service) MoveHost2RecycleModule(ctx *rest.Contexts) {
	s.moveHostToDefaultModule(ctx, common.DefaultRecycleModuleFlag)
}

func (s *Service) MoveHostToResourcePool(ctx *rest.Contexts) {
	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := ctx.DecodeInto(&conf); nil != err {
		ctx.RespAutoError(err)
		return
	}


	if 0 == len(conf.HostIDs) {
		ctx.RespEntity(nil)
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.MoveHostFromModuleToResPool, conf.HostIDs...); err != nil {
		blog.Errorf("check host authorization failed, hosts: %+v, err: %v", conf.HostIDs, err)
		if err != ac.NoAuthorizeError {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenMoveBizHostToResPoolNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, conf.HostIDs)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, auth.NoAuthorizeError)
		return
	}

	var exceptionArr []metadata.ExceptionResult
	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		exceptionArr, err = lgc.MoveHostToResourcePool(ctx.Kit.Ctx, conf)
		if err != nil {
			blog.Errorf("move host to resource pool failed, err:%s, input:%#v, rid:%s", err.Error(), conf, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr,txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// AssignHostToApp transfer host from resource pool to idle module
func (s *Service) AssignHostToApp(ctx *rest.Contexts) {

	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := ctx.DecodeInto(&conf); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var exceptionArr []metadata.ExceptionResult
	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		exceptionArr, err = lgc.AssignHostToApp(ctx.Kit.Ctx, conf)
		if err != nil {
			blog.Errorf("assign host to app, but assign to app http do error. err: %v, input:%+v,rid:%s", err, conf, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr,txnErr)
		return
	}
	ctx.RespEntity(nil)
}

func (s *Service) AssignHostToAppModule(ctx *rest.Contexts) {
	data := new(metadata.HostToAppModule)
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	appID, _, moduleID, err := lgc.GetTopoIDByName(ctx.Kit.Ctx, data)
	if nil != err {
		blog.Errorf("get app  topology id by name error:%s, msg: applicationName:%s, setName:%s, moduleName:%s,input;%+v,rid:%s", err.Error(), data.AppName, data.SetName, data.ModuleName, data, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrAddHostToModuleFailStr, "search application module not found "))
		return
	}

	if 0 == appID || 0 == moduleID {
		// get default app
		ownerAppID, err := lgc.GetDefaultAppID(ctx.Kit.Ctx)
		if err != nil {
			blog.Errorf("assign host to app module, but get resource pool failed, err: %v,input:%+v,rid:%s", err, data, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		if 0 == ownerAppID {
			blog.Errorf("assign host to app module, but get resource pool failed, err: %v,input:%+v,rid:%s", err, data, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrAddHostToModuleFailStr, "not found resource pool business"))
			return
		}

		// get idle module
		mConds := mapstr.New()
		mConds.Set(common.BKDefaultField, common.DefaultResModuleFlag)
		mConds.Set(common.BKModuleNameField, common.DefaultResModuleName)
		mConds.Set(common.BKAppIDField, ownerAppID)
		ownerModuleID, err := lgc.GetResourcePoolModuleID(ctx.Kit.Ctx, mConds)
		if nil != err {
			blog.Errorf("assign host to app module, but get unused host pool failed, ownerid[%v], err: %v,input:%+v,param:%+v,rid:%s", ownerModuleID, err, data, mConds, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrAddHostToModuleFailStr, err.Error()))
			return
		}
		appID = ownerAppID
		moduleID = ownerModuleID
		data.AppName = common.DefaultAppName
		data.SetName = ""
		data.ModuleName = common.DefaultResModuleName

	}

	// TODO host can not exist, not exist create
	// check authorization
	hostIDArr := make([]int64, 0)
	existNewAddHost := false
	for _, ip := range data.Ips {
		hostID, err := s.ip2hostID(ctx, ip, data.PlatID)
		if err != nil {
			blog.Errorf("invalid ip:%v, err:%v, rid:%s", ip, err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrAddHostToModuleFailStr, err.Error()))
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
			if err := s.AuthManager.AuthorizeAddToResourcePool(ctx.Kit.Ctx, ctx.Kit.Header); err != nil {
				blog.Errorf("check host authorization for add to resource pool failed, err: %v, rid: %s", err, ctx.Kit.Rid)
				resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed)})
				return
			}
		*/
		// 检查转移主机到目标业务的权限
		// auth: check target business update priority
		// if err := s.AuthManager.AuthorizeByBusinessID(ctx.Kit.Ctx, ctx.Kit.Header, authmeta.Update, appID); err != nil {
		// 	blog.Errorf("AssignHostToApp failed, authorize on business update failed, business: %d, err: %v, rid:%s", appID, err, ctx.Kit.Rid)
		// 	resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed)})
		// 	return
		// }
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
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
			err := lgc.EnterIP(ctx.Kit.Ctx, util.GetOwnerID(ctx.Request.Request.Header), appID, moduleID, ip, data.PlatID, host, data.IsIncrement)
			if nil != err {
				blog.Errorf("%s add host error: %s,input:%+v,rid:%s", ip, err.Error(), data, ctx.Kit.Rid)
				errmsg = append(errmsg, fmt.Sprintf("%s add host error: %s", ip, err.Error()))
			}
		}
		if 0 == len(errmsg) {
			return nil
		}

		blog.Errorf("assign host to app module failed, err: %v,rid:%s", errmsg, ctx.Kit.Rid)
		return ctx.Kit.CCError.Error(common.CCErrAddHostToModuleFailStr)
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// GetHostModuleRelation  query host and module relation,
// hostID can empty
func (s *Service) GetHostModuleRelation(ctx *rest.Contexts) {
	data := new(metadata.HostModuleRelationParameter)
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var cond metadata.HostModuleRelationRequest
	if data.AppID != 0 {
		cond.ApplicationID = data.AppID
	}
	pageSize := 500
	if len(data.HostID) > 0 {
		if len(data.HostID) > pageSize {
			blog.Errorf("GetHostModuleRelation host id length %d exceeds 500, rid: %s", len(data.HostID), ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommXXExceedLimit, common.BKHostIDField, pageSize))
		}
		cond.HostIDArr = data.HostID
	}
	if data.Page.Limit == 0 {
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "page.limit"))
	}
	if data.Page.Limit > pageSize {
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommPageLimitIsExceeded))
	}
	cond.Page = data.Page
	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	moduleHostConfig, err := lgc.GetHostModuleRelation(ctx.Kit.Ctx, cond)
	if err != nil {
		blog.Errorf("GetHostModuleRelation logcis err:%s,cond:%#v,rid:%s", err.Error(), cond, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(moduleHostConfig.Info)
	return
}

// TransferHostAcrossBusiness  Transfer host across business,
// delete old business  host and module relation
func (s *Service) TransferHostAcrossBusiness(ctx *rest.Contexts) {
	data := new(metadata.TransferHostAcrossBusinessParameter)
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		err := lgc.TransferHostAcrossBusiness(ctx.Kit.Ctx, data.SrcAppID, data.DstAppID, data.HostID, data.DstModuleIDArr)
		if err != nil {
			blog.Errorf("TransferHostAcrossBusiness logcis err:%s,input:%#v,rid:%s", err.Error(), data, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
	return
}

// DeleteHostFromBusiness delete host from business
// dangerous operation
func (s *Service) DeleteHostFromBusiness(ctx *rest.Contexts) {
	
	data := new(metadata.DeleteHostFromBizParameter)
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	var exceptionArr []metadata.ExceptionResult
	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {
		var err error
		exceptionArr, err = lgc.DeleteHostFromBusiness(ctx.Kit.Ctx, data.AppID, data.HostIDArr)
		if err != nil {
			blog.Errorf("DeleteHostFromBusiness logcis err:%s,input:%#v,rid:%s", err.Error(), data, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr,txnErr)
		return
	}
	ctx.RespEntity(nil)
	return
}

// move host to idle, fault or recycle module under the same business.
func (s *Service) moveHostToDefaultModule(ctx *rest.Contexts, defaultModuleFlag int) {

	defErr := ctx.Kit.CCError
	rid := ctx.Kit.Rid
	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := ctx.DecodeInto(&conf); nil != err {
		ctx.RespAutoError(err)
		return
	}

	bizID := conf.ApplicationID

	moduleFilter := make(map[string]interface{})
	var action authmeta.Action
	if defaultModuleFlag == common.DefaultResModuleFlag {
		// 空闲机
		action = authmeta.MoveHostToBizIdleModule
		moduleFilter[common.BKDefaultField] = common.DefaultResModuleFlag
		moduleFilter[common.BKModuleNameField] = common.DefaultResModuleName
	} else if defaultModuleFlag == common.DefaultFaultModuleFlag {
		// 故障机器
		action = authmeta.MoveHostToBizFaultModule
		moduleFilter[common.BKDefaultField] = common.DefaultFaultModuleFlag
		moduleFilter[common.BKModuleNameField] = common.DefaultFaultModuleName
	} else if defaultModuleFlag == common.DefaultRecycleModuleFlag {
		// 待回收
		action = authmeta.MoveHostToBizRecycleModule
		moduleFilter[common.BKDefaultField] = common.DefaultRecycleModuleFlag
		moduleFilter[common.BKModuleNameField] = common.DefaultRecycleModuleName
	} else {
		blog.Errorf("move host to default module failed, unexpected flag, bizID: %d, defaultModuleFlag: %d, rid: %s", bizID, defaultModuleFlag, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	moduleFilter[common.BKAppIDField] = bizID
	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	moduleID, err := lgc.GetResourcePoolModuleID(ctx.Kit.Ctx, moduleFilter)
	if err != nil {
		blog.ErrorJSON("move host to default module failed, get default module id failed, filter: %s, err: %s, rid: %s", moduleFilter, err, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.Errorf(common.CCErrAddHostToModuleFailStr, moduleFilter[common.BKModuleNameField].(string)+" not foud "))
		return
	}
	
	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(ctx.Kit.Ctx, ctx.Kit.Header, action, conf.HostIDs...); err != nil {
		blog.Errorf("auth host from iam failed, hosts: %+v, err: %v, rid: %s", conf.HostIDs, err, ctx.Kit.Rid)
		if err != ac.NoAuthorizeError {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(ctx.Kit.Ctx, ctx.Kit.Header, conf.HostIDs)
		if err != nil {
			ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommAuthorizeFailed))
			return
		}
		ctx.RespEntityWithError(perm, auth.NoAuthorizeError)
		return
	}
	
	audit := lgc.NewHostModuleLog(conf.HostIDs)
	if err := audit.WithPrevious(ctx.Kit.Ctx); err != nil {
		blog.Errorf("move host to default module s failed, get prev module host config failed, hostIDs: %v, err: %s, rid: %s", conf.HostIDs, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}
	
	var result *metadata.OperaterException
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, s.EnableTxn, ctx.Kit.Header, func() error {

		transferInput := &metadata.TransferHostToInnerModule{
			ApplicationID: conf.ApplicationID,
			HostID:        conf.HostIDs,
			ModuleID:      moduleID,
		}
		var err error
		result, err = s.CoreAPI.CoreService().Host().TransferToInnerModule(ctx.Kit.Ctx, ctx.Kit.Header, transferInput)
		if err != nil {
			blog.ErrorJSON("move host to default module failed, TransferHostToDefaultModule http do error. input:%s, condition:%s, err:%s, rid:%s", conf, transferInput, err.Error(), rid)
			return defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		if !result.Result {
			blog.ErrorJSON("move host to default module failed, TransferHostToDefaultModule response failed. input:%s, transferInput:%s, response:%s, rid:%s", conf, transferInput, result, rid)
			return defErr.New(result.Code, result.ErrMsg)
		}
		
		if err := audit.SaveAudit(ctx.Kit.Ctx); err != nil {
			blog.ErrorJSON("move host to default module failed, save audit log failed, input:%s, err:%s, rid:%s", conf, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server")
		}
		return nil
		})

	if txnErr != nil {
		ctx.RespEntityWithError(result.Data,txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// GetAppHostTopoRelation  query host and module relation,
// hostID can empty
func (s *Service) GetAppHostTopoRelation(ctx *rest.Contexts) {
	data := new(metadata.HostModuleRelationRequest)
	if err := ctx.DecodeInto(&data); nil != err {
		ctx.RespAutoError(err)
		return
	}

	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	result, err := lgc.GetHostModuleRelation(ctx.Kit.Ctx, *data)
	if err != nil {
		blog.Errorf("GetHostModuleRelation logic failed, cond:%#v, err:%s, rid:%s", data, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(result)
	return
}

func (s *Service) TransferHostResourceDirectory(ctx *rest.Contexts) {
	input := new(metadata.TransferHostResourceDirectory)
	if err := ctx.DecodeInto(&input); nil != err {
		ctx.RespAutoError(err)
		return
	}

	lgc := logics.NewLogics(s.Engine, ctx.Kit.Header, s.CacheDB, s.AuthManager)
	audit := lgc.NewHostModuleLog(input.HostID)
	if err := audit.WithPrevious(ctx.Kit.Ctx); err != nil {
		blog.Errorf("TransferHostResourceDirectory, but get prev module host config failed, err: %v, hostIDs:%#v,rid:%s", err, input.HostID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	err := s.CoreAPI.CoreService().Host().TransferHostResourceDirectory(ctx.Kit.Ctx, ctx.Kit.Header, input)
	if err != nil {
		blog.Errorf("TransferHostResourceDirectory failed with coreservice http failed, input: %v, err: %v, rid: %s", input, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if err := audit.SaveAudit(ctx.Kit.Ctx); err != nil {
		blog.Errorf("move host to resource pool, but save audit log failed, err: %v, input:%+v,rid:%s", err, input.HostID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	ctx.RespEntity(nil)
	return
}
