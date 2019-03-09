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
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ClassificationOperationInterface classification opoeration methods
type ClassificationOperationInterface interface {
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)

	FindSingleClassification(params types.ContextParams, classificationID string) (model.Classification, error)
	CreateClassification(params types.ContextParams, data mapstr.MapStr) (model.Classification, error)
	DeleteClassification(params types.ContextParams, id int64, data mapstr.MapStr, cond condition.Condition) error
	FindClassification(params types.ContextParams, cond condition.Condition) ([]model.Classification, error)
	FindClassificationWithObjects(params types.ContextParams, cond condition.Condition) ([]metadata.ClassificationWithObject, error)
	UpdateClassification(params types.ContextParams, data mapstr.MapStr, id int64, cond condition.Condition) error
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

func (c *classification) CreateClassification(params types.ContextParams, data mapstr.MapStr) (model.Classification, error) {

	cls := c.modelFactory.CreateClassification(params)

	_, err := cls.Parse(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to parse the params, error info is %s", err.Error())
		return nil, err
	}

	//	if nil != params.MetaData {
	//		data.Set(metadata.BKMetadata, *params.MetaData)
	//	}
	err = cls.Create()
	if nil != err {
		blog.Errorf("[operation-cls]failed to save the classification(%#v), error info is %s", cls, err.Error())
		return nil, err
	}

	return cls, nil
}

func (c *classification) DeleteClassification(params types.ContextParams, id int64, data mapstr.MapStr, cond condition.Condition) error {

	if 0 < id {
		if nil == cond {
			cond = condition.CreateCondition()
		}
		cond.Field(metadata.ClassificationFieldID).Eq(id)
	}

	if nil != params.MetaData {
		cond.Field(metadata.BKMetadata).Eq(*params.MetaData)
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
			blog.Errorf("[operation-cls] the classification(%s) has some obejcts, forbidden to delete", cls.Classify().ClassificationID)
			return params.Err.Error(common.CCErrTopoObjectClassificationHasObject)
		}
	}

	rsp, err := c.clientSet.CoreService().Model().DeleteModelClassification(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to delete the classification, error info is %s", rsp.ErrMsg)
		return params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (c *classification) FindClassificationWithObjects(params types.ContextParams, cond condition.Condition) ([]metadata.ClassificationWithObject, error) {
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := c.clientSet.CoreService().Model().ReadModelClassification(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s", fCond, rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	datas := []metadata.ClassificationWithObject{}
	for _, cls := range rsp.Data.Info {
		clsItem := metadata.ClassificationWithObject{
			Classification: cls,
			Objects:        []metadata.Object{},
			AsstObjects:    map[string][]metadata.Object{},
		}
		queryObjectCond := condition.CreateCondition().Field(common.BKClassificationIDField).Eq(cls.ClassificationID)
		queryObjectResp, err := c.clientSet.CoreService().Model().ReadModel(context.Background(), params.Header, &metadata.QueryCondition{Condition: queryObjectCond.ToMapStr()})
		if nil != err {
			blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
			return nil, err
		}

		if !queryObjectResp.Result {
			blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s", fCond, queryObjectResp.ErrMsg)
			return nil, params.Err.New(queryObjectResp.Code, queryObjectResp.ErrMsg)
		}

		for _, info := range queryObjectResp.Data.Info {
			clsItem.Objects = append(clsItem.Objects, info.Spec)
		}

		datas = append(datas, clsItem)
	}

	for idx, clsItem := range datas {
		for _, objItem := range clsItem.Objects {
			asstItems, err := c.asst.SearchObjectAssociation(params, objItem.ObjectID)
			if nil != err {
				return nil, params.Err.New(common.CCErrTopoObjectClassificationSelectFailed, err.Error())
			}

			for _, asstItem := range asstItems {

				searchObjCond := condition.CreateCondition()
				searchObjCond.Field(common.BKObjIDField).Eq(asstItem.AsstObjID)
				asstObjs, err := c.obj.FindObject(params, searchObjCond)
				if nil != err {
					return nil, err
				}

				for _, obj := range asstObjs {
					datas[idx].AsstObjects[objItem.ObjectID] = append(datas[idx].AsstObjects[objItem.ObjectID], obj.Object())
				}

			}
		}

	}

	return datas, nil
}

func (c *classification) FindClassification(params types.ContextParams, cond condition.Condition) ([]model.Classification, error) {
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}
	rsp, err := c.clientSet.CoreService().Model().ReadModelClassification(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s", err.Error())
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[operation-cls] failed to search the clssificaiton by the condition(%#v), error info is %s", cond.ToMapStr(), rsp.ErrMsg)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	clsItems := model.CreateClassification(params, c.clientSet, rsp.Data.Info)
	return clsItems, nil
}

func (c *classification) UpdateClassification(params types.ContextParams, data mapstr.MapStr, id int64, cond condition.Condition) error {

	cls := c.modelFactory.CreateClassification(params)
	data.Set("id", id)
	if _, err := cls.Parse(data); err != nil {
		blog.Errorf("update classification, but parse classification failed, err：%v", err)
		return err
	}

	err := cls.Update(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to update the classification(%#v), error info is %s", cls, err.Error())
		return err
	}

	return nil
}
