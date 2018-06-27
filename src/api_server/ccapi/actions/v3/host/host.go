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

package host

import (
	"io"
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
)

var host = &hostAction{}

type hostAction struct {
	cc *api.APIResource
	base.BaseAction
}

// GetHosts get hosts by id without assocate object detaill
func (cli *hostAction) GetHosts(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/search"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)

}

// GetHostsAsstDetail get hosts by id with assocate object detaill
func (cli *hostAction) GetHostsAsstDetail(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/search/asstdetail"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)

}

// UpdateHosts batch to update host info
func (cli *hostAction) UpdateHosts(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/host/batch"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPUpdate)
	io.WriteString(resp, rsp)
}

// DeleteHosts batch to delete hosts
func (cli *hostAction) DeleteHosts(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/host/batch"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPDelete)
	io.WriteString(resp, rsp)
}

// GetHostDetail search host's detail infomation
func (cli *hostAction) GetHostDetail(req *restful.Request, resp *restful.Response) {
	pathParams := req.PathParameters()
	hostID := pathParams[common.BKHostIDField]
	ownerID := pathParams[common.BKOwnerIDField]
	url := cli.cc.HostAPI() + "/host/v1/hosts/" + ownerID + "/" + hostID
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)

}

// HostAddModulesRelation move host to another module
func (cli *hostAction) HostAddModulesRelation(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/modules"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// HostMoveToIDleModules move host to IDLE module
func (cli *hostAction) HostMoveToIDleModules(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/emptymodule"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// HostMoveToFaultModules move host to fault module
func (cli *hostAction) HostMoveToFaultModules(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/faultmodule"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// HostMoveToResoucePool move host to resource pool
func (cli *hostAction) HostMoveToResoucePool(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/resource"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

// AddHost add some new host into the resource
func (cli *hostAction) AddHost(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/addhost"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// AssginHostToApp assign some host to a application
func (cli *hostAction) AssginHostToApp(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/assgin"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

// Snapshot search the host snap
func (cli *hostAction) Snapshot(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/host/snapshot/" + req.PathParameter(common.BKHostIDField)
	blog.Debug("request url %s", url)
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	io.WriteString(resp, rsp)
}

func (cli *hostAction) addHostFromAgent(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/host/add/agent"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)
}

func (cli *hostAction) addHostModuleMutiple(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/modules/biz/mutiple"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPCreate)
	io.WriteString(resp, rsp)
}

func (cli *hostAction) setHostMoveToIdleModules(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/faultmodule/set"
	defErr := cli.cc.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	rsp, err := httpcli.ReqForward(req, url, common.HTTPCreate)
	if nil != err {
		blog.Errorf("newHostSyncAppTopo  http do err, url:%s, err:%s", url, err.Error())
		cli.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommHTTPDoRequestFailed, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	io.WriteString(resp, rsp)
}

func (cli *hostAction) cloneHostProperty(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/propery/clone"
	defErr := cli.cc.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	rsp, err := httpcli.ReqForward(req, url, common.HTTPUpdate)
	if nil != err {
		blog.Errorf("cloneHostProperty  http do err, url:%s, err:%s", url, err.Error())
		cli.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommHTTPDoRequestFailed, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	io.WriteString(resp, rsp)
}

func (cli *hostAction) newHostSyncAppTopo(req *restful.Request, resp *restful.Response) {
	url := cli.cc.HostAPI() + "/host/v1/hosts/sync/new/host"
	defErr := cli.cc.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	rsp, err := httpcli.ReqForward(req, url, common.HTTPCreate)
	if nil != err {
		blog.Errorf("newHostSyncAppTopo  http do err, url:%s, err:%s", url, err.Error())
		cli.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommHTTPDoRequestFailed, defErr.Errorf(common.CCErrCommHTTPDoRequestFailed).Error(), resp)
		return
	}
	io.WriteString(resp, rsp)
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/search", Params: nil, Handler: host.GetHosts, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/search/asstdetail", Params: nil, Handler: host.GetHostsAsstDetail, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/hosts/batch", Params: nil, Handler: host.DeleteHosts, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/batch", Params: nil, Handler: host.UpdateHosts, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/hosts/{bk_supplier_account}/{bk_host_id}", Params: nil, Handler: host.GetHostDetail, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/modules", Params: nil, Handler: host.HostAddModulesRelation, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/modules/idle", Params: nil, Handler: host.HostMoveToIDleModules, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/modules/fault", Params: nil, Handler: host.HostMoveToFaultModules, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/modules/resource", Params: nil, Handler: host.HostMoveToResoucePool, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/add", Params: nil, Handler: host.AddHost, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/hosts/modules/resource/idle", Params: nil, Handler: host.AssginHostToApp, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/hosts/snapshot/{bk_host_id}", Params: nil, Handler: host.Snapshot, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/host/add/agent", Params: nil, Handler: host.addHostFromAgent, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/modules/biz/mutilple", Params: nil, Handler: host.addHostModuleMutiple, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/modules/idle/set", Params: nil, Handler: host.setHostMoveToIdleModules, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/propery/clone", Params: nil, Handler: host.cloneHostProperty, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/sync/new/host", Params: nil, Handler: host.newHostSyncAppTopo, FilterHandler: nil, Version: v3.APIVersion})
	host.cc = api.NewAPIResource()
}
