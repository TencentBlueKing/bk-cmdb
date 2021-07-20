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
	"errors"
	"fmt"
	"time"

	"configcenter/src/ac/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/host_server/app/options"
	"configcenter/src/scene_server/host_server/logics"
	hostsvc "configcenter/src/scene_server/host_server/service"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		blog.Errorf("wrap server info failed, err: %v", err)
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(hostsvc.Service)
	hostSrv := new(HostServer)

	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: hostSrv.onHostConfigUpdate,
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		blog.Errorf("new backbone failed, err: %v", err)
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil != hostSrv.Config {
			configReady = true
			break
		}
		blog.Infof("waiting for config ready ...")
		time.Sleep(time.Second)
	}
	if false == configReady {
		blog.Infof("waiting config timeout.")
		return errors.New("configuration item not found")
	}

	hostSrv.Config.Redis, err = engine.WithRedis()
	if err != nil {
		return err
	}

	cacheDB, err := redis.NewFromConfig(hostSrv.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
		return fmt.Errorf("new redis client failed, err: %s", err.Error())
	}

	authManager := extensions.NewAuthManager(engine.CoreAPI)
	service.AuthManager = authManager
	service.Engine = engine
	service.Config = hostSrv.Config
	service.CacheDB = cacheDB
	service.Logic = logics.NewLogics(engine, cacheDB, authManager)
	hostSrv.Core = engine
	hostSrv.Service = service

	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		blog.Errorf("start backbone failed, err: %+v", err)
		return err
	}

	select {
	case <-ctx.Done():
	}
	return nil
}

type HostServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *hostsvc.Service
}

func (h *HostServer) WebService() *restful.Container {
	return h.Service.WebService()
}

func (h *HostServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	if h.Config == nil {
		h.Config = new(options.Config)
	}
}
