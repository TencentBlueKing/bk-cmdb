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

package esbserver

import (
	"fmt"
	"sync"

	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/thirdpartyclient/esbserver/esbutil"
	"configcenter/src/thirdpartyclient/esbserver/gse"
)

type EsbClientInterface interface {
	GseSrv() gse.GseClientInterface
}

type esbsrv struct {
	client rest.ClientInterface
	gseSrv gse.GseClientInterface
	sync.RWMutex
	esbConfig *esbutil.EsbConfigServ
	c         *util.Capability
}

func NewEsb(apiMachineryConfig *util.APIMachineryConfig, config chan esbutil.EsbConfig) (EsbClientInterface, error) {
	base := fmt.Sprintf("/api/c/compapi")

	client, err := util.NewClient(apiMachineryConfig.TLSConfig)
	if nil != err {
		return nil, err
	}
	flowcontrol := flowctrl.NewRateLimiter(apiMachineryConfig.QPS, apiMachineryConfig.Burst)
	esbConfig := esbutil.NewEsbConfigServ(config)

	esbCapability := &util.Capability{
		Client:   client,
		Discover: esbConfig,
		Throttle: flowcontrol,
	}
	esb := &esbsrv{
		client:    rest.NewRESTClient(esbCapability, base),
		esbConfig: esbConfig,
	}
	return esb, nil
}

func (e *esbsrv) GseSrv() gse.GseClientInterface {
	e.RLock()
	srv := e.gseSrv
	e.RUnlock()
	if nil == srv {
		e.Lock()
		e.gseSrv = gse.NewGsecClientInterface(e.client, e.esbConfig)
		srv = e.gseSrv
		e.Unlock()
	}
	return srv
}

func (e *esbsrv) GetEsbConfigSrv() *esbutil.EsbConfigServ {
	return e.esbConfig
}
