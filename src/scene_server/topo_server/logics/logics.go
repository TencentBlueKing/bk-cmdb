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
	AttributeOperation() model.AttributeOperationInterface
	ObjectOperation() model.ObjectOperationInterface
	IdentifierOperation() operation.IdentifierOperationInterface
	AssociationOperation() model.AssociationOperationInterface
	InstAssociationOperation() inst.AssociationOperationInterface
	GraphicsOperation() operation.GraphicsOperationInterface
	GroupOperation() model.GroupOperationInterface
}

type logics struct {
	classification  model.ClassificationOperationInterface
	attribute       model.AttributeOperationInterface
	object          model.ObjectOperationInterface
	identifier      operation.IdentifierOperationInterface
	association     model.AssociationOperationInterface
	instassociation inst.AssociationOperationInterface
	graphics        operation.GraphicsOperationInterface
	group           model.GroupOperationInterface
}

// New create a logics manager
func New(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager,
	languageIf language.CCLanguageIf) Logics {
	classificationOperation := model.NewClassificationOperation(client, authManager)
	attributeOperation := model.NewAttributeOperation(client, authManager, languageIf)
	objectOperation := model.NewObjectOperation(client, authManager)
	IdentifierOperation := operation.NewIdentifier(client)
	associationOperation := model.NewAssociationOperation(client, authManager)
	instAssociationOperation := inst.NewAssociationOperation(client, authManager)
	graphicsOperation := operation.NewGraphics(client, authManager)
	groupOperation := model.NewGroupOperation(client)
	groupOperation.SetProxy(objectOperation)
	attributeOperation.SetProxy(groupOperation)

	return &logics{
		classification:  classificationOperation,
		attribute:       attributeOperation,
		object:          objectOperation,
		identifier:      IdentifierOperation,
		association:     associationOperation,
		instassociation: instAssociationOperation,
		graphics:        graphicsOperation,
		group:           groupOperation,
	}
}

// AttributeOperation return a attribute provide AttributeOperationInterface
func (l *logics) AttributeOperation() model.AttributeOperationInterface {
	return l.attribute
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
