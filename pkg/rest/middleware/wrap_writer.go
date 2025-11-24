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

package middleware

import (
	"net/http"
)

// WrapResponseWriter is a proxy around an http.ResponseWriter that allows you to hook
// into various parts of the response process.
type WrapResponseWriter interface {
	http.ResponseWriter
	Status() int
	BytesWritten() int
}

// NewWrapResponseWriter wraps an http.ResponseWriter, returning a proxy that allows you to
// hook into various parts of the response process.
func NewWrapResponseWriter(w http.ResponseWriter) WrapResponseWriter {
	bw := basicWriter{
		ResponseWriter:     w,
		ResponseController: http.NewResponseController(w),
	}
	return &bw
}

// basicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type basicWriter struct {
	http.ResponseWriter
	*http.ResponseController
	code  int
	bytes int
}

// Header implement the http.ResponseWriter interface Header method
func (bw *basicWriter) Header() http.Header {
	return bw.ResponseWriter.Header()
}

// WriteHeader implement the http.ResponseWriter interface WriteHeader method
func (bw *basicWriter) WriteHeader(code int) {
	if bw.code != 0 {
		return
	}

	bw.code = code
	bw.ResponseWriter.WriteHeader(code)
}

// Write implement the http.ResponseWriter interface Write method
func (bw *basicWriter) Write(buf []byte) (int, error) {
	if bw.code == 0 {
		bw.WriteHeader(http.StatusOK)
	}

	n, err := bw.ResponseWriter.Write(buf)
	bw.bytes += n
	return n, err
}

// Status return the status code of the response
func (bw *basicWriter) Status() int {
	return bw.code
}

// BytesWritten return the bytes written
func (bw *basicWriter) BytesWritten() int {
	return bw.bytes
}
