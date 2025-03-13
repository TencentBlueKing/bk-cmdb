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

package level

import (
	"sort"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	nlgc "configcenter/src/source_controller/cacheservice/cache/biz-topo/logics/node"
	"configcenter/src/source_controller/cacheservice/cache/biz-topo/types"
	"configcenter/src/storage/driver/mongodb"
)

var briefTopLevel = newBriefTopoLevel()

// GetBriefTopLevel get the top level of brief biz topology
func GetBriefTopLevel() LevelI {
	return briefTopLevel
}

type briefTopoLevel struct{}

func newBriefTopoLevel() *briefTopoLevel {
	return &briefTopoLevel{}
}

// GetNodesByDB get all nodes that belongs to the topology level
func (l *briefTopoLevel) GetNodesByDB(kit *rest.Kit, bizID int64, _ []mapstr.MapStr) ([]types.Node, error) {
	parentObjMap, err := l.getMainlineObjMap(kit)
	if err != nil {
		return nil, err
	}

	topNodes := make([]types.Node, 0)
	prevNodeMap := make(map[int64][]types.Node)

	for objID := common.BKInnerObjIDModule; objID != common.BKInnerObjIDApp; objID = parentObjMap[objID] {
		nodes, err := l.getBriefTopoNodesByObj(kit, bizID, objID)
		if err != nil {
			return nil, err
		}

		if len(nodes) == 0 {
			prevNodeMap = make(map[int64][]types.Node)
			continue
		}

		for i, node := range nodes {
			l.sortNodes(prevNodeMap[node.ID])
			nodes[i].SubNodes = prevNodeMap[node.ID]
		}

		if objID == common.BKInnerObjIDSet {
			normalNodes := make([]types.Node, 0)
			for i, node := range nodes {
				defaultVal, err := util.GetIntByInterface(node.Extra)
				if err != nil {
					blog.Errorf("parse brief set node(%+v) failed, err: %v, rid: %s", node, err, kit.Rid)
					return nil, err
				}
				if defaultVal == common.DefaultResSetFlag {
					topNodes = append(topNodes, nodes[i])
					continue
				}
				normalNodes = append(normalNodes, nodes[i])
			}
			nodes = normalNodes
		}

		if parentObjMap[objID] == common.BKInnerObjIDApp {
			topNodes = append(topNodes, nodes...)
			break
		}

		prevNodeMap = make(map[int64][]types.Node)
		for _, node := range nodes {
			prevNodeMap[node.ParentID] = append(prevNodeMap[node.ParentID], node)
		}
	}

	l.sortNodes(topNodes)
	return topNodes, nil
}

// GetNodesByCache get all nodes that belongs to the brief biz topo level
func (l *briefTopoLevel) GetNodesByCache(kit *rest.Kit, bizID int64) ([]types.Node, error) {
	return l.GetNodesByDB(kit, bizID, nil)
}

// getMainlineObjMap get mainline object to parent object map
func (l *briefTopoLevel) getMainlineObjMap(kit *rest.Kit) (map[string]string, error) {
	relations := make([]metadata.Association, 0)
	filter := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(filter).Fields(common.BKObjIDField,
		common.BKAsstObjIDField).All(kit.Ctx, &relations)
	if err != nil {
		blog.Errorf("get mainline topology association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	mainlineMap := make(map[string]string)
	for _, relation := range relations {
		if relation.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineMap[relation.ObjectID] = relation.AsstObjID
	}

	return mainlineMap, nil
}

func (l *briefTopoLevel) getBriefTopoNodesByObj(kit *rest.Kit, bizID int64, objID string) ([]types.Node, error) {
	tableName := common.GetInstTableName(objID, kit.TenantID)
	cond := mapstr.MapStr{common.BKAppIDField: bizID}

	idField := common.GetInstIDField(objID)
	nameField := common.GetInstNameField(objID)
	parentField := common.BKParentIDField
	fields := []string{idField, nameField}
	switch objID {
	case common.BKInnerObjIDSet:
		fields = append(fields, common.BKDefaultField)
	case common.BKInnerObjIDModule:
		fields = append(fields, common.BKDefaultField)
		parentField = common.BKSetIDField
	}
	fields = append(fields, parentField)

	nodes, err := nlgc.PagedGetNodes(kit, tableName, cond, fields, l.nodeParser(objID, idField, nameField, parentField))
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (l *briefTopoLevel) nodeParser(objID, idField, nameField, parentIDField string) nlgc.NodeParser {
	return func(kit *rest.Kit, data []mapstr.MapStr) ([]types.Node, error) {
		nodes := make([]types.Node, len(data))
		for i, item := range data {
			id, err := util.GetInt64ByInterface(item[idField])
			if err != nil {
				blog.Errorf("parse %s brief node id failed, err: %v, item: %+v, rid: %s", objID, err, item, kit.Rid)
				return nil, err
			}

			parentID, err := util.GetInt64ByInterface(item[parentIDField])
			if err != nil {
				blog.Errorf("parse %s brief node parent id failed, err: %v, item: %+v, rid: %s", objID, err, item,
					kit.Rid)
				return nil, err
			}

			nodes[i] = types.Node{
				Kind:     objID,
				ID:       id,
				Name:     util.GetStrByInterface(item[nameField]),
				ParentID: parentID,
			}

			switch objID {
			case common.BKInnerObjIDSet, common.BKInnerObjIDModule:
				defaultVal, err := util.GetInt64ByInterface(item[common.BKDefaultField])
				if err != nil {
					blog.Errorf("parse %s brief node default value failed, err: %v, item: %+v, rid: %s", objID, err,
						item, kit.Rid)
					return nil, err
				}
				nodes[i].Extra = defaultVal
			}
		}

		return nodes, nil
	}
}

// sortNodes sort nodes by name
func (l *briefTopoLevel) sortNodes(nodes []types.Node) {
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Name < nodes[j].Name
	})
}
