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

package validator

import (
	"context"
	"regexp"
	"strconv"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/storage/dal"
)

// NewValidator returns new Validator
func NewValidator(ownerID, objID string, db dal.RDB, ctx context.Context, defLang language.DefaultCCLanguageIf, errif errors.DefaultCCErrorIf) *Validator {
	return &Validator{
		ownerID: ownerID,
		objID:   objID,
		db:      db.Clone(),
		ctx:     ctx,
		defLang: defLang,
		errif:   errif,

		propertys:    map[string]metadata.Attribute{},
		require:      map[string]bool{},
		idToProperty: map[int64]metadata.Attribute{},
		shouldIgnore: map[string]bool{},
	}
}

// Init init
func (valid *Validator) Init(attrs []metadata.Attribute) {
	for _, attr := range attrs {
		if attr.PropertyID == common.BKChildStr || attr.PropertyID == common.BKParentStr {
			continue
		}
		valid.propertys[attr.PropertyID] = attr
		valid.idToProperty[attr.ID] = attr
		valid.propertyslice = append(valid.propertyslice, attr)
		if attr.IsRequired {
			valid.require[attr.PropertyID] = true
			valid.requirefields = append(valid.requirefields, attr.PropertyID)
		}
	}
}

// valid create request
func (valid *Validator) ValidateCreate(valData map[string]interface{}) error {
	ignoreKeys := []string{
		common.BKOwnerIDField,
		common.BKDefaultField,
		common.BKInstParentStr,
		common.BKOwnerIDField,
		common.BKAppIDField,
		common.BKSupplierIDField,
		common.BKInstIDField,
	}
	for _, item := range ignoreKeys {
		valid.shouldIgnore[item] = true
	}
	FillLostedFieldValue(valData, valid.propertyslice, valid.requirefields)
	for _, key := range valid.requirefields {
		if _, ok := valData[key]; !ok {
			blog.Errorf("params in need, valid %s, data: %+v", valid.objID, valData)
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
		}
	}
	err := valid.ValidateMap(valData)
	if err != nil {
		blog.Errorf("ValidateMap, err: %v", err)
		return err
	}
	return valid.validCreateUnique(valData)
}

// valid update request
func (valid *Validator) ValidateUpdate(valData map[string]interface{}, originalData map[string]interface{}) error {
	ignoreKeys := []string{
		common.BKOwnerIDField,
		common.BKDefaultField,
		common.BKInstParentStr,
		common.BKOwnerIDField,
		common.BKAppIDField,
		common.BKDataStatusField,
		common.BKDataStatusField,
		common.BKSupplierIDField,
		common.BKInstIDField,
	}
	for _, item := range ignoreKeys {
		valid.shouldIgnore[item] = true
	}

	err := valid.ValidateMap(valData)
	if err != nil {
		return err
	}
	return valid.validUpdateUnique(valData, originalData)
}

// ValidateMap basic valid
func (valid *Validator) ValidateMap(valData map[string]interface{}) error {
	var err error
	for key, val := range valData {
		if valid.shouldIgnore[key] {
			// ignore the key field
			continue
		}

		property, ok := valid.propertys[key]
		if !ok {
			blog.Errorf("params is not valid, the key is %s, properties is %#v", key, valid.propertys)
			return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
		}
		fieldType := property.PropertyType
		switch fieldType {
		case common.FieldTypeSingleChar:
			err = valid.validChar(val, key)
		case common.FieldTypeLongChar:
			err = valid.validLongChar(val, key)
		case common.FieldTypeInt:
			err = valid.validInt(val, key)
		case common.FieldTypeEnum:
			err = valid.validEnum(val, key)
		case common.FieldTypeDate:
			err = valid.validDate(val, key)
		case common.FieldTypeTime:
			err = valid.validTime(val, key)
		case common.FieldTypeTimeZone:
			err = valid.validTimeZone(val, key)
		case common.FieldTypeBool:
			err = valid.validBool(val, key)
		case common.FieldTypeForeignKey:
			err = valid.validForeignKey(val, key)
		case common.FieldTypeFloat:
			err = valid.validFloat(val, key)
		case common.FieldTypeUser:
			err = valid.validUser(val, key)
		default:
			continue
		}
		if nil != err {
			return err
		}
	}
	return nil
}

//valid char
func (valid *Validator) validChar(val interface{}, key string) error {
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

//valid long char
func (valid *Validator) validLongChar(val interface{}, key string) error {
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

// validInt valid int
func (valid *Validator) validInt(val interface{}, key string) error {
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
	intObjOption := parseMinMaxOption(property.Option)
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

// validForeignKey valid foreign key
func (valid *Validator) validForeignKey(val interface{}, key string) error {
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

//valid char
func (valid *Validator) validTimeZone(val interface{}, key string) error {
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

//valid char
func (valid *Validator) validUser(val interface{}, key string) error {
	if nil == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	switch val.(type) {
	case string:
	default:
		blog.Error("params should be string")
		return valid.errif.Errorf(common.CCErrCommParamsNeedString, key)
	}
	return nil
}

//validBool
func (valid *Validator) validBool(val interface{}, key string) error {
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

//validFloat
func (valid *Validator) validFloat(val interface{}, key string) error {
	if nil == val || "" == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
		}
		return nil
	}

	value, err := util.GetFloat64ByInterface(val)
	if nil != err {
		blog.Error("params should be float, but found [%#v]", val)
		return valid.errif.Errorf(common.CCErrCommParamsNeedFloat, key)
	}

	property, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	floatObjOption := parseMinMaxOption(property.Option)
	if 0 == len(floatObjOption.Min) || 0 == len(floatObjOption.Max) {
		return nil
	}

	maxValue, err := strconv.ParseFloat(floatObjOption.Max, 64)
	if nil != err {
		maxValue = common.MaxFloat64
	}
	minValue, err := strconv.ParseFloat(floatObjOption.Min, 64)
	if nil != err {
		minValue = common.MinFloat64
	}
	if value > maxValue || value < minValue {
		blog.Errorf("params %s:%v not valid", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

// validEnum valid enum
func (valid *Validator) validEnum(val interface{}, key string) error {
	// validate require
	if nil == val || val == "" {
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

//valid date
func (valid *Validator) validDate(val interface{}, key string) error {
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

//valid time
func (valid *Validator) validTime(val interface{}, key string) error {
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
