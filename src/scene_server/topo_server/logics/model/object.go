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
	"fmt"

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

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	// CreateObject create common object
	CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr) (*metadata.Object, error)
	// DeleteObject delete model by query condition
	DeleteObject(kit *rest.Kit, cond mapstr.MapStr, needCheckInst bool) (*metadata.Object, error)
	// FindObjectTopo search object topo by condition
	FindObjectTopo(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.ObjectTopo, error)
	// FindSingleObject find a object by objectID
	FindSingleObject(kit *rest.Kit, field []string, objectID string) (*metadata.Object, error)
	// UpdateObject update a common object by id
	UpdateObject(kit *rest.Kit, data mapstr.MapStr, id int64) error
	// IsObjectExist check whether objID is a real model's bk_obj_id field in backend
	IsObjectExist(kit *rest.Kit, objID string) (bool, error)
	// CreateObjectByImport create object by import yaml
	CreateObjectByImport(kit *rest.Kit, data []metadata.YamlObject) ([]metadata.Object, error)
	// SearchObjectsWithTotalInfo search object with it's attribute and association
	SearchObjectsWithTotalInfo(kit *rest.Kit, ids, excludedAsst []int64) (*metadata.TotalObjectInfo, error)
}

// NewObjectOperation create a new object operation instance
func NewObjectOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) ObjectOperationInterface {

	return &object{
		clientSet:   client,
		authManager: authManager,
	}
}

type object struct {
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// IsObjectExist check whether objID is a real model's bk_obj_id field in backend
func (o *object) IsObjectExist(kit *rest.Kit, objID string) (bool, error) {

	checkObjCond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}

	objItems, err := o.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjDes, []map[string]interface{}{checkObjCond})
	if err != nil {
		blog.Errorf("failed to search object(%s), err: %v, rid: %s", objID, err, kit.Rid)
		return false, err
	}

	if objItems[0] == 0 {
		return false, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return true, nil
}

// FindSingleObject find a object by objectID
func (o *object) FindSingleObject(kit *rest.Kit, field []string, objectID string) (*metadata.Object, error) {

	queryCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKObjIDField: objectID},
		Fields:         field,
		DisableCounter: true,
	}
	objs, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("get model failed, failed to get model objects(%s), err: %v, rid: %s", objectID, err, kit.Rid)
		return nil, err
	}

	if len(objs.Info) == 0 {
		blog.Errorf("get model failed, objects(%s) not found, result: %+v, rid: %s", objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommNotFound).Error())
	}

	if len(objs.Info) > 1 {
		blog.Errorf("get model failed, objects(%s) get multiple, result: %+v, rid: %s", objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommGetMultipleObject).Error())
	}

	return &objs.Info[0], nil
}

// CreateObject create common object
func (o *object) CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr) (*metadata.Object, error) {

	obj, err := o.isValid(kit, false, data)
	if err != nil {
		blog.Errorf("valid data(%#v) failed, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	exist, err := o.isClassificationExist(kit, obj.ObjCls)
	if err != nil {
		blog.Errorf("check classification failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if !exist {
		blog.Errorf("classification (%s) is non-exist, cannot create object, rid: %s", obj.ObjCls, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectClassificationSelectFailed)
	}

	if len(obj.ObjIcon) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIconField)
	}

	objRsp, err := o.clientSet.CoreService().Model().CreateModel(kit.Ctx, kit.Header, &metadata.CreateModel{Spec: *obj})
	if err != nil {
		blog.Errorf("create object(%s) failed, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	obj.ID = int64(objRsp.Created.ID)
	// create the default group
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    NewGroupID(true),
		ObjectID:   obj.ObjectID,
		OwnerID:    obj.OwnerID,
	}

	_, err = o.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header,
		obj.ObjectID, metadata.CreateModelAttributeGroup{Data: groupData})
	if err != nil {
		blog.Errorf("create attribute group[%s] failed, err: %v, rid: %s", groupData.GroupID, err, kit.Rid)
		return nil, err
	}

	attrIDs, err := o.createDefaultAttrs(kit, isMainline, obj, groupData)
	if err != nil {
		return nil, err
	}

	keys := make([]metadata.UniqueKey, 0)
	for _, id := range attrIDs {
		keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: id})
	}

	uni := metadata.ObjectUnique{
		ObjID:   obj.ObjectID,
		OwnerID: kit.SupplierAccount,
		Keys:    keys,
		Ispre:   false,
	}
	// NOTICE: 唯一索引与index.MainLineInstanceUniqueIndex,index.InstanceUniqueIndex定义强依赖
	// 原因：建立模型之前要将表和表中的索引提前建立，mongodb 4.2.6(4.4之前)事务中不能建表，事务操作表中数据操作和建表，建立索引为互斥操作。
	_, err = o.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header,
		uni.ObjID, metadata.CreateModelAttrUnique{Data: uni})
	if err != nil {
		blog.Errorf("create unique for %s failed, err: %v, rid: %s", uni.ObjID, err, kit.Rid)
		return nil, err
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewObjectAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return obj, nil
}

