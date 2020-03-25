/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package auditlog

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/require"
)

func TestCmpData(t *testing.T) {
	type testData struct {
		src       map[string]interface{}
		dst       map[string]interface{}
		ignoreKey map[string]interface{}
		result    bool
		desc      string
	}

	testDataArr := []testData{
		testData{
			src:       map[string]interface{}{"_id": 1, "id": 2},
			dst:       map[string]interface{}{"id": 2},
			ignoreKey: map[string]interface{}{"_id": ""},
			result:    true,
			desc:      "测试忽略字段在目标值中不存在",
		},
		testData{
			src:       map[string]interface{}{"prefix": "2", "_id": 1, "id": 2},
			dst:       map[string]interface{}{"id": 2},
			ignoreKey: map[string]interface{}{"_id": "", "prefix": ""},
			result:    true,
			desc:      "测试多个忽略字段在目标值中不存在",
		},
		testData{
			src:       map[string]interface{}{"_id": 1, "id": 2},
			dst:       map[string]interface{}{"_id": 2, "id": 2},
			ignoreKey: map[string]interface{}{"_id": ""},
			result:    true,
			desc:      "测试忽略字段在目标值中不同",
		},
		testData{
			src:       map[string]interface{}{"id": 2},
			dst:       map[string]interface{}{"_id": 2, "id": 2},
			ignoreKey: map[string]interface{}{"_id": ""},
			result:    true,
			desc:      "测试忽略字段在源数据不存在",
		},
	}

	for _, item := range testDataArr {

		option := cmpopts.IgnoreMapEntries(ignorePath(item.ignoreKey))
		bl := cmp.Equal(item.src, item.dst, option)
		require.Equal(t, item.result, bl)

	}

}

func TestInstNotChange(t *testing.T) {
	type testData struct {
		content map[string]interface{}
		result  bool
		desc    string
	}

	testDataArr := []testData{
		testData{
			content: map[string]interface{}{
				"pre_data": map[string]interface{}{"_id": 1, "id": 2},
				"cur_data": map[string]interface{}{"id": 2},
			},
			result: true,
			desc:   "测试忽略字段在目标值中不存在",
		},
		testData{
			content: map[string]interface{}{
				"pre_data": map[string]interface{}{"last_time": "2", "_id": 1, "id": 2},
				"cur_data": map[string]interface{}{"id": 2},
			},
			result: true,
			desc:   "测试多个忽略字段在目标值中不存在",
		},
		testData{
			content: map[string]interface{}{
				"pre_data": map[string]interface{}{"_id": 1, "id": 2},
				"cur_data": map[string]interface{}{"_id": 2, "id": 2},
			},
			result: true,
			desc:   "测试忽略字段在目标值中不同",
		},
		testData{
			content: map[string]interface{}{
				"pre_data": map[string]interface{}{"id": 2},
				"cur_data": map[string]interface{}{"_id": 2, "id": 2},
			},
			result: true,
			desc:   "测试忽略字段在源数据不存在",
		},

		testData{
			content: map[string]interface{}{
				"pre_data": map[string]interface{}{"id": 3},
				"cur_data": map[string]interface{}{"_id": 2, "id": 2},
			},
			result: false,
			desc:   "值不同比较",
		},
		testData{
			content: map[string]interface{}{
				"pre_data": map[string]interface{}{"id——bak": 3},
				"cur_data": map[string]interface{}{"id": 2},
			},
			result: false,
			desc:   "字段不同比较",
		},
	}

	for _, item := range testDataArr {

		bl := instNotChange(context.Background(), item.content, "")
		require.Equal(t, item.result, bl)

	}

}

func ignorePath(kvMap map[string]interface{}) func(string, interface{}) bool {
	funcHandle := func(key string, value interface{}) bool {
		if _, ok := kvMap[key]; ok {
			return true
		}
		return false
	}

	return funcHandle

}
