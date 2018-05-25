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
	frtypes "configcenter/src/framework/core/types"

	"configcenter/src/scene_server/core/inst"
	"configcenter/src/scene_server/core/model"
)

func (cli *core) CreateClassification(data frtypes.MapStr) (model.Classification, error) {

	cls := cli.modelMgr.CreaetClassification()

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

func (cli *core) CreateObject(data frtypes.MapStr) (model.Object, error) {

	obj := cli.modelMgr.CreaetObject()

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

func (cli *core) CreateObjectAttribute(data frtypes.MapStr) (model.Attribute, error) {

	att := cli.modelMgr.CreateAttribute()

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

func (cli *core) CreateObjectGroup(data frtypes.MapStr) (model.Group, error) {

	grp := cli.modelMgr.CreateObjectGroup()

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

func (cli *core) CreateInst(data frtypes.MapStr) (inst.Inst, error) {

	return nil, nil
}

func (cli *core) CreateAssociation(data frtypes.MapStr) (model.Association, error) {
	return nil, nil
}
