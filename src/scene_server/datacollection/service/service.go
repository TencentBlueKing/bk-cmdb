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
	"gopkg.in/redis.v5"
)

type Service struct {
	*backbone.Engine
	db      dal.RDB
	cache   *redis.Client
	snapcli *redis.Client
	disCli  *redis.Client
	netCli  *redis.Client
	*logics.Logics
}

func (s *Service) SetDB(db dal.RDB) {
	s.db = db
}

func (s *Service) SetCache(db *redis.Client) {
	s.cache = db
}

func (s *Service) SetSnapcli(db *redis.Client) {
	s.snapcli = db
}

func (s *Service) SetDisCli(db *redis.Client) {
	s.disCli = db
}

func (s *Service) SetNetCli(db *redis.Client) {
	s.netCli = db
}

func (s *Service) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}

	api.Path("/collector/v3").Filter(s.Engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	api.Route(api.POST("/netcollect/device/action/create").To(s.CreateDevice))
	api.Route(api.POST("/netcollect/device/{device_id}/action/update").To(s.UpdateDevice))
	api.Route(api.POST("/netcollect/device/action/batch").To(s.BatchCreateDevice))
	api.Route(api.POST("/netcollect/device/action/search").To(s.SearchDevice))
	api.Route(api.DELETE("/netcollect/device/action/delete").To(s.DeleteDevice))

	api.Route(api.POST("/netcollect/property/action/create").To(s.CreateProperty))
	api.Route(api.POST("/netcollect/property/{netcollect_property_id}/action/update").To(s.UpdateProperty))
	api.Route(api.POST("/netcollect/property/action/batch").To(s.BatchCreateProperty))
	api.Route(api.POST("/netcollect/property/action/search").To(s.SearchProperty))
	api.Route(api.DELETE("/netcollect/property/action/delete").To(s.DeleteProperty))

	api.Route(api.POST("/netcollect/summary/action/search").To(s.SearchReportSummary))
	api.Route(api.POST("/netcollect/report/action/search").To(s.SearchReport))
	api.Route(api.POST("/netcollect/report/action/confirm").To(s.ConfirmReport))
	api.Route(api.POST("/netcollect/history/action/search").To(s.SearchHistory))

	api.Route(api.POST("/netcollect/collector/action/search").To(s.SearchCollector))
	api.Route(api.POST("/netcollect/collector/action/update").To(s.UpdateCollector))
	api.Route(api.POST("/netcollect/collector/action/discover").To(s.DiscoverNetDevice))

	container.Add(api)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.Healthz))
	container.Add(healthzAPI)

	return container
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

	// mongodb
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityMongo, s.db.Ping()))

	// redis
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityRedis, s.cache.Ping().Err()))
	if s.snapcli != nil {
		meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityRedis+" - snapcli", s.snapcli.Ping().Err()))
	}
	if s.disCli != nil {
		meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityRedis+" - disCli", s.disCli.Ping().Err()))
	}
	if s.netCli != nil {
		meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityRedis+" - netCli", s.netCli.Ping().Err()))
	}

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
	_ = resp.WriteEntity(answer)
}
