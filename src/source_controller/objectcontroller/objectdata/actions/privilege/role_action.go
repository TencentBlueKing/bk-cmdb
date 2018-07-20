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
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
)

var role = &roleAction{}

// roleAction
type roleAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", Params: nil, Handler: role.CreateRolePri})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", Params: nil, Handler: role.GetRolePri})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/role/{bk_supplier_account}/{bk_obj_id}/{bk_property_id}", Params: nil, Handler: role.UpdateRolePri})
	// set cc api interface
	role.CreateAction()
}

// GetRolePri get role privilege
func (cli *roleAction) GetRolePri(req *restful.Request, resp *restful.Response) {
	// get the language
	//language := util.GetActionLanguage(req)
	// get the error factory by the language
	//defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)
	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()
		pathParams := req.PathParameters()
		objID := pathParams["bk_obj_id"]
		propertyID := pathParams["bk_property_id"]
		cond := make(map[string]interface{})
		cond[common.BKOwnerIDField] = ownerID
		cond[common.BKObjIDField] = objID
		cond[common.BKPropertyIDField] = propertyID
		cond = util.SetModOwner(cond, ownerID)
		var result map[string]interface{}
		err := cli.CC.InstCli.GetOneByCondition(common.BKTableNamePrivilege, []string{}, cond, &result)
		if nil != err {
			blog.Error("get role pri field error :%v", err)
			info := make([]string, 0)
			return http.StatusNotFound, info, errors.New("not found")

		}
		privilege, ok := result["privilege"]
		if false == ok {
			blog.Error("get role pri field error :%v", err)
			info := make(map[string]interface{})
			return http.StatusOK, info, nil

		}
		return http.StatusOK, privilege, nil
	}, resp)
}

//CreateRolePri create role privilege
func (cli *roleAction) CreateRolePri(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		objID := pathParams["bk_obj_id"]
		propertyID := pathParams["bk_property_id"]
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		var roleJSON []string
		err = json.Unmarshal([]byte(value), &roleJSON)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input := make(map[string]interface{})
		input[common.BKOwnerIDField] = ownerID
		input[common.BKObjIDField] = objID
		input[common.BKPropertyIDField] = propertyID
		input[common.BKPrivilegeField] = roleJSON
		input = util.SetModOwner(input, ownerID)
		_, err = cli.CC.InstCli.Insert(common.BKTableNamePrivilege, input)
		if nil != err {
			blog.Error("create role privilege error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//UpdateRolePri update role privilege
func (cli *roleAction) UpdateRolePri(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()
		pathParams := req.PathParameters()
		objID := pathParams["bk_obj_id"]
		propertyID := pathParams["bk_property_id"]
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		var roleJSON []string
		err = json.Unmarshal([]byte(value), &roleJSON)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		input := make(map[string]interface{})
		cond := make(map[string]interface{})
		cond[common.BKOwnerIDField] = ownerID
		cond[common.BKObjIDField] = objID
		cond[common.BKPropertyIDField] = propertyID
		input[common.BKPrivilegeField] = roleJSON
		cond = util.SetModOwner(cond, ownerID)
		input = util.SetModOwner(input, ownerID)
		err = cli.CC.InstCli.UpdateByCondition(common.BKTableNamePrivilege, input, cond)
		if nil != err {
			blog.Error("update role privilege error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}
