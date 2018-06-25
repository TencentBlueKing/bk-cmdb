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
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
	"errors"
	//"fmt"
)

var _ SetInterface = (*set)(nil)

// SetInterface the set interface
type SetInterface interface {
	Maintaince

	GetModel() model.Model

	GetInstID() int
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)
}

type set struct {
	target model.Model
	datas  types.MapStr
}

func (cli *set) GetModel() model.Model {
	return cli.target
}

func (cli *set) IsMainLine() bool {
	// TODO：判断当前实例是否为主线实例
	return true
}

func (cli *set) GetAssociationModels() ([]model.Model, error) {
	// TODO:需要读取此实例关联的实例，所对应的所有模型
	return nil, nil
}

func (cli *set) GetInstID() int {
	id, err := cli.datas.Int(SetID)
	if nil != err {
		log.Errorf("failed to get the inst id, %s", err.Error())
	}
	return id
}
func (cli *set) GetInstName() string {
	return cli.datas.String(SetName)
}

func (cli *set) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *set) SetValue(key string, value interface{}) error {

	// TODO:需要根据model 的定义对输入的key 及value 进行校验

	cli.datas.Set(key, value)

	return nil
}
func (cli *set) search() ([]model.Attribute, []types.MapStr, error) {
	businessID := cli.datas.String(BusinessID)

	// get the attributes
	attrs, err := cli.target.Attributes()
	if nil != err {
		return nil, nil, err
	}

	// construct the condition which is used to check the if it is exists
	cond := common.CreateCondition().Field(BusinessID).Eq(businessID)

	// extract the required id
	for _, attrItem := range attrs {
		if attrItem.GetKey() {
			if !cli.datas.Exists(attrItem.GetID()) {
				return attrs, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
			}
			attrVal := cli.datas.String(attrItem.GetID())
			cond.Field(attrItem.GetID()).Eq(attrVal)
		}
	}

	//fmt.Println("cond:", cond.ToMapStr())

	// search by condition
	existItems, err := client.GetClient().CCV3().Set().SearchSets(cond)

	return attrs, existItems, err
}
func (cli *set) IsExists() (bool, error) {

	_, existItems, err := cli.search()
	if nil != err {
		return false, err
	}

	return 0 != len(existItems), nil
}

func (cli *set) Create() error {

	setID, err := client.GetClient().CCV3().Set().CreateSet(cli.datas)
	if nil != err {
		return err
	}
	cli.datas.Set(SetID, setID)
	return nil
}
func (cli *set) Update() error {

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

	for _, existItem := range existItems {

		instID, err := existItem.Int64(SetID)
		if nil != err {
			return err
		}

		updateCond := common.CreateCondition()
		updateCond.Field(SetID).Eq(instID)

		err = client.GetClient().CCV3().Set().UpdateSet(cli.datas, updateCond)
		if nil != err {
			return err
		}

	}
	return nil
}
func (cli *set) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()

}
