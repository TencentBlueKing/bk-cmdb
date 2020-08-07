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
	w.Config.Site.DomainUrl = current.ConfigMap["site.domain_url"] + "/"
	w.Config.Site.HtmlRoot = current.ConfigMap["site.html_root"]
	w.Config.Site.ResourcesPath = current.ConfigMap["site.resources_path"]
	w.Config.Site.BkLoginUrl = current.ConfigMap["site.bk_login_url"]
	w.Config.Site.AppCode = current.ConfigMap["site.app_code"]
	w.Config.Site.CheckUrl = current.ConfigMap["site.check_url"]
	w.Config.Site.AuthScheme = current.ConfigMap["site.authscheme"]
	if w.Config.Site.AuthScheme == "" {
		w.Config.Site.AuthScheme = "internal"
	}
	w.Config.Site.FullTextSearch = current.ConfigMap["site.full_text_search"]
	if w.Config.Site.FullTextSearch == "" {
		w.Config.Site.FullTextSearch = "off"
	}
	w.Config.Site.AccountUrl = current.ConfigMap["site.bk_account_url"]
	w.Config.Site.BkHttpsLoginUrl = current.ConfigMap["site.bk_https_login_url"]
	w.Config.Site.HttpsDomainUrl = current.ConfigMap["site.https_domain_url"]
	w.Config.Site.PaasDomainUrl = current.ConfigMap["site.paas_domain_url"]

	w.Config.Session.Name = current.ConfigMap["session.name"]
	w.Config.Session.MultipleOwner = current.ConfigMap["session.multiple_owner"]
	w.Config.Session.DefaultLanguage = current.ConfigMap["session.defaultlanguage"]
	w.Config.LoginVersion = current.ConfigMap["login.version"]
	if "" == w.Config.Session.DefaultLanguage {
		w.Config.Session.DefaultLanguage = "zh-cn"
	}

	w.Config.Version = current.ConfigMap["api.version"]
	w.Config.AgentAppUrl = current.ConfigMap["app.agent_app_url"]
	w.Config.AuthCenter.AppCode = current.ConfigMap["app.auth_app_code"]
	w.Config.AuthCenter.URL = current.ConfigMap["app.auth_url"]
	w.Config.LoginUrl = fmt.Sprintf(w.Config.Site.BkLoginUrl, w.Config.Site.AppCode, w.Config.Site.DomainUrl)
	w.Config.ConfigMap = current.ConfigMap

	if esbConfig, err := esb.ParseEsbConfig(current.ConfigMap); err == nil {
		esb.UpdateEsbConfig(*esbConfig)
	}

}

//Stop the ccapi server
func (ccWeb *WebServer) Stop() error {
	return nil
}
