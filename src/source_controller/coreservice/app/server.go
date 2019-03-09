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
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/rdapi"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/source_controller/coreservice/app/options"
	coresvr "configcenter/src/source_controller/coreservice/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/redis"

	restful "github.com/emicklei/go-restful"
)

// CoreServer the core server
type CoreServer struct {
	Core    *backbone.Engine
	Config  options.Config
	Service coresvr.CoreServiceInterface
}

func (t *CoreServer) onCoreServiceConfigUpdate(previous, current cc.ProcessConfig) {

	t.Config.Mongo = mongo.ParseConfigFromKV("mongodb", current.ConfigMap)
	t.Config.Redis = redis.ParseConfigFromKV("redis", current.ConfigMap)

	blog.V(3).Infof("the new cfg:%#v the origin cfg:%#v", t.Config, current.ConfigMap)
	for t.Core == nil {
		time.Sleep(time.Second)
		blog.V(3).Info("sleep for engine")
	}
	t.Service.SetConfig(t.Config, t.Core, nil, nil)
}

// Run main function
func Run(ctx context.Context, op *options.ServerOption) error {
	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	discover, err := discovery.NewDiscoveryInterface(op.ServConf.RegDiscover)
	if err != nil {
		return fmt.Errorf("connect zookeeper [%s] failed: %v", op.ServConf.RegDiscover, err)
	}

	c := &util.APIMachineryConfig{
		QPS:       1000,
		Burst:     2000,
		TLSConfig: nil,
	}

	machinery, err := apimachinery.NewApiMachinery(c, discover)
	if err != nil {
		return fmt.Errorf("new api machinery failed, err: %v", err)
	}

	regPath := fmt.Sprintf("%s/%s/%s", types.CC_SERV_BASEPATH, types.CC_MODULE_CORESERVICE, svrInfo.IP)
	coreSvr := new(CoreServer)

	coreService := coresvr.New()
	coreSvr.Service = coreService

	webhandler := restful.NewContainer().Add(coreService.WebService())
	webhandler.ServiceErrorHandler(rdapi.ServiceErrorHandler)
	server := backbone.Server{
		ListenAddr: svrInfo.IP,
		ListenPort: svrInfo.Port,
		Handler:    webhandler,
		TLS:        backbone.TLSConfig{},
	}

	bonC := &backbone.Config{
		RegisterPath: regPath,
		RegisterInfo: *svrInfo,
		CoreAPI:      machinery,
		Server:       server,
	}

	engine, err := backbone.NewBackbone(
		ctx,
		op.ServConf.RegDiscover,
		types.CC_MODULE_CORESERVICE,
		op.ServConf.ExConfig,
		coreSvr.onCoreServiceConfigUpdate,
		discover,
		bonC)

	if nil != err {
		return fmt.Errorf("new engine failed, error is %s", err.Error())
	}

	coreSvr.Core = engine
	coreService.SetConfig(coreSvr.Config, engine, engine.CCErr, engine.Language)

	select {
	case <-ctx.Done():
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
