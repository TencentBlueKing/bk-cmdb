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

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/eventclient"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/source_controller/proccontroller/app/options"
	"configcenter/src/source_controller/proccontroller/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	dalredis "configcenter/src/storage/dal/redis"
)

//Run ccapi server
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	coreService := new(service.ProctrlServer)
	procCtr := new(ProcController)
	procCtr.ProctrlServer = coreService

	input := &backbone.BackboneParameter{
		ConfigUpdate: procCtr.onProcConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == procCtr.Config {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}
	coreService.Core = engine
	procCtr.ProctrlServer.Core = engine
	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(coreService.WebService())); err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	}
	return nil
}

type ProcController struct {
	*service.ProctrlServer
	Config *options.Config
}

func (h *ProcController) onProcConfigUpdate(previous, current cc.ProcessConfig) {

	h.Config = &options.Config{
		Mongo: mongo.ParseConfigFromKV("mongodb", current.ConfigMap),
		Redis: dalredis.ParseConfigFromKV("redis", current.ConfigMap),
	}

	instance, err := local.NewMgo(h.Config.Mongo.BuildURI(), time.Minute)
	if err != nil {
		blog.Errorf("new mongo client failed, err: %v", err)
		return
	}
	h.ProctrlServer.Instance = instance

	cache, err := dalredis.NewFromConfig(h.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %v", err)
		return
	}
	h.ProctrlServer.Cache = cache

	ec := eventclient.NewClientViaRedis(cache, instance)
	h.ProctrlServer.EventC = ec
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
