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

// Package rest framework
package rest

import (
	"context"
	"net/http"
	"time"

	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// UnaryFunc Unary or ClientStreaming handle function
type UnaryFunc[Req, Resp any] func(*kit.Kit, *Req) (*Resp, error)

// StreamingServer server or bidi streaming server
type StreamingServer interface {
	http.ResponseWriter
	Context() context.Context
}

// StreamFunc ServerStreaming or BidiStreaming handle function
type StreamFunc[Req any] func(*Req, StreamingServer) error

// Handle Composable HTTP Handlers using generics
func Handle[Req, Resp any](fn UnaryFunc[Req, Resp]) func(w http.ResponseWriter, r *http.Request) {
	handleName := getHandleName(fn)

	f := func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		var err error
		defer func() {
			collectHandleMetrics(handleName, r.Method, st, err)
		}()

		kt := kit.NewKitFromHeader(r.Context(), r.Header)

		// 反序列化
		in, err := decodeReq[Req](r)
		if err != nil {
			log.Error(kt, "handle decode request failed", log.E(err))
			_ = APIError(kt, cerr.Wrap(cerr.InvalidRequest, err)).Render(w)
			return
		}

		// 参数校验
		if err = validateReq(kt, in); err != nil {
			log.Error(kt, "validate req failed", log.E(err))
			_ = APIError(kt, cerr.Wrap(cerr.InvalidRequest, err)).Render(w)
			return
		}

		out, respErr := fn(kt, in)
		if respErr != nil {
			_ = APIError(kt, respErr).Render(w)
			return
		}

		_ = APIOK(out).Render(w)
	}
	return f
}

type streamingServer struct {
	http.ResponseWriter
	*http.ResponseController
	ctx context.Context
}

// Context return svr's context
func (s *streamingServer) Context() context.Context {
	return s.ctx
}

// Stream Composable HTTP Handlers using generics
func Stream[Req any](fn StreamFunc[Req]) func(w http.ResponseWriter, r *http.Request) {
	handleName := getHandleName(fn)

	f := func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		var err error
		defer func() {
			collectHandleMetrics(handleName, r.Method, st, err)
		}()

		ctx := r.Context()
		kt := kit.NewKitFromHeader(ctx, r.Header)

		// 反序列化
		in, err := decodeReq[Req](r)
		if err != nil {
			log.Error(kt, "handle decode stream request failed", log.E(err))
			_ = APIError(kt, cerr.Wrap(cerr.InvalidRequest, err)).Render(w)
			return
		}

		// 参数校验
		if err = validateReq(kt, in); err != nil {
			log.Error(kt, "validate stream req failed", log.E(err))
			_ = APIError(kt, cerr.Wrap(cerr.InvalidRequest, err)).Render(w)
			return
		}

		svr := &streamingServer{
			ResponseWriter:     w,
			ResponseController: http.NewResponseController(w),
			ctx:                kt,
		}

		if err := fn(in, svr); err != nil {
			_ = APIError(kt, err).Render(w)
			return
		}
	}

	return f
}

// EmptyReq 空的请求
type EmptyReq struct{}

// EmptyResp 空的返回
type EmptyResp struct{}

// PaginationReq 分页接口通用请求
type PaginationReq struct {
	Offset int `json:"offset" in:"query=offset" validate:"gte=0"`
	Limit  int `json:"limit" in:"query=limit" validate:"gte=0"`
}

// PaginationResp 分页接口通用返回
type PaginationResp[T any] struct {
	Count int64 `json:"count"`
	Items []T   `json:"items"`
}