// DeleteObject delete model by query condition
func (o *object) DeleteObject(kit *rest.Kit, cond mapstr.MapStr, needCheckInst bool) (*metadata.Object, error) {

	// get model by conditon
	query := &metadata.QueryCondition{
		Condition:      cond,
		Fields:         make([]string, 0),
		DisableCounter: true,
	}

	objs, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("failed to find objects by query(%#v), err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}
	// shouldn't return nil, 404 error here, legacy implements just ignore not found error
	if len(objs.Info) == 0 {
		blog.V(3).Infof("object not found, rid: %s", kit.Rid)
		return nil, nil
	}

	if len(objs.Info) > 1 {
		return nil, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}

	obj := objs.Info[0]

	// check whether it can be deleted
	if needCheckInst {
		if err = o.canDelete(kit, obj.ObjectID); err != nil {
			return nil, err
		}
	}

	// generate audit log of object.
	audit := auditlog.NewObjectAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, &obj)
	if err != nil {
		blog.Errorf("generate audit log failed before delete object, objName: %s, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return nil, err
	}

	// DeleteModelCascade 将会删除模型/模型属性/属性分组/唯一校验
	_, err = o.clientSet.CoreService().Model().DeleteModelCascade(kit.Ctx, kit.Header, obj.ID)
	if err != nil {
		blog.Errorf("delete the object by the id(%d) failed, err: %v, rid: %s", obj.ID, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("delete object %s success, save audit log failed, err: %v, rid: %s", obj.ObjectName, err,
			kit.Rid)
		return nil, err
	}

	return &obj, nil
}

// canDelete return nil only when:
// 1. not inner model
// 2. model has no instances
// 3. model has no association with other model
func (o *object) canDelete(kit *rest.Kit, objID string) error {
	// step 1. ensure not inner model
	if common.IsInnerModel(objID) {
		return kit.CCError.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	cond := mapstr.MapStr{common.BKObjIDField: objID}

	// step 2. ensure model has no instances
	findInstResponse, err := o.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header, objID,
		&metadata.Condition{Condition: cond})

	if err != nil {
		blog.Errorf("failed to check if object (%s) has insts, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	if findInstResponse.Count != 0 {
		blog.Errorf("the object [%s] has been instantiated and cannot be deleted, rid: %s", objID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete, objID)
	}

	// step 3. ensure model has no association with other model
	condition := []map[string]interface{}{{
		common.BKDBOR: []mapstr.MapStr{{common.BKObjIDField: objID}, {common.AssociatedObjectIDField: objID}}},
	}

	assocCnt, err := o.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjAsst, condition)
	if err != nil {
		blog.Errorf("get object[%s] associate info failed, err: %v, rid: %s", objID, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(assocCnt) != 1 {
		blog.Errorf("get assoc num by filter failed, return answer is not only one, rid: %s", objID, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoObjectAssociationNotExist)
	}

	if assocCnt[0] != 0 {
		blog.Errorf("object[%s] has already associate to another one, rid: %s", objID, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoObjectHasAlreadyAssociated)
	}

	return nil
}

// FindObjectTopo search object topo by condition
func (o *object) FindObjectTopo(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.ObjectTopo, error) {

	// search object by objID
	queryObj := &metadata.QueryCondition{
		Condition: cond,
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKClassificationIDField,
			common.BkSupplierAccount, "position"},
		DisableCounter: true,
	}
	objs, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryObj)
	if err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if len(objs.Info) == 0 {
		return make([]metadata.ObjectTopo, 0), nil
	}

	objectIDs := make([]string, 0)
	objMap := make(map[string]metadata.Object, 0)
	for _, item := range objs.Info {
		objectIDs = append(objectIDs, item.ObjectID)
		objMap[item.ObjectID] = item
	}

	// search object association
	queryAsst := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs}},
				{common.BKAsstObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs}},
			},
		},
		Fields:         []string{common.AssociationKindIDField, common.BKObjIDField, common.BKAsstObjIDField},
		DisableCounter: true,
	}
	asstItems, err := o.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, queryAsst)
	if err != nil {
		blog.Errorf("failed to search the objects(%#v) association info , err: %v, rid: %s", objectIDs, err, kit.Rid)
		return nil, err
	}

	if len(asstItems.Info) == 0 {
		return []metadata.ObjectTopo{}, nil
	}

	asstKinds := make([]string, 0)
	asstObjIDs := make([]string, 0)
	assocObjsMap := map[string]map[string]struct{}{}
	for _, item := range asstItems.Info {
		asstKinds = append(asstKinds, item.AsstKindID)
		assocObjsMap[item.ObjectID] = map[string]struct{}{item.AsstObjID: {}}
		if _, exist := objMap[item.ObjectID]; exist {
			asstObjIDs = append(asstObjIDs, item.AsstObjID)
		}
	}

	asstObjIDs = util.RemoveDuplicatesAndEmptyByMap(asstObjIDs)

	// search association type
	query := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKinds}},
		Fields:         []string{common.AssociationKindNameField, common.AssociationKindIDField},
		DisableCounter: true,
	}
	assocType, err := o.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("get association kind[%#v] failed, err: %v, rid: %s", asstKinds, err, kit.Rid)
		return nil, err
	}

	if len(assocType.Info) == 0 {
		blog.Errorf("get association kind[%#v] failed, err: can not find this association kind, rid: %s",
			asstKinds, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoAsstKindIsNotExist)
	}

	assocTypeMap := make(map[string]*metadata.AssociationKind, 0)
	for _, assoType := range assocType.Info {
		assocTypeMap[assoType.AssociationKindID] = assoType
	}

	// search asst object by asstObjID
	cond = mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: asstObjIDs}}
	asstObjs, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: cond,
			Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKClassificationIDField,
				common.BkSupplierAccount, "position"},
		},
	)

	if err != nil {
		blog.Errorf("find object failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	asstObjMap := make(map[string]metadata.Object, 0)
	for _, asstObj := range asstObjs.Info {
		asstObjMap[asstObj.ObjectID] = asstObj
	}

	results := make([]metadata.ObjectTopo, 0)
	for _, assoc := range asstItems.Info {
		if _, exist := objMap[assoc.ObjectID]; !exist {
			continue
		}

		tmp := metadata.ObjectTopo{}
		tmp.Label = assocTypeMap[assoc.AsstKindID].AssociationKindName
		tmp.LabelName = assocTypeMap[assoc.AsstKindID].AssociationKindName
		tmp.From.ObjID = objMap[assoc.ObjectID].ObjectID
		tmp.From.ClassificationID = objMap[assoc.ObjectID].ObjCls
		tmp.From.Position = objMap[assoc.ObjectID].Position
		tmp.From.OwnerID = objMap[assoc.ObjectID].OwnerID
		tmp.From.ObjName = objMap[assoc.ObjectID].ObjectName
		tmp.To.OwnerID = asstObjMap[assoc.AsstObjID].OwnerID
		tmp.To.ObjID = asstObjMap[assoc.AsstObjID].ObjectID
		tmp.To.ClassificationID = asstObjMap[assoc.AsstObjID].ObjCls
		tmp.To.Position = asstObjMap[assoc.AsstObjID].Position
		tmp.To.ObjName = asstObjMap[assoc.AsstObjID].ObjectName

		if result, exist := assocObjsMap[assoc.AsstObjID]; exist {
			if _, exist := result[assoc.ObjectID]; exist {
				tmp.Arrows = "to,from"
			}
		} else {
			tmp.Arrows = "to"
		}

		results = append(results, tmp)
	}

	return results, nil
}

func (o *object) isClassificationValid(kit *rest.Kit, data mapstr.MapStr) error {

	if !data.Exists(metadata.ModelFieldObjCls) {
		return nil
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.ModelFieldObjCls: data[metadata.ModelFieldObjCls],
		},
	}
	rsp, err := o.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("failed to read model classification, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if len(rsp.Info) <= 0 {
		blog.Errorf("no model classification founded, err: %s, rid: %s",
			kit.CCError.CCError(common.CCErrorModelClassificationNotFound), kit.Rid)
		return kit.CCError.CCError(common.CCErrorModelClassificationNotFound)
	}
	return nil
}

