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
	SetOperation() operation.SetOperationInterface
	ModuleOperation() operation.ModuleOperationInterface
	BusinessOperation() operation.BusinessOperationInterface
	AssociationOperation() operation.AssociationOperationInterface
	AttributeOperation() operation.AttributeOperationInterface
	ClassificationOperation() operation.ClassificationOperationInterface
	GroupOperation() operation.GroupOperationInterface
	InstOperation() operation.InstOperationInterface
	ObjectOperation() operation.ObjectOperationInterface
	PermissionOperation() operation.PermissionOperationInterface
	CompatibleV2Operation() operation.CompatibleV2OperationInterface
}

type core struct {
	business       operation.BusinessOperationInterface
	set            operation.SetOperationInterface
	module         operation.ModuleOperationInterface
	association    operation.AssociationOperationInterface
	attribute      operation.AttributeOperationInterface
	classification operation.ClassificationOperationInterface
	group          operation.GroupOperationInterface
	inst           operation.InstOperationInterface
	object         operation.ObjectOperationInterface
	permission     operation.PermissionOperationInterface
	compatibleV2   operation.CompatibleV2OperationInterface
}

// New create a core manager
func New(client apimachinery.ClientSetInterface) Core {

	targetModel := model.New(client)
	targetInst := inst.New(client)

	inst := operation.NewInstOperation(client, targetModel, targetInst)
	set := operation.NewSetOperation(client, inst)
	module := operation.NewModuleOperation(client, inst)
	business := operation.NewBusinessOperation(set, module, client, inst)

	attribute := operation.NewAttributeOperation(client, targetModel, targetInst)
	classification := operation.NewClassificationOperation(client, targetModel, targetInst)
	group := operation.NewGroupOperation(client, targetModel, targetInst)
	object := operation.NewObjectOperation(client, targetModel, targetInst)

	association := operation.NewAssociationOperation(client, classification, object, attribute, inst, targetModel, targetInst)
	permission := operation.NewPermissionOperation(client)
	compatibleV2 := operation.NewCompatibleV2Operation(client)

	return &core{
		set:            set,
		module:         module,
		business:       business,
		inst:           inst,
		association:    association,
		attribute:      attribute,
		classification: classification,
		group:          group,
		object:         object,
		permission:     permission,
		compatibleV2:   compatibleV2,
	}
}

func (c *core) SetOperation() operation.SetOperationInterface {
	return c.set
}

func (c *core) ModuleOperation() operation.ModuleOperationInterface {
	return c.module
}

func (c *core) BusinessOperation() operation.BusinessOperationInterface {
	return c.business
}

func (c *core) AssociationOperation() operation.AssociationOperationInterface {
	return c.association
}
func (c *core) AttributeOperation() operation.AttributeOperationInterface {
	return c.attribute
}
func (c *core) ClassificationOperation() operation.ClassificationOperationInterface {
	return c.classification
}
func (c *core) GroupOperation() operation.GroupOperationInterface {
	return c.group
}
func (c *core) InstOperation() operation.InstOperationInterface {
	return c.inst
}
func (c *core) ObjectOperation() operation.ObjectOperationInterface {
	return c.object
}
func (c *core) PermissionOperation() operation.PermissionOperationInterface {
	return c.permission
}
func (c *core) CompatibleV2Operation() operation.CompatibleV2OperationInterface {
	return c.compatibleV2
}
