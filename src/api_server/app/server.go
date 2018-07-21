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

	"github.com/emicklei/go-restful"

	"configcenter/src/api_server/app/options"
	"configcenter/src/api_server/ccapi/logics/v2"
	apisvc "configcenter/src/api_server/service"
	"configcenter/src/api_server/service/v3"
	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
)

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

	v2Service := new(apisvc.Service)
	v3Service := new(v3.Service)
	v3Service.Client, err = util.NewClient(&util.TLSClientConfig{})
	if err != nil {
		return fmt.Errorf("new proxy client failed, err: %v", err)
	}

	v3Service.Disc, err = discovery.NewDiscoveryInterface(op.ServConf.RegDiscover, "")
	if err != nil {
		return fmt.Errorf("new proxy discovery instance failed, err: %v", err)
	}

	ctnr := restful.NewContainer()
	ctnr.Router(restful.CurlyRouter{})
	ctnr.Add(v2Service.V2WebService())
	ctnr.Add(v3Service.V3WebService())
	ctnr.Add(v3Service.V3Healthz())
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

	apiSvr := new(APIServer)
	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_APISERVER,
		op.ServConf.ExConfig,
		apiSvr.onHostConfigUpdate,
		bonC)

	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	v2Service.Engine = engine
	v2Service.Logics = &logics.Logics{Engine: engine}
	v3Service.Engine = engine
	apiSvr.Core = engine
	apiSvr.Service = v2Service
	apiSvr.Logic = v2Service.Logics
	select {}
	return nil
}

type APIServer struct {
	Core    *backbone.Engine
	Service *apisvc.Service
	Logic   *logics.Logics
}

func (h *APIServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {

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
