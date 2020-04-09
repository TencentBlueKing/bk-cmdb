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

package service

import (
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/topo_server/core/operation"
)

// CreateObjectGroup create a new object group

func (s *Service) CreateObjectGroup(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	rsp, err := s.Core.GroupOperation().CreateObjectGroup(ctx.Kit, dataWithMetadata.Data, dataWithMetadata.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	// auth: register attribute group
	if err := s.AuthManager.RegisterModelAttributeGroup(ctx.Kit.Ctx, ctx.Kit.Header, rsp.Group()); err != nil {
		blog.Errorf("create object group success, but register attribute group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
		return
	}

	retData := rsp.ToMapStr()
	id, err := retData.Int64("id")
	if err != nil {
		blog.Errorf("create object group success, but get response id failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
	}
	objAttrGroupAuditLog := operation.NewObjectAttrGroupAudit(s.Engine.CoreAPI)

	//get CurData for auditLog
	err = objAttrGroupAuditLog.MakeCurrent(rsp.ToMapStr())
	if err != nil {
		blog.Errorf("[operation-obj] make Current object failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//get AuditLog objectName
	err = objAttrGroupAuditLog.GetObjectInfo(ctx.Kit, rsp.Group().ObjectID)
	if err != nil {
		blog.Errorf("[api-att] find objectInfo failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//package audit response
	err = objAttrGroupAuditLog.SaveAuditLog(ctx.Kit, metadata.AuditCreate)
	if err != nil {
		ctx.RespAutoError(err)
	}
	ctx.RespEntity(retData)
}

// UpdateObjectGroup update the object group information
func (s *Service) UpdateObjectGroup(ctx *rest.Contexts) {
	cond := &metadata.UpdateGroupCondition{}
	err := ctx.DecodeInto(cond)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	objAttrGroupAuditLog := operation.NewObjectAttrGroupAudit(s.Engine.CoreAPI)
	//get AuditLog PreData
	err = objAttrGroupAuditLog.WithPrevious(ctx.Kit, cond.Condition.ID)
	if err != nil {
		blog.Errorf("[operation-obj] find Previous objectGroup failed, id: %+v, err: %s, rid: %s", cond.Condition.ID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//get AuditLog objectName
	err = objAttrGroupAuditLog.GetObjectInfo(ctx.Kit, "")
	if err != nil {
		blog.Errorf("[api-att] find objectInfo failed, id: %+v, err: %s, rid: %s", cond.Condition.ID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	err = s.Core.GroupOperation().UpdateObjectGroup(ctx.Kit, cond)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	// query attribute groups with given condition, so that update them to iam after updated
	searchCondition := condition.CreateCondition()
	if cond.Condition.ID != 0 {
		searchCondition.Field(common.BKFieldID).Eq(cond.Condition.ID)
	}
	result, err := s.Core.GroupOperation().FindObjectGroup(ctx.Kit, searchCondition, cond.Metadata)
	if err != nil {
		blog.Errorf("search attribute group by condition failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	attributeGroups := make([]metadata.Group, 0)
	for _, item := range result {
		attributeGroups = append(attributeGroups, item.Group())
	}

	// auth: register attribute group
	if err := s.AuthManager.UpdateRegisteredModelAttributeGroup(ctx.Kit.Ctx, ctx.Kit.Header, attributeGroups...); err != nil {
		blog.Errorf("update object group success, but update attribute group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
		return
	}

	//get AuditLog CurData
	err = objAttrGroupAuditLog.WithCurrent(ctx.Kit, cond.Condition.ID)
	if err != nil {
		blog.Errorf("[operation-obj] find Current object group failed, id: %+v, err: %s, rid: %s", cond.Condition.ID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//package audit response
	err = objAttrGroupAuditLog.SaveAuditLog(ctx.Kit, metadata.AuditUpdate)
	if err != nil {
		ctx.RespAutoError(err)
	}
	ctx.RespEntity(nil)
}

// DeleteObjectGroup delete the object group
func (s *Service) DeleteObjectGroup(ctx *rest.Contexts) {
	gid, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	objAttrGroupAuditLog := operation.NewObjectAttrGroupAudit(s.Engine.CoreAPI)
	//get AuditLog PreData
	err = objAttrGroupAuditLog.WithPrevious(ctx.Kit, gid)
	if err != nil {
		blog.Errorf("[operation-obj] find Previous object group failed, id: %+v, err: %s, rid: %s", gid, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}
	//get AuditLog objectName
	err = objAttrGroupAuditLog.GetObjectInfo(ctx.Kit, "")
	if err != nil {
		blog.Errorf("[api-att] find objectInfo failed, id: %+v, err: %s, rid: %s", gid, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	err = s.Core.GroupOperation().DeleteObjectGroup(ctx.Kit, gid)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	// auth: deregister attribute group
	if err := s.AuthManager.DeregisterModelAttributeGroupByID(ctx.Kit.Ctx, ctx.Kit.Header, gid); err != nil {
		blog.Errorf("delete object group failed, deregister attribute group to iam failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed))
		return
	}

	//package audit response
	err = objAttrGroupAuditLog.SaveAuditLog(ctx.Kit, metadata.AuditDelete)
	if err != nil {
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(nil)
}

// UpdateObjectAttributeGroupProperty update the object attribute belongs to group information
func (s *Service) UpdateObjectAttributeGroupProperty(ctx *rest.Contexts) {
	requestBody := struct {
		Data               []metadata.PropertyGroupObjectAtt `json:"data" field:"json"`
		*metadata.Metadata `json:"metadata"`
	}{}
	if err := ctx.DecodeInto(&requestBody); err != nil {
		ctx.RespAutoError(err)
		return
	}

	objectAtt := requestBody.Data
	if objectAtt == nil {
		ctx.RespAutoError(ctx.Kit.CCError.New(common.CCErrCommParamsIsInvalid, "param not set"))
		return
	}

	err := s.Core.GroupOperation().UpdateObjectAttributeGroup(ctx.Kit, objectAtt, requestBody.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

// DeleteObjectAttributeGroup delete the object attribute belongs to group information

func (s *Service) DeleteObjectAttributeGroup(ctx *rest.Contexts) {
	err := s.Core.GroupOperation().DeleteObjectAttributeGroup(ctx.Kit, ctx.Request.PathParameter("bk_object_id"), ctx.Request.PathParameter("property_id"), ctx.Request.PathParameter("group_id"))
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(nil)
}

// SearchGroupByObject search the groups by the object
func (s *Service) SearchGroupByObject(ctx *rest.Contexts) {
	cond := condition.CreateCondition()

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	resp, err := s.Core.GroupOperation().FindGroupByObject(ctx.Kit, ctx.Request.PathParameter("bk_obj_id"), cond, md.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)

}
