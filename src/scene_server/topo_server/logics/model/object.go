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
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	// CreateObject create common object
	CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr) (*metadata.Object, error)
	// DeleteObject delete common object
	DeleteObject(kit *rest.Kit, id int64, needCheckInst bool) error
	// FindObject find object by condition
	FindObject(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Object, error)
	// FindObjectTopo search object topo by condition
	FindObjectTopo(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.ObjectTopo, error)
	// FindSingleObject find a object by objectID
	FindSingleObject(kit *rest.Kit, objectID string) (*metadata.Object, error)
	// UpdateObject update a common object by id
	UpdateObject(kit *rest.Kit, data mapstr.MapStr, id int64) error
	// IsValidObject check whether objID is a real model's bk_obj_id field in backend
	IsValidObject(kit *rest.Kit, objID string) error
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
	lang        language.DefaultCCLanguageIf
}

// IsValidObject check whether objID is a real model's bk_obj_id field in backend
func (o *object) IsValidObject(kit *rest.Kit, objID string) error {

	checkObjCond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}

	objItems, err := o.FindObject(kit, checkObjCond)
	if err != nil {
		blog.Errorf("failed to check the object repeated, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, err.Error())
	}

	if len(objItems) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return nil
}

func (o *object) FindSingleObject(kit *rest.Kit, objectID string) (*metadata.Object, error) {

	cond := mapstr.MapStr{
		common.BKObjIDField: objectID,
	}

	objs, err := o.FindObject(kit, cond)
	if err != nil {
		blog.Errorf("get model failed, failed to get model by supplier account(%s) objects(%s), "+
			"err: %v, rid: %s", kit.SupplierAccount, objectID, err, kit.Rid)
		return nil, err
	}

	if len(objs) == 0 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) not found, result: %+v, "+
			"rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommNotFound).Error())
	}

	if len(objs) > 1 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) get multiple, result: %+v, "+
			"rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed,
			kit.CCError.Error(common.CCErrCommGetMultipleObject).Error())
	}

	return &objs[0], nil
}

