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

	// 根据节点前缀发下节点下面的服务， 这样保证依赖的服务可以自动发现，无需配置。
	currSvrComponent := common.GetIdentification()
	for _, component := range types.LayerModuleMap[currSvrComponent.Layer] {
		path := fmt.Sprintf("%s/%s/%s", types.CCSvrBasePath, component.Layer.String(), component.Name)
		svr, err := newServerDiscover(disc, path, component.Name)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", component, err)
		}

		d.servers[component.Name] = svr
	}
	// 如果要支持第三方服务自动发现，
	// 需要watch  types.CCSvrBasePath 节点。发现有新的节点加入，
	//  对改节点执行newServerDiscover 方法。 这个操作对象需要加锁

	//  如果当前服务不是标准服务，发现自己的服务其他节点
	if strings.HasPrefix(currSvrComponent.Name, types.CC_DISCOVERY_PREFIX) {
		path := fmt.Sprintf("%s/%s/%s", types.CCSvrBasePath, currSvrComponent.Layer.String(), currSvrComponent.Name)
		svr, err := newServerDiscover(disc, path, currSvrComponent.Name)
		if err != nil {
			return nil, fmt.Errorf("discover %s failed, err: %v", currSvrComponent, err)
		}

		d.servers[currSvrComponent.Name] = svr

	}

	return d, nil
}

type discover struct {
	servers map[string]*server
}

func (d *discover) ApiServer() Interface {
	return d.servers[types.CCModuleAPIServer.Name]
}

func (d *discover) MigrateServer() Interface {
	return d.servers[types.CCModuleMigrate.Name]
}

func (d *discover) EventServer() Interface {
	return d.servers[types.CCModuleEventServer.Name]
}

func (d *discover) HostServer() Interface {
	return d.servers[types.CCModuleHost.Name]
}

func (d *discover) ProcServer() Interface {
	return d.servers[types.CCModuleProc.Name]
}

func (d *discover) TopoServer() Interface {
	return d.servers[types.CCModuleTop.Name]
}

func (d *discover) DataCollect() Interface {
	return d.servers[types.CCModuleDataCollection.Name]
}

func (d *discover) GseProcServer() Interface {
	return d.servers[types.GSEModuleProcServer.Name]
}

func (d *discover) CoreService() Interface {
	return d.servers[types.CCModuleCoerService.Name]
}

func (d *discover) TMServer() Interface {
	return d.servers[types.CCModuleTXC.Name]
}

func (d *discover) OperationServer() Interface {
	return d.servers[types.CCModuleOperation.Name]
}

func (d *discover) TaskServer() Interface {
	return d.servers[types.CCModuleTask.Name]
}

// IsMaster check whether current is master
func (d *discover) IsMaster() bool {
	return d.servers[common.GetIdentificationName()].IsMaster(common.GetServerInfo().Address())
}