// UpdateObject update a common object by id
func (o *object) UpdateObject(kit *rest.Kit, data mapstr.MapStr, id int64) error {

	obj, err := o.isValid(kit, true, data)
	if err != nil {
		blog.Errorf("valid data failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return err
	}

	obj.ID = id

	// remove unchangeable fields.
	data.Remove(metadata.ModelFieldObjectID)
	data.Remove(metadata.ModelFieldID)

	if err := o.isClassificationValid(kit, data); err != nil {
		return err
	}
	// generate audit log of object attribute group.
	audit := auditlog.NewObjectAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).
		WithUpdateFields(data)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update object, objName: %s, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return err
	}

	input := metadata.UpdateOption{
		Condition: mapstr.MapStr{common.BKFieldID: obj.ID},
		Data:      data,
	}

	_, err = o.clientSet.CoreService().Model().UpdateModel(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("update object failed, id: %d, err: %v, rid: %s", obj.ID, err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object %s success, but save audit log failed, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return err
	}

	return nil
}

func (o *object) isValid(kit *rest.Kit, isUpdate bool, data mapstr.MapStr) (*metadata.Object, error) {

	obj := new(metadata.Object)
	if err := mapstruct.Decode2Struct(data, obj); err != nil {
		blog.Errorf("parse object failed, err: %v, input: %#v, rid: %s", err, data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectID) {
		if err := util.ValidModelIDField(data[metadata.ModelFieldObjectID],
			metadata.ModelFieldObjectID, kit.CCError); err != nil {
			blog.Errorf("failed to valid the object id(%s), rid: %s", metadata.ModelFieldObjectID, kit.Rid)
			return nil, err
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectName) {
		if err := util.ValidModelNameField(data[metadata.ModelFieldObjectName],
			metadata.ModelFieldObjectName, kit.CCError); err != nil {
			blog.Errorf("failed to valid the object name(%s), rid: %s", metadata.ModelFieldObjectName, kit.Rid)
			return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid, metadata.ModelFieldObjectName+" "+err.Error())
		}
	}

	if !isUpdate && !data.Exists(metadata.ModelFieldObjCls) {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, metadata.ModelFieldObjCls)
	}

	if !isUpdate && !metadata.IsCommon(obj.ObjectID) {
		return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid,
			fmt.Sprintf("'%s' the built-in object id, please use a new one", obj.ObjectID))
	}

	obj.OwnerID = kit.SupplierAccount
	return obj, nil
}

