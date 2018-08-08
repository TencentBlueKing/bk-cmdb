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
	"net/http"

	"github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cfnc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/errors"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
)

type ProcServer struct {
	*backbone.Engine
}

func (ps *ProcServer) WebService() http.Handler {
	getErrFun := func() errors.CCErrorIf {
		return ps.Engine.CCErr
	}

	container := restful.NewContainer()
	// v3
	ws := new(restful.WebService)
	ws.Path("/process/{version}").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON)
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

	ws.Route(ws.GET("/{" + common.BKOwnerIDField + "}/{" + common.BKAppIDField + "}/{" + common.BKProcIDField + "}").To(ps.GetProcessDetailByID))

	ws.Route(ws.POST("/openapi/GetProcessPortByApplicationID/{" + common.BKAppIDField + "}").To(ps.GetProcessPortByApplicationID))
	ws.Route(ws.POST("/openapi/GetProcessPortByIP").To(ps.GetProcessPortByIP))

	ws.Route(ws.POST("/operate/{namespace}/process").To(ps.OperateProcessInstance))
	ws.Route(ws.POST("/operate/{namespace}/process/taskresult").To(ps.QueryProcessOperateResult))

	ws.Route(ws.POST("/conftemp").To(ps.CreateConfigTemp))
	ws.Route(ws.PUT("/conftemp").To(ps.UpdateConfigTemp))
	ws.Route(ws.DELETE("/conftemp").To(ps.DeleteConfigTemp))
	ws.Route(ws.POST("/conftemp/search").To(ps.QueryConfigTemp))
	ws.Route(ws.GET("/healthz").To(ps.Healthz))

	container.Add(ws)

	return container
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

func (ps *ProcServer) OnProcessConfigUpdate(previous, current cfnc.ProcessConfig) {
	//
}
