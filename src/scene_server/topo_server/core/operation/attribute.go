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

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(params types.LogicParams, data frtypes.MapStr) (model.Attribute, error)
	DeleteObjectAttribute(params types.LogicParams, cond condition.Condition) error
	FindObjectAttribute(params types.LogicParams, cond condition.Condition) ([]model.Attribute, error)
	UpdateObjectAttribute(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
}

type attribute struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

// NewAttributeOperation create a new attribute operation instance
func NewAttributeOperation(client apimachinery.ClientSetInterface, modelFactory model.Factory, instFactory inst.Factory) AttributeOperationInterface {
	return &attribute{
		clientSet:    client,
		modelFactory: modelFactory,
		instFactory:  instFactory,
	}
}

func (cli *attribute) CreateObjectAttribute(params types.LogicParams, data frtypes.MapStr) (model.Attribute, error) {
	att := cli.modelFactory.CreateAttribute(params)

	_, err := att.Parse(data)
	if nil != err {
		return nil, err
	}

	err = att.Save()
	if nil != err {
		return nil, err
	}

	return att, nil
}

func (cli *attribute) DeleteObjectAttribute(params types.LogicParams, cond condition.Condition) error {
	return nil
}

func (cli *attribute) FindObjectAttribute(params types.LogicParams, cond condition.Condition) ([]model.Attribute, error) {
	return nil, nil
}

func (cli *attribute) UpdateObjectAttribute(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
