/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package model

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	httpheader "configcenter/src/common/http/header"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
	attrvalid "configcenter/src/common/valid/attribute"
	"configcenter/src/common/valid/attribute/manager"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	// notAddAttrModel 不允许新加属性的模型
	notAddAttrModel = map[string]bool{
		common.BKInnerObjIDPlat: true,
		common.BKInnerObjIDProc: true,
	}

	// RequiredFieldUnchangeableModels 模型的属性描述，的必填字段不允许修改
	// example: 禁止如下修改
	// db.getCollection('cc_ObjAttDes').update(
	//     {bk_obj_id: {$in: ['biz', 'host', 'set', 'module', 'plat', 'process']}},
	//     {$set: {isrequired: true}}
	// )
	RequiredFieldUnchangeableModels = map[string]bool{
		common.BKInnerObjIDApp:    true,
		common.BKInnerObjIDHost:   true,
		common.BKInnerObjIDSet:    true,
		common.BKInnerObjIDModule: true,
		common.BKInnerObjIDPlat:   true,
		common.BKInnerObjIDProc:   true,
	}
)

// Count TODO
func (m *modelAttribute) Count(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {
	cnt, err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).Count(kit.Ctx)
	return cnt, err
}

func (m *modelAttribute) saveTableAttr(kit *rest.Kit, attribute metadata.Attribute) (id uint64, err error) {

	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjAttDes)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}

	index, err := m.GetAttrLastIndex(kit, attribute)
	if err != nil {
		return id, err
	}

	attribute.PropertyIndex = index
	attribute.ID = int64(id)
	attribute.OwnerID = kit.SupplierAccount

	if attribute.CreateTime == nil {
		attribute.CreateTime = &metadata.Time{}
		attribute.CreateTime.Time = time.Now()
	}

	if attribute.LastTime == nil {
		attribute.LastTime = &metadata.Time{}
		attribute.LastTime.Time = time.Now()
	}

	if attribute.IsMultiple == nil {
		isMultiple := false
		attribute.IsMultiple = &isMultiple
	}

	if err = m.saveTableAttrCheck(kit, attribute); err != nil {
		blog.Errorf("save table attr failed, attribute: %v, err: %v, rid: %s", attribute, err, kit.Rid)
		return 0, err
	}
	if err = mongodb.Client().Table(common.BKTableNameObjAttDes).Insert(kit.Ctx, attribute); err != nil {
		blog.Errorf("save table attr failed, attr: %v, err: %v, rid: %s", attribute, err, kit.Rid)
		return id, err
	}

	return id, nil
}

func (m *modelAttribute) save(kit *rest.Kit, attribute metadata.Attribute) (id uint64, err error) {
	id, err = mongodb.Client().NextSequence(kit.Ctx, common.BKTableNameObjAttDes)
	if err != nil {
		return id, kit.CCError.New(common.CCErrObjectDBOpErrno, err.Error())
	}
	index, err := m.GetAttrLastIndex(kit, attribute)
	if err != nil {
		return id, err
	}

	attribute.PropertyIndex = index
	attribute.ID = int64(id)
	attribute.OwnerID = kit.SupplierAccount

	if attribute.CreateTime == nil {
		attribute.CreateTime = &metadata.Time{}
		attribute.CreateTime.Time = time.Now()
	}
	if attribute.LastTime == nil {
		attribute.LastTime = &metadata.Time{}
		attribute.LastTime.Time = time.Now()
	}
	if attribute.IsMultiple == nil {
		switch attribute.PropertyType {
		case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeInt, common.FieldTypeFloat,
			common.FieldTypeEnum, common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeTimeZone,
			common.FieldTypeBool, common.FieldTypeList, common.FieldTypeIDRule:
			isMultiple := false
			attribute.IsMultiple = &isMultiple
		case common.FieldTypeUser, common.FieldTypeOrganization, common.FieldTypeEnumQuote, common.FieldTypeEnumMulti:
			isMultiple := true
			attribute.IsMultiple = &isMultiple
		default:
			return 0, kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.AttributeFieldPropertyType)
		}
	}
	// 对于枚举，枚举多选，枚举引用字段, 默认值是放在option中的，需要将default置为nil
	if attribute.Default != nil && (attribute.PropertyType == common.FieldTypeEnum ||
		attribute.PropertyType == common.FieldTypeEnumMulti || attribute.PropertyType == common.FieldTypeEnumQuote) {

		attribute.Default = nil
	}
	if err = m.saveCheck(kit, attribute); err != nil {
		return 0, err
	}

	if err = mongodb.Client().Table(common.BKTableNameObjAttDes).Insert(kit.Ctx, attribute); err != nil {
		return 0, err
	}
	if attribute.PropertyType == common.FieldTypeIDRule {
		idx := types.Index{
			Name:       common.CCLogicIndexNamePrefix + attribute.PropertyID,
			Keys:       bson.D{{attribute.PropertyID, 1}},
			Background: true,
			Unique:     true,
			PartialFilterExpression: map[string]interface{}{
				attribute.PropertyID: map[string]string{common.BKDBType: "string", common.BKDBGT: ""},
			},
		}
		table := common.GetInstTableName(attribute.ObjectID, kit.SupplierAccount)
		if err = mongodb.Client().Table(table).CreateIndex(kit.Ctx, idx); err != nil {
			blog.Errorf("create index failed, index: %+v, err: %v, rid: %s", idx, err, kit.Rid)
			return 0, kit.CCError.Error(common.CCErrObjectDBOpErrno)
		}

		unique := metadata.ObjectUnique{
			ID:       id,
			ObjID:    attribute.ObjectID,
			Keys:     []metadata.UniqueKey{{Kind: metadata.UniqueKeyKindProperty, ID: uint64(attribute.ID)}},
			Ispre:    false,
			OwnerID:  kit.SupplierAccount,
			LastTime: metadata.Now(),
		}
		err = mongodb.Client().Table(common.BKTableNameObjUnique).Insert(kit.Ctx, &unique)
		if nil != err {
			blog.Errorf("create unique failed, val: %+v, err: %v, rid: %s", &unique, err, kit.Rid)
			return 0, kit.CCError.Error(common.CCErrObjectDBOpErrno)
		}
	}
	return id, nil
}

func (m *modelAttribute) checkUnique(kit *rest.Kit, isCreate bool, objID, propertyID, propertyName string,
	modelBizID int64) error {
	cond := map[string]interface{}{
		common.BKObjIDField: objID,
	}

	andCond := make([]map[string]interface{}, 0)
	if isCreate {
		nameFieldCond := map[string]interface{}{common.BKPropertyNameField: propertyName}
		idFieldCond := map[string]interface{}{common.BKPropertyIDField: propertyID}
		andCond = append(andCond, map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{nameFieldCond, idFieldCond},
		})
	} else {
		// update attribute. not change name, 无需判断
		if propertyName == "" {
			return nil
		}
		cond[common.BKPropertyIDField] = map[string]interface{}{common.BKDBNE: propertyID}
		cond[common.BKPropertyNameField] = propertyName
	}

	if modelBizID > 0 {
		// search special business model and global shared model
		andCond = append(andCond, map[string]interface{}{
			common.BKDBOR: []map[string]interface{}{
				{common.BKAppIDField: modelBizID},
				{common.BKAppIDField: 0},
				{common.BKAppIDField: map[string]interface{}{common.BKDBExists: false}},
			},
		})
	}

	if len(andCond) > 0 {
		cond[common.BKDBAND] = andCond
	}
	util.SetModOwner(cond, kit.SupplierAccount)

	resultAttrs := make([]metadata.Attribute, 0)
	err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &resultAttrs)
	blog.V(5).Infof("checkUnique db cond:%#v, result:%#v, rid:%s", cond, resultAttrs, kit.Rid)
	if err != nil {
		blog.ErrorJSON("checkUnique select error. err:%s, cond:%s, rid:%s", err.Error(), cond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	language := httpheader.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	for _, attrItem := range resultAttrs {
		if attrItem.PropertyID == propertyID {
			blog.ErrorJSON("check unique attribute id duplicate. attr: %s, rid: %s", attrItem, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommDuplicateItem, lang.Language("model_attr_bk_property_id"))
		}
		if attrItem.PropertyName == propertyName {
			blog.ErrorJSON("check unique attribute id duplicate. attr: %s, rid: %s", attrItem, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommDuplicateItem, lang.Language("model_attr_bk_property_name"))
		}
	}

	return nil
}

func (m *modelAttribute) checkTableAttributeMustNotEmpty(kit *rest.Kit, attribute metadata.Attribute) error {
	if attribute.PropertyID == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyID)
	}
	if attribute.PropertyName == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyName)
	}
	if attribute.PropertyType != common.FieldTypeInnerTable {
		return kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.AttributeFieldPropertyType)
	}
	return nil
}

