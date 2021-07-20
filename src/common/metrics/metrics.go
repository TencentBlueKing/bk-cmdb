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

package metrics

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"github.com/emicklei/go-restful"
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

const Namespace = "cmdb"

// labels
const (
	LabelHandler     = "handler"
	LabelHTTPStatus  = "status_code"
	LabelOrigin      = "origin"
	LabelProcessName = "process_name"
	LabelAppCode     = "app_code"
	LabelHost        = "host"
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
}

// NewService returns new metrics service
func NewService(conf Config) *Service {
	registry := prometheus.NewRegistry()
	register := prometheus.WrapRegistererWith(prometheus.Labels{LabelProcessName: conf.ProcessName, LabelHost: strings.Split(conf.ProcessInstance, ":")[0]}, registry)

	// set up global register
	globalRegister = register

	srv := Service{conf: conf, registry: register}

	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: Namespace + "_http_request_total",
			Help: "http request total.",
		},
		[]string{LabelHandler, LabelHTTPStatus, LabelOrigin, LabelAppCode},
	)
	register.MustRegister(srv.requestTotal)

	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    Namespace + "_http_request_duration_millisecond",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{LabelHandler, LabelAppCode},
	)
	register.MustRegister(srv.requestDuration)
	register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	register.MustRegister(prometheus.NewGoCollector())

	srv.httpHandler = promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)
	return &srv
}

// Registry returns the prometheus.Registerer
func (s *Service) Registry() prometheus.Registerer {
	return s.registry
}

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

		s.requestDuration.With(s.label(LabelHandler, uri, LabelAppCode, r.Header.Get(common.BKHTTPRequestAppCode))).
			Observe(float64(time.Since(before) / time.Millisecond))

		s.requestTotal.With(s.label(
			LabelHandler, uri,
			LabelHTTPStatus, strconv.Itoa(resp.StatusCode()),
			LabelOrigin, getOrigin(r.Header),
			LabelAppCode, r.Header.Get(common.BKHTTPRequestAppCode),
		)).Inc()
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

func (s *Service) label(labelKVs ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for index := 0; index < len(labelKVs); index += 2 {
		labels[labelKVs[index]] = labelKVs[index+1]
	}
	return labels
}

func getOrigin(header http.Header) string {
	if header.Get(common.BKHTTPOtherRequestID) != "" {
		return "ESB"
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
