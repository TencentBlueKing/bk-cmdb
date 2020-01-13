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

// GetModelDefaultCustom 获取模型在列表页面展示字段
func (s *Service) GetModelDefaultCustom(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	result, err := s.CoreAPI.CoreService().Host().GetDefaultUserCustom(srvData.ctx, srvData.header)
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
	// ensure return {} by json decode
	if result.Data == nil {
		result.Data = make(map[string]interface{}, 0)
	}
	_ = resp.WriteEntity(metadata.GetUserCustomResult{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}

// SaveModelDefaultCustom 设置模型在列表页面展示字段
func (s *Service) SaveModelDefaultCustom(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	objID := req.PathParameter("obj_id")

	input := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&input); err != nil {
		blog.Errorf("save user custom failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if len(input) == 0 {
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPBodyEmpty)})
		return
	}

	userCustomInput := make(map[string]interface{}, 0)
	// add prefix all key
	for key, val := range input {
		userCustomInput[fmt.Sprintf("%s_%s", objID, key)] = val
	}

	result, err := s.CoreAPI.CoreService().Host().UpdateDefaultUserCustom(srvData.ctx, srvData.header, userCustomInput)
	if err != nil {
		blog.ErrorJSON("SaveUserCustom GetUserCustomByUser http do error,err:%s,input:%s, rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if err := result.CCError(); err != nil {
		blog.ErrorJSON("SaveUserCustom GetUserCustomByUser http reply error. result: %s, input: %s, rid: %s", result, input, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}

	_ = resp.WriteEntity(result)
}
