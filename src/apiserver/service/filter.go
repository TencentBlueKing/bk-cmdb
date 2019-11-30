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
	"fmt"
	"net/http"
	"strings"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/parser"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

type RequestType string

const (
	UnknownType     RequestType = "unknown"
	TopoType        RequestType = "topo"
	HostType        RequestType = "host"
	ProcType        RequestType = "proc"
	EventType       RequestType = "event"
	DataCollectType RequestType = "collect"
	OperationType   RequestType = "operation"
	TaskType        RequestType = "task"
)

func (s *service) URLFilterChan(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	rid := util.GetHTTPCCRequestID(req.Request.Header)

	var kind RequestType
	var err error
	kind, err = URLPath(req.Request.RequestURI).FilterChain(req)
	if err != nil {
		blog.Errorf("rewrite request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
		if err := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
			Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
			ErrCode: common.CCErrRewriteRequestUriFailed,
			Data:    nil,
		}); err != nil {
			blog.Errorf("response request[url: %s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			return
		}
		return
	}

	defer func() {
		if err != nil {
			blog.Errorf("proxy request url[%s] failed, err: %v, rid: %s", req.Request.RequestURI, err, rid)
			if rerr := resp.WriteError(http.StatusInternalServerError, &metadata.RespError{
				Msg:     fmt.Errorf("rewrite request failed, %s", err.Error()),
				ErrCode: common.CCErrRewriteRequestUriFailed,
				Data:    nil,
			}); rerr != nil {
				blog.Errorf("proxy request[url: %s] failed, err: %v, rid: %s", req.Request.RequestURI, rerr, rid)
				return
			}
			return
		}
	}()

	servers := make([]string, 0)
	switch kind {
	case TopoType:
		servers, err = s.discovery.TopoServer().GetServers()

	case ProcType:
		servers, err = s.discovery.ProcServer().GetServers()

	case EventType:
		servers, err = s.discovery.EventServer().GetServers()

	case HostType:
		servers, err = s.discovery.HostServer().GetServers()

	case DataCollectType:
		servers, err = s.discovery.DataCollect().GetServers()

	case OperationType:
		servers, err = s.discovery.OperationServer().GetServers()
	case TaskType:
		servers, err = s.discovery.TaskServer().GetServers()
	}

	if err != nil {
		return
	}

	if strings.HasPrefix(servers[0], "https://") {
		req.Request.URL.Host = servers[0][8:]
		req.Request.URL.Scheme = "https"
	} else {
		req.Request.URL.Host = servers[0][7:]
		req.Request.URL.Scheme = "http"
	}

	chain.ProcessFilter(req, resp)
}

func (s *service) authFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		rid := util.GetHTTPCCRequestID(req.Request.Header)
		path := req.Request.URL.Path

		blog.V(7).Infof("authFilter on url: %s, rid: %s", path, rid)
		if s.authorizer.Enabled() == false {
			blog.V(7).Infof("auth disabled, skip auth filter, rid: %s", rid)
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/verify" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/business_list" {
			fchain.ProcessFilter(req, resp)
			return
		}
		if path == "/api/v3/auth/admin_entrance" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/skip_url" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if path == "/api/v3/auth/convert" {
			fchain.ProcessFilter(req, resp)
			return
		}

		// if common.BKSuperOwnerID == util.GetOwnerID(req.Request.Header) {
		// 	blog.Errorf("authFilter failed, can not use super supplier account, rid: %s", rid)
		// 	rsp := metadata.BaseResp{
		// 		Code:   common.CCErrCommParseAuthAttributeFailed,
		// 		ErrMsg: "invalid supplier account.",
		// 		Result: false,
		// 	}
		// 	resp.WriteHeaderAndJson(http.StatusBadRequest, rsp, restful.MIME_JSON)
		// 	return
		// }

		language := util.GetLanguage(req.Request.Header)
		attribute, err := parser.ParseAttribute(req, s.engine)
		if err != nil {
			blog.Errorf("authFilter failed, caller: %s, parse auth attribute for %s %s failed, err: %v, rid: %s", req.Request.RemoteAddr, req.Request.Method, req.Request.URL.Path, err, rid)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommParseAuthAttributeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommParseAuthAttributeFailed).Error(),
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		// check if authorize is nil or not, which means to check if the authorize instance has
		// already been initialized or not. if not, api server should not be used.
		if nil == s.authorizer {
			blog.Errorf("authorize instance has not been initialized, rid: %s", rid)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommCheckAuthorizeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		blog.V(7).Infof("auth filter parse attribute result: %v, rid: %s", attribute, rid)
		decision, err := s.authorizer.Authorize(req.Request.Context(), attribute)
		if err != nil {
			blog.Errorf("authFilter failed, authorized request failed, url: %s, err: %v, rid: %s", path, err, rid)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommCheckAuthorizeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
				Result: false,
			}
			resp.WriteAsJson(rsp)
			return
		}

		if !decision.Authorized {
			blog.V(4).Infof("authcenter.AdoptPermissions attribute: %+v, rid: %s", attribute, rid)
			permissions, err := authcenter.AdoptPermissions(req.Request.Header, s.engine.CoreAPI, attribute.Resources)
			if err != nil {
				blog.Errorf("adopt permission failed, err: %v, rid: %s", err, rid)
				rsp := metadata.BaseResp{
					Code:   common.CCErrCommCheckAuthorizeFailed,
					ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
					Result: false,
				}
				resp.WriteAsJson(rsp)
				return
			}
			blog.Warnf("authFilter failed, url: %s, reason: %+v, permissions: %+v, rid: %s", path, decision, permissions, rid)
			rsp := metadata.BaseResp{
				Code:        9900403,
				ErrMsg:      errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommAuthNotHavePermission).Error(),
				Result:      false,
				Permissions: permissions,
			}
			resp.WriteAsJson(rsp)
			return
		}

		blog.V(7).Infof("authFilter authorize on url:%s success, rid: %s", path, rid)
		fchain.ProcessFilter(req, resp)
		return
	}
}
