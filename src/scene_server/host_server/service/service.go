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
	api.Path("/host/v3").Filter(s.Engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	api.Route(api.DELETE("/hosts/batch").To(s.DeleteHostBatchFromResourcePool))
	api.Route(api.GET("/hosts/{bk_supplier_account}/{bk_host_id}").To(s.GetHostInstanceProperties))
	api.Route(api.GET("/hosts/snapshot/{bk_host_id}").To(s.HostSnapInfo))
	api.Route(api.POST("/hosts/add").To(s.AddHost))
	// api.Route(api.POST("/host/add/agent").To(s.AddHostFromAgent))
	api.Route(api.POST("/hosts/sync/new/host").To(s.NewHostSyncAppTopo))
	api.Route(api.PUT("/updatemany/hosts/cloudarea_field").To(s.UpdateHostCloudAreaField))

	// host favorites
	api.Route(api.POST("/hosts/favorites/search").To(s.ListHostFavourites))
	api.Route(api.POST("/hosts/favorites").To(s.AddHostFavourite))
	api.Route(api.PUT("/hosts/favorites/{id}").To(s.UpdateHostFavouriteByID))
	api.Route(api.DELETE("/hosts/favorites/{id}").To(s.DeleteHostFavouriteByID))
	api.Route(api.PUT("/hosts/favorites/{id}/incr").To(s.IncrHostFavouritesCount))

	api.Route(api.POST("/hosts/modules").To(s.TransferHostModule))
	api.Route(api.POST("/hosts/modules/idle").To(s.MoveHost2IdleModule))
	api.Route(api.POST("/hosts/modules/fault").To(s.MoveHost2FaultModule))
	api.Route(api.POST("/hosts/modules/recycle").To(s.MoveHost2RecycleModule))
	api.Route(api.POST("/hosts/modules/resource").To(s.MoveHostToResourcePool))
	api.Route(api.POST("/hosts/modules/resource/idle").To(s.AssignHostToApp))
	api.Route(api.POST("/host/add/module").To(s.AssignHostToAppModule))
	api.Route(api.POST("/host/transfer_with_auto_clear_service_instance/bk_biz_id/{bk_biz_id}/").To(s.TransferHostWithAutoClearServiceInstance))
	api.Route(api.POST("/host/transfer_with_auto_clear_service_instance/bk_biz_id/{bk_biz_id}/preview/").To(s.TransferHostWithAutoClearServiceInstancePreview))
	api.Route(api.POST("/usercustom").To(s.SaveUserCustom))
	api.Route(api.POST("/usercustom/user/search").To(s.GetUserCustom))
	api.Route(api.POST("/usercustom/default/search").To(s.GetDefaultCustom))
	api.Route(api.POST("/hosts/search").To(s.SearchHost))
	api.Route(api.POST("/hosts/search/asstdetail").To(s.SearchHostWithAsstDetail))
	api.Route(api.PUT("/hosts/batch").To(s.UpdateHostBatch))
	api.Route(api.PUT("/hosts/property/batch").To(s.UpdateHostPropertyBatch))
	api.Route(api.PUT("/hosts/property/clone").To(s.CloneHostProperty))
	api.Route(api.POST("/hosts/modules/idle/set").To(s.MoveSetHost2IdleModule))
	// get host module relation in app
	api.Route(api.POST("/hosts/modules/read").To(s.GetHostModuleRelation))
	api.Route(api.POST("/host/topo/relation/read").To(s.GetAppHostTopoRelation))
	// transfer host to other business
	api.Route(api.POST("/hosts/modules/across/biz").To(s.TransferHostAcrossBusiness))
	//  delete host from business, used for framework
	api.Route(api.DELETE("/hosts/module/biz/delete").To(s.DeleteHostFromBusiness))

	// next generation host search api
	api.Route(api.POST("/hosts/list_hosts_without_app").To(s.ListHostsWithNoBiz))
	api.Route(api.POST("/hosts/app/{appid}/list_hosts").To(s.ListBizHosts))
	api.Route(api.POST("/hosts/app/{bk_biz_id}/list_hosts_topo").To(s.ListBizHostsTopo))

	api.Route(api.POST("/userapi").To(s.AddUserCustomQuery))
	api.Route(api.PUT("/userapi/{bk_biz_id}/{id}").To(s.UpdateUserCustomQuery))
	api.Route(api.DELETE("/userapi/{bk_biz_id}/{id}").To(s.DeleteUserCustomQuery))
	api.Route(api.POST("/userapi/search/{bk_biz_id}").To(s.GetUserCustomQuery))
	api.Route(api.GET("/userapi/detail/{bk_biz_id}/{id}").To(s.GetUserCustomQueryDetail))
	api.Route(api.GET("/userapi/data/{bk_biz_id}/{id}/{start}/{limit}").To(s.GetUserCustomQueryResult))

	api.Route(api.POST("/host/lock").To(s.LockHost))
	api.Route(api.DELETE("/host/lock").To(s.UnlockHost))
	api.Route(api.POST("/host/lock/search").To(s.QueryHostLock))
	api.Route(api.POST("/host/count_by_topo_node/bk_biz_id/{bk_biz_id}").To(s.CountTopoNodeHosts))

	api.Route(api.POST("/findmany/modulehost").To(s.FindModuleHost))

	// cloud sync
	api.Route(api.POST("/hosts/cloud/add").To(s.AddCloudTask))
	api.Route(api.DELETE("/hosts/cloud/delete/{taskID}").To(s.DeleteCloudTask))
	api.Route(api.POST("/hosts/cloud/search").To(s.SearchCloudTask))
	api.Route(api.PUT("/hosts/cloud/update").To(s.UpdateCloudTask))
	api.Route(api.POST("/hosts/cloud/startSync").To(s.StartCloudSync))
	api.Route(api.POST("/hosts/cloud/resourceConfirm").To(s.CreateResourceConfirm))
	api.Route(api.POST("/hosts/cloud/searchConfirm").To(s.SearchConfirm))
	api.Route(api.POST("/hosts/cloud/confirmHistory/add").To(s.AddConfirmHistory))
	api.Route(api.POST("/hosts/cloud/confirmHistory/search").To(s.SearchConfirmHistory))
	api.Route(api.POST("/hosts/cloud/accountSearch").To(s.SearchAccount))
	api.Route(api.POST("/hosts/cloud/syncHistory").To(s.SearchCloudSyncHistory))

	api.Route(api.POST("/findmany/cloudarea").To(s.FindManyCloudArea))
	api.Route(api.POST("/create/cloudarea").To(s.CreatePlat))
	api.Route(api.PUT("/update/cloudarea/{bk_cloud_id}").To(s.UpdatePlat))
	api.Route(api.DELETE("/delete/cloudarea/{bk_cloud_id}").To(s.DelPlat))

	// first install use api
	api.Route(api.POST("/host/install/bk").To(s.BKSystemInstall))

	// 主机属性自动应用
	api.Route(api.POST("/create/host_apply_rule/bk_biz_id/{bk_biz_id}").To(s.CreateHostApplyRule))
	api.Route(api.PUT("/update/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}").To(s.UpdateHostApplyRule))
	api.Route(api.DELETE("/deletemany/host_apply_rule/bk_biz_id/{bk_biz_id}").To(s.DeleteHostApplyRule))
	api.Route(api.GET("/find/host_apply_rule/{host_apply_rule_id}/bk_biz_id/{bk_biz_id}/").To(s.GetHostApplyRule))
	api.Route(api.POST("/findmany/host_apply_rule/bk_biz_id/{bk_biz_id}").To(s.ListHostApplyRule))
	api.Route(api.POST("/createmany/host_apply_rule/bk_biz_id/{bk_biz_id}/batch_create_or_update").To(s.BatchCreateOrUpdateHostApplyRule))
	api.Route(api.POST("/createmany/host_apply_plan/bk_biz_id/{bk_biz_id}/preview").To(s.GenerateApplyPlan))
	api.Route(api.POST("/updatemany/host_apply_plan/bk_biz_id/{bk_biz_id}/run").To(s.RunHostApplyRule))
	api.Route(api.POST("/findmany/host_apply_rule/bk_biz_id/{bk_biz_id}/host_related_rules").To(s.ListHostRelatedApplyRule))

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
	_ = resp.WriteEntity(answer)
}

func (s *Service) InitBackground() {
	header := make(http.Header, 0)
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKProcInstanceOpUser)
	}
	s.CacheDB.FlushDb()

	srvData := s.newSrvComm(header)
	go srvData.lgc.TimerTriggerCheckStatus(srvData.ctx)
}
