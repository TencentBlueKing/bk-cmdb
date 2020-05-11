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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/coreservice/app/options"
	"configcenter/src/source_controller/coreservice/cache"
	cacheop "configcenter/src/source_controller/coreservice/cache"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/source_controller/coreservice/core/association"
	"configcenter/src/source_controller/coreservice/core/auditlog"
	"configcenter/src/source_controller/coreservice/core/datasynchronize"
	"configcenter/src/source_controller/coreservice/core/host"
	"configcenter/src/source_controller/coreservice/core/hostapplyrule"
	"configcenter/src/source_controller/coreservice/core/instances"
	"configcenter/src/source_controller/coreservice/core/label"
	"configcenter/src/source_controller/coreservice/core/mainline"
	"configcenter/src/source_controller/coreservice/core/model"
	"configcenter/src/source_controller/coreservice/core/operation"
	"configcenter/src/source_controller/coreservice/core/process"
	"configcenter/src/source_controller/coreservice/core/settemplate"
	dbSystem "configcenter/src/source_controller/coreservice/core/system"
	watchEvent "configcenter/src/source_controller/coreservice/event"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	dalredis "configcenter/src/storage/dal/redis"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream"

	"github.com/emicklei/go-restful"
	"gopkg.in/redis.v5"
)

// CoreServiceInterface the topo service methods used to init
type CoreServiceInterface interface {
	WebService() *restful.Container
	SetConfig(cfg options.Config, engine *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error
}

// New create topo service instance
func New() CoreServiceInterface {
	return &coreService{}
}

// coreService topo service
type coreService struct {
	engine      *backbone.Engine
	langFactory map[common.LanguageType]language.DefaultCCLanguageIf
	language    language.CCLanguageIf
	err         errors.CCErrorIf
	cfg         options.Config
	core        core.Core
	db          dal.RDB
	rds         *redis.Client
	cacheSet    *cache.ClientSet
}

func (s *coreService) SetConfig(cfg options.Config, engine *backbone.Engine, err errors.CCErrorIf, lang language.CCLanguageIf) error {

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

	db, dbErr := local.NewMgo(s.cfg.Mongo.GetMongoConf(), time.Minute)
	if dbErr != nil {
		blog.Errorf("failed to connect the txc server, error info is %s", dbErr.Error())
		return dbErr
	}

	cache, cacheRrr := dalredis.NewFromConfig(cfg.Redis)
	if cacheRrr != nil {
		blog.Errorf("new redis client failed, err: %v", cacheRrr)
		return cacheRrr
	}

	initErr := db.InitTxnManager(cache)
	if initErr != nil {
		blog.Errorf("failed to init txn manager, error info is %v", initErr)
		return initErr
	}

	s.db = db
	s.rds = cache

	// connect the remote mongodb
	instance := instances.New(db, s, cache, lang)
	hostApplyRuleCore := hostapplyrule.New(db, instance)
	s.core = core.New(
		model.New(db, s, lang, cache),
		instance,
		association.New(db, s),
		datasynchronize.New(db, s),
		mainline.New(db, lang),
		host.New(db, cache, s, hostApplyRuleCore),
		auditlog.New(db),
		process.New(db, s, cache),
		label.New(db),
		settemplate.New(db),
		operation.New(db),
		hostApplyRuleCore,
		dbSystem.New(db),
	)

	event, eventErr := reflector.NewReflector(s.cfg.Mongo.GetMongoConf())
	if eventErr != nil {
		blog.Errorf("new reflector failed, err: %v", eventErr)
		return eventErr
	}

	c, cacheErr := cacheop.NewCache(cache, db, event)
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

	if err := watchEvent.NewEvent(s.db, s.rds, watcher, engine.ServiceManageInterface); err != nil {
		blog.Errorf("new watch event failed, err: %v", err)
		return err
	}

	return nil
}

// WebService the web service
func (s *coreService) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	api.Path("/api/v3").Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)

	// init service actions
	s.initService(api)
	container.Add(api)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.Healthz))
	container.Add(healthzAPI)

	return container
}

func (s *coreService) Language(header http.Header) language.DefaultCCLanguageIf {
	lang := util.GetLanguage(header)
	l, exist := s.langFactory[common.LanguageType(lang)]
	if !exist {
		return s.langFactory[common.Chinese]
	}
	return l
}
