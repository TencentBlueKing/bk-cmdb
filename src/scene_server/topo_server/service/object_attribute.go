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
	"configcenter/src/scene_server/topo_server/core/model"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
)

// CreateObjectAttribute create a new object attribute
func (s *Service) CreateObjectAttribute(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	attr, err := s.Core.AttributeOperation().CreateObjectAttribute(ctx.Kit, dataWithMetadata.Data, dataWithMetadata.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	blog.Debug(attr.ToMapStr())
	// auth: register resource
	attribute := attr.Attribute()
	if err := s.AuthManager.RegisterModelAttribute(ctx.Kit.Ctx, ctx.Kit.Header, *attribute); err != nil {
		blog.Errorf("create object attribute success, but register model attribute to auth failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
		return
	}

	cond := condition.CreateCondition()
	cond.Field("id").Eq(attribute.ID)
	attrInfo, err := s.Core.AttributeOperation().FindObjectAttributeWithDetail(ctx.Kit, cond, dataWithMetadata.Metadata)
	if err != nil {
		blog.Errorf("create object attribute success, but get attributes detail failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoSearchModelAttriFailedPleaseRefresh))
		return
	}
	if len(attrInfo) <= 0 {
		blog.Errorf("create object attribute success, but get attributes detail failed, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoSearchModelAttriFailedPleaseRefresh))
		return
	}

	objAuditLog := model.NewObjectAuditLog(s.Engine.CoreAPI, metadata.ModelAttributeType, metadata.ModelAttributeRes)
	//get CurData
	err = objAuditLog.WithCurrent(ctx.Kit, attribute.ID)
	if err != nil {
		blog.Errorf("find Current modelAttribute failed, id: %+v, err: %s, rid: %s", attribute.ID, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//package audit response
	err = objAuditLog.SaveAuditLog(ctx.Kit, metadata.AuditCreate)
	if err != nil {
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(attrInfo[0])
}

// SearchObjectAttribute search the object attributes
func (s *Service) SearchObjectAttribute(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	data := dataWithMetadata.Data

	cond := condition.CreateCondition()
	if data.Exists(metadata.PageName) {
		page, err := data.MapStr(metadata.PageName)
		if err != nil {
			blog.Errorf("SearchObjectAttribute failed, page info convert to mapstr failed, page: %v, err: %v, rid: %s", data[metadata.PageName], err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		if err := cond.SetPage(page); err != nil {
			blog.Errorf("SearchObjectAttribute, cond set page failed, page: %v, err: %v, rid: %v", page, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		data.Remove(metadata.PageName)
	}

	if err := cond.Parse(data); nil != err {
		blog.Errorf("search object attribute, but failed to parse the data into condition, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	cond.Field(metadata.AttributeFieldIsSystem).NotEq(true)
	cond.Field(metadata.AttributeFieldIsAPI).NotEq(true)

	resp, err := s.Core.AttributeOperation().FindObjectAttributeWithDetail(ctx.Kit, cond, dataWithMetadata.Metadata)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}
	ctx.RespEntity(resp)
}

// UpdateObjectAttribute update the object attribute
func (s *Service) UpdateObjectAttribute(ctx *rest.Contexts) {
	data := make(mapstr.MapStr)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params id(%s), error info is %s, rid: %s", ctx.Request.PathParameter("id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	objAuditLog := model.NewObjectAuditLog(s.Engine.CoreAPI, metadata.ModelAttributeType, metadata.ModelAttributeRes)
	//get AuditLog PreData
	err = objAuditLog.WithPrevious(ctx.Kit, id)
	if err != nil {
		blog.Errorf("[api-att] find Previous objectAttribute failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	// TODO: why does remove this????
	data.Remove(metadata.BKMetadata)

	// UpdateObjectAttribute should not update bk_property_index、bk_property_group
	data.Remove(common.BKPropertyIndexField)
	data.Remove(common.BKPropertyGroupField)

	err = s.Core.AttributeOperation().UpdateObjectAttribute(ctx.Kit, data, id)

	// auth: update registered resource
	if err := s.AuthManager.UpdateRegisteredModelAttributeByID(ctx.Kit.Ctx, ctx.Kit.Header, id); err != nil {
		blog.Errorf("update object attribute success , but update registered model attribute to auth failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.Error(common.CCErrCommRegistResourceToIAMFailed))
		return
	}

	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	//get AuditLog CurData
	err = objAuditLog.WithCurrent(ctx.Kit, id)
	if err != nil {
		blog.Errorf("[api-att] find Current objectAttribute failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	//package audit response
	err = objAuditLog.SaveAuditLog(ctx.Kit, metadata.AuditUpdate)
	if err != nil {
		ctx.RespAutoError(err)
	}

	ctx.RespEntity(nil)
}

// DeleteObjectAttribute delete the object attribute
func (s *Service) DeleteObjectAttribute(ctx *rest.Contexts) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("id", ctx.Request.PathParameter("id"))
	id, err := paramPath.Int64("id")
	if nil != err {
		blog.Errorf("[api-att] failed to parse the path params id(%s), error info is %s , rid: %s", ctx.Request.PathParameter("id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	objAuditLog := model.NewObjectAuditLog(s.Engine.CoreAPI, metadata.ModelAttributeType, metadata.ModelAttributeRes)
	//get AuditLog PreData
	err = objAuditLog.WithPrevious(ctx.Kit, id)
	if err != nil {
		blog.Errorf("[api-att] find Previous object attribute failed, id: %+v, err: %s, rid: %s", id, err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
	}

	cond := condition.CreateCondition()
	cond.Field(metadata.AttributeFieldID).Eq(id)

	listRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: []int64{id},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	ruleResult, err := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, 0, listRuleOption)
	if err != nil {
		blog.Errorf("delete object attribute failed, ListHostApplyRule failed, listRuleOption: %+v, err: %+v, rid: %s", listRuleOption, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ruleIDs := make([]int64, 0)
	for _, item := range ruleResult.Info {
		ruleIDs = append(ruleIDs, item.ID)
	}

	// auth: update registered resource
	if err := s.AuthManager.DeregisterModelAttributeByID(ctx.Kit.Ctx, ctx.Kit.Header, id); err != nil {
		blog.Errorf("delete object attribute failed, deregister model attribute to auth failed, err: %+v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrCommUnRegistResourceToIAMFailed))
		return
	}

	md := new(MetaShell)
	if err := ctx.DecodeInto(md); err != nil {
		ctx.RespAutoError(err)
		return
	}
	err = s.Core.AttributeOperation().DeleteObjectAttribute(ctx.Kit, cond, md.Metadata)
	if err != nil {
		blog.Errorf("delete object attribute failed, DeleteObjectAttribute failed, params: %+v, err: %+v, rid: %s", ctx.Kit, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(ruleIDs) > 0 {
		deleteRuleOption := metadata.DeleteHostApplyRuleOption{
			RuleIDs: ruleIDs,
		}
		if err := s.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, 0, deleteRuleOption); err != nil {
			blog.Errorf("delete object attribute success, but DeleteHostApplyRule failed, params: %+v, err: %+v, rid: %s", deleteRuleOption, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}

	ctx.RespEntity(nil)
}

func (s *Service) UpdateObjectAttributeIndex(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	objID := ctx.Request.PathParameter(common.BKObjIDField)

	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if nil != err {
		blog.Errorf("[api-att] failed to parse the params id(%s), error info is %s , rid: %s", ctx.Request.PathParameter("id"), err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoPathParamPaserFailed))
		return
	}
	result, err := s.Core.AttributeOperation().UpdateObjectAttributeIndex(ctx.Kit, objID, data, id)
	if err != nil {
		blog.Errorf("UpdateObjectAttributeIndex failed, error info is %s , rid: %s", err.Error(), ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(result)
}

// ListHostModelAttribute list host model's attributes
func (s *Service) ListHostModelAttribute(ctx *rest.Contexts) {
	dataWithMetadata := MapStrWithMetadata{}
	if err := ctx.DecodeInto(&dataWithMetadata); err != nil {
		ctx.RespAutoError(err)
		return
	}
	data := dataWithMetadata.Data
	cond := condition.CreateCondition()
	data.Remove(metadata.PageName)
	if err := cond.Parse(data); nil != err {
		blog.Errorf("search object attribute, but failed to parse the data into condition, err: %v, rid: %s", err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	cond.Field(metadata.AttributeFieldIsSystem).NotEq(true)
	cond.Field(metadata.AttributeFieldIsAPI).NotEq(true)
	cond.Field(common.BKObjIDField).Eq(common.BKInnerObjIDHost)

	attributes, err := s.Core.AttributeOperation().FindObjectAttributeWithDetail(ctx.Kit, cond, dataWithMetadata.Metadata)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	hostAttributes := make([]metadata.HostObjAttDes, 0)
	for _, item := range attributes {
		if item == nil {
			continue
		}
		hostApplyEnabled := metadata.CheckAllowHostApplyOnField(item.PropertyID)
		hostAttribute := metadata.HostObjAttDes{
			ObjAttDes:        *item,
			HostApplyEnabled: hostApplyEnabled,
		}
		hostAttributes = append(hostAttributes, hostAttribute)
	}
	ctx.RespEntity(hostAttributes)
}
