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
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
)

// BaseInstOperation return the base inst operation interface
func BaseInstOperation() inst.OperationInterface {
	return mgr.OutputerMgr.InstOperation()
}

// CreateCommonInst create a common inst object
func CreateCommonInst(target model.Model) (inst.CommonInstInterface, error) {
	return mgr.OutputerMgr.InstOperation().CreateCommonInst(target), nil
}

// FindInstsLikeName find all insts by the name
func FindInstsLikeName(target model.Model, instName string) (inst.Iterator, error) {
	return mgr.OutputerMgr.InstOperation().FindCommonInstLikeName(target, instName)
}

// FindInstsByCondition find all insts by the condition
func FindInstsByCondition(target model.Model, cond common.Condition) (inst.Iterator, error) {
	return mgr.OutputerMgr.InstOperation().FindCommonInstByCondition(target, cond)
}
