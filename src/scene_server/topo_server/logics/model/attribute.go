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
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"configcenter/src/ac/extensions"
	"configcenter/src/apimachinery"
	"configcenter/src/common"
	"configcenter/src/common/auditlog"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
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
	FindObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, objID string) ([]metadata.Attribute, error)
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
	FieldValid
	clientSet   apimachinery.ClientSetInterface
	authManager *extensions.AuthManager
}

// FieldValid field valid method
type FieldValid struct {
	lang language.DefaultCCLanguageIf
}

// Valid valid the field
func (f *FieldValid) Valid(kit *rest.Kit, data mapstr.MapStr, fieldID string) (string, error) {

	val, err := data.String(fieldID)
	if nil != err {
		return val, kit.CCError.New(common.CCErrCommParamsIsInvalid, fieldID+" "+err.Error())
	}
	if 0 == len(val) {
		return val, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, fieldID)
	}

	return val, nil
}

// ValidID check the property ID
func (f *FieldValid) ValidID(kit *rest.Kit, value string) error {
	if common.AttributeIDMaxLength < utf8.RuneCountInString(value) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			f.lang.Language("model_attr_bk_property_id"), common.AttributeIDMaxLength)
	}
	match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, value)
	if nil != err {
		return err
	}

	if !match {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, value)
	}

	return nil
}

// ValidName check the name
func (f *FieldValid) ValidName(kit *rest.Kit, value string) error {
	if common.AttributeNameMaxLength < utf8.RuneCountInString(value) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			f.lang.Language("model_attr_bk_property_name"), common.AttributeNameMaxLength)
	}
	value = strings.TrimSpace(value)
	return nil
}

// ValidPlaceHolder check the PlaceHolder
func (f *FieldValid) ValidPlaceHolder(kit *rest.Kit, value string) error {
	if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(value) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			f.lang.Language("model_attr_placeholder"), common.AttributePlaceHolderMaxLength)
	}
	return nil
}

