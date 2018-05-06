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

// SetWrapper the set wrapper
type SetWrapper struct {
	set inst.Inst
}

// SetValue set the key value
func (cli *SetWrapper) SetValue(key string, val interface{}) error {
	return cli.set.SetValue(key, val)
}

// SetIntroduction set the introducrtion of the set
func (cli *SetWrapper) SetIntroduction(intro string) error {
	return cli.set.SetValue(fieldSetDesc, intro)
}

// SetDescription set the description of the set
func (cli *SetWrapper) SetDescription(desc string) error {
	return cli.set.SetValue(fieldDescription, desc)
}

// SetEnv set the env of the set
func (cli *SetWrapper) SetEnv(env string) error {
	return cli.set.SetValue(fieldSetEnv, env)
}

// SetServiceStatus the service status of the set
func (cli *SetWrapper) SetServiceStatus(status string) error {
	return cli.set.SetValue(fieldServiceStatus, status)
}

// SetCapacity set the capacity of the set
func (cli *SetWrapper) SetCapacity(capacity int64) error {
	return cli.set.SetValue(fieldCapacity, capacity)
}

// SetBussinessID set the business id of the set
func (cli *SetWrapper) SetBussinessID(businessID string) error {
	return cli.set.SetValue(fieldBusinessID, businessID)
}

// SetSupplierAccount set the supplier account code of the set
func (cli *SetWrapper) SetSupplierAccount(supplierAccount string) error {
	return cli.set.SetValue(fieldSupplierAccount, supplierAccount)
}

// SetID the id of the set
func (cli *SetWrapper) SetID(id string) error {
	return cli.set.SetValue(fieldSetID, id)
}

// SetParent set the parent id of the set
func (cli *SetWrapper) SetParent(parentInstID int64) error {
	return cli.set.SetValue(fieldParentID, parentInstID)
}

// SetName the name of the set
func (cli *SetWrapper) SetName(name string) error {
	return cli.set.SetValue(fieldSetName, name)
}

// Save save the data
func (cli *SetWrapper) Save() error {
	return cli.set.Save()
}
