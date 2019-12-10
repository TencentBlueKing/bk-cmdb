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
	"context"
	"fmt"
	"net/http"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

type ModelMainline struct {
	root         *metadata.TopoModelNode
	dbProxy      dal.RDB
	associations []metadata.Association
}

func NewModelMainline(proxy dal.RDB) (*ModelMainline, error) {
	modelMainline := &ModelMainline{dbProxy: proxy}
	modelMainline.associations = make([]metadata.Association, 0)
	return modelMainline, nil
}

func (mm *ModelMainline) loadMainlineAssociations(ctx context.Context, header http.Header) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	filter := map[string]interface{}{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BkSupplierAccount:      util.GetOwnerID(header),
	}
	err := mm.dbProxy.Table(common.BKTableNameObjAsst).Find(filter).All(ctx, &mm.associations)
	if err != nil {
		blog.Errorf("query topo model mainline association from db failed, %+v, rid: %s", err, rid)
		return fmt.Errorf("query topo model mainline association from db failed, %+v", err)
	}
	blog.V(5).Infof("get topo model mainline associations result: %+v, rid: %s", mm.associations, rid)
	return nil
}

func (mm *ModelMainline) constructTopoTree(ctx context.Context) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	// step2: construct a tree fro associations
	topoModelNodeMap := map[string]*metadata.TopoModelNode{}
	for _, association := range mm.associations {
		blog.V(5).Infof("association: %+v, rid: %s", association, rid)
		parentObjectID := association.AsstObjID
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

		childObjectID := association.ObjectID
		if _, exist := topoModelNodeMap[childObjectID]; exist == false {
			topoModelNodeMap[childObjectID] = &metadata.TopoModelNode{
				ObjectID: childObjectID,
				Children: []*metadata.TopoModelNode{},
			}
		}
		parentTopoModelNode.Children = append(parentTopoModelNode.Children, topoModelNodeMap[childObjectID])
	}
	blog.V(2).Infof("bizTopoModelNode: %+v, rid: %s", mm.root, rid)
	return nil
}

func (mm *ModelMainline) GetRoot(ctx context.Context, header http.Header, withDetail bool) (*metadata.TopoModelNode, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	if err := mm.loadMainlineAssociations(ctx, header); err != nil {
		blog.Errorf("get topo model failed, load model mainline associations failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get topo model failed, load model mainline associations failed, err: %+v", err)
	}

	if err := mm.constructTopoTree(ctx); err != nil {
		blog.Errorf("get topo model failed, construct tree from model mainline associations failed, err: %+v, rid: %s", err, rid)
		return nil, fmt.Errorf("get topo model failed, construct tree from model mainline associations failed, err: %+v", err)
	}
	if withDetail == true {
		// thinking what's detail actually
		panic("detail option not implemented yet.")
	}
	return mm.root, nil
}
