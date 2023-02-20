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
		if _, ok := valData[field.PropertyID]; !ok {
			switch field.PropertyType {
			case common.FieldTypeSingleChar, common.FieldTypeLongChar:
				fillLostedStringFieldValue(valData, field, rid)
			case common.FieldTypeEnum:
				fillLostedEnumFieldValue(ctx, valData, field, rid)
			case common.FieldTypeEnumMulti:
				fillLostedEnumMultiFieldValue(ctx, valData, field, rid)
			case common.FieldTypeEnumQuote:
				fillLostedEnumQuoteFieldValue(ctx, valData, field, rid)
			case common.FieldTypeDate:
				fillLostedDateFieldValue(valData, field, rid)
			case common.FieldTypeFloat:
				fillLostedFloatFieldValue(valData, field, rid)
			case common.FieldTypeInt:
				fillLostedIntFieldValue(valData, field, rid)
			case common.FieldTypeTime:
				fillLostedTimeFieldValue(valData, field, rid)
			case common.FieldTypeUser:
				fillLostedUserFieldValue(valData, field, rid)
			case common.FieldTypeOrganization:
				fillLostedOrganizationFieldValue(valData, field, rid)
			case common.FieldTypeTimeZone:
				fillLostedTimeZoneFieldValue(valData, field, rid)
			case common.FieldTypeList:
				fillLostedListFieldValue(valData, field, rid)
			case common.FieldTypeBool:
				fillLostedBoolFieldValue(valData, field, rid)
			default:
				valData[field.PropertyID] = nil
			}
		}
	}
}

func fillLostedStringFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = ""
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		blog.Errorf("single char default value not string, value: %+v, rid: %s", field.Default, rid)
		return
	}
	valData[field.PropertyID] = defaultVal
}

func fillLostedEnumFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute, rid string) {
	defaultOptions, err := getEnumOption(ctx, field.Option, rid)
	if err != nil {
		blog.Errorf("get enum option failed, err: %v, rid: %s", err, rid)
		valData[field.PropertyID] = nil
		return
	}

	if defaultOptions == nil {
		valData[field.PropertyID] = nil
		return
	}

	if len(defaultOptions) == 1 {
		valData[field.PropertyID] = defaultOptions[0].ID
		return
	}
	valData[field.PropertyID] = nil
}

func fillLostedEnumMultiFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute, rid string) {
	defaultOptions, err := getEnumOption(ctx, field.Option, rid)
	if err != nil {
		blog.Errorf("get enum option failed, err: %v, rid: %s", err, rid)
		valData[field.PropertyID] = nil
		return
	}

	if defaultOptions == nil {
		valData[field.PropertyID] = nil
		return
	}

	ids := make([]string, 0)
	for _, k := range defaultOptions {
		ids = append(ids, k.ID)
	}
	if len(ids) == 0 {
		valData[field.PropertyID] = nil
		return
	}

	valData[field.PropertyID] = ids
}

func fillLostedEnumQuoteFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute, rid string) {
	enumQuoteOptions, err := metadata.ParseEnumQuoteOption(ctx, field.Option)
	if err != nil {
		blog.Errorf("parse enum quote option failed, err: %v, rid: %s", err, rid)
		valData[field.PropertyID] = nil
		return
	}
	if len(enumQuoteOptions) == 0 {
		valData[field.PropertyID] = nil
		return
	}
	instIDs := make([]int64, 0)
	for _, k := range enumQuoteOptions {
		instIDs = append(instIDs, k.InstID)
	}
	if len(instIDs) == 0 {
		valData[field.PropertyID] = nil
		return
	}

	valData[field.PropertyID] = instIDs
}

func fillLostedDateFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		blog.Errorf("date type field default value not string, value: %+v, rid: %s", field.Default, rid)
		return
	}

	if ok := util.IsDate(defaultVal); !ok {
		blog.Errorf("date type field default value format is err, defaultVal:%+v, rid: %s", defaultVal, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedFloatFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, err := util.GetFloat64ByInterface(field.Default)
	if err != nil {
		blog.Errorf("float type field default value is not number, value: %+v, rid: %s", field.Default, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedIntFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, err := util.GetInt64ByInterface(field.Default)
	if err != nil {
		blog.Errorf("int type field default value is not number, value: %+v, rid: %s", field.Default, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedTimeFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		blog.Errorf("time type field default value not string, value: %+v, rid: %s", field.Default, rid)
		return
	}

	if _, ok := util.IsTime(defaultVal); !ok {
		blog.Errorf("time type field default value format is err, defaultVal: %+v, rid: %s", defaultVal, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedUserFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		blog.Errorf("user type field default value not string, value: %+v, rid: %s", field.Default, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedOrganizationFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.([]interface{})
	if !ok {
		blog.Errorf("organization type field default value not array, val: %+v, rid: %s", field.Default, rid)
		return
	}

	for _, orgID := range defaultVal {
		if !util.IsNumeric(orgID) {
			blog.Errorf("orgID params not int, type: %T, rid: %s", orgID, rid)
			return
		}
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedTimeZoneFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		blog.Errorf("time zone type field default value not string, value: %+v, rid: %s", field.Default, rid)
		return
	}

	if ok := util.IsTimeZone(defaultVal); !ok {
		blog.Errorf("time zone type field default value format is err, defaultVal: %+v, rid: %s", defaultVal, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func fillLostedListFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return
	}

	arrOption, ok := field.Option.([]interface{})
	if !ok || len(arrOption) == 0 {
		blog.Errorf("list type field default value not array, val: %+v, rid: %s", field.Default, rid)
		return
	}

	defaultVal := util.GetStrByInterface(field.Default)
	for _, value := range arrOption {
		val := util.GetStrByInterface(value)
		if defaultVal == val {
			valData[field.PropertyID] = defaultVal
			return
		}
	}
	valData[field.PropertyID] = nil
}

func fillLostedBoolFieldValue(valData mapstr.MapStr, field metadata.Attribute, rid string) {
	valData[field.PropertyID] = false
	if field.Default == nil {
		return
	}

	defaultVal, ok := field.Default.(bool)
	if !ok {
		blog.Errorf("bool type field default value not bool, val: %+v, rid: %s", field.Default, rid)
		return
	}

	valData[field.PropertyID] = defaultVal
}

func getEnumOption(ctx context.Context, val interface{}, rid string) ([]metadata.EnumVal, error) {
	enumOptions, err := metadata.ParseEnumOption(ctx, val)
	if err != nil {
		blog.Errorf("parse enum option failed, err: %v, rid: %s", err, rid)
		return nil, err
	}

	if len(enumOptions) == 0 {
		return nil, nil
	}

	defaultOptions := make([]metadata.EnumVal, 0)
	for _, k := range enumOptions {
		if k.IsDefault {
			defaultOptions = append(defaultOptions, k)
		}
	}

	if len(defaultOptions) == 0 {
		return nil, nil
	}
	return defaultOptions, nil
}
