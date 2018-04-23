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
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"fmt"
	"io"

	httpcli "configcenter/src/common/http/httpclient"

	restful "github.com/emicklei/go-restful"
)

var history *historyAction = &historyAction{}

type historyAction struct {
	base.BaseAction
}

func (cli *historyAction) AddHistory(req *restful.Request, resp *restful.Response) {
	url := history.CC.HostAPI() + "/host/v1/history"
	blog.Info("create add history  URL:%s", url)

	reply, err := httpcli.ReqForward(req, url, common.HTTPCreate)
	if nil != err {
		rsp, rsperr := cli.CC.CreateAPIRspStr(common.CCErrCommHTTPDoRequestFailed, err.Error())
		if nil != rsperr {
			blog.Error("create response failed, error information is %v", rsperr)
		} else {
			// TODO: 暂时不设置 resp.WriteHeader(httpcode)
			io.WriteString(resp, rsp)
		}
		return
	}

	resp.Write([]byte(reply))
}

func (cli *historyAction) GetHistorys(req *restful.Request, resp *restful.Response) {
	skip := req.PathParameter("skip")
	limit := req.PathParameter("limit")
	url := history.CC.HostAPI() + fmt.Sprintf("/host/v1/history/%s/%s", skip, limit)
	blog.Info("create get history  URL:%s", url)

	reply, err := httpcli.ReqForward(req, url, common.HTTPSelectGet)
	if nil != err {
		rsp, rsperr := cli.CC.CreateAPIRspStr(common.CCErrCommHTTPDoRequestFailed, err.Error())
		if nil != rsperr {
			blog.Error("create response failed, error information is %v", rsperr)
		} else {
			// TODO: 暂时不设置 resp.WriteHeader(httpcode)
			io.WriteString(resp, rsp)
		}
		return
	}
	resp.Write([]byte(reply))
}

func init() {
	history.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/hosts/history", Params: nil, Handler: history.AddHistory, FilterHandler: nil, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/hosts/history/{skip}/{limit}", Params: nil, Handler: history.GetHistorys, FilterHandler: nil, Version: v3.APIVersion})

}
