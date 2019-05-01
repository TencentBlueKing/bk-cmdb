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

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	//"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/synchronize_server/app/options"
	synchronizeService "configcenter/src/scene_server/synchronize_server/service"
	//"configcenter/src/storage/dal/redis"
	synchronizeClient "configcenter/src/apimachinery/synchronize"
	synchronizeUtil "configcenter/src/apimachinery/synchronize/util"
)

func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(synchronizeService.Service)
	synchronSrv := &SynchronizeServer{
		synchronizeClientConfig: make(chan synchronizeUtil.SychronizeConfig, 10),
	}
	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: synchronSrv.onSynchronizeServerConfigUpdate,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if synchronSrv.Config == nil {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}
	service.Engine = engine
	service.Config = synchronSrv.Config
	synchronSrv.Service = service
	synchronizeClientInst, err := synchronizeClient.NewSynchronize(engine.ApiMachineryConfig(), synchronSrv.synchronizeClientConfig)
	if err != nil {
		return fmt.Errorf("new NewSynchronize failed, err: %v", err)
	}
	service.SetSynchronizeServer(synchronizeClientInst)
	go synchronSrv.Service.InitBackground()
	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(service.WebService())); err != nil {
		return err
	}
	select {}
}

type SynchronizeServer struct {
	Core                    *backbone.Engine
	Config                  *options.Config
	Service                 *synchronizeService.Service
	synchronizeClientConfig chan synchronizeUtil.SychronizeConfig
}

func (s *SynchronizeServer) onSynchronizeServerConfigUpdate(previous, current cc.ProcessConfig) {
	configInfo := &options.Config{}
	names := current.ConfigMap["synchronize.name"]
	configInfo.Names = strings.Split(names, ",")

	configInfo.Trigger.TriggerType = current.ConfigMap["trigger.type"]
	// role  unit minute.
	// type = timing, ervery day  role minute trigger
	// type = interval, interval role  minute trigger
	configInfo.Trigger.Role = current.ConfigMap["trigger.role"]

	for _, name := range configInfo.Names {
		if strings.TrimSpace(name) == "" {
			continue
		}
		configItem := &options.ConfigItem{}
		appNames := current.ConfigMap[name+".AppNames"]
		syncResource := current.ConfigMap[name+".SynchronizeResource"]
		targetHost := current.ConfigMap[name+".Host"]
		fieldSign := current.ConfigMap[name+".FieldSign"]
		dataSign := current.ConfigMap[name+".Flag"]
		supplerAccount := current.ConfigMap[name+".SupplerAccount"]
		witeList := current.ConfigMap[name+".WiteList"]
		objectIDs := current.ConfigMap[name+".ObjectID"]

		configItem.AppNames = strings.Split(appNames, ",")
		if syncResource == "1" {
			configItem.SyncResource = true
		}
		if witeList == "1" {
			configItem.WiteList = true
		}
		configItem.ObjectIDArr = strings.Split(objectIDs, ",")
		configItem.Name = name
		configItem.TargetHost = targetHost
		configItem.FieldSign = fieldSign
		configItem.SynchronizeFlag = dataSign
		configItem.SupplerAccount = strings.Split(supplerAccount, ",")
		configInfo.ConifgItemArray = append(configInfo.ConifgItemArray, configItem)
		if targetHost != "" {
			s.synchronizeClientConfig <- synchronizeUtil.SychronizeConfig{
				Name:  name,
				Addrs: []string{targetHost},
			}
		}
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
