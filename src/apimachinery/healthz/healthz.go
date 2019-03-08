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

package healthz

import (
	"context"
	"errors"
	"fmt"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
	"configcenter/src/common/metric"
	"configcenter/src/common/types"
)

type HealthzInterface interface {
	HealthCheck(moduleName string) (healthy bool, err error)
}

func NewHealthzClient(capability *util.Capability, disc discovery.DiscoveryInterface) HealthzInterface {
	return &health{
		capability: capability,
		disc:       disc,
	}
}

type health struct {
	capability *util.Capability
	disc       discovery.DiscoveryInterface
}

func (h *health) HealthCheck(moduleName string) (healthy bool, err error) {
	var name string
	switch moduleName {
	case types.CC_MODULE_AUDITCONTROLLER:
		h.capability.Discover = h.disc.AuditCtrl()
		name = "audit"

	case types.CC_MODULE_HOSTCONTROLLER:
		h.capability.Discover = h.disc.HostCtrl()
		name = "host"

	case types.CC_MODULE_OBJECTCONTROLLER:
		h.capability.Discover = h.disc.ObjectCtrl()
		name = "object"

	case types.CC_MODULE_PROCCONTROLLER:
		h.capability.Discover = h.disc.ProcCtrl()
		name = "process"

	case types.CC_MODULE_DATACOLLECTION:
		h.capability.Discover = h.disc.DataCollect()
		name = "collector"

	case types.CC_MODULE_HOST:
		h.capability.Discover = h.disc.HostServer()
		name = "host"

	case types.CC_MODULE_MIGRATE:
		h.capability.Discover = h.disc.MigrateServer()
		name = "migrate"

	case types.CC_MODULE_PROC:
		h.capability.Discover = h.disc.ProcServer()
		name = "process"

	case types.CC_MODULE_TOPO:
		h.capability.Discover = h.disc.TopoServer()
		name = "topo"

	case types.CC_MODULE_EVENTSERVER:
		h.capability.Discover = h.disc.EventServer()
		name = "event"

	default:
		return false, fmt.Errorf("unsupported health module: %s", moduleName)
	}

	resp := new(metric.HealthResponse)
	client := rest.NewRESTClient(h.capability, fmt.Sprintf("/%s/v3", name))
	err = client.Get().
		WithContext(context.Background()).
		SubResource("/healthz").
		Body(nil).
		Do().
		Into(resp)

	if err != nil {
		return false, err
	}

	if !resp.Result {
		return false, errors.New(resp.Message)
	}

	return true, nil
}
