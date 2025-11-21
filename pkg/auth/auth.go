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

// Package auth defines authorization related operations.
package auth

import (
	"context"

	"google.golang.org/grpc"

	"github.com/TencentBlueKing/bk-cmdb/pkg/auth/meta"
	grpccli "github.com/TencentBlueKing/bk-cmdb/pkg/client/grpc"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	authpb "github.com/TencentBlueKing/bk-cmdb/pkg/proto/auth-server"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// Authorizer defines authorize related operations.
type Authorizer interface {
	// Authorize checks whether the user has permission to operate resources.
	Authorize(ctx context.Context, resources ...meta.ResourceAttribute) ([]meta.Decision, error)
	// ListAuthorizedResources lists the resources that the user has permission to operate.
	ListAuthorizedResources(ctx context.Context, opts *meta.ListAuthResOptions) (*meta.AuthResInfo, error)
}

// NewAuthorizer creates a new authorizer.
func NewAuthorizer(ctx context.Context, sd sd.Discovery, tls *config.TLSConfig) (Authorizer, error) {
	opt := &grpccli.Options{
		ServiceName: config.AuthServer,
		TLSConf:     tls,
		Builder:     sd,
	}
	conn, err := grpccli.NewGrpcClient(ctx, opt)
	if err != nil {
		log.Error(ctx, "new grpc client failed", log.E(err))
		return nil, err
	}

	return NewAuthorizerWithCli(conn), nil
}

// NewAuthorizerWithCli creates a new authorizer with grpc client.
func NewAuthorizerWithCli(conn *grpc.ClientConn) Authorizer {
	return &authorizer{
		client: authpb.NewAuthClient(conn),
	}
}

// authorizer is the cmdb authorizer.
type authorizer struct {
	// client is the auth server grpc client.
	client authpb.AuthClient
}

// Authorize checks whether the user has permission to operate resources.
func (a *authorizer) Authorize(ctx context.Context, resources ...meta.ResourceAttribute) ([]meta.Decision, error) {
	authReq := &authpb.AuthorizeReq{
		Resources: make([]*authpb.ResourceAttribute, len(resources)),
	}
	for i, resource := range resources {
		authReq.Resources[i] = authpb.ConvertToPBAuthAttr(&resource)
	}

	authRes, err := a.client.Authorize(ctx, authReq)
	if err != nil {
		log.Error(ctx, "authorize failed", log.E(err), "req", authReq)
		return nil, err
	}

	decisions := make([]meta.Decision, len(authRes.Decisions))
	for i, decision := range authRes.Decisions {
		decisions[i] = meta.Decision{Authorized: decision.Authorized}
	}
	return decisions, nil
}

// ListAuthorizedResources lists the resources that the user has permission to operate.
func (a *authorizer) ListAuthorizedResources(ctx context.Context, opts *meta.ListAuthResOptions) (*meta.AuthResInfo,
	error) {

	req := &authpb.ListAuthResReq{
		ResourceType: string(opts.ResourceType),
		Action:       string(opts.Action),
	}
	resp, err := a.client.ListAuthorizedResources(ctx, req)
	if err != nil {
		log.Error(ctx, "list authorized resources failed", log.E(err), "req", req)
		return nil, err
	}

	return &meta.AuthResInfo{
		IDs:   resp.Ids,
		IsAny: resp.IsAny,
	}, nil
}
