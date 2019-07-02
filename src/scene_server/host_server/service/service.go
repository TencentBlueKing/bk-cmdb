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
	"context"
	"net/http"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/host_server/app/options"
	"configcenter/src/scene_server/host_server/logics"

	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"
)

type Service struct {
	*options.Config
	*backbone.Engine
	disc        discovery.DiscoveryInterface
	CacheDB     *redis.Client
	AuthManager *extensions.AuthManager
}

type srvComm struct {
	header        http.Header
	rid           string
	ccErr         errors.DefaultCCErrorIf
	ccLang        language.DefaultCCLanguageIf
	ctx           context.Context
	ctxCancelFunc context.CancelFunc
	user          string
	ownerID       string
	lgc           *logics.Logics
}

func (s *Service) newSrvComm(header http.Header) *srvComm {
	rid := util.GetHTTPCCRequestID(header)
	lang := util.GetLanguage(header)
	user := util.GetUser(header)
	ctx, cancel := s.Engine.CCCtx.WithCancel()
	ctx = context.WithValue(ctx, common.ContextRequestIDField, rid)
	ctx = context.WithValue(ctx, common.ContextRequestUserField, user)

	return &srvComm{
		header:        header,
		rid:           rid,
		ccErr:         s.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:        s.Language.CreateDefaultCCLanguageIf(lang),
		ctx:           ctx,
		ctxCancelFunc: cancel,
		user:          util.GetUser(header),
		ownerID:       util.GetOwnerID(header),
		lgc:           logics.NewLogics(s.Engine, header, s.CacheDB, s.AuthManager),
	}
}

