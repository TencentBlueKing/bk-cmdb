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
	"configcenter/src/scene_server/topo_server/core/operation"

	"github.com/rs/xid"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(kit *rest.Kit, data *metadata.Attribute) (*metadata.Attribute, error)
	CreateObjectAttributeBatch(kit *rest.Kit, data map[string]operation.ImportObjectData) (mapstr.MapStr, error)
	FindObjectAttributeBatch(kit *rest.Kit, objIDs []string) (mapstr.MapStr, error)
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

type rowInfo struct {
	Row  int64  `json:"row"`
	Info string `json:"info"`
	// value can empty, eg:parse error
	PropertyID string `json:"bk_property_id"`
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
	if !isUpdate && a.isPropertyTypeIntEnumListSingleLong(data.PropertyType) {
		if err := util.ValidPropertyOption(data.PropertyType, data.Option, kit.CCError); nil != err {
			return err
		}
	}

	if data.Placeholder != "" && common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(data.Placeholder) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
			a.lang.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Language("model_attr_placeholder"),
			common.AttributePlaceHolderMaxLength)
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
		blog.Errorf("failed to search the attribute group data (%#v), err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(groupResult) > 0 {
		data.PropertyGroupName = groupResult[0].GroupName
		return nil
	}

	if data.BizID == 0 {
		data.PropertyGroup = common.BKDefaultField
		data.PropertyGroupName = common.BKDefaultField
		return nil
	}

	// create the biz default group
	bizDefaultGroupCond := mapstr.MapStr{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: common.BKBizDefault,
	}
	bizDefaultGroupResult, err := a.grp.FindGroupByObject(kit, data.ObjectID, bizDefaultGroupCond, data.BizID)
	if err != nil {
		blog.Errorf("failed to search the attr group data: %#v, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(bizDefaultGroupResult) == 0 {
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
			blog.Errorf("failed to create the default group, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	}

	data.PropertyGroup = common.BKBizDefault
	data.PropertyGroupName = common.BKBizDefault
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
		blog.Errorf("obj id is not exist, obj id: %s, rid: %s", data.ObjectID, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	if err = a.checkAttributeGroupExist(kit, data); err != nil {
		blog.Errorf("failed to create the default group, err: %v, rid: %s", err, kit.Rid)
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

// setCreateObjectBatchObjResult set create object batch obj result
func (a *attribute) setCreateObjectBatchObjResult(objID string, result mapstr.MapStr, itemErr, addErr, setErr,
	succInfo []rowInfo) {
	subResult := mapstr.New()
	if len(itemErr) > 0 {
		subResult["errors"] = itemErr
	}
	if len(addErr) > 0 {
		subResult["insert_failed"] = addErr
	}
	if len(setErr) > 0 {
		subResult["update_failed"] = setErr
	}
	if len(succInfo) > 0 {
		subResult["success"] = succInfo
	}
	if len(subResult) > 0 {
		result[objID] = subResult
	}
}

// CreateObjectAttributeBatch manipulate model in cc_ObjDes
// this method does'nt act as it's name, it create or update model's attributes indeed.
// it only operate on model already exist, that it to say no new model will be created.
func (a *attribute) CreateObjectAttributeBatch(kit *rest.Kit, inputDataMap map[string]operation.ImportObjectData) (
	mapstr.MapStr, error) {
	result := make(mapstr.MapStr)
	hasError := false

	for objID, inputData := range inputDataMap {
		subResult := make(mapstr.MapStr)
		exist, err := a.obj.IsObjectExist(kit, objID)
		if !exist {
			blog.Errorf("create model patch, obj id not exist, obj id: %s, rid: %s", objID, kit.Rid)
			continue
		}
		if err != nil {
			blog.Errorf("create model patch, obj id not exist, obj id: %s, err: %v, rid: %s", objID, err, kit.Rid)
			subResult["error"] = fmt.Sprintf("the obj id: %s is invalid", objID)
			result[objID] = subResult
			hasError = true
			continue
		}

		// 异常错误，比如接卸该行数据失败， 查询数据失败
		var itemErr []rowInfo
		// 新加属性失败
		var addErr []rowInfo
		// 更新属性失败
		var setErr []rowInfo
		// 成功数据的信息
		var succInfo []rowInfo

		// update the object's attribute
		for idx, attr := range inputData.Attr {
			targetAttr := new(metadata.Attribute)
			if err := mapstruct.Decode2Struct(attr, targetAttr); err != nil {
				blog.Errorf("unmarshal mapstr data into module failed, module: %s, err: %s, rid: %s", attr, err,
					kit.Rid)
				itemErr = append(itemErr, rowInfo{Row: idx, Info: err.Error()})
				hasError = true
				continue
			}
			targetAttr.OwnerID = kit.SupplierAccount
			targetAttr.ObjectID = objID

			if targetAttr.PropertyID == common.BKInstParentStr {
				continue
			}

			if len(targetAttr.PropertyGroupName) == 0 {
				targetAttr.PropertyGroup = "Default"
			}

			// find group
			itemErr, setErr, err = a.findGroupExist(kit, objID, targetAttr, itemErr, setErr, idx)
			if err != nil {
				hasError = true
				continue
			}

			// create or update the attribute
			addErr, setErr, err = a.createOrUpdateAttr(kit, objID, targetAttr, addErr, setErr, idx)
			if err != nil {
				hasError = true
				continue
			}

			succInfo = append(succInfo, rowInfo{Row: idx, Info: "", PropertyID: targetAttr.PropertyID})
		}

		// 将需要返回的信息更新到result中。 这个函数会修改result参数的值
		a.setCreateObjectBatchObjResult(objID, result, itemErr, addErr, setErr, succInfo)
	}

	if hasError {
		return result, kit.CCError.Error(common.CCErrCommNotAllSuccess)
	}

	return result, nil
}

//
func (a *attribute) createOrUpdateAttr(kit *rest.Kit, objID string, targetAttr *metadata.Attribute, addErr,
	setErr []rowInfo, idx int64) ([]rowInfo, []rowInfo, error) {
	attrCond := mapstr.MapStr{
		metadata.AttributeFieldObjectID:   objID,
		metadata.AttributeFieldPropertyID: targetAttr.ID,
	}
	util.AddModelBizIDCondition(attrCond, targetAttr.BizID)
	queryCond := &metadata.QueryCondition{
		Condition: attrCond,
		Page: metadata.BasePage{
			Limit: common.BKNoLimit,
		},
	}
	attrs, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		addErr = append(addErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
		return addErr, setErr, err
	}

	if len(attrs.Info) == 0 {
		if _, err = a.CreateObjectAttribute(kit, targetAttr); err != nil {
			addErr = append(addErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
			return addErr, setErr, err
		}
	}

	for _, attr := range attrs.Info {
		if err := a.UpdateObjectAttribute(kit, attr.ToMapStr(), attr.ID, attr.BizID); nil != err {
			setErr = append(setErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
			return addErr, setErr, err
		}
	}

	return addErr, setErr, nil
}

// findGroupExist if group not exist, create one group
func (a *attribute) findGroupExist(kit *rest.Kit, objID string, targetAttr *metadata.Attribute, itemErr,
	setErr []rowInfo, idx int64) ([]rowInfo, []rowInfo, error) {
	grpCond := mapstr.MapStr{
		metadata.GroupFieldObjectID:  objID,
		metadata.GroupFieldGroupName: targetAttr.PropertyGroupName,
	}
	grps, err := a.grp.FindGroupByObject(kit, objID, grpCond, targetAttr.BizID)
	if err != nil {
		blog.Errorf("find object group failed, obj id: %s, group name: %s, err: %v, rid: %s", objID,
			targetAttr.PropertyGroupName, err, kit.Rid)
		itemErr = append(itemErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
		return itemErr, setErr, err
	}

	if len(grps) > 0 {
		targetAttr.PropertyGroup = grps[0].GroupID // should be only one group
	} else {
		g := metadata.Group{
			GroupName: targetAttr.PropertyGroupName,
			GroupID:   xid.New().String(),
			ObjectID:  objID,
			OwnerID:   kit.SupplierAccount,
			BizID:     targetAttr.BizID,
		}
		group, err := a.grp.CreateObjectGroup(kit, &g)
		if err != nil {
			setErr = append(setErr, rowInfo{Row: idx, Info: err.Error(), PropertyID: targetAttr.PropertyID})
			return itemErr, setErr, err
		}

		targetAttr.PropertyGroup = group.GroupID
	}

	return itemErr, setErr, nil
}

// getNonInnerAttributes get non inner object attributes
func (a *attribute) getNonInnerAttributes(kit *rest.Kit, objID string) (*metadata.QueryModelAttributeDataResult, error) {
	query := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			metadata.AttributeFieldObjectID: objID,
			metadata.AttributeFieldIsSystem: false,
			metadata.AttributeFieldIsAPI:    false,
			common.BKAppIDField:             0,
		},
	}

	rsp, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, query)
	if err != nil {
		blog.Errorf("get module attr failed, err: %v, rid: %s", err, kit.Rid)
		return nil, err
	}

	return rsp, nil
}

// FindObjectAttributeBatch batch get objects attributes
func (a *attribute) FindObjectAttributeBatch(kit *rest.Kit, objIDs []string) (mapstr.MapStr, error) {
	result := make(mapstr.MapStr)

	for _, objID := range objIDs {
		attrs, err := a.getNonInnerAttributes(kit, objID)
		if err != nil {
			return nil, err
		}

		for _, attr := range attrs.Info {
			cond := mapstr.MapStr{
				metadata.GroupFieldGroupID:  attr.PropertyGroup,
				metadata.GroupFieldObjectID: attr.ObjectID,
			}
			grps, err := a.grp.FindGroupByObject(kit, attr.ObjectID, cond, attr.BizID)
			if err != nil {
				return nil, err
			}

			for _, grp := range grps {
				// should be only one
				attr.PropertyGroupName = grp.GroupName
			}
		}

		result.Set(objID, mapstr.MapStr{
			"attr": attrs,
		})
	}

	return result, nil
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
	attrIDMap := make(map[string][]int64, 0)
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
		attrIDMap[attrItem.ObjectID] = append(attrIDMap[attrItem.ObjectID], attrItem.ID)
	}

	for objID, instIDs := range attrIDMap {
		// delete the attribute.
		deleteCond := &metadata.DeleteOption{
			Condition: mapstr.MapStr{
				common.BKFieldID: mapstr.MapStr{common.BKDBIN: instIDs},
			},
		}
		_, err = a.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, objID, deleteCond)
		if err != nil {
			blog.Errorf("delete object attribute failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}
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
		blog.Errorf("unmarshal mapstr data into module failed, module: %s, err: %v, rid: %s", attr, err, kit.Rid)
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
		blog.Errorf("failed to update module attr, err: %v, rid: %s", err, kit.Rid)
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
