// Package service TODO
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
package service

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/operation_server/app/options"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// OperationServer TODO
type OperationServer struct {
	*backbone.Engine
	Config      *options.Config
	AuthManager *extensions.AuthManager
}

// WebService TODO
func (o *OperationServer) WebService() *restful.Container {

	getErrFunc := func() errors.CCErrorIf {
		return o.Engine.CCErr
	}

	api := new(restful.WebService)
	api.Path("/operation/v3").Filter(o.Engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).
		Produces(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	o.newOperationService(api)
	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	container.Add(api)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(o.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

func (o *OperationServer) newOperationService(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  o.Engine.CCErr,
		Language: o.Engine.Language,
	})

	utility.AddToRestfulWebService(web)
}

// Healthz TODO
func (o *OperationServer) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := o.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// coreservice
	coreSrv := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_CORESERVICE}
	if _, err := o.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_CORESERVICE); err != nil {
		coreSrv.IsHealthy = false
		coreSrv.Message = err.Error()
	}
	meta.Items = append(meta.Items, coreSrv)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "operation server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_OPERATION,
		HealthMeta: meta,
		AtTime:     metadata.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	answer.SetCommonResponse()
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}

// OnOperationConfigUpdate TODO
func (o *OperationServer) OnOperationConfigUpdate(previous, current cc.ProcessConfig) {
	o.Config = &options.Config{}
}
