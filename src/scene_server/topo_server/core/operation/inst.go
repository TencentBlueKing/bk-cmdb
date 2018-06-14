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

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	CreateInst(params types.LogicParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error)
	DeleteInst(params types.LogicParams, cond condition.Condition) error
	FindInst(params types.LogicParams, cond condition.Condition) ([]inst.Inst, error)
	UpdateInst(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
}

type commonInst struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

// NewInstOperation create a new inst operation instance
func NewInstOperation(client apimachinery.ClientSetInterface, modelFactory model.Factory, instFactory inst.Factory) InstOperationInterface {
	return &commonInst{
		clientSet:    client,
		modelFactory: modelFactory,
		instFactory:  instFactory,
	}
}

func (cli *commonInst) CreateInst(params types.LogicParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error) {
	item := cli.instFactory.CreateInst(params, obj)

	err := item.SetValues(data)
	if nil != err {
		return nil, err
	}

	err = item.Save()
	if nil != err {
		return nil, err
	}

	return item, nil
}

func (cli *commonInst) DeleteInst(params types.LogicParams, cond condition.Condition) error {
	return nil
}

func (cli *commonInst) FindInst(params types.LogicParams, cond condition.Condition) ([]inst.Inst, error) {
	return nil, nil
}

func (cli *commonInst) UpdateInst(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
