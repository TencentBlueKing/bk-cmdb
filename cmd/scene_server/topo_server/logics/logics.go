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

// Package logics TODO
package logics

import (
	"configcenter/api"
	"configcenter/cmd/scene_server/topo_server/logics/container"
	inst2 "configcenter/cmd/scene_server/topo_server/logics/inst"
	model2 "configcenter/cmd/scene_server/topo_server/logics/model"
	operation2 "configcenter/cmd/scene_server/topo_server/logics/operation"
	"configcenter/cmd/scene_server/topo_server/logics/settemplate"
	"configcenter/pkg/ac/extensions"
	"configcenter/pkg/language"
)

// Logics provides management interface for operations of model and instance and related resources like association
type Logics interface {
	ClassificationOperation() model2.ClassificationOperationInterface
	SetOperation() inst2.SetOperationInterface
	ModuleOperation() inst2.ModuleOperationInterface
	AttributeOperation() model2.AttributeOperationInterface
	InstOperation() inst2.InstOperationInterface
	ObjectOperation() model2.ObjectOperationInterface
	IdentifierOperation() operation2.IdentifierOperationInterface
	AssociationOperation() model2.AssociationOperationInterface
	InstAssociationOperation() inst2.AssociationOperationInterface
	ImportAssociationOperation() operation2.AssociationOperationInterface
	GraphicsOperation() operation2.GraphicsOperationInterface
	GroupOperation() model2.GroupOperationInterface
	BusinessOperation() inst2.BusinessOperationInterface
	BusinessSetOperation() inst2.BusinessSetOperationInterface
	SetTemplateOperation() settemplate.SetTemplate
	ContainerOperation() container.ClusterOperationInterface
}

type logics struct {
	classification    model2.ClassificationOperationInterface
	set               inst2.SetOperationInterface
	object            model2.ObjectOperationInterface
	identifier        operation2.IdentifierOperationInterface
	module            inst2.ModuleOperationInterface
	attribute         model2.AttributeOperationInterface
	inst              inst2.InstOperationInterface
	association       model2.AssociationOperationInterface
	instassociation   inst2.AssociationOperationInterface
	graphics          operation2.GraphicsOperationInterface
	group             model2.GroupOperationInterface
	importassociation operation2.AssociationOperationInterface
	business          inst2.BusinessOperationInterface
	businessSet       inst2.BusinessSetOperationInterface
	setTemplate       settemplate.SetTemplate
	container         container.ClusterOperationInterface
}

// New create a logics manager
func New(client api.ClientSetInterface, authManager *extensions.AuthManager,
	languageIf language.CCLanguageIf) Logics {
	classificationOperation := model2.NewClassificationOperation(client, authManager)
	setOperation := inst2.NewSetOperation(client, languageIf)
	moduleOperation := inst2.NewModuleOperation(client, authManager)
	attributeOperation := model2.NewAttributeOperation(client, authManager, languageIf)
	objectOperation := model2.NewObjectOperation(client, authManager)
	IdentifierOperation := operation2.NewIdentifier(client)
	associationOperation := model2.NewAssociationOperation(client, authManager)
	instAssociationOperation := inst2.NewAssociationOperation(client, authManager)
	importAssociationOperation := operation2.NewAssociationOperation(client, authManager)
	instOperation := inst2.NewInstOperation(client, languageIf, authManager)
	graphicsOperation := operation2.NewGraphics(client, authManager)
	groupOperation := model2.NewGroupOperation(client)
	businessOperation := inst2.NewBusinessOperation(client, authManager)
	businessSetOperation := inst2.NewBusinessSetOperation(client, authManager)
	containerOperation := container.NewClusterOperation(client, authManager)
	setTemplate := settemplate.NewSetTemplate(client)

	instOperation.SetProxy(instAssociationOperation)
	instAssociationOperation.SetProxy(instOperation)
	associationOperation.SetProxy(objectOperation, instOperation, instAssociationOperation)
	groupOperation.SetProxy(objectOperation)
	setOperation.SetProxy(instOperation, moduleOperation)
	moduleOperation.SetProxy(instOperation)
	attributeOperation.SetProxy(groupOperation, objectOperation)
	businessOperation.SetProxy(instOperation, moduleOperation, setOperation)
	businessSetOperation.SetProxy(instOperation)
	containerOperation.SetProxy(containerOperation)
	return &logics{
		classification:    classificationOperation,
		set:               setOperation,
		object:            objectOperation,
		identifier:        IdentifierOperation,
		inst:              instOperation,
		association:       associationOperation,
		module:            moduleOperation,
		attribute:         attributeOperation,
		instassociation:   instAssociationOperation,
		graphics:          graphicsOperation,
		group:             groupOperation,
		importassociation: importAssociationOperation,
		business:          businessOperation,
		businessSet:       businessSetOperation,
		setTemplate:       setTemplate,
		container:         containerOperation,
	}
}

// SetOperation return a setOperation provide SetOperationInterface
func (l *logics) SetOperation() inst2.SetOperationInterface {
	return l.set
}

// ModuleOperation return a module provide ModuleOperationInterface
func (l *logics) ModuleOperation() inst2.ModuleOperationInterface {
	return l.module
}

// AttributeOperation return a attribute provide AttributeOperationInterface
func (l *logics) AttributeOperation() model2.AttributeOperationInterface {
	return l.attribute
}

// ClassificationOperation return a classification provide ClassificationOperationInterface
func (l *logics) ClassificationOperation() model2.ClassificationOperationInterface {
	return l.classification
}

// ObjectOperation return a object provide ObjectOperationInterface
func (l *logics) ObjectOperation() model2.ObjectOperationInterface {
	return l.object
}

// IdentifierOperation return a identifier provide IdentifierOperationInterface
func (l *logics) IdentifierOperation() operation2.IdentifierOperationInterface {
	return l.identifier
}

// InstOperation return a inst provide InstOperationInterface
func (l *logics) InstOperation() inst2.InstOperationInterface {
	return l.inst
}

// AssociationOperation return a association provide AssociationOperationInterface
func (l *logics) AssociationOperation() model2.AssociationOperationInterface {
	return l.association
}

// InstAssociationOperation return a instance association provide AssociationOperationInterface
func (l *logics) InstAssociationOperation() inst2.AssociationOperationInterface {
	return l.instassociation
}

// ImportAssociationOperation return a import association provide AssociationOperationInterface
func (l *logics) ImportAssociationOperation() operation2.AssociationOperationInterface {
	return l.importassociation
}

// GraphicsOperation return a inst provide GraphicsOperation
func (l *logics) GraphicsOperation() operation2.GraphicsOperationInterface {
	return l.graphics
}

// GroupOperation return a inst provide GroupOperationInterface
func (l *logics) GroupOperation() model2.GroupOperationInterface {
	return l.group
}

// BusinessOperation return a inst provide BusinessOperation
func (l *logics) BusinessOperation() inst2.BusinessOperationInterface {
	return l.business
}

// BusinessSetOperation return a inst provide BusinessOperation
func (l *logics) BusinessSetOperation() inst2.BusinessSetOperationInterface {
	return l.businessSet
}

// ContainerOperation return a inst provide ContainerOperation
func (l *logics) ContainerOperation() container.ClusterOperationInterface {
	return l.container
}

// SetTemplateOperation set template operation
func (l *logics) SetTemplateOperation() settemplate.SetTemplate {
	return l.setTemplate
}
