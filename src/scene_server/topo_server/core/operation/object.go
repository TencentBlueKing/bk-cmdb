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

package operation

import (
	"configcenter/src/apimachinery"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObject(params types.LogicParams, data frtypes.MapStr) (model.Object, error)
	DeleteObject(params types.LogicParams, cond condition.Condition) error
	FindObject(params types.LogicParams, cond condition.Condition) ([]model.Object, error)
	UpdateObject(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
}

type object struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

// NewObjectOperation create a new object operation instance
func NewObjectOperation(client apimachinery.ClientSetInterface, modelFactory model.Factory, instFactory inst.Factory) ObjectOperationInterface {
	return &object{
		clientSet:    client,
		modelFactory: modelFactory,
		instFactory:  instFactory,
	}
}

func (cli *object) CreateObject(params types.LogicParams, data frtypes.MapStr) (model.Object, error) {
	obj := cli.modelFactory.CreaetObject(params)

	_, err := obj.Parse(data)
	if nil != err {
		return nil, err
	}

	err = obj.Save()
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (cli *object) DeleteObject(params types.LogicParams, cond condition.Condition) error {
	return nil
}

func (cli *object) FindObject(params types.LogicParams, cond condition.Condition) ([]model.Object, error) {
	return nil, nil
}

func (cli *object) UpdateObject(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
