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
	coreInst "configcenter/src/scene_server/topo_server/core/inst"
	coreModel "configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/operation"
	"configcenter/src/scene_server/topo_server/logics/inst"
	"configcenter/src/scene_server/topo_server/logics/model"
)

// Logics provides management interface for operations of model and instance and related resources like association
type Logics interface {
	ClassificationOperation() model.ClassificationOperationInterface
	BusinessOperation() inst.BusinessOperationInterface
}

type logics struct {
	classification    model.ClassificationOperationInterface
	businessOperation inst.BusinessOperationInterface
}

// New create a logics manager , languageIf language.CCLanguageIf
func New(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) Logics {
	classificationOperation := model.NewClassificationOperation(client, authManager)

	// TODO 临时借调
	targetModel := coreModel.New(client, nil)
	targetInst := coreInst.New(client)
	instOperation := operation.NewInstOperation(client, nil, authManager)
	associationOperation := operation.NewAssociationOperation(client, authManager)
	objectOperation := operation.NewObjectOperation(client, authManager)
	setOperation := operation.NewSetOperation(client, nil)
	moduleOperation := operation.NewModuleOperation(client, authManager)
	instOperation.SetProxy(targetModel, targetInst, associationOperation, objectOperation)
	setOperation.SetProxy(objectOperation, instOperation, moduleOperation)
	moduleOperation.SetProxy(instOperation)

	businessOperationOperation := inst.NewBusinessOperation(client, authManager)
	businessOperationOperation.SetProxy(instOperation, objectOperation, setOperation, moduleOperation)

	return &logics{
		classification:    classificationOperation,
		businessOperation: businessOperationOperation,
	}
}

// ClassificationOperation return a classification provide ClassificationOperationInterface
func (c *logics) ClassificationOperation() model.ClassificationOperationInterface {
	return c.classification
}

// BusinessOperation return a inst provide InstOperationInterface
func (c *logics) BusinessOperation() inst.BusinessOperationInterface {
	return c.businessOperation
}
