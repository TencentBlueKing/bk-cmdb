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

package hosts

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"net/http"
	"strings"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/emicklei/go-restful"
)

var hostModuleConfig *hostModuleConfigAction = &hostModuleConfigAction{}

type hostModuleConfigAction struct {
	base.BaseAction
}

type moduleHostConfigParams struct {
	ApplicationID int   `json:"bk_biz_id"`
	HostID        []int `json:"bk_host_id"`
	ModuleID      []int `json:"bk_module_id"`
	IsIncrement   bool  `json:"is_increment"`
}

type defaultModuleHostConfigParams struct {
	ApplicationID int   `json:"bk_biz_id"`
	HostID        []int `json:"bk_host_id"`
}

func init() {
	hostModuleConfig.CreateAction()

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/modules", Params: nil, Handler: hostModuleConfig.HostModuleRelation})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/emptymodule", Params: nil, Handler: hostModuleConfig.MoveHost2EmptyModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/faultmodule", Params: nil, Handler: hostModuleConfig.MoveHost2FaultModule})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/resource", Params: nil, Handler: hostModuleConfig.MoveHostToResourcePool})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/assgin", Params: nil, Handler: hostModuleConfig.AssignHostToApp})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/host/add/module", Params: nil, Handler: hostModuleConfig.AssignHostToAppModule})

}

// HostModuleRelation add host module relation
func (m *hostModuleConfigAction) HostModuleRelation(req *restful.Request, resp *restful.Response) {
	value, err := ioutil.ReadAll(req.Request.Body)
	var data moduleHostConfigParams
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	m.CallResponseEx(func() (int, interface{}, error) {
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		for _, moduleID := range data.ModuleID {
			//校验目标模块是否存在
			module, err := logics.GetModuleByModuleID(req, data.ApplicationID, moduleID, m.CC.ObjCtrl())
			if nil != err {
				blog.Error("get dstmdouel info error, params:%v, error:%v", data.ModuleID, err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoModuleSelectFailed)
			}
			if 0 == len(module) {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoMulueIDNotfoundFailed)
			}
		}

		logClient, err := logics.NewHostModuleConfigLog(req, data.HostID, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommResourceInitFailed)
		}

		for _, hostID := range data.HostID {
			bl, err := logics.IsExistHostIDInApp(m.CC, req, data.ApplicationID, hostID)
			if nil != err {
				blog.Error("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", data.ApplicationID, hostID, err.Error())
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostNotINAPPFail, hostID)

			}
			if false == bl {
				blog.Error("Host does not belong to the current application; error, params:{appid:%d, hostid:%s}", data.ApplicationID, hostID)
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostNotINAPP, hostID)
			}

			params := make(map[string]interface{})
			delModulesURL := ""
			params[common.BKAppIDField] = data.ApplicationID
			params[common.BKHostIDField] = hostID

			if data.IsIncrement {
				delModulesURL = m.CC.HostCtrl() + "/host/v1/meta/hosts/defaultmodules"
			} else {
				delModulesURL = m.CC.HostCtrl() + "/host/v1/meta/hosts/modules"

			}
			isSuccess, errMsg, _ := logics.GetHttpResult(req, delModulesURL, common.HTTPDelete, params)
			if !isSuccess {
				blog.Error("remove hosthostconfig error, params:%v, error:%s", params, errMsg)
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostDELResourcePool, hostID)
			}

			addModulesURL := m.CC.HostCtrl() + "/host/v1/meta/hosts/modules"

			params[common.BKModuleIDField] = data.ModuleID
			isSuccess, errMsg, _ = logics.GetHttpResult(req, addModulesURL, common.HTTPCreate, params)
			if !isSuccess {
				blog.Error("add hosthostconfig error, params:%v, error:%s", params, errMsg)
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostAddRelationFail, hostID)

			}
		}
		user := util.GetActionUser(req)
		logClient.SaveLog(fmt.Sprintf("%d", data.ApplicationID), user)

		return http.StatusOK, nil, nil
	}, resp)

}

