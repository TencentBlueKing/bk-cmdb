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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/storage/driver/mongodb"
)

// ModelMainline TODO
type ModelMainline struct {
	root         *metadata.TopoModelNode
	associations []metadata.Association
}

// NewModelMainline TODO
func NewModelMainline() (*ModelMainline, error) {
	modelMainline := &ModelMainline{}
	modelMainline.associations = make([]metadata.Association, 0)
	return modelMainline, nil
}

func (mm *ModelMainline) loadMainlineAssociations(kit *rest.Kit) error {
	filter := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	err := mongodb.Shard(kit.ShardOpts()).Table(common.BKTableNameObjAsst).Find(filter).All(kit.Ctx, &mm.associations)
	if err != nil {
		blog.Errorf("query topo model mainline association from db failed, %+v, rid: %s", err, kit.Rid)
		return fmt.Errorf("query topo model mainline association from db failed, %+v", err)
	}
	blog.V(5).Infof("get topo model mainline associations result: %+v, rid: %s", mm.associations, kit.Rid)
	return nil
}

func (mm *ModelMainline) constructTopoTree(kit *rest.Kit) error {
	// step2: construct a tree fro associations
	topoModelNodeMap := map[string]*metadata.TopoModelNode{}
	for _, association := range mm.associations {
		blog.V(5).Infof("association: %+v, rid: %s", association, kit.Rid)
		parentObjectID := association.AsstObjID
		if _, exist := topoModelNodeMap[parentObjectID]; !exist {
			topoModelNodeMap[parentObjectID] = &metadata.TopoModelNode{
				ObjectID: parentObjectID,
				Children: []*metadata.TopoModelNode{},
			}
		}

		parentTopoModelNode := topoModelNodeMap[parentObjectID]

		// extract tree root node pointer
		if parentObjectID == common.BKInnerObjIDApp {
			mm.root = parentTopoModelNode
		}

		childObjectID := association.ObjectID
		if _, exist := topoModelNodeMap[childObjectID]; !exist {
			topoModelNodeMap[childObjectID] = &metadata.TopoModelNode{
				ObjectID: childObjectID,
				Children: []*metadata.TopoModelNode{},
			}
		}
		parentTopoModelNode.Children = append(parentTopoModelNode.Children, topoModelNodeMap[childObjectID])
	}
	return nil
}

// GetRoot TODO
func (mm *ModelMainline) GetRoot(kit *rest.Kit, withDetail bool) (*metadata.TopoModelNode,
	error) {
	if err := mm.loadMainlineAssociations(kit); err != nil {
		blog.Errorf("load model mainline associations failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, fmt.Errorf("load model mainline associations failed, err: %+v", err)
	}

	if err := mm.constructTopoTree(kit); err != nil {
		blog.Errorf("construct tree from model mainline associations failed, err: %+v, rid: %s", err, kit.Rid)
		return nil, fmt.Errorf("construct tree from model mainline associations failed, err: %+v", err)
	}
	if withDetail {
		// thinking what's detail actually
		panic("detail option not implemented yet.")
	}
	return mm.root, nil
}
