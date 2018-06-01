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
 
package hosts

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	restful "github.com/emicklei/go-restful"
)

var history *historyAction = &historyAction{}

type historyAction struct {
	base.BaseAction
}

func init() {
	history.CreateAction()

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/history", Params: nil, Handler: history.AddHistory})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/history/{start}/{limit}", Params: nil, Handler: history.GetHistorys})
}

// AddHistory  add host history
func (cli *historyAction) AddHistory(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, _ := ioutil.ReadAll(req.Request.Body)
		bodyData := make(map[string]interface{})
		err := json.Unmarshal([]byte(value), &bodyData)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		content, ok := bodyData["content"].(string)
		if !ok || "" == content {
			blog.Error("参数content不能为空")
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		params := make(map[string]interface{}, 1)
		params["content"] = content
		user := util.GetActionUser(req)

		URL := history.CC.HostCtrl() + "/host/v1/history/" + user
		isSuccess, errMsg, data := logics.GetHttpResult(req, URL, common.HTTPCreate, params)
		if !isSuccess {
			blog.Error("add host history   error, params:%v, error:%s", params, errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostHisCreateFail)

		}

		return http.StatusOK, data, nil
	}, resp)
}

// GetHistorys get host historys
func (cli *historyAction) GetHistorys(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		start := req.PathParameter("start")
		limit := req.PathParameter("limit")
		user := util.GetActionUser(req)

		url := history.CC.HostCtrl() + fmt.Sprintf("/host/v1/history/%s/%s/%s", user, start, limit)

		reply, err := httpcli.ReqForward(req, url, common.HTTPSelectGet)
		if nil != err {
			blog.Error("get host history   error, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostHisGetFail)

		}
		return http.StatusOK, reply, nil
	}, resp)
}
