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
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	restful "github.com/emicklei/go-restful"
)

var host *hostAction = &hostAction{}

type hostAction struct {
	base.BaseAction
}

// FavouriteParms user request params
type FavouriteParms struct {
	ID          string `json:"id"`
	Info        string `json:"info"`
	QueryParams string `json:"query_params"`
	Name        string `json:"name"`
	IsDefault   int    `json:"is_default"`
	Count       int    `json:"count"`
}

func init() {
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "hosts/favorites/search", Params: nil, Handler: host.GetHostFavourites})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "hosts/favorites", Params: nil, Handler: host.AddHostFavourite})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "hosts/favorites/{id}", Params: nil, Handler: host.UpdateHostFavouriteByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "hosts/favorites/{id}", Params: nil, Handler: host.DeleteHostFavouriteByID})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/hosts/favorites/{id}/incr", Params: nil, Handler: host.IncrHostFavouritesCount})

	// create CC object
	host.CreateAction()
}

// AddHostFavourite  add host favourite
func (cli *hostAction) AddHostFavourite(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		user := util.GetActionUser(req)
		value, err := ioutil.ReadAll(req.Request.Body)
		data := FavouriteParms{}
		err = json.Unmarshal([]byte(value), &data)

		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)
		}
		if "" == data.Name {
			blog.Error("get unmarshall json value %v error:名字不能为空")
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostEmptyFavName)

		}

		url := cli.CC.HostCtrl() + "/host/v1/hosts/favorites/" + user
		isSuccess, errMsg, retData := logics.GetHttpResult(req, url, common.HTTPCreate, data)
		if !isSuccess {
			blog.Error("add host favorites error, params:%v, error:%s", string(value), errMsg)
			return http.StatusInternalServerError, nil, defErr.Error(common.CCErrHostFavCreateFail)

		}
		return http.StatusOK, retData, nil
	}, resp)
}

// UpdateHostFavouriteByID update host favourite
func (cli *hostAction) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		ID := req.PathParameter("id")

		if "" == ID || "0" == ID {
			blog.Error("host favourite id  %id", ID)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		user := util.GetActionUser(req)

		value, err := ioutil.ReadAll(req.Request.Body)
		data := make(map[string]interface{})
		err = json.Unmarshal([]byte(value), &data)
		if err != nil {
			blog.Error("get unmarshall json value %v error:%v", string(value), err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		if nil == data["name"] || "" == data["name"] {
			blog.Error("get unmarshall json value %v error:名字不能为空")
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostEmptyFavName)

		}

		url := cli.CC.HostCtrl() + "/host/v1/hosts/favorites/" + user
		url = fmt.Sprintf("%s/%s", url, ID)
		isSuccess, errMsg, retData := logics.GetHttpResult(req, url, common.HTTPUpdate, data)
		if !isSuccess {
			blog.Error("Edit host favourite error, params:%v, error:%s", ID, errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavUpdateFail)

		}

		return http.StatusOK, retData, nil
	}, resp)

}

// DeleteHostFavouriteByID delete host favourite
func (cli *hostAction) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		ID := req.PathParameter("id")

		if "" == ID || "0" == ID {
			blog.Error("host favourite id  %id", ID)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		user := util.GetActionUser(req)

		URL := cli.CC.HostCtrl() + "/host/v1/hosts/favorites"
		URL = fmt.Sprintf("%s/%s/%s", URL, user, ID)
		isSuccess, errMsg, retData := logics.GetHttpResult(req, URL, common.HTTPDelete, "")
		if !isSuccess {
			blog.Error("delete host favourite error, params:id:%s, error:%s", ID, errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavDeleteFail)

		}

		return http.StatusOK, retData, nil

	}, resp)

}

// GetHostFavourites get host favourites
func (cli *hostAction) GetHostFavourites(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {
		value, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			blog.Error("get params value  error:%v", err)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		user := util.GetActionUser(req)

		url := cli.CC.HostCtrl() + "/host/v1/hosts/favorites/search/" + user
		isSuccess, errMsg, retData := logics.GetHttpResult(req, url, common.HTTPSelectPost, string(value))
		if !isSuccess {
			blog.Error("query host favourite error, params:%s, error:%s", string(value), errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavGetFail)

		}
		return http.StatusOK, retData, nil

	}, resp)
}

// IncrHostFavouritesCount increase host favourites count
func (cli *hostAction) IncrHostFavouritesCount(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	// get the error factory by the language
	defErr := cli.CC.Error.CreateDefaultCCErrorIf(language)
	cli.CallResponseEx(func() (int, interface{}, error) {

		ID := req.PathParameter("id")

		if "" == ID || "0" == ID {
			blog.Error("host favourite id  %id", ID)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrCommHTTPInputInvalid)

		}
		user := util.GetActionUser(req)

		url := cli.CC.HostCtrl() + "/host/v1/hosts/favorites/search"
		url = fmt.Sprintf("%s/%s/%s", url, user, ID)
		isSuccess, errMsg, retData := logics.GetHttpResult(req, url, common.HTTPSelectPost, "")
		if !isSuccess {
			blog.Error("get host favourite error, params:%v, error:%s", ID, errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavGetFail)

		}

		row := retData.(map[string]interface{})
		count, ok := util.GetIntByInterface(row["count"])
		if nil != ok {
			blog.Error("get host favourite error, params:%v", row)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavGetFail)

		}
		count = count + 1
		newData := make(map[string]interface{}, 1)
		newData["count"] = count

		url = cli.CC.HostCtrl() + "/host/v1/hosts/favorites"
		url = fmt.Sprintf("%s/%s/%s", url, user, ID)
		isSuccess, errMsg, _ = logics.GetHttpResult(req, url, common.HTTPUpdate, newData)
		if !isSuccess {
			blog.Error("Edit host favourite error, params:%v, error:%s", ID, errMsg)
			return http.StatusBadRequest, nil, defErr.Error(common.CCErrHostFavUpdateFail)

		}
		info := make(map[string]interface{}, 3)
		info["id"] = ID
		info["count"] = count
		return http.StatusOK, retData, nil

	}, resp)
}
