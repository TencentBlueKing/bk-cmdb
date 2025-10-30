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

	"github.com/google/uuid"

	"github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
)

// UnaryFunc Unary or ClientStreaming handle function
type UnaryFunc[Req, Resp any] func(context.Context, *Req) (*Resp, error)

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

		rid := uuid.New().String()
		ctx := r.Context()
		ctx = log.WithAttr(ctx, log.RidAttr(rid))

		// 反序列化
		in, err := decodeReq[Req](r)
		if err != nil {
			log.Error(ctx, "handle decode request failed", log.E(err))
			ApiRespError(err, w, r, cerr.INVALID_REQUEST)
			return
		}

		// 参数校验
		if err = validateReq(r.Context(), in); err != nil {
			log.Error(ctx, "validate req failed", log.E(err))
			ApiRespError(err, w, r, cerr.INVALID_REQUEST)
			return
		}

		out, respErr := fn(ctx, in)
		if respErr != nil {
			ApiRespError(respErr, w, r, "")
			return
		}
		_ = APIOK(out).Render(w, r)
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

		// 反序列化
		in, err := decodeReq[Req](r)
		if err != nil {
			log.Error(ctx, "handle decode stream request failed", log.E(err))
			ApiRespError(err, w, r, cerr.INVALID_REQUEST)
			return
		}

		// 参数校验
		if err = validateReq(r.Context(), in); err != nil {
			log.Error(ctx, "validate stream req failed", log.E(err))
			ApiRespError(err, w, r, cerr.INVALID_REQUEST)
			return
		}

		svr := &streamingServer{
			ResponseWriter:     w,
			ResponseController: http.NewResponseController(w),
			ctx:                r.Context(),
		}

		if err := fn(in, svr); err != nil {
			ApiRespError(err, w, r, "")
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

// ApiRespError return api response error
func ApiRespError(err error, w http.ResponseWriter, r *http.Request, errorCode cerr.ErrorCode) {
	var respErr *cerr.RespError
	if errorCode != "" {
		respErr = cerr.GetDefaultErrorManager().ConvToRespError(err, cerr.WithCode(errorCode))
	} else {
		// convert error to response error
		respErr = cerr.GetDefaultErrorManager().ConvToRespError(err)
	}
	// translate error message
	respErr = i18n.GetDefaultManager().RespError(r.Context(), respErr)
	_ = APIError(respErr).Render(w, r)
}
