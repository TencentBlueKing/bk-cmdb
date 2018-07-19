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

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/source_controller/proccontroller/app/options"
	"configcenter/src/source_controller/proccontroller/service"
	"configcenter/src/storage/mgoclient"
	"configcenter/src/storage/redisclient"
)

//Run ccapi server
func Run(ctx context.Context, op *options.ServerOption) error {
	// clientset
	apiMachConf := &util.APIMachineryConfig{
		ZkAddr:    op.ServConf.RegDiscover,
		QPS:       op.ServConf.Qps,
		Burst:     op.ServConf.Burst,
		TLSConfig: nil,
	}

	apiMachinery, err := apimachinery.NewApiMachinery(apiMachConf)
	if err != nil {
		return fmt.Errorf("create api machinery object failed. err: %v", err)
	}
	// server
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("creae server info object failed. err: %v", err)
	}

	proctrlSvr := new(service.ProctrlServer)

	bksvr := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    proctrlSvr.WebService(),
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_PROCCONTROLLER, svrInfo.IP)
	bkConf := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      apiMachinery,
		Server:       bksvr,
	}

	proctrlSvr.Core, err = backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_PROCCONTROLLER,
		op.ServConf.ExConfig,
		proctrlSvr.OnProcessConfUpdate,
		bkConf)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}

	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == proctrlSvr.MongoCfg {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}

	proctrlSvr.DbInstance, err = mgoclient.NewMgoCli(proctrlSvr.MongoCfg.Address, proctrlSvr.MongoCfg.Port, proctrlSvr.MongoCfg.User, proctrlSvr.MongoCfg.Password, proctrlSvr.MongoCfg.Mechanism, proctrlSvr.MongoCfg.Database)
	if err != nil {
		return fmt.Errorf("new mongo client failed, err: %v", err)
	}
	err = proctrlSvr.DbInstance.Open()
	if err != nil {
		return fmt.Errorf("new mongo client failed, err: %v", err)
	}

	proctrlSvr.CacheDI, err = redisclient.NewFromConfig(*proctrlSvr.RedisCfg)
	if err != nil {
		return fmt.Errorf("new redis client failed, err: %v", err)
	}
	select {
	case <-ctx.Done():
		blog.Errorf("processctroller will exit!")
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
