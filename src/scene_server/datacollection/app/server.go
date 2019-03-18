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
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/datacollection/app/options"
	"configcenter/src/scene_server/datacollection/datacollection"
	"configcenter/src/scene_server/datacollection/logics"
	svc "configcenter/src/scene_server/datacollection/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(svc.Service)
	process := new(DCServer)

	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
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

		instance, err := local.NewMgo(process.Config.MongoDB.BuildURI(), time.Minute)
		if err != nil {
			return fmt.Errorf("new mongo client failed, err: %s", err.Error())
		}

		esbChan := make(chan esbutil.EsbConfig, 1)
		esbChan <- process.Config.Esb
		esb, err := esbserver.NewEsb(engine.ApiMachineryConfig(), esbChan)
		if err != nil {
			return fmt.Errorf("new esb client failed, err: %s", err.Error())
		}

		process.Service.Logics = logics.NewLogics(ctx, service.Engine, instance, esb)

		err = datacollection.NewDataCollection(ctx, process.Config, process.Core).Run()
		if err != nil {
			return fmt.Errorf("run datacollection routine failed %s", err.Error())
		}
		break
	}

	blog.InfoJSON("process started with info %s", svrInfo)
	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(service.WebService())); err != nil {
		return err
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

		out, _ := json.MarshalIndent(current.ConfigMap, "", "  ") //ignore err, cause ConfigMap is map[string]string
		blog.V(3).Infof("config updated: \n%s", out)

		dbprefix := "mongodb"
		mongoConf := mongo.ParseConfigFromKV(dbprefix, current.ConfigMap)
		h.Config.MongoDB = mongoConf

		ccredisPrefix := "redis"
		redisConf := redis.ParseConfigFromKV(ccredisPrefix, current.ConfigMap)
		h.Config.CCRedis = redisConf

		snapPrefix := "snap-redis"
		snapredisConf := redis.ParseConfigFromKV(snapPrefix, current.ConfigMap)
		h.Config.SnapRedis.Config = snapredisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[snapPrefix+".enable"]

		discoverPrefix := "discover-redis"
		discoverRedisConf := redis.ParseConfigFromKV(discoverPrefix, current.ConfigMap)
		h.Config.DiscoverRedis.Config = discoverRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[discoverPrefix+".enable"]

		netcollectPrefix := "netcollect-redis"
		netcollectRedisConf := redis.ParseConfigFromKV(netcollectPrefix, current.ConfigMap)
		h.Config.NetcollectRedis.Config = netcollectRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[netcollectPrefix+".enable"]

		esbPrefix := "esb"
		h.Config.Esb.Addrs = current.ConfigMap[esbPrefix+".addr"]
		h.Config.Esb.AppCode = current.ConfigMap[esbPrefix+".appCode"]
		h.Config.Esb.AppSecret = current.ConfigMap[esbPrefix+".appSecret"]
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
