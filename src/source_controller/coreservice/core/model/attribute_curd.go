/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
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
	"configcenter/src/common/http/rest"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/universalsql"
	"configcenter/src/common/universalsql/mongo"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/types"
	"configcenter/src/storage/driver/mongodb"
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

	if nil == attribute.CreateTime {
		attribute.CreateTime = &metadata.Time{}
		attribute.CreateTime.Time = time.Now()
	}

	if nil == attribute.LastTime {
		attribute.LastTime = &metadata.Time{}
		attribute.LastTime.Time = time.Now()
	}

	if err = m.saveCheck(kit, attribute); err != nil {
		return 0, err
	}

	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Insert(kit.Ctx, attribute)
	return id, err
}

func (m *modelAttribute) checkUnique(kit *rest.Kit, isCreate bool, objID, propertyID, propertyName string, modelBizID int64) error {
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

	language := util.GetLanguage(kit.Header)
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

func (m *modelAttribute) checkAttributeValidity(kit *rest.Kit, attribute metadata.Attribute, isUpdate bool) error {
	language := util.GetLanguage(kit.Header)
	lang := m.language.CreateDefaultCCLanguageIf(language)
	if attribute.PropertyID != "" {
		attribute.PropertyID = strings.TrimSpace(attribute.PropertyID)
		if common.AttributeIDMaxLength < utf8.RuneCountInString(attribute.PropertyID) {
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, lang.Language("model_attr_bk_property_id"),
				common.AttributeIDMaxLength)
		}

		if !SatisfyMongoFieldLimit(attribute.PropertyID) {
			blog.Errorf("attribute.PropertyID:%s not SatisfyMongoFieldLimit", attribute.PropertyID)
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
		switch attribute.PropertyType {
		case common.FieldTypeSingleChar, common.FieldTypeLongChar, common.FieldTypeInt, common.FieldTypeFloat,
			common.FieldTypeEnum, common.FieldTypeEnumMulti, common.FieldTypeDate, common.FieldTypeTime,
			common.FieldTypeUser, common.FieldTypeOrganization, common.FieldTypeTimeZone, common.FieldTypeBool,
			common.FieldTypeList, common.FieldTypeEnumQuote:
		default:
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
		}
	}

	if attribute.Default != nil {
		// 更新场景时，如果需要置空默认值，此时传递的参数为统一处理为default:""
		if isUpdate {
			if val, ok := attribute.Default.(string); ok {
				if val == "" {
					return nil
				}
			}
		}
		if err := m.checkAttributeDefaultValue(kit, attribute); err != nil {
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

func (m *modelAttribute) checkAttributeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error {
	switch attribute.PropertyType {
	case common.FieldTypeEnum, common.FieldTypeEnumMulti, common.FieldTypeEnumQuote:
	case common.FieldTypeSingleChar, common.FieldTypeLongChar:
		if err := m.checkStringTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeInt:
		if err := m.checkIntTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeFloat:
		if err := m.checkFloatTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeDate:
		if err := m.checkDateTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeTime:
		if err := m.checkTimeTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeUser:
		if err := m.checkUserTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeOrganization:
		if err := m.checkOrganizationTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeTimeZone:
		if err := m.checkTimeZoneTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeBool:
		if err := m.checkBoolTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	case common.FieldTypeList:
		if err := m.checkListTypeDefaultValue(kit, attribute); err != nil {
			return err
		}
	default:
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldPropertyType)
	}

	return nil
}

func (m *modelAttribute) checkStringTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error{
	if err := util.ValidateStringType(attribute.Default); err != nil {
		blog.Errorf("single char or long char default value not string, err: %v, rid: %s", err, kit.Rid)
		return err
	}

	defaultStr, ok := attribute.Default.(string)
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "attribute.default")
	}

	if attribute.Option == nil {
		return nil
	}
	optStr, ok := attribute.Option.(string)
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "attribute.Option")
	}

	match, err := regexp.MatchString(optStr, defaultStr)
	if err != nil || !match {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	return nil
}

func (m *modelAttribute) checkIntTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error {
	if ok := util.IsNumeric(attribute.Default); !ok {
		blog.Errorf("int type default value not numeric, type: %T, rid: %s", attribute.Default, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}
	if attribute.Option == nil {
		return nil
	}

	defaultVal, err := util.GetInt64ByInterface(attribute.Default)
	if err != nil {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedInt, "attribute.default")
	}

	tmp, ok := attribute.Option.(map[string]interface{})
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	min, ok := tmp["min"]
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "min")
	}
	if ok := util.IsNumeric(min); ok {
		minVal, err := util.GetInt64ByInterface(min)
		if err != nil {
			return kit.CCError.Errorf(common.CCErrCommParamsNeedInt, "attribute.Option.min")
		}
		if defaultVal < minVal {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "int default value")
		}
	}

	max, ok := tmp["max"]
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "max")
	}
	if ok := util.IsNumeric(max); ok {
		maxVal, err := util.GetInt64ByInterface(max)
		if err != nil {
			return kit.CCError.Errorf(common.CCErrCommParamsNeedInt, "attribute.Option.max")
		}
		if defaultVal > maxVal {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "int default value")
		}
	}
	return nil
}

