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
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/source_controller/coreservice/core"
	"fmt"
)

// SearchMainlineModelTopo get topo tree of model on mainline
func (m *topoManager) SearchMainlineModelTopo() (*metadata.TopoModelNode, error) {
	// TODO support withDetail option
	// step1: get all model associations
	mongoCondition := mongo.NewCondition()
	mongoCondition.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: common.AssociationKindMainline})

	ctx := core.ContextParams{}
	associations := make([]mapstr.MapStr, 0)
	err := m.dbProxy.Table(common.BKTableNameObjAsst).Find(mongoCondition.ToMapStr()).All(ctx, &associations)
	if err != nil {
		blog.Errorf("query topo model mainline association from db failed, %+v", err)
		return nil, fmt.Errorf("query topo model mainline association from db failed, %+v", err)
	}
	blog.V(2).Infof("get topo model mainline associations result: %+v", associations)

	// step2: construct a tree fro associations
	var bizTopoModelNode *metadata.TopoModelNode
	topoModelNodeMap := map[string]*metadata.TopoModelNode{}
	for _, association := range associations {
		blog.V(5).Infof("association: %+v", association)
		parentObjectID := association[common.AssociatedObjectIDField].(string)
		if _, exist := topoModelNodeMap[parentObjectID]; exist == false {
			topoModelNodeMap[parentObjectID] = &metadata.TopoModelNode{
				ObjectID: parentObjectID,
				Children: []*metadata.TopoModelNode{},
			}
		}

		parentTopoModelNode := topoModelNodeMap[parentObjectID]

		// extract tree root node pointer
		if parentObjectID == common.BKInnerObjIDApp {
			bizTopoModelNode = parentTopoModelNode
		}

		childObjectID := association[common.BKObjIDField].(string)
		if _, exist := topoModelNodeMap[childObjectID]; exist == false {
			topoModelNodeMap[childObjectID] = &metadata.TopoModelNode{
				ObjectID: childObjectID,
				Children: []*metadata.TopoModelNode{},
			}
		}
		parentTopoModelNode.Children = append(parentTopoModelNode.Children, topoModelNodeMap[childObjectID])
	}
	blog.V(2).Infof("bizTopoModelNode: %+v", bizTopoModelNode)
	return bizTopoModelNode, nil
}
