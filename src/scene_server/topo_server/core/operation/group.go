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
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// GroupOperationInterface group operation methods
type GroupOperationInterface interface {
	CreateObjectGroup(kit *rest.Kit, data mapstr.MapStr, modelBizID int64) (model.GroupInterface, error)
	DeleteObjectGroup(kit *rest.Kit, groupID int64) error
	FindObjectGroup(kit *rest.Kit, cond condition.Condition, modelBizID int64) ([]model.GroupInterface, error)
	FindGroupByObject(kit *rest.Kit, objID string, cond condition.Condition, modelBizID int64) ([]model.GroupInterface, error)
	UpdateObjectGroup(kit *rest.Kit, cond *metadata.UpdateGroupCondition) error
	UpdateObjectAttributeGroup(kit *rest.Kit, cond []metadata.PropertyGroupObjectAtt, modelBizID int64) error
	DeleteObjectAttributeGroup(kit *rest.Kit, objID, propertyID, groupID string) error
	SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface)
}

// NewGroupOperation create a new group operation instance
func NewGroupOperation(client apimachinery.ClientSetInterface) GroupOperationInterface {
	return &group{
		clientSet: client,
	}
}

type group struct {
	clientSet    apimachinery.ClientSetInterface
	modelFactory model.Factory
	instFactory  inst.Factory
	obj          ObjectOperationInterface
}

func (g *group) SetProxy(modelFactory model.Factory, instFactory inst.Factory, obj ObjectOperationInterface) {
	g.modelFactory = modelFactory
	g.instFactory = instFactory
	g.obj = obj
}

