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
	"errors"
	"reflect"
	"testing"
)

func TestNewParseInterface(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name string
		args args
		want *ParseInterface
	}{
		{"", args{nil}, &ParseInterface{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewParseInterface(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewParseInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInterface_Get(t *testing.T) {
	type fields struct {
		data interface{}
		err  error
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *ParseInterface
	}{
		{"", fields{nil, nil}, args{"key"}, &ParseInterface{nil, errors.New("key not found")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ParseInterface{
				data: tt.fields.data,
				err:  tt.fields.err,
			}
			if got := p.Get(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInterface.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInterface_Interface(t *testing.T) {
	type fields struct {
		data interface{}
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		want    interface{}
		wantErr bool
	}{
		{"", fields{1, nil}, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ParseInterface{
				data: tt.fields.data,
				err:  tt.fields.err,
			}
			got, err := p.Interface()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterface.Interface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInterface.Interface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInterface_String(t *testing.T) {
	type fields struct {
		data interface{}
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"", fields{"string", nil}, "string", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ParseInterface{
				data: tt.fields.data,
				err:  tt.fields.err,
			}
			got, err := p.String()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterface.String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseInterface.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseInterface_ArrayInterface(t *testing.T) {
	type fields struct {
		data interface{}
		err  error
	}
	tests := []struct {
		name    string
		fields  fields
		want    []interface{}
		wantErr bool
	}{
		{"", fields{[]interface{}{1, 2}, nil}, []interface{}{1, 2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ParseInterface{
				data: tt.fields.data,
				err:  tt.fields.err,
			}
			got, err := p.ArrayInterface()
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseInterface.ArrayInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseInterface.ArrayInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}
