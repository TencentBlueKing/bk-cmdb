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

package search

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/search", Params: nil, Handler: host.HostSearch})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/search/asstdetail", Params: nil, Handler: host.HostSearchWithAsstDetail})
	// create CC object
	host.CreateAction()
}

// HostSearch search host by multi condition
func (cli *hostAction) HostSearch(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("get  http body value %v error:%s", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var data hostParse.HostCommonSearch
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%s", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		reply, err := logics.HostSearch(req, data, false, cli.CC.HostCtrl(), cli.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		return http.StatusOK, reply, nil
	}, resp)
}

// HostSearch search host by multi condition
func (cli *hostAction) HostSearchWithAsstDetail(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("get  http body value %v error:%s", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		var data hostParse.HostCommonSearch
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%s", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		reply, err := logics.HostSearch(req, data, true, cli.CC.HostCtrl(), cli.CC.ObjCtrl())
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostGetFail)
		}
		return http.StatusOK, reply, nil
	}, resp)
}
