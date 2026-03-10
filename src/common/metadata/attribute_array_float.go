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
	register.Register(&arrayFloat{})
}

// arrayFloat represents a float array attribute type.
type arrayFloat struct{}

// Name returns the name of the arrayDocument attribute.
func (a *arrayFloat) Name() string {
	return "array_float"
}

// DisplayName returns the display name for user.
func (a *arrayFloat) DisplayName() string {
	return "浮点数数组"
}

// RealType returns the db type of the document attribute.
func (a *arrayFloat) RealType() string {
	return common.FieldTypeLongChar
}

// Info returns the tips for user.
func (a *arrayFloat) Info() string {
	return "浮点数数组"
}

// Validate validates the arrayFloat attribute value.
func (a *arrayFloat) Validate(ctx context.Context, objID, propertyType string, required bool,
	option, value interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	if value == nil {
		if required {
			blog.Errorf("array_float %s.%s required, rid: %s", objID, propertyType, rid)
			return fmt.Errorf("array_float %s.%s required", objID, propertyType)
		}
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(value)
	if !ok {
		blog.Errorf("array_float %s.%s not []interface{}, got %T, rid: %s",
			objID, propertyType, value, rid)
		return fmt.Errorf("array_float %s.%s must be array", objID, propertyType)
	}

	opts, err := a.parseArrayFloatOption(option)
	if err != nil {
		blog.Errorf("array_float parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_float invalid option: %v", err)
	}

	return a.validateFloatArray(rid, objID, propertyType, arr, opts)
}

// FillLostValue fills missing values with default value.
func (a *arrayFloat) FillLostValue(ctx context.Context, valData mapstr.MapStr,
	propertyID string, defaultValue, option interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	valData[propertyID] = nil
	if defaultValue == nil {
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(defaultValue)
	if !ok {
		blog.Errorf("array_float default not []interface{}, rid: %s", rid)
		return fmt.Errorf("array_float default must be array")
	}

	opts, err := a.parseArrayFloatOption(option)
	if err != nil {
		blog.Errorf("array_float parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_float invalid option: %v", err)
	}

	if err := a.validateFloatArray(rid, "", "", arr, opts); err != nil {
		return err
	}

	valData[propertyID] = arr
	return nil
}

// ValidateOption validates the option field.
func (a *arrayFloat) ValidateOption(ctx context.Context, option, defaultVal interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	opts, err := a.parseArrayFloatOption(option)
	if err != nil {
		blog.Errorf("array_float parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_float invalid option: %v", err)
	}

	if opts.Option.Min > opts.Option.Max {
		blog.Errorf("array_float min %f > max %f, rid: %s",
			opts.Option.Min, opts.Option.Max, rid)
		return fmt.Errorf("array_float min must not exceed max")
	}

	if defaultVal == nil {
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(defaultVal)
	if !ok {
		blog.Errorf("array_float default not []interface{}, rid: %s", rid)
		return fmt.Errorf("array_float default must be array")
	}

	return a.validateFloatArray(rid, "", "", arr, opts)
}

// validateFloatArray validates all floats in array are within range.
func (a *arrayFloat) validateFloatArray(rid, objID, prop string,
	arr []interface{}, opts ArrayOption[FloatOption]) error {

	if opts.Cap < len(arr) {
		return fmt.Errorf("array_float invalid cap %d, rid: %s", opts.Cap, rid)
	}
	for i, v := range arr {
		floatVal, err := util.GetFloat64ByInterface(v)
		if err != nil {
			if objID != "" {
				blog.Errorf("array_float %s.%s item [%d] not float64, rid: %s",
					objID, prop, i, rid)
				return fmt.Errorf("array_float %s.%s item [%d] not float64", objID, prop, i)
			}
			blog.Errorf("array_float item [%d] not float64, rid: %s", i, rid)
			return fmt.Errorf("array_float item [%d] not float64", i)
		}
		if floatVal < opts.Option.Min || floatVal > opts.Option.Max {
			if objID != "" {
				blog.Errorf("array_float %s.%s item [%d] %f not in [%f,%f], rid: %s",
					objID, prop, i, floatVal, opts.Option.Min, opts.Option.Max, rid)
				return fmt.Errorf("array_float %s.%s item [%d] not in [%f,%f]",
					objID, prop, i, opts.Option.Min, opts.Option.Max)
			}
			blog.Errorf("array_float item [%d] %f not in [%f,%f], rid: %s",
				i, floatVal, opts.Option.Min, opts.Option.Max, rid)
			return fmt.Errorf("array_float item [%d] not in [%f,%f]",
				i, opts.Option.Min, opts.Option.Max)
		}
	}
	return nil
}

// parseArrayFloatOption parses the option into FloatOption.
func (a *arrayFloat) parseArrayFloatOption(option interface{}) (ArrayOption[FloatOption], error) {
	arrayOption, err := ParseArrayOption[FloatOption](option, ParseFloatOption)
	if err != nil {
		return ArrayOption[FloatOption]{}, err
	}
	return arrayOption, nil
}

var _ register.AttributeTypeI = (*arrayFloat)(nil)
