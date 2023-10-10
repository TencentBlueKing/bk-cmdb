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
	"fmt"
	"reflect"
	"strconv"
	"strings"

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
	if err := ctx.DecodeInto(attr); err != nil {
		ctx.RespAutoError(err)
		return
	}

	// if the attribute created here is a table attribute field,
	// it needs to go through a separate process.
	if attr.PropertyType == common.FieldTypeInnerTable {
		attrs, err := s.createTableAttribute(ctx, attr)
		if err != nil {
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntity(attrs)
		return
	}

	// 新建组织字段时，默认为多选，当api接口创建模型属性时，没有传ismultiple，默认置为true，支持多选
	if ok := checkJsonTagContainIsMultipleField(*attr); !ok {
		if attr.PropertyType == common.FieldTypeOrganization {
			isMultiple := true
			attr.IsMultiple = &isMultiple
		} else {
			isMultiple := false
			attr.IsMultiple = &isMultiple
		}
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
			blog.Errorf("create model attribute failed, attr: %+v, err: %v, rid: %s", attr, err, ctx.Kit.Rid)
			return err
		}
		if attribute == nil {
			blog.Errorf("return the created model attribute is empty, attr: %+v, rid: %s", attr, ctx.Kit.Rid)
			return fmt.Errorf("created model attribute is empty, attr: %v", attr)
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

func (s *Service) createTableAttribute(ctx *rest.Contexts, attr *metadata.Attribute) (*metadata.ObjAttDes, error) {

	bizID, err := parseRequestBizID(ctx)
	if err != nil {
		return nil, err
	}

	if attr.TemplateID != 0 {
		return nil, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKTemplateID)
	}

	isBizCustomField := false
	if bizID > 0 {
		attr.BizID = bizID
		isBizCustomField = true
	}

	if err := s.createTableObjectTable(ctx, attr.ObjectID, attr.PropertyID); err != nil {
		blog.Errorf("create table object table failed, attr: %+v, err: %v, rid: %s", *attr, err, ctx.Kit.Rid)
		return nil, err
	}

	attrInfo := new(metadata.ObjAttDes)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		attribute, err := s.Logics.AttributeOperation().CreateTableObjectAttribute(ctx.Kit, attr)
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
		return nil, txnErr
	}
	return attrInfo, nil
}

func parseRequestBizID(ctx *rest.Contexts) (int64, error) {
	bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
	if bizIDStr == "" {
		return 0, nil
	}

	bizID, err := strconv.ParseInt(bizIDStr, 10, 64)
	if err != nil {
		blog.Errorf("create biz custom field, but parse biz ID failed, error: %v, rid: %s", err, ctx.Kit.Rid)
		return 0, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	if bizID == 0 {
		blog.Errorf("create biz custom field, but biz ID is 0, rid: %s", ctx.Kit.Rid)
		return 0, ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKAppIDField)
	}
	return bizID, nil
}

