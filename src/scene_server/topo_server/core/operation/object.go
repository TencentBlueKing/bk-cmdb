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
	"fmt"

	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	frtypes "configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObject(params types.ContextParams, data frtypes.MapStr) (model.Object, error)
	DeleteObject(params types.ContextParams, id int64, cond condition.Condition) error
	FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error)
	FindSingleObject(params types.ContextParams, objectID string) (model.Object, error)
	UpdateObject(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface)
	IsValidObject(params types.ContextParams, objID string) error
}

// NewObjectOperation create a new object operation instance
func NewObjectOperation(client apimachinery.ClientSetInterface) ObjectOperationInterface {
	return &object{
		clientSet: client,
	}
}

type object struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	cls          ClassificationOperationInterface
}

func (o *object) SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface) {
	o.modelFactory = modelFactory
	o.instFactory = instFactory
}

func (o *object) IsValidObject(params types.ContextParams, objID string) error {

	checkObjCond := condition.CreateCondition()
	checkObjCond.Field(metadata.AttributeFieldObjectID).Eq(objID)
	checkObjCond.Field(metadata.AttributeFieldSupplierAccount).Eq(params.SupplierAccount)

	objItems, err := o.FindObject(params, checkObjCond)
	if nil != err {
		blog.Errorf("[opeartion-attr] failed to check the object repeated, error info is %s", err.Error())
		return params.Err.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if 0 == len(objItems) {
		return params.Err.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("the object id  '%s' is invalid", objID))
	}

	return nil
}
func (o *object) FindSingleObject(params types.ContextParams, objectID string) (model.Object, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
	cond.Field(common.BKObjIDField).Eq(objectID)

	objs, err := o.FindObject(params, cond)
	if nil != err {
		blog.Errorf("[api-inst] failed to find the supplier account(%s) objects(%s), error info is %s", params.SupplierAccount, objectID, err.Error())
		return nil, err
	}
	for _, item := range objs {
		return item, nil
	}
	return nil, params.Err.Error(common.CCErrTopoObjectSelectFailed)
}
func (o *object) CreateObject(params types.ContextParams, data frtypes.MapStr) (model.Object, error) {
	obj := o.modelFactory.CreaetObject(params)

	_, err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	// check the classification
	_, err = obj.GetClassification()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object, error info is %s", err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	// check repeated

	exists, err := obj.IsExists()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object(%#v), error info is %s", data, err.Error())
		return nil, params.Err.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	if exists {
		blog.Errorf("[operation-obj] the object(%#v) is repeated", data)
		return nil, params.Err.Error(common.CCErrCommDuplicateItem)
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), error info is %s", data, err.Error())
		return nil, err
	}

	// create the default group
	grp := obj.CreateGroup()
	grp.SetDefault(true)
	grp.SetIndex(-1)
	grp.SetName("Default")
	grp.SetID("default")

	if err = grp.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, error info is %s", err.Error())
	}

	// create the default inst name
	attr := obj.CreateAttribute()
	attr.SetIsOnly(true)
	attr.SetIsPre(true)
	attr.SetCreator("user")
	attr.SetIsEditable(true)
	attr.SetGroupIndex(-1)
	attr.SetGroup(grp)
	attr.SetIsRequired(true)
	attr.SetType(common.FieldTypeSingleChar)
	attr.SetID(obj.GetInstNameFieldName())
	attr.SetName(obj.GetDefaultInstPropertyName())

	if err = attr.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, error info is %s", err.Error())
	}

	return obj, nil
}

func (o *object) DeleteObject(params types.ContextParams, id int64, cond condition.Condition) error {

	rsp, err := o.clientSet.ObjectController().Meta().DeleteObject(context.Background(), id, params.Header, cond.ToMapStr())

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

func (o *object) FindObject(params types.ContextParams, cond condition.Condition) ([]model.Object, error) {

	rsp, err := o.clientSet.ObjectController().Meta().SelectObjects(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, error info is %s", err.Error())
		return nil, params.Err.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	return model.CreateObject(params, o.clientSet, rsp.Data), nil
}

func (o *object) UpdateObject(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error {

	obj := o.modelFactory.CreaetObject(params)
	obj.SetRecordID(id)
	_, err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), error info is %s", data, err.Error())
		return err
	}

	if err = obj.Update(data); nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%d), the new data(%#v), error info is %s", id, data, err.Error())
		return params.Err.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	return nil
}
