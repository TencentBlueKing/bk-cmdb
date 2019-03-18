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
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/auth/authcenter"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/admin_server/synchronizer"
	"configcenter/src/scene_server/auth_synchronizer/app/options"
	webservice "configcenter/src/scene_server/auth_synchronizer/pkg/service"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

// Run start http service and synchroizer
func Run(ctx context.Context, serverOptions *options.ServerOption) error {
	blog.Info("AuthSynchronizer app.server start...")
	httpServerConfig, err := NewHTTPServerConfig(serverOptions)
	if err != nil {
		return fmt.Errorf("extract http server config failed, err: %v", err)
	}

	discoverClient, err := discovery.NewDiscoveryInterface(serverOptions.ServConf.RegDiscover)
	if err != nil {
		return fmt.Errorf("connect zookeeper [%s] failed: %v", serverOptions.ServConf.RegDiscover, err)
	}

	apiMachineryConfig := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: &util.TLSClientConfig{InsecureSkipVerify: true},
	}

	apiMachineryClient, err := apimachinery.NewApiMachinery(apiMachineryConfig, discoverClient)
	if err != nil {
		return fmt.Errorf("new api machinery client failed, err: %v", err)
	}

	h.Config.Auth, err = authcenter.ParseConfigFromKV("auth", current.ConfigMap)
	if err != nil {
		blog.Warnf("parse authcenter config failed: %v", err)
	}
	serviceContainer := new(webservice.Service)
	server := backbone.Server{
		ListenAddr: httpServerConfig.IP,
		ListenPort: httpServerConfig.Port,
		Handler:    restful.NewContainer().Add(serviceContainer.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	// regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_AUTH_SYNCHROIZER, httpServerConfig.IP)
	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_DATACOLLECTION, httpServerConfig.IP)
	backbonConfig := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *httpServerConfig,
		CoreAPI:      apiMachineryClient,
		Server:       server,
	}

	synchronizerConfig := new(SynchronizerConfig)

	// note: NewBackbone will run a server in backend to sync newest config to
	engine, err := backbone.NewBackbone(
		ctx,
		serverOptions.ServConf.RegDiscover,
		// types.CC_MODULE_AUTH_SYNCHROIZER,
		types.CC_MODULE_DATACOLLECTION,
		serverOptions.ServConf.ExConfig,
		synchronizerConfig.onHostConfigUpdate,
		discoverClient,
		backbonConfig,
	)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	synchronizerConfig.Core = engine
	synchronizerConfig.Service = serviceContainer
	// wait for synchronizerConfig.Config
	for {
		if synchronizerConfig.Config != nil {
			break
		}
		time.Sleep(time.Second * 2)
		blog.Info("config not found, retry 2s later")
	}

	esbChan := make(chan esbutil.EsbConfig, 1)
	esbChan <- synchronizerConfig.Config.Esb
	sync := synchronizer.NewSynchronizer(ctx, synchronizerConfig.Config, synchronizerConfig.Core)
	blog.Info("begin to start synchronizer...")
	err = sync.Run()
	if err != nil {
		return fmt.Errorf("run auth synchronizer routine failed %s", err.Error())
	}

	blog.InfoJSON("process started with info %+v", httpServerConfig)

	<-ctx.Done()
	blog.V(0).Info("process stoped")
	return nil
}

// NewHTTPServerConfig new ServerInfo for running a http service
func NewHTTPServerConfig(serverOptions *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := serverOptions.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := serverOptions.ServConf.GetPort()
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
