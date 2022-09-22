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

package api

import (
	"configcenter/api/adminserver"
	"configcenter/api/apiserver"
	"configcenter/api/authserver"
	"configcenter/api/cacheservice"
	"configcenter/api/cloudserver"
	"configcenter/api/coreservice"
	"configcenter/api/discovery"
	"configcenter/api/eventserver"
	"configcenter/api/flowctrl"
	"configcenter/api/healthz"
	"configcenter/api/hostserver"
	"configcenter/api/procserver"
	"configcenter/api/taskserver"
	"configcenter/api/toposerver"
	"configcenter/api/util"
)

// ClientSetInterface TODO
type ClientSetInterface interface {
	HostServer() hostserver.HostServerClientInterface
	TopoServer() toposerver.TopoServerClientInterface
	ProcServer() procserver.ProcServerClientInterface
	AdminServer() adminserver.AdminServerClientInterface
	ApiServer() apiserver.ApiServerClientInterface
	EventServer() eventserver.EventServerClientInterface

	CoreService() coreservice.CoreServiceClientInterface
	TaskServer() taskserver.TaskServerClientInterface
	CloudServer() cloudserver.CloudServerClientInterface
	AuthServer() authserver.AuthServerClientInterface

	CacheService() cacheservice.CacheServiceClientInterface

	Healthz() healthz.HealthzInterface
}

// NewApiMachinery TODO
func NewApiMachinery(c *util.APIMachineryConfig, discover discovery.DiscoveryInterface) (ClientSetInterface, error) {
	client, err := util.NewClient(c.TLSConfig)
	if err != nil {
		return nil, err
	}

	flowcontrol := flowctrl.NewRateLimiter(c.QPS, c.Burst)
	return NewClientSet(client, discover, flowcontrol), nil
}

// NewClientSet TODO
func NewClientSet(client util.HttpClient, discover discovery.DiscoveryInterface, throttle flowctrl.RateLimiter) ClientSetInterface {
	return &ClientSet{
		version:  "v3",
		client:   client,
		discover: discover,
		throttle: throttle,
	}
}

// NewMockClientSet TODO
func NewMockClientSet() *ClientSet {
	return &ClientSet{
		version:  "unit_test",
		client:   nil,
		discover: discovery.NewMockDiscoveryInterface(),
		throttle: flowctrl.NewMockRateLimiter(),
		Mock:     util.MockInfo{Mocked: true},
	}
}

// ClientSet TODO
type ClientSet struct {
	version  string
	client   util.HttpClient
	discover discovery.DiscoveryInterface
	throttle flowctrl.RateLimiter
	Mock     util.MockInfo
}

// HostServer TODO
func (cs *ClientSet) HostServer() hostserver.HostServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.HostServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	cs.Mock.SetMockData = false
	return hostserver.NewHostServerClientInterface(c, cs.version)
}

// TopoServer TODO
func (cs *ClientSet) TopoServer() toposerver.TopoServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.TopoServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	cs.Mock.SetMockData = false
	return toposerver.NewTopoServerClient(c, cs.version)
}

// ProcServer TODO
func (cs *ClientSet) ProcServer() procserver.ProcServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.ProcServer(),
		Throttle: cs.throttle,
	}
	cs.Mock.SetMockData = false
	return procserver.NewProcServerClientInterface(c, cs.version)
}

// AdminServer TODO
func (cs *ClientSet) AdminServer() adminserver.AdminServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.MigrateServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	cs.Mock.SetMockData = false
	return adminserver.NewAdminServerClientInterface(c, cs.version)
}

// ApiServer TODO
func (cs *ClientSet) ApiServer() apiserver.ApiServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.ApiServer(),
		Throttle: cs.throttle,
	}
	return apiserver.NewApiServerClientInterface(c, cs.version)
}

// EventServer TODO
func (cs *ClientSet) EventServer() eventserver.EventServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.EventServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	cs.Mock.SetMockData = false
	return eventserver.NewEventServerClientInterface(c, cs.version)
}

// Healthz TODO
func (cs *ClientSet) Healthz() healthz.HealthzInterface {
	c := &util.Capability{
		Client:   cs.client,
		Throttle: cs.throttle,
	}
	return healthz.NewHealthzClient(c, cs.discover)
}

// CoreService TODO
func (cs *ClientSet) CoreService() coreservice.CoreServiceClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.CoreService(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	return coreservice.NewCoreServiceClient(c, cs.version)
}

// TaskServer TODO
func (cs *ClientSet) TaskServer() taskserver.TaskServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.TaskServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	return taskserver.NewProcServerClientInterface(c, cs.version)
}

// CloudServer TODO
func (cs *ClientSet) CloudServer() cloudserver.CloudServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.CloudServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	return cloudserver.NewCloudServerClientInterface(c, cs.version)
}

// AuthServer TODO
func (cs *ClientSet) AuthServer() authserver.AuthServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.AuthServer(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	return authserver.NewAuthServerClientInterface(c, cs.version)
}

// CacheService TODO
func (cs *ClientSet) CacheService() cacheservice.CacheServiceClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.CacheService(),
		Throttle: cs.throttle,
		Mock:     cs.Mock,
	}
	return cacheservice.NewCacheServiceClient(c, cs.version)
}
