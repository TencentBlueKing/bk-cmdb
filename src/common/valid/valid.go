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

package valid

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
	"configcenter/src/common/util"
)

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
	if util.IsNumeric(value) {
		return nil
	}

	if util.IsDate(value) {
		return nil
	}

	// string type with time format is supported
	if _, ok := util.IsTime(value); ok {
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
