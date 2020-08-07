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
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

type rowInfo struct {
	Row  int64  `json:"row"`
	Info string `json:"info"`
	// value can empty, eg:parse error
	PropertyID string `json:"bk_property_id"`
}

// ObjectOperationInterface object operation methods
type ObjectOperationInterface interface {
	CreateObjectBatch(kit *rest.Kit, data mapstr.MapStr, metadata *metadata.Metadata) (mapstr.MapStr, error)
	FindObjectBatch(kit *rest.Kit, objIDs []string, metaData *metadata.Metadata) (mapstr.MapStr, error)
	CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr, metadata *metadata.Metadata) (model.Object, error)
	CanDelete(kit *rest.Kit, targetObj model.Object) error
	DeleteObject(kit *rest.Kit, id int64, needCheckInst bool, metaData *metadata.Metadata) error
	FindObject(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.Object, error)
	FindObjectTopo(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]metadata.ObjectTopo, error)
	// Deprecated: not allowed to use anymore.
	FindSingleObject(kit *rest.Kit, objectID string, metaData *metadata.Metadata) (model.Object, error)
	FindObjectWithID(kit *rest.Kit, object string, objectID int64, metaData *metadata.Metadata) (model.Object, error)
	UpdateObject(kit *rest.Kit, data mapstr.MapStr, id int64) error

	SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface, grp GroupOperationInterface, unique UniqueOperationInterface)
	IsValidObject(kit *rest.Kit, objID string, metaData *metadata.Metadata) error
}

// NewObjectOperation create a new object operation instance
func NewObjectOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager) ObjectOperationInterface {
	return &object{
		clientSet:   client,
		authManager: authManager,
	}
}

type object struct {
	clientSet    apimachinery.ClientSetInterface
	authManager  *extensions.AuthManager
	modelFactory model.Factory
	instFactory  inst.Factory
	cls          ClassificationOperationInterface
	grp          GroupOperationInterface
	unique       UniqueOperationInterface
	asst         AssociationOperationInterface
	inst         InstOperationInterface
	attr         AttributeOperationInterface
}

func (o *object) SetProxy(modelFactory model.Factory, instFactory inst.Factory, cls ClassificationOperationInterface, asst AssociationOperationInterface, inst InstOperationInterface, attr AttributeOperationInterface, grp GroupOperationInterface, unique UniqueOperationInterface) {
	o.modelFactory = modelFactory
	o.instFactory = instFactory
	o.asst = asst
	o.inst = inst
	o.attr = attr
	o.grp = grp
	o.unique = unique
}

