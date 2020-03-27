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
	"unicode/utf8"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
)

// ValidPropertyOption valid property field option
func ValidPropertyOption(propertyType string, option interface{}, errProxy errors.DefaultCCErrorIf) error {
	switch propertyType {
	case common.FieldTypeEnum:
		return ValidFieldTypeEnumOption(option, errProxy)
	case common.FieldTypeInt:
		return ValidFieldTypeIntOption(option, errProxy)
	case common.FieldTypeList:
		return ValidFieldTypeListOption(option, errProxy)
	}
	return nil
}

func ValidFieldTypeEnumOption(option interface{}, errProxy errors.DefaultCCErrorIf) error {
	if nil == option {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	arrOption, ok := option.([]interface{})
	if false == ok {
		blog.Errorf(" option %v not enum option", option)
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}

	if len(arrOption) > common.AttributeOptionArrayMaxLength {
		blog.Errorf(" option array length %d exceeds max length %d", len(arrOption), common.AttributeOptionArrayMaxLength)
		return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionArrayMaxLength)
	}
	for _, o := range arrOption {
		mapOption, ok := o.(map[string]interface{})
		if false == ok || mapOption == nil {
			blog.Errorf(" option %v not enum option, enum option item must id and name", option)
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
		}
		idVal, idOk := mapOption["id"]
		if !idOk || idVal == "" {
			blog.Errorf("enum option id can't be empty", option)
			return errProxy.Errorf(common.CCErrCommParamsNeedSet, "option id")
		}
		if idValStr, ok := idVal.(string); !ok {
			blog.Errorf("idVal %v not string", idVal)
			return errProxy.Errorf(common.CCErrCommParamsNeedString, "option id")
		} else if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(idValStr) {
			blog.Errorf(" option id %s length %d exceeds max length %d", idValStr, utf8.RuneCountInString(idValStr), common.AttributeOptionValueMaxLength)
			return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option id", common.AttributeOptionValueMaxLength)
		}
		nameVal, nameOk := mapOption["name"]
		if !nameOk || nameVal == "" {
			blog.Errorf("enum option name can't be empty", option)
			return errProxy.Errorf(common.CCErrCommParamsNeedSet, "option name")
		}
		switch mapOption["type"] {
		case "text":
			if nameValStr, ok := nameVal.(string); !ok {
				blog.Errorf(" nameVal %v not string", nameVal)
				return errProxy.Errorf(common.CCErrCommParamsNeedString, "option name")
			} else if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(nameValStr) {
				blog.Errorf(" option name %s length %d exceeds max length %d", nameValStr, utf8.RuneCountInString(nameValStr), common.AttributeOptionValueMaxLength)
				return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option name", common.AttributeOptionValueMaxLength)
			}
		default:
			blog.Errorf("enum option type must be 'text', current: %v", mapOption["type"])
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option type")
		}
	}

	return nil
}

func ValidFieldTypeIntOption(option interface{}, errProxy errors.DefaultCCErrorIf) error {
	if nil == option {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	tmp, ok := option.(map[string]interface{})
	if false == ok {
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
				if 0 == len(d) {
					isPass = true
				}
				if 11 < len(d) {
					return errProxy.Errorf(common.CCErrCommOverLimit, "option.min")
				}
			}

			if !isPass {
				if ok := IsNumeric(min); !ok {
					return errProxy.Errorf(common.CCErrCommParamsNeedInt, "option.min")
				}
				minVal, err = GetIntByInterface(min)
				if nil != err {
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
				if 0 == len(d) {
					isPass = true
				}
				if 11 < len(d) {
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
	}

	return nil
}

func ValidFieldTypeListOption(option interface{}, errProxy errors.DefaultCCErrorIf) error {
	if nil == option {
		return errProxy.Errorf(common.CCErrCommParamsLostField, "option")
	}

	arrOption, ok := option.([]interface{})
	if false == ok {
		blog.Errorf(" option %v not string type list option", option)
		return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "option")
	}
	if len(arrOption) > common.AttributeOptionArrayMaxLength {
		blog.Errorf(" option array length %d exceeds max length %d", len(arrOption), common.AttributeOptionArrayMaxLength)
		return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionArrayMaxLength)
	}

	for _, val := range arrOption {
		switch value := val.(type) {
		case string: // 只可以是字符类型
			if common.AttributeOptionValueMaxLength < utf8.RuneCountInString(value) {
				blog.Errorf(" option value %s length %d exceeds max length %d", value, utf8.RuneCountInString(value), common.AttributeOptionValueMaxLength)
				return errProxy.Errorf(common.CCErrCommValExceedMaxFailed, "option", common.AttributeOptionValueMaxLength)
			}
		default:
			blog.Errorf(" option %v not string type list option", option)
			return errProxy.Errorf(common.CCErrCommParamsIsInvalid, "list option need string type item")
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
