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
	"configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/api"

	"github.com/emicklei/go-restful"
)

var module = &moduleAction{}

type moduleAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/module/{app_id}/{set_id}", Params: nil, Handler: module.CreateModule, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/module/{app_id}/{set_id}/{module_id}", Params: nil, Handler: module.DeleteModule, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/module/{app_id}/{set_id}/{module_id}", Params: nil, Handler: module.UpdateModule, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/module/search/{owner_id}/{app_id}/{set_id}", Params: nil, Handler: module.SelectModule, Version: v3.APIVersion})

	// set cc api interface
	module.CreateAction()

}

// CreateModule create a module of the set
func (cli *moduleAction) CreateModule(req *restful.Request, resp *restful.Response) {

	blog.Info("create object")

	appID := req.PathParameter("app_id")
	setID := req.PathParameter("set_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaModule(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID, setID),
		resp)

}

// DeleteModule delete the module
func (cli *moduleAction) DeleteModule(req *restful.Request, resp *restful.Response) {

	blog.Info("delete module")

	appID := req.PathParameter("app_id")
	setID := req.PathParameter("set_id")
	moduleID := req.PathParameter("module_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaModule(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID, setID, moduleID),
		resp)

}

// UpdateModule update the module information
func (cli *moduleAction) UpdateModule(req *restful.Request, resp *restful.Response) {

	blog.Info("update module")

	appID := req.PathParameter("app_id")
	setID := req.PathParameter("set_id")
	moduleID := req.PathParameter("module_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaModule(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID, setID, moduleID),
		resp)

}

// SelectModule search the module detail information
func (cli *moduleAction) SelectModule(req *restful.Request, resp *restful.Response) {

	blog.Info("select module ")
	ownerID := req.PathParameter("owner_id")
	appID := req.PathParameter("app_id")
	setID := req.PathParameter("set_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaModule(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, appID, setID),
		resp)
}