// IsValidObject check whether objID is a real model's bk_obj_id field in backend
func (o *object) IsValidObject(kit *rest.Kit, objID string, metaData *metadata.Metadata) error {

	checkObjCond := condition.CreateCondition()
	checkObjCond.Field(metadata.AttributeFieldObjectID).Eq(objID)

	objItems, err := o.FindObject(kit, checkObjCond, metaData)
	if nil != err {
		blog.Errorf("failed to check the object repeated, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, err.Error())
	}

	if 0 == len(objItems) {
		return kit.CCError.New(common.CCErrCommParamsIsInvalid, fmt.Sprintf("the object id  '%s' is invalid", objID))
	}

	return nil
}

// CreateObjectBatch manipulate model in cc_ObjDes
// this method does'nt act as it's name, it create or update model's attributes indeed.
// it only operate on model already exist, that it to say no new model will be created.
func (o *object) CreateObjectBatch(kit *rest.Kit, data mapstr.MapStr, metaData *metadata.Metadata) (mapstr.MapStr, error) {
	inputDataMap := map[string]ImportObjectData{}
	if err := data.MarshalJSONInto(&inputDataMap); nil != err {
		return nil, err
	}

	result := mapstr.New()
	hasError := false
	for objID, inputData := range inputDataMap {
		subResult := mapstr.New()
		if err := o.IsValidObject(kit, objID, metaData); nil != err {
			blog.Errorf("create model patch, but not a valid model id, model id: %s, rid: %s", objID, kit.Rid)
			subResult["error"] = fmt.Sprintf("the model(%s) is invalid", objID)
			result[objID] = subResult
			hasError = true
			continue
		}

		// 异常错误，比如接卸该行数据失败， 查询数据失败
		var itemErr []rowInfo
		// 新加属性失败
		var addErr []rowInfo
		// 更新属性失败
		var setErr []rowInfo
		// 成功数据的信息
		var succInfo []rowInfo

		// update the object's attribute
		for idx, attr := range inputData.Attr {

			metaAttr := metadata.Attribute{}
			targetAttr, err := metaAttr.Parse(attr)
			if nil != err {
				blog.Errorf("create object batch, but got invalid object attribute, object id: %s, rid: %s", objID, kit.Rid)
				itemErr = append(itemErr, rowInfo{Row: idx, Info: err.Error()})
				hasError = true
				continue
			}
			targetAttr.OwnerID = kit.SupplierAccount
			targetAttr.ObjectID = objID
			if metaData != nil {
				targetAttr.Metadata = *metaData
			}

			if targetAttr.PropertyID == common.BKChildStr || targetAttr.PropertyID == common.BKInstParentStr {
				continue
			}

			if 0 == len(targetAttr.PropertyGroupName) {
				targetAttr.PropertyGroup = "Default"
			}

			// find group
			grpCond := condition.CreateCondition()
			grpCond.Field(metadata.GroupFieldObjectID).Eq(objID)
			grpCond.Field(metadata.GroupFieldGroupName).Eq(targetAttr.PropertyGroupName)
			grps, err := o.grp.FindObjectGroup(kit, grpCond, metaData)
			if nil != err {
				blog.Errorf("create object patch, but find object group failed, object id: %s, group: %s, rid: %s", objID, targetAttr.PropertyGroupName, kit.Rid)
				itemErr = append(itemErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
				hasError = true
				continue
			}

			if 0 != len(grps) {
				targetAttr.PropertyGroup = grps[0].Group().GroupID // should be only one group
			} else {
				newGrp := o.modelFactory.CreateGroup(kit, metaData)
				g := metadata.Group{
					GroupName: targetAttr.PropertyGroupName,
					GroupID:   model.NewGroupID(false),
					ObjectID:  objID,
					OwnerID:   kit.SupplierAccount,
				}
				if metaData != nil {
					g.Metadata = *metaData
				}
				newGrp.SetGroup(g)
				err := newGrp.Save(nil)
				if nil != err {
					setErr = append(setErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
					hasError = true
					continue
				}

				targetAttr.PropertyGroup = newGrp.Group().GroupID
			}

			// create or update the attribute
			attrID, err := attr.String(metadata.AttributeFieldPropertyID)
			if nil != err {
				addErr = append(addErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
				hasError = true
				continue
			}
			attrCond := condition.CreateCondition()
			attrCond.Field(metadata.AttributeFieldObjectID).Eq(objID)
			attrCond.Field(metadata.AttributeFieldPropertyID).Eq(attrID)
			attrs, err := o.attr.FindObjectAttribute(kit, attrCond, metaData)
			if nil != err {
				addErr = append(addErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
				hasError = true
				continue
			}

			if 0 == len(attrs) {

				newAttr := o.modelFactory.CreateAttribute(kit)
				if err = newAttr.Save(targetAttr.ToMapStr()); nil != err {
					addErr = append(addErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
					hasError = true
					continue
				}

			}

			for _, newAttr := range attrs {
				if err := newAttr.Update(targetAttr.ToMapStr()); nil != err {
					setErr = append(setErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
					hasError = true
					continue
				}

			}

			succInfo = append(succInfo, rowInfo{Row: idx, Info: "", PropertyID: targetAttr.PropertyID})
		}

		// 将需要返回的信息更新到result中。 这个函数会修改result参数的值
		o.setCreateObjectBatchObjResult(kit, objID, result, itemErr, addErr, setErr, succInfo)
	}

	if hasError {
		return result, kit.CCError.Error(common.CCErrCommNotAllSuccess)
	}
	return result, nil
}

// setCreateObjectBatchObjResult
func (o *object) setCreateObjectBatchObjResult(kit *rest.Kit, objID string, result mapstr.MapStr, itemErr, addErr, setErr, succInfo []rowInfo) {
	subResult := mapstr.New()
	if len(itemErr) > 0 {
		subResult["errors"] = itemErr
	}
	if len(addErr) > 0 {
		subResult["insert_failed"] = addErr
	}
	if len(setErr) > 0 {
		subResult["update_failed"] = setErr
	}
	if len(succInfo) > 0 {
		subResult["success"] = succInfo
	}
	if len(subResult) > 0 {
		result[objID] = subResult
	}
}

func (o *object) FindObjectBatch(kit *rest.Kit, objIDs []string, metaData *metadata.Metadata) (mapstr.MapStr, error) {
	result := mapstr.New()

	for _, objID := range objIDs {
		obj, err := o.FindSingleObject(kit, objID, metaData)
		if nil != err {
			return nil, err
		}

		attrs, err := obj.GetNonInnerAttributes()
		if nil != err {
			return nil, err
		}

		for _, attr := range attrs {
			attribute := attr.Attribute()
			grpCond := condition.CreateCondition()
			grpCond.Field(metadata.GroupFieldGroupID).Eq(attribute.PropertyGroup)
			grpCond.Field(metadata.GroupFieldObjectID).Eq(attribute.ObjectID)
			grps, err := o.grp.FindObjectGroup(kit, grpCond, metaData)
			if nil != err {
				return nil, err
			}

			for _, grp := range grps {
				// should be only one
				attribute.PropertyGroupName = grp.Group().GroupName
			}
		}

		result.Set(objID, mapstr.MapStr{
			"attr": attrs,
		})
	}

	return result, nil
}

// Deprecated:
// It's not allowed to find a object with only bk_obj_id field, because it
// may find the wrong object, as the public object id is unique, but the private business
// object id can be same across the different business.
func (o *object) FindSingleObject(kit *rest.Kit, objectID string, metaData *metadata.Metadata) (model.Object, error) {

	cond := condition.CreateCondition()
	cond.Field(common.BKObjIDField).Eq(objectID)

	objs, err := o.FindObject(kit, cond, metaData)
	if nil != err {
		blog.Errorf("get model failed, failed to get model by supplier account(%s) objects(%s), err: %s, rid: %s", kit.SupplierAccount, objectID, err.Error(), kit.Rid)
		return nil, err
	}

	if len(objs) == 0 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) not found, result: %+v, rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Error(common.CCErrCommNotFound).Error())
	}

	if len(objs) > 1 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) get multiple, result: %+v, rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Error(common.CCErrCommGetMultipleObject).Error())
	}

	objects := make([]metadata.Object, 0)
	for _, obj := range objs {
		objects = append(objects, obj.Object())
	}

	for _, item := range objs {
		return item, nil
	}
	return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, objectID).Error())
}

