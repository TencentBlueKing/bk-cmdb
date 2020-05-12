/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package service

import (
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/rdapi"
	"configcenter/src/scene_server/container_server/app/options"
	"configcenter/src/scene_server/container_server/core"

	"github.com/emicklei/go-restful"
)

// ContainerServiceInterface the container service method used to init
type ContainerServiceInterface interface {
	// SetConfig set configs for container service
	SetConfig(config options.Config, engine *backbone.Engine, core core.Interface, ccErrIf errors.CCErrorIf, language language.CCLanguageIf)
	// WebService return restful Container
	WebService() *restful.Container
}

// ContainerService container service
type ContainerService struct {
	engine      *backbone.Engine
	core        core.Interface
	langFactory map[common.LanguageType]language.DefaultCCLanguageIf
	language    language.CCLanguageIf
	ccErrIf     errors.CCErrorIf
	cfg         options.Config
}

// New create container service instance
func New() *ContainerService {
	return &ContainerService{}
}

// SetConfig implements ContainerServiceInterface
func (s *ContainerService) SetConfig(
	cfg options.Config, engine *backbone.Engine, core core.Interface,
	ccErrIf errors.CCErrorIf, lang language.CCLanguageIf) {

	s.cfg = cfg
	s.engine = engine
	s.core = core

	if ccErrIf != nil {
		s.ccErrIf = ccErrIf
	}

	if lang != nil {
		s.langFactory = make(map[common.LanguageType]language.DefaultCCLanguageIf)
		s.langFactory[common.Chinese] = lang.CreateDefaultCCLanguageIf(string(common.Chinese))
		s.langFactory[common.English] = lang.CreateDefaultCCLanguageIf(string(common.English))
	}
}

// WebService implements ContainerServiceInterface
func (s *ContainerService) WebService() *restful.Container {
	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.ccErrIf
	}

	api.Path("/api/v3").
		Filter(s.engine.Metric().RestfulMiddleWare).
		Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	s.initService(api)

	healthz := new(restful.WebService).Produces(restful.MIME_JSON)
	healthz.Route(healthz.GET("/healthz").To(s.Healthz))
	container := restful.NewContainer().Add(api)
	container.Add(healthz)

	return container
}
