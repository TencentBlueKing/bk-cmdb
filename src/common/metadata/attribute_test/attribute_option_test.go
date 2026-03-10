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
	"configcenter/src/common/metadata"
	"context"
	"testing"
)

func TestIntOptionParse(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		expectErr bool
	}{
		{
			name: "normal",
			option: map[string]interface{}{
				"min": 1,
				"max": 10,
			},
		},
		{
			name: "min only",
			option: map[string]interface{}{
				"min": 1,
			}, expectErr: true,
		},
		{
			name: "max only",
			option: map[string]interface{}{
				"max": 10,
			}, expectErr: true,
		},
		{
			name: "type error",
			option: map[string]interface{}{
				"min": "1",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: common.FieldTypeInt,
				Option:       tt.option,
			}

			err := attr.Validate(context.Background(), 5, "test")

			if (err.ErrCode != 0) != tt.expectErr {
				t.Fatalf("expect err=%v got=%v", tt.expectErr, err)
			}

		})
	}
}

func TestFloatOption(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		expectErr bool
	}{
		{
			name: "normal",
			option: map[string]interface{}{
				"min": 1.4,
				"max": 10.4,
			},
		},
		{
			name: "min only",
			option: map[string]interface{}{
				"min": 1.3,
			}, expectErr: true,
		},
		{
			name: "max only",
			option: map[string]interface{}{
				"max": 10.5,
			}, expectErr: true,
		},
		{
			name: "type error",
			option: map[string]interface{}{
				"min": "1.4",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: common.FieldTypeFloat,
				Option:       tt.option,
			}

			err := attr.Validate(context.Background(), 5, "test")

			if (err.ErrCode != 0) != tt.expectErr {
				t.Fatalf("expect err=%v got=%v", tt.expectErr, err)
			}

		})
	}
}

func TestDocumentOption(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		expectErr bool
	}{
		{
			name: "normal",
			option: map[string]interface{}{
				"allow_suffixes": nil,
				"allow_size":     1024,
				"regex":          ".+",
				"type":           "image",
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "document",
				Option:       tt.option,
			}

			err := attr.Validate(context.Background(), metadata.DocumentValueDocument{
				Value: "ok",
				Name:  "ok",
			}, "test")

			if (err.ErrCode != 0) != tt.expectErr {
				t.Fatalf("expect err=%v got=%v", tt.expectErr, err)
			}

		})
	}
}
