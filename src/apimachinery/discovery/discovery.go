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

package discovery

import (
	"fmt"
	"time"

	regd "configcenter/src/common/RegisterDiscover"
	"configcenter/src/common/types"
)

type DiscoveryInterface interface {
	MigrateServer() Interface
	EventServer() Interface
	HostServer() Interface
	ProcServer() Interface
	TopoServer() Interface
	DataCollect() Interface
	AuditCtrl() Interface
	HostCtrl() Interface
	ObjectCtrl() Interface
	ProcCtrl() Interface
	GseProcServ() Interface
}

type Interface interface {
	GetServers() ([]string, error)
}

func NewDiscoveryInterface(zkAddr string) (DiscoveryInterface, error) {
	disc := regd.NewRegDiscoverEx(zkAddr, 10*time.Second)
	if err := disc.Start(); nil != err {
		return nil, err
	}

	d := &discover{
		servers: make(map[string]Interface),
	}
	for component, _ := range types.AllModule {
		if component == types.CC_MODULE_APISERVER || component == types.CC_MODULE_WEBSERVER {
			continue
		}
		path := fmt.Sprintf("%s/%s", types.CC_SERV_BASEPATH, component)
		svr, err := newServerDiscover(disc, path)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", component, err)
		}

		d.servers[component] = svr
	}

	return d, nil
}

type discover struct {
	servers map[string]Interface
}

func (d *discover) MigrateServer() Interface {
	return d.servers[types.CC_MODULE_MIGRATE]
}

func (d *discover) EventServer() Interface {
	return d.servers[types.CC_MODULE_EVENTSERVER]
}

func (d *discover) HostServer() Interface {
	return d.servers[types.CC_MODULE_HOST]
}

func (d *discover) ProcServer() Interface {
	return d.servers[types.CC_MODULE_PROC]
}

func (d *discover) TopoServer() Interface {
	return d.servers[types.CC_MODULE_TOPO]
}

func (d *discover) DataCollect() Interface {
	return d.servers[types.CC_MODULE_DATACOLLECTION]
}

func (d *discover) AuditCtrl() Interface {
	return d.servers[types.CC_MODULE_AUDITCONTROLLER]
}

func (d *discover) HostCtrl() Interface {
	return d.servers[types.CC_MODULE_HOSTCONTROLLER]
}

func (d *discover) ObjectCtrl() Interface {
	return d.servers[types.CC_MODULE_OBJECTCONTROLLER]
}

func (d *discover) ProcCtrl() Interface {
	return d.servers[types.CC_MODULE_PROCCONTROLLER]
}

func (d *discover) GseProcServ() Interface {
	return d.servers[types.GSE_MODULE_PROCSERVER]
}
