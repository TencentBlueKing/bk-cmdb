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
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/framework/core/option"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

// NewValidMap returns new NewValidMap
func NewValidMap(ownerID, objID string, pheader http.Header, engine *backbone.Engine) *ValidMap {
	return &ValidMap{
		Engine:  engine,
		pheader: pheader,
		ownerID: ownerID,
		objID:   objID,

		propertys:    map[string]metadata.Attribute{},
		require:      map[string]bool{},
		isOnly:       map[string]bool{},
		shouldIgnore: map[string]bool{},
	}
}

// NewValidMapWithKeyFields returns new NewValidMap
func NewValidMapWithKeyFields(ownerID, objID string, ignoreFields []string, pheader http.Header, engine *backbone.Engine) *ValidMap {
	tmp := NewValidMap(ownerID, objID, pheader, engine)

	for _, item := range ignoreFields {
		tmp.shouldIgnore[item] = true
	}
	return tmp
}

// Init init
func (valid *ValidMap) Init() error {
	valid.errif = valid.CCErr.CreateDefaultCCErrorIf(util.GetLanguage(valid.pheader))
	m := map[string]interface{}{
		common.BKObjIDField:   valid.objID,
		common.BKOwnerIDField: valid.ownerID,
	}
	result, err := valid.CoreAPI.ObjectController().Meta().SelectObjectAttWithParams(valid.ctx, valid.pheader, m)
	if nil != err {
		return err
	}
	if !result.Result {
		return valid.errif.Error(result.Code)
	}
	for _, attr := range result.Data {
		if attr.PropertyID == common.BKChildStr || attr.PropertyID == common.BKParentStr {
			continue
		}
		valid.propertys[attr.PropertyID] = attr
		if attr.IsRequired {
			valid.require[attr.PropertyID] = true
		}
		if attr.IsOnly {
			valid.isOnly[attr.PropertyID] = true
		}
	}
	return nil
}

// ValidMap basic valid
func (valid *ValidMap) ValidMap(valData map[string]interface{}, validType string, instID int64) error {
	err := valid.Init()
	if nil != err {
		blog.Errorf("init validator faile %s", err.Error())
		return err
	}

	// valRule.GetObjAttrByID(valid.pheader, valid.objID)
	// valid.IsRequireArr = valRule.IsRequireArr
	// valid.IsOnlyArr = valRule.IsOnlyArr
	// valid.PropertyKv = valRule.PropertyKv
	// keyDataArr := make([]string, 0)
	// var result bool
	// var err error
	// blog.Infof("valid rule:%v \nvalid data:%v", valRule, valData)

	// for key := range valid.KeyFileds {
	// 	// set the key field
	// 	keyDataArr = append(keyDataArr, key)
	// }

	//set default value
	setEnumDefault(valData, valid.propertys)

	//valid create request
	if validType == common.ValidCreate {
		fillLostedFieldValue(valData, valid.propertys)
	}

	for key, val := range valData {

		if valid.shouldIgnore[key] {
			// ignore the key field
			continue
		}

		property, ok := valid.propertys[key]
		if !ok {
			blog.Error("params is not valid, the key is %s", key)
			return valid.errif.Errorf(common.CCErrCommParamsIsInvalid, key)
		}
		fieldType := property.PropertyType
		option := property.Option
		switch fieldType {
		case common.FieldTypeSingleChar:
			err = valid.validChar(val, key)
		case common.FieldTypeLongChar:
			err = valid.validLongChar(val, key)
		case common.FieldTypeInt:
			err = valid.validInt(val, key, option)
		case common.FieldTypeEnum:
			err = valid.validEnum(val, key, option)
		case common.FieldTypeDate:
			err = valid.validDate(val, key)
		case common.FieldTypeTime:
			err = valid.validTime(val, key)
		case common.FieldTypeTimeZone:
			err = valid.validTimeZone(val, key)
		case common.FieldTypeBool:
			err = valid.validBool(val, key)
		default:
			continue
		}
		if nil != err {
			return err
		}
	}

	//fmt.Printf("valdata:%+v\n", valData)
	//valid unique
	if validType == common.ValidCreate {
		err = valid.validCreateUnique(valData)
		return err
	} else {
		err = valid.validUpdateUnique(valData, valid.objID, instID)
		return err
	}

}

//valid char
func (valid *ValidMap) validChar(val interface{}, key string) error {
	if nil == val {
		blog.Error("params in need")
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	switch value := val.(type) {
	case string:
		if len(value) > common.FieldTypeSingleLenChar {
			blog.Errorf("params over length %d", common.FieldTypeSingleLenChar)
			return valid.errif.Errorf(common.CCErrCommOverLimit, key, common.FieldTypeSingleLenChar)
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
func (valid *ValidMap) validLongChar(val interface{}, key string) error {
	if nil == val {
		blog.Error("params in need")
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	switch value := val.(type) {
	case string:
		if len(value) > common.FieldTypeLongLenChar {
			blog.Errorf("params over length %d", common.FieldTypeSingleLenChar)
			return valid.errif.Errorf(common.CCErrCommOverLimit, key, common.FieldTypeSingleLenChar)
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
func (valid *ValidMap) validInt(val interface{}, key string) error {
	var value int64
	if nil == val || "" == val {
		if valid.require[key] {
			blog.Error("params can not be null")
			return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return nil
	}

	var id int
	switch value := val.(type) {
	case int:
		id = value
	case int32:
		id = int(value)
	case int64:
		id = int(value)
	case float64:
		id = int(value)
	case float32:
		id = int(value)
	default:
		blog.Errorf("params %s:%#v not int", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsNeedInt, key)
	}

	option, ok := valid.propertys[key]
	if !ok {
		return nil
	}
	intObjOption := parseIntOption(option)
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return nil
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
		blog.Errorf("params %s:%#v not valid", key, val)
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

//valid char
func (valid *ValidMap) validTimeZone(val interface{}, key string) error {
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

//validBool
func (valid *ValidMap) validBool(val interface{}, key string) error {
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

// validEnum valid enum
func (valid *ValidMap) validEnum(val interface{}, key string, option interface{}) error {
	// validate require
	if nil == val || "" == val {
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
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	return nil
}

//valid date
func (valid *ValidMap) validDate(val interface{}, key string) error {
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
func (valid *ValidMap) validTime(val interface{}, key string) error {
	isIn := util.InArray(key, valid.IsRequireArr)
	if !isIn && nil == val {
		return true, nil
	}

	if isIn && nil == val {
		blog.Error("params in need")
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
	}

	valStr, ok := val.(string)
	if false == ok {
		blog.Error("date can shoule be string")
		return valid.errif.Errorf(common.CCErrCommParamsShouldBeString, key)

	}
	if isIn && 0 == len(valStr) {
		blog.Error("params  can not be empty")
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	result := util.IsTime(valStr)
	if !result {
		blog.Error("params   not valid")
		return valid.errif.Errorf(common.CCErrCommParamsInvalid, key)
	}
	isIn = util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(valStr) {
		blog.Error("params  can not be empty")
		return valid.errif.Errorf(common.CCErrCommParamsNeedSet, key)
	}
	return true, nil

}
