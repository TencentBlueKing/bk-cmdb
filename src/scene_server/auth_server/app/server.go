/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2019 THL A29 Limited, a Tencent company. All rights reserved.
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
	"sync"
	"time"

	"configcenter/src/ac/iam"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/resource/esb"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/auth_server/app/options"
	"configcenter/src/scene_server/auth_server/logics"
	"configcenter/src/scene_server/auth_server/sdk/auth"
	"configcenter/src/scene_server/auth_server/sdk/client"
	sdktypes "configcenter/src/scene_server/auth_server/sdk/types"
	"configcenter/src/scene_server/auth_server/service"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	// init esb client
	esb.InitEsbClient(nil)

	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap authServer info failed, err: %v", err)
	}

	authServer := new(AuthServer)

	input := &backbone.BackboneParameter{
		ConfigUpdate: authServer.onAuthConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	authServer.Core = engine
	for {
		if authServer.Config == nil {
			time.Sleep(time.Second * 2)
			blog.V(3).Info("config not found, retry 2s later")
			continue
		}

		authConf := authServer.Config.Auth
		iamConf := sdktypes.IamConfig{
			Address:   authConf.Address,
			AppCode:   authConf.AppCode,
			AppSecret: authConf.AppSecret,
			SystemID:  authConf.SystemID,
			TLS:       authServer.Config.TLS,
		}
		opt := sdktypes.Options{
			Metric: engine.Metric().Registry(),
		}

		iamCli, err := client.NewClient(iamConf, opt)
		if err != nil {
			blog.Errorf("new iam client, err: %s", err.Error())
			return err
		}

		authConfig := sdktypes.Config{
			Iam:     iamConf,
			Options: opt,
		}
		lgc := logics.NewLogics(engine.CoreAPI)
		authorizer, err := auth.NewAuth(authConfig, lgc)
		if err != nil {
			return fmt.Errorf("new authorize failed, err: %v", err)
		}

		authServer.Service = service.NewAuthService(engine, iamCli, lgc, authorizer)
		break
	}
	err = backbone.StartServer(ctx, cancel, engine, authServer.Service.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("auth server will exit!")
	}

	return nil
}

type AuthServer struct {
	Core    *backbone.Engine
	Config  *options.Config
	Service *service.AuthService
}

var configLock sync.Mutex

func (a *AuthServer) onAuthConfigUpdate(previous, current cc.ProcessConfig) {
	configLock.Lock()
	defer configLock.Unlock()
	if len(current.ConfigData) > 0 {
		if a.Config == nil {
			a.Config = new(options.Config)
		}
		blog.InfoJSON("config updated: \n%s", string(current.ConfigData))
		var err error
		a.Config.Auth, err = iam.ParseConfigFromKV("authServer", nil)
		if err != nil {
			blog.Warnf("parse auth center config failed: %v", err)
		}

		a.Config.TLS, err = util.NewTLSClientConfigFromConfig("authServer", nil)
		if err != nil {
			blog.Warnf("parse auth center tls config failed: %v", err)
		}

		if esbConfig, err := esb.ParseEsbConfig("authServer"); err == nil {
			esb.UpdateEsbConfig(*esbConfig)
		}
	}
}
