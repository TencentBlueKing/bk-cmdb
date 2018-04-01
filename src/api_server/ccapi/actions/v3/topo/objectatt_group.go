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

	restful "github.com/emicklei/go-restful"
)

var objattgroup = &objectAttGroupAction{}

type objectAttGroupAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/objectatt/group/new", Params: nil, Handler: objattgroup.CreatePropertyGroup, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/objectatt/group/update", Params: nil, Handler: objattgroup.UpdatePropertyGroup, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/objectatt/group/groupid/{id}", Params: nil, Handler: objattgroup.DeletePropertyGroup, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/objectatt/group/property", Params: nil, Handler: objattgroup.UpdatePropertyGroupObjectAtt, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/objectatt/group/owner/{owner_id}/object/{object_id}/propertyids/{property_id}/groupids/{group_id}", Params: nil, Handler: objattgroup.DeletePropertyGroupObjectAtt, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/objectatt/group/property/owner/{owner_id}/object/{object_id}", Params: nil, Handler: objattgroup.SelectPropertyGroupByObjectID, Version: v3.APIVersion})

	// init
	objattgroup.CreateAction()
}

// CreatePropertyGroup to create property group
func (cli *objectAttGroupAction) CreatePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group, create")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaObjAttGroup(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)
}

// UpdatePropertyGroup to update property group
func (cli *objectAttGroupAction) UpdatePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group, update")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaObjAttGroup(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)
}

// DeletePropertyGroup to update property group
func (cli *objectAttGroupAction) DeletePropertyGroup(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group, create")

	senceCLI := api.NewClient(module.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaObjAttGroup(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, req.PathParameter("id")), resp)
}

// UpdatePropertyGroupObjectAtt to update property group object attribute
func (cli *objectAttGroupAction) UpdatePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group property, update")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaObjAttGroupProperty(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)
}

// DeletePropertyGroupObjectAtt to delete property group object attribute
func (cli *objectAttGroupAction) DeletePropertyGroupObjectAtt(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group property, delete")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaObjAttGroupProperty(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		},
			req.PathParameter("owner_id"),
			req.PathParameter("object_id"),
			req.PathParameter("property_id"),
			req.PathParameter("group_id"),
		), resp)
}

func (cli *objectAttGroupAction) SelectPropertyGroupByObjectID(req *restful.Request, resp *restful.Response) {
	blog.Info("object attribute group property, select")
	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectPropertyGroupByObjectID(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		},
			req.PathParameter("owner_id"),
			req.PathParameter("object_id"),
		), resp)
}
