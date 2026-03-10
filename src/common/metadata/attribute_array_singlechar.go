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

package metadata

import (
	"context"
	"fmt"
	"regexp"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/valid/attribute/manager/register"
)

func init() {
	// Register the arraySinglechar attribute type
	register.Register(arraySinglechar{})
}

type arraySinglechar struct {
}

// Name returns the name of the arraySinglechar attribute.
func (a arraySinglechar) Name() string {
	return "array_singlechar"
}

// DisplayName returns the display name for user.
func (a arraySinglechar) DisplayName() string {
	return "短字符数组"
}

// RealType returns the db type of the arraySinglechar attribute.
// Flattened array uses LongChar as storage type
func (a arraySinglechar) RealType() string {
	return common.FieldTypeLongChar
}

// Info returns the tips for user.
func (a arraySinglechar) Info() string {
	return "短字符数组"
}

// Validate validates the arraySinglechar attribute value
func (a arraySinglechar) Validate(ctx context.Context, objID string, propertyType string, required bool,
	option, value interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	if value == nil {
		if required {
			blog.Errorf("array_singlechar attribute %s.%s value is required but got nil, rid: %s",
				objID, propertyType, rid)
			return fmt.Errorf("array_singlechar attribute %s.%s value is required but got nil",
				objID, propertyType)
		}
		return nil
	}

	// Validate that value is a slice of any
	strArray, ok := util.ConvertAnyToSlice(value)
	if !ok {
		blog.Errorf("array_singlechar attribute %s.%s value must be []interface{}, got %T, rid: %s",
			objID, propertyType, value, rid)
		return fmt.Errorf("array_singlechar attribute %s.%s value must be []interface{}, got %T",
			objID, propertyType, value)
	}

	// Parse option for regex pattern
	regex := common.FieldTypeSingleCharRegexp
	arrayOpt, err := ParseArrayOption[string](option, func(v any) (string, error) {
		if v == nil {
			return "", nil
		}
		s, ok := v.(string)
		if !ok {
			return s, fmt.Errorf("invalid type %T for option %s, rid: %s", v, option, rid)
		}
		return s, nil
	})
	if err != nil {
		return err
	}
	if arrayOpt.Cap < len(strArray) {
		return fmt.Errorf("array_singlechar invalid cap %d, rid: %s",
			arrayOpt.Cap, rid)
	}
	regex = arrayOpt.Option
	// Compile regex pattern
	pattern, err := regexp.Compile(regex)
	if err != nil {
		blog.Errorf("array_singlechar invalid regex pattern %s, err: %v, rid: %s", regex, err, rid)
		return fmt.Errorf("array_singlechar invalid regex pattern: %v", err)
	}

	// Validate each item in the array
	for i, item := range strArray {
		strVal, ok := item.(string)
		if !ok {
			blog.Errorf("array_singlechar attribute %s.%s array item [%d] type %T is not string, rid: %s",
				objID, propertyType, i, item, rid)
			return fmt.Errorf("array_singlechar attribute %s.%s array item [%d] type %T is not string",
				objID, propertyType, i, item)
		}

		// Validate length
		if len(strVal) > common.FieldTypeSingleLenChar {
			blog.Errorf("array_singlechar attribute %s.%s array item [%d] length %d exceeds max %d, rid: %s",
				objID, propertyType, i, len(strVal), common.FieldTypeSingleLenChar, rid)
			return fmt.Errorf("array_singlechar attribute %s.%s array item [%d] length exceeds max %d",
				objID, propertyType, i, common.FieldTypeSingleLenChar)
		}

		// Validate regex
		if !pattern.MatchString(strVal) {
			blog.Errorf("array_singlechar attribute %s.%s array item [%d] value '%s' does not match regex, rid: %s",
				objID, propertyType, i, strVal, rid)
			return fmt.Errorf("array_singlechar attribute %s.%s array item [%d] does not match regex pattern",
				objID, propertyType, i)
		}
	}

	return nil
}

