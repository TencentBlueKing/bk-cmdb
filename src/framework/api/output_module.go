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

// CreateModule create a new module object
func CreateModule(supplierAccount string) (*ModuleWrapper, error) {

	targetModel, err := GetModuleModel(supplierAccount)
	if nil != err {
		return nil, err
	}

	moduleInst := mgr.OutputerMgr.InstOperation().CreateModuleInst(targetModel)
	return &ModuleWrapper{module: moduleInst}, err
}

// FindModuleLikeName find all insts by the name
func FindModuleLikeName(supplierAccount, moduleName string) (*ModuleIteratorWrapper, error) {
	targetModel, err := GetModuleModel(supplierAccount)
	if nil != err {
		return nil, err
	}

	iter, err := mgr.OutputerMgr.InstOperation().FindModulesLikeName(targetModel, moduleName)
	return &ModuleIteratorWrapper{module: iter}, err
}

// FindModuleByCondition find all insts by the condition
func FindModuleByCondition(supplierAccount string, cond common.Condition) (*ModuleIteratorWrapper, error) {
	targetModel, err := GetModuleModel(supplierAccount)
	if nil != err {
		return nil, err
	}
	iter, err := mgr.OutputerMgr.InstOperation().FindModulesByCondition(targetModel, cond)
	return &ModuleIteratorWrapper{module: iter}, err
}
