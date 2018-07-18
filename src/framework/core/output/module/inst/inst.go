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

var _ CommonInstInterface = (*inst)(nil)

// CommonInstInterface the inst interface
type CommonInstInterface interface {
	Maintaince

	GetModel() model.Model

	GetInstID() int
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)
}

type inst struct {
	target model.Model
	datas  types.MapStr
}

func (cli *inst) GetModel() model.Model {
	return cli.target
}

func (cli *inst) GetInstID() int {
	id, err := cli.datas.Int(InstID)
	if nil != err {
		log.Errorf("failed to get the inst id, %s", err.Error())
	}
	return id
}
func (cli *inst) GetInstName() string {
	return cli.datas.String(InstName)
}

func (cli *inst) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *inst) SetValue(key string, value interface{}) error {

	// TODO:需要增加对输入的key 的校验

	cli.datas.Set(key, value)

	return nil
}

func (cli *inst) search() ([]model.Attribute, []types.MapStr, error) {

	// get the attributes
	attrs, err := cli.target.Attributes()
	if nil != err {
		return nil, nil, err
	}

	// construct the condition which is used to check the if it is exists
	cond := common.CreateCondition()

	// extract the object id
	objID := cli.target.GetID()
	cond.Field(model.ObjectID).Eq(objID)
	//log.Infof("attrs:%#v", attrs)
	// extract the required id
	for _, attrItem := range attrs {
		//log.Infof("attrs:%#v", attrItem)
		if attrItem.GetKey() {

			attrVal := cli.datas.String(attrItem.GetID())
			if 0 == len(attrVal) {
				return attrs, nil, errors.New("the key field(" + attrItem.GetID() + ") is not set")
			}

			cond.Field(attrItem.GetID()).Eq(attrVal)
		}
	}

	// search by condition
	existItems, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).CommonInst().SearchInst(cond)
	return attrs, existItems, err
}

func (cli *inst) IsExists() (bool, error) {

	_, existItems, err := cli.search()
	if nil != err {
		return false, err
	}

	return 0 != len(existItems), nil
}
func (cli *inst) Create() error {

	instID, err := client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).CommonInst().CreateCommonInst(cli.datas)
	if nil != err {
		return err
	}
	cli.datas.Set(InstID, instID)
	return nil
}
func (cli *inst) Update() error {

	attrs, existItems, err := cli.search()
	if nil != err {
		return err
	}

	targetInstID := InstID
	objID := cli.target.GetID()
	switch objID {
	case Plat:
		targetInstID = PlatID
	}

	// update the exists
	for _, existItem := range existItems {

		cli.datas.ForEach(func(key string, val interface{}) {
			existItem.Set(key, val)
		})

		instID, err := existItem.Int(targetInstID)
		if nil != err {
			return err
		}

		updateCond := common.CreateCondition().Field(targetInstID).Eq(instID).Field(model.ObjectID).Eq(objID)

		// clear the invalid field
		existItem.ForEach(func(key string, val interface{}) {
			for _, attrItem := range attrs {
				if attrItem.GetID() == key {
					return
				}
			}
			existItem.Remove(key)
		})

		err = client.GetClient().CCV3(client.Params{SupplierAccount: cli.target.GetSupplierAccount()}).CommonInst().UpdateCommonInst(existItem, updateCond)
		if nil != err {
			return err
		}
		cli.datas.Set(InstID, instID)
	}

	return nil
}
func (cli *inst) Save() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return cli.Update()
	}

	return cli.Create()
}