// SearchObjectAttributeForWeb search form field attributes provided to the front end.
func (s *Service) SearchObjectAttributeForWeb(ctx *rest.Contexts) {

	queryCond, modelBizID, err := combinationSearchObjectAttrCond(ctx)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttrsWithTableByCondition(ctx.Kit.Ctx, ctx.Kit.Header,
		modelBizID, queryCond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	attrs := make([]*metadata.ObjAttDes, 0)
	if len(resp.Info) == 0 {
		ctx.RespEntity(attrs)
		return
	}

	grpMap, err := s.getPropertyGroupName(ctx.Kit, resp.Info, modelBizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	for _, attr := range resp.Info {
		attrInfo := &metadata.ObjAttDes{
			Attribute: attr,
		}
		grpName, ok := grpMap[attr.PropertyGroup]
		if !ok {
			blog.Errorf("failed to get property group name, attr: %+v, propertyGroup: %v, rid: %s",
				attr, attr.PropertyGroup, ctx.Kit.Rid)
			ctx.RespAutoError(ctx.Kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, common.BKPropertyNameField))
			return
		}
		attrInfo.PropertyGroupName = grpName
		attrs = append(attrs, attrInfo)
	}

	ctx.RespEntity(attrs)
}

func combinationSearchObjectAttrCond(ctx *rest.Contexts) (*metadata.QueryCondition, int64, error) {
	option := new(MapStrWithModelBizID)
	if err := ctx.DecodeInto(&option); err != nil {
		return nil, 0, err
	}
	data := option.Data
	util.AddModelBizIDCondition(data, option.ModelBizID)
	data[metadata.AttributeFieldIsSystem] = false
	data[metadata.AttributeFieldIsAPI] = false

	basePage := metadata.BasePage{}
	if data.Exists(metadata.PageName) {
		page, err := data.MapStr(metadata.PageName)
		if err != nil {
			blog.Errorf("page info convert to mapstr failed, err: %v, rid: %s", err, ctx.Kit.Rid)
			return nil, 0, err
		}
		if err := mapstruct.Decode2Struct(page, &basePage); err != nil {
			blog.Errorf("page info convert to struct failed, page: %v, err: %v, rid: %s", page, err, ctx.Kit.Rid)
			return nil, 0, err
		}
		data.Remove(metadata.PageName)
	}

	queryCond := &metadata.QueryCondition{
		Condition:      data,
		Page:           basePage,
		DisableCounter: true,
	}
	return queryCond, option.ModelBizID, nil
}

// SearchObjectAttribute search the object attributes
func (s *Service) SearchObjectAttribute(ctx *rest.Contexts) {

	queryCond, modelBizID, err := combinationSearchObjectAttrCond(ctx)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	resp, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttrByCondition(ctx.Kit.Ctx, ctx.Kit.Header, queryCond)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	grpMap, err := s.getPropertyGroupName(ctx.Kit, resp.Info, modelBizID)
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
			blog.Errorf("failed to get property group name, attr: %v, property: %s, rid: %s", attr, attr.PropertyGroup,
				ctx.Kit.Rid)
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
	// adapt input path param with bk_biz_id and attr id.
	id, bizID, err := getAttrIDAndBizID(ctx)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}

	// table field updates need to go through a separate process.
	if util.GetStrByInterface(data[common.BKPropertyTypeField]) == common.FieldTypeInnerTable {
		if err := s.updateObjectTableAttribute(ctx, id, bizID, data); err != nil {
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntity(nil)
		return
	}
	data = removeImmutableFields(data)
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.AttributeOperation().UpdateObjectAttribute(ctx.Kit, data, id, bizID, false)
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

// updateObjectTableAttribute update the table object attribute
func (s *Service) updateObjectTableAttribute(ctx *rest.Contexts, id, bizID int64, data mapstr.MapStr) error {

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err := s.Logics.AttributeOperation().UpdateTableObjectAttr(ctx.Kit, data, id, bizID)
		if err != nil {
			return err
		}
		return nil
	})

	if txnErr != nil {
		return txnErr
	}
	return nil
}

func removeImmutableFields(data mapstr.MapStr) mapstr.MapStr {
	// TODO: why does remove this????
	data.Remove(metadata.BKMetadata)
	data.Remove(common.BKAppIDField)

	// UpdateObjectAttribute should not update bk_property_index、bk_property_group
	data.Remove(common.BKPropertyIndexField)
	data.Remove(common.BKPropertyGroupField)
	return data
}

func getAttrIDAndBizID(ctx *rest.Contexts) (int64, int64, error) {

	id, err := strconv.ParseInt(ctx.Request.PathParameter("id"), 10, 64)
	if err != nil {
		blog.Errorf("failed to parse the path params id: %s, err: %v, rid: %s", ctx.Request.PathParameter("id"),
			err, ctx.Kit.Rid)
		return 0, 0, err
	}
	// adapt input path param with bk_biz_id
	var bizID int64
	if bizIDStr := ctx.Request.PathParameter(common.BKAppIDField); bizIDStr != "" {
		bizIDStr := ctx.Request.PathParameter(common.BKAppIDField)
		bizID, err = strconv.ParseInt(bizIDStr, 10, 64)
		if err != nil {
			blog.Errorf("create biz custom field, but parse biz ID failed, error: %v, rid: %s", err, ctx.Kit.Rid)
			return 0, 0, err
		}
		if bizID == 0 {
			blog.Errorf("create biz custom field, but biz ID is 0, rid: %s", ctx.Kit.Rid)
			return 0, 0, err
		}
	}
	return id, bizID, nil
}

