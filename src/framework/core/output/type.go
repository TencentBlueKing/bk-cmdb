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

package output

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

// OutputerKey the output name
type OutputerKey string

// Puter send the data input
type Puter interface {
	// save into the storage
	Put(data types.MapStr) error
}

// Outputer is the interface that must be implemented by every Outputer.
type Outputer interface {

	// Name the Inputer description.
	// This information will be printed when the Inputer is abnormal, which is convenient for debugging.
	Name() string

	// Run the output main loop. This should block until singnalled to stop by invocation of the Stop() method.
	Put(data types.MapStr) error

	// Stop is the invoked to signal that the Run() method should its execution.
	// It will be invoked at most once.
	Stop() error
}

// ModelOutputer the interface which used to maintence the model
type ModelOutputer interface {

	// GetModel return the model
	GetModel(supplierAccount, classificationID, objID string) (model.Model, error)

	// CreateClassification create a new classification
	CreateClassification(name string) model.Classification

	// FindClassificationsLikeName find a array of the classification by the name
	FindClassificationsLikeName(supplierAccount, name string) (model.ClassificationIterator, error)

	// FindClassificationsByCondition find a array of the classification by the condition
	FindClassificationsByCondition(supplierAccount string, cond common.Condition) (model.ClassificationIterator, error)
}

// CustomOutputer the interface which used to maintence the custom outputer
type CustomOutputer interface {
	// AddOutputer add a new outputer instance
	AddOutputer(target Outputer) OutputerKey

	// RemoveOutputer delete the outputer instace by the outputer key
	RemoveOutputer(key OutputerKey)

	// FetchOutputer find and return the puter by the outputer key
	FetchOutputer(key OutputerKey) Puter

	// CreateCustomOutputer create a new custom outputer
	CreateCustomOutputer(name string, run func(data types.MapStr) error) (OutputerKey, Puter)
}

// Manager is the interface that must be implemented by every output manager.
type Manager interface {
	// Model interface
	ModelOutputer

	// Custom outputer
	CustomOutputer

	// InstOperation operation
	InstOperation() inst.OperationInterface
}
