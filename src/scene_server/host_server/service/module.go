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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	hutil "configcenter/src/scene_server/host_server/util"
	"github.com/emicklei/go-restful"
)

func (s *Service) AddHostMultiAppModuleRelation(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	// user := util.GetUser(pheader)

	result, err := s.CoreAPI.ObjectController().Privilege().GetSystemFlag(context.Background(), common.BKDefaultOwnerID, common.HostCrossBizField, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add host multiple app module relation failed, err: %v, result err: %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostNotAllowedToMutiBiz)})
		return
	}

	params := new(metadata.CloudHostModuleParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("add host multiple app module relation failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	module, err := s.Logics.GetModuleByModuleID(pheader, params.ApplicationID, params.ModuleID)
	if err != nil {
		blog.Errorf("add host multiple app module relation, but get module[%v] failed, err: %v", params.ModuleID, err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoModuleSelectFailed)})
		return
	}

	if len(module) == 0 {
		blog.Errorf("add host multiple app module relation, but get invalid module")
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
		return
	}

	defaultAppID, err := s.Logics.GetDefaultAppID(common.BKDefaultOwnerID, pheader)
	if err != nil {
		blog.Errorf("add host multiple app module relation, but get default appid failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoAppSearchFailed)})
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
		hResult, err := s.CoreAPI.HostController().Host().GetHosts(context.Background(), pheader, query)
		if err != nil || (err == nil && !hResult.Result) {
			blog.Errorf("add host multiple app module relation, but get hosts failed, err: %v, %v", err, hResult.ErrMsg)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		hostList := hResult.Data.Info
		if len(hostList) == 0 {
			blog.Errorf("add host multiple app module relation, but get 0 hosts ")
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		//check if host in this module
		hostID, err := util.GetInt64ByInterface(hostList[0][common.BKHostIDField])
		if nil != err {
			blog.Error("add host multiple app module relation, but get invalid host id[%v], err:%v", hostList[0][common.BKHostIDField], err.Error())
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}
		moduleHostCond := map[string][]int64{common.BKHostIDField: []int64{hostID}}
		confs, err := s.Logics.GetConfigByCond(pheader, moduleHostCond)
		if err != nil {
			blog.Error("add host multiple app module relation, but get host config failed, err:%v", err)
			errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
			continue
		}

		for _, conf := range confs {
			if conf[common.BKAppIDField] == defaultAppID {
				p := metadata.ModuleHostConfigParams{
					ApplicationID: defaultAppID,
					HostID:        hostID,
				}
				hResult, err := s.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(context.Background(), pheader, &p)
				if err != nil || (err == nil && !hResult.Result) {
					blog.Errorf("add host multiple app module relation, but delete default module host conf failed, err: %v, %v", err, hResult.ErrMsg)
					errMsg = append(errMsg, s.Language.Languagef("host_ip_not_exist", hostInfo.IP))
					continue
				}
			}

			if conf[common.BKModuleIDField] == params.ModuleID {
				blog.Errorf("add host multiple app module relation, but host already exist in module")
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
		result, err := s.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), pheader, &opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("add host multiple app module relation, but add module host config failed, err: %v, %v", err, result.ErrMsg)
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
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAddHostToModule), Data: detail})
	}

	// TODO: add audit log later.
	hostModuleLog := s.Logics.NewHostModuleLog(req.Request.Header, hostIDArr)
	hostModuleLog.WithCurrent()
	hostModuleLog.SaveAudit(fmt.Sprintf("%d", params.ApplicationID), util.GetUser(req.Request.Header), "")
	resp.WriteEntity(metadata.NewSuccessResp(nil))

}

