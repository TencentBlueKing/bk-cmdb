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
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"net/http"

	"encoding/json"
	"io/ioutil"

	"github.com/emicklei/go-restful"
)

func init() {
	hostModuleConfig.CreateAction()

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/addhost", Params: nil, Handler: hostModuleConfig.AddHost})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/addhost", Params: nil, Handler: hostModuleConfig.AddHost})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/host", Params: nil, Handler: hostModuleConfig.AddHost})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/host/add/agent", Params: nil, Handler: hostModuleConfig.AddHostFromAgent})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/sync/new/host", Params: nil, Handler: hostModuleConfig.NewHostSyncAppTopo})

}

// AddHost add host
func (m *hostModuleConfigAction) AddHost(req *restful.Request, resp *restful.Response) {
	type hostList struct {
		ApplicationID int                            `json:"bk_biz_id"`
		HostInfo      map[int]map[string]interface{} `json:"host_info"`
		SupplierID    int                            `json:"bk_supplier_id"`
		InputType     string                         `json:"input_type"`
	}
	ownerID := common.BKDefaultOwnerID
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	defLang := m.CC.Lang.CreateDefaultCCLanguageIf(util.GetActionLanguage(req))
	m.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		var data hostList

		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		if nil == data.HostInfo {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostInfo")
		}

		//get default biz
		var appID = data.ApplicationID
		if 0 == appID {
			appID, err = logics.GetDefaultAppIDBySupplierID(req, data.SupplierID, common.BKAppIDField, m.CC.ObjCtrl(), defLang)

			if nil != err {
				return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, common.DefaultAppName)
			}
		}

		//get empty set
		conds := make(map[string]interface{})
		conds[common.BKDefaultField] = common.DefaultResModuleFlag
		conds[common.BKModuleNameField] = common.DefaultResModuleName
		conds[common.BKAppIDField] = appID

		moduleID, err := logics.GetSingleModuleID(req, conds, m.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, common.DefaultResModuleName)
		}

		err, succ, updateErrRow, errRow := logics.AddHost(req, ownerID, appID, data.HostInfo, data.InputType, []int{moduleID}, m.CC)

		retData := make(map[string]interface{})
		retData["success"] = succ

		if nil == err {
			return http.StatusOK, retData, nil
		} else {

			retData["error"] = errRow
			retData["update_error"] = updateErrRow

			return http.StatusInternalServerError, retData, defErr.Error(common.CCErrHostCreateFail)
		}
	}, resp)
}

// AddHost add host
func (m *hostModuleConfigAction) NewHostSyncAppTopo(req *restful.Request, resp *restful.Response) {
	type hostList struct {
		ApplicationID int                            `json:"bk_biz_id"`
		ModuleID      []int                          `json:"bk_module_id"`
		HostInfo      map[int]map[string]interface{} `json:"host_info"`
		SupplierID    int                            `json:"bk_supplier_id"`
		//InputType     string                         `json:"input_type"`
	}
	ownerID := common.BKDefaultOwnerID
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
	defLang := m.CC.Lang.CreateDefaultCCLanguageIf(util.GetActionLanguage(req))
	m.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		var data hostList

		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		if nil == data.HostInfo {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "host_info")
		}
		if common.BatchHostAddMaxRow < len(data.HostInfo) {
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommXXExceedLimit, "host_info ", common.BatchHostAddMaxRow)
		}
		if nil == data.ModuleID || 0 == len(data.ModuleID) {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKModuleIDField)
		}

		appConds := map[string]interface{}{
			common.BKAppIDField: data.ApplicationID,
		}
		_, err = logics.GetAppInfo(req, "", appConds, m.CC.ObjCtrl(), defLang)
		if nil != err {
			blog.Errorf("host sync app %d error:%s", data.ApplicationID, err.Error())
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrTopoGetAppFaild, err.Error())
		}

		data.ModuleID, err = logics.NewHostSyncValidModule(req, data.ApplicationID, data.ModuleID, m.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrTopoGetModuleFailed, err.Error())
		}

		err, succ, updateErrRow, errRow := logics.AddHost(req, ownerID, data.ApplicationID, data.HostInfo, common.InputTypeApiNewHostSync, data.ModuleID, m.CC)

		retData := make(map[string]interface{})
		retData["success"] = succ

		if nil == err {
			return http.StatusOK, retData, nil
		} else {

			retData["error"] = errRow
			retData["update_error"] = updateErrRow

			return http.StatusInternalServerError, retData, defErr.Error(common.CCErrHostCreateFail)
		}
	}, resp)
}

// AddHostFromAgent import host
func (m *hostModuleConfigAction) AddHostFromAgent(req *restful.Request, resp *restful.Response) {
	type hostList struct {
		HostInfo map[string]interface{}
		//ImportFrom string
	}
	ownerID := common.BKDefaultOwnerID

	var data hostList

	language := util.GetActionLanguage(req)
	defErr := m.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := m.CC.Lang.CreateDefaultCCLanguageIf(language)

	m.CallResponseEx(func() (int, interface{}, error) {

		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPBodyEmpty)
		}
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		if nil == data.HostInfo {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostInfo")
			//m.ResponseFailed(common.CC_Err_Comm_http_Input_Params, "主机参数不能为空", resp)
		}

		//get default app
		appID, err := logics.GetDefaultAppID(req, ownerID, common.BKAppIDField, m.CC.ObjCtrl(), defLang)

		if 0 == appID || nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrAddHostToModule, err.Error())
		}

		//get internal set
		conds := make(map[string]interface{})
		conds[common.BKDefaultField] = common.DefaultResModuleFlag
		conds[common.BKModuleNameField] = common.DefaultResModuleName
		conds[common.BKAppIDField] = appID

		moduleID, err := logics.GetSingleModuleID(req, conds, m.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrAddHostToModule, err.Error())
		}

		// get language
		language := util.GetActionLanguage(req)
		addHost := make(map[int]map[string]interface{})
		data.HostInfo["import_from"] = common.HostAddMethodAgent
		addHost[1] = data.HostInfo

		defErr := m.CC.Error.CreateDefaultCCErrorIf(language)

		err, _, updateErrRow, errRow := logics.AddHost(req, ownerID, appID, addHost, "", []int{moduleID}, m.CC)

		if nil == err {
			return http.StatusOK, nil, nil
		} else {
			var errString string
			if 0 < len(updateErrRow) {
				errString = updateErrRow[0]
			} else if 0 < len(errRow) {
				errString = errRow[0]
			}
			return http.StatusInternalServerError, resp, defErr.Errorf(common.CCErrAddHostToModuleFailStr, errString)

		}
	}, resp)
}
