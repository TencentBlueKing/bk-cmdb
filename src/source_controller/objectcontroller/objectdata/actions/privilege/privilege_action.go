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
	"io/ioutil"
	"net/http"

	"configcenter/src/common/util"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var privilege = &privilegeAction{}

//privilege Action
type privilegeAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/privilege/group/detail/{bk_supplier_account}/{group_id}", Params: nil, Handler: privilege.CreateUserGroupPrivi})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/privilege/group/detail/{bk_supplier_account}/{group_id}", Params: nil, Handler: privilege.UpdateUserGroupPrivi})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/privilege/group/detail/{bk_supplier_account}/{group_id}", Params: nil, Handler: privilege.GetUserGroupPrivi})
	// set cc api interface
	privilege.CreateAction()
}

//CreateUserGroupPrivi create group privi
func (cli *privilegeAction) CreateUserGroupPrivi(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()

		pathParams := req.PathParameters()
		groupID := pathParams["group_id"]
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("insert user group privi failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		info, err := js.Map()
		if err != nil {
			blog.Error("insert user group privi failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		data := make(map[string]interface{})
		data[common.BKOwnerIDField] = ownerID
		data[common.BKUserGroupIDField] = groupID
		data[common.BKPrivilegeField] = info
		data = util.SetModOwner(data, ownerID)
		_, err = cli.CC.InstCli.Insert(common.BKTableNameUserGroupPrivilege, data)
		if nil != err {
			blog.Error("insert user group privi error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//UpdateUserGroupPrivi update group privi
func (cli *privilegeAction) UpdateUserGroupPrivi(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()

		pathParams := req.PathParameters()
		groupID := pathParams["group_id"]
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("update user group privi failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		info, err := js.Map()
		if err != nil {
			blog.Error("update user group privi failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		cond := make(map[string]interface{})
		data := make(map[string]interface{})
		cond[common.BKOwnerIDField] = ownerID
		cond[common.BKUserGroupIDField] = groupID
		data[common.BKPrivilegeField] = info
		cond = util.SetModOwner(cond, ownerID)
		data = util.SetModOwner(data, ownerID)
		err = cli.CC.InstCli.UpdateByCondition(common.BKTableNameUserGroupPrivilege, data, cond)
		if nil != err {
			blog.Error("update user group privi error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//GetUserGroupPrivi get group privi
func (cli *privilegeAction) GetUserGroupPrivi(req *restful.Request, resp *restful.Response) {

	//get the language
	language := util.GetActionLanguage(req)
	//get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()

		pathParams := req.PathParameters()
		groupID := pathParams["group_id"]

		cond := make(map[string]interface{})
		cond[common.BKOwnerIDField] = ownerID
		cond[common.BKUserGroupIDField] = groupID
		cond = util.SetModOwner(cond, ownerID)
		var result interface{}
		err := cli.CC.InstCli.GetOneByCondition(common.BKTableNameUserGroupPrivilege, []string{}, cond, &result)
		if nil != err {
			data := make(map[string]interface{})
			data[common.BKOwnerIDField] = ownerID
			data[common.BKUserGroupIDField] = groupID
			data[common.BKPrivilegeField] = common.KvMap{}
			blog.Error("get user group privi error :%v", err)
			return http.StatusNotFound, result, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, result, nil
	}, resp)
}
