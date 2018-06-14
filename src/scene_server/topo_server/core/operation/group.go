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

// GroupOperationInterface group operation methods
type GroupOperationInterface interface {
	CreateObjectGroup(params types.LogicParams, data frtypes.MapStr) (model.Group, error)
	DeleteObjectGroup(params types.LogicParams, cond condition.Condition) error
	FindObjectGroup(params types.LogicParams, cond condition.Condition) ([]model.Group, error)
	FindGroupByObject(params types.LogicParams, cond condition.Condition) ([]model.Group, error)
	UpdateObjectGroup(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error
}

type group struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

// NewGroupOperation create a new group operation instance
func NewGroupOperation(client apimachinery.ClientSetInterface, modelFactory model.Factory, instFactory inst.Factory) GroupOperationInterface {
	return &group{
		clientSet:    client,
		modelFactory: modelFactory,
		instFactory:  instFactory,
	}
}

func (cli *group) CreateObjectGroup(params types.LogicParams, data frtypes.MapStr) (model.Group, error) {
	grp := cli.modelFactory.CreateGroup(params)

	_, err := grp.Parse(data)
	if nil != err {
		return nil, err
	}

	err = grp.Save()
	if nil != err {
		return nil, err
	}

	return grp, nil
}

func (cli *group) DeleteObjectGroup(params types.LogicParams, cond condition.Condition) error {
	return nil
}

func (cli *group) FindObjectGroup(params types.LogicParams, cond condition.Condition) ([]model.Group, error) {
	return nil, nil
}

func (cli *group) FindGroupByObject(params types.LogicParams, cond condition.Condition) ([]model.Group, error) {
	return nil, nil
}
func (cli *group) UpdateObjectGroup(params types.LogicParams, data frtypes.MapStr, cond condition.Condition) error {
	return nil
}
