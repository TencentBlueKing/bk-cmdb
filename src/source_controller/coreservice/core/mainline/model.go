/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mainline

import (
	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
)

// SearchMainlineModelTopo get topo tree of model on mainline
func (m *topoManager) SearchMainlineModelTopo() (*metadata.TopoModelNode, error) {
	// TODO support withDetail option
	// step1: get all model associations
	mongoCondition := mongo.NewCondition()

	ctx := core.ContextParams{}
	ossociations := make([]mapstr.MapStr, 0)
	err := m.dbProxy.Table(common.BKTableNameObjAsst).Find(mongoCondition.ToMapStr()).All(ctx, ossociations)
	if err != nil {
		return nil, err
	}

	// step2: construct a tree fro associations
	var bizTopoModelNode *metadata.TopoModelNode
	topoNodeMap := map[string]*metadata.TopoModelNode{}
	for _, associaction := range ossociations {
		parentObjectID := associaction[common.AssociatedObjectIDField].(string)
		if _, exist := topoNodeMap[parentObjectID]; exist == false {
			topoModelNode := &metadata.TopoModelNode{ObjectID: parentObjectID}
			topoNodeMap[parentObjectID] = topoModelNode
		}

		parentNode := topoNodeMap[parentObjectID]

		// extract tree root node pointer
		if parentObjectID == "biz" {
			bizTopoModelNode = parentNode
		}

		childObjectID := associaction[common.BKObjIDField].(string)
		if _, exist := topoNodeMap[childObjectID]; exist == false {
			topoModelNode := &metadata.TopoModelNode{ObjectID: childObjectID}
			topoNodeMap[childObjectID] = topoModelNode
		}
		childNode, _ := topoNodeMap[childObjectID]
		topoNodeMap[parentObjectID].Children = append(topoNodeMap[parentObjectID].Children, childNode)
	}
	return bizTopoModelNode, nil
}
