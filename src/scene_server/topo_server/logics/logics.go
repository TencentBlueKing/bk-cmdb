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

package logics

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common/language"
	"configcenter/src/scene_server/topo_server/logics/inst"
	"configcenter/src/scene_server/topo_server/logics/model"
	"configcenter/src/scene_server/topo_server/logics/operation"
)

// Logics provides management interface for operations of model and instance and related resources like association
type Logics interface {
	ClassificationOperation() model.ClassificationOperationInterface
	ModuleOperation() inst.ModuleOperationInterface
	AttributeOperation() model.AttributeOperationInterface
	InstOperation() inst.InstOperationInterface
	ObjectOperation() model.ObjectOperationInterface
	IdentifierOperation() operation.IdentifierOperationInterface
	AssociationOperation() model.AssociationOperationInterface
	InstAssociationOperation() inst.AssociationOperationInterface
	GraphicsOperation() operation.GraphicsOperationInterface
	GroupOperation() model.GroupOperationInterface
	BusinessOperation() inst.BusinessOperationInterface
}

type logics struct {
	classification  model.ClassificationOperationInterface
	module          inst.ModuleOperationInterface
	attribute       model.AttributeOperationInterface
	inst            inst.InstOperationInterface
	object          model.ObjectOperationInterface
	identifier      operation.IdentifierOperationInterface
	association     model.AssociationOperationInterface
	instassociation inst.AssociationOperationInterface
	graphics        operation.GraphicsOperationInterface
	group           model.GroupOperationInterface
	business        inst.BusinessOperationInterface
}

// New create a logics manager
func New(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager,
	languageIf language.CCLanguageIf) Logics {
	classificationOperation := model.NewClassificationOperation(client, authManager)
	moduleOperation := inst.NewModuleOperation(client, authManager)
	attributeOperation := model.NewAttributeOperation(client, authManager, languageIf)
	objectOperation := model.NewObjectOperation(client, authManager)
	IdentifierOperation := operation.NewIdentifier(client)
	associationOperation := model.NewAssociationOperation(client, authManager)
	instAssociationOperation := inst.NewAssociationOperation(client, authManager)
	instOperation := inst.NewInstOperation(client, languageIf, authManager)
	graphicsOperation := operation.NewGraphics(client, authManager)
	groupOperation := model.NewGroupOperation(client)

	instOperation.SetProxy(instAssociationOperation)
	instAssociationOperation.SetProxy(instOperation)
	groupOperation.SetProxy(objectOperation)
	moduleOperation.SetProxy(instOperation)
	attributeOperation.SetProxy(groupOperation, objectOperation)
	businessOperation := inst.NewBusinessOperation(client, authManager)
	businessOperation.SetProxy(objectOperation, instOperation, moduleOperation)

	return &logics{
		classification:  classificationOperation,
		module:          moduleOperation,
		attribute:       attributeOperation,
		inst:            instOperation,
		object:          objectOperation,
		identifier:      IdentifierOperation,
		association:     associationOperation,
		instassociation: instAssociationOperation,
		graphics:        graphicsOperation,
		group:           groupOperation,
		business:        businessOperation,
	}
}

// ModuleOperation return a module provide ModuleOperationInterface
func (c *logics) ModuleOperation() inst.ModuleOperationInterface {
	return c.module
}

// AttributeOperation return a attribute provide AttributeOperationInterface
func (c *logics) AttributeOperation() model.AttributeOperationInterface {
	return c.attribute
}

// ClassificationOperation return a classification provide ClassificationOperationInterface
func (c *logics) ClassificationOperation() model.ClassificationOperationInterface {
	return c.classification
}

// ObjectOperation return a object provide ObjectOperationInterface
func (c *logics) ObjectOperation() model.ObjectOperationInterface {
	return c.object
}

// IdentifierOperation return a identifier provide IdentifierOperationInterface
func (c *logics) IdentifierOperation() operation.IdentifierOperationInterface {
	return c.identifier
}

// InstOperation return a inst provide InstOperationInterface
func (c *logics) InstOperation() inst.InstOperationInterface {
	return c.inst
}

// AssociationOperation return a association provide AssociationOperationInterface
func (c *logics) AssociationOperation() model.AssociationOperationInterface {
	return c.association
}

// InstAssociationOperation return a instance association provide AssociationOperationInterface
func (c *logics) InstAssociationOperation() inst.AssociationOperationInterface {
	return c.instassociation
}

// GraphicsOperation return a inst provide GraphicsOperation
func (c *logics) GraphicsOperation() operation.GraphicsOperationInterface {
	return c.graphics
}

// GroupOperation return a inst provide GroupOperationInterface
func (c *logics) GroupOperation() model.GroupOperationInterface {
	return c.group
}

// BusinessOperation return a inst provide BusinessOperation
func (c *logics) BusinessOperation() inst.BusinessOperationInterface {
	return c.business
}
