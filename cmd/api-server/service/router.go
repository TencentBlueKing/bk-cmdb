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
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/TencentBlueKing/bk-cmdb/pkg/log"
	"github.com/TencentBlueKing/bk-cmdb/pkg/rest"
	"github.com/TencentBlueKing/bk-cmdb/pkg/runtime/server/middleware"
)

// NewRouter creates a new api-server router.
func (s *Service) NewRouter(ctx context.Context) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(I18nMiddleware)

	r.Use(Authentication) // 统一鉴权中间件
	r.Use(middleware.ConvHttpMiddleware(middleware.PrintHttpLog, s.metric.HTTPMiddleware))

	// register grpc gateway http handlers
	grpcMux, err := s.newGrpcMux(ctx)
	if err != nil {
		log.Error(ctx, "new grpc mux failed", log.E(err), "addr", grpcMux)
		return nil, err
	}
	r.Mount("/", grpcMux)

	// register restful http handlers
	r.Post("/api/v4/user/info", rest.Handle(s.UserInfo))
	r.Post("/api/v4/authorized/users", rest.Handle(s.ListAuthorizedUsers))
	r.Post("/api/v4/translations/reload", rest.Handle(s.ReloadTranslation))
	return r, nil
}
