/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

// Package service provides sync server's web service
package service

import (
	"net/http"

	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/sync_server/logics"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service defines sync server's web service
type Service struct {
	engine *backbone.Engine
	lgc    *logics.Logics
}

// New Service
func New(engine *backbone.Engine, lgc *logics.Logics) *Service {
	return &Service{
		engine: engine,
		lgc:    lgc,
	}
}

// WebService provides web service
func (s *Service) WebService() *restful.Container {
	errors.SetGlobalCCError(s.engine.CCErr)
	getErrFunc := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	api := new(restful.WebService)
	api.Path("/sync/v3/").Filter(s.engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).
		Produces(restful.MIME_JSON)

	// init service actions
	s.initService(api)

	container := restful.NewContainer().Add(api)

	opentelemetry.AddOtlpFilter(container)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

func (s *Service) initService(api *restful.WebService) {
	u := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	u.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/sync/full/text/search/data",
		Handler: s.SyncFullTextSearchData})
	u.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/migrate/full/text/search",
		Handler: s.MigrateFullTextSearch})

	u.AddToRestfulWebService(api)
}
