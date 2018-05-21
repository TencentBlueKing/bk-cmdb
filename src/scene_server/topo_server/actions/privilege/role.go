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
 
package privilege

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
)

var role = &roleAction{}

type roleAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/privilege/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", Params: nil, Handler: role.CreatePrivilege})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/privilege/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", Params: nil, Handler: role.GetPrivilege})
	role.CreateAction()
}

//GetPrivilege get privilege
func (cli *roleAction) GetPrivilege(req *restful.Request, resp *restful.Response) {
	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		objID, _ := pathParams["bk_obj_id"]
		propertyID, _ := pathParams["bk_property_id"]

		//get privilege
		sPriURL := cli.CC.ObjCtrl() + "/object/v1/role/" + ownerID + "/" + objID + "/" + propertyID
		blog.Info("get privilege url: %s", sPriURL)
		priInfo, err := httpcli.ReqHttp(req, sPriURL, common.HTTPSelectGet, nil)
		if nil != err {
			blog.Error("get role pri error :%v", err)
			return http.StatusOK, make([]string, 0), nil
		}
		blog.Info("get privilege return: %s", priInfo)
		var result params.CommonResult
		err = json.Unmarshal([]byte(priInfo), &result)
		if nil != err || !result.Result {
			return http.StatusOK, make([]string, 0), nil
		}
		return http.StatusOK, result.Data, nil
	}, resp)
}

//CreatePrivilege create privilege
func (cli *roleAction) CreatePrivilege(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		objID, _ := pathParams["bk_obj_id"]
		propertyID, _ := pathParams["bk_property_id"]
		value, err := ioutil.ReadAll(req.Request.Body)
		//get privilege
		sPriURL := cli.CC.ObjCtrl() + "/object/v1/role/" + ownerID + "/" + objID + "/" + propertyID
		blog.Info("create privilege get url: %s", sPriURL)
		priInfo, err := httpcli.ReqHttp(req, sPriURL, common.HTTPSelectGet, nil)
		blog.Info("create privilege get return: %s", priInfo)
		if nil != err {
			blog.Error("create role pri error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoRolePrivilegeCreateFailed)
		}
		var result params.CommonResult
		err = json.Unmarshal([]byte(priInfo), &result)
		var createInfo string
		if nil == err && result.Result {
			if result.Result {
				blog.Info("create privilege return: update: %v", result.Data)
				createInfo, err = httpcli.ReqHttp(req, sPriURL, common.HTTPUpdate, value)
			} else {
				blog.Info("create privilege return: create1 %v", result.Data)
				createInfo, err = httpcli.ReqHttp(req, sPriURL, common.HTTPCreate, value)
			}
		} else {
			blog.Info("create privilege return: create2")
			createInfo, err = httpcli.ReqHttp(req, sPriURL, common.HTTPCreate, value)
		}

		if nil != err {
			blog.Error("create role pri error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoRolePrivilegeCreateFailed)
		}
		blog.Info("create privilege return: %s", createInfo)
		return http.StatusOK, createInfo, nil
	}, resp)

}
