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
	"fmt"
	"regexp"
	"strconv"

	"configcenter/src/common"
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

// parseIntOption  parse int data in option
func parseIntOption(val interface{}) (IntOption, error) {
	intOption := IntOption{}
	if val == nil || val == "" {
		return intOption, fmt.Errorf("int type field option is null")
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
		return intOption, fmt.Errorf("unknow val type: %T", val)
	}

	return intOption, nil
}

// parseFloatOption  parse float data in option
func parseFloatOption(val interface{}) (FloatOption, error) {
	floatOption := FloatOption{}
	if nil == val || "" == val {
		return floatOption, fmt.Errorf("float type field option is null")
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
		return floatOption, fmt.Errorf("unknow val type: %T", val)
	}
	return floatOption, nil
}

// FillLostFieldValue fill the value in inst map data
func FillLostFieldValue(ctx context.Context, valData mapstr.MapStr, propertys []metadata.Attribute) error {
	for _, field := range propertys {
		if _, ok := valData[field.PropertyID]; !ok {
			switch field.PropertyType {
			case common.FieldTypeSingleChar, common.FieldTypeLongChar:
				if err := fillLostStringFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeEnum:
				if err := fillLostEnumFieldValue(ctx, valData, field); err != nil {
					return err
				}
			case common.FieldTypeEnumMulti:
				if err := fillLostEnumMultiFieldValue(ctx, valData, field); err != nil {
					return err
				}
			case common.FieldTypeEnumQuote:
				if err := fillLostEnumQuoteFieldValue(ctx, valData, field); err != nil {
					return err
				}
			case common.FieldTypeDate:
				if err := fillLostDateFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeFloat:
				if err := fillLostFloatFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeInt:
				if err := fillLostIntFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeTime:
				if err := fillLostTimeFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeUser:
				if err := fillLostUserFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeOrganization:
				if err := fillLostOrganizationFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeTimeZone:
				if err := fillLostTimeZoneFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeList:
				if err := fillLostListFieldValue(valData, field); err != nil {
					return err
				}
			case common.FieldTypeBool:
				if err := fillLostBoolFieldValue(valData, field); err != nil {
					return err
				}
			default:
				valData[field.PropertyID] = nil
			}
		}
	}
	return nil
}

func fillLostStringFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = ""
	if field.Default == nil {
		return nil
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		return fmt.Errorf("single char default value not string, value: %v", field.Default)
	}

	if len(defaultVal) == 0 {
		return nil
	}

	option, ok := field.Option.(string)
	if !ok {
		return fmt.Errorf("single char regular verification rules is illegal, value: %v", field.Option)
	}
	if len(option) == 0 {
		return  nil
	}

	match, err := regexp.MatchString(option, defaultVal)
	if err != nil || !match {
		return fmt.Errorf("the current string does not conform to regular verification rules")
	}
	valData[field.PropertyID] = defaultVal
	return nil
}

func fillLostEnumFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	defaultOptions, err := getEnumOption(ctx, field.Option)
	if err != nil {
		return err
	}

	if defaultOptions == nil {
		return nil
	}

	if len(defaultOptions) == 1 {
		valData[field.PropertyID] = defaultOptions[0].ID
		return nil
	}

	return fmt.Errorf("there are multiple default values for enum fields, value: %v", field.Option)
}

func fillLostEnumMultiFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	defaultOptions, err := getEnumOption(ctx, field.Option)
	if err != nil {
		return err
	}

	if defaultOptions == nil {
		return nil
	}

	ids := make([]interface{}, 0)
	for _, k := range defaultOptions {
		ids = append(ids, k.ID)
	}
	if len(ids) == 0 {
		return nil
	}

	valData[field.PropertyID] = ids
	return nil
}

func fillLostEnumQuoteFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	enumQuoteOptions, err := metadata.ParseEnumQuoteOption(ctx, field.Option)
	if err != nil {
		return err
	}
	if len(enumQuoteOptions) == 0 {
		return nil
	}

	instIDs := make([]interface{}, 0)
	for _, k := range enumQuoteOptions {
		instIDs = append(instIDs, k.InstID)
	}
	if len(instIDs) == 0 {
		return nil
	}

	valData[field.PropertyID] = instIDs
	return nil
}

func fillLostDateFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	if ok := util.IsDate(field.Default); !ok {
		return fmt.Errorf("date type field default value format is err, defaultVal: %v", field.Default)
	}

	valData[field.PropertyID] = field.Default
	return nil
}

func fillLostFloatFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	defaultVal, err := util.GetFloat64ByInterface(field.Default)
	if err != nil {
		return err
	}

	intObjOption, err := parseFloatOption(field.Option)
	if err != nil {
		return err
	}
	if len(intObjOption.Min) == 0 || len(intObjOption.Max) == 0 {
		return fmt.Errorf("float type field max or min value is wrong")
	}

	maxValue, err := strconv.ParseFloat(intObjOption.Max, 64)
	if err != nil {
		maxValue = float64(common.MaxInt64)
	}
	minValue, err := strconv.ParseFloat(intObjOption.Min, 64)
	if err != nil {
		minValue = float64(common.MinInt64)
	}

	if defaultVal > maxValue || defaultVal < minValue {
		return fmt.Errorf("float type field default value is illegal, value: %v", field.Default)
	}

	valData[field.PropertyID] = defaultVal
	return nil
}

func fillLostIntFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	defaultVal, err := util.GetInt64ByInterface(field.Default)
	if err != nil {
		return err
	}

	intObjOption, err := parseIntOption(field.Option)
	if err != nil {
		return err
	}
	if len(intObjOption.Min) == 0 || len(intObjOption.Max) == 0 {
		return fmt.Errorf("int type field max or min value is wrong")
	}

	maxValue, err := strconv.ParseInt(intObjOption.Max, 10, 64)
	if err != nil {
		maxValue = common.MaxInt64
	}
	minValue, err := strconv.ParseInt(intObjOption.Min, 10, 64)
	if err != nil {
		minValue = common.MinInt64
	}

	if defaultVal > maxValue || defaultVal < minValue {
		return fmt.Errorf("int type field default value is illegal, value: %v", field.Default)
	}

	valData[field.PropertyID] = defaultVal
	return nil
}

func fillLostTimeFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	if _, ok := util.IsTime(field.Default); !ok {
		return fmt.Errorf("time type field default value format is err, defaultVal: %v", field.Default)
	}

	valData[field.PropertyID] = field.Default
	return nil
}

func fillLostUserFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	defaultVal, ok := field.Default.(string)
	if !ok {
		return fmt.Errorf("user type field default value not string, value: %v", field.Default)
	}

	if ok := util.IsUser(defaultVal); !ok {
		return fmt.Errorf("user type field default value not user type, value: %s", defaultVal)
	}

	valData[field.PropertyID] = defaultVal
	return nil
}

func fillLostOrganizationFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	defaultVal, ok := field.Default.([]interface{})
	if !ok {
		return  fmt.Errorf("organization type field default value not array, type: %T", field.Default)
	}

	for _, orgID := range defaultVal {
		if !util.IsInteger(orgID) {
			return fmt.Errorf("orgID params not int, type: %T", orgID)
		}
	}

	valData[field.PropertyID] = defaultVal
	return nil
}

func fillLostTimeZoneFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	if ok := util.IsTimeZone(field.Default); !ok {
		return fmt.Errorf("the default value of time zone type is not in time zone format %v", field.Default)
	}

	valData[field.PropertyID] = field.Default
	return nil
}

func fillLostListFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = nil
	if field.Default == nil {
		return nil
	}

	arrOption, ok := field.Option.([]interface{})
	if !ok || len(arrOption) == 0 {
		return fmt.Errorf("list type field option is null, option: %v", field.Option)
	}

	defaultVal := util.GetStrByInterface(field.Default)
	for _, value := range arrOption {
		val := util.GetStrByInterface(value)
		if defaultVal == val {
			valData[field.PropertyID] = defaultVal
			return nil
		}
	}

	return fmt.Errorf("list type default value is error, default value: %v", field.Default)
}

func fillLostBoolFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = false
	if field.Default == nil {
		return nil
	}

	if err := util.ValidateBoolType(field.Default); err != nil {
		return err
	}

	valData[field.PropertyID] = field.Default
	return nil
}

func getEnumOption(ctx context.Context, val interface{}) ([]metadata.EnumVal, error) {
	enumOptions, err := metadata.ParseEnumOption(ctx, val)
	if err != nil {
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
