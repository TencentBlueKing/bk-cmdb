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

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/source_controller/hostcontroller/app/options"
	"configcenter/src/source_controller/hostcontroller/logics"
	"configcenter/src/source_controller/hostcontroller/service"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/mongo/remote"
	dalredis "configcenter/src/storage/dal/redis"
)

//Run ccapi server
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	coreService := new(service.Service)
	hostCtrl := new(HostController)
	hostCtrl.Service = coreService
	coreService.Logics = &logics.Logics{Instance: nil}

	input := &backbone.BackboneParameter{
		ConfigUpdate: hostCtrl.onHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	hostCtrl.Core = engine

	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == hostCtrl.Config {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}

	coreService.Logics.Engine = coreService.Core
	err = hostCtrl.initComponents()
	if err != nil {
		blog.Errorf("initializing dependent components error. err:%s", err.Error())
		return err
	}

	if err := backbone.StartServer(ctx, coreService.Core, coreService.WebService()); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
		break
	}
	return nil
}

type HostController struct {
	*service.Service
	Config *options.Config
}

func (h *HostController) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	h.Config = &options.Config{
		Mongo: mongo.ParseConfigFromKV("mongodb", current.ConfigMap),
		Redis: dalredis.ParseConfigFromKV("redis", current.ConfigMap),
	}
}

// initCompents Initialize dependent components
func (h *HostController) initComponents() error {
	var err error
	var instance dal.RDB
	if h.Config.Mongo.Enable == "true" {
		instance, err = local.NewMgo(h.Config.Mongo.BuildURI(), time.Minute)
	} else {
		instance, err = remote.NewWithDiscover(h.Core.ServiceManageInterface.TMServer().GetServers)
	}
	if err != nil {
		blog.Errorf("new mongo client failed, err: %v", err)
		return err
	}

	cache, err := dalredis.NewFromConfig(h.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %v", err)
		return err
	}
	ec := eventclient.NewClientViaRedis(cache, instance)

	h.Service.Instance = instance
	h.Service.Logics.Instance = instance
	h.Service.Logics.Cache = cache
	h.Service.Logics.EventC = ec

	h.Cache = cache
	h.Service.Cache = cache
	h.Service.EventC = ec
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
