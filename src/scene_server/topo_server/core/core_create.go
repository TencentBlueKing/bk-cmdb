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

package core

import (
	//frcommon "configcenter/src/framework/common"
	frtypes "configcenter/src/common/types"
	"configcenter/src/scene_server/topo_server/core/types"

	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

func (cli *core) CreateClassification(params types.LogicParams, data frtypes.MapStr) (model.Classification, error) {

	cls := cli.modelFactory.CreaetClassification(params)

	err := cls.Parse(data)
	if nil != err {
		return nil, err
	}

	err = cls.Save()
	if nil != err {
		return nil, err
	}

	return cls, nil
}

func (cli *core) CreateObject(params types.LogicParams, data frtypes.MapStr) (model.Object, error) {

	obj := cli.modelFactory.CreaetObject(params)

	err := obj.Parse(data)
	if nil != err {
		return nil, err
	}

	err = obj.Save()
	if nil != err {
		return nil, err
	}

	return obj, nil
}

func (cli *core) CreateObjectAttribute(params types.LogicParams, data frtypes.MapStr) (model.Attribute, error) {

	att := cli.modelFactory.CreateAttribute(params)

	err := att.Parse(data)
	if nil != err {
		return nil, err
	}

	err = att.Save()
	if nil != err {
		return nil, err
	}

	return att, nil
}

func (cli *core) CreateObjectGroup(params types.LogicParams, data frtypes.MapStr) (model.Group, error) {

	grp := cli.modelFactory.CreateGroup(params)

	err := grp.Parse(data)
	if nil != err {
		return nil, err
	}

	err = grp.Save()
	if nil != err {
		return nil, err
	}

	return grp, nil
}

func (cli *core) CreateInst(params types.LogicParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error) {

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

func (cli *core) CreateAssociation(params types.LogicParams, data frtypes.MapStr) (model.Association, error) {
	return nil, nil
}
