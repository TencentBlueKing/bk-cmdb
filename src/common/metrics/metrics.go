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

// Package metrics TODO
package metrics

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/types"

	"github.com/emicklei/go-restful/v3"
	"github.com/mssola/user_agent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// a global register which is used to collect metrics we need.
// it will be initialized when process is up for safe usage.
// and then be revised later when service is initialized.
var globalRegister prometheus.Registerer

func init() {
	// set default global register
	globalRegister = prometheus.DefaultRegisterer
}

// Register must only be called after backbone engine is started.
func Register() prometheus.Registerer {
	return globalRegister
}

// Namespace TODO
const Namespace = "cmdb"

// labels
const (
	LabelHandler     = "handler"
	LabelHTTPStatus  = "status_code"
	LabelOrigin      = "origin"
	LabelProcessName = "process_name"
	LabelAppCode     = "app_code"
	LabelHost        = "host"
	LabelUser        = "user"
	LabelTenantId    = "tenant_id" //租户ID
)

// labels
const (
	KeySelectedRoutePath string = "SelectedRoutePath"
)

// Config metrics config
type Config struct {
	ProcessName     string
	ProcessInstance string
}

// Service an http service
type Service struct {
	conf Config

	httpHandler http.Handler

	registry        prometheus.Registerer
	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	userTotal       *prometheus.CounterVec
}

// NewService returns new metrics service
func NewService(conf Config) *Service {
	registry := prometheus.NewRegistry()
	register := prometheus.WrapRegistererWith(prometheus.Labels{LabelProcessName: conf.ProcessName,
		LabelHost: strings.Split(conf.ProcessInstance, ":")[0]}, registry)

	// set up global register
	globalRegister = register

	srv := Service{conf: conf, registry: register}

	requestTotalLabels := []string{LabelHandler, LabelHTTPStatus, LabelOrigin, LabelAppCode, LabelTenantId}

	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: Namespace + "_http_request_total",
			Help: "http request total.",
		},
		requestTotalLabels,
	)
	register.MustRegister(srv.requestTotal)

	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    Namespace + "_http_request_duration_millisecond",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{LabelHandler, LabelAppCode, LabelTenantId},
	)
	register.MustRegister(srv.requestDuration)
	register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	register.MustRegister(prometheus.NewGoCollector())

	// add user metrics for api-server
	if conf.ProcessName == types.CC_MODULE_APISERVER {
		srv.userTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: Namespace + "_user_http_request_total",
				Help: "user http request total.",
			},
			[]string{LabelUser, LabelOrigin, LabelTenantId},
		)
		register.MustRegister(srv.userTotal)
	}

	srv.httpHandler = promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)
	return &srv
}

// Registry returns the prometheus.Registerer
func (s *Service) Registry() prometheus.Registerer {
	return s.registry
}

// ServeHTTP 用于HTTP服务
func (s *Service) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.httpHandler.ServeHTTP(resp, req)
}

// RestfulWebService is the http WebService for go-restful framework
func (s *Service) RestfulWebService() *restful.WebService {
	ws := restful.WebService{}
	ws.Path("/metrics")
	ws.Route(ws.GET("/").To(func(req *restful.Request, resp *restful.Response) {
		blog.Info("metrics")
		s.httpHandler.ServeHTTP(resp, req.Request)
	}))

	return &ws
}

// HTTPMiddleware is the http middleware for go-restful framework
func (s *Service) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.RequestURI == "/metrics" || r.RequestURI == "/metrics/" {
			s.ServeHTTP(w, r)
			return
		}

		var uri string
		req := r.WithContext(context.WithValue(r.Context(), KeySelectedRoutePath, &uri))
		resp := restful.NewResponse(w)
		before := time.Now()
		next.ServeHTTP(resp, req)
		if uri == "" {
			requestUrl, err := url.ParseRequestURI(r.RequestURI)
			if err != nil {
				return
			}
			uri = requestUrl.Path
		}

		if !utf8.ValidString(uri) {
			blog.Errorf("uri: %s not utf-8", uri)
			return
		}

		appCode := httpheader.GetAppCode(r.Header)
		tenantID := httpheader.GetTenantID(r.Header)
		s.requestDuration.With(s.label(LabelHandler, uri, LabelAppCode, appCode, LabelTenantId, tenantID)).
			Observe(float64(time.Since(before) / time.Millisecond))

		requestTotalLabels := []string{
			LabelHandler, uri,
			LabelHTTPStatus, strconv.Itoa(resp.StatusCode()),
			LabelOrigin, getOrigin(r.Header),
			LabelAppCode, appCode,
			LabelTenantId, tenantID,
		}
		s.requestTotal.With(s.label(requestTotalLabels...)).Inc()

		// add user metrics for api-server
		if s.conf.ProcessName == types.CC_MODULE_APISERVER {
			s.userTotal.With(s.label(
				LabelUser, httpheader.GetUser(r.Header),
				LabelOrigin, getOrigin(r.Header),
				LabelTenantId, tenantID,
			)).Inc()
		}
	})
}

// RestfulMiddleWare is the http middleware for go-restful framework
func (s *Service) RestfulMiddleWare(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if v := req.Request.Context().Value(KeySelectedRoutePath); v != nil {
		if selectedRoutePath, ok := v.(*string); ok {
			*selectedRoutePath = req.SelectedRoutePath()
		}
	}
	chain.ProcessFilter(req, resp)
}

// 越界问题：奇数情况在下标越界 -->奇数情况忽略最后一项KV
func (s *Service) label(labelKVs ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for index := 0; index < len(labelKVs)-1; index += 2 {
		labels[labelKVs[index]] = labelKVs[index+1]
	}
	return labels
}

func getOrigin(header http.Header) string {
	if httpheader.IsReqFromWeb(header) {
		return "webserver"
	}

	if httpheader.GetBkJWT(header) != "" {
		return "api-gateway"
	}

	if header.Get(httpheader.BKHTTPHeaderUser) != "" {
		return "esb"
	}

	if userString := header.Get("User-Agent"); userString != "" {
		ua := user_agent.New(userString)
		browser, _ := ua.Browser()
		if browser != "" {
			return "browser"
		}
	}

	return "Unknown"
}
