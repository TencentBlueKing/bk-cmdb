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

// Package kit is designed to manage and store business metadata, context and trace info
package kit

import (
	"context"
	"runtime"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

const (
	defaultScopeName = "kit"
)

// Metadata the biz metadata
type Metadata struct {
	User     string // 操作人
	AppCode  string // 来源AppCode
	TenantID string // 来源多租户ID
}

// Kit a biz metadata and context kit
type Kit struct {
	context.Context
	Metadata
}

// WithCancel same as context.WithCancel but return Kit
func (kt *Kit) WithCancel() (*Kit, context.CancelFunc) {
	ctx, cancel := context.WithCancel(kt.Context)
	return newKit(ctx, kt.Metadata), cancel
}

// WithTimeout same as context.WithTimeout but return Kit
func (kt *Kit) WithTimeout(timeout time.Duration) (*Kit, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(kt.Context, timeout)
	return newKit(ctx, kt.Metadata), cancel
}

// StartSpan return a new kit and within a span
func (kt *Kit) StartSpan(name string, opts ...trace.SpanStartOption) (*Kit, trace.Span) {
	scopeName := defaultScopeName

	// auto use caller lib and pkg.func name
	if name == "" {
		scopeName, name = getCaller()
	}

	ctx, span := otel.Tracer(scopeName).Start(kt.Context, name, opts...)

	// auto set metadata attr
	span.SetAttributes(
		attribute.String("user", kt.User),
		attribute.String("appCode", kt.AppCode),
		attribute.String("tenantID", kt.TenantID),
	)

	ctx = log.WithSpan(ctx, span)
	return newKit(ctx, kt.Metadata), span
}

// SetSpanAttr set attr to current span
func (kt *Kit) SetSpanAttr(attr map[string]string) {
	span := trace.SpanFromContext(kt.Context)

	for k, v := range attr {
		span.SetAttributes(attribute.String(k, v))
	}
}

// Rid aka traceID/request_id
func (kt *Kit) Rid() string {
	span := trace.SpanContextFromContext(kt.Context)
	if span.HasTraceID() {
		return span.TraceID().String()
	}

	return ""
}

// NewKit with metadata
func NewKit(ctx context.Context, md Metadata) *Kit {
	return newKit(ctx, md)
}

func newKit(ctx context.Context, md Metadata) *Kit {
	return &Kit{Context: ctx, Metadata: md}
}

// getCaller 获取调用函数名
func getCaller() (string, string) {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "", ""
	}

	fn := runtime.FuncForPC(pc).Name()
	idx := strings.LastIndexByte(fn, '/')
	if idx < 0 || len(fn) <= idx+1 {
		return fn, ""
	}

	// 返回 libname pkg.func
	return fn[:idx], fn[idx+1:]
}
