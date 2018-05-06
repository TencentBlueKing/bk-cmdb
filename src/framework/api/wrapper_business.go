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

// BusinessWrapper the business wrapper
type BusinessWrapper struct {
	business inst.Inst
}

// SetValue set the key value
func (cli *BusinessWrapper) SetValue(key string, val interface{}) error {
	return cli.business.SetValue(key, val)
}

// Save save the data
func (cli *BusinessWrapper) Save() error {
	return cli.business.Save()
}

// SetDeveloper set the biz developer
func (cli *BusinessWrapper) SetDeveloper(developer string) error {
	return cli.business.SetValue(fieldBizDeveloper, developer)
}

// SetMaintainer set the maintainer for the biz
func (cli *BusinessWrapper) SetMaintainer(maintainer string) error {
	return cli.business.SetValue(fieldBizMaintainer, maintainer)
}

// SetName set the  business name
func (cli *BusinessWrapper) SetName(name string) error {
	return cli.business.SetValue(fieldBizName, name)
}

// SetProductor set the productor name
func (cli *BusinessWrapper) SetProductor(productor string) error {
	return cli.business.SetValue(fieldBizProductor, productor)
}

// SetTester set the tester name
func (cli *BusinessWrapper) SetTester(tester string) error {
	return cli.business.SetValue(fieldBizTester, tester)
}

// SetSupplierAccount set the supplier account
func (cli *BusinessWrapper) SetSupplierAccount(supplierAccount string) error {
	return cli.business.SetValue(fieldSupplierAccount, supplierAccount)
}

// SetLifeCycle set the life cycle
func (cli *BusinessWrapper) SetLifeCycle(lifeCycle string) error {
	return cli.business.SetValue(fieldLifeCycle, lifeCycle)
}

// SetOperator set the operator
func (cli *BusinessWrapper) SetOperator(operator string) error {
	return cli.business.SetValue(fieldBizOperator, operator)
}
