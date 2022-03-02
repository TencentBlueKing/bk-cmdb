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

package opentelemetry

import (
	"fmt"
	"net/http"

	"configcenter/src/common"

	"github.com/emicklei/go-restful/v3"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/emicklei/go-restful/otelrestful"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// AddOtlpFilter add OpenTelemetry Protocol filter
func AddOtlpFilter(container *restful.Container) {
	if container != nil && openTelemetryCfg.enable {
		container.Filter(otelrestful.OTelFilter(serviceName()))
	}
}

// UseOtlpMiddleware use OpenTelemetry Protocol middleware
func UseOtlpMiddleware(ws *gin.Engine) {
	if ws != nil && openTelemetryCfg.enable {
		ws.Use(otelgin.Middleware(serviceName()))
	}
}

// WrapperTraceClient wrapper client to record trace
func WrapperTraceClient(client *http.Client) {
	if client != nil && openTelemetryCfg.enable {
		client.Transport = otelhttp.NewTransport(client.Transport)
	}
}

func serviceName() string {
	return fmt.Sprintf("%s_%s", "cmdb", common.GetIdentification())
}