func (o *object) isClassificationExist(kit *rest.Kit, clsID string) (bool, error) {

	filter := []map[string]interface{}{{common.BKClassificationIDField: clsID}}
	objCls, err := o.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjClassification, filter)
	if err != nil {
		blog.Errorf("get object classification failed, err: %v, rid: %s", err, kit.Rid)
		return false, err
	}

	if objCls[0] == 0 {
		return false, nil
	}

	return true, nil
}

func (o *object) createDefaultAttrs(kit *rest.Kit, isMainline bool, obj *metadata.Object,
	groupData metadata.Group) ([]uint64, error) {

	attrs := make([]metadata.Attribute, 0)
	attrs = append(attrs, metadata.Attribute{
		ObjectID:          obj.ObjectID,
		IsOnly:            true,
		IsPre:             true,
		Creator:           "user",
		IsEditable:        true,
		PropertyIndex:     -1,
		PropertyGroup:     groupData.GroupID,
		PropertyGroupName: groupData.GroupName,
		IsRequired:        true,
		PropertyType:      common.FieldTypeSingleChar,
		PropertyID:        common.GetInstNameField(obj.ObjectID),
		PropertyName:      common.DefaultInstName,
		OwnerID:           kit.SupplierAccount,
	})

	if isMainline {
		attrs = append(attrs, metadata.Attribute{
			ObjectID:          obj.ObjectID,
			IsOnly:            true,
			IsPre:             true,
			Creator:           "system",
			IsEditable:        true,
			IsSystem:          true,
			PropertyIndex:     -1,
			PropertyGroup:     groupData.GroupID,
			PropertyGroupName: groupData.GroupName,
			IsRequired:        true,
			PropertyType:      common.FieldTypeInt,
			PropertyID:        common.BKInstParentStr,
			PropertyName:      common.BKInstParentStr,
			OwnerID:           kit.SupplierAccount,
		})
	}

	param := &metadata.CreateModelAttributes{Attributes: attrs}
	rspAttr, err := o.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, obj.ObjectID, param)
	if err != nil {
		blog.Errorf("create model(%s) attrs failed, input: %#v, err: %v, rid: %s", obj.ObjectID, param, err, kit.Rid)
		return nil, err
	}

	for _, exception := range rspAttr.Exceptions {
		return nil, kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(rspAttr.Repeated) > 0 {
		blog.Errorf("attr(%#v) is duplicated, objID: %s, rid: %s", rspAttr.Repeated, obj.ObjectID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	result := make([]uint64, 0)
	for _, item := range rspAttr.Created {
		result = append(result, item.ID)
	}

	return result, nil
}

// CreateObjectByImport create object by import yaml
func (o *object) CreateObjectByImport(kit *rest.Kit, data []metadata.YamlObject) ([]metadata.Object, error) {

	assts := make([]metadata.AsstWithAsstObjInfo, 0)
	objs := make([]metadata.Object, 0)
	for _, objInfo := range data {

		if err := objInfo.Validate(); err.ErrCode != 0 {
			blog.Errorf("validate object info failed, objectinfo: %+v, rid: %s", objInfo, kit.Rid)
			return nil, err.ToCCError(kit.CCError)
		}

		object := mapstr.MapStr{
			common.BKObjIDField:            objInfo.ObjectID,
			common.BKObjNameField:          objInfo.ObjectName,
			common.BKObjIconField:          objInfo.ObjIcon,
			common.BKClassificationIDField: objInfo.ClsID,
			common.CreatorField:            kit.User,
		}

		obj, err := o.isValid(kit, false, object)
		if err != nil {
			blog.Errorf("valid data(%#v) failed, err: %v, rid: %s", objInfo, err, kit.Rid)
			return nil, err
		}

		exist, err := o.isClassificationExist(kit, obj.ObjCls)
		if err != nil {
			blog.Errorf("check classification failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if !exist {
			cls := metadata.Classification{
				ClassificationID:   objInfo.ClsID,
				ClassificationName: objInfo.ClsName,
			}
			if err := o.createClassification(kit, cls); err != nil {
				blog.Errorf("create classification failed, err: %v, rid: %s", err, kit.Rid)
				return nil, err
			}
		}

		input := &metadata.CreateModel{Spec: *obj}
		objRsp, err := o.clientSet.CoreService().Model().CreateModel(kit.Ctx, kit.Header, input)
		if err != nil {
			blog.Errorf("create object(%s) failed, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
			return nil, err
		}
		obj.ID = int64(objRsp.Created.ID)

		if err := o.createObjectAttr(kit, obj.ObjectID, objInfo.ObjectAttr, objInfo.ObjectAttrUnique); err != nil {
			blog.Errorf("create object(%s) failed, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
			return nil, err
		}

		assts = append(assts, objInfo.ObjectAsst...)
		objs = append(objs, *obj)

		// generate audit log of object.
		audit := auditlog.NewObjectAuditLog(o.clientSet.CoreService())
		generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, int64(objRsp.Created.ID), nil)
		if err != nil {
			blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		// save audit log.
		if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
			blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
	}

	if err := o.createObjectAssociation(kit, assts); err != nil {
		blog.Errorf("create object associations failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return objs, nil
}

func (o *object) createClassification(kit *rest.Kit, data metadata.Classification) error {
	// create classification
	input := metadata.CreateOneModelClassification{Data: data}
	rsp, err := o.clientSet.CoreService().Model().CreateModelClassification(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("create classification %s failed, err: %v, rid: %s", data.ClassificationName, err, kit.Rid)
		return err
	}

	// get created model classification data by id
	clsReq := &metadata.QueryCondition{Condition: mapstr.MapStr{metadata.ClassificationFieldID: int64(rsp.Created.ID)}}
	clsResp, err := o.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, clsReq)
	if err != nil {
		blog.Errorf("get created model classification by id(%d) failed, err: %v, rid: %s", rsp.Created.ID, err, kit.Rid)
		return err
	}

	if len(clsResp.Info) != 1 {
		blog.Errorf("get created model classification by id(%d) returns not one cls, rid: %s", rsp.Created.ID, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommNotFound)
	}

	cls := &clsResp.Info[0]
	// generate audit log of object classification.
	audit := auditlog.NewObjectClsAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, cls.ID, cls)
	if err != nil {
		blog.Errorf("create object classification %s success, but generate audit log failed, err: %v, rid: %s",
			cls.ClassificationName, err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object classification %s success, but save audit log failed, err: %v, rid: %s",
			cls.ClassificationName, err, kit.Rid)
		return err
	}

	return nil
}

func (o *object) createObjectAttr(kit *rest.Kit, objID string, attr []metadata.Attribute, attrUniq [][]string) error {

	var groupIndex int64 = 1
	createdGroup := map[string]struct{}{}
	attrs := make([]metadata.Attribute, 0)
	for _, item := range attr {
		if _, exist := createdGroup[item.PropertyGroupName]; !exist {

			err := o.createObjectAttrGroup(kit, objID, item.PropertyGroup, item.PropertyGroupName, groupIndex)
			if err != nil {
				blog.Errorf("create attribute group[%s] failed, err: %v, rid: %s", item.PropertyGroup, err, kit.Rid)
				return err
			}

			groupIndex += 1
			createdGroup[item.PropertyGroupName] = struct{}{}
		}

		item.Creator = kit.User
		attrs = append(attrs, item)
	}

	param := &metadata.CreateModelAttributes{Attributes: attrs}
	rspAttr, err := o.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, objID, param)
	if err != nil {
		blog.Errorf("create model(%s) attrs failed, input: %#v, err: %v, rid: %s", objID, param, err, kit.Rid)
		return err
	}

	for _, exception := range rspAttr.Exceptions {
		return kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(rspAttr.Repeated) > 0 {
		blog.Errorf("attr(%#v) is duplicated, objID: %s, rid: %s", rspAttr.Repeated, objID, kit.Rid)
		return kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	attrIDMap := make(map[string]uint64)
	for _, attr := range rspAttr.Created {
		attrIDMap[attrs[attr.OriginIndex].PropertyID] = attr.ID
	}

	for _, item := range attrUniq {
		keys := make([]metadata.UniqueKey, 0)
		for _, uniqAttrID := range item {
			keys = append(keys, metadata.UniqueKey{
				Kind: metadata.UniqueKeyKindProperty,
				ID:   attrIDMap[uniqAttrID],
			})
		}

		cond := metadata.CreateModelAttrUnique{Data: metadata.ObjectUnique{
			ObjID:   objID,
			OwnerID: kit.SupplierAccount,
			Keys:    keys,
			Ispre:   false,
		}}
		_, err = o.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header, objID, cond)
		if err != nil {
			blog.Errorf("create unique for %s failed, err: %v, rid: %s", objID, err, kit.Rid)
			return err
		}
	}

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	for _, item := range rspAttr.Created {
		attrs[item.OriginIndex].ID = int64(item.ID)
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, int64(item.ID), &attrs[item.OriginIndex])
		if err != nil {
			blog.Errorf("generate audit log after creating attr %s failed, err: %v, rid: %s",
				attrs[item.OriginIndex].PropertyName, err, kit.Rid)
			return err
		}

		// save audit log.
		if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
			blog.Errorf("save audit log after creating attr %s failed, err: %v, rid: %s",
				attrs[item.OriginIndex].PropertyName, err, kit.Rid)
			return err
		}
	}

	return nil
}

func (o *object) createObjectAttrGroup(kit *rest.Kit, objID, groupID, groupName string, groupIndex int64) error {
	group := metadata.Group{
		GroupID:    groupID,
		GroupName:  groupName,
		GroupIndex: groupIndex,
		ObjectID:   objID,
		OwnerID:    kit.SupplierAccount,
	}
	groupParams := metadata.CreateModelAttributeGroup{Data: group}
	rsp, err := o.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header, objID, groupParams)
	if err != nil {
		blog.Errorf("create attribute group[%s] failed, err: %v, rid: %s", group.GroupID, err, kit.Rid)
		return err
	}

	group.ID = int64(rsp.Created.ID)

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, group.ID, &group)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but generate audit log failed, err: %v, rid: %s",
			group.GroupName, err, kit.Rid)
		return err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute group %s success, but save audit log failed, err: %v, rid: %s",
			group.GroupName, err, kit.Rid)
		return err
	}

	return nil
}