// FillLostValue fills the lost value with default value
func (a arraySinglechar) FillLostValue(ctx context.Context, valData mapstr.MapStr, propertyId string,
	defaultValue, option interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	valData[propertyId] = nil
	if defaultValue == nil {
		return nil
	}

	// Validate default value
	strArray, ok := util.ConvertAnyToSlice(defaultValue)
	if !ok {
		blog.Errorf("array_singlechar default value must be []interface{}, got %T, rid: %s", defaultValue, rid)
		return fmt.Errorf("array_singlechar default value must be []interface{}, got %T", defaultValue)
	}

	// Parse option for regex pattern
	regex := common.FieldTypeSingleCharRegexp

	arrayOpt, err := ParseArrayOption[string](option, func(v any) (string, error) {
		if v == nil {
			return "", nil
		}
		s, ok := v.(string)
		if !ok {
			return s, fmt.Errorf("invalid type %T for option %s, rid: %s", v, option, rid)
		}
		return s, nil
	})
	if err != nil {
		return err
	}
	if arrayOpt.Cap < len(strArray) {
		return fmt.Errorf("array_singlechar invalid cap %d, rid: %s",
			arrayOpt.Cap, rid)
	}

	// Compile regex pattern
	pattern, err := regexp.Compile(regex)
	if err != nil {
		blog.Errorf("array_singlechar invalid regex pattern %s, err: %v, rid: %s", regex, err, rid)
		return fmt.Errorf("array_singlechar invalid regex pattern: %v", err)
	}

	// Validate each item in default array
	for i, item := range strArray {
		strVal, ok := item.(string)
		if !ok {
			blog.Errorf("array_singlechar default value array item [%d] type %T is not string, rid: %s", i, item, rid)
			return fmt.Errorf("array_singlechar default value array item [%d] type %T is not string", i, item)
		}

		if len(strVal) > common.FieldTypeSingleLenChar {
			blog.Errorf("array_singlechar default value array item [%d] length %d exceeds max %d, rid: %s",
				i, len(strVal), common.FieldTypeSingleLenChar, rid)
			return fmt.Errorf("array_singlechar default value array item [%d] length exceeds max %d",
				i, common.FieldTypeSingleLenChar)
		}

		if !pattern.MatchString(strVal) {
			blog.Errorf("array_singlechar default value array item [%d] value '%s' does not match regex, rid: %s",
				i, strVal, rid)
			return fmt.Errorf("array_singlechar default value array item [%d] does not match regex pattern", i)
		}
	}

	valData[propertyId] = strArray
	return nil
}

// ValidateOption validates the option field
func (a arraySinglechar) ValidateOption(ctx context.Context, option, defaultVal interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	arrayOpt, err := ParseArrayOption[string](option, func(v any) (string, error) {
		if v == nil {
			return "", nil
		}
		s, ok := v.(string)
		if !ok {
			return s, fmt.Errorf("invalid type %T for option %s, rid: %s", v, option, rid)
		}
		return s, nil
	})
	if err != nil {
		return err
	}

	if defaultVal == nil {
		return nil
	}

	// Validate default value
	strArray, ok := util.ConvertAnyToSlice(defaultVal)
	if !ok {
		blog.Errorf("array_singlechar default value must be []interface{}, got %T, rid: %s", defaultVal, rid)
		return fmt.Errorf("array_singlechar default value must be []interface{}, got %T", defaultVal)
	}

	// Get regex pattern
	regex := arrayOpt.Option

	pattern, err := regexp.Compile(regex)
	if err != nil {
		blog.Errorf("array_singlechar invalid regex pattern %s, err: %v, rid: %s", regex, err, rid)
		return fmt.Errorf("array_singlechar invalid regex pattern: %v", err)
	}

	// Validate each item in default array
	for i, item := range strArray {
		strVal, ok := item.(string)
		if !ok {
			blog.Errorf("array_singlechar default value array item [%d] type %T is not string, rid: %s", i, item, rid)
			return fmt.Errorf("array_singlechar default value array item [%d] type %T is not string", i, item)
		}

		if len(strVal) > common.FieldTypeSingleLenChar {
			blog.Errorf("array_singlechar default value array item [%d] length exceeds max %d, rid: %s",
				i, common.FieldTypeSingleLenChar, rid)
			return fmt.Errorf("array_singlechar default value array item [%d] length exceeds max %d",
				i, common.FieldTypeSingleLenChar)
		}

		if !pattern.MatchString(strVal) {
			blog.Errorf("array_singlechar default value array item [%d] does not match regex, rid: %s", i, rid)
			return fmt.Errorf("array_singlechar default value array item [%d] does not match regex pattern", i)
		}
	}

	return nil
}

var _ register.AttributeTypeI = &arraySinglechar{}
