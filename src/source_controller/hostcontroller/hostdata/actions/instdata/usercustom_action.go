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
	_ "configcenter/src/common/util"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

var (
	userCustomTableName    string = "cc_UserCustom"
	mgo_on_not_found_error string = "not found"
)

var userCostomAction = &userCustomAction{}

type userCustomAction struct {
	base.BaseAction
}

//保存用户自定义的配置
//表字段User用户名。OwnerID 供应商,IsDefault是否默认配置，HostQueryColumn,HostDisplayColumn主机的自定义配置，

//AddUserCustom add user custom config
func (cli *userCustomAction) AddUserCustom(req *restful.Request, resp *restful.Response) {
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
		data := make(map[string]interface{})
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value  %s error %v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		ID := xid.New()
		data["id"] = ID.String()
		data["bk_user"] = req.PathParameter("bk_user") //libraries.GetOperateUser(req)
		data = util.SetModOwner(data, ownerID)
		_, err = cc.InstCli.Insert(userCustomTableName, data)
		if nil != err {
			blog.Error("Create  user custom fail, error information is %s, params:%v", err.Error(), data)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCreateUserCustom)

		}
		return http.StatusOK, nil, nil
	}, resp)
}

//UpdateUserCustomByID update user custom
func (cli *userCustomAction) UpdateUserCustomByID(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()

		conditons := make(map[string]interface{})
		conditons["id"] = req.PathParameter("id")
		conditons["bk_user"] = req.PathParameter("bk_user") //libraries.GetOperateUser(req)
		conditons = util.SetModOwner(conditons, ownerID)

		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}
		data := make(map[string]interface{})
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value  %s error %v", value, err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)

		}
		err = cc.InstCli.UpdateByCondition(userCustomTableName, data, conditons)
		if nil != err {
			blog.Error("add  user custom fail, error information is %s, params:%v", err.Error(), data)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommDBUpdateFailed)
		}
		return http.StatusOK, nil, nil
	}, resp)

}

//GetUserCustomByUser get user custom
func (cli *userCustomAction) GetUserCustomByUser(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		cc := api.NewAPIResource()

		conds := make(map[string]interface{})
		user := req.PathParameter("bk_user")
		conds["bk_user"] = user
		conds = util.SetModOwner(conds, ownerID)

		result := make(map[string]interface{})
		err := cc.InstCli.GetOneByCondition(userCustomTableName, nil, conds, &result)
		if nil != err && mgo_on_not_found_error != err.Error() { //get one row from mgo, not found是未找到数据
			blog.Error("add  user custom fail, error information is %s, params:%v", err.Error(), conds)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		blog.InfoJSON("GetUserCustomByUser result %s", result)
		return http.StatusOK, result, nil
	}, resp)
}

//GetDefaultUserCustom get default user custom
func (cli *userCustomAction) GetDefaultUserCustom(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		cc := api.NewAPIResource()

		conds := make(map[string]interface{})
		conds["is_default"] = 1
		conds = util.SetModOwner(conds, ownerID)

		result := make(map[string]interface{})
		err := cc.InstCli.GetOneByCondition(userCustomTableName, nil, conds, &result)
		if nil != err {
			blog.Error("add  user custom fail, error information is %s, params:%v", err.Error(), conds)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}
		return http.StatusOK, result, nil
	}, resp)
}

func init() {
	userCostomAction.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/usercustom/{bk_user}", Params: nil, Handler: userCostomAction.AddUserCustom})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/usercustom/{bk_user}/{id}", Params: nil, Handler: userCostomAction.UpdateUserCustomByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/user/search/{bk_user}", Params: nil, Handler: userCostomAction.GetUserCustomByUser})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/default/search/{bk_user}", Params: nil, Handler: userCostomAction.GetDefaultUserCustom})

}
