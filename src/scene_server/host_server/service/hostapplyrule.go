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
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	meta "configcenter/src/common/metadata"

	"github.com/emicklei/go-restful"
)

func (s *Service) CreateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("CreateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.CreateHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("CreateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().CreateHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("CreateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(rule))
}

func (s *Service) UpdateHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleIDStr := req.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("UpdateHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.UpdateHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("UpdateHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().UpdateHostApplyRule(srvData.ctx, srvData.header, bizID, ruleID, option)
	if err != nil {
		blog.ErrorJSON("UpdateHostApplyRule failed, core service CreateHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(rule))
}

func (s *Service) DeleteHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("DeleteHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.DeleteHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("DeleteHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	if err := s.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(srvData.ctx, srvData.header, bizID, option); err != nil {
		blog.ErrorJSON("DeleteHostApplyRule failed, core service DeleteHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(make(map[string]interface{})))
}

func (s *Service) GetHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleIDStr := req.PathParameter(common.HostApplyRuleIDField)
	ruleID, err := strconv.ParseInt(ruleIDStr, 10, 64)
	if err != nil {
		blog.Errorf("GetHostApplyRule failed, parse biz id failed, ruleIDStr: %s, err: %v,rid:%s", ruleIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.HostApplyRuleIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	rule, err := s.CoreAPI.CoreService().HostApplyRule().GetHostApplyRule(srvData.ctx, srvData.header, bizID, ruleID)
	if err != nil {
		blog.ErrorJSON("GetHostApplyRule failed, core service GetHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(rule))
}

func (s *Service) ListHostApplyRule(req *restful.Request, resp *restful.Response) {
	srvData := s.newSrvComm(req.Request.Header)
	rid := srvData.rid

	bizIDStr := req.PathParameter(common.BKAppIDField)
	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("ListHostApplyRule failed, parse biz id failed, bizIDStr: %s, err: %v,rid:%s", bizIDStr, err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Errorf(common.CCErrCommParamsInvalid, common.BKAppIDField)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	option := meta.ListHostApplyRuleOption{}
	if err := json.NewDecoder(req.Request.Body).Decode(&option); err != nil {
		blog.Errorf("ListHostApplyRule failed, decode request body failed, err: %v,rid:%s", err, rid)
		result := &meta.RespError{Msg: srvData.ccErr.Error(common.CCErrCommJSONUnmarshalFailed)}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	ruleResult, err := s.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(srvData.ctx, srvData.header, bizID, option)
	if err != nil {
		blog.ErrorJSON("ListHostApplyRule failed, core service ListHostApplyRule failed, bizID: %s, option: %s, err: %s, rid: %s", bizID, option, err.Error(), rid)
		result := &meta.RespError{Msg: err}
		_ = resp.WriteError(http.StatusBadRequest, result)
		return
	}

	_ = resp.WriteEntity(meta.NewSuccessResp(ruleResult))
}
