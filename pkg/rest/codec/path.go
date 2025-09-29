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
	"net/http"
	"net/url"
	"reflect"
)

type pathCodec struct {
	values    url.Values
	cachedTag map[string]*Tag
}

// NewPathCodec ...
func NewPathCodec(r *http.Request, rt reflect.Type) (*pathCodec, error) {
	c := &pathCodec{cachedTag: make(map[string]*Tag)}

	if err := c.makePathValues(r, rt); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *pathCodec) makePathValues(r *http.Request, rt reflect.Type) error {
	values := url.Values{}
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		// 非导出需要跳过, 无法设置值
		if !field.IsExported() {
			continue
		}

		pathTag, ok := field.Tag.Lookup("path")
		if !ok {
			continue
		}

		tag, err := ParseTag(pathTag)
		if err != nil {
			return err
		}
		c.cachedTag[pathTag] = tag

		values[tag.Name] = []string{r.PathValue(tag.Name)}
	}

	c.values = values
	return nil
}

// Decode ...
func (c *pathCodec) Decode(field reflect.StructField, fv reflect.Value) error {
	pathTag, ok := field.Tag.Lookup("path")
	if !ok {
		return nil
	}

	tag := c.cachedTag[pathTag]

	v, ok := c.values[tag.Name]
	if !ok {
		return nil
	}

	rv, err := getFieldValue(field.Type, tag, v)
	if err != nil {
		return err
	}
	if !rv.IsValid() {
		return fmt.Errorf("%s not valid", rv)
	}

	fv.Set(rv)
	return nil
}
