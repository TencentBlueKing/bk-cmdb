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
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common/language"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/operation"
	"configcenter/src/scene_server/topo_server/core/settemplate"
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
	GraphicsOperation() operation.GraphicsOperationInterface
	IdentifierOperation() operation.IdentifierOperationInterface
	AuditOperation() operation.AuditOperationInterface
	UniqueOperation() operation.UniqueOperationInterface
	SetTemplateOperation() settemplate.SetTemplate
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
	graphics       operation.GraphicsOperationInterface
	audit          operation.AuditOperationInterface
	identifier     operation.IdentifierOperationInterface
	unique         operation.UniqueOperationInterface
	setTemplate    settemplate.SetTemplate
}

// New create a logics manager
func New(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager, languageIf language.CCLanguageIf) Core {
	// create instances
	attributeOperation := operation.NewAttributeOperation(client, authManager)
	classificationOperation := operation.NewClassificationOperation(client, authManager)
	groupOperation := operation.NewGroupOperation(client)
	objectOperation := operation.NewObjectOperation(client, authManager)
	instOperation := operation.NewInstOperation(client, languageIf, authManager)
	moduleOperation := operation.NewModuleOperation(client, authManager)
	setOperation := operation.NewSetOperation(client, languageIf)
	businessOperation := operation.NewBusinessOperation(client, authManager)
	associationOperation := operation.NewAssociationOperation(client, authManager)
	graphics := operation.NewGraphics(client, authManager)
	identifier := operation.NewIdentifier(client)
	audit := operation.NewAuditOperation(client)
	unique := operation.NewUniqueOperation(client, authManager)
	setTemplate := settemplate.NewSetTemplate(client)

	targetModel := model.New(client, languageIf)
	targetInst := inst.New(client)

	// set the operation
	objectOperation.SetProxy(targetModel, targetInst, classificationOperation, associationOperation, instOperation, attributeOperation, groupOperation, unique)
	groupOperation.SetProxy(targetModel, targetInst, objectOperation)
	attributeOperation.SetProxy(targetModel, targetInst, objectOperation, associationOperation, groupOperation)
	classificationOperation.SetProxy(targetModel, targetInst, associationOperation, objectOperation)
	associationOperation.SetProxy(classificationOperation, objectOperation, groupOperation, attributeOperation, instOperation, targetModel, targetInst)

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
		graphics:       graphics,
		audit:          audit,
		identifier:     identifier,
		unique:         unique,
		setTemplate:    setTemplate,
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
func (c *core) GraphicsOperation() operation.GraphicsOperationInterface {
	return c.graphics
}
func (c *core) AuditOperation() operation.AuditOperationInterface {
	return c.audit
}
func (c *core) IdentifierOperation() operation.IdentifierOperationInterface {
	return c.identifier
}
func (c *core) UniqueOperation() operation.UniqueOperationInterface {
	return c.unique
}
func (c *core) SetTemplateOperation() settemplate.SetTemplate {
	return c.setTemplate
}
