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
)

// ModuleWrapper the module wrapper
type ModuleWrapper struct {
	module inst.Inst
}

// SetValue set the key value
func (cli *ModuleWrapper) SetValue(key string, val interface{}) error {
	return cli.module.SetValue(key, val)
}

// Save save the data
func (cli *ModuleWrapper) Save() error {
	return cli.module.Save()
}

// SetOperator set the operator
func (cli *ModuleWrapper) SetOperator(operator string) error {
	return cli.module.SetValue(fieldOperator, operator)
}

// SetBakOperator set the bak operator
func (cli *ModuleWrapper) SetBakOperator(bakOperator string) error {
	return cli.module.SetValue(fieldBakOperator, bakOperator)
}

// SetBussinessID set the business for the module
func (cli *ModuleWrapper) SetBussinessID(businessID int64) error {
	return cli.module.SetValue(fieldBusinessID, businessID)
}

// SetSupplierAccount set the supplier account
func (cli *ModuleWrapper) SetSupplierAccount(supplierAccount string) error {
	return cli.module.SetValue(fieldSupplierAccount, supplierAccount)
}

// SetParent set the parent inst
func (cli *ModuleWrapper) SetParent(parentInstID int64) error {
	if err := cli.module.SetValue(fieldSetID, parentInstID); nil != err {
		return err
	}
	return cli.module.SetValue(fieldParentID, parentInstID)
}

// SetName set the module name
func (cli *ModuleWrapper) SetName(name string) error {
	return cli.module.SetValue(fieldModuleName, name)
}

// SetID set the module id
func (cli *ModuleWrapper) SetID(id string) error {
	return cli.module.SetValue(fieldModuleID, id)
}
