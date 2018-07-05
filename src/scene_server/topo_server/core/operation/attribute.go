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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
	"context"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(params types.ContextParams, data frtypes.MapStr) (model.Attribute, error)
	DeleteObjectAttribute(params types.ContextParams, id int64, cond condition.Condition) error
	FindObjectAttribute(params types.ContextParams, cond condition.Condition) ([]model.Attribute, error)
	UpdateObjectAttribute(params types.ContextParams, data frtypes.MapStr, attID int64, cond condition.Condition) error
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

func (cli *attribute) CreateObjectAttribute(params types.ContextParams, data frtypes.MapStr) (model.Attribute, error) {
	att := cli.modelFactory.CreateAttribute(params)

	_, err := att.Parse(data)
	if nil != err {
		blog.Errorf("[operation-attr] failed to parse the attribute data (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	err = att.Create()
	if nil != err {
		blog.Errorf("[operation-attr] failed to save the attribute data (%#v), error info is %s", data, err.Error())
		return nil, err
	}

	return att, nil
}

func (cli *attribute) DeleteObjectAttribute(params types.ContextParams, id int64, cond condition.Condition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().DeleteObjectAttByID(context.Background(), id, params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-attr] failed to delete the attribute by the id(%d) or the condition(%#v), error info is %s", id, cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (cli *attribute) FindObjectAttribute(params types.ContextParams, cond condition.Condition) ([]model.Attribute, error) {

	rsp, err := cli.clientSet.ObjectController().Meta().SelectObjectAttWithParams(context.Background(), params.Header.ToHeader(), cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-attr] failed to search attribute by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return model.CreateAttribute(params, cli.clientSet, rsp.Data), nil
}

func (cli *attribute) UpdateObjectAttribute(params types.ContextParams, data frtypes.MapStr, attID int64, cond condition.Condition) error {

	rsp, err := cli.clientSet.ObjectController().Meta().UpdateObjectAttByID(context.Background(), attID, params.Header.ToHeader(), data)

	if nil != err {
		blog.Errorf("[operation-attr] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-attr] failed to update the attribute by the condition(%#v) or the attr-id(%d), error info is %s", cond.ToMapStr(), attID, rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}
