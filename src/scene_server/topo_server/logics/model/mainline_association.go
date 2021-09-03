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
	"regexp"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateMainlineAssociation create mainline object association
func (assoc *association) CreateMainlineAssociation(kit *rest.Kit, data *metadata.MainlineAssociation,
	maxTopoLevel int) (*metadata.Object, error) {

	if data.AsstObjID == "" {
		blog.Errorf("bk_asst_obj_id empty, input: %#v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAsstObjIDField)
	}

	if data.ClassificationID == "" {
		blog.Errorf("bk_classification_id empty, input: %#v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKClassificationIDField)
	}

	mlCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
	}
	mainlineAsst, ccErr := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mlCond)
	if ccErr != nil {
		blog.Errorf("failed to check the mainline topo level, err: %v, rid: %s", ccErr, kit.Rid)
		return nil, ccErr
	}

	// 总层数等于关联关系数加1，通过count查出的数量与实际主线模型数量差1个
	if len(mainlineAsst.Info)+1 >= maxTopoLevel {
		blog.Errorf("the mainline topo level is %d, the max limit is %d, rid: %s", len(mainlineAsst.Info)+1,
			maxTopoLevel, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoBizTopoLevelOverLimit)
	}

	objData := mapstr.MapStr{
		common.BKObjIDField:            data.ObjectID,
		common.BKObjNameField:          data.ObjectName,
		common.BKObjIconField:          data.ObjectIcon,
		common.BKClassificationIDField: data.ClassificationID,
		common.BkSupplierAccount:       data.OwnerID,
	}

	// find the mainline parent object
	parentObjID := data.AsstObjID
	var childObjID string
	for _, item := range mainlineAsst.Info {
		if item.AsstObjID == parentObjID {
			childObjID = item.ObjectID
			break
		}
	}

	if len(childObjID) == 0 {
		blog.Errorf("object(%s) has not got any mainline child object, rid: %s", parentObjID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, parentObjID)
	}

	currentObj, err := assoc.CreateObject(kit, true, objData)
	if err != nil {
		return nil, err
	}

	// update the mainline topo inst association
	if _, err = assoc.SetMainlineInstAssociation(kit, parentObjID, childObjID, currentObj.ObjectID,
		currentObj.ObjectName); err != nil {
		blog.Errorf("failed set the mainline inst association, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err = assoc.createMainlineObjectAssociation(kit, currentObj.ObjectID, parentObjID); err != nil {
		blog.Errorf("create mainline object[%s] association related to object[%s] failed, err: %v, rid: %s",
			currentObj.ObjectID, parentObjID, err, kit.Rid)
		return nil, err
	}

	if err = assoc.setMainlineParentObject(kit, childObjID, currentObj.ObjectID); err != nil {
		blog.Errorf("update mainline current object's[%s] child object[%s] association to current failed, "+
			"err: %v, rid: %s", currentObj.ObjectID, childObjID, err, kit.Rid)
		return nil, err
	}

	return currentObj, nil
}

// DeleteMainlineAssociation delete mainline association by objID
func (assoc *association) DeleteMainlineAssociation(kit *rest.Kit, targetObjID string) error {

	if common.IsInnerModel(targetObjID) {
		return kit.CCError.Errorf(common.CCErrCommForbiddenOperateInnerModelInstanceWithCommonAPI)
	}

	mainlineCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: targetObjID},
				{common.BKAsstObjIDField: targetObjID},
			}},
	}
	rsp, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mainlineCond)
	if err != nil {
		blog.Errorf("search mainline association failed, cond: %#v, err: %v, rid: %s", mainlineCond, err, kit.Rid)
		return err
	}

	var parentObjID, childObjID string
	for _, assoc := range rsp.Info {
		// delete this object related association.
		// a pre-defined association can not be updated.
		if assoc.IsPre != nil && *assoc.IsPre {
			blog.Errorf("it's a pre-defined association, can not be deleted, cond: %#v, rid: %s", mainlineCond, kit.Rid)
			return kit.CCError.CCError(common.CCErrorTopoDeletePredefinedAssociation)
		}

		if assoc.AsstObjID == targetObjID {
			childObjID = assoc.ObjectID
		}

		if assoc.ObjectID == targetObjID {
			parentObjID = assoc.AsstObjID
		}
	}

	if len(childObjID) == 0 || len(parentObjID) == 0 {
		blog.Errorf("target object(%s) has no parent or child, rid: %s", targetObjID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, targetObjID)
	}

	if err = assoc.ResetMainlineInstAssociation(kit, targetObjID, childObjID); err != nil {
		blog.Errorf("failed to delete the object(%s)'s instance, err: %v, rid: %s", targetObjID, err, kit.Rid)
		return err
	}

	if err = assoc.createMainlineObjectAssociation(kit, childObjID, parentObjID); err != nil {
		blog.Errorf("failed to update the association, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// delete the object association
	_, err = assoc.clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header,
		&metadata.DeleteOption{Condition: mainlineCond.Condition})
	if err != nil {
		blog.Errorf("delete object association failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// object的删除函数通过object的id进行删除，需要在这里查一次object的id
	objIDCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKObjIDField: targetObjID},
		Fields:         []string{common.BKFieldID},
		DisableCounter: true,
	}
	//delete objects
	if err = assoc.obj.DeleteObject(kit, objIDCond, false); err != nil {
		blog.Errorf("failed to delete the object(%s), err: %v, rid: %s", targetObjID, err, kit.Rid)
		return err
	}

	return nil
}

