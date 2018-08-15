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

	"configcenter/src/common"
	"configcenter/src/common/metadata"

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
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)

	FindSingleClassification(params types.ContextParams, classificationID string) (model.Classification, error)
	CreateClassification(params types.ContextParams, data frtypes.MapStr) (model.Classification, error)
	DeleteClassification(params types.ContextParams, id int64, data frtypes.MapStr, cond condition.Condition) error
	FindClassification(params types.ContextParams, cond condition.Condition) ([]model.Classification, error)
	FindClassificationWithObjects(params types.ContextParams, cond condition.Condition) ([]metadata.ClassificationWithObject, error)
	UpdateClassification(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error
}

// NewClassificationOperation create a new classification operation instance
func NewClassificationOperation(client apimachinery.ClientSetInterface) ClassificationOperationInterface {
	return &classification{
		clientSet: client,
	}
}

type classification struct {
	clientSet    apimachinery.ClientSetInterface
	asst         AssociationOperationInterface
	obj          ObjectOperationInterface
	modelFactory model.Factory
	instFactory  inst.Factory
}

func (c *classification) SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface) {
	c.modelFactory = modelFactory
	c.instFactory = instFactory
	c.asst = asst
	c.obj = obj
}

func (c *classification) FindSingleClassification(params types.ContextParams, classificationID string) (model.Classification, error) {

	cond := condition.CreateCondition()
	cond.Field(metadata.ClassFieldClassificationID).Eq(classificationID)

	objs, err := c.FindClassification(params, cond)
	if nil != err {
		blog.Errorf("[operation-cls] failed to find the supplier account(%s) classification(%s), error info is %s", params.SupplierAccount, classificationID, err.Error())
		return nil, err
	}
	for _, item := range objs {
		return item, nil
	}
	return nil, params.Err.Error(common.CCErrTopoObjectClassificationSelectFailed)
}

func (c *classification) CreateClassification(params types.ContextParams, data frtypes.MapStr) (model.Classification, error) {

	cls := c.modelFactory.CreaetClassification(params)

	_, err := cls.Parse(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to parse the params, error info is %s", err.Error())
		return nil, err
	}

	err = cls.Create()
	if nil != err {
		blog.Errorf("[operation-cls]failed to save the classification(%#v), error info is %s", cls, err.Error())
		return nil, err
	}

	return cls, nil
}

func (c *classification) DeleteClassification(params types.ContextParams, id int64, data frtypes.MapStr, cond condition.Condition) error {

	if 0 < id {
		if nil == cond {
			cond = condition.CreateCondition()
		}
		cond.Field(metadata.ClassificationFieldID).Eq(id)
	}

	clsItems, err := c.FindClassification(params, cond)
	if nil != err {
		return err
	}

	for _, cls := range clsItems {
		objs, err := cls.GetObjects()
		if nil != err {
			return err
		}

		if 0 != len(objs) {
			blog.Errorf("[operation-cls] the classification(%s) has some obejcts, forbidden to delete", cls.GetID())
			return params.Err.Error(common.CCErrTopoObjectClassificationHasObject)
		}
	}

	rsp, err := c.clientSet.ObjectController().Meta().DeleteClassification(context.Background(), id, params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("failed to delete the classification, error info is %s", rsp.ErrMsg)
		return params.Err.Error(rsp.Code)
	}

	return nil
}

func (c *classification) FindClassificationWithObjects(params types.ContextParams, cond condition.Condition) ([]metadata.ClassificationWithObject, error) {

	rsp, err := c.clientSet.ObjectController().Meta().SelectClassificationWithObject(context.Background(), params.SupplierAccount, params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-cls] failed to search the clssificaiton by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}
	//fmt.Println("rsp.data:", rsp.Data)
	for idx, clsItem := range rsp.Data {
		rsp.Data[idx].AsstObjects = make(map[string][]metadata.Object)
		for _, objItem := range clsItem.Objects {
			asstItems, err := c.asst.SearchObjectAssociation(params, objItem.ObjectID)
			if nil != err {
				return nil, params.Err.New(common.CCErrTopoObjectClassificationSelectFailed, err.Error())
			}

			for _, asstItem := range asstItems {
				if common.BKChildStr == asstItem.ObjectAttID {
					continue
				}

				searchObjCond := condition.CreateCondition()
				searchObjCond.Field(common.BKOwnerIDField).Eq(params.SupplierAccount)
				searchObjCond.Field(common.BKObjIDField).Eq(asstItem.AsstObjID)
				asstObjs, err := c.obj.FindObject(params, searchObjCond)
				if nil != err {
					return nil, err
				}

				for _, obj := range asstObjs {
					rsp.Data[idx].AsstObjects[objItem.ObjectID] = append(rsp.Data[idx].AsstObjects[objItem.ObjectID], obj.Origin())
				}

			}
		}

	}

	return rsp.Data, nil
}

func (c *classification) FindClassification(params types.ContextParams, cond condition.Condition) ([]model.Classification, error) {

	rsp, err := c.clientSet.ObjectController().Meta().SelectClassifications(context.Background(), params.Header, cond.ToMapStr())
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if common.CCSuccess != rsp.Code {
		blog.Errorf("[operation-cls] failed to search the clssificaiton by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.Error(rsp.Code)
	}

	clsItems := model.CreateClassification(params, c.clientSet, rsp.Data)
	return clsItems, nil
}

func (c *classification) UpdateClassification(params types.ContextParams, data frtypes.MapStr, id int64, cond condition.Condition) error {

	cls := c.modelFactory.CreaetClassification(params)
	data.Set("id", id)
	cls.Parse(data)

	err := cls.Update(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to update the classification(%#v), error info is %s", cls, err.Error())
		return err
	}

	return nil
}
