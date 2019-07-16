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
	"configcenter/src/common/metadata"
	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/source_controller/hostcontroller/logics"
	"configcenter/src/storage/dal"
)

type Service struct {
	Core     *backbone.Engine
	Instance dal.RDB
	Cache    *redis.Client
	Logics   *logics.Logics
}

func (s *Service) WebService() *restful.Container {
	container := restful.NewContainer()

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.Core.CCErr
	}
	api.Path("/host/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	api.Route(api.POST("/hosts/favorites/{user}").To(s.AddHostFavourite))
	api.Route(api.PUT("/hosts/favorites/{user}/{id}").To(s.UpdateHostFavouriteByID))
	api.Route(api.DELETE("/hosts/favorites/{user}/{id}").To(s.DeleteHostFavouriteByID))
	api.Route(api.POST("/hosts/favorites/search/{user}").To(s.GetHostFavourites))
	api.Route(api.POST("/hosts/favorites/search/{user}/{id}").To(s.GetHostFavouriteByID))
	api.Route(api.POST("/history/{user}").To(s.AddHistory))
	api.Route(api.GET("/history/{user}/{start}/{limit}").To(s.GetHistorys))
	api.Route(api.GET("/host/{bk_host_id}").To(s.GetHostByID))
	api.Route(api.POST("/hosts/search").To(s.GetHosts))
	api.Route(api.POST("/insts").To(s.AddHost))
	api.Route(api.GET("/host/snapshot/{bk_host_id}").To(s.GetHostSnap))
	api.Route(api.POST("/meta/hosts/modules/search").To(s.GetHostModulesIDs))
	api.Route(api.POST("/meta/hosts/modules").To(s.AddModuleHostConfig))
	api.Route(api.DELETE("/meta/hosts/modules").To(s.DelModuleHostConfig))
	api.Route(api.DELETE("/meta/hosts/defaultmodules").To(s.DelDefaultModuleHostConfig))
	api.Route(api.PUT("/meta/hosts/resource").To(s.MoveHost2ResourcePool))
	api.Route(api.POST("/meta/hosts/assign").To(s.AssignHostToApp))
	api.Route(api.POST("/meta/hosts/module/config/search").To(s.GetModulesHostConfig))
	api.Route(api.POST("/userapi").To(s.AddUserConfig))
	api.Route(api.PUT("/userapi/{bk_biz_id}/{id}").To(s.UpdateUserConfig))
	api.Route(api.DELETE("/userapi/{bk_biz_id}/{id}").To(s.DeleteUserConfig))
	api.Route(api.POST("/userapi/search").To(s.GetUserConfig))
	api.Route(api.GET("/userapi/detail/{bk_biz_id}/{id}").To(s.UserConfigDetail))
	api.Route(api.POST("/usercustom/{bk_user}").To(s.AddUserCustom))
	api.Route(api.PUT("/usercustom/{bk_user}/{id}").To(s.UpdateUserCustomByID))
	api.Route(api.POST("/usercustom/user/search/{bk_user}").To(s.GetUserCustomByUser))
	api.Route(api.POST("/usercustom/default/search/{bk_user}").To(s.GetDefaultUserCustom))
	api.Route(api.POST("/transfer/host/default/module").To(s.TransferHostToDefaultModuleConfig))

	api.Route(api.POST("/host/lock").To(s.LockHost))
	api.Route(api.DELETE("/host/lock").To(s.UnlockHost))
	api.Route(api.POST("/host/lock/search").To(s.QueryLockHost))

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
	resp.WriteEntity(answer)
}
