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

// Package app defines the server main entry.
package app

import (
	"context"
	goflag "flag"

	"github.com/spf13/cobra"

	"github.com/TencentBlueKing/bk-cmdb/cmd/api-server/options"
	"github.com/TencentBlueKing/bk-cmdb/cmd/api-server/service"
	cc "github.com/TencentBlueKing/bk-cmdb/pkg/config-center"
	"github.com/TencentBlueKing/bk-cmdb/pkg/config-center/config"
	opt "github.com/TencentBlueKing/bk-cmdb/pkg/config-center/options"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/runtime/server"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand() *cobra.Command {
	opts := options.NewOptions()

	cmd := &cobra.Command{
		Use:   "apiserver",
		Short: "A http service for handle unified http request",
		RunE: func(c *cobra.Command, args []string) error {
			if err := opts.Validate(); err != nil {
				return err
			}

			ctx := c.Context()
			svr, err := newApiServer(ctx, opts)
			if err != nil {
				log.Error(ctx, "new api server failed", log.E(err), "opt", opts)
				return err
			}

			router, err := svr.service.NewRouter(ctx)
			if err != nil {
				log.Error(ctx, "new router failed", log.E(err))
				return err
			}

			runOpts := &server.RunOptions{
				CommonInfo: svr.serverInfo,
				Registry:   svr.sd,
				Router:     router,
				Finalizer:  svr.service.Close,
			}

			return server.Run(ctx, runOpts)
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs, false)
	fs.AddGoFlagSet(goflag.CommandLine)

	return cmd
}

// apiServer defines the api-server.
type apiServer struct {
	serverInfo *server.CommonInfo
	sd         sd.ServiceDiscovery
	service    *service.Service
}

// newApiServer creates a new api-server.
func newApiServer(ctx context.Context, opts *opt.Options) (*apiServer, error) {
	newSvrOpts := &server.Options{
		Options:  opts,
		ConfOpts: &server.ConfigOptions{NeededConfigs: []cc.ConfigType{cc.CommonConfType}},
	}
	discServices := []config.ServiceName{config.CoreServer, config.AuthServer, config.AdminServer, config.CDCServer,
		config.Collector, config.Governancer, config.TaskServer}

	svrInfo, sd, err := server.NewServerInfoWithSvcDisc(ctx, newSvrOpts, discServices)
	if err != nil {
		return nil, err
	}

	// create api service
	svc, err := service.NewService(ctx, sd, svrInfo.TLSConf, svrInfo.Metrics)
	if err != nil {
		return nil, err
	}

	return &apiServer{
		serverInfo: svrInfo,
		sd:         sd,
		service:    svc,
	}, nil
}
