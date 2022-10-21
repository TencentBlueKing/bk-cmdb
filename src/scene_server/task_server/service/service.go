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

// Package service TODO
package service

import (
	"context"
	"net/http"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/task_server/app/options"
	"configcenter/src/scene_server/task_server/logics"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service TODO
type Service struct {
	*options.Config
	*backbone.Engine
	disc    discovery.DiscoveryInterface
	CacheDB redis.Client
	DB      dal.RDB
	Logics  *logics.Logics
}

// WebService TODO
func (s *Service) WebService() *restful.Container {

	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}
	api.Path("/task/v3").Filter(s.Engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	s.addAPIService(api)
	container.Add(api)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	container.Add(commonAPI)

	return container
}

func (s *Service) addAPIService(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/create", Handler: s.CreateTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/task", Handler: s.CreateTaskBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/findmany/list/{name}", Handler: s.ListTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/findmany/list/latest/{name}", Handler: s.ListLatestTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/findone/detail/{task_id}", Handler: s.DetailTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/deletemany", Handler: s.DeleteTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/task/set/status/sucess/id/{task_id}/sub_id/{sub_task_id}", Handler: s.StatusToSuccess})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/task/set/status/failure/id/{task_id}/sub_id/{sub_task_id}", Handler: s.StatusToFailure})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/latest/sync_status",
		Handler: s.ListLatestSyncStatus})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/sync_status_history",
		Handler: s.ListSyncStatusHistory})

	utility.AddToRestfulWebService(web)

}

// Healthz TODO
func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// mongodb status
	mongoItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityMongo}
	if s.DB == nil {
		mongoItem.IsHealthy = false
		mongoItem.Message = "not connected"
	} else if err := s.DB.Ping(); err != nil {
		mongoItem.IsHealthy = false
		mongoItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, mongoItem)

	// redis status
	redisItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityRedis}
	if s.CacheDB == nil {
		redisItem.IsHealthy = false
		redisItem.Message = "not connected"
	} else if err := s.CacheDB.Ping(context.Background()).Err(); err != nil {
		redisItem.IsHealthy = false
		redisItem.Message = err.Error()
	}

	// zk health status
	zkItem := metric.HealthItem{IsHealthy: true, Name: types.CCFunctionalityServicediscover}
	if err := s.Engine.Ping(); err != nil {
		zkItem.IsHealthy = false
		zkItem.Message = err.Error()
	}
	meta.Items = append(meta.Items, zkItem)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "task server is unhealthy"
			break
		}
	}

	info := metric.HealthInfo{
		Module:     types.CC_MODULE_TASK,
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
	_ = resp.WriteEntity(answer)
}
