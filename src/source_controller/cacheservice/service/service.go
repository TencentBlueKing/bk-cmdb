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

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/app/options"
	"configcenter/src/source_controller/cacheservice/cache"
	cacheop "configcenter/src/source_controller/cacheservice/cache"
	watchEvent "configcenter/src/source_controller/cacheservice/event"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream"

	"github.com/emicklei/go-restful"
)

// CacheServiceInterface the cache service methods used to init
type CacheServiceInterface interface {
	WebService() *restful.Container
	SetConfig(cfg options.Config, engine *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error
}

// New create cache service instance
func New() CacheServiceInterface {
	return &cacheService{}
}

// cacheService cache service
type cacheService struct {
	engine      *backbone.Engine
	langFactory map[common.LanguageType]language.DefaultCCLanguageIf
	language    language.CCLanguageIf
	err         errors.CCErrorIf
	cfg         options.Config
	core        core.Core
	cacheSet    *cache.ClientSet
}

func (s *cacheService) SetConfig(cfg options.Config, engine *backbone.Engine, err errors.CCErrorIf, lang language.CCLanguageIf) error {

	s.cfg = cfg
	s.engine = engine

	if nil != err {
		s.err = err
	}

	if nil != lang {
		s.langFactory = make(map[common.LanguageType]language.DefaultCCLanguageIf)
		s.langFactory[common.Chinese] = lang.CreateDefaultCCLanguageIf(string(common.Chinese))
		s.langFactory[common.English] = lang.CreateDefaultCCLanguageIf(string(common.English))
	}

	event, eventErr := reflector.NewReflector(s.cfg.Mongo.GetMongoConf())
	if eventErr != nil {
		blog.Errorf("new reflector failed, err: %v", eventErr)
		return eventErr
	}

	c, cacheErr := cacheop.NewCache(event)
	if cacheErr != nil {
		blog.Errorf("new cache instance failed, err: %v", cacheErr)
		return cacheErr
	}
	s.cacheSet = c

	watcher, watchErr := stream.NewStream(s.cfg.Mongo.GetMongoConf())
	if watchErr != nil {
		blog.Errorf("new watch stream failed, err: %v", watchErr)
		return watchErr
	}

	if err := watchEvent.NewEvent(watcher, engine.ServiceManageInterface); err != nil {
		blog.Errorf("new watch event failed, err: %v", err)
		return err
	}

	return nil
}

// WebService the web service
func (s *cacheService) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	api.Path("/cache/v3").Filter(s.engine.Metric().RestfulMiddleWare).Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)

	// init service actions
	s.initService(api)
	container.Add(api)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.Healthz))
	container.Add(healthzAPI)

	return container
}

func (s *cacheService) Language(header http.Header) language.DefaultCCLanguageIf {
	lang := util.GetLanguage(header)
	l, exist := s.langFactory[common.LanguageType(lang)]
	if !exist {
		return s.langFactory[common.Chinese]
	}
	return l
}
