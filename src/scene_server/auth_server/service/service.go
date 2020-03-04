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
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/rdapi"
	"configcenter/src/scene_server/auth_server/logics"

	"github.com/emicklei/go-restful"
)

type AuthService struct {
	engine *backbone.Engine
	auth   ac.AuthInterface
	lgc    *logics.Logics
}

func NewAuthService(engine *backbone.Engine, auth ac.AuthInterface) *AuthService {
	return &AuthService{
		engine: engine,
		auth:   auth,
		lgc:    logics.NewLogics(engine.CoreAPI),
	}
}

func (s *AuthService) checkRequestFromIamFilter() func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		if !auth.IsAuthed() {
			chain.ProcessFilter(req, resp)
			return
		}

		isAuthorized, err := s.auth.CheckRequestAuthorization(req.Request)
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			rsp := IamBaseResp{
				Code:    common.CCErrCommHTTPDoRequestFailed,
				Message: err.Error(),
			}
			_ = resp.WriteAsJson(rsp)
			return
		}
		if !isAuthorized {
			resp.WriteHeader(http.StatusUnauthorized)
			rsp := IamBaseResp{
				Code:    common.CCErrCommAuthNotHavePermission,
				Message: "request not from iam",
			}
			_ = resp.WriteAsJson(rsp)
			return
		}

		chain.ProcessFilter(req, resp)
		return
	}
}

func (s *AuthService) WebService() *restful.Container {
	api := new(restful.WebService)
	api.Path("/auth/v3")
	api.Filter(s.engine.Metric().RestfulMiddleWare)
	api.Filter(rdapi.AllGlobalFilter(func() errors.CCErrorIf {
		return s.engine.CCErr
	}))
	// only allows iam to pull resource using these api
	api.Filter(s.checkRequestFromIamFilter())
	api.Produces(restful.MIME_JSON)

	s.initResourcePull(api)

	container := restful.NewContainer()
	container.Add(api)

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

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/auth/find/empty/resource", Handler: s.PullNoRelatedInstanceResource})

	utility.AddToRestfulWebService(api)
}
