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
	"configcenter/src/thirdparty/esbserver/gse"
	"configcenter/src/thirdparty/esbserver/iam"
	"configcenter/src/thirdparty/esbserver/nodeman"
	"configcenter/src/thirdparty/esbserver/user"

	"github.com/prometheus/client_golang/prometheus"
)

type EsbClientInterface interface {
	GseSrv() gse.GseClientInterface
	User() user.UserClientInterface
	NodemanSrv() nodeman.NodeManClientInterface
	IamSrv() iam.IamClientInterface
}

type esbsrv struct {
	client     rest.ClientInterface
	gseSrv     gse.GseClientInterface
	userSrv    user.UserClientInterface
	nodemanSrv nodeman.NodeManClientInterface
	iamSrv     iam.IamClientInterface
	sync.RWMutex
	esbConfig *esbutil.EsbConfigSrv
	c         *util.Capability
}

// NewEsb new a esb client
//
func NewEsb(apiMachineryConfig *util.APIMachineryConfig, cfgChan chan esbutil.EsbConfig, defaultCfg *esbutil.EsbConfig, reg prometheus.Registerer) (EsbClientInterface, error) {
	base := fmt.Sprintf("/api/c/compapi")

	client, err := util.NewClient(apiMachineryConfig.TLSConfig)
	if nil != err {
		return nil, err
	}
	flowControl := flowctrl.NewRateLimiter(apiMachineryConfig.QPS, apiMachineryConfig.Burst)
	esbConfig := esbutil.NewEsbConfigSrv(cfgChan, defaultCfg)

	esbCapability := &util.Capability{
		Client:   client,
		Discover: esbConfig,
		Throttle: flowControl,
		Reg:      reg,
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

func (e *esbsrv) NodemanSrv() nodeman.NodeManClientInterface {
	e.RLock()
	srv := e.nodemanSrv
	e.RUnlock()
	if nil == srv {
		e.Lock()
		e.nodemanSrv = nodeman.NewNodeManClientInterface(e.client, e.esbConfig)
		srv = e.nodemanSrv
		e.Unlock()
	}
	return srv
}

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

func (e *esbsrv) GetEsbConfigSrv() *esbutil.EsbConfigSrv {
	return e.esbConfig
}

func (e *esbsrv) User() user.UserClientInterface {
	e.RLock()
	srv := e.userSrv
	e.RUnlock()
	if nil == srv {
		e.Lock()
		e.userSrv = user.NewUserClientInterface(e.client, e.esbConfig)
		srv = e.userSrv
		e.Unlock()
	}
	return srv
}
