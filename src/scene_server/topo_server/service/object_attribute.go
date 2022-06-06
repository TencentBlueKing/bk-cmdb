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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// CreateObjectAttribute create a new object attribute
func (s *Service) CreateObjectAttribute(ctx *rest.Contexts) {
	attr := new(metadata.Attribute)
	if err := ctx.DecodeInto(&attr); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// do not support add preset attribute by api
	attr.IsPre = false
	isBizCustomField := false
	// adapt input path param with bk_biz_id
	if bizIDStr := ctx.Request.PathParameter(common.BKAppIDField); bizIDStr != "" {
		bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
		bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
		if err != nil {
			blog.Errorf("create biz custom field, but parse biz ID failed, error: %s, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
			return
		}
		if bizID == 0 {
			blog.Errorf("create biz custom field, but biz ID is 0, rid: %s", ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
			return
		}
		attr.BizID = bizID
		isBizCustomField = true
	}

	attrInfo := new(metadata.ObjAttDes)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		attribute, err := s.Logics.AttributeOperation().CreateObjectAttribute(ctx.Kit, attr)
		if err != nil {
			return err
		}
		if attribute == nil {
			return err
		}
		attrInfo.Attribute = *attribute
		attrInfo.PropertyGroupName = attribute.PropertyGroupName

		if isBizCustomField {
			attrInfo.BizID = attr.BizID
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(attrInfo)
}

// SearchObjectAttribute search the object attributes
func (s *Service) SearchObjectAttribute(ctx *rest.Contexts) {
	dataWithModelBizID := MapStrWithModelBizID{}
	if err := ctx.DecodeInto(&dataWithModelBizID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	data := dataWithModelBizID.Data
	util.AddModelBizIDCondition(data, dataWithModelBizID.ModelBizID)
	data[metadata.AttributeFieldIsSystem] = false
	data[metadata.AttributeFieldIsAPI] = false

	basePage := metadata.BasePage{}
	if data.Exists(metadata.PageName) {
		page, err := data.MapStr(metadata.PageName)
		if err != nil {
			blog.Errorf("page info convert to mapstr failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		if err := mapstruct.Decode2Struct(page, &basePage); err != nil {
			blog.Errorf("page info convert to struct failed, page: %v, err: %v, rid: %s", page, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		data.Remove(metadata.PageName)
	}

	queryCond := &metadata.QueryCondition{
		Condition: data,
		Page:      basePage,
	}
	resp, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttrByCondition(ctx.Kit.Ctx, ctx.Kit.Header, queryCond)
	if nil != err {
		ctx.RespAutoError(err)
		return
	}

	grpMap, err := s.getPropertyGroupName(ctx, resp.Info, dataWithModelBizID.ModelBizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	attrInfos := make([]*metadata.ObjAttDes, 0)
	for _, attr := range resp.Info {
		attrInfo := &metadata.ObjAttDes{
			Attribute: attr,
		}
		grpName, ok := grpMap[attr.PropertyGroup]
		if !ok {
			blog.Errorf("failed to get property group name, attr: %s, property: %s", attr, attr.PropertyGroup)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKPropertyNameField))
			return
		}
		attrInfo.PropertyGroupName = grpName
		attrInfos = append(attrInfos, attrInfo)
	}

	ctx.RespEntity(attrInfos)
}

// UpdateObjectAttribute update the object attribute
func (s *Service) UpdateObjectAttribute(ctx *rest.Contexts) {
	data := make(mapstr.MapStr)
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the path params id: %s, err: %s, rid: %s", ctx.Request.PathParameter("id"),
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	// adapt input path param with bk_biz_id
	var bizID int64
	if bizIDStr := ctx.Request.PathParameter(common.BKAppIDField); bizIDStr != "" {
		bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
		bizID, err = strconv.ParseInt(bizIDStr, 10, 64)
		if err != nil {
			blog.Errorf("create biz custom field, but parse biz ID failed, error: %s, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
			return
		}
		if bizID == 0 {
			blog.Errorf("create biz custom field, but biz ID is 0, rid: %s", ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField))
			return
		}
	}
	// TODO: why does remove this????
	data.Remove(metadata.BKMetadata)
	data.Remove(common.BKAppIDField)

	// UpdateObjectAttribute should not update bk_property_index、bk_property_group
	data.Remove(common.BKPropertyIndexField)
	data.Remove(common.BKPropertyGroupField)

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.AttributeOperation().UpdateObjectAttribute(ctx.Kit, data, id, bizID)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(nil)
}

// DeleteObjectAttribute delete the object attribute
func (s *Service) DeleteObjectAttribute(ctx *rest.Contexts) {

	paramPath := mapstr.MapStr{}
	paramPath.Set("id", ctx.Request.PathParameter("id"))
	id, err := paramPath.Int64("id")
	if err != nil {
		blog.Errorf("failed to parse the path params id: %s, err: %s , rid: %s", ctx.Request.PathParameter("id"),
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	listRuleOption := metadata.ListHostApplyRuleOption{
		ModuleIDs: []int64{id},
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	ruleResult, err := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header,
		0, listRuleOption)
	if err != nil {
		blog.Errorf("get host apply rule failed, listRuleOption: %+v, err: %+v, rid: %s", listRuleOption, err,
			ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}
	ruleIDs := make([]int64, 0)
	for _, item := range ruleResult.Info {
		ruleIDs = append(ruleIDs, item.ID)
	}

	modelType := new(ModelType)
	if err := ctx.DecodeInto(modelType); err != nil {
		ctx.RespAutoError(err)
		return
	}

	cond := mapstr.MapStr{metadata.AttributeFieldID: id}
	if err = s.Logics.AttributeOperation().DeleteObjectAttribute(ctx.Kit, cond, modelType.BizID); err != nil {
		blog.Errorf("delete object attribute failed, params: %+v, err: %+v, rid: %s", cond, err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if len(ruleIDs) > 0 {
		deleteRuleOption := metadata.DeleteHostApplyRuleOption{
			RuleIDs: ruleIDs,
		}
		if err := s.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(ctx.Kit.Ctx, ctx.Kit.Header, 0,
			deleteRuleOption); err != nil {
			blog.Errorf("delete host apply rule failed, params: %+v, err: %+v, rid: %s", deleteRuleOption, err,
				ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
	}
	ctx.RespEntity(nil)
}

// UpdateObjectAttributeIndex update object attribute index
func (s *Service) UpdateObjectAttributeIndex(ctx *rest.Contexts) {
	data := make(map[string]interface{})
	if err := ctx.DecodeInto(&data); err != nil {
		ctx.RespAutoError(err)
		return
	}
	objID := ctx.Request.PathParameter(common.BKObjIDField)

	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the id from path params, id: %s, err: %s , rid: %s",
			ctx.Request.PathParameter("id"), err, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoPathParamPaserFailed))
		return
	}

	var result *metadata.UpdateAttrIndexData
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		var err error
		input := metadata.UpdateOption{
			Condition: mapstr.MapStr{common.BKFieldID: id},
			Data:      data,
		}
		result, err = s.Engine.CoreAPI.CoreService().Model().UpdateModelAttrsIndex(ctx.Kit.Ctx, ctx.Kit.Header, objID,
			&input)
		if err != nil {
			blog.Errorf("update object attribute index failed, err: %v , rid: %s", err, ctx.Kit.Rid)
			return err
		}
		return nil
	})

	if txnErr != nil {
		ctx.RespAutoError(txnErr)
		return
	}
	ctx.RespEntity(result)
}

// ListHostModelAttribute list host model's attributes
func (s *Service) ListHostModelAttribute(ctx *rest.Contexts) {
	dataWithModelBizID := MapStrWithModelBizID{}
	if err := ctx.DecodeInto(&dataWithModelBizID); err != nil {
		ctx.RespAutoError(err)
		return
	}
	data := dataWithModelBizID.Data
	data[metadata.AttributeFieldIsSystem] = false
	data[metadata.AttributeFieldIsAPI] = false
	util.AddModelBizIDCondition(data, dataWithModelBizID.ModelBizID)

	basePage := metadata.BasePage{}
	if data.Exists(metadata.PageName) {
		page, err := data.MapStr(metadata.PageName)
		if err != nil {
			blog.Errorf("page info convert to mapstr failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		if err := mapstruct.Decode2Struct(page, &basePage); err != nil {
			blog.Errorf("page info convert to struct failed, page: %v, err: %v, rid: %s", page, err, ctx.Kit.Rid)
			ctx.RespAutoError(err)
			return
		}
		data.Remove(metadata.PageName)
	}

	queryCond := &metadata.QueryCondition{
		Condition: data,
		Page:      basePage,
	}

	result, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttr(ctx.Kit.Ctx, ctx.Kit.Header,
		common.BKInnerObjIDHost, queryCond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	grpMap, err := s.getPropertyGroupName(ctx, result.Info, dataWithModelBizID.ModelBizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	hostAttributes := make([]metadata.HostObjAttDes, 0)
	for _, item := range result.Info {
		hostApplyEnabled := metadata.CheckAllowHostApplyOnField(&item)
		hostAttribute := metadata.HostObjAttDes{
			ObjAttDes: metadata.ObjAttDes{
				Attribute: item,
			},
			HostApplyEnabled: hostApplyEnabled,
		}
		grpName, ok := grpMap[item.PropertyGroup]
		if !ok {
			blog.Errorf("failed to get property group name, attr: %s, property: %s", item, item.PropertyGroup)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKPropertyNameField))
			return
		}
		hostAttribute.ObjAttDes.PropertyGroupName = grpName
		hostAttributes = append(hostAttributes, hostAttribute)
	}
	ctx.RespEntity(hostAttributes)
}

func (s *Service) getPropertyGroupName(ctx *rest.Contexts, attrs []metadata.Attribute,
	modelBizID int64) (map[string]string, error) {
	if len(attrs) == 0 {
		return make(map[string]string), nil
	}

	grpOrCond := make([]map[string]interface{}, 0)
	for _, attr := range attrs {
		grpOrCond = append(grpOrCond, map[string]interface{}{
			metadata.GroupFieldGroupID:  attr.PropertyGroup,
			metadata.GroupFieldObjectID: attr.ObjectID,
		})
	}
	grpCond := map[string]interface{}{
		common.BKDBOR: grpOrCond,
	}
	util.AddModelBizIDCondition(grpCond, modelBizID)
	cond := metadata.QueryCondition{
		Condition:      grpCond,
		DisableCounter: true,
	}
	rsp, err := s.Engine.CoreAPI.CoreService().Model().ReadAttributeGroupByCondition(ctx.Kit.Ctx, ctx.Kit.Header, cond)
	if err != nil {
		blog.Errorf("failed to get attr group, err: %s, rid: %s", err.Error(), ctx.Kit.Rid)
		return nil, err
	}

	grpMap := make(map[string]string)
	for _, grp := range rsp.Info {
		grpMap[grp.GroupID] = grp.GroupName
	}

	return grpMap, nil
}
