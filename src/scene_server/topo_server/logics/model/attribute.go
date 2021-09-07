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
	SetProxy(grp GroupOperationInterface, obj ObjectOperationInterface)
}

// NewAttributeOperation create a new attribute operation instance
func NewAttributeOperation(client apimachinery.ClientSetInterface, authManager *extensions.AuthManager,
	languageIf language.CCLanguageIf) AttributeOperationInterface {
	return &attribute{
		clientSet:   client,
		authManager: authManager,
		lang:        languageIf,
	}
}

type attribute struct {
	lang        language.CCLanguageIf
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
	grp         GroupOperationInterface
	obj         ObjectOperationInterface
}

// SetProxy SetProxy
func (a *attribute) SetProxy(grp GroupOperationInterface, obj ObjectOperationInterface) {
	a.grp = grp
	a.obj = obj
}

// IsValid check is valid
func (a *attribute) IsValid(kit *rest.Kit, isUpdate bool, data *metadata.Attribute) error {
	if data.PropertyID == common.BKInstParentStr {
		return nil
	}

	// check if property type for creation is valid, can't update property type
	if !isUpdate && data.PropertyType == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
	}

	if !isUpdate || data.PropertyID != "" {
		if common.AttributeIDMaxLength < utf8.RuneCountInString(data.PropertyID) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Language(
					"model_attr_bk_property_id"), common.AttributeIDMaxLength)
		}
		match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, data.PropertyID)
		if err != nil {
			return err
		}

		if !match {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, data.PropertyID)
		}
	}

	if !isUpdate || data.PropertyName != "" {
		if common.AttributeNameMaxLength < utf8.RuneCountInString(data.PropertyName) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Language(
					"model_attr_bk_property_name"), common.AttributeNameMaxLength)
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
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Language("model_attr_placeholder"),
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

// checkAttributeGroupExist check attribute group exist, not exist create default group
func (a *attribute) checkAttributeGroupExist(kit *rest.Kit, data *metadata.Attribute) error {
	cond := mapstr.MapStr{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: data.PropertyGroup,
	}
	groupResult, err := a.grp.FindGroupByObject(kit, data.ObjectID, cond, data.BizID)
	if err != nil {
		blog.Errorf("failed to search the attribute group data (%#v), err: %s, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(groupResult) > 0 {
		data.PropertyGroupName = groupResult[0].GroupName
		return nil
	}

	// create the default group
	if data.BizID > 0 {
		cond := mapstr.MapStr{
			common.BKObjIDField:           data.ObjectID,
			common.BKPropertyGroupIDField: common.BKBizDefault,
		}
		groupResult, err := a.grp.FindGroupByObject(kit, data.ObjectID, cond, data.BizID)
		if err != nil {
			blog.Errorf("failed to search the attr group data: %#v, err: %s, rid: %s", cond, err, kit.Rid)
			return err
		}
		if len(groupResult) == 0 {
			group := metadata.Group{
				IsDefault:  true,
				GroupIndex: -1,
				GroupName:  common.BKBizDefault,
				GroupID:    common.BKBizDefault,
				ObjectID:   data.ObjectID,
				OwnerID:    data.OwnerID,
				BizID:      data.BizID,
			}

			if _, err := a.grp.CreateObjectGroup(kit, &group); err != nil {
				blog.Errorf("failed to create the default group, err: %s, rid: %s", err, kit.Rid)
				return err
			}
			data.PropertyGroup = common.BKBizDefault
			data.PropertyGroupName = common.BKBizDefault
		} else {
			data.PropertyGroup = common.BKDefaultField
			data.PropertyGroupName = common.BKDefaultField
		}
	}

	return nil
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
	exist, err := a.obj.IsObjectExist(kit, data.ObjectID)
	if err != nil {
		return nil, err
	}
	if !exist {
		blog.Errorf("obj id is not valid, obj id: %+v, rid: %d.", data.ObjectID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	if err = a.checkAttributeGroupExist(kit, data); err != nil {
		blog.Errorf("failed to create the default group, err: %s, rid: %s", err, kit.Rid)
		return nil, err
	}

	if err := a.IsValid(kit, false, data); err != nil {
		return nil, err
	}

	data.OwnerID = kit.SupplierAccount
	input := metadata.CreateModelAttributes{Attributes: []metadata.Attribute{*data}}
	resp, err := a.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, data.ObjectID, &input)
	if err != nil {
		blog.Errorf("failed to create model attrs, err: %v, input: %#v, rid: %s", err, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}
	for _, exception := range resp.Exceptions {
		return nil, kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(resp.Repeated) > 0 {
		blog.Errorf("create model attrs failed, the attr is duplicated, ObjectID: %s, input: %s, rid: %s",
			data.ObjectID, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	if len(resp.Created) != 1 {
		blog.Errorf("create model attrs created amount error, ObjectID: %s, input: %s, rid: %s", data.ObjectID,
			input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}
	data.ID = int64(resp.Created[0].ID)

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, int64(resp.Created[0].ID), nil)
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
	util.AddModelBizIDCondition(cond, modelBizID)
	queryCond := &metadata.QueryCondition{
		Condition: cond,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	attrItems, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("failed to find the attributes by the cond(%v), err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(attrItems.Info) == 0 {
		blog.Errorf("not find the attributes by the cond(%v), rid: %s", cond, kit.Rid)
		return nil
	}

	auditLogArr := make([]metadata.AuditLog, 0)
	attrID := make([]int64, 0)
	var objID string
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	for _, attrItem := range attrItems.Info {
		// generate audit log of model attribute.
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attrItem.ID, &attrItem)
		if err != nil {
			blog.Errorf("generate audit log failed, model attribute %s, err: %v, rid: %s", attrItem.PropertyName,
				err, kit.Rid)
			return err
		}
		auditLogArr = append(auditLogArr, *auditLog)
		attrID = append(attrID, attrItem.ID)
		objID = attrItem.ObjectID
	}

	// delete the attribute.
	deleteCond := &metadata.DeleteOption{
		Condition: mapstr.MapStr{
			common.BKFieldID: mapstr.MapStr{common.BKDBIN: attrID},
		},
	}
	_, err = a.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, objID, deleteCond)
	if err != nil {
		blog.Errorf("delete object attribute failed, err: %v, rid: %s", err, kit.Rid)
		return err
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

	attr := new(metadata.Attribute)
	if err := mapstruct.Decode2Struct(data, attr); err != nil {
		blog.Errorf("unmarshal mapstr data into module failed, module: %s, err: %s, rid: %s", attr, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}

	if err := a.IsValid(kit, true, attr); err != nil {
		return err
	}

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
	cond := mapstr.MapStr{
		common.BKFieldID: attID,
	}
	util.AddModelBizIDCondition(cond, modelBizID)
	input := metadata.UpdateOption{
		Condition: cond,
		Data:      data,
	}
	if _, err := a.clientSet.CoreService().Model().UpdateModelAttrsByCondition(kit.Ctx, kit.Header, &input); err != nil {
		blog.Errorf("failed to update module attr, err: %s, rid: %s", err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("update object attribute success, but save audit log failed, attID: %d, err: %v, rid: %s",
			attID, err, kit.Rid)
		return err
	}

	return nil
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

	if len(asst.Info) <= 0 {
		return false, fmt.Errorf("model association [%+v] not found", cond)
	}

	for _, mainline := range asst.Info {
		if mainline.ObjectID == modelID {
			return true, nil
		}
	}

	return false, nil
}
