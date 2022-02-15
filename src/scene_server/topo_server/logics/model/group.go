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
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/rs/xid"
)

// GroupOperationInterface group operation methods
type GroupOperationInterface interface {
	// CreateObjectGroup create object group
	CreateObjectGroup(kit *rest.Kit, data *metadata.Group) (*metadata.Group, error)
	// DeleteObjectGroup delete object group
	DeleteObjectGroup(kit *rest.Kit, groupID int64) error
	// FindObjectGroup Use FindGroupByObject function instead
	// Deprecated: 后续将逐步替换为带bk_obj_id参数的 FindGroupByObject
	FindObjectGroup(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) ([]metadata.Group, error)
	// FindGroupByObject find group by object
	FindGroupByObject(kit *rest.Kit, objID string, cond mapstr.MapStr, modelBizID int64) ([]metadata.Group, error)
	// UpdateObjectGroup update object group
	UpdateObjectGroup(kit *rest.Kit, cond *metadata.UpdateGroupCondition) error
	// UpdateObjectAttributeGroup Update Object Attribute group
	UpdateObjectAttributeGroup(kit *rest.Kit, cond []metadata.PropertyGroupObjectAtt, modelBizID int64) error
	// ExchangeObjectGroupIndex exchange the group index of two groups
	ExchangeObjectGroupIndex(kit *rest.Kit, ids []int64) error
	// DeleteObjectAttributeGroup delete object attribute group
	DeleteObjectAttributeGroup(kit *rest.Kit, objID, propertyID, groupID string) error
	// SetProxy SetProxy
	SetProxy(obj ObjectOperationInterface)
}

// NewGroupOperation create a new group operation instance
func NewGroupOperation(client apimachinery.ClientSetInterface) GroupOperationInterface {
	return &group{
		clientSet: client,
	}
}

// group group
type group struct {
	clientSet apimachinery.ClientSetInterface
	obj       ObjectOperationInterface
}

// SetProxy SetProxy
func (g *group) SetProxy(obj ObjectOperationInterface) {
	g.obj = obj
}

// CreateObjectGroup create object group
func (g *group) CreateObjectGroup(kit *rest.Kit, data *metadata.Group) (*metadata.Group, error) {

	if len(data.GroupName) == 0 {
		blog.Errorf("failed to valid the group id(%s), rid: %s", metadata.GroupFieldGroupID, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupID)
	}
	if err := util.ValidModelNameField(data.GroupName, metadata.GroupFieldGroupName, kit.CCError); err != nil {
		blog.Errorf("failed to valid the group name(%s), rid: %s", metadata.GroupFieldGroupName, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupName)
	}

	// check the object
	isObjExists, err := g.obj.IsObjectExist(kit, data.ObjectID)
	if err != nil {
		blog.Errorf("check if object(%s) exists failed, err: %v, rid: %s", data.ObjectID, err, kit.Rid)
		return nil, err
	}
	if !isObjExists {
		blog.Errorf("object (%s) does not exist, rid: %s", data.ObjectID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	// create a new group
	rsp, err := g.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header, data.ObjectID,
		metadata.CreateModelAttributeGroup{Data: *data})
	if err != nil {
		blog.Errorf("create attribute group %s failed, err: %v, rid: %s", data.GroupName, err, kit.Rid)
		return nil, err
	}

	data.ID = int64(rsp.Created.ID)

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(g.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, data.ID, data)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but generate audit log failed, err: %v, rid: %s",
			data.GroupName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute group %s success, but save audit log failed, err: %v, rid: %s",
			data.GroupName, err, kit.Rid)
		return nil, err
	}

	return data, nil
}

