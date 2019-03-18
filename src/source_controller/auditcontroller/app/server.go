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

	restful "github.com/emicklei/go-restful"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/common/version"
	"configcenter/src/source_controller/auditcontroller/app/options"
	"configcenter/src/source_controller/auditcontroller/logics"
	"configcenter/src/source_controller/auditcontroller/service"
	"configcenter/src/storage/dal/mongo"
	"configcenter/src/storage/dal/mongo/local"
)

//Run ccapi server
func Run(ctx context.Context, op *options.ServerOption) error {

	svrInfo, err := newServerInfo(op)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %s", err.Error())
	}

	coreService := new(service.Service)
	audit := new(AuditController)
	audit.Service = coreService
	coreService.Logics = &logics.Logics{Instance: audit.Instance, Engine: audit.Service.Engine}
	input := &backbone.BackboneParameter{
		ConfigUpdate: audit.onAduitConfigUpdate,
		ConfigPath:   op.ServConf.ExConfig,
		Regdiscv:     op.ServConf.RegDiscover,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	coreService.Engine = engine
	coreService.Logics.Engine = engine

	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if nil == audit.Config {
			time.Sleep(time.Second)
			continue
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Failed to get configuration")
	}
	if err := backbone.StartServer(ctx, engine, restful.NewContainer().Add(coreService.WebService())); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		break
	}
	return nil
}

// AuditController  audit controller config
type AuditController struct {
	*service.Service
	Config *options.Config
}

func (h *AuditController) onAduitConfigUpdate(previous, current cc.ProcessConfig) {
	h.Config = &options.Config{
		Mongo: mongo.ParseConfigFromKV("mongodb", current.ConfigMap),
	}

	instance, err := local.NewMgo(h.Config.Mongo.BuildURI(), time.Minute)
	if err != nil {
		blog.Errorf("new mongo client failed, err: %s", err.Error())
		return
	}
	h.Service.Instance = instance
	h.Service.Logics.Instance = instance

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
