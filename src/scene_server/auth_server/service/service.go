/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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

	"configcenter/src/ac"
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/auth_server/logics"
	sdkauth "configcenter/src/scene_server/auth_server/sdk/auth"
	"configcenter/src/scene_server/auth_server/types"

	"github.com/emicklei/go-restful"
)

type AuthService struct {
	engine     *backbone.Engine
	auth       ac.AuthInterface
	lgc        *logics.Logics
	authorizer sdkauth.Authorizer
}

func NewAuthService(engine *backbone.Engine, auth ac.AuthInterface, lgc *logics.Logics, authorizer sdkauth.Authorizer) *AuthService {
	return &AuthService{
		engine:     engine,
		auth:       auth,
		lgc:        lgc,
		authorizer: authorizer,
	}
}

func (s *AuthService) checkRequestFromIamFilter() func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		if !auth.EnableAuthorize() {
			chain.ProcessFilter(req, resp)
			return
		}

		isAuthorized, err := s.auth.CheckRequestAuthorization(req.Request)
		if err != nil {
			rsp := types.BaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			}
			_ = resp.WriteAsJson(rsp)
			return
		}
		if !isAuthorized {
			rsp := types.BaseResp{
				Code:    types.UnauthorizedErrorCode,
				Message: "request not from iam",
			}
			_ = resp.WriteAsJson(rsp)
			return
		}

		// use iam request id as cc rid
		rid := req.Request.Header.Get(iam.IamRequestHeader)
		resp.Header().Set(iam.IamRequestHeader, rid)
		if rid != "" {
			req.Request.Header.Set(common.BKHTTPCCRequestID, rid)
		} else if rid = util.GetHTTPCCRequestID(req.Request.Header); rid == "" {
			rid = util.GenerateRID()
			req.Request.Header.Set(common.BKHTTPCCRequestID, rid)
		}
		resp.Header().Set(common.BKHTTPCCRequestID, rid)

		// use iam language as cc language
		req.Request.Header.Set(common.BKHTTPLanguage, req.Request.Header.Get("Blueking-Language"))

		chain.ProcessFilter(req, resp)
		return
	}
}

func (s *AuthService) WebService() *restful.Container {
	api := new(restful.WebService)
	api.Path("/auth/v3")
	api.Filter(s.engine.Metric().RestfulMiddleWare)
	// only allows iam to pull resource using these api
	api.Filter(s.checkRequestFromIamFilter())
	api.Produces(restful.MIME_JSON)

	s.initResourcePull(api)

	container := restful.NewContainer()
	container.Add(api)

	authAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	authAPI.Path("/ac/v3")
	s.initAuth(authAPI)
	container.Add(authAPI)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.Healthz))
	container.Add(healthzAPI)

	return container
}

func (s *AuthService) initResourcePull(api *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/empty/resource", Handler: s.PullNoRelatedInstanceResource})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/instance/resource", Handler: s.PullInstanceResource})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/system/resource", Handler: s.PullSystemResource})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/business/resource", Handler: s.PullBusinessResource})

	utility.AddToRestfulWebService(api)
}

func (s *AuthService) initAuth(api *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/authorize", Handler: s.Authorize})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/authorize/batch", Handler: s.AuthorizeBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/authorize/any/batch", Handler: s.AuthorizeAnyBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/authorized_resource", Handler: s.ListAuthorizedResources})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/no_auth_skip_url", Handler: s.GetNoAuthSkipUrl})

	utility.AddToRestfulWebService(api)
}
