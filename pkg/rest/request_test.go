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

	"github.com/stretchr/testify/require"
)

// NewMockRequest creates a new mock request.
func NewMockRequest(t *testing.T, method string, body io.ReadCloser) *http.Request {
	req, err := http.NewRequest(method, "/vm/xxx?name=alice", body)
	require.NoError(t, err)
	return req
}

func BenchmarkDecodeReq(b *testing.B) {
	r := &http.Request{
		Method: http.MethodPost,
	}

	for i := 0; i < b.N; i++ {
		result, err := decodeReq[http.Request](r)
		if err != nil {
			b.Error(err)
		}
		if result.Method != http.MethodPost {
			b.Error("invalid result")
		}
	}
}
