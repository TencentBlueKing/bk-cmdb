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
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
    "configcenter/src/scene_server/datacollection/logics"
    "configcenter/src/storage/dal"

    "github.com/emicklei/go-restful"
    redis "gopkg.in/redis.v5"

)

type Service struct {
	*backbone.Engine
	db    dal.RDB
	cache *redis.Client
	*logics.Logics
}

func (s *Service) SetDB(db dal.RDB) {
	s.db = db
}

func (s *Service) SetCache(db *redis.Client) {
	s.cache = db
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}

	ws.Path("/collector/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/netcollect/device/action/create").To(s.CreateDevice))
	ws.Route(ws.POST("/netcollect/device/{device_id}/action/update").To(s.UpdateDevice))
	ws.Route(ws.POST("/netcollect/device/action/batch").To(s.BatchCreateDevice))
	ws.Route(ws.POST("/netcollect/device/action/search").To(s.SearchDevice))
	ws.Route(ws.DELETE("/netcollect/device/action/delete").To(s.DeleteDevice))

	ws.Route(ws.POST("/netcollect/property/action/create").To(s.CreateProperty))
	ws.Route(ws.POST("/netcollect/property/{netcollect_property_id}/action/update").To(s.UpdateProperty))
	ws.Route(ws.POST("/netcollect/property/action/batch").To(s.BatchCreateProperty))
	ws.Route(ws.POST("/netcollect/property/action/search").To(s.SearchProperty))
	ws.Route(ws.DELETE("/netcollect/property/action/delete").To(s.DeleteProperty))

	ws.Route(ws.POST("/netcollect/summary/action/search").To(s.SearchReportSummary))
	ws.Route(ws.POST("/netcollect/report/action/search").To(s.SearchReport))
	ws.Route(ws.POST("/netcollect/report/action/confirm").To(s.ConfirmReport))
	ws.Route(ws.POST("/netcollect/history/action/search").To(s.SearchHistory))

	ws.Route(ws.POST("/netcollect/collector/action/search").To(s.SearchCollector))
	ws.Route(ws.POST("/netcollect/collector/action/update").To(s.UpdateCollector))
	ws.Route(ws.POST("/netcollect/collector/action/discover").To(s.DiscoverNetDevice))

	ws.Path("/collector/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
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

	// topo server
	objCtr := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_TOPO}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_TOPO); err != nil {
		objCtr.IsHealthy = false
		objCtr.Message = err.Error()
	}
	meta.Items = append(meta.Items, objCtr)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "datacolection server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_DATACOLLECTION,
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
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}
