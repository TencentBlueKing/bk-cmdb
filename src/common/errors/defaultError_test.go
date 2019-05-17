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

package errors

import (
	"testing"
)

func Test_ccDefaultErrorHelper_New(t *testing.T) {
	type fields struct {
		language  string
		errorStr  func(language string, ErrorCode int) error
		errorStrf func(language string, ErrorCode int, args ...interface{}) error
	}
	type args struct {
		errorCode int
		msg       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"", fields{"zh", func(l string, c int) error {
			return nil
		}, func(l string, c int, a ...interface{}) error {
			return nil
		}}, args{0, ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &ccDefaultErrorHelper{
				language:  tt.fields.language,
				errorStr:  tt.fields.errorStr,
				errorStrf: tt.fields.errorStrf,
			}
			if err := cli.New(tt.args.errorCode, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("ccDefaultErrorHelper.New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ccDefaultErrorHelper_Error(t *testing.T) {
	type fields struct {
		language  string
		errorStr  func(language string, ErrorCode int) error
		errorStrf func(language string, ErrorCode int, args ...interface{}) error
	}
	type args struct {
		errCode int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"", fields{"zh", func(l string, c int) error {
			return nil
		}, func(l string, c int, a ...interface{}) error {
			return nil
		}}, args{0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &ccDefaultErrorHelper{
				language:  tt.fields.language,
				errorStr:  tt.fields.errorStr,
				errorStrf: tt.fields.errorStrf,
			}
			if err := cli.Error(tt.args.errCode); (err != nil) != tt.wantErr {
				t.Errorf("ccDefaultErrorHelper.Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_ccDefaultErrorHelper_Errorf(t *testing.T) {
	type fields struct {
		language  string
		errorStr  func(language string, ErrorCode int) error
		errorStrf func(language string, ErrorCode int, args ...interface{}) error
	}
	type args struct {
		errCode int
		args    []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cli := &ccDefaultErrorHelper{
				language:  tt.fields.language,
				errorStr:  tt.fields.errorStr,
				errorStrf: tt.fields.errorStrf,
			}
			if err := cli.Errorf(tt.args.errCode, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("ccDefaultErrorHelper.Errorf() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
