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

package host

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"fmt"
	"io"

	restful "github.com/emicklei/go-restful"
)

var userAPI *userAPIAction = &userAPIAction{}

type userAPIAction struct {
	base.BaseAction
}

//AddUserConfig 新加userapi
func (u *userAPIAction) Add(req *restful.Request, resp *restful.Response) {

	url := userAPI.CC.HostAPI() + "/host/v1/userapi"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)

	io.WriteString(resp, rsp)
	return
}

//UpdateUserConfig 更新user api 内容
func (u *userAPIAction) Update(req *restful.Request, resp *restful.Response) {

	url := userAPI.CC.HostAPI() + fmt.Sprintf("/host/v1/userapi/%s/%s", req.PathParameter("app_id"), req.PathParameter("id"))
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)

	io.WriteString(resp, rsp)
	return
}

func (u *userAPIAction) Delete(req *restful.Request, resp *restful.Response) {

	url := userAPI.CC.HostAPI() + fmt.Sprintf("/host/v1/userapi/%s/%s", req.PathParameter("app_id"), req.PathParameter("id"))
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)

	io.WriteString(resp, rsp)
	return
}

func (u *userAPIAction) Get(req *restful.Request, resp *restful.Response) {

	url := userAPI.CC.HostAPI() + fmt.Sprintf("/host/v1/userapi/search/%s", req.PathParameter("app_id"))
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)

	io.WriteString(resp, rsp)
	return
}

func (u *userAPIAction) Detail(req *restful.Request, resp *restful.Response) {

	url := userAPI.CC.HostAPI() + fmt.Sprintf("/host/v1/userapi/detail/%s/%s", req.PathParameter("app_id"), req.PathParameter("id"))
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)

	io.WriteString(resp, rsp)
	return
}

//GetUserAPIData 获取API的数据
func (u *userAPIAction) GetUserAPIData(req *restful.Request, resp *restful.Response) {

	url := userAPI.CC.HostAPI() + fmt.Sprintf("/host/v1/userapi/data/%s/%s/%s/%s", req.PathParameter("app_id"), req.PathParameter("id"), req.PathParameter("skip"), req.PathParameter("limit"))
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)

	io.WriteString(resp, rsp)
	return
}

func init() {
	userAPI.CreateAction()

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/userapi", Params: nil, Handler: userAPI.Add, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/userapi/{app_id}/{id}", Params: nil, Handler: userAPI.Update, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/userapi/{app_id}/{id}", Params: nil, Handler: userAPI.Delete, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/userapi/search/{app_id}", Params: nil, Handler: userAPI.Get, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/userapi/detail/{app_id}/{id}", Params: nil, Handler: userAPI.Detail, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/userapi/data/{app_id}/{id}/{skip}/{limit}", Params: nil, Handler: userAPI.GetUserAPIData, Version: v3.APIVersion})

}
