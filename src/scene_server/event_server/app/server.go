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
	"sync"
	"time"

	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/event_server/app/options"
	"configcenter/src/scene_server/event_server/distribution"
	svc "configcenter/src/scene_server/event_server/service"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := svc.NewService(ctx)

	process := new(EventServer)
	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	service.Engine = engine
	process.Core = engine
	process.Service = service
	errCh := make(chan error, 1)
	for {
		if process.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}

		process.Config.MongoDB, err = engine.WithMongo()
		if err != nil {
			return err
		}
		process.Config.Redis, err = engine.WithRedis()
		if err != nil {
			return err
		}

		db, err := local.NewMgo(process.Config.MongoDB.GetMongoConf(), time.Minute)
		if err != nil {
			return fmt.Errorf("connect mongo server failed, err: %s", err.Error())
		}
		process.Service.SetDB(db)

		cache, err := redis.NewFromConfig(process.Config.Redis)
		if err != nil {
			return fmt.Errorf("connect redis server failed, err: %s", err.Error())
		}
		process.Service.SetCache(cache)

		subCli, err := redis.NewFromConfig(process.Config.Redis)
		if err != nil {
			return fmt.Errorf("connect subcli redis server failed, err: %s", err.Error())
		}

		go func() {
			errCh <- distribution.SubscribeChannel(subCli)
		}()

		go func() {
			errCh <- distribution.Start(ctx, cache, db, engine.CoreAPI, engine.ServiceManageInterface)
		}()

		break
	}
	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	case err = <-errCh:
		blog.Errorf("distribution routine stopped, err: %v", err)
		return err
	}

	return nil
}

type EventServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *svc.Service
}

var configLock sync.Mutex

func (h *EventServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigData) > 0 {
		if h.Config == nil {
			h.Config = new(options.Config)
		}
		// ignore err, cause ConfigMap is map[string]string
		blog.Infof("config updated: \n%s", string(current.ConfigData))
	}
}
