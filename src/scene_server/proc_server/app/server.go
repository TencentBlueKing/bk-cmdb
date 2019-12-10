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
	"os"
	"time"

	"configcenter/src/auth"
	"configcenter/src/auth/authcenter"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/proc_server/app/options"
	"configcenter/src/scene_server/proc_server/logics"
	"configcenter/src/scene_server/proc_server/service"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {

	svrInfo, err := newServerInfo(op)
	if err != nil {
		blog.Errorf("fail to new server information. err: %s", err.Error())
		return fmt.Errorf("make server information failed, err:%v", err)
	}

	procSvr := new(service.ProcServer)
	procSvr.EsbConfigChn = make(chan esbutil.EsbConfig, 0)

	input := &backbone.BackboneParameter{
		ConfigUpdate: procSvr.OnProcessConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == procSvr.Config {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("configuration item not found")
	}

	// transaction client
	txn, err := procSvr.Config.Mongo.GetTransactionClient(engine)
	if err != nil {
		blog.Errorf("new transaction client failed, err: %+v", err)
		return fmt.Errorf("new transaction client failed, err: %+v", err)
	}
	procSvr.TransactionClient = txn

	authConf, err := authcenter.ParseConfigFromKV("auth", procSvr.ConfigMap)
	if err != nil {
		return err
	}

	authorize, err := auth.NewAuthorize(nil, authConf, engine.Metric().Registry())
	if err != nil {
		return fmt.Errorf("new authorize failed, err: %v", err)
	}

	esbSrv, err := esbserver.NewEsb(engine.ApiMachineryConfig(), procSvr.EsbConfigChn, nil, engine.Metric().Registry())
	if err != nil {
		return fmt.Errorf("create esb api  object failed. err: %v", err)
	}
	procSvr.AuthManager = extensions.NewAuthManager(engine.CoreAPI, authorize)
	procSvr.Engine = engine
	procSvr.EsbSrv = esbSrv
	procSvr.Logic = &logics.Logic{
		Engine: procSvr.Engine,
	}

	err = backbone.StartServer(ctx, cancel, engine, procSvr.WebService(), true)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		blog.Infof("process will exit!")
	}

	return nil
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

	svrInfo := &types.ServerInfo{
		IP:       ip,
		Port:     port,
		HostName: hostname,
		Scheme:   "http",
		Version:  version.GetVersion(),
		Pid:      os.Getpid(),
	}

	return svrInfo, nil
}
