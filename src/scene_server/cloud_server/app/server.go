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
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/cloud_server/app/options"
	"configcenter/src/scene_server/cloud_server/cloudsync"
	"configcenter/src/scene_server/cloud_server/logics"
	svc "configcenter/src/scene_server/cloud_server/service"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := svc.NewService(ctx)

	process := new(CloudServer)
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
	for {
		if process.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}

		redisConf, err := engine.WithRedis()
		if nil != err {
			blog.Errorf("get redis conf failed: %s", err.Error())
			return err
		}

		mongoConfig, err := engine.WithMongo()
		if nil != err {
			blog.Errorf("get mongo conf failed: %s", err.Error())
			return err
		}

		mongoConf := mongoConfig.GetMongoConf()
		db, err := local.NewMgo(mongoConf, time.Minute)
		if err != nil {
			return fmt.Errorf("connect mongo server failed, err: %s", err.Error())
		}
		process.Service.SetDB(db)

		cache, err := redis.NewFromConfig(redisConf)
		if err != nil {
			return fmt.Errorf("connect redis server failed, err: %s", err.Error())
		}
		process.Service.SetCache(cache)

		blog.Infof("enable auth center: %v", auth.IsAuthed())

		process.Service.Logics = logics.NewLogics(service.Engine, db, cache)

		syncConf := cloudsync.SyncConf{
			ZKClient:  service.Engine.ServiceManageClient().Client(),
			DB:        db,
			Logics:    process.Service.Logics,
			AddrPort:  input.SrvInfo.Instance(),
			MongoConf: mongoConf,
		}
		err = cloudsync.CloudSync(&syncConf)
		if err != nil {
			return fmt.Errorf("ProcessTask failed: %v", err)
		}

		break
	}
	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("process will exit!")
	}

	return nil
}

type CloudServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *svc.Service
}

var configLock sync.Mutex

func (c *CloudServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigMap) > 0 {
		if c.Config == nil {
			c.Config = new(options.Config)
		}
		// ignore err, cause ConfigMap is map[string]string
		out, _ := json.MarshalIndent(current.ConfigMap, "", "  ")
		blog.Infof("config updated: \n%s", out)
	}
}