func (g *group) CreateObjectGroup(kit *rest.Kit, data mapstr.MapStr, modelBizID int64) (model.GroupInterface, error) {
	grp := g.modelFactory.CreateGroup(kit, modelBizID)

	_, err := grp.Parse(data)
	if nil != err {
		blog.Errorf("[operation-grp] failed to parse the group data(%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	//  check the object
	if err = g.obj.IsValidObject(kit, grp.Group().ObjectID); nil != err {
		blog.Errorf("[operation-grp] the group (%#v) is in valid, rid: %s", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// create a new group
	err = grp.Create()
	if nil != err {
		blog.Errorf("[operation-grp] failed to save the group data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(g.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, grp.Group().ID, nil)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but generate audit log failed, err: %v, rid: %s",
			grp.Group().GroupName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute group %s success, but save audit log failed, err: %v, rid: %s",
			grp.Group().GroupName, err, kit.Rid)
		return nil, err
	}

	return grp, nil
}

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
	cond := condition.CreateCondition().Field(common.BKFieldID).Eq(groupID)
	rsp, err := g.clientSet.CoreService().Model().DeleteAttributeGroupByCondition(context.Background(), kit.Header, metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-grp]failed to request object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return err
	}
	if !rsp.Result {
		blog.Errorf("[operation-grp]failed to delete the group(%d), err: %s, rid: %s", groupID, rsp.ErrMsg, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectGroupDeleteFailed)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("delete object attribute group success, but save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func (g *group) FindObjectGroup(kit *rest.Kit, cond condition.Condition, modelBizID int64) ([]model.GroupInterface, error) {
	fCond := cond.ToMapStr()
	util.AddModelBizIDConditon(fCond, modelBizID)

	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroupByCondition(context.Background(), kit.Header, metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[opeartion-grp] failed to search the group by the condition(%#v), error info is %s, rid: %s", fCond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return model.CreateGroup(kit, g.clientSet, rsp.Data.Info), nil
}

func (g *group) FindGroupByObject(kit *rest.Kit, objID string, cond condition.Condition, modelBizID int64) ([]model.GroupInterface, error) {
	fCond := cond.ToMapStr()
	util.AddModelBizIDConditon(fCond, modelBizID)

	rsp, err := g.clientSet.CoreService().Model().ReadAttributeGroup(context.Background(), kit.Header, objID, metadata.QueryCondition{Condition: fCond})
	if nil != err {
		blog.Errorf("[operation-grp] failed to request the object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to search the group of the object(%s) by the condition (%#v), error info is %s, rid: %s", objID, fCond, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return model.CreateGroup(kit, g.clientSet, rsp.Data.Info), nil
}

func (g *group) UpdateObjectAttributeGroup(kit *rest.Kit, conds []metadata.PropertyGroupObjectAtt, modelBizID int64) error {
	for _, cond := range conds {
		// if the target group doesn't exist, don't change the original group
		grpCond := condition.CreateCondition()
		grpCond.Field(metadata.GroupFieldGroupID).Eq(cond.Data.PropertyGroupID)
		grpCond.Field(metadata.GroupFieldObjectID).Eq(cond.Condition.ObjectID)
		grps, err := g.FindObjectGroup(kit, grpCond, modelBizID)
		if nil != err {
			blog.Errorf("[operation-grp] failed to get the group  by the condition (%#v), error info is %s , rid: %s", cond, err.Error(), kit.Rid)
			return err
		}
		if len(grps) != 1 {
			blog.Errorf("[operation-grp] failed to set the group  by the condition (%#v), error info is group is invalid, rid: %s", cond, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.GroupFieldGroupID)
		}
	}

	for _, cond := range conds {
		input := metadata.UpdateOption{
			Condition: mapstr.NewFromStruct(cond.Condition, "json"),
			Data:      mapstr.NewFromStruct(cond.Data, "json"),
		}

		rsp, err := g.clientSet.CoreService().Model().UpdateModelAttrsByCondition(context.Background(), kit.Header, &input)
		if nil != err {
			blog.Errorf("[operation-grp] failed to set the group  by the condition (%#v), error info is %s , rid: %s", cond, err.Error(), kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		if !rsp.Result {
			blog.Errorf("[operation-grp] failed to set the group  by the condition (%#v), error info is %s , rid: %s", cond, rsp.ErrMsg, kit.Rid)
			return kit.CCError.New(rsp.Code, rsp.ErrMsg)
		}
	}

	return nil
}

func (g *group) DeleteObjectAttributeGroup(kit *rest.Kit, objID, propertyID, groupID string) error {
	input := metadata.UpdateOption{
		Condition: condition.CreateCondition().
			Field(common.BKObjIDField).Eq(objID).
			Field(common.BKPropertyIDField).Eq(propertyID).
			Field(common.BKPropertyGroupField).Eq(groupID).ToMapStr(),
		Data: mapstr.MapStr{
			"bk_property_index":         -1,
			common.BKPropertyGroupField: "default",
		},
	}

	rsp, err := g.clientSet.CoreService().Model().UpdateModelAttrs(context.Background(), kit.Header, objID, &input)
	if nil != err {
		blog.Errorf("[operation-grp] failed to set the group , error info is %s , rid: %s", err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to set the group , error info is %s , rid: %s", rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	return nil
}

func (g *group) UpdateObjectGroup(kit *rest.Kit, cond *metadata.UpdateGroupCondition) error {
	if cond.Data.Index == nil && cond.Data.Name == nil {
		return nil
	}
	input := metadata.UpdateOption{
		Condition: mapstr.NewFromStruct(cond.Condition, "json"),
		Data:      mapstr.NewFromStruct(cond.Data, "json"),
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(g.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(input.Data)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, cond.Condition.ID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update attribute group, groupName: %s, err: %v, rid: %s",
			cond.Data.Name, err, kit.Rid)
		return err
	}

	// to update.
	rsp, err := g.clientSet.CoreService().Model().UpdateAttributeGroupByCondition(context.Background(), kit.Header, input)
	if nil != err {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s , rid: %s", cond.Data, cond.Condition, err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}
	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s , rid: %s", cond.Data, cond.Condition, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object attribute group %s success, but save audit log failed, err: %v, rid: %s",
			cond.Data.Name, err, kit.Rid)
		return err
	}
	return nil
}
