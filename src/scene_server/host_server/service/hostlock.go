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

	"configcenter/src/ac"
	"configcenter/src/ac/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

func (s *Service) LockHost(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)
	input := &metadata.HostLockRequest{}

	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("lock host , but decode body failed, err: %s, rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == len(input.IDS) {
		blog.Errorf("lock host, id_list is empty,input:%+v, rid:%s", input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "id_list")})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, meta.Update, input.IDS...); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v", input.IDS, err)
			resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(srvData.ctx, srvData.header, input.IDS)
		if err != nil {
			resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	err := srvData.lgc.LockHost(srvData.ctx, input)
	if nil != err {
		blog.Errorf("lock host, handle host lock error, error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) UnlockHost(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)
	input := &metadata.HostLockRequest{}

	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("unlock host , but decode body failed, err: %s, rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == len(input.IDS) {
		blog.Errorf("unlock host, id_list is empty, input:%+v,rid:%s", input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "id_list")})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, meta.Update, input.IDS...); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", input.IDS, err, srvData.rid)
			_ = resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(srvData.ctx, srvData.header, input.IDS)
		if err != nil {
			blog.Errorf("gen no permission response failed, err: %v, rid: %s", err, srvData.rid)
			resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	err := srvData.lgc.UnlockHost(srvData.ctx, input)
	if nil != err {
		blog.Errorf("unlock host, handle host unlock error, error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}
	_ = resp.WriteEntity(metadata.NewSuccessResp(nil))
}

func (s *Service) QueryHostLock(req *restful.Request, resp *restful.Response) {

	srvData := s.newSrvComm(req.Request.Header)
	input := &metadata.QueryHostLockRequest{}

	if err := json.NewDecoder(req.Request.Body).Decode(input); err != nil {
		blog.Errorf("query lock host , but decode body failed, err: %s, rid:%s", err.Error(), srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	if 0 == len(input.IDS) {
		blog.Errorf("query lock host, id_list is empty, input:%+v,rid:%s", input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsNeedSet, "id_list")})
		return
	}

	// auth: check authorization
	if err := s.AuthManager.AuthorizeByHostsIDs(srvData.ctx, srvData.header, meta.Update, input.IDS...); err != nil {
		if err != ac.NoAuthorizeError {
			blog.Errorf("check host authorization failed, hosts: %+v, err: %v, rid: %s", input.IDS, err, srvData.rid)
			_ = resp.WriteError(http.StatusForbidden, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		perm, err := s.AuthManager.GenEditBizHostNoPermissionResp(srvData.ctx, srvData.header, input.IDS)
		if err != nil {
			blog.Errorf("gen no permission response failed, err: %v, rid: %s", err, srvData.rid)
			resp.WriteError(http.StatusOK, &metadata.RespError{Msg: srvData.ccErr.Error(common.CCErrCommAuthorizeFailed)})
			return
		}
		resp.WriteEntity(perm)
		return
	}

	hostLockInfos, err := srvData.lgc.QueryHostLock(srvData.ctx, input)
	if nil != err {
		blog.Errorf("query lock host, handle query host lock error, error:%s, input:%+v,rid:%s", err.Error(), input, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &metadata.RespError{Msg: err})
		return
	}

	_ = resp.WriteEntity(metadata.HostLockResultResponse{
		BaseResp: metadata.SuccessBaseResp,
		Data:     hostLockInfos,
	})
}
