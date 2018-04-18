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
	"configcenter/src/common/http/httpclient"
	"configcenter/src/common/util"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var innerObject = []string{common.BKInnerObjIDApp, common.BKInnerObjIDSet, common.BKInnerObjIDModule, common.BKInnerObjIDProc, common.BKInnerObjIDHost, common.BKInnerObjIDPlat} //{"app", "set", "module", "process", "host", "plat"}

type IntOption struct {
	Min string `json:min`
	Max string `json:max`
}

type EnumVal struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default"`
}

type ValidMap struct {
	ownerID      string
	objID        string
	objCtrl      string
	IsRequireArr []string
	IsOnlyArr    []string
	KeyFileds    map[string]interface{}
	PropertyKv   map[string]string
	ccError      errors.DefaultCCErrorIf
}

type InstRst struct {
	Result  bool        `json:result`
	Code    int         `json:code`
	Message interface{} `json:message`
	Data    interface{} `json:data`
}

func NewValidMap(ownerID, objID, objCtrl string, err errors.DefaultCCErrorIf) *ValidMap {
	return &ValidMap{ownerID: ownerID, objID: objID, objCtrl: objCtrl, KeyFileds: make(map[string]interface{}, 0), ccError: err}
}

func NewValidMapWithKeyFileds(ownerID, objID, objCtrl string, keyFileds []string, err errors.DefaultCCErrorIf) *ValidMap {
	tmp := &ValidMap{ownerID: ownerID, objID: objID, objCtrl: objCtrl, KeyFileds: make(map[string]interface{}, 0), ccError: err}

	for _, item := range keyFileds {
		tmp.KeyFileds[item] = item
	}
	return tmp
}

//basic valid
func (valid *ValidMap) ValidMap(valData map[string]interface{}, validType string, instID int) (bool, error) {
	valRule := NewValRule(valid.ownerID, valid.objCtrl)

	valRule.GetObjAttrByID(valid.objID)
	valid.IsRequireArr = valRule.IsRequireArr
	valid.IsOnlyArr = valRule.IsOnlyArr
	valid.PropertyKv = valRule.PropertyKv
	keyDataArr := make([]string, 0)
	var result bool
	var err error
	blog.Infof("valid rule:%v \nvalid data:%v", valRule, valData)

	for key := range valid.KeyFileds {
		// set the key filed
		keyDataArr = append(keyDataArr, key)
	}

	//set default value
	valid.setEnumDefault(valData, valRule)
	for key, val := range valData {

		if _, keyOk := valid.KeyFileds[key]; keyOk {
			// ignore the key filed
			continue
		}

		keyDataArr = append(keyDataArr, key)

		rule, ok := valRule.FieldRule[key]
		if !ok {
			blog.Error("params is not valid, the key is %s", key)
			return false, valid.ccError.Errorf(common.CCErrCommParamsIsInvalid, key)
		}

		fieldType := rule[common.BKPropertyTypeField].(string)
		option := rule[common.BKOptionField].(string)
		switch fieldType {
		case common.FiledTypeSingleChar:
			if nil == val {
				blog.Error("params in need")
				return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			result, err = valid.validChar(val, key)
			if 0 != len(option) && result && "" != val {
				//fmt.Println(option)
				strReg := regexp.MustCompile(option)
				strVal := val.(string)
				//fmt.Println(strVal)
				result = strReg.MatchString(strVal)
				if !result {
					err = valid.ccError.Error(common.CCErrFieldRegValidFailed)
				} else {
					err = nil
				}
			}
		case common.FiledTypeLongChar:
			if nil == val {
				blog.Error("params in need")
				return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
			}
			result, err = valid.validLongChar(val, key)
			if 0 != len(option) && result && "" != val {
				//fmt.Println(option)
				strReg := regexp.MustCompile(option)
				strVal := val.(string)

				result = strReg.MatchString(strVal)
				if !result {
					err = valid.ccError.Error(common.CCErrFieldRegValidFailed)
				} else {
					err = nil
				}
			}
		case common.FiledTypeInt:
			result, err = valid.validInt(val, key, option)
		case common.FiledTypeEnum:
			result, err = valid.validEnum(val, key, option)
		case common.FiledTypeDate:
			result, err = valid.validDate(val, key)
		case common.FiledTypeTime:
			result, err = valid.validTime(val, key)
		case common.FieldTypeTimeZone:
			result, err = valid.validTimeZone(val, key)
		case common.FiledTypeBool:
			result, err = valid.validBool(val, key)
		default:
			continue
		}
		if !result {
			return result, err
		}
	}
	//valid create request
	if validType == common.ValidCreate {
		diffArr := util.StrArrDiff(valRule.NoEnumFiledArr, keyDataArr)
		if 0 != len(diffArr) {
			//			var lanDiffArr []string
			//			for _, i := range diffArr {
			//				lanDiffArr = append(lanDiffArr, valid.PropertyKv[i])
			//			}
			keyStr := strings.Join(diffArr, ",")
			blog.Error("params lost filed")
			return false, valid.ccError.Errorf(common.CCErrCommParamsLostField, keyStr)
		}
	}
	//fmt.Printf("valdata:%+v\n", valData)
	//valid unique
	if validType == common.ValidCreate {
		result, err = valid.validCreateUnique(valData)
		return result, err
	} else {
		result, err = valid.validUpdateUnique(valData, valid.objID, instID)
		return result, err
	}

	return true, nil
}

//valid create unique
func (valid *ValidMap) validCreateUnique(valData map[string]interface{}) (bool, error) {
	isInner := false
	objID := valid.objID
	if util.InArray(valid.objID, innerObject) {
		isInner = true
	} else {
		objID = "object"
	}

	if 0 == len(valid.IsOnlyArr) {
		blog.Debug("is only array is zero %+v", valid.IsOnlyArr)
		return true, nil
	}
	searchCond := make(map[string]interface{})
	for key, val := range valData {
		if util.InArray(key, valid.IsOnlyArr) {
			searchCond[key] = val
		}
	}
	if !isInner {
		searchCond[common.BKObjIDField] = valid.objID
	}

	if 0 == len(searchCond) {
		return true, nil
	}
	condition := make(map[string]interface{})
	condition["condition"] = searchCond
	info, _ := json.Marshal(condition)
	httpCli := httpclient.NewHttpClient()
	httpCli.SetHeader("Content-Type", "application/json")
	httpCli.SetHeader("Accept", "application/json")
	blog.Info("get insts by cond: %s", string(info))
	url := fmt.Sprintf("%s/object/v1/insts/%s/search", valid.objCtrl, objID)
	if !strings.HasPrefix(url, "http://") {
		url = fmt.Sprintf("http://%s", url)
	}
	blog.Info("get insts by url : %s", url)
	rst, err := httpCli.POST(url, nil, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return false, err
	}

	var rstRes InstRst
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return false, jserr
	}
	if false == rstRes.Result {
		blog.Error("get rst res error :%v", rstRes)
		return false, valid.ccError.Error(common.CCErrCommUniqueCheckFailed)
	}

	data := rstRes.Data.(map[string]interface{})
	count, err := util.GetIntByInterface(data["count"])
	if nil != err {
		blog.Error("get data error :%v", data)
		return false, valid.ccError.Error(common.CCErrCommParseDataFailed)
	}
	if 0 != count {
		blog.Error("duplicate data ")
		return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
	}
	return true, nil
}

