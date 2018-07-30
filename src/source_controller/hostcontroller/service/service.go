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
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/source_controller/hostcontroller/logics"
	"configcenter/src/storage"
)

type Service struct {
	Core     *backbone.Engine
	Instance storage.DI
	Cache    *redis.Client
	Logics   logics.Logics
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.Core.CCErr
	}
	ws.Path("/host/v3").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.POST("/hosts/favorites/{user}").To(s.AddHostFavourite))
	ws.Route(ws.PUT("/hosts/favorites/{user}/{id}").To(s.UpdateHostFavouriteByID))
	ws.Route(ws.DELETE("/hosts/favorites/{user}/{id}").To(s.DeleteHostFavouriteByID))
	ws.Route(ws.POST("/hosts/favorites/search/{user}").To(s.GetHostFavourites))
	ws.Route(ws.POST("/hosts/favorites/search/{user}/{id}").To(s.GetHostFavouriteByID))
	ws.Route(ws.POST("/history/{user}").To(s.AddHistory))
	ws.Route(ws.GET("/history/{user}/{start}/{limit}").To(s.GetHistorys))
	ws.Route(ws.GET("/host/{bk_host_id}").To(s.GetHostByID))
	ws.Route(ws.POST("/hosts/search").To(s.GetHosts))
	ws.Route(ws.POST("/insts").To(s.AddHost))
	ws.Route(ws.GET("/host/snapshot/{bk_host_id}").To(s.GetHostSnap))
	ws.Route(ws.POST("/meta/hosts/modules/search").To(s.GetHostModulesIDs))
	ws.Route(ws.POST("/meta/hosts/modules").To(s.AddModuleHostConfig))
	ws.Route(ws.DELETE("/meta/hosts/modules").To(s.DelModuleHostConfig))
	ws.Route(ws.DELETE("/meta/hosts/defaultmodules").To(s.DelDefaultModuleHostConfig))
	ws.Route(ws.PUT("/meta/hosts/resource").To(s.MoveHost2ResourcePool))
	ws.Route(ws.POST("/meta/hosts/assign").To(s.AssignHostToApp))
	ws.Route(ws.POST("/meta/hosts/module/config/search").To(s.GetModulesHostConfig))
	ws.Route(ws.POST("/userapi").To(s.AddUserConfig))
	ws.Route(ws.PUT("/userapi/{bk_biz_id}/{id}").To(s.UpdateUserConfig))
	ws.Route(ws.DELETE("/userapi/{bk_biz_id}/{id}").To(s.DeleteUserConfig))
	ws.Route(ws.POST("/userapi/search").To(s.GetUserConfig))
	ws.Route(ws.GET("/userapi/detail/{bk_biz_id}/{id}").To(s.UserConfigDetail))
	ws.Route(ws.POST("/usercustom/{bk_user}").To(s.AddUserCustom))
	ws.Route(ws.PUT("/usercustom/{bk_user}/{id}").To(s.UpdateUserCustomByID))
	ws.Route(ws.POST("/usercustom/user/search/{bk_user}").To(s.GetUserCustomByUser))
	ws.Route(ws.POST("/usercustom/default/search/{bk_user}").To(s.GetDefaultUserCustom))
	ws.Route(ws.GET("/healthz").To(s.Healthz))

	return ws
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
