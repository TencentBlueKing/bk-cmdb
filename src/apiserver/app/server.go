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
	"configcenter/src/auth"
	"configcenter/src/common/blog"
	"configcenter/src/framework/core/errors"
	"context"
	"fmt"
	"os"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/apiserver/app/options"
	"configcenter/src/apiserver/service"
	"configcenter/src/common/backbone"
	"configcenter/src/common/types"
	"configcenter/src/common/version"

	cc "configcenter/src/common/backbone/configcenter"
	"github.com/emicklei/go-restful"
)

// Run main loop function
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	discover, err := discovery.NewDiscoveryInterface(op.ServConf.RegDiscover)
	if err != nil {
		return fmt.Errorf("connect zookeeper [%s] failed: %v", op.ServConf.RegDiscover, err)
	}

	c := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}

	machinery, err := apimachinery.NewApiMachinery(c, discover)
	if err != nil {
		return fmt.Errorf("new api machinery failed, err: %v", err)
	}

	apisvr := service.NewService()

	httpClient, err := util.NewClient(&util.TLSClientConfig{})
	if err != nil {
		return fmt.Errorf("new proxy client failed, err: %v", err)
	}

	ctnr := restful.NewContainer()
	ctnr.Router(restful.CurlyRouter{})
	for _, item := range apisvr.WebServices() {
		ctnr.Add(item)
	}

	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    ctnr,
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_APISERVER, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	apiServer := new(APIServer)
	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_APISERVER,
		op.ServConf.ExConfig,
		apiServer.onApiServerConfigUpdate,
		discover,
		bonC)

	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	try := 300
	for i := 0; i <= try; i++ {
		if len(apiServer.Config) == 0 {
			blog.Info("waiting for the config items from zookeeper, will wait for another second.")
			time.Sleep(time.Second)
			continue
		} else {
			break
		}

		if i == try {
			return errors.New("wait for api server config items in zookeeper timeout")
		}
	}

	authorize, err := auth.NewAuthorize(nil, apiServer.Config)
	if err != nil {
		return err
	}
	apisvr.SetConfig(engine, httpClient, discover, authorize)
	apiServer.Core = engine

	select {
	case <-ctx.Done():
	}
	return nil
}

// APIServer apiserver struct
type APIServer struct {
	Core   *backbone.Engine
	Config map[string]string
}

func (h *APIServer) onApiServerConfigUpdate(previous, current cc.ProcessConfig) {
	h.Config = current.ConfigMap
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
