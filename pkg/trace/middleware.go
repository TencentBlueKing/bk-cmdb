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

package trace

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

const (
	// requestIDHeader is the name of the HTTP Header which contains the request id.
	requestIDHeader = "X-Request-Id"
)

// requestIDContext 自定义从X-Request-Id统一转换为traceID
type requestIDContext struct{}

// Inject ...
func (r requestIDContext) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.HasTraceID() {
		return
	}

	carrier.Set(requestIDHeader, spanCtx.TraceID().String())
}

// Extract 如果没有traceparent,且X-Request-Id是合法的,生成对应的traceID
func (r requestIDContext) Extract(ctx context.Context, carrier propagation.TextMapCarrier) context.Context {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		return ctx
	}

	rid := carrier.Get(requestIDHeader)
	if rid == "" {
		return ctx
	}

	// 兼容uuid格式
	newRid := strings.ReplaceAll(rid, "-", "")

	traceID, err := trace.TraceIDFromHex(newRid)
	if err != nil {
		log.Error(ctx, "conv rid to traceID failed", log.RidAttr(rid), log.E(err))
		return ctx
	}

	span := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
	})

	return trace.ContextWithRemoteSpanContext(ctx, span)
}

// Fields ...
func (r requestIDContext) Fields() []string {
	return []string{requestIDHeader}
}

// Middleware 扩展 otelhttp trace 中间件, 支持traceparent / X-Request-Id
func Middleware(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		span := trace.SpanFromContext(r.Context())
		ctx := log.WithSpan(r.Context(), span)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}

	otelHandler := otelhttp.NewHandler(
		http.HandlerFunc(f),
		string(config.GetServiceName()),
	)

	return otelHandler
}
