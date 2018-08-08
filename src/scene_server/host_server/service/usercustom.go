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

func (s *Service) SaveUserCustom(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)

	params := make(map[string]interface{})
	if err := json.NewDecoder(req.Request.Body).Decode(&params); err != nil {
		blog.Errorf("save user custom failed with decode body err: %v", err)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	result, err := s.CoreAPI.HostController().User().GetUserCustomByUser(context.Background(), user, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Error("save user custom, but get user custom failed, err: %v, %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CC_Err_Comm_USER_CUSTOM_SAVE_FAIL)})
		return
	}

	if len(result.Data) == 0 {
		result, err := s.CoreAPI.HostController().User().AddUserCustom(context.Background(), user, pheader, params)
		if err != nil || (err == nil && !result.Result) {
			blog.Errorf("save user custom, add failed, err: %v, %v", err, result.ErrMsg)
			resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CC_Err_Comm_USER_CUSTOM_SAVE_FAIL)})
			return
		} else {
			resp.WriteEntity(result)
			return
		}
	}
	id := result.Data["id"].(string)
	uResult, err := s.CoreAPI.HostController().User().UpdateUserCustomByID(context.Background(), user, id, pheader, params)
	if err != nil || (err == nil && !uResult.Result) {
		blog.Errorf("save user custom, but update failed, err: %v, %v", err, uResult.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CC_Err_Comm_USER_CUSTOM_SAVE_FAIL)})
		return
	}

	resp.WriteEntity(uResult)
}

func (s *Service) GetUserCustom(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)
	result, err := s.CoreAPI.HostController().User().GetUserCustomByUser(context.Background(), user, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Error("get user custom, but get user custom failed, err: %v, %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CC_Err_Comm_USER_CUSTOM_SAVE_FAIL)})
		return
	}
	resp.WriteEntity(result)
}

func (s *Service) GetDefaultCustom(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))
	user := util.GetUser(pheader)

	result, err := s.CoreAPI.HostController().User().GetDefaultUserCustom(context.Background(), user, pheader)
	if err != nil || (err == nil && !result.Result) {
		blog.Errorf("get default user custom failed, err: %v, %v", err, result.ErrMsg)
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrHostCustomGetDefaultFail)})
		return
	}
	resp.WriteEntity(metadata.GetUserCustomResult{
		BaseResp: metadata.SuccessBaseResp,
		Data:     result.Data,
	})
}
