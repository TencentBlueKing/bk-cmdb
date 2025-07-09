/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/ac/extensions"
	"configcenter/src/ac/iam"
	"configcenter/src/common"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/host_server/app/options"
	"configcenter/src/scene_server/host_server/logics"
	hostsvc "configcenter/src/scene_server/host_server/service"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful/v3"
)

// Run TODO
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		blog.Errorf("wrap server info failed, err: %v", err)
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(hostsvc.Service)
	hostSrv := new(HostServer)

	input := &backbone.BackboneParameter{
		SrvRegdiscv:  backbone.SrvRegdiscv{Regdiscv: op.ServConf.RegDiscover, TLSConfig: op.ServConf.GetTLSClientConf()},
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

	hostSrv.Config.Auth, err = iam.ParseConfigFromKV("authServer", nil)
	if err != nil {
		blog.Warnf("parse auth center config failed: %v", err)
	}

	iamCli := new(iam.IAM)
	if auth.EnableAuthorize() {
		blog.Info("enable auth center access")
		iamCli, err = iam.NewIAM(hostSrv.Config.Auth, engine.Metric().Registry())
		if err != nil {
			return fmt.Errorf("new iam client failed: %v", err)
		}
	} else {
		blog.Infof("disable auth center access")
	}
	authManager := extensions.NewAuthManager(engine.CoreAPI, iamCli)

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

// HostServer TODO
type HostServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *hostsvc.Service
}

// WebService TODO
func (h *HostServer) WebService() *restful.Container {
	return h.Service.WebService()
}

func (h *HostServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	if h.Config == nil {
		h.Config = new(options.Config)
	}
}
