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

	"configcenter/src/ac/iam"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/event_server/app/options"
	svc "configcenter/src/scene_server/event_server/service"
	"configcenter/src/scene_server/event_server/sync/hostidentifier"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdparty/gse/client"
)

const (
	// defaultInitWaitDuration is default duration for new EventServer init.
	defaultInitWaitDuration = time.Second

	// defaultDBConnectTimeout is default connect timeout of cc db.
	defaultDBConnectTimeout = 5 * time.Second

	// defaultBatchSyncIntervalHours default batch sync host identifier interval hours
	defaultBatchSyncIntervalHours = 24

	// defaultGoroutineCount default goroutine count
	defaultGoroutineCount = 10
)

// EventServer is event server.
type EventServer struct {
	ctx    context.Context
	engine *backbone.Engine

	// config for this eventserver app.
	config *options.Config

	// service main service instance.
	service *svc.Service

	// make host configs update action safe.
	hostConfigUpdateMu sync.Mutex

	// db is cc main database.
	db dal.RDB

	// redisCli is cc redis client.
	redisCli redis.Client
}

// NewEventServer creates a new EventServer object.
func NewEventServer(ctx context.Context, op *options.ServerOption) (*EventServer, error) {
	// build server info.
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return nil, fmt.Errorf("build server info, %+v", err)
	}

	// new EventServer instance.
	newEventServer := &EventServer{ctx: ctx}

	engine, err := backbone.NewBackbone(ctx, &backbone.BackboneParameter{
		ConfigUpdate: newEventServer.OnHostConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	})
	if err != nil {
		return nil, fmt.Errorf("build backbone, %+v", err)
	}

	// set global cc errors.
	errors.SetGlobalCCError(engine.CCErr)

	// set backbone engine.
	newEventServer.engine = engine
	newEventServer.service = svc.NewService(ctx, engine)

	return newEventServer, nil
}

// Engine returns engine of the EventServer instance.
func (es *EventServer) Engine() *backbone.Engine {
	return es.engine
}

// Service returns main service of the EventServer instance.
func (es *EventServer) Service() *svc.Service {
	return es.service
}

// DB returns cc database client of the EventServer instance.
func (es *EventServer) DB() dal.RDB {
	return es.db
}

// RedisCli returns cc redis client of the EventServer instance.
func (es *EventServer) RedisCli() redis.Client {
	return es.redisCli
}

// OnHostConfigUpdate is callback for updating configs.
func (es *EventServer) OnHostConfigUpdate(prev, curr cc.ProcessConfig) {
	es.hostConfigUpdateMu.Lock()
	defer es.hostConfigUpdateMu.Unlock()

	if len(curr.ConfigData) > 0 {
		// NOTE: allow to update configs with empty values?
		// NOTE: what is prev used for? build a compare logic here?

		if es.config == nil {
			es.config = &options.Config{}
		}

		if data, err := json.MarshalIndent(curr.ConfigData, "", "  "); err == nil {
			blog.Infof("on host config update event: \n%s", data)
		}

		// TODO: add your configs updates here.
	}

}

// initConfigs inits configs for new EventServer server.
func (es *EventServer) initConfigs() error {
	for {
		// wait and parse configs that async updated by backbone engine.
		es.hostConfigUpdateMu.Lock()
		if es.config == nil {
			es.hostConfigUpdateMu.Unlock()

			blog.Info("can't find configs to run the new eventserver, try again later!")
			time.Sleep(defaultInitWaitDuration)
			continue
		}

		// ready to init new instance.
		es.hostConfigUpdateMu.Unlock()
		break
	}

	var err error
	blog.Info("found configs to run the new eventserver now!")

	// mongodb.
	es.config.MongoDB, err = es.engine.WithMongo()
	if err != nil {
		return fmt.Errorf("init mongodb configs, %+v", err)
	}

	// cc redis.
	es.config.Redis, err = es.engine.WithRedis()
	if err != nil {
		return fmt.Errorf("init cc redis configs, %+v", err)
	}

	return nil
}

