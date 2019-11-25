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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"

	"github.com/emicklei/go-restful"
	"github.com/mssola/user_agent"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const ns = "cmdb_"

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
	requestInflight *Gauge
}

// NewService returns new metrics service
func NewService(conf Config) *Service {
	registry := prometheus.NewRegistry()
	register := prometheus.WrapRegistererWith(prometheus.Labels{LabelProcessName: conf.ProcessName, LabelHost: strings.Split(conf.ProcessInstance, ":")[0]}, registry)
	srv := Service{conf: conf, registry: register}

	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: ns + "http_request_total",
			Help: "http request total.",
		},
		[]string{LabelHandler, LabelHTTPStatus, LabelOrigin},
	)
	register.MustRegister(srv.requestTotal)

	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: ns + "http_request_duration_millisecond",
			Help: "Histogram of latencies for HTTP requests.",
		},
		[]string{LabelHandler},
	)
	register.MustRegister(srv.requestDuration)

	srv.requestInflight = NewGauge(
		prometheus.GaugeOpts{
			Name: ns + "http_request_in_flight",
			Help: "current number of request being served.",
		},
	)
	register.MustRegister(srv.requestInflight)

	register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	register.MustRegister(prometheus.NewGoCollector())

	srv.httpHandler = promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)
	return &srv
}

// labels
const (
	LabelHandler     = "handler"
	LabelHTTPStatus  = "status_code"
	LabelOrigin      = "origin"
	LabelProcessName = "process_name"
	LabelHost        = "host"
)

// labels
const (
	KeySelectedRoutePath string = "SelectedRoutePath"
)

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
		s.requestInflight.Inc()
		defer s.requestInflight.Dec()

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
		duration := util.ToMillisecond(time.Since(before))

		s.requestDuration.With(s.label(LabelHandler, uri)).Observe(duration)
		s.requestTotal.With(s.label(
			LabelHandler, uri,
			LabelHTTPStatus, strconv.Itoa(resp.StatusCode()),
			LabelOrigin, getOrigin(r.Header),
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
