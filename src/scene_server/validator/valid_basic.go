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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/errors"
	"configcenter/src/common/util"
	api "configcenter/src/source_controller/api/object"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// NewValidMap returns new NewValidMap
func NewValidMap(ownerID, objID, objCtrl string, forward *api.ForwardParam, err errors.DefaultCCErrorIf) *ValidMap {
	return &ValidMap{ownerID: ownerID, objID: objID, objCtrl: objCtrl, KeyFileds: make(map[string]interface{}, 0), ccError: err, forward: forward}
}

// NewValidMapWithKeyFields returns new NewValidMap
func NewValidMapWithKeyFields(ownerID, objID, objCtrl string, keyFileds []string, forward *api.ForwardParam, err errors.DefaultCCErrorIf) *ValidMap {
	tmp := &ValidMap{ownerID: ownerID, objID: objID, objCtrl: objCtrl, KeyFileds: make(map[string]interface{}, 0), ccError: err, forward: forward}

	for _, item := range keyFileds {
		tmp.KeyFileds[item] = item
	}
	return tmp
}

// ValidMap basic valid
func (valid *ValidMap) ValidMap(valData map[string]interface{}, validType string, instID int) (bool, error) {
	valRule := NewValRule(valid.ownerID, valid.objCtrl)

	valRule.GetObjAttrByID(valid.forward, valid.objID)
	valid.IsRequireArr = valRule.IsRequireArr
	valid.IsOnlyArr = valRule.IsOnlyArr
	valid.PropertyKv = valRule.PropertyKv
	keyDataArr := make([]string, 0)
	var result bool
	var err error
	blog.Infof("valid rule:%v \nvalid data:%v", valRule, valData)

	for key := range valid.KeyFileds {
		// set the key field
		keyDataArr = append(keyDataArr, key)
	}

	//set default value
	setEnumDefault(valData, valRule)

	//valid create request
	if validType == common.ValidCreate {
		fillLostedFieldValue(valData, valRule.AllFieldAttDes, valid.IsRequireArr)
	}

	for key, val := range valData {

		if _, keyOk := valid.KeyFileds[key]; keyOk {
			// ignore the key field
			continue
		}

		keyDataArr = append(keyDataArr, key)

		rule, ok := valRule.FieldRule[key]
		if !ok {
			blog.Error("params is not valid, the key is %s", key)
			return false, valid.ccError.Errorf(common.CCErrCommParamsIsInvalid, key)
		}

		fieldType := rule[common.BKPropertyTypeField].(string)
		option := rule[common.BKOptionField]
		switch fieldType {
		case common.FieldTypeSingleChar:
			if nil == val {
				blog.Error("params in need")
				return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			result, err = valid.validChar(val, key)
			if option != nil && result && "" != val {
				//fmt.Println(option)
				strReg := regexp.MustCompile(option.(string))
				strVal := val.(string)
				//fmt.Println(strVal)
				result = strReg.MatchString(strVal)
				if !result {
					err = valid.ccError.Error(common.CCErrFieldRegValidFailed)
				} else {
					err = nil
				}
			}
		case common.FieldTypeLongChar:
			if nil == val {
				blog.Error("params in need")
				return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			result, err = valid.validLongChar(val, key)
			if option != nil && result && "" != val {
				//fmt.Println(option)
				strReg := regexp.MustCompile(option.(string))
				strVal := val.(string)

				result = strReg.MatchString(strVal)
				if !result {
					err = valid.ccError.Error(common.CCErrFieldRegValidFailed)
				} else {
					err = nil
				}
			}
		case common.FieldTypeInt:
			result, err = valid.validInt(val, key, option)
		case common.FieldTypeEnum:
			result, err = valid.validEnum(val, key, option)
		case common.FieldTypeDate:
			result, err = valid.validDate(val, key)
		case common.FieldTypeTime:
			result, err = valid.validTime(val, key)
		case common.FieldTypeTimeZone:
			result, err = valid.validTimeZone(val, key)
		case common.FieldTypeBool:
			result, err = valid.validBool(val, key)
		default:
			continue
		}
		if !result {
			return result, err
		}
	}

	if validType == common.ValidCreate {
		diffArr := util.StrArrDiff(valRule.NoEnumFiledArr, keyDataArr)
		if 0 != len(diffArr) {
			keyStr := strings.Join(diffArr, ",")
			blog.Error("params lost filed")
			return false, valid.ccError.Errorf(common.CCErrCommParamsLostField, keyStr)
		}
	}

	if validType == common.ValidCreate {
		result, err = valid.validCreateUnique(valData)
		return result, err
	} else {
		result, err = valid.validUpdateUnique(valData, valid.objID, instID)
		return result, err
	}

}

//valid char
func (valid *ValidMap) validChar(val interface{}, key string) (bool, error) {
	if reflect.TypeOf(val).Kind() != reflect.String {
		blog.Error("params should be  string")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedString, key)
	}
	value := reflect.ValueOf(val).String()
	if len(value) > common.FieldTypeSingleLenChar {
		blog.Errorf("params over length %d", common.FieldTypeSingleLenChar)
		return false, valid.ccError.Errorf(common.CCErrCommOverLimit, key, common.FieldTypeSingleLenChar)
	}
	isIn := util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(value) {
		blog.Error("params can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	return true, nil
}

//valid long char
func (valid *ValidMap) validLongChar(val interface{}, key string) (bool, error) {
	if reflect.TypeOf(val).Kind() != reflect.String {
		blog.Error("params should be string")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedString, key)
	}
	value := reflect.ValueOf(val).String()
	if len(value) > 512 {
		blog.Errorf("params over length %d", common.FieldTypeLongLenChar)
		return false, valid.ccError.Errorf(common.CCErrCommOverLimit, key, common.FieldTypeLongLenChar)
	}
	isIn := util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(value) {
		blog.Error("params can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}

	return true, nil
}

// validInt valid int
func (valid *ValidMap) validInt(val interface{}, key string, option interface{}) (bool, error) {
	var value int64
	if nil == val || "" == val {
		isIn := util.InArray(key, valid.IsRequireArr)
		if true == isIn {
			blog.Error("params can not be null")
			return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return true, nil
	}

	// validate type
	value, err := strconv.ParseInt(fmt.Sprint(val), 10, 64)
	if err != nil {
		blog.Error("params not int")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	// validate by option
	if nil == option || "" == option {
		return true, nil
	}
	intObjOption := parseIntOption(option)
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return true, nil
	}

	maxValue, err := strconv.ParseInt(intObjOption.Max, 10, 64)
	if err != nil {
		maxValue = common.MaxInt64
	}
	minValue, err := strconv.ParseInt(intObjOption.Min, 10, 64)
	if err != nil {
		minValue = common.MinInt64
	}
	if value > maxValue || value < minValue {
		blog.Error("params  not valid")
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return true, nil
}

//valid char
func (valid *ValidMap) validTimeZone(val interface{}, key string) (bool, error) {

	isIn := util.InArray(key, valid.IsRequireArr)
	if isIn && nil == val {
		blog.Error("params can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	if nil == val {
		return true, nil
	}
	if reflect.TypeOf(val).Kind() != reflect.String {
		blog.Error("params should be  timezone")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedTimeZone, key)
	}
	value := reflect.ValueOf(val).String()
	isMatch := util.IsTimeZone(value)
	if false == isMatch {
		blog.Error("params should be  timezone")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedTimeZone, key)
	}
	return true, nil
}

//validBool
func (valid *ValidMap) validBool(val interface{}, key string) (bool, error) {

	isIn := util.InArray(key, valid.IsRequireArr)
	if isIn && nil == val {
		blog.Error("params can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	if nil == val {
		return true, nil
	}

	if reflect.TypeOf(val).Kind() != reflect.Bool {
		blog.Error("params should be  bool")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedBool, key)
	}
	return true, nil
}

// validEnum valid enum
func (valid *ValidMap) validEnum(val interface{}, key string, option interface{}) (bool, error) {
	// validate require
	if nil == val || "" == val {
		if util.InArray(key, valid.IsRequireArr) {
			blog.Error("params %s can not be empty", key)
			return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
		}
		return true, nil
	}

	// validate type
	valStr, ok := val.(string)
	if !ok {
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}

	// validate within enum
	enumOption := ParseEnumOption(option)
	match := false
	for _, k := range enumOption {
		if k.ID == valStr {
			match = true
			break
		}
	}
	if !match {
		blog.Error("params %s not valid, option %#v, raw option %#v, value: %#v", key, enumOption, option, val)
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return true, nil
}

//valid date
func (valid *ValidMap) validDate(val interface{}, key string) (bool, error) {
	isIn := util.InArray(key, valid.IsRequireArr)
	if !isIn && nil == val {
		return true, nil
	}
	if isIn && nil == val {
		blog.Error("params in need")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	valStr, ok := val.(string)
	if false == ok {
		blog.Error("date can shoule be string")
		return false, valid.ccError.Errorf(common.CCErrCommParamsShouldBeString, key)

	}
	if isIn && 0 == len(valStr) {
		blog.Error("date params  can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	result := util.IsDate(valStr)
	if !result {
		blog.Error("params  is not valid")
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}
	isIn = util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(valStr) {
		blog.Error("params  can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	return true, nil
}

//valid time
func (valid *ValidMap) validTime(val interface{}, key string) (bool, error) {
	isIn := util.InArray(key, valid.IsRequireArr)
	if !isIn && nil == val {
		return true, nil
	}

	if isIn && nil == val {
		blog.Error("params in need")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}

	valStr, ok := val.(string)
	if false == ok {
		blog.Error("date can shoule be string")
		return false, valid.ccError.Errorf(common.CCErrCommParamsShouldBeString, key)

	}
	if isIn && 0 == len(valStr) {
		blog.Error("params  can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	result := util.IsTime(valStr)
	if !result {
		blog.Error("params   not valid")
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}
	isIn = util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(valStr) {
		blog.Error("params  can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	return true, nil

}