func (o *object) CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr) (*metadata.Object, error) {

	obj, err := o.IsValid(kit, false, data)
	if err != nil {
		blog.Errorf("valid data(%#v) failed, err: %v, rid: %s", data, err, kit.Rid)
		return nil, err
	}

	objCls, err := o.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKClassificationIDField: obj.ObjCls}})
	if err != nil {
		blog.Errorf("get object classification failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err = objCls.CCError(); err != nil {
		blog.Errorf("get object classification failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(objCls.Data.Info) == 0 {
		blog.Errorf("can't find classification by params, classification: %s is not exist, rid: %s",
			obj.ObjCls, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKClassificationIDField)
	}

	filter := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: obj.ObjectID},
			{common.BKObjNameField: obj.ObjectName},
		}}
	cnt, err := o.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjDes,
		[]map[string]interface{}{filter})
	if err != nil {
		blog.Errorf("get object number by filter failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	for index := range cnt {
		if cnt[index] != 0 {
			return nil, kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, "bk_obj_id/bk_obj_name")
		}
	}

	if len(obj.ObjIcon) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKObjIconField)
	}

	obj.OwnerID = kit.SupplierAccount

	objRsp, err := o.clientSet.CoreService().Model().CreateModel(kit.Ctx, kit.Header, &metadata.CreateModel{Spec: *obj})
	if err != nil {
		blog.Errorf("create object failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = objRsp.CCError(); err != nil {
		blog.Errorf("failed to search the object(%s), err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	obj.ID = int64(objRsp.Data.Created.ID)

	// create the default group
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    "default",
		ObjectID:   obj.ObjectID,
		OwnerID:    obj.OwnerID,
	}

	rspGrp, err := o.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header,
		obj.ObjectID, metadata.CreateModelAttributeGroup{Data: groupData})
	if err != nil {
		blog.Errorf("create attribute group failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rspGrp.CCError(); err != nil {
		blog.Errorf("create attribute group[%s] failed, err: %v, rid: %s", groupData.GroupID, err, kit.Rid)
		return nil, err
	}

	keys := make([]metadata.UniqueKey, 0)
	// create the default inst attribute
	attr := metadata.Attribute{
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
		PropertyID:        obj.GetInstNameFieldName(),
		PropertyName:      obj.GetDefaultInstPropertyName(),
		OwnerID:           kit.SupplierAccount,
	}

	// create a new record
	rspAttr, err := o.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header,
		attr.ObjectID, &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{attr}})
	if err != nil {
		blog.Errorf("failed to request coreService to create model attrs, err: %v, ObjectID: %s, input: %#v, rid: %s",
			err, attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rspAttr.CCError(); err != nil {
		blog.Errorf("create model attrs failed, ObjectID: %s, input: %#v, err: %v, rid: %s",
			attr.ObjectID, attr, err, kit.Rid)
		return nil, rspAttr.CCError()
	}

	for _, exception := range rspAttr.Data.Exceptions {
		return nil, kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(rspAttr.Data.Repeated) > 0 {
		blog.Errorf("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %#v, rid: %s",
			attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	if len(rspAttr.Data.Created) != 1 {
		blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %#v, rid: %s",
			attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}

	attr.ID = int64(rspAttr.Data.Created[0].ID)

	keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(attr.ID)})

	if isMainline {
		pAttr := metadata.Attribute{
			ObjectID:          obj.ObjectID,
			IsOnly:            true,
			IsPre:             true,
			Creator:           "system",
			IsEditable:        true,
			PropertyIndex:     -1,
			PropertyGroup:     groupData.GroupID,
			PropertyGroupName: groupData.GroupName,
			IsRequired:        true,
			PropertyType:      common.FieldTypeInt,
			PropertyID:        common.BKInstParentStr,
			PropertyName:      common.BKInstParentStr,
			IsSystem:          true,
			OwnerID:           kit.SupplierAccount,
		}

		rsppAttr, err := o.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header,
			pAttr.ObjectID, &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{pAttr}})
		if err != nil {
			blog.Errorf("failed to request coreService to create model attrs, err: %v, ObjectID: %s, input: %#v, "+
				"rid: %s", err, pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
		}

		if err = rsppAttr.CCError(); err != nil {
			blog.Errorf("create model attrs failed, ObjectID: %s, input: %#v, rid: %s", pAttr.ObjectID, pAttr, kit.Rid)
			return nil, err
		}

		for _, exception := range rsppAttr.Data.Exceptions {
			return nil, kit.CCError.New(int(exception.Code), exception.Message)
		}

		if len(rsppAttr.Data.Repeated) > 0 {
			blog.Errorf("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %#v, rid: %s",
				pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
		}

		if len(rsppAttr.Data.Created) != 1 {
			blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %#v, rid: %s",
				pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
		}
		pAttr.ID = int64(rsppAttr.Data.Created[0].ID)

		keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(pAttr.ID)})
	}

	uni := metadata.ObjectUnique{
		ObjID:   obj.ObjectID,
		OwnerID: kit.SupplierAccount,
		Keys:    keys,
		Ispre:   false,
	}
	// NOTICE: 2021年03月29日  唯一索引与index.MainLineInstanceUniqueIndex,index.InstanceUniqueIndex定义强依赖
	// 原因：建立模型之前要将表和表中的索引提前建立，mongodb 4.2.6(4.4之前)事务中不能建表，事务操作表中数据操作和建表，建立索引为互斥操作。
	resp, err := o.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header,
		uni.ObjID, metadata.CreateModelAttrUnique{Data: uni})
	if err != nil {
		blog.Errorf("create unique for %s failed, err: %v, rid: %s", uni.ObjID, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectUniqueCreateFailed)
	}
	if err = resp.CCError(); err != nil {
		blog.Errorf("create unique for %s failed, err: %v, rid: %s", uni.ObjID, err, kit.Rid)
		return nil, err
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewObjectAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, nil)
	if err != nil {
		blog.Errorf("create object %s success, but generate audit log failed, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object %s success, but save audit log failed, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return nil, err
	}

	return obj, nil
}

