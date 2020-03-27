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
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
)

func (s *Service) SaveUserCustom(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	params := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("save user custom failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result, err := s.CoreAPI.CoreService().Host().GetUserCustomByUser(srvData.ctx, srvData.user, srvData.header)
	if err != nil {
		blog.Errorf("SaveUserCustom GetUserCustomByUser http do error,err:%s,input:%s, rid:%s", err.Error(), params, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("SaveUserCustom GetUserCustomByUser http response error,err code:%d,err msg:%s,input:%s, rid:%s", result.Code, result.ErrMsg, params, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	if len(result.Data) == 0 {
		result, err := s.CoreAPI.CoreService().Host().AddUserCustom(srvData.ctx, srvData.user, srvData.header, params)
		if err != nil {
			blog.Errorf("SaveUserCustom AddUserCustom http do error,err:%s,input:%s, rid:%s", err.Error(), params, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
			return
		}
		if !result.Result {
			blog.Errorf("SaveUserCustom AddUserCustom http response error,err code:%d,err msg:%s,input:%s, rid:%s", result.Code, result.ErrMsg, params, srvData.rid)
			_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
			return
		}
		_ = resp.WriteEntity(result)
		return

	}
	id := result.Data["id"].(string)
	uResult, err := s.CoreAPI.CoreService().Host().UpdateUserCustomByID(srvData.ctx, srvData.user, id, srvData.header, params)
	if err != nil {
		blog.Errorf("SaveUserCustom UpdateUserCustomByID http do error,err:%s,input:%s, rid:%s", err.Error(), params, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !uResult.Result {
		blog.Errorf("SaveUserCustom UpdateUserCustomByID http response error,err code:%d,err msg:%s,input:%s, rid:%s", uResult.Code, uResult.ErrMsg, params, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(uResult.Code, uResult.ErrMsg)})
		return
	}

	_ = resp.WriteEntity(uResult)
}

func (s *Service) GetUserCustom(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	result, err := s.CoreAPI.CoreService().Host().GetUserCustomByUser(srvData.ctx, srvData.user, srvData.header)
	if err != nil {
		blog.Errorf("GetUserCustom http do error,err:%s, rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("GetUserCustom http response error,err code:%d,err msg:%s, rid:%s", result.Code, result.ErrMsg, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}
	_ = resp.WriteEntity(result)
}

func (s *Service) GetDefaultCustom(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	result, err := s.CoreAPI.CoreService().Host().GetDefaultUserCustom(srvData.ctx, srvData.user, srvData.header)
	if err != nil {
		blog.Errorf("GetDefaultCustom http do error,err:%s, rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("GetDefaultCustom http response error,err code:%d,err msg:%s, rid:%s", result.Code, result.ErrMsg, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}
	_ = resp.WriteEntity(metadata.GetUserCustomResult{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}
