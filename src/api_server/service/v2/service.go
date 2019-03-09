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

package v2

import (
	"context"
	"net/http"

	"github.com/emicklei/go-restful"

	logics "configcenter/src/api_server/logics/v2"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/util"
)

type Service struct {
	*backbone.Engine
	//*logics.Logics
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
		lgc:           logics.NewLogics(s.Engine, header),
	}
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/api/v2").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("App/getapplist").To(s.getAppList))
	ws.Route(ws.POST("app/getapplist").To(s.getAppList))
	ws.Route(ws.POST("App/getAppByID").To(s.getAppByID))
	ws.Route(ws.POST("app/getAppByID").To(s.getAppByID))
	ws.Route(ws.POST("app/getappbyid").To(s.getAppByID))
	ws.Route(ws.POST("App/getappbyuin").To(s.getAppByUin))
	ws.Route(ws.POST("app/getappbyuin").To(s.getAppByUin))
	ws.Route(ws.POST("User/getUserRoleApp").To(s.getUserRoleApp))
	ws.Route(ws.POST("user/getUserRoleApp").To(s.getUserRoleApp))
	ws.Route(ws.POST("TopSetModule/getappsetmoduletreebyappid").To(s.getAppSetModuleTreeByAppId))
	ws.Route(ws.POST("app/addApp").To(s.addApp))
	ws.Route(ws.POST("app/deleteApp").To(s.deleteApp))
	ws.Route(ws.POST("app/editApp").To(s.editApp))
	ws.Route(ws.POST("App/getHostAppByCompanyId").To(s.getHostAppByCompanyId))

	ws.Route(ws.POST("Module/getmodules").To(s.getModulesByApp))
	ws.Route(ws.POST("module/editmodule").To(s.updateModule))
	ws.Route(ws.POST("module/addModule").To(s.addModule))
	ws.Route(ws.POST("module/delModule").To(s.deleteModule))

	ws.Route(ws.POST("Set/getsetsbyproperty").To(s.getSets))
	ws.Route(ws.POST("Set/getsetproperty").To(s.getsetproperty))
	ws.Route(ws.POST("Set/getmodulesbyproperty").To(s.getModulesByProperty))
	ws.Route(ws.POST("set/getmodulesbyproperty").To(s.getModulesByProperty))
	ws.Route(ws.POST("set/addset").To(s.addSet))
	ws.Route(ws.POST("set/updateset").To(s.updateSet))
	ws.Route(ws.POST("set/updateSetServiceStatus").To(s.updateSetServiceStatus))
	ws.Route(ws.POST("set/delset").To(s.delSet))
	ws.Route(ws.POST("set/delSetHost").To(s.delSetHost))

	ws.Route(ws.POST("host/addhost").To(s.addHost))
	ws.Route(ws.POST("host/enterIp").To(s.enterIP))
	ws.Route(ws.POST("host/enterip").To(s.enterIP))

	ws.Route(ws.POST("host/getAgentStatus").To(s.getAgentStatus))

	ws.Route(ws.POST("Host/gethostlistbyip").To(s.getHostListByIP))
	ws.Route(ws.POST("host/gethostlistbyip").To(s.getHostListByIP))
	ws.Route(ws.POST("Host/getsethostlist").To(s.getSetHostList))
	ws.Route(ws.POST("host/getmodulehostlist").To(s.getModuleHostList))
	ws.Route(ws.POST("Host/getmodulehostlist").To(s.getModuleHostList))
	ws.Route(ws.POST("host/getapphostlist").To(s.getAppHostList))
	ws.Route(ws.POST("Host/getapphostlist").To(s.getAppHostList))
	ws.Route(ws.POST("set/gethostsbyproperty").To(s.getHostsByProperty))
	ws.Route(ws.POST("Set/gethostsbyproperty").To(s.getHostsByProperty))
	ws.Route(ws.POST("Host/updateHostStatus").To(s.updateHostStatus))

	ws.Route(ws.POST("Host/updateHostByAppId").To(s.updateHostByAppID))
	ws.Route(ws.POST("Host/getCompanyIdByIps").To(s.getCompanyIDByIps))
	ws.Route(ws.POST("host/getCompanyIdByIps").To(s.getCompanyIDByIps))
	ws.Route(ws.POST("Host/getHostListByAppidAndField").To(s.getHostListByAppIDAndField))
	ws.Route(ws.POST("host/getHostListByAppidAndField").To(s.getHostListByAppIDAndField))
	ws.Route(ws.POST("Host/getIPAndProxyByCompany").To(s.getIPAndProxyByCompany))
	ws.Route(ws.POST("Host/updatehostmodule").To(s.updateHostModule))
	ws.Route(ws.POST("host/updatehostmodule").To(s.updateHostModule))
	ws.Route(ws.POST("host/updateCustomProperty").To(s.updateCustomProperty))
	ws.Route(ws.POST("host/cloneHostProperty").To(s.cloneHostProperty))
	ws.Route(ws.POST("host/delHostInApp").To(s.delHostInApp))
	ws.Route(ws.POST("host/getgitServerIp").To(s.getGitServerIp))
	ws.Route(ws.POST("host/hardinfo").To(s.GetHostHardInfo))

	ws.Route(ws.POST("/CustomerGroup/getContentByCustomerGroupID").To(s.getContentByCustomerGroupID))
	ws.Route(ws.POST("CustomerGroup/getContentByCustomerGroupId").To(s.getContentByCustomerGroupID))
	ws.Route(ws.POST("/CustomerGroup/getCustomerGroupList").To(s.getCustomerGroupList))

	ws.Route(ws.POST("Plat/updateHost").To(s.updateHost))
	ws.Route(ws.POST("Plat/get").To(s.getPlats))
	ws.Route(ws.POST("Plat/get").To(s.getPlats))
	ws.Route(ws.POST("Plat/delete").To(s.deletePlats))
	ws.Route(ws.POST("Plat/add").To(s.createPlats))

	ws.Route(ws.POST("process/getProcessPortByApplicationID").To(s.getProcessPortByApplicationID))
	ws.Route(ws.POST("process/getProcessPortByIP").To(s.getProcessPortByIP))

	ws.Route(ws.POST("Property/getList").To(s.getObjProperty))

	ws.Route(ws.POST("/Host/getAppOwnerHostList").To(s.getHostListByOwner))
	ws.Route(ws.POST("/App/getAppByUinExt").To(s.getAppByOwnerAndUin))

	return ws

}

func (s *Service) SetEngine(engine *backbone.Engine) {
	s.Engine = engine
}
