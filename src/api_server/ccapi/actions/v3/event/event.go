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
 
package process

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/core/cc/actions"

	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	httpcli "configcenter/src/common/http/httpclient"
	"io"

	"github.com/emicklei/go-restful"
)

var event *procAction = &procAction{}

type procAction struct {
	base.BaseAction
}

//create process
func (cli *procAction) Subscribe(req *restful.Request, resp *restful.Response) {
	blog.Info("create subscribe")
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]

	url := cli.CC.EventAPI() + "/event/v1/subscribe/" + ownerID + "/" + appID
	blog.Info("create subscribe url:%s", req.Request.URL.Path)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

//delete process
func (cli *procAction) UnSubscribe(req *restful.Request, resp *restful.Response) {
	blog.Info("delete subscribe")
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	subscribeID := pathParams["subscribe_id"]
	url := cli.CC.EventAPI() + "/event/v1/subscribe/" + ownerID + "/" + appID + "/" + subscribeID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}

//update process
func (cli *procAction) Rebook(req *restful.Request, resp *restful.Response) {
	blog.Info("update subscribe")
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	subscribeID := pathParams["subscribe_id"]
	url := cli.CC.EventAPI() + "/event/v1/subscribe/" + ownerID + "/" + appID + "/" + subscribeID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

//search process
func (cli *procAction) Query(req *restful.Request, resp *restful.Response) {
	blog.Info("search subscribe")
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.EventAPI() + "/event/v1/subscribe/search/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

//search process
func (cli *procAction) Ping(req *restful.Request, resp *restful.Response) {
	blog.Info("ping callback")
	url := cli.CC.EventAPI() + "/event/v1/subscribe/ping"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

//search process
func (cli *procAction) Telnet(req *restful.Request, resp *restful.Response) {
	blog.Info("telnet callback")
	url := cli.CC.EventAPI() + "/event/v1/subscribe/telnet"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/event/subscribe/ping", Params: nil, Handler: event.Ping, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/event/subscribe/telnet", Params: nil, Handler: event.Telnet, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/event/subscribe/search/{owner_id}/{app_id}", Params: nil, Handler: event.Query, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/event/subscribe/{owner_id}/{app_id}", Params: nil, Handler: event.Subscribe, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/event/subscribe/{owner_id}/{app_id}/{subscribe_id}", Params: nil, Handler: event.UnSubscribe, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/event/subscribe/{owner_id}/{app_id}/{subscribe_id}", Params: nil, Handler: event.Rebook, FilterHandler: nil, Version: v3.APIVersion})
	// set cc api interface
	event.CreateAction()
}
