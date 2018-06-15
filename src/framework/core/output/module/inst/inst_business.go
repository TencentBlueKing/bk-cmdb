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
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ Inst = (*business)(nil)

type business struct {
	target model.Model
	datas  types.MapStr
}

func (cli *business) GetModel() model.Model {
	return cli.target
}

func (cli *business) IsMainLine() bool {
	return false
}

func (cli *business) GetAssociationModels() ([]model.Model, error) {
	// TODO:需要读取此实例关联的实例，所对应的所有模型
	return nil, nil
}

func (cli *business) GetInstID() int {
	instID, err := cli.datas.Int(BusinessID)
	if err != nil {
		log.Errorf("get bk_biz_id faile %v", err)
	}
	return instID
}
func (cli *business) GetInstName() string {

	return cli.datas.String(BusinessNameField)
}

func (cli *business) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *business) GetAssociationsByModleID(modleID string) ([]Inst, error) {
	// TODO:获取当前实例所关联的特定模型的所有已关联的实例
	return nil, nil
}

func (cli *business) GetAllAssociations() (map[model.Model][]Inst, error) {
	// TODO:获取所有已关联的模型及对应的实例
	return nil, nil
}

func (cli *business) SetParent(parentInstID int) error {
	return errors.ErrNotSuppportedFunctionality
}

func (cli *business) GetParent() ([]Topo, error) {
	return nil, errors.ErrNotSuppportedFunctionality
}

func (cli *business) GetChildren() ([]Topo, error) {
	return nil, errors.ErrNotSuppportedFunctionality
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
	items, err := client.GetClient().CCV3().Business().SearchBusiness(cond)
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

	bizID, err := client.GetClient().CCV3().Business().CreateBusiness(cli.datas)
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

	// update the exists
	for _, existItem := range existItems {

		cli.datas.ForEach(func(key string, val interface{}) {
			existItem.Set(key, val)
		})

		instID, err := existItem.Int(BusinessID)
		if nil != err {
			return err
		}
		// clear the invalid field
		existItem.ForEach(func(key string, val interface{}) {
			for _, attrItem := range attrs {
				if attrItem.GetID() == key {
					return
				}
			}
			existItem.Remove(key)
		})
		//fmt.Println("the new:", existItem)
		err = client.GetClient().CCV3().Business().UpdateBusiness(existItem, instID)
		if nil != err {
			return err
		}
		cli.datas.Set(BusinessID, instID)

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