func (s *Service) HostModuleRelation(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	config := new(metadata.HostsModuleRelation)
	if err := json.NewDecoder(req.Request.Body).Decode(config); err != nil {
		blog.Errorf("add host and module relation failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	for _, moduleID := range config.ModuleID {
		module, err := s.Logics.GetModuleByModuleID(pheader, config.ApplicationID, moduleID)
		if err != nil {
			blog.Errorf("add host and module relation, but get module with id[%d] failed, err: %v", moduleID, err)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoModuleSelectFailed)})
			return
		}

		if len(module) == 0 {
			blog.Errorf("add host and module relation, but get empty module with id[%d] ", moduleID)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
			return
		}
	}

	audit := s.Logics.NewHostModuleLog(pheader, config.HostID)
	if err := audit.WithPrevious(); err != nil {
		blog.Errorf("host module relation, get prev module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	for _, hostID := range config.HostID {
		exist, err := s.Logics.IsHostExistInApp(config.ApplicationID, hostID, pheader)
		if err != nil {
			blog.Error("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", config.ApplicationID, hostID, err.Error())
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrHostNotINAPPFail)})
			return
		}

		if !exist {
			blog.Errorf("Host does not belong to the current application, appid: %v, hostid: %v", config.ApplicationID, hostID)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrHostNotINAPP, hostID)})
			return
		}

		opt := metadata.ModuleHostConfigParams{
			ApplicationID: config.ApplicationID,
			HostID:        hostID,
		}

		var result *metadata.BaseResp
		if config.IsIncrement {
			result, err = s.CoreAPI.HostController().Module().DelDefaultModuleHostConfig(context.Background(), pheader, &opt)
		} else {
			result, err = s.CoreAPI.HostController().Module().DelModuleHostConfig(context.Background(), pheader, &opt)
		}
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("update host module relation, but delete default config failed, err: %v, %v", err, result.ErrMsg)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrHostDELResourcePool)})
			return
		}

		opt.ModuleID = config.ModuleID
		result, err = s.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), pheader, &opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("update host module relation, but add config failed, err: %v, %v", err, result.ErrMsg)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrHostAddRelationFail)})
			return
		}
	}

	user := util.GetUser(pheader)
	if err := audit.SaveAudit(strconv.FormatInt(config.ApplicationID, 10), user, ""); err != nil {
		blog.Errorf("host module relation, save audit log failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
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
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	defLang := s.Language.CreateDefaultCCLanguageIf(util.GetLanguage(pheader))
	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("move host to resource pool failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	cond := hutil.NewOperation().WithAppID(conf.ApplicationID).Data()
	appInfo, err := s.Logics.GetAppDetails(common.BKOwnerIDField, cond, pheader)
	if err != nil {
		blog.Errorf("move host to resource pool, but get app detail failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrHostMoveResourcePoolFail, fmt.Sprintf("%v", conf.HostID))})
		return
	}

	ownerID := appInfo[common.BKOwnerIDField].(string)
	if "" == ownerID {
		blog.Errorf("move host to resource pool, but get app detail failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: errors.New(defLang.Language("host_resource_pool_not_exist"))})
		return
	}

	ownerAppID, err := s.Logics.GetDefaultAppID(ownerID, pheader)
	if err != nil {
		blog.Errorf("move host to resource pool, but get default appid failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: errors.New(defLang.Language("host_resource_pool_get_fail"))})
		return
	}
	if 0 == conf.ApplicationID {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: errors.New(defLang.Language("host_resource_pool_not_exist"))})
		return
	}
	if ownerAppID == conf.ApplicationID {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: errors.New(defLang.Language("host_belong_resource_pool"))})
		return
	}

	conds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(ownerAppID)
	moduleID, err := s.Logics.GetResoulePoolModuleID(pheader, conds.Data())
	if err != nil {
		blog.Errorf("move host to resource pool, but get module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: errors.New(defLang.Languagef("host_resource_module_get_fail", err.Error()))})
		return
	}

	param := &metadata.ParamData{
		ApplicationID:       conf.ApplicationID,
		HostID:              conf.HostID,
		OwnerModuleID:       moduleID,
		OwnerAppplicationID: ownerAppID,
	}

	audit := s.Logics.NewHostModuleLog(pheader, conf.HostID)
	if err := audit.WithPrevious(); err != nil {
		blog.Errorf("move host to resource pool, but get prev module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}
	result, err := s.CoreAPI.HostController().Module().MoveHost2ResourcePool(context.Background(), pheader, param)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("move host to resource pool, but update host module failed, err: %v, %v", err, result.ErrMsg)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrHostEditRelationPoolFail)})
		return
	}

	user := util.GetUser(pheader)
	if err := audit.SaveAudit(strconv.FormatInt(conf.ApplicationID, 10), user, "move host to resource pool"); err != nil {
		blog.Errorf("move host to resource pool, but save audit log failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("assign host to app failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	cond := hutil.NewOperation().WithAppID(conf.ApplicationID).Data()
	fields := fmt.Sprintf("%s,%s", common.BKOwnerIDField, common.BKAppNameField)
	appInfo, err := s.Logics.GetAppDetails(fields, cond, pheader)
	if err != nil {
		blog.Errorf("assign host to app failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CC_Err_Comm_APP_QUERY_FAIL)})
		return
	}

	ownerID := appInfo[common.BKOwnerIDField].(string)
	if "" == ownerID {
		blog.Errorf("assign host to app, but get app detail failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "OwnerID")})
		return
	}

	appID, err := s.Logics.GetDefaultAppID(ownerID, pheader)
	if err != nil {
		blog.Errorf("assign host to app, but get default appid failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoAppSearchFailed)})
		return
	}
	if 0 == conf.ApplicationID {
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommNotFound)})
		return
	}
	if appID == conf.ApplicationID {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	}

	conds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(appID)
	ownerModuleID, err := s.Logics.GetResoulePoolModuleID(pheader, conds.Data())
	if err != nil || (err == nil && 0 == ownerModuleID) {
		blog.Errorf("assign host to app, but get module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
		return
	}

	mConds := hutil.NewOperation().WithDefaultField(int64(common.DefaultResModuleFlag)).WithModuleName(common.DefaultResModuleName).WithAppID(conf.ApplicationID)
	moduleID, err := s.Logics.GetResoulePoolModuleID(pheader, mConds.Data())
	if err != nil || (err == nil && 0 == moduleID) {
		blog.Errorf("assign host to app, but get module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)})
		return
	}

	params := make(map[string]interface{})
	params[common.BKAppIDField] = conf.ApplicationID
	params[common.BKHostIDField] = conf.HostID
	params[common.BKModuleIDField] = moduleID
	params["bk_owner_module_id"] = ownerModuleID
	params["bk_owner_biz_id"] = appID

	audit := s.Logics.NewHostModuleLog(pheader, conf.HostID)
	audit.WithPrevious()

	result, err := s.CoreAPI.HostController().Module().AssignHostToApp(context.Background(), pheader, params)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("assign host to app, but assign to app failed, err: %v, error message:%s", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostEditRelationPoolFail)})
		return
	}

	user := util.GetUser(pheader)
	if err := audit.SaveAudit(strconv.FormatInt(conf.ApplicationID, 10), user, "assign host to app"); err != nil {
		blog.Errorf("assign host to app, but save audit failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) AssignHostToAppModule(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	data := new(metadata.HostToAppModule)
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("assign host to app module failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	appID, _, moduleID, err := s.Logics.GetTopoIDByName(pheader, data)
	if nil != err {
		blog.Errorf("get app  topology id by name error:%s, msg: applicationName:%s, setName:%s, moduleName:%s", err.Error(), data.AppName, data.SetName, data.ModuleName)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrAddHostToModuleFailStr, "search application module not found ")})
		return
	}

	if 0 == appID || 0 == moduleID {
		//get default app
		ownerAppID, err := s.Logics.GetDefaultAppID(data.OwnerID, pheader)
		if err != nil {
			blog.Errorf("assign host to app module, but get resource pool failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAddHostToModuleFailStr)})
			return
		}
		if 0 == ownerAppID {
			blog.Errorf("assign host to app module, but get resource pool failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAddHostToModuleFailStr)})
			return
		}

		//get idle module
		mConds := make(map[string]interface{})
		mConds[common.BKDefaultField] = common.DefaultResModuleFlag
		mConds[common.BKModuleNameField] = common.DefaultResModuleName
		mConds[common.BKAppIDField] = ownerAppID
		ownerModuleID, err := s.Logics.GetResoulePoolModuleID(pheader, mConds)
		if nil != err {
			blog.Errorf("assign host to app module, but get unused host pool failed, ownerid[%v], err: %v", ownerModuleID, err)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAddHostToModuleFailStr)})
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
		err := s.Logics.EnterIP(pheader, util.GetOwnerID(req.Request.Header), appID, moduleID, ip, data.PlatID, host, data.IsIncrement)
		if nil != err {
			blog.Errorf("%s add host error: %s", ip, err.Error())
			errmsg = append(errmsg, fmt.Sprintf("%s add host error: %s", ip, err.Error()))
		}
	}
	if 0 == len(errmsg) {
		resp.WriteEntity(metadata.NewSuccessResp(nil))
		return
	} else {
		blog.Errorf("assign host to app module failed, err: %v", errmsg)
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAddHostToModuleFailStr)})
		return
	}
}

