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
	"regexp"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/util"
)

// validTime valid object Attribute that is time type
func (valid *validator) validTime(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	valStr, ok := val.(string)
	if false == ok {
		blog.Error("date can shoule be string")
		return valid.errif.Errorf(common.CCErrCommParamsShouldBeString, key)
	}

	result := util.IsTime(valStr)
	if !result {
		blog.Error("params   not valid")
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validDate valid object Attribute that is date type
func (valid *validator) validDate(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}
	valStr, ok := val.(string)
	if false == ok {
		blog.Error("date can shoule be string")
		return valid.errif.Errorf(common.CCErrCommParamsShouldBeString, key)

	}
	result := util.IsDate(valStr)
	if !result {
		blog.Error("params  is not valid")
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validEnum valid object attribute that is enum type
func (valid *validator) validEnum(val interface{}, key string) error {
	// validate require
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	// validate type
	valStr, ok := val.(string)
	if !ok {
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}

	option, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	// validate within enum
	enumOption := ParseEnumOption(option.Option)
	match := false
	for _, k := range enumOption {
		if k.ID == valStr {
			match = true
			break
		}
	}
	if !match {
		blog.V(3).Infof("params %s not valid, option %#v, raw option %#v, value: %#v", key, enumOption, option, val)
		blog.Errorf("params %s not valid , enum value: %#v", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

//validBool valid object attribute that is bool type
func (valid *validator) validBool(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch val.(type) {
	case bool:
	default:
		blog.Error("params should be  bool")
		return valid.errif.Errorf(common.CCErrCommParamsNeedBool, key)
	}
	return nil
}

//valid char valid object attribute that is timezone type
func (valid *validator) validTimeZone(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch value := val.(type) {
	case string:
		isMatch := util.IsTimeZone(value)
		if false == isMatch {
			blog.Error("params should be  timezone")
			return valid.errif.Errorf(common.CCErrCommParamsNeedTimeZone, key)
		}
	default:
		blog.Error("params should be  timezone")
		return valid.errif.Errorf(common.CCErrCommParamsNeedTimeZone, key)
	}
	return nil
}

// validForeignKey valid object attribute that is foreign key type
func (valid *validator) validForeignKey(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	_, ok := util.GetTypeSensitiveUInt64(val)
	if !ok {
		blog.Errorf("params %s:%#v not int", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	return nil
}

// validInt valid object attribute that is int type
func (valid *validator) validInt(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	var value int64
	value, err := util.GetInt64ByInterface(val)
	if nil != err {
		blog.Errorf("params %s:%#v not int", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	property, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	intObjOption := parseIntOption(property.Option)
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
		blog.Errorf("params %s:%#v not valid", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validFloat valid object attribute that is float type
func (valid *validator) validFloat(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	var value float64
	value, err := util.GetFloat64ByInterface(val)
	if nil != err {
		blog.Errorf("params %s:%#v not float", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
	}

	property, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	intObjOption := parseFloatOption(property.Option)
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
		blog.Errorf("params %s:%#v not valid", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validInt valid object attribute that is long char type
func (valid *validator) validLongChar(val interface{}, key string) error {
	if nil == val || "" == val {
		if valid.require[key] {
			blog.Error("params in need")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch value := val.(type) {
	case string:
		if len(value) > common.FieldTypeLongLenChar {
			blog.Errorf("params over length %d", common.FieldTypeSingleLenChar)
			return valid.errif.Errorf(common.CCErrCommOverLimit, key)
		}
		if 0 == len(value) {
			if valid.require[key] {
				blog.Error("params can not be empty")
				return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			return nil
		}

		if property, ok := valid.propertys[key]; ok && "" != val {
			option, ok := property.Option.(string)
			if !ok {
				break
			}
			strReg, err := regexp.Compile(option)
			if nil != err {
				blog.Errorf(`params "%s" not match regexp "%s"`, val, option)
				return valid.errif.Error(common.CCErrFieldRegValidFailed)
			}
			if !strReg.MatchString(value) {
				blog.Errorf(`params "%s" not match regexp "%s"`, val, option)
				return valid.errif.Error(common.CCErrFieldRegValidFailed)
			}
		}
	default:
		blog.Error("params should be  string")
		return valid.errif.Errorf(common.CCErrCommParamsNeedString, key)
	}

	return nil
}

// validChar valid object attribute that is  char type
func (valid *validator) validChar(val interface{}, key string) error {
	if nil == val || "" == val {
		if valid.require[key] {
			blog.Error("params in need")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
		}
		return nil
	}
	switch value := val.(type) {
	case string:
		if len(value) > common.FieldTypeSingleLenChar {
			blog.Errorf("params over length %d", common.FieldTypeSingleLenChar)
			return valid.errif.Errorf(common.CCErrCommOverLimit, key)
		}
		if 0 == len(value) {
			if valid.require[key] {
				blog.Error("params can not be empty")
				return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			return nil
		}

		if property, ok := valid.propertys[key]; ok && "" != val {
			option, ok := property.Option.(string)
			if !ok {
				break
			}
			strReg, err := regexp.Compile(option)
			if nil != err {
				blog.Errorf(`params "%s" not match regexp "%s"`, val, option)
				return valid.errif.Error(common.CCErrFieldRegValidFailed)
			}
			if !strReg.MatchString(value) {
				blog.Errorf(`params "%s" not match regexp "%s"`, val, option)
				return valid.errif.Error(common.CCErrFieldRegValidFailed)
			}
		}
	default:
		blog.Error("params should be  string")
		return valid.errif.Errorf(common.CCErrCommParamsNeedString, key)
	}

	return nil
}
