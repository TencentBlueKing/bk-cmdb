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

package instdata

import (
	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/commondata"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
)

var userAPI *userAPIAction = &userAPIAction{tableName: "cc_UserAPI"}

type userAPIAction struct {
	base.BaseAction
	tableName string
}

//Add add new user api
func (u *userAPIAction) Add(req *restful.Request, resp *restful.Response) {

	value, _ := ioutil.ReadAll(req.Request.Body)

	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	params := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &params); nil != err {
		blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	name, _ := params["name"]
	if "" == name {
		blog.Error("parameter Name is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "Name").Error(), resp)
		return
	}

	appID, _ := util.GetInt64ByInterface(params[common.BKAppIDField])
	if 0 >= appID {
		blog.Error("parameter ApplicationID is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField).Error(), resp)
		return
	}
	user, _ := params["create_user"].(string)
	if "" == user {
		blog.Error("parameter CreateUser is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, "create_user").Error(), resp)
		return
	}

	queryParams := make(map[string]interface{})
	//queryParams["CreateUser"] = params["User"] //libraries.GetOperateUser(req)
	queryParams[common.BKAppIDField] = appID
	queryParams["name"] = name
	queryParams = util.SetModOwner(queryParams, ownerID)
	rowCount, err := userAPI.CC.InstCli.GetCntByCondition(u.tableName, queryParams)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), queryParams)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	if 0 != rowCount {
		blog.Error("[%s] user api is exist", name)
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommDuplicateItem, defErr.Error(common.CCErrCommDuplicateItem).Error(), resp)
		return
	}
	//mogo generate id for product
	xidDevice := xid.New()

	params[common.BKAppIDField] = appID
	params[common.BKFieldID] = xidDevice.String()
	params[common.CreateTimeField] = time.Now()
	params["modify_user"] = ""
	params[common.LastTimeField] = ""
	params = util.SetModOwner(params, ownerID)
	_, err = u.CC.InstCli.Insert(u.tableName, params)

	if err != nil {
		blog.Error("create user api  error:data:%v error:%v", params, err)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBInsertFailed, defErr.Error(common.CCErrCommDBInsertFailed).Error(), resp)
		return
	}

	info := make(map[string]interface{})
	info["id"] = xidDevice.String()

	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, info)
	io.WriteString(resp, rsp)

}

//Update update user api content
func (u *userAPIAction) Update(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	ID := req.PathParameter("id")
	appID, err := util.GetInt64ByInterface(req.PathParameter(common.BKAppIDField))

	value, _ := ioutil.ReadAll(req.Request.Body)
	data := make(map[string]interface{})
	if err := json.Unmarshal([]byte(value), &data); nil != err {
		blog.Error("fail to unmarshal json, error information is %s, msg:%s", err.Error(), string(value))
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}
	data[common.LastTimeField] = time.Now()

	params := make(map[string]interface{})
	params[common.BKFieldID] = ID
	params[common.BKAppIDField] = appID
	params = util.SetModOwner(params, ownerID)

	rowCount, err := u.CC.InstCli.GetCntByCondition(u.tableName, params)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), params)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	if 1 != rowCount {
		blog.Info("host user api not permissions or not exists, params:%v", params)
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommNotFound, defErr.Error(common.CCErrCommNotFound).Error(), resp)
		return
	}
	//edit new not duplicate
	newName, ok := data["name"]
	if ok {
		dupParams := make(map[string]interface{})
		dupParams["name"] = newName
		dupParams[common.BKAppIDField] = appID
		dupParams[common.BKFieldID] = common.KvMap{common.BKDBNE: ID}
		dupParams = util.SetModOwner(dupParams, ownerID)

		rowCount, getErr := u.CC.InstCli.GetCntByCondition(u.tableName, dupParams)
		if nil != getErr {
			blog.Error("query user api validate name duplicatie fail, error information is %s, params:%v", getErr.Error(), dupParams)
			userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
			return
		}
		if 0 < rowCount {
			blog.Info("host user api  name duplicatie , params:%v", dupParams)
			userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommDuplicateItem, defErr.Error(common.CCErrCommDuplicateItem).Error(), resp)
			return
		}
	}

	//json 中的数字会被转换未doubule， 转换未int64
	data[common.BKAppIDField] = appID
	err = u.CC.InstCli.UpdateByCondition(u.tableName, data, params)
	if nil != err {
		blog.Error("updata user api fail, error information is %s, params:%v", err.Error(), params)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBUpdateFailed, defErr.Errorf(common.CCErrCommDBUpdateFailed).Error(), resp)
		return
	}
	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, nil)
	io.WriteString(resp, rsp)
	return

}

