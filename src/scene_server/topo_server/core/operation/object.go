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
	"context"
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObject(params types.LogicParams, data frtypes.MapStr) (model.Object, error)
	DeleteObject(params types.LogicParams, id int64, cond condition.Condition) error
	FindObject(params types.LogicParams, cond condition.Condition) ([]model.Object, error)
	FindSingleObject(params types.LogicParams, objectID string) (model.Object, error)
	UpdateObject(params types.LogicParams, data frtypes.MapStr, id int64, cond condition.Condition) error
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

func (cli *object) FindSingleObject(params types.LogicParams, objectID string) (model.Object, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.Header.OwnerID)
	cond.Field(common.BKObjIDField).Eq(objectID)

	objs, err := cli.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the supplier account(%s) objects(%s), error info is %s", params.Header.OwnerID, objectID, err.Error())
		return nil, err
	}
	for _, item := range objs {
		return item, nil
	}
	return nil, fmt.Errorf("not found the object(%s)", objectID)
}
func (cli *object) CreateObject(params types.LogicParams, data frtypes.MapStr) (model.Object, error) {
	obj := cli.modelFactory.CreaetObject(params)

	_, err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	return obj, nil
}

func (cli *object) DeleteObject(params types.LogicParams, id int64, cond condition.Condition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().DeleteObject(context.Background(), id, params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[opration-obj] failed to delete the object by the condition(%#v) or the id(%d)", cond.ToMapStr(), id)
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (cli *object) FindObject(params types.LogicParams, cond condition.Condition) ([]model.Object, error) {

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return model.CreateObject(params, cli.clientSet, rsp.Data), nil
}

func (cli *object) UpdateObject(params types.LogicParams, data frtypes.MapStr, id int64, cond condition.Condition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().UpdateObject(context.Background(), id, params.Header.ToHeader(), data)

	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-obj] failed to set the object to the new data(%#v) by the condition(%#v) or the  id (%d)", data, cond.ToMapStr(), id)
		return params.Err.Error(rsp.Code)
	}

	return nil
}
