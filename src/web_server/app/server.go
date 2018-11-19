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
	"fmt"
	"os"
	"strings"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/storage/dal/redis"
	confCenter "configcenter/src/web_server/app/config"
	"configcenter/src/web_server/app/options"
	"configcenter/src/web_server/logics"
	"configcenter/src/web_server/middleware"
	websvc "configcenter/src/web_server/service"
)

type WebServer struct {
	Config options.Config
}

func Run(ctx context.Context, op *options.ServerOption) error {

	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	c := &util.APIMachineryConfig{
		ZkAddr:    op.ServConf.RegDiscover,
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}

	machinery, err := apimachinery.NewApiMachinery(c)
	if err != nil {
		return fmt.Errorf("new api machinery failed, err: %v", err)
	}

	service := new(websvc.Service)
	service.Disc, err = discovery.NewDiscoveryInterface(op.ServConf.RegDiscover)
	if err != nil {
		return fmt.Errorf("new proxy discovery instance failed, err: %v", err)
	}

	webSvr := new(WebServer)
	webSvr.getConfig(op.ServConf.RegDiscover)
	service.Config = webSvr.Config
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    service.WebService(),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_WEBSERVER, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_WEBSERVER,
		op.ServConf.ExConfig,
		webSvr.onServerConfigUpdate,
		bonC)

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
		return fmt.Errorf("Configuration item not found")
	}

	redisAddress := webSvr.Config.Session.Host
	redisSecret := strings.TrimSpace(webSvr.Config.Session.Secret)

	if !strings.Contains(redisAddress, ":") && len(webSvr.Config.Session.Port) > 0 {
		redisAddress = webSvr.Config.Session.Host + ":" + webSvr.Config.Session.Port
	}

	cacheCli, err := redis.NewFromConfig(redis.Config{
		Address:    redisAddress,
		Password:   redisSecret,
		MasterName: webSvr.Config.Session.MasterName,
		Database:   "0",
	})

	if nil != err {
		return err
	}

	service.Engine = engine
	service.Logics = &logics.Logics{Engine: engine}
	service.Config = webSvr.Config
	middleware.Engine = engine
	middleware.CacheCli = cacheCli

	select {}
	return nil

}

func (w *WebServer) getConfig(regDiscover string) error {
	cfCenter := confCenter.NewConfCenter(regDiscover)
	cfCenter.GetConfigOnce()

	/// fetch config of itselft
	var confData []byte
	_ = confData
	for {
		confData = cfCenter.GetConfigureCtx()
		if confData == nil {
			blog.Warnf("fail to get configure, will get again")
			time.Sleep(time.Second * 2)
			continue
		} else {
			blog.Infof("get configure. ctx(%s)", string(confData))
			cfCenter.Stop()
			break
		}
	}
	config := cfCenter.ParseConf(confData)
	w.Config.Site.DomainUrl = config["site.domain_url"] + "/"
	w.Config.Site.HtmlRoot = config["site.html_root"]
	w.Config.Site.ResourcesPath = config["site.resources_path"]
	w.Config.Site.BkLoginUrl = config["site.bk_login_url"]
	w.Config.Site.AppCode = config["site.app_code"]
	w.Config.Site.CheckUrl = config["site.check_url"]
	w.Config.Site.AccountUrl = config["site.bk_account_url"]
	w.Config.Site.BkHttpsLoginUrl = config["site.bk_https_login_url"]
	w.Config.Site.HttpsDomainUrl = config["site.https_domain_url"]

	w.Config.Session.Name = config["session.name"]
	w.Config.Session.Skip = config["session.skip"]
	w.Config.Session.Host = config["session.host"]
	w.Config.Session.Port = config["session.port"]
	w.Config.Session.Address = config["session.address"]
	w.Config.Session.Secret = strings.TrimSpace(config["session.secret"])
	w.Config.Session.MultipleOwner = config["session.multiple_owner"]
	w.Config.Session.DefaultLanguage = config["session.defaultlanguage"]
	w.Config.Session.MasterName = config["session.mastername"]
	w.Config.LoginVersion = config["login.version"]
	if "" == w.Config.Session.DefaultLanguage {
		w.Config.Session.DefaultLanguage = "zh-cn"
	}
	w.Config.Version = config["api.version"]
	w.Config.AgentAppUrl = config["app.agent_app_url"]
	w.Config.LoginUrl = fmt.Sprintf(w.Config.Site.BkLoginUrl, w.Config.Site.AppCode, w.Config.Site.DomainUrl)
	w.Config.ConfigMap = config
	return nil
}

func (w *WebServer) onServerConfigUpdate(previous, current cc.ProcessConfig) {
	w.Config.Site.DomainUrl = current.ConfigMap["site.domain_url"] + "/"
	w.Config.Site.HtmlRoot = current.ConfigMap["site.html_root"]
	w.Config.Site.ResourcesPath = current.ConfigMap["site.resources_path"]
	w.Config.Site.BkLoginUrl = current.ConfigMap["site.bk_login_url"]
	w.Config.Site.AppCode = current.ConfigMap["site.app_code"]
	w.Config.Site.CheckUrl = current.ConfigMap["site.check_url"]
	w.Config.Site.AccountUrl = current.ConfigMap["site.bk_account_url"]
	w.Config.Site.BkHttpsLoginUrl = current.ConfigMap["site.bk_https_login_url"]
	w.Config.Site.HttpsDomainUrl = current.ConfigMap["site.https_domain_url"]

	w.Config.Session.Name = current.ConfigMap["session.name"]
	w.Config.Session.Skip = current.ConfigMap["session.skip"]
	w.Config.Session.Host = current.ConfigMap["session.host"]
	w.Config.Session.Port = current.ConfigMap["session.port"]
	w.Config.Session.Address = current.ConfigMap["session.address"]
	w.Config.Session.MasterName = current.ConfigMap["session.mastername"]
	w.Config.Session.Secret = strings.TrimSpace(current.ConfigMap["session.secret"])
	w.Config.Session.MultipleOwner = current.ConfigMap["session.multiple_owner"]
	w.Config.Session.DefaultLanguage = current.ConfigMap["session.defaultlanguage"]
	w.Config.LoginVersion = current.ConfigMap["login.version"]
	if "" == w.Config.Session.DefaultLanguage {
		w.Config.Session.DefaultLanguage = "zh-cn"
	}

	w.Config.Version = current.ConfigMap["api.version"]
	w.Config.AgentAppUrl = current.ConfigMap["app.agent_app_url"]
	w.Config.LoginUrl = fmt.Sprintf(w.Config.Site.BkLoginUrl, w.Config.Site.AppCode, w.Config.Site.DomainUrl)
	w.Config.ConfigMap = current.ConfigMap

}

//Stop the ccapi server
func (ccWeb *WebServer) Stop() error {
	return nil
}

func newServerInfo(op *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := op.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := op.ServConf.GetPort()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}
	return info, nil
}