// DeleteObjectAttribute delete the object attribute
func (s *Service) DeleteObjectAttribute(ctx *rest.Contexts) {

	kit := ctx.Kit
	paramPath := mapstr.MapStr{}
	paramPath.Set("id", ctx.Request.PathParameter("id"))
	id, err := paramPath.Int64("id")
	if err != nil {
		blog.Errorf("failed to parse the path params id: %s, err: %s , rid: %s", ctx.Request.PathParameter("id"),
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	modelType := new(ModelType)
	if err := ctx.DecodeInto(modelType); err != nil {
		ctx.RespAutoError(err)
		return
	}

	attr, err := s.getModelAttrByID(ctx.Kit, id, modelType.BizID)
	if err != nil {
		ctx.RespAutoError(err)
		return
	}
	if attr.ID == 0 {
		ctx.RespEntity(nil)
		return
	}

	if attr.TemplateID != 0 {
		ctx.RespAutoError(kit.CCError.CCErrorf(common.CCErrorTopoFieldTemplateForbiddenDeleteAttr, id, attr.TemplateID))
		return
	}

	if attr.PropertyType == common.FieldTypeInnerTable {
		if err := s.deleteTableObject(ctx.Kit, attr.ObjectID, attr.PropertyID, attr.ID); err != nil {
			ctx.RespAutoError(err)
			return
		}
		ctx.RespEntity(nil)
		return
	}
	if err = s.Logics.AttributeOperation().DeleteObjectAttribute(kit, []metadata.Attribute{attr}); err != nil {
		blog.Errorf("delete object attribute failed, id: %,bizID: %d, err: %+v, rid: %s", id, modelType.BizID,
			err, ctx.Kit.Rid)
		ctx.RespAutoError(err)
		return
	}

	if err := s.deleteHostApplyRule(ctx.Kit, id); err != nil {
		ctx.RespAutoError(err)
		return
	}

	ctx.RespEntity(nil)
}

func (s *Service) deleteHostApplyRule(kit *rest.Kit, id int64) error {
	listRuleOption := metadata.ListHostApplyRuleOption{
		AttributeIDs: []int64{id},
		Page:         metadata.BasePage{Limit: common.BKNoLimit},
	}
	ruleResult, err := s.Engine.CoreAPI.CoreService().HostApplyRule().ListHostApplyRule(kit.Ctx, kit.Header,
		0, listRuleOption)
	if err != nil {
		blog.Errorf("get host apply rule failed, listRuleOption: %+v, err: %v, rid: %s", listRuleOption, err, kit.Rid)
		return err
	}

	ruleIDs := make([]int64, 0)
	for _, item := range ruleResult.Info {
		ruleIDs = append(ruleIDs, item.ID)
	}

	if len(ruleIDs) > 0 {
		deleteRuleOption := metadata.DeleteHostApplyRuleOption{
			RuleIDs: ruleIDs,
		}
		if err := s.Engine.CoreAPI.CoreService().HostApplyRule().DeleteHostApplyRule(kit.Ctx, kit.Header, 0,
			deleteRuleOption); err != nil {
			blog.Errorf("delete host apply rule failed, params: %+v, err: %v, rid: %s", deleteRuleOption, err, kit.Rid)
			return err
		}
	}
	return nil
}

func (s *Service) getModelAttrByID(kit *rest.Kit, id, bizID int64) (metadata.Attribute, error) {

	cond := mapstr.MapStr{metadata.AttributeFieldID: id}
	util.AddModelBizIDCondition(cond, bizID)
	queryCond := &metadata.QueryCondition{
		Condition:      cond,
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}

	attrs, err := s.Engine.CoreAPI.CoreService().Model().ReadModelAttrsWithTableByCondition(kit.Ctx, kit.Header,
		bizID, queryCond)
	if err != nil {
		blog.Errorf("failed to find the attributes by the cond(%v), err: %v, rid: %s", cond, err, kit.Rid)
		return metadata.Attribute{}, nil
	}

	if len(attrs.Info) == 0 {
		blog.Errorf("not find the attributes by the cond(%v), rid: %s", cond, kit.Rid)
		return metadata.Attribute{}, nil
	}

	if len(attrs.Info) > 1 {
		blog.Errorf("the number of attributes queried is incorrect cond(%v), rid: %s", cond, kit.Rid)
		return metadata.Attribute{}, kit.CCError.CCError(common.CCErrCommGetMultipleObject)
	}
	return attrs.Info[0], nil
}

func (s *Service) deleteTableObject(kit *rest.Kit, objID, propertyID string, id int64) error {

	cond := mapstr.MapStr{common.BKObjIDField: objID}
	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(kit.Ctx, kit.Header, func() error {
		_, err := s.Logics.ObjectOperation().DeleteTableObject(kit, objID, propertyID, id)
		if err != nil {
			blog.Errorf("delete table object failed, cond: %+v, err: %v, rid: %v", cond, err, kit.Rid)
			return err
		}
		return nil
	})
	if txnErr != nil {
		return txnErr
	}
	return nil
}

// UpdateObjectAttributeIndex update object attribute index
func (s *Service) UpdateObjectAttributeIndex(ctx *rest.Contexts) {
	data := new(metadata.UpdateAttrIndexInput)
	if err := ctx.DecodeInto(data); err != nil {
		ctx.RespAutoError(err)
		return
	}

	objID := ctx.Request.PathParameter(common.BKObjIDField)
	idStr := ctx.Request.PathParameter("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		blog.Errorf("parse id from path params failed, err: %v, id: %s, rid: %s", err, idStr, ctx.Kit.Rid)
		ctx.RespAutoError(ctx.Kit.CCError.CCError(common.CCErrorTopoPathParamPaserFailed))
		return
	}

	txnErr := s.Engine.CoreAPI.CoreService().Txn().AutoRunTxn(ctx.Kit.Ctx, ctx.Kit.Header, func() error {
		err = s.Engine.CoreAPI.CoreService().Model().UpdateModelAttrIndex(ctx.Kit.Ctx, ctx.Kit.Header, objID,
			id, data)
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
	ctx.RespEntity(nil)
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

	grpMap, err := s.getPropertyGroupName(ctx.Kit, result.Info, dataWithModelBizID.ModelBizID)
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

func (s *Service) getPropertyGroupName(kit *rest.Kit, attrs []metadata.Attribute,
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
	rsp, err := s.Engine.CoreAPI.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, cond)
	if err != nil {
		blog.Errorf("failed to get attr group, err: %s, rid: %s", err.Error(), kit.Rid)
		return nil, err
	}

	grpMap := make(map[string]string)
	for _, grp := range rsp.Info {
		grpMap[grp.GroupID] = grp.GroupName
	}

	return grpMap, nil
}

// checkJsonTagContainIsMultipleField verify whether the ismultiple field exists
// 当创建组织字段属性时，前端的默认行为为多选，ismultiple参数为true. 为了和前端保持一致的动作，通过api接口创建时组织字段时，
// 在用户没有传ismultiple字段时，需要默认给ismultiple置为true
func checkJsonTagContainIsMultipleField(data interface{}) bool {
	typeOfOption := reflect.TypeOf(data)
	valueOfOption := reflect.ValueOf(data)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tagTmp := typeOfOption.Field(i).Tag.Get("json")
		tags := strings.Split(tagTmp, ",")

		if tags[0] == "" {
			continue
		}

		if tags[0] == common.BKIsMultipleField {
			fieldValue := valueOfOption.Field(i)
			if fieldValue.IsNil() {
				return false
			}
		}
	}

	return true
}
