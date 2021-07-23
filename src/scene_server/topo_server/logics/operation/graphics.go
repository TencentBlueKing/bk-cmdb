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
func NewGraphics(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) GraphicsOperationInterface {
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

	graphCondition := &metadata.TopoGraphics{
		ScopeType: scopeType,
		ScopeID:   scopeID,
	}

	rsp, err := g.clientSet.CoreService().TopoGraphics().SearchTopoGraphics(kit.Ctx, kit.Header, graphCondition)
	if nil != err {
		blog.Errorf("search the graphics failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}
	if err = rsp.CCError(); err != nil {
		blog.Errorf("search the graphics failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	dbNodes := rsp.Data
	graphNodes := map[string]*metadata.TopoGraphics{}
	for index, node := range dbNodes {
		graphNodes[node.NodeType+node.ObjID+strconv.Itoa(node.InstID)] = &dbNodes[index]
	}

	// TODO 调用 asst 中 searchObjectAssociation
	assts, err := g.searchObjectAssociation(kit, "")
	if err != nil {
		blog.Errorf("select object asst failed, err: %s, rid: %v", err, kit.Rid)
		return nil, err
	}

	objAssts := make(map[string][]metadata.Association, 0)
	for _, asst := range assts {
		objAssts[asst.ObjectID] = append(objAssts[asst.ObjectID], asst)
	}

	// TODO obj 中的 FindObject
	objs, err := g.findObject(kit, nil)
	if err != nil {
		blog.Errorf("SelectObject failed, err: %s, rid: %v", err, kit.Rid)
		return nil, err
	}

	asstKindIDs := make([]string, 0)
	var associationKindMap = make(map[string]*metadata.AssociationKind, 0)
	for _, obj := range objs {
		for _, asst := range objAssts[obj.ObjectID] {
			if _, ok := associationKindMap[asst.AsstKindID]; !ok {
				asstKindIDs = append(asstKindIDs, asst.AsstKindID)
				associationKindMap[asst.AsstKindID] = nil
			}
		}
	}

	associationKindMap, err = g.findAssociationTypeByAsstKindID(kit, asstKindIDs)
	if err != nil {
		blog.ErrorJSON("select object topo graphics failed, err: %v, kinds: %#v, rid: %s", err, asstKindIDs,
			kit.Rid)
		return nil, err
	}

	for _, obj := range objs {
		node := g.genTopoNode(obj, kit.SupplierAccount, graphNodes)
		for _, asst := range objAssts[obj.ObjectID] {
			node.Assts = append(node.Assts, metadata.GraphAsst{
				NodeType:              "obj",
				ObjID:                 asst.AsstObjID,
				InstID:                asst.ID,
				AssociationKindInstID: associationKindMap[asst.AsstKindID].ID,
				Label:                 map[string]string{},
			})
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func (g *graphics) findAssociationTypeByAsstKindID(kit *rest.Kit, asstKindIDs []string) (
	map[string]*metadata.AssociationKind, error) {

	typeCond := mapstr.MapStr{
		common.AssociationKindIDField: map[string][]string{
			common.BKDBIN: asstKindIDs,
		},
	}
	// TODO 调用asst中SearchType
	resp, err := g.searchType(kit, &metadata.SearchAssociationTypeRequest{
		Condition: typeCond,
	})
	if err != nil {
		blog.ErrorJSON("get association kind failed, err: %v, kinds: %#v rid: %s", err, asstKindIDs, kit.Rid)
		return nil, err
	}
	if err = resp.CCError(); err != nil {
		blog.ErrorJSON("get association kind failed, err: %v, kinds: %#v rid: %s", err, asstKindIDs, kit.Rid)
		return nil, err
	}
	if len(resp.Data.Info) != len(asstKindIDs) {
		blog.ErrorJSON("get association kind failed, err: %v, kinds: %#v rid: %s", resp.ErrMsg, asstKindIDs,
			kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrTopoGetAssociationKindFailed, asstKindIDs)
	}

	var associationKindMap = make(map[string]*metadata.AssociationKind)
	for _, kind := range resp.Data.Info {
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
		InstID:          0,
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

	rsp, err := g.clientSet.CoreService().TopoGraphics().UpdateTopoGraphics(kit.Ctx, kit.Header, datas)
	if err != nil {
		blog.ErrorJSON("UpdateGraphics failed ,err: %v, datas: %#v, rid: %v", err, datas, kit.Rid)
		return err
	}
	if err = rsp.CCError(); err != nil {
		blog.ErrorJSON("UpdateGraphics failed ,err: %v, datas: %#v, rid: %v", err, datas, kit.Rid)
		return err
	}

	return nil
}

// TODO 暂时使用的函数
func (g *graphics) findObject(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Object, error) {
	rsp, err := g.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if nil != err {
		blog.ErrorJSON("find object failed, err: %v, cond: %#v, rid: %v", err, cond, kit.Rid)
		return nil, err
	}

	if err = rsp.CCError(); err != nil {
		blog.ErrorJSON("find object failed, err: %v, cond: %#v, rid: %v", err, cond, kit.Rid)
		return nil, err
	}

	return rsp.Data.Info, nil
}

func (g *graphics) searchObjectAssociation(kit *rest.Kit, objID string) ([]metadata.Association, error) {
	cond := mapstr.MapStr{}
	if len(objID) != 0 {
		cond.Set(common.BKObjIDField, objID)
	}

	rsp, err := g.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond})
	if nil != err {
		blog.ErrorJSON("search object association failed, err: %v, input: %#v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	if err = rsp.CCError(); err != nil {
		blog.ErrorJSON("search object association failed, err: %v, input: %#v, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	return rsp.Data.Info, nil
}

func (g *graphics) searchType(kit *rest.Kit, request *metadata.SearchAssociationTypeRequest) (resp *metadata.
	SearchAssociationTypeResult, err error) {
	input := metadata.QueryCondition{
		Condition: request.Condition,
		Page:      metadata.BasePage{Limit: request.Limit, Start: request.Start, Sort: request.Sort},
	}

	return g.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, &input)
}
