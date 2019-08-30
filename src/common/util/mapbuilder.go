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
	"net/http"
)

type MapBuilder struct {
	value map[string]interface{}
}

func NewMapBuilder(kvPairs ...interface{}) *MapBuilder {
	value := map[string]interface{}{}
	for i := range kvPairs {
		if i%2 == 0 {
			value[kvPairs[i].(string)] = kvPairs[i+1]
		}
	}
	return &MapBuilder{value}
}

func (m *MapBuilder) Build() map[string]interface{} {
	return m.value
}

func (m *MapBuilder) Set(k string, v interface{}) {
	m.value[k] = v
}

func (m *MapBuilder) Append(k string, vs ...interface{}) {
	_, ok := m.value[k]
	if !ok {
		m.value[k] = []interface{}{}
	}
	m.value[k] = append(m.value[k].([]interface{}), vs...)
}

func (m *MapBuilder) Delete(k string) {
	delete(m.value, k)
}

func CopyMap(data map[string]interface{}, keys []string, ignores []string) map[string]interface{} {
	newInst := make(map[string]interface{})

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
		newInst[key] = data[key]
	}
	return newInst

}

// CloneHeader clone http header
func CloneHeader(src http.Header) http.Header {
	tar := http.Header{}
	for key := range src {
		tar.Set(key, src.Get(key))
	}
	return tar
}

// CopyHeader copy http header into target
func CopyHeader(src http.Header, target http.Header) {
	for key := range src {
		target.Set(key, src.Get(key))
	}
}