// DeleteObject delete model by id
func (o *object) DeleteObject(kit *rest.Kit, id int64, needCheckInst bool) error {
	if id <= 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
	}

	// get model by id
	cond := mapstr.MapStr{
		metadata.ModelFieldID: id,
	}

	objs, err := o.FindObject(kit, cond)
	if err != nil {
		blog.Errorf("failed to find objects, the condition is (%v) err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}
	// shouldn't return 404 error here, legacy implements just ignore not found error
	if len(objs) == 0 {
		blog.V(3).Infof("object not found, condition: %v, rid: %s", cond, kit.Rid)
		return nil
	}
	if len(objs) > 1 {
		return kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	obj := objs[0]

	// check whether it can be deleted
	if needCheckInst {
		if err = o.CanDelete(kit, obj); err != nil {
			return err
		}
	}

	// generate audit log of object.
	audit := auditlog.NewObjectAuditLog(o.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, &obj)
	if err != nil {
		blog.Errorf("generate audit log failed before delete object, objName: %s, err: %v, rid: %s",
			obj.ObjectName, err, kit.Rid)
		return err
	}

	// DeleteModelCascade 将会删除模型/模型属性/属性分组/唯一校验
	rsp, err := o.clientSet.CoreService().Model().DeleteModelCascade(kit.Ctx, kit.Header, id)
	if err != nil {
		blog.Errorf("failed to request the object controller, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to delete the object by the id(%d), err: %v, rid: %s", id, err, kit.Rid)
		return err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("delete object %s success, save audit log failed, err: %v, rid: %s", obj.ObjectName, err, kit.Rid)
		return err
	}

	return nil
}

// CanDelete return nil only when:
// 1. not inner model
// 2. model has no instances
// 3. model has no association with other model
func (o *object) CanDelete(kit *rest.Kit, targetObj metadata.Object) error {
	// step 1. ensure not inner model
	if common.IsInnerModel(targetObj.GetObjectID()) {
		return kit.CCError.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	cond := mapstr.New()
	if targetObj.IsCommon() {
		cond.Set(common.BKObjIDField, targetObj.ObjectID)
	}

	// step 2. ensure model has no instances
	input := &metadata.Condition{Condition: cond}
	findInstResponse, err := o.clientSet.CoreService().Instance().CountInstances(kit.Ctx, kit.Header,
		targetObj.ObjectID, input)

	if err != nil {
		blog.Errorf("failed to check if it (%s) has some insts, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
		return err
	}

	if err = findInstResponse.CCError(); err != nil {
		blog.Errorf("failed to check if it (%s) has some insts, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
		return err
	}

	if findInstResponse.Data.Count != 0 {
		blog.Errorf("the object [%s] has been instantiated and cannot be deleted, rid: %s",
			targetObj.ObjectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete, targetObj.ObjectID)
	}

	// step 3. ensure model has no association with other model
	or := make([]interface{}, 0)
	or = append(or, mapstr.MapStr{common.BKObjIDField: targetObj.ObjectID})
	or = append(or, mapstr.MapStr{common.AssociatedObjectIDField: targetObj.ObjectID})

	assocResult, err := o.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: map[string]interface{}{common.BKDBOR: or}})
	if err != nil {
		blog.Errorf("check object[%s] can be deleted, but get object associate info failed, err: %v, rid: %s",
			targetObj.ObjectID, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = assocResult.CCError(); err != nil {
		blog.Errorf("get object[%s] associate info failed, err: %v, rid: %s", targetObj.ObjectID, err, kit.Rid)
		return kit.CCError.Error(assocResult.Code)
	}

	if len(assocResult.Data.Info) != 0 {
		blog.Errorf("object[%s] has already associate to another one., rid: %s", targetObj.ObjectID, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoObjectHasAlreadyAssociated)
	}

	return nil
}

func (o *object) isFrom(kit *rest.Kit, fromObjID []string, cond mapstr.MapStr) (map[string]bool, error) {

	asstItems, err := o.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: fromObjID}},
			Fields:    []string{common.BKAsstObjIDField, common.BKObjIDField},
		})

	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if err = asstItems.CCError(); err != nil {
		blog.Errorf("search object[%s] association failed, err: %v, rid: %s", fromObjID, err, kit.Rid)
		return nil, err
	}

	result := make(map[string]bool, 0)
	for _, asst := range asstItems.Data.Info {
		result[asst.AsstObjID] = false
		if asst.AsstObjID == cond[asst.ObjectID] {
			result[asst.AsstObjID] = true
		}
	}

	return result, nil
}

func (o *object) FindObjectTopo(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.ObjectTopo, error) {

	// search object by objID
	objs, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: cond,
			Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BKClassificationIDField,
				common.BkSupplierAccount, "position"},
		},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = objs.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s",
			cond, err, kit.Rid)
		return nil, err
	}

	if len(objs.Data.Info) == 0 {
		return []metadata.ObjectTopo{}, nil
	}

	var objectIDs []string
	objMap := make(map[string]metadata.Object, 0)
	for index := range objs.Data.Info {
		objectIDs = append(objectIDs, objs.Data.Info[index].ObjectID)
		objMap[objs.Data.Info[index].ObjectID] = objs.Data.Info[index]
	}

	// search object association
	asstItems, err := o.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: mapstr.MapStr{common.BKObjIDField: mapstr.MapStr{common.BKDBIN: objectIDs}},
			Fields:    []string{common.AssociationKindIDField, common.BKObjIDField, common.BKAsstObjIDField}})

	if err != nil {
		blog.Errorf("search object association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.New(common.CCErrCommHTTPDoRequestFailed, err.Error())
	}

	if err = asstItems.CCError(); err != nil {
		blog.Errorf("failed to search the objects(%#v) association info , err: %v, rid: %s",
			objectIDs, err, kit.Rid)
		return nil, err
	}

	if len(asstItems.Data.Info) == 0 {
		return []metadata.ObjectTopo{}, nil
	}

	var asstKinds, asstObjIDs []string
	assocObjsMap := mapstr.New()
	for index := range asstItems.Data.Info {
		asstKinds = append(asstKinds, asstItems.Data.Info[index].AsstKindID)
		asstObjIDs = append(asstObjIDs, asstItems.Data.Info[index].AsstObjID)
		assocObjsMap.Set(asstItems.Data.Info[index].AsstObjID, asstItems.Data.Info[index].ObjectID)
	}

	// search association type
	assocType, err := o.clientSet.CoreService().Association().ReadAssociationType(kit.Ctx, kit.Header,
		&metadata.QueryCondition{
			Condition: mapstr.MapStr{common.AssociationKindIDField: mapstr.MapStr{common.BKDBIN: asstKinds}},
			Fields:    []string{common.AssociationKindNameField, common.AssociationKindIDField},
		})

	if err != nil {
		blog.Errorf("get association kind[%#v] failed, err: %v, rid: %s", asstKinds, err, kit.Rid)
		return nil, err
	}

	if err = assocType.CCError(); err != nil {
		blog.Errorf("get association kind[%#v] failed, err: %v, rid: %s", asstKinds, err, kit.Rid)
		return nil, err
	}

	if len(assocType.Data.Info) == 0 {
		blog.Errorf("get association kind[%#v] failed, err: can not find this association kind., rid: %s",
			asstKinds, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorTopoAsstKindIsNotExist)
	}

	assocTypeMap := make(map[string]*metadata.AssociationKind, 0)
	for index := range assocType.Data.Info {
		assocTypeMap[assocType.Data.Info[index].AssociationKindID] = assocType.Data.Info[index]
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
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = asstObjs.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s",
			cond, err, kit.Rid)
		return nil, err
	}

	asstObjMap := make(map[string]metadata.Object, 0)
	for index := range asstObjs.Data.Info {
		asstObjMap[asstObjs.Data.Info[index].ObjectID] = asstObjs.Data.Info[index]
	}

	// search direction of association
	isFromMap, err := o.isFrom(kit, asstObjIDs, assocObjsMap)
	if err != nil {
		blog.Errorf("check direction of association failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	results := make([]metadata.ObjectTopo, 0)
	for _, assoc := range asstItems.Data.Info {
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

		if isFromMap[objMap[assoc.ObjectID].ObjectID] {
			tmp.Arrows = "to,from"
		} else {
			tmp.Arrows = "to"
		}

		results = append(results, tmp)
	}

	return results, nil
}

func (o *object) FindObject(kit *rest.Kit, cond mapstr.MapStr) ([]metadata.Object, error) {
	rsp, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: cond},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %v, rid: %s", cond, err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the objects by the condition(%#v) , err: %v, rid: %s",
			cond, err, kit.Rid)
		return nil, err
	}

	return rsp.Data.Info, nil
}

