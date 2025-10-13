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

package codec

import (
	"fmt"
	"strings"
)

const (
	// tagName 结构体tag名称
	// 格式参考 https://pkg.go.dev/encoding/json/v2#example-package-FormatFlags
	tagName      = "req"
	inOptName    = "in"
	queryOptName = "query"
	pathOptName  = "path"
)

// Tag is a struct tag
type Tag struct {
	Name   string            // tag name
	In     string            // query/path参数中
	Option map[string]string // 自定义参数
}

func parseTag(tagStr string) (*Tag, error) {
	tagStr = strings.TrimSpace(tagStr)
	if tagStr == "" {
		return nil, fmt.Errorf("tag is empty")
	}

	parts := strings.Split(tagStr, ",")
	name := strings.TrimSpace(parts[0])

	t := &Tag{
		Name:   name,
		Option: map[string]string{},
	}

	for _, part := range parts[1:] {
		part = strings.TrimSpace(part)
		if part == "" {
			return nil, fmt.Errorf("tag option not valid")
		}

		opt := strings.SplitN(part, ":", 2)
		key := opt[0]
		val := ""
		if len(opt) == 2 {
			val = opt[1]
		}

		if key == inOptName {
			t.In = val
		} else {
			t.Option[key] = val
		}
	}

	return t, nil
}
