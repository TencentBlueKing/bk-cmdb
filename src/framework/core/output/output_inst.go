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
)

// CreateInst create a instance for the model
func (cli *manager) CreateInst(target model.Model) (inst.Inst, error) {
	return inst.CreateInst(target)
}

// FindInstsLikeName find all insts by the name
func (cli *manager) FindInstsLikeName(target model.Model, instName string) (inst.Iterator, error) {
	return inst.FindInstsLikeName(target, instName)
}

// FindInstsByCondition find all insts by the condition
func (cli *manager) FindInstsByCondition(target model.Model, condition common.Condition) (inst.Iterator, error) {
	return inst.FindInstsByCondition(target, condition)
}
