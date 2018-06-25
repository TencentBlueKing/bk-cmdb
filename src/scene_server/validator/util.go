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
	"configcenter/src/common/util"
	api "configcenter/src/source_controller/api/object"
	"encoding/json"

	"gopkg.in/mgo.v2/bson"
)

func getString(val interface{}) string {
	if val == nil {
		return ""
	}
	if ret, ok := val.(string); ok {
		return ret
	}
	return ""
}
func getBool(val interface{}) bool {
	if val == nil {
		return false
	}
	if ret, ok := val.(bool); ok {
		return ret
	}
	return false
}

// fillLostedFieldValue fill the value in inst map data
func fillLostedFieldValue(valData map[string]interface{}, fields []api.ObjAttDes, isRequireArr []string) {
	for _, field := range fields {
		_, ok := valData[field.PropertyID]
		if !ok {
			if util.InStrArr(isRequireArr, field.PropertyID) {
				continue
			}
			switch field.PropertyType {
			case common.FieldTypeSingleChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeLongChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeInt:
				valData[field.PropertyID] = nil
			case common.FieldTypeEnum:
				enumOptions := ParseEnumOption(field.Option)
				v := ""
				if len(enumOptions) > 0 {
					var defaultOption *EnumVal
					for _, k := range enumOptions {
						if k.IsDefault {
							defaultOption = &k
							break
						}
					}
					if nil != defaultOption {
						v = defaultOption.ID
					}
				}
				valData[field.PropertyID] = v
			case common.FieldTypeDate:
				valData[field.PropertyID] = ""
			case common.FieldTypeTime:
				valData[field.PropertyID] = ""
			case common.FieldTypeUser:
				valData[field.PropertyID] = ""
			case common.FieldTypeMultiAsst:
				valData[field.PropertyID] = nil
			case common.FieldTypeTimeZone:
				valData[field.PropertyID] = nil
			case common.FieldTypeBool:
				valData[field.PropertyID] = nil
			default:
				valData[field.PropertyID] = nil
			}
		}
	}
}

// ParseEnumOption convert val to []EnumVal
func ParseEnumOption(val interface{}) []EnumVal {
	enumOptions := []EnumVal{}
	if nil == val || "" == val {
		return enumOptions
	}
	switch options := val.(type) {
	case string:
		json.Unmarshal([]byte(options), &enumOptions)
	case []interface{}:
		for _, optionVal := range options {
			if option, ok := optionVal.(map[string]interface{}); ok {
				enumOption := EnumVal{}
				enumOption.ID = getString(option["id"])
				enumOption.Name = getString(option["name"])
				enumOption.Type = getString(option["type"])
				enumOption.IsDefault = getBool(option["is_default"])
				enumOptions = append(enumOptions, enumOption)
			} else if option, ok := optionVal.(bson.M); ok {
				enumOption := EnumVal{}
				enumOption.ID = getString(option["id"])
				enumOption.Name = getString(option["name"])
				enumOption.Type = getString(option["type"])
				enumOption.IsDefault = getBool(option["is_default"])
				enumOptions = append(enumOptions, enumOption)
			}
		}
	}
	return enumOptions
}

//parseIntOption  parse int data in option
func parseIntOption(val interface{}) IntOption {
	intOption := IntOption{}
	if nil == val || "" == val {
		return intOption
	}
	switch option := val.(type) {
	case string:
		json.Unmarshal([]byte(option), &intOption)
	case map[string]interface{}:
		intOption.Min = getString(option["min"])
		intOption.Max = getString(option["max"])
	}
	return intOption
}

//setEnumDefault
func setEnumDefault(valData map[string]interface{}, valRule *ValRule) {

	for key, val := range valData {
		rule, ok := valRule.FieldRule[key]
		if !ok {
			continue
		}
		fieldType := rule[common.BKPropertyTypeField].(string)
		option := rule[common.BKOptionField]
		switch fieldType {
		case common.FieldTypeEnum:
			if nil != val {
				valStr, ok := val.(string)
				if false == ok {
					return
				}
				if "" != valStr {
					continue
				}
			}

			enumOption := ParseEnumOption(option)
			var defaultOption *EnumVal

			for _, k := range enumOption {
				if k.IsDefault {
					defaultOption = &k
					break
				}
			}
			if nil != defaultOption {
				valData[key] = defaultOption.ID
			}

		}

	}

	return
}
