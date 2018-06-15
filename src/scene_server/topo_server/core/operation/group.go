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
	"configcenter/src/common"
	"configcenter/src/common/metadata"
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// GroupOperationInterface group operation methods
type GroupOperationInterface interface {
	CreateObjectGroup(params types.LogicParams, data frtypes.MapStr) (model.Group, error)
	DeleteObjectGroup(params types.LogicParams, groupID string) error
	FindObjectGroup(params types.LogicParams, cond condition.Condition) ([]model.Group, error)
	FindGroupByObject(params types.LogicParams, objID string, cond condition.Condition) ([]model.Group, error)
	UpdateObjectGroup(params types.LogicParams, cond *metadata.UpdateGroupCondition) error
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
		blog.Errorf("[operation-grp] failed to parse the group data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	err = grp.Save()
	if nil != err {
		blog.Errorf("[operation-grp] failed to save the group data (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	return grp, nil
}

func (cli *group) DeleteObjectGroup(params types.LogicParams, groupID string) error {

	rsp, err := cli.clientSet.ObjectController().Meta().DeletePropertyGroup(context.Background(), groupID, params.Header.ToHeader())
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object conroller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-grp] failed to delete the group by the groupid (%s), error info is %s", groupID, err.Error())
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (cli *group) FindObjectGroup(params types.LogicParams, cond condition.Condition) ([]model.Group, error) {

	rsp, err := cli.clientSet.ObjectController().Meta().SelectGroup(context.Background(), params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[opeartion-grp] failed to search the group by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return nil, nil
}

func (cli *group) FindGroupByObject(params types.LogicParams, objID string, cond condition.Condition) ([]model.Group, error) {

	rsp, err := cli.clientSet.ObjectController().Meta().SelectPropertyGroupByObjectID(context.Background(), params.Header.OwnerID, objID, params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-grp] failed to search the group of the object(%s) by the condition (%#v), error info is %s", objID, cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}
	return nil, nil
}
func (cli *group) UpdateObjectGroup(params types.LogicParams, cond *metadata.UpdateGroupCondition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().UpdatePropertyGroup(context.Background(), params.Header.ToHeader(), cond)

	if nil != err {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s ", cond.Data, cond.Condition, err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s ", cond.Data, cond.Condition, rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}