// initModules inits modules for new EventServer.
func (es *EventServer) initModules() error {
	// create mongodb client.
	db, err := local.NewMgo(es.config.MongoDB.GetMongoConf(), defaultDBConnectTimeout)
	if err != nil {
		return fmt.Errorf("create new mongodb client, %+v", err)
	}
	es.db = db
	es.service.SetDB(db)
	blog.Info("init modules, create mongo client success[%+v]", es.config.MongoDB.GetMongoConf())

	// connect to cc redis.
	redisCli, err := redis.NewFromConfig(es.config.Redis)
	if err != nil {
		return fmt.Errorf("connect to cc redis, %+v", err)
	}
	es.redisCli = redisCli
	es.service.SetCache(redisCli)
	blog.Infof("init modules, connected to cc redis, %+v", es.config.Redis)

	// initialize auth authorizer
	es.service.SetAuthorizer(iam.NewAuthorizer(es.engine.CoreAPI))

	return nil
}

// Run runs a new EventServer.
func (es *EventServer) Run() error {
	// init configs.
	if err := es.initConfigs(); err != nil {
		return err
	}
	blog.Info("init configs success!")

	// ready to setup comms for new server instance now.
	if err := es.initModules(); err != nil {
		return err
	}
	blog.Info("init modules success!")

	startUp, err := cc.Bool("eventServer.hostIdentifier.startUp")
	if err != nil || !startUp {
		return nil
	}
	if err := es.runSyncData(); err != nil {
		return err
	}
	blog.Info("run sync data success!")
	return nil
}

func (es *EventServer) runSyncData() error {
	gseTaskServerConfig, err := client.NewGseConnConfig("gse.taskServer")
	if err != nil {
		return err
	}
	gseTaskServerClient, err := client.NewGseTaskServerClient(gseTaskServerConfig.Endpoints,
		gseTaskServerConfig.TLSConf)
	if err != nil {
		return err
	}
	gseApiServerConfig, err := client.NewGseConnConfig("gse.apiServer")
	if err != nil {
		return err
	}
	gseApiServerClient, err := client.NewGseApiServerClient(gseApiServerConfig.Endpoints, gseApiServerConfig.TLSConf)
	if err != nil {
		return err
	}
	syncData := hostidentifier.NewHostIdentifier(es.ctx, es.redisCli, es.engine,
		gseTaskServerClient, gseApiServerClient)

	// watch主机身份变化创建任务调用gse接口推送
	go syncData.WatchToSyncHostIdentifier()

	// 周期全量同步主机身份
	batchSyncIntervalHours, err := cc.Int("eventServer.hostIdentifier.batchSyncIntervalHours")
	if err != nil {
		batchSyncIntervalHours = defaultBatchSyncIntervalHours
	}
	go func() {
		for {
			if !es.engine.Discovery().IsMaster() {
				time.Sleep(time.Minute)
				continue
			}
			syncData.BatchSyncHostIdentifier()
			time.Sleep(time.Duration(batchSyncIntervalHours) * time.Hour)
		}
	}()

	for i := 0; i < defaultGoroutineCount; i++ {
		// 查询推送主机身份任务结果并处理失败主机
		go syncData.GetTaskExecutionStatus()
		// 协程将失败的主机重新变成新任务
		go syncData.MakeNewTaskFromFailHost()
	}

	return nil
}

// Run setups a new EventServer app with a context and options and runs it as server instance.
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	// create eventserver.
	eventServer, err := NewEventServer(ctx, op)
	if err != nil {
		return fmt.Errorf("create new eventserver, %+v", err)
	}

	// run new event server.
	if err := eventServer.Run(); err != nil {
		return err
	}

	// all modules is initialized success, start the new server now.
	if err := backbone.StartServer(ctx, cancel, eventServer.Engine(), eventServer.Service().WebService(), true); err != nil {
		return err
	}
	blog.Info("EventServer init and run success!")

	<-ctx.Done()
	blog.Info("EventServer stopping now!")
	return nil
}
