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
	"github.com/rs/xid"
)

var group = &groupAction{}

//group Action
type groupAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/privilege/group/{bk_supplier_account}", Params: nil, Handler: group.CreateUserGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/privilege/group/{bk_supplier_account}/{group_id}", Params: nil, Handler: group.UpdateUserGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/privilege/group/{bk_supplier_account}/{group_id}", Params: nil, Handler: group.DeleteUserGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/privilege/group/{bk_supplier_account}/search", Params: nil, Handler: group.SearchUserGroup})

	// set cc api interface
	group.CreateAction()
}

//CreateUserGroup create group
func (cli *groupAction) CreateUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("create user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		data, err := js.Map()
		if err != nil {
			blog.Error("create user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		data[common.BKOwnerIDField] = ownerID
		guid := xid.New()
		data[common.BKUserGroupIDField] = guid.String()
		data = util.SetModOwner(data, ownerID)
		_, err = cli.CC.InstCli.Insert(common.BKTableNameUserGroup, data)
		if nil != err {
			blog.Error("create user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//UpdateUserGroup create group
func (cli *groupAction) UpdateUserGroup(req *restful.Request, resp *restful.Response) {

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
			blog.Error("update user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		data, err := js.Map()
		if err != nil {
			blog.Error("update user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		cond := make(map[string]interface{})
		cond[common.BKOwnerIDField] = ownerID
		cond[common.BKUserGroupIDField] = groupID
		cond = util.SetModOwner(cond, ownerID)
		data = util.SetModOwner(data, ownerID)
		err = cli.CC.InstCli.UpdateByCondition(common.BKTableNameUserGroup, data, cond)
		if nil != err {
			blog.Error("update user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//DeleteUserGroup create group
func (cli *groupAction) DeleteUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
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
		err := cli.CC.InstCli.DelByCondition(common.BKTableNameUserGroup, cond)
		if nil != err {
			blog.Error("delete user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, nil, nil
	}, resp)
}

//SearchUserGroup create group
func (cli *groupAction) SearchUserGroup(req *restful.Request, resp *restful.Response) {

	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	ownerID := util.GetActionOnwerID(req)

	cli.CallResponseEx(func() (int, interface{}, error) {

		defer req.Request.Body.Close()
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read json data error :%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			blog.Error("get user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		cond, err := js.Map()
		if err != nil {
			blog.Error("get user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		cond[common.BKOwnerIDField] = ownerID
		cond = util.SetModOwner(cond, ownerID)
		var result []interface{}
		err = cli.CC.InstCli.GetMutilByCondition(common.BKTableNameUserGroup, []string{}, cond, &result, "", 0, 0)
		if nil != err {
			blog.Error("get user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}
		return http.StatusOK, result, nil
	}, resp)
}