func (o *object) createObjectAssociation(kit *rest.Kit, data []metadata.AsstWithAsstObjInfo) error {

	for _, item := range data {

		// if the on delete action is empty, set none as default.
		if len(item.OnDelete) == 0 {
			item.OnDelete = metadata.NoAction
		}

		// check source object exists
		queryCond := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				common.BKObjIDField: mapstr.MapStr{common.BKDBIN: []string{item.ObjectID, item.AsstObjID}},
			},
			Fields: []string{common.BKObjIDField},
		}
		objRsp, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
		if err != nil {
			blog.Errorf("read the object(%s) failed, err: %v, rid: %s", item.ObjectID, err, kit.Rid)
			return err
		}

		if len(objRsp.Info) == 0 {
			blog.Errorf("object(%s) and asstObject(%s) is invalid, rid: %s", item.ObjectID, item.AsstObjID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "bk_obj_id&bk_asst_obj_id")
		}

		checkMap := make(map[string]struct{})
		for _, item := range objRsp.Info {
			checkMap[item.ObjectID] = struct{}{}
		}

		if _, exist := checkMap[item.ObjectID]; !exist {
			blog.Errorf("object(%s) is invalid, rid: %s", item.ObjectID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
		}

		if _, exist := checkMap[item.AsstObjID]; !exist {
			blog.Errorf("object(%s) is invalid, rid: %s", item.AsstObjID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKAsstObjIDField)
		}

		cond := &metadata.CreateModelAssociation{Spec: item.Association}
		_, err = o.clientSet.CoreService().Association().CreateModelAssociation(kit.Ctx, kit.Header, cond)
		if err != nil {
			blog.Errorf("create object association failed, param: %#v , err: %v, rid: %s", cond, err, kit.Rid)
			return err
		}
	}

	return nil
}

