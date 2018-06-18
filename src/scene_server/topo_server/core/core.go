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

package core

import (
	"configcenter/src/apimachinery"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/operation"
)

// Core Provides management interfaces for models and instances
type Core interface {
	AssociationOperation() operation.AssociationOperationInterface
	AttributeOperation() operation.AttributeOperationInterface
	ClassificationOperation() operation.ClassificationOperationInterface
	GroupOperation() operation.GroupOperationInterface
	InstOperation() operation.InstOperationInterface
	ObjectOperation() operation.ObjectOperationInterface
}

type core struct {
	association    operation.AssociationOperationInterface
	attribute      operation.AttributeOperationInterface
	classification operation.ClassificationOperationInterface
	group          operation.GroupOperationInterface
	inst           operation.InstOperationInterface
	object         operation.ObjectOperationInterface
}

// New create a core manager
func New(client apimachinery.ClientSetInterface) Core {
	targetModel := model.New(client)
	targetInst := inst.New(client)
	return &core{
		inst:           operation.NewInstOperation(client, targetModel, targetInst),
		association:    operation.NewAssociationOperation(client, targetModel, targetInst),
		attribute:      operation.NewAttributeOperation(client, targetModel, targetInst),
		classification: operation.NewClassificationOperation(client, targetModel, targetInst),
		group:          operation.NewGroupOperation(client, targetModel, targetInst),
		object:         operation.NewObjectOperation(client, targetModel, targetInst),
	}
}

func (cli *core) AssociationOperation() operation.AssociationOperationInterface {
	return cli.association
}
func (cli *core) AttributeOperation() operation.AttributeOperationInterface {
	return cli.attribute
}
func (cli *core) ClassificationOperation() operation.ClassificationOperationInterface {
	return cli.classification
}
func (cli *core) GroupOperation() operation.GroupOperationInterface {
	return cli.group
}
func (cli *core) InstOperation() operation.InstOperationInterface {
	return cli.inst
}
func (cli *core) ObjectOperation() operation.ObjectOperationInterface {
	return cli.object
}
