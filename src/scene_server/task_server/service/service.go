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
	"fmt"
	"net/http"
	"time"

	"configcenter/pkg/tenant"
	tenanttype "configcenter/pkg/tenant/types"
	hosttype "configcenter/pkg/types/host"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/scene_server/task_server/app/options"
	"configcenter/src/scene_server/task_server/logics"
	"configcenter/src/storage/dal/mongo/sharding"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// Service TODO
type Service struct {
	*options.Config
	*backbone.Engine
	disc    discovery.DiscoveryInterface
	CacheDB redis.Client
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
	commonAPI.Route(commonAPI.POST("/refresh/tenants").To(s.RefreshTenants))
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
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/field_template/task",
		Handler: s.CreateFieldTemplateTask})

	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/findmany/list/{name}", Handler: s.ListTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/findmany/list/latest/{name}",
		Handler: s.ListLatestTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/findone/detail/{task_id}",
		Handler: s.DetailTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/deletemany", Handler: s.DeleteTask})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/latest/sync_status",
		Handler: s.ListLatestSyncStatus})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/findmany/sync_status_history",
		Handler: s.ListSyncStatusHistory})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/task/find/field_template/task_sync_result",
		Handler: s.ListFieldTmplTaskSyncResult})

	utility.AddToRestfulWebService(web)

}

// Healthz TODO
func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}

	// mongodb status
	meta.Items = append(meta.Items, mongodb.Healthz()...)

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

// RefreshTenants refresh tenants
func (s *Service) RefreshTenants(req *restful.Request, resp *restful.Response) {
	rHeader := req.Request.Header
	defErr := s.CCErr.CreateDefaultCCErrorIf(httpheader.GetLanguage(rHeader))
	kit := rest.NewKitFromHeader(rHeader, s.CCErr)

	tenants := make([]tenanttype.Tenant, 0)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameTenant).Find(mapstr.MapStr{}).All(kit.Ctx,
		&tenants)
	if err != nil {
		blog.Errorf("find all tenants failed, err: %v, rid: %s", err, kit.Rid)
		result := &metadata.RespError{
			Msg: defErr.Errorf(common.CCErrObjectDBOpErrno, fmt.Errorf("find all tenants failed")),
		}
		resp.WriteError(http.StatusInternalServerError, result)
	}

	tenant.SetTenant(tenants)

	resp.WriteEntity(metadata.NewSuccessResp(tenants))
}

// SyncDefaultAreaHostTask background task
func (s *Service) SyncDefaultAreaHostTask(engine *backbone.Engine) error {

	kit := rest.NewKit()
	cond := map[string]any{common.MongoMetaID: common.ShardingDBConfID}
	conf := new(sharding.ShardingDBConf)
	err := mongodb.Shard(kit.SysShardOpts()).Table(common.BKTableNameSystem).Find(cond).One(kit.Ctx, &conf)
	if err != nil {
		blog.Errorf("get sharding db config failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if len(conf.SlaveDB) == 0 {
		blog.Info("no slave db, skip default area host compare background task")
		return nil
	}

	go func() {
		for {
			// only master can run it
			if !engine.ServiceManageInterface.IsMaster() {
				blog.V(4).Infof("it is not master, skip sync default area hosts, rid: %s", kit.Rid)
				time.Sleep(time.Minute)
				continue
			}

			time.Sleep(20 * time.Minute)
			if err := s.syncDefaultAreaHosts(kit); err != nil {
				blog.Errorf("sync default area hosts failed, err: %v, rid: %s", err, kit.Rid)
				continue
			}
		}
	}()

	return nil
}

func (s *Service) syncDefaultAreaHosts(kit *rest.Kit) error {

	allTenants := tenant.GetAllTenants()
	for _, tenant := range allTenants {
		newTenantKit := kit.NewKit().WithTenant(tenant.TenantID)
		tenantCond := mapstr.MapStr{
			common.TenantID: newTenantKit.TenantID,
		}

		lastHostID := int64(0)
		for {
			host := make([]metadata.DefaultAreaHost, 0)
			tenantCond[common.BKHostIDField] = mapstr.MapStr{
				common.BKDBGT: lastHostID,
			}
			err := mongodb.Shard(newTenantKit.SysShardOpts()).Table(common.BKTableNameDefaultAreaHost).Find(tenantCond).
				Limit(common.BKMaxLimitSize).Sort(common.BKHostIDField).All(kit.Ctx, &host)
			if err != nil {
				blog.Errorf("failed to get host, err: %v, rid: %s", err, newTenantKit.Rid)
				return err
			}

			if len(host) == 0 {
				break
			}
			lastHostID = host[len(host)-1].HostID

			hostIDs := make([]int64, 0)
			hostIDMap := make(map[int64]struct{})
			for _, h := range host {
				hostIDs = append(hostIDs, h.HostID)
				hostIDMap[h.HostID] = struct{}{}
			}
			cond := mapstr.MapStr{
				common.BKHostIDField: mapstr.MapStr{
					common.BKDBIN: hostIDs,
				},
			}

			existHosts := make([]hosttype.HostBaseInfo, 0)
			err = mongodb.Shard(newTenantKit.ShardOpts()).Table(common.BKTableNameBaseHost).Find(cond).Fields(
				common.BKHostIDField).All(kit.Ctx, &existHosts)
			if err != nil {
				blog.Errorf("failed to get exist hosts, err: %v, rid: %s", err, newTenantKit.Rid)
				return err
			}
			for _, h := range existHosts {
				delete(hostIDMap, h.HostID)
			}

			redundantHost := make([]int64, 0)
			for hostID := range hostIDMap {
				redundantHost = append(redundantHost, hostID)
			}

			if len(redundantHost) == 0 {
				continue
			}
			cond = mapstr.MapStr{
				common.BKHostIDField: mapstr.MapStr{
					common.BKDBIN: redundantHost,
				},
			}
			err = mongodb.Shard(newTenantKit.SysShardOpts()).Table(common.BKTableNameDefaultAreaHost).Delete(
				newTenantKit.Ctx, cond)
			if err != nil {
				blog.Errorf("failed to delete redundant host, err: %v, rid: %s", err, newTenantKit.Rid)
				return err
			}
		}

	}
	return nil
}
