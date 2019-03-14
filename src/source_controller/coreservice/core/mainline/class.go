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
    "configcenter/src/common/mapstr"
    "configcenter/src/common/metadata"
    "configcenter/src/common/universalsql/mongo"
    "configcenter/src/source_controller/coreservice/core"
    "configcenter/src/storage/dal"
)

type ModelMainline struct {
    root *metadata.TopoModelNode
    dbProxy dal.RDB
    associations []mapstr.MapStr
}

func NewModelMainline(proxy dal.RDB) (*ModelMainline, error) {
    modelMainline := &ModelMainline{dbProxy: proxy}
    modelMainline.associations = make([]mapstr.MapStr, 0)
    return modelMainline, nil
}

func (mm *ModelMainline) loadMainlineAssociations() error  {
    mongoCondition := mongo.NewCondition()
    mongoCondition.Element(&mongo.Eq{Key: common.AssociationKindIDField, Val: common.AssociationKindMainline})

    ctx := core.ContextParams{}
    err := mm.dbProxy.Table(common.BKTableNameObjAsst).Find(mongoCondition.ToMapStr()).All(ctx, &mm.associations)
    if err != nil {
        blog.Errorf("query topo model mainline association from db failed, %+v", err)
        return fmt.Errorf("query topo model mainline association from db failed, %+v", err)
    }
    blog.V(5).Infof("get topo model mainline associations result: %+v", mm.associations)
    return nil
}

func (mm *ModelMainline) constructTopoTree() error  {
    // step2: construct a tree fro associations
    topoModelNodeMap := map[string]*metadata.TopoModelNode{}
    for _, association := range mm.associations {
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
            mm.root = parentTopoModelNode
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
    blog.V(2).Infof("bizTopoModelNode: %+v", mm.root)
    return nil
}

func (mm *ModelMainline) GetRoot(withDetail bool) (*metadata.TopoModelNode, error) {
   if err := mm.loadMainlineAssociations(); err != nil {
       blog.Errorf("get topo model failed, load model mainline associations failed, err: %+v", err)
       return nil, fmt.Errorf("get topo model failed, load model mainline associations failed, err: %+v", err)
   }
   
   if err := mm.constructTopoTree(); err != nil {
       blog.Errorf("get topo model failed, construct tree from model mainline associations failed, err: %+v", err)
       return nil, fmt.Errorf("get topo model failed, construct tree from model mainline associations failed, err: %+v", err)
   }
   if withDetail == true {
       // thinking what's detail actually
       panic("detail option not implemented yet.")
   }
   return mm.root, nil
}
