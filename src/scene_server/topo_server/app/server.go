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
	"errors"
	"fmt"
	"strconv"
	"time"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/service"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/thirdpartyclient/elasticsearch"
)

// TopoServer the topo server
type TopoServer struct {
	Core        *backbone.Engine
	Config      options.Config
	Service     *service.Service
	configReady bool
}

func (t *TopoServer) onTopoConfigUpdate(previous, current cc.ProcessConfig, confType string) {
	t.configReady = true
	if current.ConfigMap["level.businessTopoMax"] != "" {
		max, err := strconv.Atoi(current.ConfigMap["level.businessTopoMax"])
		if err != nil {
			t.Config.BusinessTopoLevelMax = common.BKTopoBusinessLevelDefault
			blog.Errorf("invalid business topo max value, err: %v", err)
		} else {
			t.Config.BusinessTopoLevelMax = max
		}
		blog.Infof("config update with max topology level: %d", t.Config.BusinessTopoLevelMax)
	}

	t.Config.ConfigMap = current.ConfigMap
	blog.Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)

	var err error
	t.Config.Es, err = elasticsearch.ParseConfigFromKV("es", current.ConfigMap)
	if err != nil {
		blog.Warnf("parse es config failed: %v", err)
	}
}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	blog.Infof("srv conf: %+v", svrInfo)

	server := new(TopoServer)
	server.Config.BusinessTopoLevelMax = common.BKTopoBusinessLevelDefault
	server.Service = new(service.Service)

	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: server.onTopoConfigUpdate,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	server.Core = engine

	if err := server.CheckForReadiness(); err != nil {
		return err
	}

	server.Config.Mongo, err = engine.WithMongo()
	if err != nil {
		return err
	}
	server.Config.Redis, err = engine.WithRedis()
	if err != nil {
		return err
	}
	server.Config.Auth, err = engine.WithAuth()
	if err != nil {
		return err
	}

	enableTxn := false
	if server.Config.Mongo.TxnEnabled == "true" {
		enableTxn = true
	}
	blog.Infof("enableTxn is %t", enableTxn)
	txn, err := local.NewTransaction(enableTxn, server.Config.Mongo.GetMongoConf(), server.Config.Redis)
	if err != nil {
		return fmt.Errorf("initial transaction failed, err: %v", err)
	}

	authorize, err := authcenter.NewAuthCenter(nil, server.Config.Auth, engine.Metric().Registry())
	if err != nil {
		blog.Errorf("it is failed to create a new auth API, err:%s", err.Error())
		return err
	}

	essrv := new(elasticsearch.EsSrv)
	if server.Config.Es.FullTextSearch == "on" {
		esClient, err := elasticsearch.NewEsClient(server.Config.Es)
		if err != nil {
			blog.Errorf("failed to create elastic search client, err:%s", err.Error())
			return fmt.Errorf("new es client failed, err: %v", err)
		}
		essrv.Client = esClient
	}

	authManager := extensions.NewAuthManager(engine.CoreAPI, authorize)
	server.Service = &service.Service{
		Language:    engine.Language,
		Engine:      engine,
		AuthManager: authManager,
		Es:          essrv,
		Core:        core.New(engine.CoreAPI, authManager, engine.Language),
		Error:       engine.CCErr,
		Txn:         txn,
		EnableTxn:   enableTxn,
		Config:      server.Config,
	}

	err = backbone.StartServer(ctx, cancel, engine, server.Service.WebService(), true)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	}
	return nil
}

const waitForSeconds = 180

func (t *TopoServer) CheckForReadiness() error {
	for i := 1; i < waitForSeconds; i++ {
		if !t.configReady {
			blog.Info("waiting for topology server configuration ready.")
			time.Sleep(time.Second)
			continue
		}
		blog.Info("topology server configuration ready.")
		return nil
	}
	return errors.New("wait for topology server configuration timeout")
}
