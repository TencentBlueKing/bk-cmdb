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
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestArrayIntOptionParse(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name:  "normal-int-array",
			value: []any{1, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
				"option": map[string]interface{}{
					"min": 1,
					"max": 10,
				},
			},
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 0,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{1},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{1, 2, 3, 4, 5, 6},
			option: map[string]interface{}{
				"len": 1,
				"cap": 5,
			},
			expectErr: true,
		},
		{
			name:  "element-min-violation",
			value: []any{0, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
				"option": map[string]interface{}{
					"min": 1,
					"max": 10,
				},
			},
			expectErr: true,
		},
		{
			name:  "element-max-violation",
			value: []any{1, 2, 20},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
				"option": map[string]interface{}{
					"min": 1,
					"max": 10,
				},
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{1, json.Number("2"), 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: false, // util.GetInt64ByInterface
		},
		{
			name:  "not-array",
			value: 123,
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{1, 2, 3},
			expectErr: true,
		}, {
			name:  "int-option-missing",
			value: []any{1, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "len-equal-boundary",
			value: []any{1},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{1, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_int",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArrayFloatOptionParse(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name:  "normal-int-array",
			value: []any{1.3, 2.1, 3.2},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
				"option": map[string]interface{}{
					"min": 1.1,
					"max": 10,
				},
			},
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 0,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{1.3},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{1.3, 2, 3, 4, 5, 6},
			option: map[string]interface{}{
				"len": 1,
				"cap": 5,
			},
			expectErr: true,
		},
		{
			name:  "element-min-violation",
			value: []any{0, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
				"option": map[string]interface{}{
					"min": 1.5,
					"max": 10,
				},
			},
			expectErr: true,
		},
		{
			name:  "element-max-violation",
			value: []any{1.3, 2, 20},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
				"option": map[string]interface{}{
					"min": 1,
					"max": 10.5,
				},
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{1.3, "2", 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: false, // util.GetFloat64ByInterface
		},
		{
			name:  "not-array",
			value: 123,
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{1.3, 2, 3},
			expectErr: true,
		}, {
			name:  "float-option-missing",
			value: []any{1.3, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "len-equal-boundary",
			value: []any{1.3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{1.3, 2, 3},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_float",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArrayBoolOptionParse(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name:  "normal-bool-array",
			value: []any{true, false, true},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{true},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{true, false, true},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:      "element-type-error",
			value:     []any{1, 2, true},
			expectErr: true,
		},
		{
			name:      "not-array",
			value:     true,
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{true, false, true},
			expectErr: true,
		}, {
			name:  "bool-option-missing",
			value: []any{true, false, true},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "len-equal-boundary",
			value: []any{true},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{true, false, true},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_bool",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArraySingleCharOptionParse(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name:  "normal-singlechar-array",
			value: []any{"abc", "中文"},
			option: map[string]interface{}{
				"len":    2,
				"cap":    10,
				"option": ".+",
			},
		}, {
			name:  "charLength-singlechar-array",
			value: []any{strings.Repeat("中文", (256/len("中文"))+1)},
			option: map[string]interface{}{
				"len":    1,
				"cap":    10,
				"option": ".+",
			},
			expectErr: true,
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{"中文"},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{"abc", "中文", "!@#$%^&*()_+"},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{1, 2, true},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "not-array",
			value: true,
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{"abc", "中文"},
			expectErr: true,
		}, {
			name:  "singlechar-option-missing",
			value: []any{"abc", "中文"},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "len-equal-boundary",
			value: []any{"abc", "中文"},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{"abc", "中文", "!@#$%^&*()_+"},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_singlechar",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArrayLongCharOptionParse(t *testing.T) {

	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name:  "normal-longchar-array",
			value: []any{"abc", "中文"},
			option: map[string]interface{}{
				"len":    2,
				"cap":    10,
				"option": ".+",
			},
		},
		{
			name:  "charLength-longchar-array",
			value: []any{strings.Repeat("中文", (2000/len("中文"))+1)},
			option: map[string]interface{}{
				"len":    1,
				"cap":    10,
				"option": ".+",
			},
			expectErr: true,
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{"中文"},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{"abc", "中文", "!@#$%^&*()_+"},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{1, 2, true},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "not-array",
			value: true,
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{"abc", "中文"},
			expectErr: true,
		}, {
			name:  "longchar-option-missing",
			value: []any{"abc", "中文"},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
		{
			name:  "len-equal-boundary",
			value: []any{"abc", "中文"},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{"abc", "中文", "!@#$%^&*()_+"},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_longchar",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArrayDateOptionParse(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name:  "normal-date-array",
			value: []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			option: map[string]interface{}{
				"len":    2,
				"cap":    10,
				"option": ".+",
			},
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			option: map[string]interface{}{
				"len": 3,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 1,
			},
			expectErr: true,
		},
		{
			name:  "element-format-error",
			value: []any{now.Format(time.DateOnly), now.Format(time.RFC3339)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{now.Format(time.DateOnly), now.Unix()},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "not-array",
			value: now.Format(time.DateOnly),
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			expectErr: true,
		},
		{
			name:  "date-option-missing",
			value: []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "len-equal-boundary",
			value: []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{now.Format(time.DateOnly), now.Add(20 * time.Hour).Format(time.DateOnly)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_date",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArrayTimeOptionParse(t *testing.T) {
	now := time.Now()
	format := "2006-01-02T15:04:05+08:00"
	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{

		{
			name: "normal-time-array",
			value: []any{now.Format(time.DateTime),
				now.Format(format),
				now.Format("2006-01-02T15:04:05.1+08:00"),
				now.Format("2006-01-02T15:04:05.12+08:00"),
				now.Format("2006-01-02T15:04:05.123+08:00"),
				//now.Format("2006-01-02T15:04:05.123+08:00:00"),//err
				//now.Format("2006-01-02T15:04:05.123-07:00"), //err
				//now.Format("2006-01-02T15:04:05Z"), //err
			},
			option: map[string]interface{}{
				"len":    2,
				"cap":    10,
				"option": ".+",
			},
		},
		{
			name:  "empty-array-valid",
			value: []any{},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "len-not-enough",
			value: []any{now.Format(time.DateTime), now.Format(format)},
			option: map[string]interface{}{
				"len": 3,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "exceed-cap",
			value: []any{now.Format(time.DateTime), now.Format(format)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 1,
			},
			expectErr: true,
		},
		{
			name:  "element-format-error",
			value: []any{now.Format(time.DateTime), now.Format(time.RFC3339)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{now.Format(time.DateTime), now.Unix()},
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:  "not-array",
			value: now.Format(time.DateTime),
			option: map[string]interface{}{
				"len": 1,
				"cap": 2,
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{now.Format(time.DateTime), now.Format(format)},
			expectErr: true,
		},
		{
			name:  "time-option-missing",
			value: []any{now.Format(time.DateTime), now.Format(format)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "len-equal-boundary",
			value: []any{now.Format(time.DateTime), now.Format(format)},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{now.Format(time.DateTime), now.Format(format)},
			option: map[string]interface{}{
				"len": 1,
				"cap": 3,
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_time",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestArrayDocumentOptionParse(t *testing.T) {
	tests := []struct {
		name      string
		option    map[string]interface{}
		value     any
		expectErr bool
	}{
		{
			name:  "normal-document-array",
			value: []any{metadata.DocumentValueDocument{}},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
		},
		{
			name:  "empty-array-valid",
			value: []any{},

			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
		},
		{
			name:  "len-not-enough",
			value: []any{metadata.DocumentValueDocument{}},

			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
			expectErr: false,
		},
		{
			name: "exceed-cap",
			value: []any{metadata.DocumentValueDocument{},
				metadata.DocumentValueDocument{},
				metadata.DocumentValueDocument{},
				metadata.DocumentValueDocument{}},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
			expectErr: true,
		},
		{
			name:  "element-type-error",
			value: []any{metadata.DocumentValueDocument{}, 1},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
			expectErr: true,
		},
		{
			name:  "not-array",
			value: metadata.DocumentValueDocument{},
			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
			expectErr: true,
		},
		{
			name:      "option-missing",
			value:     []any{metadata.DocumentValueDocument{}},
			expectErr: true,
		},
		{
			name:  "document-option-missing",
			value: []any{metadata.DocumentValueDocument{}},

			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
			},
			expectErr: false,
		},
		{
			name:  "len-equal-boundary",
			value: []any{metadata.DocumentValueDocument{}},

			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
		},
		{
			name:  "cap-equal-boundary",
			value: []any{metadata.DocumentValueDocument{}},

			option: map[string]interface{}{
				"len": 2,
				"cap": 3,
				"option": metadata.DocumentOption{
					AllowSuffixes: nil,
					AllowSize:     1024,
					Regex:         "",
					Type:          "image",
				},
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			attr := metadata.Attribute{
				PropertyType: "array_document",
				Option:       tt.option,
			}
			err := attr.Validate(context.Background(), tt.value, "test")

			if tt.expectErr {
				if err.ErrCode == 0 {
					t.Fatalf("expect error:%v", err)
				}
				return
			}

			if err.ErrCode != 0 {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
