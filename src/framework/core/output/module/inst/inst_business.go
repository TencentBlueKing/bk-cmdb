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

package inst

import (
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/errors"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ BusinessInterface = (*business)(nil)

// BusinessInterface the business interface
type BusinessInterface interface {
	Maintaince

	GetModel() model.Model

	GetInstID() (int, error)
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)
}

type business struct {
	target model.Model
	datas  types.MapStr
}

func (cli *business) GetModel() model.Model {
	return cli.target
}

func (cli *business) GetInstID() (int, error) {
	return cli.datas.Int(BusinessID)
}
func (cli *business) GetInstName() string {

	return cli.datas.String(BusinessNameField)
}

func (cli *business) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *business) SetValue(key string, value interface{}) error {

	// TODO:需要根据model 的定义对输入的key 及value 进行校验

	cli.datas[key] = value

	return nil
}
func (cli *business) search() ([]model.Attribute, []types.MapStr, error) {
	attrs, err := cli.target.Attributes()
	if nil != err {
		return nil, nil, err
	}

	cond := common.CreateCondition()
	for _, attrItem := range attrs {
		if attrItem.GetKey() {

			attrVal := cli.datas.String(attrItem.GetID())
			if 0 == len(attrVal) {
				return nil, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
			}

			cond.Field(attrItem.GetID()).Eq(attrVal)
		}
	}

	// search by condition
	items, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Business().SearchBusiness(cond)
	return attrs, items, err
}
func (cli *business) IsExists() (bool, error) {

	// search by condition
	_, existItems, err := cli.search()
	if nil != err {
		return false, err
	}

	return 0 != len(existItems), nil
}

func (cli *business) Create() error {

	bizID, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Business().CreateBusiness(cli.datas)
	if err == nil {
		cli.datas.Set(BusinessID, bizID)
		return nil
	}
	cli.datas.Set(BusinessID, bizID)
	return err

}
func (cli *business) Update() error {

	attrs, existItems, err := cli.search()
	if nil != err {
		return err
	}

	// clear the invalid field
	cli.datas.ForEach(func(key string, val interface{}) {
		for _, attrItem := range attrs {
			if attrItem.GetID() == key {
				return
			}
		}
		cli.datas.Remove(key)
	})

	cli.datas.Remove("create_time") //invalid check , need to delete

	supplierAccount := cli.target.GetSupplierAccount()
	// update the exists
	for _, existItem := range existItems {

		bizID, err := existItem.Int(BusinessID)
		if nil != err {
			return err
		}

		cli.datas.Remove(BusinessID)

		//fmt.Println("the new:", existItem)
		err = client.GetClient().CCV3(client.Params{SupplierAccount: supplierAccount}).Business().UpdateBusiness(cli.datas, bizID)
		if nil != err {
			return err
		}

		cli.datas.Set(BusinessID, bizID)

	}

	return nil
}

func (cli *business) Save() error {
	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()
}
