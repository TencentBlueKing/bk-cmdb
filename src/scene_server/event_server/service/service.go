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
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/storage"

	"github.com/emicklei/go-restful"

	redis "gopkg.in/redis.v5"
)

type Service struct {
	*backbone.Engine
	db    storage.DI
	cache *redis.Client
}

func (s *Service) SetDB(db storage.DI) {
	s.db = db
}

func (s *Service) SetCache(db *redis.Client) {
	s.cache = db
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/event/v3").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/subscribe/search/{ownerID}/{appID}").To(s.Query))
	ws.Route(ws.POST("/subscribe/ping").To(s.Ping))
	ws.Route(ws.POST("/subscribe/telnet").To(s.Telnet))
	ws.Route(ws.POST("/subscribe/{ownerID}/{appID}").To(s.Subscribe))
	ws.Route(ws.DELETE("/subscribe/{ownerID}/{appID}/{subscribeID}").To(s.UnSubscribe))
	ws.Route(ws.PUT("/subscribe/{ownerID}/{appID}/{subscribeID}").To(s.Rebook))

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

	// mongodb
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityServicediscover, s.db.Ping()))

	// redis
	meta.Items = append(meta.Items, metric.NewHealthItem(types.CCFunctionalityServicediscover, s.cache.Ping().Err()))

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "event server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_EVENTSERVER,
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
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}
