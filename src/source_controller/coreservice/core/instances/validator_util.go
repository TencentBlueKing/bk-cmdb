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
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

// EnumOption enum option
type EnumOption []EnumVal

// IntOption integer option
type IntOption struct {
	Min string `bson:"min" json:"min"`
	Max string `bson:"max" json:"max"`
}

// FloatOption float option
type FloatOption struct {
	Min string `bson:"min" json:"min"`
	Max string `bson:"max" json:"max"`
}

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

// GetDefault returns EnumOption's default value
func (opt EnumOption) GetDefault() *EnumVal {
	for index := range opt {
		if opt[index].IsDefault {
			return &opt[index]
		}
	}
	return nil
}

// EnumVal enum option val
type EnumVal struct {
	ID        string `bson:"id"           json:"id"`
	Name      string `bson:"name"         json:"name"`
	Type      string `bson:"type"         json:"type"`
	IsDefault bool   `bson:"is_default"   json:"is_default"`
}

// ParseEnumOption convert val to []EnumVal
func ParseEnumOption(ctx context.Context, val interface{}) (EnumOption, error) {
	rid := util.ExtractRequestIDFromContext(ctx)
	enumOptions := []EnumVal{}
	if nil == val || "" == val {
		return enumOptions, nil
	}
	switch options := val.(type) {
	case []EnumVal:
		return options, nil
	case string:
		err := json.Unmarshal([]byte(options), &enumOptions)
		if nil != err {
			blog.Errorf("ParseEnumOption error : %s, rid: %s", err.Error(), rid)
			return nil, err
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
			} else {
				return nil, fmt.Errorf("unknow val type: %#v", val)
			}
		}
	case bson.A:
		for _, optionVal := range options {
			if option, ok := optionVal.(map[string]interface{}); ok {
				enumOption := EnumVal{}
				enumOption.ID = getString(option["id"])
				enumOption.Name = getString(option["name"])
				enumOption.Type = getString(option["type"])
				enumOption.IsDefault = getBool(option["is_default"])
				enumOptions = append(enumOptions, enumOption)
			} else if option, ok := optionVal.(bson.D); ok {
				opt := option.Map()
				enumOption := EnumVal{}
				enumOption.ID = getString(opt["id"])
				enumOption.Name = getString(opt["name"])
				enumOption.Type = getString(opt["type"])
				enumOption.IsDefault = getBool(opt["is_default"])
				enumOptions = append(enumOptions, enumOption)
			} else {
				return nil, fmt.Errorf("unknow val type: %#v", val)
			}
		}
	default:
		return nil, fmt.Errorf("unknow val type: %#v", val)
	}
	return enumOptions, nil
}

// parseIntOption  parse int data in option
func parseIntOption(ctx context.Context, val interface{}) IntOption {
	rid := util.ExtractRequestIDFromContext(ctx)
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
	case bson.M:
		intOption.Min = getString(option["min"])
		intOption.Max = getString(option["max"])
	case bson.D:
		opt := option.Map()
		intOption.Min = getString(opt["min"])
		intOption.Max = getString(opt["max"])
	default:
		blog.Warnf("unknow val type: %#v, rid: %s", val, rid)
	}
	return intOption
}

// parseFloatOption  parse float data in option
func parseFloatOption(ctx context.Context, val interface{}) FloatOption {
	rid := util.ExtractRequestIDFromContext(ctx)
	floatOption := FloatOption{}
	if nil == val || "" == val {
		return floatOption
	}
	switch option := val.(type) {
	case string:
		floatOption.Min = gjson.Get(option, "min").Raw
		floatOption.Max = gjson.Get(option, "max").Raw
	case map[string]interface{}:
		floatOption.Min = getString(option["min"])
		floatOption.Max = getString(option["max"])
	case bson.M:
		floatOption.Min = getString(option["min"])
		floatOption.Max = getString(option["max"])
	case bson.D:
		opt := option.Map()
		floatOption.Min = getString(opt["min"])
		floatOption.Max = getString(opt["max"])
	default:
		blog.Warnf("unknow val type: %#v, rid: %s", val, rid)
	}
	return floatOption
}

// FillLostedFieldValue fill the value in inst map data
func FillLostedFieldValue(ctx context.Context, valData mapstr.MapStr, propertys []metadata.Attribute) {
	rid := util.ExtractRequestIDFromContext(ctx)
	for _, field := range propertys {
		if field.PropertyID == common.BKChildStr || field.PropertyID == common.BKParentStr {
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
				enumOptions, err := metadata.ParseEnumOption(ctx, field.Option)
				if err != nil {
					blog.Warnf("ParseEnumOption failed: %v, rid: %s", err, rid)
					valData[field.PropertyID] = nil
					continue
				}
				if len(enumOptions) > 0 {
					var defaultOption *metadata.EnumVal
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

func isEmpty(value interface{}) bool {
	return value == nil || value == ""
}
