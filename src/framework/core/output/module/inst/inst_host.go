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
	cccommon "configcenter/src/common"
	"configcenter/src/framework/common"
	"configcenter/src/framework/core/errors"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ HostInterface = (*host)(nil)

// HostInterface the host interface
type HostInterface interface {
	IsExists() (bool, error)
	Create() error
	Update() error
	Save() error

	GetModel() model.Model

	GetInstID() int
	GetInstName() string

	SetValue(key string, value interface{}) error
	GetValues() (types.MapStr, error)
}

type host struct {
	bizID     int64
	moduleIDS []int64
	target    model.Model
	datas     types.MapStr
}

func (cli *host) GetModel() model.Model {
	return cli.target
}

func (cli *host) GetInstID() int {
	instID, err := cli.datas.Int(cccommon.BKHostIDField)
	if err != nil {
		log.Errorf("get bk_host_id faile %v", err)
	}
	return instID
}

func (cli *host) GetInstName() string {
	return cli.datas.String(cccommon.BKHostIDField)
}

func (cli *host) GetValues() (types.MapStr, error) {
	return cli.datas, nil
}

func (cli *host) SetValue(key string, value interface{}) error {
	cli.datas.Set(key, value)
	return nil
}
func (cli *host) search() ([]model.Attribute, []types.MapStr, error) {

	return nil, nil, nil
}
func (cli *host) IsExists() (bool, error) {

	attrs, err := cli.target.Attributes()
	if nil != err {
		return false, err
	}

	cond := common.CreateCondition()

	for _, attrItem := range attrs {

		if !attrItem.GetKey() {
			continue
		}

		if !cli.datas.Exists(attrItem.GetID()) {
			continue
		}

		cond.Field(attrItem.GetID()).Eq(cli.datas[attrItem.GetID()])

	}

	items, err := client.GetClient().CCV3().Host().SearchHost(cond)
	if nil != err {
		return false, err
	}

	return 0 != len(items), nil
}
func (cli *host) Create() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if exists {
		return nil
	}

	return cli.Save()

}
func (cli *host) Update() error {

	if exists, err := cli.IsExists(); nil != err {
		return err
	} else if !exists {
		return nil
	}

	return cli.Save()
}
func (cli *host) Save() error {

	attrs, err := cli.target.Attributes()
	if nil != err {
		return err
	}

	bizID := 0
	if cli.datas.Exists(BusinessID) {
		id, err := cli.datas.Int(BusinessID)
		if nil != err {
			return err
		}
		bizID = id
	}

	log.Errorf("set bizid:%d", bizID)

	// clear the invalid data
	for _, attrItem := range attrs {
		if !cli.datas.Exists(attrItem.GetID()) {
			cli.datas.Remove(attrItem.GetID())
		}
	}

	// clear the undefined data
	cli.datas.ForEach(func(key string, val interface{}) {
		for _, attrItem := range attrs {
			if attrItem.GetID() == key {
				return
			}
		}

		cli.datas.Remove(key)
	})

	hostID, err := client.GetClient().CCV3().Host().CreateHostBatch(cli.bizID, cli.moduleIDS, cli.datas)
	if err != nil {
		log.Errorf("failed to create host, error info is %s", err.Error())
		return err
	}
	if len(hostID) != 1 {
		return errors.New("incorrect id number received")
	}
	cli.datas.Set(cccommon.BKHostIDField, hostID[0])
	return nil
}