//MoveHost2EmptyModule move host to empty module
func (m *hostModuleConfigAction) MoveHost2EmptyModule(req *restful.Request, resp *restful.Response) {

	m.moveHostToModuleByName(req, resp, common.DefaultResModuleName)
}

//MoveHost2FaultModule move host 2 fault module
func (m *hostModuleConfigAction) MoveHost2FaultModule(req *restful.Request, resp *restful.Response) {
	m.moveHostToModuleByName(req, resp, common.DefaultFaultModuleName)
}

//MoveHostToResourcePool move host to resource pool
func (m *hostModuleConfigAction) MoveHostToResourcePool(req *restful.Request, resp *restful.Response) {
	value, err := ioutil.ReadAll(req.Request.Body)
	var data defaultModuleHostConfigParams
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	m.CallResponseEx(func() (int, interface{}, error) {
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		reply, err := logics.MoveHost2ResourcePool(m.CC, req, data.ApplicationID, data.HostID)

		if err != nil {
			return http.StatusInternalServerError, reply, defErr.Errorf(common.CCErrHostMoveResourcePoolFail, err.Error())

		} else {
			return http.StatusOK, nil, nil
		}
	}, resp)

}

//AssignHostToApp assign host to app
func (m *hostModuleConfigAction) AssignHostToApp(req *restful.Request, resp *restful.Response) {
	value, err := ioutil.ReadAll(req.Request.Body)
	var data defaultModuleHostConfigParams
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	m.CallResponseEx(func() (int, interface{}, error) {
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		conds := make(map[string]interface{}, 1)
		conds[common.BKAppIDField] = data.ApplicationID
		fields := fmt.Sprintf("%s,%s", common.BKOwnerIDField, common.BKAppNameField)
		appinfo, err := logics.GetAppInfo(req, fields, conds, m.CC.ObjCtrl())
		if err != nil {
			m.ResponseFailed(common.CC_Err_Comm_APP_QUERY_FAIL, err.Error(), resp)
		}
		ownerID := appinfo[common.BKOwnerIDField].(string)
		if "" == ownerID {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "OwnerID")

		}

		//get default app
		appID, err := logics.GetDefaultAppID(req, ownerID, common.BKAppIDField, m.CC.ObjCtrl())
		blog.Infof("ownerid %s default appid %d", ownerID, appID)
		if err != nil {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoAppSearchFailed)
		}
		if 0 == appID {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommNotFound)
		}
		if appID == data.ApplicationID {
			return http.StatusOK, nil, nil
		}
		user := util.GetActionUser(req)

		//get resource empty set
		mConds := make(map[string]interface{})
		mConds[common.BKDefaultField] = common.DefaultResModuleFlag
		mConds[common.BKModuleNameField] = common.DefaultResModuleName
		mConds[common.BKAppIDField] = appID
		ownerModuleID, err := logics.GetSingleModuleID(req, mConds, m.CC.ObjCtrl())
		blog.Infof("ownerid %s default appid %d idle moduleID %d", ownerID, appID, ownerModuleID)

		if nil != err {
			blog.Errorf("ownerid %s default appid %d idle moduleID not found", ownerID, appID)
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrTopoMulueIDNotfoundFailed)
		}
		if 0 == ownerModuleID {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrTopoMulueIDNotfoundFailed)
		}

		//current app empty set
		mConds = make(map[string]interface{})
		mConds[common.BKDefaultField] = common.DefaultResModuleFlag
		mConds[common.BKModuleNameField] = common.DefaultResModuleName
		mConds[common.BKAppIDField] = data.ApplicationID
		moduleID, err := logics.GetSingleModuleID(req, mConds, m.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrTopoMulueIDNotfoundFailed)
		}
		assignModulesURL := m.CC.HostCtrl() + "/host/v1/meta/hosts/assign"
		params := make(map[string]interface{})
		params[common.BKAppIDField] = data.ApplicationID
		params[common.BKHostIDField] = data.HostID
		params[common.BKModuleIDField] = moduleID
		params["bk_owner_module_id"] = ownerModuleID
		params["bk_owner_biz_id"] = appID
		isSuccess, errMsg, _ := logics.GetHttpResult(req, assignModulesURL, common.HTTPCreate, params)
		if !isSuccess {
			blog.Error("add hostconfig error, params:%v, error:%s", params, errMsg)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostEditRelationPoolFail)
		}
		logClient, err := logics.NewHostModuleConfigLog(req, nil, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl())
		logClient.SetDesc(fmt.Sprintf("分配主机到业务[%s]", appinfo[common.BKAppNameField].(string)))
		logClient.SetHostID(data.HostID)
		logClient.SaveLog(fmt.Sprintf("%d", data.ApplicationID), user)

		return http.StatusOK, nil, nil
	}, resp)

}

