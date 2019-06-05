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
	"math"
	"net/http"
	"strconv"
	"sync/atomic"
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

	registry        prometheus.Registerer
	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec

	inFlight    uint64
	inFlightMax uint64
}

// NewService returns new metrics service
func NewService(conf Config) *Service {
	registry := prometheus.NewRegistry()
	register := prometheus.WrapRegistererWith(prometheus.Labels{LableProcessName: conf.ProcessName, LableProcessInstance: conf.ProcessInstance}, registry)
	srv := Service{conf: conf, registry: register}

	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: ns + "http_request_total",
			Help: "http request total.",
		},
		[]string{LableHandler, LableHTTPStatus, LableOrigin},
	)
	register.MustRegister(srv.requestTotal)

	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: ns + "http_request_duration_seconds",
			Help: "Histogram of latencies for HTTP requests.",
		},
		[]string{LableHandler},
	)
	register.MustRegister(srv.requestDuration)

	register.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: ns + "http_request_in_flight",
			Help: "current number of request being served.",
		},
		func() float64 { return math.Float64frombits((atomic.LoadUint64(&srv.inFlight))) },
	))

	register.MustRegister(prometheus.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: ns + "http_request_in_flight_max",
			Help: "max number of request being served.",
		},
		func() float64 { return math.Float64frombits((atomic.LoadUint64(&srv.inFlightMax))) },
	))

	register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	register.MustRegister(prometheus.NewGoCollector())

	srv.httpHandler = promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)
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

// lables
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
		s.increaseInFlight()
		defer s.decreaseInFlight()

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
			uri = r.RequestURI
		}
		duration := time.Since(before).Seconds()

		s.requestDuration.With(s.lable(LableHandler, uri)).Observe(duration)
		s.requestTotal.With(s.lable(
			LableHandler, uri,
			LableHTTPStatus, strconv.Itoa(resp.StatusCode()),
			LableOrigin, getOrigin(r.Header),
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

func (s *Service) increaseInFlight() {
	new := s.addInFlight(1)
	old := atomic.LoadUint64(&s.inFlightMax)
	if new > old {
		atomic.CompareAndSwapUint64(&s.inFlightMax, old, new)
	}
}

func (s *Service) decreaseInFlight() {
	s.addInFlight(-1)

}

func (s *Service) addInFlight(val float64) uint64 {
	for {
		oldBits := atomic.LoadUint64(&s.inFlight)
		newBits := math.Float64bits(math.Float64frombits(oldBits) + val)
		if atomic.CompareAndSwapUint64(&s.inFlight, oldBits, newBits) {
			return newBits
		}
	}
}

func (s *Service) lable(lableKVs ...string) prometheus.Labels {
	lables := prometheus.Labels{}
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
