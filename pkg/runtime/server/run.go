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
	"fmt"
	"net/http"
	"time"

	"github.com/oklog/run"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/runtime/cli"
	sd "github.com/TencentBlueKing/bk-cmdb/pkg/service-discovery"
	"github.com/TencentBlueKing/bk-cmdb/pkg/trace"
)

// RunOptions defines the run server options.
type RunOptions struct {
	// CommonInfo is the common server info.
	*CommonInfo
	// Registry is the service registry.
	Registry sd.Registry
	// Router is the http router.
	Router http.Handler
	// GrpcOpts is the grpc server options.
	GrpcOpts *GrpcOptions
	// Finalizer defines the finalize logics before exit.
	Finalizer func()
}

// Run runs the server.
func Run(ctx context.Context, opts *RunOptions) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	g := new(run.Group)
	g.Add(func() error {
		// wait until context is canceled or shutdown signal is received
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s := <-cli.SignalChan:
			log.Warn(ctx, "Signal received", "signal", s)
			return fmt.Errorf("%w %s received", cli.ErrSignal, s)
		}
	}, func(reason error) {
		cancel()

		// deregister service
		if err := opts.Registry.Deregister(context.Background()); err != nil {
			log.Error(context.Background(), "deregister service failed", log.E(err))
		}
	})

	// convert tls config to server tls config
	serverTls, tlsEnabled, err := opts.TLSConf.ToServerConf()
	if err != nil {
		log.Error(ctx, "convert tls config to server config failed", log.E(err), "conf", opts.TLSConf)
		return err
	}

	if opts.GrpcOpts != nil {
		// run grpc server
		if err = runGrpcServer(ctx, opts, g, serverTls, tlsEnabled); err != nil {
			return err
		}
	}

	// run http server
	runHTTPServer(ctx, opts, g, serverTls, tlsEnabled)

	// register service
	if err = opts.Registry.Register(ctx); err != nil {
		return err
	}

	g.Add(func() error {
		// wait until context is canceled
		<-ctx.Done()
		return nil
	}, func(reason error) {
		// do finalize logics before exit
		if opts.Finalizer != nil {
			opts.Finalizer()
		}

		// shutdown trace exporter
		timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer timeoutCancel()

		if err := trace.Shutdown(timeoutCtx); err != nil {
			log.Error(ctx, "shutdown trace exporter with error", "reason", reason, log.E(err))
			return
		}
		log.Info(ctx, "shutdown trace exporter done", "reason", reason)
	})

	return g.Run()
}
