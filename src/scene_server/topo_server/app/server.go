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
	"net/http"
	"time"

	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/topo_server/app/options"
	"configcenter/src/scene_server/topo_server/core"
	"configcenter/src/scene_server/topo_server/service"
	"configcenter/src/storage/driver/redis"
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

	t.Config.ConfigMap = current.ConfigMap
	blog.Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)

	var err error
	t.Config.Es, err = elasticsearch.ParseConfigFromKV("es", current.ConfigMap)
	if err != nil {
		blog.Warnf("parse es config failed: %v", err)
	}
}

func (t *TopoServer) setBusinessTopoLevelMax() error {
	tryCnt := 30
	for i := 1; i <= tryCnt; i++ {
		time.Sleep(time.Second * 2)
		res, err := t.Core.CoreAPI.CoreService().System().SearchConfigAdmin(context.Background(), http.Header{})
		if err != nil {
			blog.Warnf("setBusinessTopoLevelMax failed,  try count:%d, SearchConfigAdmin err: %v", i, err)
			continue
		}
		if res.Result == false {
			blog.Warnf("setBusinessTopoLevelMax failed,  try count:%d, SearchConfigAdmin err: %s", i, res.ErrMsg)
			continue
		}
		t.Config.BusinessTopoLevelMax = int(res.Data.Backend.MaxBizTopoLevel)
		break
	}

	if t.Config.BusinessTopoLevelMax == 0 {
		blog.Errorf("setBusinessTopoLevelMax failed, BusinessTopoLevelMax is 0, check the coreservice and the value in table cc_System")
		return fmt.Errorf("setBusinessTopoLevelMax failed")
	}

	blog.Infof("setBusinessTopoLevelMax successfully, BusinessTopoLevelMax is %d", t.Config.BusinessTopoLevelMax)
	return nil
}

// Run main function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	blog.Infof("srv conf: %+v", svrInfo)
	blog.Infof("enableTxn is %t", op.EnableTxn)

	server := new(TopoServer)
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

	if err := server.setBusinessTopoLevelMax(); err != nil {
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

	authorize, err := authcenter.NewAuthCenter(nil, server.Config.Auth, engine.Metric().Registry())
	if err != nil {
		blog.Errorf("it is failed to create a new auth API, err:%s", err.Error())
		return err
	}
	// TODO  redis, auth 可以在backbone 完成
	if err := redis.InitClient("redis", &server.Config.Redis); err != nil {
		blog.Errorf("it is failed to connect reids. err:%s", err.Error())
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
		EnableTxn:   op.EnableTxn,
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