//valid update unique
func (valid *ValidMap) validUpdateUnique(valData map[string]interface{}, objID string, instID int) (bool, error) {
	isInner := false
	urlID := valid.objID
	if util.InArray(valid.objID, innerObject) {
		isInner = true
	} else {
		urlID = "object"
	}

	if 0 == len(valid.IsOnlyArr) {
		return true, nil
	}
	searchCond := make(map[string]interface{})
	for key, val := range valData {
		if util.InArray(key, valid.IsOnlyArr) {
			searchCond[key] = val
		}
	}
	if !isInner {
		searchCond[common.BKObjIDField] = valid.objID
	}
	if 0 == len(searchCond) {
		return true, nil
	}
	condition := make(map[string]interface{})
	condition["condition"] = searchCond
	info, _ := json.Marshal(condition)
	httpCli := httpclient.NewHttpClient()
	httpCli.SetHeader("Content-Type", "application/json")
	httpCli.SetHeader("Accept", "application/json")
	blog.Info("get insts by cond: %s", string(info))
	blog.Info("get insts by cond instID: %v", instID)
	rst, err := httpCli.POST(fmt.Sprintf("%s/object/v1/insts/%s/search", valid.objCtrl, urlID), nil, []byte(info))
	blog.Info("get insts by return: %s", string(rst))
	if nil != err {
		blog.Error("request failed, error:%v", err)
		return false, valid.ccError.Error(common.CCErrCommHTTPDoRequestFailed)
	}

	var rstRes InstRst
	if jserr := json.Unmarshal(rst, &rstRes); nil != jserr {
		blog.Error("can not unmarshal the result , error information is %v", jserr)
		return false, valid.ccError.Error(common.CCErrCommJSONUnmarshalFailed)
	}
	if false == rstRes.Result {
		blog.Error("valid update unique false: %v", rstRes)
		return false, valid.ccError.Error(common.CCErrCommUniqueCheckFailed)
	}
	data := rstRes.Data.(map[string]interface{})
	count, err := util.GetIntByInterface(data["count"])
	if nil != err {
		err := "data false"
		blog.Error("data struct false %v", err)
		return false, valid.ccError.Error(common.CCErrCommParseDataFailed)
	}
	if 0 == count {
		return true, nil
	} else if 1 == count {
		info, ok := data["info"]
		if false == ok {
			blog.Error("data struct false lack info %v", data)
			return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
		}
		infoMap, ok := info.([]interface{})
		if false == ok {
			blog.Error("data struct false lack info is not array%v", data)
			return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
		}
		for _, j := range infoMap {
			i := j.(map[string]interface{})
			objIDName := util.GetObjIDByType(objID)
			instIDc, ok := i[objIDName]
			if false == ok {
				blog.Error("data struct false no objID%v", objIDName)
				return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
			}
			instIDci, err := util.GetIntByInterface(instIDc)

			if nil != err {
				blog.Error("instID not int , error info is %s", err.Error())
				return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
			}
			if instIDci == instID {
				return true, nil
			}
			return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
		}
	} else {
		//err := "duplicate data "
		return false, valid.ccError.Error(common.CCErrCommDuplicateItem)
	}
	return true, nil
}

