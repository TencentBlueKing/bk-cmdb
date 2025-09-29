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
	"fmt"
	"strings"
)

// Tag is a struct tag
type Tag struct {
	Name   string
	Option map[string]string
}

// ParseTag ...
func ParseTag(tagStr string) (*Tag, error) {
	if tagStr == "" {
		return nil, fmt.Errorf("tag is empty")
	}

	tagStr = strings.TrimSpace(tagStr)
	parts := strings.Split(tagStr, ",")
	name := strings.TrimSpace(parts[0])
	if name == "" {
		return nil, fmt.Errorf("tag name is empty")
	}

	t := &Tag{
		Name:   name,
		Option: map[string]string{},
	}

	for _, part := range parts[1:] {
		if part == "" {
			return nil, fmt.Errorf("tag option not valid")
		}

		part = strings.TrimSpace(part)
		opt := strings.SplitN(part, "=", 2)
		key := opt[0]
		val := ""
		if len(opt) == 2 {
			val = opt[1]
		}
		t.Option[key] = val
	}

	return t, nil
}
