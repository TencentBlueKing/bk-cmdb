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

var graphics = &graphicsAction{}

type graphicsAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/search", Params: nil, Handler: graphics.SelectObjectTopoGraphics, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/update", Params: nil, Handler: graphics.UpdateObjectTopoGraphics, Version: v3.APIVersion})

	// init
	graphics.CreateAction()
}

// CreateObject create some objects
func (cli *graphicsAction) SelectObjectTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("SelectObjectTopoGraphics")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	scopeType := req.PathParameter("scope_type")
	scopeID := req.PathParameter("scope_id")
	cli.CallResponse(
		senceCLI.ReForwardSelectObjectTopoGraphics(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, scopeType, scopeID),
		resp)

}

// CreateObjectBatch batch create some objects
func (cli *graphicsAction) UpdateObjectTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("UpdateObjectTopoGraphics ")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	scopeType := req.PathParameter("scope_type")
	scopeID := req.PathParameter("scope_id")
	cli.CallResponse(
		senceCLI.ReForwardUpdateObjectTopoGraphics(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, scopeType, scopeID),
		resp)

}
