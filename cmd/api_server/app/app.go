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
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/oklog/run"
	"github.com/spf13/cobra"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/TencentBlueKing/bk-cmdb/cmd/api_server/options"
	"github.com/TencentBlueKing/bk-cmdb/cmd/api_server/service"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/runtime/cli"
)

// NewAPIServerCommand creates a *cobra.Command object with default parameters
func NewAPIServerCommand() *cobra.Command {
	opts := options.NewOptions()
	handlerOpts := log.NewHandlerOptions()

	cmd := &cobra.Command{
		Use:   "apiserver",
		Short: "A http service for handle unified http request",
		RunE: func(c *cobra.Command, args []string) error {
			if err := handlerOpts.Validate(); err != nil {
				return err
			}

			handler := log.NewContextualHandler(handlerOpts)
			log.SetDefault(handler)

			return runHTTPServer(c.Context(), opts)
		},
	}

	fs := cmd.Flags()
	opts.AddFlags(fs)
	handlerOpts.AddFlags(fs)
	fs.AddGoFlagSet(goflag.CommandLine)

	return cmd
}

func runHTTPServer(ctx context.Context, opts *options.Options) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var g run.Group

	router := service.NewRouter()
	registerHTTPServer(ctx, &g, router, opts)

	// 监听信号
	g.Add(func() error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-cli.SignalChan:
			log.Warn(ctx, "Signal received", "signal", s)
			return fmt.Errorf("%w %s received", cli.ErrSignal, s)
		}
	}, func(err error) {
		cancel()
	})

	// block here
	return g.Run()
}

func registerHTTPServer(ctx context.Context, g *run.Group, router http.Handler, opts *options.Options) {
	addr := net.JoinHostPort(opts.Address, strconv.Itoa(opts.Port))

	h2cHandler := h2c.NewHandler(router, &http2.Server{})
	svr := http.Server{Addr: addr, Handler: h2cHandler}

	g.Add(func() error {
		log.Info(ctx, "listening for http requests and metrics", "addr", addr)
		return svr.ListenAndServe()

	}, func(err error) {
		st := time.Now()
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()

		if e := svr.Shutdown(timeoutCtx); e != nil {
			log.Error(ctx, "shutdown http server with error",
				"reason", err, "duration", time.Since(st), log.E(err))
			return
		}
		log.Info(ctx, "shutdown http server done", "reason", err, "duration", time.Since(st))
	})
}
