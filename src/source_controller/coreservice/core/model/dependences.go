/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package model

import (
	"configcenter/src/common/http/rest"
)

// ATTENTIONS: the dependent methods of the other module

// OperationDependences methods definition
type OperationDependences interface {

	// HasInstance used to check if the model has some instances
	HasInstance(kit *rest.Kit, objIDS []string) (exists bool, err error)

	// HasAssociation used to check if the model has some associations
	HasAssociation(kit *rest.Kit, objIDS []string) (exists bool, err error)

	// CascadeDeleteAssociation cascade delete all associated data (included instances, model association, instance association) associated with modelObjID
	CascadeDeleteAssociation(kit *rest.Kit, objIDS []string) error

	// CascadeDeleteInstances cascade delete all instances(included instances, instance association) associated with modelObjID
	CascadeDeleteInstances(kit *rest.Kit, objIDS []string) error
}