func (m *modelAttribute) checkAttributeMustNotEmpty(kit *rest.Kit, attribute metadata.Attribute) error {
	if attribute.PropertyID == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyID)
	}
	if attribute.PropertyName == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyName)
	}
	if attribute.PropertyType == "" {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
	}

	return nil
}

func (m *modelAttribute) checkTableAttributeValidity(kit *rest.Kit, attribute metadata.Attribute) error {

	lang := m.language.CreateDefaultCCLanguageIf(httpheader.GetLanguage(kit.Header))

	if attribute.PropertyID != "" {
		attribute.PropertyID = strings.TrimSpace(attribute.PropertyID)
		if common.AttributeIDMaxLength < utf8.RuneCountInString(attribute.PropertyID) {

			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_id"),
				common.AttributeIDMaxLength)
		}

		if !SatisfyMongoFieldLimit(attribute.PropertyID) {
			blog.Errorf("attribute property id:(%s) not satisfy mongo field limit, rid: %s",
				attribute.PropertyID, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
		}

		// check only preset attribute's property id can start with bk_ or _bk
		if !attribute.IsPre {
			if strings.HasPrefix(attribute.PropertyID, "bk_") ||
				strings.HasPrefix(attribute.PropertyID, "_bk") {
				return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
			}
		}
	}

	attribute.PropertyName = strings.TrimSpace(attribute.PropertyName)
	if common.AttributeNameMaxLength < utf8.RuneCountInString(attribute.PropertyName) {
		blog.Errorf("attribute property name is exceed max limit:(%d), rid: %s", common.AttributeNameMaxLength, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_name"),
			common.AttributeNameMaxLength)
	}

	if attribute.Placeholder != "" {
		attribute.Placeholder = strings.TrimSpace(attribute.Placeholder)

		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(attribute.Placeholder) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_placeholder"),
				common.AttributePlaceHolderMaxLength)
		}
		match, err := regexp.MatchString(common.FieldTypeLongCharRegexp, attribute.Placeholder)
		if nil != err || !match {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPlaceHolder)
		}
	}

	if attribute.Unit != "" {
		attribute.Unit = strings.TrimSpace(attribute.Unit)
		if common.AttributeUnitMaxLength < utf8.RuneCountInString(attribute.Unit) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_uint"),
				common.AttributeUnitMaxLength)
		}
	}

	if attribute.PropertyType != common.FieldTypeInnerTable {
		blog.Errorf("attr property type is error, property is : %s", attribute.PropertyType, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
	}

	tableOption, err := metadata.ParseTableAttrOption(attribute.Option)
	if err != nil {
		blog.Errorf("get attribute option failed, error: %v, option: %v, rid: %s", err, kit.Rid)
		return err
	}

	if len(tableOption.Header) == 0 {
		blog.Errorf("table attribute option invalid, header is nil, tableOption: %+v, rid: %s", tableOption, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "table header")
	}

	if err := m.checkTableAttr(kit, attribute.PropertyID, attribute.ObjectID, tableOption); err != nil {
		blog.Errorf("check table attribute failed, tableOption: %+v, err: %v, rid: %s", tableOption, err, kit.Rid)
		return err
	}
	return nil
}

func (m *modelAttribute) validAndGetTableAttrHeaderDetail(kit *rest.Kit, header []metadata.Attribute) (
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
			return nil, kit.CCError.Errorf(common.CCErrCommXXExceedLimit, "table header long char",
				metadata.TableLongCharMaxNum)
		}

		// check if property type for creation is valid, can't update property type
		if header[index].PropertyType == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyType)
		}

		if header[index].PropertyID == "" {
			return nil, kit.CCError.Errorf(common.CCErrCommParamsNeedSet, metadata.AttributeFieldPropertyID)
		}

		if common.AttributeIDMaxLength < utf8.RuneCountInString(header[index].PropertyID) {
			return nil, kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, common.AttributeIDMaxLength)
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
			return nil, kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, common.AttributeNameMaxLength)
		}

		if err = attrvalid.ValidTableFieldOption(kit, header[index].PropertyType, header[index].Option,
			header[index].Default, header[index].IsMultiple, header[index].ObjectID); err != nil {
			return nil, err
		}
		propertyAttr[header[index].PropertyID] = &header[index]
	}
	return propertyAttr, nil
}

