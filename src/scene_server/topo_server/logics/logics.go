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
	"configcenter/src/scene_server/topo_server/logics/model"
)

// Logics provides management interface for operations of model and instance and related resources like association
type Logics interface {
	ClassificationOperation() model.ClassificationOperationInterface
	AttributeOperation() model.AttributeOperationInterface
}

type logics struct {
	classification model.ClassificationOperationInterface
	attribute      model.AttributeOperationInterface
}

// New create a logics manager
func New(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) Logics {
	classificationOperation := model.NewClassificationOperation(client, authManager)
	attributeOperation := model.NewAttributeOperation(client, authManager)

	return &logics{
		classification: classificationOperation,
		attribute:      attributeOperation,
	}
}

func (l *logics) ClassificationOperation() model.ClassificationOperationInterface {
	return l.classification
}

func (l *logics) AttributeOperation() model.AttributeOperationInterface {
	return l.attribute
}