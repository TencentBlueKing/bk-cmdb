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
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
)

func (s *Service) GetHostFavourites(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)
	query := new(metadata.QueryInput)
	if err := json.NewDecoder(req.Request.Body).Decode(query); err != nil {
		blog.Errorf("get host favourite failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	result, err := s.CoreAPI.HostController().Favorite().GetHostFavourites(context.Background(), user, pheader, query)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add host favorite failed, query: %v, err: %v, %v", query, err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostFavGetFail)})
		return
	}

	resp.WriteEntity(metadata.GetHostFavoriteResult{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) AddHostFavourite(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)

	para := new(metadata.FavouriteParms)
	if err := json.NewDecoder(req.Request.Body).Decode(para); err != nil {
		blog.Errorf("add host favourite failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if para.Name == "" {
		blog.Errorf("add host favorite, but got empty favorite name, param: %v", para)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostEmptyFavName)})
		return
	}

	result, err := s.CoreAPI.HostController().Favorite().AddHostFavourite(context.Background(), user, pheader, para)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("add host favorite failed, para: %v, err: %v, %v", para, err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostEmptyFavName)})
		return
	}

	resp.WriteEntity(metadata.Response{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}

func (s *Service) UpdateHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)
	ID := req.PathParameter("id")

	if "" == ID || "0" == ID {
		blog.Error("update host favourite failed, with id  %id", ID)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	data := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&data); err != nil {
		blog.Errorf("update host favorite failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	if data["name"] == nil || data["name"].(string) == "" {
		blog.Errorf("update host favorite, but got empty name, data: %v", data)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostFavUpdateFail)})
		return
	}

	result, err := s.CoreAPI.HostController().Favorite().UpdateHostFavouriteByID(context.Background(), user, ID, pheader, data)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("update host favorite failed, data: %v, err: %v, %v", data, err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostFavUpdateFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) DeleteHostFavouriteByID(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)
	ID := req.PathParameter("id")
	if "" == ID || "0" == ID {
		blog.Error("delete host favourite failed, with id  %id", ID)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	result, err := s.CoreAPI.HostController().Favorite().DeleteHostFavouriteByID(context.Background(), user, ID, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("delete host favorite failed, ID: %v, err: %v, %v", ID, err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostFavDeleteFail)})
		return
	}

	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) IncrHostFavouritesCount(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)
	ID := req.PathParameter("id")
	if "" == ID || "0" == ID {
		blog.Error("delete host favourite failed, with id  %id", ID)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPInputInvalid)})
		return
	}

	result, err := s.CoreAPI.HostController().Favorite().GetHostFavouriteByID(context.Background(), user, ID, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("increase host favorite count failed, ID: %v, err: %v, %v", ID, err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostFavGetFail)})
		return
	}

	count := result.Data.Count + 1
	data := map[string]interface{}{"count": count}
	uResult, err := s.CoreAPI.HostController().Favorite().UpdateHostFavouriteByID(context.Background(), user, ID, pheader, data)
	if err != nil || (err == nil && !uResult.Result) {
		blog.Errorf("increase host favorite count failed, ID: %v, err: %v, %v", ID, err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostFavUpdateFail)})
		return
	}

	info := make(map[string]interface{})
	info["id"] = ID
	info["count"] = count
	resp.WriteEntity(metadata.NewSuccessResp(info))
}
