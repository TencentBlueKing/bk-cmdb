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
	"configcenter/src/framework/core/types"
)

// PlatIteratorWrapper the plat iterator wrapper
type PlatIteratorWrapper struct {
	plat inst.Iterator
}

// Next next the plat
func (cli *PlatIteratorWrapper) Next() (*PlatWrapper, error) {

	plat, err := cli.plat.Next()

	return &PlatWrapper{plat: plat}, err

}

// ForEach the foreach function
func (cli *PlatIteratorWrapper) ForEach(callback func(plat *PlatWrapper) error) error {

	return cli.plat.ForEach(func(item inst.CommonInstInterface) error {
		return callback(&PlatWrapper{plat: item})
	})
}

// PlatWrapper the plat wrapper
type PlatWrapper struct {
	plat inst.CommonInstInterface
}

// SetValue set the key value
func (cli *PlatWrapper) SetValue(key string, val interface{}) error {
	return cli.plat.SetValue(key, val)
}

// IsExists check the set
func (cli *PlatWrapper) IsExists() (bool, error) {
	return cli.plat.IsExists()
}

// Create only to create
func (cli *PlatWrapper) Create() error {
	return cli.plat.Create()
}

// Update only to update
func (cli *PlatWrapper) Update() error {
	return cli.plat.Update()
}

// Save save the data
func (cli *PlatWrapper) Save() error {
	return cli.plat.Save()
}

// GetValues return the values
func (cli *PlatWrapper) GetValues() (types.MapStr, error) {
	return cli.plat.GetValues()
}

// GetModel get the model for the plat
func (cli *PlatWrapper) GetModel() model.Model {
	return cli.plat.GetModel()
}

// SetSupplierAccount set the supplier account code of the set
func (cli *PlatWrapper) SetSupplierAccount(supplierAccount string) error {
	return cli.plat.SetValue(fieldSupplierAccount, supplierAccount)
}

// GetSupplierAccount get the supplier account
func (cli *PlatWrapper) GetSupplierAccount() (string, error) {
	vals, err := cli.plat.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldSupplierAccount), nil
}

// GetID get the set id
func (cli *PlatWrapper) GetID() (int, error) {
	vals, err := cli.plat.GetValues()
	if nil != err {
		return 0, err
	}
	return vals.Int(fieldPlatID)
}

// SetName the name of the set
func (cli *PlatWrapper) SetName(name string) error {
	return cli.plat.SetValue(fieldPlatName, name)
}

// GetName get the set name
func (cli *PlatWrapper) GetName() (string, error) {
	vals, err := cli.plat.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldPlatName), nil
}
