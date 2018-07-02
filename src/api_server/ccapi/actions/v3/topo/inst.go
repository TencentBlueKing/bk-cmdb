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

var inst = &instAction{}

type instAction struct {
	base.BaseAction
	//sencecli *api.Client
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/inst/{owner_id}/{obj_id}", Params: nil, Handler: inst.CreateInst, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", Params: nil, Handler: inst.DeleteInst, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/inst/{owner_id}/{obj_id}/batch", Params: nil, Handler: inst.DeleteInstBatch, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/inst/{owner_id}/{obj_id}/{inst_id}", Params: nil, Handler: inst.UpdateInst, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/inst/{owner_id}/{obj_id}/batch", Params: nil, Handler: inst.UpdateInstBatch, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/{owner_id}/{obj_id}", Params: nil, Handler: inst.SelectInsts, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}/detail", Params: nil, Handler: inst.SelectInstsAndAsstDetail, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/association/search/owner/{owner_id}/object/{obj_id}", Params: nil, Handler: inst.SelectInstsByAssociation, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/owner/{owner_id}/object/{obj_id}", Params: nil, Handler: inst.SelectInstsByObject, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/{owner_id}/{obj_id}/{inst_id}", Params: nil, Handler: inst.SelectInst, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/search/topo/owner/{owner_id}/object/{object_id}/inst/{inst_id}", Params: nil, Handler: inst.SelectTopo, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/inst/association/topo/search/owner/{owner_id}/object/{object_id}/inst/{inst_id}", Params: nil, Handler: inst.SelectInstAssociationTopo, FilterHandler: nil, Version: v3.APIVersion})
	// set cc api interface
	inst.CreateAction()
	//inst.sencecli = api.NewClient(inst.CC.TopoAPI())
}

// CreateInst create inst object
func (cli *instAction) CreateInst(req *restful.Request, resp *restful.Response) {

	blog.Info("create inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardCreateMetaInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID),
		resp)

}

// DeleteInst delete some inst object
func (cli *instAction) DeleteInstBatch(req *restful.Request, resp *restful.Response) {

	blog.Info("delete inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID, "-1"), // TODO:  -1 means batch operation, need to do implement a new api
		resp)

}

// DeleteInst delete some inst object
func (cli *instAction) DeleteInst(req *restful.Request, resp *restful.Response) {

	blog.Info("delete inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")
	instID := req.PathParameter("inst_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID, instID),
		resp)

}

// UpdateInstBatch update some inst object info
func (cli *instAction) UpdateInstBatch(req *restful.Request, resp *restful.Response) {

	blog.Info("update inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID, "-1"), // TODO: need to implement a new api, -1 means to batch operation
		resp)

}

// UpdateInst update some inst object info
func (cli *instAction) UpdateInst(req *restful.Request, resp *restful.Response) {

	blog.Info("update inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")
	instID := req.PathParameter("inst_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(
		senceCLI.ReForwardUpdateMetaInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID, instID),
		resp)

}

// SelectTopo search inst topo
func (cli *instAction) SelectTopo(req *restful.Request, resp *restful.Response) {

	blog.Info("select inst topo")

	ownerID := req.PathParameter("owner_id")
	objectID := req.PathParameter("object_id")
	instID := req.PathParameter("inst_id")

	sencecli := api.NewClient(inst.CC.TopoAPI())
	//sencecli.SetAddress(cli.CC.TopoAPI)
	cli.CallResponse(
		sencecli.ReForwardSelectMetaInstsTopo(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objectID, instID),
		resp)
}

// SelectTopo search inst topo
func (cli *instAction) SelectInstAssociationTopo(req *restful.Request, resp *restful.Response) {

	blog.Info("select inst topo")

	ownerID := req.PathParameter("owner_id")
	objectID := req.PathParameter("object_id")
	instID := req.PathParameter("inst_id")

	sencecli := api.NewClient(inst.CC.TopoAPI())
	//sencecli.SetAddress(cli.CC.TopoAPI)
	cli.CallResponse(
		sencecli.ReForwardSelectInstAssociationTopo(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objectID, instID),
		resp)
}

// SelectInst search some insts
func (cli *instAction) SelectInst(req *restful.Request, resp *restful.Response) {

	blog.Info("select inst")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")
	instID := req.PathParameter("inst_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())

	cli.CallResponse(
		senceCLI.ReForwardSelectMetaInst(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID, instID),
		resp)

}

// SelectInsts search instss
func (cli *instAction) SelectInsts(req *restful.Request, resp *restful.Response) {

	blog.Info("select insts")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(

		senceCLI.ReForwardSelectMetaInsts(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID),
		resp)

}

// SelectInsts search instss
func (cli *instAction) SelectInstsAndAsstDetail(req *restful.Request, resp *restful.Response) {

	blog.Info("select insts")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(

		senceCLI.ReForwardSelectMetaInstsAndAsstDetail(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID),
		resp)

}

// SelectInstsByAssociation search instss
func (cli *instAction) SelectInstsByAssociation(req *restful.Request, resp *restful.Response) {

	blog.Info("select association insts")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(
		senceCLI.ReForwardSelectInstByAssociation(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID),
		resp)

}

func (cli *instAction) SelectInstsByObject(req *restful.Request, resp *restful.Response) {
	blog.Info("select insts by object")

	ownerID := req.PathParameter("owner_id")
	objID := req.PathParameter("obj_id")

	senceCLI := api.NewClient(inst.CC.TopoAPI())
	senceCLI.SetAddress(cli.CC.TopoAPI()) // TODO: need to be removed

	cli.CallResponse(
		senceCLI.ReForwardSelectInstByObject(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, ownerID, objID),
		resp)

}
