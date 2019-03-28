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

	"github.com/emicklei/go-restful"
	redis "gopkg.in/redis.v5"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cfnc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/errors"
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
	ws.Path("/process/{version}").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.POST("/{bk_supplier_account}/{bk_biz_id}").To(ps.CreateProcess))
	ws.Route(ws.DELETE("/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.DeleteProcess))
	ws.Route(ws.POST("/search/{bk_supplier_account}/{bk_biz_id}").To(ps.SearchProcess))
	ws.Route(ws.PUT("/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.UpdateProcess))
	ws.Route(ws.PUT("/{bk_supplier_account}/{bk_biz_id}").To(ps.BatchUpdateProcess))

	ws.Route(ws.GET("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.GetProcessBindModule))
	ws.Route(ws.PUT("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}").To(ps.BindModuleProcess))
	ws.Route(ws.DELETE("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}").To(ps.DeleteModuleProcessBind))

	ws.Route(ws.GET("/{" + common.BKOwnerIDField + "}/{" + common.BKAppIDField + "}/{" + common.BKProcessIDField + "}").To(ps.GetProcessDetailByID))

	ws.Route(ws.POST("/operate/process").To(ps.OperateProcessInstance))
	ws.Route(ws.GET("/operate/process/taskresult/{taskID}").To(ps.QueryProcessOperateResult))

	ws.Route(ws.POST("/template/{bk_supplier_account}/{bk_biz_id}").To(ps.CreateTemplate))
	ws.Route(ws.PUT("/template/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.UpdateTemplate))
	ws.Route(ws.DELETE("/template/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.DeleteTemplate))
	ws.Route(ws.POST("/template/search/{bk_supplier_account}/{bk_biz_id}").To(ps.SearchTemplate))
	ws.Route(ws.POST("/template/version/search/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.SearchTemplateVersion))
	ws.Route(ws.POST("/template/version/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.CreateTemplateVersion))
	ws.Route(ws.PUT("/template/vesrion/{bk_supplier_account}/{bk_biz_id}/{template_id}/{version_id}").To(ps.UpdateTemplateVersion))
	ws.Route(ws.GET("/template/proc/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.GetProcBindTemplate))
	ws.Route(ws.PUT("/template/proc/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{template_id}").To(ps.BindProc2Template))
	ws.Route(ws.DELETE("/template/proc/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{template_id}").To(ps.DeleteProc2Template))
	ws.Route(ws.POST("/template/preview/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.PreviewCfg))
	ws.Route(ws.POST("/template/create/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.CreateCfg))
	ws.Route(ws.POST("/template/push/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.PushCfg))
	ws.Route(ws.POST("/template/getremote/{bk_supplier_account}/{bk_biz_id}/{template_id}").To(ps.GetRemoteCfg))
	ws.Route(ws.GET("/template/group/{bk_supplier_account}/{bk_biz_id}").To(ps.GetTemplateGroup))

	//v2
	ws.Route(ws.POST("/openapi/GetProcessPortByApplicationID/{" + common.BKAppIDField + "}").To(ps.GetProcessPortByApplicationID))
	ws.Route(ws.POST("/openapi/GetProcessPortByIP").To(ps.GetProcessPortByIP))

	ws.Route(ws.GET("/healthz").To(ps.Healthz))
	ws.Route(ws.POST("/process/refresh/hostinstnum").To(ps.RefreshProcHostInstByEvent))
	return ws
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

}
