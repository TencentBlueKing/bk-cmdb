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
	"errors"

	"configcenter/src/framework/common"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ ModuleInterface = (*module)(nil)

// ModuleInterface the module interface
type ModuleInterface interface {
	Maintaince

	SetTopo(bizID, setID int64)

	GetModel() model.Model

	GetInstID() (int64, error)
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)
}

type module struct {
	bizID  int64
	setID  int64
	target model.Model
	datas  types.MapStr
}

func (cli *module) SetTopo(bizID, setID int64) {
	cli.setID = setID
	cli.bizID = bizID

}
func (cli *module) GetModel() model.Model {
	return cli.target
}

func (cli *module) GetInstID() (int64, error) {
	return cli.datas.Int64(ModuleID)
}
func (cli *module) GetInstName() string {
	return cli.datas.String(ModuleName)
}

func (cli *module) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *module) SetValue(key string, value interface{}) error {

	// TODO:需要根据model 的定义对输入的key 及value 进行校验

	cli.datas.Set(key, value)

	return nil
}

func (cli *module) search() ([]model.Attribute, []types.MapStr, error) {

	// get the attributes
	attrs, err := cli.target.Attributes()
	if nil != err {
		return nil, nil, err
	}

	// construct the condition which is used to check the if it is exists
	cond := common.CreateCondition().Field(BusinessID).Eq(cli.bizID).Field(SetID).Eq(cli.setID)

	// extract the required id
	for _, attrItem := range attrs {
		if attrItem.GetKey() {

			if attrItem.GetID() == BusinessID {
				if 0 >= cli.bizID {
					return nil, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
				}
				cond.Field(BusinessID).Eq(cli.bizID)
				continue
			}

			if attrItem.GetID() == SetID {
				if 0 >= cli.setID {
					return nil, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
				}
				cond.Field(SetID).Eq(cli.setID)
				continue
			}

			attrVal := cli.datas.String(attrItem.GetID())
			if 0 == len(attrVal) {
				return attrs, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
			}

			cond.Field(attrItem.GetID()).Eq(attrVal)
		}
	}

	//log.Infof("the module search condition:%#v", cond.ToMapStr())
	// search by condition
	existItems, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Module().SearchModules(cond)

	return attrs, existItems, err

}

func (cli *module) IsExists() (bool, error) {

	_, items, err := cli.search()
	if nil != err {
		return false, err
	}
	return 0 != len(items), nil
}
func (cli *module) Create() error {

	if 0 <= cli.setID {
		cli.datas.Set(ParentID, cli.setID)
		cli.datas.Set(SetID, cli.setID)
	}

	if 0 < cli.bizID {
		cli.datas.Set(BusinessID, cli.bizID)
	}

	moduleID, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Module().CreateModule(cli.bizID, cli.setID, cli.datas)
	if nil != err {
		return err
	}
	cli.datas.Set(ModuleID, moduleID)
	return nil
}
func (cli *module) Update() error {

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

	if 0 < cli.setID {
		cli.datas.Set(ParentID, cli.setID)
		cli.datas.Set(SetID, cli.setID)
	}
	if 0 < cli.bizID {
		cli.datas.Set(BusinessID, cli.bizID)
	}
	cli.datas.Remove("create_time") //invalid check , need to delete
	for _, existItem := range existItems {

		instID, err := existItem.Int64(ModuleID)
		if nil != err {
			return err
		}

		cli.datas.Remove(ModuleID)

		err = client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).Module().UpdateModule(cli.bizID, cli.setID, instID, cli.datas)
		if nil != err {
			log.Infof("failed to  update the module (%#v), error info is %s", existItem, err.Error())
			return err
		}

		cli.datas.Set(ModuleID, instID)

	}
	return nil
}
func (cli *module) Save() error {

	//fmt.Println("bizID:", cli.bizID, "setID:", cli.setID)
	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()
}