// SearchObjectsWithTotalInfo search object with it's attribute and association
func (o *object) SearchObjectsWithTotalInfo(kit *rest.Kit, ids, excludedAsst []int64) (*metadata.TotalObjectInfo,
	error) {

	if len(ids) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, "object_id")
	}

	objCond := metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: ids}},
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKClassificationIDField,
			common.BKObjIconField, common.BKIsPre},
		DisableCounter: true,
	}
	objs, err := o.searchObjectByCondition(kit, objCond)
	if err != nil {
		blog.Errorf("get object by ids failed, ids[%v], err: %v, rid: %s", ids, err, kit.Rid)
		return nil, err
	}

	objIDs := make([]string, 0)
	clsIDs := make([]string, 0)
	for _, obj := range objs {
		objIDs = append(objIDs, obj.ObjectID)
		clsIDs = append(clsIDs, obj.ObjCls)
	}

	input := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKClassificationIDField: mapstr.MapStr{common.BKDBIN: clsIDs}},
		Fields:         []string{common.BKClassificationNameField, common.BKClassificationIDField},
		DisableCounter: true,
	}
	clsRsp, err := o.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("find classification failed, cond: %v, err: %v, rid: %s", input, err, kit.Rid)
		return nil, err
	}

	if len(clsRsp.Info) == 0 {
		blog.Errorf("find classification failed, result is empty, cond: %v, rid: %s", input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorModelClassificationNotFound)
	}

	clsMap := make(map[string]string)
	for _, item := range clsRsp.Info {
		clsMap[item.ClassificationID] = item.ClassificationName
	}

	attrCond := metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}},
		DisableCounter: true,
	}
	attrRsp, attrUni, err := o.searchObjAttrByCondition(kit, attrCond)
	if err != nil {
		blog.Errorf("find classification failed, cond: %v, err: %v, rid: %s", attrCond, err, kit.Rid)
		return nil, err
	}

	asstCond := metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs},
			common.BKFieldID: mapstr.MapStr{common.BKDBNIN: excludedAsst}},
		DisableCounter: true,
	}
	asstRsp, asstKind, err := o.searchObjAsstByCondition(kit, asstCond)
	if err != nil {
		blog.Errorf("find objects associations failed, cond: %v, err: %v, rid: %s", asstCond, err, kit.Rid)
		return nil, err
	}

	result := make(map[string]interface{}, 0)
	for objID, obj := range objs {
		objInfo := mapstr.MapStr{
			common.BKObjIDField:            obj.ObjectID,
			common.BKObjNameField:          obj.ObjectName,
			common.BKObjIconField:          obj.ObjIcon,
			common.BKClassificationIDField: obj.ObjCls,
			common.BKIsPre:                 obj.IsPre,
		}
		objInfo[common.BKClassificationNameField] = clsMap[obj.ObjCls]
		objInfo["object_attr"] = attrRsp[objID]
		objInfo["object_attr_unique"] = attrUni[objID]
		objInfo["object_asst"] = asstRsp[objID]
		result[objID] = objInfo
	}

	return &metadata.TotalObjectInfo{Object: result, Asst: asstKind}, nil
}

