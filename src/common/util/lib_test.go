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
	"configcenter/src/common"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/require"
)

func TestInArray(t *testing.T) {
	type args struct {
		obj    interface{}
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				target: []string{"a", "b", "c"},
				obj:    "a",
			},
			want: true,
		},
		{
			args: args{
				target: []string{"a", "b", "c"},
				obj:    "d",
			},
			want: false,
		},
		{
			args: args{
				target: []interface{}{"a", "b", "c", 1},
				obj:    1,
			},
			want: true,
		},
		{
			args: args{
				target: []interface{}{"a", "b", "c", 1},
				obj:    int64(1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InArray(tt.args.obj, tt.args.target); got != tt.want {
				t.Errorf("InArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArrayUnique(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantRet []interface{}
	}{
		{
			args: args{
				[]interface{}{"a", "b", "c", 1},
			},
			wantRet: []interface{}{"a", "b", "c", 1},
		},
		{
			args: args{
				[]interface{}{"a", "b", "c", 1, 1, "a", ""},
			},
			wantRet: []interface{}{"a", "b", "c", 1, ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := ArrayUnique(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("ArrayUnique() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestRemoveDuplicatesAndEmpty(t *testing.T) {
	type args struct {
		a []string
	}
	tests := []struct {
		name    string
		args    args
		wantRet []string
	}{
		{
			args: args{
				[]string{"a", "b", "c", "a", ""},
			},
			wantRet: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRet := RemoveDuplicatesAndEmpty(tt.args.a); !reflect.DeepEqual(gotRet, tt.wantRet) {
				t.Errorf("RemoveDuplicatesAndEmpty() = %v, want %v", gotRet, tt.wantRet)
			}
		})
	}
}

func TestStrArrDiff(t *testing.T) {
	type args struct {
		slice1 []string
		slice2 []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{
				[]string{"a", "b", "c", "a", ""},
				[]string{"a", "b"},
			},
			want: []string{"c", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StrArrDiff(tt.args.slice1, tt.args.slice2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StrArrDiff() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetActionLanguage(t *testing.T) {
	req := httptest.NewRequest("POST", "http://127.0.0.1/call", nil)

	language := GetActionLanguage(restful.NewRequest(req))
	//require.Empty(t, language)

	req.Header.Set(common.BKHTTPLanguage, "cn")
	language = GetActionLanguage(restful.NewRequest(req))
	require.Equal(t, "cn", language)

	req.Header.Set(common.BKHTTPLanguage, "cnn")
	language = GetActionLanguage(restful.NewRequest(req))
	require.NotEqual(t, "cn", language)
}

func TestInStrArr(t *testing.T) {
	type args struct {
		arr []string
		key string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{[]string{"key"}, "key"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InStrArr(tt.args.arr, tt.args.key); got != tt.want {
				t.Errorf("InStrArr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetLanguage(t *testing.T) {
	type args struct {
		header http.Header
	}
	header := http.Header{}
	header.Set(common.BKHTTPLanguage, "zh")
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{header}, "zh"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetLanguage(tt.args.header); got != tt.want {
				t.Errorf("GetLanguage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetActionUser(t *testing.T) {
	type args struct {
		req *restful.Request
	}
	req := &http.Request{Header: http.Header{}}
	req.Header.Set(common.BKHTTPHeaderUser, "user")
	r := restful.NewRequest(req)
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{r}, "user"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActionUser(tt.args.req); got != tt.want {
				t.Errorf("GetActionUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetActionOnwerID(t *testing.T) {
	type args struct {
		req *restful.Request
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActionOnwerID(tt.args.req); got != tt.want {
				t.Errorf("GetActionOnwerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUser(tt.args.header); got != tt.want {
				t.Errorf("GetUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOwnerID(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOwnerID(tt.args.header); got != tt.want {
				t.Errorf("GetOwnerID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetOwnerIDAndUser(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetOwnerIDAndUser(tt.args.header)
			if got != tt.want {
				t.Errorf("GetOwnerIDAndUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetOwnerIDAndUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetActionOnwerIDAndUser(t *testing.T) {
	type args struct {
		req *restful.Request
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetActionOnwerIDAndUser(tt.args.req)
			if got != tt.want {
				t.Errorf("GetActionOnwerIDAndUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetActionOnwerIDAndUser() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetActionLanguageByHTTPHeader(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActionLanguageByHTTPHeader(tt.args.header); got != tt.want {
				t.Errorf("GetActionLanguageByHTTPHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetActionOnwerIDByHTTPHeader(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetActionOnwerIDByHTTPHeader(tt.args.header); got != tt.want {
				t.Errorf("GetActionOnwerIDByHTTPHeader() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHTTPCCRequestID(t *testing.T) {
	type args struct {
		header http.Header
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHTTPCCRequestID(tt.args.header); got != tt.want {
				t.Errorf("GetHTTPCCRequestID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Slice_Len(t *testing.T) {
	tests := []struct {
		name string
		p    Int64Slice
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Len(); got != tt.want {
				t.Errorf("Int64Slice.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Slice_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		p    Int64Slice
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("Int64Slice.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInt64Slice_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		p    Int64Slice
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Swap(tt.args.i, tt.args.j)
		})
	}
}
