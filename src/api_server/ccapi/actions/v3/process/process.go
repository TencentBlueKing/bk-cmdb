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

package process

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

var proc = &procAction{}

type procAction struct {
	base.BaseAction
}

//CreateProcess create process
func (cli *procAction) CreateProcess(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.ProcAPI() + "/process/v1/" + ownerID + "/" + appID
	blog.Info("Create App url:%s", req.Request.URL.Path)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// DeleteProcess delete process
func (cli *procAction) DeleteProcess(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	procID := pathParams["proc_id"]
	url := cli.CC.ProcAPI() + "/process/v1/" + ownerID + "/" + appID + "/" + procID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}

// UpdateProcess update process
func (cli *procAction) UpdateProcess(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	procID := pathParams["proc_id"]
	url := cli.CC.ProcAPI() + "/process/v1/" + ownerID + "/" + appID + "/" + procID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// BatchUpdateProcess batch update process
func (cli *procAction) BatchUpdateProcess(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.ProcAPI() + "/process/v1/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// SearchProcess search process
func (cli *procAction) SearchProcess(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	url := cli.CC.ProcAPI() + "/process/v1/search/" + ownerID + "/" + appID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// GetProcess get process
func (cli *procAction) GetProcess(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	procID := pathParams["proc_id"]
	url := cli.CC.ProcAPI() + "/process/v1/" + ownerID + "/" + appID + "/" + procID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

// BindProcModule bind proc module
func (cli *procAction) BindProcModule(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	procID := pathParams["proc_id"]
	moduleName := pathParams["module_name"]
	url := cli.CC.ProcAPI() + "/process/v1/module/" + ownerID + "/" + appID + "/" + procID + "/" + moduleName
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// GetProcBindModule get proc bind module
func (cli *procAction) GetProcBindModule(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	procID := pathParams["proc_id"]
	url := cli.CC.ProcAPI() + "/process/v1/module/" + ownerID + "/" + appID + "/" + procID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

// DeleteProcBindModule delete proc bind module
func (cli *procAction) DeleteProcBindModule(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	ownerID := pathParams["owner_id"]
	appID := pathParams["app_id"]
	procID := pathParams["proc_id"]
	moduleName := pathParams["module_name"]
	url := cli.CC.ProcAPI() + "/process/v1/module/" + ownerID + "/" + appID + "/" + procID + "/" + moduleName
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}
func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/proc/{owner_id}/{app_id}", Params: nil, Handler: proc.CreateProcess, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/proc/{owner_id}/{app_id}/{proc_id}", Params: nil, Handler: proc.DeleteProcess, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/proc/{owner_id}/{app_id}/{proc_id}", Params: nil, Handler: proc.UpdateProcess, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/proc/{owner_id}/{app_id}", Params: nil, Handler: proc.BatchUpdateProcess, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/proc/search/{owner_id}/{app_id}", Params: nil, Handler: proc.SearchProcess, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/proc/{owner_id}/{app_id}/{proc_id}", Params: nil, Handler: proc.GetProcess, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/proc/module/{owner_id}/{app_id}/{proc_id}/{module_name}", Params: nil, Handler: proc.BindProcModule, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/proc/module/{owner_id}/{app_id}/{proc_id}", Params: nil, Handler: proc.GetProcBindModule, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/proc/module/{owner_id}/{app_id}/{proc_id}/{module_name}", Params: nil, Handler: proc.DeleteProcBindModule, Version: v3.APIVersion})
	// set cc api interface
	proc.CreateAction()
}