// IsValid check is valid
func (a *attribute) IsValid(kit *rest.Kit, isUpdate bool, data *metadata.Attribute) error {

	if data.PropertyID == common.BKInstParentStr {
		return nil
	}

	// check if property type for creation is valid, can't update property type
	if !isUpdate {
		if _, err := a.FieldValid.Valid(kit, data.ToMapStr(), metadata.AttributeFieldPropertyType); nil != err {
			return err
		}
	}

	if !isUpdate || data.ToMapStr().Exists(metadata.AttributeFieldPropertyID) {
		val, err := a.FieldValid.Valid(kit, data.ToMapStr(), metadata.AttributeFieldPropertyID)
		if nil != err {
			return err
		}
		if err = a.FieldValid.ValidID(kit, val); nil != err {
			return err
		}
	}

	if !isUpdate || data.ToMapStr().Exists(metadata.AttributeFieldPropertyName) {
		val, err := a.FieldValid.Valid(kit, data.ToMapStr(), metadata.AttributeFieldPropertyName)
		if nil != err {
			return err
		}
		if err = a.FieldValid.ValidName(kit, val); nil != err {
			return err
		}
	}

	// check option validity for creation,
	// update validation is in coreservice cause property type need to be obtained from db
	if !isUpdate {
		propertyType, err := data.ToMapStr().String(metadata.AttributeFieldPropertyType)
		if nil != err {
			return kit.CCError.New(common.CCErrCommParamsIsInvalid, err.Error())
		}

		option, exists := data.ToMapStr().Get(metadata.AttributeFieldOption)
		if exists && a.isPropertyTypeIntEnumListSingleLong(propertyType) {
			if err := util.ValidPropertyOption(propertyType, option, kit.CCError); nil != err {
				return err
			}
		}
	}

	if val, ok := data.ToMapStr()[metadata.AttributeFieldPlaceHolder]; ok && val != "" {
		if placeholder, ok := val.(string); ok {
			if err := a.FieldValid.ValidPlaceHolder(kit, placeholder); nil != err {
				return err
			}
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

// checkObjectExist 检查当前objID在数据库中是否存在
func (a *attribute) checkObjectExist(kit *rest.Kit, objID string) error {

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
		blog.Errorf("[operation-grp] failed to parse the group data(%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, err
	}

	//  check the object
	if err = a.checkObjectExist(kit, grp.ObjectID); nil != err {
		blog.Errorf("[operation-grp] the group (%#v) is in valid, rid: %s", data, kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}

	// create a new group
	rsp, err := a.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header, grp.ObjectID,
		metadata.CreateModelAttributeGroup{Data: grp})
	if nil != err {
		blog.Errorf("[operation-grp] failed to save the group data (%#v), error info is %s, rid: %s", data, err.Error(), kit.Rid)
		return nil, kit.CCError.New(common.CCErrTopoObjectGroupCreateFailed, err.Error())
	}
	if !rsp.Result {
		blog.Errorf("[model-grp] failed to create the group(%s), err: is %s, rid: %s", grp.GroupID, rsp.ErrMsg, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectGroupCreateFailed)
	}

	grp.ID = int64(rsp.Data.Created.ID)

	// generate audit log of object attribute group.
	audit := auditlog.NewAttributeGroupAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, grp.ID, &grp)
	if err != nil {
		blog.Errorf("create object attribute group %s success, but generate audit log failed, err: %v, rid: %s", grp.GroupName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err = audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("create object attribute group %s success, but save audit log failed, err: %v, rid: %s", grp.GroupName, err, kit.Rid)
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
		return nil, kit.CCError.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	if yes {
		if data.IsRequired {
			return nil, kit.CCError.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
		}
	}

	// check the object id
	if err = a.checkObjectExist(kit, data.ObjectID); err != nil {
		return nil, kit.CCError.New(common.CCErrTopoObjectAttributeCreateFailed, err.Error())
	}

	cond := mapstr.MapStr{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: data.PropertyGroup,
		common.BKAppIDField:           data.BizID,
	}
	queryCond := metadata.QueryCondition{
		Condition:      cond,
		DisableCounter: true,
	}
	groupResult, err := a.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header, queryCond)
	if nil != err {
		blog.Errorf("failed to request the attr group, err: %s, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	// create the default group
	if len(groupResult.Data.Info) == 0 {
		if data.BizID > 0 {
			cond := mapstr.MapStr{
				common.BKObjIDField:           data.ObjectID,
				common.BKPropertyGroupIDField: common.BKBizDefault,
			}
			queryCond := metadata.QueryCondition{
				Condition:      cond,
				DisableCounter: true,
			}
			groupResult, err := a.clientSet.CoreService().Model().ReadAttributeGroupByCondition(kit.Ctx, kit.Header,
				queryCond)
			if nil != err {
				blog.Errorf("failed to request the attr group, err: %s, rid: %s", err, kit.Rid)
				return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
			}
			if len(groupResult.Data.Info) == 0 {
				group := metadata.Group{
					IsDefault:  true,
					GroupIndex: -1,
					GroupName:  common.BKBizDefault,
					GroupID:    common.BKBizDefault,
					ObjectID:   data.ObjectID,
					OwnerID:    data.OwnerID,
					BizID:      data.BizID,
				}

				if _, err := a.CreateObjectGroup(kit, mapstr.MapStr{"field": group}); err != nil {
					blog.Errorf("failed to create the default group, err: %s, rid: %s", err, kit.Rid)
					return nil, kit.CCError.Error(common.CCErrTopoObjectGroupCreateFailed)
				}
			}
			data.PropertyGroup = common.BKBizDefault
		} else {
			data.PropertyGroup = common.BKDefaultField
		}
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
	cond[common.BKAppIDField] = modelBizID
	attrItems, err := a.FindObjectAttribute(kit, cond, "")
	if nil != err {
		blog.Errorf("failed to find the attributes by the cond(%v), err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.New(common.CCErrTopoObjectAttributeDeleteFailed, err.Error())
	}

	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)

	deleteCond, exist := cond.Get(metadata.AttributeFieldID)
	if !exist {
		blog.Errorf("failed to get delete cond, err: %v, rid: %s", cond, err, kit.Rid)
		return errors.New("get attribute field id failed")
	}

	for _, attrItem := range attrItems {
		// generate audit log of model attribute.
		auditLog, err := audit.GenerateAuditLog(generateAuditParameter, attrItem.ID, &attrItem)
		if err != nil {
			blog.Errorf("generate audit log failed, model attribute %s, err: %v, rid: %s", attrItem.PropertyName,
				err, kit.Rid)
			return err
		}

		// delete the attribute.
		if _, err := a.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, attrItem.ObjectID,
			&metadata.DeleteOption{Condition: mapstr.MapStr{metadata.AttributeFieldID: deleteCond}}); err != nil {
			blog.Errorf("delete object attribute failed, err: %v, rid: %s", err, kit.Rid)
			return kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
		}

		// save audit log.
		if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
			blog.Errorf("delete object attribute %s success, but save audit log failed, err: %v, rid: %s",
				attrItem.PropertyName, err, kit.Rid)
			return err
		}
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
	if nil != err {
		blog.Errorf("failed to update module attr index, err: %s, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return rsp.Data, nil
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

	for _, mainline := range asst.Data.Info {
		if mainline.ObjectID == modelID {
			return true, nil
		}
	}

	return false, nil
}

// FindObjectAttribute find attribute by obj id
func (a *attribute) FindObjectAttribute(kit *rest.Kit, cond mapstr.MapStr, objID string) ([]metadata.Attribute, error) {
	var limit, start int64
	var sort string
	if cond.Exists(metadata.PageName) {
		page, err := cond.MapStr(metadata.PageName)
		if err != nil {
			blog.Errorf("page info convert to mapstr failed, page: %v, err: %v, rid: %s", cond[metadata.PageName],
				err, kit.Rid)
			return nil, err
		}

		limit, err = page.Int64(metadata.PageLimit)
		if err != nil {
			blog.Errorf("get limit from page failed, page: %v, err: %v, rid: %s", page, err, kit.Rid)
			return nil, err
		}

		start, err = page.Int64(metadata.PageStart)
		if err != nil {
			blog.Errorf("get start from page failed, page: %v, err: %v, rid: %s", page, err, kit.Rid)
			return nil, err
		}

		s, exist := page.Get(metadata.PageSort)
		if !exist {
			sort = common.BKFieldID
		}
		sort = s.(string)
		cond.Remove(metadata.PageName)
	}

	opt := &metadata.QueryCondition{
		Condition:      cond,
		DisableCounter: true,
		Page:           metadata.BasePage{Limit: int(limit), Start: int(start), Sort: sort},
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, opt)
	if err != nil {
		blog.Errorf("find business attributes failed, err: %v, rid: %s", err, kit.Rid)
		return nil, kit.CCError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	return resp.Data.Info, nil
}
