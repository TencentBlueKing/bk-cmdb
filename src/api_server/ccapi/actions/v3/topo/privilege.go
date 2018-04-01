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

var pri = &priAction{}

type priAction struct {
	base.BaseAction
}

// CreateRolePri create process
func (cli *priAction) CreateRolePri(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	propertyID := pathParams["property_id"]
	objID := pathParams["obj_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/" + ownerID + "/" + objID + "/" + propertyID
	blog.Info("createRolePri url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// GetRolePri get process
func (cli *priAction) GetRolePri(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	propertyID := pathParams["property_id"]
	objID := pathParams["obj_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/" + ownerID + "/" + objID + "/" + propertyID
	blog.Info("GetRolePri url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

// CreateUserGroup create user group
func (cli *priAction) CreateUserGroup(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/group/" + ownerID
	blog.Info("CreateUserGroup url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// UpdateUserGroup update user group
func (cli *priAction) UpdateUserGroup(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	groupID := pathParams["group_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/group/" + ownerID + "/" + groupID
	blog.Info("UpdateUserGroup url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// DeleteUserGroup delete user group
func (cli *priAction) DeleteUserGroup(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	groupID := pathParams["group_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/group/" + ownerID + "/" + groupID
	blog.Info("UpdateUserGroup url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}

// SearchUserGroup search user group
func (cli *priAction) SearchUserGroup(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/group/" + ownerID + "/search"
	blog.Info("SearchUserGroup url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// GetUserGroupPri get user group privilege
func (cli *priAction) GetUserGroupPri(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	groupID := pathParams["group_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/group/detail/" + ownerID + "/" + groupID
	blog.Info("SearchUserGroup url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

// UpdateUserGroupPri update user group privilege
func (cli *priAction) UpdateUserGroupPri(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	groupID := pathParams["group_id"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/group/detail/" + ownerID + "/" + groupID
	blog.Info("UpdateUserGroupPri url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// GetUserPri update user group privilege
func (cli *priAction) GetUserPri(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	userName := pathParams["user_name"]
	url := cli.CC.TopoAPI() + "/topo/v1/privilege/user/detail/" + ownerID + "/" + userName
	blog.Info("GetUserPri url:%s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/topo/privilege/{owner_id}/{obj_id}/{property_id}", Params: nil, Handler: pri.CreateRolePri, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/privilege/{owner_id}/{obj_id}/{property_id}", Params: nil, Handler: pri.GetRolePri, Version: v3.APIVersion})
	//user group action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/topo/privilege/group/{owner_id}", Params: nil, Handler: pri.CreateUserGroup, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/topo/privilege/group/{owner_id}/{group_id}", Params: nil, Handler: pri.UpdateUserGroup, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/topo/privilege/group/{owner_id}/search", Params: nil, Handler: pri.SearchUserGroup, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/topo/privilege/group/{owner_id}/{group_id}", Params: nil, Handler: pri.DeleteUserGroup, Version: v3.APIVersion})
	//user group privilege
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/privilege/group/detail/{owner_id}/{group_id}", Params: nil, Handler: pri.GetUserGroupPri, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/topo/privilege/group/detail/{owner_id}/{group_id}", Params: nil, Handler: pri.UpdateUserGroupPri, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/topo/privilege/user/detail/{owner_id}/{user_name}", Params: nil, Handler: pri.GetUserPri, Version: v3.APIVersion})
	// set cc api interface
	pri.CreateAction()
}