func (s *Service) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}
	api.Path("/host/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	api.Route(api.DELETE("/hosts/batch").To(s.DeleteHostBatchFromResourcePool))
	api.Route(api.GET("/hosts/{bk_supplier_account}/{bk_host_id}").To(s.GetHostInstanceProperties))
	api.Route(api.GET("/hosts/snapshot/{bk_host_id}").To(s.HostSnapInfo))
	api.Route(api.POST("/hosts/add").To(s.AddHost))
	// api.Route(api.POST("/host/add/agent").To(s.AddHostFromAgent))
	api.Route(api.POST("/hosts/sync/new/host").To(s.NewHostSyncAppTopo))
	api.Route(api.POST("hosts/favorites/search").To(s.GetHostFavourites))
	api.Route(api.POST("hosts/favorites").To(s.AddHostFavourite))
	api.Route(api.PUT("hosts/favorites/{id}").To(s.UpdateHostFavouriteByID))
	api.Route(api.DELETE("hosts/favorites/{id}").To(s.DeleteHostFavouriteByID))
	api.Route(api.PUT("/hosts/favorites/{id}/incr").To(s.IncrHostFavouritesCount))
	api.Route(api.POST("/hosts/modules").To(s.HostModuleRelation))
	api.Route(api.POST("/hosts/modules/idle").To(s.MoveHost2EmptyModule))
	api.Route(api.POST("/hosts/modules/fault").To(s.MoveHost2FaultModule))
	api.Route(api.POST("/hosts/modules/resource").To(s.MoveHostToResourcePool))
	api.Route(api.POST("/hosts/modules/resource/idle").To(s.AssignHostToApp))
	api.Route(api.POST("/host/add/module").To(s.AssignHostToAppModule))
	api.Route(api.POST("/usercustom").To(s.SaveUserCustom))
	api.Route(api.POST("/usercustom/user/search").To(s.GetUserCustom))
	api.Route(api.POST("/usercustom/default/search").To(s.GetDefaultCustom))
	api.Route(api.POST("/hosts/search").To(s.SearchHost))
	api.Route(api.POST("/hosts/search/asstdetail").To(s.SearchHostWithAsstDetail))
	api.Route(api.PUT("/hosts/batch").To(s.UpdateHostBatch))
	api.Route(api.PUT("/hosts/property/clone").To(s.CloneHostProperty))
	api.Route(api.POST("/hosts/modules/idle/set").To(s.MoveSetHost2IdleModule))
	// get host module relation in app
	api.Route(api.POST("/hosts/modules/read").To(s.GetHostModuleRelation))
	// transfer host to other business
	api.Route(api.POST("/hosts/modules/across/biz").To(s.TransferHostAcrossBusiness))
	//  delete host from business, used for framework
	api.Route(api.DELETE("/hosts/module/biz/delete").To(s.DeleteHostFromBusiness))

	api.Route(api.POST("/userapi").To(s.AddUserCustomQuery))
	api.Route(api.PUT("/userapi/{bk_biz_id}/{id}").To(s.UpdateUserCustomQuery))
	api.Route(api.DELETE("/userapi/{bk_biz_id}/{id}").To(s.DeleteUserCustomQuery))
	api.Route(api.POST("/userapi/search/{bk_biz_id}").To(s.GetUserCustomQuery))
	api.Route(api.GET("/userapi/detail/{bk_biz_id}/{id}").To(s.GetUserCustomQueryDetail))
	api.Route(api.GET("/userapi/data/{bk_biz_id}/{id}/{start}/{limit}").To(s.GetUserCustomQueryResult))

	api.Route(api.POST("/host/lock").To(s.LockHost))
	api.Route(api.DELETE("/host/lock").To(s.UnlockHost))
	api.Route(api.POST("/host/lock/search").To(s.QueryHostLock))

	api.Route(api.GET("/host/getHostListByAppidAndField/{" + common.BKAppIDField + "}/{field}").To(s.getHostListByAppidAndField))
	api.Route(api.PUT("/openapi/host/{" + common.BKAppIDField + "}").To(s.UpdateHost))
	api.Route(api.PUT("/host/updateHostByAppID/{appid}").To(s.UpdateHostByAppID))
	api.Route(api.POST("/gethostlistbyip").To(s.HostSearchByIP))
	api.Route(api.POST("/gethostlistbyconds").To(s.HostSearchByConds))
	api.Route(api.POST("/getmodulehostlist").To(s.HostSearchByModuleID))
	api.Route(api.POST("/getsethostlist").To(s.HostSearchBySetID))
	api.Route(api.POST("/getapphostlist").To(s.HostSearchByAppID))
	api.Route(api.POST("/gethostsbyproperty").To(s.HostSearchByProperty))
	api.Route(api.POST("/getIPAndProxyByCompany").To(s.GetIPAndProxyByCompany))
	api.Route(api.PUT("/openapi/updatecustomproperty").To(s.UpdateCustomProperty))
	api.Route(api.POST("/openapi/host/getHostAppByCompanyId").To(s.GetHostAppByCompanyId))
	api.Route(api.DELETE("/openapi/host/delhostinapp").To(s.DelHostInApp))
	api.Route(api.POST("/openapi/host/getGitServerIp").To(s.GetGitServerIp))
	api.Route(api.GET("/plat").To(s.GetPlat))
	api.Route(api.POST("/plat").To(s.CreatePlat))
	api.Route(api.DELETE("/plat/{bk_cloud_id}").To(s.DelPlat))

	api.Route(api.POST("/findmany/modulehost").To(s.FindModuleHost))

	// cloud sync
	api.Route(api.POST("/hosts/cloud/add").To(s.AddCloudTask))
	api.Route(api.DELETE("/hosts/cloud/delete/{taskID}").To(s.DeleteCloudTask))
	api.Route(api.POST("/hosts/cloud/search").To(s.SearchCloudTask))
	api.Route(api.PUT("/hosts/cloud/update").To(s.UpdateCloudTask))
	api.Route(api.POST("/hosts/cloud/startSync").To(s.StartCloudSync))
	api.Route(api.POST("/hosts/cloud/resourceConfirm").To(s.ResourceConfirm))
	api.Route(api.POST("/hosts/cloud/searchConfirm").To(s.SearchConfirm))
	api.Route(api.POST("/hosts/cloud/confirmHistory/add").To(s.AddConfirmHistory))
	api.Route(api.POST("/hosts/cloud/confirmHistory/search").To(s.SearchConfirmHistory))
	api.Route(api.POST("/hosts/cloud/accountSearch").To(s.SearchAccount))
	api.Route(api.POST("/hosts/cloud/syncHistory").To(s.CloudSyncHistory))

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

	// object controller
	objCtr := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_OBJECTCONTROLLER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_OBJECTCONTROLLER); err != nil {
		objCtr.IsHealthy = false
		objCtr.Message = err.Error()
	}
	meta.Items = append(meta.Items, objCtr)

	// host controller
	hostCtrl := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_HOSTCONTROLLER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_HOSTCONTROLLER); err != nil {
		hostCtrl.IsHealthy = false
		hostCtrl.Message = err.Error()
	}
	meta.Items = append(meta.Items, hostCtrl)

	// coreservice
	coreSrv := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_CORESERVICE}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_CORESERVICE); err != nil {
		coreSrv.IsHealthy = false
		coreSrv.Message = err.Error()
	}
	meta.Items = append(meta.Items, coreSrv)

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

func (s *Service) InitBackground() {
	header := make(http.Header, 0)
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}

	srvData := s.newSrvComm(header)
	go srvData.lgc.TimerTriggerCheckStatus(srvData.ctx)
}
