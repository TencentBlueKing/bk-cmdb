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
)

func TestCheckLen(t *testing.T) {
	type args struct {
		sInput string
		min    int
		max    int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{"123", 0, 3},
			want: true,
		},
		{
			args: args{"123", 1, 2},
			want: false,
		},
		{
			args: args{"123", -1, 3},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckLen(tt.args.sInput, tt.args.min, tt.args.max); got != tt.want {
				t.Errorf("CheckLen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsChar(t *testing.T) {
	type args struct {
		sInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"c"}, want: true},
		{args: args{" c"}, want: false},
		{args: args{"c "}, want: false},
		{args: args{"和"}, want: false},
		{args: args{"_"}, want: false},
		{args: args{"3"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsChar(tt.args.sInput); got != tt.want {
				t.Errorf("IsChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNumChar(t *testing.T) {
	type args struct {
		sInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"1"}, want: true},
		{args: args{"aA1"}, want: true},
		{args: args{" 1"}, want: false},
		{args: args{"1 "}, want: false},
		{args: args{"和"}, want: false},
		{args: args{"_"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumChar(tt.args.sInput); got != tt.want {
				t.Errorf("IsNumChar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDate(t *testing.T) {
	type args struct {
		sInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"2018-10-10"}, want: true},
		{args: args{"2018/10/10"}, want: false},
		{args: args{`2018\10\10`}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsDate(tt.args.sInput); got != tt.want {
				t.Errorf("IsDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTime(t *testing.T) {
	type args struct {
		sInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{args: args{"2018-10-10 10:56:67"}, want: true},
		{args: args{"105667"}, want: false},
		{args: args{`10-56-67`}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTime(tt.args.sInput); got != tt.want {
				t.Errorf("IsTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
