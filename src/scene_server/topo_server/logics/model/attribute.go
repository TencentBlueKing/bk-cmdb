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
	"unicode/utf8"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/mapstruct"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(kit *rest.Kit, data *metadata.Attribute) (*metadata.Attribute, error)
	DeleteObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) error
	UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64) error
	UpdateObjectAttributeIndex(kit *rest.Kit, objID string, data mapstr.MapStr,
		attID int64) (*metadata.UpdateAttrIndexData, error)
}

// NewAttributeOperation create a new attribute operation instance
func NewAttributeOperation(client apimachinery.ClientSetInterface,
	authManager *extensions.AuthManager) AttributeOperationInterface {
	return &attribute{
		clientSet:   client,
		authManager: authManager,
	}
}

type attribute struct {
	lang        language.DefaultCCLanguageIf
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// IsValid check is valid
func (a *attribute) IsValid(kit *rest.Kit, isUpdate bool, data *metadata.Attribute) error {

	if data.PropertyID == common.BKInstParentStr {
		return nil
	}

	// check if property type for creation is valid, can't update property type
	if !isUpdate {
		if data.PropertyType == "" {
			return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
		}
	}

	if !isUpdate || data.ToMapStr().Exists(metadata.AttributeFieldPropertyID) {
		if common.AttributeIDMaxLength < utf8.RuneCountInString(data.PropertyID) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.Language("model_attr_bk_property_id"), common.AttributeIDMaxLength)
		}
		match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, data.PropertyID)
		if err != nil {
			return err
		}

		if !match {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, data.PropertyID)
		}
	}

	if !isUpdate || data.ToMapStr().Exists(metadata.AttributeFieldPropertyName) {
		if common.AttributeNameMaxLength < utf8.RuneCountInString(data.PropertyName) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.Language("model_attr_bk_property_name"), common.AttributeNameMaxLength)
		}
	}

	// check option validity for creation,
	// update validation is in coreservice cause property type need to be obtained from db
	if !isUpdate {
		if a.isPropertyTypeIntEnumListSingleLong(data.PropertyType) {
			if err := util.ValidPropertyOption(data.PropertyType, data.Option, kit.CCError); nil != err {
				return err
			}
		}
	}

	if data.Placeholder != "" {
		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(data.Placeholder) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, a.lang.Language("model_attr_placeholder"),
				common.AttributePlaceHolderMaxLength)
		}
	}

	return nil
}

// isPropertyTypeIntEnumListSingleLong check is property type in enum list single long
func (a *attribute) isPropertyTypeIntEnumListSingleLong(propertyType string) bool {
	switch propertyType {
	case common.FieldTypeInt, common.FieldTypeEnum, common.FieldTypeList:
		return true
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		return true
	default:
		return false
	}
}

// isObjExists 检查当前objID在数据库中是否存在
func (a *attribute) isObjExists(kit *rest.Kit, objID string) error {
	checkObjCond := mapstr.MapStr{
		common.BKObjIDField: objID,
	}
	cond := make([]map[string]interface{}, 0)
	cond = append(cond, checkObjCond)

	resp, e := a.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNameObjDes, cond)
	if e != nil {
		blog.Errorf("failed to check the object repeated, err: %s, rid: %s", e, kit.Rid)
		return e
	}

	if resp[0] > 0 {
		return fmt.Errorf("check the object repeated, objID:%s, rid:%s", objID, kit.Rid)
	}

	return nil
}

