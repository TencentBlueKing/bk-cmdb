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
	"configcenter/src/common/valid"

	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EnumOption enum option
type EnumOption []EnumVal

// IntOption integer option
type IntOption struct {
	Min int64 `bson:"min" json:"min"`
	Max int64 `bson:"max" json:"max"`
}

// FloatOption float option
type FloatOption struct {
	Min float64 `bson:"min" json:"min"`
	Max float64 `bson:"max" json:"max"`
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
	if val == nil || val == "" {
		return IntOption{}, fmt.Errorf("int type field option is null")
	}

	var optMap map[string]interface{}

	switch option := val.(type) {
	case string:
		return parseIntOptionMaxMin(gjson.Get(option, "max").Raw, gjson.Get(option, "min").Raw)
	case map[string]interface{}:
		optMap = option
	case bson.M:
		optMap = option
	case bson.D:
		optMap = option.Map()
	default:
		return IntOption{}, fmt.Errorf("unknow val type: %T", val)
	}

	return parseIntOptionMaxMin(optMap["max"], optMap["min"])
}

func parseIntOptionMaxMin(maxVal, minVal interface{}) (IntOption, error) {
	intOption := IntOption{
		Min: common.MinInt64,
		Max: common.MaxInt64,
	}

	switch max := maxVal.(type) {
	case string:
		if len(max) != 0 && max != `""` {
			maxValue, err := strconv.ParseInt(max, 10, 64)
			if err != nil {
				return IntOption{}, fmt.Errorf("parse max option %+v failed, err: %v", max, err)
			}
			intOption.Max = maxValue
		}
	default:
		maxValue, err := util.GetInt64ByInterface(max)
		if err != nil {
			return IntOption{}, fmt.Errorf("parse max option %+v failed, err: %v", max, err)
		}
		intOption.Max = maxValue
	}

	switch min := minVal.(type) {
	case string:
		if len(min) != 0 && min != `""` {
			minValue, err := strconv.ParseInt(min, 10, 64)
			if err != nil {
				return IntOption{}, fmt.Errorf("parse min option %+v failed, err: %v", min, err)
			}
			intOption.Max = minValue
		}
	default:
		minValue, err := util.GetInt64ByInterface(min)
		if err != nil {
			return IntOption{}, fmt.Errorf("parse min option %+v failed, err: %v", min, err)
		}
		intOption.Min = minValue
	}
	return intOption, nil
}

// parseFloatOption  parse float data in option
func parseFloatOption(val interface{}) (FloatOption, error) {
	if val == nil || val == "" {
		return FloatOption{}, fmt.Errorf("float type field option is null")
	}

	var optMap map[string]interface{}

	switch option := val.(type) {
	case string:
		return parseFloatOptionMaxMin(gjson.Get(option, "max").Raw, gjson.Get(option, "min").Raw)
	case map[string]interface{}:
		optMap = option
	case bson.M:
		optMap = option
	case bson.D:
		optMap = option.Map()
	default:
		return FloatOption{}, fmt.Errorf("unknow val type: %T", val)
	}

	return parseFloatOptionMaxMin(optMap["max"], optMap["min"])
}

func parseFloatOptionMaxMin(maxVal, minVal interface{}) (FloatOption, error) {
	floatOption := FloatOption{
		Min: float64(common.MinInt64),
		Max: float64(common.MaxInt64),
	}

	switch max := maxVal.(type) {
	case string:
		if len(max) != 0 && max != `""` {
			maxValue, err := strconv.ParseFloat(max, 64)
			if err != nil {
				return FloatOption{}, fmt.Errorf("parse max option %+v failed, err: %v", max, err)
			}
			floatOption.Max = maxValue
		}
	default:
		maxValue, err := util.GetFloat64ByInterface(max)
		if err != nil {
			return FloatOption{}, fmt.Errorf("parse max option %+v failed, err: %v", max, err)
		}
		floatOption.Max = maxValue
	}

	switch min := minVal.(type) {
	case string:
		if len(min) != 0 && min != `""` {
			minValue, err := strconv.ParseFloat(min, 64)
			if err != nil {
				return FloatOption{}, fmt.Errorf("parse min option %+v failed, err: %v", min, err)
			}
			floatOption.Max = minValue
		}
	default:
		minValue, err := util.GetFloat64ByInterface(min)
		if err != nil {
			return FloatOption{}, fmt.Errorf("parse min option %+v failed, err: %v", min, err)
		}
		floatOption.Min = minValue
	}
	return floatOption, nil
}