// AssignHostToAppModule 将某一个ip分配到具体业务下的模块， enterip使用
func (m *hostModuleConfigAction) AssignHostToAppModule(req *restful.Request, resp *restful.Response) {

	type inputStruct struct {
		Ips        []string `json:"ips"`
		HostName   []string `json:"bk_host_name"`
		ModuleName string   `json:"bk_module_name"`
		SetName    string   `json:"bk_set_name"`
		AppName    string   `json:"bk_biz_name"`
		OsType     string   `json:"bk_os_type"`
		OwnerID    string   `json:"bk_supplier_account"`
	}
	language := util.GetActionLanguage(req)
	defErr := m.CC.Error.CreateDefaultCCErrorIf(language)

	m.CallResponseEx(func() (int, interface{}, error) {
		value, _ := ioutil.ReadAll(req.Request.Body)
		var data inputStruct
		err := json.Unmarshal([]byte(value), &data)
		if nil != err {
			blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		appID, _, moduleID, err := logics.GetTopoIDByName(req, data.OwnerID, data.AppName, data.SetName, data.ModuleName, m.CC.ObjCtrl(), defErr)
		if nil != err {
			blog.Error("get app  topology id by name error:%s, msg: applicationName:%s, setName:%s, moduleName:%s", err.Error(), data.AppName, data.SetName, data.ModuleName)
			return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, "search appliaction module not foud ")
		}

		var strHostName string
		if 0 == appID || 0 == moduleID {
			//get default app
			ownerAppID, err := logics.GetDefaultAppID(req, data.OwnerID, common.BKAppIDField, m.CC.ObjCtrl())
			if err != nil {
				blog.Infof("ownerid %s 资源池未找到", ownerAppID)
				return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, "search "+common.DefaultAppName+" not foud ")
			}
			if 0 == ownerAppID {
				blog.Infof("ownerid %s 资源池未找到", ownerAppID)
				return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, common.DefaultAppName+" not foud ")
			}

			//get idle module
			mConds := make(map[string]interface{})
			mConds[common.BKDefaultField] = common.DefaultResModuleFlag
			mConds[common.BKModuleNameField] = common.DefaultResModuleName
			mConds[common.BKAppIDField] = ownerAppID
			ownerModuleID, err := logics.GetSingleModuleID(req, mConds, m.CC.ObjCtrl())
			if nil != err {
				blog.Infof("ownerid %s 资源池业务空闲机未找到", ownerAppID)
				return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, common.DefaultResModuleName+" not foud ")
			}
			appID = ownerAppID
			moduleID = ownerModuleID
			data.AppName = common.DefaultAppName
			data.SetName = ""
			data.ModuleName = common.DefaultResModuleName

		}
		var errmsg []string
		for index, ip := range data.Ips {
			if index < len(data.HostName) {
				strHostName = data.HostName[index]
			} else {
				strHostName = ""
			}

			//dispatch to app
			err := logics.EnterIP(req, data.OwnerID, appID, moduleID, ip, data.OsType, strHostName, data.AppName, data.SetName, data.ModuleName, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl(), defErr)
			if nil != err {
				blog.Errorf("%s add host error: %s", ip, err.Error())
				errmsg = append(errmsg, fmt.Sprintf("%s add host error: %s", ip, err.Error()))
			}
		}
		if 0 == len(errmsg) {
			return http.StatusOK, nil, nil
		} else {
			return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, strings.Join(errmsg, ","))
		}
	}, resp)

}

