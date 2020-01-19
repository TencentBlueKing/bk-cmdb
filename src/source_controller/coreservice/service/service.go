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
	"time"

	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/rdapi"
	"configcenter/src/source_controller/coreservice/app/options"
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
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/remote"
	dalredis "configcenter/src/storage/dal/redis"

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
	engine   *backbone.Engine
	language language.CCLanguageIf
	err      errors.CCErrorIf
	cfg      options.Config
	core     core.Core
	db       dal.RDB
	cache    *redis.Client
}

func (s *coreService) SetConfig(cfg options.Config, engine *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error {

	s.cfg = cfg
	s.engine = engine

	if nil != err {
		s.err = err
	}

	if nil != language {
		s.language = language
	}

	var dbErr error
	var db dal.RDB
	if s.cfg.Mongo.Enable == "true" {
		db, dbErr = local.NewMgo(s.cfg.Mongo.BuildURI(), time.Minute)
	} else {
		db, dbErr = remote.NewWithDiscover(s.engine)
	}
	if dbErr != nil {
		blog.Errorf("failed to connect the txc server, error info is %s", dbErr.Error())
		return dbErr
	}

	cache, cacheRrr := dalredis.NewFromConfig(cfg.Redis)
	if cacheRrr != nil {
		blog.Errorf("new redis client failed, err: %v", cacheRrr)
		return cacheRrr
	}

	s.db = db
	s.cache = cache

	// connect the remote mongodb
	instance := instances.New(db, s, cache, language)
	hostApplyRuleCore := hostapplyrule.New(db, instance)
	s.core = core.New(
		model.New(db, s, language, cache),
		instance,
		association.New(db, s),
		datasynchronize.New(db, s),
		mainline.New(db),
		host.New(db, cache, s, hostApplyRuleCore),
		auditlog.New(db),
		process.New(db, s, cache),
		label.New(db),
		settemplate.New(db),
		operation.New(db),
		hostApplyRuleCore,
		dbSystem.New(db),
	)
	return nil
}

// WebService the web service
func (s *coreService) WebService() *restful.Container {

	container := restful.NewContainer()

	api := new(restful.WebService)
	getErrFunc := func() errors.CCErrorIf {
		return s.err
	}
	api.Path("/api/v3").Filter(rdapi.AllGlobalFilter(getErrFunc)).Produces(restful.MIME_JSON)

	// init service actions
	s.initService(api)
	container.Add(api)

	healthzAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	healthzAPI.Route(healthzAPI.GET("/healthz").To(s.Healthz))
	container.Add(healthzAPI)

	return container
}