// DeleteObjectGroup delete object group
func (g *group) DeleteObjectGroup(kit *rest.Kit, groupID int64) error {

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(g.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, groupID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before delete attribute group, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// to delete.
	cond := metadata.DeleteOption{Condition: mapstr.MapStr{
		common.BKFieldID: groupID,
	}}

	_, err = g.clientSet.CoreService().Model().DeleteAttributeGroupByCondition(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("delete object group %d failed, err: %v, rid: %s", groupID, err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("delete object attribute group success, but save audit log failed, err: %v, rid: %s",
			err, kit.Rid)
		return err
	}
	return nil
}

// FindObjectGroup Use FindGroupByObject function instead
// Deprecated: 后续将逐步替换为带bk_obj_id参数的 FindGroupByObject
func (g *group) FindObjectGroup(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) ([]metadata.Group, error) {

	util.AddModelBizIDCondition(cond, modelBizID)

	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header,
		metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("find object group but do http request failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return rsp.Info, nil
}

// FindGroupByObject find group by object
func (g *group) FindGroupByObject(kit *rest.Kit, objID string, cond mapstr.MapStr, modelBizID int64) ([]metadata.Group,
	error) {

	util.AddModelBizIDCondition(cond, modelBizID)

	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroup(kit.Ctx, kit.Header, objID,
		metadata.QueryCondition{Condition: cond})
	if err != nil {
		blog.Errorf("find group by object but do http request failed, err: %v, cond: %#v, rid: %s", err, cond,
			kit.Rid)
		return nil, err
	}

	return rsp.Info, nil
}

// UpdateObjectGroup update object group
func (g *group) UpdateObjectGroup(kit *rest.Kit, cond *metadata.UpdateGroupCondition) error {

	if cond.Data.Name != nil && len(*cond.Data.Name) != 0 {
		if err := util.ValidModelNameField(*cond.Data.Name, metadata.GroupFieldGroupName, kit.CCError); err != nil {
			blog.Errorf("failed to valid the group name(%s), rid: %s", cond.Data.Name, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupName)
		}
	}

	input := metadata.UpdateOption{
		Condition: mapstr.MapStr{
			common.BKFieldID: cond.Condition.ID,
		},
		Data: mapstr.MapStr{},
	}
	if cond.Data.Name != nil {
		input.Data.Set(common.BKPropertyGroupNameField, cond.Data.Name)
	}
	if cond.Data.IsCollapse != nil {
		input.Data.Set(common.BKIsCollapseField, cond.Data.IsCollapse)
	}
	if cond.Data.Index != nil {
		input.Data.Set(common.BKPropertyGroupIndexField, cond.Data.Index)
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(g.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit,
		metadata.AuditUpdate).WithUpdateFields(input.Data)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, cond.Condition.ID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update attribute group, groupName: %s, err: %v, rid: %s",
			cond.Data.Name, err, kit.Rid)
		return err
	}

	// to update.
	_, err = g.clientSet.CoreService().Model().UpdateAttributeGroupByCondition(kit.Ctx, kit.Header, input)
	if err != nil {
		blog.Errorf("update object group failed, err: %v, input: %#v, rid: %s", err, cond, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object attribute group %s success, but save audit log failed, err: %v, rid: %s",
			cond.Data.Name, err, kit.Rid)
		return err
	}
	return nil
}

// UpdateObjectAttributeGroup Update Object Attribute group
func (g *group) UpdateObjectAttributeGroup(kit *rest.Kit, conds []metadata.PropertyGroupObjectAtt,
	modelBizID int64) error {

	propertyGroupIDs := make([]string, 0)
	objIds := make([]string, 0)
	for _, cond := range conds {
		propertyGroupIDs = append(propertyGroupIDs, cond.Data.PropertyGroupID)
		objIds = append(objIds, cond.Condition.ObjectID)
	}
	grpCond := mapstr.MapStr{
		metadata.GroupFieldGroupID: map[string][]string{
			common.BKDBIN: propertyGroupIDs,
		},
		metadata.GroupFieldObjectID: map[string][]string{
			common.BKDBIN: objIds,
		},
	}

	grps, err := g.FindObjectGroup(kit, grpCond, modelBizID)
	if err != nil {
		blog.Errorf("update object attribute group failed, err: %v, input: %#v, rid: %s", err, grpCond, kit.Rid)
		return err
	}

	grpMap := make(map[string]struct{})
	for _, grp := range grps {
		if _, ok := grpMap[grp.GroupID]; ok {
			blog.Errorf("there is more than one group (%s), rid: %s", grp.GroupID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupID)
		}
		grpMap[grp.GroupID] = struct{}{}
	}

	for _, propertyGroupID := range propertyGroupIDs {
		if _, ok := grpMap[propertyGroupID]; !ok {
			blog.Errorf("group (%s) does not exist, rid: %s", propertyGroupID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupID)
		}
	}

	for _, cond := range conds {
		input := metadata.UpdateOption{
			Condition: mapstr.MapStr{
				common.BKOwnerIDField:    cond.Condition.OwnerID,
				common.BKObjIDField:      cond.Condition.ObjectID,
				common.BKPropertyIDField: cond.Condition.PropertyID,
			},
			Data: mapstr.MapStr{
				common.BKPropertyGroupField: cond.Data.PropertyGroupID,
				common.BKPropertyIndexField: cond.Data.PropertyIndex,
			},
		}

		_, err = g.clientSet.CoreService().Model().UpdateModelAttrsByCondition(kit.Ctx, kit.Header, &input)
		if err != nil {
			blog.Errorf("update model attrs failed, err: %v, input: %#v, rid: %s", err, cond, kit.Rid)
			return err
		}

	}

	return nil
}

// DeleteObjectAttributeGroup delete object attribute group
func (g *group) DeleteObjectAttributeGroup(kit *rest.Kit, objID, propertyID, groupID string) error {

	input := metadata.UpdateOption{
		Data: mapstr.MapStr{
			common.BKPropertyIndexField: -1,
			common.BKPropertyGroupField: common.BKDefaultField,
		},
		Condition: mapstr.MapStr{
			common.BKObjIDField:         objID,
			common.BKPropertyIDField:    propertyID,
			common.BKPropertyGroupField: groupID,
		},
	}

	_, err := g.clientSet.CoreService().Model().UpdateModelAttrs(kit.Ctx, kit.Header, objID, &input)
	if err != nil {
		blog.Errorf("delete object attribute group failed, err: %v, input: %#v, rid: %s", err, input, kit.Rid)
		return err
	}

	return nil
}

// NewGroupID generate new group id, default group has a specific id
func NewGroupID(isDefault bool) string {
	if isDefault {
		return "default"
	} else {
		return xid.New().String()
	}
}

// ExchangeObjectGroupIndex exchange the group index of two groups
// related issue: https://github.com/Tencent/bk-cmdb/issues/5873
func (g *group) ExchangeObjectGroupIndex(kit *rest.Kit, ids []int64) error {

	if len(ids) != 2 {
		blog.Errorf("group ids must be two, now is %d, rid: %s", len(ids), kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKFieldID)
	}

	cond := metadata.QueryCondition{
		Condition: mapstr.MapStr{common.BKFieldID: mapstr.MapStr{common.BKDBIN: ids}},
		Fields: []string{common.BKPropertyGroupIndexField, common.BKFieldID, common.BKObjIDField,
			common.BKAppIDField},
	}
	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("search group info by ids failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if len(rsp.Info) != 2 {
		blog.Errorf("search group info by ids failed, result not two group info, number: %d, rid: %s",
			len(rsp.Info), kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectGroupUpdateFailed, "result not two group info")
	}

	// custom object create will create default group, index is -1, -2 has't been used
	// once update any group_index will cause Mongo duplicate error without temp
	tempIndex := int64(-2)
	groupA, groupB := rsp.Info[0], rsp.Info[1]
	if groupA.BizID != groupB.BizID || groupA.ObjectID != groupB.ObjectID {
		blog.Errorf("two groups not the same business or object, rid: %s", kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectGroupUpdateFailed, "two groups not the same business/object")
	}

	// update group A to a temp group_index
	updateCond := &metadata.UpdateGroupCondition{}
	updateCond.Condition.ID = groupA.ID
	updateCond.Data.Index = &tempIndex
	if err := g.UpdateObjectGroup(kit, updateCond); err != nil {
		blog.Errorf("failed to set the group to the new data (%#v) by the condition (%#v), err: %v , rid: %s",
			updateCond.Data, updateCond.Condition, err, kit.Rid)
		return err
	}

	// update group B to group A origin group_index
	updateCond.Condition.ID = groupB.ID
	updateCond.Data.Index = &groupA.GroupIndex
	if err := g.UpdateObjectGroup(kit, updateCond); err != nil {
		blog.Errorf("failed to set the group to the new data (%#v) by the condition (%#v), err: %v , rid: %s",
			updateCond.Data, updateCond.Condition, err, kit.Rid)
		return err
	}

	// update group A to group B origin group_index
	updateCond.Condition.ID = groupA.ID
	updateCond.Data.Index = &groupB.GroupIndex
	if err := g.UpdateObjectGroup(kit, updateCond); err != nil {
		blog.Errorf("failed to set the group to the new data (%#v) by the condition (%#v), err: %v , rid: %s",
			updateCond.Data, updateCond.Condition, err, kit.Rid)
		return err
	}

	return nil
}
