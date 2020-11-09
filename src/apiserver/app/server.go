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

	"configcenter/src/apimachinery/util"
	"configcenter/src/apiserver/app/options"
	"configcenter/src/apiserver/service"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/storage/dal/redis"

	"github.com/emicklei/go-restful"
)

// Run main loop function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	client, err := util.NewClient(&util.TLSClientConfig{})
	if err != nil {
		return fmt.Errorf("new proxy client failed, err: %v", err)
	}

	svc := service.NewService()

	apiSvr := new(APIServer)
	input := &backbone.BackboneParameter{
		ConfigUpdate: apiSvr.onApiServerConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}

	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	redisConf, err := engine.WithRedis()
	if err != nil {
		return err
	}
	cache, err := redis.NewFromConfig(redisConf)
	if err != nil {
		return fmt.Errorf("connect redis server failed, err: %s", err.Error())
	}

	limiter := service.NewLimiter(engine.ServiceManageClient().Client())
	err = limiter.SyncLimiterRules()
	if err != nil {
		blog.Infof("SyncLimiterRules failed, err: %v", err)
		return err
	}

	svc.SetConfig(engine, client, engine.Discovery(), engine.CoreAPI, cache, limiter)

	ctnr := restful.NewContainer()
	ctnr.Router(restful.CurlyRouter{})
	for _, item := range svc.WebServices() {
		ctnr.Add(item)
	}
	apiSvr.Core = engine

	err = backbone.StartServer(ctx, cancel, engine, ctnr, false)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
	}
	return nil
}

type APIServer struct {
	Core        *backbone.Engine
	Config      map[string]string
	configReady bool
}

func (h *APIServer) onApiServerConfigUpdate(previous, current cc.ProcessConfig) {
	h.configReady = true
}

const waitForSeconds = 180

