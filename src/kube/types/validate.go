/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
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

package types

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"

	"configcenter/src/common"
	ccErr "configcenter/src/common/errors"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal/table"
)

// NumericSettings numeric type check parameter setting.
type NumericSettings struct {
	// Min Minimum value allowed for numeric types.
	Min int64

	// Max maximum value allowed for numeric types.
	Max int64
}

// ValidateNumeric 1、judgment is a number type. 2、the judgment is that they are all within the specified range.
func ValidateNumeric(data interface{}, param NumericSettings) error {

	if data == nil {
		return errors.New("data is nil")
	}

	v, err := util.GetInt64ByInterface(data)
	if err != nil {
		return err
	}

	if v > param.Max || v < param.Min {
		return fmt.Errorf("data: %d out of range [min: %d - max: %d]", v, param.Min, param.Max)
	}
	return nil
}

const (
	// FirstLevel the first level of key value type data.
	FirstLevel = 1

	defaultStringLength   = 256
	defaultKeyValueNumber = 10
	defaultMaximumLevel   = 3
	defaultMaxArrayLength = 10
)

// StringSettings string type check parameter setting.
type StringSettings struct {
	// MaxLength maximum length allowed for string type.
	MaxLength int

	// RegularCheck regular expressions involved in strings
	RegularCheck string
}

// ValidateString judgment of string data:
// 1、judgment type.
// 2、judgment length.
// 3、check if the regular expression is satisfied if necessary.
// 4、if the maximum length of the string is not set, the default maximum length is 256.
func ValidateString(data interface{}, param StringSettings) error {

	if data == nil {
		return errors.New("data is nil")
	}

	tmpType := reflect.TypeOf(data)

	if tmpType.Kind() != reflect.String {
		return errors.New("data type is not string")
	}

	if param.MaxLength == 0 {
		param.MaxLength = defaultStringLength
	}

	v := data.(string)
	if len(v) > param.MaxLength {
		return fmt.Errorf("data length(%d) is exceeded max length: %d", len(v), param.MaxLength)
	}

	if len(param.RegularCheck) < 0 {
		return nil
	}

	if !regexp.MustCompile(param.RegularCheck).MatchString(v) {
		return fmt.Errorf("invalid data %s, regular is %s", v, param.RegularCheck)
	}

	return nil
}

// ValidateBoolen Boolean type judgment
func ValidateBoolen(data interface{}) error {
	if data == nil {
		return errors.New("data is nil")
	}

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Bool {
		return fmt.Errorf("data is not of type boolean, type: %v", v.Kind())
	}

	return nil
}

// ValidateMapString mapstring type judgment：
// 1、type must be map.
// 2、the type of key and value must be string.
// 3、check the number of key-value pairs.
// 4、the default maximum number of key-value pairs allowed is 10.
func ValidateMapString(data interface{}, length int) error {
	if data == nil {
		return errors.New("data is nil")
	}

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.Map {
		return fmt.Errorf("data type is error, type: %v", v.Kind())
	}
	if length == 0 {
		length = defaultKeyValueNumber
	}

	mapKeys := v.MapKeys()
	if len(mapKeys) > length {
		return fmt.Errorf("data length is exceeded max length %d", length)
	}

	for _, key := range mapKeys {
		if key.Kind() != reflect.String || v.MapIndex(key).Kind() != reflect.String {
			return fmt.Errorf("data key or value type is not string")
		}
	}
	return nil
}

// ArraySettingsParam parameter settings of the array type,
// including fine-grained check parameter settings for each element of the array
type ArraySettingsParam struct {
	// ArrayMaxLength maximum length of an array allowed.
	ArrayMaxLength int
	MapObjectParam MapObjectSettings `json:"map_object_param"`
	NumericParam   NumericSettings   `json:"numeric_param"`
	StringParam    StringSettings    `json:"string_param"`
}

// ValidateArray 1、the length of the check array. If maxLength is set to 0,
// it means that the length will not be checked.
// 2、determines the type of array elements. Currently only bool, numeric,
// string, and map are supported.
// The rest of the types are not supported, nor are multidimensional arrays supported.
// 3、when the elements of the array are of type map, the maximum nesting
// level maxDeep of the map needs to be set.
func ValidateArray(data interface{}, param *ArraySettingsParam) error {
	if data == nil {
		return errors.New("data is nil")
	}

	v := reflect.ValueOf(data)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		if param.ArrayMaxLength == 0 {
			param.ArrayMaxLength = defaultMaxArrayLength
		}
		if v.Len() > param.ArrayMaxLength {
			return fmt.Errorf("data len exceed max length: %d", param.ArrayMaxLength)
		}

		for i := 0; i < v.Len(); i++ {
			switch v.Index(i).Kind() {
			case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint16,
				reflect.Uint32, reflect.Uint64, reflect.Uint8:
				if err := ValidateNumeric(v.Index(i).Interface(), param.NumericParam); err != nil {
					return err
				}
			case reflect.String:
				if err := ValidateString(v.Index(i).Interface(), param.StringParam); err != nil {
					return err
				}

			case reflect.Map:
				if err := ValidateKVObject(v.Index(i).Interface(), param.MapObjectParam, 1); err != nil {
					return err
				}

			case reflect.Bool:
				if err := ValidateBoolen(v.Index(i).Interface()); err != nil {
					return err
				}

			default:
				return fmt.Errorf("unsupported type: %v", v.Index(i).Kind())
			}
		}
	default:
		return errors.New("data type is not array")
	}
	return nil
}

