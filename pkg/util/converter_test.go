/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - CMDB) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
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

package util

import (
	"reflect"
	"testing"

	"github.com/samber/lo"
)

func TestGetString(t *testing.T) {
	type args struct {
		value any
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
		wantOk  bool
	}{
		{
			name: "str",
			args: args{
				value: "abc",
			},
			wantStr: "abc",
			wantOk:  true,
		},
		{
			name: "int",
			args: args{
				value: 1,
			},
			wantStr: "",
			wantOk:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStr, gotOk := GetString(tt.args.value)
			if gotStr != tt.wantStr {
				t.Errorf("GetString() gotStr = %v, want %v", gotStr, tt.wantStr)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetString() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

var anySliceReflectType = reflect.TypeFor[[]any]()

// BenchmarkGetAnySliceWarp test result :
// BenchmarkGetAnySliceWarp/direct_set
// BenchmarkGetAnySliceWarp/direct_set-16         	 3861835	       281.7 ns/op
// BenchmarkGetAnySliceWarp/reflect_set
// BenchmarkGetAnySliceWarp/reflect_set-16        	 2394316	       493.6 ns/op
// BenchmarkGetAnySliceWarp/lo_convert_to_any
// BenchmarkGetAnySliceWarp/lo_convert_to_any-16  	15207632	        77.99 ns/op
func BenchmarkGetAnySliceWarp(b *testing.B) {
	var targetSlice = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	targetVal := reflect.ValueOf(targetSlice)
	b.Run("direct set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := make([]any, len(targetSlice))
			for idx := range len(targetSlice) {
				slice[idx] = targetVal.Index(idx).Interface()
			}
			slice = nil
		}
	})
	b.Run("reflect set", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			slice := reflect.MakeSlice(anySliceReflectType, 10, 10)
			for idx := range len(targetSlice) {
				slice.Index(idx).Set(targetVal.Index(idx))
			}
			_ = slice.Interface().([]any)
		}
	})
	b.Run("lo convert to any", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// fastest way
			_ = lo.ToAnySlice(targetSlice)
		}
	})

}
