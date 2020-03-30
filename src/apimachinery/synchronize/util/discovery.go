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

type SychronizeConfig struct {
	Name  string
	Addrs []string
}

type SynchronizeConfigServ struct {
	addrs map[string][]string
	sync.RWMutex
}

type SynchronizeServDiscoveryInterace interface {
	GetServers() ([]string, error)
	GetServersChan() chan []string
}

var (
	synchronize = &SynchronizeConfigServ{
		addrs: make(map[string][]string, 0),
	}
)

func NewSynchronizeConfigServ(srvChan chan SychronizeConfig) *SynchronizeConfigServ {
	go func() {
		if nil == srvChan {
			return
		}
		for {
			config := <-srvChan
			synchronize.Lock()
			synchronize.addrs[config.Name] = config.Addrs
			synchronize.Unlock()
		}
	}()

	return synchronize
}

type synchronizeConfig struct {
	flag string
}

func NewSyncrhonizeConfig(flag string) SynchronizeServDiscoveryInterace {
	return &synchronizeConfig{
		flag: flag,
	}
}

func (s *synchronizeConfig) GetEsbServDiscoveryInterace(flag string) SynchronizeServDiscoveryInterace {
	// mabye will deal some logics about server
	return s
}

func (s *synchronizeConfig) GetServers() ([]string, error) {
	// mabye will deal some logics about server
	synchronize.RLock()
	defer synchronize.RUnlock()

	return synchronize.addrs[s.flag], nil
}

func (s *synchronizeConfig) GetServersChan() chan []string {
	return nil
}
