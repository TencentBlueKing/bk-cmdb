/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

// Package topo defines the topology cache related logics
package topo

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/level"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/tree"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	"configcenter/src/storage/driver/mongodb"
)

// GenBizTopo generate business topology tree
func GenBizTopo(kit *rest.Kit, bizID int64, topoType types.TopoType, byCache bool) (any, error) {
	// read from secondary node in mongodb cluster
	kit.Ctx = util.SetDBReadPreference(kit.Ctx, common.SecondaryPreferredMode)

	// get biz info
	filter := mapstr.MapStr{
		common.BKAppIDField: bizID,
	}
	biz := new(metadata.BizInst)
	if err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameBaseApp).Find(filter).Fields(common.BKAppIDField,
		common.BKAppNameField, common.BKDefaultField).One(kit.Ctx, biz); err != nil {
		blog.Errorf("get biz %d failed, err: %v, rid: %s", bizID, err, kit.Rid)
		return nil, err
	}

	// get topology generator
	topology, err := GetTopology(kit, topoType)
	if err != nil {
		blog.Errorf("get %s topology generator failed, err: %v, rid: %v", topoType, err, kit.Rid)
		return nil, err
	}

	// get topology nodes and generate biz topology tree
	var nodes []types.Node
	if byCache {
		nodes, err = topology.TopLevel.GetNodesByCache(kit, bizID)
	} else {
		nodes, err = topology.TopLevel.GetNodesByDB(kit, bizID, nil)
	}
	if err != nil {
		blog.Errorf("get biz %d %s topo nodes failed, by cache: %v, err: %v, rid: %s", biz.BizID, topoType, byCache,
			err, kit.Rid)
		return nil, err
	}

	bizTopo := &types.BizTopo{
		Biz: &types.BizInfo{
			ID:   biz.BizID,
			Name: biz.BizName,
		},
		Nodes: nodes,
	}

	topo, err := topology.Tree.RearrangeBizTopo(kit, biz, nodes)
	if err != nil {
		blog.Errorf("rearrange biz %d %s topo failed, err: %v, topo: %+v, rid: %s", biz, topoType, err, bizTopo,
			kit.Rid)
		return nil, err
	}

	return topo, nil
}

// Topology defines the topology generator
type Topology struct {
	// Tree is the topology tree generator
	Tree tree.TreeI
	// TopLevel is the topology tree's top level generator
	TopLevel level.LevelI
}

// TopologyGetter defines the function to get topology generator
type TopologyGetter func(kit *rest.Kit) (*Topology, error)

// topoGetterMap is the mapping of topology type to TopologyGetter
var topoGetterMap = map[types.TopoType]TopologyGetter{}

// GetTopology get topology generator
func GetTopology(kit *rest.Kit, topoType types.TopoType) (*Topology, error) {
	getter, exists := topoGetterMap[topoType]
	if !exists {
		blog.Errorf("%s topology getter not exists, rid: %v", topoType, kit.Rid)
		return nil, fmt.Errorf("topology type %s is invalid", topoType)
	}

	topology, err := getter(kit)
	if err != nil {
		blog.Errorf("get %s topology generator failed, rid: %v", topoType, kit.Rid)
		return nil, err
	}

	return topology, nil
}
