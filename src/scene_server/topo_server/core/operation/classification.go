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

	"configcenter/src/apimachinery"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ClassificationOperationInterface classification opoeration methods
type ClassificationOperationInterface interface {
	CreateClassification(params types.LogicParams, data frtypes.MapStr) (model.Classification, error)
	DeleteClassification(params types.LogicParams, id int64, cond condition.Condition) error
	FindClassification(params types.LogicParams, cond condition.Condition) ([]model.Classification, error)
	UpdateClassification(params types.LogicParams, data frtypes.MapStr, id int64, cond condition.Condition) error
}

type classification struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

// NewClassificationOperation create a new classification operation instance
func NewClassificationOperation(client apimachinery.ClientSetInterface, modelFactory model.Factory, instFactory inst.Factory) ClassificationOperationInterface {
	return &classification{
		clientSet:    client,
		modelFactory: modelFactory,
		instFactory:  instFactory,
	}
}

func (cli *classification) CreateClassification(params types.LogicParams, data frtypes.MapStr) (model.Classification, error) {

	cls := cli.modelFactory.CreaetClassification(params)

	_, err := cls.Parse(data)
	if nil != err {
		return nil, err
	}

	err = cls.Save()
	if nil != err {
		return nil, err
	}

	return cls, nil
}

func (cli *classification) DeleteClassification(params types.LogicParams, id int64, cond condition.Condition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().DeleteClassification(context.Background(), id, params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-cls] failed to delete the classification by id (%d) or the condition (%#v), error info is %s ", id, cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (cli *classification) FindClassification(params types.LogicParams, cond condition.Condition) ([]model.Classification, error) {

	rsp, err := cli.clientSet.ObjectController().Meta().SelectClassifications(context.Background(), params.Header.ToHeader(), cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-cls] failed to search the clssificaiton by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	clsItems := model.CreateClassification(params, cli.clientSet, rsp.Data)
	return clsItems, nil
}

func (cli *classification) UpdateClassification(params types.LogicParams, data frtypes.MapStr, id int64, cond condition.Condition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().UpdateClassification(context.Background(), id, params.Header.ToHeader(), data)

	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-cls] failed to update the classification by the id (%d) or the condtion (%#v), error info is %s", id, cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}
