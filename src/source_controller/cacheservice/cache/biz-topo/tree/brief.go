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

package tree

import (
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
)

// BriefTopoTree defines brief biz topology tree type
type BriefTopoTree struct{}

// RearrangeBizTopo rearrange brief business topology tree
func (t *BriefTopoTree) RearrangeBizTopo(kit *rest.Kit, biz *metadata.BizInst, nodes []types.Node) (any, error) {
	parsedNodes, idleSets := make([]*types.BriefNode, 0), make([]*types.BriefNode, 0)
	for _, node := range nodes {
		parsedNode, isIdle, err := parseBriefNode(kit, &node)
		if err != nil {
			return nil, err
		}

		if isIdle {
			idleSets = append(idleSets, parsedNode)
			continue
		}
		parsedNodes = append(parsedNodes, parsedNode)
	}

	return &types.BizBriefTopology{
		Biz: &types.BriefBizInfo{
			ID:      biz.BizID,
			Name:    biz.BizName,
			Default: biz.Default,
		},
		Idle:  idleSets,
		Nodes: parsedNodes,
	}, nil
}

func parseBriefNode(kit *rest.Kit, node *types.Node) (*types.BriefNode, bool, error) {
	parsedNode := &types.BriefNode{
		Object: node.Kind,
		ID:     node.ID,
		Name:   node.Name,
	}

	isIdle := false

	switch node.Kind {
	case common.BKInnerObjIDSet:
		defaultVal, err := util.GetIntByInterface(node.Extra)
		if err != nil {
			blog.Errorf("parse brief set node(%+v) failed, err: %v, rid: %s", node, err, kit.Rid)
			return nil, false, err
		}
		if defaultVal == common.DefaultResSetFlag {
			isIdle = true
		}
		parsedNode.Default = &defaultVal

	case common.BKInnerObjIDModule:
		defaultVal, err := util.GetIntByInterface(node.Extra)
		if err != nil {
			blog.Errorf("parse brief set node(%+v) failed, err: %v, rid: %s", node, err, kit.Rid)
			return nil, false, err
		}
		parsedNode.Default = &defaultVal
	}

	for _, subNode := range node.SubNodes {
		parsedSubNode, _, err := parseBriefNode(kit, &subNode)
		if err != nil {
			return nil, false, err
		}
		parsedNode.SubNodes = append(parsedNode.SubNodes, parsedSubNode)
	}

	return parsedNode, isIdle, nil
}