func (m *modelAttribute) checkTableAttr(kit *rest.Kit, propertyID, objectID string,
	tableOption *metadata.TableAttributesOption) error {

	if len(tableOption.Header) == 0 {
		blog.Errorf("table attribute option invalid, header is nil, tableOption: %+v, rid: %s", tableOption, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "table header")
	}

	headerMap, err := m.validAndGetTableAttrHeaderDetail(kit, tableOption.Header)
	if err != nil {
		blog.Errorf("failed to valid the header of the table, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	// 这里还得获取一下数据库中的内容 因为可能这个header只是新加的，在更新场景下还得把以前的header内容拿出来 没有加进来
	input := mapstr.MapStr{
		common.BKPropertyIDField:   propertyID,
		common.BKObjIDField:        objectID,
		common.BKPropertyTypeField: common.FieldTypeInnerTable,
	}
	attrResult, err := m.newSearch(kit, input)
	if err != nil {
		blog.Errorf("failed to search the attr of the model, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	if len(attrResult) != 1 {
		blog.Errorf("failed to search the attributes of the model, input: %+v, err: %v, rid: %s", input, err, kit.Rid)
		return err
	}

	op, err := metadata.ParseTableAttrOption(attrResult[0].Option)
	if err != nil {
		blog.Errorf("get attribute option failed, error: %v, option: %v, rid: %s", err, kit.Rid)
		return err
	}

	for index := range op.Header {
		if _, ok := headerMap[op.Header[index].PropertyID]; !ok {
			headerMap[op.Header[index].PropertyID] = &op.Header[index]
		}
	}

	if tableOption.Default != nil {
		for _, value := range tableOption.Default {
			for k, v := range value {
				if err := m.checkTableAttributeDefaultValue(kit, headerMap[k].Option, v,
					headerMap[k].PropertyType); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

var validAttrPropertyTypes = map[string]struct{}{
	common.FieldTypeSingleChar:   {},
	common.FieldTypeLongChar:     {},
	common.FieldTypeInt:          {},
	common.FieldTypeFloat:        {},
	common.FieldTypeEnum:         {},
	common.FieldTypeEnumMulti:    {},
	common.FieldTypeDate:         {},
	common.FieldTypeTime:         {},
	common.FieldTypeUser:         {},
	common.FieldTypeOrganization: {},
	common.FieldTypeTimeZone:     {},
	common.FieldTypeBool:         {},
	common.FieldTypeList:         {},
	common.FieldTypeEnumQuote:    {},
	common.FieldTypeIDRule:       {},
}

func (m *modelAttribute) checkAttributeValidity(kit *rest.Kit, attribute metadata.Attribute,
	propertyType string) error {
	language := httpheader.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	if attribute.PropertyID != "" {
		if err := m.validateAttrPropertyID(kit, attribute, lang); err != nil {
			return err
		}
	}

	if attribute.PropertyName = strings.TrimSpace(attribute.PropertyName); common.AttributeNameMaxLength <
		utf8.RuneCountInString(attribute.PropertyName) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_name"),
			common.AttributeNameMaxLength)
	}

	if attribute.Placeholder != "" {
		attribute.Placeholder = strings.TrimSpace(attribute.Placeholder)

		if common.AttributePlaceHolderMaxLength < utf8.RuneCountInString(attribute.Placeholder) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_placeholder"),
				common.AttributePlaceHolderMaxLength)
		}
		match, err := regexp.MatchString(common.FieldTypeLongCharRegexp, attribute.Placeholder)
		if err != nil || !match {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPlaceHolder)

		}
	}

	if attribute.Unit != "" {
		attribute.Unit = strings.TrimSpace(attribute.Unit)
		if common.AttributeUnitMaxLength < utf8.RuneCountInString(attribute.Unit) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_uint"),
				common.AttributeUnitMaxLength)
		}
	}

	if attribute.PropertyType != "" {
		if _, exists := validAttrPropertyTypes[attribute.PropertyType]; !exists {
			if _, ok := manager.Get(attribute.PropertyType); !ok {
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
			}
		}
	}

	if attribute.Default != nil && propertyType != common.FieldTypeEnum && propertyType != common.FieldTypeEnumMulti &&
		propertyType != common.FieldTypeEnumQuote && propertyType != common.FieldTypeIDRule {

		if err := m.checkAttributeDefaultValue(kit, attribute, propertyType); err != nil {
			return err
		}
	}

	if opt, ok := attribute.Option.(string); ok && opt != "" {
		if common.AttributeOptionMaxLength < utf8.RuneCountInString(opt) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_option_regex"),
				common.AttributeOptionMaxLength)
		}
	}

	return nil
}

func (m *modelAttribute) validateAttrPropertyID(kit *rest.Kit, attribute metadata.Attribute,
	lang language.DefaultCCLanguageIf) error {

	attribute.PropertyID = strings.TrimSpace(attribute.PropertyID)
	if common.AttributeIDMaxLength < utf8.RuneCountInString(attribute.PropertyID) {
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_id"),
			common.AttributeIDMaxLength)
	}

	if !SatisfyMongoFieldLimit(attribute.PropertyID) {
		blog.Errorf("attribute.PropertyID: %s not satisfy mongo field limit", attribute.PropertyID)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
	}

	// check only preset attribute's property id can start with bk_ or _bk
	if !attribute.IsPre {
		if strings.HasPrefix(attribute.PropertyID, "bk_") || strings.HasPrefix(attribute.PropertyID, "_bk") {
			return kit.CCError.CCErrorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyID)
		}
	}
	return nil
}

