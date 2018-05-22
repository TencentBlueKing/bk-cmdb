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

package metadata

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/api/metadata"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
)

var graphics = &graphicsAction{}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/topographics/search", Params: nil, Handler: graphics.SearchTopoGraphics})

	// set cc api resource
	objaction.CC = api.NewAPIResource()
}

// ObjectAction
type graphicsAction struct {
	base.BaseAction
}

// CreateClassification create object's classification
func (cli *graphicsAction) SearchTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("create classification")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		selector := map[string]interface{}{}
		if jsErr := json.Unmarshal([]byte(value), selector); nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", string(value), jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		results := []metadata.TopoGraphics{}
		if selErr := cli.CC.InstCli.GetMutilByCondition(metadata.ObjectDes{}.TableName(), nil, selector, &results, "", -1, -1); nil != selErr {
			blog.Error("select data failed, error information is %s", selErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrObjectDBOpErrno)
		}

		return http.StatusOK, results, nil
	}, resp)
}
