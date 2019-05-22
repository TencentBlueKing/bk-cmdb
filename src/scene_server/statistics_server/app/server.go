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
	"os"
	"time"

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/statistics_server/app/options"
	svc "configcenter/src/scene_server/statistics_server/service"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful"
)

func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		blog.Errorf("wrap server info failed, err: %v", err)
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(svc.Service)
	statisticSrv := new(StatisticServer)

	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: statisticSrv.onStatisticalConfigUpdate,
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		blog.Errorf("new backbone failed, err: %v", err)
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if "" != statisticSrv.Config.Redis.Address {
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
	cacheDB, err := redis.NewFromConfig(statisticSrv.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
		return fmt.Errorf("new redis client failed, err: %s", err.Error())
	}

	authorizer, err := auth.NewAuthorize(nil, statisticSrv.Config.Auth)
	if err != nil {
		blog.Errorf("new host authorizer failed, err: %+v", err)
		return fmt.Errorf("new host authorizer failed, err: %+v", err)
	}
	authManager := extensions.NewAuthManager(engine.CoreAPI, authorizer)
	service.AuthManager = authManager
	service.Engine = engine
	service.CacheDB = cacheDB
	statisticSrv.Core = engine
	statisticSrv.Service = service

	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(service.WebService())); err != nil {
		blog.Errorf("start backbone failed, err: %+v", err)
		return err
	}

	select {}
}

type StatisticServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service *svc.Service
}

func (h *StatisticServer) onStatisticalConfigUpdate(previous, current cc.ProcessConfig) {
	var err error

	h.Config.Redis.Address = current.ConfigMap["redis.host"]
	h.Config.Redis.Database = current.ConfigMap["redis.database"]
	h.Config.Redis.Password = current.ConfigMap["redis.pwd"]
	h.Config.Redis.Port = current.ConfigMap["redis.port"]
	h.Config.Redis.MasterName = current.ConfigMap["redis.user"]

	h.Config.Auth, err = authcenter.ParseConfigFromKV("auth", current.ConfigMap)
	if err != nil {
		blog.Warnf("parse auth center config failed: %v", err)
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