func (o *object) searchObjectByCondition(kit *rest.Kit, cond metadata.QueryCondition) (map[string]metadata.Object,
	error) {
	objs, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, &cond)
	if err != nil {
		blog.Errorf("get object by ids failed, cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, err
	}

	result := make(map[string]metadata.Object)
	for _, item := range objs.Info {
		result[item.ObjectID] = item
	}

	return result, nil
}

func (o *object) searchObjAttrByCondition(kit *rest.Kit, cond metadata.QueryCondition) (map[string][]mapstr.MapStr,
	map[string][][]string, error) {

	attrRsp, err := o.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, &cond)
	if err != nil {
		blog.Errorf("find object attribute failed, cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, nil, err
	}

	objIDs := make([]string, 0)
	attrMap := make(map[string][]metadata.Attribute)
	attrIDMap := make(map[int64]string)
	for _, item := range attrRsp.Info {
		if _, exist := attrMap[item.ObjectID]; !exist {
			attrMap[item.ObjectID] = make([]metadata.Attribute, 0)
		}

		attrMap[item.ObjectID] = append(attrMap[item.ObjectID], item)
		objIDs = append(objIDs, item.ObjectID)
		attrIDMap[item.ID] = item.PropertyID
	}

	queryCond := metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}},
		DisableCounter: true,
	}
	grpRsp, err := o.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("find objects attributes group failed, cond: %v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, nil, err
	}

	grpMap := make(map[string]map[string]string)
	for _, item := range grpRsp.Info {
		if _, exist := grpMap[item.ObjectID]; !exist {
			grpMap[item.ObjectID] = make(map[string]string)
		}
		grpMap[item.ObjectID][item.GroupID] = item.GroupName
	}

	uniCond := metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objIDs}},
		DisableCounter: true,
	}
	uniRsp, err := o.clientSet.CoreService().Model().ReadModelAttrUnique(kit.Ctx, kit.Header, uniCond)
	if err != nil {
		blog.Errorf("find objects attributes group failed, cond: %v, err: %v, rid: %s", uniCond, err, kit.Rid)
		return nil, nil, err
	}

	attrID := make(map[string][][]string)
	for _, item := range uniRsp.Info {
		if _, exist := attrID[item.ObjID]; !exist {
			attrID[item.ObjID] = make([][]string, 0)
		}

		uniqIDs := make([]string, 0)
		for _, id := range item.Keys {
			uniqIDs = append(uniqIDs, attrIDMap[int64(id.ID)])
		}

		attrID[item.ObjID] = append(attrID[item.ObjID], uniqIDs)
	}

	result := make(map[string][]mapstr.MapStr)
	for objID, attrs := range attrMap {
		for _, attr := range attrs {
			attrInfo := attr.ToMapStr()
			attrInfo["bk_property_group_name"] = grpMap[objID][attr.PropertyGroup]
			delete(attrInfo, common.BKFieldID)
			delete(attrInfo, common.CreateTimeField)
			delete(attrInfo, "bk_issystem")
			delete(attrInfo, "bk_isapi")
			delete(attrInfo, common.LastTimeField)
			delete(attrInfo, common.BkSupplierAccount)
			delete(attrInfo, common.BKAppIDField)
			delete(attrInfo, common.CreatorField)
			result[objID] = append(result[objID], attrInfo)
		}
	}

	return result, attrID, nil
}

