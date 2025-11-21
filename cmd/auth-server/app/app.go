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

// Package app defines the auth server main entry.
package app

import (
	"context"
	"flag"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/TencentBlueKing/bk-cmdb/cmd/auth-server/options"
	"github.com/TencentBlueKing/bk-cmdb/cmd/auth-server/service"
	cc "github.com/TencentBlueKing/bk-cmdb/pkg/config-center"
	opt "github.com/TencentBlueKing/bk-cmdb/pkg/config-center/options"
	authpb "github.com/TencentBlueKing/bk-cmdb/pkg/proto/auth-server"
	"github.com/TencentBlueKing/bk-cmdb/pkg/runtime/server"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// NewAuthServerCommand creates a *cobra.Command object with default parameters
func NewAuthServerCommand() *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "auth-server",
		Short: "authorization server",
		RunE: func(c *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			ctx := c.Context()
			svr, err := newAuthServer(ctx, opts)
			if err != nil {
				return err
			}

			runOpts := &server.RunOptions{
				CommonInfo: svr.serverInfo,
				Registry:   svr.sd,
				GrpcOpts: &server.GrpcOptions{
					RegisterServer: func(server *grpc.Server) {
						authpb.RegisterAuthServer(server, svr.service)
					},
				},
				Router:    svr.service.NewRouter(),
				Finalizer: svr.finalizer,
			}

			return server.Run(ctx, runOpts)
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs, true)
	fs.AddGoFlagSet(flag.CommandLine)

	return cmd
}

// authServer defines the auth-server.
type authServer struct {
	serverInfo *server.CommonInfo
	sd         sd.Registry
	service    *service.Service
}

// newAuthServer creates a new auth-server.
func newAuthServer(ctx context.Context, opts *opt.Options) (*authServer, error) {
	newSvrOpts := &server.Options{
		Options: opts,
		// TODO: add db config when dao is ready
		ConfOpts: &server.ConfigOptions{NeededConfigs: []cc.ConfigType{cc.CommonConfType}},
	}
	svrInfo, sd, err := server.NewServerInfoWithRegistry(ctx, newSvrOpts)
	if err != nil {
		return nil, err
	}

	// create auth service
	svc, err := service.NewService(svrInfo.Metrics)
	if err != nil {
		return nil, err
	}

	return &authServer{
		serverInfo: svrInfo,
		sd:         sd,
		service:    svc,
	}, nil
}

// finalizer defines the finalize logics before exit.
func (a *authServer) finalizer() {}
