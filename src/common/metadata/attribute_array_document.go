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
	register.Register(&arrayDocument{})
}

// arrayDocument represents a Document array attribute type.
type arrayDocument struct{}

// Name returns the name of the arrayDocument attribute.
func (a *arrayDocument) Name() string {
	return "array_document"
}

// DisplayName returns the display name for user.
func (a *arrayDocument) DisplayName() string {
	return "附件数组"
}

// RealType returns the db type of the document attribute.
func (a *arrayDocument) RealType() string {
	return common.FieldTypeLongChar
}

// Info returns the tips for user.
func (a *arrayDocument) Info() string {
	return "附件数组"
}

// Validate validates the arrayDocument attribute value.
func (a *arrayDocument) Validate(ctx context.Context, objID, propertyType string, required bool,
	option, value interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	if value == nil {
		if required {
			blog.Errorf("array_document %s.%s required, rid: %s", objID, propertyType, rid)
			return fmt.Errorf("array_document %s.%s required", objID, propertyType)
		}
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(value)
	if !ok {
		blog.Errorf("array_document %s.%s not []interface{}, got %T, rid: %s",
			objID, propertyType, value, rid)
		return fmt.Errorf("array_document %s.%s must be array", objID, propertyType)
	}

	opts, err := a.parseArrayDocumentOption(option)
	if err != nil {
		blog.Errorf("array_document parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_document invalid option: %v", err)
	}

	return a.validateDocumentArray(rid, objID, propertyType, arr, opts)
}

// FillLostValue fills missing values with default value.
func (a *arrayDocument) FillLostValue(ctx context.Context, valData mapstr.MapStr,
	propertyID string, defaultValue, option interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	valData[propertyID] = nil
	if defaultValue == nil {
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(defaultValue)
	if !ok {
		blog.Errorf("array_document default not []interface{}, rid: %s", rid)
		return fmt.Errorf("array_document default must be array")
	}

	opts, err := a.parseArrayDocumentOption(option)
	if err != nil {
		blog.Errorf("array_document parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_document invalid option: %v", err)
	}

	if err := a.validateDocumentArray(rid, "", "", arr, opts); err != nil {
		return err
	}

	valData[propertyID] = arr
	return nil
}

// ValidateOption validates the option field.
func (a *arrayDocument) ValidateOption(ctx context.Context, option, defaultVal interface{}) error {

	rid := util.ExtractRequestIDFromContext(ctx)

	opts, err := a.parseArrayDocumentOption(option)
	if err != nil {
		blog.Errorf("array_document parse option failed: %v, rid: %s", err, rid)
		return fmt.Errorf("array_document invalid option: %v", err)
	}
	err = document{}.ValidateOption(ctx, opts.Option, nil)
	if err != nil {
		return err
	}
	if defaultVal == nil {
		return nil
	}

	arr, ok := util.ConvertAnyToSlice(defaultVal)
	if !ok {
		blog.Errorf("array_document default not []interface{}, rid: %s", rid)
		return fmt.Errorf("array_document default must be array")
	}

	return a.validateDocumentArray(rid, "", "", arr, opts)
}

// validateDocumentArray validates all Documents in array are within range.
func (a *arrayDocument) validateDocumentArray(rid, objID, prop string,
	arr []interface{}, opts ArrayOption[DocumentOption]) error {

	if opts.Cap < len(arr) {
		return fmt.Errorf("array_document invalid cap %d, rid: %s", opts.Cap, rid)
	}
	for _, v := range arr {
		_, err := document{}.ParseValue(v)
		if err != nil {
			return err
		}
	}
	return nil
}

// parseArrayDocumentOption parses the option into DocumentOption.
func (a *arrayDocument) parseArrayDocumentOption(option interface{}) (ArrayOption[DocumentOption], error) {
	arrayOption, err := ParseArrayOption[DocumentOption](option, document{}.ParseOption)
	if err != nil {
		return ArrayOption[DocumentOption]{}, err
	}
	return arrayOption, nil
}

var _ register.AttributeTypeI = (*arrayDocument)(nil)
