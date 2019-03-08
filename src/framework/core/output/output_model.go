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
	"configcenter/src/framework/core/output/module/model"
)

// GetModel get the model
func (cli *manager) GetModel(supplierAccount, classificationID, objID string) (model.Model, error) {
	condInner := common.CreateCondition().Field(model.ClassificationID).Eq(classificationID)
	clsIter, err := cli.FindClassificationsByCondition(supplierAccount, condInner)
	if nil != err {
		return nil, err
	}
	//fmt.Println("owner:", supplierAccount)
	var targetModel model.Model
	err = clsIter.ForEach(func(item model.Classification) error {

		condInner = common.CreateCondition().Field(model.ObjectID).Eq(objID)
		condInner.Field(model.SupplierAccount).Eq(supplierAccount)
		condInner.Field(model.ClassificationID).Eq(item.GetID())

		modelIter, err := item.FindModelsByCondition(supplierAccount, condInner)
		if nil != err {
			return err
		}

		err = modelIter.ForEach(func(modelItem model.Model) error {
			targetModel = modelItem
			return nil
		})

		return nil
	})

	if nil != err {
		return nil, err
	}
	//fmt.Println("model:", targetModel)
	return targetModel, err
}

// CreateClassification create a new classification
func (cli *manager) CreateClassification(name string) model.Classification {
	return model.CreateClassification(name)
}

// FindClassificationsLikeName find a array of the classification by the name
func (cli *manager) FindClassificationsLikeName(supplierAccount, name string) (model.ClassificationIterator, error) {
	return model.FindClassificationsLikeName(supplierAccount, name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func (cli *manager) FindClassificationsByCondition(supplierAccount string, cond common.Condition) (model.ClassificationIterator, error) {
	return model.FindClassificationsByCondition(supplierAccount, cond)
}
