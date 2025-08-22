/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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

// Package app defines admin server logics
package app

import (
	"context"
	"fmt"
	"time"

	iamcli "configcenter/src/ac/iam"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/configures"
	"configcenter/src/scene_server/admin_server/iam"
	"configcenter/src/scene_server/admin_server/logics"
	svc "configcenter/src/scene_server/admin_server/service"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/thirdparty/monitor"
)

// Run start server
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	process, err := parseSeverConfig(ctx, op)
	if err != nil {
		return err
	}

	// adminserver conf not depend discovery
	err = process.ConfigCenter.Start(process.Config.Configures.Dir, process.Config.Errors.Res,
		process.Config.Language.Res)
	if err != nil {
		return err
	}

	service := svc.NewService(ctx)
	service.Engine = process.Core
	service.Config = *process.Config
	service.ConfigCenter = process.ConfigCenter
	process.Service = service

	if dbErr := mongodb.SetShardingCli("", &process.Config.MongoDB, process.Config.Crypto); dbErr != nil {
		return fmt.Errorf("connect mongo server failed %s", dbErr.Error())
	}
	db := mongodb.Dal()
	process.Service.SetDB(db)

	if err = mongodb.SetWatchCli("watch", &process.Config.WatchDB, process.Config.Crypto); err != nil {
		return fmt.Errorf("connect watch mongo server failed, err: %v", err)
	}
	process.Service.SetWatchDB(mongodb.Dal("watch"))

	// init old migrate db for old version migration before v3.15.1, remove this after v3.15.2
	oldMigrateDB, err := local.NewOldMgo(process.Config.MongoDB.GetMongoConf(), time.Minute)
	if err != nil {
		return fmt.Errorf("new mongodb client for previous version failed, err: %v", err)
	}
	process.Service.SetOldMigrateDB(oldMigrateDB)

	cache, err := redis.NewFromConfig(process.Config.Redis)
	if err != nil {
		return fmt.Errorf("connect redis server failed, err: %s", err.Error())
	}
	process.Service.SetCache(cache)

	process.Service.Logics = logics.NewLogics(process.Core)

	if err := service.InitClients(); err != nil {
		return err
	}

	var iamCli *iamcli.IAM
	if auth.EnableAuthorize() {
		iamCli, err = iamcli.NewIAM()
		if err != nil {
			return fmt.Errorf("new iam client failed: %v", err)
		}
		process.Service.SetIam(iamCli)
	} else {
		blog.Infof("disable auth center access.")
	}
	if err = service.InitCrypto(); err != nil {
		return err
	}

	if err = service.BackgroundTask(*process.Config); err != nil {
		return err
	}

	err = backbone.StartServer(ctx, cancel, process.Core, service.WebService(), true)
	if err != nil {
		return err
	}

	errors.SetGlobalCCError(process.Core.CCErr)

	syncor := iam.NewSyncor()
	syncor.SetDB(mongodb.Dal())
	syncor.SetSyncIAMPeriod(process.Config.SyncIAMPeriodMinutes)
	go syncor.SyncIAM(iamCli, cache, service.Logics)

	select {
	case <-ctx.Done():
	}
	blog.V(0).Info("process stopped")
	return nil
}

