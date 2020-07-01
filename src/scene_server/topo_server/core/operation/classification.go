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

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// ClassificationOperationInterface classification operation methods
type ClassificationOperationInterface interface {
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, asst AssociationOperationInterface, obj ObjectOperationInterface)

	// FindSingleClassification(kit *rest.Kit, classificationID string) (model.Classification, error)
	CreateClassification(kit *rest.Kit, data mapstr.MapStr) (model.Classification, error)
	DeleteClassification(kit *rest.Kit, id int64, cond condition.Condition, metadata *metadata.Metadata) error
	FindClassification(kit *rest.Kit, cond condition.Condition, metadata *metadata.Metadata) ([]model.Classification, error)
	FindClassificationWithObjects(kit *rest.Kit, cond condition.Condition, metadata *metadata.Metadata) ([]metadata.ClassificationWithObject, error)
	UpdateClassification(kit *rest.Kit, data mapstr.MapStr, id int64, cond condition.Condition) error
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

func (c *classification) CreateClassification(kit *rest.Kit, data mapstr.MapStr) (model.Classification, error) {
	cls := c.modelFactory.CreateClassification(kit)
	_, err := cls.Parse(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to parse the kit, error info is %s, rid: %s", err.Error(), kit.Rid)
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
	// if err := c.authManager.AuthorizeResourceCreate(kit.Ctx, kit.Header, businessID, metaData.ModelClassification); err != nil {
	// 	blog.V(2).Infof("create classification %+v failed, authorization failed, err: %+v, rid: %s", class, err, kit.Rid)
	// 	return nil, err
	// }

	err = cls.Create()
	if nil != err {
		blog.Errorf("[operation-cls]failed to save the classification(%#v), error info is %s, rid: %s", cls, err.Error(), kit.Rid)
		return nil, err
	}

	class := cls.Classify()

	//package audit response
	err = NewObjectClsAudit(kit, c.clientSet, class.ID).buildSnapshotForCur().SaveAuditLog(metadata.AuditCreate)
	if err != nil {
		blog.Errorf("create object attribute %s success, but update to auditLog failed, err: %v, rid: %s", class.ClassificationName, err, kit.Rid)
		return nil, err
	}
	return cls, nil
}

func (c *classification) DeleteClassification(kit *rest.Kit, id int64, cond condition.Condition, metaData *metadata.Metadata) error {

	if 0 < id {
		if nil == cond {
			cond = condition.CreateCondition()
		}
		cond.Field(metadata.ClassificationFieldID).Eq(id)
	}

	if nil != metaData {
		cond.Field(metadata.BKMetadata).Eq(*metaData)
	}

	clsItems, err := c.FindClassification(kit, cond, metaData)
	if nil != err {
		return err
	}

	for _, cls := range clsItems {
		objs, err := cls.GetObjects()
		if nil != err {
			return err
		}

		if 0 != len(objs) {
			blog.Warnf("[operation-cls] the classification(%s) has some objects, forbidden to delete, rid: %s", cls.Classify().ClassificationID, kit.Rid)
			return kit.CCError.Error(common.CCErrTopoObjectClassificationHasObject)
		}

	}

	//get PreData
	objAudit := NewObjectClsAudit(kit, c.clientSet, id).buildSnapshotForPre()

	rsp, err := c.clientSet.CoreService().Model().DeleteModelClassification(context.Background(), kit.Header, &metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	if !rsp.Result {
		blog.Errorf("failed to delete the classification, error info is %s, rid: %s", rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	//saveAuditLog
	err = objAudit.SaveAuditLog(metadata.AuditDelete)
	if err != nil {
		blog.Errorf("Delete object attribute group success, but update to auditLog failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

func (c *classification) FindClassificationWithObjects(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]metadata.ClassificationWithObject, error) {
	fCond := cond.ToMapStr()
	if nil != metaData {
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := c.clientSet.CoreService().Model().ReadModelClassification(context.Background(), kit.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	clsIDs := make([]string, 0)
	for _, cls := range rsp.Data.Info {
		clsIDs = append(clsIDs, cls.ClassificationID)
	}
	clsIDs = util.StrArrayUnique(clsIDs)
	queryObjectCond := condition.CreateCondition().Field(common.BKClassificationIDField).In(clsIDs)
	queryObjectResp, err := c.clientSet.CoreService().Model().ReadModel(context.Background(), kit.Header, &metadata.QueryCondition{Condition: queryObjectCond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}
	if !queryObjectResp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", fCond, queryObjectResp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(queryObjectResp.Code, queryObjectResp.ErrMsg)
	}
	objMap := make(map[string][]metadata.Object)
	objIDs := make([]string, 0)
	for _, info := range queryObjectResp.Data.Info {
		objIDs = append(objIDs, info.Spec.ObjectID)
		objMap[info.Spec.ObjCls] = append(objMap[info.Spec.ObjCls], info.Spec)
	}

	datas := make([]metadata.ClassificationWithObject, 0)
	for _, cls := range rsp.Data.Info {
		clsItem := metadata.ClassificationWithObject{
			Classification: cls,
			Objects:        []metadata.Object{},
		}
		if obj, ok := objMap[cls.ClassificationID]; ok {
			clsItem.Objects = obj
		}
		datas = append(datas, clsItem)
	}

	return datas, nil
}

func (c *classification) FindClassification(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.Classification, error) {
	fCond := cond.ToMapStr()
	if nil != metaData {
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		// fCond.Merge(metadata.BizLabelNotExist)
	}

	rsp, err := c.clientSet.CoreService().Model().ReadModelClassification(context.Background(), kit.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-cls]failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	if !rsp.Result {
		blog.Errorf("[operation-cls] failed to search the classification by the condition(%#v), error info is %s, rid: %s", cond.ToMapStr(), rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	clsItems := model.CreateClassification(kit, c.clientSet, rsp.Data.Info, metaData)
	return clsItems, nil
}

func (c *classification) UpdateClassification(kit *rest.Kit, data mapstr.MapStr, id int64, cond condition.Condition) error {
	cls := c.modelFactory.CreateClassification(kit)
	data.Set("id", id)
	if _, err := cls.Parse(data); err != nil {
		blog.Errorf("update classification, but parse classification failed, err：%v, rid: %s", err, kit.Rid)
		return err
	}

	class := cls.Classify()
	class.ID = id

	// auth: check authorization
	// if err := c.authManager.AuthorizeByClassification(kit.Ctx, kit.Header, metaData.Update, class); err != nil {
	// 	blog.V(2).Infof("update classification %s failed, authorization failed, err: %+v, rid: %s", class, err, kit.Rid)
	// 	return err
	// }

	//get PreData
	objAudit := NewObjectClsAudit(kit, c.clientSet, id).buildSnapshotForPre()

	err := cls.Update(data)
	if nil != err {
		blog.Errorf("[operation-cls]failed to update the classification(%#v), error info is %s, rid: %s", cls, err.Error(), kit.Rid)
		return err
	}

	//get CurData and saveAuditLog
	err = objAudit.buildSnapshotForCur().SaveAuditLog(metadata.AuditUpdate)
	if err != nil {
		blog.Errorf("update object attribute-id %s success, but update to auditLog failed, err: %v, rid: %s", id, err, kit.Rid)
		return err
	}
	return nil
}
