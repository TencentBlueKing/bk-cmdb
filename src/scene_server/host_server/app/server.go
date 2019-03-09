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

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/host_server/app/options"
	hostsvc "configcenter/src/scene_server/host_server/service"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful"
)

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

	service := new(hostsvc.Service)
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    restful.NewContainer().Add(service.WebService()),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_HOST, svrInfo.IP)
	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	hostSvr := new(HostServer)
	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_HOST,
		op.ServConf.ExConfig,
		hostSvr.onHostConfigUpdate,
		discover,
		bonC)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if "" == hostSvr.Config.Redis.Address {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}
	cacheDB, err := redis.NewFromConfig(hostSvr.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
	}

	service.Engine = engine
	service.Config = &hostSvr.Config
	service.CacheDB = cacheDB
	hostSvr.Core = engine
	hostSvr.Service = service

	go hostSvr.Service.InitBackground()
	select {}
	return nil
}

type HostServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service *hostsvc.Service
}

func (h *HostServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	h.Config.Gse.ZkAddress = current.ConfigMap["gse.addr"]
	h.Config.Gse.ZkUser = current.ConfigMap["gse.user"]
	h.Config.Gse.ZkPassword = current.ConfigMap["gse.pwd"]
	h.Config.Gse.RedisPort = current.ConfigMap["gse.port"]
	h.Config.Gse.RedisPassword = current.ConfigMap["gse.redis_pwd"]

	h.Config.Redis.Address = current.ConfigMap["redis.host"]
	h.Config.Redis.Database = current.ConfigMap["redis.database"]
	h.Config.Redis.Password = current.ConfigMap["redis.pwd"]
	h.Config.Redis.Port = current.ConfigMap["redis.port"]
	h.Config.Redis.MasterName = current.ConfigMap["redis.user"]
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
