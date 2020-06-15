/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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

	"configcenter/src/ac/iam"
	"configcenter/src/common/backbone"
	"configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/auth_server/app/options"
	"configcenter/src/scene_server/auth_server/service"
	"configcenter/src/storage/dal/redis"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap authServer info failed, err: %v", err)
	}

	authServer := new(AuthServer)

	input := &backbone.BackboneParameter{
		ConfigUpdate: authServer.onAuthConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	authServer.Core = engine
	for {
		if authServer.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}

		iamCli, err := iam.NewIam(nil, authServer.Config.Auth, engine.Metric().Registry())
		if err != nil {
			blog.Errorf("new iam client, err: %s", err.Error())
			return err
		}

		redisConf, err := engine.WithRedis()
		if nil != err {
			blog.Errorf("get redis conf failed: %s", err.Error())
			return err
		}
		redisCli, err := redis.NewFromConfig(redisConf)
		if nil != err {
			blog.Errorf("new redis client failed: %s", err.Error())
			return err
		}

		// TODO use unified cache
		//listDone := make(chan bool, 1)
		//errChan := make(chan error, 1)
		//mongoConf, err := engine.WithMongo()
		//if nil != err {
		//	blog.Errorf("get mongo conf failed: %s", err.Error())
		//	return err
		//}
		//go func() {
		//	cache.SyncDatabaseToRedis(ctx, mongoConf.GetMongoConf(), redisCli, listDone, errChan)
		//}()
		//select {
		//case err = <-errChan:
		//	if nil != err {
		//		blog.Errorf("sync database to redis failed: %s", err.Error())
		//		return err
		//	}
		//case <-listDone:
		//	blog.V(5).Info("cache current mongo database into redis done")
		//}

		authServer.Service = service.NewAuthService(engine, iamCli, redisCli)
		break
	}
	err = backbone.StartServer(ctx, cancel, engine, authServer.Service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("auth server will exit!")
	}

	return nil
}

type AuthServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *service.AuthService
}

var configLock sync.Mutex

func (a *AuthServer) onAuthConfigUpdate(previous, current configcenter.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigMap) > 0 {
		if a.Config == nil {
			a.Config = new(options.Config)
		}
		blog.InfoJSON("config updated: \n%s", current.ConfigMap)

		var err error
		a.Config.Auth, err = iam.ParseConfigFromKV("auth", current.ConfigMap)
		if err != nil {
			blog.Warnf("parse auth center config failed: %v", err)
		}
	}
}
