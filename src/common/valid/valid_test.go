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

package util

import (
	"testing"

	"configcenter/src/common"
	"configcenter/src/common/errors"
)

type errif struct {
}

func (ei errif) Error(errCode int) error {
	return nil
}

func (ei errif) Errorf(errCode int, args ...interface{}) error {
	return nil
}

func (ei errif) New(errCode int, msg string) error {
	return nil
}

func TestValidPropertyOption(t *testing.T) {
	type args struct {
		propertyType string
		option       interface{}
		errProxy     errors.DefaultCCErrorIf
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"", args{"property", "option", errif{}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidPropertyOption(tt.args.propertyType, tt.args.option, tt.args.errProxy); (err != nil) != tt.wantErr {
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

func TestIsInnerObject(t *testing.T) {
	type args struct {
		objID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"id"}, false},
		{"", args{common.BKInnerObjIDApp}, true},
		{"", args{common.BKInnerObjIDHost}, true},
		{"", args{common.BKInnerObjIDModule}, true},
		{"", args{common.BKInnerObjIDPlat}, true},
		{"", args{common.BKInnerObjIDProc}, true},
		{"", args{common.BKInnerObjIDSet}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInnerObject(tt.args.objID); got != tt.want {
				t.Errorf("IsInnerObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
