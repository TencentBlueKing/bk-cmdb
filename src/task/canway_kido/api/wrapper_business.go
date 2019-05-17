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

// BusinessIteratorWrapper the business iterator wrapper
type BusinessIteratorWrapper struct {
	business inst.BusinessIterator
}

// Next next the business
func (cli *BusinessIteratorWrapper) Next() (*BusinessWrapper, error) {

	busi, err := cli.business.Next()

	return &BusinessWrapper{business: busi}, err

}

// ForEach the foreach function
func (cli *BusinessIteratorWrapper) ForEach(callback func(business *BusinessWrapper) error) error {

	return cli.business.ForEach(func(item inst.BusinessInterface) error {
		return callback(&BusinessWrapper{business: item})
	})
}

// BusinessWrapper the business wrapper
type BusinessWrapper struct {
	business inst.BusinessInterface
}

// SetValue set the key value
func (cli *BusinessWrapper) SetValue(key string, val interface{}) error {
	return cli.business.SetValue(key, val)
}

// IsExists check the set
func (cli *BusinessWrapper) IsExists() (bool, error) {
	return cli.business.IsExists()
}

// Create only to create
func (cli *BusinessWrapper) Create() error {
	return cli.business.Create()
}

// Update only to update
func (cli *BusinessWrapper) Update() error {
	return cli.business.Update()
}

// Save save the data
func (cli *BusinessWrapper) Save() error {
	return cli.business.Save()
}

// GetValues return the values
func (cli *BusinessWrapper) GetValues() (types.MapStr, error) {
	return cli.business.GetValues()
}

// GetModel get the model for the business
func (cli *BusinessWrapper) GetModel() model.Model {
	return cli.business.GetModel()
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

// GetBusinessID get the id for the business
func (cli *BusinessWrapper) GetBusinessID() (int64, error) {
	vals, err := cli.business.GetValues()
	if nil != err {
		return 0, err
	}
	if !vals.Exists(fieldBusinessID) {
		return 0, fmt.Errorf("the business id is not set")
	}
	val, err := vals.Int(fieldBusinessID)
	return int64(val), err
}

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
