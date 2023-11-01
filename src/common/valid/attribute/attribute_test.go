/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package attrvalid

import (
	"fmt"
	"net/http"
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/metadata"
)

type errif struct {
}

func (ei errif) CreateDefaultCCErrorIf(language string) errors.DefaultCCErrorIf {
	return defErrIf{}
}

func (ei errif) Error(language string, errCode int) error {
	return errors.New(errCode, fmt.Sprintf("%d", errCode))
}

func (ei errif) Errorf(language string, errCode int, args ...interface{}) error {
	return errors.New(errCode, fmt.Sprintf("%v", args))
}

func (ei errif) Load(res map[string]errors.ErrorCode) {
	return
}

func (ei errif) New(errCode int, msg string) error {
	return errors.New(errCode, msg)
}

type defErrIf struct {
}

func (d defErrIf) Error(errCode int) error {
	return errors.New(errCode, fmt.Sprintf("%d", errCode))
}

func (d defErrIf) Errorf(errCode int, args ...interface{}) error {
	return errors.New(errCode, fmt.Sprintf("%v", args))
}

func (d defErrIf) CCError(errCode int) errors.CCErrorCoder {
	return errors.New(errCode, fmt.Sprintf("%d", errCode))
}

func (d defErrIf) CCErrorf(errCode int, args ...interface{}) errors.CCErrorCoder {
	return errors.New(errCode, fmt.Sprintf("%v", args))
}

func (d defErrIf) New(errorCode int, msg string) error {
	return errors.New(errorCode, msg)
}

func TestValidPropertyOption(t *testing.T) {
	type args struct {
		propertyType string
		option       interface{}
		isMultiple   bool
		defaultVal   interface{}
	}

	kit := rest.NewKitFromHeader(http.Header{}, errif{})

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"enum", args{common.FieldTypeEnum, `[{"id":"a","name":"a","type":"text","is_default":true}]`, false, "a"},
			false},
		{"enum_multi", args{common.FieldTypeEnumMulti, metadata.EnumOption{{
			ID:        "a",
			Name:      "a",
			Type:      "text",
			IsDefault: false,
		}, {
			ID:        "b",
			Name:      "b",
			Type:      "text",
			IsDefault: true,
		}, {
			ID:        "c",
			Name:      "c",
			Type:      "text",
			IsDefault: true,
		}}, true, []string{"c", "b"}}, false},
		{"int", args{common.FieldTypeInt, metadata.PrevIntOption{
			Min: 1,
			Max: 100,
		}, false, 1}, false},
		{"int_default", args{common.FieldTypeInt, "{}", false, 1}, false},
		{"float", args{common.FieldTypeFloat, map[string]interface{}{
			"min": -100.1,
			"max": 100.2,
		}, false, 1.23}, false},
		{"float_default", args{common.FieldTypeFloat, "{}", false, 1}, false},
		{"list", args{common.FieldTypeList, []interface{}{"a", "b", "c"}, false, "a"}, false},
		{"char", args{common.FieldTypeSingleChar, "a.*", false, "abc"}, false},
		{"long_char", args{common.FieldTypeLongChar, "", false, "abc"}, false},
		{"bool", args{common.FieldTypeBool, false, false, true}, false},
		// invalid test
		{"enum", args{common.FieldTypeEnum, `[{"name":"a","type":"text","is_default":true}]`, false, "b"},
			true},
		{"enum_multi", args{common.FieldTypeEnum, metadata.EnumOption{{
			ID:        "a",
			Name:      "a",
			Type:      "aaa",
			IsDefault: false,
		}}, false, nil}, true},
		{"int", args{common.FieldTypeInt, metadata.PrevIntOption{
			Min: 101,
			Max: 100,
		}, false, 100}, true},
		{"int_default", args{common.FieldTypeInt, `{"min":1}`, false, -1}, true},
		{"float", args{common.FieldTypeFloat, map[string]interface{}{
			"min": -100.1,
			"max": 100.2,
		}, false, 111.23}, true},
		{"float_default", args{common.FieldTypeFloat, "{}", false, "a"}, true},
		{"list", args{common.FieldTypeList, []interface{}{"a", "b", "c"}, false, "d"}, true},
		{"char", args{common.FieldTypeSingleChar, "a.*", false, "bc"}, true},
		{"long_char", args{common.FieldTypeLongChar, "*", false, "abc"}, true},
		{"bool", args{common.FieldTypeBool, "false", false, true}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidPropertyOption(kit, tt.args.propertyType, tt.args.option, tt.args.isMultiple,
				tt.args.defaultVal)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidPropertyOption() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIsStrProperty(t *testing.T) {
	type args struct {
		propertyType string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"property"}, false},
		{"", args{common.FieldTypeLongChar}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsStrProperty(tt.args.propertyType); got != tt.want {
				t.Errorf("IsStrProperty() = %v, want %v", got, tt.want)
			}
		})
	}
}
