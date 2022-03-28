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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateMainlineAssociation create mainline object association
func (assoc *association) CreateMainlineAssociation(kit *rest.Kit, data *metadata.MainlineAssociation) (
	*metadata.Object, error) {

	if data.AsstObjID == "" {
		blog.Errorf("bk_asst_obj_id empty, input: %#v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKAsstObjIDField)
	}

	if data.ClassificationID == "" {
		blog.Errorf("bk_classification_id empty, input: %#v, rid: %s", data, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKClassificationIDField)
	}

	if err := assoc.checkMaxBizTopoLevel(kit); err != nil {
		return nil, err
	}

	mlCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{common.AssociationKindIDField: common.AssociationKindMainline},
	}
	mainlineAsst, ccErr := assoc.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, mlCond)
	if ccErr != nil {
		blog.Errorf("failed to check the mainline topo level, err: %v, rid: %s", ccErr, kit.Rid)
		return nil, ccErr
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

	objData := mapstr.MapStr{
		common.BKObjIDField:            data.ObjectID,
		common.BKObjNameField:          data.ObjectName,
		common.BKObjIconField:          data.ObjectIcon,
		common.BKClassificationIDField: data.ClassificationID,
		common.BkSupplierAccount:       data.OwnerID,
	}
	currentObj, err := assoc.obj.CreateObject(kit, true, objData)
	if err != nil {
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

	// update the mainline topo inst association
	if _, err = assoc.instasst.SetMainlineInstAssociation(kit, parentObjID, childObjID, currentObj.ObjectID,
		currentObj.ObjectName); err != nil {
		blog.Errorf("failed set the mainline inst association, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}

	return currentObj, nil
}

func (assoc *association) checkMaxBizTopoLevel(kit *rest.Kit) error {

	items, err := assoc.SearchMainlineAssociationTopo(kit, common.BKInnerObjIDApp)
	if nil != err {
		blog.Errorf("[operation-asst] failed to check the mainline topo level, error info is %s, rid: %s", err.Error(),
			kit.Rid)
		return err
	}

	res, err := assoc.clientSet.CoreService().System().SearchPlatformSetting(kit.Ctx, kit.Header)
	if err != nil {
		blog.Errorf("get business topo level max failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.CCErrTopoObjectSelectFailed)
	}
	if res.Result == false {
		blog.Errorf("get business topo level max failed, search config admin err: %s, rid: %s", res.ErrMsg, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, common.CCErrTopoObjectSelectFailed)
	}

	if len(items) >= int(res.Data.Backend.MaxBizTopoLevel) {
		blog.Errorf("[operation-asst] the mainline topo level is %d, the max limit is %d, rid: %s", len(items),
			res.Data.Backend.MaxBizTopoLevel, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoBizTopoLevelOverLimit)
	}
	return nil
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

	if err = assoc.instasst.ResetMainlineInstAssociation(kit, targetObjID, childObjID); err != nil {
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
	objIDCond := mapstr.MapStr{common.BKObjIDField: targetObjID}
	if _, err = assoc.obj.DeleteObject(kit, objIDCond, false); err != nil {
		blog.Errorf("failed to delete the object(%s), err: %v, rid: %s", targetObjID, err, kit.Rid)
		return err
	}

	return nil
}

// SearchMainlineAssociationTopo get mainline topo of special model
// result is a list with targetObj as head, so if you want a full topo, target must be biz model.
func (assoc *association) SearchMainlineAssociationTopo(kit *rest.Kit, targetObjID string) (
	[]*metadata.MainlineObjectTopo, error) {

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
	for _, asso := range mainlineAssoc.Info {
		childMap[asso.AsstObjID] = asso.ObjectID
		parentMap[asso.ObjectID] = asso.AsstObjID
	}

	// 遍历获取以targetObj为头的topo切片
	needFind := []string{targetObjID}
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
	// judge whether it is an inner mainline model
	if common.IsInnerMainlineModel(objID) {
		return true, nil
	}

	filter := []map[string]interface{}{{
		common.AssociationKindIDField: common.AssociationKindMainline,
		common.BKAsstObjIDField:       objID,
	}}
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
