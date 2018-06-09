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

package instdata

import (
	"encoding/json"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/common/instdata"
	"configcenter/src/source_controller/hostcontroller/hostdata/logics"
	"github.com/emicklei/go-restful"
)

func init() {
	moduleHostConfigActionCli.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/hosts/modules/search", Params: nil, Handler: moduleHostConfigActionCli.GetHostModulesIDs})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/hosts/modules", Params: nil, Handler: moduleHostConfigActionCli.AddModuleHostConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/hosts/modules", Params: nil, Handler: moduleHostConfigActionCli.DelModuleHostConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/meta/hosts/defaultmodules", Params: nil, Handler: moduleHostConfigActionCli.DelDefaultModuleHostConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/meta/hosts/resource", Params: nil, Handler: moduleHostConfigActionCli.MoveHost2ResourcePool})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/meta/hosts/assign", Params: nil, Handler: moduleHostConfigActionCli.AssignHostToApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/meta/hosts/module/config/search", Params: nil, Handler: moduleHostConfigActionCli.GetModulesHostConfig})
}

var moduleBaseTaleName = "cc_ModuleBase"
var moduleHostConfigActionCli *moduleHostConfigAction = &moduleHostConfigAction{}

// HostAction
type moduleHostConfigAction struct {
	base.BaseAction
}