// SearchMainlineAssociationTopo get mainline topo of special model
// result is a list with targetObj as head, so if you want a full topo, target must be biz model.
func (assoc *association) SearchMainlineAssociationTopo(kit *rest.Kit,
	targetObjID string) ([]*metadata.MainlineObjectTopo, error) {

	if len(targetObjID) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommInstDataNil, common.BKObjIDField)
	}

	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
	}
	mainlineAssoc, err := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search mainline association failed, cond: %#v, err: %v, rid: %s", query, err, kit.Rid)
		return nil, err
	}

	childMap := make(map[string]string)
	parentMap := make(map[string]string)
	needFind := []string{targetObjID}
	for _, asso := range mainlineAssoc.Info {
		childMap[asso.AsstObjID] = asso.ObjectID
		parentMap[asso.ObjectID] = asso.AsstObjID
	}

	// 遍历获取以targetObj为头的topo切片
	cursor := targetObjID
	for {
		child, exist := childMap[cursor]
		if !exist {
			break
		}
		needFind = append(needFind, child)
		cursor = child
	}

	queryCond := &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKObjNameField, common.BkSupplierAccount},
		Condition: mapstr.MapStr{
			common.BKObjIDField: mapstr.MapStr{common.BKDBIN: needFind},
		},
	}
	objMsg, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("search objects[%#v] failed, cond: %#v, err: %v, rid: %s", needFind, queryCond, err, kit.Rid)
		return nil, err
	}

	objMap := map[string]metadata.Object{}
	for _, obj := range objMsg.Info {
		objMap[obj.ObjectID] = obj
	}

	result := make([]*metadata.MainlineObjectTopo, 0)
	for _, objID := range needFind {

		result = append(result, &metadata.MainlineObjectTopo{
			ObjID:      objID,
			ObjName:    objMap[objID].ObjectName,
			OwnerID:    objMap[objID].OwnerID,
			NextObj:    childMap[objID],
			NextName:   objMap[childMap[objID]].ObjectName,
			PreObjID:   parentMap[objID],
			PreObjName: objMap[parentMap[objID]].ObjectName,
		})
	}

	return result, nil
}

