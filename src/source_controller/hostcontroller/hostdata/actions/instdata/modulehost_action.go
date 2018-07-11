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
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/eventdata"
	"configcenter/src/source_controller/common/instdata"
	"configcenter/src/source_controller/hostcontroller/hostdata/logics"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
)

type ModuleHostConfigParams struct {
	ApplicationID int   `json:"bk_biz_id"`
	HostID        int   `json:"bk_host_id"`
	ModuleID      []int `json:"bk_module_id"`
}

var (
	moduleBaseTaleName = "cc_ModuleBase"
)

var moduleHostConfigActionCli *moduleHostConfigAction = &moduleHostConfigAction{}

// HostAction
type moduleHostConfigAction struct {
	base.BaseAction
}

//AddModuleHostConfig add module host config
func (cli *moduleHostConfigAction) AddModuleHostConfig(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := ModuleHostConfigParams{}
		if err := json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		hostID := params.HostID

		//add new relation ship
		ec := eventdata.NewEventContextByReq(req)
		for _, moduleID := range params.ModuleID {
			_, err := logics.AddSingleHostModuleRelation(ec, cc, hostID, moduleID, params.ApplicationID, ownerID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostTransferModule)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

// DelDefaultModuleHostConfig delete default module host config
func (cli *moduleHostConfigAction) DelDefaultModuleHostConfig(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := ModuleHostConfigParams{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		defaultModuleIDs, err := logics.GetDefaultModuleIDs(cc, params.ApplicationID, ownerID)
		if nil != err {
			blog.Errorf("defaultModuleIds appID:%d, error:%v", params.ApplicationID, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetModule)
		}

		hostID := params.HostID

		//delete default host module relation
		ec := eventdata.NewEventContextByReq(req)
		for _, defaultModuleID := range defaultModuleIDs {
			_, err := logics.DelSingleHostModuleRelation(ec, cc, hostID, defaultModuleID, params.ApplicationID, ownerID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrDelDefaultModuleHostConfig)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

//DelModuleHostConfig delete module host config
func (cli *moduleHostConfigAction) DelModuleHostConfig(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := ModuleHostConfigParams{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		getModuleParams := make(map[string]interface{}, 2)
		getModuleParams[common.BKHostIDField] = params.HostID
		getModuleParams[common.BKAppIDField] = params.ApplicationID
		getModuleParams = util.SetModOwner(getModuleParams, ownerID)
		moduleIDs, err := logics.GetModuleIDsByHostID(cc, getModuleParams) //params.HostID, params.ApplicationID)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetOriginHostModuelRelationship)
		}

		ec := eventdata.NewEventContextByReq(req)
		for _, moduleID := range moduleIDs {
			_, err := logics.DelSingleHostModuleRelation(ec, cc, params.HostID, moduleID, params.ApplicationID, ownerID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrDelOriginHostModuelRelationship)
			}
		}

		return http.StatusOK, nil, nil
	}, resp)
}

//GetHostModulesIDs get host module ids
func (cli *moduleHostConfigAction) GetHostModulesIDs(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := ModuleHostConfigParams{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		condition := map[string]interface{}{common.BKAppIDField: params.ApplicationID, common.BKHostIDField: params.HostID}
		condition = util.SetModOwner(condition, ownerID)
		moduleIDs, err := logics.GetModuleIDsByHostID(cc, condition) //params.HostID, params.ApplicationID)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetModule)
		}
		return http.StatusOK, moduleIDs, nil
	}, resp)
}