func (s *Service) moveHostToModuleByName(req *restful.Request, resp *restful.Response, moduleName string) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	conf := new(metadata.DefaultModuleHostConfigParams)
	if err := json.NewDecoder(req.Request.Body).Decode(&conf); err != nil {
		blog.Errorf("move host to module %s failed with decode body err: %v", moduleName, err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	conds := make(map[string]interface{})
	moduleNameLogKey := "idle"
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
	moduleID, err := s.Logics.GetResoulePoolModuleID(pheader, conds)
	if err != nil {
		blog.Errorf("move host to module %s, get module id err: %v", moduleName, err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrAddHostToModuleFailStr, conds[common.BKModuleNameField].(string)+" not foud ")})
		return
	}

	audit := s.Logics.NewHostModuleLog(pheader, conf.HostID)
	if err := audit.WithPrevious(); err != nil {
		blog.Errorf("move host to module %s, get prev module host config failed, err: %v", moduleName, err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
		return
	}

	for _, hostID := range conf.HostID {
		exist, err := s.Logics.IsHostExistInApp(conf.ApplicationID, hostID, pheader)
		if err != nil {
			blog.Error("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", conf.ApplicationID, hostID, err.Error())
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrHostNotINAPPFail)})
			return
		}

		if !exist {
			blog.Errorf("Host does not belong to the current application, appid: %v, hostid: %v", conf.ApplicationID, hostID)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Errorf(common.CCErrHostNotINAPP, hostID)})
			return
		}

		opt := metadata.ModuleHostConfigParams{
			ApplicationID: conf.ApplicationID,
			HostID:        hostID,
		}
		result, err := s.CoreAPI.HostController().Module().DelModuleHostConfig(context.Background(), pheader, &opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("move host to module %s, but delete old failed, err: %v, %v", err, result.ErrMsg)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrHostDELResourcePool)})
			return
		}

		opt.ModuleID = []int64{moduleID}
		result, err = s.CoreAPI.HostController().Module().AddModuleHostConfig(context.Background(), pheader, &opt)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("move host to module %s, but delete old failed, err: %v, %v", err, result.ErrMsg)
			resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: defErr.Error(common.CCErrAddHostToModuleFailStr)})
			return
		}

		user := util.GetUser(pheader)
		if err := audit.SaveAudit(strconv.FormatInt(conf.ApplicationID, 10), user, "host to "+moduleNameLogKey+" module"); err != nil {
			blog.Errorf("move host to module %s, save audit log failed, err: %v", moduleName, err)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")})
			return
		}

	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}
