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
	"configcenter/src/scene_server/host_server/app/options"
	"configcenter/src/scene_server/host_server/logics"
	"github.com/emicklei/go-restful"
)

type Service struct {
	*options.Config
	*backbone.Engine
	*logics.Logics
}

func (s *Service) WebService(filter restful.FilterFunction) *restful.WebService {
	ws := new(restful.WebService)
	ws.Path("/host/v3").Filter(filter).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
    restful.DefaultRequestContentType(restful.MIME_JSON)
    restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.POST("/host/batch").To(s.DeleteHostBatch))
	ws.Route(ws.GET("/hosts/{bk_supplier_account}/{bk_host_id}").To(s.GetHostInstanceProperties))
	ws.Route(ws.GET("/host/snapshot/{bk_host_id}").To(s.HostSnapInfo))
	ws.Route(ws.POST("/hosts/addhost").To(s.AddHost))
	ws.Route(ws.POST("/host/add/agent").To(s.AddHostFromAgent))
	ws.Route(ws.POST("hosts/favorites/search").To(s.GetHostFavourites))
	ws.Route(ws.POST("hosts/favorites").To(s.AddHostFavourite))
	ws.Route(ws.PUT("hosts/favorites/{id}").To(s.UpdateHostFavouriteByID))
	ws.Route(ws.DELETE("hosts/favorites/{id}").To(s.DeleteHostFavouriteByID))
	ws.Route(ws.PUT("/hosts/favorites/{id}/incr").To(s.IncrHostFavouritesCount))
	ws.Route(ws.POST("/history").To(s.AddHistory))
	ws.Route(ws.GET("/history/{start}/{limit}").To(s.GetHistorys))
	ws.Route(ws.POST("/hosts/modules/biz/mutiple").To(s.AddHostMultiAppModuleRelation))
	ws.Route(ws.POST("/hosts/modules").To(s.HostModuleRelation))
	ws.Route(ws.POST("/hosts/emptymodule").To(s.MoveHost2EmptyModule))
	ws.Route(ws.POST("/hosts/faultmodule").To(s.MoveHost2FaultModule))
	ws.Route(ws.POST("/hosts/resource").To(s.MoveHostToResourcePool))
	ws.Route(ws.POST("/hosts/assgin").To(s.AssignHostToApp))
	ws.Route(ws.POST("/host/add/module").To(s.AssignHostToAppModule))
	ws.Route(ws.POST("/usercustom").To(s.SaveUserCustom))
	ws.Route(ws.POST("/usercustom/user/search").To(s.GetUserCustom))
	ws.Route(ws.POST("/usercustom/default/search").To(s.GetDefaultCustom))
	ws.Route(ws.GET("getAgentStatus/{appid}").To(s.GetAgentStatus))
	ws.Route(ws.PUT("/openapi/host/{" + common.BKAppIDField + "}").To(s.UpdateHost))
	ws.Route(ws.PUT("/host/updateHostByAppID/{appid}").To(s.UpdateHostByAppID))
	ws.Route(ws.POST("/gethostlistbyip").To(s.HostSearchByIP))
	ws.Route(ws.POST("/gethostlistbyconds").To(s.HostSearchByConds))
	ws.Route(ws.POST("/getmodulehostlist").To(s.HostSearchByModuleID))
	ws.Route(ws.POST("/getsethostlist").To(s.HostSearchBySetID))
	ws.Route(ws.POST("/getapphostlist").To(s.HostSearchByAppID))
	ws.Route(ws.POST("/gethostsbyproperty").To(s.HostSearchByProperty))
	ws.Route(ws.POST("/getIPAndProxyByCompany").To(s.GetIPAndProxyByCompany))
	ws.Route(ws.PUT("/openapi/updatecustomproperty").To(s.UpdateCustomProperty))
	ws.Route(ws.PUT("openapi/host/clonehostproperty").To(s.CloneHostProperty))
	ws.Route(ws.POST("/openapi/host/getHostAppByCompanyId").To(s.GetHostAppByCompanyId))
	ws.Route(ws.DELETE("/openapi/host/delhostinapp").To(s.DelHostInApp))
	ws.Route(ws.POST("/openapi/host/getGitServerIp").To(s.GetGitServerIp))
	ws.Route(ws.GET("/plat").To(s.GetPlat))
	ws.Route(ws.POST("/plat").To(s.CreatePlat))
	ws.Route(ws.DELETE("/plat/{bk_cloud_id}").To(s.DelPlat))
	ws.Route(ws.POST("/search").To(s.SearchHost))
	ws.Route(ws.POST("/search/asstdetail").To(s.SearchHostWithAsstDetail))
	ws.Route(ws.PUT("/host/batch").To(s.UpdateHostBatch))
	ws.Route(ws.POST("/userapi").To(s.AddUserCustomQuery))
	ws.Route(ws.PUT("/userapi/{bk_biz_id}/{id}").To(s.UpdateUserCustomQuery))
	ws.Route(ws.DELETE("/userapi/{bk_biz_id}/{id}").To(s.DeleteUserCustomQuery))
	ws.Route(ws.POST("/userapi/search/{bk_biz_id}").To(s.GetUserCustomQuery))
	ws.Route(ws.GET("/userapi/detail/{bk_biz_id}/{id}").To(s.GetUserCustomQueryDetail))
	ws.Route(ws.GET("/userapi/data/{bk_biz_id}/{id}/{start}/{limit}").To(s.GetUserCustomQueryResult))

	return ws
}
