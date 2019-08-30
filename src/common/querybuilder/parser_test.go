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

package querybuilder_test

import (
	"encoding/json"
	"testing"

	"configcenter/src/common/querybuilder"

	"github.com/stretchr/testify/assert"
)

func TestNormalParser(t *testing.T) {
	cases := []map[string]interface{}{
		{
			"condition": "AND",
			"rules": []map[string]interface{}{
				{
					"operator": "equal",
					"field":    "field",
					"value":    1,
				},
			},
		}, {
			"condition": "AND",
			"rules": []map[string]interface{}{
				{
					"operator": "equal",
					"field":    "field",
					"value":    1,
				}, {
					"condition": "AND",
					"rules": []map[string]interface{}{
						{
							"operator": "equal",
							"field":    "field",
							"value":    1,
						},
					},
				},
			},
		}, {
			"operator": "equal",
			"field":    "field",
			"value":    1,
		},
	}
	for idx, data := range cases {
		t.Logf("running normal parser, idx: %d, data: %+v", idx, data)
		filter, errKey, err := querybuilder.ParseRule(data)
		assert.Nil(t, err)
		assert.Empty(t, errKey)
		assert.NotNil(t, filter)
	}

	data := map[string]interface{}(nil)
	filter, errKey, err := querybuilder.ParseRule(data)
	assert.Nil(t, err)
	assert.Empty(t, errKey)
	assert.Nil(t, filter)
}

func TestAsStructField(t *testing.T) {
	type Foo struct {
		QueryFilter querybuilder.QueryFilter `json:"query_filter"`
		Key         string                   `json:"key"`
	}
	foo := new(Foo)
	input := `{"key": "test", "query_filter": {"operator":"equal", "value":"1", "field":"field"}}`
	err := json.Unmarshal([]byte(input), &foo)
	assert.Nil(t, err)

	output, err := json.Marshal(foo)
	assert.Nil(t, err)
	t.Logf("output: %s", output)
}
