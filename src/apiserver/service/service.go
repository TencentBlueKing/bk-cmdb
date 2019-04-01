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
	"net/http"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apiserver/core"
	compatiblev2 "configcenter/src/apiserver/core/compatiblev2/service"
	"configcenter/src/auth"
	"configcenter/src/auth/parser"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
)

// Service service methods
type Service interface {
	WebServices() []*restful.WebService
	SetConfig(enableAuth bool, engine *backbone.Engine, httpClient HTTPClient, discovery discovery.DiscoveryInterface, authorize auth.Authorize)
}

// NewService create a new service instance
func NewService() Service {
	return &service{
		core: core.New(nil, compatiblev2.New(nil)),
	}
}

type service struct {
	enableAuth bool
	engine     *backbone.Engine
	client     HTTPClient
	core       core.Core
	discovery  discovery.DiscoveryInterface
	authorizer auth.Authorizer
}

func (s *service) SetConfig(enableAuth bool, engine *backbone.Engine, httpClient HTTPClient, discovery discovery.DiscoveryInterface, authorize auth.Authorize) {
	s.enableAuth = false
	s.enableAuth = enableAuth
	s.engine = engine
	s.client = httpClient
	s.discovery = discovery
	s.core.CompatibleV2Operation().SetConfig(engine)
	s.authorizer = authorize
}

func (s *service) WebServices() []*restful.WebService {
	getErrFun := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	ws := &restful.WebService{}
	ws.Path(rootPath).Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Filter(s.authFilter(getErrFun))
	ws.Route(ws.POST("/auth/verify").To(s.AuthVerify))
	ws.Route(ws.GET("/auth/business-list").To(s.GetAuthorizedAppList))
	ws.Route(ws.GET("{.*}").Filter(s.URLFilterChan).To(s.Get))
	ws.Route(ws.POST("{.*}").Filter(s.URLFilterChan).To(s.Post))
	ws.Route(ws.PUT("{.*}").Filter(s.URLFilterChan).To(s.Put))
	ws.Route(ws.DELETE("{.*}").Filter(s.URLFilterChan).To(s.Delete))

	allWebServices := make([]*restful.WebService, 0)
	allWebServices = append(allWebServices, ws, s.core.CompatibleV2Operation().WebService())

	return allWebServices
}

func (s *service) authFilter(errFunc func() errors.CCErrorIf) func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, fchain *restful.FilterChain) {
		if req.Request.URL.Path == "/api/v3/auth/verify" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if req.Request.URL.Path == "/api/v3/auth/business-list" {
			fchain.ProcessFilter(req, resp)
			return
		}

		if common.BKSuperOwnerID == util.GetOwnerID(req.Request.Header) {
			blog.Errorf("request id: %s, can not use super supplier account", util.GetHTTPCCRequestID(req.Request.Header))
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommParseAuthAttributeFailed,
				ErrMsg: "invalid supplier account.",
				Result: false,
			}
			resp.WriteHeaderAndJson(http.StatusBadRequest, rsp, restful.MIME_JSON)
			return
		}

		if !s.enableAuth {
			fchain.ProcessFilter(req, resp)
			return
		}

		language := util.GetLanguage(req.Request.Header)
		attribute, err := parser.ParseAttribute(req, s.engine)
		if err != nil {
			blog.Errorf("request id: %s, parse auth attribute for %s %s failed, err: %v", util.GetHTTPCCRequestID(req.Request.Header), req.Request.Method, req.Request.URL.Path, err)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommParseAuthAttributeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommParseAuthAttributeFailed).Error(),
				Result: false,
			}
			resp.WriteHeaderAndJson(http.StatusBadRequest, rsp, restful.MIME_JSON)
			return
		}

		// check if authorize is nil or not, which means to check if the authorize instance has
		// already been initialized or not. if not, api server should not be used.
		if nil == s.authorizer {
			blog.Error("authorize instance has not been initialized")
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommCheckAuthorizeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
				Result: false,
			}
			resp.WriteHeaderAndJson(http.StatusInternalServerError, rsp, restful.MIME_JSON)
		}

		blog.InfoJSON("attr: %s", attribute)
		decision, err := s.authorizer.Authorize(req.Request.Context(), attribute)
		if err != nil {
			blog.Errorf("request id: %s, url: %s, authorized failed, because authorize this request failed, err: %v", util.GetHTTPCCRequestID(req.Request.Header), req.Request.URL.Path, err)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommCheckAuthorizeFailed,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
				Result: false,
			}
			resp.WriteHeaderAndJson(http.StatusInternalServerError, rsp, restful.MIME_JSON)
			return
		}

		if !decision.Authorized {
			blog.Errorf("request id: %s, url: %s, auth failed. reason: %s ", util.GetHTTPCCRequestID(req.Request.Header), req.Request.URL.Path, decision.Reason)
			rsp := metadata.BaseResp{
				Code:   common.CCErrCommAuthNotHavePermission,
				ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommAuthNotHavePermission).Error(),
				Result: false,
			}
			resp.WriteHeaderAndJson(http.StatusForbidden, rsp, restful.MIME_JSON)
			return
		}

		fchain.ProcessFilter(req, resp)
		return
	}
}
