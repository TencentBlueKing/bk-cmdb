/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under,
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package app

import (
	"context"
	"fmt"
	"time"

	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/container_server/app/options"
	"configcenter/src/scene_server/container_server/core"
	cntsvr "configcenter/src/scene_server/container_server/service"
	"configcenter/src/storage/dal/redis"
)

const waitSeconds = 180

// ContainerServer the container server which manages container data
type ContainerServer struct {
	Core        *backbone.Engine
	Config      options.Config
	Service     cntsvr.ContainerServiceInterface
	configReady bool
}

func (t *ContainerServer) onConfigUpdate(previous, current cc.ProcessConfig) {
	t.configReady = true
	// TODO:
	blog.Warnf("config change previous %#v, current %#v", previous, current)
}

func (t *ContainerServer) checkForReadiness() error {
	for i := 1; i < waitSeconds; i++ {
		if !t.configReady {
			time.Sleep(time.Second)
			continue
		}
		blog.Info("container server config ready.")
		return nil
	}
	return fmt.Errorf("wait for container server config timeout after %ds", waitSeconds)
}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %+v", err)
	}

	blog.Infof("srv conf: %+v", svrInfo)

	containerSvr := new(ContainerServer)

	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: containerSvr.onConfigUpdate,
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	if err := containerSvr.checkForReadiness(); err != nil {
		return err
	}

	containerSvr.Config.Redis, err = engine.WithRedis()
	if err != nil {
		return err
	}

	cacheDB, err := redis.NewFromConfig(containerSvr.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err %s", err.Error())
		return fmt.Errorf("new redis client failed, err %s", err.Error())
	}

	coreIf := core.New(engine.CoreAPI, engine.Language, cacheDB)
	svc := cntsvr.New()
	svc.SetConfig(containerSvr.Config, engine, coreIf, engine.CCErr, engine.Language)
	containerSvr.Service = svc
	containerSvr.Core = engine

	err = backbone.StartServer(ctx, cancel, engine, containerSvr.Service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("context cancelled")
	}

	return nil
}
