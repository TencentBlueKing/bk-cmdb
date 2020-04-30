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

package cache

import (
	"fmt"

	"configcenter/src/source_controller/coreservice/cache/business"
	"configcenter/src/source_controller/coreservice/cache/topo_tree"
	"configcenter/src/storage/dal"
	"configcenter/src/storage/reflector"
	"gopkg.in/redis.v5"
)

type Interface interface {
	SearchTopologyTree(opt *topo_tree.SearchOption) ([]*topo_tree.Topology, error)
}

func NewCache(rds *redis.Client, db dal.DB, event reflector.Interface) (Interface, error) {
	if err := business.NewCache(event, rds, db); err != nil {
		return nil, fmt.Errorf("new business cache failed, err: %v", err)
	}

	bizClient := business.NewClient(rds, db)

	cache := &cache{
		TopologyTree: topo_tree.NewTopologyTree(bizClient),
	}
	return cache, nil
}

type cache struct {
	*topo_tree.TopologyTree
}
