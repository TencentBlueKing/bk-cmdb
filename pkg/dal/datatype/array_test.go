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

package datatype

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type arrValueTestCase[T ArrayElem] struct {
	name string
	arr  Array[T]
	want string
}

type arrScanFailTestCase[T ArrayElem] struct {
	name          string
	arr           Array[T]
	value         string
	containsError string
}

func TestInt32Array(t *testing.T) {
	tests := []arrValueTestCase[int32]{
		{
			name: "simple",
			arr:  NewArray([]int32{1, 2, 3}),
			want: "{1,2,3}",
		},
		{
			name: "neg",
			arr:  NewArray([]int32{-1, 2, -3}),
			want: "{-1,2,-3}",
		},
		{
			name: "empty array",
			arr:  NewArray([]int32{}),
			want: "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arr.Value()
			assert.Nilf(t, err, "Values() fail")
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
			newArr := NewArray([]int32{})

			err = newArr.Scan(tt.want)
			assert.Nilf(t, err, "Scan() fail")
			if !assert.Equal(t, tt.arr, newArr) {
				t.Errorf("Scan() got = %v, want %v", newArr, tt.arr)
			}

		})
	}

	errCases := []arrScanFailTestCase[int32]{
		{
			name:          "overflow",
			arr:           NewArray([]int32{}),
			value:         "{1234567890123456789012345678901234567890123456789012345678901234567890}",
			containsError: strconv.ErrRange.Error(),
		},
		{
			name:          "syntax error",
			arr:           NewArray([]int32{}),
			value:         "[]int{1,2,4}",
			containsError: "expected '{'",
		},
		{
			name:          "unclosed",
			arr:           NewArray([]int32{}),
			value:         "{1,2,4",
			containsError: "expected '}'",
		},
	}
	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.arr.Scan(tt.value)
			assert.NotNilf(t, err, "Scan() success but should be overflow")
			assert.ErrorContains(t, err, tt.containsError)
		})
	}

}

func TestInt64Array(t *testing.T) {
	tests := []arrValueTestCase[int64]{
		{
			name: "simple",
			arr:  NewArray([]int64{1, 2, 3}),
			want: "{1,2,3}",
		},
		{
			name: "neg",
			arr:  NewArray([]int64{-1, 2, -3}),
			want: "{-1,2,-3}",
		},
		{
			name: "empty array",
			arr:  NewArray([]int64{}),
			want: "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arr.Value()
			assert.Nilf(t, err, "Values() fail")
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
			newArr := NewArray([]int64{})

			err = newArr.Scan(tt.want)
			assert.Nilf(t, err, "Scan() fail")
			if !assert.Equal(t, tt.arr, newArr) {
				t.Errorf("Scan() got = %v, want %v", newArr, tt.arr)
			}

		})
	}
	errCases := []arrScanFailTestCase[int64]{
		{
			name:          "overflow",
			arr:           NewArray([]int64{}),
			value:         "{1234567890123456789012345678901234567890123456789012345678901234567890}",
			containsError: strconv.ErrRange.Error(),
		},
		{
			name:          "syntax error",
			arr:           NewArray([]int64{}),
			value:         "[]int{1,2,4}",
			containsError: "expected '{'",
		},
		{
			name:          "unclosed",
			arr:           NewArray([]int64{}),
			value:         "{1,2,4",
			containsError: "expected '}'",
		},
	}
	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.arr.Scan(tt.value)
			assert.NotNilf(t, err, "Scan() success but should be overflow")
			assert.ErrorContains(t, err, tt.containsError)
		})
	}
}

func TestFloat32Array(t *testing.T) {
	tests := []arrValueTestCase[float32]{
		{
			name: "simple",
			arr:  NewArray([]float32{1.1, 2.2, 3.3}),
			want: "{1.1,2.2,3.3}",
		},
		{
			name: "empty array",
			arr:  NewArray([]float32{}),
			want: "{}",
		},
		{
			name: "negative",
			arr:  NewArray([]float32{-1.1, 2.2, -3.3}),
			want: "{-1.1,2.2,-3.3}",
		},
		{
			name: "zero",
			arr:  NewArray([]float32{0.0, 0.0, 0.0}),
			want: "{0,0,0}",
		},
		{
			name: "truncate",
			arr:  NewArray([]float32{0.00000047683715820312532532564654}),
			want: "{0.00000047683716}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arr.Value()
			assert.Nilf(t, err, "Values() fail")
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
			newArr := NewArray([]float32{})

			err = newArr.Scan(tt.want)
			assert.Nilf(t, err, "Scan() fail")
			if !assert.Equal(t, tt.arr, newArr) {
				t.Errorf("Scan() got = %v, want %v", newArr, tt.arr)
			}

		})
	}
	errCases := []arrScanFailTestCase[float32]{
		{
			name:          "overflow",
			arr:           NewArray([]float32{}),
			value:         "{1234567890123456789012345678901234567890123456789012345678901234567890}",
			containsError: strconv.ErrRange.Error(),
		},
	}

	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.arr.Scan(tt.value)
			assert.NotNilf(t, err, "Scan() success but should be overflow")
			assert.ErrorContains(t, err, tt.containsError)
		})
	}
}

