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

// Package service TODO
package service

import (
	"fmt"
	"net/http"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/language"
	"configcenter/src/common/rdapi"
	apigwcli "configcenter/src/common/resource/apigw"
	"configcenter/src/common/webservice/restfulservice"
	"configcenter/src/source_controller/cacheservice/app/options"
	"configcenter/src/source_controller/cacheservice/cache"
	cacheop "configcenter/src/source_controller/cacheservice/cache"
	"configcenter/src/source_controller/cacheservice/event/bsrelation"
	"configcenter/src/source_controller/cacheservice/event/flow"
	"configcenter/src/source_controller/cacheservice/event/identifier"
	"configcenter/src/source_controller/coreservice/core"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/stream/scheduler"
	"configcenter/src/thirdparty/apigw"
	"configcenter/src/thirdparty/logplatform/opentelemetry"

	"github.com/emicklei/go-restful/v3"
)

// CacheServiceInterface the cache service methods used to init
type CacheServiceInterface interface {
	WebService() *restful.Container
	SetConfig(cfg options.Config, engine *backbone.Engine, err errors.CCErrorIf, language language.CCLanguageIf) error
	Scheduler() *scheduler.Scheduler
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
	authManager *extensions.AuthManager
	scheduler   *scheduler.Scheduler
}

// SetConfig TODO
func (s *cacheService) SetConfig(cfg options.Config, engine *backbone.Engine, errf errors.CCErrorIf,
	lang language.CCLanguageIf) error {

	s.cfg = cfg
	s.engine = engine

	if errf != nil {
		s.err = errf
	}

	s.cfg.DisableVerifyTenant, _ = cc.Bool("tenant.disableVerifyTenant")
	s.cfg.EnableMultiTenantMode, _ = cc.Bool("tenant.enableMultiTenantMode")
	if nil != lang {
		s.langFactory = make(map[common.LanguageType]language.DefaultCCLanguageIf)
		s.langFactory[common.Chinese] = lang.CreateDefaultCCLanguageIf(string(common.Chinese))
		s.langFactory[common.English] = lang.CreateDefaultCCLanguageIf(string(common.English))
	}

	ccErr := s.InitClients()
	if ccErr != nil {
		fmt.Errorf("new apigw client failed: %v", ccErr)
	}
	iamCli := new(iam.IAM)
	if auth.EnableAuthorize() {
		var rawErr error
		iamCli, rawErr = iam.NewIAM()
		if rawErr != nil {
			return fmt.Errorf("new iam client failed: %v", rawErr)
		}
		s.authManager = extensions.NewAuthManager(engine.CoreAPI, iamCli)
	}

	taskScheduler, err := scheduler.New(mongodb.Dal(), mongodb.Dal("watch"), engine.ServiceManageInterface)
	if err != nil {
		blog.Errorf("new watch task scheduler instance failed, err: %v", err)
		return err
	}
	s.scheduler = taskScheduler

	c, cacheErr := cacheop.NewCache(engine.ServiceManageInterface)
	if cacheErr != nil {
		blog.Errorf("new cache instance failed, err: %v", cacheErr)
		return cacheErr
	}
	s.cacheSet = c
	if err = taskScheduler.AddTasks(c.GetWatchTasks()...); err != nil {
		return err
	}

	flowEvent, flowErr := flow.NewEvent()
	if flowErr != nil {
		blog.Errorf("new watch event failed, err: %v", flowErr)
		return flowErr
	}
	if err = taskScheduler.AddTasks(flowEvent.GetWatchTasks()...); err != nil {
		return err
	}

	hostIdentity, err := identifier.NewIdentity()
	if err != nil {
		blog.Errorf("new host identity event failed, err: %v", err)
		return err
	}
	if err = taskScheduler.AddTasks(hostIdentity.GetWatchTasks()...); err != nil {
		return err
	}

	bsRelation, err := bsrelation.NewBizSetRelation()
	if err != nil {
		blog.Errorf("new biz set relation event failed, err: %v", err)
		return err
	}
	if err = taskScheduler.AddTasks(bsRelation.GetWatchTasks()...); err != nil {
		return err
	}

	if err = taskScheduler.Start(); err != nil {
		blog.Errorf("start event watch task scheduler failed, err: %v", err)
		return err
	}

	return nil
}

// WebService the web service
func (s *cacheService) WebService() *restful.Container {

	container := restful.NewContainer()

	opentelemetry.AddOtlpFilter(container)

	getErrFunc := func() errors.CCErrorIf { return s.err }

	api := new(restful.WebService)
	api.Path("/cache/v3").Filter(s.engine.Metric().RestfulMiddleWare).Filter(rdapi.AllGlobalFilter(getErrFunc)).
		Produces(restful.MIME_JSON).Consumes(restful.MIME_JSON)

	// init service actions
	s.initService(api)
	container.Add(api)

	// common api
	commonAPI := new(restful.WebService).Produces(restful.MIME_JSON)
	commonAPI.Route(commonAPI.GET("/healthz").To(s.Healthz))
	commonAPI.Route(commonAPI.GET("/version").To(restfulservice.Version))
	commonAPI.Route(commonAPI.POST("/refresh/tenants").To(s.RefreshTenant))
	container.Add(commonAPI)

	return container
}

// Scheduler returns the watch task scheduler
func (s *cacheService) Scheduler() *scheduler.Scheduler {
	return s.scheduler
}

// Language TODO
func (s *cacheService) Language(header http.Header) language.DefaultCCLanguageIf {
	lang := httpheader.GetLanguage(header)
	l, exist := s.langFactory[common.LanguageType(lang)]
	if !exist {
		return s.langFactory[common.Chinese]
	}
	return l
}

// InitClients init apiGW client
func (s *cacheService) InitClients() error {

	var clients []apigw.ClientType
	if s.cfg.EnableMultiTenantMode && !s.cfg.DisableVerifyTenant {
		clients = []apigw.ClientType{apigw.User}
	}

	if auth.EnableAuthorize() {
		clients = append(clients, apigw.Iam)
	}

	if len(clients) > 0 {
		err := apigwcli.Init("apiGW", s.engine.Metric().Registry(), clients)
		if err != nil {
			blog.Errorf("init gse api gateway client failed, err: %v", err)
			return err
		}
	}

	return nil
}
