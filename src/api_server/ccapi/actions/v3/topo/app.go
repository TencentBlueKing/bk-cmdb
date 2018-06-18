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

package topo

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"io"

	"github.com/emicklei/go-restful"
)

var app = &appAction{}

type appAction struct {
	base.BaseAction
}

// CreateApp create application
func (cli *appAction) CreateApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/" + ownerID
	//	req.Request.URL.Path = "/topo/v1/app/" + ownerID
	blog.Info("Create App url:%s", req.Request.URL.Path)
	//	httpcli.ProxyRestHttp(req, resp, url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// DeleteApp delete application
func (cli *appAction) DeleteApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/" + ownerID + "/" + appID
	//req.Request.URL.Path = "/topo/v1/app/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}

// UpdateApp update application
func (cli *appAction) UpdateApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// UpdateAppDataStatus update application data status
func (cli *appAction) UpdateAppDataStatus(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	flag := pathParams["flag"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/status/" + flag + "/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// SearchApp search application
func (cli *appAction) SearchApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/search/" + ownerID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// GetInternalTopo get internal topo
func (cli *appAction) GetInternalTopo(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/topo/internal/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

// CreateDefaultApp create default application
func (cli *appAction) CreateDefaultApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/default/" + ownerID
	//	req.Request.URL.Path = "/topo/v1/app/" + ownerID
	blog.Info("Create default App url:%s", req.Request.URL.Path)
	//	httpcli.ProxyRestHttp(req, resp, url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// SearchDefaultApp search default application
func (cli *appAction) SearchDefaultApp(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/app/default/" + ownerID + "/search"
	//	req.Request.URL.Path = "/topo/v1/app/" + ownerID
	blog.Info("search default App url:%s", req.Request.URL.Path)
	//	httpcli.ProxyRestHttp(req, resp, url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/biz/default/{owner_id}", Params: nil, Handler: app.CreateDefaultApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/biz/default/{owner_id}/search", Params: nil, Handler: app.SearchDefaultApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/biz/{owner_id}", Params: nil, Handler: app.CreateApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/biz/{owner_id}/{app_id}", Params: nil, Handler: app.DeleteApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/biz/{owner_id}/{app_id}", Params: nil, Handler: app.UpdateApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/biz/search/{owner_id}", Params: nil, Handler: app.SearchApp, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/biz/status/{flag}/{owner_id}/{app_id}", Params: nil, Handler: app.UpdateAppDataStatus, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/internal/{owner_id}/{app_id}", Params: nil, Handler: app.GetInternalTopo, Version: v3.APIVersion})
	// set cc api interface
	app.CreateAction()
}
