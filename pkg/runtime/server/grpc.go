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

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/oklog/run"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/runtime/server/middleware"
)

// GrpcOptions defines the options for grpc server.
type GrpcOptions struct {
	// RegisterServer registers grpc server.
	RegisterServer func(server *grpc.Server)
}

// runGrpcServer runs grpc server.
func runGrpcServer(ctx context.Context, opts *RunOptions, g *run.Group, tlsConf *tls.Config, tlsEnabled bool) error {
	if opts.GrpcOpts == nil {
		return nil
	}

	if opts.GrpcOpts.RegisterServer == nil {
		return errors.New("grpc register server function is not set")
	}

	// create grpc server with prometheus and open-telemetry and cmdb interceptors
	grpcOpts := []grpc.ServerOption{
		grpc.UnaryInterceptor(middleware.ConvGrpcUnaryInterceptor(middleware.PrintGrpcLog,
			opts.Metrics.GrpcUnaryServerInterceptor)),
		grpc.StreamInterceptor(middleware.ConvGrpcStreamInterceptor(middleware.PrintGrpcStreamLog,
			opts.Metrics.GrpcStreamServerInterceptor)),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	if tlsEnabled {
		grpcOpts = append(grpcOpts, grpc.Creds(credentials.NewTLS(tlsConf)))
	}
	server := grpc.NewServer(grpcOpts...)

	// register grpc server
	opts.GrpcOpts.RegisterServer(server)

	// Register reflection service on gRPC server.
	reflection.Register(server)

	// listen and serve grpc requests
	addr := net.JoinHostPort(opts.Server.IP, strconv.Itoa(opts.Server.RpcPort))
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error(ctx, "listen grpc address failed", "addr", addr, log.E(err))
		return err
	}

	g.Add(func() error {
		log.Info(ctx, "listen and serve grpc requests", "addr", addr)
		return server.Serve(listener)
	}, func(err error) {
		// shutdown grpc server
		st := time.Now()
		server.GracefulStop()
		log.Info(ctx, "shutdown grpc server done", "reason", err, "duration", time.Since(st))
	})

	return nil
}
