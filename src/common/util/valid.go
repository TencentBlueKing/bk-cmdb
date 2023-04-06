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

package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
)

// ValidPropertyOption valid property field option
func ValidPropertyOption(propertyType string, option interface{}, isMultiple bool, defaultVal interface{},
	errProxy ccErr.DefaultCCErrorIf) error {
	switch propertyType {
	case common.FieldTypeEnum, common.FieldTypeEnumMulti:
		return ValidFieldTypeEnumOption(option, isMultiple, errProxy)
	case common.FieldTypeInt:
		return ValidFieldTypeInt(option, defaultVal, errProxy)
	case common.FieldTypeFloat:
		return ValidFieldTypeFloat(option, defaultVal, errProxy)
	case common.FieldTypeList:
		return ValidFieldTypeList(option, defaultVal, errProxy)
	case common.FieldTypeLongChar, common.FieldTypeSingleChar:
		return ValidFieldTypeString(option, defaultVal, errProxy)
	}
	return nil
}

// ValidFieldTypeEnumOption validate enum field type's option
func ValidFieldTypeEnumOption(option interface{}, isMultiple bool, errProxy ccErr.DefaultCCErrorIf) error {
	if option == nil {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	arrOption, ok := option.([]interface{})
	if !ok {
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	if len(arrOption) > common.AttributeOptionArrayMaxLength {
		return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionArrayMaxLength)
	}

	var count int
	for _, o := range arrOption {
		mapOption, ok := o.(map[string]interface{})
		if !ok || mapOption == nil {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
		}
		idVal, idOk := mapOption["id"]
		if !idOk || idVal == "" {
			return errProxy.Errorf(common.CCErrCommParamsNeedSet, "option id")
		}
		if idValStr, ok := idVal.(string); !ok {
			return errProxy.Errorf(common.CCErrCommParamsNeedString, "option id")
		} else if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(idValStr) {
			return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option id",
				common.AttributeOptionValueMaxLength)
		}

		nameVal, nameOk := mapOption["name"]
		if !nameOk || nameVal == "" {
			return errProxy.Errorf(common.CCErrCommParamsNeedSet, "option name")
		}

		isDefault, ok := mapOption["is_default"]
		if !ok {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option default")
		}
		isDefaultVal, ok := isDefault.(bool)
		if !ok {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option default")
		}
		if isDefaultVal {
			count += 1
		}

		switch mapOption["type"] {
		case "text":
			if nameValStr, ok := nameVal.(string); !ok {
				return errProxy.Errorf(common.CCErrCommParamsNeedString, "option name")
			} else if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(nameValStr) {
				return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option name",
					common.AttributeOptionValueMaxLength)
			}
		default:
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option type")
		}
	}

	if !isMultiple && count != 1 {
		return errProxy.CCError(common.CCErrCommParamsNeedSingleChoice)
	}

	return nil
}

