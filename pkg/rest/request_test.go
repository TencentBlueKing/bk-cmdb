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
	"io"
	"net/http"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockRequest creates a new mock request.
func newMockRequest(t testing.TB, method string, body io.ReadCloser) *http.Request {
	req, err := http.NewRequest(method, "/vm/xxx?name=alices&age=20&age_ptr=21&slice_str=1&slice_str=2", body)
	require.NoError(t, err)
	return req
}

func TestDecode(t *testing.T) {
	r := newMockRequest(t, http.MethodGet, nil)
	type Req struct {
		Name     string   `json:"name" query:"name"`
		Age      int32    `json:"age" query:"age"`
		Bool     bool     `json:"bool" query:"bool"`
		AgePtr   *int     `json:"agePtr" query:"age_ptr"`
		SliceStr []string `json:"sliceStr" query:"slice_str"`
	}

	req, err := decodeReq[Req](r)
	assert.NoError(t, err)
	assert.Equal(t, "alices", req.Name)
	assert.Equal(t, int32(20), req.Age)
	assert.Equal(t, lo.ToPtr(21), req.AgePtr)
	assert.Equal(t, []string{"1", "2"}, req.SliceStr)
}

func BenchmarkDecodeReq(b *testing.B) {
	type Req struct {
		Name     string   `json:"name" query:"name"`
		Age      int32    `json:"age" query:"age" in:"query=age"`
		Bool     bool     `json:"bool" query:"bool"`
		AgePtr   *int     `json:"agePtr" query:"age_ptr"`
		SliceStr []string `json:"sliceStr" query:"slice_str"`
	}

	for b.Loop() {
		r := newMockRequest(b, http.MethodGet, nil)
		req, err := decodeReq[Req](r)
		if err != nil {
			b.Fatal(err)
		}
		if req.Age != int32(20) {
			b.Fatal("age not equal")
		}
	}
}