func (m *modelAttribute) checkTableAttributeDefaultValue(kit *rest.Kit, option, defautValue interface{},
	propertyType string) error {

	switch propertyType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		if err := attrvalid.ValidFieldTypeString(kit, option, defautValue); err != nil {
			return err
		}
	case common.FieldTypeInt:
		if err := attrvalid.ValidFieldTypeInt(kit, option, defautValue); err != nil {
			return err
		}
	case common.FieldTypeFloat:
		if err := attrvalid.ValidFieldTypeFloat(kit, option, defautValue); err != nil {
			return err
		}
	case common.FieldTypeBool:
		if err := valid.ValidateBoolType(defautValue); err != nil {
			blog.Errorf("bool type default value not bool, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	case common.FieldTypeEnumMulti:
		// 默认值相关的检查都是按照最宽松的进行校验
		isMulti := true
		if err := attrvalid.ValidFieldTypeEnumOption(kit, option, &isMulti); err != nil {
			blog.Errorf("enum multi type default value not enum multi, err: %v, rid: %s", err, kit.Rid)
			return err
		}
	default:
		blog.Errorf("property type is error, propertyType: %v, rid: %s", propertyType, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
	}

	return nil
}

// checkAttributeDefaultValue 校验属性的default字段，对于枚举，枚举多选，枚举引用字段, 默认值是放在option中的，不能调用该函数校验
func (m *modelAttribute) checkAttributeDefaultValue(kit *rest.Kit, attribute metadata.Attribute,
	propertyType string) error {

	var err error
	switch propertyType {
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		err = attrvalid.ValidFieldTypeString(kit, attribute.Option, attribute.Default)
	case common.FieldTypeInt:
		err = attrvalid.ValidFieldTypeInt(kit, attribute.Option, attribute.Default)
	case common.FieldTypeFloat:
		err = attrvalid.ValidFieldTypeFloat(kit, attribute.Option, attribute.Default)
	case common.FieldTypeDate:
		if ok := util.IsDate(attribute.Default); !ok {
			return fmt.Errorf("date default value is not date type, type: %T", attribute.Default)
		}
	case common.FieldTypeTime:
		if _, ok := util.IsTime(attribute.Default); !ok {
			return fmt.Errorf("time default value formart is not time string, type: %T", attribute.Default)
		}
	case common.FieldTypeUser:
		err = m.checkUserTypeDefaultValue(kit, attribute)
	case common.FieldTypeOrganization:
		err = m.checkOrganizationTypeDefaultValue(kit, attribute)
	case common.FieldTypeTimeZone:
		if ok := util.IsTimeZone(attribute.Default); !ok {
			return fmt.Errorf("time zone default value is not time zone type, type: %T", attribute.Default)
		}
	case common.FieldTypeBool:
		err = valid.ValidateBoolType(attribute.Default)
		blog.Errorf("bool type default value not bool, err: %v, rid: %s", err, kit.Rid)

	case common.FieldTypeList:
		err = attrvalid.ValidFieldTypeList(kit, attribute.Option, attribute.Default)

	default:
		if propertyType == common.FieldTypeEnum || propertyType == common.FieldTypeEnumMulti ||
			propertyType == common.FieldTypeEnumQuote {
			return fmt.Errorf("enum, enummulti, enumquote type default field is nil")
		}
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
	}

	if err != nil {
		return err
	}

	return nil
}

func (m *modelAttribute) checkUserTypeDefaultValue(kit *rest.Kit, attribute metadata.Attribute) error {
	switch value := attribute.Default.(type) {
	case string:
		value = strings.TrimSpace(value)
		if len(value) > common.FieldTypeUserLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeUserLenChar, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
		}

		if len(value) == 0 {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
		}

		// regex check
		match := util.IsUser(value)
		if !match {
			blog.Errorf(`value "%s" not match regexp, rid: %s`, value, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
		}
	default:
		blog.Errorf("user type default value is invalid, defaultVal: %+v, rid: %s", attribute.Default, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}
	return nil
}

func (m *modelAttribute) checkOrganizationTypeDefaultValue(kit *rest.Kit, attribute metadata.Attribute) error {
	switch org := attribute.Default.(type) {
	case []interface{}:
		if len(org) == 0 {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
		}

		for _, orgID := range org {
			if !util.IsInteger(orgID) {
				blog.Errorf("orgID params not int, type: %T, rid: %s", orgID, kit.Rid)
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
			}
		}
	default:
		blog.Errorf("org type default value is invalid, defaultVal: %+v, rid: %s", attribute.Default, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}
	return nil
}

func (m *modelAttribute) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition, isSync bool) (
	cnt uint64, err error) {

	err = m.checkUpdate(kit, data, cond, isSync)
	if err != nil {
		blog.ErrorJSON("checkUpdate error. data:%s, cond:%s, rid:%s", data, cond, kit.Rid)
		return cnt, err
	}
	cnt, err = mongodb.Client().Table(common.BKTableNameObjAttDes).UpdateMany(kit.Ctx, cond.ToMapStr(), data)
	if nil != err {
		blog.Errorf("request(%s): database operation is failed, error info is %s", kit.Rid, err.Error())
		return 0, err
	}

	return cnt, err
}

func (m *modelAttribute) newSearch(kit *rest.Kit, cond mapstr.MapStr) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) search(kit *rest.Kit, cond universalsql.Condition) (resultAttrs []metadata.Attribute,
	err error) {
	resultAttrs = []metadata.Attribute{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) searchWithSort(kit *rest.Kit, cond metadata.QueryCondition) (resultAttrs []metadata.Attribute,
	err error) {
	resultAttrs = []metadata.Attribute{}

	instHandler := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.Condition)
	err = instHandler.Start(uint64(cond.Page.Start)).Limit(uint64(cond.Page.Limit)).Sort(cond.Page.Sort).All(kit.Ctx,
		&resultAttrs)

	return resultAttrs, err
}

func (m *modelAttribute) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) (resultAttrs []mapstr.MapStr,
	err error) {

	resultAttrs = []mapstr.MapStr{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

// delete delete the model scene isMode is true, no need to check whether
// the isFromModel field inherits from the field template
func (m *modelAttribute) delete(kit *rest.Kit, cond universalsql.Condition, isFromModel bool) (cnt uint64, err error) {

	resultAttrs := make([]metadata.Attribute, 0)
	fields := []string{common.BKFieldID, common.BKPropertyIDField, common.BKPropertyTypeField,
		common.BKObjIDField, common.BKAppIDField}

	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(condMap).Fields(fields...).All(kit.Ctx, &resultAttrs)
	if nil != err {
		blog.Errorf("request(%s): database count operation is failed, error info is %s", kit.Rid, err.Error())
		return 0, err
	}

	cnt = uint64(len(resultAttrs))
	if cnt == 0 {
		return cnt, nil
	}

	objIDArrMap := make(map[string][]int64, 0)
	idRuleAttrMap := make(map[string][]int64, 0)
	for _, attr := range resultAttrs {
		if attr.PropertyType == common.FieldTypeInnerTable {
			blog.Errorf("property is error, attrItem: %+v, rid: %s", attr, kit.Rid)
			return 0, kit.CCError.New(common.CCErrTopoObjectSelectFailed, common.BKPropertyTypeField)
		}

		if !isFromModel && attr.TemplateID != 0 {
			return 0, kit.CCError.CCErrorf(common.CCErrorTopoFieldTemplateForbiddenDeleteAttr, attr.ID, attr.TemplateID)
		}
		if attr.PropertyType == common.FieldTypeIDRule {
			idRuleAttrMap[attr.ObjectID] = append(idRuleAttrMap[attr.ObjectID], attr.ID)
			continue
		}

		objIDArrMap[attr.ObjectID] = append(objIDArrMap[attr.ObjectID], attr.ID)
	}

	if err := m.cleanAttributeFieldInInstances(kit, resultAttrs); err != nil {
		blog.Errorf("delete object attributes with cond: %v, but delete these attribute in instance failed, "+
			"err: %v, rid: %s", condMap, err, kit.Rid)
		return 0, err
	}

	// delete template attribute when delete model attribute
	if err := m.cleanAttrTemplateRelation(kit, resultAttrs); err != nil {
		blog.Errorf("delete the relation between attributes and templates failed, attr: %v, err: %v, rid: %s",
			resultAttrs, err, kit.Rid)
		return 0, err
	}

	if len(objIDArrMap) != 0 {
		exist, err := m.checkAttributeInUnique(kit, objIDArrMap)
		if err != nil {
			blog.ErrorJSON("check attribute in unique error. err:%s, input:%s, rid:%s", err.Error(), condMap, kit.Rid)
			return 0, err
		}
		// delete field in module unique. not allow delete
		if exist {
			blog.ErrorJSON("delete field in unique. delete cond:%s, field:%s, rid:%s", condMap, resultAttrs, kit.Rid)
			return 0, kit.CCError.Error(common.CCErrCoreServiceNotAllowUniqueAttr)
		}
	}

	if len(idRuleAttrMap) != 0 {
		if err = m.delIDRuleUnique(kit, idRuleAttrMap); err != nil {
			return 0, err
		}
	}

	deleteCnt, err := mongodb.Client().Table(common.BKTableNameObjAttDes).DeleteMany(kit.Ctx, condMap)
	if nil != err {
		blog.Errorf("request(%s): database deletion operation is failed, error info is %s", kit.Rid, err.Error())
		return deleteCnt, err
	}

	return cnt, err
}

type bizObjectFields struct {
	bizID  int64
	fields []string
}

// cleanAttributeFieldInInstances remove attribute filed in this object's instances
func (m *modelAttribute) cleanAttributeFieldInInstances(kit *rest.Kit, attrs []metadata.Attribute) error {
	// this operation may take a long time, do not use transaction
	kit.Ctx = context.Background()

	objectFields, hostApplyFields, err := m.getObjAndHostApplyFields(kit, attrs)
	if err != nil {
		return err
	}

	// delete these attribute's fields in the model instance
	var hitError error
	wg := sync.WaitGroup{}
	for object, objFields := range objectFields {
		if len(objFields) == 0 {
			// no fields need to be removed, skip directly.
			continue
		}

		for _, objField := range objFields {
			fields := objField.fields
			existConds := make([]map[string]interface{}, len(fields))

			for index, field := range fields {
				existConds[index] = map[string]interface{}{
					field: map[string]interface{}{
						common.BKDBExists: true,
					},
				}
			}

			cond := map[string]interface{}{common.BKDBOR: existConds}

			if objField.bizID > 0 {
				if !isBizObject(object) {
					return fmt.Errorf("unsupported object %s's clean instance field operation", object)
				}

				if object == common.BKInnerObjIDHost {
					if err := m.cleanHostAttributeField(kit.Ctx, kit.SupplierAccount, objField); err != nil {
						return err
					}
					continue
				}

				cond[common.BKAppIDField] = objField.bizID
			} else {
				if isBizObject(object) {
					if object == common.BKInnerObjIDHost {
						ele := bizObjectFields{
							bizID:  0,
							fields: fields,
						}
						if err := m.cleanHostAttributeField(kit.Ctx, kit.SupplierAccount, ele); err != nil {
							return err
						}
						continue
					}
				} else {
					cond[common.BKObjIDField] = object
				}
			}

			cond = util.SetQueryOwner(cond, kit.SupplierAccount)

			collectionName := common.GetInstTableName(object, kit.SupplierAccount)
			wg.Add(1)
			go func(collName string, filter types.Filter, fields []string) {
				defer wg.Done()
				hitError = m.dropColumns(kit, object, collName, filter, fields)

			}(collectionName, cond, fields)
		}
	}
	// wait for all the public object routine is done.
	wg.Wait()
	if hitError != nil {
		return hitError
	}

	// wait for all the public object routine is done.
	wg.Wait()
	if hitError != nil {
		return hitError
	}

	// step 3: clean host apply fields
	if err := m.cleanHostApplyField(kit.Ctx, kit.SupplierAccount, hostApplyFields); err != nil {
		return err
	}

	return nil
}

// getObjAndHostApplyFields TODO: now, we only support set, module, host model's biz attribute clean operation.
func (m *modelAttribute) getObjAndHostApplyFields(kit *rest.Kit, attrs []metadata.Attribute) (
	map[string][]bizObjectFields, map[int64][]int64, error) {
	objectFields := make(map[string][]bizObjectFields, 0)
	hostApplyFields := make(map[int64][]int64)
	for _, attr := range attrs {
		if attr.PropertyType == common.FieldTypeInnerTable {
			blog.Errorf("property is error, attrItem: %+v, rid: %s", attr, kit.Rid)
			return nil, nil, kit.CCError.New(common.CCErrTopoObjectSelectFailed, common.BKPropertyTypeField)
		}
		biz := attr.BizID
		if biz != 0 {
			if !isBizObject(attr.ObjectID) {
				return nil, nil, fmt.Errorf("unsupported object %s's clean instance field operation", attr.ObjectID)
			}
		}

		_, exist := objectFields[attr.ObjectID]
		if !exist {
			objectFields[attr.ObjectID] = make([]bizObjectFields, 0)
		}
		objectFields[attr.ObjectID] = append(objectFields[attr.ObjectID], bizObjectFields{
			bizID:  biz,
			fields: []string{attr.PropertyID},
		})

		if attr.ObjectID == common.BKInnerObjIDHost {
			hostApplyFields[biz] = append(hostApplyFields[biz], attr.ID)
		}
	}
	return objectFields, hostApplyFields, nil
}

func (m *modelAttribute) dropColumns(kit *rest.Kit, object, collName string, filter types.Filter,
	fields []string) error {
	instCount, err := mongodb.Client().Table(collName).Find(filter).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count instances with the attribute to delete failed, table: %s, cond: %v, fields: %v, err: %v, "+
			"rid: %s", collName, filter, fields, err, kit.Rid)
		return err
	}

	instIDField := common.GetInstIDField(object)
	for start := uint64(0); start < instCount; start += pageSize {
		insts := make([]map[string]interface{}, 0)
		err := mongodb.Client().Table(collName).Find(filter).Start(0).Limit(pageSize).Fields(instIDField).
			All(kit.Ctx, &insts)
		if err != nil {
			blog.Errorf("get instance ids with the attr to delete failed, table: %s, cond: %v, fields: %v, err: %v, "+
				"rid: %s", collName, filter, fields, err, kit.Rid)
			return err
		}

		if len(insts) == 0 {
			return nil
		}

		instIDs := make([]int64, len(insts))
		for index, inst := range insts {
			instID, err := util.GetInt64ByInterface(inst[instIDField])
			if err != nil {
				blog.Errorf("get instance id failed, inst: %+v, err: %v, rid: %s", inst, err, kit.Rid)
				return err
			}
			instIDs[index] = instID
		}

		instFilter := map[string]interface{}{
			instIDField: map[string]interface{}{
				common.BKDBIN: instIDs,
			},
		}

		if err := mongodb.Client().Table(collName).DropColumns(kit.Ctx, instFilter, fields); err != nil {
			blog.Error("delete object's attribute from instance failed, table: %s, cond: %v, fields: %v, err: %v, "+
				"rid: %s", collName, instFilter, fields, err, kit.Rid)
			return err
		}
	}
	return nil
}

func (m *modelAttribute) cleanAttrTemplateRelation(kit *rest.Kit, attrs []metadata.Attribute) error {

	if len(attrs) == 0 {
		return nil
	}

	attrMap := make(map[string][]int64)
	for _, attr := range attrs {
		if attr.PropertyType == common.FieldTypeInnerTable {
			blog.Errorf("property is error, attrItem: %+v, rid: %s", attr, kit.Rid)
			return kit.CCError.New(common.CCErrTopoObjectSelectFailed, common.BKPropertyTypeField)
		}
		attrMap[attr.ObjectID] = append(attrMap[attr.ObjectID], attr.ID)
	}

	for objID, attrIDs := range attrMap {
		cond := mapstr.MapStr{
			common.BKAttributeIDField: mapstr.MapStr{common.BKDBIN: attrIDs},
		}
		cond = util.SetQueryOwner(cond, kit.SupplierAccount)
		switch objID {
		case common.BKInnerObjIDSet:
			if err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Delete(kit.Ctx, cond); err != nil {
				return err
			}

		case common.BKInnerObjIDModule:
			if err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Delete(kit.Ctx, cond); err != nil {
				return err
			}
		}
	}

	return nil
}

const pageSize = 2000

func (m *modelAttribute) cleanHostAttributeField(ctx context.Context, ownerID string, info bizObjectFields) error {
	cond := mapstr.MapStr{}
	cond = util.SetQueryOwner(cond, ownerID)
	// biz id = 0 means all the hosts.
	// TODO: optimize when the filed is a public filed in all the host instances. handle with page
	if info.bizID != 0 {
		// find hosts in this biz
		cond = mapstr.MapStr{
			common.BKAppIDField: info.bizID,
		}
	}

	hostCount, err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(cond).Count(ctx)
	if err != nil {
		return err
	}

	type hostInst struct {
		HostID int64 `bson:"bk_host_id"`
	}

	fields := info.fields
	existConds := make([]map[string]interface{}, len(fields))

	for index, field := range fields {
		existConds[index] = map[string]interface{}{
			field: map[string]interface{}{
				common.BKDBExists: true,
			},
		}
	}

	for start := uint64(0); start < hostCount; start += pageSize {
		hostList := make([]hostInst, 0)
		err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(cond).Start(start).Limit(pageSize).Fields(common.BKHostIDField).All(ctx,
			&hostList)
		if err != nil {
			return err
		}

		if len(hostList) == 0 {
			return nil
		}

		ids := make([]int64, len(hostList))
		for index, host := range hostList {
			ids[index] = host.HostID
		}

		hostFilter := mapstr.MapStr{
			common.BKHostIDField: mapstr.MapStr{common.BKDBIN: ids},
			common.BKDBOR:        existConds,
		}
		if err := mongodb.Client().Table(common.BKTableNameBaseHost).DropColumns(ctx, hostFilter,
			info.fields); err != nil {
			return fmt.Errorf("clean host biz attribute %v failed, err: %v", info.fields, err)
		}
	}

	return nil

}

func (m *modelAttribute) cleanHostApplyField(ctx context.Context, ownerID string,
	hostApplyFields map[int64][]int64) error {
	orCond := make([]map[string]interface{}, 0)
	for bizID, attrIDs := range hostApplyFields {
		attrCond := map[string]interface{}{
			common.BKAttributeIDField: map[string]interface{}{
				common.BKDBIN: attrIDs,
			},
		}
		// global attribute requires removing host apply rules for all biz
		if bizID != 0 {
			attrCond[common.BKAppIDField] = bizID
		}
		orCond = append(orCond, attrCond)
	}
	if len(orCond) == 0 {
		return nil
	}
	cond := make(map[string]interface{})
	cond = util.SetQueryOwner(cond, ownerID)
	cond[common.BKDBOR] = orCond
	if err := mongodb.Client().Table(common.BKTableNameHostApplyRule).Delete(ctx, cond); err != nil {
		blog.ErrorJSON("cleanHostApplyField failed, err: %s, cond: %s", err, cond)
		return err
	}
	return nil

}

// isBizObject TODO
// now, we only support set, module, host model's biz attribute clean operation.
func isBizObject(objectID string) bool {
	switch objectID {
	// biz is a special object, but it can not have biz attribute obviously.
	case common.BKInnerObjIDApp:
		return true
	case common.BKInnerObjIDHost:
		return true
	case common.BKInnerObjIDModule:
		return true
	case common.BKInnerObjIDSet:
		return true
	default:
		// TODO: remove this when the common object support biz attribute and biz instance field.
		return false

	}
}

// saveTableAttrCheck form new field check
func (m *modelAttribute) saveTableAttrCheck(kit *rest.Kit, attribute metadata.Attribute) error {
	if err := m.checkTableAttributeMustNotEmpty(kit, attribute); err != nil {
		return err
	}
	if err := m.checkTableAttributeValidity(kit, attribute); err != nil {
		return err
	}
	return nil
}

// saveCheck 新加字段检查
func (m *modelAttribute) saveCheck(kit *rest.Kit, attr metadata.Attribute) error {
	if err := m.checkAddField(kit, attr); err != nil {
		return err
	}

	if err := m.checkAttributeMustNotEmpty(kit, attr); err != nil {
		return err
	}
	if err := m.checkAttributeValidity(kit, attr, attr.PropertyType); err != nil {
		return err
	}

	// check name duplicate
	if err := m.checkUnique(kit, true, attr.ObjectID, attr.PropertyID, attr.PropertyName, attr.BizID); err != nil {
		blog.Errorf("save attribute check unique input: %+v, err: %v, rid: %s", attr, err, kit.Rid)
		return err
	}

	if attr.PropertyType != common.FieldTypeIDRule {
		return nil
	}

	dbAttrs := make([]metadata.Attribute, 0)
	cond := mapstr.MapStr{common.BKObjIDField: attr.ObjectID}
	util.SetQueryOwner(cond, kit.SupplierAccount)
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &dbAttrs); err != nil {
		blog.Errorf("get %s attributes failed, err: %v, rid: %s", attr.ObjectID, err, kit.Rid)
		return err
	}

	attrTypeMap := make(map[string]string)
	for _, dbAttr := range dbAttrs {
		attrTypeMap[dbAttr.PropertyID] = dbAttr.PropertyType
	}

	err := attrvalid.ValidPropertyOption(kit, attr.PropertyType, attr.Option, attrTypeMap)
	if err != nil {
		blog.ErrorJSON("valid property option failed, err: %s, data: %s, rid: %s", err, attr, kit.Ctx)
		return err
	}

	if err = checkAddIDRule(kit, attr.ObjectID); err != nil {
		blog.ErrorJSON("check add asset id, err: %s, data: %s, rid: %s", err, attr, kit.Ctx)
		return err
	}

	return nil
}

