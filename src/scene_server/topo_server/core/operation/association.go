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
	"configcenter/src/apimachinery"

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

// NewAssociationOperation create a new association operation instance
func NewAssociationOperation(client apimachinery.ClientSetInterface, cls ClassificationOperationInterface, obj ObjectOperationInterface, attr AttributeOperationInterface, inst InstOperationInterface, targetModel model.Factory, targetInst inst.Factory) AssociationOperationInterface {
	return &association{
		clientSet:    client,
		cls:          cls,
		obj:          obj,
		attr:         attr,
		inst:         inst,
		modelFactory: targetModel,
		instFactory:  targetInst,
	}
}

func (cli *association) CreateCommonAssociation(params types.ContextParams, data *metadata.Association) (model.Association, error) {

	//  check the association
	//	cond := condition.CreateCondition()
	//	cond.Field(metadata.AssociationFieldAssociationObjectID).Eq(data.AsstObjID)
	//	cond.Field(metadata.AssociationFieldObjectAttributeID).Eq(data.ObjectAttID)

	//asst := cli.modelFactory.(params)
	//asst.Parse(data)

	//cli.clientSet.ObjectController().Meta().SelectObjectAssociations()

	return nil, nil
}
func (cli *association) DeleteAssociation(params types.ContextParams, cond condition.Condition) error {
	return nil
}
func (cli *association) UpdateAssociation(params types.ContextParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
