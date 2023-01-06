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
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common/language"
	"configcenter/src/scene_server/topo_server/logics/inst"
	"configcenter/src/scene_server/topo_server/logics/kube"
	"configcenter/src/scene_server/topo_server/logics/model"
	"configcenter/src/scene_server/topo_server/logics/operation"
	"configcenter/src/scene_server/topo_server/logics/settemplate"
)

// Logics provides management interface for operations of model and instance and related resources like association
type Logics interface {
	ClassificationOperation() model.ClassificationOperationInterface
	SetOperation() inst.SetOperationInterface
	ModuleOperation() inst.ModuleOperationInterface
	AttributeOperation() model.AttributeOperationInterface
	InstOperation() inst.InstOperationInterface
	ObjectOperation() model.ObjectOperationInterface
	IdentifierOperation() operation.IdentifierOperationInterface
	AssociationOperation() model.AssociationOperationInterface
	InstAssociationOperation() inst.AssociationOperationInterface
	ImportAssociationOperation() operation.AssociationOperationInterface
	GraphicsOperation() operation.GraphicsOperationInterface
	GroupOperation() model.GroupOperationInterface
	BusinessOperation() inst.BusinessOperationInterface
	BusinessSetOperation() inst.BusinessSetOperationInterface
	SetTemplateOperation() settemplate.SetTemplate
	KubeOperation() kube.KubeOperationInterface
}

type logics struct {
	classification    model.ClassificationOperationInterface
	set               inst.SetOperationInterface
	object            model.ObjectOperationInterface
	identifier        operation.IdentifierOperationInterface
	module            inst.ModuleOperationInterface
	attribute         model.AttributeOperationInterface
	inst              inst.InstOperationInterface
	association       model.AssociationOperationInterface
	instassociation   inst.AssociationOperationInterface
	graphics          operation.GraphicsOperationInterface
	group             model.GroupOperationInterface
	importassociation operation.AssociationOperationInterface
	business          inst.BusinessOperationInterface
	businessSet       inst.BusinessSetOperationInterface
	setTemplate       settemplate.SetTemplate
	kube              kube.KubeOperationInterface
}

// New create a logics manager
func New(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager,
	languageIf language.CCLanguageIf) Logics {
	classificationOperation := model.NewClassificationOperation(client, authManager)
	setOperation := inst.NewSetOperation(client, languageIf)
	moduleOperation := inst.NewModuleOperation(client, authManager)
	attributeOperation := model.NewAttributeOperation(client, authManager, languageIf)
	objectOperation := model.NewObjectOperation(client, authManager)
	IdentifierOperation := operation.NewIdentifier(client)
	associationOperation := model.NewAssociationOperation(client, authManager)
	instAssociationOperation := inst.NewAssociationOperation(client, authManager)
	importAssociationOperation := operation.NewAssociationOperation(client, authManager)
	instOperation := inst.NewInstOperation(client, languageIf, authManager)
	graphicsOperation := operation.NewGraphics(client, authManager)
	groupOperation := model.NewGroupOperation(client)
	businessOperation := inst.NewBusinessOperation(client, authManager)
	businessSetOperation := inst.NewBusinessSetOperation(client, authManager)
	kubeOperation := kube.NewClusterOperation(client, authManager)
	setTemplate := settemplate.NewSetTemplate(client)

	instOperation.SetProxy(instAssociationOperation)
	instAssociationOperation.SetProxy(instOperation)
	associationOperation.SetProxy(objectOperation, instOperation, instAssociationOperation)
	importAssociationOperation.SetProxy(instAssociationOperation)
	groupOperation.SetProxy(objectOperation)
	setOperation.SetProxy(instOperation, moduleOperation)
	moduleOperation.SetProxy(instOperation)
	attributeOperation.SetProxy(groupOperation, objectOperation)
	businessOperation.SetProxy(instOperation, moduleOperation, setOperation)
	businessSetOperation.SetProxy(instOperation)
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
		kube:              kubeOperation,
	}
}

// SetOperation return a setOperation provide SetOperationInterface
func (l *logics) SetOperation() inst.SetOperationInterface {
	return l.set
}

// ModuleOperation return a module provide ModuleOperationInterface
func (l *logics) ModuleOperation() inst.ModuleOperationInterface {
	return l.module
}

// AttributeOperation return a attribute provide AttributeOperationInterface
func (l *logics) AttributeOperation() model.AttributeOperationInterface {
	return l.attribute
}

// ClassificationOperation return a classification provide ClassificationOperationInterface
func (l *logics) ClassificationOperation() model.ClassificationOperationInterface {
	return l.classification
}

// ObjectOperation return a object provide ObjectOperationInterface
func (l *logics) ObjectOperation() model.ObjectOperationInterface {
	return l.object
}

// IdentifierOperation return a identifier provide IdentifierOperationInterface
func (l *logics) IdentifierOperation() operation.IdentifierOperationInterface {
	return l.identifier
}

// InstOperation return a inst provide InstOperationInterface
func (l *logics) InstOperation() inst.InstOperationInterface {
	return l.inst
}

// AssociationOperation return a association provide AssociationOperationInterface
func (l *logics) AssociationOperation() model.AssociationOperationInterface {
	return l.association
}

// InstAssociationOperation return a instance association provide AssociationOperationInterface
func (l *logics) InstAssociationOperation() inst.AssociationOperationInterface {
	return l.instassociation
}

// ImportAssociationOperation return a import association provide AssociationOperationInterface
func (l *logics) ImportAssociationOperation() operation.AssociationOperationInterface {
	return l.importassociation
}

// GraphicsOperation return a inst provide GraphicsOperation
func (l *logics) GraphicsOperation() operation.GraphicsOperationInterface {
	return l.graphics
}

// GroupOperation return a inst provide GroupOperationInterface
func (l *logics) GroupOperation() model.GroupOperationInterface {
	return l.group
}

// BusinessOperation return a inst provide BusinessOperation
func (l *logics) BusinessOperation() inst.BusinessOperationInterface {
	return l.business
}

// BusinessSetOperation return a inst provide BusinessOperation
func (l *logics) BusinessSetOperation() inst.BusinessSetOperationInterface {
	return l.businessSet
}

// kubeOperation return a inst provide kubeOperation
func (l *logics) KubeOperation() kube.KubeOperationInterface {
	return l.kube
}

// SetTemplateOperation set template operation
func (l *logics) SetTemplateOperation() settemplate.SetTemplate {
	return l.setTemplate
}
