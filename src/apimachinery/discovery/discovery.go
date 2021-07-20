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

	"configcenter/src/common"
	"configcenter/src/common/backbone/service_mange/zk"
	"configcenter/src/common/blog"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

type ServiceManageInterface interface {
	// 判断当前进程是否为master 进程， 服务注册节点的第一个节点
	IsMaster() bool
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
	CloudServer() Interface
	AuthServer() Interface
	Server(name string) Interface
	CacheService() Interface
	ServiceManageInterface
}

type Interface interface {
	// 获取注册在zk上的所有服务节点
	GetServers() ([]string, error)
	// 最新的服务节点信息存放在该channel里，可被用来消费，以监听服务节点的变化
	GetServersChan() chan []string
}

// NewServiceDiscovery new a simple discovery module which can be used to get alive server address
func NewServiceDiscovery(client *zk.ZkClient) (DiscoveryInterface, error) {
	disc := registerdiscover.NewRegDiscoverEx(client)

	d := &discover{
		servers: make(map[string]*server),
	}

	curServiceName := common.GetIdentification()
	services := types.GetDiscoveryService()
	// 将当前服务也放到需要发现中
	services[curServiceName] = struct{}{}
	for component := range services {
		// 如果所有服务都按需发现服务。这个地方时不需要配置
		if component == types.CC_MODULE_WEBSERVER && curServiceName != types.CC_MODULE_WEBSERVER {
			continue
		}
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

func (d *discover) OperationServer() Interface {
	return d.servers[types.CC_MODULE_OPERATION]
}

func (d *discover) TaskServer() Interface {
	return d.servers[types.CC_MODULE_TASK]
}

func (d *discover) CloudServer() Interface {
	return d.servers[types.CC_MODULE_CLOUD]
}

func (d *discover) AuthServer() Interface {
	return d.servers[types.CC_MODULE_AUTH]
}

func (d *discover) CacheService() Interface {
	return d.servers[types.CC_MODULE_CACHESERVICE]
}

// IsMaster check whether current is master
func (d *discover) IsMaster() bool {
	return d.servers[common.GetIdentification()].IsMaster(common.GetServerInfo().UUID)
}

// Server 根据服务名获取服务再服务发现组件中的相关信息
func (d *discover) Server(name string) Interface {
	if svr, ok := d.servers[name]; ok {
		return svr
	}
	blog.V(5).Infof("not found server. name: %s", name)

	return emptyServerInst
}
