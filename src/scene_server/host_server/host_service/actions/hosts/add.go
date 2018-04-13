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
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/host/add/agent", Params: nil, Handler: hostModuleConfig.AddHostFromAgent})

}

// AddHost add host
func (m *hostModuleConfigAction) AddHost(req *restful.Request, resp *restful.Response) {
	type hostList struct {
		ApplicationID int                            `json:"bk_biz_id"`
		HostInfo      map[int]map[string]interface{} `json:"host_info"`
		SupplierID    int                            `json:"bk_supplier_id"`
	}
	ownerID := common.BKDefaultOwnerID
	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))
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
		appID, err := logics.GetDefaultAppIDBySupplierID(req, data.SupplierID, "bk_biz_id", m.CC.ObjCtrl())

		if 0 == appID || nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, common.DefaultAppName)
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

		err, succ, updateErrRow, errRow := logics.AddHost(req, ownerID, appID, data.HostInfo, moduleID, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl(), defErr)

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
	value, err := ioutil.ReadAll(req.Request.Body)
	var data hostList

	defErr := m.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	m.CallResponseEx(func() (int, interface{}, error) {

		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		if nil == data.HostInfo {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "HostInfo")

			m.ResponseFailed(common.CC_Err_Comm_http_Input_Params, "主机参数不能为空", resp)
		}

		//get default app
		appID, err := logics.GetDefaultAppID(req, ownerID, common.BKAppIDField, m.CC.ObjCtrl())

		if 0 == appID || nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, err.Error())
		}

		//get internal set
		conds := make(map[string]interface{})
		conds[common.BKDefaultField] = common.DefaultResModuleFlag
		conds[common.BKModuleNameField] = common.DefaultResModuleName
		conds[common.BKAppIDField] = appID

		moduleID, err := logics.GetSingleModuleID(req, conds, m.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, err.Error())
		}

		// get language
		language := util.GetActionLanguage(req)
		addHost := make(map[int]map[string]interface{})
		data.HostInfo["import_from"] = common.HostAddMethodAgent
		addHost[1] = data.HostInfo

		defErr := m.CC.Error.CreateDefaultCCErrorIf(language)

		err, _, updateErrRow, errRow := logics.AddHost(req, ownerID, appID, addHost, moduleID, m.CC.HostCtrl(), m.CC.ObjCtrl(), m.CC.AuditCtrl(), defErr)

		if nil == err {
			return http.StatusOK, nil, nil
		} else {
			var errString string
			if 0 < len(updateErrRow) {
				errString = updateErrRow[0]
			} else if 0 < len(errRow) {
				errString = errRow[0]
			}
			return http.StatusInternalServerError, resp, defErr.Errorf(common.CCErrHostModuleRelationAddFailed, errString)

		}
	}, resp)
}
