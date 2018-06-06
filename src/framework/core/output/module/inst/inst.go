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
)

var _ Inst = (*inst)(nil)

type inst struct {
	target model.Model
	datas  types.MapStr
}

func (cli *inst) GetModel() model.Model {
	return cli.target
}

func (cli *inst) IsMainLine() bool {
	// TODO：判断当前实例是否为主线实例
	return true
}

func (cli *inst) GetAssociationModels() ([]model.Model, error) {
	// TODO:需要读取此实例关联的实例，所对应的所有模型
	return nil, nil
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

func (cli *inst) GetAssociationsByModleID(modleID string) ([]Inst, error) {
	// TODO:获取当前实例所关联的特定模型的所有已关联的实例
	return nil, nil
}

func (cli *inst) GetAllAssociations() (map[model.Model][]Inst, error) {
	// TODO:获取所有已关联的模型及对应的实例
	return nil, nil
}

func (cli *inst) SetParent(parentInstID int) error {
	return nil
}

func (cli *inst) GetParent() ([]Topo, error) {
	return nil, nil
}

func (cli *inst) GetChildren() ([]Topo, error) {
	return nil, nil
}

func (cli *inst) SetValue(key string, value interface{}) error {

	// TODO:需要增加对输入的key 的校验

	cli.datas.Set(key, value)

	return nil
}

func (cli *inst) Save() error {
	// get the attributes
	attrs, err := cli.target.Attributes()
	if nil != err {
		return err
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
				return errors.New("the key field(" + attrItem.GetID() + ") is not set")
			}

			cond.Field(attrItem.GetID()).Eq(attrVal)
		}
	}

	log.Infof("cond:%#v", cond.ToMapStr())

	// search by condition
	existItems, err := client.GetClient().CCV3().CommonInst().SearchInst(cond)
	if nil != err {
		return err
	}

	log.Infof("the exists:%#v", cli.target.GetID())

	// create a new
	if 0 == len(existItems) {
		_, err = client.GetClient().CCV3().CommonInst().CreateCommonInst(cli.datas)
		return err
	}

	targetInstID := InstID
	switch cli.target.GetID() {
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

		err = client.GetClient().CCV3().CommonInst().UpdateCommonInst(existItem, updateCond)
		if nil != err {
			return err
		}

	}

	return nil
}