//AssignHostToApp assign host to app
func (cli *moduleHostConfigAction) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		type paramsStruct struct {
			ApplicationID      int   `json:"bk_biz_id"`
			HostID             []int `json:"bk_host_id"`
			ModuleID           int   `json:"bk_module_id"`
			OwnerApplicationID int   `json:"bk_owner_biz_id"`
			OwnerModuleID      int   `json:"bk_owner_module_id"`
		}

		cc := api.NewAPIResource()
		ec := eventdata.NewEventContextByReq(req)
		//instdata.DataH = cc.InstCli
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		params := paramsStruct{}
		if err := json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		getModuleParams := make(map[string]interface{})
		for _, hostID := range params.HostID {
			//delete relation in default app module
			_, err := logics.DelSingleHostModuleRelation(ec, cc, hostID, params.OwnerModuleID, params.OwnerApplicationID, ownerID)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTransferHostFromPool)
			}
			getModuleParams[common.BKHostIDField] = hostID
			getModuleParams = util.SetModOwner(getModuleParams, ownerID)
			moduleIDs, err := logics.GetModuleIDsByHostID(cc, getModuleParams)
			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetModule)
			}
			//delete from empty module, no relation
			if 0 < len(moduleIDs) {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrAlreadyAssign)
			}

			//add new host
			_, err = logics.AddSingleHostModuleRelation(ec, cc, hostID, params.ModuleID, params.ApplicationID, ownerID)
			if nil != err {
			}
		}
		return http.StatusOK, nil, nil
	}, resp)
	return

}

//GetModulesHostConfig  get module host config
func (cli *moduleHostConfigAction) GetModulesHostConfig(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		var params = make(map[string][]int)
		cc := api.NewAPIResource()
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		query := make(map[string]interface{})
		for key, val := range params {
			conditon := make(map[string]interface{})
			conditon[common.BKDBIN] = val
			query[key] = conditon
		}
		query = util.SetModOwner(query, ownerID)
		fields := []string{common.BKAppIDField, common.BKHostIDField, common.BKSetIDField, common.BKModuleIDField}
		var result []interface{}
		err = cc.InstCli.GetMutilByCondition("cc_ModuleHostConfig", fields, query, &result, common.BKHostIDField, 0, 100000)
		if err != nil {
			blog.Error("fail to get module host config %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		return http.StatusOK, result, nil
	}, resp)
}

//MoveHostToSourcePool move host 2 resource pool
func (cli *moduleHostConfigAction) MoveHost2ResourcePool(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		type paramsStruct struct {
			ApplicationID       int   `json:"bk_biz_id"`
			HostID              []int `json:"bk_host_id"`
			OwnerModuleID       int   `json:"bk_owner_module_id"`
			OwnerAppplicationID int   `json:"bk_owner_biz_id"`
		}

		cc := api.NewAPIResource()
		ec := eventdata.NewEventContextByReq(req)
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		instdata.DataH = cc.InstCli
		params := paramsStruct{}
		if err = json.Unmarshal([]byte(value), &params); nil != err {
			blog.Error("fail to unmarshal json, error information is %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		idleModuleID, err := logics.GetIDleModuleID(cc, params.ApplicationID)
		if nil != err {
			blog.Error("获取业务默认模块失败 error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetModule)
		}
		errHostIDs, faultHostIDs, err := logics.CheckHostInIDle(cc, params.ApplicationID, idleModuleID, params.HostID)
		if nil != err {
			blog.Error("获取主机模块关系失败， error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrGetModule)
		}
		if 0 != len(errHostIDs) {
			data := common.KvMap{common.BKHostIDField: errHostIDs}
			blog.Errorf("主机属于空闲机以外的模块 %v", data)
			return http.StatusInternalServerError, data, defErr.Error(common.CCErrNotBelongToIdleModule)
		}
		var succ, addErr, delErr []int
		for _, hostID := range params.HostID {

			//host not belong to other biz, add new host
			if !util.ContainsInt(faultHostIDs, hostID) {
				_, err = logics.AddSingleHostModuleRelation(ec, cc, hostID, params.OwnerModuleID, params.OwnerAppplicationID, ownerID)
				if nil != err {
					addErr = append(addErr, hostID)
					continue
				}
			}

			//delete origin relation
			_, err := logics.DelSingleHostModuleRelation(ec, cc, hostID, idleModuleID, params.ApplicationID, ownerID)
			if nil != err {
				delErr = append(delErr, hostID)
				continue
			}

			succ = append(succ, hostID)
		}
		if 0 != len(addErr) || 0 != len(delErr) {
			addErr = append(addErr, delErr...)
			data := common.KvMap{"成功": succ, "失败": addErr}
			blog.Errorf("主机属于空闲机以外的模块 %v", data)
			return http.StatusInternalServerError, data, defErr.Error(common.CCErrTransfer2ResourcePool)
		}

		return http.StatusOK, nil, nil
	}, resp)
}

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
