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
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/proc_server/app/options"
	"configcenter/src/scene_server/proc_server/proc_service/service"
	"configcenter/src/storage/dal/redis"
	"configcenter/src/thirdpartyclient/esbserver"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"

	"github.com/emicklei/go-restful"
)

//Run ccapi server
func Run(ctx context.Context, op *options.ServerOption) error {

	discover, err := discovery.NewDiscoveryInterface(op.ServConf.RegDiscover)
	if err != nil {
		return fmt.Errorf("connect zookeeper [%s] failed: %v", op.ServConf.RegDiscover, err)
	}

	// clientset
	apiMachConf := &util.APIMachineryConfig{
		QPS:       op.ServConf.Qps,
		Burst:     op.ServConf.Burst,
		TLSConfig: nil,
	}

	apiMachinery, err := apimachinery.NewApiMachinery(apiMachConf, discover)
	if err != nil {
		return fmt.Errorf("create api machinery object failed. err: %v", err)
	}

	svrInfo, err := newServerInfo(op)
	if err != nil {
		blog.Errorf("fail to new server information. err: %s", err.Error())
		return fmt.Errorf("make server information failed, err:%v", err)
	}

	procSvr := new(service.ProcServer)
	procSvr.EsbConfigChn = make(chan esbutil.EsbConfig, 0)
	container := restful.NewContainer()
	container.Add(procSvr.WebService())

	bkbsvr := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    container,
		TLS:        backbone.TLSConfig{},
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_PROC, svrInfo.IP)
	bkbCfg := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      apiMachinery,
		Server:       bkbsvr,
	}

	engine, err := backbone.NewBackbone(ctx, op.ServConf.RegDiscover,
		types.CC_MODULE_PROC,
		op.ServConf.ExConfig,
		procSvr.OnProcessConfigUpdate,
		discover,
		bkbCfg)
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
		return fmt.Errorf("Configuration item not found")
	}
	cacheDB, err := redis.NewFromConfig(*procSvr.Config.Redis)
	if err != nil {
		blog.Errorf("new redis client failed, err: %s", err.Error())
		return fmt.Errorf("new redis client failed, err: %s", err)
	}

	esbSrv, err := esbserver.NewEsb(apiMachConf, procSvr.EsbConfigChn)
	if err != nil {
		return fmt.Errorf("create esb api  object failed. err: %v", err)
	}
	procSvr.Engine = engine
	procSvr.EsbServ = esbSrv
	procSvr.Cache = cacheDB
	go procSvr.InitFunc()

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
