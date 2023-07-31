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
	"sort"
	"unicode/utf8"

	"configcenter/pkg/filter"
	filtertools "configcenter/pkg/tools/filter"
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
	"configcenter/src/common/valid"
)

// AttributeOperationInterface attribute operation methods
type AttributeOperationInterface interface {
	CreateObjectAttribute(kit *rest.Kit, data *metadata.Attribute) (*metadata.Attribute, error)
	BatchCreateObjectAttr(kit *rest.Kit, objID string, attrs []*metadata.Attribute, fromTemplate bool) error
	CreateTableObjectAttribute(kit *rest.Kit, data *metadata.Attribute) (*metadata.Attribute, error)
	DeleteObjectAttribute(kit *rest.Kit, attrItems []metadata.Attribute) error
	UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64, isSync bool) error
	// CreateObjectBatch upsert object attributes
	CreateObjectBatch(kit *rest.Kit, data map[string]metadata.ImportObjectData) (mapstr.MapStr, error)
	UpdateTableObjectAttr(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64) error
	// FindObjectBatch find object to attributes mapping
	FindObjectBatch(kit *rest.Kit, objIDs []string) (mapstr.MapStr, error)
	ValidObjIDAndInstID(kit *rest.Kit, objID string, option interface{}, isMultiple bool) error
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

