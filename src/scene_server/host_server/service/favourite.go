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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

type infoParam struct {
	ExactSearch bool     `json:"exact_search"`
	InnerIP     bool     `json:"bk_host_innerip"`
	OuterIP     bool     `json:"bk_host_outerip"`
	IPList      []string `json:"ip_list"`
}

type queryParams []queryParam
type queryParam struct {
	ObjID    string      `json:"bk_obj_id"`
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
	OuterIP  bool        `json:"bk_host_outerip"`
	IPList   []string    `json:"ip_list"`
}

func (s *Service) ListHostFavourites(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	query := new(metadata.QueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(query); err != nil {
		blog.Errorf("get host favourite failed with decode body err: %v,rid:%s", err, srvData.rid)
		responseData := &metadata.RespError{
			Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed),
		}
		_ = resp.WriteError(http.StatusBadRequest, responseData)
		return
	}
	result, err := s.CoreAPI.CoreService().Host().ListHostFavourites(srvData.ctx, srvData.user, srvData.header, query)
	if err != nil {
		blog.Errorf("GetHostFavourites http do error,err:%s,input:%+v,rid:%s", err.Error(), query, srvData.rid)
		responseData := &metadata.RespError{
			Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed),
		}
		_ = resp.WriteError(http.StatusInternalServerError, responseData)
		return
	}
	if !result.Result {
		blog.Errorf("GetHostFavourites http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, query, srvData.rid)
		responseData := &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)}
		_ = resp.WriteError(http.StatusInternalServerError, responseData)
		return
	}

	_ = resp.WriteEntity(metadata.GetHostFavoriteResult{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) AddHostFavourite(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)

	param := new(metadata.FavouriteParms)
	if err := json.NewDecoder(req.Request.Body).Decode(param); err != nil {
		blog.Errorf("add host favourite failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if param.Name == "" {
		blog.Errorf("add host favorite, but got empty favorite name, param: %+v,rid:%s", param, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostEmptyFavName)})
		return
	}

	if param.Info != "" {
		// check if the info string matches the required structure
		err := json.Unmarshal([]byte(param.Info), &infoParam{})
		if err != nil {
			blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), param.Info, srvData.rid)
			_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "info")})
			return
		}
	}
	if param.QueryParams != "" {
		err := json.Unmarshal([]byte(param.QueryParams), &queryParams{})
		if err != nil {
			blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), param.QueryParams, srvData.rid)
			_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "query params")})
			return
		}
	}

	result, err := s.CoreAPI.CoreService().Host().AddHostFavourite(srvData.ctx, srvData.user, srvData.header, param)
	if err != nil {
		blog.Errorf("AddHostFavourite http do error,err:%s,input:%+v,rid:%s", err.Error(), param, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("AddHostFavourite http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, param, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	_ = resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	ID := req.PathParameter("id")

	if "" == ID || "0" == ID {
		blog.Errorf("update host favourite failed, with id  %d,rid:%s", ID, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update host favorite failed with decode body err: %v,rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if data["name"] == nil || data["name"].(string) == "" {
		blog.Errorf("update host favorite, but got empty name, data: %+v, rid:%s", data, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrHostFavUpdateFail)})
		return
	}

	if info, exists := data["info"]; exists {
		info := info.(string)
		if info != "" {
			// check if the info string matches the required structure
			err := json.Unmarshal([]byte(info), &infoParam{})
			if err != nil {
				blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), info, srvData.rid)
				_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "info")})
				return
			}
		}
	}
	if queryParam, exists := data["query_params"]; exists {
		queryParam := queryParam.(string)
		if queryParam != "" {
			// check if the info string matches the required structure
			err := json.Unmarshal([]byte(queryParam), &queryParams{})
			if err != nil {
				blog.Errorf("AddHostFavourite info unmarshal failed, err: %v, input:%+v, rid:%s", err.Error(), queryParam, srvData.rid)
				_ = resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, "query params")})
				return
			}
		}
	}

	result, err := s.CoreAPI.CoreService().Host().UpdateHostFavouriteByID(srvData.ctx, srvData.user, ID, srvData.header, data)
	if err != nil {
		blog.Errorf("UpdateHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("UpdateHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	ID := req.PathParameter("id")

	if "" == ID || "0" == ID {
		blog.Errorf("delete host favourite failed, with id  %d,rid:%s", ID, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	result, err := s.CoreAPI.CoreService().Host().DeleteHostFavouriteByID(srvData.ctx, srvData.user, ID, srvData.header)
	if err != nil {
		blog.Errorf("DeleteHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), ID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("DeleteHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, ID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) IncrHostFavouritesCount(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	ID := req.PathParameter("id")
	if "" == ID || "0" == ID {
		blog.Errorf("delete host favourite failed, with id  %s, rid:%s", ID, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	result, err := s.CoreAPI.CoreService().Host().GetHostFavouriteByID(srvData.ctx, srvData.user, ID, srvData.header)
	if err != nil {
		blog.Errorf("IncrHostFavouritesCount GetHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), ID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !result.Result {
		blog.Errorf("IncrHostFavouritesCount GetHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", result.Code, result.ErrMsg, ID, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(result.Code, result.ErrMsg)})
		return
	}

	count := result.Data.Count + 1
	data := map[string]interface{}{"count": count}
	uResult, err := s.CoreAPI.CoreService().Host().UpdateHostFavouriteByID(srvData.ctx, srvData.user, ID, srvData.header, data)
	if err != nil {
		blog.Errorf("IncrHostFavouritesCount UpdateHostFavouriteByID http do error,err:%s,input:%+v,rid:%s", err.Error(), data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommHTTPDoRequestFailed)})
		return
	}
	if !uResult.Result {
		blog.Errorf("IncrHostFavouritesCount UpdateHostFavouriteByID http response error,err code:%d,err msg:%s,input:%+v,rid:%s", uResult.Code, uResult.ErrMsg, data, srvData.rid)
		_ = resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: srvData.ccErr.New(uResult.Code, uResult.ErrMsg)})
		return
	}

	info := make(map[string]interface{})
	info["id"] = ID
	info["count"] = count
	_ = resp.WriteEntity(metadata.NewSuccessResp(info))
}
