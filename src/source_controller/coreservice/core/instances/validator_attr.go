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

package instances

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

// validTime valid object Attribute that is time type
func (valid *validator) validTime(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params can not be null, rid: %s", rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	valStr, ok := val.(string)
	if false == ok {
		blog.Errorf("date can should be string, rid: %s", rid)
		return valid.errif.Errorf(common.CCErrCommParamsShouldBeString, key)
	}

	result := util.IsTime(valStr)
	if !result {
		blog.Errorf("params not valid, rid: %s", rid)
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validDate valid object Attribute that is date type
func (valid *validator) validDate(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params key: %s can not be null, rid: %s", key, rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}
	valStr, ok := val.(string)
	if false == ok {
		blog.Errorf("date should be string, rid: %s", rid)
		return valid.errif.Errorf(common.CCErrCommParamsShouldBeString, key)

	}
	result := util.IsDate(valStr)
	if !result {
		blog.Errorf("params key: %s is not valid, rid: %s", valStr, rid)
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validEnum valid object attribute that is enum type
func (valid *validator) validEnum(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	// validate require
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params key :%s, can not be null, rid: %s", key, rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	// validate type
	valStr, ok := val.(string)
	if !ok {
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}

	option, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	// validate within enum
	enumOption, err := metadata.ParseEnumOption(ctx, option.Option)
	if err != nil {
		blog.Warnf("ParseEnumOption failed: %v, rid: %s", err, rid)
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	match := false
	for _, k := range enumOption {
		if k.ID == valStr {
			match = true
			break
		}
	}
	if !match {
		blog.Errorf("params %s not valid, option %#v, raw option %#v, value: %#v, rid: %s", key, enumOption, option, val, rid)
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validBool valid object attribute that is bool type
func (valid *validator) validBool(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params key: %s can not be null, rid: %s", key, rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch val.(type) {
	case bool:
	default:
		blog.Errorf("params key: %s should be bool, rid: %s", key, rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedBool, key)
	}
	return nil
}

// valid char valid object attribute that is timezone type
func (valid *validator) validTimeZone(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params key: %s can not be null, rid: %s", key, rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch value := val.(type) {
	case string:
		isMatch := util.IsTimeZone(value)
		if false == isMatch {
			blog.Errorf("params key: %s should be timezone, rid: %s", key, rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedTimeZone, key)
		}
	default:
		blog.Errorf("params key: %s should be timezone, rid: %s", key, rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedTimeZone, key)
	}
	return nil
}

// validInt valid object attribute that is int type
func (valid *validator) validInt(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params can not be null, rid: %s", rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	if !valid.isNumeric(val) {
		blog.Errorf("params %s:%#v not int, rid: %s", key, val, rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	var value int64
	value, err := util.GetInt64ByInterface(val)
	if nil != err {
		blog.Errorf("params %s:%#v not int, rid: %s", key, val, rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	property, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	intObjOption := metadata.ParseIntOption(ctx, property.Option)
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return nil
	}

	maxValue, err := strconv.ParseInt(intObjOption.Max, 10, 64)
	if nil != err {
		maxValue = common.MaxInt64
	}
	minValue, err := strconv.ParseInt(intObjOption.Min, 10, 64)
	if nil != err {
		minValue = common.MinInt64
	}
	if value > maxValue || value < minValue {
		blog.Errorf("params %s:%#v not valid, rid: %s", key, val, rid)
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validFloat valid object attribute that is float type
func (valid *validator) validFloat(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val {
		if valid.require[key] {
			blog.Errorf("params can not be null, rid: %s", rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	if !valid.isNumeric(val) {
		blog.Errorf("params %s:%#v not int, rid: %s", key, val, rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	var value float64
	value, err := util.GetFloat64ByInterface(val)
	if nil != err {
		blog.Errorf("params %s:%#v not float, rid: %s", key, val, rid)
		return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
	}

	property, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	intObjOption := metadata.ParseFloatOption(ctx, property.Option)
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return nil
	}

	maxValue, err := strconv.ParseFloat(intObjOption.Max, 64)
	if nil != err {
		maxValue = float64(common.MaxInt64)
	}
	minValue, err := strconv.ParseFloat(intObjOption.Min, 64)
	if nil != err {
		minValue = float64(common.MinInt64)
	}
	if value > maxValue || value < minValue {
		blog.Errorf("params %s:%#v not valid, rid: %s", key, val, rid)
		return valid.errif.CCErrorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validInt valid object attribute that is long char type
func (valid *validator) validLongChar(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val || "" == val {
		if valid.require[key] {
			blog.Errorf("params in need, rid: %s", rid)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch value := val.(type) {
	case string:
		value = strings.TrimSpace(value)
		if len(value) > common.FieldTypeLongLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeSingleLenChar, rid)
			return valid.errif.Errorf(common.CCErrCommOverLimit, key)
		}
		if 0 == len(value) {
			if valid.require[key] {
				blog.Errorf("params can not be empty, rid: %s", rid)
				return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			return nil
		}

		match, err := regexp.MatchString(common.FieldTypeLongCharRegexp, value)
		if nil != err || !match {
			blog.Errorf(`params "%s" not match longchar regexp, rid:  %s`, val, rid)
			return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
		}

		if property, ok := valid.propertys[key]; ok && "" != val {
			option, ok := property.Option.(string)
			if !ok {
				break
			}
			strReg, err := regexp.Compile(option)
			if nil != err {
				blog.Errorf(`params "%s" not match regexp "%s", rid: %s`, val, option, rid)
				return valid.errif.Errorf(common.CCErrFieldRegValidFailed, key)
			}
			if !strReg.MatchString(value) {
				blog.Errorf(`params "%s" not match regexp "%s", rid: %s`, val, option, rid)
				return valid.errif.Errorf(common.CCErrFieldRegValidFailed, key)
			}
		}
	default:
		blog.Errorf("params should be string, rid: %s", rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedString, key)
	}

	return nil
}

// validChar valid object attribute that is  char type
func (valid *validator) validChar(ctx context.Context, val interface{}, key string) error {
	rid := util.ExtractRequestIDFromContext(ctx)
	if nil == val || "" == val {
		if valid.require[key] {
			blog.Errorf("params in need, rid: %s", rid)
			return valid.errif.CCErrorf(common.CCErrCommParamsNeedSet, key)
		}
		return nil
	}
	switch value := val.(type) {
	case string:
		if len(value) > common.FieldTypeSingleLenChar {
			blog.Errorf("params over length %d, rid: %s", common.FieldTypeSingleLenChar, rid)
			return valid.errif.CCErrorf(common.CCErrCommOverLimit, key)
		}
		if 0 == len(value) {
			if valid.require[key] {
				blog.Errorf("params can not be empty, rid: %s", rid)
				return valid.errif.CCErrorf(common.CCErrCommParamsNeedSet, key)
			}
			return nil
		}

		value = strings.TrimSpace(value)
		match, err := regexp.MatchString(common.FieldTypeSingleCharRegexp, value)
		if nil != err || !match {
			blog.Errorf(`params "%s" not match singlechar regexp, rid:  %s`, val, rid)
			return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
		}

		if property, ok := valid.propertys[key]; ok && "" != val {
			option, ok := property.Option.(string)
			if !ok {
				break
			}
			strReg, err := regexp.Compile(option)
			if nil != err {
				blog.Errorf(`params "%s" not match regexp "%s", rid:  %s`, val, option, rid)
				return valid.errif.CCErrorf(common.CCErrFieldRegValidFailed, key)
			}
			if !strReg.MatchString(value) {
				blog.Errorf(`params "%s" not match regexp "%s", rid: %s`, val, option, rid)
				return valid.errif.CCErrorf(common.CCErrFieldRegValidFailed, key)
			}
		}
	default:
		blog.Errorf("params should be string, rid: %s", rid)
		return valid.errif.Errorf(common.CCErrCommParamsNeedString, key)
	}

	return nil
}

//valid list
func (valid *validator) validList(ctx context.Context, val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null, list field key: %s", key)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
		}
		return nil
	}

	strVal, ok := val.(string)
	if !ok {
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}

	option, ok := valid.propertys[key]
	if !ok {
		blog.Errorf("option %v invalid, not string type list option", option)
		return valid.errif.Errorf(common.CCErrCollectNetDeviceObjPropertyNotExist, key)
	}
	listOption, ok := option.Option.([]interface{})
	if false == ok {
		blog.Errorf("option %v invalid, not string type list option", option)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	match := false
	for _, inVal := range listOption {
		inValStr, ok := inVal.(string)
		if !ok {
			blog.Errorf("inner list option convert to string  failed, params %s not valid , list field value: %#v", key, val)
			return valid.errif.Errorf(common.CCErrParseAttrOptionListFailed, key)
		}
		if strVal == inValStr {
			match = true
			break
		}
	}
	if !match {
		blog.Errorf("params %s not valid, option %#v, raw option %#v, value: %#v", key, listOption, option, val)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// isNumeric judges if value is a number
func (valid *validator) isNumeric(val interface{}) bool {
	switch val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, json.Number:
		return true
	}

	return false
}
