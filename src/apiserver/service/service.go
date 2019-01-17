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
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apiserver/core"
	compatiblev2 "configcenter/src/apiserver/core/compatiblev2/service"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/rdapi"

	"github.com/emicklei/go-restful"
)

// Service service methods
type Service interface {
	WebServices() []*restful.WebService
	SetConfig(engine *backbone.Engine, httpClient HTTPClient, discovery discovery.DiscoveryInterface)
}

// NewService create a new service instance
func NewService() Service {
	return &service{
		core: core.New(nil, compatiblev2.New(nil)),
	}
}

type service struct {
	engine    *backbone.Engine
	client    HTTPClient
	core      core.Core
	discovery discovery.DiscoveryInterface
}

func (s *service) SetConfig(engine *backbone.Engine, httpClient HTTPClient, discovery discovery.DiscoveryInterface) {
	s.engine = engine
	s.client = httpClient
	s.discovery = discovery
	s.core.CompatibleV2Operation().SetConfig(engine)
}

func (s *service) WebServices() []*restful.WebService {

	allWebServices := []*restful.WebService{}

	getErrFun := func() errors.CCErrorIf {
		return s.engine.CCErr
	}

	// init V3
	ws := &restful.WebService{}

	ws.Path(rootPath).Filter(rdapi.AllGlobalFilter(getErrFun)).Produces(restful.MIME_JSON)
	ws.Route(ws.GET("{.*}").Filter(s.URLFilterChan).To(s.Get))
	ws.Route(ws.POST("{.*}").Filter(s.URLFilterChan).To(s.Post))
	ws.Route(ws.PUT("{.*}").Filter(s.URLFilterChan).To(s.Put))
	ws.Route(ws.DELETE("{.*}").Filter(s.URLFilterChan).To(s.Delete))

	allWebServices = append(allWebServices, ws)

	// init v2
	allWebServices = append(allWebServices, s.core.CompatibleV2Operation().WebService())

	return allWebServices
}
