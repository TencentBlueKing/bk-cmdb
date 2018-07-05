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

package detail

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/scene_server/host_server/host_service/instapi"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"encoding/json"

	"configcenter/src/common/util"
	"net/http"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/emicklei/go-restful"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/hosts/{bk_supplier_account}/{bk_host_id}", Params: nil, Handler: host.GetHostDetailByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/host/snapshot/{bk_host_id}", Params: nil, Handler: host.HostSnapInfo})

	// create CC object
	host.CreateAction()
}

//GetHostDetailByID get host detail
func (cli *hostAction) GetHostDetailByID(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		pathParams := req.PathParameters()
		hostID := pathParams["bk_host_id"]

		//ownerID := pathParams["bk_supplier_account"]
		ownerID := util.GetActionOnwerID(req)

		//gHostURL := "http://" + cli.CC.HostCtrl + "/host/v1/host/" + hostID
		gHostURL := cli.CC.HostCtrl() + "/host/v1/host/" + hostID

		gHostRe, err := httpcli.ReqHttp(req, gHostURL, common.HTTPSelectGet, nil)
		if nil != err {
			blog.Error("GetHostDetailByID info error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostDetailFail)

		}

		// deal the association id
		instapi.Inst.InitInstHelper(cli.CC.HostCtrl(), cli.CC.ObjCtrl())
		gHostRe, retStrErr := instapi.Inst.GetInstDetails(req, common.BKInnerObjIDHost, ownerID, gHostRe, map[string]interface{}{
			"start": 0,
			"limit": common.BKNoLimit,
			"sort":  "",
		})

		if common.CCSuccess != retStrErr {
			blog.Error("failed to replace association object, error code is %d", retStrErr)
		}
		//
		js, err := simplejson.NewJson([]byte(gHostRe))
		gHostData, _ := js.Map()
		gResult := gHostData["result"].(bool)
		if false == gResult {
			blog.Error("GetHostDetailByID  info error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostDetailFail)
		}

		//
		hostData := gHostData["data"].(map[string]interface{})

		//gHostAttrURL := "http://" + cli.CC.ObjCtrl + "/object/v1/meta/objectatts"
		gHostAttrURL := cli.CC.ObjCtrl() + "/object/v1/meta/objectatts"
		searchBody := make(map[string]interface{})
		searchBody[common.BKObjIDField] = common.BKInnerObjIDHost
		searchBody[common.BKOwnerIDField] = ownerID
		searchJson, _ := json.Marshal(searchBody)
		gHostAttrRe, err := httpcli.ReqHttp(req, gHostAttrURL, common.HTTPSelectPost, []byte(searchJson))
		if nil != err {
			blog.Error("GetHostDetailByID  attr error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostDetailFail)
		}

		js, err = simplejson.NewJson([]byte(gHostAttrRe))
		gHostAttr, _ := js.Map()
		gAttrResult := gHostAttr["result"].(bool)
		if false == gAttrResult {
			blog.Error("GetHostDetailByID  attr error :%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostDetailFail)
		}
		hostAttrArr := gHostAttr["data"].([]interface{})
		reResult := make([]map[string]interface{}, 0)
		for _, i := range hostAttrArr {
			attr := i.(map[string]interface{})
			data := make(map[string]interface{})
			propertyID := attr[common.BKPropertyIDField].(string)
			if propertyID == common.BKChildStr {
				continue
			}
			data[common.BKPropertyIDField] = propertyID
			data[common.BKPropertyNameField] = attr[common.BKPropertyNameField]

			data[common.BKPropertyValueField] = hostData[propertyID]
			reResult = append(reResult, data)
		}
		return http.StatusOK, reResult, nil
	}, resp)
}

func (cli *hostAction) HostSnapInfo(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)

	cli.CallResponseEx(func() (int, interface{}, error) {
		url := host.CC.HostCtrl() + "/host/v1/host/snapshot/" + req.PathParameter("bk_host_id")
		reply, err := httpcli.ReqHttp(req, url, common.HTTPSelectGet, nil)
		if nil != err {
			blog.Error("query host snapshot   error, error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSnap)
		}
		js, err := simplejson.NewJson([]byte(reply))
		if nil != err {
			blog.Error("query host snapshot   error, error:%s", err.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSnap)

		}
		output, _ := js.Map()
		isSucc, _ := output["result"].(bool)
		if !isSucc {
			blog.Error("query host snapshot   error, error:%s", string(reply))
			errcode, err := js.Get("bk_error_code").Int()
			if err != nil {
				return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSnap)
			}
			return http.StatusInternalServerError, nil, defErr.Error(errcode)
		} else {
			snap, err := js.Get("data").Get("data").String()

			if nil == err && "" != snap {
				ret, err := logics.ParseHostSnap(snap)
				if nil != err {
					blog.Error("query host snapshot   error, error:%s", string(reply))
					return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostSnap)

				} else {
					return http.StatusOK, ret, nil
				}

				return http.StatusOK, "", nil

			} else {
				return http.StatusOK, "", nil
			}

		}

	}, resp)

}
