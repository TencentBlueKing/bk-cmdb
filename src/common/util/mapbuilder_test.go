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
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapBuilder(t *testing.T) {
	builder := NewMapBuilder("a", 1)
	assert.Equal(t, map[string]interface{}{"a": 1}, builder.Build())
	builder.Set("b", 2)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2}, builder.Build())
	builder.Append("c", 3)
	assert.Equal(t, map[string]interface{}{"a": 1, "b": 2, "c": []interface{}{3}}, builder.Build())
	builder.Delete("a")
	assert.Equal(t, map[string]interface{}{"b": 2, "c": []interface{}{3}}, builder.Build())
}

func TestMapFromJSON(t *testing.T) {
	info := `{
		"InnerIP" : "127.0.0.1"
	}`
	assert.Equal(t, map[string]interface{}{"InnerIP": "127.0.0.1"}, NewMapFromJSON(info))
}
