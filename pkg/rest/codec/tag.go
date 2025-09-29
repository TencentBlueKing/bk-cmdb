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
	"reflect"
	"strings"
)

const (
	// tagName 结构体tag名称
	// 格式参考 https://pkg.go.dev/encoding/json/v2#example-package-FormatFlags
	tagName = "req"
)

// Tag is a struct tag
type Tag struct {
	Option map[string]string
}

func getStructTags(rt reflect.Type) ([]*Tag, error) {
	tags := make([]*Tag, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		tags[i] = new(Tag)

		// 非导出需要跳过, 无法设置值
		if !field.IsExported() {
			continue
		}

		tagStr, ok := field.Tag.Lookup(tagName)
		if !ok {
			continue
		}

		tag, err := parseTag(tagStr)
		if err != nil {
			return nil, err
		}
		tags[i] = tag
	}

	return tags, nil
}

func parseTag(tagStr string) (*Tag, error) {
	tagStr = strings.TrimSpace(tagStr)
	if tagStr == "" {
		return nil, fmt.Errorf("tag is empty")
	}

	parts := strings.Split(tagStr, ",")
	t := &Tag{
		Option: map[string]string{},
	}

	for _, part := range parts {
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
		t.Option[key] = val
	}

	return t, nil
}
