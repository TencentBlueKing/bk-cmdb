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

package attrvalid

import (
	"fmt"
	"regexp"
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
	"configcenter/src/common/valid/attribute/manager"
)

// ValidPropertyOption valid property field option
func ValidPropertyOption(kit *rest.Kit, propertyType string, option interface{}, extraOpt interface{}) error {
	switch propertyType {
	case common.FieldTypeEnum, common.FieldTypeEnumMulti:
		isMultiple, ok := extraOpt.(*bool)
		if !ok {
			blog.Errorf("extra opt(%+v) type %T is invalid, rid: %s", extraOpt, extraOpt, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
		}
		return ValidFieldTypeEnumOption(kit, option, isMultiple)
	case common.FieldTypeInt:
		return ValidFieldTypeInt(kit, option, extraOpt)
	case common.FieldTypeFloat:
		return ValidFieldTypeFloat(kit, option, extraOpt)
	case common.FieldTypeList:
		return ValidFieldTypeList(kit, option, extraOpt)
	case common.FieldTypeLongChar, common.FieldTypeSingleChar:
		return ValidFieldTypeString(kit, option, extraOpt)
	case common.FieldTypeBool:
		return valid.ValidateBoolType(option)
	case common.FieldTypeIDRule:
		attrTypeMap, ok := extraOpt.(map[string]string)
		if !ok {
			blog.Errorf("extra opt(%+v) type %T is invalid, rid: %s", extraOpt, extraOpt, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
		}
		return ValidIDRuleOption(kit, option, attrTypeMap)
	}

	if handle, ok := manager.Get(propertyType); ok {
		if err := handle.ValidateOption(kit.Ctx, option, extraOpt); err != nil {
			blog.Errorf("valid property option failed, property type: %s, option: %+v, extra opt: %+v, err: %v, rid: %s",
				propertyType, option, extraOpt, err, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, err.Error())
		}
	}

	return nil
}

// ValidFieldTypeEnumOption validate enum field type's option
func ValidFieldTypeEnumOption(kit *rest.Kit, option interface{}, isMultiple *bool) error {
	if option == nil {
		return kit.CCError.Errorf(common.CCErrCommParamsLostField, "option")
	}

	if isMultiple == nil {
		return kit.CCError.Errorf(common.CCErrCommParamsInvalid, common.BKIsMultipleField)
	}

	enumOption, err := metadata.ParseEnumOption(option)
	if err != nil {
		blog.Errorf("parse enum option %+v failed, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	if len(enumOption) > common.AttributeOptionArrayMaxLength {
		blog.Errorf("enum option array length %d exceeds max length %d, rid: %s", len(enumOption),
			common.AttributeOptionArrayMaxLength, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionArrayMaxLength)
	}

	var count int
	for _, o := range enumOption {
		id := o.ID
		if o.ID == "" {
			blog.Errorf("enum option id can't be empty, option: %+v, rid: %s", option, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "option id")
		}

		if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(id) {
			blog.Errorf("option id %s length %d exceeds max length %d, rid: %s", id, utf8.RuneCountInString(id),
				common.AttributeOptionValueMaxLength, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option id",
				common.AttributeOptionValueMaxLength)
		}

		name := o.Name
		if name == "" {
			blog.Errorf("enum option name can't be empty, option: %+v, rid: %s", option, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsNeedSet, "option name")
		}

		if o.IsDefault {
			count += 1
		}

		switch o.Type {
		case "text":
			if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(name) {
				blog.Errorf("option name %s length %d exceeds max length %d, rid: %s", name,
					utf8.RuneCountInString(name), common.AttributeOptionValueMaxLength, kit.Rid)
				return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option name",
					common.AttributeOptionValueMaxLength)
			}
		default:
			blog.Errorf("enum option type must be 'text', current: %v, rid: %s", o.Type, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option type")
		}
	}

	if !*isMultiple && count != 1 {
		blog.Errorf("field type is single choice, but default value is multiple, count: %d, rid: %s", count,
			kit.Rid)
		return kit.CCError.CCError(common.CCErrCommParamsNeedSingleChoice)
	}

	return nil
}

// ValidFieldTypeInt validate int or float field type's option and default value
func ValidFieldTypeInt(kit *rest.Kit, option, defaultVal interface{}) error {
	if option == nil {
		return kit.CCError.Errorf(common.CCErrCommParamsLostField, "option")
	}

	// validate maximum & minimum option
	intOption, err := metadata.ParseIntOption(option)
	if err != nil {
		blog.Errorf("parse int option %+v failed, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	maxVal := intOption.Max
	minVal := intOption.Min

	if minVal > maxVal {
		blog.Errorf("option min value %d is greater than max value %d, rid: %s", minVal, maxVal, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	// validate default value
	if defaultVal == nil {
		return nil
	}
	defaultValue, err := util.GetInt64ByInterface(defaultVal)
	if err != nil {
		blog.Errorf("parse int field default value %+v failed, err: %v, rid: %s", defaultVal, err, kit.Rid)
		return err
	}

	if defaultValue < minVal || defaultValue > maxVal {
		return fmt.Errorf("int type field default value over limit")
	}

	return nil
}

// ValidFieldTypeFloat validate int or float field type's option default value
func ValidFieldTypeFloat(kit *rest.Kit, option, defaultVal interface{}) error {
	if option == nil {
		return kit.CCError.Errorf(common.CCErrCommParamsLostField, "option")
	}

	// validate maximum & minimum option
	floatOption, err := metadata.ParseFloatOption(option)
	if err != nil {
		blog.Errorf("parse float option %+v failed, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	maxVal := floatOption.Max
	minVal := floatOption.Min

	if minVal > maxVal {
		blog.Errorf("option min value %d is greater than max value %d, rid: %s", minVal, maxVal, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	// validate default value
	if defaultVal == nil {
		return nil
	}

	defaultValue, err := util.GetFloat64ByInterface(defaultVal)
	if err != nil {
		blog.Errorf("parse float field default value %+v failed, err: %v, rid: %s", defaultVal, err, kit.Rid)
		return err
	}

	if defaultValue < minVal || defaultValue > maxVal {
		return fmt.Errorf("float type field default value over limit")
	}
	return nil
}

// ValidFieldTypeList validate list field type's option and default value
func ValidFieldTypeList(kit *rest.Kit, option, defaultVal interface{}) error {
	if option == nil {
		return kit.CCError.Errorf(common.CCErrCommParamsLostField, "option")
	}

	listOption, err := metadata.ParseListOption(option)
	if err != nil {
		blog.Errorf("parse list option %+v failed, err: %v, rid: %s", option, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	if len(listOption) > common.AttributeOptionArrayMaxLength {
		blog.Errorf("option array length %d exceeds maximum %d", len(listOption), common.AttributeOptionArrayMaxLength)
		return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionArrayMaxLength)
	}

	isDefaultValid := false
	if defaultVal == nil {
		isDefaultValid = true
	}
	listDefaultVal := util.GetStrByInterface(defaultVal)

	for _, value := range listOption {
		if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(value) {
			blog.Errorf("option value %s length %d exceeds max length %d, rid: %s", value,
				utf8.RuneCountInString(value), common.AttributeOptionValueMaxLength, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommValExceedMaxFailed, "option",
				common.AttributeOptionValueMaxLength)
		}

		if listDefaultVal == value {
			isDefaultValid = true
		}
	}

	if !isDefaultValid {
		blog.Errorf("default list value %+v is invalid, option: %+v, rid: %s", defaultVal, listOption, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "list default value")
	}

	return nil
}

// ValidFieldTypeString validate string field type's regex option and default value
func ValidFieldTypeString(kit *rest.Kit, option, defaultVal interface{}) error {
	if option == nil {
		return nil
	}

	// 校验正则是否合法
	regular, ok := option.(string)
	if !ok {
		blog.Errorf("string type option %+v type %T is invalid, rid: %s", option, option, kit.Rid)
		return kit.CCError.Errorf(common.CCIllegalRegularExpression, "option")
	}

	strDefVal := ""
	if defaultVal != nil {
		strDefVal, ok = defaultVal.(string)
		if !ok {
			blog.Errorf("string type default value %+v type %T is invalid, rid: %s", defaultVal, defaultVal, kit.Rid)
			return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "string default value")
		}
	}

	if len(regular) == 0 {
		return nil
	}

	if _, err := regexp.Compile(regular); err != nil {
		blog.Errorf("regular expression %s is invalid, err: %v, rid: %s", regular, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrorCheckRegularFailed, "regular is wrong")
	}

	if defaultVal == nil {
		return nil
	}

	match, err := regexp.MatchString(regular, strDefVal)
	if err != nil || !match {
		blog.Errorf("default value %s not matches string option %s, err: %v, rid: %s", strDefVal, regular, err, kit.Rid)
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "string default value")
	}

	return nil
}

var validTableFieldType = map[string]struct{}{
	common.FieldTypeInt:        {},
	common.FieldTypeEnumMulti:  {},
	common.FieldTypeLongChar:   {},
	common.FieldTypeSingleChar: {},
	common.FieldTypeFloat:      {},
	common.FieldTypeBool:       {},
}

// ValidTableFieldOption judging the legitimacy of the basic type of the form field
func ValidTableFieldOption(kit *rest.Kit, propertyType string, option, defaultValue interface{},
	isMultiple *bool, objID string) error {

	_, exists := validTableFieldType[propertyType]
	if !exists {
		return fmt.Errorf("type(%s) is not among the underlying types supported by the table field", propertyType)
	}

	switch propertyType {
	case common.FieldTypeEnumMulti:
		return ValidFieldTypeEnumOption(kit, option, isMultiple)
	default:
		return ValidPropertyOption(kit, propertyType, option, defaultValue)
	}
}

// ValidPropertyTypeIsMultiple valid object attr field type is multiple
func ValidPropertyTypeIsMultiple(kit *rest.Kit, propertyType string, isMultiple *bool) error {
	switch propertyType {
	case common.FieldTypeSingleChar, common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeEnum,
		common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeLongChar, common.FieldTypeTimeZone,
		common.FieldTypeBool, common.FieldTypeList:
		if isMultiple != nil && *isMultiple {
			return kit.CCError.Errorf(common.CCErrCommFieldTypeNotSupportMultiple, propertyType)
		}
	}
	return nil
}

// IsStrProperty  is string property
func IsStrProperty(propertyType string) bool {
	if common.FieldTypeLongChar == propertyType || common.FieldTypeSingleChar == propertyType {
		return true
	}

	return false
}

// ValidIDRuleOption validate id rule field type's option
func ValidIDRuleOption(kit *rest.Kit, val interface{}, attrTypeMap map[string]string) error {
	rules, err := metadata.ParseSubIDRules(val)
	if err != nil {
		return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, err.Error())
	}

	for _, rule := range rules {
		switch rule.Kind {
		case metadata.Attr:
			attrType, exists := attrTypeMap[rule.Val]
			if !exists {
				blog.Errorf("attr val %s is invalid, attribute not exists, rid: %s", rule.Val, kit.Rid)
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
			}
			if !metadata.IsValidAttrRuleType(attrType) {
				blog.Errorf("attr val %s type %s is invalid, rid: %s", rule.Val, attrType, kit.Rid)
				return kit.CCError.Errorf(common.CCErrCommParamsIsInvalid, "option")
			}
		}
	}

	return nil
}
