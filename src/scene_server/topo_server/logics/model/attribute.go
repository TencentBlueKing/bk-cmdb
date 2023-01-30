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
	// CreateObjectBatch upsert object attributes
	CreateObjectBatch(kit *rest.Kit, data map[string]metadata.ImportObjectData) (mapstr.MapStr, error)
	// FindObjectBatch find object to attributes mapping
	FindObjectBatch(kit *rest.Kit, objIDs []string) (mapstr.MapStr, error)
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
	cond := []map[string]interface{}{{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: data.PropertyGroup,
		common.BKAppIDField:           data.BizID,
	}}

	defCntRes, err := a.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNamePropertyGroup, cond)
	if err != nil {
		blog.Errorf("get attribute group count by cond(%#v) failed, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(defCntRes) != 1 {
		blog.Errorf("get attr group count by cond(%#v) returns not one result, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommNotFound)
	}

	if defCntRes[0] > 0 {
		return nil
	}

	if data.BizID == 0 {
		data.PropertyGroup = common.BKDefaultField
		return nil
	}

	// create the biz default group if it is not exist
	bizDefaultGroupCond := []map[string]interface{}{{
		common.BKObjIDField:           data.ObjectID,
		common.BKPropertyGroupIDField: common.BKBizDefault,
		common.BKAppIDField:           data.BizID,
	}}

	bizDefCntRes, err := a.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
		common.BKTableNamePropertyGroup, bizDefaultGroupCond)
	if err != nil {
		blog.Errorf("get attribute group count by cond(%#v) failed, err: %v, rid: %s", cond, err, kit.Rid)
		return err
	}

	if len(bizDefCntRes) != 1 {
		blog.Errorf("get attr group count by cond(%#v) returns not one result, err: %v, rid: %s", cond, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommNotFound)
	}

	if bizDefCntRes[0] == 0 {
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
	}

	data.PropertyGroup = common.BKBizDefault
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
		blog.Errorf("create model(%s) attr but it is duplicated, input: %#v, rid: %s", data.ObjectID, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	if len(resp.Created) != 1 {
		blog.Errorf("created model(%s) attr amount is not one, input: %#v, rid: %s", data.ObjectID, input, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}

	// get current model attribute data by id.
	attrReq := &metadata.QueryCondition{Condition: mapstr.MapStr{metadata.AttributeFieldID: int64(resp.Created[0].ID)}}
	attrRes, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, data.ObjectID, attrReq)
	if err != nil {
		blog.Errorf("get created model attribute by id(%d) failed, err: %v, rid: %s", resp.Created[0].ID, err, kit.Rid)
		return nil, err
	}

	if len(attrRes.Info) != 1 {
		blog.Errorf("get created model attribute by id(%d) returns not one attr, rid: %s", resp.Created[0].ID, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommNotFound)
	}

	data = &attrRes.Info[0]

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, data.ID, data)
	if err != nil {
		blog.Errorf("gen audit log after creating attr %s failed, err: %v, rid: %s", data.PropertyName, err, kit.Rid)
		return nil, err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log after creating attr %s failed, err: %v, rid: %s", data.PropertyName, err, kit.Rid)
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
		rsp, err := a.clientSet.CoreService().Model().DeleteModelAttr(kit.Ctx, kit.Header, objID, deleteCond)
		if err != nil {
			blog.Errorf("delete object attribute failed, err: %v, rid: %s", err, kit.Rid)
			return err
		}

		if ccErr := rsp.CCError(); ccErr != nil {
			blog.Errorf("delete object attribute failed, err: %v, rid: %s", ccErr, kit.Rid)
			return ccErr
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

type rowInfo struct {
	Row  int64  `json:"row"`
	Info string `json:"info"`
	// value can empty, eg: parse error
	PropID string `json:"bk_property_id"`
}

type createObjectBatchObjResult struct {
	// 异常错误，比如接卸该行数据失败， 查询数据失败
	Errors []rowInfo `json:"errors,omitempty"`
	// 新加属性失败
	InsertFailed []rowInfo `json:"insert_failed,omitempty"`
	// 更新属性失败
	UpdateFailed []rowInfo `json:"update_failed,omitempty"`
	// 成功数据的信息
	Success []rowInfo `json:"success,omitempty"`
	// 失败信息，如模型不存在
	Error string `json:"error,omitempty"`
}

// CreateObjectBatch this method doesn't act as it's name, it create or update model's attributes indeed.
// it only operate on model already exist, that is to say no new model will be created.
func (a *attribute) CreateObjectBatch(kit *rest.Kit, inputDataMap map[string]metadata.ImportObjectData) (mapstr.MapStr,
	error) {

	result := mapstr.New()
	hasError := false
	for objID, inputData := range inputDataMap {
		// check if the object exists
		isObjExists, err := a.obj.IsObjectExist(kit, objID)
		if err != nil {
			result[objID] = createObjectBatchObjResult{
				Error: fmt.Sprintf("check if object(%s) exists failed, err: %v", objID, err),
			}
			hasError = true
			continue
		}
		if !isObjExists {
			result[objID] = createObjectBatchObjResult{Error: fmt.Sprintf("object (%s) does not exist", objID)}
			hasError = true
			continue
		}

		// get group name to property id map
		groupNames := make([]string, 0)
		for _, attr := range inputData.Attr {
			if len(attr.PropertyGroupName) == 0 {
				continue
			}
			groupNames = append(groupNames, attr.PropertyGroupName)
		}

		grpNameIDMap := make(map[string]string)
		if len(groupNames) > 0 {
			grpCond := metadata.QueryCondition{
				Condition: mapstr.MapStr{
					metadata.GroupFieldGroupName: mapstr.MapStr{common.BKDBIN: groupNames},
					metadata.GroupFieldObjectID:  objID,
				},
				Fields: []string{metadata.GroupFieldGroupID, metadata.GroupFieldGroupName},
				Page:   metadata.BasePage{Limit: common.BKNoLimit},
			}

			grpRsp, err := a.clientSet.CoreService().Model().ReadAttributeGroup(kit.Ctx, kit.Header, objID, grpCond)
			if err != nil {
				result[objID] = createObjectBatchObjResult{
					Error: fmt.Sprintf("find object group failed, err: %v, cond: %#v", err, grpCond),
				}
				hasError = true
				continue
			}

			for _, grp := range grpRsp.Info {
				grpNameIDMap[grp.GroupName] = grp.GroupID
			}
		}

		// upsert the object's attribute
		result[objID], hasError = a.upsertObjectAttrBatch(kit, objID, inputData.Attr, grpNameIDMap)
	}

	if hasError {
		return result, kit.CCError.Error(common.CCErrCommNotAllSuccess)
	}
	return result, nil
}

func (a *attribute) upsertObjectAttrBatch(kit *rest.Kit, objID string, attributes map[int64]metadata.Attribute,
	grpNameIDMap map[string]string) (createObjectBatchObjResult, bool) {

	objRes := createObjectBatchObjResult{}
	hasError := false
	for idx, attr := range attributes {
		propID := attr.PropertyID
		if propID == common.BKInstParentStr {
			continue
		}

		if err := a.IsValid(kit, true, &attr); err != nil {
			blog.Errorf("attribute(%#v) is invalid, rid: %s", attr, err, kit.Rid)
			objRes.Errors = append(objRes.Errors, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
			hasError = true
			continue
		}

		attr.OwnerID = kit.SupplierAccount
		attr.ObjectID = objID

		if len(attr.PropertyGroupName) != 0 {
			groupID, exists := grpNameIDMap[attr.PropertyGroupName]
			if exists {
				attr.PropertyGroup = groupID
			} else {
				grp := metadata.CreateModelAttributeGroup{
					Data: metadata.Group{GroupName: attr.PropertyGroupName, GroupID: NewGroupID(false), ObjectID: objID,
						OwnerID: kit.SupplierAccount, BizID: attr.BizID,
					}}

				_, err := a.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header, objID, grp)
				if err != nil {
					blog.Errorf("create attribute group[%#v] failed, err: %v, rid: %s", grp, err, kit.Rid)
					objRes.Errors = append(objRes.Errors, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
					hasError = true
					continue
				}
				attr.PropertyGroup = grp.Data.GroupID
			}
		} else {
			attr.PropertyGroup = NewGroupID(true)
		}

		// check if attribute exists, if exists, update these attributes, otherwise, create the attribute
		attrCond := mapstr.MapStr{metadata.AttributeFieldObjectID: objID, metadata.AttributeFieldPropertyID: propID}
		util.AddModelBizIDCondition(attrCond, attr.BizID)

		attrCnt, err := a.clientSet.CoreService().Count().GetCountByFilter(kit.Ctx, kit.Header,
			common.BKTableNameObjAttDes, []map[string]interface{}{attrCond})
		if err != nil {
			blog.Errorf("count attribute failed, err: %v, cond: %#v, rid: %s", err, attrCond, kit.Rid)
			objRes.Errors = append(objRes.Errors, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
			hasError = true
			continue
		}

		if attrCnt[0] == 0 {
			// create attribute
			createAttrOpt := &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{attr}}
			_, err := a.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, objID, createAttrOpt)
			if err != nil {
				blog.Errorf("create attribute(%#v) failed, ObjID: %s, err: %v, rid: %s", attr, objID, err, kit.Rid)
				objRes.InsertFailed = append(objRes.InsertFailed, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
				hasError = true
				continue
			}
		} else {
			// update attribute
			updateData := attr.ToMapStr()
			updateData.Remove(metadata.AttributeFieldPropertyID)
			updateData.Remove(metadata.AttributeFieldObjectID)
			updateData.Remove(metadata.AttributeFieldID)
			updateAttrOpt := metadata.UpdateOption{Condition: attrCond, Data: updateData}
			_, err := a.clientSet.CoreService().Model().UpdateModelAttrs(kit.Ctx, kit.Header, objID, &updateAttrOpt)
			if err != nil {
				blog.Errorf("failed to update module attr, err: %s, rid: %s", err, kit.Rid)
				objRes.UpdateFailed = append(objRes.UpdateFailed, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
				hasError = true
				continue
			}
		}

		objRes.Success = append(objRes.Success, rowInfo{Row: idx, PropID: attr.PropertyID})
	}

	return objRes, hasError
}

// FindObjectBatch find object to attribute mapping batch
func (a *attribute) FindObjectBatch(kit *rest.Kit, objIDs []string) (mapstr.MapStr, error) {
	result := mapstr.New()

	for _, objID := range objIDs {
		attrCond := &metadata.QueryCondition{
			Condition: mapstr.MapStr{
				metadata.AttributeFieldObjectID: objID,
				metadata.AttributeFieldIsSystem: false,
				metadata.AttributeFieldIsAPI:    false,
				common.BKAppIDField:             0,
			},
			Page: metadata.BasePage{Limit: common.BKNoLimit},
		}
		attrRsp, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, attrCond)
		if err != nil {
			blog.Errorf("get object(%s) not inner attributes failed, err: %v, rid: %s", err, kit.Rid)
			return nil, err
		}

		if len(attrRsp.Info) == 0 {
			result.Set(objID, mapstr.MapStr{"attr": attrRsp.Info})
			continue
		}

		groupIDs := make([]string, 0)
		for _, attr := range attrRsp.Info {
			groupIDs = append(groupIDs, attr.PropertyGroup)
		}

		grpCond := metadata.QueryCondition{
			Condition: mapstr.MapStr{
				metadata.GroupFieldGroupID:  mapstr.MapStr{common.BKDBIN: groupIDs},
				metadata.GroupFieldObjectID: objID,
			},
			Fields: []string{metadata.GroupFieldGroupID, metadata.GroupFieldGroupName},
			Page:   metadata.BasePage{Limit: common.BKNoLimit},
		}

		grpRsp, err := a.clientSet.CoreService().Model().ReadAttributeGroup(kit.Ctx, kit.Header, objID, grpCond)
		if err != nil {
			blog.Errorf("find object group failed, err: %v, cond: %#v, rid: %s", err, grpCond, kit.Rid)
			return nil, err
		}

		grpIDNameMap := make(map[string]string)
		for _, grp := range grpRsp.Info {
			grpIDNameMap[grp.GroupID] = grp.GroupName
		}

		for idx, attr := range attrRsp.Info {
			attrRsp.Info[idx].PropertyGroupName = grpIDNameMap[attr.PropertyGroup]
		}

		result.Set(objID, mapstr.MapStr{"attr": attrRsp.Info})
	}

	return result, nil
}