func (m *modelAttribute) checkFloatTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error {
	if ok := util.IsNumeric(attribute.Default); !ok {
		blog.Errorf("float type default value not numeric, type: %T, rid: %s", attribute.Default, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}
	if attribute.Option == nil {
		return nil
	}

	defaultVal, err := util.GetFloat64ByInterface(attribute.Default)
	if err != nil {
		return kit.CCError.Errorf(common.CCErrCommParamsNeedFloat, "attribute.default")
	}

	tmp, ok := attribute.Option.(map[string]interface{})
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}
	min, ok := tmp["min"]
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "min")
	}
	if ok := util.IsNumeric(min); ok {
		minVal, err := util.GetFloat64ByInterface(min)
		if err != nil {
			return kit.CCError.Errorf(common.CCErrCommParamsNeedFloat, "attribute.Option.min")
		}
		if defaultVal < minVal {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "float default value")
		}
	}

	max, ok := tmp["max"]
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "max")
	}
	if ok := util.IsNumeric(max); ok {
		maxVal, err := util.GetFloat64ByInterface(max)
		if err != nil {
			return kit.CCError.Errorf(common.CCErrCommParamsNeedFloat, "attribute.Option.max")
		}
		if defaultVal > maxVal {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "float default value")
		}
	}
	return nil
}

func (m *modelAttribute) checkListTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error {
	arrOption, ok := attribute.Option.([]interface{})
	if !ok || len(arrOption) == 0 {
		blog.Errorf("option %v not string type list option", attribute.Option)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "attribute.Option")
	}

	defaultVal := util.GetStrByInterface(attribute.Default)
	for _, value := range arrOption {
		val := util.GetStrByInterface(value)
		if defaultVal == val {
			return nil
		}
	}

	return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "list default value")
}

func (m *modelAttribute) checkUserTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error {
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
	}
	return nil
}

func (m *modelAttribute) checkOrganizationTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error {
	switch org := attribute.Default.(type) {
	case []interface{}:
		if len(org) == 0 {
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
		}

		for _, orgID := range org {
			if !util.IsNumeric(orgID) {
				blog.Errorf("orgID params not int, type: %T, rid: %s", orgID, kit.Rid)
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
			}
		}
	}
	return nil
}

