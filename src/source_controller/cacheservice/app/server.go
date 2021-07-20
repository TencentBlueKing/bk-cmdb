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
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/source_controller/cacheservice/app/options"
	cachesvr "configcenter/src/source_controller/cacheservice/service"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// CoreServer the core server
type CacheServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service cachesvr.CacheServiceInterface
}

func (c *CacheServer) onCacheServiceConfigUpdate(previous, current cc.ProcessConfig) {
	if c.Config == nil {
		c.Config = new(options.Config)
	}

	blog.V(3).Infof("the new cfg:%#v the origin cfg:%#v", c.Config, string(current.ConfigData))

}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	cacheSvr := new(CacheServer)
	cacheService := cachesvr.New()
	cacheSvr.Service = cacheService

	input := &backbone.BackboneParameter{
		ConfigUpdate: cacheSvr.onCacheServiceConfigUpdate,
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
		if nil == cacheSvr.Config {
			time.Sleep(time.Second)
			continue
		}

		configReady = true
		break

	}

	if false == configReady {
		return fmt.Errorf("configuration item not found")
	}

	cacheSvr.Core = engine

	if err := initResource(cacheSvr); err != nil {
		return nil
	}

	err = cacheService.SetConfig(*cacheSvr.Config, engine, engine.CCErr, engine.Language)
	if err != nil {
		return err
	}

	err = backbone.StartServer(ctx, cancel, engine, cacheService.WebService(), true)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	}
	return nil
}

func initResource(cacheSvr *CacheServer) error {
	var err error
	cacheSvr.Config.Mongo, err = cacheSvr.Core.WithMongo()
	if err != nil {
		return err
	}

	cacheSvr.Config.WatchMongo, err = cacheSvr.Core.WithMongo("watch")
	if err != nil {
		return err
	}

	cacheSvr.Config.Redis, err = cacheSvr.Core.WithRedis()
	if err != nil {
		return err
	}

	dbErr := mongodb.InitClient("", &cacheSvr.Config.Mongo)
	if dbErr != nil {
		blog.Errorf("failed to connect the db server, error info is %s", dbErr.Error())
		return dbErr
	}

	cacheRrr := redis.InitClient("redis", &cacheSvr.Config.Redis)
	if cacheRrr != nil {
		blog.Errorf("new redis client failed, err: %v", cacheRrr)
		return cacheRrr
	}

	return nil
}
