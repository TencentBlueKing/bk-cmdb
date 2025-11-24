/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

// Package metrics defines metrics collecting related logics.
package metrics

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest/middleware"
	rtmiddleware "github.com/TencentBlueKing/bk-cmdb/pkg/runtime/server/middleware"
)

// registerer is a global register which is used to collect metrics we need.
// it will be initialized when process is up for safe usage.
// and then be revised later when service is initialized.
var registerer prometheus.Registerer

func init() {
	// set default global register
	registerer = prometheus.DefaultRegisterer
}

// Registerer returns the global prometheus registerer, must only be called after metrics service is initialized.
func Registerer() prometheus.Registerer {
	return registerer
}

const (
	// Namespace is the namespace of cmdb metrics.
	Namespace = "cmdb"
	// OrmCmdSubsystem defines the subsystem of the orm command
	OrmCmdSubsystem = "orm"
)

// labels
const (
	LabelHandler     = "handler"
	LabelHTTPStatus  = "status_code"
	LabelProcessName = "process_name"
	LabelAppCode     = "app_code"
	LabelHost        = "host"
	LabelTenantId    = "tenant_id"
	LabelReqProtocol = "protocol"
)

// Config metrics config
type Config struct {
	ProcessName config.ServiceName
	Host        string
}

// Service defines the common metrics info for cmdb service.
type Service struct {
	service config.ServiceName

	httpHandler http.Handler

	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

// NewService returns new metrics service.
func NewService(conf *Config) *Service {
	registry := prometheus.NewRegistry()

	// set up global register
	registerer = prometheus.WrapRegistererWith(prometheus.Labels{LabelProcessName: string(conf.ProcessName),
		LabelHost: conf.Host}, registry)

	srv := Service{service: conf.ProcessName}

	// register request total metrics
	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: Namespace + "_request_total",
			Help: "total requests.",
		},
		[]string{LabelHandler, LabelHTTPStatus, LabelAppCode, LabelTenantId, LabelReqProtocol},
	)
	registerer.MustRegister(srv.requestTotal)

	// register request duration metrics
	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    Namespace + "_request_duration_ms",
			Help:    "histogram of latencies for requests.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{LabelHandler, LabelAppCode, LabelTenantId, LabelReqProtocol},
	)
	registerer.MustRegister(srv.requestDuration)

	// register process and golang collectors
	registerer.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	registerer.MustRegister(collectors.NewGoCollector())

	// create http handler for metrics
	srv.httpHandler = promhttp.InstrumentMetricHandler(registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	return &srv
}

// ServeHTTP returns the ServeHTTP method of metrics http handler.
func (s *Service) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.httpHandler.ServeHTTP(resp, req)
}

// HTTPMiddleware returns the HTTP middleware that records user metrics.
func (s *Service) HTTPMiddleware(kt *kit.Kit, w middleware.WrapResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.RequestURI == "/metrics" || r.RequestURI == "/metrics/" {
		s.ServeHTTP(w, r)
		return
	}

	before := time.Now()
	next.ServeHTTP(w, r)

	// get request uri, use route pattern if it's not empty
	uri := rest.RoutePattern(r)
	if uri == "" {
		requestUrl, err := url.ParseRequestURI(r.RequestURI)
		if err != nil {
			log.Error(r.Context(), "uri: %s is invalid", r.RequestURI)
			return
		}
		uri = requestUrl.Path
	}

	if !utf8.ValidString(uri) {
		log.Error(r.Context(), "uri: %s is invalid", r.RequestURI)
		return
	}

	s.recordMetrics(kt, before, "http", uri, strconv.Itoa(w.Status()))
}

func (s *Service) recordMetrics(kt *kit.Kit, before time.Time, protocol, handler, status string) {
	// collect request duration metrics
	s.requestDuration.With(prometheus.Labels{
		LabelHandler:     handler,
		LabelAppCode:     kt.AppCode,
		LabelTenantId:    kt.TenantID,
		LabelReqProtocol: protocol,
	}).Observe(float64(time.Since(before) / time.Millisecond))

	// collect request total metrics
	s.recordRequestTotal(kt, protocol, handler, status)
}

func (s *Service) recordRequestTotal(kt *kit.Kit, protocol, handler, status string) {
	s.requestTotal.With(prometheus.Labels{
		LabelHandler:     handler,
		LabelHTTPStatus:  status,
		LabelAppCode:     kt.AppCode,
		LabelTenantId:    kt.TenantID,
		LabelReqProtocol: protocol,
	}).Inc()
}

// GrpcUnaryServerInterceptor is the gRPC unary server interceptor that records user metrics.
func (s *Service) GrpcUnaryServerInterceptor(kt *kit.Kit, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (any, error) {

	before := time.Now()

	// get grpc method name
	method := strings.TrimPrefix(info.FullMethod, "/")
	if i := strings.Index(method, "/"); i >= 0 {
		method = method[i+1:]
	}

	// call grpc handler
	resp, err := handler(kt, req)

	// record grpc metrics
	s.recordMetrics(kt, before, "grpc", method, s.getHttpStatusByGrpcErr(err))

	return resp, err
}

func (s *Service) getHttpStatusByGrpcErr(err error) string {
	if err == nil {
		return "200"
	}

	statusCode := http.StatusInternalServerError
	st, ok := status.FromError(err)
	if ok {
		statusCode = runtime.HTTPStatusFromCode(st.Code())
	} else {
		var e cerr.CodeError
		if errors.As(err, &e) {
			statusCode = cerr.GetHTTPStatus(e.GetCode())
		}
	}
	return strconv.Itoa(statusCode)
}

// GrpcStreamServerInterceptor is the gRPC stream server interceptor that records user metrics.
func (s *Service) GrpcStreamServerInterceptor(kt *kit.Kit, srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	before := time.Now()
	// get grpc method name
	method := strings.TrimPrefix(info.FullMethod, "/")
	if i := strings.Index(method, "/"); i >= 0 {
		method = method[i+1:]
	}

	// record grpc total request metrics for message send and receive
	ssWrapper := rtmiddleware.NewServerStreamWrapper(kt, ss, method, func(m any) error {
		err := ss.SendMsg(m)
		s.recordRequestTotal(kt, "grpc", method+"_send", s.getHttpStatusByGrpcErr(err))
		return err
	}, func(m any) error {
		err := ss.RecvMsg(m)
		if errors.Is(err, io.EOF) {
			return err
		}
		s.recordRequestTotal(kt, "grpc", method+"_recv", s.getHttpStatusByGrpcErr(err))
		return err
	})

	// call grpc handler and record metrics
	err := handler(srv, ssWrapper)
	s.recordMetrics(kt, before, "grpc", method, s.getHttpStatusByGrpcErr(err))
	return err
}
