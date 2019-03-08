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
)

type EsbConfigServ struct {
	addrs     string
	appCode   string
	appSecret string
	sync.RWMutex
}

type EsbServDiscoveryInterace interface {
	GetServers() ([]string, error)
}

func NewEsbConfigServ(srvChan chan EsbConfig) *EsbConfigServ {
	esb := &EsbConfigServ{}
	go func() {
		if nil == srvChan {
			return
		}
		for {
			config := <-srvChan
			esb.Lock()
			esb.addrs = config.Addrs
			esb.appCode = config.AppCode
			esb.appSecret = config.AppSecret
			esb.Unlock()
		}
	}()

	return esb
}

func (esb *EsbConfigServ) GetEsbServDiscoveryInterace() EsbServDiscoveryInterace {
	// mabye will deal some logic about server
	return esb
}

func (esb *EsbConfigServ) GetServers() ([]string, error) {
	// mabye will deal some logic about server
	esb.RLock()
	defer esb.RUnlock()
	return []string{esb.addrs}, nil
}

func (esb *EsbConfigServ) GetConfig() EsbConfig {
	esb.RLock()
	defer esb.RUnlock()
	return EsbConfig{
		Addrs:     esb.addrs,
		AppCode:   esb.appCode,
		AppSecret: esb.appSecret,
	}
}
