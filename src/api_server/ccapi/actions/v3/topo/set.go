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

var set = &setAction{}

type setAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/set/{app_id}", Params: nil, Handler: set.CreateSet, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/set/{app_id}/{set_id}", Params: nil, Handler: set.DeleteSet, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/set/{app_id}/batch", Params: nil, Handler: set.DeleteSetBatch, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/set/{app_id}/{set_id}", Params: nil, Handler: set.UpdateSet, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/set/search/{owner_id}/{app_id}", Params: nil, Handler: set.SelectSet, Version: v3.APIVersion})

	// init
	set.CreateAction()
}

// CreateSet create a set
func (cli *setAction) CreateSet(req *restful.Request, resp *restful.Response) {

	blog.Info("create set")
	appID := req.PathParameter("app_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardCreateMetaSet(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID),
		resp)

}

// DeleteSet delete a set
func (cli *setAction) DeleteSetBatch(req *restful.Request, resp *restful.Response) {

	appID := req.PathParameter("app_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaSet(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID, "-1"), // TODO: -1 means to batch operation, need to implement a new api
		resp)
}

// DeleteSet delete a set
func (cli *setAction) DeleteSet(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object")

	appID := req.PathParameter("app_id")
	setID := req.PathParameter("set_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaSet(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID, setID),
		resp)
}

// UpdateSet update a set
func (cli *setAction) UpdateSet(req *restful.Request, resp *restful.Response) {

	blog.Info("update set")

	appID := req.PathParameter("app_id")
	setID := req.PathParameter("set_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaSet(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, appID, setID),
		resp)

}

// SelectSet search some sets
func (cli *setAction) SelectSet(req *restful.Request, resp *restful.Response) {

	blog.Info("select set")

	ownerID := req.PathParameter("owner_id")
	appID := req.PathParameter("app_id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaSet(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, appID),
		resp)

}
