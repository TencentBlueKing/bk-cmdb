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

package service

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/TencentBlueKing/bk-cmdb/pkg/healthz"
	"github.com/TencentBlueKing/bk-cmdb/pkg/i18n"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
)

// NewRouter ...
func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(i18n.Middleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", healthz.HealthzHandler)
	r.Get("/-/healthy", healthz.HealthyHandler)
	r.Get("/-/ready", healthz.ReadyHandler)

	// pprof
	r.Mount("/debug", middleware.Profiler())

	// metrics 配置
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	svr := service{}
	r.Post("/user/info", rest.Handle(svr.UserInfo))

	return r
}
