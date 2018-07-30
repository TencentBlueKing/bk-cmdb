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
	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/source_controller/auditcontroller/logics"
	"configcenter/src/storage"
)

type Service struct {
	*backbone.Engine
	*logics.Logics
	Instance storage.DI
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/audit/{version}").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.POST("/host/{owner_id}/{biz_id}/{user}").To(s.AddHostLog))
	ws.Route(ws.POST("/hosts/{owner_id}/{biz_id}/{user}").To(s.AddHostLogs))
	ws.Route(ws.POST("/obj/{owner_id}/{biz_id}/{user}").To(s.AddObjectLog))
	ws.Route(ws.POST("/objs/{owner_id}/{biz_id}/{user}").To(s.AddObjectLogs))
	ws.Route(ws.POST("/proc/{owner_id}/{biz_id}/{user}").To(s.AddProcLog))
	ws.Route(ws.POST("/procs/{owner_id}/{biz_id}/{user}").To(s.AddProcLogs))
	ws.Route(ws.POST("/module/{owner_id}/{biz_id}/{user}").To(s.AddModuleLog))
	ws.Route(ws.POST("/modules/{owner_id}/{biz_id}/{user}").To(s.AddModuleLogs))
	ws.Route(ws.POST("/app/{owner_id}/{biz_id}/{user}").To(s.AddAppLog))
	ws.Route(ws.POST("set/{owner_id}/{biz_id}/{user}").To(s.AddSetLog))
	ws.Route(ws.POST("/sets/{owner_id}/{biz_id}/{user}").To(s.AddSetLogs))
	ws.Route(ws.POST("/search").To(s.Get))
	ws.Route(ws.GET("/healthz").To(s.Healthz))

	return ws
}

func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	// mongodb status
	mongoItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityMongo}
	if err := s.Instance.Ping(); err != nil {
		mongoItem.IsHealthy = false
		mongoItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, mongoItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "audit controller is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_AUDITCONTROLLER,
		HealthMeta: meta,
		AtTime:     types.Now(),
	}

	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		Data:    info,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	resp.WriteJson(answer, "application/json")
}
