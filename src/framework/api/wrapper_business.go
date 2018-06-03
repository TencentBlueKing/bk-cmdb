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

// BusinessIteratorWrapper the business iterator wrapper
type BusinessIteratorWrapper struct {
	business inst.Iterator
}

// Next next the business
func (cli *BusinessIteratorWrapper) Next() (*BusinessWrapper, error) {

	busi, err := cli.business.Next()

	return &BusinessWrapper{business: busi}, err

}

// ForEach the foreach function
func (cli *BusinessIteratorWrapper) ForEach(callback func(business *BusinessWrapper) error) error {

	return cli.business.ForEach(func(item inst.Inst) error {
		return callback(&BusinessWrapper{business: item})
	})
}

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

// GetDeveloper get the developer name
func (cli *BusinessWrapper) GetDeveloper() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldBizDeveloper), nil
}

// SetMaintainer set the maintainer for the biz
func (cli *BusinessWrapper) SetMaintainer(maintainer string) error {
	return cli.business.SetValue(fieldBizMaintainer, maintainer)
}

// GetMaintainer get the maintaner name
func (cli *BusinessWrapper) GetMaintainer() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldBizMaintainer), nil
}

// SetName set the  business name
func (cli *BusinessWrapper) SetName(name string) error {
	return cli.business.SetValue(fieldBizName, name)
}

// GetName get the business name
func (cli *BusinessWrapper) GetName() (string, error) {
	return cli.business.GetInstName(), nil
}

// SetProductor set the productor name
func (cli *BusinessWrapper) SetProductor(productor string) error {
	return cli.business.SetValue(fieldBizProductor, productor)
}

// GetProductor get the productor for the business
func (cli *BusinessWrapper) GetProductor() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldBizProductor), nil
}

// SetTester set the tester name
func (cli *BusinessWrapper) SetTester(tester string) error {
	return cli.business.SetValue(fieldBizTester, tester)
}

// GetTester get the tester name
func (cli *BusinessWrapper) GetTester() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldBizTester), nil
}

/* TODO need to delete the follow code
// SetSupplierAccount set the supplier account
func (cli *BusinessWrapper) SetSupplierAccount(supplierAccount string) error {
	id, _ := strconv.Atoi(supplierAccount)
	cli.SetValue(fieldSupplierID, id)
	return cli.business.SetValue(fieldSupplierAccount, supplierAccount)
}

// GetSupplierAccount get the supplier account
func (cli *BusinessWrapper) GetSupplierAccount() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldSupplierAccount), nil
}
*/

// SetLifeCycle set the life cycle
func (cli *BusinessWrapper) SetLifeCycle(lifeCycle string) error {
	return cli.business.SetValue(fieldLifeCycle, lifeCycle)
}

// GetLifeCycle get the life cycle
func (cli *BusinessWrapper) GetLifeCycle() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldLifeCycle), nil
}

// SetOperator set the operator
func (cli *BusinessWrapper) SetOperator(operator string) error {
	return cli.business.SetValue(fieldBizOperator, operator)
}

// GetOperator get the operator
func (cli *BusinessWrapper) GetOperator() (string, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return "", err
	}
	return vals.String(fieldBizOperator), nil
}
