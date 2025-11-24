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
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/oklog/run"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/TencentBlueKing/bk-cmdb/pkg/healthz"
	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest/middleware"
	"github.com/TencentBlueKing/bk-cmdb/pkg/trace"
)

// runHTTPServer runs http server.
func runHTTPServer(ctx context.Context, opts *RunOptions, g *run.Group, tlsConf *tls.Config, tlsEnabled bool) {
	// generate http handlers
	router := newRouter(opts)

	// use extra port for grpc gateway
	addr := net.JoinHostPort(opts.Server.IP, strconv.Itoa(opts.Server.HttpPort))
	svr := http.Server{Addr: addr}

	g.Add(func() error {
		// if tls is enabled, listen and serve http requests with tls
		if tlsEnabled {
			log.Info(ctx, "listen and serve http requests with tls", "addr", addr)
			svr.Handler = router
			svr.TLSConfig = tlsConf
			if err := svr.ListenAndServeTLS("", ""); err != nil {
				log.Error(ctx, "listen and serve http requests with tls failed", log.E(err))
				return err
			}
			return nil
		}

		// if tls is not enabled, listen and serve http requests without tls
		log.Info(ctx, "listen and serve http requests without tls", "addr", addr)
		svr.Handler = h2c.NewHandler(router, &http2.Server{})
		return svr.ListenAndServe()
	}, func(err error) {
		// shutdown http server
		st := time.Now()
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()

		if e := svr.Shutdown(timeoutCtx); e != nil {
			log.Error(ctx, "shutdown http server with error", "reason", err, "duration", time.Since(st), log.E(e))
			return
		}
		log.Info(ctx, "shutdown http server done", "reason", err, "duration", time.Since(st))
	})
}

// newRouter creates a new server router instance with common handlers and middlewares.
func newRouter(opts *RunOptions) http.Handler {
	r := rest.NewRouter()

	// add http middlewares
	r.Use(middleware.Recoverer)

	// register http handlers
	if opts.Router != nil {
		r.Group(func(r rest.Router) {
			r.Use(trace.Middleware)
			r.Mount("/", opts.Router)
		})
	}

	// healthz
	r.Get("/healthz", healthz.HealthzHandler)
	r.Get("/-/healthy", healthz.HealthyHandler)
	r.Get("/-/ready", healthz.ReadyHandler)

	// pprof
	r.Mount("/debug", middleware.Profiler())

	// metrics
	r.Get("/metrics", opts.Metrics.ServeHTTP)

	return r
}
