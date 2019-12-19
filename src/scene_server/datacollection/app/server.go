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

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	enableauth "configcenter/src/common/auth"
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

	re "gopkg.in/redis.v5"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
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

		mgoCli, err := local.NewMgo(process.Config.MongoDB.BuildURI(), time.Minute)
		if err != nil {
			return fmt.Errorf("new mongo client failed, err: %s", err.Error())
		}

		esbChan := make(chan esbutil.EsbConfig, 1)
		esbChan <- process.Config.Esb
		esb, err := esbserver.NewEsb(engine.ApiMachineryConfig(), esbChan, nil, engine.Metric().Registry())
		if err != nil {
			return fmt.Errorf("new esb client failed, err: %s", err.Error())
		}

		process.Service.SetDB(mgoCli)
		process.Service.Logics = logics.NewLogics(ctx, service.Engine, mgoCli, esb)
		datacollection := datacollection.NewDataCollection(ctx, process.Core, mgoCli, engine.Metric().Registry())

		blog.Infof("[data-collection][RUN]connecting to cc redis %+v", process.Config.CCRedis)
		redisCli, err := redis.NewFromConfig(process.Config.CCRedis)
		if nil != err {
			blog.Errorf("[data-collection][RUN] connect cc redis failed: %v", err)
			return err
		}
		blog.Infof("[data-collection][RUN]connected to cc redis %+v", process.Config.CCRedis)
		process.Service.SetCache(redisCli)

		var snapcli, disCli, netCli *re.Client
		if process.Config.SnapRedis.Enable != "false" {
			blog.Infof("[data-collection][RUN]connecting to snap-redis %+v", process.Config.SnapRedis.Config)
			snapcli, err = redis.NewFromConfig(process.Config.SnapRedis.Config)
			if nil != err {
				blog.Errorf("[data-collection][RUN] connect snap-redis failed: %v", err)
				return err
			}
			process.Service.SetSnapcli(snapcli)
		}
		if process.Config.DiscoverRedis.Enable != "false" {
			blog.Infof("[data-collection][RUN]connecting to discover-redis %+v", process.Config.DiscoverRedis.Config)
			disCli, err = redis.NewFromConfig(process.Config.DiscoverRedis.Config)
			if nil != err {
				blog.Errorf("[data-collection][RUN] connect discover-redis failed: %v", err)
				return err
			}
			blog.Infof("[data-collection][RUN]connected to discover-redis %+v", process.Config.DiscoverRedis.Config)
			process.Service.SetDisCli(disCli)
		}
		if process.Config.NetCollectRedis.Enable != "false" {
			blog.Infof("[data-collection][RUN]connecting to netcollect-redis %+v", process.Config.NetCollectRedis.Config)
			netCli, err = redis.NewFromConfig(process.Config.NetCollectRedis.Config)
			if nil != err {
				blog.Errorf("[data-collection][RUN] connect netcollect-redis failed: %v", err)
				return err
			}
			blog.Infof("[data-collection][RUN]connected to netcollect-redis %+v", process.Config.NetCollectRedis.Config)
			process.Service.SetNetCli(netCli)
		}
		if enableauth.IsAuthed() {
			blog.Info("[data-collection] auth enabled")
			authorize, err := auth.NewAuthorize(nil, process.Config.AuthConfig, engine.Metric().Registry())
			if err != nil {
				return fmt.Errorf("[data-collection] new authorize failed, err: %v", err)
			}
			datacollection.AuthManager = *extensions.NewAuthManager(engine.CoreAPI, authorize)
		}

		err = datacollection.Run(redisCli, snapcli, disCli, netCli)
		if err != nil {
			return fmt.Errorf("run datacollection routine failed %s", err.Error())
		}
		break
	}

	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	}
	blog.V(0).Info("process stopped")
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
		// ignore err, cause ConfigMap is map[string]string
		out, _ := json.MarshalIndent(current.ConfigMap, "", "  ")
		blog.V(3).Infof("config updated: \n%s", out)

		dbPrefix := "mongodb"
		mongoConf := mongo.ParseConfigFromKV(dbPrefix, current.ConfigMap)
		h.Config.MongoDB = mongoConf

		ccRedisPrefix := "redis"
		redisConf := redis.ParseConfigFromKV(ccRedisPrefix, current.ConfigMap)
		h.Config.CCRedis = redisConf

		snapPrefix := "snap-redis"
		snapRedisConf := redis.ParseConfigFromKV(snapPrefix, current.ConfigMap)
		h.Config.SnapRedis.Config = snapRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[snapPrefix+".enable"]

		discoverPrefix := "discover-redis"
		discoverRedisConf := redis.ParseConfigFromKV(discoverPrefix, current.ConfigMap)
		h.Config.DiscoverRedis.Config = discoverRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[discoverPrefix+".enable"]

		netCollectPrefix := "netcollect-redis"
		netCollectRedisConf := redis.ParseConfigFromKV(netCollectPrefix, current.ConfigMap)
		h.Config.NetCollectRedis.Config = netCollectRedisConf
		h.Config.SnapRedis.Enable = current.ConfigMap[netCollectPrefix+".enable"]

		esbPrefix := "esb"
		h.Config.Esb.Addrs = current.ConfigMap[esbPrefix+".addr"]
		h.Config.Esb.AppCode = current.ConfigMap[esbPrefix+".appCode"]
		h.Config.Esb.AppSecret = current.ConfigMap[esbPrefix+".appSecret"]

		var err error
		authPrefix := "auth"
		h.Config.AuthConfig, err = authcenter.ParseConfigFromKV(authPrefix, current.ConfigMap)
		if err != nil {
			blog.Fatalf("auth config invalid, err: %+v", err)
		}
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
