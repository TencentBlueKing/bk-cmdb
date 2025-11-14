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

// Package service defines auth-server service.
package service

import (
	"context"
	"errors"
	"io"

	"google.golang.org/grpc"

	"github.com/TencentBlueKing/bk-cmdb/pkg/kit"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/metrics"
	authpb "github.com/TencentBlueKing/bk-cmdb/pkg/proto/auth-server"
)

// Service is auth-server service.
type Service struct {
	authpb.UnimplementedAuthServer
	metric *metrics.Service
}

// NewService creates a new service.
func NewService(metric *metrics.Service) (*Service, error) {
	return &Service{metric: metric}, nil
}

// Authorize checks whether the user has permission to operate resources.
func (s *Service) Authorize(ctx context.Context, req *authpb.AuthorizeReq) (*authpb.AuthorizeResp, error) {
	kt := kit.GetGrpcKit(ctx)

	decisions := make([]*authpb.Decision, 0, len(req.Resources))
	for _, resource := range req.Resources {
		authorized := false
		if resource.Basic.Type == "skip" || kt.User == "admin" {
			authorized = true
		}

		decisions = append(decisions, &authpb.Decision{Authorized: authorized})
	}
	return &authpb.AuthorizeResp{Decisions: decisions}, nil
}

// ListAuthorizedResources lists the resources that the user has permission to operate.
func (s *Service) ListAuthorizedResources(ctx context.Context, req *authpb.ListAuthResReq) (*authpb.ListAuthResResp,
	error) {

	if req.ResourceType == "skip" {
		return &authpb.ListAuthResResp{IsAny: true}, nil
	}
	return &authpb.ListAuthResResp{Ids: []string{"1", "2", "3"}}, nil
}

// AuthorizeStream this is only a demo for streaming.
func (s *Service) AuthorizeStream(stream grpc.BidiStreamingServer[authpb.AuthorizeReq, authpb.AuthorizeResp]) error {
	kt := kit.GetGrpcKit(stream.Context())

	for {
		req, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			log.Error(kt, "stream receive failed", log.E(err))
			return err
		}

		resp, err := s.Authorize(kt, req)
		if err != nil {
			log.Error(kt, "authorize failed", log.E(err))
			return err
		}

		if err = stream.Send(resp); err != nil {
			log.Error(kt, "stream send failed", log.E(err))
			return err
		}
	}
}
