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

package app

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/logics"
	websvc "configcenter/src/web_server/service"

	"github.com/holmeswang/contrib/sessions"
)

type WebServer struct {
	Config options.Config
}

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {

	// init esb client
	esb.InitEsbClient(nil)

	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(websvc.Service)

	webSvr := new(WebServer)
	service.Config = &webSvr.Config
	input := &backbone.BackboneParameter{
		ConfigUpdate: webSvr.onServerConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if "" == webSvr.Config.Site.DomainUrl {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return errors.New("configuration item not found")
	}

	webSvr.Config.Redis, err = engine.WithRedis()
	if err != nil {
		return err
	}
	var redisErr error
	if webSvr.Config.Redis.MasterName == "" {
		// MasterName 为空，表示使用直连redis 。 使用Host,Port 做链接redis参数
		service.Session, redisErr = sessions.NewRedisStore(10, "tcp", webSvr.Config.Redis.Address, webSvr.Config.Redis.Password, []byte("secret"))
		if redisErr != nil {
			return fmt.Errorf("failed to create new redis store, error info is %v", redisErr)
		}
	} else {
		// MasterName 不为空，表示使用哨兵模式的redis。MasterName 是Master标记
		address := strings.Split(webSvr.Config.Redis.Address, ";")
		service.Session, redisErr = sessions.NewRedisStoreWithSentinel(address, 10, webSvr.Config.Redis.MasterName, "tcp", webSvr.Config.Redis.Password, []byte("secret"))
		if redisErr != nil {
			return fmt.Errorf("failed to create new redis store, error info is %v", redisErr)
		}
	}
	cacheCli, err := redis.NewFromConfig(webSvr.Config.Redis)

	if nil != err {
		return err
	}

	service.Engine = engine
	service.CacheCli = cacheCli
	service.Logics = &logics.Logics{Engine: engine}
	service.Config = &webSvr.Config

	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), false)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	}

	return nil
}

func (w *WebServer) onServerConfigUpdate(previous, current cc.ProcessConfig) {
	domainUrl, _ := cc.String("webServer.site.domainUrl")
	w.Config.Site.DomainUrl = domainUrl + "/"
	w.Config.Site.HtmlRoot, _ = cc.String("webServer.site.htmlRoot")
	w.Config.Site.ResourcesPath, _ = cc.String("webServer.site.resourcesPath")
	w.Config.Site.BkLoginUrl, _ = cc.String("webServer.site.bkLoginUrl")
	w.Config.Site.AppCode, _ = cc.String("webServer.site.appCode")
	w.Config.Site.CheckUrl, _ = cc.String("webServer.site.checkUrl")

	authscheme, err := cc.String("webServer.site.authscheme")
	if err != nil {
		w.Config.Site.AuthScheme = "internal"
	} else {
		w.Config.Site.AuthScheme = authscheme
	}

	fullTextSearch, err := cc.String("es.fullTextSearch")
	if err != nil {
		w.Config.Site.FullTextSearch = "off"
	} else {
		w.Config.Site.FullTextSearch = fullTextSearch
	}

	w.Config.Site.AccountUrl, _ = cc.String("webServer.site.bkAccountUrl")
	w.Config.Site.BkHttpsLoginUrl, _ = cc.String("webServer.site.bkHttpsLoginUrl")
	w.Config.Site.HttpsDomainUrl, _ = cc.String("webServer.site.httpsDomainUrl")
	w.Config.Site.PaasDomainUrl, _ = cc.String("webServer.site.paasDomainUrl")
	w.Config.Site.HelpDocUrl, _ = cc.String("webServer.site.helpDocUrl")

	w.Config.Session.Name, _ = cc.String("webServer.session.name")
	w.Config.Session.MultipleOwner, _ = cc.String("webServer.session.multipleOwner")
	w.Config.Session.DefaultLanguage, _ = cc.String("webServer.session.defaultlanguage")
	w.Config.LoginVersion, _ = cc.String("webServer.login.version")
	if "" == w.Config.Session.DefaultLanguage {
		w.Config.Session.DefaultLanguage = "zh-cn"
	}

	w.Config.Version, _ = cc.String("webServer.api.version")
	w.Config.AgentAppUrl, _ = cc.String("webServer.app.agentAppUrl")
	w.Config.AuthCenter.AppCode, _ = cc.String("webServer.app.authAppCode")
	w.Config.AuthCenter.URL, _ = cc.String("webServer.app.authUrl")
	w.Config.LoginUrl = fmt.Sprintf(w.Config.Site.BkLoginUrl, w.Config.Site.AppCode, w.Config.Site.DomainUrl)
	if esbConfig, err := esb.ParseEsbConfig("webServer"); err == nil {
		esb.UpdateEsbConfig(*esbConfig)
	}

}

//Stop the ccapi server
func (ccWeb *WebServer) Stop() error {
	return nil
}