func TestFloat64Array(t *testing.T) {
	tests := []arrValueTestCase[float64]{
		{
			name: "simple",
			arr:  NewArray([]float64{1.1, 2.2, 3.3}),
			want: "{1.1,2.2,3.3}",
		},
		{
			name: "empty array",
			arr:  NewArray([]float64{}),
			want: "{}",
		},
		{
			name: "nil array",
			arr:  NewArray([]float64(nil)),
			want: "{}",
		},
		{
			name: "negative",
			arr:  NewArray([]float64{-1.1, 2.2, -3.3}),
			want: "{-1.1,2.2,-3.3}",
		},
		{
			name: "zero",
			arr:  NewArray([]float64{0.0, 0.0, 0.0}),
			want: "{0,0,0}",
		},
		{
			name: "truncate",
			arr:  NewArray([]float64{0.00000047683715820312532532564654}),
			want: "{0.0000004768371582031253}",
		},
		{
			name: "overflow",
			arr:  NewArray([]float64{1234567890123456789012345678901234567890123456789012345678901234567890}),
			want: "{1234567890123456700000000000000000000000000000000000000000000000000000}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arr.Value()
			assert.Nilf(t, err, "Values() fail")
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
			newArr := NewArray([]float64{})

			err = newArr.Scan(tt.want)
			assert.Nilf(t, err, "Scan() fail")
			if !assert.Equal(t, tt.arr, newArr) {
				t.Errorf("Scan() got = %v, want %v", newArr, tt.arr)
			}

		})
	}
}

func TestStringArray(t *testing.T) {
	tests := []arrValueTestCase[string]{
		{
			name: "simple",
			arr:  NewArray([]string{"1,2,3", "{inside brace}", `"{}"`}),
			want: `{"1,2,3","{inside brace}","\"{}\""}`,
		},
		{
			name: "empty array",
			arr:  NewArray([]string{}),
			want: "{}",
		},
		{
			name: "negative",
			arr:  NewArray([]string{"{{{{{{{{{{{{{{{{{"}),
			want: `{"{{{{{{{{{{{{{{{{{"}`,
		},
		{
			name: "zero",
			arr:  NewArray([]string{""}),
			want: `{""}`,
		},
		{
			name: "negative num",
			arr:  NewArray([]string{"-1.1", "2.2", "-3.3"}),
			want: `{"-1.1","2.2","-3.3"}`,
		},
		{
			name: "zero num",
			arr:  NewArray([]string{"0.0", "0.0", "0.0"}),
			want: `{"0.0","0.0","0.0"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arr.Value()
			assert.Nilf(t, err, "Values() fail")
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
			newArr := NewArray([]string{})

			err = newArr.Scan(tt.want)
			assert.Nilf(t, err, "Scan() fail")
			if !assert.Equal(t, tt.arr, newArr) {
				t.Errorf("Scan() got = %v, want %v", newArr, tt.arr)
			}

		})
	}
	errCases := []arrScanFailTestCase[string]{
		{
			name:          "syntax error",
			arr:           NewArray([]string{}),
			value:         `{", "`,
			containsError: "expected '}'",
		},
		{
			name:          "syntax error",
			arr:           NewArray([]string{}),
			value:         `{"}`,
			containsError: "expected '}'",
		},
		{
			name:          "unclosed",
			arr:           NewArray([]string{}),
			value:         "{1,2,4,",
			containsError: "expected '}'",
		},
	}

	for _, tt := range errCases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.arr.Scan(tt.value)
			assert.NotNilf(t, err, "Scan() success but should be overflow")
			assert.ErrorContains(t, err, tt.containsError)
		})
	}
}

func TestByteaArray(t *testing.T) {
	tests := []arrValueTestCase[[]byte]{
		{
			name: "simple",
			arr:  NewArray([][]byte{[]byte("1,2,3"), []byte("{inside brace}"), []byte(`"{}"`)}),
			want: `{"\\x312c322c33","\\x7b696e736964652062726163657d","\\x227b7d22"}`,
		},
		{
			name: "empty array",
			arr:  NewArray([][]byte{}),
			want: "{}",
		},
		{
			name: "negative",
			arr:  NewArray([][]byte{[]byte("{{{{{{{{{{{{{{{{{{{{")}),
			want: `{"\\x7b7b7b7b7b7b7b7b7b7b7b7b7b7b7b7b7b7b7b7b"}`,
		},
		{
			name: "zero",
			arr:  NewArray([][]byte{[]byte("")}),
			want: `{"\\x"}`,
		},
		{
			name: "negative num",
			arr:  NewArray([][]byte{[]byte("-1.1"), []byte("2.2"), []byte("-3.3")}),
			want: `{"\\x2d312e31","\\x322e32","\\x2d332e33"}`,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.arr.Value()
			assert.Nilf(t, err, "Values() fail")
			if !assert.Equal(t, tt.want, got) {
				t.Errorf("Value() got = %v, want %v", got, tt.want)
			}
			newArr := NewArray([][]byte{})

			err = newArr.Scan(tt.want)
			assert.Nilf(t, err, "Scan() fail")
			if !assert.Equal(t, tt.arr, newArr) {
				t.Errorf("Scan() got = %v, want %v", newArr, tt.arr)
			}

		})
	}
}
func TestNullArray(t *testing.T) {
	intArrayRaw := `{1,2,3}`
	var intArr *Array[int64]
	err := intArr.Scan(intArrayRaw)
	// not support null array scan
	assert.ErrorContains(t, err, "can not scan to nil array", "Scan() should failed with nil array")
}

func TestNullElement(t *testing.T) {
	nullStrArrayRaw := `{"1,2,3",NULL,"{}"}`
	strArray := NewArray([]string{})
	err := strArray.Scan(nullStrArrayRaw)
	// currently not support null
	assert.NotNilf(t, err, "Scan() success but should be failed")
}
