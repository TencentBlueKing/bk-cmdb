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
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/base"
	"configcenter/src/common/blog"
	"configcenter/src/common/core/cc/actions"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/common/commondata"
	"github.com/emicklei/go-restful"
	"github.com/rs/xid"
	"strconv"
)

func init() {
	userAPI.CreateAction()
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPCreate, Path: "/userapi", Params: nil, Handler: userAPI.AddUserConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPUpdate, Path: "/userapi/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.UpdateUserConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPDelete, Path: "/userapi/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.DeleteUserConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectPost, Path: "/userapi/search", Params: nil, Handler: userAPI.GetUserConfig})
	actions.RegisterNewAction(actions.Action{Verb: common.HTTPSelectGet, Path: "/userapi/detail/{bk_biz_id}/{id}", Params: nil, Handler: userAPI.GetUserConfigDetail})
}

var userAPI *userAPIAction = &userAPIAction{tableName: "cc_UserAPI"}

type userAPIAction struct {
	base.BaseAction
	tableName string
}

//AddUserConfig add new user api
func (u *userAPIAction) AddUserConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	c := new(meta.UserConfig)
	if err := json.NewDecoder(req.Request.Body).Decode(c); err != nil {
		blog.Errorf("add user config failed, decode err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(c.Name) == 0 || c.AppID == 0 || len(c.CreateUser) == 0 {
		blog.Errorf("add user config failed, err: invalid user config. name , appid and create user can not be empty.")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
	}

	queryParams := make(map[string]interface{})
	queryParams[common.BKAppIDField] = c.AppID
	queryParams["name"] = c.Name

	rowCount, err := userAPI.CC.InstCli.GetCntByCondition(u.tableName, queryParams)
	if nil != err {
		blog.Error("query user api fail, error information is %s, params:%v", err.Error(), queryParams)

		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	if 0 != rowCount {
		blog.Error("[%s] user api is exist", c.Name)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
		return
	}
	//mogo generate id for product
	xidDevice := xid.New()

	conf := make(map[string]interface{})
	conf[common.BKAppIDField] = c.AppID
	conf[common.BKFieldID] = xidDevice.String()
	conf[common.CreateTimeField] = time.Now().UTC()
	conf["modify_user"] = ""
	conf[common.LastTimeField] = ""
	_, err = u.CC.InstCli.Insert(u.tableName, conf)
	if err != nil {
		blog.Error("create user api failed, data:%v error:%v", conf, err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBInsertFailed)})
		return
	}

	resp.WriteEntity(meta.IDResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.ID{ID: xidDevice.String()},
	})
}

//UpdateUserConfig update user api content
func (u *userAPIAction) UpdateUserConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetActionLanguage(req)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)
	ID := req.PathParameter("id")

	appID, err := strconv.ParseInt(req.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(data); err != nil {
		blog.Errorf("del module host config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	data[common.LastTimeField] = time.Now().UTC()

	params := make(map[string]interface{})
	params[common.BKFieldID] = ID
	params[common.BKAppIDField] = appID
	rowCount, err := u.CC.InstCli.GetCntByCondition(u.tableName, params)
	if nil != err {
		blog.Errorf("query user api fail, error information is %s, params: %v", err.Error(), params)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	if 1 != rowCount {
		blog.Errorf("host user api not permissions or not exists, params: %v", params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommNotFound)})
		return
	}
	//edit new not duplicate
	newName, ok := data["name"]
	if ok {
		dupParams := make(map[string]interface{})
		dupParams["name"] = newName
		dupParams[common.BKAppIDField] = appID
		dupParams[common.BKFieldID] = common.KvMap{common.BKDBNE: ID}

		rowCount, getErr := u.CC.InstCli.GetCntByCondition(u.tableName, dupParams)
		if nil != getErr {
			blog.Error("query user api validate name duplicate fail, error information is %s, params:%v", getErr.Error(), dupParams)
			resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
			return
		}
		if 0 < rowCount {
			blog.Errorf("host user api  name duplicate , params: %v", dupParams)
			resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDuplicateItem)})
			return
		}
	}

	//json 中的数字会被转换为double， 转换未int64
	data[common.BKAppIDField] = appID
	err = u.CC.InstCli.UpdateByCondition(u.tableName, data, params)
	if nil != err {
		blog.Errorf("update user api fail, error information is %s, params:%v", err.Error(), params)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBUpdateFailed)})
		return
	}
	resp.WriteEntity(meta.NewSuccessResp(nil))
}

