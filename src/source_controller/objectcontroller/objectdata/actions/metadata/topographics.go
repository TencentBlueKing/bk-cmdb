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
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/topographics/update", Params: nil, Handler: graphics.UpdateTopoGraphics})

	// set cc api resource
	graphics.CC = api.NewAPIResource()
}

// ObjectAction
type graphicsAction struct {
	base.BaseAction
}

// CreateClassification create object's classification
func (cli *graphicsAction) SearchTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("SearchTopoGraphics")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		blog.Infof("search param %s", value)

		selector := metadata.TopoGraphics{}
		if jsErr := json.Unmarshal(value, &selector); nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", value, jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		results := []metadata.TopoGraphics{}
		if selErr := cli.CC.InstCli.GetMutilByCondition(metadata.TopoGraphics{}.TableName(), nil, util.SetModOwner(selector, ownerID), &results, "", -1, -1); nil != selErr {
			blog.Error("select data failed, error information is %s", selErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBSelectFailed)
		}

		return http.StatusOK, results, nil
	}, resp)
}

func (cli *graphicsAction) UpdateTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("UpdateTopoGraphics")

	// get the language
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("read http request body failed, error:%s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		datas := []metadata.TopoGraphics{}
		if jsErr := json.Unmarshal(value, &datas); nil != jsErr {
			blog.Error("failed to unmarshal the data, data is %s, error info is %s ", value, jsErr.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		for index := range datas {
			blog.InfoJSON("update graphic %s", datas[index])
			datas[index].OwnerID = ownerID
			_, err = cli.CC.InstCli.Insert(metadata.TopoGraphics{}.TableName(), datas[index].FillBlank())
			if cli.CC.InstCli.IsDuplicateErr(err) {
				condition := metadata.TopoGraphics{}
				condition.SetScopeType(*datas[index].ScopeType)
				condition.SetScopeID(*datas[index].ScopeID)
				condition.SetNodeType(*datas[index].NodeType)
				condition.SetObjID(*datas[index].ObjID)
				condition.SetInstID(*datas[index].InstID)
				condition.OwnerID = ownerID
				if err = cli.CC.InstCli.UpdateByCondition(metadata.TopoGraphics{}.TableName(), datas[index], condition); err != nil {
					blog.Error("update data failed, error information is %s", err.Error())
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBUpdateFailed)
				}
			} else if err != nil {
				blog.Error("insert data failed, error information is %s", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommDBInsertFailed)
			}
		}
		return http.StatusOK, nil, nil
	}, resp)
}
