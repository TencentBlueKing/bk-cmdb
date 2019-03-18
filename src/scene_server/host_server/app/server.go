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
    "github.com/emicklei/go-restful"
    "os"
    "strconv"
    "time"
    "errors"

    "configcenter/src/common"
    "configcenter/src/common/backbone"
    cc "configcenter/src/common/backbone/configcenter"
    "configcenter/src/common/blog"
    "configcenter/src/common/types"
    "configcenter/src/common/version"
    "configcenter/src/scene_server/host_server/app/options"
    "configcenter/src/scene_server/host_server/authorize"
    hostsvc "configcenter/src/scene_server/host_server/service"
    "configcenter/src/storage/dal/redis"
)

func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
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
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if "" != hostSrv.Config.Redis.Address && hostSrv.Config.Auth.Address != ""{
            configReady = true
		    break
		}
        blog.Infof("waiting for config ready ...")
        time.Sleep(time.Second)
	}
	if false == configReady {
        blog.Infof("waiting config timeout.")
		return errors.New("Configuration item not found")
	}
	cacheDB, err := redis.NewFromConfig(hostSrv.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
	}

	config := hostSrv.Config.Auth
	blog.Info("host server auth config is: %+v", config)
    authorizer, err := authorize.NewHostAuthorizer(nil, config)
    if err != nil {
        blog.Errorf("make host authorizer failed, err: %+v", err)
        return fmt.Errorf("make host authorizer failed, err: %+v", err)
    }
    service.Authorizer = *authorizer
    
	service.Engine = engine
	service.Config = &hostSrv.Config
	service.CacheDB = cacheDB
	hostSrv.Core = engine
	hostSrv.Service = service

	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(service.WebService())); err != nil {
		return err
	}

	go hostSrv.Service.InitBackground()
	select {}
}

type HostServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service *hostsvc.Service
}

func (h *HostServer) WebService() *restful.WebService {
	return h.Service.WebService()
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
	
    h.Config.Auth.Address = current.ConfigMap["auth.address"]
    h.Config.Auth.AppCode = current.ConfigMap["auth.appCode"]
    h.Config.Auth.AppSecret = current.ConfigMap["auth.appSecret"]

    enable, err := strconv.ParseBool(current.ConfigMap["auth.enable"])
    if err == nil {
        blog.Errorf("parse auth enable field failed, set default to true, err: %+v", err)
        enable = true
    }
    h.Config.Auth.Enable = enable
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
