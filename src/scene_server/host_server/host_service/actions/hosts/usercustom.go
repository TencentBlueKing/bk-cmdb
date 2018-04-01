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
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	restful "github.com/emicklei/go-restful"
)

var user *userAction = &userAction{}

type userAction struct {
	base.BaseAction
}

//表字段User用户名,OwnerID 供应商,IsDefault是否默认配置，HostQueryColumn,HostDisplayColumn主机的自定义配置，

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/usercustom", Params: nil, Handler: user.SaveUserCustom})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/user/search", Params: nil, Handler: user.GetUserCustom})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/default/search", Params: nil, Handler: user.GetDefaultCustom})
	user.CreateAction()
}

//SaveUserCustom save user custom
func (cli *userAction) SaveUserCustom(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		cc := api.NewAPIResource()

		value, err := ioutil.ReadAll(req.Request.Body)

		user := util.GetActionUser(req)

		params := make(map[string]interface{})
		err = json.Unmarshal([]byte(value), &params)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}

		URL := cc.HostCtrl() + "/host/v1/usercustom/user/search/" + user
		isSuccess, errMsg, data := logics.GetHttpResult(req, URL, common.HTTPSelectPost, nil)
		if !isSuccess {
			blog.Error("get user custom info eror, error:%s", errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostCustomCreateFail)

		}

		userCustom := data.(map[string]interface{})
		method := common.HTTPUpdate
		if 0 == len(userCustom) {
			URL = cc.HostCtrl() + "/host/v1/usercustom/" + user
			method = common.HTTPCreate
		} else {
			URL = cc.HostCtrl() + fmt.Sprintf("/host/v1/usercustom/%s/%s", user, userCustom["id"].(string))

		}
		isSuccess, errMsg, data = logics.GetHttpResult(req, URL, method, params)
		if !isSuccess {
			blog.Error("get user custom info eror, error:%s", errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostCustomCreateFail)

		}
		return http.StatusOK, "", nil

	}, resp)

}

//GetUserCustom   get user custom config
func (cli *userAction) GetUserCustom(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		cc := api.NewAPIResource()
		user := util.GetActionUser(req)

		url := cc.HostCtrl() + "/host/v1/usercustom/user/search/" + user
		isSuccess, errMsg, data := logics.GetHttpResult(req, url, common.HTTPSelectPost, nil)
		if !isSuccess {
			blog.Error("get user custom info eror, error:%s", errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostCustomGetFail)

		}

		return http.StatusOK, data, nil

	}, resp)

}

//GetDefaultCustom get default custom
func (cli *userAction) GetDefaultCustom(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		cc := api.NewAPIResource()
		user := util.GetActionUser(req)
		url := cc.HostCtrl() + "/host/v1/usercustom/default/search/" + user
		isSuccess, errMsg, data := logics.GetHttpResult(req, url, common.HTTPSelectPost, nil)
		if !isSuccess {
			blog.Error("get default user custom info eror, error:%s", errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostCustomGetDefaultFail)

		}
		return http.StatusOK, data, nil

	}, resp)

}