//moveHostToModuleName translate module to idle and fault module relation
func (m *hostModuleConfigAction) moveHostToModuleByName(req *restful.Request, resp *restful.Response, moduleName string) {
	value, err := ioutil.ReadAll(req.Request.Body)
	var data defaultModuleHostConfigParams
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	m.CallResponseEx(func() (int, interface{}, error) {
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		//fmt.Println(moduleURL)
		conds := make(map[string]interface{})
		if common.DefaultResModuleName == moduleName {
			//空闲机
			conds[common.BKDefaultField] = common.DefaultResModuleFlag
			conds[common.BKModuleNameField] = common.DefaultResModuleName
		} else {
			//故障机器
			conds[common.BKDefaultField] = common.DefaultFaultModuleFlag
			conds[common.BKModuleNameField] = common.DefaultFaultModuleName
		}

		conds[common.BKAppIDField] = data.ApplicationID
		moduleID, err := logics.GetSingleModuleID(req, conds, m.CC.ObjCtrl())
		if nil != err {
			return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, conds[common.BKModuleNameField].(string)+" not foud ")

		}
		moduleHostConfigParams := make(map[string]interface{})
		moduleHostConfigParams[common.BKAppIDField] = data.ApplicationID
		logClient, err := logics.NewHostModuleConfigLog(req, data.HostID, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl())
		if nil != err {

			return http.StatusBadGateway, nil, defErr.Errorf(common.CCErrCommResourceInitFailed, "audit server")
		}

		for _, hostID := range data.HostID {
			bl, err := logics.IsExistHostIDInApp(m.CC, req, data.ApplicationID, hostID)
			if nil != err {
				blog.Error("check host is exist in app error, params:{appid:%d, hostid:%s}, error:%s", data.ApplicationID, hostID, err.Error())
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostNotINAPPFail, hostID)
			}
			if false == bl {
				blog.Error("Host does not belong to the current application; error, params:{appid:%d, hostid:%s}", data.ApplicationID, hostID)
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostNotINAPP, hostID)
			}

			moduleHostConfigParams[common.BKHostIDField] = hostID
			delModulesURL := m.CC.HostCtrl() + "/host/v1/meta/hosts/modules"
			isSuccess, errMsg, _ := logics.GetHttpResult(req, delModulesURL, common.HTTPDelete, moduleHostConfigParams)
			if !isSuccess {
				blog.Error("remove hosthostconfig error, params:%v, error:%s", moduleHostConfigParams, errMsg)
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed)
			}
			moduleHostConfigParams[common.BKModuleIDField] = []int{moduleID}
			addModulesURL := m.CC.HostCtrl() + "/host/v1/meta/hosts/modules"

			isSuccess, errMsg, _ = logics.GetHttpResult(req, addModulesURL, common.HTTPCreate, moduleHostConfigParams)
			if !isSuccess {
				blog.Error("add hosthostconfig error, params:%v, error:%s", moduleHostConfigParams, errMsg)
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, errMsg)
			}
		}
		user := util.GetActionUser(req)
		logClient.SetDesc("转移主机到" + moduleName)
		logClient.SaveLog(fmt.Sprintf("%d", data.ApplicationID), user)

		return http.StatusOK, nil, nil
	}, resp)
}
