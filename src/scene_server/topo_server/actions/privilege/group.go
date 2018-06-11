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
	_ "fmt"
	"io/ioutil"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var group = &groupAction{}

type groupAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/privilege/group/{bk_supplier_account}", Params: nil, Handler: group.CreateUserGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/privilege/group/{bk_supplier_account}/{group_id}", Params: nil, Handler: group.DeleteUserGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/privilege/group/{bk_supplier_account}/{group_id}", Params: nil, Handler: group.UpdateUserGroup})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/privilege/group/{bk_supplier_account}/search", Params: nil, Handler: group.SearchUserGroup})

	group.CreateAction()
}

//CreateUserGroup create user group
func (cli *groupAction) CreateUserGroup(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]

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
			cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		groupName, ok := data["group_name"]
		if !ok {
			blog.Error("group_name not found %v")
			return http.StatusBadRequest, nil, defErr.Errorf(common.CCErrCommParamsNeedSet, "group_name")
		}
		//get user group url
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID + "/" + "search"
		cond := make(map[string]interface{})
		cond["group_name"] = groupName
		blog.Info("get user group url: %s", groupURL)

		byteConds, err := json.Marshal(cond)
		if nil != err {
			blog.Errorf("json marshal error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONMarshalFailed)
		}
		blog.Info("get user group content: %s ", string(byteConds))
		groupInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPSelectPost, byteConds)
		if nil != err {
			blog.Error("get user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
		}
		blog.Info("get user group return: %s", groupInfo)
		var group params.SearchGroup
		err = json.Unmarshal([]byte(groupInfo), &group)
		if nil != err || !group.Result {
			blog.Error("create user group json Unmarshal error data:%s error:%v", groupInfo, err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupCreateFailed)
		}
		groupList, ok := group.Data.([]interface{})
		if ok && len(groupList) > 0 {
			blog.Error("create user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDuplicateItem)
		}
		cgroupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID
		blog.Info("create user group url: %v", cgroupURL)
		createInfo, err := httpcli.ReqHttp(req, cgroupURL, common.HTTPCreate, value)
		if nil != err {
			blog.Error("create user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupCreateFailed)
		}
		return http.StatusOK, createInfo, nil
	}, resp)
}

//DeleteUserGroup delete user group
func (cli *groupAction) DeleteUserGroup(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		groupID, _ := pathParams["group_id"]

		//delete user group url
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID + "/" + groupID

		blog.Info("delete privilege url: %s", groupURL)
		deleteInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPDelete, nil)
		if nil != err {
			blog.Error("create user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupDeleteFailed)
		}
		return http.StatusOK, deleteInfo, nil
	}, resp)
}

//UpdateUserGroup delete user group
func (cli *groupAction) UpdateUserGroup(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]
		groupID, _ := pathParams["group_id"]

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
			cli.ResponseFailed(common.CC_Err_Comm_http_Input_Params, common.CC_Err_Comm_http_Input_Params_STR, resp)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		groupName, ok := data["group_name"]
		if ok { //has group name name not duplicate data
			//get user group url
			groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID + "/" + "search"
			cond := make(map[string]interface{})
			cond["group_name"] = groupName
			cond[common.BKOwnerIDField] = ownerID
			cond["group_id"] = map[string]interface{}{common.BKDBNE: groupID}
			blog.Info("get user group url: %s", groupURL)

			byteConds, innerErr := json.Marshal(cond)
			if nil != innerErr {
				blog.Errorf("json marshal error:%s", innerErr.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONMarshalFailed)
			}
			blog.Info("get user group content: %s ", string(byteConds))
			groupInfo, innerErr := httpcli.ReqHttp(req, groupURL, common.HTTPSelectPost, byteConds)
			if nil != innerErr {
				blog.Error("get user group error :%v", innerErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			blog.Info("get user group return: %s", groupInfo)
			var group params.SearchGroup
			innerErr = json.Unmarshal([]byte(groupInfo), &group)
			if nil != innerErr || !group.Result {
				blog.Error("create user group json Unmarshal error data:%s error:%v", groupInfo, innerErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupCreateFailed)
			}
			groupList, ok := group.Data.([]interface{})
			if ok && len(groupList) > 0 {
				blog.Error("create user group error :%v", innerErr)
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDuplicateItem)
			}

		}

		//update user group url
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID + "/" + groupID

		blog.Info("update privilege url: %s", groupURL)
		updateInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPUpdate, value)
		if nil != err {
			blog.Error("create user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupUpdateFailed)
		}
		return http.StatusOK, updateInfo, nil
	}, resp)
}

//SearchUserGroup search user group
func (cli *groupAction) SearchUserGroup(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		ownerID, _ := pathParams["bk_supplier_account"]

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
		data, err := js.Map()
		if err != nil {
			blog.Error("get user group failed, err msg : %v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		cond := make(map[string]interface{})
		for i, k := range data {
			c := make(map[string]interface{})
			c[common.BKDBLIKE] = k
			cond[i] = c
		}
		jsonStr, _ := json.Marshal(cond)
		//search user group url
		groupURL := cli.CC.ObjCtrl() + "/object/v1/privilege/group/" + ownerID + "/" + "search"

		blog.Info("search user group url: %s", groupURL)
		blog.Info("search user group info: %s", value)
		searchInfo, err := httpcli.ReqHttp(req, groupURL, common.HTTPSelectPost, jsonStr)
		blog.Info("search user group result: %s", searchInfo)
		if nil != err {
			blog.Error("search user group error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoUserGroupSelectFailed)
		}
		return http.StatusOK, searchInfo, nil
	}, resp)
}