// IsMainlineObject check whether objID is mainline object or not
func (assoc *association) IsMainlineObject(kit *rest.Kit, objID string) (bool, error) {
	filter := []map[string]interface{}{
		{common.AssociationKindIDField: common.AssociationKindMainline,
			common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: objID},
				{common.BKAsstObjIDField: objID},
			}},
	}
	asstCnt, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNameObjAsst, filter)
	if err != nil {
		blog.Errorf("check object(%s) if is mainline object failed, err: %v, rid: %s", objID, err, kit.Rid)
		return false, err
	}

	if len(asstCnt) <= 0 {
		blog.Errorf("get association by filter: %#v failed, return is empty, rid: %s", filter, kit.Rid)
		return false, kit.CCError.CCError(common.CCErrorTopoObjectAssociationNotExist)
	}

	if asstCnt[0] == 0 {
		return false, nil
	}

	return true, nil
}

// TODO after merge , delete this func and use SetMainlineInstAssociation in inst/mainline_association
func (assoc *association) SetMainlineInstAssociation(kit *rest.Kit, parentObjID, childObjID, currentObjID,
	currentObjName string) ([]int64, error) {

	defaultCond := &metadata.QueryInput{}
	cond := mapstr.New()
	if metadata.IsCommon(parentObjID) {
		cond.Set(common.BKObjIDField, parentObjID)
	}
	defaultCond.Condition = cond
	// fetch all parent instances.
	// TODO replace to FindInst in inst/inst.go after merge
	parentInsts, err := assoc.FindInst(kit, parentObjID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find parent object(%s) inst, err: %v, rid: %s", parentObjID, err, kit.Rid)
		return nil, err
	}

	createdInstIDs := make([]int64, len(parentInsts.Info))
	expectParent2Children := make(map[int64][]int64)
	// filters out special character for mainline instances
	re, _ := regexp.Compile(`[#/,><|]`)
	instanceName := re.ReplaceAllString(currentObjName, "")
	var parentInstIDs []int64
	// create current object instance for each parent instance and insert the current instance to
	for _, parentInst := range parentInsts.Info {
		id, err := parentInst.Int64(common.GetInstIDField(parentObjID))
		if err != nil {
			blog.Errorf("failed to find the inst id, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
		parentInstIDs = append(parentInstIDs, id)

		// we create the current object's instance for each parent instance belongs to the parent object.
		currentInst := mapstr.MapStr{common.BKObjIDField: currentObjID}
		currentInst.Set(common.GetInstNameField(currentObjID), instanceName)
		// set current instance's parent id to parent instance's id, so that they can be chained.
		currentInst.Set(common.BKInstParentStr, id)

		if parentObjID == common.BKInnerObjIDApp {
			currentInst.Set(common.BKAppIDField, id)
		} else {
			if bizID, ok := parentInst.Get(common.BKAppIDField); ok {
				currentInst.Set(common.BKAppIDField, bizID)
			}
		}

		// create the instance now.
		currentInstID, err := assoc.createInst(kit, currentObjID, currentInst)
		if err != nil {
			blog.Errorf("failed to create object(%s) default inst, err: %v, rid: %s", currentObjID, err, kit.Rid)
			return nil, err
		}

		createdInstIDs = append(createdInstIDs, int64(currentInstID))
	}

	// reset the child's parent instance's parent id to current instance's id.
	childInst, err := assoc.getMainlineChildInst(kit, parentObjID, childObjID, parentInstIDs)
	if err != nil {
		blog.Errorf("failed to get the object(%s) mainline child inst, err: %v, rid: %s",
			parentObjID, err, kit.Rid)
		return nil, err
	}

	for _, child := range childInst {
		childID, err := child.Int64(common.GetInstIDField(childObjID))
		if err != nil {
			blog.Errorf("failed to get the inst id from the inst(%#v), rid: %s", child, kit.Rid)
			continue
		}

		parentID, err := child.Int64(common.BKParentIDField)
		if err != nil {
			blog.Errorf("failed to get the object(%s) mainline parent id, "+
				"the current inst(%v), err: %v, rid: %s", childObjID, child, err, kit.Rid)
			continue
		}
		expectParent2Children[parentID] = append(expectParent2Children[parentID], childID)
	}

	for parentID, childIDs := range expectParent2Children {
		if err = assoc.updateMainlineAssociation(kit, childIDs, childObjID, parentID); err != nil {
			blog.Errorf("failed to set the object mainline child inst, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}
	}

	// create audit log for the created instances.
	audit := auditlog.NewInstanceAudit(assoc.clientSet.CoreService())

	cond = map[string]interface{}{
		metadata.GetInstIDFieldByObjID(currentObjID): map[string]interface{}{common.BKDBIN: createdInstIDs},
	}
	// generate audit log.
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, currentObjID, cond)
	if err != nil {
		blog.Errorf(" creat inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	err = audit.SaveAuditLog(kit, auditLog...)
	if err != nil {
		blog.Errorf("creat inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrAuditSaveLogFailed)
	}

	return createdInstIDs, nil
}

// TODO after merge this func need to be deleted, replace to the func of depend on mainline instance association
func (assoc *association) getMainlineChildInst(kit *rest.Kit, objID, childObjID string, instIDs []int64) (
	[]mapstr.MapStr, error) {

	cond := mapstr.MapStr{common.BKInstParentStr: mapstr.MapStr{common.BKDBIN: instIDs}}
	if metadata.IsCommon(childObjID) {
		cond.Set(metadata.ModelFieldObjectID, childObjID)
	} else if childObjID == common.BKInnerObjIDSet {
		cond.Set(common.BKDefaultField, mapstr.MapStr{common.BKDBNE: common.DefaultResSetFlag})
	}

	instCond := &metadata.QueryInput{
		Condition: cond,
		Fields:    fmt.Sprintf("%s,%s", common.GetInstIDField(childObjID), common.BKParentIDField),
	}
	instRsp, err := assoc.FindInst(kit, childObjID, instCond)
	if err != nil {
		blog.Errorf("search inst by object(%s) failed, err: %v, rid: %s", childObjID, err, kit.Rid)
		return nil, err
	}

	return instRsp.Info, nil
}

// IsValidObject check whether objID is a real model's bk_obj_id field in backend
// TODO this function should be delete.
// TODO every function which use this logic need to replace to IsValidObject in model/object.go.
func (assoc *association) IsValidObject(kit *rest.Kit, objID string) error {

	checkObjCond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}

	objItems, err := assoc.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: checkObjCond},
	)
	if err != nil {
		blog.Errorf("find object failed, cond: %+v, err: %v, rid: %s", checkObjCond, err, kit.Rid)
		return err
	}

	if len(objItems.Info) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	return nil
}

// TODO this function should be delete.
// TODO every function which use this logic need to replace to IsValidObject in model/object.go.
func (assoc *association) CreateObject(kit *rest.Kit, isMainline bool,
	data mapstr.MapStr) (*metadata.Object, error) {

	obj, err := IsValid(kit, false, data)
	if err != nil {
		blog.Errorf("valid data(%#v) failed, err: %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	objCls, err := assoc.clientSet.CoreService().Model().ReadModelClassification(kit.Ctx, kit.Header,
		&metadata.QueryCondition{Condition: mapstr.MapStr{common.BKClassificationIDField: obj.ObjCls}})
	if err != nil {
		blog.Errorf("get object classification failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	if len(objCls.Info) == 0 {
		blog.Errorf("can't find classification by params, classification: %s is not exist, rid: %s",
			obj.ObjCls, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKClassificationIDField)
	}

	filter := mapstr.MapStr{
		common.BKDBOR: []mapstr.MapStr{
			{common.BKObjIDField: obj.ObjectID},
			{common.BKObjNameField: obj.ObjectName},
		}}
	cnt, err := assoc.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjDes,
		[]map[string]interface{}{filter})
	if err != nil {
		blog.Errorf("get object number by filter failed, err: %s, rid: %s", err.Error(), kit.Rid)
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

	objRsp, err := assoc.clientSet.CoreService().Model().CreateModel(kit.Ctx, kit.Header,
		&metadata.CreateModel{Spec: *obj})
	if err != nil {
		blog.Errorf("failed to search the object(%s), err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return nil, err
	}

	obj.ID = int64(objRsp.Created.ID)

	// create the default group
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    "default",
		ObjectID:   obj.ObjectID,
		OwnerID:    obj.OwnerID,
	}

	_, err = assoc.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header,
		obj.ObjectID, metadata.CreateModelAttributeGroup{Data: groupData})
	if err != nil {
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
	rspAttr, err := assoc.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header,
		attr.ObjectID, &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{attr}})
	if err != nil {
		blog.Errorf("create model attrs failed, ObjectID: %s, input: %#v, err: %v, rid: %s",
			attr.ObjectID, attr, err, kit.Rid)
		return nil, err
	}

	for _, exception := range rspAttr.Exceptions {
		return nil, kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(rspAttr.Repeated) > 0 {
		blog.Errorf("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %#v, rid: %s",
			attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	if len(rspAttr.Created) != 1 {
		blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %#v, rid: %s",
			attr.ObjectID, attr, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}

	attr.ID = int64(rspAttr.Created[0].ID)

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

		rsppAttr, err := assoc.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header,
			pAttr.ObjectID, &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{pAttr}})
		if err != nil {
			blog.Errorf("create model attrs failed, ObjectID: %s, input: %#v, rid: %s", pAttr.ObjectID, pAttr, kit.Rid)
			return nil, err
		}

		for _, exception := range rsppAttr.Exceptions {
			return nil, kit.CCError.New(int(exception.Code), exception.Message)
		}

		if len(rsppAttr.Repeated) > 0 {
			blog.Errorf("create model attrs failed, the attr is duplicated, "+
				"ObjectID: %s, input: %#v, rid: %s", pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
		}

		if len(rsppAttr.Created) != 1 {
			blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %#v, rid: %s",
				pAttr.ObjectID, pAttr, kit.Rid)
			return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
		}
		pAttr.ID = int64(rsppAttr.Created[0].ID)

		keys = append(keys, metadata.UniqueKey{Kind: metadata.UniqueKeyKindProperty, ID: uint64(pAttr.ID)})
	}

	uni := metadata.ObjectUnique{
		ObjID:   obj.ObjectID,
		OwnerID: kit.SupplierAccount,
		Keys:    keys,
		Ispre:   false,
	}
	// NOTICE: 2021年03月29日  唯一索引与index.MainLineInstanceUniqueIndex,index.InstanceUniqueIndex定义强依赖
	// 原因：建立模型之前要将表和表中的索引提前建立，mongodb 4.2.6(4.4之前)事务中不能建表
	// 事务操作表中数据操作和建表，建立索引为互斥操作。
	_, err = assoc.clientSet.CoreService().Model().CreateModelAttrUnique(kit.Ctx, kit.Header,
		uni.ObjID, metadata.CreateModelAttrUnique{Data: uni})
	if err != nil {
		blog.Errorf("create unique for %s failed, err: %v, rid: %s", uni.ObjID, err, kit.Rid)
		return nil, err
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewObjectAuditLog(assoc.clientSet.CoreService())
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

// TODO need to be deleted after merge
func IsValid(kit *rest.Kit, isUpdate bool, data mapstr.MapStr) (*metadata.Object, error) {

	obj := new(metadata.Object)
	if err := mapstruct.Decode2Struct(data, obj); err != nil {
		blog.Errorf("parse object failed, err: %v, input: %#v, rid: %s", err, data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed)
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectID) {

		if err := util.ValidModelIDField(data[metadata.ModelFieldObjectID],
			metadata.ModelFieldObjectID, kit.CCError); err != nil {
			blog.Errorf("failed to valid the object id(%s), rid: %s",
				metadata.ModelFieldObjectID, kit.Rid)
			return nil, err
		}
	}

	if !isUpdate || data.Exists(metadata.ModelFieldObjectName) {
		if err := util.ValidModelNameField(data[metadata.ModelFieldObjectName],
			metadata.ModelFieldObjectName, kit.CCError); err != nil {
			blog.Errorf("failed to valid the object name(%s), rid: %s",
				metadata.ModelFieldObjectName, kit.Rid)
			return nil, kit.CCError.New(common.CCErrCommParamsIsInvalid,
				metadata.ModelFieldObjectName+" "+err.Error())
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

// FindInst search instance by condition
// TODO need to delete after merge
func (assoc *association) FindInst(kit *rest.Kit, objID string,
	cond *metadata.QueryInput) (*metadata.InstResult, error) {

	result := new(metadata.InstResult)
	switch objID {
	case common.BKInnerObjIDHost:
		rsp, err := assoc.clientSet.CoreService().Host().GetHosts(kit.Ctx, kit.Header, cond)
		if err != nil {
			blog.Errorf("get host failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Count
		result.Info = rsp.Info
		return result, nil

	default:
		input := &metadata.QueryCondition{Condition: cond.Condition, TimeCondition: cond.TimeCondition}
		input.Page.Start = cond.Start
		input.Page.Limit = cond.Limit
		input.Page.Sort = cond.Sort
		input.Fields = strings.Split(cond.Fields, ",")
		rsp, err := assoc.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, objID, input)
		if err != nil {
			blog.Errorf("search instance failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		result.Count = rsp.Count
		result.Info = rsp.Info
		return result, nil
	}
}

// TODO need to be deleted after merge
func (assoc *association) createInst(kit *rest.Kit, objID string, data mapstr.MapStr) (uint64, error) {
	cond := &metadata.CreateModelInstance{
		Data: data,
	}
	rsp, err := assoc.clientSet.CoreService().Instance().CreateInstance(kit.Ctx, kit.Header, objID, cond)
	if err != nil {
		blog.Errorf("failed to create object instance, err: %v, rid: %s", err, kit.Rid)
		return 0, err
	}

	return rsp.Created.ID, nil
}

// TODO after merge this func need to be deleted, replace to the func of depend on mainline instance association
func (assoc *association) updateMainlineAssociation(kit *rest.Kit, childID []int64, childObjID string,
	parentID int64) error {

	cond := mapstr.MapStr{common.GetInstIDField(childObjID): mapstr.MapStr{common.BKDBIN: childID}}
	if metadata.IsCommon(childObjID) {
		cond.Set(metadata.ModelFieldObjectID, childObjID)
	}

	input := metadata.UpdateOption{
		Data: mapstr.MapStr{
			common.BKInstParentStr: parentID,
		},
		Condition: cond,
	}
	_, err := assoc.clientSet.CoreService().Instance().UpdateInstance(kit.Ctx, kit.Header, childObjID, &input)
	if err != nil {
		blog.Errorf("failed to update instance of object(%s), err: %v, rid: %s", childObjID, err, kit.Rid)
		return err
	}

	return nil
}

func (assoc *association) setMainlineParentObject(kit *rest.Kit, childObjID, parentObjID string) error {
	cond := mapstr.MapStr{
		common.BKObjIDField:           childObjID,
		common.AssociationKindIDField: common.AssociationKindMainline,
	}

	_, err := assoc.clientSet.CoreService().Association().DeleteModelAssociation(kit.Ctx, kit.Header,
		&metadata.DeleteOption{Condition: cond})
	if err != nil {
		blog.Errorf("delete object(%s) association failed, err: %v, rid: %s", childObjID, err, kit.Rid)
		return err
	}

	return assoc.createMainlineObjectAssociation(kit, childObjID, parentObjID)
}

func (assoc *association) createMainlineObjectAssociation(kit *rest.Kit, childObjID, parentObjID string) error {
	objAsstID := fmt.Sprintf("%s_%s_%s", childObjID, common.AssociationKindMainline, parentObjID)
	defined := false
	association := metadata.Association{
		OwnerID:              kit.SupplierAccount,
		AssociationName:      objAsstID,
		AssociationAliasName: objAsstID,
		ObjectID:             childObjID,
		// related to it's parent object id
		AsstObjID:  parentObjID,
		AsstKindID: common.AssociationKindMainline,
		Mapping:    metadata.OneToOneMapping,
		OnDelete:   metadata.NoAction,
		IsPre:      &defined,
	}

	_, err := assoc.clientSet.CoreService().Association().CreateMainlineModelAssociation(kit.Ctx, kit.Header,
		&metadata.CreateModelAssociation{Spec: association})
	if err != nil {
		blog.Errorf("create mainline object association failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// TODO this func need to replace to ResetMainlineInstAssociation in inst/mainline_association
func (assoc *association) ResetMainlineInstAssociation(kit *rest.Kit, currentObjID, childObjID string) error {

	cond := mapstr.New()
	if metadata.IsCommon(currentObjID) {
		cond.Set(common.BKObjIDField, currentObjID)
	}
	defaultCond := &metadata.QueryInput{Condition: cond}

	// 获取 current 模型的所有实例
	currentInsts, err := assoc.FindInst(kit, currentObjID, defaultCond)
	if err != nil {
		blog.Errorf("failed to find current object(%s) inst, err: %v, rid: %s", currentObjID, err, kit.Rid)
		return err
	}

	instIDs := make([]int64, len(currentInsts.Info))
	instParentMap := map[int64]int64{}
	for _, currInst := range currentInsts.Info {
		instID, err := currInst.Int64(common.GetInstIDField(currentObjID))
		if err != nil {
			blog.Errorf("get inst id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		instIDs = append(instIDs, instID)

		currInstParentID, err := currInst.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get bk_parent_id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		instParentMap[instID] = currInstParentID
	}

	children, err := assoc.getMainlineChildInst(kit, currentObjID, childObjID, instIDs)
	if err != nil {
		return err
	}

	// 检查实例删除后，会不会出现重名冲突
	canReset, repeatedInstName, err := assoc.checkInstNameRepeat(kit, instParentMap, children)
	if err != nil {
		blog.Errorf("can not be reset, err: %+v, rid: %s", err, kit.Rid)
		return err
	}

	if canReset == false {
		blog.Errorf("can not be reset, inst name repeated, inst: %s, rid: %s", repeatedInstName, kit.Rid)
		errMsg := kit.CCError.CCError(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat).Error() +
			" " + repeatedInstName
		return kit.CCError.New(common.CCErrTopoDeleteMainLineObjectAndInstNameRepeat, errMsg)
	}

	// NEED FIX: 下面循环中的continue ，会在处理实例异常的时候跳过当前拓扑的处理，此方式可能会导致某个业务拓扑失败，但是不会影响所有。
	// 修改 currentInsts 所有孩子结点的父节点，为 currentInsts 的父节点，并删除 currentInsts

	childIDMap := make(map[int64][]int64)
	for _, child := range children {
		childInstID, err := child.Int64(common.GetInstIDField(currentObjID))
		if err != nil {
			blog.Errorf("get inst id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		childParentID, err := child.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get bk_parent_id in current insts failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		childIDMap[childParentID] = append(childIDMap[childParentID], childInstID)
	}

	parentChildMap := map[int64][]int64{}
	for inst, parent := range instParentMap {
		parentChildMap[parent] = append(parentChildMap[parent], childIDMap[inst]...)
	}

	// set the child's parent
	for parent, child := range parentChildMap {
		if len(child) == 0 {
			continue
		}

		if err = assoc.updateMainlineAssociation(kit, child, childObjID, parent); err != nil {
			blog.Errorf("failed to set the object mainline child inst, parent: %d, child: %v, err: %v, rid: %s",
				parent, child, err, kit.Rid)
			return err
		}
	}

	// delete the current inst
	if err := assoc.DeleteMainlineInstWithID(kit, currentObjID, instIDs); err != nil {
		blog.Errorf("failed to delete the current inst, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// checkInstNameRepeat 检查如果将 currentInsts 都删除之后，拥有共同父节点的孩子结点会不会出现名字冲突
// 如果有冲突，返回 (false, 冲突实例名, nil)
// TODO after merge this func need to be deleted, replace to the func of depend on mainline instance association
func (assoc *association) checkInstNameRepeat(kit *rest.Kit, instParentMap map[int64]int64, children []mapstr.MapStr) (
	canReset bool, repeatedInstName string, err error) {

	parentChildName := map[int64]map[string]struct{}{}
	for _, child := range children {
		instName, err := child.String(common.BKInstNameField)
		if err != nil {
			return false, "", err
		}

		childInstID, err := child.Int64(common.BKInstParentStr)
		if err != nil {
			blog.Errorf("get child parent id in child insts failed, err: %v, rid: %s", err, kit.Rid)
			return false, "", err
		}

		for instID, parentID := range instParentMap {
			if instID != childInstID {
				continue
			}

			childName, exist := parentChildName[parentID]
			if !exist {
				childName = make(map[string]struct{})
				parentChildName[parentID] = childName
			}

			if _, exist := childName[instName]; exist {
				return false, instName, nil
			}

			childName[instName] = struct{}{}
		}
	}

	return true, "", nil
}

// TODO should be deleted after merge, and which call this func use DeleteMainlineInstWithID in inst/inst.go to replace
// DeleteMainlineInstWithID delete mainline instance by it's bk_inst_id
func (assoc *association) DeleteMainlineInstWithID(kit *rest.Kit, objID string, instIDs []int64) error {

	// if this instance has been bind to a instance by the association, then this instance should not be deleted.
	cnt, err := assoc.clientSet.CoreService().Association().CountInstanceAssociations(kit.Ctx, kit.Header, objID,
		&metadata.Condition{
			Condition: mapstr.MapStr{common.BKDBOR: []mapstr.MapStr{
				{common.BKObjIDField: objID, common.BKInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}},
				{common.BKAsstObjIDField: objID, common.BKAsstInstIDField: mapstr.MapStr{common.BKDBIN: instIDs}},
			}},
		})
	if err != nil {
		blog.Errorf("count association by object(%s) failed, err: %s, rid: %s", objID, err, kit.Rid)
		return err
	}

	if cnt.Count > 0 {
		return kit.CCError.CCError(common.CCErrorInstHasAsst)
	}

	// delete this instance now.
	delCond := mapstr.MapStr{metadata.GetInstIDFieldByObjID(objID): mapstr.MapStr{common.BKDBIN: instIDs}}
	if metadata.IsCommon(objID) {
		delCond.Set(common.BKObjIDField, objID)
	}

	// generate audit log.
	audit := auditlog.NewInstanceAudit(assoc.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLogByCondGetData(generateAuditParameter, objID, delCond)
	if err != nil {
		blog.Errorf(" delete inst, generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// to delete.
	ops := metadata.DeleteOption{
		Condition: delCond,
	}
	_, err = assoc.clientSet.CoreService().Instance().DeleteInstance(kit.Ctx, kit.Header, objID, &ops)
	if err != nil {
		blog.Errorf("failed to delete the object(%s) inst by the condition(%#v), err: %v", objID, ops, err)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("delete inst, save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrAuditSaveLogFailed)
	}

	return nil
}
