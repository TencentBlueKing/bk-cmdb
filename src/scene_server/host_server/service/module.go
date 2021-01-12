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
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

// HostModuleRelation transfer host to module specify by bk_module_id (in the same business)
// move a business host to a module.
func (s *Service) TransferHostModule(ctx *rest.Contexts) {
	config := new(metadata.HostsModuleRelation)
	if err := ctx.DecodeInto(&config); nil != err {
		ctx.RespAutoError(err)
		return
	}

	for _, moduleID := range config.ModuleID {
		module, err := s.Logic.GetNormalModuleByModuleID(ctx.Kit, config.ApplicationID, moduleID)
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

	audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), config.HostID)
	if err := audit.WithPrevious(ctx.Kit); err != nil {
		blog.Errorf("host module relation, get prev module host config failed, err: %v,param:%+v,rid:%s", err, config, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	var result *metadata.OperaterException
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		result, err = s.CoreAPI.CoreService().Host().TransferToNormalModule(ctx.Kit.Ctx, ctx.Kit.Header, config)
		if err != nil {
			blog.Errorf("add host module relation, but add config failed, err: %v, %v,input:%+v,rid:%s", err, result.ErrMsg, config, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !result.Result {
			blog.Errorf("add host module relation, but add config failed, err: %v, %v.input:%+v,rid:%s", err, result.ErrMsg, config, ctx.Kit.Rid)
			return result.CCError()
		}

		if err := audit.SaveAudit(ctx.Kit); err != nil {
			blog.Errorf("host module relation, save audit log failed, err: %v,input:%+v,rid:%s", err, config, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(result.Data, txnErr)
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

	var exceptionArr []metadata.ExceptionResult
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		exceptionArr, err = s.Logic.MoveHostToResourcePool(ctx.Kit, conf)
		if err != nil {
			blog.Errorf("move host to resource pool failed, err:%s, input:%#v, rid:%s", err.Error(), conf, ctx.Kit.Rid)
			return ctx.Kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr, txnErr)
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
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		exceptionArr, err = s.Logic.AssignHostToApp(ctx.Kit, conf)
		if err != nil {
			blog.Errorf("assign host to app, but assign to app http do error. err: %v, input:%+v,rid:%s", err, conf, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr, txnErr)
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

	rawErr := data.Validate()
	if rawErr.ErrCode != 0 {
		ctx.RespAutoError(rawErr.ToCCError(ctx.Kit.CCError))
		return
	}

	cond := metadata.HostModuleRelationRequest{
		HostIDArr: data.HostID,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	if data.AppID != 0 {
		cond.ApplicationID = data.AppID
	}

	moduleHostConfig, err := s.Logic.GetHostModuleRelation(ctx.Kit, cond)

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

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logic.TransferHostAcrossBusiness(ctx.Kit, data.SrcAppID, data.DstAppID, data.HostID, data.DstModuleID)
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
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		exceptionArr, err = s.Logic.DeleteHostFromBusiness(ctx.Kit, data.AppID, data.HostIDArr)
		if err != nil {
			blog.Errorf("DeleteHostFromBusiness logcis err:%s,input:%#v,rid:%s", err.Error(), data, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(exceptionArr, txnErr)
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
	if defaultModuleFlag == common.DefaultResModuleFlag {
		// 空闲机
		moduleFilter[common.BKDefaultField] = common.DefaultResModuleFlag
		moduleFilter[common.BKModuleNameField] = common.DefaultResModuleName
	} else if defaultModuleFlag == common.DefaultFaultModuleFlag {
		// 故障机器
		moduleFilter[common.BKDefaultField] = common.DefaultFaultModuleFlag
		moduleFilter[common.BKModuleNameField] = common.DefaultFaultModuleName
	} else if defaultModuleFlag == common.DefaultRecycleModuleFlag {
		// 待回收
		moduleFilter[common.BKDefaultField] = common.DefaultRecycleModuleFlag
		moduleFilter[common.BKModuleNameField] = common.DefaultRecycleModuleName
	} else {
		blog.Errorf("move host to default module failed, unexpected flag, bizID: %d, defaultModuleFlag: %d, rid: %s", bizID, defaultModuleFlag, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	moduleFilter[common.BKAppIDField] = bizID
	moduleID, err := s.Logic.GetResourcePoolModuleID(ctx.Kit, moduleFilter)
	if err != nil {
		blog.ErrorJSON("move host to default module failed, get default module id failed, filter: %s, err: %s, rid: %s", moduleFilter, err, ctx.Kit.Rid)
		ctx.RespAutoError(defErr.Errorf(common.CCErrAddHostToModuleFailStr, moduleFilter[common.BKModuleNameField].(string)+" not foud "))
		return
	}

	audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), conf.HostIDs)
	if err := audit.WithPrevious(ctx.Kit); err != nil {
		blog.Errorf("move host to default module s failed, get prev module host config failed, hostIDs: %v, err: %s, rid: %s", conf.HostIDs, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	var result *metadata.OperaterException
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {

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

		if err := audit.SaveAudit(ctx.Kit); err != nil {
			blog.ErrorJSON("move host to default module failed, save audit log failed, input:%s, err:%s, rid:%s", conf, err, ctx.Kit.Rid)
			return ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server")
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespEntityWithError(result.Data, txnErr)
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

	result, err := s.Logic.GetHostModuleRelation(ctx.Kit, *data)
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

	audit := auditlog.NewHostModuleLog(s.CoreAPI.CoreService(), input.HostID)
	if err := audit.WithPrevious(ctx.Kit); err != nil {
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

	if err := audit.SaveAudit(ctx.Kit); err != nil {
		blog.Errorf("move host to resource pool, but save audit log failed, err: %v, input:%+v,rid:%s", err, input.HostID, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Errorf(common.CCErrCommResourceInitFailed, "audit server"))
		return
	}

	ctx.RespEntity(nil)
	return
}
