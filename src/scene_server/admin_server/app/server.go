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
	"time"

	"configcenter/src/auth/authcenter"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	"configcenter/src/common/backbone/configcenter"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/authsynchronizer"
	"configcenter/src/scene_server/admin_server/configures"
	svc "configcenter/src/scene_server/admin_server/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	process := new(MigrateServer)
	config, err := configcenter.ParseConfigWithFile(op.ServConf.ExConfig)
	if nil != err {
		return fmt.Errorf("parse config file error %s", err.Error())
	}

	process.Config = new(options.Config)
	// ignore err, because ConfigMap is map[string]string
	out, _ := json.MarshalIndent(config.ConfigMap, "", "  ")
	blog.V(3).Infof("config updated: \n%s", out)

	mongoConf := mongo.ParseConfigFromKV("mongodb", config.ConfigMap)
	process.Config.MongoDB = mongoConf

	redisConf := redis.ParseConfigFromKV("redis", config.ConfigMap)
	process.Config.Redis = redisConf

	process.Config.Errors.Res = config.ConfigMap["errors.res"]
	process.Config.Language.Res = config.ConfigMap["language.res"]
	process.Config.Configures.Dir = config.ConfigMap["confs.dir"]

	process.Config.Register.Address = config.ConfigMap["register-server.addrs"]

	process.Config.ProcSrvConfig.CCApiSrvAddr, _ = config.ConfigMap["procsrv.cc_api"]

	process.Config.AuthCenter, err = authcenter.ParseConfigFromKV("auth", config.ConfigMap)
	if err != nil && auth.IsAuthed() {
		blog.Errorf("parse authcenter error: %v, config: %+v", err, config.ConfigMap)
	}
	service := svc.NewService(ctx)

	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     process.Config.Register.Address,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	service.Engine = engine
	service.Config = *process.Config
	process.Core = engine
	process.Service = service
	process.ConfigCenter = configures.NewConfCenter(ctx, engine.ServiceManageClient())

	// adminserver conf not depend discovery
	err = process.ConfigCenter.Start(
		config.ConfigMap["confs.dir"],
		config.ConfigMap["errors.res"],
		config.ConfigMap["language.res"],
	)
	if err != nil {
		return err
	}

	for {
		if process.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}

		db, err := local.NewMgo(process.Config.MongoDB.GetMongoConf(), time.Minute)
		if err != nil {
			return fmt.Errorf("connect mongo server failed %s", err.Error())
		}
		process.Service.SetDB(db)
		cache, err := redis.NewFromConfig(process.Config.Redis)
		if err != nil {
			return fmt.Errorf("connect redis server failed, err: %s", err.Error())
		}
		process.Service.SetCache(cache)
		process.Service.SetApiSrvAddr(process.Config.ProcSrvConfig.CCApiSrvAddr)

		if auth.IsAuthed() {
			blog.Info("enable auth center access.")
			authCli, err := authcenter.NewAuthCenter(nil, process.Config.AuthCenter, engine.Metric().Registry())
			if err != nil {
				return fmt.Errorf("new authcenter client failed: %v", err)
			}
			process.Service.SetAuthCenter(authCli)

			if process.Config.AuthCenter.EnableSync {
				authSynchronizer := authsynchronizer.NewSynchronizer(ctx, &process.Config.AuthCenter, engine.CoreAPI, engine.Metric().Registry(), service.Engine)
				authSynchronizer.Run()
				blog.Info("enable auth center and enable auth sync function.")
			}

		} else {
			blog.Infof("disable auth center access.")
		}
		break
	}
	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	}
	blog.V(0).Info("process stopped")
	return nil
}

type MigrateServer struct {
	Core         *backbone.Engine
	Config       *options.Config
	Service      *svc.Service
	ConfigCenter *configures.ConfCenter
}

func (h *MigrateServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {}
