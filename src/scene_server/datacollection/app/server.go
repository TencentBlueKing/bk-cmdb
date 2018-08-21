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
	"sync"
	"time"

	"github.com/emicklei/go-restful"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/datacollection/app/options"
	"configcenter/src/scene_server/datacollection/datacollection"
	svc "configcenter/src/scene_server/datacollection/service"
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

	service := new(svc.Service)
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    restful.NewContainer().Add(service.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_DATACOLLECTION, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	process := new(DCServer)
	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_DATACOLLECTION,
		op.ServConf.ExConfig,
		process.onHostConfigUpdate,
		bonC)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	service.Engine = engine
	process.Core = engine
	process.Service = service
	for {
		if process.Config == nil {
			time.Sleep(time.Second * 2)
			blog.Info("config not found, retry 2s later")
			continue
		}

		err := datacollection.NewDataCollection(process.Config, process.Core).Run()
		if err != nil {
			return fmt.Errorf("run datacollection routine failed %s", err.Error())
		}
		break
	}

	<-ctx.Done()
	blog.V(0).Info("process stoped")
	return nil
}

type DCServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *svc.Service
}

var configLock sync.Mutex

func (h *DCServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigMap) > 0 {
		if h.Config == nil {
			h.Config = new(options.Config)
		}
		dbprefix := "mongodb"
		h.Config.MongoDB.Address = current.ConfigMap[dbprefix+".host"]
		h.Config.MongoDB.User = current.ConfigMap[dbprefix+".usr"]
		h.Config.MongoDB.Password = current.ConfigMap[dbprefix+".pwd"]
		h.Config.MongoDB.Database = current.ConfigMap[dbprefix+".database"]
		h.Config.MongoDB.Port = current.ConfigMap[dbprefix+".port"]
		h.Config.MongoDB.MaxOpenConns = current.ConfigMap[dbprefix+".maxOpenConns"]
		h.Config.MongoDB.MaxIdleConns = current.ConfigMap[dbprefix+".maxIDleConns"]

		ccredisPrefix := "redis"
		h.Config.CCRedis.Address = current.ConfigMap[ccredisPrefix+".host"]
		h.Config.CCRedis.Password = current.ConfigMap[ccredisPrefix+".pwd"]
		h.Config.CCRedis.Database = current.ConfigMap[ccredisPrefix+".database"]
		h.Config.CCRedis.Port = current.ConfigMap[ccredisPrefix+".port"]
		h.Config.CCRedis.MasterName = current.ConfigMap[ccredisPrefix+".mastername"]

		snapPrefix := "snap-redis"
		h.Config.SnapRedis.Address = current.ConfigMap[snapPrefix+".host"]
		h.Config.SnapRedis.Password = current.ConfigMap[snapPrefix+".pwd"]
		h.Config.SnapRedis.Database = current.ConfigMap[snapPrefix+".database"]
		h.Config.SnapRedis.Port = current.ConfigMap[snapPrefix+".port"]
		h.Config.SnapRedis.MasterName = current.ConfigMap[snapPrefix+".mastername"]

		discoverPrefix := "discover-redis"
		h.Config.DiscoverRedis.Address = current.ConfigMap[discoverPrefix+".host"]
		h.Config.DiscoverRedis.Password = current.ConfigMap[discoverPrefix+".pwd"]
		h.Config.DiscoverRedis.Database = current.ConfigMap[discoverPrefix+".database"]
		h.Config.DiscoverRedis.Port = current.ConfigMap[discoverPrefix+".port"]
		h.Config.DiscoverRedis.MasterName = current.ConfigMap[discoverPrefix+".mastername"]
	}
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
