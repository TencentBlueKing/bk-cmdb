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
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apiserver/core"
	compatiblev2 "configcenter/src/apiserver/core/compatiblev2/service"
	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
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
	WebServices(auth authcenter.AuthConfig) []*restful.WebService
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
	s.enableAuth = enableAuth
	s.engine = engine
	s.client = httpClient
	s.discovery = discovery
	s.core.CompatibleV2Operation().SetConfig(engine)
	s.authorizer = authorize
}

func (s *service) WebServices(auth authcenter.AuthConfig) []*restful.WebService {
	getErrFun := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	ws := &restful.WebService{}
	ws.Path(rootPath).Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON)
	if s.authorizer.Enabled() == true {
		ws.Filter(s.authFilter(getErrFun))
	}
	ws.Route(ws.POST("/auth/verify").To(s.AuthVerify))
	ws.Route(ws.GET("/auth/business-list").To(s.GetAnyAuthorizedAppList))
	ws.Route(ws.GET("/auth/admin-entrance").To(s.GetAdminEntrance))
	ws.Route(ws.GET("{.*}").Filter(s.URLFilterChan).To(s.Get))
	ws.Route(ws.POST("{.*}").Filter(s.URLFilterChan).To(s.Post))
	ws.Route(ws.PUT("{.*}").Filter(s.URLFilterChan).To(s.Put))
	ws.Route(ws.DELETE("{.*}").Filter(s.URLFilterChan).To(s.Delete))

	allWebServices := make([]*restful.WebService, 0)
	allWebServices = append(allWebServices, ws, s.core.CompatibleV2Operation().WebService())
	allWebServices = append(allWebServices, s.V3Healthz())
	return allWebServices
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

		if path == "/api/v3/auth/business-list" {
			fchain.ProcessFilter(req, resp)
			return
		}
		if path == "/api/v3/auth/admin-entrance" {
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

		blog.V(7).Infof("auth filter parse attribute result: %s, rid: %s", attribute, rid)
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
			permissions, err := authcenter.AdoptPermissions(attribute.Resources)
			if err != nil {
				rsp := metadata.BaseResp{
					Code:   common.CCErrCommCheckAuthorizeFailed,
					ErrMsg: errFunc().CreateDefaultCCErrorIf(language).Error(common.CCErrCommCheckAuthorizeFailed).Error(),
					Result: false,
				}
				resp.WriteAsJson(rsp)
				return
			}
			blog.Warnf("authFilter failed, url: %s, reason: %s, rid: %s", path, decision.Reason, rid)
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
