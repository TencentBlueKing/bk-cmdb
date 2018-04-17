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

package openapi

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/core/cc/api"
	httpcli "configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/validator"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

var plat *platAction = &platAction{}

type platAction struct {
	base.BaseAction
}

func init() {

	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/plat", Params: nil, Handler: plat.GetPlat})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/plat", Params: nil, Handler: plat.CreatePlat})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/plat/{" + common.BKCloudIDField + "}", Params: nil, Handler: plat.DelPlat})
	// create CC object
	plat.CreateAction()
}

// GetPlat: 获取所有子网
func (cli *platAction) GetPlat(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		searchParams := make(map[string]interface{})
		searchParams["condition"] = nil
		searchParams["fields"] = ""
		searchParams["start"] = 0
		searchParams["limit"] = 0
		searchParams["sort"] = ""

		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/plat/search"
		inputJson, _ := json.Marshal(searchParams)
		res, err := httpcli.ReqHttp(req, sURL, common.HTTPSelectPost, []byte(inputJson))

		if nil != err {
			blog.Error("GetPlat error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatQueryFailed)

		}

		// deal result
		var rst api.BKAPIRsp
		if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
			cli.Response(&rst, resp)
			return http.StatusOK, rst.Data, nil
		} else {
			blog.Error("unmarshal the json failed, error information is %v", jserr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatQueryFailed)
		}
		return http.StatusOK, rst.Data, nil

	}, resp)
}

// DelPlat: 删除子网
func (cli *platAction) DelPlat(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		platID, convErr := strconv.Atoi(req.PathParameter(common.BKCloudIDField))
		if nil != convErr {
			blog.Error("the platID is invalid, error info is %s", convErr.Error())
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatDeleteFailed)
		}

		url := host.CC.HostCtrl() + "/host/v1/hosts/search"
		searchParams := map[string]interface{}{
			"fields": common.BKHostIDField,
			"condition": map[string]interface{}{
				common.BKCloudIDField: platID,
			},
		}
		inputJson, _ := json.Marshal(searchParams)

		hostInfo, err := httpcli.ReqHttp(req, url, common.HTTPSelectPost, []byte(inputJson))
		if nil != err {
			blog.Error("search host error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatDeleteFailed)
		}

		hostResMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(hostInfo), &hostResMap)
		if nil != err {
			blog.Error("search host error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatDeleteFailed)
		}

		hostResDataMap := hostResMap["data"].(map[string]interface{})
		hostCount := hostResDataMap["count"].(float64)

		if hostCount > 0 {
			blog.Error("plat [%d] has host data, can not delete", platID)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoHostInPlatFailed)

		}

		param := make(map[string]interface{})
		param[common.BKCloudIDField] = platID

		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/plat"
		paramJson, _ := json.Marshal(param)
		res, err := httpcli.ReqHttp(req, sURL, common.HTTPDelete, []byte(paramJson))

		if nil != err {
			blog.Error("DelPlat error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatDeleteFailed)
		}

		// deal result
		var rst api.BKAPIRsp
		if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
			cli.Response(&rst, resp)
			return http.StatusOK, rst.Data, nil
		} else {
			blog.Error("unmarshal the json failed, error information is %v", jserr)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatDeleteFailed)
		}
	}, resp)
}

// CreatePlat: 添加子网
func (cli *platAction) CreatePlat(req *restful.Request, resp *restful.Response) {
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(util.GetActionLanguage(req))

	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if nil != err {
			blog.Error("read request body failed, error:%v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrCommHTTPReadBodyFailed)
		}

		input := make(map[string]interface{})
		err = json.Unmarshal(value, &input)
		if nil != err {
			blog.Error("Unmarshal json failed, error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}

		//param := make(map[string]interface{})
		//param["PlatName"] = input["PlatName"]

		input[common.BKOwnerIDField] = common.BKDefaultOwnerID

		sURL := cli.CC.ObjCtrl() + "/object/v1/insts/plat"
		language := util.GetActionLanguage(req)
		valid := validator.NewValidMap(input[common.BKOwnerIDField].(string), common.BKInnerObjIDPlat, cli.CC.ObjCtrl(), cli.CC.Error.CreateDefaultCCErrorIf(language))
		ok, validErr := valid.ValidMap(input, common.ValidCreate, 0)
		if false == ok || nil != validErr {
			blog.Error("CreatePlat error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatCreateFailed)

		}
		inputJson, _ := json.Marshal(input)
		res, err := httpcli.ReqHttp(req, sURL, common.HTTPCreate, []byte(inputJson))

		if nil != err {
			blog.Error("CreatePlat error: %v", err)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrTopoPlatCreateFailed)
		}

		// deal result
		var rst api.BKAPIRsp
		if jserr := json.Unmarshal([]byte(res), &rst); nil == jserr {
			// cli.Response(&rst, resp)
			return http.StatusOK, rst.Data, nil
		} else {
			blog.Error("unmarshal the json failed, error information is %v", jserr)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommJSONUnmarshalFailed)
		}
	}, resp)
}
