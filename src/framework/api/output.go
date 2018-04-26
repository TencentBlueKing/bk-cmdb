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
 
package api

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/output"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
	"errors"
)

// CreateCustomOutputer create a new custom outputer
func CreateCustomOutputer(name string, runFunc func(data types.MapStr) error) (output.OutputerKey, output.Puter, error) {

	if 0 == len(name) {
		return output.OutputerKey(""), nil, errors.New("the name parmeter must be set")
	}

	if nil == runFunc {
		return output.OutputerKey(""), nil, errors.New("the run function must be set")
	}

	key, sender := mgr.OutputerMgr.CreateCustomOutputer(name, runFunc)
	return key, sender, nil
}

// CreateClassification create a new classification
func CreateClassification() model.Classification {
	return mgr.OutputerMgr.CreateClassification()
}

// FindClassificationsLikeName find a array of the classification by the name
func FindClassificationsLikeName(name string) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsLikeName(name)
}

// FindClassificationsByCondition find a array of the classification by the condition
func FindClassificationsByCondition(condition *common.Condition) (model.ClassificationIterator, error) {
	return mgr.OutputerMgr.FindClassificationsByCondition(condition)
}