// checkTableAttrUpdate delete the field that cannot be updated, check whether the field is repeated
func (m *modelAttribute) checkTableAttrUpdate(kit *rest.Kit, data mapstr.MapStr,
	cond universalsql.Condition) (err error) {

	dbAttributeArr, err := m.search(kit, cond)
	if err != nil {
		blog.Errorf("find nothing by the condition: %+v, err: %v, rid: %s", cond.ToMapStr(), err, kit.Rid)
		return err
	}
	if len(dbAttributeArr) == 0 {
		blog.Errorf("find nothing by the condition(%#v), rid: %s", cond.ToMapStr(), kit.Rid)
		return nil
	}

	// is there a predefined field for the updated attribute
	hasIsPreProperty := false
	for _, dbAttribute := range dbAttributeArr {
		if dbAttribute.IsPre {
			hasIsPreProperty = true
			break
		}
	}

	// 预定义字段，只能更新分组、分组内排序、单位、提示语和option
	if hasIsPreProperty {
		_ = data.ForEach(func(key string, val interface{}) error {
			if key != metadata.AttributeFieldPropertyGroup &&
				key != metadata.AttributeFieldPropertyIndex &&
				key != metadata.AttributeFieldUnit &&
				key != metadata.AttributeFieldPlaceHolder &&
				key != metadata.AttributeFieldOption {
				data.Remove(key)
			}
			return nil
		})
	}

	if err := checkAttrOption(kit, data, dbAttributeArr); err != nil {
		return err
	}

	if err = checkPropertyGroup(kit, data, dbAttributeArr); err != nil {
		return err
	}

	attr := metadata.Attribute{}
	if err = data.MarshalJSONInto(&attr); err != nil {
		blog.Errorf("marshal json into attribute failed, data: %+v, err: %v, rid: %s", data, err, kit.Rid)
		return err
	}

	if err = m.checkTableAttributeValidity(kit, attr); err != nil {
		blog.Errorf("check attribute validity failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	for _, dbAttr := range dbAttributeArr {
		if err = m.checkChangeField(kit, dbAttr, data); err != nil {
			return err
		}
	}

	// 删除不可更新字段， 避免由于传入数据，修改字段
	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldSupplierAccount)
	data.Remove(metadata.AttributeFieldPropertyType)
	data.Remove(metadata.AttributeFieldCreateTime)
	data.Remove(metadata.AttributeFieldIsPre)
	data.Remove(common.BKTemplateID)

	data.Set(metadata.AttributeFieldLastTime, time.Now())
	return err
}

