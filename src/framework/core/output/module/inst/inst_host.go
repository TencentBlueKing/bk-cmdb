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
	"configcenter/src/framework/core/errors"
	"configcenter/src/framework/core/log"
	"configcenter/src/framework/core/output/module/client"
	"configcenter/src/framework/core/output/module/model"
	"configcenter/src/framework/core/types"
)

var _ Inst = (*host)(nil)

type host struct {
	target model.Model
	datas  types.MapStr
}

func (cli *host) GetModel() model.Model {
	return cli.target
}

func (cli *host) IsMainLine() bool {
	return false
}

func (cli *host) GetAssociationModels() ([]model.Model, error) {
	// TODO:需要读取此实例关联的实例，所对应的所有模型
	return nil, nil
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

func (cli *host) GetAssociationsByModleID(modleID string) ([]Inst, error) {
	// TODO:获取当前实例所关联的特定模型的所有已关联的实例
	return nil, nil
}

func (cli *host) GetAllAssociations() (map[model.Model][]Inst, error) {
	// TODO:获取所有已关联的模型及对应的实例
	return nil, nil
}

func (cli *host) SetParent(parentInstID int) error {
	return errors.ErrNotSuppportedFunctionality
}

func (cli *host) GetParent() ([]Topo, error) {
	return nil, errors.ErrNotSuppportedFunctionality
}

func (cli *host) GetChildren() ([]Topo, error) {
	return nil, errors.ErrNotSuppportedFunctionality
}

func (cli *host) SetValue(key string, value interface{}) error {
	cli.datas.Set(key, value)
	return nil
}
func (cli *host) search() ([]model.Attribute, []types.MapStr, error) {
	return nil, nil, nil
}
func (cli *host) IsExists() (bool, error) {
	return true, nil
}
func (cli *host) Create() error {
	return nil
}
func (cli *host) Update() error {
	return nil
}
func (cli *host) Save() error {

	attrs, err := cli.target.Attributes()
	if nil != err {
		return err
	}

	// clear the invalid data
	for _, attrItem := range attrs {
		if !cli.datas.Exists(attrItem.GetID()) {
			cli.datas.Remove(attrItem.GetID())
		}
	}

	hostID, err := client.GetClient().CCV3().Host().CreateHostBatch(cli.datas)
	if err != nil {
		return err
	}
	if len(hostID) != 1 {
		return errors.New("incorrect id number received")
	}
	cli.datas.Set(cccommon.BKHostIDField, hostID[0])
	return nil
}