//DeleteUserConfig delete user api
func (u *userAPIAction) DeleteUserConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	ID := req.PathParameter("id")
	appID, err := strconv.ParseInt(req.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("update user config failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	params := make(map[string]interface{})
	params[common.BKAppIDField] = appID
	params["id"] = ID

	rowCount, err := u.CC.InstCli.GetCntByCondition(u.tableName, params)
	if nil != err {
		blog.Errorf("query user api fail, err: %v, params: %v", err, params)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}
	if 1 != rowCount {
		blog.Errorf("host user api not permissions or not exists, params:%v", params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommNotFound)})
		return
	}
	err = u.CC.InstCli.DelByCondition(u.tableName, params)
	if nil != err {
		blog.Errorf("delete user api fail, err: %v, params:%v", err, params)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBDeleteFailed)})
		return
	}

	resp.WriteEntity(meta.NewSuccessResp(nil))
	return
}

//GetUserConfig get user api
func (u *userAPIAction) GetUserConfig(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	var dat commondata.ObjQueryInput
	if err := json.NewDecoder(req.Request.Body).Decode(&dat); err != nil {
		blog.Errorf("get user config, but decode body failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	condition := make(map[string]interface{})
	if nil != dat.Condition {
		condition = dat.Condition.(map[string]interface{})
	}

	appID, err := util.GetInt64ByInterface(condition[common.BKAppIDField])
	if err != nil {
		blog.Errorf("get user config, but parse appid failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	if 0 >= appID {
		blog.Error("The ApplicationID of the Condition parameter is required")
		userAPI.ResponseFailedEx(http.StatusBadRequest, common.CCErrCommParamsNeedSet, defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField).Error(), resp)
		return
	}
	count, err := u.CC.InstCli.GetCntByCondition(u.tableName, condition)
	if err != nil {
		blog.Error("get user api information failed, err:%v", err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	var fieldArr []string
	if "" != dat.Fields {
		fieldArr = strings.Split(",", dat.Fields)

	}

	skip, limit, sort := dat.Start, dat.Limit, dat.Sort
	if 0 == limit {
		limit = 20
	}
	if "" == sort {
		sort = common.CreateTimeField
	}

	var result []interface{}
	err = u.CC.InstCli.GetMutilByCondition(u.tableName, fieldArr, condition, &result, sort, skip, limit)
	if err != nil {
		blog.Error("get user api information failed, err: %v", err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.GetUserConfigResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     meta.UserConfigResult{Count: count, Info: result},
	})
}

//Detail use api detail
func (u *userAPIAction) GetUserConfigDetail(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := u.CC.Error.CreateDefaultCCErrorIf(language)

	appID, err := strconv.ParseInt(req.PathParameter(common.BKAppIDField), 10, 64)
	if err != nil {
		blog.Errorf("get user config detail failed, err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommParamsIsInvalid)})
		return
	}

	ID := req.PathParameter("id")
	params := make(map[string]interface{})
	params[common.BKAppIDField] = appID
	params["id"] = ID

	var fieldArr []string
	result := make(map[string]interface{})
	err = u.CC.InstCli.GetOneByCondition(u.tableName, fieldArr, params, &result)
	if err != nil && mgo_on_not_found_error != err.Error() {
		blog.Errorf("get user api detail failed, id: %v err: %v", ID, err)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommDBSelectFailed)})
		return
	}

	resp.WriteEntity(meta.GetUserConfigDetailResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     result,
	})
}
