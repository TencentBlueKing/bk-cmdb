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

package topo

import (
	"configcenter/src/api_server/ccapi/actions/v3"
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/api"
	"net/http"
	"github.com/emicklei/go-restful"
	"fmt"
	"encoding/json"
	"io"
)

var objatt = &objectAttAction{}

type objectAttAction struct {
	base.BaseAction
}

func init() {

	// register actions
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/object/attr", Params: nil, Handler: objatt.CreateObjectAtt, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/object/attr/{attr_id}", Params: nil, Handler: objatt.DeleteObjectAtt, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/object/attr/{attr_id}", Params: nil, Handler: objatt.UpdateObjectAtt, Version: v3.APIVersion})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/object/attr/search", Params: nil, Handler: objatt.SelectObjectAttWithParams, Version: v3.APIVersion})

	// init
	objatt.CreateAction()
}

// AddObjectAtt create some object's attributes
func (cli *objectAttAction) CreateObjectAtt(req *restful.Request, resp *restful.Response) {
	blog.Info("create objectatt")

	url := fmt.Sprintf("%s/topo/v1/objectattr", module.CC.TopoAPI())
	res ,err := httpclient.ReqForward(req, url, common.HTTPCreate)
	if nil != err {
		blog.Error("httpclient ReqForward err :%v", err)
		cli.ResponseFailed(http.StatusInternalServerError, err.Error(), resp)
		return
	}
	reply := make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &reply)
	if nil != err {
		blog.Error("json unmarshal err :%v", err)
		cli.ResponseFailed(http.StatusInternalServerError, fmt.Sprintf("json unmarshal err %v", err.Error()), resp)
		return
	}
	io.WriteString(resp,res)
	return
	//senceCLI := api.NewClient(module.CC.TopoAPI())
	//cli.CallResponse(
	//	senceCLI.ReForwardCreateMetaObjAtt(func(url, method string) (string, error) {
	//		return httpclient.ReqForward(req, url, method)
	//	}), resp)
}

// DeleteObjectAtt delete some object's attributes
func (cli *objectAttAction) DeleteObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("delete objectatt")

	attrID := req.PathParameter("attr_id")
	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardDeleteMetaObjAtt(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}, attrID),
		resp)

}

// UpdateObjectAtt update some object's attributes
func (cli *objectAttAction) UpdateObjectAtt(req *restful.Request, resp *restful.Response) {

	blog.Info("update objectatt")

	attrID := req.PathParameter("attr_id")
	url := fmt.Sprintf("%s/topo/v1/objectattr/%s", module.CC.TopoAPI(), attrID)
	res ,err := httpclient.ReqForward(req, url, common.HTTPUpdate)
	if nil != err {
		blog.Error("httpclient ReqForward err :%v", err)
		cli.ResponseFailed(http.StatusInternalServerError, err.Error(), resp)
		return
	}
	reply := make(map[string]interface{})
	err = json.Unmarshal([]byte(res), &reply)
	if nil != err {
		blog.Error("json unmarshal err :%v", err)
		cli.ResponseFailed(http.StatusInternalServerError, fmt.Sprintf("json unmarshal err %v", err.Error()), resp)
		return
	}
	io.WriteString(resp,res)

	//attrID := req.PathParameter("attr_id")
	//senceCLI := api.NewClient(module.CC.TopoAPI())
	//cli.CallResponse(
	//	senceCLI.ReForwardUpdateMetaObjAtt(func(url, method string) (string, error) {
	//		return httpclient.ReqForward(req, url, method)
	//	}, attrID),
	//	resp)

}

// SelectObjectAttWithParams search object's attributes with params
func (cli *objectAttAction) SelectObjectAttWithParams(req *restful.Request, resp *restful.Response) {

	blog.Info("select objectatt whith params")

	senceCLI := api.NewClient(module.CC.TopoAPI())
	cli.CallResponse(
		senceCLI.ReForwardSelectMetaObjAtt(func(url, method string) (string, error) {
			return httpclient.ReqForward(req, url, method)
		}), resp)

}
