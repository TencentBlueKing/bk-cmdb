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
	"context"
	"time"

	"google.golang.org/grpc"

	cerr "github.com/TencentBlueKing/bk-cmdb/pkg/errors"
	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/validator"
)

// GrpcUnaryInterceptor is the grpc unary server interceptor with kit.
type GrpcUnaryInterceptor func(kt *kit.Kit, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (
	any, error)

// ConvGrpcUnaryInterceptor returns a grpc unary server interceptor with the given cmdb grpc interceptors.
func ConvGrpcUnaryInterceptor(interceptors ...GrpcUnaryInterceptor) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// TODO confirm request timeout
		ctx, cancel := context.WithTimeout(ctx, 25*time.Second)
		defer cancel()

		kt, err := kit.NewKitFromGrpcCtx(ctx)
		if err != nil {
			log.Error(ctx, "new grpc kit failed", log.E(err))
			return nil, cerr.Wrap(cerr.INVALID_REQUEST, err)
		}

		nextHandler := handler
		handler = func(ctx context.Context, req any) (any, error) {
			// validate request
			if v, ok := req.(validator.Validator); ok {
				if err = v.Validate(ctx); err != nil {
					return nil, cerr.Wrap(cerr.INVALID_REQUEST, err)
				}
			}
			return nextHandler(ctx, req)
		}

		// call all interceptors
		for i := len(interceptors) - 1; i >= 0; i-- {
			nextHandler := handler
			interceptor := interceptors[i]
			handler = func(ctx context.Context, req any) (any, error) {
				return interceptor(kt, req, info, nextHandler)
			}
		}
		return handler(kt, req)
	}
}

// GrpcStreamInterceptor is the grpc stream server interceptor with kit.
type GrpcStreamInterceptor func(kt *kit.Kit, srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error

// GrpcMsgHandler is the grpc stream message handler.
type GrpcMsgHandler func(any) error

// ConvGrpcStreamInterceptor returns a grpc stream server interceptor with the given cmdb grpc interceptors.
func ConvGrpcStreamInterceptor(interceptors ...GrpcStreamInterceptor) grpc.StreamServerInterceptor {
	return func(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// TODO confirm request timeout
		ctx, cancel := context.WithTimeout(ss.Context(), 25*time.Second)
		defer cancel()

		kt, err := kit.NewKitFromGrpcCtx(ctx)
		if err != nil {
			log.Error(ss.Context(), "new grpc kit failed", log.E(err))
			return cerr.Wrap(cerr.INVALID_REQUEST, err)
		}

		// call all interceptors, generate grpc stream handler and server stream wrapper
		for i := len(interceptors) - 1; i >= 0; i-- {
			interceptor := interceptors[i]
			nextHandler := handler

			handler = func(srv any, stream grpc.ServerStream) error {
				return interceptor(kt, srv, stream, info, nextHandler)
			}
		}

		return handler(srv, ss)
	}
}

// serverStreamWrapper wraps grpc.ServerStream allowing each Sent/Recv of message to increment counters.
type serverStreamWrapper struct {
	kit *kit.Kit
	grpc.ServerStream
	sendMsg GrpcMsgHandler
	recvMsg GrpcMsgHandler
}

// NewServerStreamWrapper creates a new server stream wrapper.
func NewServerStreamWrapper(kt *kit.Kit, ss grpc.ServerStream, method string,
	sendMsg, recvMsg GrpcMsgHandler) grpc.ServerStream {
	return &serverStreamWrapper{
		kit:          kt,
		ServerStream: ss,
		sendMsg:      sendMsg,
		recvMsg:      recvMsg,
	}
}

// Context return the server stream's context
func (s *serverStreamWrapper) Context() context.Context {
	return s.kit
}

// SendMsg send message.
func (s *serverStreamWrapper) SendMsg(m any) error {
	return s.sendMsg(m)
}

// RecvMsg receive message.
func (s *serverStreamWrapper) RecvMsg(m any) error {
	return s.recvMsg(m)
}
