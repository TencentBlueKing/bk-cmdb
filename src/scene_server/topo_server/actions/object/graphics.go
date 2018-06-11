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
	"configcenter/src/source_controller/api/metadata"
	api "configcenter/src/source_controller/api/object"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"github.com/emicklei/go-restful"
	"io/ioutil"
	"net/http"
	"strconv"
)

func init() {

	// register action
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/search", Params: nil, Handler: obj.SelectObjectTopoGraphics})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/objects/topographics/scope_type/{scope_type}/scope_id/{scope_id}/action/update", Params: nil, Handler: obj.UpdateObjectTopoGraphics})

}

func (cli *objectAction) SelectObjectTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("select object topo graphics")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	defLang := cli.CC.Lang.CreateDefaultCCLanguageIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		scopeType := req.PathParameter("scope_type")
		scopeID := req.PathParameter("scope_id")

		graphcondition := api.TopoGraphics{}
		graphcondition.SetScopeType(scopeType)
		graphcondition.SetScopeID(scopeID)
		dbnodes, err := cli.mgr.SearchGraphics(forward, &graphcondition, defErr)
		if err != nil {
			blog.Errorf("SearchGraphics failed %v", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoGraphicsUpdateFailed)
		}

		graphnodes := map[string]*api.TopoGraphics{}
		for index, node := range dbnodes {
			graphnodes[*node.NodeType+*node.ObjID+strconv.Itoa(*node.InstID)] = &dbnodes[index]
		}

		nodes := []api.TopoGraphics{}
		if scopeType == "global" {
			objs, err := cli.mgr.SelectObject(forward, []byte("{}"), defErr)
			if err != nil {
				blog.Errorf("SelectObject failed %v", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoGraphicsSearchFailed)
			}

			assts, err := cli.mgr.SelectObjectAsst(forward, map[string]interface{}{}, defErr)
			if err != nil {
				blog.Errorf("SelectObjectAsst failed %v", err.Error())
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoGraphicsSearchFailed)
			}

			objAssts := map[string][]api.ObjAsstDes{}
			for _, asst := range assts {
				objAssts[asst.ObjectID] = append(objAssts[asst.ObjectID], asst)
			}

			for _, obj := range objs {
				node := api.TopoGraphics{}
				node.SetNodeType("obj")
				node.SetObjID(obj.ObjectID)
				node.SetInstID(0)
				node.SetNodeName(obj.ObjectName)
				node.SetScopeType("global")
				node.SetScopeID("0")
				node.SetBizID(0)
				node.SetSupplierAccount("0")
				node.SetIsPre(obj.IsPre)
				node.SetIcon(obj.ObjIcon)
				commondata.TranslateObjectName(defLang, &obj.ObjectDes)

				oldnode := graphnodes[*node.NodeType+*node.ObjID+strconv.Itoa(*node.InstID)]
				if oldnode != nil {
					node.SetPosition(oldnode.Position)
					node.SetExt(oldnode.Ext)
				} else {
					node.SetPosition(&metadata.Position{})
					node.SetExt(map[string]interface{}{})
				}

				for _, asst := range objAssts[obj.ObjectID] {
					node.Assts = append(node.Assts, api.GraphAsst{
						AsstType: "",
						NodeType: "obj",
						ObjID:    asst.AsstObjID,
						InstID:   0,
						ObjAtt:   asst.ObjectAttID,
						Lable:    map[string]string{},
					})
				}
				nodes = append(nodes, node)
			}
		}

		return http.StatusOK, nodes, nil
	}, resp)
}

func (cli *objectAction) UpdateObjectTopoGraphics(req *restful.Request, resp *restful.Response) {

	blog.Info("update object topo graphics")

	// get the language
	language := util.GetActionLanguage(req)

	// get the error info by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	forward := &api.ForwardParam{Header: req.Request.Header}

	// execute
	cli.CallResponseEx(func() (int, interface{}, error) {
		scopeType := req.PathParameter("scope_type")
		scopeID := req.PathParameter("scope_id")

		// read body
		val, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("failed to read request body, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		datas := []api.TopoGraphics{}

		err = json.Unmarshal(val, &datas)
		if nil != err {
			blog.Error("unmarshal the json, error info is %s", err.Error())
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		for index := range datas {
			datas[index].SetScopeType(scopeType)
			datas[index].SetScopeID(scopeID)
		}

		err = cli.mgr.UpdateGraphics(forward, datas, defErr)
		if err != nil {
			blog.Errorf("UpdateGraphics failed %v", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoGraphicsUpdateFailed)
		}

		return http.StatusOK, nil, nil
	}, resp)
}
