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
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/source_controller/coreservice/app/options"
	coresvr "configcenter/src/source_controller/coreservice/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"
)

// CoreServer the core server
type CoreServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service coresvr.CoreServiceInterface
}

func (t *CoreServer) onCoreServiceConfigUpdate(previous, current cc.ProcessConfig) {

	t.Config.Mongo = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)
	t.Config.Redis = redis.ParseConfigFromKV("redis", current.ConfigMap)

	blog.V(3).Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)

}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	coreSvr := new(CoreServer)
	coreService := coresvr.New()
	coreSvr.Service = coreService

	input := &backbone.BackboneParameter{
		ConfigUpdate: coreSvr.onCoreServiceConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	var configReady bool
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		// redis not found
		if "" == coreSvr.Config.Redis.Address {
			time.Sleep(time.Second)
			continue
		}

		configReady = true
		break

	}

	if false == configReady {
		return fmt.Errorf("configuration item not found")
	}

	coreSvr.Core = engine
	err = coreService.SetConfig(coreSvr.Config, engine, engine.CCErr, engine.Language)
	if err != nil {
		return err
	}
	err = backbone.StartServer(ctx, cancel, engine, coreService.WebService(), true)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	}
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
