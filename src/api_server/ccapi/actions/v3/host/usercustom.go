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
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"io"

	restful "github.com/emicklei/go-restful"
)

// SaveUserCustom save user configure
func SaveUserCustom(req *restful.Request, resp *restful.Response) {
	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/usercustom"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// GetUser get useer configure
func GetUser(req *restful.Request, resp *restful.Response) {
	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/usercustom/user/search"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// GetDefault get user default configure
func GetDefault(req *restful.Request, resp *restful.Response) {
	cc := api.NewAPIResource()
	url := cc.HostAPI() + "/host/v1/usercustom/default/search"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/usercustom", Params: nil, Handler: SaveUserCustom, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/user/search", Params: nil, Handler: GetUser, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/usercustom/default/search", Params: nil, Handler: GetDefault, Version: v3.APIVersion})

}
