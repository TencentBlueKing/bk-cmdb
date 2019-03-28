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

package synchronize

import (
	"net/http"
	"sync"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/synchronize/synchronizeserver"
	synchronizeUtil "configcenter/src/apimachinery/synchronize/util"
	"configcenter/src/apimachinery/util"
)

type SynchronizeClientInterface interface {
	SynchronizeSrv(flag string) synchronizeserver.SynchronizeClientInterface
}

type synchronizeSrv struct {
	apiMachineryConfig *util.APIMachineryConfig
	baseClient         *http.Client
	synchronizeClient  map[string]synchronizeserver.SynchronizeClientInterface
	sync.RWMutex
	synchronizeConfig *synchronizeUtil.SynchronizeConfigServ
}

func NewSynchronize(apiMachineryConfig *util.APIMachineryConfig, config chan synchronizeUtil.SychronizeConfig) (SynchronizeClientInterface, error) {

	client, err := util.NewClient(apiMachineryConfig.TLSConfig)
	if nil != err {
		return nil, err
	}
	synchronizeConfig := synchronizeUtil.NewSynchronizeConfigServ(config)

	synchronize := &synchronizeSrv{
		synchronizeClient:  make(map[string]synchronizeserver.SynchronizeClientInterface, 0),
		baseClient:         client,
		apiMachineryConfig: apiMachineryConfig,
		synchronizeConfig:  synchronizeConfig,
	}
	return synchronize, nil
}

func (s *synchronizeSrv) SynchronizeSrv(flag string) synchronizeserver.SynchronizeClientInterface {
	s.RLock()
	srv, ok := s.synchronizeClient[flag]
	s.RUnlock()
	if nil == srv || !ok {
		s.Lock()
		s.synchronizeClient[flag] = synchronizeserver.NewSychronizeClientInterface(s.getSrvClent(flag))
		srv = s.synchronizeClient[flag]
		s.Unlock()
	}
	return srv
}

func (s *synchronizeSrv) getSrvClent(flag string) rest.ClientInterface {
	flowcontrol := flowctrl.NewRateLimiter(s.apiMachineryConfig.QPS, s.apiMachineryConfig.Burst)
	config := synchronizeUtil.NewSyncrhonizeConfig(flag)

	capability := &util.Capability{
		Client:   s.baseClient,
		Discover: config,
		Throttle: flowcontrol,
	}

	return rest.NewRESTClient(capability, "/synchronize/v3")
}
