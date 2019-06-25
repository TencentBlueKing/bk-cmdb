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
	"configcenter/src/common/eventclient"
	"configcenter/src/common/metadata"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/source_controller/hostcontroller/logics"
	"configcenter/src/storage/dal"

	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"
)

type Service struct {
	Core     *backbone.Engine
	Instance dal.RDB
	EventC   eventclient.Client
	Cache    *redis.Client
	Logics   *logics.Logics
}

func (s *Service) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.Core.CCErr
	}
	api.Path("/host/v3").Filter(s.Core.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	//Cloud host resource sync
	api.Route(api.POST("/hosts/cloud/add").To(s.AddCloudTask))
	api.Route(api.POST("/hosts/cloud/confirm").To(s.ResourceConfirm))
	api.Route(api.POST("/hosts/cloud/nameCheck").To(s.TaskNameCheck))
	api.Route(api.DELETE("/hosts/cloud/delete/{taskID}").To(s.DeleteCloudTask))
	api.Route(api.POST("/hosts/cloud/search/").To(s.SearchCloudTask))
	api.Route(api.PUT("/hosts/cloud/update").To(s.UpdateCloudTask))
	api.Route(api.DELETE("/hosts/cloud/confirm/delete/{resourceID}").To(s.DeleteConfirm))
	api.Route(api.POST("/hosts/cloud/confirm/search").To(s.SearchConfirm))
	api.Route(api.POST("/hosts/cloud/syncHistory/add").To(s.AddSyncHistory))
	api.Route(api.POST("/hosts/cloud/syncHistory/search").To(s.SearchSyncHistory))
	api.Route(api.POST("/hosts/cloud/confirmHistory/add").To(s.AddConfirmHistory))
	api.Route(api.POST("/hosts/cloud/confirmHistory/search").To(s.SearchConfirmHistory))

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
	if err := s.Core.Ping(); err != nil {
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

	// redis status
	redisItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityRedis}
	if err := s.Cache.Ping().Err(); err != nil {
		redisItem.IsHealthy = false
		redisItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, redisItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "host controller is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_HOST,
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
	resp.WriteAsJson(answer)
}
