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

// Package cache TODO
package cache

import (
	"fmt"

	"configcenter/src/apimachinery/discovery"
	"configcenter/src/source_controller/cacheservice/cache/host"
	"configcenter/src/source_controller/cacheservice/cache/mainline"
	"configcenter/src/source_controller/cacheservice/cache/topology"
	"configcenter/src/source_controller/cacheservice/cache/topotree"
	"configcenter/src/source_controller/cacheservice/event/watch"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/driver/mongodb"
	"configcenter/src/storage/driver/redis"
	"configcenter/src/storage/reflector"
	"configcenter/src/storage/stream"
)

// NewCache TODO
func NewCache(reflector reflector.Interface, loopW stream.LoopInterface, isMaster discovery.ServiceManageInterface,
	watchDB dal.DB) (*ClientSet, error) {

	if err := mainline.NewMainlineCache(loopW); err != nil {
		return nil, fmt.Errorf("new business cache failed, err: %v", err)
	}

	if err := host.NewCache(reflector); err != nil {
		return nil, fmt.Errorf("new host cache failed, err: %v", err)
	}

	topo, err := topology.NewTopology(isMaster, loopW)
	if err != nil {
		return nil, err
	}

	mainlineClient := mainline.NewMainlineClient()
	hostClient := host.NewClient()

	cache := &ClientSet{
		Tree:     topotree.NewTopologyTree(mainlineClient),
		Host:     hostClient,
		Business: mainlineClient,
		Topology: topo,
		Event:    watch.NewClient(watchDB, mongodb.Client(), redis.Client()),
	}
	return cache, nil
}

// ClientSet TODO
type ClientSet struct {
	Tree     *topotree.TopologyTree
	Topology *topology.Topology
	Host     *host.Client
	Business *mainline.Client
	Event    *watch.Client
}
