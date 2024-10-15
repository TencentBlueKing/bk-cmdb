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

// Package app runs transfer service
package app

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/source_controller/transfer-service/app/options"
	"configcenter/src/source_controller/transfer-service/service"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"

	"github.com/spf13/viper"
)

// TransferService is the cmdb sync service that transfers data from one cmdb to another
type TransferService struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *service.Service
}

func (s *TransferService) onConfigUpdate(previous, current cc.ProcessConfig) {
	if s.Config == nil {
		s.Config = new(options.Config)
	}

	blog.V(4).Infof("the new conf: %#v, the origin conf: %s", s.Config, string(current.ConfigData))
}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	svr := new(TransferService)

	input := &backbone.BackboneParameter{
		ConfigUpdate: svr.onConfigUpdate,
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
		if svr.Config == nil {
			time.Sleep(time.Second)
			continue
		}

		configReady = true
		break

	}

	if !configReady {
		return fmt.Errorf("configuration item not found")
	}

	svr.Core = engine

	if err = svr.initResource(op.ExSyncConfFile); err != nil {
		return err
	}

	svr.Service, err = service.New(svr.Config, engine)
	if err != nil {
		return err
	}

	err = backbone.StartServer(ctx, cancel, engine, svr.Service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	}
	return nil
}

func (s *TransferService) initResource(exSyncConfFile string) error {
	var err error

	s.Config.Mongo, err = s.Core.WithMongo()
	if err != nil {
		return fmt.Errorf("get mongo config failed, err: %v", err)
	}

	if err = mongodb.InitClient("", &s.Config.Mongo); err != nil {
		return fmt.Errorf("init mongo client failed, err: %v", err)
	}

	s.Config.WatchMongo, err = s.Core.WithMongo("watch")
	if err != nil {
		return fmt.Errorf("get watch mongo config failed, err: %v", err)
	}

	if err = mongodb.InitClient("watch", &s.Config.WatchMongo); err != nil {
		return fmt.Errorf("init watch mongo client failed, err: %v", err)
	}

	s.Config.Redis, err = s.Core.WithRedis()
	if err != nil {
		return fmt.Errorf("get redis config failed, err: %v", err)
	}

	err = redis.InitClient("redis", &s.Config.Redis)
	if err != nil {
		return fmt.Errorf("init redis client failed, err: %v", err)
	}

	s.Config.Sync = new(options.SyncConfig)
	if err = cc.UnmarshalKey("transferService", s.Config.Sync); err != nil {
		return fmt.Errorf("parse sync config failed, err: %v", err)
	}

	if !s.Config.Sync.EnableSync {
		return nil
	}

	if err = s.Config.Sync.Validate(); err != nil {
		return fmt.Errorf("sync config is invalid, err: %v", err)
	}

	if s.Config.Sync.Role != options.SyncRoleDest {
		return nil
	}

	if exSyncConfFile == "" {
		return fmt.Errorf("extra sync config file path is not set")
	}

	filePath := strings.TrimSuffix(exSyncConfFile, ".yaml")
	parser := viper.New()
	parser.SetConfigName(path.Base(filePath))
	parser.AddConfigPath(path.Dir(filePath))
	if err = parser.ReadInConfig(); err != nil {
		return fmt.Errorf("read dest extra sync config failed, err: %v", err)
	}

	s.Config.DestExConf = new(options.DestExSyncConf)
	if err = parser.Unmarshal(s.Config.DestExConf); err != nil {
		return fmt.Errorf("parse dest extra sync config failed, err: %v", err)
	}

	if err = s.Config.DestExConf.Validate(); err != nil {
		return fmt.Errorf("dest extra sync config is invalid, err: %v", err)
	}

	return nil
}
