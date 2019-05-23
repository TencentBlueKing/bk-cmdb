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
	"time"

	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cfnc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/proc_server/app/options"
	"configcenter/src/scene_server/proc_server/logics"
	ccRedis "configcenter/src/storage/dal/redis"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"
)

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

type ProcServer struct {
	*backbone.Engine
	EsbConfigChn       chan esbutil.EsbConfig
	Config             *options.Config
	EsbServ            esbserver.EsbClientInterface
	Cache              *redis.Client
	procHostInstConfig logics.ProcHostInstConfig
	ConfigMap          map[string]string
	AuthManager        *extensions.AuthManager
	Logic              *logics.Logic
}

func (s *ProcServer) newSrvComm(header http.Header) *srvComm {
	lang := util.GetLanguage(header)
	ctx, cancel := s.Engine.CCCtx.WithCancel()
	return &srvComm{
		header:        header,
		rid:           util.GetHTTPCCRequestID(header),
		ccErr:         s.CCErr.CreateDefaultCCErrorIf(lang),
		ccLang:        s.Language.CreateDefaultCCLanguageIf(lang),
		ctx:           ctx,
		ctxCancelFunc: cancel,
		user:          util.GetUser(header),
		ownerID:       util.GetOwnerID(header),
		lgc:           logics.NewLogics(s.Engine, header, s.Cache, s.EsbServ, &s.procHostInstConfig),
	}
}

func (ps *ProcServer) WebService() *restful.WebService {
	getErrFunc := func() errors.CCErrorIf {
		return ps.Engine.CCErr
	}

	// v3
	ws := new(restful.WebService)
	ws.Path("/process/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	// ws.Route(ws.POST("/{bk_supplier_account}/{bk_biz_id}").To(ps.CreateProcess))
	// ws.Route(ws.DELETE("/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.DeleteProcess))
	// ws.Route(ws.POST("/search/{bk_supplier_account}/{bk_biz_id}").To(ps.SearchProcess))
	// ws.Route(ws.PUT("/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.UpdateProcess))
	// ws.Route(ws.PUT("/{bk_supplier_account}/{bk_biz_id}").To(ps.BatchUpdateProcess))
	//
	// ws.Route(ws.GET("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.GetProcessBindModule))
	// ws.Route(ws.PUT("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}").To(ps.BindModuleProcess))
	// ws.Route(ws.DELETE("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}").To(ps.DeleteModuleProcessBind))
	//
	// ws.Route(ws.GET("/{" + common.BKOwnerIDField + "}/{" + common.BKAppIDField + "}/{" + common.BKProcessIDField + "}").To(ps.GetProcessDetailByID))

	//v2
	ws.Route(ws.POST("/openapi/GetProcessPortByApplicationID/{" + common.BKAppIDField + "}").To(ps.GetProcessPortByApplicationID))
	ws.Route(ws.POST("/openapi/GetProcessPortByIP").To(ps.GetProcessPortByIP))

	ws.Route(ws.GET("/healthz").To(ps.Healthz))
	return ws
}

func (s *ProcServer) WebService2(web *restful.WebService) {
	utility := rest.NewRestUtility(rest.Config{
		ErrorIf:  s.Engine.CCErr,
		Language: s.Engine.Language,
	})

	// service category
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/findmany/proc/service_category", Handler: s.GetServiceCategory})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/proc/service_category", Handler: s.CreateServiceCategory})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/proc/service_category", Handler: s.DeleteServiceCategory})

	// service template
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/proc/service_template", Handler: s.CreateServiceTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/findmany/proc/service_template", Handler: s.ListServiceTemplates})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/proc/service_template", Handler: s.DeleteServiceTemplate})

	// process template
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/createmany/proc/proc_template/for_service_template", Handler: s.CreateProcessTemplateBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodPut, Path: "/update/proc/proc_template/for_service_template", Handler: s.UpdateProcessTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/deletemany/proc/proc_template/for_service_template", Handler: s.DeleteProcessTemplateBatch})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/proc/proc_template/id/{processTemplateID}", Handler: s.GetProcessTemplate})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/findmany/proc/proc_template", Handler: s.ListProcessTemplate})

	// service instance
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/proc/service_instance/with_template", Handler: s.CreateServiceInstances})
	utility.AddHandler(rest.Action{Verb: http.MethodPost, Path: "/create/proc/service_instance/with_raw", Handler: s.CreateServiceInstances})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/proc/service_instance/{service_instance_id}/process", Handler: s.DeleteProcessInstanceInServiceInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/proc/service_instance", Handler: s.GetServiceInstancesInModule})
	utility.AddHandler(rest.Action{Verb: http.MethodDelete, Path: "/delete/proc/service_instance", Handler: s.DeleteServiceInstance})
	utility.AddHandler(rest.Action{Verb: http.MethodGet, Path: "/find/proc/service_instance/difference", Handler: s.FindDifferencesBetweenServiceAndProcessInstance})

	utility.AddToRestfulWebService(web)
}

func (s *ProcServer) Healthz(req *restful.Request, resp *restful.Response) {
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

	// host controller
	procCtrl := metric.HealthItem{IsHealthy: true, Name: types.CC_MODULE_PROCCONTROLLER}
	if _, err := s.Engine.CoreAPI.Healthz().HealthCheck(types.CC_MODULE_PROCCONTROLLER); err != nil {
		procCtrl.IsHealthy = false
		procCtrl.Message = err.Error()
	}
	meta.Items = append(meta.Items, procCtrl)

	for _, item := range meta.Items {
		if item.IsHealthy == false {
			meta.IsHealthy = false
			meta.Message = "proc server is unhealthy"
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

func (ps *ProcServer) OnProcessConfigUpdate(previous, current cfnc.ProcessConfig) {

	//
	esbAddr, addrOk := current.ConfigMap["esb.addr"]
	esbAppCode, appCodeOk := current.ConfigMap["esb.appCode"]
	esbAppSecret, appSecretOk := current.ConfigMap["esb.appSecret"]
	if addrOk && appCodeOk && appSecretOk {
		go func() {
			ps.EsbConfigChn <- esbutil.EsbConfig{Addrs: esbAddr, AppCode: esbAppCode, AppSecret: esbAppSecret}
		}()
	}

	cfg := ccRedis.ParseConfigFromKV("redis", current.ConfigMap)
	ps.Config = &options.Config{
		Redis: &cfg,
	}

	hostInstPrefix := "host instance"
	procHostInstConfig := &ps.procHostInstConfig
	if val, ok := current.ConfigMap[hostInstPrefix+".maxEventCount"]; ok {
		eventCount, err := util.GetIntByInterface(val)
		if nil == err {
			procHostInstConfig.MaxEventCount = eventCount
		}
	}
	if val, ok := current.ConfigMap[hostInstPrefix+".maxModuleIDCount"]; ok {
		mid_count, err := util.GetIntByInterface(val)
		if nil == err {
			procHostInstConfig.MaxRefreshModuleCount = mid_count
		}
	}
	if val, ok := current.ConfigMap[hostInstPrefix+".getModuleIDInterval"]; ok {
		get_mid_interval, err := util.GetIntByInterface(val)
		if nil == err {
			procHostInstConfig.GetModuleIDInterval = time.Duration(get_mid_interval) * time.Second
		}
	}
	ps.ConfigMap = current.ConfigMap
}