func (o *object) FindObjectWithID(kit *rest.Kit, object string, objectID int64, metaData *metadata.Metadata) (model.Object, error) {
	cond := condition.CreateCondition()
	cond.Field("id").Eq(objectID)

	objs, err := o.FindObject(kit, cond, metaData)
	if nil != err {
		blog.Errorf("get model failed, failed to get model by supplier account(%s) objects(%s), err: %s, rid: %s", kit.SupplierAccount, objectID, err.Error(), kit.Rid)
		return nil, err
	}

	if len(objs) == 0 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%d) not found, result: %+v, rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Error(common.CCErrCommNotFound).Error())
	}

	if len(objs) > 1 {
		blog.Errorf("get model failed, get model by supplier account(%s) objects(%s) get multiple, result: %+v, rid: %s", kit.SupplierAccount, objectID, objs, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Error(common.CCErrCommGetMultipleObject).Error())
	}

	objects := make([]metadata.Object, 0)
	for _, obj := range objs {
		objects = append(objects, obj.Object())
	}

	for _, item := range objs {
		return item, nil
	}
	return nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, objectID).Error())
}

func (o *object) CreateObject(kit *rest.Kit, isMainline bool, data mapstr.MapStr, metaData *metadata.Metadata) (model.Object, error) {
	obj := o.modelFactory.CreateObject(kit)
	err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), err: %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	// check the classification
	_, err = obj.GetClassification()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	// check repeated
	exists, err := obj.IsExists()
	if nil != err {
		blog.Errorf("[operation-obj] failed to create the object(%#v), err: %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectCreateFailed, err.Error())
	}

	if exists {
		blog.Errorf("[operation-obj] the object(%#v) is repeated, rid: %s", data, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommDuplicateItem, obj.Object().ObjectID+"/"+obj.Object().ObjectName)
	}

	err = obj.Create()
	if nil != err {
		blog.Errorf("[operation-obj] failed to save the data(%#v), err: %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	object := obj.Object()
	// create the default group
	grp := obj.CreateGroup(metaData)
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    model.NewGroupID(true),
		ObjectID:   object.ObjectID,
		OwnerID:    object.OwnerID,
	}
	if nil != metaData {
		groupData.Metadata = *metaData
	}
	grp.SetGroup(groupData)

	if err = grp.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default group, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrTopoObjectGroupCreateFailed)
	}

	keys := make([]metadata.UniqueKey, 0)
	// create the default inst name
	group := grp.Group()
	attr := obj.CreateAttribute()
	attr.SetAttribute(metadata.Attribute{
		ObjectID:          object.ObjectID,
		IsOnly:            true,
		IsPre:             true,
		Creator:           "user",
		IsEditable:        true,
		PropertyIndex:     -1,
		PropertyGroup:     group.GroupID,
		PropertyGroupName: group.GroupName,
		IsRequired:        true,
		PropertyType:      common.FieldTypeSingleChar,
		PropertyID:        obj.GetInstNameFieldName(),
		PropertyName:      obj.GetDefaultInstPropertyName(),
	})
	if nil != metaData {
		attr.Attribute().Metadata = *metaData
	}
	if err = attr.Create(); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(attr.Attribute().ID)})

	if isMainline {
		pAttr := obj.CreateAttribute()
		pAttr.SetAttribute(metadata.Attribute{
			ObjectID:          object.ObjectID,
			IsOnly:            true,
			IsPre:             true,
			Creator:           "system",
			IsEditable:        true,
			PropertyIndex:     -1,
			PropertyGroup:     group.GroupID,
			PropertyGroupName: group.GroupName,
			IsRequired:        true,
			PropertyType:      common.FieldTypeInt,
			PropertyID:        common.BKInstParentStr,
			PropertyName:      common.BKInstParentStr,
			IsSystem:          true,
		})

		if err = pAttr.Create(); nil != err {
			blog.Errorf("[operation-obj] failed to create the default inst name field, err: %s, rid: %s", err.Error(), kit.Rid)
			return nil, kit.CCError.Error(common.CCErrTopoObjectAttributeCreateFailed)
		}
		keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(pAttr.Attribute().ID)})
	}

	uni := obj.CreateUnique()
	uni.SetKeys(keys)
	uni.SetIsPre(false)
	uni.SetMustCheck(true)
	if err = uni.Save(nil); nil != err {
		blog.Errorf("[operation-obj] failed to create the default inst name field, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	// auth: register model to iam
	if err := o.authManager.RegisterObject(kit.Ctx, kit.Header, object); err != nil {
		return nil, kit.CCError.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
	}

	//package audit response
	err = NewObjectAudit(kit, o.clientSet, obj.Object().ID).buildSnapshotForCur().SaveAuditLog(metadata.AuditCreate)
	if err != nil {
		blog.Errorf("update object %s success, but update to auditLog failed, err: %v, rid: %s", object.ObjectName, err, kit.Rid)
		return nil, err
	}

	return obj, nil
}

// CanDelete return nil only when:
// 1. not inner model
// 2. model has no instances
// 3. model has no association with other model
func (o *object) CanDelete(kit *rest.Kit, targetObj model.Object) error {
	// step 1. ensure not inner model
	if common.IsInnerModel(targetObj.GetObjectID()) {
		return kit.CCError.Error(common.CCErrTopoForbiddenToDeleteModelFailed)
	}

	tObject := targetObj.Object()
	cond := condition.CreateCondition()
	if targetObj.IsCommon() {
		cond.Field(common.BKObjIDField).Eq(tObject.ObjectID)
	}

	// step 2. ensure model has no instances
	query := &metadata.QueryInput{}
	query.Condition = cond.ToMapStr()
	findInstResponse, err := o.inst.FindOriginInst(kit, targetObj.GetObjectID(), query)
	if nil != err {
		blog.Errorf("[operation-obj] failed to check if it (%s) has some insts, err: %s, rid: %s", tObject.ObjectID, err.Error(), kit.Rid)
		return err
	}
	if 0 != findInstResponse.Count {
		blog.Errorf("the object [%s] has been instantiated and cannot be deleted, rid: %s", tObject.ObjectID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrTopoObjectHasSomeInstsForbiddenToDelete, tObject.ObjectID)
	}

	// step 3. ensure model has no association with other model
	or := make([]interface{}, 0)
	or = append(or, mapstr.MapStr{common.BKObjIDField: tObject.ObjectID})
	or = append(or, mapstr.MapStr{common.AssociatedObjectIDField: tObject.ObjectID})

	cond = condition.CreateCondition()
	cond.NewOR().Array(or)

	assocResult, err := o.asst.SearchObject(kit, &metadata.SearchAssociationObjectRequest{Condition: cond.ToMapStr()})
	if err != nil {
		blog.Errorf("check object[%s] can be deleted, but get object associate info failed, err: %v, rid: %s", tObject.ObjectID, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !assocResult.Result {
		blog.Errorf("check if object[%s] can be deleted, but get object associate info failed, err: %v, rid: %s", tObject.ObjectID, err, kit.Rid)
		return kit.CCError.Error(assocResult.Code)
	}

	if len(assocResult.Data) != 0 {
		blog.Errorf("check if object[%s] can be deleted, but object has already associate to another one., rid: %s", tObject.ObjectID, kit.Rid)
		return kit.CCError.Error(common.CCErrorTopoObjectHasAlreadyAssociated)
	}

	return nil
}

// DeleteObject delete model by id
func (o *object) DeleteObject(kit *rest.Kit, id int64, needCheckInst bool, metaData *metadata.Metadata) error {
	if id <= 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKFieldID)
	}

	// get model by id
	cond := condition.CreateCondition()
	cond.Field(metadata.ModelFieldID).Eq(id)
	objs, err := o.FindObject(kit, cond, metaData)
	if nil != err {
		blog.Errorf("[operation-obj] failed to find objects, the condition is (%v) err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return err
	}
	// shouldn't return 404 error here, legacy implements just ignore not found error
	if len(objs) == 0 {
		blog.V(3).Infof("[operation-obj] object not found, condition: %v, err: %s, rid: %s", cond, err.Error(), kit.Rid)
		return nil
	}
	if len(objs) > 1 {
		return kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	obj := objs[0]
	object := obj.Object()

	// check whether it can be deleted
	if needCheckInst {
		if err := o.CanDelete(kit, obj); nil != err {
			return err
		}
	}

	//get PreData
	objAudit := NewObjectAudit(kit, o.clientSet, id).buildSnapshotForPre()

	// DeleteModelCascade 将会删除模型/模型属性/属性分组/唯一校验
	rsp, err := o.clientSet.CoreService().Model().DeleteModelCascade(kit.Ctx, kit.Header, id)
	if nil != err {
		blog.Errorf("[operation-obj] failed to request the object controller, err: %s, rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("[operation-obj] failed to delete the object by the id(%d), rid: %s", id, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	// auth: deregister models
	if err := o.authManager.DeregisterObject(kit.Ctx, kit.Header, object); err != nil {
		blog.ErrorJSON("Delete Object success, but deregister object from iam failed, objects: %s, err: %s, rid: %s", object, err.Error(), kit.Rid)
		return kit.CCError.New(common.CCErrCommUnRegistResourceToIAMFailed, err.Error())
	}

	//saveAuditLog
	err = objAudit.SaveAuditLog(metadata.AuditDelete)
	if err != nil {
		blog.Errorf("Delete object %s success, but update to auditLog failed, err: %v, rid: %s", object.ObjectName, err, kit.Rid)
		return err
	}
	return nil
}

func (o *object) isFrom(kit *rest.Kit, fromObjID, toObjID string, metaData *metadata.Metadata) (bool, error) {
	asstItems, err := o.asst.SearchObjectAssociation(kit, fromObjID, metaData)
	if nil != err {
		return false, err
	}

	for _, asst := range asstItems {
		if asst.AsstObjID == toObjID {
			return true, nil
		}
	}

	return false, nil
}

func (o *object) FindObjectTopo(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]metadata.ObjectTopo, error) {
	objs, err := o.FindObject(kit, cond, metaData)
	if nil != err {
		blog.Errorf("[operation-obj] failed to find object, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	results := make([]metadata.ObjectTopo, 0)
	for _, obj := range objs {
		object := obj.Object()
		asstItems, err := o.asst.SearchObjectAssociation(kit, object.ObjectID, metaData)
		if nil != err {
			return nil, err
		}

		for _, asst := range asstItems {

			// find association kind with association kind id.
			typeCond := condition.CreateCondition()
			typeCond.Field(common.AssociationKindIDField).Eq(asst.AsstKindID)
			request := &metadata.SearchAssociationTypeRequest{
				Condition: typeCond.ToMapStr(),
			}

			resp, err := o.asst.SearchType(kit, request)
			if err != nil {
				blog.Errorf("find object topo failed, because get association kind[%s] failed, err: %v, rid: %s", asst.AsstKindID, err, kit.Rid)
				return nil, kit.CCError.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
			}
			if !resp.Result {
				blog.Errorf("find object topo failed, because get association kind[%s] failed, err: %v, rid: %s", asst.AsstKindID, resp.ErrMsg, kit.Rid)
				return nil, kit.CCError.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
			}

			// should only be one association kind.
			if len(resp.Data.Info) == 0 {
				blog.Errorf("find object topo failed, because get association kind[%s] failed, err: can not find this association kind., rid: %s", asst.AsstKindID, kit.Rid)
				return nil, kit.CCError.Errorf(common.CCErrTopoGetAssociationKindFailed, asst.AsstKindID)
			}

			cond = condition.CreateCondition()
			cond.Field(common.BKObjIDField).Eq(asst.AsstObjID)

			asstObjs, err := o.FindObject(kit, cond, metaData)
			if nil != err {
				blog.Errorf("[operation-obj] failed to find object, err: %s, rid: %s", err.Error(), kit.Rid)
				return nil, err
			}

			for _, asstObj := range asstObjs {
				assocObject := asstObj.Object()
				tmp := metadata.ObjectTopo{}
				tmp.Label = resp.Data.Info[0].AssociationKindName
				tmp.LabelName = resp.Data.Info[0].AssociationKindName
				tmp.From.ObjID = object.ObjectID
				cls, err := obj.GetClassification()
				if nil != err {
					return nil, err
				}
				tmp.From.ClassificationID = cls.Classify().ClassificationID
				tmp.From.Position = object.Position
				tmp.From.OwnerID = object.OwnerID
				tmp.From.ObjName = object.ObjectName
				tmp.To.OwnerID = assocObject.OwnerID
				tmp.To.ObjID = assocObject.ObjectID

				cls, err = asstObj.GetClassification()
				if nil != err {
					return nil, err
				}
				tmp.To.ClassificationID = cls.Classify().ClassificationID
				tmp.To.Position = assocObject.Position
				tmp.To.ObjName = assocObject.ObjectName
				ok, err := o.isFrom(kit, assocObject.ObjectID, object.ObjectID, metaData)
				if nil != err {
					return nil, err
				}

				if ok {
					tmp.Arrows = "to,from"
				} else {
					tmp.Arrows = "to"
				}

				results = append(results, tmp)
			}
		}

	}

	return results, nil
}

func (o *object) FindObject(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.Object, error) {
	fCond := cond.ToMapStr()
	if nil != metaData {
		// search model from special business
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		// search global shared model
		fCond.Merge(metadata.BizLabelNotExist)
	}
	rsp, err := o.clientSet.CoreService().Model().ReadModel(context.Background(), kit.Header, &metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-obj] find object failed, cond: %+v, err: %s, rid: %s", fCond, err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-obj] failed to search the objects by the condition(%#v) , error info is %s, rid: %s", fCond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	models := make([]metadata.Object, 0)
	for index := range rsp.Data.Info {
		models = append(models, rsp.Data.Info[index].Spec)
	}
	return model.CreateObject(kit, o.clientSet, models), nil
}

func (o *object) UpdateObject(kit *rest.Kit, data mapstr.MapStr, id int64) error {

	obj := o.modelFactory.CreateObject(kit)
	obj.SetRecordID(id)
	err := obj.Parse(data)
	if nil != err {
		blog.Errorf("[operation-obj] failed to parse the data(%#v), err: %s, rid: %s", data, err.Error(), kit.Rid)
		return err
	}

	object := obj.Object()

	/*
		// auth: check authorization
		if err := o.authManager.AuthorizeObjectByRawID(kit.Ctx, kit.Header, meta.Update, object.ObjectID); err != nil {
			blog.V(2).Infof("update model %s failed, authorization failed, err: %+v, rid: %s", object.ObjectID, err, kit.Rid)
			return err
		}
	*/

	//get PreData
	objAudit := NewObjectAudit(kit, o.clientSet, id).buildSnapshotForPre()

	// check repeated
	exists, err := obj.IsExists()
	if nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%#v), err: %s, rid: %s", data, err.Error(), kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	if exists {
		blog.Errorf("[operation-obj] the object(%#v) is repeated, rid: %s", data, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommDuplicateItem, obj.Object().ObjectName)
	}
	if err = obj.Update(data); nil != err {
		blog.Errorf("[operation-obj] failed to update the object(%d), the new data(%#v), err: %s, rid: %s", id, data, err.Error(), kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	bizID, err := metadata.BizIDFromMetadata(object.Metadata)
	if err != nil {
		blog.Error("update object: %s, but parse business id failed, err: %v, rid: %s", object.ObjectID, err, kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectUpdateFailed, err.Error())
	}

	// auth update register info
	if err := o.authManager.UpdateRegisteredObjectsByRawIDs(kit.Ctx, kit.Header, bizID, id); err != nil {
		blog.Errorf("update object %s success, but update to auth failed, err: %v, rid: %s", object.ObjectName, err, kit.Rid)
		return kit.CCError.New(common.CCErrCommRegistResourceToIAMFailed, err.Error())
	}

	//get CurData and saveAuditLog
	err = objAudit.buildSnapshotForCur().SaveAuditLog(metadata.AuditUpdate)
	if err != nil {
		blog.Errorf("update object %s success, but update to auditLog failed, err: %v, rid: %s", object.ObjectName, err, kit.Rid)
		return err
	}
	return nil
}
