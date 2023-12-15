/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

// Package service defines gin web service with middlewares
package service

import (
	"github.com/gin-gonic/gin"
)

// MiddlewareService is gin web service with middlewares
type MiddlewareService struct {
	ws          *gin.Engine
	middlewares gin.HandlersChain
}

// NewMiddlewareService new MiddlewareService
func NewMiddlewareService(ws *gin.Engine, middlewares ...gin.HandlerFunc) *MiddlewareService {
	return &MiddlewareService{
		ws:          ws,
		middlewares: middlewares,
	}
}

// Post handle post request
func (m *MiddlewareService) Post(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.ws.POST(relativePath, append(m.middlewares, handlers...)...)
}

// Get handle get request
func (m *MiddlewareService) Get(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.ws.GET(relativePath, append(m.middlewares, handlers...)...)
}

// Put handle put request
func (m *MiddlewareService) Put(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.ws.PUT(relativePath, append(m.middlewares, handlers...)...)
}

// Delete handle delete request
func (m *MiddlewareService) Delete(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return m.ws.DELETE(relativePath, append(m.middlewares, handlers...)...)
}
