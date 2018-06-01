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
 
package httpclient

import (
	"net/http"
	"testing"
)

const (
	uri = "http://127.0.0.1:8081/"
)

func Test_RequestEx(t *testing.T) {
	cli := NewHttpClient()

	type Args struct {
		url    string
		method string
		header http.Header
		data   []byte
	}

	tests := []struct {
		name string
		args Args
		want int
	}{
		{
			args: Args{
				url:    uri + "testnode",
				method: "GET",
				header: nil,
				data:   nil,
			},
			want: 200,
		},
		{
			args: Args{
				url:    uri + "testnode",
				method: "POST",
				header: nil,
				data:   []byte("test"),
			},
			want: 200,
		},
		{
			args: Args{
				url:    uri + "testnode",
				method: "PUT",
				header: nil,
				data:   []byte("test"),
			},
			want: 200,
		},
		{
			args: Args{
				url:    uri + "testnode",
				method: "DELETE",
				header: nil,
				data:   nil,
			},
			want: 200,
		},
		{
			args: Args{
				url:    uri + "testnode",
				method: "PATCH",
				header: nil,
				data:   []byte("test"),
			},
			want: 200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, _, err := cli.RequestEx(tt.args.url, tt.args.method, tt.args.header, tt.args.data)
			if err != nil {
				t.Errorf("fail to do http request to url(%s), err:%s", tt.args.url, err.Error())
			}

			if code != tt.want {
				t.Errorf("RequestEx() return code: %d, but want: %d", code, tt.want)
			}
		})
	}
}
