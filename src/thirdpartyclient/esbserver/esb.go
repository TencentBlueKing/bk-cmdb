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

func NewEsb(c *util.Capability, config chan esbutil.EsbConfig, version string) EsbClientInterface {
	base := fmt.Sprintf("/api/c/compapi")
	esbCapability := *c
	esbConfig := esbutil.NewEsbServConfig(config)
	esbCapability.Discover = esbConfig

	esb := &esbsrv{
		client:    rest.NewRESTClient(&esbCapability, base),
		esbConfig: esbConfig,
	}
	return esb
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
