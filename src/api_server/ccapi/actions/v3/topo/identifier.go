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

var identifier = &identifierAction{}

type identifierAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/identifier/{obj_type}/search", Params: nil, Handler: identifier.SearchIdentifier, Version: v3.APIVersion})

	// init
	identifier.CreateAction()
}

// CreateObject create some objects
func (cli *identifierAction) SearchIdentifier(req *restful.Request, resp *restful.Response) {

	blog.Info("SelectObjectTopoGraphics")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	objType := req.PathParameter("obj_type")
	cli.CallResponse(
		senceCLI.ReForwardSearchIdentifier(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, objType),
		resp)
}
