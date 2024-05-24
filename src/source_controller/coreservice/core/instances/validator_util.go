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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"configcenter/src/common/valid"
	"configcenter/src/storage/driver/mongodb"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FillLostFieldValue fill the value in inst map data
func FillLostFieldValue(ctx context.Context, valData mapstr.MapStr, propertys []metadata.Attribute) error {
	var idRuleField *metadata.Attribute
	for idx, field := range propertys {
		val, ok := valData[field.PropertyID]
		if ok && (field.PropertyType != common.FieldTypeIDRule || val != "") {
			continue
		}

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
		case common.FieldTypeIDRule:
			idRuleField = &propertys[idx]
		default:
			valData[field.PropertyID] = nil
		}
	}

	// 由于id规则字段可能会来自实例的其他字段，所以需要在最后进行填充
	if idRuleField != nil {
		if err := fillLostIDRuleFieldValue(ctx, valData, *idRuleField); err != nil {
			return err
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

	floatOption, err := metadata.ParseFloatOption(field.Option)
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

	intOption, err := metadata.ParseIntOption(field.Option)
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

	listVal, err := metadata.ParseListOption(field.Option)
	if err != nil {
		return err
	}

	defaultVal := util.GetStrByInterface(field.Default)
	for _, value := range listVal {
		if defaultVal == value {
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
	enumOptions, err := metadata.ParseEnumOption(val)
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

func fillLostIDRuleFieldValue(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute) error {
	val, err := GetIDRuleVal(ctx, valData, field)
	if err != nil {
		return err
	}

	valData[field.PropertyID] = val
	return nil
}

// GetIDRuleVal get id rule value
func GetIDRuleVal(ctx context.Context, valData mapstr.MapStr, field metadata.Attribute) (string, error) {
	rid := util.ExtractRequestIDFromContext(ctx)

	rules, err := metadata.ParseSubIDRules(field.Option)
	if err != nil {
		blog.Errorf("parse sub id rule failed, field: %+v, err: %v, rid: %s", field, err, rid)
		return "", err
	}

	var val string
	for _, rule := range rules {
		switch rule.Kind {
		case metadata.Const:
			val += rule.Val

		case metadata.Attr:
			val += util.GetStrByInterface(valData[rule.Val])

		case metadata.GlobalID:
			seqName := metadata.GetIDRule(common.GlobalIDRule)
			id, err := mongodb.Client().NextSequence(ctx, seqName)
			if err != nil {
				blog.Errorf("get next sequence failed, seq name: %s, err: %v, rid: %s", seqName, err, rid)
				return "", err
			}
			idStr, err := metadata.MakeUpDigit(id, rule.Len)
			if err != nil {
				blog.Errorf("make up the id failed, id: %d, len: %d, err: %v, rid: %s", id, rule.Len, err, rid)
				return "", err
			}
			val += idStr

		case metadata.LocalID:
			seqName := metadata.GetIDRule(field.ObjectID)
			id, err := mongodb.Client().NextSequence(ctx, seqName)
			if err != nil {
				blog.Errorf("get next sequence failed, seq name: %s, err: %v, rid: %s", seqName, err, rid)
				return "", err
			}
			idStr, err := metadata.MakeUpDigit(id, rule.Len)
			if err != nil {
				blog.Errorf("make up the id failed, id: %d, len: %d, err: %v, rid: %s", id, rule.Len, err, rid)
				return "", err
			}
			val += idStr

		case metadata.RandomID:
			val += metadata.GetIDRuleRandomID(rule.Len)

		default:
			blog.Errorf("option is invalid, val: %v, rid: %s", field.Option, rid)
			return "", fmt.Errorf("option is invalid, val: %v", field.Option)
		}
	}

	return val, nil
}
