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

package apimachinery

import (
	"configcenter/src/apimachinery/adminserver"
	"configcenter/src/apimachinery/auditcontroller"
	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/eventserver"
	"configcenter/src/apimachinery/flowctrl"
	"configcenter/src/apimachinery/gseprocserver"
	"configcenter/src/apimachinery/healthz"
	"configcenter/src/apimachinery/hostcontroller"
	"configcenter/src/apimachinery/hostserver"
	"configcenter/src/apimachinery/objcontroller"
	"configcenter/src/apimachinery/proccontroller"
	"configcenter/src/apimachinery/procserver"
	"configcenter/src/apimachinery/toposerver"
	"configcenter/src/apimachinery/util"
)

type ClientSetInterface interface {
	HostServer() hostserver.HostServerClientInterface
	TopoServer() toposerver.TopoServerClientInterface
	ProcServer() procserver.ProcServerClientInterface
	AdminServer() adminserver.AdminServerClientInterface
	EventServer() eventserver.EventServerClientInterface

	ObjectController() objcontroller.ObjControllerClientInterface
	AuditController() auditcontroller.AuditCtrlInterface
	ProcController() proccontroller.ProcCtrlClientInterface
	HostController() hostcontroller.HostCtrlClientInterface

	GseProcServer() gseprocserver.GseProcClientInterface
	Healthz() healthz.HealthzInterface
}

func NewApiMachinery(c *util.APIMachineryConfig) (ClientSetInterface, error) {
	client, err := util.NewClient(c.TLSConfig)
	if err != nil {
		return nil, err
	}

	discover, err := discovery.NewDiscoveryInterface(c.ZkAddr, c.GseProcServ)
	if err != nil {
		return nil, err
	}

	flowcontrol := flowctrl.NewRateLimiter(c.QPS, c.Burst)
	return NewClientSet(client, discover, flowcontrol), nil
}

func NewClientSet(client util.HttpClient, discover discovery.DiscoveryInterface, throttle flowctrl.RateLimiter) ClientSetInterface {
	return &ClientSet{
		version:  "v3",
		client:   client,
		discover: discover,
		throttle: throttle,
	}
}

type ClientSet struct {
	version  string
	client   util.HttpClient
	discover discovery.DiscoveryInterface
	throttle flowctrl.RateLimiter
}

func (cs *ClientSet) HostServer() hostserver.HostServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.HostServer(),
		Throttle: cs.throttle,
	}
	return hostserver.NewHostServerClientInterface(c, cs.version)
}

func (cs *ClientSet) TopoServer() toposerver.TopoServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.TopoServer(),
		Throttle: cs.throttle,
	}
	return toposerver.NewTopoServerClient(c, cs.version)
}

func (cs *ClientSet) ObjectController() objcontroller.ObjControllerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.ObjectCtrl(),
		Throttle: cs.throttle,
	}
	return objcontroller.NewObjectControllerInterface(c, cs.version)
}

func (cs *ClientSet) ProcServer() procserver.ProcServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.ProcServer(),
		Throttle: cs.throttle,
	}
	return procserver.NewProcServerClientInterface(c, cs.version)
}

func (cs *ClientSet) AdminServer() adminserver.AdminServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.MigrateServer(),
		Throttle: cs.throttle,
	}
	return adminserver.NewAdminServerClientInterface(c, cs.version)
}

func (cs *ClientSet) EventServer() eventserver.EventServerClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.EventServer(),
		Throttle: cs.throttle,
	}
	return eventserver.NewEventServerClientInterface(c, cs.version)
}

func (cs *ClientSet) AuditController() auditcontroller.AuditCtrlInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.AuditCtrl(),
		Throttle: cs.throttle,
	}
	return auditcontroller.NewAuditCtrlInterface(c, cs.version)
}

func (cs *ClientSet) ProcController() proccontroller.ProcCtrlClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.ProcCtrl(),
		Throttle: cs.throttle,
	}
	return proccontroller.NewProcCtrlClientInterface(c, cs.version)
}

func (cs *ClientSet) HostController() hostcontroller.HostCtrlClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.HostCtrl(),
		Throttle: cs.throttle,
	}
	return hostcontroller.NewHostCtrlClientInterface(c, cs.version)
}

func (cs *ClientSet) GseProcServer() gseprocserver.GseProcClientInterface {
	c := &util.Capability{
		Client:   cs.client,
		Discover: cs.discover.GseProcServ(),
		Throttle: cs.throttle,
	}
	return gseprocserver.NewGseProcClientInterface(c, "v1")
}

func (cs *ClientSet) Healthz() healthz.HealthzInterface {
	c := &util.Capability{
		Client:   cs.client,
		Throttle: cs.throttle,
	}
	return healthz.NewHealthzClient(c, cs.discover)
}
