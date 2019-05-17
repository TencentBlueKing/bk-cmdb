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
)

// CreateSet create a new set object
func CreateSet(supplierAccount string) (*SetWrapper, error) {

	targetModel, err := GetSetModel(supplierAccount)
	if nil != err {
		return nil, err
	}

	setInst := mgr.OutputerMgr.InstOperation().CreateSetInst(targetModel)
	return &SetWrapper{set: setInst}, err
}

// FindSetLikeName find all insts by the name
func FindSetLikeName(supplierAccount, setName string) (*SetIteratorWrapper, error) {
	targetModel, err := GetSetModel(supplierAccount)
	if nil != err {
		return nil, err
	}

	iter, err := mgr.OutputerMgr.InstOperation().FindSetsLikeName(targetModel, setName)
	return &SetIteratorWrapper{set: iter}, err
}

// FindSetByCondition find all insts by the condition
func FindSetByCondition(supplierAccount string, cond common.Condition) (*SetIteratorWrapper, error) {
	targetModel, err := GetSetModel(supplierAccount)
	if nil != err {
		return nil, err
	}
	iter, err := mgr.OutputerMgr.InstOperation().FindSetsByCondition(targetModel, cond)
	return &SetIteratorWrapper{set: iter}, err
}
