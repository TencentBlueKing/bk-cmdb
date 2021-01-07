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

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

func (s *Service) LockHost(req *restful.Request, resp *restful.Response) {

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	input := &metadata.HostLockRequest{}

	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("lock host , but decode body failed, err: %s, logID:%s", err.Error(), util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == len(input.IPS) {
		blog.Errorf("lock host, ip_list is empty, logID:%s", util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "ip_list")})
		return
	}
	err := s.Logics.LockHost(context.Background(), req.Request.Header, input)
	if nil != err {
		blog.Errorf("lock host, handle host lock error, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) UnlockHost(req *restful.Request, resp *restful.Response) {

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	input := &metadata.HostLockRequest{}

	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("unlock host , but decode body failed, err: %s, logID:%s", err.Error(), util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == len(input.IPS) {
		blog.Errorf("unlock host, ip_list is empty, logID:%s", util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "ip_list")})
		return
	}
	err := s.Logics.UnlockHost(context.Background(), req.Request.Header, input)
	if nil != err {
		blog.Errorf("unlock host, handle host unlock error, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) QueryHostLock(req *restful.Request, resp *restful.Response) {

	defErr := s.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(req.Request.Header))
	input := &metadata.QueryHostLockRequest{}

	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("query lock host , but decode body failed, err: %s, logID:%s", err.Error(), util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == len(input.IPS) {
		blog.Errorf("query lock host, ip_list is empty, logID:%s", util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Errorf(common.CCErrCommParamsNeedSet, "ip_list")})
		return
	}

	hostLockInfos, err := s.Logics.QueryHostLock(context.Background(), req.Request.Header, input)
	if nil != err {
		blog.Errorf("query lock host, handle query host lock error, error:%s, logID:%s", err.Error(), util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	resp.WriteEntity(metadata.HostLockResultResponse{
		BaseResp: metadata.SuccessBaseResp,
		Data:     hostLockInfos,
	})
}
