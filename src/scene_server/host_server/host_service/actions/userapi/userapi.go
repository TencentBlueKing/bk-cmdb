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

package userapi

import (
	"configcenter/src/common"
	"configcenter/src/common/bkbase"
	"configcenter/src/common/blog"
	hostParse "configcenter/src/common/paraparse"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/host_service/logics"
	userAPISdk "configcenter/src/source_controller/api/userapi"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"configcenter/src/common/core/cc/actions"

	restful "github.com/emicklei/go-restful"
)

var userAPI *userAPIAction = &userAPIAction{}

type userAPIAction struct {
	base.BaseAction
}

func init() {
	userAPI.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/userapi", Params: nil, Handler: userAPI.Add})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/userapi/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.Update})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/userapi/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.Delete})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/userapi/search/{bk_biz_id}", Params: nil, Handler: userAPI.Get})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/userapi/detail/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.Detail})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/userapi/data/{bk_biz_id}/{id}/{start}/{limit}", Params: nil, Handler: userAPI.GetUserAPIData})

}

//Add add new user api
func (u *userAPIAction) Add(req *restful.Request, resp *restful.Response) {

	value, _ := ioutil.ReadAll(req.Request.Body)

	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	URL, err := u.CC.AddrSrv.GetServer(types.CC_MODULE_HOSTCONTROLLER)
	if nil != err {
		blog.Errorf("get host addr from service discovery module error: %s", err.Error())
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, "hostcontroller").Error(), nil, resp)
		return
	}
	client := userAPISdk.NewClient(URL)

	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &params); nil != err {
		blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil, resp)
		return
	}

	name, _ := params["name"]
	if "" == name {
		blog.Error("parameter name is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "bk_name").Error(), nil, resp)
		return
	}

	appID, _ := util.GetInt64ByInterface(params[common.BKAppIDField])
	if 0 >= appID {
		blog.Error("parameter ApplicationID is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField).Error(), nil, resp)
		return
	}
	params["create_user"] = util.GetActionUser(req)
	code, reply, err := client.Create(params)
	if nil != err {
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CC_Err_Comm_Host_Get_FAIL, err.Error(), nil, resp)
		return
	}
	if code != http.StatusOK {
		userAPI.ResponseFailedEx(code, reply.Code, reply.Message, nil, resp)
		return
	}

	u.ResponseSuccess(reply.Data, resp)
	return

}

//Update update user api content
func (u *userAPIAction) Update(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	URL, err := u.CC.AddrSrv.GetServer(types.CC_MODULE_HOSTCONTROLLER)
	if nil != err {
		blog.Errorf("get host addr from service discovery module error: %s", err.Error())
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, "hostcontroller").Error(), nil, resp)
		return
	}
	value, _ := ioutil.ReadAll(req.Request.Body)

	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &params); nil != err {
		blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil, resp)
		return
	}
	params["modify_user"] = util.GetActionUser(req)

	client := userAPISdk.NewClient(URL)
	code, reply, err := client.Update(params, req.PathParameter("bk_biz_id"), req.PathParameter("id"))
	if nil != err {
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CC_Err_Comm_http_DO, err.Error(), nil, resp)
		return
	}
	if code != http.StatusOK {
		userAPI.ResponseFailedEx(code, reply.Code, reply.Message, nil, resp)
		return
	}
	u.ResponseSuccess(reply.Data, resp)
	return

}

//Delete delete user api
func (u *userAPIAction) Delete(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	ID := req.PathParameter("id")
	appID := (req.PathParameter("bk_biz_id"))

	URL, err := u.CC.AddrSrv.GetServer(types.CC_MODULE_HOSTCONTROLLER)
	if nil != err {
		blog.Errorf("get host addr from service discovery module error: %s", err.Error())
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, "hostcontroller").Error(), nil, resp)
		return
	}
	client := userAPISdk.NewClient(URL)
	code, reply, err := client.Delete(appID, ID)
	if nil != err {
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CC_Err_Comm_http_DO, err.Error(), nil, resp)
		return
	}
	if code != http.StatusOK {
		userAPI.ResponseFailedEx(code, reply.Code, reply.Message, nil, resp)
		return
	}

	u.ResponseSuccess(nil, resp)
	return

}

