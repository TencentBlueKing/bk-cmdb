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
	"strings"
)

type formCodec struct {
	values url.Values
	isForm bool
}

// NewFormCodec ...
func NewFormCodec(r *http.Request) (*formCodec, error) {
	isForm := false

	contentType := r.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		if err := r.ParseForm(); err != nil {
			return nil, err
		}
		isForm = true
	}

	c := &formCodec{
		values: r.PostForm,
		isForm: isForm,
	}
	return c, nil
}

// Decode ...
func (c *formCodec) Decode(field reflect.StructField, fv reflect.Value, tag *Tag) error {
	if !c.isForm {
		return nil
	}

	formTag, ok := tag.Option["form"]
	if !ok {
		return nil
	}

	v := c.values[formTag]
	if len(v) == 0 {
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