// ValidFieldTypeInt validate int or float field type's option and default value
func ValidFieldTypeInt(option, defaultVal interface{}, errProxy ccErr.DefaultCCErrorIf) error {
	if option == nil {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	tmp, ok := option.(map[string]interface{})
	if !ok {
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	{
		// min
		min, ok := tmp["min"]
		maxVal := 99999999999 // default
		minVal := -9999999999 // default
		err := errProxy.Error(common.CCErrCommParamsNeedInt)

		isPass := false
		if ok {
			switch d := min.(type) {
			case string:
				if len(d) == 0 {
					isPass = true
				}
				if len(d) > 11 {
					return errProxy.Errorf(common.CCErrCommOverLimit, "option.min")
				}
			}

			if !isPass {
				if ok := IsNumeric(min); !ok {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.min")
				}
				minVal, err = GetIntByInterface(min)
				if err != nil {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.min")
				}
			}
		}

		// max
		max, ok := tmp["max"]
		if ok {
			isPass := false
			switch d := max.(type) {
			case string:
				if len(d) == 0 {
					isPass = true
				}
				if len(d) > 11 {
					return errProxy.Errorf(common.CCErrCommOverLimit, "option.max")
				}
			}
			if !isPass {
				if ok := IsNumeric(max); !ok {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.max")
				}
				maxVal, err = GetIntByInterface(max)
				if nil != err {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.max")
				}
			}
		}

		if minVal > maxVal {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option.max")
		}
		if defaultVal == nil {
			return nil
		}
		defaultValue, err := GetIntByInterface(defaultVal)
		if err != nil {
			return err
		}
		if defaultValue < minVal || defaultValue > maxVal {
			return fmt.Errorf("int type field default value over limit")
		}
	}

	return nil
}

// ValidFieldTypeFloat validate int or float field type's option default value
func ValidFieldTypeFloat(option, defaultVal interface{}, errProxy ccErr.DefaultCCErrorIf) error {
	if option == nil {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	tmp, ok := option.(map[string]interface{})
	if !ok {
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	{
		// min
		min, ok := tmp["min"]
		maxVal := float64(common.MaxInt64) // default
		minVal := float64(common.MinInt64) // default
		err := errProxy.Error(common.CCErrCommParamsNeedInt)

		isPass := false
		if ok {
			switch d := min.(type) {
			case string:
				if len(d) == 0 {
					isPass = true
				}
				if len(d) > 11 {
					return errProxy.Errorf(common.CCErrCommOverLimit, "option.min")
				}
			}

			if !isPass {
				if ok := IsNumeric(min); !ok {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.min")
				}
				minVal, err = GetFloat64ByInterface(min)
				if err != nil {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.min")
				}
			}
		}

		// max
		max, ok := tmp["max"]
		if ok {
			isPass := false
			switch d := max.(type) {
			case string:
				if len(d) == 0 {
					isPass = true
				}
				if len(d) > 11 {
					return errProxy.Errorf(common.CCErrCommOverLimit, "option.max")
				}
			}
			if !isPass {
				if ok := IsNumeric(max); !ok {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.max")
				}
				maxVal, err = GetFloat64ByInterface(max)
				if nil != err {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.max")
				}
			}
		}

		if minVal > maxVal {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option.max")
		}
		if defaultVal == nil {
			return nil
		}
		defaultValue, err := GetFloat64ByInterface(defaultVal)
		if err != nil {
			return err
		}
		if defaultValue < minVal || defaultValue > maxVal {
			return fmt.Errorf("float type field default value over limit")
		}
	}
	return nil
}

// ValidFieldTypeList validate list field type's option and default value
func ValidFieldTypeList(option, defaultVal interface{}, errProxy ccErr.DefaultCCErrorIf) error {
	if option == nil {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	arrOption, ok := option.([]interface{})
	if !ok {
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}
	if len(arrOption) > common.AttributeOptionArrayMaxLength {
		return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionArrayMaxLength)
	}

	valueList := make([]string, len(arrOption))
	for _, val := range arrOption {
		switch value := val.(type) {
		case string: // 只可以是字符类型
			if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(value) {
				return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option",
					common.AttributeOptionValueMaxLength)
			}

			valueList = append(valueList, value)
		default:
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "list option need string type item")
		}
	}

	// 没有默认值，直接返回
	if defaultVal == nil {
		return nil
	}

	listDefaultVal := GetStrByInterface(defaultVal)
	for _, value := range valueList {
		if listDefaultVal == value {
			return nil
		}
	}

	return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "list default value")
}

// ValidFieldTypeString validate string field type's regex option and default value
func ValidFieldTypeString(option, defaultVal interface{}, errProxy ccErr.DefaultCCErrorIf) error {
	if option == nil {
		return nil
	}

	// 校验正则是否合法
	regular, ok := option.(string)
	if !ok {
		return errProxy.Errorf(common.CCIllegalRegularExpression, "option")
	}

	if len(regular) == 0 && defaultVal == nil {
		return nil
	}

	if len(regular) == 0 && defaultVal != nil {
		if _, ok := defaultVal.(string); !ok {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "string default value")
		}
	}

	if len(regular) > 0 && defaultVal == nil {
		if _, err := regexp.Compile(regular); err != nil {
			return errProxy.Errorf(common.CCErrorCheckRegularFailed, "regular is wrong")
		}
	}

	if len(regular) > 0 && defaultVal != nil {
		if _, err := regexp.Compile(regular); err != nil {
			return errProxy.Errorf(common.CCErrorCheckRegularFailed, "regular is wrong")
		}

		stringDefaultVal, ok := defaultVal.(string)
		if !ok {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "string default value")
		}

		match, err := regexp.MatchString(regular, stringDefaultVal)
		if err != nil || !match {
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "string default value")
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

// IsInnerObject is inner object model
func IsInnerObject(objID string) bool {
	switch objID {
	case common.BKInnerObjIDApp:
		return true
	case common.BKInnerObjIDBizSet:
		return true
	case common.BKInnerObjIDProject:
		return true
	case common.BKInnerObjIDHost:
		return true
	case common.BKInnerObjIDModule:
		return true
	case common.BKInnerObjIDPlat:
		return true
	case common.BKInnerObjIDProc:
		return true
	case common.BKInnerObjIDSet:
		return true
	}

	return false
}

// IsNumeric judges if value is a number
func IsNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		return true
	}

	return false
}

// IsInteger judges if value is a integer
func IsInteger(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, json.Number:
		return true
	}

	return false
}

// IsBasicValue test if an interface is the basic supported golang type or not.
func IsBasicValue(value interface{}) bool {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Float32,
		reflect.Float64,
		reflect.String:
		return true
	default:
		return false
	}
}

// ValidateStringType validate if the value is a string type
func ValidateStringType(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.String {
		return fmt.Errorf("value(%+v) is not of string type", value)
	}
	return nil
}

// ValidateBoolType validate if the value is a bool type
func ValidateBoolType(value interface{}) error {
	if reflect.TypeOf(value).Kind() != reflect.Bool {
		return fmt.Errorf("value(%+v) is not of bool type", value)
	}
	return nil
}

