/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/require"
)

func TestGetDailAddress(t *testing.T) {
	type args struct {
		URL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"", args{"http://localhost:80/path?q=a"}, "localhost:80", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDailAddress(tt.args.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDailAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDailAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeekRequest(t *testing.T) {
	expectbody := `{"name":"john"}`
	raw, err := http.NewRequest("POST", "/test/1", bytes.NewBufferString(expectbody))
	require.NoError(t, err)
	req := restful.NewRequest(raw)
	content, err := PeekRequest(req.Request)
	require.NoError(t, err)
	require.Equal(t, expectbody, string(content))

	ncontent, err := ioutil.ReadAll(req.Request.Body)
	require.NoError(t, err)
	require.Equal(t, expectbody, string(ncontent))

	ncontent, err = ioutil.ReadAll(req.Request.Body)
	require.NoError(t, err)
	require.Equal(t, "", string(ncontent))
}
