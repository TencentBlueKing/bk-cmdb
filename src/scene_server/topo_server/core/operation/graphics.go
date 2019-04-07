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
	"context"
	"strconv"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/auth/meta"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/types"
)

type GraphicsOperationInterface interface {
	SelectObjectTopoGraphics(params types.ContextParams, scopeType, scopeID string) ([]metadata.TopoGraphics, error)
	UpdateObjectTopoGraphics(params types.ContextParams, scopeType, scopeID string, datas []metadata.TopoGraphics) error

	SetProxy(obj ObjectOperationInterface, asst AssociationOperationInterface)
}

func NewGraphics(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) GraphicsOperationInterface {
	return &graphics{
		clientSet:   client,
		authManager: authManager,
	}
}

type graphics struct {
	clientSet   apimachinery.ClientSetInterface
	obj         ObjectOperationInterface
	asst        AssociationOperationInterface
	authManager *extensions.AuthManager
}

func (g *graphics) SetProxy(obj ObjectOperationInterface, asst AssociationOperationInterface) {
	g.obj = obj
	g.asst = asst
}

func (g *graphics) SelectObjectTopoGraphics(params types.ContextParams, scopeType, scopeID string) ([]metadata.TopoGraphics, error) {

	graphcondition := &metadata.TopoGraphics{}
	graphcondition.SetScopeType(scopeType)
	graphcondition.SetScopeID(scopeID)
	if nil != params.MetaData {
		graphcondition.SetMetaData(*params.MetaData)
	}
	rsp, err := g.clientSet.ObjectController().Meta().SearchTopoGraphics(context.Background(), params.Header, graphcondition)
	if nil != err {
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[graphics] failed to search the graphics , error info is %s", rsp.ErrMsg)
		return nil, params.Err.New(common.CCErrTopoGraphicsSearchFailed, rsp.ErrMsg)
	}

	dbnodes := rsp.Data

	graphnodes := map[string]*metadata.TopoGraphics{}
	for index, node := range dbnodes {
		graphnodes[node.NodeType+node.ObjID+strconv.Itoa(node.InstID)] = &dbnodes[index]
	}

	nodes := []metadata.TopoGraphics{}
	if scopeType == "global" {

		objs, err := g.obj.FindObject(params, condition.CreateCondition())
		if err != nil {
			blog.Errorf("SelectObject failed %v", err.Error())
			return nil, params.Err.New(common.CCErrTopoGraphicsSearchFailed, err.Error())
		}

		assts, err := g.asst.SearchObjectAssociation(params, "")
		if err != nil {
			blog.Errorf("SelectObjectAsst failed %v", err.Error())
			return nil, params.Err.New(common.CCErrTopoGraphicsSearchFailed, err.Error())
		}

		objAssts := map[string][]metadata.Association{}
		for _, asst := range assts {
			objAssts[asst.ObjectID] = append(objAssts[asst.ObjectID], asst)
		}

		for _, obj := range objs {
			object := obj.Object()
			node := metadata.TopoGraphics{}
			node.SetNodeType("obj")
			node.SetObjID(object.ObjectID)
			node.SetInstID(0)
			node.SetNodeName(object.ObjectName)
			node.SetScopeType("global")
			node.SetScopeID("0")
			node.SetSupplierAccount("0")
			node.SetIsPre(object.IsPre)
			node.SetIcon(object.ObjIcon)

			oldnode := graphnodes[node.NodeType+node.ObjID+strconv.Itoa(node.InstID)]
			if oldnode != nil {
				node.SetPosition(oldnode.Position)
				node.SetExt(oldnode.Ext)
			} else {
				node.SetPosition(metadata.Position{})
				node.SetExt(map[string]interface{}{})
			}

			for _, asst := range objAssts[object.ObjectID] {

				typeCond := condition.CreateCondition()
				typeCond.Field(common.AssociationKindIDField).Eq(asst.AsstKindID)
				request := &metadata.SearchAssociationTypeRequest{
					Condition: typeCond.ToMapStr(),
				}

				resp, err := g.asst.SearchType(params, request)
				if err != nil {
					blog.Errorf("select object topo graph failed, because get association kind[%s] failed, err: %v", asst.AsstKindID, err)
					return nil, params.Err.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
				}
				if !resp.Result {
					blog.Errorf("select object topo graph failed, because get association kind[%s] failed, err: %v", asst.AsstKindID, resp.ErrMsg)
					return nil, params.Err.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
				}

				// should only be one association kind.
				if len(resp.Data.Info) == 0 {
					blog.Errorf("select object topo graph failed, because get association kind[%s] failed, err: can not find this association kind.", asst.AsstKindID)
					return nil, params.Err.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
				}

				node.Assts = append(node.Assts, metadata.GraphAsst{
					AsstType:              "",
					NodeType:              "obj",
					ObjID:                 asst.AsstObjID,
					InstID:                asst.ID,
					AssociationKindInstID: resp.Data.Info[0].ID,
					Label:                 map[string]string{},
				})
			}
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func (g *graphics) UpdateObjectTopoGraphics(params types.ContextParams, scopeType, scopeID string, datas []metadata.TopoGraphics) error {

	for index := range datas {
		datas[index].SetScopeType(scopeType)
		datas[index].SetScopeID(scopeID)
		if nil != params.MetaData {
			datas[index].SetMetaData(*params.MetaData)
		}
	}

	// check authorization
	objectIDs := make([]string, 0)
	for _, data := range datas {
		objectIDs = append(objectIDs, data.ObjID)
	}
	if err := g.authManager.AuthorizeByObjectID(params.Context, params.Header, meta.ModelTopologyOperation, objectIDs...); err != nil {
		blog.Errorf("UpdateGraphics failed, authorization failed, err: %+v", err)
		return params.Err.Error(common.CCErrCommAuthorizeFailed)
	}

	rsp, err := g.clientSet.ObjectController().Meta().UpdateTopoGraphics(context.Background(), params.Header, datas)
	if err != nil {
		blog.Errorf("UpdateGraphics failed %v", err.Error())
		return params.Err.New(common.CCErrTopoGraphicsUpdateFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[graphics] failed to update the graphics, error info is %s", rsp.ErrMsg)
		return params.Err.New(common.CCErrTopoGraphicsUpdateFailed, rsp.ErrMsg)
	}

	return nil
}