func parseSeverConfig(ctx context.Context, op *options.ServerOption) (*MigrateServer, error) {
	process := new(MigrateServer)
	process.Config = new(options.Config)
	if err := cc.SetLocalFile(op.ServConf.ExConfig); err != nil {
		return nil, fmt.Errorf("parse config file error %s", err.Error())
	}

	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return nil, fmt.Errorf("wrap server info failed, err: %v", err)
	}

	process.Config.Errors.Res, _ = cc.String("errors.res")
	process.Config.Language.Res, _ = cc.String("language.res")
	process.Config.Configures.Dir, _ = cc.String("confs.dir")
	process.Config.Register.Address, _ = cc.String("registerServer.addrs")
	process.Config.Register.TLS, _ = cc.NewTLSClientConfigFromConfig("registerServer.tls")
	process.Config.MigrateDataID, _ = cc.Bool("hostsnap.migrateDataID")
	process.Config.SyncIAMPeriodMinutes, _ = cc.Int("adminServer.syncIAMPeriodMinutes")
	process.Config.DisableVerifyTenant, _ = cc.Bool("tenant.disableVerifyTenant")
	process.Config.EnableMultiTenantMode, _ = cc.Bool("tenant.enableMultiTenantMode")

	// load mongodb, redis and common config from configure directory
	mongodbPath := process.Config.Configures.Dir + "/" + types.CCConfigureMongo
	if err := cc.SetMongodbFromFile(mongodbPath); err != nil {
		return nil, fmt.Errorf("parse mongodb config from file[%s] failed, err: %v", mongodbPath, err)
	}

	redisPath := process.Config.Configures.Dir + "/" + types.CCConfigureRedis
	if err := cc.SetRedisFromFile(redisPath); err != nil {
		return nil, fmt.Errorf("parse redis config from file[%s] failed, err: %v", redisPath, err)
	}

	commonPath := process.Config.Configures.Dir + "/" + types.CCConfigureCommon
	if err := cc.SetCommonFromFile(commonPath); err != nil {
		return nil, fmt.Errorf("parse common config from file[%s] failed, err: %v", commonPath, err)
	}

	process.Config.SnapReportMode, _ = cc.String("datacollection.hostsnap.reportMode")
	process.Config.SnapKafka, _ = cc.Kafka("kafka.snap")

	if err := monitor.InitMonitor(); err != nil {
		return nil, fmt.Errorf("init monitor failed, err: %v", err)
	}

	mongoConf, err := cc.Mongo("mongodb")
	if err != nil {
		return nil, err
	}
	process.Config.MongoDB = mongoConf

	watchDBConf, err := cc.Mongo("watch")
	if err != nil {
		return nil, err
	}
	process.Config.WatchDB = watchDBConf

	redisConf, err := cc.Redis("redis")
	if err != nil {
		return nil, err
	}
	process.Config.Redis = redisConf

	snapRedisConf, err := cc.Redis("redis.snap")
	if err != nil {
		return nil, fmt.Errorf("get host snapshot redis configuration failed, err: %v", err)
	}
	process.Config.SnapRedis = snapRedisConf

	process.Config.Crypto, err = cc.Crypto("crypto")
	if err != nil {
		return nil, fmt.Errorf("get crypto config failed, err: %v", err)
	}

	if err = parseShardingTableConfig(process); err != nil {
		return nil, err
	}

	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onMigrateConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		SrvRegdiscv: backbone.SrvRegdiscv{
			Regdiscv:  process.Config.Register.Address,
			TLSConfig: &process.Config.Register.TLS,
		},
		SrvInfo: svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("new backbone failed, err: %v", err)
	}

	process.Core = engine
	process.ConfigCenter = configures.NewConfCenter(ctx, engine.ServiceManageClient())
	return process, nil
}

// MigrateServer TODO
type MigrateServer struct {
	Core         *backbone.Engine
	Config       *options.Config
	Service      *svc.Service
	ConfigCenter *configures.ConfCenter
}

func (h *MigrateServer) onMigrateConfigUpdate(previous, current cc.ProcessConfig) {}

func parseShardingTableConfig(process *MigrateServer) error {
	if cc.IsExist("shardingTable.indexInterval") {
		val, err := cc.Int64("shardingTable.indexInterval")
		if err != nil {
			blog.Errorf("config shardingTable.indexInterval parse error. err: %s", err)
			return fmt.Errorf("config shardingTable.indexInterval parse error. err: %s", err)
		}
		if val < 30 || val > 720 {
			blog.Errorf("config shardingTable.indexInterval value illegal. must be in 20-720(minute), "+
				"but now val is %d", val)
			return fmt.Errorf("config shardingTable.indexInterval parse error. err: %s", err)
		}
		process.Config.ShardingTable.IndexesInterval = val
	} else {
		blog.Infof("config shardingTable.index not set. use default value(30m)")
		// IndexesInterval 表中同步索引间隔时间，单位分钟， 最小30分钟， 默认60分钟， 最大720分钟
		process.Config.ShardingTable.IndexesInterval = 60
	}

	// TableInterval模型shardingTable 对比和处理， 单位秒， 最小60秒，默认 120秒， 最大1800s
	if cc.IsExist("shardingTable.tableInterval") {
		val, err := cc.Int64("shardingTable.tableInterval")
		if err != nil {
			blog.Errorf("config shardingTable.tableInterval parse error. err: %s", err)
			return fmt.Errorf("config shardingTable.tableInterval parse error. err: %s", err)
		}
		if val < 30 || val > 720 {
			blog.Errorf("config shardingTable.tableInterval value illegal. must be in 60-1800(second), "+
				"but now val is %d", val)
			return fmt.Errorf("config shardingTable.tableInterval parse error. err: %s", err)
		}
		process.Config.ShardingTable.TableInterval = val

	} else {
		blog.Infof("config shardingTable.tableInterval not set. use default value(120s)")
		// TableInterval模型shardingTable 对比和处理， 单位秒， 最小60秒，默认 120秒， 最大1800s
		process.Config.ShardingTable.TableInterval = 120

	}

	return nil
}
