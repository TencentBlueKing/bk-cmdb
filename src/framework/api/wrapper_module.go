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
	"fmt"

	"configcenter/src/framework/core/output/module/inst"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

// ModuleIteratorWrapper the module iterator wrapper
type ModuleIteratorWrapper struct {
	module inst.ModuleIterator
}

// Next next the module
func (cli *ModuleIteratorWrapper) Next() (*ModuleWrapper, error) {

	module, err := cli.module.Next()

	return &ModuleWrapper{module: module}, err

}

// ForEach the foreach function
func (cli *ModuleIteratorWrapper) ForEach(callback func(module *ModuleWrapper) error) error {

	return cli.module.ForEach(func(item inst.ModuleInterface) error {
		return callback(&ModuleWrapper{module: item})
	})
}

// ModuleWrapper the module wrapper
type ModuleWrapper struct {
	module inst.ModuleInterface
}

// SetValue set the key value
func (cli *ModuleWrapper) SetValue(key string, val interface{}) error {
	return cli.module.SetValue(key, val)
}

// IsExists check the set
func (cli *ModuleWrapper) IsExists() (bool, error) {
	return cli.module.IsExists()
}

// Create only to create
func (cli *ModuleWrapper) Create() error {
	return cli.module.Create()
}

// Update only to update
func (cli *ModuleWrapper) Update() error {
	return cli.module.Update()
}

// GetModel get the model for the module
func (cli *ModuleWrapper) GetModel() model.Model {
	return cli.module.GetModel()
}

// Save save the data
func (cli *ModuleWrapper) Save() error {
	return cli.module.Save()
}

// GetValues return the values
func (cli *ModuleWrapper) GetValues() (types.MapStr, error) {
	return cli.module.GetValues()
}

// GetModuleID get the id for the module
func (cli *ModuleWrapper) GetModuleID() (int64, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return 0, err
	}
	if !vals.Exists(fieldModuleID) {
		return 0, fmt.Errorf("the module id is not set")
	}
	val, err := vals.Int(fieldModuleID)
	return int64(val), err
}

// SetOperator set the operator
func (cli *ModuleWrapper) SetOperator(operator string) error {
	return cli.module.SetValue(fieldOperator, operator)
}

// GetOperator get the operator for the host
func (cli *ModuleWrapper) GetOperator() (string, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldOperator), nil
}

// SetBakOperator set the bak operator
func (cli *ModuleWrapper) SetBakOperator(bakOperator string) error {
	return cli.module.SetValue(fieldBakOperator, bakOperator)
}

// GetBakOperator get the bak operator
func (cli *ModuleWrapper) GetBakOperator() (string, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldBakOperator), nil
}

// GetBusinessID get the business id
func (cli *ModuleWrapper) GetBusinessID() (int, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return 0, err
	}
	return vals.Int(fieldBusinessID)
}

// SetSupplierAccount set the supplier account
func (cli *ModuleWrapper) SetSupplierAccount(supplierAccount string) error {
	return cli.module.SetValue(fieldSupplierAccount, supplierAccount)
}

// GetSupplierAccount get the supplier account
func (cli *ModuleWrapper) GetSupplierAccount() (string, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldSupplierAccount), nil
}

// SetTopo set the parent inst
func (cli *ModuleWrapper) SetTopo(bizID, setID int64) error {
	cli.module.SetTopo(bizID, setID)
	return nil
}

// SetName set the module name
func (cli *ModuleWrapper) SetName(name string) error {
	return cli.module.SetValue(fieldModuleName, name)
}

// GetName get the module name
func (cli *ModuleWrapper) GetName() (string, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldModuleName), nil
}

// GetID get the id for the host
func (cli *ModuleWrapper) GetID() (int64, error) {
	vals, err := cli.module.GetValues()
	if nil != err {
		return 0, err
	}
	return vals.Int64(fieldModuleID)
}
