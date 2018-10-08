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

	"configcenter/src/apimachinery"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/proc_server/app/options"
	"configcenter/src/scene_server/proc_server/proc_service/service"
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

	svrInfo, err := newServerInfo(op)
	if err != nil {
		blog.Errorf("fail to new server information. err: %s", err.Error())
		return fmt.Errorf("make server information failed, err:%v", err)
	}

	procSvr := new(service.ProcServer)

	bkbsvr := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    procSvr.WebService(),
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
		bkbCfg)

	procSvr.Engine = engine

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

/*
//Run ccapi server
func Run(op *options.ServerOption) error {

	setConfig(op)

	serv, err := ccapi.NewCCAPIServer(op.ServConf)
	if err != nil {
		blog.Error("fail to create ccapi server. err:%s", err.Error())
		return err
	}

	//pid
	if err := common.SavePid(); err != nil {
		blog.Error("fail to save pid: err:%s", err.Error())
	}

	serv.Start()

	return nil
}

func setConfig(op *options.ServerOption) {
	//server cert directory

}
*/
