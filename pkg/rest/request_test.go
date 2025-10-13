/*
 * Tencent is pleased to support the open source community by making
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

package rest

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TencentBlueKing/bk-cmdb/pkg/rest/codec"
)

// reqStruct for rest
type reqStruct struct {
	Org      string   `json:"-" req:"org,in:path" in:"path=org"`
	Name     string   `json:"-" req:"name,in:query" in:"query=name"`
	Age      int32    `json:"age" req:"-,in:form"`
	Bool     bool     `json:"-" req:"bool,in:query" in:"query=bool"`
	AgePtr   *int     `json:"-" req:"age_ptr,in:query" in:"query=age_ptr"`
	SliceStr []string `json:"-" req:"slice_str,in:query" in:"query=slice_str"`
	Page     int64    `json:"page"`
}

// newMockRequest creates a new mock request.
func newMockRequest(t testing.TB, method string, header map[string]string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, "/{org}/vm/xxx?name=alices&age=20&age_ptr=21&slice_str=1&slice_str=2", body)
	req.SetPathValue("org", "myOrg")
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Body = io.NopCloser(body)

	require.NoError(t, err)
	return req
}

func TestDecode(t *testing.T) {
	header := map[string]string{
		"age_ptr": "21",
	}
	r := newMockRequest(t, http.MethodGet, header, nil)

	req, err := decodeReq[reqStruct](r)
	assert.NoError(t, err)
	assert.Equal(t, "myOrg", req.Org)
	assert.Equal(t, "alices", req.Name)
	assert.Equal(t, lo.ToPtr(21), req.AgePtr)
	assert.Equal(t, []string{"1", "2"}, req.SliceStr)
}

func TestJsonDecode(t *testing.T) {
	header := map[string]string{
		"age_ptr":      "21",
		"Content-Type": "application/json",
	}

	jsonData := `{"page": 64}`
	r := newMockRequest(t, http.MethodPost, header, bytes.NewBufferString(jsonData))

	req, err := decodeReq[reqStruct](r)
	assert.NoError(t, err)
	assert.Equal(t, "myOrg", req.Org)
	assert.Equal(t, "alices", req.Name)
	assert.Equal(t, lo.ToPtr(21), req.AgePtr)
	assert.Equal(t, []string{"1", "2"}, req.SliceStr)
	assert.Equal(t, int64(64), req.Page)
}

func TestDecodeErr(t *testing.T) {
	header := map[string]string{
		"age_ptr": "21",
	}

	// array not support
	type Req2 struct {
		SliceStr [1]string `json:"-" req:"slice_str,in:query"`
	}
	r := newMockRequest(t, http.MethodGet, header, nil)
	_, err := decodeReq[Req2](r)
	assert.ErrorIs(t, err, codec.ErrUnsupportedType)
}

func BenchmarkDecodeReq(b *testing.B) {
	for b.Loop() {
		r := newMockRequest(b, http.MethodGet, nil, nil)
		req, err := decodeReq[reqStruct](r)
		if err != nil {
			b.Fatal(err)
		}
		if req.Name != "alices" {
			b.Fatal("name not equal")
		}
	}
}
