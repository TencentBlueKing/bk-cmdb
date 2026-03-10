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
	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
	"context"
	"fmt"
	"testing"

	"configcenter/src/common"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestAttributeValidate_Int(t *testing.T) {

	tests := []struct {
		name    string
		attr    metadata.Attribute
		value   interface{}
		wantErr bool
	}{
		{
			name: "int success",
			attr: metadata.Attribute{
				PropertyType: common.FieldTypeInt,
				Option: map[string]interface{}{
					"min": 1,
					"max": 10,
				},
				IsRequired: true,
			},
			value:   5,
			wantErr: false,
		},
		{
			name: "int out of range",
			attr: metadata.Attribute{
				PropertyType: common.FieldTypeInt,
				Option: map[string]interface{}{
					"min": 1,
					"max": 10,
				},
			},
			value:   20,
			wantErr: true,
		},
		{
			name: "int type error",
			attr: metadata.Attribute{
				PropertyType: common.FieldTypeInt,
			},
			value:   "abc",
			wantErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			err := tt.attr.Validate(context.Background(), tt.value, "test")

			if (err.ErrCode != 0) != tt.wantErr {
				t.Fatalf("expect err %v got %v", tt.wantErr, err)
			}

		})
	}
}

func TestAttributeValidate_Float(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeFloat,
		Option: map[string]interface{}{
			"min": 1.0,
			"max": 10.0,
		},
	}

	if err := attr.Validate(context.Background(), 5.5, "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := attr.Validate(context.Background(), 20.1, "test"); err.ErrCode == 0 {
		t.Fatalf("expect error")
	}
}

func TestAttributeValidate_Bool(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeBool,
	}

	if err := attr.Validate(context.Background(), true, "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error")
	}

	if err := attr.Validate(context.Background(), "true", "test"); err.ErrCode == 0 {
		t.Fatalf("expect error")
	}
}

func TestAttributeValidate_SingleChar(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeSingleChar,
		IsRequired:   true,
	}

	if err := attr.Validate(context.Background(), "hello", "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error")
	}

	if err := attr.Validate(context.Background(), 123, "test"); err.ErrCode == 0 {
		t.Fatalf("expect error")
	}
}

func TestAttributeValidate_LongChar(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeLongChar,
	}

	if err := attr.Validate(context.Background(), "long text", "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error")
	}
}

func TestAttributeValidate_Date(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeDate,
	}

	if err := attr.Validate(context.Background(), "2024-01-01", "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error")
	}

	if err := attr.Validate(context.Background(), "invalid-date", "test"); err.ErrCode == 0 {
		t.Fatalf("expect error")
	}
}

func TestAttributeValidate_Time(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeTime,
	}

	if err := attr.Validate(context.Background(), "2024-01-01 10:00:00", "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error")
	}
}

func TestAttributeValidate_Enum(t *testing.T) {

	attr := metadata.Attribute{
		PropertyType: common.FieldTypeEnum,
		Option: []interface{}{
			map[string]interface{}{
				"id":   "1",
				"name": "test1",
				"type": "text",
			},
			map[string]interface{}{
				"id":   "2",
				"name": "test2",
				"type": "text",
			},
		},
	}

	if err := attr.Validate(context.Background(), "1", "test"); err.ErrCode != 0 {
		t.Fatalf("unexpected error")
	}

	if err := attr.Validate(context.Background(), "3", "test"); err.ErrCode == 0 {
		t.Fatalf("expect error")
	}
}

func TestName(t *testing.T) {
	var a any = primitive.A{1}
	fmt.Println(util.ConvertAnyToSlice(a))
	fmt.Println(util.ConvertToInterfaceSlice(a))
	fmt.Println(util.ConvertToInterfaceSlice(1)) //conv to single item
}