// ValidateNotEmptyStringType validate if the value is a not empty string type
func ValidateNotEmptyStringType(value interface{}) error {
	strVal, ok := value.(string)
	if !ok {
		return fmt.Errorf("value(%+v) is not of string type", value)
	}

	if len(strVal) == 0 {
		return errors.New("value is empty")
	}
	return nil
}

// ValidateDatetimeType validate if the value is a datetime type
func ValidateDatetimeType(value interface{}) error {
	// time type is supported
	if _, ok := value.(time.Time); ok {
		return nil
	}

	// timestamp type is supported
	if IsNumeric(value) {
		return nil
	}

	// string type with time format is supported
	if _, ok := IsTime(value); ok {
		return nil
	}
	return fmt.Errorf("value(%+v) is not of time type", value)
}

// ValidateSliceOfBasicType validate if the value is a slice of basic type
func ValidateSliceOfBasicType(value interface{}, limit uint) error {
	if value == nil {
		return errors.New("value is nil")
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.Array:
	case reflect.Slice:
	default:
		return fmt.Errorf("value(%+v) is not of array type", value)
	}

	v := reflect.ValueOf(value)
	length := v.Len()
	if length == 0 {
		return errors.New("value is empty")
	}

	if length > int(limit) {
		return fmt.Errorf("elements length %d exceeds maximum %d", length, limit)
	}

	// each element in the array or slice should be of the same basic type.
	var typ string
	for i := 0; i < length; i++ {
		item := v.Index(i).Interface()

		var itemType string
		switch item.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
			itemType = "numeric"
		case bool:
			itemType = "bool"
		case string:
			itemType = "string"
		default:
			return fmt.Errorf("array element index(%d) value(%+v) is not of basic type", i, item)
		}

		if i == 0 {
			typ = itemType
			continue
		}

		if typ != itemType {
			return fmt.Errorf("array element index(%d) value(%+v) type is not %s", i, item, typ)
		}
	}

	return nil
}

var mainlineNameRegexp = regexp.MustCompile(common.FieldTypeMainlineRegexp)

// ValidTopoNameField validate business topology name, including set and service templates that may generate them
func ValidTopoNameField(name string, nameField string, errProxy ccErr.DefaultCCErrorIf) (string, error) {
	name = strings.Trim(name, " ")

	if len(name) == 0 {
		return name, errProxy.CCErrorf(common.CCErrCommParamsNeedSet, nameField)
	}

	if utf8.RuneCountInString(name) > common.MainlineNameFieldMaxLength {
		return name, errProxy.CCErrorf(common.CCErrCommValExceedMaxFailed, nameField, common.MainlineNameFieldMaxLength)
	}

	match := mainlineNameRegexp.MatchString(name)
	if !match {
		return name, errProxy.CCErrorf(common.CCErrCommParamsInvalid, nameField)
	}

	return name, nil
}

// ValidMustSetStringField valid if the value is of string type and is not empty
func ValidMustSetStringField(value interface{}, field string, errProxy ccErr.DefaultCCErrorIf) (string, error) {
	switch val := value.(type) {
	case string:
		if len(val) == 0 {
			return val, errProxy.Errorf(common.CCErrCommParamsNeedSet, field)
		}
		return val, nil
	default:
		return "", errProxy.New(common.CCErrCommParamsNeedString, field)
	}
}

// ValidModelIDField validate model related id field, like classification id, attribute id, group id...
func ValidModelIDField(value interface{}, field string, errProxy ccErr.DefaultCCErrorIf) error {
	strValue, err := ValidMustSetStringField(value, field, errProxy)
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(strValue) > common.AttributeIDMaxLength {
		return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, field, common.AttributeIDMaxLength)
	}

	match, err := regexp.MatchString(common.FieldTypeStrictCharRegexp, strValue)
	if nil != err {
		return err
	}
	if !match {
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, field)
	}
	return nil
}

// ValidModelNameField validate model related name field, like classification name, attribute name, group name...
func ValidModelNameField(value interface{}, field string, errProxy ccErr.DefaultCCErrorIf) error {
	strValue, err := ValidMustSetStringField(value, field, errProxy)
	if err != nil {
		return err
	}

	if utf8.RuneCountInString(strValue) > common.AttributeNameMaxLength {
		return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, field, common.AttributeNameMaxLength)
	}
	return nil
}

// ValidPropertyTypeIsMultiple valid object attr field type is multiple
func ValidPropertyTypeIsMultiple(propertyType string, isMultiple bool, errProxy ccErr.DefaultCCErrorIf) error {
	switch propertyType {
	case common.FieldTypeSingleChar, common.FieldTypeInt, common.FieldTypeFloat, common.FieldTypeEnum,
		common.FieldTypeDate, common.FieldTypeTime, common.FieldTypeLongChar, common.FieldTypeTimeZone,
		common.FieldTypeBool, common.FieldTypeList:
		if isMultiple {
			return errProxy.Errorf(common.CCErrCommFieldTypeNotSupportMultiple, propertyType)
		}
	}
	return nil
}
