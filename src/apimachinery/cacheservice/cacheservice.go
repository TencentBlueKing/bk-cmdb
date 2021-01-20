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

package cacheservice

import (
	"fmt"

	"configcenter/src/apimachinery/cacheservice/cache/event"
	"configcenter/src/apimachinery/cacheservice/cache/host"
	"configcenter/src/apimachinery/cacheservice/cache/topology"
	"configcenter/src/apimachinery/rest"
	"configcenter/src/apimachinery/util"
)

type Cache interface {
	Host() host.Interface
	Topology() topology.Interface
	Event() event.Interface
}

type CacheServiceClientInterface interface {
	Cache() Cache
}

func NewCacheServiceClient(c *util.Capability, version string) CacheServiceClientInterface {
	base := fmt.Sprintf("/cache/%s", version)
	return &cacheService{
		restCli: rest.NewRESTClient(c, base),
	}
}

type cacheService struct {
	restCli rest.ClientInterface
}

type cache struct {
	restCli rest.ClientInterface
}

func (c *cacheService) Cache() Cache {
	return &cache{
		restCli: c.restCli,
	}
}

func (c *cache) Host() host.Interface {
	return host.NewCacheClient(c.restCli)
}

func (c *cache) Topology() topology.Interface {
	return topology.NewCacheClient(c.restCli)
}

func (c *cache) Event() event.Interface {
	return event.NewCacheClient(c.restCli)
}
