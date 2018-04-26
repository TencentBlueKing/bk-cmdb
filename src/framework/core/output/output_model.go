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

// CreateClassification create a new classification
func (cli *manager) CreateClassification() model.Classification {
	return model.CreateClassification()
}

// FindClassificationsLikeName find a array of the classification by the name
func (cli *manager) FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return model.FindClassificationsLikeName(name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func (cli *manager) FindClassificationsByCondition(condition *common.Condition) (model.ClassificationIterator, error) {
	return model.FindClassificationsByCondition(condition)
}
