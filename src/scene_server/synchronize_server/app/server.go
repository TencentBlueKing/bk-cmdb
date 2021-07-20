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
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/backbone"
	cc "configcenter/src/common/backbone/configcenter"
	//"configcenter/src/common/blog"
	"configcenter/src/common/types"
	"configcenter/src/scene_server/synchronize_server/app/options"
	synchronizeService "configcenter/src/scene_server/synchronize_server/service"
	//"configcenter/src/storage/dal/redis"
	synchronizeClient "configcenter/src/apimachinery/synchronize"
	synchronizeUtil "configcenter/src/apimachinery/synchronize/util"
)

func Run(ctx context.Context, cancel context.CancelFunc, op *options.ServerOption) error {
	svrInfo, err := types.NewServerInfo(op.ServConf)
	if err != nil {
		return fmt.Errorf("wrap server info failed, err: %v", err)
	}

	service := new(synchronizeService.Service)
	synchronSrv := &SynchronizeServer{
		synchronizeClientConfig: make(chan synchronizeUtil.SychronizeConfig, 10),
	}
	input := &backbone.BackboneParameter{
		Regdiscv:     op.ServConf.RegDiscover,
		ConfigPath:   op.ServConf.ExConfig,
		ConfigUpdate: synchronSrv.onSynchronizeServerConfigUpdate,
		SrvInfo:      svrInfo,
	}
	engine, err := backbone.NewBackbone(ctx, input)
	if err != nil {
		return fmt.Errorf("new backbone failed, err: %v", err)
	}
	configReady := false
	for sleepCnt := 0; sleepCnt < common.APPConfigWaitTime; sleepCnt++ {
		if synchronSrv.Config == nil {
			time.Sleep(time.Second)
		} else {
			configReady = true
			break
		}
	}
	if false == configReady {
		return fmt.Errorf("Configuration item not found")
	}
	service.Engine = engine
	service.Config = synchronSrv.Config
	synchronSrv.Service = service
	synchronizeClientInst, err := synchronizeClient.NewSynchronize(engine.ApiMachineryConfig(), synchronSrv.synchronizeClientConfig)
	if err != nil {
		return fmt.Errorf("new NewSynchronize failed, err: %v", err)
	}
	service.SetSynchronizeServer(synchronizeClientInst)
	go synchronSrv.Service.InitBackground()
	err = backbone.StartServer(ctx, cancel, engine, service.WebService(), true)
	if err != nil {
		return err
	}
	select {
	case <-ctx.Done():
	}
	return nil
}

type SynchronizeServer struct {
	Core                    *backbone.Engine
	Config                  *options.Config
	Service                 *synchronizeService.Service
	synchronizeClientConfig chan synchronizeUtil.SychronizeConfig
}

func (s *SynchronizeServer) onSynchronizeServerConfigUpdate(previous, current cc.ProcessConfig) {

	configInfo := &options.Config{}
	names, _ := cc.String("synchronizeServer.name")
	configInfo.Names = SplitFilter(names, ",")

	configInfo.Trigger.TriggerType, _ = cc.String("synchronizeServer.trigger.type")
	// role  unit minute.
	// type = timing, ervery day  role minute trigger
	// type = interval, interval role  minute trigger
	configInfo.Trigger.Role, _ = cc.String("synchronizeServer.trigger.role")

	for _, name := range configInfo.Names {
		if strings.TrimSpace(name) == "" {
			continue
		}
		configItem := &options.ConfigItem{}
		appNames, _ := cc.String("synchronizeServer." + name + ".AppNames")
		syncResource, _ := cc.String("synchronizeServer." + name + ".SynchronizeResource")
		targetHost, _ := cc.String("synchronizeServer." + name + ".Host")
		fieldSign, _ := cc.String("synchronizeServer." + name + ".FieldSign")
		dataSign, _ := cc.String("synchronizeServer." + name + ".Flag")
		supplerAccount, _ := cc.String("synchronizeServer." + name + ".SupplerAccount")
		whiteList, _ := cc.String("synchronizeServer." + name + ".WhiteList")
		objectIDs, _ := cc.String("synchronizeServer." + name + ".ObjectID")
		ignoreModelAttr, _ := cc.String("synchronizeServer." + name + ".IgnoreModelAttribute")
		strEnableInstFilter, _ := cc.String("synchronizeServer." + name + ".EnableInstFilter")

		configItem.AppNames = SplitFilter(appNames, ",")
		if syncResource == "1" {
			configItem.SyncResource = true
		}
		if whiteList == "1" {
			configItem.WhiteList = true
		}
		// 使用忽略模型属性变的模式。 模型属性，模型分组 将不做同步
		// 但是数据源cmdb中满足条件的实例会同步到目标cmdb。
		// 目标cmdb中新建相同的唯一标识模型或者模型的字段。内容会自动展示出来
		if ignoreModelAttr == "1" {
			configItem.IgnoreModelAttr = true
		}

		configItem.ObjectIDArr = SplitFilter(objectIDs, ",")
		configItem.Name = name
		configItem.TargetHost = targetHost
		configItem.FieldSign = fieldSign
		configItem.SynchronizeFlag = dataSign
		configItem.SupplerAccount = SplitFilter(supplerAccount, ",")
		if strEnableInstFilter == "1" {
			configItem.EnableInstFilter = true
		}

		configInfo.ConifgItemArray = append(configInfo.ConifgItemArray, configItem)
		if targetHost != "" {
			s.synchronizeClientConfig <- synchronizeUtil.SychronizeConfig{
				Name:  name,
				Addrs: []string{targetHost},
			}
		}

	}
	s.Config = configInfo

}

// SplitFilter split string with sep. remove blanks for blank item children and children
func SplitFilter(s, sep string) []string {
	itemArr := strings.Split(s, sep)
	var strArr []string
	for _, item := range itemArr {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		strArr = append(strArr, item)
	}
	return strArr
}
