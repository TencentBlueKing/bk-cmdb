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

package apigw

import (
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/metadata"
	"configcenter/src/web_server/capability"
	"configcenter/src/web_server/middleware"
	"configcenter/src/web_server/middleware/service"

	"github.com/gin-gonic/gin"
)

type apigw struct {
	ws     *service.MiddlewareService
	engine *backbone.Engine
}

// Init api gateway service
func Init(c *capability.Capability) {
	a := &apigw{
		ws:     service.NewMiddlewareService(c.Ws, middleware.ApiGWMiddleware),
		engine: c.Engine,
	}

	a.ws.Post("/demo", a.Demo)
}

// Demo api gateway http handler demo
func (a *apigw) Demo(c *gin.Context) {
	c.JSON(http.StatusOK, metadata.BkBaseResp{Code: common.CCSuccess, Message: common.CCSuccessStr})
}
