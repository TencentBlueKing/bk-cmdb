/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package app starts sync server
package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/sync_server/app/options"
	"configcenter/src/scene_server/sync_server/logics"
	fulltextsearch "configcenter/src/scene_server/sync_server/logics/full-text-search"
	"configcenter/src/scene_server/sync_server/service"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/stream"
	"configcenter/src/thirdparty/elasticsearch"
)

// SyncServer is the sync server
type SyncServer struct {
	Core    *backbone.Engine
	Config  *logics.Config
	Service *service.Service
	Logics  *logics.Logics
}

// Run sync server
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	server := new(SyncServer)

	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: server.onConfigUpdate,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	server.Core = engine

	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if server.Config != nil {
			configReady = true
			break
		}
		blog.Infof("waiting for config ready ...")
		time.Sleep(time.Second)
	}
	if !configReady {
		blog.Infof("waiting config timeout.")
		return errors.New("configuration item not found")
	}

	watcher, err := initClient(engine)
	if err != nil {
		return err
	}

	// init sync server logics, then start web service
	server.Logics, err = logics.New(engine, server.Config, watcher)
	if err != nil {
		return fmt.Errorf("new logics failed, err: %v", err)
	}

	server.Service = service.New(engine, server.Logics)

	err = backbone.StartServer(ctx, cancel, engine, server.Service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	}

	return nil
}

func initClient(engine *backbone.Engine) (stream.LoopInterface, error) {
	// init mongo and redis client
	mongoConf, err := engine.WithMongo()
	if err != nil {
		return nil, err
	}

	if err = mongodb.InitClient("", &mongoConf); err != nil {
		blog.Errorf("init mongo client failed, err: %v, conf: %+v", err, mongoConf)
		return nil, err
	}

	watchMongoConf, dbErr := engine.WithMongo("watch")
	if dbErr != nil {
		blog.Errorf("new watch mongo client failed, err: %v", dbErr)
		return nil, dbErr
	}

	if err = mongodb.InitClient("watch", &watchMongoConf); err != nil {
		blog.Errorf("init watch mongo client failed, err: %v, conf: %+v", err, watchMongoConf)
		return nil, err
	}

	redisConf, err := engine.WithRedis()
	if err != nil {
		return nil, err
	}

	if err = redis.InitClient("redis", &redisConf); err != nil {
		blog.Errorf("init redis client failed, err: %v, conf: %+v", err, redisConf)
		return nil, err
	}

	watcher, err := stream.NewLoopStream(mongoConf.GetMongoConf(), engine.ServiceManageInterface)
	if err != nil {
		blog.Errorf("new loop watch stream failed, err: %v", err)
		return nil, err
	}
	return watcher, nil
}

func (s *SyncServer) onConfigUpdate(previous, current cc.ProcessConfig) {
	s.Config = new(logics.Config)
	s.Config.FullTextSearch = new(fulltextsearch.Config)
	blog.Infof("config updated, new config: %s", string(current.ConfigData))

	err := cc.UnmarshalKey("syncServer", s.Config)
	if err != nil {
		return
	}

	s.Config.FullTextSearch.Es, err = elasticsearch.ParseConfig("es")
	if err != nil {
		blog.Warnf("parse es config failed: %v", err)
	}
}
