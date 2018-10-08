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
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	input := new(metadata.HostLockRequest)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("lock host, but decode body failed, err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}

	err := s.Logics.LockHost(context.Background(), req.Request.Header, input)
	if nil != err {
		blog.Errorf("lock host, lock host handle failed, err: %s, input:%+v, logID:%s", err.Error(), input, util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}

	resp.WriteEntity(metadata.HostLockResponse{
		BaseResp: metadata.SuccessBaseResp,
	})

}

func (s *Service) UnlockHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	input := new(metadata.HostLockRequest)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("unlock host, but decode body failed, err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	err := s.Logics.UnlockHost(context.Background(), req.Request.Header, input)
	if nil != err {
		blog.Errorf("unlock host, unlock host handle failed, err: %s, input:%+v, logID:%s", err.Error(), input, util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}

	resp.WriteEntity(metadata.HostLockResponse{
		BaseResp: metadata.SuccessBaseResp,
	})

}

func (s *Service) QueryLockHost(req *restful.Request, resp *restful.Response) {
	pheader := req.Request.Header
	defErr := s.Core.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(pheader))

	input := new(metadata.QueryHostLockRequest)
	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("unlock host, but decode body failed, err: %s", err.Error())
		resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: defErr.Error(common.CCErrCommHTTPReadBodyFailed)})
		return
	}
	hostLockArr, err := s.Logics.QueryHostLock(context.Background(), req.Request.Header, input)
	if nil != err {
		blog.Errorf("unlock host, unlock host handle failed, err: %s, input:%+v, logID:%s", err.Error(), input, util.GetHTTPCCRequestID(req.Request.Header))
		resp.WriteError(http.StatusInternalServerError, &metadata.RespError{Msg: err})
		return
	}
	result := metadata.HostLockQueryResponse{
		BaseResp: metadata.SuccessBaseResp,
	}
	result.Data.Info = hostLockArr
	result.Data.Count = int64(len(hostLockArr))
	resp.WriteEntity(result)

}
