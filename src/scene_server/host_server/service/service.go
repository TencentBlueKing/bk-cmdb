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

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/host_server/app/options"
	"configcenter/src/scene_server/host_server/logics"
)

type Service struct {
	*options.Config
	*backbone.Engine
	*logics.Logics
	disc discovery.DiscoveryInterface
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/host/{version}").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON)
	//restful.DefaultRequestContentType(restful.MIME_JSON)
	//restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.DELETE("/hosts/batch").To(s.DeleteHostBatch))
	ws.Route(ws.GET("/hosts/{bk_supplier_account}/{bk_host_id}").To(s.GetHostInstanceProperties))
	ws.Route(ws.GET("/hosts/snapshot/{bk_host_id}").To(s.HostSnapInfo))
	ws.Route(ws.POST("/hosts/add").To(s.AddHost))
	ws.Route(ws.POST("/host/add/agent").To(s.AddHostFromAgent))
	ws.Route(ws.POST("/hosts/sync/new/host").To(s.NewHostSyncAppTopo))
	ws.Route(ws.POST("hosts/favorites/search").To(s.GetHostFavourites))
	ws.Route(ws.POST("hosts/favorites").To(s.AddHostFavourite))
	ws.Route(ws.PUT("hosts/favorites/{id}").To(s.UpdateHostFavouriteByID))
	ws.Route(ws.DELETE("hosts/favorites/{id}").To(s.DeleteHostFavouriteByID))
	ws.Route(ws.PUT("/hosts/favorites/{id}/incr").To(s.IncrHostFavouritesCount))
	ws.Route(ws.POST("/hosts/history").To(s.AddHistory))
	ws.Route(ws.GET("/hosts/history/{start}/{limit}").To(s.GetHistorys))
	ws.Route(ws.POST("/hosts/modules/biz/mutiple").To(s.AddHostMultiAppModuleRelation))
	ws.Route(ws.POST("/hosts/modules").To(s.HostModuleRelation))
	ws.Route(ws.POST("/hosts/modules/idle").To(s.MoveHost2EmptyModule))
	ws.Route(ws.POST("/hosts/modules/fault").To(s.MoveHost2FaultModule))
	ws.Route(ws.POST("/hosts/modules/resource").To(s.MoveHostToResourcePool))
	ws.Route(ws.POST("/hosts/modules/resource/idle").To(s.AssignHostToApp))
	ws.Route(ws.POST("/host/add/module").To(s.AssignHostToAppModule))
	ws.Route(ws.POST("/usercustom").To(s.SaveUserCustom))
	ws.Route(ws.POST("/usercustom/user/search").To(s.GetUserCustom))
	ws.Route(ws.POST("/usercustom/default/search").To(s.GetDefaultCustom))
	ws.Route(ws.POST("/hosts/search").To(s.SearchHost))
	ws.Route(ws.POST("/hosts/search/asstdetail").To(s.SearchHostWithAsstDetail))
	ws.Route(ws.PUT("/hosts/batch").To(s.UpdateHostBatch))
	ws.Route(ws.PUT("/hosts/property/clone").To(s.CloneHostProperty))
	ws.Route(ws.POST("/hosts/modules/idle/set").To(s.MoveSetHost2IdleModule))

	ws.Route(ws.POST("/userapi").To(s.AddUserCustomQuery))
	ws.Route(ws.PUT("/userapi/{bk_biz_id}/{id}").To(s.UpdateUserCustomQuery))
	ws.Route(ws.DELETE("/userapi/{bk_biz_id}/{id}").To(s.DeleteUserCustomQuery))
	ws.Route(ws.POST("/userapi/search/{bk_biz_id}").To(s.GetUserCustomQuery))
	ws.Route(ws.GET("/userapi/detail/{bk_biz_id}/{id}").To(s.GetUserCustomQueryDetail))
	ws.Route(ws.GET("/userapi/data/{bk_biz_id}/{id}/{start}/{limit}").To(s.GetUserCustomQueryResult))

	ws.Route(ws.GET("/host/getHostListByAppidAndField/{" + common.BKAppIDField + "}/{field}").To(s.getHostListByAppidAndField))
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
	ws.Route(ws.POST("/openapi/host/getHostAppByCompanyId").To(s.GetHostAppByCompanyId))
	ws.Route(ws.DELETE("/openapi/host/delhostinapp").To(s.DelHostInApp))
	ws.Route(ws.POST("/openapi/host/getGitServerIp").To(s.GetGitServerIp))
	ws.Route(ws.GET("/plat").To(s.GetPlat))
	ws.Route(ws.POST("/plat").To(s.CreatePlat))
	ws.Route(ws.DELETE("/plat/{bk_cloud_id}").To(s.DelPlat))
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

	// object controller
	objCtr := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_OBJECTCONTROLLER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_OBJECTCONTROLLER); err != nil {
		objCtr.IsHealthy = false
		objCtr.Message = err.Error()
	}
	meta.Items = append(meta.Items, objCtr)

	// audit controller
	auditCtrl := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_AUDITCONTROLLER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_AUDITCONTROLLER); err != nil {
		auditCtrl.IsHealthy = false
		auditCtrl.Message = err.Error()
	}
	meta.Items = append(meta.Items, auditCtrl)

	// host controller
	hostCtrl := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_HOSTCONTROLLER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_HOSTCONTROLLER); err != nil {
		hostCtrl.IsHealthy = false
		hostCtrl.Message = err.Error()
	}
	meta.Items = append(meta.Items, hostCtrl)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "host server is unhealthy"
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
