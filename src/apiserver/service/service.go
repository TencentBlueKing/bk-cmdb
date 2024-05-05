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

// Package service TODO
package service

import (
	"configcenter/src/ac"
	"configcenter/src/ac/iam"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metrics"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful/v3"
	"github.com/prometheus/client_golang/prometheus"
)

// Service service methods
type Service interface {
	WebServices() []*restful.WebService
	SetConfig(engine *backbone.Engine, httpClient HTTPClient, discovery discovery.DiscoveryInterface,
		clientSet apimachinery.ClientSetInterface, cache redis.Client, limiter *Limiter)
}

// NewService create a new service instance
func NewService() Service {
	return new(service)
}

type service struct {
	engine     *backbone.Engine
	client     HTTPClient
	discovery  discovery.DiscoveryInterface
	clientSet  apimachinery.ClientSetInterface
	authorizer ac.AuthorizeInterface
	cache      redis.Client
	limiter    *Limiter
	// noPermissionRequestTotal is the total number of request without permission
	noPermissionRequestTotal *prometheus.CounterVec
}

// SetConfig set config
func (s *service) SetConfig(engine *backbone.Engine, httpClient HTTPClient, discovery discovery.DiscoveryInterface,
	clientSet apimachinery.ClientSetInterface, cache redis.Client, limiter *Limiter) {
	s.engine = engine
	s.client = httpClient
	s.discovery = discovery
	s.clientSet = clientSet
	s.cache = cache
	s.limiter = limiter
	s.authorizer = iam.NewAuthorizer(clientSet)
}

// WebServices TODO
func (s *service) WebServices() []*restful.WebService {
	getErrFun := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	ws := &restful.WebService{}
	ws.Path(rootPath)
	ws.Filter(s.JwtFilter())
	ws.Filter(s.engine.Metric().RestfulMiddleWare)
	ws.Filter(rdapi.AllGlobalFilter(getErrFun))
	ws.Filter(rdapi.RequestLogFilter())
	ws.Filter(s.LimiterFilter())
	ws.Produces(restful.MIME_JSON)

	// route skip auth api
	s.routeSkipAuthAPI(ws)

	// route need auth api
	s.routeNeedAuthAPI(ws, getErrFun)

	allWebServices := make([]*restful.WebService, 0)
	allWebServices = append(allWebServices, ws)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	allWebServices = append(allWebServices, commonAPI)

	return allWebServices
}

// routeSkipAuthAPI route apis that need skip api server authorization, and authorize in its scene server logics
// note: this is only temporary, delete the api server authorize logic when all api is updated
func (s *service) routeSkipAuthAPI(ws *restful.WebService) {
	ws.Route(ws.POST("/auth/verify").To(s.AuthVerify))
	ws.Route(ws.GET("/auth/business_list").To(s.GetAnyAuthorizedAppList))
	ws.Route(ws.POST("/auth/skip_url").To(s.GetUserNoAuthSkipURL))

	ws.Route(ws.POST("/biz/{.*}").Filter(s.BizFilterChan).To(s.Post))
	ws.Route(ws.POST("/biz/search/{.*}").Filter(s.BizFilterChan).To(s.Post))

	ws.Route(ws.POST("/findmany/hosts/by_service_templates/biz/{.*}").Filter(s.HostFilterChan).To(s.Post))
	ws.Route(ws.POST("/findmany/module_relation/bk_biz_id/{.*}").Filter(s.HostFilterChan).To(s.Post))
	ws.Route(ws.POST("/findmany/hosts/relation/with_topo").Filter(s.HostFilterChan).To(s.Post))
	ws.Route(ws.PUT("/updatemany/hosts/all/property").Filter(s.HostFilterChan).To(s.Put))
	ws.Route(ws.POST("/check/objectattr/host_apply_enabled").Filter(s.HostFilterChan).To(s.Post))

	ws.Route(ws.POST("/update/transaction/commit").Filter(s.TxnFilterChan).To(s.Post))
	ws.Route(ws.POST("/update/transaction/abort").Filter(s.TxnFilterChan).To(s.Post))

	ws.Route(ws.POST("/count/{bk_obj_id}/instances").To(s.CountInstance))
	ws.Route(ws.POST("/group/related/{kind}/resource/by_ids").Filter(s.WebCoreFilterChan).To(s.Post))
}

func (s *service) routeNeedAuthAPI(ws *restful.WebService, errFunc func() errors.CCErrorIf) {
	if auth.EnableAuthorize() {
		s.noPermissionRequestTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cmdb_no_permission_request_total",
				Help: "total number of request without permission.",
			},
			[]string{metrics.LabelHandler, metrics.LabelAppCode},
		)
		s.engine.Metric().Registry().MustRegister(s.noPermissionRequestTotal)

	}

	ws.Route(ws.GET("{.*}").Filter(s.authFilter(errFunc)).Filter(s.URLFilterChan).To(s.Get))
	ws.Route(ws.POST("{.*}").Filter(s.authFilter(errFunc)).Filter(s.URLFilterChan).To(s.Post))
	ws.Route(ws.PUT("{.*}").Filter(s.authFilter(errFunc)).Filter(s.URLFilterChan).To(s.Put))
	ws.Route(ws.DELETE("{.*}").Filter(s.authFilter(errFunc)).Filter(s.URLFilterChan).To(s.Delete))
}
