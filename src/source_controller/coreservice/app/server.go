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
	"configcenter/src/source_controller/coreservice/app/options"
	coresvr "configcenter/src/source_controller/coreservice/service"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
)

// CoreServer the core server
type CoreServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service coresvr.CoreServiceInterface
}

func (t *CoreServer) onCoreServiceConfigUpdate(previous, current cc.ProcessConfig) {
	if t.Config == nil {
		t.Config = new(options.Config)
	}

	blog.V(3).Infof("the new cfg:%#v the origin cfg:%#v", t.Config, string(current.ConfigData))

}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
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
		if nil == coreSvr.Config {
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

	if err := initResource(coreSvr); err != nil {
		return err
	}

	err = coreService.SetConfig(*coreSvr.Config, engine, engine.CCErr, engine.Language)
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

func initResource(coreSvr *CoreServer) error {
	var err error
	coreSvr.Config.Mongo, err = coreSvr.Core.WithMongo()
	if err != nil {
		return err
	}
	coreSvr.Config.Redis, err = coreSvr.Core.WithRedis()
	if err != nil {
		return err
	}

	dbErr := mongodb.InitClient("", &coreSvr.Config.Mongo)
	if dbErr != nil {
		blog.Errorf("failed to connect the db server, error info is %s", dbErr.Error())
		return dbErr
	}

	cacheRrr := redis.InitClient("redis", &coreSvr.Config.Redis)
	if cacheRrr != nil {
		blog.Errorf("new redis client failed, err: %v", cacheRrr)
		return cacheRrr
	}

	initErr := mongodb.Client().InitTxnManager(redis.Client())
	if initErr != nil {
		blog.Errorf("failed to init txn manager, error info is %v", initErr)
		return initErr
	}

	return nil
}
