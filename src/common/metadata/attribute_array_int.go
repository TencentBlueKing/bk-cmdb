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

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/util"
	"configcenter/src/common/valid/attribute/manager/register"
)

func init() {
	register.Register(&arrayInt{})
}

// arrayInt represents an integer array attribute type.
type arrayInt struct{}

// Name returns the name of the arrayDocument attribute.
func (a *arrayInt) Name() string {
	return "array_int"
}

// DisplayName returns the display name for user.
func (a *arrayInt) DisplayName() string {
	return "整数数组"
}

// RealType returns the db type of the document attribute.
func (a *arrayInt) RealType() string {
	return common.FieldTypeLongChar
}

// Info returns the tips for user.
func (a *arrayInt) Info() string {
	return "整数数组"
}

// Validate validates the arrayInt attribute value.
func (a *arrayInt) Validate(ctx context.Context, objID, propertyType string, required bool,
	option, value interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	if value == nil {
		if required {
			blog.Errorf("array_int %s.%s required, rid: %s", objID, propertyType, rid)
			return fmt.Errorf("array_int %s.%s required", objID, propertyType)
		}
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(value)
	if !ok {
		blog.Errorf("array_int %s.%s not []interface{}, got %T, rid: %s",
			objID, propertyType, value, rid)
		return fmt.Errorf("array_int %s.%s must be array", objID, propertyType)
	}

	opts, err := a.parseArrayIntOption(option)
	if err != nil {
		blog.Errorf("array_int parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_int invalid option: %v", err)
	}

	return a.validateIntArray(rid, objID, propertyType, arr, opts)
}

// FillLostValue fills missing values with default value.
func (a *arrayInt) FillLostValue(ctx context.Context, valData mapstr.MapStr,
	propertyID string, defaultValue, option interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	valData[propertyID] = nil
	if defaultValue == nil {
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(defaultValue)
	if !ok {
		blog.Errorf("array_int default not []interface{}, rid: %s", rid)
		return fmt.Errorf("array_int default must be array")
	}

	opts, err := a.parseArrayIntOption(option)
	if err != nil {
		blog.Errorf("array_int parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_int invalid option: %v", err)
	}

	if err := a.validateIntArray(rid, "", "", arr, opts); err != nil {
		return err
	}

	valData[propertyID] = arr
	return nil
}

// ValidateOption validates the option field.
func (a *arrayInt) ValidateOption(ctx context.Context, option, defaultVal interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	opts, err := a.parseArrayIntOption(option)
	if err != nil {
		blog.Errorf("array_int parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_int invalid option: %v", err)
	}

	if opts.Option.Min > opts.Option.Max {
		blog.Errorf("array_int min %d > max %d, rid: %s",
			opts.Option.Min, opts.Option.Max, rid)
		return fmt.Errorf("array_int min must not exceed max")
	}

	if defaultVal == nil {
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(defaultVal)
	if !ok {
		blog.Errorf("array_int default not []interface{}, rid: %s", rid)
		return fmt.Errorf("array_int default must be array")
	}

	return a.validateIntArray(rid, "", "", arr, opts)
}

// validateIntArray validates all integers in array are within range.
func (a *arrayInt) validateIntArray(rid, objID, prop string,
	arr []interface{}, opts ArrayOption[IntOption]) error {

	if opts.Cap < len(arr) {
		return fmt.Errorf("array_int invalid cap %d, rid: %s", opts.Cap, rid)
	}
	for i, v := range arr {
		intVal, err := util.GetInt64ByInterface(v)
		if err != nil {
			if objID != "" {
				blog.Errorf("array_int %s.%s item [%d] not int64, rid: %s",
					objID, prop, i, rid)
				return fmt.Errorf("array_int %s.%s item [%d] not int64", objID, prop, i)
			}
			blog.Errorf("array_int item [%d] not int64, rid: %s", i, rid)
			return fmt.Errorf("array_int item [%d] not int64", i)
		}

		if intVal < opts.Option.Min || intVal > opts.Option.Max {
			if objID != "" {
				blog.Errorf("array_int %s.%s item [%d] %d not in [%d,%d], rid: %s",
					objID, prop, i, intVal, opts.Option.Min, opts.Option.Max, rid)
				return fmt.Errorf("array_int %s.%s item [%d] not in [%d,%d]",
					objID, prop, i, opts.Option.Min, opts.Option.Max)
			}
			blog.Errorf("array_int item [%d] %d not in [%d,%d], rid: %s",
				i, intVal, opts.Option.Min, opts.Option.Max, rid)
			return fmt.Errorf("array_int item [%d] not in [%d,%d]",
				i, opts.Option.Min, opts.Option.Max)
		}
	}
	return nil
}

// parseArrayIntOption parses the option into IntOption.
func (a *arrayInt) parseArrayIntOption(option interface{}) (ArrayOption[IntOption], error) {
	arrayOption, err := ParseArrayOption[IntOption](option, ParseIntOption)
	if err != nil {
		return ArrayOption[IntOption]{}, err
	}
	fmt.Println(arrayOption)
	return arrayOption, nil
}

var _ register.AttributeTypeI = (*arrayInt)(nil)
