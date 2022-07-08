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

package operation

import (
	"strconv"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// GraphicsOperationInterface graphics operation methods
type GraphicsOperationInterface interface {
	// SelectObjectTopoGraphics select object topographics
	SelectObjectTopoGraphics(kit *rest.Kit, scopeType, scopeID string) ([]metadata.TopoGraphics, error)
	// UpdateObjectTopoGraphics update object topographics
	UpdateObjectTopoGraphics(kit *rest.Kit, scopeType, scopeID string, datas []metadata.TopoGraphics) error
}

// NewGraphics create a new graphics operation instance
func NewGraphics(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) GraphicsOperationInterface {
	return &graphics{
		clientSet:   client,
		authManager: authManager,
	}
}

// graphics graphics
type graphics struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// SelectObjectTopoGraphics select object topographics
func (g *graphics) SelectObjectTopoGraphics(kit *rest.Kit, scopeType, scopeID string) ([]metadata.TopoGraphics, error) {

	nodes := make([]metadata.TopoGraphics, 0)
	if scopeType != "global" {
		return nodes, nil
	}

	graphCondition := new(metadata.TopoGraphics)
	graphCondition.ScopeType = scopeType
	graphCondition.ScopeID = scopeID

	rsp, err := g.clientSet.CoreService().TopoGraphics().SearchTopoGraphics(kit.Ctx, kit.Header, graphCondition)
	if err != nil {
		blog.Errorf("search the graphics failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	graphNodes := make(map[string]*metadata.TopoGraphics, 0)
	for index, node := range rsp {
		graphNodes[node.NodeType+node.ObjID+strconv.Itoa(node.InstID)] = &rsp[index]
	}

	assts, err := g.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, nil)
	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	objAssts := make(map[string][]metadata.Association, 0)
	asstKindIDs := make([]string, 0)
	var associationKindMap = make(map[string]*metadata.AssociationKind, 0)
	for _, asst := range assts.Info {
		objAssts[asst.ObjectID] = append(objAssts[asst.ObjectID], asst)
		if _, ok := associationKindMap[asst.AsstKindID]; !ok {
			asstKindIDs = append(asstKindIDs, asst.AsstKindID)
			associationKindMap[asst.AsstKindID] = nil
		}
	}

	objs, err := g.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, nil)
	if err != nil {
		blog.ErrorJSON("find object failed, err: %v, rid: %v", err, kit.Rid)
		return nil, err
	}

	associationKindMap, err = g.findAssociationTypeByAsstKindID(kit, asstKindIDs)
	if err != nil {
		blog.Errorf("select object topo graphics failed, err: %v, kinds: %#v, rid: %s", err, asstKindIDs,
			kit.Rid)
		return nil, err
	}

	for _, obj := range objs.Info {
		node := g.genTopoNode(obj, kit.SupplierAccount, graphNodes)
		for _, asst := range objAssts[obj.ObjectID] {
			tmp := metadata.GraphAsst{
				NodeType:              "obj",
				ObjID:                 asst.AsstObjID,
				InstID:                asst.ID,
				AssociationKindInstID: associationKindMap[asst.AsstKindID].ID,
				Label:                 map[string]string{},
			}
			node.Assts = append(node.Assts, tmp)
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (g *graphics) findAssociationTypeByAsstKindID(kit *rest.Kit, asstKindIDs []string) (
	map[string]*metadata.AssociationKind, error) {

	cond := mapstr.MapStr{
		common.AssociationKindIDField: map[string][]string{
			common.BKDBIN: asstKindIDs,
		},
	}
	input := metadata.QueryCondition{
		Condition: cond,
	}

	resp, err := g.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("get association kind failed, err: %v, kinds: %#v rid: %s", err, asstKindIDs, kit.Rid)
		return nil, err
	}
	if resp.Count != len(asstKindIDs) {
		blog.Errorf("get association kind failed, err: %v, kinds: %#v rid: %s", resp, asstKindIDs,
			kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrTopoGetAssociationKindFailed, asstKindIDs)
	}

	var associationKindMap = make(map[string]*metadata.AssociationKind)
	for _, kind := range resp.Info {
		associationKindMap[kind.AssociationKindID] = kind
	}

	return associationKindMap, nil
}

func (g graphics) genTopoNode(obj metadata.Object, supplierAccount string,
	graphNodes map[string]*metadata.TopoGraphics) metadata.TopoGraphics {
	node := metadata.TopoGraphics{
		ScopeType:       "global",
		ScopeID:         "0",
		NodeType:        "obj",
		ObjID:           obj.ObjectID,
		IsPre:           obj.IsPre,
		NodeName:        obj.ObjectName,
		Icon:            obj.ObjIcon,
		SupplierAccount: supplierAccount,
	}

	oldNode := graphNodes[node.NodeType+node.ObjID+strconv.Itoa(node.InstID)]
	if oldNode != nil {
		node.SetPosition(oldNode.Position)
		node.SetExt(oldNode.Ext)
	} else {
		node.SetPosition(metadata.Position{})
		node.SetExt(map[string]interface{}{})
	}

	return node
}

// UpdateObjectTopoGraphics update object topographics
func (g *graphics) UpdateObjectTopoGraphics(kit *rest.Kit, scopeType, scopeID string,
	datas []metadata.TopoGraphics) error {

	for index := range datas {
		datas[index].SetScopeType(scopeType)
		datas[index].SetScopeID(scopeID)
	}

	err := g.clientSet.CoreService().TopoGraphics().UpdateTopoGraphics(kit.Ctx, kit.Header, datas)
	if err != nil {
		blog.Errorf("UpdateGraphics failed ,err: %v, datas: %#v, rid: %s", err, datas, kit.Rid)
		return err
	}

	return nil
}