func getObjectAttrTemplateID(kit *rest.Kit, attrID int64) (int64, error) {
	cond := mapstr.MapStr{
		common.BKFieldID: attrID,
	}
	cond = util.SetQueryOwner(cond, kit.SupplierAccount)
	attrs := make([]metadata.Attribute, 0)

	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).Fields(common.BKTemplateID).
		All(kit.Ctx, &attrs); err != nil {
		blog.Errorf("find attrs failed, attrID: %d, err: %v, rid: %s", attrID, err, kit.Rid)
		return 0, err
	}

	attrsNum := len(attrs)
	if attrsNum <= 0 || attrsNum > 1 {
		blog.Errorf("attributes num error, attID: %d, num: %d, rid: %s", attrID, attrsNum, kit.Rid)
		return 0, kit.CCError.Errorf(common.CCErrCommParamsInvalid, attrID)
	}

	return attrs[0].TemplateID, nil
}

func getTemplateAttrByID(kit *rest.Kit, templateID int64, fields []string) (*metadata.FieldTemplateAttr, error) {

	attrCond := mapstr.MapStr{
		common.BKFieldID: templateID,
	}
	attrCond = util.SetQueryOwner(attrCond, kit.SupplierAccount)

	templateAttr := make([]metadata.FieldTemplateAttr, 0)
	if err := mongodb.Client().Table(common.BKTableNameObjAttDesTemplate).Find(attrCond).Fields(fields...).
		All(kit.Ctx, &templateAttr); err != nil {
		blog.Errorf("find field template attr failed, cond: %v, err: %v, rid: %s", attrCond, err, kit.Rid)
		return nil, kit.CCError.CCError(common.CCErrCommDBSelectFailed)
	}

	templateAttrNum := len(templateAttr)
	if templateAttrNum > 1 {
		blog.Errorf("attributes num error, attID: %d, num: %d, rid: %s", templateID, templateAttrNum, kit.Rid)
		return nil, kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKTemplateID)
	}
	// here is the scenario of releasing the management
	if templateAttrNum == 0 {
		return nil, nil
	}

	return &templateAttr[0], nil
}

