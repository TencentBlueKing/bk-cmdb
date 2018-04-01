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

var obj = &objectAction{}

type objectAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object/batch", Params: nil, Handler: obj.CreateObjectBatch, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/search/batch", Params: nil, Handler: obj.SelectObjectBatch, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object", Params: nil, Handler: obj.CreateObject, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/object/{id}", Params: nil, Handler: obj.DeleteObject, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/object/{id}", Params: nil, Handler: obj.UpdateObject, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objects", Params: nil, Handler: obj.SelectObjectWithParams, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objects/topo", Params: nil, Handler: obj.SelectObjectTopo, Version: v3.APIVersion})

	// init
	obj.CreateAction()
}

// CreateObject create some objects
func (cli *objectAction) CreateObject(req *restful.Request, resp *restful.Response) {

	blog.Info("create object")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaObject(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}),
		resp)

}

// CreateObjectBatch batch create some objects
func (cli *objectAction) CreateObjectBatch(req *restful.Request, resp *restful.Response) {

	blog.Info("batch create object ")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaObjectBatch(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}),
		resp)

}

// SelectObjectBatch batch select some objects
func (cli *objectAction) SelectObjectBatch(req *restful.Request, resp *restful.Response) {

	blog.Info("batch select object ")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardSelectMetaObjectBatch(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}),
		resp)

}

// DeleteObject delete some objects
func (cli *objectAction) DeleteObject(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object")

	id := req.PathParameter("id")
	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaObject(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, id),
		resp)

}

// UpdateObject update some object information
func (cli *objectAction) UpdateObject(req *restful.Request, resp *restful.Response) {

	blog.Info("update object")

	id := req.PathParameter("id")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaObject(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, id),
		resp)

}

// SelectObjectTopo search object topo tree
func (cli *objectAction) SelectObjectTopo(req *restful.Request, resp *restful.Response) {

	blog.Info("select object topo whith params")
	sencecli := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		sencecli.ReForwardSelectMetaObjectTopo(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}),
		resp)

}

// SelectObjectAttWithParams search object with params
func (cli *objectAction) SelectObjectWithParams(req *restful.Request, resp *restful.Response) {

	blog.Info("select object whith params")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaObject(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}),
		resp)

}
