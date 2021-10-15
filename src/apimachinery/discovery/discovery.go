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
	"configcenter/src/common/blog"
	"configcenter/src/common/registerdiscover"
	"configcenter/src/common/types"
)

type ServiceManageInterface interface {
	// 判断当前进程是否为master进程
	IsMaster() bool
}

type DiscoveryInterface interface {
	ApiServer() Interface
	AdminServer() Interface
	EventServer() Interface
	HostServer() Interface
	ProcServer() Interface
	TopoServer() Interface
	DataCollect() Interface
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
	// 获取注册在服务发现上的所有服务节点
	GetServers() ([]string, error)
	// 最新的服务节点信息存放在该channel里，可被用来消费，以监听服务节点的变化
	GetServersChan() chan []string
}

// NewServiceDiscovery new a simple discovery module which can be used to get alive server address
func NewServiceDiscovery(rd *registerdiscover.RegDiscv) (DiscoveryInterface, error) {
	d := &discover{
		servers: make(map[string]*server),
	}

	curServiceName := common.GetIdentification()
	services := types.GetDiscoveryService()
	// 将当前服务也放到需要发现中
	services[curServiceName] = struct{}{}
	for component := range services {
		// 如果所有服务都按需发现服务。这个地方时不需要配置
		if component == types.CCModuleWeb && curServiceName != types.CCModuleWeb {
			continue
		}
		path := fmt.Sprintf("%s/%s", types.CCDiscoverBaseEndpoint, component)
		svr, err := newServerDiscover(rd, path, component)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", component, err)
		}

		d.servers[component] = svr
	}

	electPath := fmt.Sprintf("%s/%s", types.CCDiscoverBaseElection, curServiceName)
	master := newServerMaster(rd, electPath, curServiceName)
	d.master = master

	return d, nil
}

type discover struct {
	servers map[string]*server
	master  *master
}

// ApiServer returns apiserver info
func (d *discover) ApiServer() Interface {
	return d.servers[types.CCModuleApi]
}

// AdminServer returns adminserver info
func (d *discover) AdminServer() Interface {
	return d.servers[types.CCModuleAdmin]
}

// EventServer returns eventserver info
func (d *discover) EventServer() Interface {
	return d.servers[types.CCModuleEvent]
}

// HostServer returns hostserver info
func (d *discover) HostServer() Interface {
	return d.servers[types.CCModuleHost]
}

// ProcServer returns procserver info
func (d *discover) ProcServer() Interface {
	return d.servers[types.CCModuleProc]
}

// TopoServer returns toposerver info
func (d *discover) TopoServer() Interface {
	return d.servers[types.CCModuleTopo]
}

// DataCollect returns datacollection info
func (d *discover) DataCollect() Interface {
	return d.servers[types.CCModuleDataCollection]
}

// CoreService returns coreservice info
func (d *discover) CoreService() Interface {
	return d.servers[types.CCModuleCoreService]
}

// OperationServer returns operationserver info
func (d *discover) OperationServer() Interface {
	return d.servers[types.CCModuleOperation]
}

// TaskServer returns taskserver info
func (d *discover) TaskServer() Interface {
	return d.servers[types.CCModuleTask]
}

// CloudServer returns cloudserver info
func (d *discover) CloudServer() Interface {
	return d.servers[types.CCModuleCloud]
}

// AuthServer returns authserver info
func (d *discover) AuthServer() Interface {
	return d.servers[types.CCModuleAuth]
}

// CacheService returns cacheservice info
func (d *discover) CacheService() Interface {
	return d.servers[types.CCModuleCacheService]
}

// IsMaster checks whether current instance is master
func (d *discover) IsMaster() bool {
	return d.master.IsMaster()
}

// Server 根据服务名获取服务再服务发现组件中的相关信息
func (d *discover) Server(name string) Interface {
	if svr, ok := d.servers[name]; ok {
		return svr
	}
	blog.V(5).Infof("not found server. name: %s", name)

	return emptyServerInst
}