func (o *object) searchObjAsstByCondition(kit *rest.Kit, cond metadata.QueryCondition) (
	map[string][]mapstr.MapStr, []mapstr.MapStr, error) {

	asstRsp, err := o.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, &cond)
	if err != nil {
		blog.Errorf("find objects associations failed, cond: %v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, nil, err
	}

	asstObjID := make([]string, 0)
	asstKindID := make([]string, 0)
	for _, asst := range asstRsp.Info {
		asstObjID = append(asstObjID, asst.AsstObjID)
		asstKindID = append(asstKindID, asst.AsstKindID)
	}

	asstObjCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: asstObjID}},
		DisableCounter: true,
	}
	asstObjRsp, err := o.searchObjectByCondition(kit, *asstObjCond)
	if err != nil {
		blog.Errorf("find asst objects failed, cond: %v, err: %v, rid: %s", asstObjCond, err, kit.Rid)
		return nil, nil, err
	}

	asstMap := make(map[string][]mapstr.MapStr)
	for _, item := range asstRsp.Info {
		if _, exist := asstMap[item.ObjectID]; !exist {
			asstMap[item.ObjectID] = make([]mapstr.MapStr, 0)
		}

		asstInfo := item.ToMapStr()
		delete(asstInfo, common.BKFieldID)
		delete(asstInfo, common.BkSupplierAccount)
		delete(asstInfo, "on_delete")
		asstInfo["bk_asst_obj_name"] = asstObjRsp[item.AsstObjID].ObjectName
		asstInfo["bk_asst_obj_icon"] = asstObjRsp[item.AsstObjID].ObjIcon
		asstInfo["bk_asst_obj_ispre"] = asstObjRsp[item.AsstObjID].IsPre

		asstMap[item.ObjectID] = append(asstMap[item.ObjectID], asstInfo)
	}

	asstKindCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKindID}},
		DisableCounter: true,
	}

	asstKind, err := o.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header, asstKindCond)
	if err != nil {
		blog.Errorf("find asst kind failed, cond: %v, err: %v, rid: %s", asstKindCond, err, kit.Rid)
		return nil, nil, err
	}

	asstKindMapstr := make([]mapstr.MapStr, 0)
	for _, item := range asstKind.Info {
		itemMap := item.ToMapStr()
		delete(itemMap, common.BKFieldID)
		delete(itemMap, common.BkSupplierAccount)
		asstKindMapstr = append(asstKindMapstr, itemMap)
	}

	return asstMap, asstKindMapstr, nil
}
