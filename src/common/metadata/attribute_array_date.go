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
	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/valid/attribute/manager/register"
	"context"
	"fmt"
)

func init() {
	// Register the arrayDate attribute type
	register.Register(arrayDate{})
}

type arrayDate struct {
}

// Name returns the name of the arrayDate attribute.
func (a arrayDate) Name() string {
	return "array_date"
}

// DisplayName returns the display name for user.
func (a arrayDate) DisplayName() string {
	return "日期数组"
}

// RealType returns the db type of the arrayDate attribute.
// Flattened array uses LongChar as storage type
func (a arrayDate) RealType() string {
	return common.FieldTypeLongChar
}

// Info returns the tips for user.
func (a arrayDate) Info() string {
	return "日期数组，格式为YYYY-MM-DD"
}

// Validate validates the arrayDate attribute value
func (a arrayDate) Validate(ctx context.Context, objID string, propertyType string, required bool,
	option interface{}, value interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	if value == nil {
		if required {
			blog.Errorf("array_date attribute %s.%s value is required but got nil, rid: %s",
				objID, propertyType, rid)
			return fmt.Errorf("array_date attribute %s.%s value is required but got nil",
				objID, propertyType)
		}
		return nil
	}

	// Validate that value is a slice of any
	dateArray, ok := util.ConvertAnyToSlice(value)
	if !ok {
		blog.Errorf("array_date attribute %s.%s value must be []interface{}, got %T, rid: %s",
			objID, propertyType, value, rid)
		return fmt.Errorf("array_date attribute %s.%s value must be []interface{}, got %T",
			objID, propertyType, value)
	}

	opts, err := ParseArrayOption[string](option, nil)
	if err != nil {
		blog.Errorf("array_date parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_date invalid option: %v", err)
	}
	if opts.Cap < len(dateArray) {
		return fmt.Errorf("array_date invalid cap %d, rid: %s", opts.Cap, rid)
	}
	// Validate each item in the array
	for i, item := range dateArray {
		// Validate date format
		if !util.IsDate(item) {
			blog.Errorf("array_date attribute %s.%s array item [%d] type %T is not a valid date, rid: %s",
				objID, propertyType, i, item, rid)
			return fmt.Errorf("array_date attribute %s.%s array item [%d] is not a valid date",
				objID, propertyType, i)
		}
	}

	return nil
}

// FillLostValue fills the lost value with default value
func (a arrayDate) FillLostValue(ctx context.Context, valData mapstr.MapStr, propertyId string,
	defaultValue, option interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	valData[propertyId] = nil
	if defaultValue == nil {
		return nil
	}

	// Validate default value
	defaultArray, ok := util.ConvertAnyToSlice(defaultValue)
	if !ok {
		blog.Errorf("array_date default value must be []interface{}, got %T, rid: %s", defaultValue, rid)
		return fmt.Errorf("array_date default value must be []interface{}, got %T", defaultValue)
	}

	// Validate each item in default array
	for i, item := range defaultArray {
		if !util.IsDate(item) {
			blog.Errorf("array_date default value array item [%d] type %T is not a valid date, rid: %s", i, item, rid)
			return fmt.Errorf("array_date default value array item [%d] is not a valid date", i)
		}
	}

	valData[propertyId] = defaultArray
	return nil
}

// ValidateOption validates the option field
func (a arrayDate) ValidateOption(ctx context.Context, option interface{}, defaultVal interface{}) error {
	rid := util.ExtractRequestIDFromContext(ctx)

	_, err := ParseArrayOption[string](option, nil)
	if err != nil {
		return err
	}
	if defaultVal == nil {
		return nil
	}

	// Validate default value
	defaultArray, ok := util.ConvertAnyToSlice(defaultVal)
	if !ok {
		blog.Errorf("array_date default value must be []interface{}, got %T, rid: %s", defaultVal, rid)
		return fmt.Errorf("array_date default value must be []interface{}, got %T", defaultVal)
	}

	// Validate each item in default array
	for i, item := range defaultArray {
		if !util.IsDate(item) {
			blog.Errorf("array_date default value array item [%d] type %T is not a valid date, rid: %s", i, item, rid)
			return fmt.Errorf("array_date default value array item [%v] is not a valid date", item)
		}
	}

	return nil
}

var _ register.AttributeTypeI = &arrayDate{}
