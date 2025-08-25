/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package service TODO
package service

import (
	"context"
	"net/http"
	"time"

	"configcenter/pkg/tenant/tools"
	iamtypes "configcenter/src/ac/iam/types"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/auth_server/app/options"
	"configcenter/src/scene_server/auth_server/logics"
	sdkauth "configcenter/src/scene_server/auth_server/sdk/auth"
	"configcenter/src/scene_server/auth_server/sdk/client"
	"configcenter/src/scene_server/auth_server/types"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// AuthService TODO
type AuthService struct {
	engine     *backbone.Engine
	iamClient  client.Interface
	lgc        *logics.Logics
	authorizer sdkauth.Authorizer
	config     *options.Config
}

// NewAuthService TODO
func NewAuthService(engine *backbone.Engine, iamClient client.Interface, lgc *logics.Logics,
	authorizer sdkauth.Authorizer, config *options.Config) *AuthService {
	return &AuthService{
		engine:     engine,
		iamClient:  iamClient,
		lgc:        lgc,
		authorizer: authorizer,
		config:     config,
	}
}

func (s *AuthService) checkRequestFromIamFilter() func(req *restful.Request, resp *restful.Response,
	chain *restful.FilterChain) {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		if !auth.EnableAuthorize() {
			chain.ProcessFilter(req, resp)
			return
		}

		// set tenant id
		tenantID, err := tools.ValidateDisableTenantMode(httpheader.GetTenantID(req.Request.Header),
			s.config.EnableMultiTenantMode)
		if err != nil {
			_ = resp.WriteAsJson(metadata.BkBaseResp{Code: types.InternalServerErrorCode, Message: err.Error()})
			return
		}
		httpheader.SetTenantID(req.Request.Header, tenantID)

		isAuthorized, err := checkRequestAuthorization(s.iamClient, req.Request)
		if err != nil {
			rsp := metadata.BkBaseResp{
				Code:    types.InternalServerErrorCode,
				Message: err.Error(),
			}
			_ = resp.WriteAsJson(rsp)
			return
		}
		if !isAuthorized {
			rsp := metadata.BkBaseResp{
				Code:    types.UnauthorizedErrorCode,
				Message: "request not from iam",
			}
			_ = resp.WriteAsJson(rsp)
			return
		}

		// use iam request id as cc rid
		rid := req.Request.Header.Get(iamtypes.IamRequestHeader)
		resp.Header().Set(iamtypes.IamRequestHeader, rid)
		if rid == "" {
			if rid = httpheader.GetRid(req.Request.Header); rid == "" {
				rid = util.GenerateRID()
			}
		}
		httpheader.SetRid(req.Request.Header, rid)
		httpheader.SetRid(resp.Header(), rid)

		// use iam language as cc language
		httpheader.SetLanguage(req.Request.Header, req.Request.Header.Get("Blueking-Language"))

		req.Request.Header = util.SetHTTPReadPreference(req.Request.Header, common.SecondaryPreferredMode)

		user := httpheader.GetUser(req.Request.Header)
		if len(user) == 0 {
			httpheader.SetUser(req.Request.Header, "auth")
		}

		chain.ProcessFilter(req, resp)
		return
	}
}

var iamToken = struct {
	token            string
	tokenRefreshTime time.Time
}{}

func checkRequestAuthorization(iamClient client.Interface, req *http.Request) (bool, error) {
	rid := req.Header.Get(iamtypes.IamRequestHeader)
	name, pwd, ok := req.BasicAuth()
	if !ok || name != iamtypes.SystemIDIAM {
		blog.Errorf("request have no basic authorization, rid: %s", rid)
		return false, nil
	}
	// if cached token is set within a minute, use it to check request authorization
	if iamToken.token != "" && time.Since(iamToken.tokenRefreshTime) <= time.Minute && pwd == iamToken.token {
		return true, nil
	}
	var err error
	iamToken.token, err = iamClient.GetSystemToken(context.Background(), req.Header)
	if err != nil {
		blog.Errorf("check request authorization get system token failed, error: %s, rid: %s", err.Error(), rid)
		return false, err
	}
	iamToken.tokenRefreshTime = time.Now()
	if pwd == iamToken.token {
		return true, nil
	}
	blog.Errorf("request password not match system token, rid: %s", rid)
	return false, nil
}

// WebService TODO
func (s *AuthService) WebService() *restful.Container {
	api := new(restful.WebService)
	api.Path("/auth/v3")
	api.Filter(s.engine.Metric().RestfulMiddleWare)
	// only allows iam to pull resource using these api
	api.Filter(s.checkRequestFromIamFilter())
	api.Produces(restful.MIME_JSON)

	s.initResourcePull(api)

	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	container.Add(api)

	authAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	authAPI.Path("/ac/v3")
	s.initAuth(authAPI)
	container.Add(authAPI)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

func (s *AuthService) initResourcePull(api *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/resource", Handler: s.PullResource})

	utility.AddToRestfulWebService(api)
}

func (s *AuthService) initAuth(api *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/authorize/batch", Handler: s.AuthorizeBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/authorize/any/batch", Handler: s.AuthorizeAnyBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/authorized_resource",
		Handler: s.ListAuthorizedResources})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/no_auth_skip_url", Handler: s.GetNoAuthSkipUrl})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/find/permission_to_apply",
		Handler: s.GetPermissionToApply})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/register/resource_creator_action",
		Handler: s.RegisterResourceCreatorAction})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/register/batch_resource_creator_action",
		Handler: s.BatchRegisterResourceCreatorAction})

	utility.AddToRestfulWebService(api)
}