//AddModuleHostConfig add module host config
func (cli *moduleHostConfigAction) AddModuleHostConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("add module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	cc := api.NewAPIResource()
	ec := eventdata.NewEventContextByReq(req)
	for _, moduleID := range params.ModuleID {
		_, err := logics.AddSingleHostModuleRelation(ec, cc, params.HostID, moduleID, params.ApplicationID)
		if nil != err {
			resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostTransferModule)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

//DelDefaultModuleHostConfig delete default module host config
func (cli *moduleHostConfigAction) DelDefaultModuleHostConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("del default module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	cc := api.NewAPIResource()
	defaultModuleIDs, err := logics.GetDefaultModuleIDs(cc, params.ApplicationID)
	if nil != err {
		blog.Errorf("defaultModuleIds appID:%d, error:%v", params.ApplicationID, err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	//delete default host module relation
	ec := eventdata.NewEventContextByReq(req)
	for _, defaultModuleID := range defaultModuleIDs {
		_, err := logics.DelSingleHostModuleRelation(ec, cc, params.HostID, defaultModuleID, params.ApplicationID)
		if nil != err {
			blog.Errorf("del default module host config failed, with relation, err:%v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrDelDefaultModuleHostConfig)})
			return
		}
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
}

//DelModuleHostConfig delete module host config
func (cli *moduleHostConfigAction) DelModuleHostConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("del module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	getModuleParams := make(map[string]interface{}, 2)
	getModuleParams[common.BKHostIDField] = params.HostID
	getModuleParams[common.BKAppIDField] = params.ApplicationID
	moduleIDs, err := logics.GetModuleIDsByHostID(cc, getModuleParams) //params.HostID, params.ApplicationID)
	if nil != err {
		blog.Errorf("delete module host config failed, %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetOriginHostModuelRelationship)})
		return
	}

	ec := eventdata.NewEventContextByReq(req)
	for _, moduleID := range moduleIDs {
		_, err := logics.DelSingleHostModuleRelation(ec, cc, params.HostID, moduleID, params.ApplicationID)
		if nil != err {
			blog.Errorf("delete module host config, but delete module relation failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrDelOriginHostModuelRelationship)})
			return
		}
	}

	resp.WriteAsJson(meta.NewSuccessResp(nil))
}

//GetHostModulesIDs get host module ids
func (cli *moduleHostConfigAction) GetHostModulesIDs(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	params := meta.ModuleHostConfigParams{}
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Error("get host module id failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	moduleIDs, err := logics.GetModuleIDsByHostID(cc, map[string]interface{}{common.BKAppIDField: params.ApplicationID, common.BKHostIDField: params.HostID}) //params.HostID, params.ApplicationID)
	if nil != err {
		blog.Errorf("get host module id failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	resp.WriteAsJson(meta.GetHostModuleIDsResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     moduleIDs,
	})
}

//AssignHostToApp assign host to app
func (cli *moduleHostConfigAction) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	ec := eventdata.NewEventContextByReq(req)
	params := new(meta.AssignHostToAppParams)
	if err := json.NewDecoder(req.Request.Body).Decode(params); err != nil {
		blog.Errorf("assign host to app failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	getModuleParams := make(map[string]interface{})
	for _, hostID := range params.HostID {
		// delete relation in default app module
		_, err := logics.DelSingleHostModuleRelation(ec, cc, hostID, params.OwnerModuleID, params.OwnerApplicationID)
		if nil != err {
			blog.Errorf("assign host to app, but delete host module relationship failed, err: %v")
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrTransferHostFromPool)})
			return
		}

		getModuleParams[common.BKHostIDField] = hostID
		moduleIDs, err := logics.GetModuleIDsByHostID(cc, getModuleParams)
		if nil != err {
			blog.Errorf("assign host to app, but get module failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
			return
		}

		// delete from empty module, no relation
		if 0 < len(moduleIDs) {
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrAlreadyAssign)})
			return
		}

		// add new host
		_, err = logics.AddSingleHostModuleRelation(ec, cc, hostID, params.ModuleID, params.ApplicationID)
		if nil != err {
			blog.Errorf("assign host to app, but add single host module relation failed, err: %v", err)
			resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrTransferHostFromPool)})
		}
	}

	resp.WriteAsJson(meta.SuccessBaseResp)
}

//GetModulesHostConfig  get module host config
func (cli *moduleHostConfigAction) GetModulesHostConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	var params = make(map[string][]int)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("del module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	query := make(map[string]interface{})
	for key, val := range params {
		conditon := make(map[string]interface{})
		conditon[common.BKDBIN] = val
		query[key] = conditon
	}

	fields := []string{common.BKAppIDField, common.BKHostIDField, common.BKSetIDField, common.BKModuleIDField}
	cc := api.NewAPIResource()
	var result []interface{}
	err := cc.InstCli.GetMutilByCondition("cc_ModuleHostConfig", fields, query, &result, common.BKHostIDField, 0, 100000)
	if err != nil {
		blog.Error("get module host config failed, err: %v", err)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteAsJson(meta.HostConfig{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}

//MoveHostToSourcePool move host 2 resource pool
func (cli *moduleHostConfigAction) MoveHost2ResourcePool(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cc := api.NewAPIResource()
	ec := eventdata.NewEventContextByReq(req)
	instdata.DataH = cc.InstCli
	params := new(meta.ParamData)
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("move host to resourece pool failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	idleModuleID, err := logics.GetIDleModuleID(cc, params.ApplicationID)
	if nil != err {
		blog.Error("get default module failed, error:%s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	errHostIDs, faultHostIDs, err := logics.CheckHostInIDle(cc, params.ApplicationID, idleModuleID, params.HostID)
	if nil != err {
		blog.Error("get host relationship failed, err: %s", err.Error())
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrGetModule)})
		return
	}

	if 0 != len(errHostIDs) {
		blog.Errorf("move host to resource pool, but it does not belongs to free module, hostid: %v", errHostIDs)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrNotBelongToIdleModule)})
		return
	}

	var succ, addErr, delErr []int
	for _, hostID := range params.HostID {
		//host not belong to other biz, add new host
		if !util.ContainsInt(faultHostIDs, hostID) {
			_, err = logics.AddSingleHostModuleRelation(ec, cc, hostID, params.OwnerModuleID, params.OwnerAppplicationID)
			if nil != err {
				addErr = append(addErr, hostID)
				continue
			}
		}

		//delete origin relation
		_, err := logics.DelSingleHostModuleRelation(ec, cc, hostID, idleModuleID, params.ApplicationID)
		if nil != err {
			delErr = append(delErr, hostID)
			continue
		}
		succ = append(succ, hostID)
	}

	if 0 != len(addErr) || 0 != len(delErr) {
		addErr = append(addErr, delErr...)
		blog.Errorf("move host to resource pool, success: %v, failed: %v", succ, addErr)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: defErr.Error(common.CCErrTransfer2ResourcePool)})
		return
	}

	resp.WriteAsJson(meta.NewSuccessResp(nil))
}