//Delete delete user api
func (u *userAPIAction) Delete(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	ID := req.PathParameter("id")
	appID, _ := util.GetInt64ByInterface(req.PathParameter(common.BKAppIDField))

	params := make(map[string]interface{})
	params[common.BKAppIDField] = appID
	params["id"] = ID
	params = util.SetModOwner(params, ownerID)

	rowCount, err := u.CC.InstCli.GetCntByCondition(u.tableName, params)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), params)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	if 1 != rowCount {
		blog.Info("host user api not permissions or not exists, params:%v", params)
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommNotFound, defErr.Error(common.CCErrCommNotFound).Error(), resp)
		return
	}
	err = u.CC.InstCli.DelByCondition(u.tableName, params)
	if nil != err {
		blog.Error("delete user api fail, error information is %s, params:%v", err.Error(), params)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBDeleteFailed, defErr.Errorf(common.CCErrCommDBDeleteFailed).Error(), resp)
		return
	}
	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, nil)
	io.WriteString(resp, rsp)
	return
}

//Get get user api
func (u *userAPIAction) Get(req *restful.Request, resp *restful.Response) {

	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	var dat commondata.ObjQueryInput
	value, err := ioutil.ReadAll(req.Request.Body)
	//
	if nil == err && nil == value {
		value = []byte("{}")
	}
	err = json.Unmarshal([]byte(value), &dat)
	if err != nil {
		blog.Error("fail to unmarshal json, error information is,input:%v error:%v", string(value), err)
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommJSONUnmarshalFailed, defErr.Error(common.CCErrCommJSONUnmarshalFailed).Error(), resp)
		return
	}

	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}

	appID, _ := util.GetInt64ByInterface(condition[common.BKAppIDField])

	if 0 >= appID {
		blog.Error("The ApplicationID of the Condition parameter is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField).Error(), resp)
		return
	}

	skip := dat.Start
	limit := dat.Limit
	sort := dat.Sort
	var fieldArr []string
	if "" != dat.Fields {
		fieldArr = strings.Split(",", dat.Fields)

	}

	if 0 == limit {
		limit = 20
	}
	if "" == sort {
		sort = common.CreateTimeField
	}

	condition[common.BKAppIDField] = appID
	condition = util.SetModOwner(condition, ownerID)
	//result := make([]interface{}, 0)
	count, err := u.CC.InstCli.GetCntByCondition(u.tableName, condition)
	if err != nil {
		blog.Error("get user api infomation error,input:%v error:%v", string(value), err)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	var result []interface{}
	err = u.CC.InstCli.GetMutilByCondition(u.tableName, fieldArr, condition, &result, sort, skip, limit)
	if err != nil {
		blog.Error("get user api infomation error,input:%v error:%v", string(value), err)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}
	info := make(map[string]interface{})
	info["count"] = count
	info["info"] = result
	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, info)
	io.WriteString(resp, rsp)
	return
}

//Detail use api detail
func (u *userAPIAction) Detail(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	ownerID := util.GetActionOnwerID(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	appID, _ := util.GetInt64ByInterface(req.PathParameter(common.BKAppIDField))
	ID := req.PathParameter("id")

	params := make(map[string]interface{})
	params[common.BKAppIDField] = appID
	params["id"] = ID
	params = util.SetModOwner(params, ownerID)
	var fieldArr []string

	result := make(map[string]interface{})
	err := u.CC.InstCli.GetOneByCondition(u.tableName, fieldArr, params, &result)
	if err != nil && mgo_on_not_found_error != err.Error() {
		blog.Error("get user api infomation error,input:%v error:%v", ID, err)
		userAPI.ResponseFailedEx(http.StatusBadGateway, common.CCErrCommDBSelectFailed, defErr.Error(common.CCErrCommDBSelectFailed).Error(), resp)
		return
	}

	rsp, _ := u.CC.CreateAPIRspStr(common.CCSuccess, result)
	io.WriteString(resp, rsp)
	return

}

func init() {
	userAPI.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/userapi", Params: nil, Handler: userAPI.Add})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/userapi/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.Update})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/userapi/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.Delete})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/userapi/search", Params: nil, Handler: userAPI.Get})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/userapi/detail/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.Detail})
}