// CreateObjectGroup create object groupdataValidationFormulaStrLen
func (a *attribute) CreateObjectGroup(kit *rest.Kit, data mapstr.MapStr) (*metadata.Group, error) {
	grp := metadata.Group{}

	err := mapstr.SetValueToStructByTags(&grp, data)
	if nil != err {
		blog.Errorf("failed to parse the group data(%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	//  check the object
	if err = a.isObjExists(kit, grp.ObjectID); nil != err {
		blog.Errorf("the group (%#v) is in valid, rid: %s", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// create a new group
	rsp, err := a.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header, grp.ObjectID,
		metadata.CreateModelAttributeGroup{Data: grp})
	if err != nil {
		blog.Errorf("failed to save the group data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	grp.ID = int64(rsp.Created.ID)

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, grp.ID, &grp)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but generate audit log failed, err: %v, rid: %s",
			grp.GroupName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute group %s success, but save audit log failed, err: %v, rid: %s",
			grp.GroupName, err, kit.Rid)
		return nil, err
	}

	return &grp, nil
}

// CreateObjectAttribute create object attribute
func (a *attribute) CreateObjectAttribute(kit *rest.Kit, data *metadata.Attribute) (*metadata.Attribute, error) {
	if data.IsOnly {
		data.IsRequired = true
	}

	if len(data.PropertyGroup) == 0 {
		data.PropertyGroup = "default"
	}

	// check if the object is mainline object, if yes. then user can not create required attribute.
	yes, err := a.isMainlineModel(kit, data.ObjectID)
	if err != nil {
		blog.Errorf("not allow to add required attribute to mainline object: %+v. "+"rid: %d.", data, kit.Rid)
		return nil, err
	}

	if yes && data.IsRequired {
		return nil, kit.CCError.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
	}

	// check the object id
	if err = a.isObjExists(kit, data.ObjectID); err != nil {
		return nil, err
	}

	filters := make([]map[string]interface{}, 0)
	filters = append(filters, mapstr.MapStr{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: data.PropertyGroup,
		common.BKAppIDField:           data.BizID,
	}, mapstr.MapStr{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: common.BKBizDefault,
		common.BKAppIDField:           data.BizID,
	})

	resp, e := a.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header, common.BKTableNamePropertyGroup,
		filters)
	if e != nil {
		blog.Errorf("get property group failed, filters: %+v., err: %s, rid: %d.", data, e, kit.Rid)
		return nil, e
	}

	if resp[0] == 0 && data.BizID > 0 {
		if resp[1] == 0 {
			group := metadata.Group{
				IsDefault:  true,
				GroupIndex: -1,
				GroupName:  common.BKBizDefault,
				GroupID:    common.BKBizDefault,
				ObjectID:   data.ObjectID,
				OwnerID:    data.OwnerID,
				BizID:      data.BizID,
			}
			// TODO 替换依赖直接传递&group参数
			if _, err := a.CreateObjectGroup(kit, mapstr.MapStr{"field": group}); err != nil {
				blog.Errorf("failed to create the default group, err: %s, rid: %s", err, kit.Rid)
				return nil, err
			}
		}
		data.PropertyGroup = common.BKBizDefault
	} else {
		data.PropertyGroup = common.BKDefaultField
	}

	if err := a.IsValid(kit, false, data); nil != err {
		return nil, err
	}

	// check the property id repeated
	data.OwnerID = kit.SupplierAccount

	input := metadata.CreateModelAttributes{Attributes: []metadata.Attribute{*data}}
	if _, err := a.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, data.ObjectID,
		&input); err != nil {
		blog.Errorf("failed to create model attrs, err: %s, input: %s, rid: %s", err, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, data.ID, nil)
	if err != nil {
		blog.Errorf("create object attribute %s success, but generate audit log failed, err: %v, rid: %s",
			data.PropertyName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute %s success, but save audit log failed, err: %v, rid: %s",
			data.PropertyName, err, kit.Rid)
		return nil, err
	}

	return data, nil
}

// DeleteObjectAttribute delete object attribute
func (a *attribute) DeleteObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, modelBizID int64) error {
	attr := new(metadata.Attribute)
	if err := mapstruct.Decode2Struct(cond, attr); err != nil {
		blog.Errorf("unmarshal mapstr data into module failed, module: %s, err: %s, rid: %s", cond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}
	if err := a.IsValid(kit, false, attr); nil != err {
		return err
	}

	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)

	util.AddModelBizIDCondition(cond, modelBizID)
	queryCond := &metadata.QueryCondition{
		Condition: cond,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
		DisableCounter: true,
	}
	attrItems, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("failed to find the attributes by the cond(%v), err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectAttributeDeleteFailed, err.Error())
	}

	auditLogArr := make([]metadata.AuditLog, 0)
	attrID := make([]string, 0)
	var objID string
	for _, attrItem := range attrItems.Info {
		// generate audit log of model attribute.
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attrItem.ID, &attrItem)
		if err != nil {
			blog.Errorf("generate audit log failed, model attribute %s, err: %v, rid: %s", attrItem.PropertyName,
				err, kit.Rid)
			return err
		}
		auditLogArr = append(auditLogArr, *auditLog)
		attrID = append(attrID, attrItem.PropertyID)
		objID = attrItem.ObjectID
	}

	deleteCond := &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKDBIN: attrID,
		},
	}
	// delete the attribute.
	if _, err := a.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, objID, deleteCond); err != nil {
		blog.Errorf("delete object attribute failed, err: %v, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLogArr...); err != nil {
		blog.Errorf("delete object attribute success, but save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	return nil
}

// UpdateObjectAttribute update object attribute
func (a *attribute) UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64) error {
	cond := make(map[string]interface{})
	util.AddModelBizIDCondition(cond, modelBizID)

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update model attribute, attID: %d, err: %v, rid: %s",
			attID, err, kit.Rid)
		return err
	}

	// to update.
	cond[common.BKFieldID] = attID
	input := metadata.UpdateOption{
		Condition: cond,
		Data:      data,
	}
	if _, err := a.clientSet.CoreService().Model().UpdateModelAttrsByCondition(kit.Ctx, kit.Header, &input); err != nil {
		blog.Errorf("failed to update module attr, err: %s, rid: %s", err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object attribute success, but save audit log failed, attID: %d, err: %v, rid: %s",
			attID, err, kit.Rid)
		return err
	}

	return nil
}

// UpdateObjectAttributeIndex update object attribute index by obj id
func (a *attribute) UpdateObjectAttributeIndex(kit *rest.Kit, objID string, data mapstr.MapStr,
	attID int64) (*metadata.UpdateAttrIndexData, error) {
	input := metadata.UpdateOption{
		Condition: mapstr.MapStr{common.BKFieldID: attID},
		Data:      data,
	}

	rsp, err := a.clientSet.CoreService().Model().UpdateModelAttrsIndex(kit.Ctx, kit.Header, objID, &input)
	if err != nil {
		blog.Errorf("failed to update module attr index, err: %s, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return rsp, nil
}

// isMainlineModel check is mainline model by module id
func (a *attribute) isMainlineModel(kit *rest.Kit, modelID string) (bool, error) {
	cond := mapstr.MapStr{
		common.AssociationKindIDField: common.AssociationKindMainline,
	}
	queryCond := &metadata.QueryCondition{
		Condition:      cond,
		DisableCounter: true,
	}
	asst, err := a.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		return false, err
	}

	for _, mainline := range asst.Info {
		if mainline.ObjectID == modelID {
			return true, nil
		}
	}

	return false, nil
}
