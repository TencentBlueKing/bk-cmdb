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
	"configcenter/src/apimachinery/synchronize"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metric"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/synchronize_server/app/options"
	"configcenter/src/scene_server/synchronize_server/logics"

	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"
)

type Service struct {
	*options.Config
	*backbone.Engine
	disc           discovery.DiscoveryInterface
	CacheDB        *redis.Client
	synchronizeSrv synchronize.SynchronizeClientInterface
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

func (s *Service) SetSynchronizeServer(synchronizeSrv synchronize.SynchronizeClientInterface) {
	s.synchronizeSrv = synchronizeSrv
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
		lgc:           logics.NewLogics(s.Engine, header, s.CacheDB, s.synchronizeSrv),
	}
}

func (s *Service) WebService() *restful.WebService {
	ws := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.CCErr
	}
	ws.Path("/synchronize/{version}").Filter(s.Engine.Metric().RestfulMiddleWare).Filter(rdapi.HTTPRequestIDFilter(getErrFunc)).Produces(restful.MIME_JSON)

	ws.Route(ws.POST("/search").To(s.Find))
	ws.Route(ws.POST("/set/identifier/flag").To(s.SetIdentifierFlag))

	return ws
}

func (s *Service) Healthz(req *restful.Request, resp *restful.Response) {
	meta := metric.HealthMeta{IsHealthy: true}
	answer := metric.HealthResponse{
		Code:    common.CCSuccess,
		OK:      meta.IsHealthy,
		Result:  meta.IsHealthy,
		Message: meta.Message,
	}
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteEntity(answer)
}

// InitBackground  initialization backgroud task
func (s *Service) InitBackground() {
	header := make(http.Header, 0)
	if "" == util.GetOwnerID(header) {
		header.Set(common.BKHTTPOwnerID, common.BKSuperOwnerID)
		header.Set(common.BKHTTPHeaderUser, common.BKSynchronizeDataTaskDefaultUser)
	}

	srvData := s.newSrvComm(header)
	go srvData.lgc.TriggerSynchronize(srvData.ctx, s.Config)
}
