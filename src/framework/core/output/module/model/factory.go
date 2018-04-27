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

package model

import "configcenter/src/framework/common"

// CreateClassification create a new Classification instance
func CreateClassification(name string) Classification {
	return &classification{classificationName: name}
}

// FindClassificationsLikeName find a array of the classification by the name
func FindClassificationsLikeName(name string) (ClassificationIterator, error) {
	// TODO: 按照名字模糊查找
	return nil, nil
}

// FindClassificationsByCondition find a array of the classification by the condition
func FindClassificationsByCondition(condition common.Condition) (ClassificationIterator, error) {
	// TODO: 按照条件搜索
	return nil, nil
}
