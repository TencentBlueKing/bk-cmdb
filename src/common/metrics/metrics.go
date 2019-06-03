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
	"net/http"
	"strconv"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"

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

	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestInFlight *prometheus.GaugeVec
}

// NewService returns new metrics service
func NewService(conf Config) *Service {
	srv := Service{}
	registry := prometheus.NewRegistry()

	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: ns + "http_request_total",
			Help: "http request total.",
		},
		[]string{LableProcessName, LableProcessInstance, LableHandler, LableHTTPStatus},
	)
	registry.MustRegister(srv.requestTotal)

	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: ns + "http_request_duration_seconds",
			Help: "Histogram of latencies for HTTP requests.",
		},
		[]string{LableProcessName, LableProcessInstance, LableHandler},
	)
	registry.MustRegister(srv.requestDuration)

	srv.requestInFlight = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: ns + "http_request_in_flight",
			Help: "current number of request being served.",
		},
		[]string{LableProcessName, LableProcessInstance, LableHandler},
	)
	registry.MustRegister(srv.requestInFlight)

	srv.httpHandler = promhttp.HandlerFor(registry, promhttp.HandlerOpts{ErrorLog: blog.GlogWriter{}})
	return &srv
}

// lables
const (
	LableHandler         = "handler"
	LableHTTPStatus      = "status_code"
	LableOrigin          = "origin"
	LableProcessName     = "process_name"
	LableProcessInstance = "process_instance"
)

// MiddleWareFunc is the http middleware for go-restful framework
func (s *Service) MiddleWareFunc(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	uri := req.SelectedRoutePath()
	s.requestInFlight.With(s.lable(LableHandler, uri)).Inc()

	before := time.Now()
	chain.ProcessFilter(req, resp)
	duration := time.Since(before).Seconds()

	s.requestTotal.With(s.lable(
		LableHandler, uri,
		LableHTTPStatus, strconv.Itoa(resp.StatusCode()),
		LableOrigin, getOrigin(req.Request.Header),
	)).Inc()
	s.requestDuration.With(s.lable(LableHandler, uri)).Observe(duration)
	s.requestInFlight.With(s.lable(LableHandler, uri)).Inc()
}

func (s *Service) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.httpHandler.ServeHTTP(resp, req)
}

func (s *Service) lable(lableKVs ...string) prometheus.Labels {
	lables := prometheus.Labels{LableProcessName: s.conf.ProcessName, LableProcessInstance: s.conf.ProcessInstance}
	for index := 0; index < len(lableKVs); index += 2 {
		lables[lableKVs[index]] = lableKVs[index+1]
	}
	return lables
}

func getOrigin(header http.Header) string {
	if header.Get(common.BKHTTPOtherRequestID) != "" {
		return "ESB"
	}
	if uastring := header.Get("User-Agent"); uastring != "" {
		ua := user_agent.New(uastring)
		browser, _ := ua.Browser()
		if browser != "" {
			return "browser"
		}
	}

	return "Unknow"
}
