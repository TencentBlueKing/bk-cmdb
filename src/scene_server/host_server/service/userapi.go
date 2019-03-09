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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/paraparse"
	"configcenter/src/common/util"
)

func (s *Service) AddUserCustomQuery(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	ucq := new(meta.UserConfig)
	if err := json.NewDecoder(req.Request.Body).Decode(ucq); nil != err {
		blog.Errorf("AddUserCustomQuery add user custom query failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if "" == ucq.Name {
		blog.Error("AddUserCustomQuery add user custom query parameter name is required,input:%+v,rid:%s", ucq, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "name")})
		return
	}

	if 0 >= ucq.AppID {
		blog.Error("AddUserCustomQuery add user custom query parameter ApplicationID is required,input:%+v,rid:%s", ucq, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, common.BKAppIDField)})
		return
	}

	ucq.CreateUser = srvData.user
	result, err := s.CoreAPI.HostController().User().AddUserConfig(srvData.ctx, srvData.header, ucq)
	if err != nil {
		blog.Errorf("GetUserCustom http do error,err:%s, input:%+v,rid:%s", err.Error(), ucq, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustom http response error,err code:%d,err msg:%s, input:%+v,rid:%s", result.Code, result.ErrMsg, ucq, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) UpdateUserCustomQuery(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	params := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&params); nil != err {
		blog.Errorf("update user custom query failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	params["modify_user"] = srvData.user
	params[common.LastTimeField] = time.Now().UTC()
	bizID := req.PathParameter("bk_biz_id")
	result, err := s.CoreAPI.HostController().User().UpdateUserConfig(srvData.ctx, bizID, req.PathParameter("id"), srvData.header, params)
	if err != nil {
		blog.Errorf("UpdateUserCustomQuery http do error,err:%s, biz:%v,input:%+v,rid:%s", err.Error(), bizID, params, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("UpdateUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,input:%+v,rid:%s", result.Code, result.ErrMsg, bizID, params, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     nil,
	})
	return
}

func (s *Service) DeleteUserCustomQuery(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)

	ID := req.PathParameter("id")
	appID := req.PathParameter("bk_biz_id")

	result, err := s.CoreAPI.HostController().User().DeleteUserConfig(srvData.ctx, appID, ID, srvData.header)
	if err != nil {
		blog.Errorf("DeleteUserCustomQuery http do error,err:%s, biz:%v:%+v,rid:%s", err.Error(), appID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("DeleteUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,rid:%s", result.Code, result.ErrMsg, appID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     nil,
	})

}

func (s *Service) GetUserCustomQuery(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)

	input := &meta.QueryInput{}
	if err := json.NewDecoder(req.Request.Body).Decode(input); nil != err {
		blog.Errorf("get user custom query failed with decode body err: %v,rid:%s", err, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
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
		blog.Error("GetUserCustomQuery query user custom query parameter ApplicationID not integer in url,bizID:%s,rid:%s", req.PathParameter("bk_biz_id"), srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, common.BKAppIDField)})
		return
	}
	input.Condition = condition

	result, err := s.CoreAPI.HostController().User().GetUserConfig(srvData.ctx, srvData.header, input)
	if err != nil {
		blog.Errorf("GetUserCustomQuery http do error,err:%s, biz:%v,input:%+v,rid:%s", err.Error(), req.PathParameter("bk_biz_id"), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustomQuery http response error,err code:%d,err msg:%s, bizID:%v,input:%+v,rid:%s", result.Code, result.ErrMsg, req.PathParameter("bk_biz_id"), input, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) GetUserCustomQueryDetail(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)

	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	result, err := s.CoreAPI.HostController().User().GetUserConfigDetail(srvData.ctx, appID, ID, srvData.header)
	if err != nil {
		blog.Errorf("GetUserCustomQueryDetail http do error,err:%s, biz:%v,ID:%+v,rid:%s", err.Error(), appID, ID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustomQueryDetail http response error,err code:%d,err msg:%s, bizID:%v,ID:%+v,rid:%s", result.Code, result.ErrMsg, appID, ID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	resp.WriteEntity(meta.Response{
		BaseResp: meta.SuccessBaseResp,
		Data:     result.Data,
	})

}

func (s *Service) GetUserCustomQueryResult(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)

	appID := req.PathParameter("bk_biz_id")
	ID := req.PathParameter("id")

	intAppID, err := util.GetInt64ByInterface(appID)
	if nil != err {
		blog.Errorf("UserAPIResult custom query failed,  err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedInt, "ApplicationID")})
		return
	}

	result, err := s.CoreAPI.HostController().User().GetUserConfigDetail(srvData.ctx, appID, ID, srvData.header)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("UserAPIResult custom query failed,  err: %v, appid:%s, id:%s", err.Error(), appID, ID)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrGetUserCustomQueryDetailFaild, err.Error())})
		return
	}

	if "" == result.Data.Name {
		blog.Errorf("UserAPIResult custom query not found, appid:%s, id:%s, logID:%s", appID, ID, srvData.rid)
		resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommNotFound)})
		return
	}

	var input meta.HostCommonSearch
	input.AppID = intAppID

	err = json.Unmarshal([]byte(result.Data.Info), &input)
	if nil != err {
		blog.Errorf("UserAPIResult custom unmarshal failed,  err: %v, appid:%s, id:%s, logID:%s", err.Error(), appID, ID, srvData.rid)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	input.Page.Start, _ = util.GetIntByInterface(req.PathParameter("start"))
	input.Page.Limit, _ = util.GetIntByInterface(req.PathParameter("limit"))

	retData, err := srvData.lgc.SearchHost(srvData.ctx, &input, false)
	if nil != err || (nil == err && !result.Result) {
		if nil == err {
			err = fmt.Errorf(result.ErrMsg)
		}
		blog.Errorf("UserAPIResult custom query search host failed,  err: %v, appid:%s, id:%s", err.Error(), appID, ID)
		resp.WriteError(http.StatusInternalServerError, &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrGetUserCustomQueryDetailFaild, err.Error())})
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
