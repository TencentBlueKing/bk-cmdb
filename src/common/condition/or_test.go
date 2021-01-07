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

package condition

import (
	"encoding/json"
	"testing"

	//"configcenter/src/common/blog"
	types "configcenter/src/common/mapstr"
)

func TestORField(t *testing.T) {
	type testData struct {
		input  []types.MapStr
		output string
	}

	testDataArr := []testData{
		testData{
			input:  []types.MapStr{types.MapStr{"a": "c"}},
			output: `{"$or":[{"a":"c"}]}`,
		},
		testData{
			input:  []types.MapStr{types.MapStr{"a": "c"}, types.MapStr{"b": "c"}},
			output: `{"$or":[{"a":"c"},{"b":"c"}]}`,
		},
	}
	for _, item := range testDataArr {
		f := &orField{
			condition: CreateCondition(),
		}
		for _, input := range item.input {
			f.Item(input)

		}
		byteArr, err := json.Marshal(f.ToMapStr())
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
		if string(byteArr) != item.output {
			t.Errorf("expected %s not %s", item.output, string(byteArr))
			return
		}

	}

}

func TestORArrField(t *testing.T) {
	type testData struct {
		input  []interface{}
		output string
	}

	testDataArr := []testData{
		testData{
			input:  []interface{}{types.MapStr{"a": "c"}},
			output: `{"$or":[{"a":"c"}]}`,
		},
		testData{
			input:  []interface{}{types.MapStr{"a": "c"}, types.MapStr{"b": "c"}},
			output: `{"$or":[{"a":"c"},{"b":"c"}]}`,
		},
	}
	for _, item := range testDataArr {
		f := &orField{
			condition: CreateCondition(),
		}
		f.Array(item.input)
		byteArr, err := json.Marshal(f.ToMapStr())
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
		if string(byteArr) != item.output {
			t.Errorf("expected %s not %s", item.output, string(byteArr))
			return
		}

	}

}

func TestORMapStrArrField(t *testing.T) {
	type testData struct {
		input  []types.MapStr
		output string
	}

	testDataArr := []testData{
		testData{
			input:  []types.MapStr{types.MapStr{"a": "c"}},
			output: `{"$or":[{"a":"c"}]}`,
		},
		testData{
			input:  []types.MapStr{types.MapStr{"a": "c"}, types.MapStr{"b": "c"}},
			output: `{"$or":[{"a":"c"},{"b":"c"}]}`,
		},
	}
	for _, item := range testDataArr {
		f := &orField{
			condition: CreateCondition(),
		}
		f.MapStrArr(item.input)
		byteArr, err := json.Marshal(f.ToMapStr())
		if err != nil {
			t.Errorf("%s", err.Error())
			return
		}
		if string(byteArr) != item.output {
			t.Errorf("expected %s not %s", item.output, string(byteArr))
			return
		}

	}

}
