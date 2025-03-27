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
	"configcenter/src/thirdparty/esbserver/esbutil"
	"configcenter/src/thirdparty/esbserver/iam"

	"github.com/prometheus/client_golang/prometheus"
)

// EsbClientInterface TODO
type EsbClientInterface interface {
	IamSrv() iam.IamClientInterface
}

type esbsrv struct {
	client rest.ClientInterface
	iamSrv iam.IamClientInterface
	sync.RWMutex
	esbConfig *esbutil.EsbConfigSrv
	c         *util.Capability
}

// NewEsb new a esb client
func NewEsb(apiMachineryConfig *util.APIMachineryConfig, cfgChan chan esbutil.EsbConfig, defaultCfg *esbutil.EsbConfig,
	reg prometheus.Registerer) (EsbClientInterface, error) {
	base := fmt.Sprintf("/api/c/compapi")

	client, err := util.NewClient(apiMachineryConfig.TLSConfig)
	if nil != err {
		return nil, err
	}
	flowControl := flowctrl.NewRateLimiter(apiMachineryConfig.QPS, apiMachineryConfig.Burst)
	esbConfig := esbutil.NewEsbConfigSrv(cfgChan, defaultCfg)

	esbCapability := &util.Capability{
		Client:     client,
		Discover:   esbConfig,
		Throttle:   flowControl,
		MetricOpts: util.MetricOption{Register: reg},
	}

	esb := &esbsrv{
		client:    rest.NewRESTClient(esbCapability, base),
		esbConfig: esbConfig,
	}
	return esb, nil
}

// IamSrv TODO
func (e *esbsrv) IamSrv() iam.IamClientInterface {
	e.RLock()
	srv := e.iamSrv
	e.RUnlock()
	if nil == srv {
		e.Lock()
		e.iamSrv = iam.NewIamClientInterface(e.client, e.esbConfig)
		srv = e.iamSrv
		e.Unlock()
	}
	return srv
}

// GetEsbConfigSrv TODO
func (e *esbsrv) GetEsbConfigSrv() *esbutil.EsbConfigSrv {
	return e.esbConfig
}
