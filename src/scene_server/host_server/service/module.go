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
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"

    "github.com/emicklei/go-restful"

    "configcenter/src/common"
    "configcenter/src/common/blog"
    "configcenter/src/common/mapstr"
    "configcenter/src/common/metadata"
    "configcenter/src/common/util"
    hutil "configcenter/src/scene_server/host_server/util"
)

func (s *Service) AddHostMultiAppModuleRelation(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	result, err := s.CoreAPI.ObjectController().Privilege().GetSystemFlag(srvData.ctx, common.BKDefaultOwnerID, common.HostCrossBizField, srvData.header)
	if err != nil {
		blog.Errorf("AddHostMultiAppModuleRelation GetSystemFlag http do error,err:%s,rid:%s", err.Error(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("AddHostMultiAppModuleRelation GetSystemFlag http response error,err code:%d,err msg:%s,rid:%s", result.Code, result.ErrMsg, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	params := new(metadata.CloudHostModuleParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("add host multiple app module relation failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	module, err := srvData.lgc.GetNormalModuleByModuleID(srvData.ctx, params.ApplicationID, params.ModuleID)
	if err != nil {
		blog.Errorf("add host multiple app module relation, but get module[%v] failed, err: %v,input:%+v,rid:%s", params.ModuleID, err, params, srvData.ctx)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoModuleSelectFailed)})
		return
	}

	if len(module) == 0 {
		blog.Errorf("add host multiple app module relation, but get invalid module.input:%+v,rid:%s", params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
		return
	}

	defaultAppID, err := srvData.lgc.GetDefaultAppID(srvData.ctx, srvData.ownerID)
	if err != nil {
		blog.Errorf("add host multiple app module relation, but get default appid failed, err: %v,param:%+v,rid:%s", err, params, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	var errMsg, succ []string
	var hostIDArr []int64

	for index, hostInfo := range params.HostInfoArr {
		cond := hutil.NewOperation().WithHostInnerIP(hostInfo.IP).WithCloudID(int64(hostInfo.CloudID)).Data()
		query := &metadata.QueryInput{
			Condition: cond,
			Start:     0,
			Limit:     common.BKNoLimit,
			Sort:      common.BKHostIDField,
		}
		hResult, err := s.CoreAPI.HostController().Host().GetHosts(srvData.ctx, srvData.header, query)
		if err != nil || (err == nil && !hResult.Result) {
			blog.Errorf("add host multiple app module relation, but get hosts failed, err: %v, %v,param:%+v,rid:%s", err, hResult.ErrMsg, params, srvData.rid)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		hostList := hResult.Data.Info
		if len(hostList) == 0 {
			blog.Errorf("add host multiple app module relation, but get 0 hosts.params:%+v,rid:%s", params, srvData.rid)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		//check if host in this module
		hostID, err := util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if nil != err {
			blog.Errorf("add host multiple app module relation, but get invalid host id[%v], err:%v.params:%+v,rid:%s", hostList[0][common.BKHostIDField], err.Error(), params, srvData.rid)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}
		moduleHostCond := map[string][]int64{common.BKHostIDField: []int64{hostID}}
		confs, err := srvData.lgc.GetConfigByCond(srvData.ctx, moduleHostCond)
		if err != nil {
			blog.Errorf("add host multiple app module relation, but get host config failed, err:%v.param:%+v,rid:%s", err, params, srvData.rid)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		for _, conf := range confs {
			if conf[common.BKAppIDField] == defaultAppID {
				p := metadata.ModuleHostConfigParams{
					ApplicationID: defaultAppID,
					HostID:        hostID,
				}
				hResult, err := s.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(srvData.ctx, srvData.header, &p)
				if err != nil || (err == nil && !hResult.Result) {
					blog.Errorf("add host multiple app module relation, but delete default module host conf failed, err: %v, %v.params:%+v,rid:%s", err, hResult.ErrMsg, params, srvData.rid)
					errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
					continue
				}
			}

			if conf[common.BKModuleIDField] == params.ModuleID {
				blog.Errorf("add host multiple app module relation, but host already exist in module.params:%+v,rid:%s", params, srvData.rid)
				errMsg = append(errMsg, s.Language.Languagef("host_str_belong_module", hostInfo.IP))
				continue
			}
		}

		//add host to this module
		opt := metadata.ModuleHostConfigParams{
			ApplicationID: params.ApplicationID,
			HostID:        hostID,
			ModuleID:      []int64{params.ModuleID},
		}
		result, err := s.CoreAPI.HostController().Module().AddModuleHostConfig(srvData.ctx, srvData.header, &opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("add host multiple app module relation, but add module host config failed, err: %v, %v.params:%+v,rid:%s", err, result.ErrMsg, params, srvData.rid)
			errMsg = append(errMsg, s.Language.Languagef("host_str_add_module_relation_fail", hostInfo.IP))
			continue
		}

		hostIDArr = append(hostIDArr, hostID)
		succ = append(succ, strconv.Itoa(index))
	}

	if 0 != len(errMsg) {
		detail := make(map[string]interface{})
		detail["success"] = succ
		detail["error"] = errMsg
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrAddHostToModule), Data: detail})
	}

	// TODO: add audit log later.
	hostModuleLog := srvData.lgc.NewHostModuleLog(hostIDArr)
	hostModuleLog.WithCurrent(srvData.ctx)
	hostModuleLog.SaveAudit(srvData.ctx, fmt.Sprintf("%d", params.ApplicationID), util.GetUser(req.Request.Header), "")
	resp.WriteEntity(metadata.NewSuccessResp(nil))

}

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

	for _, hostID := range config.HostID {
		exist, err := srvData.lgc.IsHostExistInApp(srvData.ctx, config.ApplicationID, hostID)
		if err != nil {
			blog.Errorf("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s,input:%+v,rid:%s", config.ApplicationID, hostID, err.Error(), config, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotINAPPFail)})
			return
		}

		if !exist {
			blog.Errorf("Host does not belong to the current application, appid: %v, hostid: %v,input:%+v,rid:%s", config.ApplicationID, hostID, config, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostNotINAPP, hostID)})
			return
		}

		opt := metadata.ModuleHostConfigParams{
			ApplicationID: config.ApplicationID,
			HostID:        hostID,
		}

		var result *metadata.BaseResp
		if config.IsIncrement {
			result, err = s.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(srvData.ctx, srvData.header, &opt)
		} else {
			result, err = s.CoreAPI.HostController().Module().DelModuleHostConfig(srvData.ctx, srvData.header, &opt)
		}
		if err != nil {
			blog.Errorf("update host module relation, but delete default config failed, err: %v, %v,input:%+v,param:%+v,rid:%s", err, result.ErrMsg, config, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !result.Result {
			blog.Errorf("update host module relation, but delete default config failed, err: %v, %v.input:%+v,param:%+v,rid:%s", err, result.ErrMsg, config, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}

		opt.ModuleID = config.ModuleID
		result, err = s.CoreAPI.HostController().Module().AddModuleHostConfig(srvData.ctx, srvData.header, &opt)
		if err != nil {
			blog.Errorf("add host module relation, but add config failed, err: %v, %v,input:%+v,param:%+v,rid:%s", err, result.ErrMsg, config, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !result.Result {
			blog.Errorf("add host module relation, but add config failed, err: %v, %v.input:%+v,param:%+v,rid:%s", err, result.ErrMsg, config, opt, srvData.rid)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}
	}

	if err := audit.SaveAudit(srvData.ctx, strconv.FormatInt(config.ApplicationID, 10), srvData.user, ""); err != nil {
		blog.Errorf("host module relation, save audit log failed, err: %v,input:%+v,rid:%s", err, config, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
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

	cond := hutil.NewOperation().WithAppID(conf.ApplicationID).Data()
	appInfo, err := srvData.lgc.GetAppDetails(srvData.ctx, common.BKOwnerIDField, cond)
	if err != nil {
		blog.Errorf("move host to resource pool, but get app detail failed, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostMoveResourcePoolFail, fmt.Sprintf("%v", conf.HostID))})
		return
	}
	if 0 == len(appInfo) {
		blog.Errorf("assign host to app error, not found app appID: %d, input:%#v,rid:%s", conf.ApplicationID, conf, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommNotFound)})
		return
	}

	ownerID, err := appInfo.String(common.BKOwnerIDField)
	if nil != err {
		blog.Errorf("move host to resource pool , but get app detail failed, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetResourceFail, "app info OwnerID not string")})
		return
	}

	if "" == ownerID {
		blog.Errorf("move host to resource pool, but get app detail failed, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetResourceFail, "app info OwnerID empty")})
		return
	}

	ownerAppID, err := srvData.lgc.GetDefaultAppID(srvData.ctx, srvData.ownerID)
	if err != nil {
		blog.Errorf("move host to resource pool, but get default appid failed, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	if 0 == conf.ApplicationID {
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostNotResourceFail)})
		return
	}
	if ownerAppID == conf.ApplicationID {
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostBelongResourceFail)})
		return
	}

	conds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(ownerAppID)
	moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, conds.MapStr())
	if err != nil {
		blog.Errorf("move host to resource pool, but get module id failed, err: %v, input:%+v,param:%+v,rid:%s", err, conf, conds.Data(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}

	param := &metadata.ParamData{
		ApplicationID:       conf.ApplicationID,
		HostID:              conf.HostID,
		OwnerModuleID:       moduleID,
		OwnerAppplicationID: ownerAppID,
	}

	audit := srvData.lgc.NewHostModuleLog(conf.HostID)
	if err := audit.WithPrevious(srvData.ctx); err != nil {
		blog.Errorf("move host to resource pool, but get prev module host config failed, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}
	result, err := s.CoreAPI.HostController().Module().MoveHost2ResourcePool(srvData.ctx, srvData.header, param)
	if err != nil {
		blog.Errorf("move host to resource pool, but update host module http do error, err: %v, input:%+v,query:%+v,rid:%v", err, conf, param, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("move host to resource pool, but update host module http response error, err code:%d, err messge:%s, input:%+v,query:%+v,rid:%v", result.Code, result.ErrMsg, conf, param, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	if err := audit.SaveAudit(srvData.ctx, strconv.FormatInt(conf.ApplicationID, 10), srvData.user, "move host to resource pool"); err != nil {
		blog.Errorf("move host to resource pool, but save audit log failed, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}
	businessMetadata := conf.Metadata
	if businessMetadata.Label == nil {
		businessMetadata.Label = make(metadata.Label)
	}
	businessMetadata.Label.SetBusinessID(conf.ApplicationID)
	if err := srvData.lgc.DeleteHostBusinessAttributes(srvData.ctx, conf.HostID, &businessMetadata); err != nil {
		blog.Errorf("move host to resource pool, delete host bussiness private, err: %v, input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("assign host to app failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	cond := hutil.NewOperation().WithAppID(conf.ApplicationID).Data()
	fields := fmt.Sprintf("%s,%s", common.BKOwnerIDField, common.BKAppNameField)
	appInfo, err := srvData.lgc.GetAppDetails(srvData.ctx, fields, cond)
	if err != nil {
		blog.Errorf("assign host to app failed, err: %v,input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	if 0 == len(appInfo) {
		blog.Errorf("assign host to app error, not foud app appID: %d,input:%+v,rid:%s", conf.ApplicationID, conf, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommNotFound)})
		return
	}

	ownerID, err := appInfo.String(common.BKOwnerIDField)
	if nil != err {
		blog.Errorf("assign host to app, but get app detail failed, err: %v,input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "OwnerID")})
		return
	}

	appID, err := srvData.lgc.GetDefaultAppID(srvData.ctx, ownerID)
	if err != nil {
		blog.Errorf("assign host to app, but get default appid failed, err: %v,input:%+v,rid:%s", err, conf, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	if 0 == conf.ApplicationID {
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostGetResourceFail, "not found")})
		return
	}
	if appID == conf.ApplicationID {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	conds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(appID)
	ownerModuleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, conds.MapStr())
	if err != nil {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,rid:%s", err, conds.MapStr(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	if 0 == ownerModuleID {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,rid:%s", err, conds.MapStr(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleNotExist, common.DefaultResModuleName)})
		return
	}

	mConds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(conf.ApplicationID)
	moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, mConds.MapStr())
	if err != nil {
		blog.Errorf("assign host to app, but get module id failed, err: %v,input:%+v,params:%+v,rid:%s", err, conf, mConds.MapStr(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
		return
	}
	if moduleID == 0 {
		blog.Errorf("assign host to app, but get module id failed, %s not found: %v,input:%+v,params:%+v,rid:%s", common.DefaultResModuleName, conf, mConds.MapStr(), srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrHostModuleNotExist, common.DefaultResModuleName)})
		return
	}

	params := make(map[string]interface{})
	params[common.BKAppIDField] = conf.ApplicationID
	params[common.BKHostIDField] = conf.HostID
	params[common.BKModuleIDField] = moduleID
	params["bk_owner_module_id"] = ownerModuleID
	params["bk_owner_biz_id"] = appID

	audit := srvData.lgc.NewHostModuleLog(conf.HostID)
	audit.WithPrevious(srvData.ctx)

	result, err := s.CoreAPI.HostController().Module().AssignHostToApp(srvData.ctx, srvData.header, params)
	if err != nil {
		blog.Errorf("assign host to app, but assign to app http do error. err: %v, input:%+v,param:%+v,rid:%s", err, conf, params)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostEditRelationPoolFail)})
		return
	}
	if !result.Result {
		blog.Errorf("assign host to app, but assign to app http response error. err code:%d, err msg:%s,input:%+v,param:%+v,rid:%s", result.Code, result.ErrMsg, conf, params)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	if err := audit.SaveAudit(srvData.ctx, strconv.FormatInt(conf.ApplicationID, 10), srvData.user, "assign host to app"); err != nil {
		blog.Errorf("assign host to app, but save audit failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) AssignHostToAppModule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	data := new(metadata.HostToAppModule)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("assign host to app module failed with decode body err: %v", err)
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
		ownerAppID, err := srvData.lgc.GetDefaultAppID(srvData.ctx, data.OwnerID)
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

		//dispatch to app
		err := srvData.lgc.EnterIP(srvData.ctx, util.GetOwnerID(req.Request.Header), appID, moduleID, ip, data.PlatID, host, data.IsIncrement)
		if nil != err {
			blog.Errorf("%s add host error: %s,input:%+v,rid:%s", ip, err.Error(), data, srvData.rid)
			errmsg = append(errmsg, fmt.Sprintf("%s add host error: %s", ip, err.Error()))
		}
	}
	if 0 == len(errmsg) {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	} else {
		blog.Errorf("assign host to app module failed, err: %v,rid:%s", errmsg, srvData)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrAddHostToModuleFailStr)})
		return
	}
}

// GetHostModuleRelation  query host and module relation,
// hostID can emtpy
func (s *Service) GetHostModuleRelation(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	data := new(metadata.HostModuleRelationParameter)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("Transfer host across business failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	cond := make(map[string][]int64, 0)
	if data.AppID != 0 {
		cond[common.BKAppIDField] = []int64{data.AppID}
	}
	if len(data.HostID) > 0 {
		cond[common.BKHostIDField] = data.HostID
	}
	if len(cond) == 0 {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
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
// delete old business  host and module reltaion
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

func (s *Service) moveHostToModuleByName(req *restful.Request, resp *restful.Response, moduleName string) {
    pheader := req.Request.Header
    srvData := s.newSrvComm(pheader )
    defErr := srvData.ccErr
    ctx := srvData.ctx
    rid :=srvData.rid
    conf := new(metadata.DefaultModuleHostConfigParams)
    if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
        blog.Errorf("move host to module %s failed with decode body err: %v,rid: %s", moduleName, err,rid)
        resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
        return
    }

    conds := make(map[string]interface{})
    var moduleNameLogKey string
    if common.DefaultResModuleName == moduleName {
        //空闲机
        moduleNameLogKey = "idle"
        conds[common.BKDefaultField] = common.DefaultResModuleFlag
        conds[common.BKModuleNameField] = common.DefaultResModuleName
    } else {
        //故障机器
        moduleNameLogKey = "fault"
        conds[common.BKDefaultField] = common.DefaultFaultModuleFlag
        conds[common.BKModuleNameField] = common.DefaultFaultModuleName
    }
    conds[common.BKAppIDField] = conf.ApplicationID
    moduleID, err := srvData.lgc.GetResoulePoolModuleID(srvData.ctx, conds)
    if err != nil {
        blog.Errorf("move host to module %s, get module id err: %v", moduleName, err)
        resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrAddHostToModuleFailStr, conds[common.BKModuleNameField].(string)+" not foud ")})
        return
    }

    audit := srvData.lgc.NewHostModuleLog( conf.HostID)
    if err := audit.WithPrevious(srvData.ctx); err != nil {
        blog.Errorf("move host to module %s, get prev module host config failed, err: %v", moduleName, err)
        resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
        return
    }

    notExistHostID, err := srvData.lgc.ExistHostIDSInApp(ctx, conf.ApplicationID, conf.HostID)
    if err != nil {
        blog.Errorf("moveHostToModuleByName ExistHostIDSInApp error, err:%s,input:%#v,rid:%s", err.Error(), conf, rid)
        resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
        return
    }
    if len(notExistHostID) > 0 {
        blog.Errorf("Host does not belong to the current application, appid: %v, hostid: %#v, not exist in app:%#v,rid:%s", conf.ApplicationID, conf.HostID, notExistHostID, rid)
        notTipStrHostID := ""
        for _, hostID := range notExistHostID {
            notTipStrHostID = fmt.Sprintf("%s,%s", notTipStrHostID, hostID)
        }
        resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrHostNotINAPP, strings.Trim(notTipStrHostID, ","))})
        return
    }
    transferInput := &metadata.TransferHostToDefaultModuleConfig{
        ApplicationID: conf.ApplicationID,
        HostID:        conf.HostID,
        ModuleID:      moduleID,
    }
    result, err := s.CoreAPI.HostController().Module().TransferHostToDefaultModule(ctx, pheader, transferInput)
    if err != nil {
        blog.Errorf("moveHostToModuleByName TransferHostToDefaultModule http do error. input:%#v,condition:%#v,err:%v,rid:%s", conf, transferInput, err.Error(), rid)
        resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
        return
    }
    if !result.Result {
        blog.Errorf("moveHostToModuleByName TransferHostToDefaultModule http reply error. input:%#v,condition:%#v,err:%#v,rid:%s", conf, transferInput, result, rid)
        resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPDoRequestFailed)})
        return
    }

    user := util.GetUser(pheader)
    if err := audit.SaveAudit(srvData.ctx, strconv.FormatInt(conf.ApplicationID, 10), user, "host to "+moduleNameLogKey+" module"); err != nil {
        blog.Errorf("move host to module %s, save audit log failed, err: %v, rid: %s", moduleName, err, rid)
        resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
        return
    }
    resp.WriteEntity(metadata.NewSuccessResp(nil))
}
