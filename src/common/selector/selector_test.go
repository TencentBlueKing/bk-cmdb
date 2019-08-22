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

package selector_test

import (
	"fmt"
	"strings"
	"testing"

	"configcenter/src/common/selector"

	"github.com/stretchr/testify/assert"
)

func TestLabel(t *testing.T) {
	// assert normal
	var errorKey string
	var err error
	var labels selector.Labels

	key := "key"
	value := "value"
	for _, key = range []string{"key", "0key", "key-.k", strings.Repeat("k", 63)} {
		labels := selector.Labels{
			key: value,
		}
		errorKey, err := labels.Validate()
		assert.Nil(t, err)
		assert.Empty(t, errorKey)
	}

	// assert key err
	for _, key = range []string{"-key", ".key", "_key", "key-", "key_", "key.", strings.Repeat("k", 64)} {
		labels := selector.Labels{
			key: value,
		}
		errorKey, err = labels.Validate()
		assert.NotNil(t, err)
		assert.Equal(t, key, errorKey)
	}

	// assert value err
	key = "key"
	for _, value = range []string{"value", "0value", "value-_.v", strings.Repeat("v", 63)} {
		for _, value = range []string{"-value", ".value", "_value", "value-", "value_", "value.", strings.Repeat("v", 64)} {
			labels := selector.Labels{
				key: value,
			}
			errorKey, err = labels.Validate()
			assert.NotNil(t, err)
			assert.Equal(t, fmt.Sprintf("%s:%s", key, value), errorKey)
		}
	}

	// assert update
	labels = selector.Labels{
		"key1": "value1",
	}
	labels2 := selector.Labels{
		"key2": "value2",
	}
	labels.AddLabel(labels2)
	labels.RemoveLabel([]string{"key1", "key2", "key3"})
	assert.Empty(t, labels)
}

func TestSelector(t *testing.T) {
	// assert normal
	sl := selector.Selectors{
		{
			Key:      "key",
			Operator: "=",
			Values:   []string{"value"},
		}, {
			Key:      "key",
			Operator: "!=",
			Values:   []string{"value"},
		}, {
			Key:      "key",
			Operator: "in",
			Values:   []string{"value", "value1"},
		}, {
			Key:      "key",
			Operator: "notin",
			Values:   []string{"value", "value1"},
		}, {
			Key:      "key",
			Operator: "exists",
			Values:   []string{},
		}, {
			Key:      "key",
			Operator: "!",
			Values:   []string{},
		},
	}
	errKey, err := sl.Validate()
	assert.Nil(t, err)
	assert.Empty(t, errKey)
	filter, err := sl.ToMgoFilter()
	assert.Nil(t, err)
	assert.NotEmpty(t, filter)

	// assert abnormal
	ss := selector.Selectors{
		{
			Key:      "key",
			Operator: "=",
			Values:   []string{"value", "value2"},
		}, {
			Key:      "key",
			Operator: "!=",
			Values:   []string{"value", "value2"},
		}, {
			Key:      "key",
			Operator: "in",
			Values:   []string{},
		}, {
			Key:      "key",
			Operator: "notin",
			Values:   []string{},
		}, {
			Key:      "key",
			Operator: "exists",
			Values:   []string{"value"},
		}, {
			Key:      "key",
			Operator: "!",
			Values:   []string{"value"},
		}, {
			Key:      ".key",
			Operator: "=",
			Values:   []string{"value"},
		}, {
			Key:      "key",
			Operator: "?",
			Values:   []string{"value"},
		},
	}
	for _, sl := range ss {
		errKey, err = sl.Validate()
		assert.NotNil(t, err)
		assert.NotEmpty(t, errKey)
		sl.ToMgoFilter()
	}

	// assert to filter abnormal
	ss = selector.Selectors{
		{
			Key:      "key",
			Operator: "=",
			Values:   []string{},
		}, {
			Key:      "key",
			Operator: "!=",
			Values:   []string{},
		},
	}
	for _, sl := range ss {
		filter, err := sl.ToMgoFilter()
		assert.NotNil(t, err)
		assert.Empty(t, filter)
	}
	key, err := ss.Validate()
	assert.NotEmpty(t, key)
	assert.NotNil(t, err)
	filter, err = ss.ToMgoFilter()
	assert.Empty(t, filter)
	assert.NotNil(t, err)
}
