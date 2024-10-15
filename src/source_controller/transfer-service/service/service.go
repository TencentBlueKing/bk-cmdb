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

// Package service provides transfer service's web service
package service

import (
	"net/http"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/sync"
	"configcenter/src/storage/stream"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service defines transfer service's web service
type Service struct {
	engine *backbone.Engine
	syncer *sync.Syncer
}

// New Service
func New(conf *options.Config, engine *backbone.Engine) (*Service, error) {
	loopW, err := stream.NewLoopStream(conf.Mongo.GetMongoConf(), engine.ServiceManageInterface)
	if err != nil {
		blog.Errorf("new loop stream failed, err: %v", err)
		return nil, err
	}

	syncer, err := sync.NewSyncer(conf, engine.ServiceManageInterface, loopW, engine.CoreAPI.CacheService(),
		engine.Metric().Registry())
	if err != nil {
		blog.Errorf("new syncer failed, err: %v", err)
		return nil, err
	}

	return &Service{
		engine: engine,
		syncer: syncer,
	}, nil
}

// WebService provides web service
func (s *Service) WebService() *restful.Container {
	errors.SetGlobalCCError(s.engine.CCErr)
	getErrFunc := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	api := new(restful.WebService)
	api.Path("/transfer/v3/").Filter(s.engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).
		Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)

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
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.engine.CCErr,
		Language: s.engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/sync/cmdb/data", Handler: s.SyncCmdbData})

	utility.AddToRestfulWebService(api)
}
