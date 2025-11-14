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

// Package service define apiserver service
package service

import (
	"context"

	"google.golang.org/grpc"

	"github.com/TencentBlueKing/bk-cmdb/pkg/auth"
	grpccli "github.com/TencentBlueKing/bk-cmdb/pkg/client/grpc"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/metrics"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// Service is api-server service.
type Service struct {
	grpcClients map[config.ServiceName]*grpc.ClientConn
	authorizer  auth.Authorizer
	metric      *metrics.Service
}

// NewService creates a new service.
func NewService(ctx context.Context, sd sd.Discovery, tls *config.TLSConfig, metric *metrics.Service) (*Service,
	error) {

	// new grpc clients
	grpcServices := []config.ServiceName{config.AuthServer}
	grpcClients := make(map[config.ServiceName]*grpc.ClientConn)
	for _, service := range grpcServices {
		opt := &grpccli.Options{
			ServiceName: service,
			TLSConf:     tls,
			Builder:     sd,
		}
		conn, err := grpccli.NewGrpcClient(ctx, opt)
		if err != nil {
			log.Error(ctx, "new grpc client failed", log.E(err))
			return nil, err
		}

		grpcClients[service] = conn
	}

	// create authorizer
	authorizer := auth.NewAuthorizerWithCli(grpcClients[config.AuthServer])

	return &Service{
		grpcClients: grpcClients,
		authorizer:  authorizer,
		metric:      metric,
	}, nil
}

// Close closes api service.
func (s *Service) Close() {
	for _, conn := range s.grpcClients {
		_ = conn.Close()
	}
}