// checkAttrTemplateInfo the topo server has similar judgment logic. If it needs to be modified,
// both sides need to be modified at the same time. The function name in topo is: canAttrsUpdate
func checkAttrTemplateInfo(kit *rest.Kit, input mapstr.MapStr, attrID int64, isSync bool) error {
	// 1. 来自字段组合模版同步操作，都可以进行修改，直接正常返回
	if isSync {
		return nil
	}

	// 2. 不是同步操作，更新属性的bk_template_id为非0时，需要报错
	data := input.Clone()
	if newTmplID, ok := data[common.BKTemplateID]; ok {
		id, err := util.GetIntByInterface(newTmplID)
		if err != nil {
			blog.Errorf("get int by interface failed, val: %v, err: %s, rid: %s", newTmplID, err, kit.Rid)
			return err
		}

		if id != 0 {
			blog.Errorf("modify field %s forbidden, val: %s, rid: %s", common.BKTemplateID, newTmplID, kit.Rid)
			return kit.CCError.CCErrorf(common.CCErrCommModifyFieldForbidden, common.BKTemplateID)
		}
	}

	// 3. 不是同步操作，更新模型自己的属性，正常返回
	tmplID, err := getObjectAttrTemplateID(kit, attrID)
	if err != nil {
		return err
	}
	if tmplID == 0 {
		return nil
	}

	// 4. 验证来自模版的属性，是否可以正常更新
	return validTmplAttrCanUpdate(kit, data, tmplID)
}

func validTmplAttrCanUpdate(kit *rest.Kit, data mapstr.MapStr, tmplID int64) error {
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

	templateAttr, err := getTemplateAttrByID(kit, tmplID, fields)
	if err != nil {
		return err
	}
	if templateAttr == nil {
		return nil
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

func checkAttrOption(kit *rest.Kit, data mapstr.MapStr, dbAttributeArr []metadata.Attribute) error {
	option, exists := data.Get(metadata.AttributeFieldOption)
	if !exists {
		return nil
	}
	propertyType := dbAttributeArr[0].PropertyType
	for _, dbAttribute := range dbAttributeArr {
		if dbAttribute.PropertyType != propertyType {
			blog.Errorf("update option, but property type not the same, db attributes: %s, rid:%s",
				dbAttributeArr, kit.Ctx)
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, "cond")
		}
	}

	// 属性更新时，如果没有传入ismultiple参数，则使用数据库中的ismultiple值进行校验，如果传了ismultiple参数，则使用更新时的参数
	isMultiple := dbAttributeArr[0].IsMultiple
	if val, ok := data.Get(common.BKIsMultipleField); ok {
		ismultiple, ok := val.(bool)
		if !ok {
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKIsMultipleField)
		}
		isMultiple = &ismultiple
	}

	var extraOpt interface{}
	switch propertyType {
	case common.FieldTypeEnum, common.FieldTypeEnumMulti:
		extraOpt = isMultiple
	case common.FieldTypeIDRule:
		dbAttrs := make([]metadata.Attribute, 0)
		cond := mapstr.MapStr{common.BKObjIDField: dbAttributeArr[0].ObjectID}
		util.SetQueryOwner(cond, kit.SupplierAccount)
		if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).All(kit.Ctx, &dbAttrs); err != nil {
			blog.Errorf("get %s attributes failed, err: %v, rid: %s", dbAttributeArr[0].ObjectID, err, kit.Rid)
			return err
		}

		attrTypeMap := make(map[string]string)
		for _, dbAttr := range dbAttrs {
			attrTypeMap[dbAttr.PropertyID] = dbAttr.PropertyType
		}
		extraOpt = attrTypeMap
	default:
		extraOpt = data[common.BKDefaultFiled]
	}

	err := attrvalid.ValidPropertyOption(kit, propertyType, option, extraOpt)
	if err != nil {
		blog.ErrorJSON("valid property option failed, err: %s, data: %s, rid:%s", err, data, kit.Ctx)
		return err
	}

	return nil
}

func checkPropertyGroup(kit *rest.Kit, data mapstr.MapStr, dbAttributeArr []metadata.Attribute) error {

	grp, exists := data.Get(metadata.AttributeFieldPropertyGroup)
	if !exists {
		return nil
	}

	if grp == "" {
		data.Remove(metadata.AttributeFieldPropertyGroup)
	}

	// check if property group exists in object
	objIDs := make([]string, 0)
	for _, dbAttribute := range dbAttributeArr {
		objIDs = append(objIDs, dbAttribute.ObjectID)
	}
	objIDs = util.StrArrayUnique(objIDs)
	cond := map[string]interface{}{
		common.BKObjIDField: map[string]interface{}{
			common.BKDBIN: objIDs,
		},
	}
	if grp != "" {
		cond[common.BKPropertyGroupIDField] = grp
	}

	cnt, err := mongodb.Client().Table(common.BKTableNamePropertyGroup).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("property group count failed, err: %s, condition: %s, rid: %s", err, cond, kit.Rid)
		return err
	}
	if cnt != uint64(len(objIDs)) {
		blog.Errorf("property group invalid, objIDs: %s have %d property groups, rid: %s", objIDs, cnt, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsInvalid, metadata.AttributeFieldPropertyGroup)
	}
	return nil
}

