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
	GraphicsOperation() operation.GraphicsOperationInterface
	IdentifierOperation() operation.IdentifierOperationInterface
	AuditOperation() operation.AuditOperationInterface
	HealthOperation() operation.HealthOperationInterface
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
	graphics       operation.GraphicsOperationInterface
	audit          operation.AuditOperationInterface
	identifier     operation.IdentifierOperationInterface
	health         operation.HealthOperationInterface
}

// New create a core manager
func New(client apimachinery.ClientSetInterface) Core {

	// health
	healthOpeartion := operation.NewHealthOperation(client)

	// create insts
	attributeOperation := operation.NewAttributeOperation(client)
	classificationOperation := operation.NewClassificationOperation(client)
	groupOperation := operation.NewGroupOperation(client)
	objectOperation := operation.NewObjectOperation(client)
	instOperation := operation.NewInstOperation(client)
	moduleOperation := operation.NewModuleOperation(client)
	setOperation := operation.NewSetOperation(client)
	businessOperation := operation.NewBusinessOperation(client)
	associationOperation := operation.NewAssociationOperation(client)
	permissionOperation := operation.NewPermissionOperation(client)
	compatibleV2Operation := operation.NewCompatibleV2Operation(client)
	graphics := operation.NewGraphics(client)
	identifier := operation.NewIdentifier(client)
	audit := operation.NewAuditOperation(client)

	targetModel := model.New(client)
	targetInst := inst.New(client)

	// set the operation
	objectOperation.SetProxy(targetModel, targetInst, classificationOperation, associationOperation, instOperation, attributeOperation, groupOperation)
	groupOperation.SetProxy(targetModel, targetInst, objectOperation)
	attributeOperation.SetProxy(targetModel, targetInst, objectOperation, associationOperation, groupOperation)
	classificationOperation.SetProxy(targetModel, targetInst, associationOperation, objectOperation)
	associationOperation.SetProxy(classificationOperation, objectOperation, attributeOperation, instOperation, targetModel, targetInst)

	instOperation.SetProxy(targetModel, targetInst, associationOperation, objectOperation)
	moduleOperation.SetProxy(instOperation)
	setOperation.SetProxy(objectOperation, instOperation, moduleOperation)
	businessOperation.SetProxy(setOperation, moduleOperation, instOperation, objectOperation)

	graphics.SetProxy(objectOperation, associationOperation)

	return &core{
		set:            setOperation,
		module:         moduleOperation,
		business:       businessOperation,
		inst:           instOperation,
		association:    associationOperation,
		attribute:      attributeOperation,
		classification: classificationOperation,
		group:          groupOperation,
		object:         objectOperation,
		permission:     permissionOperation,
		compatibleV2:   compatibleV2Operation,
		graphics:       graphics,
		audit:          audit,
		identifier:     identifier,
		health:         healthOpeartion,
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
func (c *core) GraphicsOperation() operation.GraphicsOperationInterface {
	return c.graphics
}
func (c *core) AuditOperation() operation.AuditOperationInterface {
	return c.audit
}
func (c *core) IdentifierOperation() operation.IdentifierOperationInterface {
	return c.identifier
}
func (c *core) HealthOperation() operation.HealthOperationInterface {
	return c.health
}