func (m *modelAttribute) checkDateTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error{
	dateStr, ok := attribute.Default.(string)
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	if ok := util.IsDate(dateStr); !ok {
		blog.Errorf("date type default value not date format, type: %T, rid: %s", attribute.Default, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	return nil
}

func (m *modelAttribute) checkTimeTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error{
	timeStr, ok := attribute.Default.(string)
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	if _, ok := util.IsTime(timeStr); !ok {
		blog.Errorf("time type default value not time format, rid: %s", kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	return nil
}

func (m *modelAttribute) checkTimeZoneTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error{
	timeZoneStr, ok := attribute.Default.(string)
	if !ok {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	if ok := util.IsTimeZone(timeZoneStr); !ok {
		blog.Errorf("time zone type default value not time zone format, rid: %s", kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, metadata.AttributeFieldDefault)
	}

	return nil
}

func (m *modelAttribute) checkBoolTypeDefaultValue (kit *rest.Kit, attribute metadata.Attribute) error{
	if err := util.ValidateBoolType(attribute.Default); err != nil {
		blog.Errorf("bool type default value not bool, err: %v, rid: %s", err, kit.Rid)
		return err
	}
	return nil
}

func (m *modelAttribute) update(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (cnt uint64, err error) {
	err = m.checkUpdate(kit, data, cond)
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

func (m *modelAttribute) search(kit *rest.Kit, cond universalsql.Condition) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) searchWithSort(kit *rest.Kit, cond metadata.QueryCondition) (resultAttrs []metadata.Attribute, err error) {
	resultAttrs = []metadata.Attribute{}

	instHandler := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.Condition)
	err = instHandler.Start(uint64(cond.Page.Start)).Limit(uint64(cond.Page.Limit)).Sort(cond.Page.Sort).All(kit.Ctx, &resultAttrs)

	return resultAttrs, err
}

func (m *modelAttribute) searchReturnMapStr(kit *rest.Kit, cond universalsql.Condition) (resultAttrs []mapstr.MapStr, err error) {

	resultAttrs = []mapstr.MapStr{}
	err = mongodb.Client().Table(common.BKTableNameObjAttDes).Find(cond.ToMapStr()).All(kit.Ctx, &resultAttrs)
	return resultAttrs, err
}

func (m *modelAttribute) delete(kit *rest.Kit, cond universalsql.Condition) (cnt uint64, err error) {

	resultAttrs := make([]metadata.Attribute, 0)
	fields := []string{common.BKFieldID, common.BKPropertyIDField, common.BKObjIDField, common.BKAppIDField}

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
	for _, attr := range resultAttrs {
		objIDArrMap[attr.ObjectID] = append(objIDArrMap[attr.ObjectID], attr.ID)
	}

	if err := m.cleanAttributeFieldInInstances(kit.Ctx, kit.SupplierAccount, resultAttrs); err != nil {
		blog.Errorf("delete object attributes with cond: %v, but delete these attribute in instance failed, "+
			"err: %v, rid: %s", condMap, err, kit.Rid)
		return 0, err
	}

	// delete template attribute when delete model attribute
	if err := m.cleanAttrTemplateRelation(kit.Ctx, kit.SupplierAccount, resultAttrs); err != nil {
		blog.Errorf("delete the relation between attributes and templates failed, attr: %v, err: %v, rid: %s",
			resultAttrs, err, kit.Rid)
		return 0, err
	}

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

// cleanAttributeFieldInInstances TODO
// remove attribute filed in this object's instances
func (m *modelAttribute) cleanAttributeFieldInInstances(ctx context.Context, ownerID string, attrs []metadata.Attribute) error {
	// this operation may take a long time, do not use transaction
	ctx = context.Background()

	objectFields := make(map[string][]bizObjectFields, 0)
	hostApplyFields := make(map[int64][]int64)

	// TODO: now, we only support set, module, host model's biz attribute clean operation.
	for _, attr := range attrs {
		biz := attr.BizID
		if biz != 0 {
			if !isBizObject(attr.ObjectID) {
				return fmt.Errorf("unsupported object %s's clean instance field operation", attr.ObjectID)
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

			cond := map[string]interface{}{
				common.BKDBOR: existConds,
			}

			if objField.bizID > 0 {
				if !isBizObject(object) {
					return fmt.Errorf("unsupported object %s's clean instance field operation", object)
				}

				if object == common.BKInnerObjIDHost {
					if err := m.cleanHostAttributeField(ctx, ownerID, objField); err != nil {
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
						if err := m.cleanHostAttributeField(ctx, ownerID, ele); err != nil {
							return err
						}
						continue
					}
				} else {
					cond[common.BKObjIDField] = object
				}
			}

			cond = util.SetQueryOwner(cond, ownerID)

			collectionName := common.GetInstTableName(object, ownerID)
			wg.Add(1)
			go func(collName string, filter types.Filter, fields []string) {
				defer wg.Done()

				instCount, err := mongodb.Client().Table(collName).Find(filter).Count(ctx)
				if err != nil {
					blog.Error("count instances with the attribute to delete failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, filter, fields, err)
					hitError = err
					return
				}

				instIDField := common.GetInstIDField(object)
				for start := uint64(0); start < instCount; start += pageSize {
					insts := make([]map[string]interface{}, 0)
					err := mongodb.Client().Table(collName).Find(filter).Start(0).Limit(pageSize).Fields(instIDField).All(ctx, &insts)
					if err != nil {
						blog.Error("get instance ids with the attribute to delete failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, filter, fields, err)
						hitError = err
						return
					}

					if len(insts) == 0 {
						return
					}

					instIDs := make([]int64, len(insts))
					for index, inst := range insts {
						instID, err := util.GetInt64ByInterface(inst[instIDField])
						if err != nil {
							blog.Error("get instance id failed, inst: %+v, err: %v", inst, err)
							hitError = err
							return
						}
						instIDs[index] = instID
					}

					instFilter := map[string]interface{}{
						instIDField: map[string]interface{}{
							common.BKDBIN: instIDs,
						},
					}

					if err := mongodb.Client().Table(collName).DropColumns(ctx, instFilter, fields); err != nil {
						blog.Error("delete object's attribute from instance failed, table: %s, cond: %v, fields: %v, err: %v", collectionName, instFilter, fields, err)
						hitError = err
						return
					}
				}
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
	if err := m.cleanHostApplyField(ctx, ownerID, hostApplyFields); err != nil {
		return err
	}

	return nil
}

func (m *modelAttribute) cleanAttrTemplateRelation(ctx context.Context, ownerID string,
	attrs []metadata.Attribute) error {

	if len(attrs) == 0 {
		return nil
	}

	attrMap := make(map[string][]int64)
	for _, attr := range attrs {
		attrMap[attr.ObjectID] = append(attrMap[attr.ObjectID], attr.ID)
	}

	for objID, attrIDs := range attrMap {
		cond := mapstr.MapStr{
			common.BKAttributeIDField: mapstr.MapStr{common.BKDBIN: attrIDs},
		}
		cond = util.SetQueryOwner(cond, ownerID)
		switch objID {
		case common.BKInnerObjIDSet:
			if err := mongodb.Client().Table(common.BKTableNameSetTemplateAttr).Delete(ctx, cond); err != nil {
				return err
			}

		case common.BKInnerObjIDModule:
			if err := mongodb.Client().Table(common.BKTableNameServiceTemplateAttr).Delete(ctx, cond); err != nil {
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
		err := mongodb.Client().Table(common.BKTableNameModuleHostConfig).Find(cond).Start(start).Limit(pageSize).Fields(common.BKHostIDField).All(ctx, &hostList)
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
		if err := mongodb.Client().Table(common.BKTableNameBaseHost).DropColumns(ctx, hostFilter, info.fields); err != nil {
			return fmt.Errorf("clean host biz attribute %v failed, err: %v", info.fields, err)
		}
	}

	return nil

}

func (m *modelAttribute) cleanHostApplyField(ctx context.Context, ownerID string, hostApplyFields map[int64][]int64) error {
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

// saveCheck TODO
//  saveCheck 新加字段检查
func (m *modelAttribute) saveCheck(kit *rest.Kit, attribute metadata.Attribute) error {

	if err := m.checkAddField(kit, attribute); err != nil {
		return err
	}

	if err := m.checkAttributeMustNotEmpty(kit, attribute); err != nil {
		return err
	}
	if err := m.checkAttributeValidity(kit, attribute, false); err != nil {
		return err
	}

	// check name duplicate
	if err := m.checkUnique(kit, true, attribute.ObjectID, attribute.PropertyID, attribute.PropertyName, attribute.BizID); err != nil {
		blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, kit.Rid)
		return err
	}

	return nil
}

// checkUpdate 删除不可以更新字段，检验字段是否重复， 返回更新的行数，错误
func (m *modelAttribute) checkUpdate(kit *rest.Kit, data mapstr.MapStr, cond universalsql.Condition) (err error) {
	attribute := metadata.Attribute{}
	if err = data.MarshalJSONInto(&attribute); err != nil {
		blog.Errorf("request(%s): MarshalJSONInto(%#v), error is %v", kit.Rid, data, err)
		return err
	}
	if err = m.checkAttributeValidity(kit, attribute, true); err != nil {
		return err
	}

	dbAttributeArr, err := m.search(kit, cond)
	if err != nil {
		blog.Errorf("request(%s): find nothing by the condition(%#v)  error(%s)", kit.Rid, cond.ToMapStr(), err)
		return err
	}
	if len(dbAttributeArr) == 0 {
		blog.Errorf("request(%s): find nothing by the condition(%#v)", kit.Rid, cond.ToMapStr())
		return nil
	}

	// 更新的属性是否存在预定义字段。
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

	if option, exists := data.Get(metadata.AttributeFieldOption); exists {
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

		if isMultiple == nil {
			return kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKIsMultipleField)
		}

		if err := util.ValidPropertyOption(propertyType, option, *isMultiple, kit.CCError); err != nil {
			blog.ErrorJSON("valid property option failed, err: %s, data: %s, rid:%s", err, data, kit.Ctx)
			return err
		}
	}

	// 删除不可更新字段， 避免由于传入数据，修改字段
	// TODO: 改成白名单方式
	data.Remove(metadata.AttributeFieldPropertyID)
	data.Remove(metadata.AttributeFieldSupplierAccount)
	data.Remove(metadata.AttributeFieldPropertyType)
	data.Remove(metadata.AttributeFieldCreateTime)
	data.Remove(metadata.AttributeFieldIsPre)
	data.Set(metadata.AttributeFieldLastTime, time.Now())

	if grp, exists := data.Get(metadata.AttributeFieldPropertyGroup); exists {
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
			common.BKPropertyGroupIDField: grp,
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
	}

	for _, dbAttribute := range dbAttributeArr {
		err = m.checkUnique(kit, false, dbAttribute.ObjectID, dbAttribute.PropertyID, attribute.PropertyName, attribute.BizID)
		if err != nil {
			blog.ErrorJSON("save attribute check unique err:%s, input:%s, rid:%s", err.Error(), attribute, kit.Rid)
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
	language := util.GetLanguage(kit.Header)
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
	if err := mongodb.Client().Table(common.BKTableNameObjAttDes).Find(opt).Sort(sortCond).Limit(1).All(kit.Ctx, &attrs); err != nil {
		blog.Error("GetAttrLastIndex, database operation is failed, err: %v, rid: %s", err, kit.Rid)
		return 0, kit.CCError.Error(common.CCErrCommDBSelectFailed)
	}

	if len(attrs) <= 0 {
		return 0, nil
	}
	return attrs[0].PropertyIndex + 1, nil
}
