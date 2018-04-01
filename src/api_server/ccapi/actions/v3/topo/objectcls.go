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

var objcls = &objClsAction{}

type objClsAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object/classification", Params: nil, Handler: objcls.CreateClassification, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/object/classification/{id}", Params: nil, Handler: objcls.DeleteClassification, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/object/classification/{id}", Params: nil, Handler: objcls.UpdateClassification, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/classifications", Params: nil, Handler: objcls.SelectClassification, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/classification/{owner_id}/objects", Params: nil, Handler: objcls.SelectClassificationWithObjects, Version: v3.APIVersion})

	// init
	objcls.CreateAction()
}

// CreateClassification create some object's classification
func (cli *objClsAction) CreateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("create object classification")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardCreateMetaObjCls(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)
}

// DeleteClassification delete some object's classification
func (cli *objClsAction) DeleteClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("delete object classification")

	id := req.PathParameter("id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaObjCls(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, id),
		resp)

}

// UpdateObjectClassification update some object's classification
func (cli *objClsAction) UpdateClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("update object classification")

	id := req.PathParameter("id")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaObjCls(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, id),
		resp)

}

// SelectClassification select all classifications
func (cli *objClsAction) SelectClassification(req *restful.Request, resp *restful.Response) {

	blog.Info("select object classification")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaObjCls(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)

}

// SelectClassificationWithObjects search all classification include all object belong to it
func (cli *objClsAction) SelectClassificationWithObjects(req *restful.Request, resp *restful.Response) {

	blog.Info("select object classification")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	ownerID := req.PathParameter("owner_id")
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaObjClsWithObjects(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID), resp)

}