//Get get user api
func (u *userAPIAction) Get(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	url, err := u.CC.AddrSrv.GetServer(types.CC_MODULE_HOSTCONTROLLER)
	if nil != err {
		blog.Errorf("get host addr from service discovery module error: %s", err.Error())
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, "hostcontroller").Error(), nil, resp)
		return
	}

	var dat commondata.ObjQueryInput
	value, err := ioutil.ReadAll(req.Request.Body)

	//no default value
	if nil == err && nil == value {
		value = []byte("{}")
	}
	err = json.Unmarshal([]byte(value), &dat)
	if err != nil {
		blog.Error("fail to unmarshal json, error information is,input:%v error:%v", string(value), err)
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), nil, resp)
		return
	}

	var condition map[string]interface{}
	if nil != dat.Condition {
		condition, _ = dat.Condition.(map[string]interface{})
	} else {
		condition = make(map[string]interface{})
	}
	//if name in condition , add like search
	name, ok := condition["name"].(string)
	if ok && "" != name {
		condition["name"] = common.KvMap{common.BKDBLIKE: hostParse.SpeceialCharChange(name)}
	}

	condition[common.BKAppIDField], _ = util.GetInt64ByInterface(req.PathParameter("bk_biz_id"))
	dat.Condition = condition

	client := userAPISdk.NewClient(url)
	code, reply, err := client.GetUserAPI(dat)
	if nil != err {
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CC_Err_Comm_http_DO, err.Error(), nil, resp)
		return
	}
	if code != http.StatusOK {
		userAPI.ResponseFailedEx(code, reply.Code, reply.Message, nil, resp)
		return
	}

	u.ResponseSuccess(reply.Data, resp)
	return
}

//Detail get user api detail
func (u *userAPIAction) Detail(req *restful.Request, resp *restful.Response) {

	hostCtrl := u.CC.HostCtrl()

	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	client := userAPISdk.NewClient(hostCtrl)
	code, reply, err := client.Detail(appID, ID)
	if nil != err {
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CC_Err_Comm_http_DO, err.Error(), nil, resp)
		return
	}
	if code != http.StatusOK {
		userAPI.ResponseFailedEx(code, reply.Code, reply.Message, nil, resp)
		return
	}

	u.ResponseSuccess(reply.Data, resp)
	return

}

//GetUserAPIData get api data
func (u *userAPIAction) GetUserAPIData(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)

	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)
	URL, err := u.CC.AddrSrv.GetServer(types.CC_MODULE_HOSTCONTROLLER)
	if nil != err {
		blog.Errorf("get host addr from service discovery module error: %s", err.Error())
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommRelyOnServerAddressFailed, defErr.Errorf(common.CCErrCommRelyOnServerAddressFailed, "hostcontroller").Error(), nil, resp)
		return
	}

	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	client := userAPISdk.NewClient(URL)
	code, reply, err := client.Detail(appID, ID)
	if nil != err {
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CC_Err_Comm_http_DO, err.Error(), nil, resp)
		return
	}
	if code != http.StatusOK {
		userAPI.ResponseFailedEx(code, reply.Code, reply.Message, nil, resp)
		return
	}

	data, _ := reply.Data.(map[string]interface{})

	cond, _ := data["info"].(string)
	if "" == cond {
		blog.Error("user api detail return, code:%d , data:%v", code, reply)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommNotFound, defErr.Error(common.CCErrCommNotFound).Error(), nil, resp)
		return
	}

	var input hostParse.HostCommonSearch

	err = json.Unmarshal([]byte(cond), &input)
	if nil != err {
		userAPI.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, err.Error(), resp)
		return
	}

	input.AppID, _ = util.GetIntByInterface(data[common.BKAppIDField])
	if fmt.Sprintf("%d", input.AppID) != appID {
		userAPI.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, "请输入正确的业务ID", resp)
		return
	}

	input.Page.Start, _ = util.GetIntByInterface(req.PathParameter("start"))
	input.Page.Limit, _ = util.GetIntByInterface(req.PathParameter("limit"))

	retData, err := logics.HostSearch(req, input, false, u.CC.HostCtrl(), u.CC.ObjCtrl())
	if nil != err {
		userAPI.ResponseFailed(common.CC_Err_Comm_Host_Get_FAIL, err.Error(), resp)
		return
	}

	u.ResponseSuccess(retData, resp)
	return

}
