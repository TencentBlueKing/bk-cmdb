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
    
    "configcenter/src/common/backbone"
    cfnc "configcenter/src/common/backbone/configcenter"
    "configcenter/src/common"

    "github.com/emicklei/go-restful"
)

type ProcServer struct {
    *backbone.Engine
}

func (ps *ProcServer) WebService(filter restful.FilterFunction) http.Handler {
    container := new(restful.Container)
    // v3
    v3WS := new(restful.WebService)
    v3WS.Path("/process/v3").Filter(filter).Produces(restful.MIME_JSON)
    
    v3WS.Route(v3WS.POST("/{bk_supplier_account}/{bk_biz_id}").To(ps.CreateProcess))
    v3WS.Route(v3WS.DELETE("/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.DeleteProcess))
    v3WS.Route(v3WS.POST("/search/{bk_supplier_account}/{bk_biz_id}").To(ps.SearchProcess))
    v3WS.Route(v3WS.PUT("/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.UpdateProcess))
    v3WS.Route(v3WS.PUT("/{bk_supplier_account}/{bk_biz_id}").To(ps.BatchUpdateProcess))
    
    v3WS.Route(v3WS.GET("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}").To(ps.GetProcessBindModule))
    v3WS.Route(v3WS.PUT("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}").To(ps.BindModuleProcess))
    v3WS.Route(v3WS.DELETE("/module/{bk_supplier_account}/{bk_biz_id}/{bk_process_id}/{bk_module_name}").To(ps.DeleteModuleProcessBind))
    
    v3WS.Route(v3WS.GET("/{" + common.BKOwnerIDField + "}/{" + common.BKAppIDField + "}/{" + common.BKProcIDField + "}").To(ps.GetProcessDetailByID))
    
    v3WS.Route(v3WS.POST("/openapi/GetProcessPortByApplicationID/{" + common.BKAppIDField + "}").To(ps.GetProcessPortByApplicationID))
    v3WS.Route(v3WS.POST("/openapi/GetProcessPortByIP").To(ps.GetProcessPortByIP))
    
    container.Add(v3WS)
    
    return container
}

func (ps *ProcServer) OnProcessConfigUpdate(previous, current cfnc.ProcessConfig) {
    //
}