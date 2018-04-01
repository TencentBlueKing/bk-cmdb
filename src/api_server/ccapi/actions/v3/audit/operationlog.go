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
 
package controllers

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"io"

	"github.com/emicklei/go-restful"
)

var audit = &auditAction{}

type auditAction struct {
	cc *api.APIResource
}

// Search search application by id
func (cli *auditAction) Search(req *restful.Request, resp *restful.Response) {
	url := cli.cc.TopoAPI() + "/topo/v1/audit/search"
	rsp, _ := httpcli.ReqForward(req, url, common.HTTPSelectPost)
	io.WriteString(resp, rsp)

}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/audit/search", Params: nil, Handler: audit.Search, Version: v3.APIVersion})
	audit.cc = api.NewAPIResource()
}
