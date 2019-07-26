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
	meta "configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"github.com/emicklei/go-restful"
)

func (s *Service) FindModuleHost(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	defErr := srvData.ccErr

	body := new(meta.HostModuleFind)
	if err := json.NewDecoder(req.Request.Body).Decode(body); err != nil {
		blog.Errorf("find host failed with decode body err: %#v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	host, err := srvData.lgc.FindHostByModuleIDs(srvData.ctx, body, false)
	if err != nil {
		blog.Errorf("find host failed, err: %#v, input:%#v, rid:%s", err, body, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	_ = resp.WriteEntity(meta.SearchHostResult{
		BaseResp: meta.SuccessBaseResp,
		Data:     *host,
	})
}

func (s *Service) ListHostByTopoNode(req *restful.Request, resp *restful.Response) {
	header := req.Request.Header
	ctx := util.NewContextFromHTTPHeader(header)
	rid := util.ExtractRequestIDFromContext(ctx)
	srvData := s.newSrvComm(header)
	defErr := srvData.ccErr

	parameter := &meta.ListHostByTopoNodeParameter{}
	if err := json.NewDecoder(req.Request.Body).Decode(parameter); err != nil {
		blog.Errorf("ListHostByTopoNode failed, decode body failed, err: %#v, rid:%s", err, srvData.rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}
	bizID, err := meta.BizIDFromMetadata(parameter.Metadata)
	if err != nil {
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrCommJSONUnmarshalFailed)})
		return
	}

	s.CoreAPI.CoreService()
	option := meta.ListHostByTopoNodeOption{
		BizID:     bizID,
		SetIDs:    parameter.SetIDs,
		ModuleIDs: parameter.ModuleIDs,
	}
	host, err := s.CoreAPI.CoreService().Host().ListHostByTopoNode(ctx, header, option)
	if err != nil {
		blog.Errorf("find host failed, err: %s, input:%#v, rid:%s", err.Error(), parameter, rid)
		_ = resp.WriteError(http.StatusBadRequest, &meta.RespError{Msg: defErr.Error(common.CCErrHostGetFail)})
		return
	}

	result := meta.NewSuccessResponse(host)
	_ = resp.WriteEntity(result)
}
