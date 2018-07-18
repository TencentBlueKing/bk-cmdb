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

	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/rdapi"
	"configcenter/src/source_controller/auditcontroller/logics"
	"configcenter/src/storage"
)

type Service struct {
	*backbone.Engine
	*logics.Logics
	Instance storage.DI
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFun := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/audit/{version}").Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)

	ws.Route(ws.POST("/host/{owner_id}/{biz_id}/{user}").To(s.AddHostLog))
	ws.Route(ws.POST("/hosts/{owner_id}/{biz_id}/{user}").To(s.AddHostLogs))
	ws.Route(ws.POST("/obj/{owner_id}/{biz_id}/{user}").To(s.AddObjectLog))
	ws.Route(ws.POST("/objs/{owner_id}/{biz_id}/{user}").To(s.AddObjectLogs))
	ws.Route(ws.POST("/proc/{owner_id}/{biz_id}/{user}").To(s.AddProcLog))
	ws.Route(ws.POST("/procs/{owner_id}/{biz_id}/{user}").To(s.AddProcLogs))
	ws.Route(ws.POST("/module/{owner_id}/{biz_id}/{user}").To(s.AddModuleLog))
	ws.Route(ws.POST("/modules/{owner_id}/{biz_id}/{user}").To(s.AddModuleLogs))
	ws.Route(ws.POST("/app/{owner_id}/{biz_id}/{user}").To(s.AddAppLog))
	ws.Route(ws.POST("set/{owner_id}/{biz_id}/{user}").To(s.AddSetLog))
	ws.Route(ws.POST("/sets/{owner_id}/{biz_id}/{user}").To(s.AddSetLogs))
	ws.Route(ws.POST("/search").To(s.Get))
	return ws
}
