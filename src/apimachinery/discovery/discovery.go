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
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

type ServiceManageInterface interface {
	// 判断当前进程是否为master 进程， 服务注册节点的第一个节点
	IsMaster() bool
	TMServer() Interface
}

type DiscoveryInterface interface {
	ApiServer() Interface
	MigrateServer() Interface
	EventServer() Interface
	HostServer() Interface
	ProcServer() Interface
	TopoServer() Interface
	DataCollect() Interface
	GseProcServer() Interface
	CoreService() Interface
	OperationServer() Interface
	TaskServer() Interface
	ServiceManageInterface
}

type Interface interface {
	GetServers() ([]string, error)
}

// NewServiceDiscovery new a simple discovery module which can be used to get alive server address
func NewServiceDiscovery(client *zk.ZkClient) (DiscoveryInterface, error) {
	disc := registerdiscover.NewRegDiscoverEx(client)

	d := &discover{
		servers: make(map[string]*server),
	}
	for component := range types.AllModule {
		if component == types.CC_MODULE_WEBSERVER {
			continue
		}
		path := fmt.Sprintf("%s/%s", types.CC_SERV_BASEPATH, component)
		svr, err := newServerDiscover(disc, path, component)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", component, err)
		}

		d.servers[component] = svr
	}
	// 如果要支持第三方服务自动发现，
	// 需要watch  types.CC_SERV_BASEPATH 节点。发现有新的节点加入，
	//  对改节点执行newServerDiscover 方法。 这个操作d 对象需要加锁

	//  如果当前服务不是标准服务，发现自己的服务其他节点
	component := common.GetIdentification()
	if strings.HasPrefix(common.GetIdentification(), types.CC_DISCOVERY_PREFIX) {
		path := fmt.Sprintf("%s/%s", types.CC_SERV_BASEPATH, component)
		svr, err := newServerDiscover(disc, path, component)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", component, err)
		}

		d.servers[component] = svr

	}

	return d, nil
}

type discover struct {
	servers map[string]*server
}

func (d *discover) ApiServer() Interface {
	return d.servers[types.CC_MODULE_APISERVER]
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

func (d *discover) GseProcServer() Interface {
	return d.servers[types.GSE_MODULE_PROCSERVER]
}

func (d *discover) CoreService() Interface {
	return d.servers[types.CC_MODULE_CORESERVICE]
}

func (d *discover) TMServer() Interface {
	return d.servers[types.CC_MODULE_TXC]
}

func (d *discover) OperationServer() Interface {
	return d.servers[types.CC_MODULE_OPERATION]
}

func (d *discover) TaskServer() Interface {
	return d.servers[types.CC_MODULE_TASK]
}

// IsMaster check whether current is master
func (d *discover) IsMaster() bool {
	return d.servers[common.GetIdentification()].IsMaster(common.GetServerInfo().Address())
}
