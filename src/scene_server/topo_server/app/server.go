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
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	toposvr "configcenter/src/scene_server/topo_server/service"
	"configcenter/src/storage/dal/mongo"
)

// TopoServer the topo server
type TopoServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service toposvr.TopoServiceInterface
}

func (t *TopoServer) onTopoConfigUpdate(previous, current cc.ProcessConfig) {
	topoMax := common.BKTopoBusinessLevelDefault
	var err error
	if current.ConfigMap["level.businessTopoMax"] != "" {
		topoMax, err = strconv.Atoi(current.ConfigMap["level.businessTopoMax"])
		if err != nil {
			blog.Errorf("invalid business topo max value, err: %v", err)
			return
		}
	}
	t.Config.BusinessTopoLevelMax = topoMax
	t.Config.Mongo = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)

	blog.V(3).Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)
	for t.Core == nil {
		time.Sleep(time.Second)
		blog.V(3).Info("sleep for engine")
	}
	t.Service.SetConfig(t.Config, t.Core)
}

// Run main function
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	blog.V(5).Infof("srv conf:", svrInfo)

	topoSvr := new(TopoServer)
	topoSvr.Config.BusinessTopoLevelMax = common.BKTopoBusinessLevelDefault

	topoService := toposvr.New()
	topoSvr.Service = topoService

	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: topoSvr.onTopoConfigUpdate,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	topoSvr.Core = engine

	topoService.SetOperation(core.New(engine.CoreAPI), engine.CCErr, engine.Language)
	// topoService.SetConfig(topoSvr.Config, engine)
	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(topoService.WebService())); err != nil {
		return err
	}
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
