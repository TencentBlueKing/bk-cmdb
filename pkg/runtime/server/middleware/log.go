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
	"bytes"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/grpc"

	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest/middleware"
	"github.com/TencentBlueKing/bk-cmdb/pkg/validator"
)

// PrintHttpLog print http request log.
func PrintHttpLog(kt *kit.Kit, w middleware.WrapResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	body := make([]byte, 0)
	if r.Body != nil {
		byt, err := io.ReadAll(r.Body)
		if err == nil {
			r.Body = io.NopCloser(bytes.NewBuffer(byt))
		}
		body = byt
	}

	next(w, r)

	log.Info(kt, "http request", "appCode", kt.AppCode, "user", kt.User, "method", r.Method, "uri",
		r.RequestURI, "body", truncateBody(body), "cost", time.Since(start), "status", w.Status())
}

func truncateBody(body []byte) string {
	if len(body) > 2048 {
		return fmt.Sprintf("%s...(Total %dB)", body[:2048], len(body))
	}
	return string(body)
}

// PrintGrpcLog print grpc call log.
func PrintGrpcLog(kt *kit.Kit, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()
	body, _ := json.Marshal(req)

	resp, err := handler(kt, req)

	log.Info(kt, "grpc call", "appCode", kt.AppCode, "user", kt.User, "method", info.FullMethod, "body",
		truncateBody(body), "cost", time.Since(start))

	return resp, err
}

// PrintGrpcStreamLog print grpc stream call log.
func PrintGrpcStreamLog(kt *kit.Kit, srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {

	log.Info(kt, "grpc stream", "method", info.FullMethod, "isClient", info.IsClientStream, "isServer",
		info.IsServerStream)

	ssWrapper := &serverStreamWrapper{
		kit:          kt,
		ServerStream: ss,
		sendMsg: func(m any) error {
			start := time.Now()
			body, _ := json.Marshal(m)
			log.Info(kt, "grpc stream send message", "appCode", kt.AppCode, "user", kt.User, "method",
				info.FullMethod, "body", truncateBody(body), "cost", time.Since(start))
			return ss.SendMsg(m)
		},
		recvMsg: func(m any) error {
			start := time.Now()
			if err := ss.RecvMsg(m); err != nil {
				if errors.Is(err, io.EOF) {
					log.Info(kt, "grpc stream receive message end", "method", info.FullMethod)
					return err
				}
				log.Error(kt, "grpc stream receive message failed", "method", info.FullMethod, log.E(err))
				return err
			}

			body, _ := json.Marshal(m)
			defer func() {
				log.Info(kt, "grpc receive message", "appCode", kt.AppCode, "user", kt.User, "method",
					info.FullMethod, "body", truncateBody(body), "cost", time.Since(start))
			}()

			// validate request
			if v, ok := m.(validator.Validator); ok {
				if err := v.Validate(kt); err != nil {
					return cerr.Wrap(cerr.INVALID_REQUEST, err)
				}
			}
			return nil
		},
	}
	return handler(srv, ssWrapper)
}
