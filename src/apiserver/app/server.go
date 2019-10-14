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
	"time"

	"configcenter/src/apimachinery/util"
	"configcenter/src/apiserver/app/options"
	"configcenter/src/apiserver/service"
	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"

	"github.com/emicklei/go-restful"
)

// Run main loop function
func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
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

	if err := apiSvr.CheckForReadiness(); err != nil {
		return err
	}

	authConf, err := authcenter.ParseConfigFromKV("auth", apiSvr.Config)
	if err != nil {
		return err
	}
	authorize, err := auth.NewAuthorize(nil, authConf, engine.Metric().Registry())
	if err != nil {
		return fmt.Errorf("new authorize failed, err: %v", err)
	}
	blog.Infof("enable authcenter: %v", authorize.Enabled())

	svc.SetConfig(engine, client, engine.Discovery(), authorize)

	ctnr := restful.NewContainer()
	ctnr.Router(restful.CurlyRouter{})
	for _, item := range svc.WebServices(authConf) {
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
	h.Config = current.ConfigMap
}

const waitForSeconds = 180

func (h *APIServer) CheckForReadiness() error {
	for i := 1; i < waitForSeconds; i++ {
		if !h.configReady {
			blog.Info("waiting for api server configuration ready.")
			time.Sleep(time.Second)
			continue
		}
		return nil
	}
	return errors.New("wait for api server configuration timeout")
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
