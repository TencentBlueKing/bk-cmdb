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

	restful "github.com/emicklei/go-restful"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	toposvr "configcenter/src/scene_server/topo_server/service"
)

// TopoServer the topo server
type TopoServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service toposvr.TopoServiceInterface
}

func (t *TopoServer) onTopoConfigUpdate(previous, current cc.ProcessConfig) {

	cfg, err := mapstr.NewFromInterface(current.ConfigMap)
	if nil != err {
		blog.Errorf("failed to update config, error info is %s", err.Error())
	}

	if err := cfg.MarshalJSONInto(&t.Config); nil != err {
		blog.Errorf("failed to update config, error info is %s", err.Error())
	}
	blog.V(3).Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)
	t.Service.SetConfig(t.Config, t.Core)
}

// Run main function
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	blog.V(3).Infof("srv conf:", svrInfo)

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
	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_TOPO, svrInfo.IP)

	topoSvr := new(TopoServer)
	topoSvr.Config.BusinessTopoLevelMax = common.BKTopoBusinessLevelDefault

	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	topoService := toposvr.New()
	topoSvr.Service = topoService

	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    restful.NewContainer().Add(topoService.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	engine, err := backbone.NewBackbone(
		ctx,
		op.ServConf.RegDiscover,
		types.CC_MODULE_TOPO,
		op.ServConf.ExConfig,
		topoSvr.onTopoConfigUpdate,
		bonC)

	if nil != err {
		return fmt.Errorf("new engine failed, error is %s", err.Error())
	}

	topoSvr.Core = engine

	topoService.SetOperation(core.New(engine.CoreAPI), engine.CCErr, engine.Language)
	topoService.SetConfig(topoSvr.Config, engine)

	select {
	case <-ctx.Done():
	}
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