func (o *object) UpdateObject(kit *rest.Kit, data mapstr.MapStr, id int64) error {

	obj, err := o.IsValid(kit, true, data)
	if err != nil {
		blog.Errorf("valid data failed, data: %#v, err: %v, rid: %s", data, err, kit.Rid)
		return err
	}

	obj.ID = id

	if len(obj.ObjectName) != 0 {
		result, err := o.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
			&metadata.QueryCondition{
				Condition: mapstr.MapStr{common.BKObjNameField: obj.ObjectName},
				Fields:    []string{metadata.ModelFieldID},
			})

		if err != nil {
			blog.Errorf("search object by bk_obj_name(%s) failed, err: %v, rid: %s",
				obj.ObjectName, err, kit.Rid)
			return err
		}

		if err = result.CCError(); err != nil {
			blog.Errorf("search object by bk_obj_name(%s) failed, err: %v, rid: %s",
				obj.ObjectName, err, kit.Rid)
			return err
		}

		if len(result.Data.Info) > 1 {
			blog.Errorf("bk_obj_name: %s exist, and get duplicate object, rid: %s", obj.ObjectName, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKObjNameField)
		}

		for _, item := range result.Data.Info {
			if item.ID != obj.ID {
				blog.Errorf("bk_obj_name: %s exist, rid: %s", obj.ObjectName, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommDuplicateItem, common.BKObjNameField)
			}
		}
	}

	// remove unchangeable fields.
	data.Remove(metadata.ModelFieldObjectID)
	data.Remove(metadata.ModelFieldID)
	data.Remove(metadata.ModelFieldObjCls)

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

	rsp, err := o.clientSet.CoreService().Model().UpdateModel(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("update object failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if err = rsp.CCError(); err != nil {
		blog.Errorf("failed to search the object(%s), err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
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

func (o *object) IsValid(kit *rest.Kit, isUpdate bool, data mapstr.MapStr) (*metadata.Object, error) {

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

	return obj, nil
}
