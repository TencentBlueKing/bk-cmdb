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
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/inst"
	"configcenter/src/scene_server/topo_server/core/model"
)

// GroupOperationInterface group operation methods
type GroupOperationInterface interface {
	CreateObjectGroup(kit *rest.Kit, data mapstr.MapStr, metaData *metadata.Metadata) (model.GroupInterface, error)
	DeleteObjectGroup(kit *rest.Kit, groupID int64) error
	FindObjectGroup(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.GroupInterface, error)
	FindGroupByObject(kit *rest.Kit, objID string, cond condition.Condition, metaData *metadata.Metadata) ([]model.GroupInterface, error)
	UpdateObjectGroup(kit *rest.Kit, cond *metadata.UpdateGroupCondition) error
	UpdateObjectAttributeGroup(kit *rest.Kit, cond []metadata.PropertyGroupObjectAtt, metaData *metadata.Metadata) error
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

func (g *group) CreateObjectGroup(kit *rest.Kit, data mapstr.MapStr, metaData *metadata.Metadata) (model.GroupInterface, error) {
	grp := g.modelFactory.CreateGroup(kit, metaData)

	_, err := grp.Parse(data)
	if nil != err {
		blog.Errorf("[operation-grp] failed to parse the group data(%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	//  check the object
	if err = g.obj.IsValidObject(kit, grp.Group().ObjectID, metaData); nil != err {
		blog.Errorf("[operation-grp] the group (%#v) is in valid, rid: %s", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// create a new group
	err = grp.Create()
	if nil != err {
		blog.Errorf("[operation-grp] failed to save the group data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	//package audit response
	err = NewObjectAudit(g.clientSet, metadata.ModelGroupRes).buildObjAttrGroupData(kit, grp.Group().ID).WithCurrent().SaveAuditLog(kit, metadata.AuditCreate)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but update to auditLog failed, err: %v, rid: %s", grp.Group().GroupName, err, kit.Rid)
		return nil, err
	}

	return grp, nil
}

func (g *group) DeleteObjectGroup(kit *rest.Kit, groupID int64) error {
	cond := condition.CreateCondition().Field(common.BKFieldID).Eq(groupID)

	//get PreData
	objAudit := NewObjectAudit(g.clientSet, metadata.ModelGroupRes).buildObjAttrGroupData(kit, groupID).WithPrevious()

	rsp, err := g.clientSet.CoreService().Model().DeleteAttributeGroupByCondition(context.Background(), kit.Header, metadata.DeleteOption{Condition: cond.ToMapStr()})
	if nil != err {
		blog.Errorf("[operation-grp]failed to request object controller, error info is %s, rid: %s", err.Error(), kit.Rid)
		return err
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp]failed to delete the group(%d), err: %s, rid: %s", groupID, rsp.ErrMsg, kit.Rid)
		return kit.CCError.Error(common.CCErrTopoObjectGroupDeleteFailed)
	}

	//saveAuditLog
	err = objAudit.SaveAuditLog(kit, metadata.AuditDelete)
	if err != nil {
		blog.Errorf("Delete object attribute group success, but update to auditLog failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func (g *group) FindObjectGroup(kit *rest.Kit, cond condition.Condition, metaData *metadata.Metadata) ([]model.GroupInterface, error) {
	fCond := cond.ToMapStr()
	if nil != metaData {
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}
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

func (g *group) FindGroupByObject(kit *rest.Kit, objID string, cond condition.Condition, metaData *metadata.Metadata) ([]model.GroupInterface, error) {
	fCond := cond.ToMapStr()
	if nil != metaData {
		fCond.Merge(metadata.PublicAndBizCondition(*metaData))
		fCond.Remove(metadata.BKMetadata)
	} else {
		fCond.Merge(metadata.BizLabelNotExist)
	}

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

func (g *group) UpdateObjectAttributeGroup(kit *rest.Kit, conds []metadata.PropertyGroupObjectAtt, metaData *metadata.Metadata) error {
	for _, cond := range conds {
		// if the target group doesn't exist, don't change the original group
		grpCond := condition.CreateCondition()
		grpCond.Field(metadata.GroupFieldGroupID).Eq(cond.Data.PropertyGroupID)
		grpCond.Field(metadata.GroupFieldObjectID).Eq(cond.Condition.ObjectID)
		grps, err := g.FindObjectGroup(kit, grpCond, metaData)
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

	//get PreData
	objAudit := NewObjectAudit(g.clientSet, metadata.ModelGroupRes).buildObjAttrGroupData(kit, cond.Condition.ID).WithPrevious()

	rsp, err := g.clientSet.CoreService().Model().UpdateAttributeGroupByCondition(context.Background(), kit.Header, input)
	if nil != err {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s , rid: %s", cond.Data, cond.Condition, err.Error(), kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	if !rsp.Result {
		blog.Errorf("[operation-grp] failed to set the group to the new data (%#v) by the condition (%#v), error info is %s , rid: %s", cond.Data, cond.Condition, rsp.ErrMsg, kit.Rid)
		return kit.CCError.New(rsp.Code, rsp.ErrMsg)
	}

	//get CurData and saveAuditLog
	err = objAudit.buildObjAttrGroupData(kit, cond.Condition.ID).WithCurrent().SaveAuditLog(kit, metadata.AuditUpdate)
	if err != nil {
		blog.Errorf("update object attribute group %s success, but update to auditLog failed, err: %v, rid: %s", cond.Data.Name, err, kit.Rid)
		return err
	}

	return nil
}
