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

	"fmt"
)

var topo = &topoAction{}

type topoAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/topo/model/mainline", Params: nil, Handler: topo.CreateTopoModel, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/topo/model/mainline/owners/{owner_id}/objectids/{obj_id}", Params: nil, Handler: topo.DeleteTopoModel, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/model/{owner_id}", Params: nil, Handler: topo.SelectTopoModel, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/model/{owner_id}/{cls_id}/{obj_id}", Params: nil, Handler: topo.SelectTopoModelByClsID, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/inst/{owner_id}/{app_id}", Params: nil, Handler: topo.SelectTopoInst, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/inst/child/{owner_id}/{obj_id}/{app_id}/{inst_id}", Params: nil, Handler: topo.SelectTopoInstChild, Version: v3.APIVersion})

	// set cc api interface
	topo.CreateAction()
}

// CreateTopoModel create a main line association map
func (cli *topoAction) CreateTopoModel(req *restful.Request, resp *restful.Response) {

	blog.Info("create topo model")

	senceCLI := api.NewClient(topo.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaTopo(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)

}

// DeleteTopoModel delete a main line association map
func (cli *topoAction) DeleteTopoModel(req *restful.Request, resp *restful.Response) {

	blog.Info("delete topo model")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(topo.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaTopo(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID),
		resp)
}

// SelectTopoModel search the main line object topo tree
func (cli *topoAction) SelectTopoModel(req *restful.Request, resp *restful.Response) {

	blog.Info("select topo model")

	ownerID := req.PathParameter("owner_id")

	senceCLI := api.NewClient(topo.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaTopo(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID),
		resp)

}

// SelectTopoModelByClsID search the object topo tree by the classification
func (cli *topoAction) SelectTopoModelByClsID(req *restful.Request, resp *restful.Response) {

	blog.Info("select topo model by clsid")

	ownerID := req.PathParameter("owner_id")
	clsID := req.PathParameter("cls_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(topo.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaTopoByClsID(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, clsID, objID),
		resp)

}

// SelectTopoInst search the inst topo tree
func (cli *topoAction) SelectTopoInst(req *restful.Request, resp *restful.Response) {

	blog.Info("select topo inst ")

	ownerID := req.PathParameter("owner_id")
	appID := req.PathParameter("app_id")

	senceCLI := api.NewClient(topo.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardSelectMetaTopoInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, fmt.Sprintf("%s?level=%s", url, req.QueryParameter("level")), method)
		}, ownerID, appID),
		resp)
}

// SelectTopoInstChild search the children inst of the inst
func (cli *topoAction) SelectTopoInstChild(req *restful.Request, resp *restful.Response) {

	blog.Info("select topo child inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")
	appID := req.PathParameter("app_id")
	instID := req.PathParameter("inst_id")

	senceCLI := api.NewClient(topo.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaTopoInstChild(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID, appID, instID),
		resp)

}
