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

// Package middleware defines the middleware for http and grpc server
package middleware

import (
	"net/http"

	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// HttpMiddleware is the http middleware with kit.
type HttpMiddleware func(kt *kit.Kit, w *ResponseWriter, r *http.Request, next http.HandlerFunc)

// ResponseWriter is the cmdb http response writer that records the status code.
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newResponseWriter creates a new ResponseWriter.
func newResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, http.StatusOK}
}

// WriteHeader records the status code when WriteHeader is called.
func (w *ResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// GetStatusCode returns the status code that was recorded when WriteHeader was called.
func (w *ResponseWriter) GetStatusCode() int {
	return w.statusCode
}

// ConvHttpMiddleware returns a http middleware with the given cmdb http interceptors.
func ConvHttpMiddleware(middlewares ...HttpMiddleware) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			kt := kit.NewKitFromHeader(r.Context(), r.Header)

			if err := kt.Validate(); err != nil {
				log.Error(r.Context(), "http kit is invalid", log.E(err))
				_ = rest.APIError(r.Context(), cerr.Wrap(cerr.INVALID_REQUEST, err)).Render(w)
				return
			}

			resp := newResponseWriter(w)

			// use all middlewares
			serveHttp := next.ServeHTTP
			for i := len(middlewares) - 1; i >= 0; i-- {
				nextServeHttp := serveHttp
				middleware := middlewares[i]
				serveHttp = func(writer http.ResponseWriter, request *http.Request) {
					middleware(kt, resp, r, nextServeHttp)
				}
			}
			serveHttp(resp, r)
		})
	}
}
