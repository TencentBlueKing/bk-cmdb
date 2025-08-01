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
	ccerr "configcenter/src/common/errors"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/sync_server/app/options"
	"configcenter/src/scene_server/sync_server/logics"
	fulltextsearch "configcenter/src/scene_server/sync_server/logics/full-text-search"
	"configcenter/src/scene_server/sync_server/service"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
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
		SrvRegdiscv: backbone.SrvRegdiscv{Regdiscv: op.ServConf.RegDiscover,
			TLSConfig: op.ServConf.GetTLSClientConf()},
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: server.onConfigUpdate,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	server.Core = engine
	ccerr.SetGlobalCCError(engine.CCErr)

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

	if err = initClient(engine); err != nil {
		return err
	}

	// init sync server logics, then start web service
	server.Logics, err = logics.New(engine, server.Config)
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
		server.Logics.Stop()
	}

	return nil
}

func initClient(engine *backbone.Engine) error {
	// init mongo and redis client
	mongoConf, err := engine.WithMongo()
	if err != nil {
		return err
	}

	cryptoConf, err := cc.Crypto("crypto")
	if err != nil {
		blog.Errorf("get crypto config failed, err: %v", err)
		return err
	}

	if err = mongodb.SetShardingCli("", &mongoConf, cryptoConf); err != nil {
		blog.Errorf("init mongo client failed, err: %v, conf: %+v", err, mongoConf)
		return err
	}

	watchMongoConf, err := engine.WithMongo("watch")
	if err != nil {
		blog.Errorf("new watch mongo client failed, err: %v", err)
		return err
	}

	if err = mongodb.SetWatchCli("watch", &watchMongoConf, cryptoConf); err != nil {
		blog.Errorf("init watch mongo client failed, err: %v, conf: %+v", err, watchMongoConf)
		return err
	}

	redisConf, err := engine.WithRedis()
	if err != nil {
		return err
	}

	if err = redis.InitClient("redis", &redisConf); err != nil {
		blog.Errorf("init redis client failed, err: %v, conf: %+v", err, redisConf)
		return err
	}

	return nil
}

func (s *SyncServer) onConfigUpdate(previous, current cc.ProcessConfig) {
	blog.Infof("config updated, new config: %s", string(current.ConfigData))

	config := &logics.Config{
		FullTextSearch: new(fulltextsearch.Config),
	}
	err := cc.UnmarshalKey("syncServer", config)
	if err != nil {
		blog.Errorf("parse syncServer config failed, err: %v", err)
		return
	}

	config.FullTextSearch.Es, err = elasticsearch.ParseConfig("es")
	if err != nil {
		blog.Errorf("parse es config failed, err: %v", err)
		return
	}

	s.Config = config
}
