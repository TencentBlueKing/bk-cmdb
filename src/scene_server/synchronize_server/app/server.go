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

	"github.com/emicklei/go-restful"

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/synchronize_server/app/options"
	//hostsvc "configcenter/src/scene_server/synchronize_server/service"
	"configcenter/src/storage/dal/redis"
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

type SynchronizeServer struct {
	Core   *backbone.Engine
	Config *options.Config
	//Service *hostsvc.Service
}

func (s *SynchronizeServer) onSynchronizeServerConfigUpdate(previous, current cc.ProcessConfig) {
	configInfo := &options.Config{}
	names := current.ConfigMap["synchronze.name"]
	configInfo.Names = strings.Split(names, ",")

	for _, name := range configInfo.Names {
		configItem := options.ConfigItem{}
		ignoreAppNames := current.ConfigMap[name+".IgnoreAppNames"]
		syncResource := current.ConfigMap[name+".SynchronizeResource"]
		targetHost := current.ConfigMap[name+".TargetHost"]
		fieldSign := current.ConfigMap[name+".FieldSign"]
		dataSign := current.ConfigMap[name+".DataSign"]
		supplerAccount := current.ConfigMap[name+".SupplerAccount"]

		configItem.IgnoreAppNames = strings.Split(ignoreAppNames, ",")
		if syncResource == "1" {
			configItem.SyncResource = true
		}
		configItem.Name = name
		configItem.TargetHost = targetHost
		configItem.FieldSign = fieldSign
		configItem.DataSign = dataSign
		configItem.SupplerAccount = strings.Split(supplerAccount, ",")
		configInfo.ConifgItemArray = append(configInfo.ConifgItemArray, configItem)
	}
	s.Config = configInfo

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