// checkUpdate delete the field that cannot be updated, check whether the field is repeated
func (m *modelAttribute) checkUpdate(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition,
	isSync bool) (err error) {

	dbAttributeArr, err := m.search(kit, cond)
	if err != nil {
		blog.Errorf("request(%s): find nothing by the condition(%#v)  error(%s)", kit.Rid, cond.ToMapStr(), err)
		return err
	}
	if len(dbAttributeArr) == 0 {
		blog.Errorf("request(%s): find nothing by the condition(%#v)", kit.Rid, cond.ToMapStr())
		return nil
	}

	// is there a predefined field for the updated attribute。
	hasIsPreProperty := false
	for _, dbAttribute := range dbAttributeArr {
		if dbAttribute.IsPre {
			hasIsPreProperty = true
			break
		}
	}

	// 预定义字段，只能更新分组、分组内排序、单位、提示语和option
	if hasIsPreProperty {
		_ = data.ForEach(func(key string, val interface{}) error {
			if key != metadata.AttributeFieldPropertyGroup && key != metadata.AttributeFieldPropertyIndex &&
				key != metadata.AttributeFieldUnit && key != metadata.AttributeFieldPlaceHolder &&
				key != metadata.AttributeFieldOption {
				data.Remove(key)
			}
			return nil
		})
	}

	if err := checkAttrOption(kit, data, dbAttributeArr); err != nil {
		return err
	}

	if err = checkPropertyGroup(kit, data, dbAttributeArr); err != nil {
		return err
	}

	propertyType := dbAttributeArr[0].PropertyType
	// 对于枚举，枚举多选，枚举引用字段, 默认值是放在option中的，需要将default置为nil
	if data[metadata.AttributeFieldDefault] != nil && (propertyType == common.FieldTypeEnum ||
		propertyType == common.FieldTypeEnumMulti || propertyType == common.FieldTypeEnumQuote) {

		data.Remove(metadata.AttributeFieldDefault)
	}

	// 删除不可更新字段， 避免由于传入数据，修改字段
	// TODO: 改成白名单方式
	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldSupplierAccount)
	data.Remove(metadata.AttributeFieldPropertyType)
	data.Remove(metadata.AttributeFieldCreateTime)
	data.Remove(metadata.AttributeFieldIsPre)

	data.Set(metadata.AttributeFieldLastTime, time.Now())

	attribute := metadata.Attribute{}
	if err = data.MarshalJSONInto(&attribute); err != nil {
		blog.Errorf("marshal json into attribute failed, data: %+v, err: %v, rid: %s", data, err, kit.Rid)
		return err
	}

	// 更新default字段时，需要使用option对default进行数据校验，当没传时需要使用当前数据库里的数据进行校验
	if attribute.Default != nil && attribute.Option == nil {
		attribute.Option = dbAttributeArr[0].Option
	}

	if err = checkAttrTemplateInfo(kit, data, dbAttributeArr[0].ID, isSync); err != nil {
		blog.Errorf("check attribute template info failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	if err = m.checkAttributeValidity(kit, attribute, dbAttributeArr[0].PropertyType); err != nil {
		blog.Errorf("check attribute validity failed, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	for _, dbAttribute := range dbAttributeArr {
		err = m.checkUnique(kit, false, dbAttribute.ObjectID, dbAttribute.PropertyID, attribute.PropertyName,
			dbAttribute.BizID)
		if err != nil {
			blog.Errorf("save attribute check unique attribute: %+v, err: %v, rid:%s", attribute, err, kit.Rid)
			return err
		}
		if err = m.checkChangeField(kit, dbAttribute, data); err != nil {
			return err
		}
	}

	return err
}

// checkAttributeInUnique 检查属性是否存在唯一校验中  objIDPropertyIDArr  属性的bk_obj_id和表中ID的集合
func (m *modelAttribute) checkAttributeInUnique(kit *rest.Kit, objIDPropertyIDArr map[string][]int64) (bool, error) {

	cond := mongo.NewCondition()

	var orCondArr []universalsql.ConditionElement
	for objID, propertyIDArr := range objIDPropertyIDArr {
		orCondItem := mongo.NewCondition()
		orCondItem.Element(mongo.Field(common.BKObjIDField).Eq(objID))
		orCondItem.Element(mongo.Field("keys.key_id").In(propertyIDArr))
		orCondItem.Element(mongo.Field("keys.key_kind").Eq("property"))
		orCondArr = append(orCondArr, orCondItem)
	}

	cond.Or(orCondArr...)
	condMap := util.SetQueryOwner(cond.ToMapStr(), kit.SupplierAccount)

	cnt, err := mongodb.Client().Table(common.BKTableNameObjUnique).Find(condMap).Count(kit.Ctx)
	if err != nil {
		blog.ErrorJSON("checkAttributeInUnique db select error. err:%s, cond:%s, rid:%s", err.Error(), condMap, kit.Rid)
		return false, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if cnt > 0 {
		return true, nil
	}

	return false, nil
}

func (m *modelAttribute) delIDRuleUnique(kit *rest.Kit, objIDPropertyIDArr map[string][]int64) error {
	cond := mongo.NewCondition()

	var orCondArr []universalsql.ConditionElement
	for objID, propertyIDArr := range objIDPropertyIDArr {
		orCondItem := mongo.NewCondition()
		orCondItem.Element(mongo.Field(common.BKObjIDField).Eq(objID))
		orCondItem.Element(mongo.Field("keys.key_id").In(propertyIDArr))
		orCondItem.Element(mongo.Field("keys.key_kind").Eq("property"))
		orCondArr = append(orCondArr, orCondItem)
	}

	cond.Or(orCondArr...)
	condMap := util.SetModOwner(cond.ToMapStr(), kit.SupplierAccount)

	if _, err := mongodb.Client().Table(common.BKTableNameObjUnique).DeleteMany(kit.Ctx, condMap); err != nil {
		blog.ErrorJSON("delete unique failed, cond: %+v, err: %v, rid: %s", condMap, err, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBDeleteFailed)
	}

	return nil
}

// checkAddField 新加模型属性的时候，如果新加的是必填字段，需要判断是否可以新加必填字段
func (m *modelAttribute) checkAddField(kit *rest.Kit, attribute metadata.Attribute) error {
	langObjID := m.getLangObjID(kit, attribute.ObjectID)
	if _, ok := notAddAttrModel[attribute.ObjectID]; ok {
		//  不允许新加字段的模型
		return kit.CCError.Errorf(common.CCErrCoreServiceNotAllowAddFieldErr, langObjID)
	}

	if _, ok := RequiredFieldUnchangeableModels[attribute.ObjectID]; ok {
		if attribute.IsRequired {
			//  不允许修改必填字段的模型
			return kit.CCError.Errorf(common.CCErrCoreServiceNotAllowAddRequiredFieldErr, langObjID)
		}

	}
	return nil
}

// checkChangeField 修改模型属性的时候，如果修改的属性包含是否为必填字段(isrequired)，需要判断该模型的必填字段是否允许被修改
func (m *modelAttribute) checkChangeField(kit *rest.Kit, attr metadata.Attribute, attrInfo mapstr.MapStr) error {
	langObjID := m.getLangObjID(kit, attr.ObjectID)
	if _, ok := RequiredFieldUnchangeableModels[attr.ObjectID]; ok {
		if attrInfo.Exists(metadata.AttributeFieldIsRequired) {
			// 不允许修改模型的必填字段
			val, ok := attrInfo[metadata.AttributeFieldIsRequired].(bool)
			if !ok {
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldIsRequired)
			}
			if val != attr.IsRequired {
				return kit.CCError.Errorf(common.CCErrCoreServiceNotAllowChangeRequiredFieldErr, langObjID)
			}
		}
	}
	return nil
}

func (m *modelAttribute) getLangObjID(kit *rest.Kit, objID string) string {
	langKey := "object_" + objID
	language := httpheader.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	langObjID := lang.Language(langKey)
	if langObjID == langKey {
		langObjID = objID
	}
	return langObjID
}

// GetAttrLastIndex TODO
func (m *modelAttribute) GetAttrLastIndex(kit *rest.Kit, attribute metadata.Attribute) (int64, error) {
	opt := make(map[string]interface{})
	opt[common.BKObjIDField] = attribute.ObjectID
	opt[common.BKPropertyGroupField] = attribute.PropertyGroup
	opt = util.SetModOwner(opt, attribute.OwnerID)
	count, err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(opt).Count(kit.Ctx)
	if err != nil {
		blog.Error("GetAttrLastIndex, request(%s): database operation is failed, error info is %v", kit.Rid, err)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}
	if count <= 0 {
		return 0, nil
	}

	attrs := make([]metadata.Attribute, 0)
	sortCond := "-bk_property_index"
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(opt).Sort(sortCond).Limit(1).All(kit.Ctx,
		&attrs); err != nil {
		blog.Error("GetAttrLastIndex, database operation is failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if len(attrs) <= 0 {
		return 0, nil
	}
	return attrs[0].PropertyIndex + 1, nil
}

func checkAddIDRule(kit *rest.Kit, objID string) error {
	cond := mapstr.MapStr{common.BKObjIDField: objID, common.BKPropertyTypeField: common.FieldTypeIDRule}

	count, err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond).Count(kit.Ctx)
	if err != nil {
		blog.Errorf("count attribute failed. err: %v, cond: %s, rid: %s", err, cond, kit.Rid)
		return kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if count >= metadata.IDRuleFieldLimit {
		return kit.CCError.Errorf(common.CCErrCommXXExceedLimit, common.FieldTypeIDRule, metadata.IDRuleFieldLimit)
	}

	return nil
}
