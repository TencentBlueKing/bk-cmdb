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
	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
)

// CreateBusiness create a new Business object
func CreateBusiness(supplierAccount string) (inst.Inst, error) {

	// TODO:需要根据supplierAccount 获取业务模型
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateSet create a new set object
func CreateSet() (inst.Inst, error) {
	// TODO:需要根据supplierAccount 获取集群模型定义
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateModule create a new module object
func CreateModule() (inst.Inst, error) {
	// TODO:需要根据supplierAccount 获取模块的定义
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateHost create a new host object
func CreateHost() (inst.Inst, error) {
	// TODO:需要根据supplierAccount 获取模块的定义
	return mgr.OutputerMgr.CreateInst(nil)
}

// CreateCommonInst create a common inst object
func CreateCommonInst(target model.Model) (inst.Inst, error) {
	// TODO:根据model 创建普通实例
	return mgr.OutputerMgr.CreateInst(nil)
}