//valid char
func (valid *ValidMap) validChar(val interface{}, key string) (bool, error) {
	if reflect.TypeOf(val).Kind() != reflect.String {
		blog.Error("params should be  string")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedString, key)
	}
	value := reflect.ValueOf(val).String()
	if len(value) > common.FiledTypeSingleLenChar {
		blog.Errorf("params over length %d", common.FiledTypeSingleLenChar)
		return false, valid.ccError.Errorf(common.CCErrCommOverLimit, key, common.FiledTypeSingleLenChar)
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
		blog.Errorf("params over length %d", common.FiledTypeLongLenChar)
		return false, valid.ccError.Errorf(common.CCErrCommOverLimit, key, common.FiledTypeLongLenChar)
	}
	isIn := util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(value) {
		blog.Error("params can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
	}

	return true, nil
}

//valid int
func (valid *ValidMap) validInt(val interface{}, key string, option string) (bool, error) {
	var value int
	if nil == val {
		isIn := util.InArray(key, valid.IsRequireArr)
		if true == isIn {
			blog.Error("params  can not be null")
			return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)

		}
		return true, nil
	}
	if reflect.TypeOf(val).Kind() == reflect.String {
		valStr := reflect.ValueOf(val).String()
		var re error
		value, re = strconv.Atoi(valStr)
		if nil != re {
			blog.Error("params  not int")
			return false, valid.ccError.Errorf(common.CCErrCommParamsNeedInt, key)
		}
	}
	var intObjOption IntOption
	if reflect.TypeOf(val).Kind() == reflect.Int {
		value2 := reflect.ValueOf(val).Int()
		value = int(value2)

	}
	if 0 == value {
		value, _ = util.GetIntByInterface(val)
	}
	if 0 == len(option) {
		return true, nil
	}
	err := json.Unmarshal([]byte(option), &intObjOption)
	if nil != err {
		return true, nil
	}
	if 0 == len(intObjOption.Min) || 0 == len(intObjOption.Max) {
		return true, nil
	}

	maxValue, err := strconv.Atoi(intObjOption.Max)
	if err != nil {
		return true, nil
	}
	minValue, err := strconv.Atoi(intObjOption.Min)
	if err != nil {
		return true, nil
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

	if reflect.TypeOf(val).Kind() != reflect.Bool {
		blog.Error("params should be  bool")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedBool, key)
	}
	return true, nil
}

//valid enum
func (valid *ValidMap) setEnumDefault(valData map[string]interface{}, valRule *ValRule) error {

	for key, val := range valData {
		rule, ok := valRule.FieldRule[key]
		if !ok {
			continue
		}
		fieldType := rule[common.BKPropertyTypeField].(string)
		option := rule[common.BKOptionField].(string)
		switch fieldType {
		case common.FiledTypeEnum:
			if nil != val {
				valStr, ok := val.(string)
				if false == ok {
					return nil
				}
				if "" != valStr {
					return nil
				}
			}

			var enumOption []EnumVal
			var defaultOption *EnumVal = nil
			re := json.Unmarshal([]byte(option), &enumOption)
			if nil != re {
				blog.Error("params  not valid")
				return valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
			}

			for _, k := range enumOption {
				if k.IsDefault {
					defaultOption = &k
					break
				}
			}
			if nil != defaultOption {
				valData[key] = defaultOption.Name
			}

		}

	}

	return nil
}

//valid enum
func (valid *ValidMap) validEnum(val interface{}, key string, option string) (bool, error) {
	valStr, ok := val.(string)
	if false == ok {
		return true, nil
	}
	var enumOption []EnumVal
	var defaultOption *EnumVal = nil
	re := json.Unmarshal([]byte(option), &enumOption)
	if nil != re {
		blog.Error("params  not valid")
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}
	match := false

	for _, k := range enumOption {
		if k.Name == valStr {
			match = true
			break
		}
		if k.IsDefault {
			defaultOption = &k
		}
	}
	if "" == valStr && nil != defaultOption {
		val = defaultOption.Name
		valStr = defaultOption.Name
	} else if !match {
		blog.Error("params  not valid")
		return false, valid.ccError.Errorf(common.CCErrCommParamsInvalid, key)
	}
	isIn := util.InArray(key, valid.IsRequireArr)
	if isIn && 0 == len(valStr) {
		blog.Error("params  can not be empty")
		return false, valid.ccError.Errorf(common.CCErrCommParamsNeedSet, key)
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
