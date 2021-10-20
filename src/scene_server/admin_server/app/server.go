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

	"configcenter/src/ac/iam"
	"configcenter/src/common/auth"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/admin_server/app/options"
	"configcenter/src/scene_server/admin_server/configures"
	svc "configcenter/src/scene_server/admin_server/service"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/monitor"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	// init esb client
	esb.InitEsbClient(nil)

	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	process := new(AdminServer)
	process.Config = new(options.Config)
	if err := cc.SetMigrateFromFile(op.ServConf.ExConfig); err != nil {
		return fmt.Errorf("parse config file error %s", err.Error())
	}

	process.Config.Errors.Res, _ = cc.String("errors.res")
	process.Config.Language.Res, _ = cc.String("language.res")
	process.Config.Configures.Dir, _ = cc.String("confs.dir")
	process.Config.Register.Address, _ = cc.String("registerServer.addrs")
	process.Config.Register.User, _ = cc.String("registerServer.usr")
	process.Config.Register.Password, _ = cc.String("registerServer.pwd")
	process.Config.Register.CertFile, _ = cc.String("registerServer.certfile")
	process.Config.Register.KeyFile, _ = cc.String("registerServer.keyfile")
	process.Config.Register.CaFile, _ = cc.String("registerServer.cafile")
	snapDataID, _ := cc.Int("hostsnap.dataID")
	process.Config.SnapDataID = int64(snapDataID)

	// load mongodb, redis and common config from configure directory
	mongodbPath := process.Config.Configures.Dir + "/" + types.CCConfigureMongo
	if err := cc.SetMongodbFromFile(mongodbPath); err != nil {
		return fmt.Errorf("parse mongodb config from file[%s] failed, err: %v", mongodbPath, err)
	}

	redisPath := process.Config.Configures.Dir + "/" + types.CCConfigureRedis
	if err := cc.SetRedisFromFile(redisPath); err != nil {
		return fmt.Errorf("parse redis config from file[%s] failed, err: %v", redisPath, err)
	}

	commonPath := process.Config.Configures.Dir + "/" + types.CCConfigureCommon
	if err := cc.SetCommonFromFile(commonPath); err != nil {
		return fmt.Errorf("parse common config from file[%s] failed, err: %v", commonPath, err)
	}

	if err := monitor.InitMonitor(); err != nil {
		return fmt.Errorf("init monitor failed, err: %v", err)
	}

	mongoConf, err := cc.Mongo("mongodb")
	if err != nil {
		return err
	}
	process.Config.MongoDB = mongoConf

	watchDBConf, err := cc.Mongo("watch")
	if err != nil {
		return err
	}
	process.Config.WatchDB = watchDBConf

	redisConf, err := cc.Redis("redis")
	if err != nil {
		return err
	}
	process.Config.Redis = redisConf

	snapRedisConf, err := cc.Redis("redis.snap")
	if err != nil {
		return fmt.Errorf("get host snapshot redis configuration failed, err: %v", err)
	}
	process.Config.SnapRedis = snapRedisConf

	process.Config.Iam, err = iam.ParseConfigFromKV("authServer", nil)
	if err != nil && auth.EnableAuthorize() {
		blog.Errorf("parse iam error: %v", err)
		return err
	}

	if err := parseShardingTableConfig(process); err != nil {
		return err
	}

	input := &backbone.BackboneParameter{
		ConfigUpdate: process.onHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     process.Config.Register.Address,
		RdUser:       process.Config.Register.User,
		RdPassword:   process.Config.Register.Password,
		RdCertFile:   process.Config.Register.CertFile,
		RdKeyFile:    process.Config.Register.KeyFile,
		RdCaFile:     process.Config.Register.CaFile,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	process.Core = engine
	process.ConfigCenter = configures.NewConfCenter(ctx, engine.RegDiscv())

	// adminserver conf not depend discovery
	err = process.ConfigCenter.Start(
		process.Config.Configures.Dir,
		process.Config.Errors.Res,
		process.Config.Language.Res,
	)

	if err != nil {
		return err
	}

	service := svc.NewService(ctx)
	service.Engine = engine
	service.Config = *process.Config
	service.ConfigCenter = process.ConfigCenter
	process.Service = service

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

		watchDB, err := local.NewMgo(process.Config.WatchDB.GetMongoConf(), time.Minute)
		if err != nil {
			return fmt.Errorf("connect watch mongo server failed, err: %v", err)
		}
		process.Service.SetWatchDB(watchDB)

		cache, err := redis.NewFromConfig(process.Config.Redis)
		if err != nil {
			return fmt.Errorf("connect redis server failed, err: %s", err.Error())
		}
		process.Service.SetCache(cache)

		if auth.EnableAuthorize() {
			blog.Info("enable auth center access.")

			iamCli, err := iam.NewIam(nil, process.Config.Iam, engine.Metric().Registry())
			if err != nil {
				return fmt.Errorf("new iam client failed: %v", err)
			}
			process.Service.SetIam(iamCli)
		} else {
			blog.Infof("disable auth center access.")
		}

		if esbConfig, err := esb.ParseEsbConfig(""); err == nil {
			esb.UpdateEsbConfig(*esbConfig)
		}
		break
	}

	if err := service.BackgroundTask(*process.Config); err != nil {
		return err
	}
	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}

	errors.SetGlobalCCError(engine.CCErr)

	select {
	case <-ctx.Done():
	}
	blog.V(0).Info("process stopped")
	return nil
}

// AdminServer is data structure of adminserver
type AdminServer struct {
	Core         *backbone.Engine
	Config       *options.Config
	Service      *svc.Service
	ConfigCenter *configures.ConfCenter
}

func (h *AdminServer) onHostConfigUpdate(previous, current cc.ProcessConfig) {}

func parseShardingTableConfig(process *AdminServer) error {
	if cc.IsExist("shardingTable.indexInterval") {
		val, err := cc.Int64("shardingTable.indexInterval")
		if err != nil {
			blog.Errorf("config shardingTable.indexInterval parse error. err: %s", err)
			return fmt.Errorf("config shardingTable.indexInterval parse error. err: %s", err)
		}
		if val < 30 || val > 720 {
			blog.Errorf("config shardingTable.indexInterval value illegal. must be in 20-720(minute), but now val is %d",
				val)
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
			blog.Errorf("config shardingTable.tableInterval value illegal. must be in 60-1800(second), but now val is %d", val)
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
