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

package inst

import (
	"io/ioutil"
	"net/http"

	"configcenter/src/common/blog"

	"github.com/tidwall/gjson"

	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

var app = &identifierAction{}

type identifierAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/identifier/{obj_type}/search", Params: nil, Handler: app.SearchIdentifier})
	// create CC object
	app.CreateAction()
}

//search application
func (cli *identifierAction) SearchIdentifier(req *restful.Request, resp *restful.Response) {
	// get the language
	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		objtype := req.PathParameter("obj_type")
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
		blog.V(3).Infof("SearchIdentifier %s", value)

		//search identifier
		sAppURL := cli.CC.ObjCtrl() + "/object/v1/identifier/" + objtype + "/search"
		appInfo, err := httpcli.ReqHttp(req, sAppURL, common.HTTPSelectPost, value)
		result := gjson.Parse(appInfo)
		if !result.Get("resule").Bool() {
			return http.StatusInternalServerError, nil, defErr.Error(int(result.Get(common.HTTPBKAPIErrorCode).Int()))
		}
		resultData := result.Get("data").Value()
		return http.StatusOK, resultData, nil
	}, resp)
}
