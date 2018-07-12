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
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	metatype "configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// InstOperationInterface inst operation methods
type InstOperationInterface interface {
	CreateInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error)
	DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition) error
	FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput) (count int, results []inst.Inst, err error)
	UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface)
}

// NewInstOperation create a new inst operation instance
func NewInstOperation(client apimachinery.ClientSetInterface) InstOperationInterface {
	return &commonInst{
		clientSet: client,
	}
}

type commonInst struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	asst         AssociationOperationInterface
}

func (c *commonInst) SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface) {
	c.modelFactory = modelFactory
	c.instFactory = instFactory
	c.asst = asst
}

func (c *commonInst) CreateInst(params types.ContextParams, obj model.Object, data frtypes.MapStr) (inst.Inst, error) {

	// extract internal data
	batchInfo := &instBatchInfo{}
	if err := data.MarshalJSONInto(batchInfo); nil != err {
		blog.Errorf("[operation-inst] failed to unmarshal the data(%#v) into the inst batch info struct, error info is %s", data, err.Error())
		return nil, params.Err.Error(common.CCErrCommJSONUnmarshalFailed)
	}

	// create association

	// create new insts
	blog.Infof("the data inst:%#v", data)
	item := c.instFactory.CreateInst(params, obj)

	item.SetValues(data)

	err := item.Create()
	if nil != err {
		blog.Errorf("[operation-inst] failed to save the object(%s) inst data (%#v), error info is %s", obj.GetID(), data, err.Error())
		return nil, err
	}

	return item, nil
}

func (c *commonInst) DeleteInst(params types.ContextParams, obj model.Object, cond condition.Condition) error {

	rsp, err := c.clientSet.ObjectController().Instance().DelObject(context.Background(), obj.GetObjectType(), params.Header, cond.ToMapStr())

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (c *commonInst) FindInst(params types.ContextParams, obj model.Object, cond *metatype.QueryInput) (count int, results []inst.Inst, err error) {

	rsp, err := c.clientSet.ObjectController().Instance().SearchObjects(context.Background(), obj.GetObjectType(), params.Header, cond)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return 0, nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to delete the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond, rsp.ErrMsg)
		return 0, nil, params.Err.Error(rsp.Code)
	}

	return rsp.Data.Count, inst.CreateInst(params, c.clientSet, obj, rsp.Data.Info), nil
}

func (c *commonInst) UpdateInst(params types.ContextParams, data frtypes.MapStr, obj model.Object, cond condition.Condition) error {

	inputParams := frtypes.New()
	inputParams.Set("data", data)
	inputParams.Set("condition", cond.ToMapStr())
	blog.Infof("data condition:%#v", inputParams)
	rsp, err := c.clientSet.ObjectController().Instance().UpdateObject(context.Background(), obj.GetObjectType(), params.Header, inputParams)

	if nil != err {
		blog.Errorf("[operation-inst] failed to request object controller, error info is %s", err.Error())
		return params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-inst] faild to set the object(%s) inst by the condition(%#v), error info is %s", obj.GetID(), cond.ToMapStr(), rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}
	return nil
}
