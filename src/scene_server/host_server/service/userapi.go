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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

func (s *Service) AddUserCustomQuery(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ucq := new(meta.UserConfig)
	if err := json.NewDecoder(req.Request.Body).Decode(ucq); nil != err {
		blog.Errorf("AddUserCustomQuery add user custom query failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if "" == ucq.Name {
		blog.Error("AddUserCustomQuery add user custom query parameter name is required")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "name")})
		return
	}

	if 0 >= ucq.AppID {
		blog.Error("AddUserCustomQuery add user custom query parameter ApplicationID is required")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)})
		return
	}

	ucq.CreateUser = util.GetUser(req.Request.Header)
	result, err := s.CoreAPI.HostController().User().AddUserConfig(context.Background(), req.Request.Header, ucq)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("AddUserCustomQuery  dd user custom query failed,  err: %v, input:%v", err.Error(), ucq)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrAddUserCustomQueryFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) UpdateUserCustomQuery(req *restful.Request, resp *restful.Response) {
	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	params := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&params); nil != err {
		blog.Errorf("update user custom query failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	params["modify_user"] = util.GetActionUser(req)
	params[common.LastTimeField] = time.Now().UTC()
	result, err := s.CoreAPI.HostController().User().UpdateUserConfig(context.Background(), req.PathParameter("bk_biz_id"), req.PathParameter("id"), req.Request.Header, params)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("update user custom query failed,  err: %v, input:%v", err, params)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrUpdateUserCustomQueryFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     nil,
	})
	return
}

func (s *Service) DeleteUserCustomQuery(req *restful.Request, resp *restful.Response) {

	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	ID := req.PathParameter("id")
	appID := req.PathParameter("bk_biz_id")

	result, err := s.CoreAPI.HostController().User().DeleteUserConfig(context.Background(), appID, ID, req.Request.Header)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("delete user custom query failed,  err: %s, input:%v", err.Error(), req.PathParameters())
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrDeleteUserCustomQueryFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     nil,
	})

}

func (s *Service) GetUserCustomQuery(req *restful.Request, resp *restful.Response) {

	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	input := &meta.QueryInput{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); nil != err {
		blog.Errorf("get user custom query failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	var condition map[string]interface{}
	if nil != input.Condition {
		condition, _ = input.Condition.(map[string]interface{})
	} else {
		condition = make(map[string]interface{})
	}
	//if name in condition , add like search
	name, ok := condition["name"].(string)
	if ok && "" != name {
		condition["name"] = common.KvMap{common.BKDBLIKE: params.SpeceialCharChange(name)}
	}

	var err error
	condition[common.BKAppIDField], err = util.GetInt64ByInterface(req.PathParameter("bk_biz_id"))
	if nil != err {
		blog.Error("UserAPIGet query user custom query parameter ApplicationID not integer in url")
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}
	input.Condition = condition

	result, err := s.CoreAPI.HostController().User().GetUserConfig(context.Background(), req.Request.Header, input)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("delete user custom query failed,err: %v, input:%v", err, input)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrSearchUserCustomQueryFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) GetUserCustomQueryDetail(req *restful.Request, resp *restful.Response) {

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	result, err := s.CoreAPI.HostController().User().GetUserConfigDetail(context.Background(), appID, ID, req.Request.Header)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("GetUserConfigDetail custom query failed,  err: %v, appid:%s, id:%s", err.Error(), appID, ID)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrGetUserCustomQueryDetailFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})

}

func (s *Service) GetUserCustomQueryResult(req *restful.Request, resp *restful.Response) {

	language := util.GetLanguage(req.Request.Header)
	defErr := s.CCErr.CreateDefaultCCErrorIf(language)

	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	intAppID, err := util.GetInt64ByInterface(appID)
	if nil != err {
		blog.Errorf("UserAPIResult custom query failed,  err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}

	result, err := s.CoreAPI.HostController().User().GetUserConfigDetail(context.Background(), appID, ID, req.Request.Header)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("UserAPIResult custom query failed,  err: %v, appid:%s, id:%s", err.Error(), appID, ID)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrGetUserCustomQueryDetailFaild, err.Error())})
		return
	}

	if "" == result.Data.Name {
		blog.Errorf("UserAPIResult custom query not found, appid:%s, id:%s, logID:%s", appID, ID, util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Errorf(common.CCErrCommNotFound)})
		return
	}

	var input meta.HostCommonSearch
	input.AppID = intAppID

	err = json.Unmarshal([]byte(result.Data.Info), &input)
	if nil != err {
		blog.Errorf("UserAPIResult custom unmarshal failed,  err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input.Page.Start, _ = util.GetIntByInterface(req.PathParameter("start"))
	input.Page.Limit, _ = util.GetIntByInterface(req.PathParameter("limit"))

	retData, err := s.Logics.SearchHost(req.Request.Header, &input, false)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("UserAPIResult custom query search host failed,  err: %v, appid:%s, id:%s", err.Error(), appID, ID)
		resp.WriteError(http.StatusBadGateway, &meta.RespError{Msg: defErr.Errorf(common.CCErrGetUserCustomQueryDetailFaild, err.Error())})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data: meta.SearchHost{
			Count: retData.Count,
			Info:  retData.Info,
		},
	})

	return
}
