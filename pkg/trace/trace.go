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

// Package trace setup otlp trace and unified the traceID, requestID header, rid log
package trace

import (
	"context"
	"fmt"
	"io"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/version"
)

// Option trace setup opt
type Option struct {
	Enabled  bool   `json:"enabled"`
	Endpoint string `json:"endpoint"`
	Token    string `json:"token"`
}

// SetupTrace setup trace by opt
func SetupTrace(ctx context.Context, opt *Option) error {
	var (
		exporter sdktrace.SpanExporter
		err      error
	)

	if opt.Enabled && opt.Endpoint != "" {
		log.Info(ctx, "trace enabled", "exporter", "otlptracehttp", "endpoint", opt.Endpoint)

		headers := map[string]string{"x-bk-token": opt.Token}

		// 配置http上报地址
		option := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(opt.Endpoint),
			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		}
		exporter, err = otlptracehttp.New(ctx, option...)
	} else {
		// 日志打印rid/span等
		exporter, err = stdouttrace.New(stdouttrace.WithWriter(io.Discard))
	}

	if err != nil {
		return err
	}

	tp, err := newTraceProvider(exporter)
	if err != nil {
		return err
	}

	// Set as global trace provider
	otel.SetTracerProvider(tp)

	// 支持传递traceparent/X-Request-Id
	// W3C Trace Context traceparent = <version>-<trace-id>-<parent-id>-<flags>
	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, // 支持从外部设置 traceparent
		requestIDContext{},         // 支持外部设置 X-Request-Id
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)
	return nil
}

// Shutdown graceful shutdown trace exporter
func Shutdown(ctx context.Context) error {
	tp, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider)
	if !ok || tp == nil {
		return nil
	}

	return tp.Shutdown(ctx)
}

// OTEL tracer provider setup
func newTraceProvider(exp sdktrace.SpanExporter) (*sdktrace.TracerProvider, error) {
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(string(config.GetServiceName())),
			semconv.ServiceVersion(version.Version),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("merge resource: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)
	return provider, nil
}

func init() {
	// 默认只做context传递, 不导出到其他输出
	opt := &Option{}
	lo.Must0(SetupTrace(context.Background(), opt))
}