// MapObjectSettings mapObject type check parameter settings
type MapObjectSettings struct {
	// MaxDeep if the array element is the maximum level allowed by the map object
	MaxDeep int

	// MaxLength the maximum number of elements per level allowed by the array element if the map object.
	MaxLength int
}

// ValidateKVObject for the verification of general kv type data. where deep represents the current level,
// The user needs to specify deep as "FirstLevel"
func ValidateKVObject(data interface{}, param MapObjectSettings, deep int) error {

	if data == nil {
		return errors.New("data is nil")
	}

	if deep == FirstLevel && param.MaxDeep == 0 {
		param.MaxDeep = defaultMaximumLevel
	}

	if deep == FirstLevel && param.MaxLength == 0 {
		param.MaxLength = defaultKeyValueNumber
	}

	if deep > param.MaxDeep {
		return fmt.Errorf("exceed max deep: %d", param.MaxDeep)
	}

	v := reflect.ValueOf(data)

	switch v.Kind() {
	case reflect.Map:
		mapKeys := v.MapKeys()
		if len(mapKeys) > param.MaxLength {
			return fmt.Errorf("keys length exceed than %d", param.MaxLength)
		}

		for _, key := range mapKeys {
			keyValue := v.MapIndex(key)
			switch keyValue.Kind() {
			// compatible with the scenario where the value is a string.
			case reflect.Interface:
				if err := convertInterfaceIntoMap(keyValue.Interface(), param, deep); err != nil {
					return err
				}
			case reflect.Map:
				if err := ValidateKVObject(keyValue.Interface(), param, deep+1); err != nil {
					return err
				}
			case reflect.String:
			case reflect.Int8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64,
				reflect.Float32, reflect.Float64:
			default:
				return errors.New("data type error")
			}
		}
	case reflect.Struct:
	case reflect.Interface:
		if err := convertInterfaceIntoMap(v.Interface(), param, deep); err != nil {
			return err
		}
	default:
		return fmt.Errorf("data type is error, type: %v", v.Kind())
	}
	return nil
}

func convertInterfaceIntoMap(target interface{}, param MapObjectSettings, deep int) error {

	value := reflect.ValueOf(target)
	switch value.Kind() {
	// compatible with the scenario where the value is a string.
	case reflect.String:
	case reflect.Map:
		if err := ValidateKVObject(value, param, deep); err != nil {
			return err
		}
	case reflect.Struct:
	default:
		return fmt.Errorf("no support the kind(%s)", value.Kind())
	}
	return nil
}

// ValidateCreate validate create data struct
// NOTE:
// 1. data must be a value type, not a pointer.
func ValidateCreate(data interface{}, field *table.Fields) ccErr.RawErrorInfo {
	if data == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if field == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"field"},
		}
	}

	typeOfOption := reflect.TypeOf(data)
	valueOfOption := reflect.ValueOf(data)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag, flag := getFieldTag(typeOfOption, JsonTag, i)
		if flag {
			continue
		}

		if !field.IsFieldRequiredByField(tag) {
			continue
		}

		if err := isRequiredField(tag, valueOfOption, i); err != nil {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{tag},
			}
		}
	}

	return ccErr.RawErrorInfo{}
}

// ValidateUpdate validate update data struct
// NOTE:
// 1. data must be a value type, not a pointer.
func ValidateUpdate(data interface{}, field *table.Fields) ccErr.RawErrorInfo {
	if data == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"data"},
		}
	}

	if field == nil {
		return ccErr.RawErrorInfo{
			ErrCode: common.CCErrCommParamsNeedSet,
			Args:    []interface{}{"field"},
		}
	}

	typeOfOption := reflect.TypeOf(data)
	valueOfOption := reflect.ValueOf(data)
	for i := 0; i < typeOfOption.NumField(); i++ {
		tag, flag := getFieldTag(typeOfOption, JsonTag, i)
		if flag {
			continue
		}

		if flag := isNotEditableField(tag, valueOfOption, i); flag {
			continue
		}

		// get whether it is an editable field based on tag
		if !field.IsFieldEditableByField(tag) {
			return ccErr.RawErrorInfo{
				ErrCode: common.CCErrCommParamsIsInvalid,
				Args:    []interface{}{tag},
			}
		}
	}

	return ccErr.RawErrorInfo{}
}
