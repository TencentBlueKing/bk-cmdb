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
	"encoding/json"
	"reflect"
	"testing"
)

func TestGetIntByInterface(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			args: args{
				a: int(1),
			},
			want: 1,
		},
		{
			args: args{
				a: int32(1),
			},
			want: 1,
		},
		{
			args: args{
				a: int64(1),
			},
			want: 1,
		},
		{
			args: args{
				a: float32(1.01),
			},
			want: 1,
		},
		{
			args: args{
				a: float64(1.01),
			},
			want: 1,
		},
		{
			args: args{
				a: "1",
			},
			want: 1,
		},
		{
			args: args{
				a: json.Number("1"),
			},
			want: 1,
		},
		{
			args: args{
				a: "a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIntByInterface(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIntByInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetIntByInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetInt64ByInterface(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			args: args{
				a: int(1),
			},
			want: 1,
		},
		{
			args: args{
				a: int32(1),
			},
			want: 1,
		},
		{
			args: args{
				a: int64(1),
			},
			want: 1,
		},
		{
			args: args{
				a: float32(1.01),
			},
			want: 1,
		},
		{
			args: args{
				a: float64(1.01),
			},
			want: 1,
		},
		{
			args: args{
				a: "1",
			},
			want: 1,
		},
		{
			args: args{
				a: json.Number("1"),
			},
			want: 1,
		},
		{
			args: args{
				a: "a",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetInt64ByInterface(tt.args.a)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInt64ByInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetInt64ByInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMapInterfaceByInerface(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			args: args{
				[]int{1, 2, 3},
			},
			want: []interface{}{1, 2, 3},
		},
		{
			args: args{
				[]int64{1, 2, 3},
			},
			want: []interface{}{int64(1), int64(2), int64(3)},
		},
		{
			args: args{
				[]int32{1, 2, 3},
			},
			want: []interface{}{int32(1), int32(2), int32(3)},
		},
		{
			args: args{
				[]string{"1", "2", "3"},
			},
			want: []interface{}{"1", "2", "3"},
		},
		{
			args: args{
				"123",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMapInterfaceByInerface(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMapInterfaceByInerface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMapInterfaceByInerface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStrByInterface(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"string"}, "string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStrByInterface(tt.args.a); got != tt.want {
				t.Errorf("GetStrByInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStrToInt(t *testing.T) {
	type args struct {
		sliceStr []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int
		wantErr bool
	}{
		{"", args{[]string{"1"}}, []int{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SliceStrToInt(tt.args.sliceStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceStrToInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceStrToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceStrToInt64(t *testing.T) {
	type args struct {
		sliceStr []string
	}
	tests := []struct {
		name    string
		args    args
		want    []int64
		wantErr bool
	}{
		{"", args{[]string{"1"}}, []int64{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SliceStrToInt64(tt.args.sliceStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceStrToInt64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SliceStrToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetStrValsFromArrMapInterfaceByKey(t *testing.T) {
	type args struct {
		arrI []interface{}
		key  string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"", args{[]interface{}{map[string]interface{}{"key": "string"}}, "key"}, []string{"string"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetStrValsFromArrMapInterfaceByKey(tt.args.arrI, tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStrValsFromArrMapInterfaceByKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
