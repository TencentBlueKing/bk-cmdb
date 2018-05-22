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

package object

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	api "configcenter/src/source_controller/api/object"
	"github.com/emicklei/go-restful"
	"net/http"
)

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/search", Params: nil, Handler: obj.SelectObjectTopoGraphics})

}

func (cli *objectAction) SelectObjectTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("select object topo")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		scopeType := req.PathParameter("scope_type")
		scopeID := req.PathParameter("scope_id")

		condition := map[string]interface{}{
			"scope_type": scopeType,
			"scope_id":   scopeID,
		}
		nodes, err := cli.mgr.SearchGraphics(forward, condition, defErr)
		if err != nil {
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoObjectAttributeCreateFailed)
		}

		if scopeType == "global" {

		}

		return 0, "", nil
	}, resp)
}