// FillLostFieldValue fill the value in inst map data
func FillLostFieldValue(ctx context.Context, valData mapstr.MapStr, propertys []metadata.Attribute) error {
	var err error
	for _, field := range propertys {
		if _, ok := valData[field.PropertyID]; ok {
			continue
		}

		switch field.PropertyType {
		case common.FieldTypeSingleChar, common.FieldTypeLongChar:
			err = fillLostStringFieldValue(valData, field)
		case common.FieldTypeEnum:
			err = fillLostEnumFieldValue(ctx, valData, field)
		case common.FieldTypeEnumMulti:
			err = fillLostEnumMultiFieldValue(ctx, valData, field)
		case common.FieldTypeEnumQuote:
			err = fillLostEnumQuoteFieldValue(ctx, valData, field)
		case common.FieldTypeDate:
			err = fillLostDateFieldValue(valData, field)
		case common.FieldTypeFloat:
			err = fillLostFloatFieldValue(valData, field)
		case common.FieldTypeInt:
			err = fillLostIntFieldValue(valData, field)
		case common.FieldTypeTime:
			err = fillLostTimeFieldValue(valData, field)
		case common.FieldTypeUser:
			err = fillLostUserFieldValue(valData, field)
		case common.FieldTypeOrganization:
			err = fillLostOrganizationFieldValue(valData, field)
		case common.FieldTypeTimeZone:
			err = fillLostTimeZoneFieldValue(valData, field)
		case common.FieldTypeList:
			err = fillLostListFieldValue(valData, field)
		case common.FieldTypeBool:
			err = fillLostBoolFieldValue(valData, field)
		default:
			valData[field.PropertyID] = nil
		}
	}

	if err != nil {
		return err
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

	// option compatible with the scene where the option is not set in the model attribute.
	option, ok := field.Option.(string)
	if field.Option != nil && !ok {
		return fmt.Errorf("single char regular verification rules is illegal, value: %v", field.Option)
	}
	if len(option) == 0 {
		valData[field.PropertyID] = defaultVal
		return nil
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
		return fmt.Errorf("parse %s default value %+v failed, err: %v", field.PropertyID, field.Default, err)
	}

	floatOption, err := parseFloatOption(field.Option)
	if err != nil {
		return fmt.Errorf("parse %s option %+v failed, err: %v", field.PropertyID, field.Option, err)
	}

	if defaultVal > floatOption.Max || defaultVal < floatOption.Min {
		return fmt.Errorf("%s default value %v is illegal", field.PropertyID, defaultVal)
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
		return fmt.Errorf("parse %s default value %+v failed, err: %v", field.PropertyID, field.Default, err)
	}

	intOption, err := parseIntOption(field.Option)
	if err != nil {
		return fmt.Errorf("parse %s option %+v failed, err: %v", field.PropertyID, field.Option, err)
	}

	if defaultVal > intOption.Max || defaultVal < intOption.Min {
		return fmt.Errorf("%s default value %v is illegal", field.PropertyID, defaultVal)
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

	ok = util.IsUser(defaultVal)
	if defaultVal != "" && !ok {
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

	var orgIDs []interface{}
	switch defaultVal := field.Default.(type) {
	case []interface{}:
		orgIDs = defaultVal
	case primitive.A:
		orgIDs = defaultVal
	default:
		return fmt.Errorf("organization type field default value not array type, propertyID: %s, type: %T",
			field.PropertyID, field.Default)
	}

	for _, orgID := range orgIDs {
		if !util.IsInteger(orgID) {
			return fmt.Errorf("orgID params not int, type: %T", orgID)
		}
	}

	valData[field.PropertyID] = orgIDs
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

	var listVal []interface{}
	switch optionVal := field.Option.(type) {
	case []interface{}:
		listVal = optionVal
	case primitive.A:
		listVal = optionVal
	default:
		return fmt.Errorf("list type field option is not array type, propertyID: %s, type: %T", field.PropertyID,
			field.Option)
	}

	if len(listVal) == 0 {
		return fmt.Errorf("list type field option is empty, propertyID: %s, option: %v", field.PropertyID, field.Option)
	}

	defaultVal := util.GetStrByInterface(field.Default)
	for _, value := range listVal {
		val := util.GetStrByInterface(value)
		if defaultVal == val {
			valData[field.PropertyID] = defaultVal
			return nil
		}
	}

	return fmt.Errorf("list type default value is error, propertyID: %s, default value: %v", field.PropertyID,
		field.Default)
}

func fillLostBoolFieldValue(valData mapstr.MapStr, field metadata.Attribute) error {
	valData[field.PropertyID] = false
	if field.Default == nil {
		return nil
	}

	if err := valid.ValidateBoolType(field.Default); err != nil {
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