// getEnumQuoteOption get enum quote field option bk_obj_id and bk_inst_id value
func (a *attribute) getEnumQuoteOption(kit *rest.Kit, option interface{}, isMultiple bool) (string, []int64, error) {
	if option == nil {
		return "", nil, kit.CCError.Errorf(common.CCErrCommParamsLostField, "option")
	}

	arrOption, ok := option.([]interface{})
	if !ok {
		blog.Errorf("option %v not enum quote option, rid: %s", option, kit.Rid)
		return "", nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}
	if len(arrOption) == 0 {
		return "", nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	if !isMultiple && len(arrOption) != 1 {
		blog.Errorf("enum option is single choice, but arr option value is multiple, rid: %s", kit.Rid)
		return "", nil, kit.CCError.CCError(common.CCErrCommParamsNeedSingleChoice)
	}

	if len(arrOption) > common.AttributeOptionArrayMaxLength {
		blog.Errorf("option array length %d exceeds max length %d, rid: %s", len(arrOption),
			common.AttributeOptionArrayMaxLength, kit.Rid)
		return "", nil, kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option",
			common.AttributeOptionArrayMaxLength)
	}

	var quoteObjID string
	instIDMap := make(map[int64]interface{}, 0)
	for _, o := range arrOption {
		mapOption, ok := o.(map[string]interface{})
		if !ok || mapOption == nil {
			blog.Errorf("enum quote option %v must contain bk_obj_id and bk_inst_id, rid: %s", option, kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
		}
		objIDVal, objIDOk := mapOption[common.BKObjIDField]
		if !objIDOk || objIDVal == "" {
			blog.Errorf("enum quote option bk_obj_id can't be empty, rid: %s", option, kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "option bk_obj_id")
		}
		objID, ok := objIDVal.(string)
		if !ok {
			blog.Errorf("objIDVal %v not string, rid: %s", objIDVal, kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommParamsNeedString, "option bk_obj_id")
		}
		if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(objID) {
			blog.Errorf("option bk_obj_id %s length %d exceeds max length %d, rid: %s", objID,
				utf8.RuneCountInString(objID), common.AttributeOptionValueMaxLength, kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option bk_obj_id",
				common.AttributeOptionValueMaxLength)
		}

		if quoteObjID == "" {
			quoteObjID = objID
		} else if quoteObjID != objID {
			blog.Errorf("enum quote objID not unique, objID: %s, rid: %s", objID, kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "quote objID")
		}

		instIDVal, instIDOk := mapOption[common.BKInstIDField]
		if !instIDOk || instIDVal == "" {
			blog.Errorf("enum quote option bk_inst_id can't be empty, rid: %s", option, kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "option bk_inst_id")
		}

		switch mapOption["type"] {
		case "int":
			instID, err := util.GetInt64ByInterface(instIDVal)
			if err != nil {
				return "", nil, err
			}
			if instID == 0 {
				return "", nil, fmt.Errorf("enum quote instID is %d, it is illegal", instID)
			}
			instIDMap[instID] = struct{}{}
		default:
			blog.Errorf("enum quote option type must be 'int', current: %v, rid: %s", mapOption["type"], kit.Rid)
			return "", nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option type")
		}
	}

	instIDs := make([]int64, 0)
	for instID := range instIDMap {
		instIDs = append(instIDs, instID)
	}

	return quoteObjID, instIDs, nil
}

// enumQuoteCanNotUseModel 校验引用模型不能为集群、模块、进程、容器、自定义层级相关模型
func (a *attribute) enumQuoteCanNotUseModel(kit *rest.Kit, objID string) error {

	// 校验引用模型不能为集群，模块，进程内置模型 TODO 容器相关的模型暂无定义，后续添加
	switch objID {
	case common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc:
		return fmt.Errorf("enum quote obj can not inner model")
	}

	// 校验引用模型不能为自定义层级模型
	query := &metadata.QueryCondition{
		Fields: []string{common.BKObjIDField, common.BKAsstObjIDField},
		Condition: mapstr.MapStr{
			common.AssociationKindIDField: common.AssociationKindMainline,
		},
		DisableCounter: true,
	}
	mainlineAsstRsp, err := a.clientSet.CoreService().Association().ReadModelAssociation(kit.Ctx, kit.Header, query)
	if err != nil {
		blog.Errorf("search mainline association failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	mainlineObjectChildMap := make(map[string]string, 0)
	for _, asst := range mainlineAsstRsp.Info {
		if asst.ObjectID == common.BKInnerObjIDHost {
			continue
		}
		mainlineObjectChildMap[asst.AsstObjID] = asst.ObjectID
	}

	objectIDs := make([]string, 0)
	for objectID := common.BKInnerObjIDApp; len(objectID) != 0; objectID = mainlineObjectChildMap[objectID] {
		if objectID == common.BKInnerObjIDApp || objectID == common.BKInnerObjIDSet ||
			objectID == common.BKInnerObjIDModule {
			continue
		}
		objectIDs = append(objectIDs, objectID)
	}

	for _, customObjID := range objectIDs {
		if objID == customObjID {
			return fmt.Errorf("enum quote obj can not custom model")
		}
	}
	return nil
}

// ValidObjIDAndInstID check obj is inner model and obj is exist, inst is exit
func (a *attribute) ValidObjIDAndInstID(kit *rest.Kit, objID string, option interface{}, isMultiple bool) error {
	quoteObjID, instIDs, err := a.getEnumQuoteOption(kit, option, isMultiple)
	if err != nil {
		blog.Errorf("get enum quote option value failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	if quoteObjID == "" || len(instIDs) == 0 {
		return fmt.Errorf("enum quote objID or instID is empty, objIDs: %s, instIDs: %v", quoteObjID, instIDs)
	}

	isObjExists, err := a.obj.IsObjectExist(kit, quoteObjID)
	if err != nil {
		blog.Errorf("check obj id is exist failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	if !isObjExists {
		blog.Errorf("enum quote option bk_obj_id is not exist, objID: %s, rid: %s", quoteObjID, kit.Rid)
		return fmt.Errorf("enum quote objID is not exist, objID: %s", quoteObjID)
	}

	if quoteObjID == objID {
		blog.Errorf("enum quote model can not model self, objID: %s, rid: %s", objID, kit.Rid)
		return fmt.Errorf("enum quote model can not model self, objID: %s", objID)
	}

	// 集群，模块，进程，容器，自定义层级模块不能被引用
	if err := a.enumQuoteCanNotUseModel(kit, quoteObjID); err != nil {
		blog.Errorf("enum quote model can not use some inner model, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	input := &metadata.QueryCondition{
		Fields: []string{common.GetInstIDField(quoteObjID)},
		Condition: mapstr.MapStr{
			common.GetInstIDField(quoteObjID): mapstr.MapStr{common.BKDBIN: instIDs},
		},
	}
	resp, err := a.clientSet.CoreService().Instance().ReadInstance(kit.Ctx, kit.Header, quoteObjID, input)
	if err != nil {
		blog.Errorf("get inst data failed, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}
	if resp.Count == 0 {
		blog.Errorf("enum quote option inst not exist, input: %+v, rid: %s", input, kit.Rid)
		return fmt.Errorf("enum quote inst not exist, instIDs: %v", instIDs)
	}

	return nil
}

func (a *attribute) validTableAttributes(kit *rest.Kit, option interface{}) error {

	if option == nil {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "option")
	}

	tableOption, err := metadata.ParseTableAttrOption(option)
	if err != nil {
		blog.Errorf("get attribute option failed, option: %+v, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	headerAttrMap, err := a.validAndGetTableAttrHeaderDetail(kit, tableOption.Header)
	if err != nil {
		return err
	}

	if err := a.ValidTableAttrDefaultValue(kit, tableOption.Default, headerAttrMap); err != nil {
		return err
	}
	return nil
}

// validAndGetTableAttrHeaderDetail in the creation and update scenarios, the full amount of header
// content needs to be passed.
func (a *attribute) validAndGetTableAttrHeaderDetail(kit *rest.Kit, header []metadata.Attribute) (
	map[string]*metadata.Attribute, error) {

	if len(header) == 0 {
		return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "table header")
	}

	if len(header) > metadata.TableHeaderMaxNum {
		return nil, kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "table header", metadata.TableHeaderMaxNum)
	}
	propertyAttr := make(map[string]*metadata.Attribute)
	var longCharNum int
	for index := range header {
		// determine whether the underlying type is legal
		if !metadata.ValidTableFieldBaseType(header[index].PropertyType) {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, header[index].PropertyType)
		}
		// the number of long characters in the basic type of the table
		// field type cannot exceed the maximum value supported by the system.
		if header[index].PropertyType == common.FieldTypeLongChar {
			longCharNum++
		}
		if longCharNum > metadata.TableLongCharMaxNum {
			return nil, kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "table header", metadata.TableLongCharMaxNum)
		}

		// check if property type for creation is valid, can't update property type
		if header[index].PropertyType == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
		}

		if header[index].PropertyID == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyID)
		}

		if common.AttributeIDMaxLength < utf8.RuneCountInString(header[index].PropertyID) {
			return nil, kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Language(
					"model_attr_bk_property_id"), common.AttributeIDMaxLength)
		}

		match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, header[index].PropertyID)
		if err != nil {
			return nil, err
		}

		if !match {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, header[index].PropertyID)
		}

		if header[index].PropertyName == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyName)
		}

		if common.AttributeNameMaxLength < utf8.RuneCountInString(header[index].PropertyName) {
			return nil, kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed,
				a.lang.CreateDefaultCCLanguageIf(util.GetLanguage(kit.Header)).Language(
					"model_attr_bk_property_name"), common.AttributeNameMaxLength)
		}

		propertyAttr[header[index].PropertyID] = &header[index]
	}
	return propertyAttr, nil
}

// ValidTableAttrDefaultValue attr: key is property_id, value is the corresponding header content.
func (a *attribute) ValidTableAttrDefaultValue(kit *rest.Kit, defaultValue []map[string]interface{},
	attr map[string]*metadata.Attribute) error {

	if len(defaultValue) == 0 {
		return nil
	}
	if len(defaultValue) > metadata.TableDefaultMaxLines {
		return kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "table.default.limit", metadata.TableDefaultMaxLines)
	}
	// judge the legality of each field of the default
	// value according to the attributes of the header.
	for _, value := range defaultValue {
		for k, v := range value {
			attr[k].IsRequired = false
			if err := attr[k].ValidTableDefaultAttr(kit.Ctx, v); err.ErrCode != 0 {
				return err.ToCCError(kit.CCError)
			}
		}
	}
	return nil
}

