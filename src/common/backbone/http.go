/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package backbone

import (
	"context"
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof" // pprof TODO
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"configcenter/src/common/blog"
	"configcenter/src/common/ssl"
	"configcenter/src/common/zkclient"
)

// ListenAndServe TODO
func ListenAndServe(c Server, svcDisc ServiceRegisterInterface, cancel context.CancelFunc) error {
	handler := c.Handler
	if c.PProfEnabled {
		rootMux := http.NewServeMux()
		rootMux.HandleFunc("/", c.Handler.ServeHTTP)
		rootMux.Handle("/debug/", http.DefaultServeMux)
		handler = rootMux
	}
	server := &http.Server{
		Addr:    net.JoinHostPort(c.ListenAddr, strconv.FormatUint(uint64(c.ListenPort), 10)),
		Handler: handler,
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM)
	go func() {
		for {
			select {
			case sig := <-exit:
				blog.Infof("receive signal %v, begin to shutdown", sig)
				svcDisc.Cancel()
				if err := svcDisc.ClearRegisterPath(); err != nil && err != zkclient.ErrNoNode {
					break
				}
				time.Sleep(time.Second * 5)
				server.SetKeepAlivesEnabled(false)
				err := server.Shutdown(context.Background())
				if err != nil {
					blog.Fatalf("Could not gracefully shutdown the server: %v \n", err)
				}
				blog.Info("server shutdown done")
				cancel()
				return
			}
		}
	}()

	if !isTLS(c.TLS) {
		blog.Infof("start insecure server on %s:%d", c.ListenAddr, c.ListenPort)
		go func() {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				blog.Fatalf("listen and serve failed, err: %v", err)
			}
		}()
		return nil
	}

	tlsC, err := ssl.ServerTLSVerifyClient(c.TLS.CAFile,
		c.TLS.CertFile,
		c.TLS.KeyFile,
		c.TLS.Password)
	if err != nil {
		return fmt.Errorf("generate tls config failed. err: %v", err)
	}

	server.TLSConfig = tlsC
	blog.Infof("start secure server on %s:%d", c.ListenAddr, c.ListenPort)
	go func() {
		if err := server.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			blog.Fatalf("listen and serve failed, err: %v", err)
		}
	}()

	return nil
}
