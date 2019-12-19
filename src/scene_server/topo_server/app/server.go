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
	"os"
	"strconv"
	"time"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdpartyclient/elasticsearch"
)

// TopoServer the topo server
type TopoServer struct {
	Core        *backbone.Engine
	Config      options.Config
	Service     *service.Service
	configReady bool
}

func (t *TopoServer) onTopoConfigUpdate(previous, current cc.ProcessConfig) {
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
	t.Config.Mongo = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)
	t.Config.Redis = redis.ParseConfigFromKV("redis", current.ConfigMap)
	t.Config.FullTextSearch = current.ConfigMap["es.full_text_search"]
	t.Config.EsUrl = current.ConfigMap["es.url"]
	t.Config.ConfigMap = current.ConfigMap
	blog.Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)

	var err error
	t.Config.Auth, err = authcenter.ParseConfigFromKV("auth", current.ConfigMap)
	if err != nil {
		blog.Warnf("parse auth center config failed: %v", err)
	}
}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
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

	cache, err := redis.NewFromConfig(server.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %v", err)
		return err
	}

	db, err := local.NewMgo(server.Config.Mongo.BuildURI(), time.Second*5)
	if err != nil {
		blog.Errorf("failed to connect the txc server, error info is %v", err)
		return err
	}
	err = db.InitTxnManager(cache)
	if err != nil {
		blog.Errorf("failed to init txn manager, error info is %v", err)
		return err
	}

	enableTxn := false
	if server.Config.Mongo.TxnEnabled == "true" {
		enableTxn = true
	}
	blog.Infof("enableTxn is %t", enableTxn)


	authorize, err := authcenter.NewAuthCenter(nil, server.Config.Auth, engine.Metric().Registry())
	if err != nil {
		blog.Errorf("it is failed to create a new auth API, err:%s", err.Error())
		return err
	}

	essrv := new(elasticsearch.EsSrv)
	if server.Config.FullTextSearch == "on" {
		// if use https, config tls.Config{xxx}, and instead NewEsClient param nil
		esclient, err := elasticsearch.NewEsClient(server.Config.EsUrl, nil)
		if err != nil {
			blog.Errorf("failed to create elastic search client, err:%s", err.Error())
			return fmt.Errorf("new es client failed, err: %v", err)
		}
		essrv.Client = esclient
	}

	authManager := extensions.NewAuthManager(engine.CoreAPI, authorize)
	server.Service = &service.Service{
		Language:    engine.Language,
		Engine:      engine,
		AuthManager: authManager,
		Es:          essrv,
		Core:        core.New(engine.CoreAPI, authManager),
		Error:       engine.CCErr,
		DB:         db,
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

func newServerInfo(op *options.ServerOption) (*types.ServerInfo, error) {
	ip, err := op.ServConf.GetAddress()
	if err != nil {
		return nil, err
	}

	port, err := op.ServConf.GetPort()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	info := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}
	return info, nil
}
