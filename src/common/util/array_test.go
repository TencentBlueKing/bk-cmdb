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
	"reflect"
	"testing"
)

func TestStrArrayUnique(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []string
	}{
		{"", args{[]string{"1", "1"}}, []string{"1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := StrArrayUnique(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("StrArrayUnique() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestIntArrayUnique(t *testing.T) {
	type args struct {
		a []int64
	}
	tests := []struct {
		name    string
		args    args
		wantRet []int64
	}{
		{"", args{[]int64{1, 1}}, []int64{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := IntArrayUnique(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("IntArrayUnique() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestIntArrIntersection(t *testing.T) {
	type args struct {
		slice1 []int64
		slice2 []int64
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{"", args{[]int64{1}, []int64{2, 1}}, []int64{1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntArrIntersection(tt.args.slice1, tt.args.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("IntArrIntersection() = %v, want %v", got, tt.want)
			}
		})
	}
}
