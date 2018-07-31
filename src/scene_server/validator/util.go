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
	"encoding/json"

	"github.com/tidwall/gjson"
	"gopkg.in/mgo.v2/bson"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
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

// FillLostedFieldValue fill the value in inst map data
func FillLostedFieldValue(valData map[string]interface{}, propertys []metadata.Attribute, ignorefields []string) {
	ignores := map[string]bool{}
	for _, field := range ignorefields {
		ignores[field] = true
	}
	for _, field := range propertys {
		if field.PropertyID == common.BKChildStr || field.PropertyID == common.BKParentStr {
			continue
		}
		if ignores[field.PropertyID] {
			continue
		}
		_, ok := valData[field.PropertyID]
		if !ok {
			switch field.PropertyType {
			case common.FieldTypeSingleChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeLongChar:
				valData[field.PropertyID] = ""
			case common.FieldTypeInt:
				valData[field.PropertyID] = nil
			case common.FieldTypeEnum:
				enumOptions := ParseEnumOption(field.Option)
				if len(enumOptions) > 0 {
					var defaultOption *EnumVal
					for _, k := range enumOptions {
						if k.IsDefault {
							defaultOption = &k
							break
						}
					}
					if nil != defaultOption {
						valData[field.PropertyID] = defaultOption.ID
					} else {
						valData[field.PropertyID] = nil
					}
				} else {
					valData[field.PropertyID] = nil
				}
			case common.FieldTypeDate:
				valData[field.PropertyID] = nil
			case common.FieldTypeTime:
				valData[field.PropertyID] = nil
			case common.FieldTypeUser:
				valData[field.PropertyID] = nil
			case common.FieldTypeMultiAsst:
				valData[field.PropertyID] = nil
			case common.FieldTypeTimeZone:
				valData[field.PropertyID] = nil
			case common.FieldTypeBool:
				valData[field.PropertyID] = false
			default:
				valData[field.PropertyID] = nil
			}
		}
	}
}

// ParseEnumOption convert val to []EnumVal
func ParseEnumOption(val interface{}) EnumOption {
	enumOptions := []EnumVal{}
	if nil == val || "" == val {
		return enumOptions
	}
	switch options := val.(type) {
	case []EnumVal:
		return options
	case string:
		err := json.Unmarshal([]byte(options), &enumOptions)
		if nil != err {
			blog.Errorf("ParseEnumOption error : %s", err.Error())
		}
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

		intOption.Min = gjson.Get(option, "min").Raw
		intOption.Max = gjson.Get(option, "max").Raw

	case map[string]interface{}:
		intOption.Min = getString(option["min"])
		intOption.Max = getString(option["max"])
	}
	return intOption
}
