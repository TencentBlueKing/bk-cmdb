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

	"configcenter/src/common/blog"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// AssociationOperationInterface association operation methods
type AssociationOperationInterface interface {
	CreateMainlineAssociation(params types.ContextParams, data *metadata.Association) (model.Association, error)
	DeleteMainlineAssociaton(params types.ContextParams, objID string) error
	SearchMainlineAssociationTopo(params types.ContextParams, targetObj model.Object) ([]*metadata.MainlineObjectTopo, error)
	SearchMainlineAssociationInstTopo(params types.ContextParams, bizID int64) ([]*metadata.TopoInstRst, error)
	CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (model.Association, error)
	DeleteAssociation(params types.ContextParams, cond condition.Condition) error
	UpdateAssociation(params types.ContextParams, data frtypes.MapStr, cond condition.Condition) error
	SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error)
	SetProxy(cls ClassificationOperationInterface, obj ObjectOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory, targetInst inst.Factory)
}

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface) AssociationOperationInterface {
	return &association{
		clientSet: client,
	}
}

type association struct {
	clientSet    apimachinery.ClientSetInterface
	cls          ClassificationOperationInterface
	obj          ObjectOperationInterface
	attr         AttributeOperationInterface
	inst         InstOperationInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

func (a *association) SetProxy(cls ClassificationOperationInterface, obj ObjectOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory, targetInst inst.Factory) {
	a.cls = cls
	a.obj = obj
	a.attr = attr
	a.inst = inst
	a.modelFactory = targetModel
	a.instFactory = targetInst
}

func (a *association) SearchObjectAssociation(params types.ContextParams, objID string) ([]metadata.Association, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(common.BKObjIDField).Eq(objID)
	rsp, err := a.clientSet.ObjectController().Meta().SelectObjectAssociations(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-asst] failed to request object controller, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if !rsp.Result {
		blog.Errorf("[operation-asst] failed to search the object(%s) association info , error info is %s", objID, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return rsp.Data, nil
}

func (a *association) CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (model.Association, error) {

	//  check the association
	//	cond := condition.CreateCondition()
	//	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(data.AsstObjID)
	//	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(data.ObjectAttID)

	//asst := a.modelFactory.(params)
	//asst.Parse(data)

	//a.clientSet.ObjectController().Meta().SelectObjectAssociations()

	return nil, nil
}
func (a *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {
	return nil
}
func (a *association) UpdateAssociation(params types.ContextParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
