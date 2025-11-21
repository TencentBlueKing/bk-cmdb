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

// Package grpc defines run cmdb client related logics.
package grpc

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"

	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// Options defines the options for grpc client.
type Options struct {
	// ServiceName is the service name.
	ServiceName config.ServiceName
	// TLSConf is the tls config.
	TLSConf *config.TLSConfig
	// Builder is the grpc resolver builder.
	Builder resolver.Builder
}

// NewGrpcClient creates a new grpc client.
func NewGrpcClient(ctx context.Context, opts *Options) (*grpc.ClientConn, error) {
	grpcOpts := []grpc.DialOption{
		// use custom resolver
		grpc.WithResolvers(opts.Builder),
		// use round-robin load balance policy
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		// use open telemetry stats handler for grpc client
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}

	// generate grpc tls dial options
	clientTls, tlsEnabled, err := opts.TLSConf.ToClientConf()
	if err != nil {
		log.Error(ctx, "convert tls config to client config failed", log.E(err), "conf", opts.TLSConf)
		return nil, err
	}

	if tlsEnabled {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(credentials.NewTLS(clientTls)))
	} else {
		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// new grpc client
	conn, err := grpc.NewClient(sd.GenGrpcServiceDiscoveryPath(opts.ServiceName), grpcOpts...)
	if err != nil {
		log.Error(ctx, "new grpc client failed", log.E(err))
		return nil, err
	}

	return conn, nil
}
