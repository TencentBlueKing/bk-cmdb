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
package esbutil

import (
	"sync"

	"configcenter/src/common/blog"
)

type EsbConfigSrv struct {
	addrs     string
	appCode   string
	appSecret string
	sync.RWMutex
}

type EsbSrvDiscoveryInterface interface {
	GetServers() ([]string, error)
}

func NewEsbConfigSrv(srvChan chan EsbConfig, defaultCfg *EsbConfig) *EsbConfigSrv {
	esb := &EsbConfigSrv{}

	if defaultCfg != nil {
		esb.addrs = defaultCfg.Addrs
		esb.appCode = defaultCfg.AppCode
		esb.appSecret = defaultCfg.AppSecret
	}

	go func() {
		if srvChan == nil {
			return
		}
		for {
			config := <-srvChan
			esb.Lock()
			esb.addrs = config.Addrs
			esb.appCode = config.AppCode
			esb.appSecret = config.AppSecret
			blog.Infof("cmdb config changed, config: %+v", config)
			esb.Unlock()
		}
	}()

	return esb
}

func (esb *EsbConfigSrv) GetEsbSrvDiscoveryInterface() EsbSrvDiscoveryInterface {
	// maybe will deal some logic about server
	return esb
}

func (esb *EsbConfigSrv) GetServers() ([]string, error) {
	// maybe will deal some logic about server
	esb.RLock()
	defer esb.RUnlock()
	return []string{esb.addrs}, nil
}

func (esb *EsbConfigSrv) GetServersChan() chan []string {
	return nil
}

func (esb *EsbConfigSrv) GetConfig() EsbConfig {
	esb.RLock()
	defer esb.RUnlock()
	return EsbConfig{
		Addrs:     esb.addrs,
		AppCode:   esb.appCode,
		AppSecret: esb.appSecret,
	}
}
