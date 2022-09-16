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

package model

import (
	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// ClassificationOperationInterface classification operation methods
type ClassificationOperationInterface interface {
	CreateClassification(kit *rest.Kit, data mapstr.MapStr) (*metadata.Classification, error)
	DeleteClassification(kit *rest.Kit, id int64) error
	FindClassification(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Classification, error)
	FindClassificationWithObjects(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.ClassificationWithObject, error)
	UpdateClassification(kit *rest.Kit, data mapstr.MapStr, id int64) error
}

// NewClassificationOperation create a new classification operation instance
func NewClassificationOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) ClassificationOperationInterface {
	return &classification{
		clientSet:   client,
		authManager: authManager,
	}
}

type classification struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// CreateClassification create classification
func (c *classification) CreateClassification(kit *rest.Kit, data mapstr.MapStr) (*metadata.Classification, error) {
	cls, err := c.isValid(kit, true, data)
	if err != nil {
		blog.ErrorJSON("validate classification update data failed, err: %s, data: %s, rid: %s",
			err, data, kit.Rid)
		return nil, err
	}

	// create classification
	input := metadata.CreateOneModelClassification{Data: *cls}
	rsp, err := c.clientSet.CoreService().Model().CreateModelClassification(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("create classification %s failed, err: %v, rid: %s", cls.ClassificationName, err, kit.Rid)
		return nil, err
	}

	// get created model classification data by id
	clsReq := &metadata.QueryCondition{Condition: mapstr.MapStr{metadata.ClassificationFieldID: int64(rsp.Created.ID)}}
	clsResp, err := c.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, clsReq)
	if err != nil {
		blog.Errorf("get created model classification by id(%d) failed, err: %v, rid: %s", rsp.Created.ID, err, kit.Rid)
		return nil, err
	}

	if len(clsResp.Info) != 1 {
		blog.Errorf("get created model classification by id(%d) returns not one cls, rid: %s", rsp.Created.ID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	cls = &clsResp.Info[0]

	// generate audit log of object classification.
	audit := auditlog.NewObjectClsAuditLog(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, cls.ID, cls)
	if err != nil {
		blog.Errorf("create object classification %s success, but generate audit log failed, err: %v, rid: %s",
			cls.ClassificationName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object classification %s success, but save audit log failed, err: %v, rid: %s",
			cls.ClassificationName, err, kit.Rid)
		return nil, err
	}

	return cls, nil
}

// DeleteClassification delete classification
func (c *classification) DeleteClassification(kit *rest.Kit, id int64) error {
	cond := map[string]interface{}{metadata.ClassificationFieldID: id}
	clsItems, err := c.FindClassification(kit, cond)
	if err != nil {
		return err
	}

	// check if the classification has objects in it, if so, can not delete it
	objects, err := c.getClassificationObjects(kit, clsItems)
	if err != nil {
		return err
	}

	if len(objects) != 0 {
		blog.Errorf("the classification(%d) has some objects, forbidden to delete, rid: %s", id, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectClassificationHasObject)
	}

	// generate audit log of object classification.
	audit := auditlog.NewObjectClsAuditLog(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, id, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before delete object classification, objClsID: %d, err: %v, rid: %s",
			id, err, kit.Rid)
		return err
	}

	// delete classification
	delOpt := &metadata.DeleteOption{Condition: cond}
	_, err = c.clientSet.CoreService().Model().DeleteModelClassification(kit.Ctx, kit.Header, delOpt)
	if err != nil {
		blog.Errorf("delete classification %d failed, err: %v, rid: %s", id, err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("delete object classification success, but save audit log failed, objClsID: %d, err: %v, rid: %s",
			id, err, kit.Rid)
		return err
	}

	return nil
}

// FindClassificationWithObjects TODO
func (c *classification) FindClassificationWithObjects(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.ClassificationWithObject, error) {
	classifications, err := c.FindClassification(kit, cond)
	if err != nil {
		return nil, err
	}

	objects, err := c.getClassificationObjects(kit, classifications)
	if err != nil {
		return nil, err
	}

	// rearrange data
	objMap := make(map[string][]metadata.Object)
	for _, object := range objects {
		objMap[object.ObjCls] = append(objMap[object.ObjCls], object)
	}

	datas := make([]metadata.ClassificationWithObject, 0)
	for _, cls := range classifications {
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

// getClassificationObjects get objects that are in the classifications
func (c *classification) getClassificationObjects(kit *rest.Kit, classifications []metadata.Classification) (
	[]metadata.Object, error) {

	clsIDs := make([]string, 0)
	for _, cls := range classifications {
		clsIDs = append(clsIDs, cls.ClassificationID)
	}
	clsIDs = util.StrArrayUnique(clsIDs)

	queryObjectCond := &metadata.QueryCondition{
		Condition: map[string]interface{}{
			common.BKClassificationIDField: map[string]interface{}{common.BKDBIN: clsIDs},
		},
	}
	rsp, err := c.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryObjectCond)
	if err != nil {
		blog.Errorf("search objects failed, err: %v, clsIDs: %+v, rid: %s", err, clsIDs, kit.Rid)
		return nil, err
	}

	return rsp.Info, nil
}

// FindClassification search classification
func (c *classification) FindClassification(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Classification, error) {
	input := &metadata.QueryCondition{Condition: cond}
	rsp, err := c.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.ErrorJSON("find classification failed, err: %s, cond: %s, rid: %s", err, cond, kit.Rid)
		return nil, err
	}

	return rsp.Info, nil
}

// UpdateClassification update classification by id
func (c *classification) UpdateClassification(kit *rest.Kit, data mapstr.MapStr, id int64) error {

	data.Remove(metadata.ClassFieldClassificationID)
	data.Remove(metadata.ClassificationFieldID)

	_, err := c.isValid(kit, true, data)
	if err != nil {
		blog.ErrorJSON("validate classification update data failed, err: %s, data: %s, rid: %s", err, data, kit.Rid)
		return err
	}

	// generate audit log of object classification.
	audit := auditlog.NewObjectClsAuditLog(c.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, id, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update object classification, objClsID: %d, err: %v, rid: %s",
			id, err, kit.Rid)
		return err
	}

	// update classification
	input := metadata.UpdateOption{
		Condition: map[string]interface{}{
			metadata.ClassificationFieldID:      id,
			metadata.ClassFieldClassificationID: mapstr.MapStr{common.BKDBNE: metadata.ClassificationUncategorizedID},
		},
		Data: data,
	}
	_, err = c.clientSet.CoreService().Model().UpdateModelClassification(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("update model classification failed, err: %v, input: %#v, rid: %s", err, data, kit.Rid)
		return err
	}

	// save audit log
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object classification success, but save audit log failed, objClsID: %d, err: %v, rid: %s",
			id, err, kit.Rid)
		return err
	}

	return nil
}

// isValid validate the classification data
func (c *classification) isValid(kit *rest.Kit, isUpdate bool, data mapstr.MapStr) (*metadata.Classification, error) {
	cls := new(metadata.Classification)
	if err := mapstruct.Decode2Struct(data, cls); err != nil {
		blog.Errorf("parse model classification failed, err: %v, input: %#v, rid: %s", err, data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}

	if !isUpdate || data.Exists(metadata.ClassFieldClassificationID) {
		if err := util.ValidModelIDField(data[metadata.ClassFieldClassificationID],
			metadata.ClassFieldClassificationID, kit.CCError); err != nil {
			return nil, err
		}
	}

	if !isUpdate || data.Exists(metadata.ClassFieldClassificationName) {
		if err := util.ValidModelNameField(data[metadata.ClassFieldClassificationName],
			metadata.ClassFieldClassificationName, kit.CCError); err != nil {
			return nil, err
		}
	}
	return cls, nil
}
