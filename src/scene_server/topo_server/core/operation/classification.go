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
	"configcenter/src/common/util"
	"context"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
	"configcenter/src/scene_server/topo_server/core/types"
)

// ClassificationOperationInterface classification operation methods
type ClassificationOperationInterface interface {
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)

	// FindSingleClassification(params types.ContextParams, classificationID string) (model.Classification, error)
	CreateClassification(params types.ContextParams, data mapstr.MapStr) (model.Classification, error)
	DeleteClassification(params types.ContextParams, id int64, data mapstr.MapStr, cond condition.Condition) error
	FindClassification(params types.ContextParams, cond condition.Condition) ([]model.Classification, error)
	FindClassificationWithObjects(params types.ContextParams, cond condition.Condition) ([]metadata.ClassificationWithObject, error)
	UpdateClassification(params types.ContextParams, data mapstr.MapStr, id int64, cond condition.Condition) error
}

// NewClassificationOperation create a new classification operation instance
func NewClassificationOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) ClassificationOperationInterface {
	return &classification{
		clientSet:   client,
		authManager: authManager,
	}
}

type classification struct {
	clientSet    apimachinery.ClientSetInterface
	authManager  *extensions.AuthManager
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

func (c *classification) CreateClassification(params types.ContextParams, data mapstr.MapStr) (model.Classification, error) {
	cls := c.modelFactory.CreateClassification(params)
	_, err := cls.Parse(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to parse the params, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	// auth: check authorization
	// class := cls.Classify()
	// var businessID int64
	// if _, exist := class.Metadata.Label[metadata.LabelBusinessID]; exist {
	// 	var err error
	// 	businessID, err = class.Metadata.Label.Int64(metadata.LabelBusinessID)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }
	//
	// if err := c.authManager.AuthorizeResourceCreate(params.Context, params.Header, businessID, meta.ModelClassification); err != nil {
	// 	blog.V(2).Infof("create classification %+v failed, authorization failed, err: %+v, rid: %s", class, err, params.ReqID)
	// 	return nil, err
	// }

	err = cls.Create()
	if nil != err {
		blog.Errorf("[operation-cls]failed to save the classification(%#v), error info is %s, rid: %s", cls, err.Error(), params.ReqID)
		return nil, err
	}

	// auth: register new created classify
	class := cls.Classify()
	if err := c.authManager.RegisterClassification(params.Context, params.Header, class); err != nil {
		return nil, params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
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
			blog.Warnf("[operation-cls] the classification(%s) has some objects, forbidden to delete, rid: %s", cls.Classify().ClassificationID, params.ReqID)
			return params.Err.Error(common.CCErrTopoObjectClassificationHasObject)
		}

		if err := c.authManager.DeregisterClassification(params.Context, params.Header, cls.Classify()); err != nil {
			return params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
		}

	}

	rsp, err := c.clientSet.CoreService().Model().DeleteModelClassification(context.Background(), params.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to delete the classification, error info is %s, rid: %s", rsp.ErrMsg, params.ReqID)
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
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	clsIDs := make([]string, 0)
	for _, cls := range rsp.Data.Info {
		clsIDs = append(clsIDs, cls.ClassificationID)
	}
	clsIDs = util.StrArrayUnique(clsIDs)
	queryObjectCond := condition.CreateCondition().Field(common.BKClassificationIDField).In(clsIDs)
	queryObjectResp, err := c.clientSet.CoreService().Model().ReadModel(context.Background(), params.Header, &metadata.QueryCondition{Condition: queryObjectCond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}
	if !queryObjectResp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", fCond, queryObjectResp.ErrMsg, params.ReqID)
		return nil, params.Err.New(queryObjectResp.Code, queryObjectResp.ErrMsg)
	}
	objMap := make(map[string][]metadata.Object)
	objIDs := make([]string, 0)
	for _, info := range queryObjectResp.Data.Info {
		objIDs = append(objIDs, info.Spec.ObjectID)
		objMap[info.Spec.ObjCls] = append(objMap[info.Spec.ObjCls], info.Spec)
	}
	objIDs = util.StrArrayUnique(objIDs)
	asstItems, err := c.asst.SearchObjectsAssociation(params, objIDs)
	if nil != err {
		return nil, params.Err.New(common.CCErrTopoObjectClassificationSelectFailed, err.Error())
	}
	asstIDs := make([]string, 0)
	for _, asstItem := range asstItems {
		asstIDs = append(asstIDs, asstItem.AsstObjID)
	}
	asstIDs = util.StrArrayUnique(asstIDs)
	searchObjCond := condition.CreateCondition()
	searchObjCond.Field(common.BKObjIDField).In(asstIDs)
	asstObjs, err := c.obj.FindObject(params, searchObjCond)
	if nil != err {
		return nil, err
	}

	asstObjsMap := make(map[string]map[string][]metadata.Object)
	for _, asstItem := range asstItems {
		asstObjMap := make(map[string][]metadata.Object)
		if asstObjs, ok := asstObjsMap[asstItem.ObjectID]; ok {
			asstObjMap = asstObjs
		}
		for _, obj := range asstObjs {
			if obj.Object().ObjectID == asstItem.AsstObjID {
				asstObjMap[asstItem.ObjectID] = append(asstObjMap[asstItem.ObjectID], obj.Object())
			}
		}
		asstObjsMap[asstItem.ObjectID] = asstObjMap
	}

	asstMap := make(map[string]map[string][]metadata.Object)
	for _, info := range queryObjectResp.Data.Info {
		asstObjMap := make(map[string][]metadata.Object)
		if asstObjs, ok := asstMap[info.Spec.ObjCls]; ok {
			asstObjMap = asstObjs
		}
		if asstObjs, ok := asstObjsMap[info.Spec.ObjectID]; ok {
			for asstObjKey, asstObj := range asstObjs {
				asstObjMap[asstObjKey] = asstObj
			}
		}
		asstMap[info.Spec.ObjCls] = asstObjMap
	}

	datas := make([]metadata.ClassificationWithObject, 0)
	for _, cls := range rsp.Data.Info {
		clsItem := metadata.ClassificationWithObject{
			Classification: cls,
			Objects:        []metadata.Object{},
			AsstObjects:    map[string][]metadata.Object{},
		}
		if obj, ok := objMap[cls.ClassificationID]; ok {
			clsItem.Objects = obj
		}
		if asst, ok := asstMap[cls.ClassificationID]; ok {
			clsItem.AsstObjects = asst
		}
		datas = append(datas, clsItem)
	}

	return datas, nil
}

func (c *classification) FindClassification(params types.ContextParams, cond condition.Condition) ([]model.Classification, error) {
	fCond := cond.ToMapStr()
	if nil != params.MetaData {
		fCond.Merge(metadata.PublicAndBizCondition(*params.MetaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		// fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := c.clientSet.CoreService().Model().ReadModelClassification(context.Background(), params.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), params.ReqID)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, params.ReqID)
		return nil, params.Err.New(rsp.Code, rsp.ErrMsg)
	}

	clsItems := model.CreateClassification(params, c.clientSet, rsp.Data.Info)
	return clsItems, nil
}

func (c *classification) UpdateClassification(params types.ContextParams, data mapstr.MapStr, id int64, cond condition.Condition) error {
	cls := c.modelFactory.CreateClassification(params)
	data.Set("id", id)
	if _, err := cls.Parse(data); err != nil {
		blog.Errorf("update classification, but parse classification failed, err：%v, rid: %s", err, params.ReqID)
		return err
	}

	class := cls.Classify()
	class.ID = id

	// auth: check authorization
	// if err := c.authManager.AuthorizeByClassification(params.Context, params.Header, meta.Update, class); err != nil {
	// 	blog.V(2).Infof("update classification %s failed, authorization failed, err: %+v, rid: %s", class, err, params.ReqID)
	// 	return err
	// }

	err := cls.Update(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to update the classification(%#v), error info is %s, rid: %s", cls, err.Error(), params.ReqID)
		return err
	}

	// auth: update registered classifications
	if len(class.ClassificationID) > 0 {
		if err := c.authManager.UpdateRegisteredClassificationByID(params.Context, params.Header, class.ClassificationID); err != nil {
			blog.Errorf("update classification %s, but update to auth failed, err: %v, rid: %s", class.ClassificationName, err, params.ReqID)
			return params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
		}
	} else {
		if err := c.authManager.UpdateRegisteredClassificationByRawID(params.Context, params.Header, class.ID); err != nil {
			blog.Errorf("update classification %s, but update to auth failed, err: %v, rid: %s", class.ClassificationName, err, params.ReqID)
			return params.Err.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
		}
	}

	return nil
}
