/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"encoding/json"
	"net/http"
)

type MapBuiler struct {
	value map[string]interface{}
}

func NewMapBuilder(kvpairs ...interface{}) *MapBuiler {
	value := map[string]interface{}{}
	for i := range kvpairs {
		if i%2 == 0 {
			value[kvpairs[i].(string)] = kvpairs[i+1]
		}
	}
	return &MapBuiler{value}
}

func (m *MapBuiler) Build() map[string]interface{} {
	return m.value
}

func (m *MapBuiler) Set(k string, v interface{}) {
	m.value[k] = v
}

func (m *MapBuiler) Append(k string, vs ...interface{}) {
	_, ok := m.value[k]
	if !ok {
		m.value[k] = []interface{}{}
	}
	m.value[k] = append(m.value[k].([]interface{}), vs...)
}

func (m *MapBuiler) Delete(k string) {
	delete(m.value, k)
}

func NewMapFromJSON(data string) map[string]interface{} {
	value := map[string]interface{}{}
	_ = json.Unmarshal([]byte(data), &value)
	return value
}

func CopyMap(data map[string]interface{}, keys []string, ignores []string) map[string]interface{} {
	newinst := map[string]interface{}{}

	ignore := map[string]bool{}
	for _, key := range ignores {
		ignore[key] = true
	}
	if len(keys) <= 0 {
		for key := range data {
			keys = append(keys, key)
		}
	}
	for _, key := range keys {
		if ignore[key] {
			continue
		}
		newinst[key] = data[key]
	}
	return newinst

}

// CopyHeader copy http header
func CopyHeader(src http.Header) http.Header {
	tar := http.Header{}
	for key := range src {
		tar.Set(key, src.Get(key))
	}
	return tar
}