// isValid check is valid
func (a *attribute) isValid(kit *rest.Kit, isUpdate bool, data *metadata.Attribute) error {
	if data.PropertyID == common.BKInstParentStr {
		return nil
	}

	if (isUpdate && data.IsMultiple != nil) || !isUpdate {
		if err := valid.ValidPropertyTypeIsMultiple(data.PropertyType, *data.IsMultiple, kit.CCError); err != nil {
			return err
		}
	}

	// 用户类型字段，在创建的时候默认是支持可多选的，而且这个字段是否可多选在页面是不可配置的,所以在创建的时候将值置为true
	if data.PropertyType == common.FieldTypeUser && !isUpdate {
		isMultiple := true
		data.IsMultiple = &isMultiple
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
		if err := valid.ValidPropertyOption(data.PropertyType, data.Option, *data.IsMultiple, data.Default, kit.Rid,
			kit.CCError); err != nil {
			return err
		}
	}

	// check enum quote field option validity creation or update
	if data.PropertyType == common.FieldTypeEnumQuote && data.IsMultiple != nil {
		if err := a.ValidObjIDAndInstID(kit, data.ObjectID, data.Option, *data.IsMultiple); err != nil {
			blog.Errorf("check objID and instID failed, err: %v, rid: %s", err, kit.Rid)
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
	case common.FieldTypeInt, common.FieldTypeEnum, common.FieldTypeList, common.FieldTypeEnumMulti:
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

func (a *attribute) validCreateTableAttribute(kit *rest.Kit, data *metadata.Attribute) error {
	// check if the object is mainline object, if yes. then user can not create required attribute.
	yes, err := a.isMainlineModel(kit, data.ObjectID)
	if err != nil {
		blog.Errorf("not allow to add required attribute to mainline object: %+v. "+"rid: %d.", data, kit.Rid)
		return err
	}

	if yes && data.IsRequired {
		return kit.CCError.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
	}

	// check the object id
	exist, err := a.obj.IsObjectExist(kit, data.ObjectID)
	if err != nil {
		return err
	}
	if !exist {
		blog.Errorf("obj id is not exist, obj id: %s, rid: %s", data.ObjectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	if err = a.checkAttributeGroupExist(kit, data); err != nil {
		blog.Errorf("failed to create the default group, err: %s, rid: %s", err, kit.Rid)
		return err
	}

	if err := a.validTableAttributes(kit, data.Option); err != nil {
		return err
	}
	return nil
}

func (a *attribute) createTableModelAndAttributeGroup(kit *rest.Kit, data *metadata.Attribute) error {

	t := metadata.Now()
	obj := metadata.Object{
		ObjCls:     metadata.ClassificationTableID,
		ObjIcon:    "icon-cc-table",
		ObjectID:   data.ObjectID,
		ObjectName: data.PropertyName,
		IsHidden:   true,
		Creator:    string(metadata.FromCCSystem),
		Modifier:   string(metadata.FromCCSystem),
		CreateTime: &t,
		LastTime:   &t,
		OwnerID:    kit.SupplierAccount,
	}

	objRsp, err := a.clientSet.CoreService().Model().CreateTableModel(kit.Ctx, kit.Header,
		&metadata.CreateModel{Spec: obj, Attributes: []metadata.Attribute{*data}})
	if err != nil {
		blog.Errorf("create table object(%s) failed, err: %v, rid: %s", obj.ObjectID, err, kit.Rid)
		return err
	}

	obj.ID = int64(objRsp.Created.ID)
	objID := metadata.GenerateModelQuoteObjName(data.ObjectID, data.PropertyID)
	// create the default group
	groupData := metadata.Group{
		IsDefault:  true,
		GroupIndex: -1,
		GroupName:  "Default",
		GroupID:    NewGroupID(true),
		ObjectID:   objID,
		OwnerID:    obj.OwnerID,
	}

	_, err = a.clientSet.CoreService().Model().CreateAttributeGroup(kit.Ctx, kit.Header,
		objID, metadata.CreateModelAttributeGroup{Data: groupData})
	if err != nil {
		blog.Errorf("create attribute group[%s] failed, err: %v, rid: %s", groupData.GroupID, err, kit.Rid)
		return err
	}

	// generate audit log of object attribute group.
	audit := auditlog.NewObjectAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)
	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, obj.ID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

// CreateTableObjectAttribute create internal form fields in a separate process
func (a *attribute) CreateTableObjectAttribute(kit *rest.Kit, data *metadata.Attribute) (*metadata.Attribute, error) {
	if data.IsOnly {
		data.IsRequired = true
	}

	if len(data.PropertyGroup) == 0 {
		data.PropertyGroup = "default"
	}

	if data.TemplateID != 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKTemplateID)
	}

	if err := a.validCreateTableAttribute(kit, data); err != nil {
		return nil, err
	}

	if err := a.createTableModelAndAttributeGroup(kit, data); err != nil {
		return nil, err
	}

	data.OwnerID = kit.SupplierAccount
	if err := a.createModelQuoteRelation(kit, data.ObjectID, data.PropertyID); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *attribute) createModelQuoteRelation(kit *rest.Kit, objectID, propertyID string) error {

	relation := metadata.ModelQuoteRelation{
		DestModel:  metadata.GenerateModelQuoteObjID(objectID, propertyID),
		SrcModel:   objectID,
		PropertyID: propertyID,
		Type:       common.ModelQuoteType(common.FieldTypeInnerTable),
	}

	if cErr := a.clientSet.CoreService().ModelQuote().CreateModelQuoteRelation(kit.Ctx, kit.Header,
		[]metadata.ModelQuoteRelation{relation}); cErr != nil {
		blog.Errorf("created quote relation failed, relation: %#v, err: %v, rid: %s", relation, cErr, kit.Rid)
		return cErr
	}

	return nil
}

func (a *attribute) preCheckObjectAttr(kit *rest.Kit, objID string, data *metadata.Attribute) error {

	if data.ObjectID != objID {
		blog.Errorf("attr object id is invalid, object: %+v, obj id: %s, rid: %s", data, objID, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsInvalid, objID)
	}

	if data.IsOnly {
		data.IsRequired = true
	}

	if len(data.PropertyGroup) == 0 {
		data.PropertyGroup = "default"
	}

	// check if the object is mainline object, if yes. then user can not create required attribute.
	yes, err := a.isMainlineModel(kit, data.ObjectID)
	if err != nil {
		blog.Errorf("not allow to add required attribute to mainline object: %+v, rid: %s", data, kit.Rid)
		return err
	}

	if yes && data.IsRequired {
		return kit.CCError.Error(common.CCErrTopoCanNotAddRequiredAttributeForMainlineModel)
	}

	// check the object id
	exist, err := a.obj.IsObjectExist(kit, data.ObjectID)
	if err != nil {
		return err
	}

	if !exist {
		blog.Errorf("obj id is not exist, obj id: %s, rid: %s", data.ObjectID, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKObjIDField)
	}

	if err = a.checkAttributeGroupExist(kit, data); err != nil {
		blog.Errorf("failed to create the default group, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := a.isValid(kit, false, data); err != nil {
		return err
	}
	return nil
}

// BatchCreateObjectAttr batch create object attributes
func (a *attribute) BatchCreateObjectAttr(kit *rest.Kit, objID string, attrs []*metadata.Attribute,
	fromTemplate bool) error {

	objAttrs := make([]metadata.Attribute, 0)
	for _, data := range attrs {
		if err := a.preCheckObjectAttr(kit, objID, data); err != nil {
			return err
		}
		objAttrs = append(objAttrs, *data)
	}

	input := metadata.CreateModelAttributes{
		Attributes:   objAttrs,
		FromTemplate: fromTemplate,
	}
	resp, err := a.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, objID, &input)
	if err != nil {
		blog.Errorf("failed to create model attrs, input: %#v, err: %v, rid: %s", input, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommHTTPDoRequestFailed)
	}

	for _, exception := range resp.Exceptions {
		return kit.CCError.New(int(exception.Code), exception.Message)
	}

	if len(resp.Repeated) > 0 {
		blog.Errorf("create model(%s) attr but it is duplicated, input: %#v, rid: %s", objID, input, kit.Rid)
		return kit.CCError.CCError(common.CCErrorAttributeNameDuplicated)
	}

	objAttrsLen := len(objAttrs)

	if len(resp.Created) != objAttrsLen {
		blog.Errorf("created model(%s) attr amount is not one, input: %#v, rid: %s", objID, input, kit.Rid)
		return kit.CCError.CCError(common.CCErrTopoObjectAttributeCreateFailed)
	}

	ids := make([]uint64, 0)
	for _, data := range resp.Created {
		ids = append(ids, data.ID)
	}

	// get current model attribute data by id.
	attrReq := &metadata.QueryCondition{
		Condition: mapstr.MapStr{metadata.AttributeFieldID: mapstr.MapStr{common.BKDBIN: ids}}}

	attrRes, err := a.clientSet.CoreService().Model().ReadModelAttr(kit.Ctx, kit.Header, objID, attrReq)
	if err != nil {
		blog.Errorf("get created model attribute by ids(%v) failed, err: %v, rid: %s", ids, err, kit.Rid)
		return err
	}

	if len(attrRes.Info) != objAttrsLen {
		blog.Errorf("get the number of model attrs based on ids(%v) does not meet expectations, rid: %s", ids, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsIsInvalid)
	}

	// generate audit log of model attribute.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	auditLog, err := audit.BatchGenerateAuditLog(generateAuditParameter, objID, attrRes.Info)
	if err != nil {
		blog.Errorf("gen audit log after creating objID %s failed, err: %v, rid: %s", objID, err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, auditLog...); err != nil {
		blog.Errorf("save audit log failed, attr ids: %v, err: %v, rid: %s", ids, err, kit.Rid)
		return err
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
		blog.Errorf("not allow to add required attribute to mainline object: %+v, rid: %s", data, kit.Rid)
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

	// the templateID is not allowed to be non-zero whether it is creating or updating the scene.
	if data.TemplateID != 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, common.BKTemplateID)
	}

	if err := a.isValid(kit, false, data); err != nil {
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
	if err := a.saveLogForCreateObjectAttribute(kit, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *attribute) saveLogForCreateObjectAttribute(kit *rest.Kit, data *metadata.Attribute) error {
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditCreate)

	auditLog, err := audit.GenerateAuditLog(generateAuditParameter, data.ID, data)
	if err != nil {
		blog.Errorf("gen audit log after creating attr %s failed, err: %v, rid: %s", data.PropertyName, err, kit.Rid)
		return err
	}

	// save audit log.
	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log after creating attr %s failed, err: %v, rid: %s", data.PropertyName, err, kit.Rid)
		return err
	}
	return nil
}

// DeleteObjectAttribute delete object attribute
func (a *attribute) DeleteObjectAttribute(kit *rest.Kit, attrItems []metadata.Attribute) error {

	auditLogArr := make([]metadata.AuditLog, 0)
	attrIDMap := make(map[string][]int64, 0)
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditDelete)
	for _, attrItem := range attrItems {
		if attrItem.PropertyType == common.FieldTypeInnerTable {
			blog.Errorf("property is error, attrItem: %+v, rid: %s", attrItem, kit.Rid)
			return kit.CCError.New(common.CCErrTopoObjectSelectFailed, common.BKPropertyTypeField)
		}
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

func (a *attribute) getTableAttrOptionFromDB(kit *rest.Kit, attID, bizID int64) (
	*metadata.TableAttributesOption, string, error) {

	cond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKFieldID: attID},
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		DisableCounter: true,
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttrsWithTableByCondition(kit.Ctx, kit.Header, bizID, cond)
	if nil != err {
		blog.Errorf("search table attr failed, cond: %+v, bizID: %d, err: %v, rid: %s", cond, bizID, err, kit.Rid)
		return nil, "", err
	}

	if len(resp.Info) == 0 {
		blog.Errorf("no table attr found, cond: %+v, bizID: %d, err: %v, rid: %s", cond, bizID, err, kit.Rid)
		return nil, "", kit.CCError.CCError(common.CCErrCommNotFound)
	}

	if len(resp.Info) > 1 {
		blog.Errorf("multi table attr found, cond: %+v, bizID: %d, err: %v, rid: %s", cond, bizID, err, kit.Rid)
		return nil, "", kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	if resp.Info[0].PropertyType != common.FieldTypeInnerTable {
		blog.Errorf("attr type error, property: %v, cond: %+v, bizID: %d, err: %v, rid: %s", resp.Info[0].PropertyType,
			cond, bizID, err, kit.Rid)
		return nil, "", kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	dbAttrsOp, err := metadata.ParseTableAttrOption(resp.Info[0].Option)
	if err != nil {
		blog.Errorf("get attribute option failed, error: %v, option: %v, rid: %s", err, kit.Rid)
		return nil, "", err
	}
	return dbAttrsOp, resp.Info[0].ObjectID, nil
}

func calcTableOptionDiffDefault(kit *rest.Kit, curAttrsOp, dbAttrsOp *metadata.TableAttributesOption, objID string) (
	*metadata.TableAttributesOption, *metadata.TableAttributesOption, []string, error) {
	// according to this map, it is judged whether it is an operation
	// to delete the table header in the update scene.
	curHeaderPropertyIDMap := make(map[string]metadata.Attribute)
	createAttrMap := make(map[string]metadata.Attribute)

	updated := new(metadata.TableAttributesOption)
	for _, header := range curAttrsOp.Header {
		// determine whether the underlying type is legal
		if !metadata.ValidTableFieldBaseType(header.PropertyType) {
			return nil, nil, nil, fmt.Errorf("table header type is invalid, type : %v", header.PropertyType)
		}
		curHeaderPropertyIDMap[header.PropertyID] = header
		header.ObjectID = metadata.GenerateModelQuoteObjID(objID, header.PropertyID)
		createAttrMap[header.PropertyID] = header
	}

	deletePropertyIDs := make([]string, 0)

	for idx := range dbAttrsOp.Header {
		value, ok := curHeaderPropertyIDMap[dbAttrsOp.Header[idx].PropertyID]
		if !ok {
			deletePropertyIDs = append(deletePropertyIDs, dbAttrsOp.Header[idx].PropertyID)
			continue
		}
		// In the update scenario, obtain the corresponding value from the DB for the unchangeable data
		value.PropertyIndex = dbAttrsOp.Header[idx].PropertyIndex
		value.BizID = dbAttrsOp.Header[idx].BizID
		value.ObjectID = metadata.GenerateModelQuoteObjID(objID, dbAttrsOp.Header[idx].PropertyID)
		value.PropertyGroup = dbAttrsOp.Header[idx].PropertyGroup
		updated.Header = append(updated.Header, value)
		// delete the update part, so that the remaining data in
		// curAttrsOp needs to be newly created.
		delete(createAttrMap, dbAttrsOp.Header[idx].PropertyID)
	}

	for id := range curAttrsOp.Default {
		for v := range curAttrsOp.Default[id] {
			if _, ok := curHeaderPropertyIDMap[v]; !ok {
				blog.Errorf("the default value tag: (%v) is not in the header, rid: %s", v, kit.Rid)
				return nil, nil, nil, kit.CCError.CCErrorf(common.CCErrCommParamsInvalid, "default")
			}
		}
		// since the default value is directly verified through attributes,
		// the default value is placed in the update part as a whole.
		updated.Default = append(updated.Default, curAttrsOp.Default[id])
	}

	header := make([]metadata.Attribute, 0)
	for _, v := range createAttrMap {
		header = append(header, v)
	}

	// the header here is the new header part
	curAttrsOp.Header = header
	if len(curAttrsOp.Header)+len(updated.Header) > metadata.TableHeaderMaxNum {
		return nil, nil, nil, kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "table header",
			metadata.TableHeaderMaxNum)
	}

	return curAttrsOp, updated, deletePropertyIDs, nil
}

// validUpdateHeader 判断更新场景下的header是否合法
func (a *attribute) validUpdateHeader(kit *rest.Kit, createHeader, updateHeader []metadata.Attribute) error {
	if len(createHeader) == 0 && len(updateHeader) == 0 {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, "table header")
	}

	allHeader := make([]metadata.Attribute, 0)

	if len(createHeader) > 0 {
		allHeader = append(allHeader, createHeader...)
	}
	if len(updateHeader) > 0 {
		allHeader = append(allHeader, updateHeader...)
	}

	_, err := a.validAndGetTableAttrHeaderDetail(kit, allHeader)
	if err != nil {
		return err
	}
	return nil
}

// UpdateTableObjectAttr update object table attribute
func (a *attribute) UpdateTableObjectAttr(kit *rest.Kit, data mapstr.MapStr, attrID, modelBizID int64) error {

	attr := new(metadata.Attribute)
	if err := mapstruct.Decode2Struct(data, attr); err != nil {
		blog.Errorf("unmarshal mapstr data into attr failed, attr: %s, err: %v, rid: %s", attr, err, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommJSONUnmarshalFailed, "data")
	}

	propertyID := util.GetStrByInterface(data[common.BKPropertyIDField])
	if propertyID == "" {
		return kit.CCError.CCErrorf(common.CCErrCommParamsNeedSet, common.BKPropertyIDField)
	}

	if err := a.canAttrsUpdate(kit, data, attrID, false, modelBizID); err != nil {
		return err
	}

	updateDataStruct, createDataStruct := *attr, *attr

	curAttrsOp, err := metadata.ParseTableAttrOption(attr.Option)
	if err != nil {
		blog.Errorf("decode attr option failed, bizID: %d, err: %v, rid: %s", modelBizID, err, kit.Rid)
		return err
	}

	dbAttrsOp, objID, err := a.getTableAttrOptionFromDB(kit, attrID, modelBizID)
	if err != nil {
		return err
	}

	// to update. here, the two parts need to be processed separately, and the underlying verification is different.
	created, updated, _, err := calcTableOptionDiffDefault(kit, curAttrsOp, dbAttrsOp, objID)
	if err != nil {
		blog.Errorf("calc table option header failed, objID: %s, curAttrsOp: %+v, dbAttrsOp: %+v, err: %v, rid: %s",
			objID, curAttrsOp, dbAttrsOp, err, kit.Rid)
		return err
	}

	if err := a.validUpdateHeader(kit, created.Header, updated.Header); err != nil {
		blog.Errorf("update header illegal, create header: %+v, update header: %+v, err: %v, rid: %s", created.Header,
			updated.Header, err, kit.Rid)
		return err
	}

	// it should be verified separately, because some headers are newly created and others are updated. different
	// scenarios correspond to different content that needs to be verified. for checking the default value, you can
	// check it as a whole, because you only need to check whether the default value conforms to the attr of the header.
	// the checksum operation of the default value is uniformly placed in the default field of the option in the update.
	headerMap := make(map[string]*metadata.Attribute)

	for idx := range created.Header {
		headerMap[created.Header[idx].PropertyID] = &created.Header[idx]
	}

	for idx := range updated.Header {
		headerMap[updated.Header[idx].PropertyID] = &updated.Header[idx]
	}

	// updated this part is to be updated
	if err := a.ValidTableAttrDefaultValue(kit, updated.Default, headerMap); err != nil {
		blog.Errorf("valid table attr default failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	updateDataStruct.Option, updateDataStruct.ObjectID = updated, objID

	updateData, err := mapstruct.Struct2Map(updateDataStruct)
	if err != nil {
		blog.Errorf("struct to map failed: data: %+v, err: %v, rid: %s", updateDataStruct, err, kit.Rid)
		return err
	}
	if len(created.Header) == 0 && len(updateData) == 0 {
		return nil
	}

	condUpdate := mapstr.MapStr{common.BKFieldID: attrID}
	util.AddModelBizIDCondition(condUpdate, modelBizID)
	input := metadata.UpdateTableOption{Condition: condUpdate}
	createDataStruct.Option = created

	if len(created.Header) > 0 {
		input.CreateData = metadata.CreatePartDataOption{Data: []metadata.Attribute{createDataStruct}, ObjID: objID}
	}

	if len(updateData) > 0 {
		input.UpdateData = updateData
	}
	err = a.clientSet.CoreService().Model().UpdateTableModelAttrsByCondition(kit.Ctx, kit.Header, &input)
	if err != nil {
		blog.Errorf("failed to update model attr, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err := a.saveUpdateTableLog(kit, data, objID, modelBizID, attrID); err != nil {
		return err
	}

	return nil
}

func (a *attribute) saveUpdateTableLog(kit *rest.Kit, data mapstr.MapStr, objID string, modelBizID,
	attrID int64) error {
	queryCond := &metadata.QueryCondition{
		Condition: mapstr.MapStr{
			common.BKObjIDField: objID,
		},
		DisableCounter: true,
		Fields:         []string{common.BKFieldID},
	}
	objResult, err := a.clientSet.CoreService().Model().ReadModel(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("[NetDevice] search net device object, search objectName fail, %v, rid: %s", err, kit.Rid)
		return err
	}
	if len(objResult.Info) == 0 {
		blog.Errorf("[NetDevice] search net device object, search objectName fail, queryCond: %+v,err: %v, rid: %s",
			queryCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}
	if len(objResult.Info) > 1 {
		blog.Errorf("[NetDevice] search net device object, search objectName fail, queryCond: %+v,err: %v, rid: %s",
			queryCond, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}
	// save audit log.
	audit := auditlog.NewObjectAttributeAuditLog(a.clientSet.CoreService())
	generateAuditParameter := auditlog.NewGenerateAuditCommonParameter(kit, metadata.AuditUpdate).WithUpdateFields(data)
	auditLog, err := audit.GenerateTableAuditLog(generateAuditParameter, objID, modelBizID, attrID, nil)
	if err != nil {
		blog.Errorf("generate audit log failed before update model attribute, attID: %d, err: %v, rid: %s",
			attrID, err, kit.Rid)
		return err
	}

	if err := audit.SaveAuditLog(kit, *auditLog); err != nil {
		blog.Errorf("save audit log failed, attID: %d, err: %v, rid: %s", attrID, err, kit.Rid)
		return err
	}
	return nil
}

func (a *attribute) getTemplateIDByObjectAttrID(kit *rest.Kit, attrID int64, fields []string) (
	*metadata.FieldTemplateAttr, error) {

	queryCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKFieldID: attrID},
		DisableCounter: true,
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		Fields:         []string{common.BKTemplateID},
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttrByCondition(kit.Ctx, kit.Header, queryCond)
	if err != nil {
		blog.Errorf("get object attr failed, cond: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, err
	}

	if len(resp.Info) == 0 {
		blog.Errorf("no object attr founded, cond: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, "object_attr")
	}
	if len(resp.Info) > 1 {
		blog.Errorf("multi object attr founded, cond: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject, "object_attr")
	}

	cond := filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, resp.Info[0].TemplateID)

	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: cond},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
		Fields:             fields,
	}
	// list field template attributes
	res, err := a.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template attributes failed, opt: %+v, err: %v, rid: %s", listOpt, err, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, "template_attr")
	}

	if len(res.Info) > 1 {
		blog.Errorf("multi object attr founded, cond: %+v, rid: %s", listOpt, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject, "field_template_attr")
	}

	return &res.Info[0], nil
}

func (a *attribute) getModelAttrByID(kit *rest.Kit, attrID int64, bizID int64) (*metadata.Attribute, error) {
	queryCond := &metadata.QueryCondition{
		Condition:      mapstr.MapStr{common.BKFieldID: attrID},
		DisableCounter: true,
		Page:           metadata.BasePage{Limit: common.BKNoLimit},
		Fields:         []string{common.BKTemplateID},
	}
	resp, err := a.clientSet.CoreService().Model().ReadModelAttrsWithTableByCondition(kit.Ctx, kit.Header, bizID,
		queryCond)
	if err != nil {
		blog.Errorf("get object attr failed, cond: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, err
	}

	if len(resp.Info) == 0 {
		blog.Errorf("no object attr founded, cond: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommNotFound, "object_attr")
	}
	if len(resp.Info) > 1 {
		blog.Errorf("multi object attr founded, cond: %+v, err: %v, rid: %s", queryCond, err, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject, "object_attr")
	}

	return &resp.Info[0], nil
}

func (a *attribute) getFieldTemplateAttr(kit *rest.Kit, templateID int64, fields []string) (
	*metadata.FieldTemplateAttr, error) {

	cond := filtertools.GenAtomFilter(common.BKFieldID, filter.Equal, templateID)
	listOpt := &metadata.CommonQueryOption{
		CommonFilterOption: metadata.CommonFilterOption{Filter: cond},
		Page:               metadata.BasePage{Limit: common.BKNoLimit},
		Fields:             fields,
	}

	// list field template attributes
	res, err := a.clientSet.CoreService().FieldTemplate().ListFieldTemplateAttr(kit.Ctx, kit.Header, listOpt)
	if err != nil {
		blog.Errorf("list field template attributes failed, opt: %+v, err: %v, rid: %s", listOpt, err, kit.Rid)
		return nil, err
	}

	if len(res.Info) == 0 {
		return &metadata.FieldTemplateAttr{}, nil
	}

	if len(res.Info) > 1 {
		blog.Errorf("multi object attr founded, cond: %+v, rid: %s", listOpt, kit.Rid)
		return nil, kit.CCError.CCErrorf(common.CCErrCommGetMultipleObject, "field_template_attr")
	}
	return &res.Info[0], nil
}

// canAttrsUpdate coreservice has a similar logical judgment. If the logic here needs to be adjusted,
// it needs to be judged whether the logic of coreservice needs to be adjusted synchronously.
// the function name is: checkAttrTemplateInfo
func (a *attribute) canAttrsUpdate(kit *rest.Kit, input mapstr.MapStr, attrID int64, isSync bool, bizID int64) error {
	// 1. 来自字段组合模版同步操作，都可以进行修改，直接正常返回
	if isSync {
		return nil
	}

	// 2. 不是同步操作，更新属性的bk_template_id为非0时，需要报错
	data := mapstr.New()
	for k, v := range input {
		data[k] = v
	}
	newTmplID, ok := data[common.BKTemplateID]
	if ok && newTmplID != 0 {
		return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKTemplateID)
	}

	// 3. 不是同步操作，更新模型自己的属性，正常返回
	attr, err := a.getModelAttrByID(kit, attrID, bizID)
	if err != nil {
		return err
	}
	if attr.TemplateID == 0 {
		return nil
	}

	// 4. 验证来自模版的属性，是否可以正常更新
	fields := make([]string, 0)
	if _, ok := data[metadata.AttributeFieldIsRequired].(bool); ok {
		fields = append(fields, metadata.AttributeFieldIsRequired)
	}
	if _, ok := data[metadata.AttributeFieldIsEditable].(bool); ok {
		fields = append(fields, metadata.AttributeFieldIsEditable)
	}

	if _, ok := data[metadata.AttributeFieldPlaceHolder].(string); ok {
		fields = append(fields, metadata.AttributeFieldPlaceHolder)
	}

	// AttributeFieldIsRequired\AttributeFieldIsEditable\AttributeFieldPlaceHolder may be allowed
	// to be modified, the update operation does not have the above attributes to return an error
	if len(fields) == 0 {
		blog.Errorf("validate attr failed, data: %+v, rid: %s", data, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, "data")
	}

	templateAttr, err := a.getFieldTemplateAttr(kit, attr.TemplateID, fields)
	if err != nil {
		return err
	}

	// whether the corresponding lock in the attribute is false, if it is false,
	// it can be updated, otherwise it cannot be updated
	for _, field := range fields {
		switch field {
		case metadata.AttributeFieldPlaceHolder:
			if templateAttr.Placeholder.Lock {
				blog.Errorf("validate attr failed, data: %+v, field: %v, rid: %s", data, field, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, metadata.AttributeFieldPlaceHolder)
			}
		case metadata.AttributeFieldIsEditable:
			if templateAttr.Editable.Lock {
				blog.Errorf("validate attr  failed, data: %+v, field: %v, rid: %s", data, field, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, metadata.AttributeFieldIsEditable)
			}
		case metadata.AttributeFieldIsRequired:
			if templateAttr.Required.Lock {
				blog.Errorf("validate attr failed, data: %+v, field: %v rid: %s", data, field, kit.Rid)
				return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, metadata.AttributeFieldIsRequired)
			}
		}
		data.Remove(field)
	}

	removeIrrelevantValues(data)

	// After removing the above irrelevant key, check whether there is a value, and report an error if there is a value.
	if len(data) > 0 {
		blog.Errorf("validate attr failed, data: %+v, rid: %s", data, kit.Rid)
		return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, "data")
	}
	return nil
}

func removeIrrelevantValues(data mapstr.MapStr) {
	// delete irrelevant keys
	data.Remove(common.CreatorField)
	data.Remove(common.CreateTimeField)
	data.Remove(common.ModifierField)
	data.Remove(common.LastTimeField)
	data.Remove(common.BkSupplierAccount)
	data.Remove(common.BKTemplateID)
	data.Remove(common.BKFieldID)
	data.Remove(common.BKPropertyTypeField)
	data.Remove(common.BKPropertyIDField)
	data.Remove(common.BKObjIDField)
}

// UpdateObjectAttribute update object attribute
func (a *attribute) UpdateObjectAttribute(kit *rest.Kit, data mapstr.MapStr, attID int64, modelBizID int64,
	isSync bool) error {

	attr := new(metadata.Attribute)
	if err := mapstruct.Decode2Struct(data, attr); err != nil {
		blog.Errorf("unmarshal mapstr data into module failed, module: %s, err: %s, rid: %s", attr, err, kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParseDBFailed)
	}

	if err := a.isValid(kit, true, attr); err != nil {
		return err
	}

	if err := a.canAttrsUpdate(kit, data, attID, isSync, modelBizID); err != nil {
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
		IsSync:    isSync,
	}
	if _, err := a.clientSet.CoreService().Model().UpdateModelAttrsByCondition(kit.Ctx, kit.Header,
		&input); err != nil {
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

type upsertResult int64

const (
	success upsertResult = iota
	undefinedFail
	insertFail
	updateFail
)

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
	ids := make([]int64, 0)
	for id := range attributes {
		ids = append(ids, id)
	}
	sort.Sort(util.Int64Slice(ids))

	for _, idx := range ids {
		attr := attributes[idx]
		propID := attr.PropertyID
		if propID == common.BKInstParentStr {
			continue
		}

		attr.OwnerID = kit.SupplierAccount
		attr.ObjectID = objID
		if err := a.isValid(kit, true, &attr); err != nil {
			blog.Errorf("attribute(%#v) is invalid, err: %v, rid: %s", attr, err, kit.Rid)
			objRes.Errors = append(objRes.Errors, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
			hasError = true
			continue
		}

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
				grpNameIDMap[attr.PropertyGroupName] = grp.Data.GroupID
			}
		} else {
			attr.PropertyGroup = NewGroupID(true)
		}

		if result, err := a.upsertObjectAttr(kit, objID, &attr); err != nil {
			switch result {
			case undefinedFail:
				objRes.Errors = append(objRes.Errors, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
			case insertFail:
				objRes.InsertFailed = append(objRes.InsertFailed, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
			case updateFail:
				objRes.UpdateFailed = append(objRes.UpdateFailed, rowInfo{Row: idx, Info: err.Error(), PropID: propID})
			}
			hasError = true
			continue
		}

		objRes.Success = append(objRes.Success, rowInfo{Row: idx, PropID: attr.PropertyID})
	}

	return objRes, hasError
}

func (a *attribute) upsertObjectAttr(kit *rest.Kit, objID string, attr *metadata.Attribute) (upsertResult, error) {
	// check if attribute exists, if exists, update these attributes, otherwise, create the attribute
	cond := mapstr.MapStr{metadata.AttributeFieldObjectID: objID, metadata.AttributeFieldPropertyID: attr.PropertyID}
	util.AddModelBizIDCondition(cond, attr.BizID)
	queryCond := &metadata.QueryCondition{Condition: cond}
	result, err := a.clientSet.CoreService().Model().ReadModelAttrsWithTableByCondition(kit.Ctx, kit.Header, attr.BizID,
		queryCond)
	if err != nil {
		blog.Errorf("find attribute failed, err: %v, cond: %#v, rid: %s", err, queryCond, kit.Rid)
		return undefinedFail, err
	}

	if len(result.Info) == 0 {
		// create attribute
		if attr.PropertyType == common.FieldTypeInnerTable {
			if _, err := a.CreateTableObjectAttribute(kit, attr); err != nil {
				blog.Errorf("create attribute(%#v) failed, ObjID: %s, err: %v, rid: %s", attr, objID, err, kit.Rid)
				return insertFail, err
			}
			return success, nil
		}

		createAttrOpt := &metadata.CreateModelAttributes{Attributes: []metadata.Attribute{*attr}}
		_, err := a.clientSet.CoreService().Model().CreateModelAttrs(kit.Ctx, kit.Header, objID, createAttrOpt)
		if err != nil {
			blog.Errorf("create attribute(%#v) failed, ObjID: %s, err: %v, rid: %s", attr, objID, err, kit.Rid)
			return insertFail, err
		}
		return success, nil
	}

	// update attribute
	updateData := attr.ToMapStr()
	if attr.PropertyType == common.FieldTypeInnerTable {
		if err := a.UpdateTableObjectAttr(kit, updateData, result.Info[0].ID, attr.BizID); err != nil {
			blog.Errorf("failed to update module attr, err: %s, rid: %s", err, kit.Rid)
			return updateFail, err
		}
		return success, nil
	}

	// 如果属性来源于字段模版，那么无法进行更新，正常返回成功
	if result.Info[0].TemplateID != 0 {
		return success, nil
	}

	updateData.Remove(metadata.AttributeFieldPropertyID)
	updateData.Remove(metadata.AttributeFieldObjectID)
	updateData.Remove(metadata.AttributeFieldID)
	updateAttrOpt := metadata.UpdateOption{Condition: cond, Data: updateData}
	_, err = a.clientSet.CoreService().Model().UpdateModelAttrs(kit.Ctx, kit.Header, objID, &updateAttrOpt)
	if err != nil {
		blog.Errorf("failed to update module attr, err: %s, rid: %s", err, kit.Rid)
		return updateFail, err
	}
	return success, nil
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
