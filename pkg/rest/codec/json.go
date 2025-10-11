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

package codec

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type jsonCodec struct {
	isJson bool
	req    *http.Request
}

// NewJsonCodec ...
func NewJsonCodec(r *http.Request) *jsonCodec {
	isJson := false

	// 限制Method, 同ParseForm的一致
	contentType := r.Header.Get("Content-Type")
	if (r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH") &&
		strings.HasPrefix(contentType, "application/json") {
		isJson = true
	}

	return &jsonCodec{req: r, isJson: isJson}
}

// Decode ...
func (j *jsonCodec) Decode(val any) error {
	if !j.isJson {
		return nil
	}

	body, err := io.ReadAll(j.req.Body)
	if err != nil {
		return err
	}

	// body等于空时，可能其他解析场景，直接正常返回
	// 如果需要判断是否有值，可通过指针处理
	if len(body) == 0 {
		return nil
	}

	if err := json.Unmarshal(body, val); err != nil {
		return fmt.Errorf("unmarshal json body: %w", err)
	}
	return nil
}
